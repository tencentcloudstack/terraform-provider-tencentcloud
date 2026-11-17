package teo

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// ptrStringPlanForZone returns a pointer to s.
func ptrStringPlanForZone(s string) *string {
	return &s
}

// fakePlanForZoneClient implements TeoPlanForZoneClient for tests. Injecting a
// fake directly does not depend on runtime monkey patching, so the tests pass
// regardless of compiler inlining (no `-gcflags=all=-l` required).
type fakePlanForZoneClient struct {
	invoke func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error)
}

func (f fakePlanForZoneClient) TeoPlanForZone(zoneId, planType string) (resourceNames, dealNames []*string, errRet error) {
	return f.invoke(zoneId, planType)
}

// newConfiguredAction returns an action with the given fake client injected.
func newConfiguredAction(fake TeoPlanForZoneClient) *TeoPlanForZoneAction {
	return &TeoPlanForZoneAction{client: fake}
}

// planForZoneConfig builds an action.InvokeRequest whose Config tfsdk.Config
// holds the given zone_id and plan_type values, matching the action schema.
func planForZoneConfig(ctx context.Context, a action.Action, zoneId, planType string) (action.InvokeRequest, error) {
	schemaResp := action.SchemaResponse{}
	a.Schema(ctx, action.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		return action.InvokeRequest{}, fmt.Errorf("schema diagnostics: %v", schemaResp.Diagnostics)
	}

	tfType := schemaResp.Schema.Type().TerraformType(ctx)
	objType, ok := tfType.(tftypes.Object)
	if !ok {
		return action.InvokeRequest{}, fmt.Errorf("unexpected schema terraform type: %T", tfType)
	}

	val := map[string]tftypes.Value{
		"zone_id":   tftypes.NewValue(tftypes.String, zoneId),
		"plan_type": tftypes.NewValue(tftypes.String, planType),
	}

	return action.InvokeRequest{
		Config: tfsdk.Config{
			Raw:    tftypes.NewValue(objType, val),
			Schema: schemaResp.Schema,
		},
	}, nil
}

// invoke runs the action with the given config and returns the response.
func invoke(t *testing.T, a *TeoPlanForZoneAction, zoneId, planType string) *action.InvokeResponse {
	t.Helper()

	invokeReq, err := planForZoneConfig(context.Background(), a, zoneId, planType)
	assert.NoError(t, err)

	invokeResp := &action.InvokeResponse{}
	a.Invoke(context.Background(), invokeReq, invokeResp)
	return invokeResp
}

// TestTeoPlanForZoneAction_Success tests Invoke passes the config values
// through and returns no diagnostics on a successful call.
func TestTeoPlanForZoneAction_Success(t *testing.T) {
	var gotZoneId, gotPlanType string

	a := newConfiguredAction(fakePlanForZoneClient{
		invoke: func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error) {
			gotZoneId, gotPlanType = zoneId, planType
			return []*string{ptrStringPlanForZone("plan-res-123456")},
				[]*string{ptrStringPlanForZone("2025122000000123")},
				nil
		},
	})

	invokeResp := invoke(t, a, "zone-27h0vbm5w1e", "sta_global")

	assert.False(t, invokeResp.Diagnostics.HasError(), "expected no diagnostics, got: %v", invokeResp.Diagnostics)
	assert.Equal(t, "zone-27h0vbm5w1e", gotZoneId)
	assert.Equal(t, "sta_global", gotPlanType)
}

// TestTeoPlanForZoneAction_MissingZoneId tests Invoke errors when zone_id is empty.
func TestTeoPlanForZoneAction_MissingZoneId(t *testing.T) {
	invoked := false

	a := newConfiguredAction(fakePlanForZoneClient{
		invoke: func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error) {
			invoked = true
			return nil, nil, nil
		},
	})

	invokeResp := invoke(t, a, "", "sta_global")

	assert.True(t, invokeResp.Diagnostics.HasError())
	assert.Contains(t, invokeResp.Diagnostics.Errors()[0].Summary(), "zone_id")
	assert.False(t, invoked, "client must not be invoked when zone_id is empty")
}

// TestTeoPlanForZoneAction_MissingPlanType tests Invoke errors when plan_type is empty.
func TestTeoPlanForZoneAction_MissingPlanType(t *testing.T) {
	invoked := false

	a := newConfiguredAction(fakePlanForZoneClient{
		invoke: func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error) {
			invoked = true
			return nil, nil, nil
		},
	})

	invokeResp := invoke(t, a, "zone-27h0vbm5w1e", "")

	assert.True(t, invokeResp.Diagnostics.HasError())
	assert.Contains(t, invokeResp.Diagnostics.Errors()[0].Summary(), "plan_type")
	assert.False(t, invoked, "client must not be invoked when plan_type is empty")
}

// TestTeoPlanForZoneAction_APIError tests Invoke surfaces API errors.
func TestTeoPlanForZoneAction_APIError(t *testing.T) {
	a := newConfiguredAction(fakePlanForZoneClient{
		invoke: func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error) {
			return nil, nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter.ZoneHasBeenBound, Message=Zone has been bound")
		},
	})

	invokeResp := invoke(t, a, "zone-invalid", "sta_global")

	assert.True(t, invokeResp.Diagnostics.HasError())
	assert.Contains(t, invokeResp.Diagnostics.Errors()[0].Summary(), "CreatePlanForZone")
}

// TestTeoPlanForZoneAction_EmptyResponse tests Invoke errors when the call
// returns neither resource names nor deal names.
func TestTeoPlanForZoneAction_EmptyResponse(t *testing.T) {
	a := newConfiguredAction(fakePlanForZoneClient{
		invoke: func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error) {
			return nil, nil, nil
		},
	})

	invokeResp := invoke(t, a, "zone-27h0vbm5w1e", "sta_global")

	assert.True(t, invokeResp.Diagnostics.HasError())
	assert.Contains(t, invokeResp.Diagnostics.Errors()[0].Summary(), "CreatePlanForZone")
}

// TestTeoPlanForZoneAction_NotConfigured tests Invoke reports a clear error
// when the action was never configured with a client.
func TestTeoPlanForZoneAction_NotConfigured(t *testing.T) {
	a := &TeoPlanForZoneAction{}

	invokeResp := invoke(t, a, "zone-27h0vbm5w1e", "sta_global")

	assert.True(t, invokeResp.Diagnostics.HasError())
	assert.Contains(t, invokeResp.Diagnostics.Errors()[0].Summary(), "not configured")
}

// TestTeoPlanForZoneAction_Metadata validates the action type name.
func TestTeoPlanForZoneAction_Metadata(t *testing.T) {
	a := &TeoPlanForZoneAction{}
	req := action.MetadataRequest{}
	resp := &action.MetadataResponse{}

	a.Metadata(context.Background(), req, resp)

	assert.Equal(t, "tencentcloud_teo_plan_for_zone", resp.TypeName)
}

// TestTeoPlanForZoneAction_Schema validates the schema definition.
func TestTeoPlanForZoneAction_Schema(t *testing.T) {
	a := &TeoPlanForZoneAction{}
	req := action.SchemaRequest{}
	resp := &action.SchemaResponse{}

	a.Schema(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Schema.Attributes, "zone_id")
	assert.Contains(t, resp.Schema.Attributes, "plan_type")
}

// TestTeoPlanForZoneAction_ConfigureNilProviderData ensures Configure with nil
// ProviderData is a silent no-op.
func TestTeoPlanForZoneAction_ConfigureNilProviderData(t *testing.T) {
	a := &TeoPlanForZoneAction{}
	req := action.ConfigureRequest{}
	resp := &action.ConfigureResponse{}

	a.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Nil(t, a.client)
}

// TestTeoPlanForZoneAction_ConfigureNilClient ensures Configure keeps the
// client unset when the shared client is nil, so Invoke fails with the
// "not configured" diagnostic instead of panicking.
func TestTeoPlanForZoneAction_ConfigureNilClient(t *testing.T) {
	a := &TeoPlanForZoneAction{}
	req := action.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{}}
	resp := &action.ConfigureResponse{}

	a.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Nil(t, a.client)
}

// TestTeoPlanForZoneAction_ConfigureWithSharedMeta ensures Configure wires the
// TeoService built from the shared client into the action.
func TestTeoPlanForZoneAction_ConfigureWithSharedMeta(t *testing.T) {
	a := &TeoPlanForZoneAction{}
	req := action.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}}}
	resp := &action.ConfigureResponse{}

	a.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.NotNil(t, a.client)
}

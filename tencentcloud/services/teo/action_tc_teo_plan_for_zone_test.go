package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// ptrStringPlanForZone returns a pointer to s.
func ptrStringPlanForZone(s string) *string {
	return &s
}

// newSharedMeta builds a *sharedmeta.ProviderMeta wrapping a non-nil
// *connectivity.TencentCloudClient for framework Action Configure.
func newSharedMeta() *sharedmeta.ProviderMeta {
	return &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}}
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

// go test ./tencentcloud/services/teo/ -run "TestTeoPlanForZoneAction" -v -count=1 -gcflags="all=-l"

// TestTeoPlanForZoneAction_Success tests Invoke calls CreatePlanForZone and
// returns no diagnostics on a successful response.
func TestTeoPlanForZoneAction_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newSharedMeta().Client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePlanForZone", func(request *teov20220901.CreatePlanForZoneRequest) (*teov20220901.CreatePlanForZoneResponse, error) {
		assert.Equal(t, "zone-27h0vbm5w1e", *request.ZoneId)
		assert.Equal(t, "sta_global", *request.PlanType)

		resp := teov20220901.NewCreatePlanForZoneResponse()
		resp.Response = &teov20220901.CreatePlanForZoneResponseParams{
			ResourceNames: []*string{
				ptrStringPlanForZone("plan-res-123456"),
			},
			DealNames: []*string{
				ptrStringPlanForZone("2025122000000123"),
			},
			RequestId: ptrStringPlanForZone("fake-request-id"),
		}
		return resp, nil
	})

	a := teo.NewTeoPlanForZoneAction()

	// Configure the action with the shared provider meta.
	configureReq := action.ConfigureRequest{ProviderData: newSharedMeta()}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	ctx := context.Background()
	invokeReq, err := planForZoneConfig(ctx, a, "zone-27h0vbm5w1e", "sta_global")
	assert.NoError(t, err)

	invokeResp := &action.InvokeResponse{}
	a.Invoke(ctx, invokeReq, invokeResp)
	assert.False(t, invokeResp.Diagnostics.HasError(), "expected no diagnostics, got: %v", invokeResp.Diagnostics)
}

// TestTeoPlanForZoneAction_MissingZoneId tests Invoke errors when zone_id is empty.
func TestTeoPlanForZoneAction_MissingZoneId(t *testing.T) {
	a := teo.NewTeoPlanForZoneAction()

	configureReq := action.ConfigureRequest{ProviderData: newSharedMeta()}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	ctx := context.Background()
	invokeReq, err := planForZoneConfig(ctx, a, "", "sta_global")
	assert.NoError(t, err)

	invokeResp := &action.InvokeResponse{}
	a.Invoke(ctx, invokeReq, invokeResp)
	assert.True(t, invokeResp.Diagnostics.HasError())
	assert.Contains(t, invokeResp.Diagnostics.Errors()[0].Summary(), "zone_id")
}

// TestTeoPlanForZoneAction_MissingPlanType tests Invoke errors when plan_type is empty.
func TestTeoPlanForZoneAction_MissingPlanType(t *testing.T) {
	a := teo.NewTeoPlanForZoneAction()

	configureReq := action.ConfigureRequest{ProviderData: newSharedMeta()}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	ctx := context.Background()
	invokeReq, err := planForZoneConfig(ctx, a, "zone-27h0vbm5w1e", "")
	assert.NoError(t, err)

	invokeResp := &action.InvokeResponse{}
	a.Invoke(ctx, invokeReq, invokeResp)
	assert.True(t, invokeResp.Diagnostics.HasError())
	assert.Contains(t, invokeResp.Diagnostics.Errors()[0].Summary(), "plan_type")
}

// TestTeoPlanForZoneAction_APIError tests Invoke surfaces API errors.
func TestTeoPlanForZoneAction_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newSharedMeta().Client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePlanForZone", func(request *teov20220901.CreatePlanForZoneRequest) (*teov20220901.CreatePlanForZoneResponse, error) {
		assert.Equal(t, "zone-invalid", *request.ZoneId)
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter.ZoneHasBeenBound, Message=Zone has been bound")
	})

	a := teo.NewTeoPlanForZoneAction()

	configureReq := action.ConfigureRequest{ProviderData: newSharedMeta()}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	ctx := context.Background()
	invokeReq, err := planForZoneConfig(ctx, a, "zone-invalid", "sta_global")
	assert.NoError(t, err)

	invokeResp := &action.InvokeResponse{}
	a.Invoke(ctx, invokeReq, invokeResp)
	assert.True(t, invokeResp.Diagnostics.HasError())
	assert.Contains(t, invokeResp.Diagnostics.Errors()[0].Summary(), "CreatePlanForZone")
}

// TestTeoPlanForZoneAction_EmptyResponse tests Invoke errors on a nil/empty response.
func TestTeoPlanForZoneAction_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newSharedMeta().Client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePlanForZone", func(request *teov20220901.CreatePlanForZoneRequest) (*teov20220901.CreatePlanForZoneResponse, error) {
		resp := teov20220901.NewCreatePlanForZoneResponse()
		resp.Response = nil
		return resp, nil
	})

	a := teo.NewTeoPlanForZoneAction()

	configureReq := action.ConfigureRequest{ProviderData: newSharedMeta()}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	ctx := context.Background()
	invokeReq, err := planForZoneConfig(ctx, a, "zone-27h0vbm5w1e", "sta_global")
	assert.NoError(t, err)

	invokeResp := &action.InvokeResponse{}
	a.Invoke(ctx, invokeReq, invokeResp)
	assert.True(t, invokeResp.Diagnostics.HasError())
}

// TestTeoPlanForZoneAction_Metadata validates the action type name.
func TestTeoPlanForZoneAction_Metadata(t *testing.T) {
	a := teo.NewTeoPlanForZoneAction()
	req := action.MetadataRequest{}
	resp := &action.MetadataResponse{}

	a.Metadata(context.Background(), req, resp)

	assert.Equal(t, "tencentcloud_teo_plan_for_zone", resp.TypeName)
}

// TestTeoPlanForZoneAction_Schema validates the schema definition.
func TestTeoPlanForZoneAction_Schema(t *testing.T) {
	a := teo.NewTeoPlanForZoneAction()
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
	a := &teo.TeoPlanForZoneAction{}
	req := action.ConfigureRequest{}
	resp := &action.ConfigureResponse{}

	a.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
}

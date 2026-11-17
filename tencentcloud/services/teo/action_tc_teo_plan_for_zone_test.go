package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringTeoPlanForZone(s string) *string {
	return &s
}

// patchTeoPlanForZone patches the TeoService.TeoPlanForZone service method with
// the given fake. Unlike patching the small SDK client getters (e.g.
// UseTeoV20220901Client), which the compiler may inline into the caller and
// silently defeat the monkey patch, the service method is far beyond the
// inlining budget, so this patch is effective no matter how the tests are
// compiled (with or without `-gcflags=all=-l`).
func patchTeoPlanForZone(patches *gomonkey.Patches, fn func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error)) {
	service := teo.NewTeoService(&connectivity.TencentCloudClient{})
	patches.ApplyMethodFunc(&service, "TeoPlanForZone", fn)
}

// teoPlanForZoneSchema returns the framework action schema used to build the
// raw config value for InvokeRequest.
func teoPlanForZoneSchema(t *testing.T) action_schema.Schema {
	a := teo.NewTeoPlanForZoneAction()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newTeoPlanForZoneInvokeRequest builds an action.InvokeRequest carrying the
// given zone_id and plan_type as raw tftypes values.
func newTeoPlanForZoneInvokeRequest(t *testing.T, zoneId, planType string) action.InvokeRequest {
	s := teoPlanForZoneSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":   tftypes.String,
			"plan_type": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"zone_id":   tftypes.NewValue(tftypes.String, zoneId),
		"plan_type": tftypes.NewValue(tftypes.String, planType),
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setTeoPlanForZoneClient injects a mock connectivity client into the action
// instance via the Configure method.
func setTeoPlanForZoneClient(a *teo.TeoPlanForZoneAction, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoPlanForZoneAction" -v -count=1

// TestTeoPlanForZoneAction_Invoke_Success tests a successful invoke passes the
// config values through to the service call.
func TestTeoPlanForZoneAction_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var gotZoneId, gotPlanType string
	patchTeoPlanForZone(patches, func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error) {
		gotZoneId, gotPlanType = zoneId, planType
		return []*string{ptrStringTeoPlanForZone("plan-res-123456")},
			[]*string{ptrStringTeoPlanForZone("2025122000000123")},
			nil
	})

	a := teo.NewTeoPlanForZoneAction()
	setTeoPlanForZoneClient(a.(*teo.TeoPlanForZoneAction), &connectivity.TencentCloudClient{})

	req := newTeoPlanForZoneInvokeRequest(t, "zone-27h0vbm5w1e", "sta_global")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.Equal(t, "zone-27h0vbm5w1e", gotZoneId)
	assert.Equal(t, "sta_global", gotPlanType)
}

// TestTeoPlanForZoneAction_Invoke_APIError tests API error handling.
func TestTeoPlanForZoneAction_Invoke_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchTeoPlanForZone(patches, func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error) {
		return nil, nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter.ZoneHasBeenBound, Message=Zone has been bound")
	})

	a := teo.NewTeoPlanForZoneAction()
	setTeoPlanForZoneClient(a.(*teo.TeoPlanForZoneAction), &connectivity.TencentCloudClient{})

	req := newTeoPlanForZoneInvokeRequest(t, "zone-invalid", "sta_global")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "purchasing TEO plan for zone")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ZoneHasBeenBound")
}

// TestTeoPlanForZoneAction_Invoke_MissingInput tests missing required inputs
// are rejected before any service call.
func TestTeoPlanForZoneAction_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	apiCalled := false
	patchTeoPlanForZone(patches, func(zoneId, planType string) (resourceNames, dealNames []*string, errRet error) {
		apiCalled = true
		return nil, nil, nil
	})

	a := teo.NewTeoPlanForZoneAction()
	setTeoPlanForZoneClient(a.(*teo.TeoPlanForZoneAction), &connectivity.TencentCloudClient{})

	// empty zone_id
	req := newTeoPlanForZoneInvokeRequest(t, "", "sta_global")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing zone_id")
	assert.False(t, apiCalled)

	// empty plan_type
	req = newTeoPlanForZoneInvokeRequest(t, "zone-27h0vbm5w1e", "")
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing plan_type")
	assert.False(t, apiCalled)
}

// TestTeoPlanForZoneAction_Invoke_ClientNotConfigured tests the nil-client
// guard.
func TestTeoPlanForZoneAction_Invoke_ClientNotConfigured(t *testing.T) {
	a := teo.NewTeoPlanForZoneAction()

	req := newTeoPlanForZoneInvokeRequest(t, "zone-27h0vbm5w1e", "sta_global")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestTeoPlanForZoneAction_MetadataAndSchema validates metadata and schema
// definition.
func TestTeoPlanForZoneAction_MetadataAndSchema(t *testing.T) {
	a := teo.NewTeoPlanForZoneAction()

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_plan_for_zone", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "zone_id")
	assert.Contains(t, attrs, "plan_type")

	zoneAttr, ok := attrs["zone_id"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, zoneAttr.Required)

	planAttr, ok := attrs["plan_type"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, planAttr.Required)
}

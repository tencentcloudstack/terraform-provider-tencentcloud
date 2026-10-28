package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoBindZoneToPlanAction" -v -count=1 -gcflags="all=-l"

// newBindZoneToPlanActionConfig builds an action.InvokeRequest.Config from the
// action schema populated with the given zone_id and plan_id values.
func newBindZoneToPlanActionConfig(t *testing.T, zoneId, planId string) tfsdk.Config {
	t.Helper()

	a := &teo.TeoBindZoneToPlanAction{}
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", schemaResp.Diagnostics)
	}

	objType := schemaResp.Schema.Type().(types.ObjectType)
	tfType := objType.TerraformType(context.Background())

	raw := tftypes.NewValue(tfType, map[string]tftypes.Value{
		"zone_id": tftypes.NewValue(tftypes.String, zoneId),
		"plan_id": tftypes.NewValue(tftypes.String, planId),
	})

	return tfsdk.Config{
		Raw:    raw,
		Schema: schemaResp.Schema,
	}
}

// TestTeoBindZoneToPlanAction_Metadata validates the action type name.
func TestTeoBindZoneToPlanAction_Metadata(t *testing.T) {
	a := &teo.TeoBindZoneToPlanAction{}
	req := action.MetadataRequest{}
	resp := &action.MetadataResponse{}

	a.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_teo_bind_zone_to_plan" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_teo_bind_zone_to_plan")
	}
}

// TestTeoBindZoneToPlanAction_Schema validates required zone_id / plan_id
// attributes on the action schema.
func TestTeoBindZoneToPlanAction_Schema(t *testing.T) {
	a := &teo.TeoBindZoneToPlanAction{}
	req := action.SchemaRequest{}
	resp := &action.SchemaResponse{}

	a.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	for name, wantRequired := range map[string]bool{
		"zone_id": true,
		"plan_id": true,
	} {
		attr, ok := resp.Schema.Attributes[name].(actionschema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q missing or not a StringAttribute", name)
			continue
		}
		if attr.Required != wantRequired {
			t.Errorf("attribute %q Required=%v, want %v", name, attr.Required, wantRequired)
		}
	}
}

// TestTeoBindZoneToPlanAction_ImplementsInterfaces validates that
// TeoBindZoneToPlanAction implements action.Action and action.ActionWithConfigure.
func TestTeoBindZoneToPlanAction_ImplementsInterfaces(t *testing.T) {
	var _ action.Action = &teo.TeoBindZoneToPlanAction{}
	var _ action.ActionWithConfigure = &teo.TeoBindZoneToPlanAction{}
}

// TestTeoBindZoneToPlanAction_Factory validates the registry factory.
func TestTeoBindZoneToPlanAction_Factory(t *testing.T) {
	got := teo.NewTeoBindZoneToPlanAction()
	if got == nil {
		t.Fatal("NewTeoBindZoneToPlanAction returned nil")
	}
	if _, ok := got.(*teo.TeoBindZoneToPlanAction); !ok {
		t.Fatalf("factory returned %T, want *teo.TeoBindZoneToPlanAction", got)
	}
}

// TestTeoBindZoneToPlanAction_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestTeoBindZoneToPlanAction_ConfigureNilProviderData(t *testing.T) {
	a := &teo.TeoBindZoneToPlanAction{}
	req := action.ConfigureRequest{}
	resp := &action.ConfigureResponse{}

	a.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// TestTeoBindZoneToPlanAction_ConfigureWrongType ensures Configure adds an
// error diagnostic when ProviderData is of an unexpected type.
func TestTeoBindZoneToPlanAction_ConfigureWrongType(t *testing.T) {
	a := &teo.TeoBindZoneToPlanAction{}
	req := action.ConfigureRequest{ProviderData: "not-a-meta"}
	resp := &action.ConfigureResponse{}

	a.Configure(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Error("expected error diagnostic for wrong ProviderData type")
	}
}

// TestTeoBindZoneToPlanAction_InvokeSuccess mocks BindZoneToPlanWithContext
// and asserts the action invokes the API with the provided zone_id and plan_id.
func TestTeoBindZoneToPlanAction_InvokeSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var capturedReq *teov20220901.BindZoneToPlanRequest
	patches.ApplyMethodFunc(teoClient, "BindZoneToPlanWithContext", func(_ context.Context, request *teov20220901.BindZoneToPlanRequest) (*teov20220901.BindZoneToPlanResponse, error) {
		capturedReq = request
		resp := teov20220901.NewBindZoneToPlanResponse()
		resp.Response = &teov20220901.BindZoneToPlanResponseParams{
			RequestId: ptrStringBindZoneToPlanAction("fake-request-id"),
		}
		return resp, nil
	})

	a := &teo.TeoBindZoneToPlanAction{Client: client}
	req := action.InvokeRequest{Config: newBindZoneToPlanActionConfig(t, "zone-12345678", "edgeone-12345678")}
	resp := &action.InvokeResponse{}

	a.Invoke(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Invoke returned unexpected diagnostics: %v", resp.Diagnostics)
	}
	assert.NotNil(t, capturedReq)
	assert.Equal(t, "zone-12345678", *capturedReq.ZoneId)
	assert.Equal(t, "edgeone-12345678", *capturedReq.PlanId)
}

// TestTeoBindZoneToPlanAction_InvokeAPIError mocks BindZoneToPlanWithContext
// to return a non-retryable error and asserts the action surfaces a diagnostic.
func TestTeoBindZoneToPlanAction_InvokeAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "BindZoneToPlanWithContext", func(_ context.Context, request *teov20220901.BindZoneToPlanRequest) (*teov20220901.BindZoneToPlanResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	a := &teo.TeoBindZoneToPlanAction{Client: client}
	req := action.InvokeRequest{Config: newBindZoneToPlanActionConfig(t, "zone-invalid", "edgeone-invalid")}
	resp := &action.InvokeResponse{}

	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostic for API error")
	}
}

// TestTeoBindZoneToPlanAction_InvokeMissingZoneId asserts the action adds a
// diagnostic when zone_id is empty.
func TestTeoBindZoneToPlanAction_InvokeMissingZoneId(t *testing.T) {
	a := &teo.TeoBindZoneToPlanAction{Client: &connectivity.TencentCloudClient{}}
	req := action.InvokeRequest{Config: newBindZoneToPlanActionConfig(t, "", "edgeone-12345678")}
	resp := &action.InvokeResponse{}

	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostic when zone_id is empty")
	}
}

// TestTeoBindZoneToPlanAction_InvokeMissingPlanId asserts the action adds a
// diagnostic when plan_id is empty.
func TestTeoBindZoneToPlanAction_InvokeMissingPlanId(t *testing.T) {
	a := &teo.TeoBindZoneToPlanAction{Client: &connectivity.TencentCloudClient{}}
	req := action.InvokeRequest{Config: newBindZoneToPlanActionConfig(t, "zone-12345678", "")}
	resp := &action.InvokeResponse{}

	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostic when plan_id is empty")
	}
}

func ptrStringBindZoneToPlanAction(s string) *string {
	return &s
}

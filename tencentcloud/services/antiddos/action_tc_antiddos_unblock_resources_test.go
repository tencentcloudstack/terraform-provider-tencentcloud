package antiddos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	antiddosv20250903 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"
)

func ptrStringAntiddosUnblockResources(s string) *string {
	return &s
}

// antiddosUnblockResourcesSchema returns the framework action schema used to
// build the raw config value for InvokeRequest.
func antiddosUnblockResourcesSchema(t *testing.T) action_schema.Schema {
	a := antiddos.NewAntiddosUnblockResources()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newAntiddosUnblockResourcesInvokeRequest builds an action.InvokeRequest
// carrying the given resources as raw tftypes values.
func newAntiddosUnblockResourcesInvokeRequest(t *testing.T, resources []string) action.InvokeRequest {
	s := antiddosUnblockResourcesSchema(t)

	resVals := make([]tftypes.Value, 0, len(resources))
	for _, r := range resources {
		resVals = append(resVals, tftypes.NewValue(tftypes.String, r))
	}

	var resourcesVal tftypes.Value
	if len(resVals) == 0 {
		resourcesVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil)
	} else {
		resourcesVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, resVals)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"resources": tftypes.List{ElementType: tftypes.String},
		},
	}, map[string]tftypes.Value{
		"resources": resourcesVal,
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// newAntiddosUnblockResourcesInvokeRequestNullResources builds an InvokeRequest
// whose resources is null.
func newAntiddosUnblockResourcesInvokeRequestNullResources(t *testing.T) action.InvokeRequest {
	s := antiddosUnblockResourcesSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"resources": tftypes.List{ElementType: tftypes.String},
		},
	}, map[string]tftypes.Value{
		"resources": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setAntiddosUnblockResourcesClient injects a mock connectivity client into
// the action instance via the Configure method.
func setAntiddosUnblockResourcesClient(a *antiddos.AntiddosUnblockResources, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/antiddos/ -run "TestAntiddosUnblockResources" -v -count=1 -gcflags="all=-l"

// TestAntiddosUnblockResources_Invoke_Success tests a successful invoke
func TestAntiddosUnblockResources_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	antiddosClient := &antiddosv20250903.Client{}
	patches.ApplyMethodReturn(client, "UseAntiddosV20250903Client", antiddosClient)

	var gotRequest *antiddosv20250903.UnblockResourcesRequest
	patches.ApplyMethodFunc(antiddosClient, "UnblockResourcesWithContext", func(_ context.Context, request *antiddosv20250903.UnblockResourcesRequest) (*antiddosv20250903.UnblockResourcesResponse, error) {
		gotRequest = request
		resp := antiddosv20250903.NewUnblockResourcesResponse()
		resp.Response = &antiddosv20250903.UnblockResourcesResponseParams{
			RequestId: ptrStringAntiddosUnblockResources("fake-request-id"),
		}
		return resp, nil
	})

	a := antiddos.NewAntiddosUnblockResources()
	setAntiddosUnblockResourcesClient(a.(*antiddos.AntiddosUnblockResources), client)

	req := newAntiddosUnblockResourcesInvokeRequest(t, []string{"117.175.94.230", "117.175.94.231"})
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, []*string{
		ptrStringAntiddosUnblockResources("117.175.94.230"),
		ptrStringAntiddosUnblockResources("117.175.94.231"),
	}, gotRequest.Resources)
}

// TestAntiddosUnblockResources_Invoke_APIError tests API error handling
func TestAntiddosUnblockResources_Invoke_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	antiddosClient := &antiddosv20250903.Client{}
	patches.ApplyMethodReturn(client, "UseAntiddosV20250903Client", antiddosClient)

	patches.ApplyMethodFunc(antiddosClient, "UnblockResourcesWithContext", func(_ context.Context, request *antiddosv20250903.UnblockResourcesRequest) (*antiddosv20250903.UnblockResourcesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=LimitExceeded, Message=unblock quota exceeded")
	})

	a := antiddos.NewAntiddosUnblockResources()
	setAntiddosUnblockResourcesClient(a.(*antiddos.AntiddosUnblockResources), client)

	req := newAntiddosUnblockResourcesInvokeRequest(t, []string{"117.175.94.230"})
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "unblocking AntiDDoS resources")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "LimitExceeded")
}

// TestAntiddosUnblockResources_Invoke_MissingInput tests missing required input
// is rejected before any API call
func TestAntiddosUnblockResources_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	antiddosClient := &antiddosv20250903.Client{}
	patches.ApplyMethodReturn(client, "UseAntiddosV20250903Client", antiddosClient)

	apiCalled := false
	patches.ApplyMethodFunc(antiddosClient, "UnblockResourcesWithContext", func(_ context.Context, request *antiddosv20250903.UnblockResourcesRequest) (*antiddosv20250903.UnblockResourcesResponse, error) {
		apiCalled = true
		return antiddosv20250903.NewUnblockResourcesResponse(), nil
	})

	a := antiddos.NewAntiddosUnblockResources()
	setAntiddosUnblockResourcesClient(a.(*antiddos.AntiddosUnblockResources), client)

	// empty resources list
	req := newAntiddosUnblockResourcesInvokeRequest(t, []string{})
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing resources")
	assert.False(t, apiCalled)

	// null resources
	req = newAntiddosUnblockResourcesInvokeRequestNullResources(t)
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing resources")
	assert.False(t, apiCalled)
}

// TestAntiddosUnblockResources_Invoke_ClientNotConfigured tests the nil-client guard
func TestAntiddosUnblockResources_Invoke_ClientNotConfigured(t *testing.T) {
	a := antiddos.NewAntiddosUnblockResources()

	req := newAntiddosUnblockResourcesInvokeRequest(t, []string{"117.175.94.230"})
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestAntiddosUnblockResources_MetadataAndSchema validates metadata and schema definition
func TestAntiddosUnblockResources_MetadataAndSchema(t *testing.T) {
	a := antiddos.NewAntiddosUnblockResources()

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_antiddos_unblock_resources", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "resources")

	assert.True(t, attrs["resources"].IsRequired())

	resAttr, ok := attrs["resources"].(action_schema.ListAttribute)
	assert.True(t, ok)
	assert.Equal(t, types.StringType, resAttr.ElementType)
}

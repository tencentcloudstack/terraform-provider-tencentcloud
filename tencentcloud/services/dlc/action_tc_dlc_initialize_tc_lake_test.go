package dlc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"
)

func ptrStringDlcInitializeTCLake(s string) *string {
	return &s
}

func ptrBoolDlcInitializeTCLake(b bool) *bool {
	return &b
}

// dlcInitializeTCLakeSchema returns the framework action schema used to build
// the raw config value for InvokeRequest.
func dlcInitializeTCLakeSchema(t *testing.T) action_schema.Schema {
	a := dlc.NewDlcInitializeTCLake()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newDlcInitializeTCLakeInvokeRequest builds an action.InvokeRequest carrying
// an empty config object (the InitializeTCLake API takes no parameters).
func newDlcInitializeTCLakeInvokeRequest(t *testing.T) action.InvokeRequest {
	s := dlcInitializeTCLakeSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{},
	}, map[string]tftypes.Value{})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setDlcInitializeTCLakeClient injects a mock connectivity client into the
// action instance via the Configure method.
func setDlcInitializeTCLakeClient(a *dlc.DlcInitializeTCLake, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/dlc/ -run "TestDlcInitializeTCLake" -v -count=1 -gcflags="all=-l"

// TestDlcInitializeTCLake_Invoke_Success tests a successful invoke
func TestDlcInitializeTCLake_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	dlcClient := &dlcv20210125.Client{}
	patches.ApplyMethodReturn(client, "UseDlcClient", dlcClient)

	var gotRequest *dlcv20210125.InitializeTCLakeRequest
	patches.ApplyMethodFunc(dlcClient, "InitializeTCLakeWithContext", func(_ context.Context, request *dlcv20210125.InitializeTCLakeRequest) (*dlcv20210125.InitializeTCLakeResponse, error) {
		gotRequest = request
		resp := dlcv20210125.NewInitializeTCLakeResponse()
		resp.Response = &dlcv20210125.InitializeTCLakeResponseParams{
			InstanceId: ptrStringDlcInitializeTCLake("dlc-instance-xxxx"),
			IsSuccess:  ptrBoolDlcInitializeTCLake(true),
			RequestId:  ptrStringDlcInitializeTCLake("fake-request-id"),
		}
		return resp, nil
	})

	a := dlc.NewDlcInitializeTCLake()
	setDlcInitializeTCLakeClient(a.(*dlc.DlcInitializeTCLake), client)

	req := newDlcInitializeTCLakeInvokeRequest(t)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotRequest)
}

// TestDlcInitializeTCLake_Invoke_APIError tests API error handling
func TestDlcInitializeTCLake_Invoke_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	dlcClient := &dlcv20210125.Client{}
	patches.ApplyMethodReturn(client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "InitializeTCLakeWithContext", func(_ context.Context, request *dlcv20210125.InitializeTCLakeRequest) (*dlcv20210125.InitializeTCLakeResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.ExternalService, Message=initialize tc lake failed")
	})

	a := dlc.NewDlcInitializeTCLake()
	setDlcInitializeTCLakeClient(a.(*dlc.DlcInitializeTCLake), client)

	req := newDlcInitializeTCLakeInvokeRequest(t)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "initializing DLC TCLake")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "FailedOperation.ExternalService")
}

// TestDlcInitializeTCLake_Invoke_ClientNotConfigured tests the nil-client guard
func TestDlcInitializeTCLake_Invoke_ClientNotConfigured(t *testing.T) {
	a := dlc.NewDlcInitializeTCLake()

	req := newDlcInitializeTCLakeInvokeRequest(t)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestDlcInitializeTCLake_MetadataAndSchema validates metadata and schema definition
func TestDlcInitializeTCLake_MetadataAndSchema(t *testing.T) {
	a := dlc.NewDlcInitializeTCLake()

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_dlc_initialize_tc_lake", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Empty(t, attrs, "expected empty attributes map for parameterless InitializeTCLake action")
}

package ssm_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	ssmsdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcssm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssm"
)

func ptrStringSsmService(s string) *string {
	return &s
}

func ptrBoolSsmService(b bool) *bool {
	return &b
}

// ssmServiceActionSchema returns the framework action schema used to build
// the raw config value for InvokeRequest.
func ssmServiceActionSchema(t *testing.T) action_schema.Schema {
	a := svcssm.NewSsmServiceAction()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newSsmServiceActionInvokeRequest builds an action.InvokeRequest carrying the
// given pay_mode as raw tftypes value. Pass an empty string to produce a null
// value (for the missing-input test).
func newSsmServiceActionInvokeRequest(t *testing.T, payMode string) action.InvokeRequest {
	s := ssmServiceActionSchema(t)

	var payModeVal tftypes.Value
	if payMode == "" {
		payModeVal = tftypes.NewValue(tftypes.String, nil)
	} else {
		payModeVal = tftypes.NewValue(tftypes.String, payMode)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"pay_mode": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"pay_mode": payModeVal,
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setSsmServiceActionClient injects a mock connectivity client into the
// action instance via the Configure method.
func setSsmServiceActionClient(a *svcssm.SsmServiceAction, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// newSsmServiceCreateResponse builds a successful CreateService response.
func newSsmServiceCreateResponse() *ssmsdk.CreateServiceResponse {
	resp := ssmsdk.NewCreateServiceResponse()
	resp.Response = &ssmsdk.CreateServiceResponseParams{
		ResourceId: ptrStringSsmService("ssm_1250000000"),
		OrderId:    ptrStringSsmService("2024102700001"),
		RequestId:  ptrStringSsmService("fake-request-id"),
	}
	return resp
}

// newSsmServiceStatusResponse builds a GetServiceStatus response whose
// ServiceEnabled field is set to the given bool.
func newSsmServiceStatusResponse(enabled bool) *ssmsdk.GetServiceStatusResponse {
	resp := ssmsdk.NewGetServiceStatusResponse()
	resp.Response = &ssmsdk.GetServiceStatusResponseParams{
		ServiceEnabled: ptrBoolSsmService(enabled),
		InvalidType:    ptrInt64SsmService(0),
		RequestId:      ptrStringSsmService("fake-request-id"),
	}
	return resp
}

func ptrInt64SsmService(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/ssm/ -run "TestSsmServiceAction" -v -count=1 -gcflags="all=-l"

// TestSsmServiceAction_Invoke_Success tests a successful invoke where
// CreateService succeeds and a subsequent GetServiceStatus poll returns
// ServiceEnabled == true.
func TestSsmServiceAction_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmsdk.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	var gotCreateRequest *ssmsdk.CreateServiceRequest
	patches.ApplyMethodFunc(ssmClient, "CreateServiceWithContext", func(_ context.Context, request *ssmsdk.CreateServiceRequest) (*ssmsdk.CreateServiceResponse, error) {
		gotCreateRequest = request
		return newSsmServiceCreateResponse(), nil
	})

	var statusCalls int32
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatusWithContext", func(_ context.Context, request *ssmsdk.GetServiceStatusRequest) (*ssmsdk.GetServiceStatusResponse, error) {
		atomic.AddInt32(&statusCalls, 1)
		return newSsmServiceStatusResponse(true), nil
	})

	a := svcssm.NewSsmServiceAction()
	setSsmServiceActionClient(a.(*svcssm.SsmServiceAction), client)

	req := newSsmServiceActionInvokeRequest(t, "PostPay")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotCreateRequest, "CreateService should have been called")
	assert.Equal(t, "PostPay", *gotCreateRequest.PayMode)
	assert.NotEmpty(t, *gotCreateRequest.ClientToken, "client_token should be generated internally via helper.BuildToken()")
	assert.Equal(t, int32(1), atomic.LoadInt32(&statusCalls), "GetServiceStatus should be polled once when already enabled")
}

// TestSsmServiceAction_Invoke_PollNotYetEnabled tests that polling continues
// when the first GetServiceStatus returns ServiceEnabled == false and a later
// poll returns ServiceEnabled == true.
func TestSsmServiceAction_Invoke_PollNotYetEnabled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmsdk.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateServiceWithContext", func(_ context.Context, request *ssmsdk.CreateServiceRequest) (*ssmsdk.CreateServiceResponse, error) {
		return newSsmServiceCreateResponse(), nil
	})

	var pollCount int32
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatusWithContext", func(_ context.Context, request *ssmsdk.GetServiceStatusRequest) (*ssmsdk.GetServiceStatusResponse, error) {
		n := atomic.AddInt32(&pollCount, 1)
		if n == 1 {
			return newSsmServiceStatusResponse(false), nil
		}
		return newSsmServiceStatusResponse(true), nil
	})

	a := svcssm.NewSsmServiceAction()
	setSsmServiceActionClient(a.(*svcssm.SsmServiceAction), client)

	req := newSsmServiceActionInvokeRequest(t, "PostPay")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.GreaterOrEqual(t, atomic.LoadInt32(&pollCount), int32(2), "expected at least two GetServiceStatus polls")
}

// TestSsmServiceAction_Invoke_MissingPayMode tests that a null pay_mode is
// rejected before any API call.
func TestSsmServiceAction_Invoke_MissingPayMode(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmsdk.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	apiCalled := false
	patches.ApplyMethodFunc(ssmClient, "CreateServiceWithContext", func(_ context.Context, request *ssmsdk.CreateServiceRequest) (*ssmsdk.CreateServiceResponse, error) {
		apiCalled = true
		return newSsmServiceCreateResponse(), nil
	})

	a := svcssm.NewSsmServiceAction()
	setSsmServiceActionClient(a.(*svcssm.SsmServiceAction), client)

	req := newSsmServiceActionInvokeRequest(t, "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing pay_mode")
	assert.False(t, apiCalled)
}

// TestSsmServiceAction_Invoke_CreateServiceError tests that a CreateService
// API error surfaces a diagnostic referencing opening SSM service.
func TestSsmServiceAction_Invoke_CreateServiceError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmsdk.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateServiceWithContext", func(_ context.Context, request *ssmsdk.CreateServiceRequest) (*ssmsdk.CreateServiceResponse, error) {
		return nil, fmt.Errorf("Code=ResourceInUse.ServiceAlreadyEnabled, Message=service already enabled")
	})

	a := svcssm.NewSsmServiceAction()
	setSsmServiceActionClient(a.(*svcssm.SsmServiceAction), client)

	req := newSsmServiceActionInvokeRequest(t, "PostPay")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "opening SSM service")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ResourceInUse.ServiceAlreadyEnabled")
}

// TestSsmServiceAction_Invoke_PollError tests that a GetServiceStatus error
// during polling surfaces a diagnostic referencing polling SSM service status.
func TestSsmServiceAction_Invoke_PollError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmsdk.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateServiceWithContext", func(_ context.Context, request *ssmsdk.CreateServiceRequest) (*ssmsdk.CreateServiceResponse, error) {
		return newSsmServiceCreateResponse(), nil
	})

	patches.ApplyMethodFunc(ssmClient, "GetServiceStatusWithContext", func(_ context.Context, request *ssmsdk.GetServiceStatusRequest) (*ssmsdk.GetServiceStatusResponse, error) {
		return nil, fmt.Errorf("internal error while querying SSM service status")
	})

	a := svcssm.NewSsmServiceAction()
	setSsmServiceActionClient(a.(*svcssm.SsmServiceAction), client)

	req := newSsmServiceActionInvokeRequest(t, "PostPay")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "polling SSM service status")
}

// TestSsmServiceAction_Invoke_ClientNotConfigured tests the nil-client guard.
func TestSsmServiceAction_Invoke_ClientNotConfigured(t *testing.T) {
	a := svcssm.NewSsmServiceAction()

	req := newSsmServiceActionInvokeRequest(t, "PostPay")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestSsmServiceAction_MetadataAndSchema validates metadata and schema
// definition. client_token must NOT be exposed as an action attribute.
func TestSsmServiceAction_MetadataAndSchema(t *testing.T) {
	a := svcssm.NewSsmServiceAction()

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_ssm_service", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Len(t, attrs, 1)
	assert.NotContains(t, attrs, "client_token", "client_token must not be exposed as an action attribute")

	payModeAttr, ok := attrs["pay_mode"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, payModeAttr.IsRequired())

	// Ensure no output attributes are exposed.
	for name := range attrs {
		assert.NotContains(t, []string{"resource_id", "order_id", "service_enabled"}, name, "action schema must not expose output attribute %q", name)
	}
}

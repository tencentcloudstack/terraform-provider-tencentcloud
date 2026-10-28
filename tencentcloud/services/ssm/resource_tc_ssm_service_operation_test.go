package ssm_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcssm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssm"
)

// ssmServiceOpMockMeta implements tccommon.ProviderMeta
type ssmServiceOpMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *ssmServiceOpMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &ssmServiceOpMockMeta{}

func newSsmServiceOpMockMeta() *ssmServiceOpMockMeta {
	return &ssmServiceOpMockMeta{client: &connectivity.TencentCloudClient{}}
}

func ssmServiceOpPtrString(s string) *string {
	return &s
}

func ssmServiceOpPtrBool(b bool) *bool {
	return &b
}

// newSsmServiceStatusResponse builds a GetServiceStatus response whose
// ServiceEnabled field is set to the given bool.
func newSsmServiceStatusResponse(enabled bool) *ssm.GetServiceStatusResponse {
	resp := ssm.NewGetServiceStatusResponse()
	resp.Response = &ssm.GetServiceStatusResponseParams{
		ServiceEnabled: ssmServiceOpPtrBool(enabled),
		RequestId:      ssmServiceOpPtrString("fake-request-id"),
	}
	return resp
}

// go test ./tencentcloud/services/ssm/ -run "TestSsmServiceOperation" -v -count=1 -gcflags="all=-l"

// TestSsmServiceOperation_Create_Success tests a successful create where
// CreateService succeeds and a subsequent GetServiceStatus poll returns
// ServiceEnabled == true.
func TestSsmServiceOperation_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := newSsmServiceOpMockMeta().client
	ssmClient := &ssm.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	var gotCreateRequest *ssm.CreateServiceRequest
	patches.ApplyMethodFunc(ssmClient, "CreateServiceWithContext", func(_ context.Context, request *ssm.CreateServiceRequest) (*ssm.CreateServiceResponse, error) {
		gotCreateRequest = request
		resp := ssm.NewCreateServiceResponse()
		resp.Response = &ssm.CreateServiceResponseParams{
			ResourceId: ssmServiceOpPtrString("ssm_1250000000"),
			OrderId:    ssmServiceOpPtrString("2024102700001"),
			RequestId:  ssmServiceOpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ssmClient, "GetServiceStatusWithContext", func(_ context.Context, request *ssm.GetServiceStatusRequest) (*ssm.GetServiceStatusResponse, error) {
		return newSsmServiceStatusResponse(true), nil
	})

	meta := newSsmServiceOpMockMeta()
	res := svcssm.ResourceTencentCloudSsmServiceOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"pay_mode": "PostPay",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id(), "resource id should be set after create")
	assert.NotNil(t, gotCreateRequest, "CreateService should have been called")
	assert.Equal(t, "PostPay", *gotCreateRequest.PayMode)
	assert.NotEmpty(t, *gotCreateRequest.ClientToken, "client_token should be generated internally")
}

// TestSsmServiceOperation_Create_PollNotYetEnabled tests that polling continues
// when the first GetServiceStatus returns ServiceEnabled == false and a later
// poll returns ServiceEnabled == true.
func TestSsmServiceOperation_Create_PollNotYetEnabled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := newSsmServiceOpMockMeta().client
	ssmClient := &ssm.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateServiceWithContext", func(_ context.Context, request *ssm.CreateServiceRequest) (*ssm.CreateServiceResponse, error) {
		resp := ssm.NewCreateServiceResponse()
		resp.Response = &ssm.CreateServiceResponseParams{
			ResourceId: ssmServiceOpPtrString("ssm_1250000000"),
			RequestId:  ssmServiceOpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	var pollCount int32
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatusWithContext", func(_ context.Context, request *ssm.GetServiceStatusRequest) (*ssm.GetServiceStatusResponse, error) {
		n := atomic.AddInt32(&pollCount, 1)
		if n == 1 {
			return newSsmServiceStatusResponse(false), nil
		}
		return newSsmServiceStatusResponse(true), nil
	})

	meta := newSsmServiceOpMockMeta()
	res := svcssm.ResourceTencentCloudSsmServiceOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"pay_mode": "PostPay",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.GreaterOrEqual(t, atomic.LoadInt32(&pollCount), int32(2), "expected at least two GetServiceStatus polls")
}

// TestSsmServiceOperation_Create_APIError tests that a CreateService API error
// surfaces as a resource create error.
func TestSsmServiceOperation_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := newSsmServiceOpMockMeta().client
	ssmClient := &ssm.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateServiceWithContext", func(_ context.Context, request *ssm.CreateServiceRequest) (*ssm.CreateServiceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceInUse.ServiceAlreadyEnabled, Message=service already enabled")
	})

	meta := newSsmServiceOpMockMeta()
	res := svcssm.ResourceTencentCloudSsmServiceOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"pay_mode": "PostPay",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ServiceAlreadyEnabled")
}

// TestSsmServiceOperation_Read tests Read is a no-op.
func TestSsmServiceOperation_Read(t *testing.T) {
	res := svcssm.ResourceTencentCloudSsmServiceOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"pay_mode": "PostPay",
	})
	d.SetId("fake-id")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestSsmServiceOperation_Delete tests Delete is a no-op.
func TestSsmServiceOperation_Delete(t *testing.T) {
	res := svcssm.ResourceTencentCloudSsmServiceOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"pay_mode": "PostPay",
	})
	d.SetId("fake-id")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestSsmServiceOperation_Schema validates schema definition.
func TestSsmServiceOperation_Schema(t *testing.T) {
	res := svcssm.ResourceTencentCloudSsmServiceOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "pay_mode")
	payMode := res.Schema["pay_mode"]
	assert.Equal(t, schema.TypeString, payMode.Type)
	assert.True(t, payMode.Required)
	assert.True(t, payMode.ForceNew)

	// Ensure client_token is NOT exposed as a schema field.
	_, hasClientToken := res.Schema["client_token"]
	assert.False(t, hasClientToken, "client_token must not be exposed as a schema field")
}

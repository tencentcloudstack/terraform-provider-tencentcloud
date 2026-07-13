package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaIncreasePlanQuota implements tccommon.ProviderMeta
type mockMetaIncreasePlanQuota struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaIncreasePlanQuota) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaIncreasePlanQuota{}

func newMockMetaIncreasePlanQuota() *mockMetaIncreasePlanQuota {
	return &mockMetaIncreasePlanQuota{client: &connectivity.TencentCloudClient{}}
}

func ptrStringIncreasePlanQuota(s string) *string {
	return &s
}

func ptrInt64IncreasePlanQuota(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestIncreasePlanQuotaOperation" -v -count=1 -gcflags="all=-l"

// TestIncreasePlanQuotaOperation_Success tests successful quota increase
func TestIncreasePlanQuotaOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaIncreasePlanQuota().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuotaWithContext", func(_ context.Context, request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		assert.NotNil(t, request.PlanId)
		assert.Equal(t, "edgeone-2unuvzjmmn2q", *request.PlanId)
		assert.NotNil(t, request.QuotaType)
		assert.Equal(t, "site", *request.QuotaType)
		assert.NotNil(t, request.QuotaNumber)
		assert.Equal(t, int64(1), *request.QuotaNumber)

		resp := teov20220901.NewIncreasePlanQuotaResponse()
		resp.Response = &teov20220901.IncreasePlanQuotaResponseParams{
			DealName:  ptrStringIncreasePlanQuota("20231025012345678901234"),
			RequestId: ptrStringIncreasePlanQuota("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaIncreasePlanQuota()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "20231025012345678901234", d.Get("deal_name"))
}

// TestIncreasePlanQuotaOperation_APIError tests API error handling
func TestIncreasePlanQuotaOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaIncreasePlanQuota().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuotaWithContext", func(_ context.Context, request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=OperationDenied.PlanIncreasePlanQuotaUnsupported, Message=Plan increase quota unsupported")
	})

	meta := newMockMetaIncreasePlanQuota()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-invalid",
		"quota_type":   "site",
		"quota_number": 1,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OperationDenied")
}

// TestIncreasePlanQuotaOperation_Read tests Read is no-op
func TestIncreasePlanQuotaOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 1,
	})
	d.SetId("test-id")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestIncreasePlanQuotaOperation_Delete tests Delete is no-op
func TestIncreasePlanQuotaOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 1,
	})
	d.SetId("test-id")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestIncreasePlanQuotaOperation_Schema validates schema definition
func TestIncreasePlanQuotaOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "plan_id")
	assert.Contains(t, res.Schema, "quota_type")
	assert.Contains(t, res.Schema, "quota_number")
	assert.Contains(t, res.Schema, "deal_name")

	planId := res.Schema["plan_id"]
	assert.Equal(t, schema.TypeString, planId.Type)
	assert.True(t, planId.Required)
	assert.True(t, planId.ForceNew)

	quotaType := res.Schema["quota_type"]
	assert.Equal(t, schema.TypeString, quotaType.Type)
	assert.True(t, quotaType.Required)
	assert.True(t, quotaType.ForceNew)

	quotaNumber := res.Schema["quota_number"]
	assert.Equal(t, schema.TypeInt, quotaNumber.Type)
	assert.True(t, quotaNumber.Required)
	assert.True(t, quotaNumber.ForceNew)

	dealName := res.Schema["deal_name"]
	assert.Equal(t, schema.TypeString, dealName.Type)
	assert.True(t, dealName.Computed)
	assert.False(t, dealName.Optional)
	assert.False(t, dealName.Required)
}

// TestIncreasePlanQuotaOperation_ResponseNil tests nil response handling
func TestIncreasePlanQuotaOperation_ResponseNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaIncreasePlanQuota().client, "UseTeoV20220901Client", teoClient)

	// Return nil response - this should trigger retry error, but since we mock retry behavior,
	// we simulate the case where response.Response is nil
	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuotaWithContext", func(_ context.Context, request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		// Return response with nil Response field
		resp := &teov20220901.IncreasePlanQuotaResponse{}
		return resp, nil
	})

	meta := newMockMetaIncreasePlanQuota()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 1,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestIncreasePlanQuotaOperation_DealNameNil tests nil DealName handling
func TestIncreasePlanQuotaOperation_DealNameNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaIncreasePlanQuota().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuotaWithContext", func(_ context.Context, request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		resp := teov20220901.NewIncreasePlanQuotaResponse()
		resp.Response = &teov20220901.IncreasePlanQuotaResponseParams{
			DealName:  nil,
			RequestId: ptrStringIncreasePlanQuota("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaIncreasePlanQuota()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	// deal_name should be empty since DealName was nil in response
	assert.Empty(t, d.Get("deal_name"))
}

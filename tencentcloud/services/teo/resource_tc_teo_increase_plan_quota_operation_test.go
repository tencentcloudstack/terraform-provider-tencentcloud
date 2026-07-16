package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestIncreasePlanQuota" -v -count=1 -gcflags="all=-l"
// TestIncreasePlanQuota_Success tests successful quota increase
func TestIncreasePlanQuota_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuota", func(request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		assert.Equal(t, "edgeone-2unuvzjmmn2q", *request.PlanId)
		assert.Equal(t, "site", *request.QuotaType)
		assert.Equal(t, int64(10), *request.QuotaNumber)

		resp := teov20220901.NewIncreasePlanQuotaResponse()
		resp.Response = &teov20220901.IncreasePlanQuotaResponseParams{
			DealName:  ptrString("202312290001"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 10,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "202312290001", d.Get("deal_name"))
}

// TestIncreasePlanQuota_APIError tests handling API errors
func TestIncreasePlanQuota_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuota", func(request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid plan_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "invalid-plan",
		"quota_type":   "site",
		"quota_number": 10,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestIncreasePlanQuota_NilResponse tests handling nil response
func TestIncreasePlanQuota_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuota", func(request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		return nil, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 10,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

// TestIncreasePlanQuota_NilDealName tests handling nil DealName in response
func TestIncreasePlanQuota_NilDealName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuota", func(request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		resp := teov20220901.NewIncreasePlanQuotaResponse()
		resp.Response = &teov20220901.IncreasePlanQuotaResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 10,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

// TestIncreasePlanQuota_MissingPlanId tests error handling when plan_id is missing
func TestIncreasePlanQuota_MissingPlanId(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"quota_type":   "site",
		"quota_number": 10,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plan_id is required")
}

// TestIncreasePlanQuota_MissingQuotaType tests error handling when quota_type is missing
func TestIncreasePlanQuota_MissingQuotaType(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_number": 10,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quota_type is required")
}

// TestIncreasePlanQuota_Read tests Read is no-op
func TestIncreasePlanQuota_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 10,
	})
	d.SetId("test-id")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestIncreasePlanQuota_Delete tests Delete is no-op
func TestIncreasePlanQuota_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoIncreasePlanQuotaOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_id":      "edgeone-2unuvzjmmn2q",
		"quota_type":   "site",
		"quota_number": 10,
	})
	d.SetId("test-id")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestIncreasePlanQuota_Schema validates schema definition
func TestIncreasePlanQuota_Schema(t *testing.T) {
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
}

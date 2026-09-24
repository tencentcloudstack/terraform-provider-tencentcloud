package config_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcconfig "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

type mockMetaForAggregateConfigDeliver struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForAggregateConfigDeliver) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForAggregateConfigDeliver{}

func newMockMetaForAggregateConfigDeliver() *mockMetaForAggregateConfigDeliver {
	return &mockMetaForAggregateConfigDeliver{client: &connectivity.TencentCloudClient{}}
}

func ptrStrAggDeliver(s string) *string    { return &s }
func ptrUint64AggDeliver(v uint64) *uint64 { return &v }
func ptrInt64AggDeliver(v int64) *int64    { return &v }

// TestAggregateConfigDeliver_Schema_GuardsPermanentDiff asserts the schema flags that
// keep `plan` convergent. Two distinct defects are guarded here:
//
//  1. The optional fields are populated by the cloud even when the configuration omits
//     them (`deliver_type` defaults to COS, `deliver_content_type` to 1, `deliver_uin`
//     to 0). Without `Computed` the populated state conflicts with the null
//     configuration, producing an "update in-place" diff on every plan that never
//     converges. This is the most likely defect for a user to hit, because it only
//     requires omitting an optional argument.
//  2. `account_group_id` must stay `ForceNew` and is written back by Read (see
//     TestAggregateConfigDeliver_Read_ImportPopulatesAccountGroupId).
func TestAggregateConfigDeliver_Schema_GuardsPermanentDiff(t *testing.T) {
	res := svcconfig.ResourceTencentCloudConfigUpdateAggregateConfigDeliver()

	// The cloud returns a default for each of these, so each must accept a
	// cloud-populated value even though the user may never set it.
	for _, name := range []string{
		"deliver_name",
		"target_arn",
		"deliver_prefix",
		"deliver_type",
		"deliver_uin",
		"deliver_content_type",
	} {
		s, ok := res.Schema[name]
		assert.True(t, ok, "%s should exist in schema", name)
		assert.True(t, s.Optional, "%s should be Optional", name)
		assert.True(t, s.Computed,
			"%s must be Computed so a cloud-populated value does not cause a permanent diff", name)
	}

	assert.True(t, res.Schema["account_group_id"].Required)
	assert.True(t, res.Schema["account_group_id"].ForceNew)
	assert.True(t, res.Schema["status"].Required)
}

// TestAggregateConfigDeliver_Read_ImportPopulatesAccountGroupId covers the import path.
// `account_group_id` is ForceNew but is NOT part of the Describe response, so during
// `terraform import` the state only carries the imported id. If Read does not write the
// field back, the next plan reports `must be replaced` (ForceNew false-positive).
func TestAggregateConfigDeliver_Read_ImportPopulatesAccountGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaForAggregateConfigDeliver().client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "DescribeAggregateConfigDeliver", func(request *configv20220802.DescribeAggregateConfigDeliverRequest) (*configv20220802.DescribeAggregateConfigDeliverResponse, error) {
		return mockDescribeAggregateConfigDeliverResponse("ca-ag-import1234"), nil
	})

	meta := newMockMetaForAggregateConfigDeliver()
	res := svcconfig.ResourceTencentCloudConfigUpdateAggregateConfigDeliver()

	// Simulate a plain `terraform import <addr> ca-ag-import1234`: the id is the only
	// thing available, no other attribute is known yet.
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("ca-ag-import1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ca-ag-import1234", d.Id())
	assert.Equal(t, "ca-ag-import1234", d.Get("account_group_id"),
		"Read must write account_group_id back, otherwise an imported resource is replaced on the next plan")
}

func mockDescribeAggregateConfigDeliverResponse(accountGroupId string) *configv20220802.DescribeAggregateConfigDeliverResponse {
	resp := configv20220802.NewDescribeAggregateConfigDeliverResponse()
	resp.Response = &configv20220802.DescribeAggregateConfigDeliverResponseParams{
		DeliverName:        ptrStrAggDeliver("tf-example-aggregate-deliver"),
		TargetArn:          ptrStrAggDeliver("qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"),
		Status:             ptrUint64AggDeliver(1),
		CreateTime:         ptrStrAggDeliver("2024-01-01 00:00:00"),
		DeliverPrefix:      ptrStrAggDeliver("config"),
		DeliverType:        ptrStrAggDeliver("COS"),
		DeliverUin:         ptrInt64AggDeliver(0),
		DeliverContentType: ptrUint64AggDeliver(3),
		RequestId:          ptrStrAggDeliver("fake-request-id"),
	}
	return resp
}

// TestAggregateConfigDeliver_Create_Success verifies Create sets the account_group_id as id and delegates to Update.
func TestAggregateConfigDeliver_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaForAggregateConfigDeliver().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *configv20220802.UpdateAggregateConfigDeliverRequest
	patches.ApplyMethodFunc(configClient, "UpdateAggregateConfigDeliverWithContext", func(ctx context.Context, request *configv20220802.UpdateAggregateConfigDeliverRequest) (*configv20220802.UpdateAggregateConfigDeliverResponse, error) {
		capturedRequest = request
		resp := configv20220802.NewUpdateAggregateConfigDeliverResponse()
		resp.Response = &configv20220802.UpdateAggregateConfigDeliverResponseParams{
			RequestId: ptrStrAggDeliver("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(configClient, "DescribeAggregateConfigDeliver", func(request *configv20220802.DescribeAggregateConfigDeliverRequest) (*configv20220802.DescribeAggregateConfigDeliverResponse, error) {
		return mockDescribeAggregateConfigDeliverResponse("ca-ag-test1234"), nil
	})

	meta := newMockMetaForAggregateConfigDeliver()
	res := svcconfig.ResourceTencentCloudConfigUpdateAggregateConfigDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"account_group_id":     "ca-ag-test1234",
		"status":               1,
		"deliver_name":         "tf-example-aggregate-deliver",
		"target_arn":           "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket",
		"deliver_prefix":       "config",
		"deliver_type":         "COS",
		"deliver_uin":          0,
		"deliver_content_type": 3,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ca-ag-test1234", d.Id())

	assert.NotNil(t, capturedRequest.AccountGroupId)
	assert.Equal(t, "ca-ag-test1234", *capturedRequest.AccountGroupId)
	assert.NotNil(t, capturedRequest.Status)
	assert.Equal(t, uint64(1), *capturedRequest.Status)
	assert.NotNil(t, capturedRequest.DeliverType)
	assert.Equal(t, "COS", *capturedRequest.DeliverType)
	assert.NotNil(t, capturedRequest.DeliverUin)
	assert.Equal(t, int64(0), *capturedRequest.DeliverUin)
	assert.NotNil(t, capturedRequest.DeliverContentType)
	assert.Equal(t, uint64(3), *capturedRequest.DeliverContentType)
}

// TestAggregateConfigDeliver_Read_Success verifies Read populates fields from Describe response.
func TestAggregateConfigDeliver_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaForAggregateConfigDeliver().client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "DescribeAggregateConfigDeliver", func(request *configv20220802.DescribeAggregateConfigDeliverRequest) (*configv20220802.DescribeAggregateConfigDeliverResponse, error) {
		return mockDescribeAggregateConfigDeliverResponse("ca-ag-read1234"), nil
	})

	meta := newMockMetaForAggregateConfigDeliver()
	res := svcconfig.ResourceTencentCloudConfigUpdateAggregateConfigDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"account_group_id": "ca-ag-read1234",
		"status":           1,
	})
	d.SetId("ca-ag-read1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ca-ag-read1234", d.Id())
	assert.Equal(t, 1, d.Get("status"))
	assert.Equal(t, "tf-example-aggregate-deliver", d.Get("deliver_name"))
	assert.Equal(t, "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket", d.Get("target_arn"))
	assert.Equal(t, "config", d.Get("deliver_prefix"))
	assert.Equal(t, "COS", d.Get("deliver_type"))
	assert.Equal(t, 0, d.Get("deliver_uin"))
	assert.Equal(t, 3, d.Get("deliver_content_type"))
	assert.Equal(t, "2024-01-01 00:00:00", d.Get("create_time"))
}

// TestAggregateConfigDeliver_Update_Success verifies Update sends correct params.
func TestAggregateConfigDeliver_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaForAggregateConfigDeliver().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *configv20220802.UpdateAggregateConfigDeliverRequest
	patches.ApplyMethodFunc(configClient, "UpdateAggregateConfigDeliverWithContext", func(ctx context.Context, request *configv20220802.UpdateAggregateConfigDeliverRequest) (*configv20220802.UpdateAggregateConfigDeliverResponse, error) {
		capturedRequest = request
		resp := configv20220802.NewUpdateAggregateConfigDeliverResponse()
		resp.Response = &configv20220802.UpdateAggregateConfigDeliverResponseParams{
			RequestId: ptrStrAggDeliver("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(configClient, "DescribeAggregateConfigDeliver", func(request *configv20220802.DescribeAggregateConfigDeliverRequest) (*configv20220802.DescribeAggregateConfigDeliverResponse, error) {
		return mockDescribeAggregateConfigDeliverResponse("ca-ag-upd1234"), nil
	})

	meta := newMockMetaForAggregateConfigDeliver()
	res := svcconfig.ResourceTencentCloudConfigUpdateAggregateConfigDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"account_group_id":     "ca-ag-upd1234",
		"status":               0,
		"deliver_content_type": 1,
	})
	d.SetId("ca-ag-upd1234")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest.AccountGroupId)
	assert.Equal(t, "ca-ag-upd1234", *capturedRequest.AccountGroupId)
	assert.NotNil(t, capturedRequest.Status)
	assert.Equal(t, uint64(0), *capturedRequest.Status)
	assert.NotNil(t, capturedRequest.DeliverContentType)
	assert.Equal(t, uint64(1), *capturedRequest.DeliverContentType)
}

// TestAggregateConfigDeliver_Delete_Success verifies Delete is a no-op and returns nil.
func TestAggregateConfigDeliver_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForAggregateConfigDeliver()
	res := svcconfig.ResourceTencentCloudConfigUpdateAggregateConfigDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"account_group_id": "ca-ag-del1234",
		"status":           1,
	})
	d.SetId("ca-ag-del1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

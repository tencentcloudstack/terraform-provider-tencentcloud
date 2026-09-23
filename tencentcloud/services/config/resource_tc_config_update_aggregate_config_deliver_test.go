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

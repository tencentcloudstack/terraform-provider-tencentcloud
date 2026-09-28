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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

type mockMetaForConfigListAggregatorsDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForConfigListAggregatorsDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForConfigListAggregatorsDS{}

func newMockMetaForConfigListAggregatorsDS() *mockMetaForConfigListAggregatorsDS {
	return &mockMetaForConfigListAggregatorsDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCLA(s string) *string { return &s }

func ptrUint64CLA(v uint64) *uint64 { return &v }

func buildAggregator(name, description string, ownerUin, accountCount uint64, createTime, aggregatorType, accountGroupId string, aggregatorStatus uint64, memberName string) *configv20220802.Aggregator {
	return &configv20220802.Aggregator{
		Name:             ptrStringCLA(name),
		Description:      ptrStringCLA(description),
		OwnerUin:         ptrUint64CLA(ownerUin),
		CreateTime:       ptrStringCLA(createTime),
		AccountCount:     ptrUint64CLA(accountCount),
		Type:             ptrStringCLA(aggregatorType),
		AccountGroupId:   ptrStringCLA(accountGroupId),
		AggregatorStatus: ptrUint64CLA(aggregatorStatus),
		MemberName:       ptrStringCLA(memberName),
	}
}

func TestConfigListAggregatorsDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigListAggregatorsDS().client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregatorsWithContext", func(_ context.Context, request *configv20220802.ListAggregatorsRequest) (*configv20220802.ListAggregatorsResponse, error) {
		resp := configv20220802.NewListAggregatorsResponse()
		resp.Response = &configv20220802.ListAggregatorsResponseParams{
			Total: ptrUint64CLA(1),
			Items: []*configv20220802.Aggregator{
				buildAggregator(
					"tf-aggregator",
					"terraform test aggregator",
					100000000000,
					2,
					"2024-01-01 00:00:00",
					"RD",
					"ca-aggregator-abcdefg",
					1,
					"member-name",
				),
			},
			RequestId: ptrStringCLA("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForConfigListAggregatorsDS()
	res := config.DataSourceTencentCloudConfigListAggregators()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	total := d.Get("total").(int)
	assert.Equal(t, 1, total)

	items := d.Get("items").([]interface{})
	assert.Len(t, items, 1)

	item0 := items[0].(map[string]interface{})
	assert.Equal(t, "tf-aggregator", item0["name"].(string))
	assert.Equal(t, "terraform test aggregator", item0["description"].(string))
	assert.Equal(t, 100000000000, item0["owner_uin"].(int))
	assert.Equal(t, "2024-01-01 00:00:00", item0["create_time"].(string))
	assert.Equal(t, 2, item0["account_count"].(int))
	assert.Equal(t, "RD", item0["type"].(string))
	assert.Equal(t, "ca-aggregator-abcdefg", item0["account_group_id"].(string))
	assert.Equal(t, 1, item0["aggregator_status"].(int))
	assert.Equal(t, "member-name", item0["member_name"].(string))
}

func TestConfigListAggregatorsDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigListAggregatorsDS().client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregatorsWithContext", func(_ context.Context, request *configv20220802.ListAggregatorsRequest) (*configv20220802.ListAggregatorsResponse, error) {
		resp := configv20220802.NewListAggregatorsResponse()
		resp.Response = &configv20220802.ListAggregatorsResponseParams{
			Total:     ptrUint64CLA(0),
			Items:     []*configv20220802.Aggregator{},
			RequestId: ptrStringCLA("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForConfigListAggregatorsDS()
	res := config.DataSourceTencentCloudConfigListAggregators()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

func TestConfigListAggregatorsDS_Schema(t *testing.T) {
	res := config.DataSourceTencentCloudConfigListAggregators()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "total")
	assert.Contains(t, res.Schema, "items")
	assert.Contains(t, res.Schema, "result_output_file")

	totalSchema := res.Schema["total"]
	assert.Equal(t, schema.TypeInt, totalSchema.Type)
	assert.True(t, totalSchema.Computed)

	itemsSchema := res.Schema["items"]
	assert.Equal(t, schema.TypeList, itemsSchema.Type)
	assert.True(t, itemsSchema.Computed)

	elemRes := itemsSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "name")
	assert.Contains(t, elemRes.Schema, "description")
	assert.Contains(t, elemRes.Schema, "owner_uin")
	assert.Contains(t, elemRes.Schema, "create_time")
	assert.Contains(t, elemRes.Schema, "account_count")
	assert.Contains(t, elemRes.Schema, "type")
	assert.Contains(t, elemRes.Schema, "account_group_id")
	assert.Contains(t, elemRes.Schema, "aggregator_status")
	assert.Contains(t, elemRes.Schema, "member_name")

	outputSchema := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputSchema.Type)
	assert.True(t, outputSchema.Optional)
}

package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoContentQuota" -v -count=1 -gcflags="all=-l"

func TestTeoContentQuotaDataSource_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeContentQuota", func(request *teov20220901.DescribeContentQuotaRequest) (*teov20220901.DescribeContentQuotaResponse, error) {
		resp := teov20220901.NewDescribeContentQuotaResponse()
		resp.Response = &teov20220901.DescribeContentQuotaResponseParams{
			PurgeQuota: []*teov20220901.Quota{
				{
					Batch:          contentQuotaPtrInt64(1000),
					Daily:          contentQuotaPtrInt64(10000),
					DailyAvailable: contentQuotaPtrInt64(8000),
					Type:           contentQuotaPtrString("purge_url"),
				},
				{
					Batch:          contentQuotaPtrInt64(1000),
					Daily:          contentQuotaPtrInt64(10000),
					DailyAvailable: contentQuotaPtrInt64(5000),
					Type:           contentQuotaPtrString("purge_prefix"),
				},
			},
			PrefetchQuota: []*teov20220901.Quota{
				{
					Batch:          contentQuotaPtrInt64(500),
					Daily:          contentQuotaPtrInt64(5000),
					DailyAvailable: contentQuotaPtrInt64(3000),
					Type:           contentQuotaPtrString("prefetch_url"),
				},
			},
			RequestId: contentQuotaPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &contentQuotaMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoContentQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2l1zk57u",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	purgeQuota := d.Get("purge_quota").([]interface{})
	assert.Len(t, purgeQuota, 2)

	purgeQuota0 := purgeQuota[0].(map[string]interface{})
	assert.Equal(t, 1000, purgeQuota0["batch"])
	assert.Equal(t, 10000, purgeQuota0["daily"])
	assert.Equal(t, 8000, purgeQuota0["daily_available"])
	assert.Equal(t, "purge_url", purgeQuota0["type"])

	purgeQuota1 := purgeQuota[1].(map[string]interface{})
	assert.Equal(t, "purge_prefix", purgeQuota1["type"])

	prefetchQuota := d.Get("prefetch_quota").([]interface{})
	assert.Len(t, prefetchQuota, 1)

	prefetchQuota0 := prefetchQuota[0].(map[string]interface{})
	assert.Equal(t, 500, prefetchQuota0["batch"])
	assert.Equal(t, 5000, prefetchQuota0["daily"])
	assert.Equal(t, 3000, prefetchQuota0["daily_available"])
	assert.Equal(t, "prefetch_url", prefetchQuota0["type"])
}

func TestTeoContentQuotaDataSource_NilQuotaFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeContentQuota", func(request *teov20220901.DescribeContentQuotaRequest) (*teov20220901.DescribeContentQuotaResponse, error) {
		resp := teov20220901.NewDescribeContentQuotaResponse()
		resp.Response = &teov20220901.DescribeContentQuotaResponseParams{
			PurgeQuota: []*teov20220901.Quota{
				{
					Type: contentQuotaPtrString("purge_all"),
				},
			},
			PrefetchQuota: nil,
			RequestId:     contentQuotaPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &contentQuotaMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoContentQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2l1zk57u",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	purgeQuota := d.Get("purge_quota").([]interface{})
	assert.Len(t, purgeQuota, 1)
	purgeQuota0 := purgeQuota[0].(map[string]interface{})
	assert.Equal(t, "purge_all", purgeQuota0["type"])

	prefetchQuota := d.Get("prefetch_quota").([]interface{})
	assert.Len(t, prefetchQuota, 0)
}

func TestTeoContentQuotaDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoContentQuota()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "purge_quota")
	assert.Contains(t, res.Schema, "prefetch_quota")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	purgeQuota := res.Schema["purge_quota"]
	assert.Equal(t, schema.TypeList, purgeQuota.Type)
	assert.True(t, purgeQuota.Computed)

	prefetchQuota := res.Schema["prefetch_quota"]
	assert.Equal(t, schema.TypeList, prefetchQuota.Type)
	assert.True(t, prefetchQuota.Computed)

	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)
}

type contentQuotaMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *contentQuotaMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &contentQuotaMockMeta{}

func contentQuotaPtrInt64(i int64) *int64 {
	return &i
}

func contentQuotaPtrString(s string) *string {
	return &s
}

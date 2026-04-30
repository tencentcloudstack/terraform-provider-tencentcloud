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

// go test ./tencentcloud/services/teo/ -run "TestTeoContentQuotaDataSource" -v -count=1 -gcflags="all=-l"

type mockMetaContentQuota struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaContentQuota) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaContentQuota{}

func newMockMetaContentQuota() *mockMetaContentQuota {
	return &mockMetaContentQuota{client: &connectivity.TencentCloudClient{}}
}

func ptrStringContentQuota(s string) *string {
	return &s
}

func ptrInt64ContentQuota(n int64) *int64 {
	return &n
}

// TestTeoContentQuotaDataSource_ReadSuccess tests successful read with quota data
func TestTeoContentQuotaDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaContentQuota().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeContentQuota", func(request *teov20220901.DescribeContentQuotaRequest) (*teov20220901.DescribeContentQuotaResponse, error) {
		resp := teov20220901.NewDescribeContentQuotaResponse()
		resp.Response = &teov20220901.DescribeContentQuotaResponseParams{
			PurgeQuota: []*teov20220901.Quota{
				{
					Type:           ptrStringContentQuota("purge_url"),
					Batch:          ptrInt64ContentQuota(1000),
					Daily:          ptrInt64ContentQuota(10000),
					DailyAvailable: ptrInt64ContentQuota(9500),
				},
				{
					Type:           ptrStringContentQuota("purge_prefix"),
					Batch:          ptrInt64ContentQuota(100),
					Daily:          ptrInt64ContentQuota(5000),
					DailyAvailable: ptrInt64ContentQuota(4500),
				},
			},
			PrefetchQuota: []*teov20220901.Quota{
				{
					Type:           ptrStringContentQuota("prefetch_url"),
					Batch:          ptrInt64ContentQuota(500),
					Daily:          ptrInt64ContentQuota(5000),
					DailyAvailable: ptrInt64ContentQuota(4000),
				},
			},
			RequestId: ptrStringContentQuota("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaContentQuota()
	res := teo.DataSourceTencentCloudTeoContentQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	purgeQuota := d.Get("purge_quota").([]interface{})
	assert.Len(t, purgeQuota, 2)
	purgeMap := purgeQuota[0].(map[string]interface{})
	assert.Equal(t, "purge_url", purgeMap["type"])
	assert.Equal(t, 1000, purgeMap["batch"])
	assert.Equal(t, 10000, purgeMap["daily"])
	assert.Equal(t, 9500, purgeMap["daily_available"])

	prefetchQuota := d.Get("prefetch_quota").([]interface{})
	assert.Len(t, prefetchQuota, 1)
	prefetchMap := prefetchQuota[0].(map[string]interface{})
	assert.Equal(t, "prefetch_url", prefetchMap["type"])
	assert.Equal(t, 500, prefetchMap["batch"])
}

// TestTeoContentQuotaDataSource_ReadNullQuota tests read with null quota fields
func TestTeoContentQuotaDataSource_ReadNullQuota(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaContentQuota().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeContentQuota", func(request *teov20220901.DescribeContentQuotaRequest) (*teov20220901.DescribeContentQuotaResponse, error) {
		resp := teov20220901.NewDescribeContentQuotaResponse()
		resp.Response = &teov20220901.DescribeContentQuotaResponseParams{
			PurgeQuota:    nil,
			PrefetchQuota: nil,
			RequestId:     ptrStringContentQuota("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaContentQuota()
	res := teo.DataSourceTencentCloudTeoContentQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	purgeQuota := d.Get("purge_quota").([]interface{})
	assert.Len(t, purgeQuota, 0)

	prefetchQuota := d.Get("prefetch_quota").([]interface{})
	assert.Len(t, prefetchQuota, 0)
}

// TestTeoContentQuotaDataSource_Schema validates schema definition
func TestTeoContentQuotaDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoContentQuota()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "purge_quota")
	assert.Contains(t, res.Schema, "prefetch_quota")
	assert.Contains(t, res.Schema, "result_output_file")

	// zone_id is Required
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	// purge_quota is Computed
	purgeQuota := res.Schema["purge_quota"]
	assert.Equal(t, schema.TypeList, purgeQuota.Type)
	assert.True(t, purgeQuota.Computed)

	// prefetch_quota is Computed
	prefetchQuota := res.Schema["prefetch_quota"]
	assert.Equal(t, schema.TypeList, prefetchQuota.Type)
	assert.True(t, prefetchQuota.Computed)

	// result_output_file is Optional
	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)
}

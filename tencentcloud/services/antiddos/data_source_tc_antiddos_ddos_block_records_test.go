package antiddos_test

import (
	"log"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	antiddosv20250903 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"
)

type mockMetaAntiddosDDoSBlockRecordsDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaAntiddosDDoSBlockRecordsDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaAntiddosDDoSBlockRecordsDS{}

func newMockMetaAntiddosDDoSBlockRecordsDS() *mockMetaAntiddosDDoSBlockRecordsDS {
	return &mockMetaAntiddosDDoSBlockRecordsDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrBlockRecords(s string) *string {
	return &s
}

func ptrUint64BlockRecords(u uint64) *uint64 {
	return &u
}

// go test ./tencentcloud/services/antiddos/ -run "TestAntiddosDDoSBlockRecordsDS" -v -count=1 -gcflags="all=-l"

func TestAntiddosDDoSBlockRecordsDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	antiddosClient := &antiddosv20250903.Client{}
	patches.ApplyMethodReturn(newMockMetaAntiddosDDoSBlockRecordsDS().client, "UseAntiddosV20250903Client", antiddosClient)

	patches.ApplyMethodFunc(antiddosClient, "DescribeDDoSBlockRecords", func(request *antiddosv20250903.DescribeDDoSBlockRecordsRequest) (*antiddosv20250903.DescribeDDoSBlockRecordsResponse, error) {
		assert.NotNil(t, request.StartTime)
		assert.NotNil(t, request.EndTime)
		assert.Equal(t, "2026-02-04T11:30:00+08:00", *request.StartTime)
		assert.Equal(t, "2026-03-04T11:30:00+08:00", *request.EndTime)

		resp := antiddosv20250903.NewDescribeDDoSBlockRecordsResponse()
		resp.Response = &antiddosv20250903.DescribeDDoSBlockRecordsResponseParams{
			TotalCount: ptrUint64BlockRecords(2),
			BlockRecords: []*antiddosv20250903.DDoSBlockRecord{
				{
					Resource:  ptrStrBlockRecords("117.175.94.231"),
					BlockTime: ptrStrBlockRecords("2026-02-05 10:00:00"),
					Status:    ptrStrBlockRecords("Blocked"),
				},
				{
					Resource:  ptrStrBlockRecords("117.175.94.232"),
					BlockTime: ptrStrBlockRecords("2026-02-06 11:00:00"),
					Status:    ptrStrBlockRecords("Unblocked"),
				},
			},
			UnblockQuotaInfo: &antiddosv20250903.DDoSUnblockQuota{
				TotalQuota:     ptrUint64BlockRecords(100),
				UsedQuota:      ptrUint64BlockRecords(5),
				QuotaStartTime: ptrStrBlockRecords("2026-02-01 00:00:00"),
				QuotaEndTime:   ptrStrBlockRecords("2026-03-01 00:00:00"),
			},
		}
		return resp, nil
	})

	meta := newMockMetaAntiddosDDoSBlockRecordsDS()
	res := svcantiddos.DataSourceTencentCloudAntiddosDDoSBlockRecords()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time": "2026-02-04T11:30:00+08:00",
		"end_time":   "2026-03-04T11:30:00+08:00",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	blockRecords := d.Get("block_records").([]interface{})
	assert.Len(t, blockRecords, 2)

	record0 := blockRecords[0].(map[string]interface{})
	assert.Equal(t, "117.175.94.231", record0["resource"].(string))
	assert.Equal(t, "2026-02-05 10:00:00", record0["block_time"].(string))
	assert.Equal(t, "Blocked", record0["status"].(string))

	record1 := blockRecords[1].(map[string]interface{})
	assert.Equal(t, "117.175.94.232", record1["resource"].(string))
	assert.Equal(t, "Unblocked", record1["status"].(string))

	unblockQuotaInfo := d.Get("unblock_quota_info").([]interface{})
	assert.Len(t, unblockQuotaInfo, 1)
	quotaMap := unblockQuotaInfo[0].(map[string]interface{})
	assert.Equal(t, 100, quotaMap["total_quota"].(int))
	assert.Equal(t, 5, quotaMap["used_quota"].(int))
	assert.Equal(t, "2026-02-01 00:00:00", quotaMap["quota_start_time"].(string))
	assert.Equal(t, "2026-03-01 00:00:00", quotaMap["quota_end_time"].(string))
}

func TestAntiddosDDoSBlockRecordsDS_ReadWithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	antiddosClient := &antiddosv20250903.Client{}
	patches.ApplyMethodReturn(newMockMetaAntiddosDDoSBlockRecordsDS().client, "UseAntiddosV20250903Client", antiddosClient)

	patches.ApplyMethodFunc(antiddosClient, "DescribeDDoSBlockRecords", func(request *antiddosv20250903.DescribeDDoSBlockRecordsRequest) (*antiddosv20250903.DescribeDDoSBlockRecordsResponse, error) {
		assert.NotNil(t, request.Filters)
		assert.Len(t, request.Filters, 1)
		assert.Equal(t, "Status", *request.Filters[0].Name)
		assert.Equal(t, "Blocked", *request.Filters[0].Values[0])

		resp := antiddosv20250903.NewDescribeDDoSBlockRecordsResponse()
		resp.Response = &antiddosv20250903.DescribeDDoSBlockRecordsResponseParams{
			TotalCount: ptrUint64BlockRecords(1),
			BlockRecords: []*antiddosv20250903.DDoSBlockRecord{
				{
					Resource:  ptrStrBlockRecords("117.175.94.231"),
					BlockTime: ptrStrBlockRecords("2026-02-05 10:00:00"),
					Status:    ptrStrBlockRecords("Blocked"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaAntiddosDDoSBlockRecordsDS()
	res := svcantiddos.DataSourceTencentCloudAntiddosDDoSBlockRecords()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time": "2026-02-04T11:30:00+08:00",
		"end_time":   "2026-03-04T11:30:00+08:00",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "Status",
				"values": schema.NewSet(schema.HashString, []interface{}{"Blocked"}),
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	blockRecords := d.Get("block_records").([]interface{})
	assert.Len(t, blockRecords, 1)

	record0 := blockRecords[0].(map[string]interface{})
	assert.Equal(t, "Blocked", record0["status"].(string))
}

func TestAntiddosDDoSBlockRecordsDS_ReadWithNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	antiddosClient := &antiddosv20250903.Client{}
	patches.ApplyMethodReturn(newMockMetaAntiddosDDoSBlockRecordsDS().client, "UseAntiddosV20250903Client", antiddosClient)

	patches.ApplyMethodFunc(antiddosClient, "DescribeDDoSBlockRecords", func(request *antiddosv20250903.DescribeDDoSBlockRecordsRequest) (*antiddosv20250903.DescribeDDoSBlockRecordsResponse, error) {
		resp := antiddosv20250903.NewDescribeDDoSBlockRecordsResponse()
		resp.Response = &antiddosv20250903.DescribeDDoSBlockRecordsResponseParams{
			BlockRecords: []*antiddosv20250903.DDoSBlockRecord{
				{
					Resource: ptrStrBlockRecords("117.175.94.233"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaAntiddosDDoSBlockRecordsDS()
	res := svcantiddos.DataSourceTencentCloudAntiddosDDoSBlockRecords()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time": "2026-02-04T11:30:00+08:00",
		"end_time":   "2026-03-04T11:30:00+08:00",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	blockRecords := d.Get("block_records").([]interface{})
	assert.Len(t, blockRecords, 1)

	record0 := blockRecords[0].(map[string]interface{})
	assert.Equal(t, "117.175.94.233", record0["resource"].(string))
}

func TestAntiddosDDoSBlockRecordsDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	antiddosClient := &antiddosv20250903.Client{}
	patches.ApplyMethodReturn(newMockMetaAntiddosDDoSBlockRecordsDS().client, "UseAntiddosV20250903Client", antiddosClient)

	patches.ApplyMethodFunc(antiddosClient, "DescribeDDoSBlockRecords", func(request *antiddosv20250903.DescribeDDoSBlockRecordsRequest) (*antiddosv20250903.DescribeDDoSBlockRecordsResponse, error) {
		resp := antiddosv20250903.NewDescribeDDoSBlockRecordsResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaAntiddosDDoSBlockRecordsDS()
	res := svcantiddos.DataSourceTencentCloudAntiddosDDoSBlockRecords()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time": "2026-02-04T11:30:00+08:00",
		"end_time":   "2026-03-04T11:30:00+08:00",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

func TestAntiddosDDoSBlockRecordsDS_ReadWithPagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	antiddosClient := &antiddosv20250903.Client{}
	patches.ApplyMethodReturn(newMockMetaAntiddosDDoSBlockRecordsDS().client, "UseAntiddosV20250903Client", antiddosClient)

	callCount := 0
	patches.ApplyMethodFunc(antiddosClient, "DescribeDDoSBlockRecords", func(request *antiddosv20250903.DescribeDDoSBlockRecordsRequest) (*antiddosv20250903.DescribeDDoSBlockRecordsResponse, error) {
		callCount++
		resp := antiddosv20250903.NewDescribeDDoSBlockRecordsResponse()
		if callCount == 1 {
			assert.Equal(t, uint64(0), *request.Offset)
			records := make([]*antiddosv20250903.DDoSBlockRecord, 0, 100)
			for i := 0; i < 100; i++ {
				records = append(records, &antiddosv20250903.DDoSBlockRecord{
					Resource:  ptrStrBlockRecords("117.175.94.231"),
					BlockTime: ptrStrBlockRecords("2026-02-05 10:00:00"),
					Status:    ptrStrBlockRecords("Blocked"),
				})
			}
			resp.Response = &antiddosv20250903.DescribeDDoSBlockRecordsResponseParams{
				TotalCount:   ptrUint64BlockRecords(101),
				BlockRecords: records,
			}
		} else {
			assert.Equal(t, uint64(100), *request.Offset)
			resp.Response = &antiddosv20250903.DescribeDDoSBlockRecordsResponseParams{
				TotalCount: ptrUint64BlockRecords(101),
				BlockRecords: []*antiddosv20250903.DDoSBlockRecord{
					{
						Resource:  ptrStrBlockRecords("117.175.94.232"),
						BlockTime: ptrStrBlockRecords("2026-02-06 11:00:00"),
						Status:    ptrStrBlockRecords("Unblocked"),
					},
				},
			}
		}
		return resp, nil
	})

	meta := newMockMetaAntiddosDDoSBlockRecordsDS()
	res := svcantiddos.DataSourceTencentCloudAntiddosDDoSBlockRecords()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time": "2026-02-04T11:30:00+08:00",
		"end_time":   "2026-03-04T11:30:00+08:00",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, 2, callCount)

	blockRecords := d.Get("block_records").([]interface{})
	assert.Len(t, blockRecords, 101)
}

func TestAntiddosDDoSBlockRecordsDS_ReadWithResultOutputFile(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	antiddosClient := &antiddosv20250903.Client{}
	patches.ApplyMethodReturn(newMockMetaAntiddosDDoSBlockRecordsDS().client, "UseAntiddosV20250903Client", antiddosClient)

	patches.ApplyMethodFunc(antiddosClient, "DescribeDDoSBlockRecords", func(request *antiddosv20250903.DescribeDDoSBlockRecordsRequest) (*antiddosv20250903.DescribeDDoSBlockRecordsResponse, error) {
		resp := antiddosv20250903.NewDescribeDDoSBlockRecordsResponse()
		resp.Response = &antiddosv20250903.DescribeDDoSBlockRecordsResponseParams{
			BlockRecords: []*antiddosv20250903.DDoSBlockRecord{
				{
					Resource:  ptrStrBlockRecords("117.175.94.231"),
					BlockTime: ptrStrBlockRecords("2026-02-05 10:00:00"),
					Status:    ptrStrBlockRecords("Blocked"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaAntiddosDDoSBlockRecordsDS()
	res := svcantiddos.DataSourceTencentCloudAntiddosDDoSBlockRecords()
	outputFile := "./test_output_block_records.json"
	defer func() {
		if err := os.Remove(outputFile); err != nil {
			log.Printf("failed to remove %s: %v", outputFile, err)
		}
	}()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time":         "2026-02-04T11:30:00+08:00",
		"end_time":           "2026-03-04T11:30:00+08:00",
		"result_output_file": outputFile,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	_, e := os.Stat(outputFile)
	assert.NoError(t, e)
}

func TestAntiddosDDoSBlockRecordsDS_Schema(t *testing.T) {
	res := svcantiddos.DataSourceTencentCloudAntiddosDDoSBlockRecords()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "start_time")
	assert.Contains(t, res.Schema, "end_time")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "block_records")
	assert.Contains(t, res.Schema, "unblock_quota_info")
	assert.Contains(t, res.Schema, "result_output_file")

	startTimeSchema := res.Schema["start_time"]
	assert.Equal(t, schema.TypeString, startTimeSchema.Type)
	assert.True(t, startTimeSchema.Required)

	endTimeSchema := res.Schema["end_time"]
	assert.Equal(t, schema.TypeString, endTimeSchema.Type)
	assert.True(t, endTimeSchema.Required)

	filtersSchema := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filtersSchema.Type)
	assert.True(t, filtersSchema.Optional)
	filtersRes := filtersSchema.Elem.(*schema.Resource)
	assert.Contains(t, filtersRes.Schema, "name")
	assert.Contains(t, filtersRes.Schema, "values")

	blockRecordsSchema := res.Schema["block_records"]
	assert.Equal(t, schema.TypeList, blockRecordsSchema.Type)
	assert.True(t, blockRecordsSchema.Computed)
	blockRecordsRes := blockRecordsSchema.Elem.(*schema.Resource)
	assert.Contains(t, blockRecordsRes.Schema, "resource")
	assert.Contains(t, blockRecordsRes.Schema, "block_time")
	assert.Contains(t, blockRecordsRes.Schema, "status")

	unblockQuotaInfoSchema := res.Schema["unblock_quota_info"]
	assert.Equal(t, schema.TypeList, unblockQuotaInfoSchema.Type)
	assert.True(t, unblockQuotaInfoSchema.Computed)
	quotaRes := unblockQuotaInfoSchema.Elem.(*schema.Resource)
	assert.Contains(t, quotaRes.Schema, "total_quota")
	assert.Contains(t, quotaRes.Schema, "used_quota")
	assert.Contains(t, quotaRes.Schema, "quota_start_time")
	assert.Contains(t, quotaRes.Schema, "quota_end_time")
}

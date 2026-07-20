package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoBillingDataDataSource" -v -count=1 -gcflags="all=-l"

type mockMetaBillingData struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBillingData) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBillingData{}

func newMockMetaBillingData() *mockMetaBillingData {
	return &mockMetaBillingData{client: &connectivity.TencentCloudClient{}}
}

func ptrStringBillingData(s string) *string {
	return &s
}

func ptrUint64BillingData(n uint64) *uint64 {
	return &n
}

// TestTeoBillingDataDataSource_ReadSuccess tests successful read with billing data
func TestTeoBillingDataDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingData().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeBillingDataWithContext", func(_ context.Context, request *teov20220901.DescribeBillingDataRequest) (*teov20220901.DescribeBillingDataResponse, error) {
		assert.Equal(t, "2025-01-01T00:00:00+08:00", *request.StartTime)
		assert.Equal(t, "2025-01-02T00:00:00+08:00", *request.EndTime)
		assert.Equal(t, "acc_flux", *request.MetricName)
		assert.Equal(t, "hour", *request.Interval)

		resp := teov20220901.NewDescribeBillingDataResponse()
		resp.Response = &teov20220901.DescribeBillingDataResponseParams{
			Data: []*teov20220901.BillingData{
				{
					Time:     ptrStringBillingData("2025-01-01T00:00:00+08:00"),
					Value:    ptrUint64BillingData(1024),
					ZoneId:   ptrStringBillingData("zone-2qtuhspy7cr6"),
					Host:     ptrStringBillingData("test.example.com"),
					ProxyId:  ptrStringBillingData("sid-2rugn89bkla9"),
					RegionId: ptrStringBillingData("CH"),
				},
				{
					Time:     ptrStringBillingData("2025-01-01T01:00:00+08:00"),
					Value:    ptrUint64BillingData(2048),
					ZoneId:   ptrStringBillingData("zone-2qtuhspy7cr6"),
					Host:     ptrStringBillingData("test.example.com"),
					ProxyId:  ptrStringBillingData("sid-2rugn89bkla9"),
					RegionId: ptrStringBillingData("CH"),
				},
			},
			RequestId: ptrStringBillingData("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBillingData()
	res := teo.DataSourceTencentCloudTeoBillingData()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time":  "2025-01-01T00:00:00+08:00",
		"end_time":    "2025-01-02T00:00:00+08:00",
		"zone_ids":    []interface{}{"zone-2qtuhspy7cr6"},
		"metric_name": "acc_flux",
		"interval":    "hour",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "acc_flux#2025-01-01T00:00:00+08:00#2025-01-02T00:00:00+08:00", d.Id())

	dataList := d.Get("data").([]interface{})
	assert.Len(t, dataList, 2)
	dataMap := dataList[0].(map[string]interface{})
	assert.Equal(t, "2025-01-01T00:00:00+08:00", dataMap["time"])
	assert.Equal(t, int(1024), dataMap["value"])
	assert.Equal(t, "zone-2qtuhspy7cr6", dataMap["zone_id"])
	assert.Equal(t, "test.example.com", dataMap["host"])
	assert.Equal(t, "sid-2rugn89bkla9", dataMap["proxy_id"])
	assert.Equal(t, "CH", dataMap["region_id"])

	dataMap2 := dataList[1].(map[string]interface{})
	assert.Equal(t, int(2048), dataMap2["value"])
}

// TestTeoBillingDataDataSource_ReadEmptyData tests read when API returns empty Data list
func TestTeoBillingDataDataSource_ReadEmptyData(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingData().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeBillingDataWithContext", func(_ context.Context, request *teov20220901.DescribeBillingDataRequest) (*teov20220901.DescribeBillingDataResponse, error) {
		resp := teov20220901.NewDescribeBillingDataResponse()
		resp.Response = &teov20220901.DescribeBillingDataResponseParams{
			Data:      nil,
			RequestId: ptrStringBillingData("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBillingData()
	res := teo.DataSourceTencentCloudTeoBillingData()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time":  "2025-01-01T00:00:00+08:00",
		"end_time":    "2025-01-02T00:00:00+08:00",
		"zone_ids":    []interface{}{"zone-2qtuhspy7cr6"},
		"metric_name": "acc_flux",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	dataList := d.Get("data").([]interface{})
	assert.Len(t, dataList, 0)
}

// TestTeoBillingDataDataSource_Schema validates schema definition
func TestTeoBillingDataDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoBillingData()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "start_time")
	assert.Contains(t, res.Schema, "end_time")
	assert.Contains(t, res.Schema, "zone_ids")
	assert.Contains(t, res.Schema, "metric_name")
	assert.Contains(t, res.Schema, "interval")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "group_by")
	assert.Contains(t, res.Schema, "data")
	assert.Contains(t, res.Schema, "result_output_file")

	startTime := res.Schema["start_time"]
	assert.Equal(t, schema.TypeString, startTime.Type)
	assert.True(t, startTime.Required)

	endTime := res.Schema["end_time"]
	assert.Equal(t, schema.TypeString, endTime.Type)
	assert.True(t, endTime.Required)

	zoneIds := res.Schema["zone_ids"]
	assert.Equal(t, schema.TypeList, zoneIds.Type)
	assert.True(t, zoneIds.Required)

	metricName := res.Schema["metric_name"]
	assert.Equal(t, schema.TypeString, metricName.Type)
	assert.True(t, metricName.Required)

	interval := res.Schema["interval"]
	assert.Equal(t, schema.TypeString, interval.Type)
	assert.True(t, interval.Optional)

	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)

	groupBy := res.Schema["group_by"]
	assert.Equal(t, schema.TypeList, groupBy.Type)
	assert.True(t, groupBy.Optional)

	data := res.Schema["data"]
	assert.Equal(t, schema.TypeList, data.Type)
	assert.True(t, data.Computed)

	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)
}

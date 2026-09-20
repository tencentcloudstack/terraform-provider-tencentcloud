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

// go test ./tencentcloud/services/teo/ -run "TestTeoInferenceServiceMonitorDataDataSource" -v -count=1 -gcflags="all=-l"

type mockMetaInferenceServiceMonitorData struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaInferenceServiceMonitorData) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaInferenceServiceMonitorData{}

func newMockMetaInferenceServiceMonitorData() *mockMetaInferenceServiceMonitorData {
	return &mockMetaInferenceServiceMonitorData{client: &connectivity.TencentCloudClient{}}
}

func ptrStringInferenceServiceMonitorData(s string) *string {
	return &s
}

func ptrFloat64InferenceServiceMonitorData(f float64) *float64 {
	return &f
}

func ptrInt64InferenceServiceMonitorData(n int64) *int64 {
	return &n
}

// TestTeoInferenceServiceMonitorDataDataSource_ReadSuccess tests successful read with monitor records
func TestTeoInferenceServiceMonitorDataDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceServiceMonitorData().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceMonitorDataWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceServiceMonitorDataRequest) (*teov20220901.DescribeInferenceServiceMonitorDataResponse, error) {
		assert.Equal(t, "zone-xxxxx", *request.ZoneId)
		assert.Equal(t, "service-xxxxx", *request.ServiceIds[0])
		assert.Equal(t, "cpu_usage_average", *request.MetricNames[0])
		assert.Equal(t, "2024-01-01T00:00:00Z", *request.StartTime)
		assert.Equal(t, "2024-01-01T01:00:00Z", *request.EndTime)
		assert.Equal(t, "hour", *request.Interval)

		resp := teov20220901.NewDescribeInferenceServiceMonitorDataResponse()
		resp.Response = &teov20220901.DescribeInferenceServiceMonitorDataResponseParams{
			TotalCount: ptrInt64InferenceServiceMonitorData(1),
			InferenceServiceMonitorRecords: []*teov20220901.InferenceServiceMonitorRecord{
				{
					ServiceId:  ptrStringInferenceServiceMonitorData("service-xxxxx"),
					MetricName: ptrStringInferenceServiceMonitorData("cpu_usage_average"),
					InferenceServiceMonitorItems: []*teov20220901.InferenceServiceMonitorItem{
						{
							Timestamp: ptrStringInferenceServiceMonitorData("2024-01-01T00:00:00Z"),
							Value:     ptrFloat64InferenceServiceMonitorData(12.5),
						},
						{
							Timestamp: ptrStringInferenceServiceMonitorData("2024-01-01T01:00:00Z"),
							Value:     ptrFloat64InferenceServiceMonitorData(25.0),
						},
					},
				},
			},
			RequestId: ptrStringInferenceServiceMonitorData("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceServiceMonitorData()
	res := teo.DataSourceTencentCloudTeoInferenceServiceMonitorData()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-xxxxx",
		"service_ids":  []interface{}{"service-xxxxx"},
		"metric_names": []interface{}{"cpu_usage_average"},
		"start_time":   "2024-01-01T00:00:00Z",
		"end_time":     "2024-01-01T01:00:00Z",
		"interval":     "hour",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	records := d.Get("inference_service_monitor_records").([]interface{})
	assert.Len(t, records, 1)
	recordMap := records[0].(map[string]interface{})
	assert.Equal(t, "service-xxxxx", recordMap["service_id"])
	assert.Equal(t, "cpu_usage_average", recordMap["metric_name"])

	items := recordMap["inference_service_monitor_items"].([]interface{})
	assert.Len(t, items, 2)
	itemMap := items[0].(map[string]interface{})
	assert.Equal(t, "2024-01-01T00:00:00Z", itemMap["timestamp"])
	assert.Equal(t, 12.5, itemMap["value"])

	itemMap2 := items[1].(map[string]interface{})
	assert.Equal(t, "2024-01-01T01:00:00Z", itemMap2["timestamp"])
	assert.Equal(t, 25.0, itemMap2["value"])
}

// TestTeoInferenceServiceMonitorDataDataSource_ReadEmptyRecords tests read when API returns empty records list
func TestTeoInferenceServiceMonitorDataDataSource_ReadEmptyRecords(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceServiceMonitorData().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceMonitorDataWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceServiceMonitorDataRequest) (*teov20220901.DescribeInferenceServiceMonitorDataResponse, error) {
		resp := teov20220901.NewDescribeInferenceServiceMonitorDataResponse()
		resp.Response = &teov20220901.DescribeInferenceServiceMonitorDataResponseParams{
			TotalCount:                     ptrInt64InferenceServiceMonitorData(0),
			InferenceServiceMonitorRecords: nil,
			RequestId:                      ptrStringInferenceServiceMonitorData("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceServiceMonitorData()
	res := teo.DataSourceTencentCloudTeoInferenceServiceMonitorData()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-xxxxx",
		"service_ids":  []interface{}{"service-xxxxx"},
		"metric_names": []interface{}{"cpu_usage_average"},
		"start_time":   "2024-01-01T00:00:00Z",
		"end_time":     "2024-01-01T01:00:00Z",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestTeoInferenceServiceMonitorDataDataSource_APIError tests read when API returns an error
func TestTeoInferenceServiceMonitorDataDataSource_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceServiceMonitorData().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceMonitorDataWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceServiceMonitorDataRequest) (*teov20220901.DescribeInferenceServiceMonitorDataResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMetaInferenceServiceMonitorData()
	res := teo.DataSourceTencentCloudTeoInferenceServiceMonitorData()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-invalid",
		"service_ids":  []interface{}{"service-xxxxx"},
		"metric_names": []interface{}{"cpu_usage_average"},
		"start_time":   "2024-01-01T00:00:00Z",
		"end_time":     "2024-01-01T01:00:00Z",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoInferenceServiceMonitorDataDataSource_Schema validates schema definition
func TestTeoInferenceServiceMonitorDataDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoInferenceServiceMonitorData()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "service_ids")
	assert.Contains(t, res.Schema, "metric_names")
	assert.Contains(t, res.Schema, "start_time")
	assert.Contains(t, res.Schema, "end_time")
	assert.Contains(t, res.Schema, "interval")
	assert.Contains(t, res.Schema, "inference_service_monitor_records")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	serviceIds := res.Schema["service_ids"]
	assert.Equal(t, schema.TypeList, serviceIds.Type)
	assert.True(t, serviceIds.Required)

	metricNames := res.Schema["metric_names"]
	assert.Equal(t, schema.TypeList, metricNames.Type)
	assert.True(t, metricNames.Required)

	startTime := res.Schema["start_time"]
	assert.Equal(t, schema.TypeString, startTime.Type)
	assert.True(t, startTime.Required)

	endTime := res.Schema["end_time"]
	assert.Equal(t, schema.TypeString, endTime.Type)
	assert.True(t, endTime.Required)

	interval := res.Schema["interval"]
	assert.Equal(t, schema.TypeString, interval.Type)
	assert.True(t, interval.Optional)

	records := res.Schema["inference_service_monitor_records"]
	assert.Equal(t, schema.TypeList, records.Type)
	assert.True(t, records.Computed)

	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)
}

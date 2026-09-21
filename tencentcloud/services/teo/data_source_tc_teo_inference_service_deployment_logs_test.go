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

type mockMetaTeoInferenceServiceDeploymentLogsDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoInferenceServiceDeploymentLogsDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoInferenceServiceDeploymentLogsDS{}

func newMockMetaTeoInferenceServiceDeploymentLogsDS() *mockMetaTeoInferenceServiceDeploymentLogsDS {
	return &mockMetaTeoInferenceServiceDeploymentLogsDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrTeoInfSvcDepLogsDS(s string) *string {
	return &s
}

// TestTeoInferenceServiceDeploymentLogsDataSource_Read_LogListPopulated tests that
// deployment_log_info_set is populated when the API returns a non-empty log list.
func TestTeoInferenceServiceDeploymentLogsDataSource_Read_LogListPopulated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoInferenceServiceDeploymentLogs", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.InferenceServiceDeploymentLogInfo, error) {
		return []*teov20220901.InferenceServiceDeploymentLogInfo{
			{
				LogMessage: ptrStrTeoInfSvcDepLogsDS("pull image success"),
				Timestamp:  ptrStrTeoInfSvcDepLogsDS("2026-09-20 10:00:00"),
			},
			{
				LogMessage: ptrStrTeoInfSvcDepLogsDS("container started"),
				Timestamp:  ptrStrTeoInfSvcDepLogsDS("2026-09-20 10:01:00"),
			},
		}, nil
	})

	meta := newMockMetaTeoInferenceServiceDeploymentLogsDS()
	res := teo.DataSourceTencentCloudTeoInferenceServiceDeploymentLogs()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"service_id": "sid-2quhspyeq8r6",
		"record_id":  "rec-2quhspyeq8r6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	deploymentLogInfoSet := d.Get("deployment_log_info_set").([]interface{})
	assert.Equal(t, 2, len(deploymentLogInfoSet))

	firstLog := deploymentLogInfoSet[0].(map[string]interface{})
	assert.Equal(t, "pull image success", firstLog["log_message"])
	assert.Equal(t, "2026-09-20 10:00:00", firstLog["timestamp"])

	secondLog := deploymentLogInfoSet[1].(map[string]interface{})
	assert.Equal(t, "container started", secondLog["log_message"])
	assert.Equal(t, "2026-09-20 10:01:00", secondLog["timestamp"])

	assert.NotEmpty(t, d.Id())
}

// TestTeoInferenceServiceDeploymentLogsDataSource_Read_EmptyResult tests that the Read
// returns an error (NonRetryableError) when the API returns an empty log list, and the
// data source id is NOT cleared by SetId("").
func TestTeoInferenceServiceDeploymentLogsDataSource_Read_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoInferenceServiceDeploymentLogs", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.InferenceServiceDeploymentLogInfo, error) {
		return []*teov20220901.InferenceServiceDeploymentLogInfo{}, nil
	})

	meta := newMockMetaTeoInferenceServiceDeploymentLogsDS()
	res := teo.DataSourceTencentCloudTeoInferenceServiceDeploymentLogs()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"service_id": "sid-2quhspyeq8r6",
		"record_id":  "rec-2quhspyeq8r6",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "response is empty")
}

// TestTeoInferenceServiceDeploymentLogsDataSource_Read_NilFieldSkipped tests that nil
// LogMessage/Timestamp fields are skipped when mapping to deployment_log_info_set.
func TestTeoInferenceServiceDeploymentLogsDataSource_Read_NilFieldSkipped(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoInferenceServiceDeploymentLogs", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.InferenceServiceDeploymentLogInfo, error) {
		return []*teov20220901.InferenceServiceDeploymentLogInfo{
			{
				LogMessage: ptrStrTeoInfSvcDepLogsDS("log without timestamp"),
				Timestamp:  nil,
			},
		}, nil
	})

	meta := newMockMetaTeoInferenceServiceDeploymentLogsDS()
	res := teo.DataSourceTencentCloudTeoInferenceServiceDeploymentLogs()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"service_id": "sid-2quhspyeq8r6",
		"record_id":  "rec-2quhspyeq8r6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	deploymentLogInfoSet := d.Get("deployment_log_info_set").([]interface{})
	assert.Equal(t, 1, len(deploymentLogInfoSet))

	logInfo := deploymentLogInfoSet[0].(map[string]interface{})
	assert.Equal(t, "log without timestamp", logInfo["log_message"])
	// When Timestamp is nil, the Read loop skips setting the key, so the value defaults to empty string.
	assert.Equal(t, "", logInfo["timestamp"])
}

// TestTeoInferenceServiceDeploymentLogsDataSource_Schema tests the schema definition.
func TestTeoInferenceServiceDeploymentLogsDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoInferenceServiceDeploymentLogs()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	// Required args
	for _, key := range []string{"zone_id", "service_id", "record_id"} {
		schemaItem, ok := res.Schema[key]
		assert.True(t, ok)
		assert.NotNil(t, schemaItem)
		assert.Equal(t, schema.TypeString, schemaItem.Type)
		assert.True(t, schemaItem.Required)
	}

	// Optional args
	for _, key := range []string{"start_time", "end_time", "sort_by", "sort_order", "result_output_file"} {
		schemaItem, ok := res.Schema[key]
		assert.True(t, ok)
		assert.NotNil(t, schemaItem)
		assert.Equal(t, schema.TypeString, schemaItem.Type)
		assert.True(t, schemaItem.Optional)
	}

	// Computed output
	deploymentLogInfoSet := res.Schema["deployment_log_info_set"]
	assert.NotNil(t, deploymentLogInfoSet)
	assert.Equal(t, schema.TypeList, deploymentLogInfoSet.Type)
	assert.True(t, deploymentLogInfoSet.Computed)
	elem := deploymentLogInfoSet.Elem.(*schema.Resource)
	for _, key := range []string{"log_message", "timestamp"} {
		fieldSchema, ok := elem.Schema[key]
		assert.True(t, ok)
		assert.NotNil(t, fieldSchema)
		assert.Equal(t, schema.TypeString, fieldSchema.Type)
		assert.True(t, fieldSchema.Computed)
	}

	// Pagination args must NOT be exposed
	_, ok := res.Schema["limit"]
	assert.False(t, ok)
	_, ok = res.Schema["offset"]
	assert.False(t, ok)
}

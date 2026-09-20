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

// go test ./tencentcloud/services/teo/ -run "TestTeoInferenceServiceDeploymentRecordsDataSource" -v -count=1 -gcflags="all=-l"

type mockMetaTeoInferenceSvcDepRecordsDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoInferenceSvcDepRecordsDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoInferenceSvcDepRecordsDS{}

func newMockMetaTeoInferenceSvcDepRecordsDS() *mockMetaTeoInferenceSvcDepRecordsDS {
	return &mockMetaTeoInferenceSvcDepRecordsDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTeoInferenceSvcDepRecordsDS(s string) *string {
	return &s
}

func ptrInt64TeoInferenceSvcDepRecordsDS(i int64) *int64 {
	return &i
}

func ptrFloat64TeoInferenceSvcDepRecordsDS(f float64) *float64 {
	return &f
}

// TestTeoInferenceServiceDeploymentRecordsDataSource_ReadSuccess tests successful read with a fully populated nested record
func TestTeoInferenceServiceDeploymentRecordsDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoInferenceServiceDeploymentRecordsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.InferenceServiceDeploymentRecord, error) {
		return []*teov20220901.InferenceServiceDeploymentRecord{
			{
				RecordId:     ptrStringTeoInferenceSvcDepRecordsDS("rec-aaaaaaaaaaaa"),
				Operation:    ptrStringTeoInferenceSvcDepRecordsDS("create"),
				Status:       ptrStringTeoInferenceSvcDepRecordsDS("succeeded"),
				Duration:     ptrInt64TeoInferenceSvcDepRecordsDS(120),
				ActiveStatus: ptrStringTeoInferenceSvcDepRecordsDS("active"),
				CreateTime:   ptrStringTeoInferenceSvcDepRecordsDS("2025-01-01T00:00:00+08:00"),
				InferenceServiceConfig: &teov20220901.InferenceServiceConfig{
					ListenPort: ptrInt64TeoInferenceSvcDepRecordsDS(8080),
					RequestPaths: []*string{
						ptrStringTeoInferenceSvcDepRecordsDS("/v1/chat"),
						ptrStringTeoInferenceSvcDepRecordsDS("/v1/embed"),
					},
					Containers: []*teov20220901.InferenceContainerConfig{
						{
							ImageType: ptrStringTeoInferenceSvcDepRecordsDS("TCR"),
							TcrRepositoryConfig: &teov20220901.InferenceTCRRepositoryConfig{
								TCRType:    ptrStringTeoInferenceSvcDepRecordsDS("Enterprise"),
								Image:      ptrStringTeoInferenceSvcDepRecordsDS("test.example.com/img:v1"),
								RegistryId: ptrStringTeoInferenceSvcDepRecordsDS("tcr-xxxxxxxx"),
								RegionName: ptrStringTeoInferenceSvcDepRecordsDS("ap-guangzhou"),
							},
							StartupCommand: ptrStringTeoInferenceSvcDepRecordsDS("/bin/start --port 8080"),
							EnvironmentVariables: []*teov20220901.InferenceEnvironmentVariable{
								{
									Key:   ptrStringTeoInferenceSvcDepRecordsDS("ENV"),
									Value: ptrStringTeoInferenceSvcDepRecordsDS("prod"),
								},
							},
						},
					},
					ResourceConfig: &teov20220901.InferenceResourceConfig{
						ScalingMode:    ptrStringTeoInferenceSvcDepRecordsDS("Auto"),
						HardwareSpecId: ptrStringTeoInferenceSvcDepRecordsDS("hsid-xxxxxxxx"),
						HardwareConfig: &teov20220901.InferenceHardwareConfig{
							GPUNum:   ptrFloat64TeoInferenceSvcDepRecordsDS(1),
							CPUNum:   ptrFloat64TeoInferenceSvcDepRecordsDS(4),
							MemSize:  ptrInt64TeoInferenceSvcDepRecordsDS(8192),
							DiskSize: ptrInt64TeoInferenceSvcDepRecordsDS(10240),
						},
						AutoScalingConfig: &teov20220901.InferenceAutoScalingConfig{
							MinInstanceCount: ptrInt64TeoInferenceSvcDepRecordsDS(1),
							ScalingPolicies: []*teov20220901.InferenceScalingPolicy{
								{
									PolicyName: ptrStringTeoInferenceSvcDepRecordsDS("sp-1"),
									PolicyType: ptrStringTeoInferenceSvcDepRecordsDS("ScheduledScaling"),
									ScheduledScalingPolicy: &teov20220901.InferenceScheduledScalingPolicy{
										ScheduledActions: []*teov20220901.InferenceScheduledScalingAction{
											{
												CronExpression:   ptrStringTeoInferenceSvcDepRecordsDS("0 0 * * *"),
												MinInstanceCount: ptrInt64TeoInferenceSvcDepRecordsDS(2),
											},
										},
									},
								},
							},
						},
						Concurrency: ptrInt64TeoInferenceSvcDepRecordsDS(10),
					},
					AffinityConfig: &teov20220901.InferenceAffinityConfig{
						Switch:       ptrStringTeoInferenceSvcDepRecordsDS("On"),
						AffinityMode: ptrStringTeoInferenceSvcDepRecordsDS("SessionId"),
						SessionIdAffinityConfig: &teov20220901.SessionIdAffinityConfig{
							Source:     ptrStringTeoInferenceSvcDepRecordsDS("Header"),
							HeaderName: ptrStringTeoInferenceSvcDepRecordsDS("EO-Infer-Session-Id"),
						},
					},
				},
			},
		}, nil
	})

	meta := newMockMetaTeoInferenceSvcDepRecordsDS()
	res := teo.DataSourceTencentCloudTeoInferenceServiceDeploymentRecords()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"service_id": "sid-xxxxxxxxxxxx",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "zone-2qtuhspy7cr6"+tccommon.FILED_SP+"sid-xxxxxxxxxxxx", d.Id())

	recordSet := d.Get("record_set").([]interface{})
	assert.Len(t, recordSet, 1)
	record := recordSet[0].(map[string]interface{})
	assert.Equal(t, "rec-aaaaaaaaaaaa", record["record_id"])
	assert.Equal(t, "create", record["operation"])
	assert.Equal(t, "succeeded", record["status"])
	assert.Equal(t, 120, record["duration"])
	assert.Equal(t, "active", record["active_status"])
	assert.Equal(t, "2025-01-01T00:00:00+08:00", record["create_time"])

	configList := record["inference_service_config"].([]interface{})
	assert.Len(t, configList, 1)
	config := configList[0].(map[string]interface{})
	assert.Equal(t, 8080, config["listen_port"])

	requestPaths := config["request_paths"].([]interface{})
	assert.Len(t, requestPaths, 2)
	assert.Equal(t, "/v1/chat", requestPaths[0].(string))

	containers := config["containers"].([]interface{})
	assert.Len(t, containers, 1)
	container := containers[0].(map[string]interface{})
	assert.Equal(t, "TCR", container["image_type"])

	tcrRepoConfigList := container["tcr_repository_config"].([]interface{})
	assert.Len(t, tcrRepoConfigList, 1)
	tcrRepoConfig := tcrRepoConfigList[0].(map[string]interface{})
	assert.Equal(t, "Enterprise", tcrRepoConfig["tcr_type"])
	assert.Equal(t, "tcr-xxxxxxxx", tcrRepoConfig["registry_id"])

	envVars := container["environment_variables"].([]interface{})
	assert.Len(t, envVars, 1)
	envVar := envVars[0].(map[string]interface{})
	assert.Equal(t, "ENV", envVar["key"])
	assert.Equal(t, "prod", envVar["value"])

	resourceConfigList := config["resource_config"].([]interface{})
	assert.Len(t, resourceConfigList, 1)
	resourceConfig := resourceConfigList[0].(map[string]interface{})
	assert.Equal(t, "Auto", resourceConfig["scaling_mode"])
	assert.Equal(t, "hsid-xxxxxxxx", resourceConfig["hardware_spec_id"])

	hardwareConfigList := resourceConfig["hardware_config"].([]interface{})
	assert.Len(t, hardwareConfigList, 1)
	hardwareConfig := hardwareConfigList[0].(map[string]interface{})
	assert.Equal(t, 1, hardwareConfig["gpu_num"])
	assert.Equal(t, 4, hardwareConfig["cpu_num"])
	assert.Equal(t, 8192, hardwareConfig["mem_size"])

	autoScalingConfigList := resourceConfig["auto_scaling_config"].([]interface{})
	assert.Len(t, autoScalingConfigList, 1)
	autoScalingConfig := autoScalingConfigList[0].(map[string]interface{})
	assert.Equal(t, 1, autoScalingConfig["min_instance_count"])

	scalingPolicies := autoScalingConfig["scaling_policies"].([]interface{})
	assert.Len(t, scalingPolicies, 1)
	policy := scalingPolicies[0].(map[string]interface{})
	assert.Equal(t, "sp-1", policy["policy_name"])

	scheduledScalingPolicyList := policy["scheduled_scaling_policy"].([]interface{})
	assert.Len(t, scheduledScalingPolicyList, 1)
	scheduledScalingPolicy := scheduledScalingPolicyList[0].(map[string]interface{})

	scheduledActions := scheduledScalingPolicy["scheduled_actions"].([]interface{})
	assert.Len(t, scheduledActions, 1)
	action := scheduledActions[0].(map[string]interface{})
	assert.Equal(t, "0 0 * * *", action["cron_expression"])
	assert.Equal(t, 2, action["min_instance_count"])

	assert.Equal(t, 10, resourceConfig["concurrency"])

	affinityConfigList := config["affinity_config"].([]interface{})
	assert.Len(t, affinityConfigList, 1)
	affinityConfig := affinityConfigList[0].(map[string]interface{})
	assert.Equal(t, "On", affinityConfig["switch"])

	sessionIdAffinityConfigList := affinityConfig["session_id_affinity_config"].([]interface{})
	assert.Len(t, sessionIdAffinityConfigList, 1)
	sessionIdAffinityConfig := sessionIdAffinityConfigList[0].(map[string]interface{})
	assert.Equal(t, "Header", sessionIdAffinityConfig["source"])
}

// TestTeoInferenceServiceDeploymentRecordsDataSource_ReadEmpty tests that an empty record set returns an error and does not clear the existing id
func TestTeoInferenceServiceDeploymentRecordsDataSource_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoInferenceServiceDeploymentRecordsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.InferenceServiceDeploymentRecord, error) {
		return nil, nil
	})

	meta := newMockMetaTeoInferenceSvcDepRecordsDS()
	res := teo.DataSourceTencentCloudTeoInferenceServiceDeploymentRecords()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"service_id": "sid-xxxxxxxxxxxx",
	})
	d.SetId("existing-state-id")

	err := res.Read(d, meta)
	assert.Error(t, err)
	// the NonRetryableError preserves existing state id; it should NOT be cleared
	assert.Equal(t, "existing-state-id", d.Id())
}

// TestTeoInferenceServiceDeploymentRecordsDataSource_NilInferenceServiceConfig tests nil safety when InferenceServiceConfig is nil
func TestTeoInferenceServiceDeploymentRecordsDataSource_NilInferenceServiceConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoInferenceServiceDeploymentRecordsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.InferenceServiceDeploymentRecord, error) {
		return []*teov20220901.InferenceServiceDeploymentRecord{
			{
				RecordId:               ptrStringTeoInferenceSvcDepRecordsDS("rec-bbbbbbbbbbbb"),
				Operation:              ptrStringTeoInferenceSvcDepRecordsDS("update"),
				Status:                 ptrStringTeoInferenceSvcDepRecordsDS("failed"),
				Duration:               ptrInt64TeoInferenceSvcDepRecordsDS(60),
				CreateTime:             ptrStringTeoInferenceSvcDepRecordsDS("2025-01-02T00:00:00+08:00"),
				ActiveStatus:           ptrStringTeoInferenceSvcDepRecordsDS("inactive"),
				InferenceServiceConfig: nil,
			},
		}, nil
	})

	meta := newMockMetaTeoInferenceSvcDepRecordsDS()
	res := teo.DataSourceTencentCloudTeoInferenceServiceDeploymentRecords()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"service_id": "sid-yyyyyyyyyyyy",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	recordSet := d.Get("record_set").([]interface{})
	assert.Len(t, recordSet, 1)
	record := recordSet[0].(map[string]interface{})
	assert.Equal(t, "rec-bbbbbbbbbbbb", record["record_id"])
	assert.Equal(t, "update", record["operation"])
	assert.Equal(t, 60, record["duration"])
	// inference_service_config should be empty (not set) since it was nil
	assert.Equal(t, []interface{}{}, record["inference_service_config"])
}

// TestTeoInferenceServiceDeploymentRecordsDataSource_Schema validates schema definition
func TestTeoInferenceServiceDeploymentRecordsDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoInferenceServiceDeploymentRecords()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "service_id")
	assert.Contains(t, res.Schema, "sort_by")
	assert.Contains(t, res.Schema, "sort_order")
	assert.Contains(t, res.Schema, "record_set")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	serviceId := res.Schema["service_id"]
	assert.Equal(t, schema.TypeString, serviceId.Type)
	assert.True(t, serviceId.Required)

	sortBy := res.Schema["sort_by"]
	assert.Equal(t, schema.TypeString, sortBy.Type)
	assert.True(t, sortBy.Optional)

	sortOrder := res.Schema["sort_order"]
	assert.Equal(t, schema.TypeString, sortOrder.Type)
	assert.True(t, sortOrder.Optional)

	recordSet := res.Schema["record_set"]
	assert.Equal(t, schema.TypeList, recordSet.Type)
	assert.True(t, recordSet.Computed)

	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)

	// verify nested record_set schema has the flattened fields
	elem := recordSet.Elem.(*schema.Resource)
	assert.Contains(t, elem.Schema, "record_id")
	assert.Contains(t, elem.Schema, "operation")
	assert.Contains(t, elem.Schema, "status")
	assert.Contains(t, elem.Schema, "duration")
	assert.Contains(t, elem.Schema, "create_time")
	assert.Contains(t, elem.Schema, "active_status")
	assert.Contains(t, elem.Schema, "inference_service_config")
}

package dts_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	dtssvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dts"
)

type mockMetaDtsSyncJobConfig struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDtsSyncJobConfig) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDtsSyncJobConfig{}

func newMockMetaDtsSyncJobConfig() *mockMetaDtsSyncJobConfig {
	return &mockMetaDtsSyncJobConfig{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDtsSyncJobConfig(s string) *string { return &s }

// TestDtsSyncJobConfig_Create_Success tests Create sets id and calls Update path
func TestDtsSyncJobConfig_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dtsClient := &dts.Client{}
	patches.ApplyMethodReturn(newMockMetaDtsSyncJobConfig().client, "UseDtsClient", dtsClient)

	patches.ApplyMethodFunc(dtsClient, "ModifySyncJobConfig", func(request *dts.ModifySyncJobConfigRequest) (*dts.ModifySyncJobConfigResponse, error) {
		assert.NotNil(t, request.JobId)
		assert.Equal(t, "sync-test-job-config", *request.JobId)
		assert.NotNil(t, request.DynamicOptions)
		resp := dts.NewModifySyncJobConfigResponse()
		resp.Response = &dts.ModifySyncJobConfigResponseParams{
			RequestId: ptrStringDtsSyncJobConfig("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dtsClient, "DescribeSyncJobs", func(request *dts.DescribeSyncJobsRequest) (*dts.DescribeSyncJobsResponse, error) {
		assert.NotNil(t, request.JobId)
		assert.Equal(t, "sync-test-job-config", *request.JobId)
		resp := dts.NewDescribeSyncJobsResponse()
		resp.Response = &dts.DescribeSyncJobsResponseParams{
			JobList: []*dts.SyncJobInfo{
				{
					JobId:   ptrStringDtsSyncJobConfig("sync-test-job-config"),
					JobName: ptrStringDtsSyncJobConfig("tf_test_sync_job_config"),
					Status:  ptrStringDtsSyncJobConfig("Running"),
					Objects: &dts.Objects{
						Mode: ptrStringDtsSyncJobConfig("Partial"),
						Databases: []*dts.Database{
							{
								DbName: ptrStringDtsSyncJobConfig("tf_ci_test"),
								DbMode: ptrStringDtsSyncJobConfig("Partial"),
							},
						},
					},
					Options: &dts.Options{
						OpTypes:            []*string{ptrStringDtsSyncJobConfig("Insert"), ptrStringDtsSyncJobConfig("Update"), ptrStringDtsSyncJobConfig("Delete"), ptrStringDtsSyncJobConfig("DDL")},
						ConflictHandleType: ptrStringDtsSyncJobConfig("ReportError"),
					},
				},
			},
			RequestId: ptrStringDtsSyncJobConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDtsSyncJobConfig()
	res := dtssvc.ResourceTencentCloudDtsSyncJobConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"job_id": "sync-test-job-config",
		"dynamic_options": []interface{}{
			map[string]interface{}{
				"op_types":             []interface{}{"Insert", "Update", "Delete", "DDL"},
				"conflict_handle_type": "ReportError",
				"ddl_options": []interface{}{
					map[string]interface{}{
						"ddl_object": "Table",
						"ddl_value":  []interface{}{"Create", "Drop"},
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sync-test-job-config", d.Id())
	assert.Equal(t, "Running", d.Get("status").(string))
}

// TestDtsSyncJobConfig_Read_Success tests Read populates state from API
func TestDtsSyncJobConfig_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dtsClient := &dts.Client{}
	patches.ApplyMethodReturn(newMockMetaDtsSyncJobConfig().client, "UseDtsClient", dtsClient)

	patches.ApplyMethodFunc(dtsClient, "DescribeSyncJobs", func(request *dts.DescribeSyncJobsRequest) (*dts.DescribeSyncJobsResponse, error) {
		resp := dts.NewDescribeSyncJobsResponse()
		resp.Response = &dts.DescribeSyncJobsResponseParams{
			JobList: []*dts.SyncJobInfo{
				{
					JobId:           ptrStringDtsSyncJobConfig("sync-read-test"),
					JobName:         ptrStringDtsSyncJobConfig("tf_test_read"),
					PayMode:         ptrStringDtsSyncJobConfig("PostPay"),
					RunMode:         ptrStringDtsSyncJobConfig("Immediate"),
					Status:          ptrStringDtsSyncJobConfig("Running"),
					TradeStatus:     ptrStringDtsSyncJobConfig("Normal"),
					SrcRegion:       ptrStringDtsSyncJobConfig("ap-guangzhou"),
					SrcDatabaseType: ptrStringDtsSyncJobConfig("mysql"),
					SrcAccessType:   ptrStringDtsSyncJobConfig("cdb"),
					Specification:   ptrStringDtsSyncJobConfig("micro"),
					Objects: &dts.Objects{
						Mode: ptrStringDtsSyncJobConfig("Partial"),
						Databases: []*dts.Database{
							{
								DbName: ptrStringDtsSyncJobConfig("tf_ci_test"),
								DbMode: ptrStringDtsSyncJobConfig("Partial"),
							},
						},
					},
					Options: &dts.Options{
						OpTypes:            []*string{ptrStringDtsSyncJobConfig("Insert"), ptrStringDtsSyncJobConfig("Update")},
						ConflictHandleType: ptrStringDtsSyncJobConfig("ReportError"),
					},
				},
			},
			RequestId: ptrStringDtsSyncJobConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDtsSyncJobConfig()
	res := dtssvc.ResourceTencentCloudDtsSyncJobConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("sync-read-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sync-read-test", d.Id())
	assert.Equal(t, "tf_test_read", d.Get("job_name").(string))
	assert.Equal(t, "Running", d.Get("status").(string))
	assert.Equal(t, "PostPay", d.Get("pay_mode").(string))
	assert.Equal(t, "Partial", d.Get("dynamic_objects.0.mode").(string))
	assert.Equal(t, "ReportError", d.Get("dynamic_options.0.conflict_handle_type").(string))
}

// TestDtsSyncJobConfig_Read_NotFound tests Read sets id to empty when not found
func TestDtsSyncJobConfig_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dtsClient := &dts.Client{}
	patches.ApplyMethodReturn(newMockMetaDtsSyncJobConfig().client, "UseDtsClient", dtsClient)

	patches.ApplyMethodFunc(dtsClient, "DescribeSyncJobs", func(request *dts.DescribeSyncJobsRequest) (*dts.DescribeSyncJobsResponse, error) {
		resp := dts.NewDescribeSyncJobsResponse()
		resp.Response = &dts.DescribeSyncJobsResponseParams{
			JobList:   []*dts.SyncJobInfo{},
			RequestId: ptrStringDtsSyncJobConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDtsSyncJobConfig()
	res := dtssvc.ResourceTencentCloudDtsSyncJobConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("sync-not-found")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestDtsSyncJobConfig_Update_Success tests Update builds ModifySyncJobConfig request
func TestDtsSyncJobConfig_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dtsClient := &dts.Client{}
	patches.ApplyMethodReturn(newMockMetaDtsSyncJobConfig().client, "UseDtsClient", dtsClient)

	patches.ApplyMethodFunc(dtsClient, "ModifySyncJobConfig", func(request *dts.ModifySyncJobConfigRequest) (*dts.ModifySyncJobConfigResponse, error) {
		assert.NotNil(t, request.JobId)
		assert.Equal(t, "sync-update-test", *request.JobId)
		assert.NotNil(t, request.DynamicObjects)
		assert.NotNil(t, request.DynamicOptions)
		assert.NotNil(t, request.DynamicOptions.OpTypes)
		assert.NotNil(t, request.DynamicOptions.DdlOptions)
		resp := dts.NewModifySyncJobConfigResponse()
		resp.Response = &dts.ModifySyncJobConfigResponseParams{
			RequestId: ptrStringDtsSyncJobConfig("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dtsClient, "DescribeSyncJobs", func(request *dts.DescribeSyncJobsRequest) (*dts.DescribeSyncJobsResponse, error) {
		resp := dts.NewDescribeSyncJobsResponse()
		resp.Response = &dts.DescribeSyncJobsResponseParams{
			JobList: []*dts.SyncJobInfo{
				{
					JobId:   ptrStringDtsSyncJobConfig("sync-update-test"),
					JobName: ptrStringDtsSyncJobConfig("tf_test_update"),
					Status:  ptrStringDtsSyncJobConfig("Running"),
					Objects: &dts.Objects{
						Mode: ptrStringDtsSyncJobConfig("Partial"),
						Databases: []*dts.Database{
							{
								DbName: ptrStringDtsSyncJobConfig("tf_ci_test"),
								DbMode: ptrStringDtsSyncJobConfig("Partial"),
							},
						},
					},
					Options: &dts.Options{
						OpTypes:            []*string{ptrStringDtsSyncJobConfig("Insert"), ptrStringDtsSyncJobConfig("Delete")},
						ConflictHandleType: ptrStringDtsSyncJobConfig("Cover"),
					},
				},
			},
			RequestId: ptrStringDtsSyncJobConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDtsSyncJobConfig()
	res := dtssvc.ResourceTencentCloudDtsSyncJobConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"job_id": "sync-update-test",
		"dynamic_objects": []interface{}{
			map[string]interface{}{
				"mode": "Partial",
				"databases": []interface{}{
					map[string]interface{}{
						"db_name": "tf_ci_test",
						"db_mode": "Partial",
					},
				},
			},
		},
		"dynamic_options": []interface{}{
			map[string]interface{}{
				"op_types":             []interface{}{"Insert", "Delete"},
				"conflict_handle_type": "Cover",
				"ddl_options": []interface{}{
					map[string]interface{}{
						"ddl_object": "Table",
						"ddl_value":  []interface{}{"Create"},
					},
				},
			},
		},
	})
	d.SetId("sync-update-test")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sync-update-test", d.Id())
	assert.Equal(t, "Running", d.Get("status").(string))
}

// TestDtsSyncJobConfig_Update_APIError tests Update handles API error
func TestDtsSyncJobConfig_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dtsClient := &dts.Client{}
	patches.ApplyMethodReturn(newMockMetaDtsSyncJobConfig().client, "UseDtsClient", dtsClient)

	patches.ApplyMethodFunc(dtsClient, "ModifySyncJobConfig", func(request *dts.ModifySyncJobConfigRequest) (*dts.ModifySyncJobConfigResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaDtsSyncJobConfig()
	res := dtssvc.ResourceTencentCloudDtsSyncJobConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"job_id": "sync-error-test",
		"dynamic_options": []interface{}{
			map[string]interface{}{
				"op_types": []interface{}{"Insert"},
				"ddl_options": []interface{}{
					map[string]interface{}{
						"ddl_object": "Table",
						"ddl_value":  []interface{}{"Create"},
					},
				},
			},
		},
	})
	d.SetId("sync-error-test")

	err := res.Update(d, meta)
	assert.Error(t, err)
}

// TestDtsSyncJobConfig_Delete tests Delete is a no-op for config resources
func TestDtsSyncJobConfig_Delete(t *testing.T) {
	meta := newMockMetaDtsSyncJobConfig()
	res := dtssvc.ResourceTencentCloudDtsSyncJobConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"job_id": "sync-delete-test",
	})
	d.SetId("sync-delete-test")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sync-delete-test", d.Id())
}

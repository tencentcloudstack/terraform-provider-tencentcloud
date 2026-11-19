package dts_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	dtssvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dts"
)

func TestAccTencentCloudDtsSyncJobsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDtsSyncJobs,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dts_sync_jobs.sync_jobs"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.0.job_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.0.job_name", "tf_dts_test"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.0.src_database_type", "mysql"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.0.dst_database_type", "cynosdbmysql"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.0.src_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.0.dst_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.0.tags.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.0.tags.0.tag_key", "aaa"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_sync_jobs.sync_jobs", "list.0.tags.0.tag_value", "bbb"),
				),
			},
		},
	})
}

const testAccDataSourceDtsSyncJobs = `

resource "tencentcloud_dts_sync_job" "job" {
	job_name = "tf_dts_test"
	pay_mode = "PostPay"
	src_database_type = "mysql"
	src_region = "ap-guangzhou"
	dst_database_type = "cynosdbmysql"
	dst_region = "ap-guangzhou"
	tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
	auto_renew = 0
   instance_class = "micro"
  }

data "tencentcloud_dts_sync_jobs" "sync_jobs" {
  job_id = tencentcloud_dts_sync_job.job.id
}

`

// ---- gomonkey-based unit tests for the schema_mode field ----
//
// These tests mock the DTS cloud API client method DescribeSyncJobs (called by
// the data source Read handler via DtsService.DescribeDtsSyncJobsByFilter) so
// they can run without real cloud credentials. They verify that the
// `schema_mode` field is correctly read back from the DescribeSyncJobs response
// into Terraform state (read path) for the tencentcloud_dts_sync_jobs data source.

type mockMetaDtsSyncJobs struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDtsSyncJobs) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDtsSyncJobs{}

func newMockMetaDtsSyncJobs() *mockMetaDtsSyncJobs {
	return &mockMetaDtsSyncJobs{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDtsSyncJobs(s string) *string { return &s }

func buildDescribeSyncJobsResponseWithSchemaMode(schemaMode *string) *dts.DescribeSyncJobsResponse {
	resp := dts.NewDescribeSyncJobsResponse()
	resp.Response = &dts.DescribeSyncJobsResponseParams{
		JobList: []*dts.SyncJobInfo{
			{
				JobId:         ptrStringDtsSyncJobs("sync-jobs-schema-mode"),
				JobName:       ptrStringDtsSyncJobs("tf_test_sync_jobs"),
				SrcAccessType: ptrStringDtsSyncJobs("cdb"),
				DstAccessType: ptrStringDtsSyncJobs("cdb"),
				RunMode:       ptrStringDtsSyncJobs("Immediate"),
				Objects: &dts.Objects{
					Mode: ptrStringDtsSyncJobs("Partial"),
					Databases: []*dts.Database{
						{
							DbName:     ptrStringDtsSyncJobs("tf_ci_test"),
							NewDbName:  ptrStringDtsSyncJobs("tf_ci_test_new"),
							DbMode:     ptrStringDtsSyncJobs("Partial"),
							SchemaMode: schemaMode,
							TableMode:  ptrStringDtsSyncJobs("All"),
						},
					},
				},
			},
		},
		RequestId: ptrStringDtsSyncJobs("fake-request-id"),
	}
	return resp
}

// TestDtsSyncJobsDataSource_Read_SchemaMode verifies that a non-nil
// Database.SchemaMode returned by DescribeSyncJobs is read back into Terraform
// state as list[].objects[].databases[].schema_mode for the data source.
func TestDtsSyncJobsDataSource_Read_SchemaMode(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dtsClient := &dts.Client{}
	patches.ApplyMethodReturn(newMockMetaDtsSyncJobs().client, "UseDtsClient", dtsClient)

	patches.ApplyMethodFunc(dtsClient, "DescribeSyncJobs", func(request *dts.DescribeSyncJobsRequest) (*dts.DescribeSyncJobsResponse, error) {
		return buildDescribeSyncJobsResponseWithSchemaMode(ptrStringDtsSyncJobs("All")), nil
	})

	meta := newMockMetaDtsSyncJobs()
	res := dtssvc.DataSourceTencentCloudDtsSyncJobs()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sync-jobs-schema-mode", d.Get("list.0.job_id").(string))
	assert.Equal(t, "All", d.Get("list.0.objects.0.databases.0.schema_mode").(string))
	assert.Equal(t, "Partial", d.Get("list.0.objects.0.databases.0.db_mode").(string))
}

// TestDtsSyncJobsDataSource_Read_SchemaModeNil verifies that when
// Database.SchemaMode is nil, the data source Read handler skips setting the
// schema_mode field (leaving it empty in state) rather than crashing.
func TestDtsSyncJobsDataSource_Read_SchemaModeNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dtsClient := &dts.Client{}
	patches.ApplyMethodReturn(newMockMetaDtsSyncJobs().client, "UseDtsClient", dtsClient)

	patches.ApplyMethodFunc(dtsClient, "DescribeSyncJobs", func(request *dts.DescribeSyncJobsRequest) (*dts.DescribeSyncJobsResponse, error) {
		return buildDescribeSyncJobsResponseWithSchemaMode(nil), nil
	})

	meta := newMockMetaDtsSyncJobs()
	res := dtssvc.DataSourceTencentCloudDtsSyncJobs()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sync-jobs-schema-mode", d.Get("list.0.job_id").(string))
	// schema_mode is left empty because the cloud API returned nil.
	assert.Equal(t, "", d.Get("list.0.objects.0.databases.0.schema_mode").(string))
	assert.Equal(t, "Partial", d.Get("list.0.objects.0.databases.0.db_mode").(string))
}

package dts_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDtsMigrateJobsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDtsMigrateJobs_all,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dts_migrate_jobs.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dts_migrate_jobs.all", "list.#"),
				),
			},
			{
				Config: testAccDataSourceDtsMigrateJobs_job,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dts_migrate_jobs.job"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_migrate_jobs.job", "list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_migrate_jobs.job", "list.0.job_name", "tf_test_migration_job"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_migrate_jobs.job", "list.0.status", "created"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_migrate_jobs.job", "list.0.trade_info.0.instance_class", "small"),
				),
			},
			{
				Config: testAccDataSourceDtsMigrateJobs_src_dest,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dts_migrate_jobs.src_dest"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dts_migrate_jobs.src_dest", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_migrate_jobs.src_dest", "list.0.src_info.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_migrate_jobs.src_dest", "list.0.src_info.0.database_type", "mysql"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_migrate_jobs.src_dest", "list.0.dst_info.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("data.tencentcloud_dts_migrate_jobs.src_dest", "list.0.dst_info.0.database_type", "cynosdbmysql"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dts_migrate_jobs.src_dest", "list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dts_migrate_jobs.src_dest", "list.0.tags.#"),
				),
			},
		},
	})
}

const testAccDataSourceDtsMigrateJob = `

resource "tencentcloud_dts_migrate_job" "migrate_job" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region = "ap-guangzhou"
  dst_region = "ap-guangzhou"
  instance_class = "small"
  job_name = "tf_test_migration_job"
  tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
}

`

const testAccDataSourceDtsMigrateJobs_all = testAccDataSourceDtsMigrateJob + `

data "tencentcloud_dts_migrate_jobs" "all" {}

`

const testAccDataSourceDtsMigrateJobs_job = testAccDataSourceDtsMigrateJob + `

data "tencentcloud_dts_migrate_jobs" "job" {
  job_id = tencentcloud_dts_migrate_job.migrate_job.id
  job_name = tencentcloud_dts_migrate_job.migrate_job.job_name
  status = ["created"]
}

`

const testAccDataSourceDtsMigrateJobs_src_dest = testAccDataSourceDtsMigrateJob + `

data "tencentcloud_dts_migrate_jobs" "src_dest" {
  
  src_region = "ap-guangzhou"
  src_database_type = ["mysql"]
  dst_region = "ap-guangzhou"
  dst_database_type = ["cynosdbmysql"]

  status = ["created"]
  tag_filters {
	tag_key = "aaa"
	tag_value = "bbb"
  }
}

`

package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDtsSyncJobsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDtsSyncJobs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dts_sync_jobs.sync_jobs"),
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

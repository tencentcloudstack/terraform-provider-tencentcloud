
package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDtsMigrateJobsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDtsMigrateJobs,
				Check: resource.ComposeTestCheckFunc(
				  testAccCheckTencentCloudDataSourceID("data.tencentcloud_dts_migrate_jobs.migrate_jobs"),
				),
			},
		},
	})
}

const testAccDataSourceDtsMigrateJobs = `

data "tencentcloud_dts_migrate_jobs" "migrate_jobs" {
  job_id = ""
  job_name = ""
  status = ""
  src_instance_id = ""
  src_region = ""
  src_database_type = ""
  src_access_type = ""
  dst_instance_id = ""
  dst_region = ""
  dst_database_type = ""
  dst_access_type = ""
  run_mode = ""
  order_seq = ""
  tag_filters {
			tag_key = ""
			tag_value = ""

  }
}

`

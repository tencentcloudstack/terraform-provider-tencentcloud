
package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				),
			},
		},
	})
}

const testAccDataSourceDtsSyncJobs = `

data "tencentcloud_dts_sync_jobs" "sync_jobs" {
  job_id = ""
  job_name = ""
  order = ""
  order_seq = ""
  status = ""
  run_mode = ""
  job_type = ""
  pay_mode = ""
  tag_filters {
			tag_key = ""
			tag_value = ""

  }
  }

`

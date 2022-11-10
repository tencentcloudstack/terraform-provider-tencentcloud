package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudPtsCronJob_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsCronJob,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_pts_cron_job.cron_job", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_pts_cron_job.cronJob",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsCronJob = `

resource "tencentcloud_pts_cron_job" "cron_job" {
  name = ""
  project_id = ""
  scenario_id = ""
  scenario_name = ""
  frequency_type = ""
  cron_expression = ""
  job_owner = ""
  end_time = ""
  notice_id = ""
  note = ""
              }

`

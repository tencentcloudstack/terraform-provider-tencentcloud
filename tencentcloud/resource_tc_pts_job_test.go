package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudPtsJob_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsJob,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_pts_job.job", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_pts_job.job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsJob = `

resource "tencentcloud_pts_job" "job" {
  scenario_id = ""
  job_owner = ""
  project_id = ""
  debug = ""
  note = ""
}

`

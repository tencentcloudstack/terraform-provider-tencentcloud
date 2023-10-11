package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixPtsCronJobRestartResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsCronJobRestart,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_cron_job_restart.cron_job_restart", "id")),
			},
		},
	})
}

const testAccPtsCronJobRestart = `

resource "tencentcloud_pts_cron_job_restart" "cron_job_restart" {
  project_id  = "project-abc"
  cron_job_id = "job-dtm93vx0"
}

`

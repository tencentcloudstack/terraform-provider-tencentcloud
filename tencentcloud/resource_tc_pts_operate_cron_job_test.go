package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsOperateCronJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsOperateCronJob,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_operate_cron_job.operate_cron_job", "id")),
			},
			{
				ResourceName:      "tencentcloud_pts_operate_cron_job.operate_cron_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsOperateCronJob = `

resource "tencentcloud_pts_operate_cron_job" "operate_cron_job" {
  project_id = "project-abc"
  cron_job_ids = 
}

`

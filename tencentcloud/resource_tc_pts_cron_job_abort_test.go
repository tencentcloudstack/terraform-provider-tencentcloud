package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsCronJobAbortResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsCronJobAbort,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_cron_job_abort.cron_job_abort", "id")),
			},
			{
				ResourceName:      "tencentcloud_pts_cron_job_abort.cron_job_abort",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsCronJobAbort = `

resource "tencentcloud_pts_cron_job_abort" "cron_job_abort" {
  project_id = ""
  cron_job_ids = 
}

`

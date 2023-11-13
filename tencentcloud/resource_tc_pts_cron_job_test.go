package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsCronJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsCronJob,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_cron_job.cron_job", "id")),
			},
			{
				ResourceName:      "tencentcloud_pts_cron_job.cron_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsCronJob = `

resource "tencentcloud_pts_cron_job" "cron_job" {
  name = &lt;nil&gt;
  project_id = &lt;nil&gt;
  scenario_id = &lt;nil&gt;
  scenario_name = &lt;nil&gt;
  frequency_type = &lt;nil&gt;
  cron_expression = &lt;nil&gt;
  job_owner = &lt;nil&gt;
  end_time = &lt;nil&gt;
  notice_id = &lt;nil&gt;
  note = &lt;nil&gt;
              }

`

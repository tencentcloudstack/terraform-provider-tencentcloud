package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsJob,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_job.job", "id")),
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
  scenario_id = &lt;nil&gt;
  job_owner = &lt;nil&gt;
  project_id = &lt;nil&gt;
  debug = &lt;nil&gt;
  note = &lt;nil&gt;
                                                        }

`

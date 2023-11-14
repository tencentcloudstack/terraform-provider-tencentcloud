package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsOperateJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsOperateJob,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_operate_job.operate_job", "id")),
			},
			{
				ResourceName:      "tencentcloud_pts_operate_job.operate_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsOperateJob = `

resource "tencentcloud_pts_operate_job" "operate_job" {
  job_id = ""
  project_id = ""
  scenario_id = ""
  abort_reason = 
}

`

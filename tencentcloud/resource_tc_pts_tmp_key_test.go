package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsTmpKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsTmpKey,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_tmp_key.tmp_key", "id")),
			},
			{
				ResourceName:      "tencentcloud_pts_tmp_key.tmp_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsTmpKey = `

resource "tencentcloud_pts_tmp_key" "tmp_key" {
  project_id = "project-abc"
  scenario_id = "scenario-abc"
}

`

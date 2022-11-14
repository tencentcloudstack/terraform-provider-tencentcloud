package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudRumProject_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumProject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_project.project", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_rum_project.project",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRumProject = `

resource "tencentcloud_rum_project" "project" {
  name = ""
  instance_i_d = ""
  rate = ""
  enable_u_r_l_group = ""
  type = ""
  repo = ""
  u_r_l = ""
  desc = ""
}

`

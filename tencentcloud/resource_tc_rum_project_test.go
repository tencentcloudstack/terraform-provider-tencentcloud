package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumProjectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumProject,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_rum_project.project", "id")),
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
  name = &lt;nil&gt;
  instance_i_d = &lt;nil&gt;
  rate = &lt;nil&gt;
  enable_u_r_l_group = &lt;nil&gt;
  type = &lt;nil&gt;
  repo = &lt;nil&gt;
  u_r_l = &lt;nil&gt;
  desc = &lt;nil&gt;
              }

`

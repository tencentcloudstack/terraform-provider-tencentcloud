package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudBiProjectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiProject,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bi_project.project", "id")),
			},
			{
				ResourceName:      "tencentcloud_bi_project.project",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBiProject = `

resource "tencentcloud_bi_project" "project" {
  name = "abc"
  color_code = "#066EFF"
  logo = &lt;nil&gt;
  mark = "abc"
  is_apply = true
  default_panel_type = 123
}

`

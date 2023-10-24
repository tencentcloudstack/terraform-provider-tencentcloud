package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudBiProjectResource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bi_project.project", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project.project", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project.project", "color_code", "#7BD936"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project.project", "logo", "TF-test"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project.project", "mark", "project mark"),
				),
			},
			{
				ResourceName:      "tencentcloud_bi_project.project",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBiProjectUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bi_project.project", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project.project", "name", "terraform_test1"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project.project", "color_code", "#066EFF"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project.project", "logo", "TF-test1"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project.project", "mark", "project mark1"),
				),
			},
		},
	})
}

const testAccBiProject = `

resource "tencentcloud_bi_project" "project" {
  name               = "terraform_test"
  color_code         = "#7BD936"
  logo               = "TF-test"
  mark               = "project mark"
}

`

const testAccBiProjectUp = `

resource "tencentcloud_bi_project" "project" {
  name               = "terraform_test1"
  color_code         = "#066EFF"
  logo               = "TF-test1"
  mark               = "project mark1"
}

`

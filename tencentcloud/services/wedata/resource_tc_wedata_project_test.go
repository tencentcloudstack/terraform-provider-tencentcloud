package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataProjectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataProject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_project.example", "id"),
				),
			},
			{
				Config: testAccWedataProjectUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_project.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_project.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataProject = `
resource "tencentcloud_wedata_project" "example" {
  project {
    project_name  = "tf_example"
    display_name  = "display_name"
    project_model = "SIMPLE"
  }

  status = 0
}
`

const testAccWedataProjectUpdate = `
resource "tencentcloud_wedata_project" "example" {
  project {
    project_name  = "tf_example"
    display_name  = "display_name_update"
    project_model = "SIMPLE"
  }

  status = 1
}
`

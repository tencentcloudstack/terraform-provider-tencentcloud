package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumProjectStatusConfigResource_basic -v
func TestAccTencentCloudRumProjectStatusConfigResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumProjectStatusConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_project_status_config.project_status_config", "id"),
				),
			},
			{
				Config: testAccRumProjectStatusConfigUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_project_status_config.project_status_config", "id"),
				),
			},
			{
				ResourceName:            "tencentcloud_rum_project_status_config.project_status_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"operate"},
			},
		},
	})
}

const testAccRumProjectStatusConfigVar = `
variable "project_id" {
  default = "` + defaultRumProjectId + `"
}
`

const testAccRumProjectStatusConfig = testAccRumProjectStatusConfigVar + `

resource "tencentcloud_rum_project_status_config" "project_status_config" {
	project_id = var.project_id
	operate = "stop"
}

`

const testAccRumProjectStatusConfigUp = testAccRumProjectStatusConfigVar + `


resource "tencentcloud_rum_project_status_config" "project_status_config" {
	project_id = var.project_id
	operate = "resume"
}

`

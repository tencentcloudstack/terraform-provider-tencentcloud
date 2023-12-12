package bh_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbCmdTemplateResource_basic -v
func TestAccTencentCloudNeedFixDasbCmdTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbCmdTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_cmd_template.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_cmd_template.example", "name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_cmd_template.example", "cmd_list", "rm -rf*"),
				),
			},
			{
				ResourceName:      "tencentcloud_dasb_cmd_template.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDasbCmdTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_cmd_template.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_cmd_template.example", "name", "tf_example_update"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_cmd_template.example", "cmd_list", "reboot"),
				),
			},
		},
	})
}

const testAccDasbCmdTemplate = `
resource "tencentcloud_dasb_cmd_template" "example" {
  name     = "tf_example"
  cmd_list = "rm -rf*"
}
`

const testAccDasbCmdTemplateUpdate = `
resource "tencentcloud_dasb_cmd_template" "example" {
  name     = "tf_example_update"
  cmd_list = "reboot"
}
`

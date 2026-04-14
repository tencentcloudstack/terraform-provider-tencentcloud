package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigRecorderConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRecorderConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_config_recorder_config.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_config_recorder_config.example", "status", "true"),
				),
			},
			{
				Config: testAccConfigRecorderConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_config_recorder_config.example", "status", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_config_recorder_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccConfigRecorderConfig = `
resource "tencentcloud_config_recorder_config" "example" {
  status = true
  resource_types = [
    "QCS::CAM::Group",
    "QCS::CAM::Role",
    "QCS::CVM::Instance",
  ]
}
`

const testAccConfigRecorderConfigUpdate = `
resource "tencentcloud_config_recorder_config" "example" {
  status = false
}
`

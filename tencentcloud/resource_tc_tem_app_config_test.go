package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTemAppConfig_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemAppConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tem_app_config.appConfig", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_app_config.appConfig",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemAppConfig = `

resource "tencentcloud_tem_app_config" "appConfig" {
  environment_id = "en-o5edaepv"
  name = "demo"
  config_data {
    key = "key"
    value = "value"
  }
  config_data {
    key = "key1"
    value = "value1"
  }
}

`

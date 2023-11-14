package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTemAppConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemAppConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tem_app_config.app_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_tem_app_config.app_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemAppConfig = `

resource "tencentcloud_tem_app_config" "app_config" {
  environment_id = "en-xxx"
  name = "xxx"
  data {
		key = "key"
		value = "value"

  }
}

`

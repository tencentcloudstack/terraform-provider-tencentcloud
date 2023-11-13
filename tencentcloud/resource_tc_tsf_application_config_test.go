package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_config.application_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_config.application_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationConfig = `

resource "tencentcloud_tsf_application_config" "application_config" {
  config_name = "my_config"
  config_version = "1.0"
  config_value = "test: 1"
  application_id = "app-123456"
  config_version_desc = "product version"
  config_type = "A"
  encode_with_base64 = true
  program_id_list = 
}

`

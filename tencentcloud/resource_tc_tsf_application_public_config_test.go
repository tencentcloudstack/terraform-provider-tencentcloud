package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationPublicConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationPublicConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_public_config.application_public_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_public_config.application_public_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationPublicConfig = `

resource "tencentcloud_tsf_application_public_config" "application_public_config" {
  config_name = "my_config"
  config_version = "1.0"
  config_value = "test: 1"
  config_version_desc = "product version"
  config_type = "P"
  encode_with_base64 = true
  program_id_list = 
}

`

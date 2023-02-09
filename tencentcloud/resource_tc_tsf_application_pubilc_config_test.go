package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationPubilcConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationPubilcConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_pubilc_config.application_pubilc_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_pubilc_config.application_pubilc_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationPubilcConfig = `

resource "tencentcloud_tsf_application_pubilc_config" "application_pubilc_config" {
  config_name = ""
  config_version = ""
  config_value = ""
  config_version_desc = ""
  config_type = ""
  encode_with_base64 = 
  program_id_list = 
                tags = {
    "createdBy" = "terraform"
  }
}

`

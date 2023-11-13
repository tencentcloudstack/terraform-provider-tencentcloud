package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationFileConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationFileConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_file_config.application_file_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_file_config.application_file_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationFileConfig = `

resource "tencentcloud_tsf_application_file_config" "application_file_config" {
  config_name = "my_config"
  config_version = "1.0"
  config_file_name = "application.yaml"
  config_file_value = "test: 1"
  application_id = "app-123456"
  config_file_path = "/etc/nginx"
  config_version_desc = "1.0"
  config_file_code = "utf-8"
  config_post_cmd = "source .bashrc"
  encode_with_base64 = true
  program_id_list = 
}

`

package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCamUserSamlConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUserSamlConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_user_saml_config.user_saml_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_user_saml_config.user_saml_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamUserSamlConfig = `

resource "tencentcloud_cam_user_saml_config" "user_saml_config" {
  saml_metadata_document = ""
}

`

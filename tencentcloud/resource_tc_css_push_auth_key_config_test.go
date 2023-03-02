package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudCssPushAuthKeyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssPushAuthKeyConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_push_auth_key_config.push_auth_key_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_css_push_auth_key_config.push_auth_key_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssPushAuthKeyConfig = `

resource "tencentcloud_css_push_auth_key_config" "push_auth_key_config" {
  domain_name = "5000.livepush.myqcloud.com"
  enable = 0
  master_auth_key = "xx"
  backup_auth_key = "xx"
  auth_delta = 60
}

`

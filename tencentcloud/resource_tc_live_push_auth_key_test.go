package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLivePushAuthKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLivePushAuthKey,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_push_auth_key.push_auth_key", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_push_auth_key.push_auth_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLivePushAuthKey = `

resource "tencentcloud_live_push_auth_key" "push_auth_key" {
  domain_name = "5000.livepush.myqcloud.com"
  enable = 0
  master_auth_key = "xx"
  backup_auth_key = "xx"
  auth_delta = 60
}

`

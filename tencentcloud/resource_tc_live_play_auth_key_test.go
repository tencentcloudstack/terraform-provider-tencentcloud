package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLivePlayAuthKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLivePlayAuthKey,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_play_auth_key.play_auth_key", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_play_auth_key.play_auth_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLivePlayAuthKey = `

resource "tencentcloud_live_play_auth_key" "play_auth_key" {
  domain_name = "5000.livepush.myqcloud.com"
  enable = 1
  auth_key = "xx"
  auth_delta = 60
  auth_back_key = "xx"
}

`

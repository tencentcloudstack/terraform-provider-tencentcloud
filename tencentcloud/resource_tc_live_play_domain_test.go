package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLivePlayDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLivePlayDomain,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_play_domain.play_domain", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_play_domain.play_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLivePlayDomain = `

resource "tencentcloud_live_play_domain" "play_domain" {
  domain_name = "5000.livepush.myqcloud.com"
  play_type = 1
}

`

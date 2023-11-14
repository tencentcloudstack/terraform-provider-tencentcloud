package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDomainRefererResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDomainReferer,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_domain_referer.domain_referer", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_domain_referer.domain_referer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveDomainReferer = `

resource "tencentcloud_live_domain_referer" "domain_referer" {
  domain_name = "5000.liveplay.myqcloud.com"
  enable = 1
  type = 1
  allow_empty = 1
  rules = ""
}

`

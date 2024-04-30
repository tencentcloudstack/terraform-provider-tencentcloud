package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafClbDomainResource_basic -v
func TestAccTencentCloudNeedFixWafClbDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafClbDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_clb_domain.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_clb_domain.example", "domain", "xxx.com"),
					resource.TestCheckResourceAttr("tencentcloud_waf_clb_domain.example", "region", "gz"),
					resource.TestCheckResourceAttr("tencentcloud_waf_clb_domain.example", "alb_type", "tsegw"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_clb_domain.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafClbDomainUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_clb_domain.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_clb_domain.example", "domain", "xxx.com"),
					resource.TestCheckResourceAttr("tencentcloud_waf_clb_domain.example", "region", "gz"),
					resource.TestCheckResourceAttr("tencentcloud_waf_clb_domain.example", "alb_type", "tsegw"),
				),
			},
		},
	})
}

const testAccWafClbDomain = `
resource "tencentcloud_waf_clb_domain" "example" {
  instance_id        = "waf_2kxtlbky11b2v4fe"
  domain             = "xxx.com"
  is_cdn             = 0
  status             = 1
  engine             = 12
  region             = "gz"
  flow_mode          = 0
  alb_type           = "tsegw"
  bot_status         = 0
  api_safe_status    = 0
  post_cls_action    = 0
  post_ckafka_action = 0
}
`

const testAccWafClbDomainUpdate = `
resource "tencentcloud_waf_clb_domain" "example" {
  instance_id        = "waf_2kxtlbky11b2v4fe"
  domain             = "xxx.com"
  is_cdn             = 0
  status             = 1
  engine             = 12
  region             = "gz"
  flow_mode          = 0
  alb_type           = "tsegw"
  bot_status         = 0
  api_safe_status    = 0
  post_cls_action    = 1
  post_ckafka_action = 1
}
`

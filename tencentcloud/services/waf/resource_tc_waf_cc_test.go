package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafCcResource_basic -v
func TestAccTencentCloudWafCcResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCc,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_cc.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "domain", "keep.qcloudwaf.com"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "name", "terraform"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "advance", "0"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "limit", "60"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "interval", "60"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "url", "/cc_demo"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "action_type", "22"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "edition", "sparta-waf"),
				),
			},
			{
				Config: testAccWafCcUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_cc.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "domain", "keep.qcloudwaf.com"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "name", "terraform"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "advance", "0"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "limit", "60"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "interval", "60"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "url", "/cc_demo"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "action_type", "22"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc.example", "edition", "sparta-waf"),
				),
			},
		},
	})
}

const testAccWafCc = `
resource "tencentcloud_waf_cc" "example" {
  domain      = "keep.qcloudwaf.com"
  name        = "terraform"
  status      = 1
  advance     = "0"
  limit       = "60"
  interval    = "60"
  url         = "/cc_demo"
  match_func  = 0
  action_type = "22"
  priority    = 50
  valid_time  = 600
  edition     = "sparta-waf"
  type        = 1
}
`

const testAccWafCcUpdate = `
resource "tencentcloud_waf_cc" "example" {
  domain      = "keep.qcloudwaf.com"
  name        = "terraform"
  status      = 0
  advance     = "0"
  limit       = "60"
  interval    = "60"
  url         = "/cc_demo"
  match_func  = 0
  action_type = "22"
  priority    = 50
  valid_time  = 600
  edition     = "sparta-waf"
  type        = 1
}
`

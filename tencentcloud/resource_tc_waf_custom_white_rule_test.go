package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafCustomWhiteRuleResource_basic -v
func TestAccTencentCloudWafCustomWhiteRuleResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCustomWhiteRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_custom_white_rule.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_custom_white_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafCustomWhiteRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_custom_white_rule.example", "id"),
				),
			},
		},
	})
}

const testAccWafCustomWhiteRule = `
resource "tencentcloud_waf_custom_white_rule" "example" {
  name        = "tf-example"
  sort_id     = "30"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "1.1.1.1"
    arg          = ""
  }

  status = "1"
  domain = "test.com"
  bypass = "geoip,cc,owasp"
}
`

const testAccWafCustomWhiteRuleUpdate = `
resource "tencentcloud_waf_custom_white_rule" "example" {
  name        = "tf-example-update"
  sort_id     = "50"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "2.2.2.2"
    arg          = ""
  }

  status = "1"
  domain = "test.com"
  bypass = "cc,owasp"
}
`

package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafCustomRuleResource_basic -v
func TestAccTencentCloudWafCustomRuleResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCustomRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_custom_rule.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_custom_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafCustomRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_custom_rule.example", "id"),
				),
			},
		},
	})
}

const testAccWafCustomRule = `
resource "tencentcloud_waf_custom_rule" "example" {
  name        = "tf-example"
  sort_id     = "50"
  redirect    = "/"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "1.1.1.1"
    arg          = ""
  }

  status      = "0"
  domain      = "test.com"
  action_type = "1"
}
`

const testAccWafCustomRuleUpdate = `
resource "tencentcloud_waf_custom_rule" "example" {
  name        = "tf-example-update"
  sort_id     = "80"
  redirect    = "/"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "2.2.2.2"
    arg          = ""
  }

  status      = "1"
  domain      = "test.com"
  action_type = "2"
}
`

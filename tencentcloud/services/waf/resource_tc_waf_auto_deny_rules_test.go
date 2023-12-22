package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafAutoDenyRulesResource_basic -v
func TestAccTencentCloudWafAutoDenyRulesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafAutoDenyRules,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_auto_deny_rules.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_auto_deny_rules.example", "attack_threshold"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_auto_deny_rules.example", "time_threshold"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_auto_deny_rules.example", "deny_time_threshold"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_auto_deny_rules.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafAutoDenyRules = `
resource "tencentcloud_waf_auto_deny_rules" "example" {
  domain              = "keep.qcloudwaf.com"
  attack_threshold    = 20
  time_threshold      = 12
  deny_time_threshold = 5
}
`

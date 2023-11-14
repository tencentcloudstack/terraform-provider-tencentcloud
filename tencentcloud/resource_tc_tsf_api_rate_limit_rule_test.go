package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApiRateLimitRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiRateLimitRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApiRateLimitRule = `

resource "tencentcloud_tsf_api_rate_limit_rule" "api_rate_limit_rule" {
  api_id = ""
  max_qps = 
  usable_status = ""
  }

`

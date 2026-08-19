package ga2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudGa2ForwardingRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGa2ForwardingRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_forwarding_rule.example", "id"),
				),
			},
			{
				Config: testAccGa2ForwardingRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_forwarding_rule.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_ga2_forwarding_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGa2ForwardingRule = `
resource "tencentcloud_ga2_forwarding_rule" "example" {
  global_accelerator_id = "ga-fhhs8w84"
  listener_id           = "lsr-dyy8jhzp"
  forwarding_policy_id  = "dm-rjssxr8k"

  rule_conditions {
    rule_condition_type  = "Path"
    rule_condition_value = ["/path"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = "epg-nt4iwozo"
  }
}
`

const testAccGa2ForwardingRuleUpdate = `
resource "tencentcloud_ga2_forwarding_rule" "example" {
  global_accelerator_id = "ga-fhhs8w84"
  listener_id           = "lsr-dyy8jhzp"
  forwarding_policy_id  = "dm-rjssxr8k"

  rule_conditions {
    rule_condition_type  = "Path"
    rule_condition_value = ["/path123"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = "epg-nt4iwozo"
  }
}
`

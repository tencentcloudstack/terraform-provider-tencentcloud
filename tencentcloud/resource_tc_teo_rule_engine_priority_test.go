package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoRuleEnginePriority_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRuleEnginePriority,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine_priority.rule_engine_priority", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_rule_engine_priority.ruleEnginePriority",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoRuleEnginePriority = `

resource "tencentcloud_teo_rule_engine_priority" "rule_engine_priority" {
  rules_priority = ""
}

`

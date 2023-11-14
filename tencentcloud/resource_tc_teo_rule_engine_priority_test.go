package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoRuleEnginePriorityResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRuleEnginePriority,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine_priority.rule_engine_priority", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_rule_engine_priority.rule_engine_priority",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoRuleEnginePriority = `

resource "tencentcloud_teo_rule_engine_priority" "rule_engine_priority" {
  rules_priority = &lt;nil&gt;
}

`

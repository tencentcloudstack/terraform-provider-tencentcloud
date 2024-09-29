package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoFunctionRulePriorityResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionRulePriority,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule_priority.teo_function_rule_priority", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule_priority.teo_function_rule_priority", "function_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule_priority.teo_function_rule_priority", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule_priority.teo_function_rule_priority", "rule_ids.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function_rule_priority.teo_function_rule_priority",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoFunctionRulePriority = `

resource "tencentcloud_teo_function_rule_priority" "teo_function_rule_priority" {
    function_id = "ef-txx7fnua"
    rule_ids    = [
        "rule-equpbht3",
        "rule-ax28n3g6",
    ]
    zone_id     = "zone-2qtuhspy7cr6"
}
`

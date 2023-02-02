package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTsfUnitRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfUnitRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_unit_rule.unit_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_unit_rule.unit_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfUnitRule = `

resource "tencentcloud_tsf_unit_rule" "unit_rule" {
  gateway_instance_id = ""
  name = ""
      description = ""
  unit_rule_item_list {
		relationship = ""
		dest_namespace_id = ""
		dest_namespace_name = ""
		name = ""
		id = ""
		unit_rule_id = ""
		priority = 
		description = ""
		unit_rule_tag_list {
			tag_type = ""
			tag_field = ""
			tag_operator = ""
			tag_value = ""
			unit_rule_item_id = ""
			id = ""
		}

  }
}

`

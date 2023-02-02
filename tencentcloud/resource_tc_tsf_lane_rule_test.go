package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTsfLaneRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfLaneRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_lane_rule.lane_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_lane_rule.lane_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfLaneRule = `

resource "tencentcloud_tsf_lane_rule" "lane_rule" {
    rule_name = ""
  remark = ""
  rule_tag_list {
		tag_id = ""
		tag_name = ""
		tag_operator = ""
		tag_value = ""
		lane_rule_id = ""
		create_time = 
		update_time = 

  }
  rule_tag_relationship = ""
  lane_id = ""
          program_id_list = 
}

`

package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudNeedFixTsfLaneRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfLaneRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfLaneRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfLaneRuleExists("tencentcloud_tsf_lane_rule.lane_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_lane_rule.lane_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "rule_name", ""),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "remark", ""),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_lane_rule.lane_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfLaneRuleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_lane_rule" {
			continue
		}

		res, err := service.DescribeTsfLaneRuleById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf lane rule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfLaneRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfLaneRuleById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf lane rule %s is not found", rs.Primary.ID)
		}

		return nil
	}
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

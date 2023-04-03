package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfLaneRuleResource_basic -v
func TestAccTencentCloudTsfLaneRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfLaneRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfLaneRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfLaneRuleExists("tencentcloud_tsf_lane_rule.lane_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_lane_rule.lane_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "rule_name", "terraform-rule-name"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "remark", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "rule_tag_relationship", "RELEATION_AND"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "enable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "rule_tag_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "rule_tag_list.0.tag_name", "xxx"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "rule_tag_list.0.tag_operator", "EQUAL"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane_rule.lane_rule", "rule_tag_list.0.tag_value", "222"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_lane_rule.lane_rule",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
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

resource "tencentcloud_tsf_lane" "lane1" {
	lane_name = "terraform-lane-1"
	remark = "lane desc1"
	lane_group_list {
		  group_id = "group-yrjkln9v"
		  entrance = true
	}
}

resource "tencentcloud_tsf_lane_rule" "lane_rule" {
	rule_name = "terraform-rule-name"
	remark = "terraform-test"
	rule_tag_list {
		  tag_name = "xxx"
		  tag_operator = "EQUAL"
		  tag_value = "222"
	}
	rule_tag_relationship = "RELEATION_AND"
	lane_id = tencentcloud_tsf_lane.lane1.id
	enable = false
  }

`

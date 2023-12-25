package tsf_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctsf "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tsf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfLaneRuleResource_basic -v
func TestAccTencentCloudTsfLaneRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

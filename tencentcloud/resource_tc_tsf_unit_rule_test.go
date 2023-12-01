package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTsfUnitRuleResource_basic -v
func TestAccTencentCloudTsfUnitRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfUnitRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfUnitRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfUnitRuleExists("tencentcloud_tsf_unit_rule.unit_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_unit_rule.unit_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "gateway_instance_id", defaultTsfGateway),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "description", "terraform-desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.relationship", "AND"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.dest_namespace_id", defaultTsfDestNamespaceId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.dest_namespace_name", "KEEP-terraform-请勿删除_default"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.name", "Rule1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.description", "rule1-desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.unit_rule_tag_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.unit_rule_tag_list.0.tag_type", "U"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.unit_rule_tag_list.0.tag_field", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.unit_rule_tag_list.0.tag_operator", "IN"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "unit_rule_item_list.0.unit_rule_tag_list.0.tag_value", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_unit_rule.unit_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfUnitRuleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_unit_rule" {
			continue
		}

		res, err := service.DescribeTsfUnitRuleById(ctx, rs.Primary.ID)
		if err != nil {
			code := err.(*sdkErrors.TencentCloudSDKError).Code
			if code == "InvalidParameterValue.GatewayParameterInvalid" {
				return nil
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf UnitRule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfUnitRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfUnitRuleById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf UnitRule %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfUnitRuleVar = `
variable "gateway_instance_id" {
	default = "` + defaultTsfGateway + `"
}
variable "dest_namespace_id" {
	default = "` + defaultTsfDestNamespaceId + `"
}
`

const testAccTsfUnitRule = testAccTsfUnitRuleVar + `

resource "tencentcloud_tsf_unit_rule" "unit_rule" {
	gateway_instance_id = var.gateway_instance_id
	name = "terraform-test"
	description = "terraform-desc"
	unit_rule_item_list {
		  relationship = "AND"
		  dest_namespace_id = var.dest_namespace_id
		  dest_namespace_name = "KEEP-terraform-请勿删除_default"
		  name = "Rule1"
		  description = "rule1-desc"
		  unit_rule_tag_list {
			  tag_type = "U"
			  tag_field = "aaa"
			  tag_operator = "IN"
			  tag_value = "1"
		}
	}
}

`

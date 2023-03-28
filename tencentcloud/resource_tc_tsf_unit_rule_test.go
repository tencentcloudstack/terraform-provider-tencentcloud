package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTsfUnitRuleResource_basic -v
func TestAccTencentCloudNeedFixTsfUnitRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfUnitRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfUnitRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfUnitRuleExists("tencentcloud_tsf_unit_rule.unit_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_unit_rule.unit_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "gateway_instance_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_rule.unit_rule", "name", ""),
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

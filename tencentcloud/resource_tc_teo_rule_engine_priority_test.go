package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoRuleEnginePriority_basic -v
func TestAccTencentCloudTeoRuleEnginePriority_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRuleEnginePriority,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleEnginePriorityExists("tencentcloud_teo_rule_engine_priority.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine_priority.basic", "zone_id", defaultZoneId),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_rule_engine_priority.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRuleEnginePriorityExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		originGroup, err := service.DescribeTeoRuleEnginePriority(ctx, rs.Primary.ID)
		if originGroup == nil {
			return fmt.Errorf("zone ruleEnginePriority %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoRuleEnginePriorityVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}
`

const testAccTeoRuleEnginePriority = testAccTeoRuleEnginePriorityVar + `

resource "tencentcloud_teo_rule_engine_priority" "basic" {
  zone_id = var.zone_id

  rules_priority {
    index = 0
    value = "rule-2a26njgck4rq"
  }
  rules_priority {
    index = 1
    value = "rule-2a26lzxhu13m"
  }
}

`

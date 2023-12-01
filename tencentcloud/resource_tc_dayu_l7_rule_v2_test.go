package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testDayuL7RuleV2ResourceName = "tencentcloud_dayu_l7_rule_v2"
var testDayuL7RuleV2ResourceKey = testDayuL7RuleV2ResourceName + ".test_rule"

func TestAccTencentCloudDayuL7RuleV2Resource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATIONAL) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuL7RuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuL7RuleV2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuL7RuleV2Exists(testDayuL7RuleV2ResourceKey),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "rule.#", "1"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "resource_id", "bgpip-000004xe"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "resource_ip", "119.28.217.162"),
				),
			},
		},
	})
}

func testAccCheckDayuL7RuleV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testDayuL4RuleV2ResourceNameTCP {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of dayu L7 rule")
		}
		business := items[0]
		domain := items[1]
		protocol := items[2]

		extendParams := make(map[string]interface{})
		extendParams["domain"] = domain
		extendParams["protocol"] = protocol
		dayuService := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		rules, _, err := dayuService.DescribeL7RulesV2(ctx, business, extendParams)
		if err != nil {
			return err
		}
		if len(rules) == 0 {
			err := fmt.Errorf("Create l7 rule failed.")
			return err
		}
		ruleItem := rules[0]
		resourceId := *ruleItem.Id
		resourceIp := *ruleItem.Ip
		ruleId := *ruleItem.RuleId
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := dayuService.DeleteL7RulesV2(ctx, business, resourceId, resourceIp, ruleId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckDayuL7RuleV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of dayu L7 rule")
		}
		business := items[0]
		domain := items[1]
		protocol := items[2]

		extendParams := make(map[string]interface{})
		extendParams["domain"] = domain
		extendParams["protocol"] = protocol
		dayuService := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		rules, _, err := dayuService.DescribeL7RulesV2(ctx, business, extendParams)

		if err != nil {
			return err
		}
		if len(rules) > 0 {
			return nil
		} else {
			return fmt.Errorf("L7 rule %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuL7RuleV2 string = `
resource "tencentcloud_dayu_l7_rule_v2" "test_rule" {
	resource_type="bgpip"
	resource_id="bgpip-000004xe"
	resource_ip="119.28.217.162"
	rule {
	  keep_enable=0
	  keeptime=0
	  source_list {
		source="1.2.3.5"
		weight=100
	  }
	  source_list {
		source="1.2.3.6"
		weight=100
	  }
	  lb_type=1
	  protocol="http"
	  source_type=2
	  domain="github.com"
	}
  }
`

package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testDayuL4RuleV2ResourceNameTCP = "tencentcloud_dayu_l4_rule_v2"
var testDayuL4RuleV2ResourceKeyTCP = testDayuL4RuleV2ResourceNameTCP + ".test"

func TestAccTencentCloudDayuL4RuleV2Resource_TCP(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATIONAL) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuL4RuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuL4RuleV2_tdp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuL4RuleV2Exists(testDayuL4RuleV2ResourceKeyTCP),
					resource.TestCheckResourceAttr(testDayuL4RuleV2ResourceKeyTCP, "business", "bgpip"),
					resource.TestCheckResourceAttr(testDayuL4RuleV2ResourceKeyTCP, "rules.0.rule_name", "test"),
					resource.TestCheckResourceAttr(testDayuL4RuleV2ResourceKeyTCP, "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testDayuL4RuleV2ResourceKeyTCP, "rules.0.source_list.#", "1"),
					resource.TestCheckResourceAttr(testDayuL4RuleV2ResourceKeyTCP, "rules.0.virtual_port", "2020"),
					resource.TestCheckResourceAttr(testDayuL4RuleV2ResourceKeyTCP, "vpn", "162.62.163.50"),
					resource.TestCheckResourceAttr(testDayuL4RuleV2ResourceKeyTCP, "rules.0.source_port", "20"),
				),
			},
		},
	})
}

func testAccCheckDayuL4RuleV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testDayuL4RuleV2ResourceNameTCP {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 4 {
			return fmt.Errorf("broken ID of L4 rule")
		}
		business := items[0]
		resourceId := items[1]
		vpn := items[2]
		virtualPortStr := items[3]
		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		extendParams := make(map[string]interface{})
		extendParams["ip"] = vpn
		virtualPort, err := strconv.Atoi(virtualPortStr)
		if err != nil {
			log.Printf("virtual_port must be int.")
		}
		extendParams["virtual_port"] = virtualPort
		result, err := service.DescribeNewL4Rules(ctx, business, extendParams)
		if len(result) > 0 && err != nil {
			delRuleId := *result[0].RuleId
			err = service.DeleteNewL4Rules(ctx, business, resourceId, vpn, []string{delRuleId})
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckDayuL4RuleV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 4 {
			return fmt.Errorf("broken ID of L4 rule")
		}
		business := items[0]
		vpn := items[2]
		virtualPortStr := items[3]
		virtualPort, err := strconv.Atoi(virtualPortStr)
		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		extendParams := make(map[string]interface{})
		extendParams["ip"] = vpn
		if err != nil {
			log.Printf("virtual_port must be int.")
		}
		extendParams["virtual_port"] = virtualPort
		result, err := service.DescribeNewL4Rules(ctx, business, extendParams)

		if err != nil {
			return err
		}
		if len(result) > 0 {
			return nil
		} else {
			return fmt.Errorf("L4 rule %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuL4RuleV2_tdp string = `
resource "tencentcloud_dayu_l4_rule_v2" "test" {
	business     = "bgpip"
	resource_id  = "bgpip-000004xg"
	vpn          = "162.62.163.50"
	virtual_port = 2020
  
	rules {
	  keep_enable   = false
	  keeptime      = 0
	  lb_type       = 1
	  protocol      = "TCP"
	  region        = 5
	  remove_switch = false
	  rule_name     = "test"
	  source_list {
		source = "1.2.3.9"
		weight = 0
	  }
	  source_port  = 20
	  source_type  = 2
	  virtual_port = 2020
	}
  }
`

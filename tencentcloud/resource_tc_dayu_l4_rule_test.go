package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testDayuL4RuleResourceName = "tencentcloud_dayu_l4_rule"
var testDayuL4RuleResourceKey = testDayuL4RuleResourceName + ".test_rule"
var testDayuL4RuleResourceNameUdp = "tencentcloud_dayu_l4_rule"
var testDayuL4RuleResourceKeyUdp = testDayuL4RuleResourceNameUdp + ".test_rule"

func TestAccTencentCloudDayuL4RuleResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuL4RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuL4Rule, defaultDayuBgpIp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuL4RuleExists(testDayuL4RuleResourceKey),
					resource.TestCheckResourceAttrSet(testDayuL4RuleResourceKey, "rule_id"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "name", "rule_test"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "source_type", "2"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "source_list.#", "2"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "protocol", "TCP"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "s_port", "80"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "d_port", "60"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "lb_type", "1"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_switch", "true"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_interval", "35"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_timeout", "30"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_health_num", "5"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_unhealth_num", "10"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "session_switch", "false"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "session_time", "300"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDayuL4RuleUpdate, defaultDayuBgpIp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuL4RuleExists(testDayuL4RuleResourceKey),
					testAccCheckDayuL4RuleExists(testDayuL4RuleResourceKey),
					testAccCheckDayuL4RuleExists(testDayuL4RuleResourceKey),
					resource.TestCheckResourceAttrSet(testDayuL4RuleResourceKey, "rule_id"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "name", "rule_test"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "source_type", "1"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "source_list.#", "3"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "protocol", "TCP"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "s_port", "800"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "d_port", "600"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_switch", "false"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_interval", "45"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_timeout", "25"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_health_num", "3"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "health_check_unhealth_num", "3"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "session_switch", "true"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKey, "session_time", "30"),
				),
			},
		},
	})
}

func TestAccTencentCloudDayuL4RuleResource_UDP(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuL4RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuL4Rule_udp, defaultDayuNet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuL4RuleExists(testDayuL4RuleResourceKeyUdp),
					resource.TestCheckResourceAttrSet(testDayuL4RuleResourceKeyUdp, "rule_id"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "resource_type", "net"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "name", "zhaoshaona.com"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "source_type", "2"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "source_list.#", "2"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "protocol", "UDP"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "s_port", "80"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "d_port", "60"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "lb_type", "1"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "health_check_switch", "true"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "health_check_interval", "35"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "health_check_timeout", "30"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "health_check_health_num", "5"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "health_check_unhealth_num", "10"),
					resource.TestCheckResourceAttr(testDayuL4RuleResourceKeyUdp, "session_switch", "false"),
				),
			},
		},
	})
}

func testAccCheckDayuL4RuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testDayuL4RuleResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of L4 rule")
		}
		resourceType := items[0]
		resourceId := items[1]
		ruleId := items[2]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, _, has, err := service.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
		if err != nil {
			_, _, has, err = service.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete L4 rule %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuL4RuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of L4 rule")
		}
		resourceType := items[0]
		resourceId := items[1]
		ruleId := items[2]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, _, has, err := service.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
		if err != nil {
			_, _, has, err = service.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("L4 rule %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuL4Rule string = `
resource "tencentcloud_dayu_l4_rule" "test_rule" {
  resource_type         = "bgpip"
  resource_id 			= "%s"
  name					= "rule_test"
  protocol				= "TCP"
  s_port			= 80
  d_port				= 60
  source_type			= 2
  health_check_switch	= true
  health_check_timeout	= 30
  health_check_interval = 35
  health_check_health_num = 5
  health_check_unhealth_num = 10
  session_switch 			= false
  session_time				= 300

  source_list{
	source = "1.1.1.1"
	weight = 100
  }
  source_list{
	source = "2.2.2.2"
	weight = 50
  }
}
`
const testAccDayuL4RuleUpdate string = `
resource "tencentcloud_dayu_l4_rule" "test_rule" {
  resource_type         = "bgpip"
  resource_id 			= "%s"
  name					= "rule_test"
  protocol				= "TCP"
  s_port			= 800
  d_port				= 600
  source_type			= 1
  health_check_switch	= false
  health_check_timeout	= 25
  health_check_interval = 45
  health_check_health_num = 3
  health_check_unhealth_num = 3
  session_switch 			= true
  session_time				= 30

  source_list{
	source = "test1.com"
	weight = 100
  }
  source_list{
	source = "test2.com"
	weight = 50
  }
  source_list{
	source = "test3.com"
	weight = 50
  }
}
`

const testAccDayuL4Rule_udp string = `
resource "tencentcloud_dayu_l4_rule" "test_rule" {
  resource_type         = "net"
  resource_id 			= "%s"
  name					= "zhaoshaona.com"
  protocol				= "UDP"
  s_port			= 80
  d_port				= 60
  source_type			= 2
  health_check_switch	= true
  health_check_timeout	= 30
  health_check_interval = 35
  health_check_health_num = 5
  health_check_unhealth_num = 10

  source_list{
	source = "1.1.1.1"
	weight = 100
  }
  source_list{
	source = "2.2.2.2"
	weight = 50
  }
}
`

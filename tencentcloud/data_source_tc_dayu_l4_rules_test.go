package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataDayuL4RulesName = "data.tencentcloud_dayu_l4_rules.id_test"

func TestAccTencentCloudDataDayuL4Rules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuL4RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataDayuL4RulesBaic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuL4RuleExists("tencentcloud_dayu_l4_rule.test_rule"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.name", "rule_testt"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.source_type", "2"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.source_port", "80"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.dest_port", "60"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.lb_type", "1"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.health_check_switch", "true"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.health_check_interval", "35"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.health_check_timeout", "30"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.health_check_health_num", "5"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.health_check_unhealth_num", "10"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.session_switch", "true"),
					resource.TestCheckResourceAttr(testDataDayuL4RulesName, "list.0.session_time", "300"),
				),
			},
		},
	})
}

const testAccTencentCloudDataDayuL4RulesBaic = `
resource "tencentcloud_dayu_l4_rule" "test_rule" {
  resource_type         = "bgpip"
  resource_id 			= "bgpip-00000294"
  name					= "rule_testt"
  protocol				= "TCP"
  source_port			= 80
  dest_port				= 60
  source_type			= 2
  health_check_switch	= true
  health_check_timeout	= 30
  health_check_interval = 35
  health_check_health_num = 5
  health_check_unhealth_num = 10
  session_switch 			= true
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

data "tencentcloud_dayu_l4_rules" "id_test" {
  resource_type = tencentcloud_dayu_l4_rule.test_rule.resource_type
  resource_id      = tencentcloud_dayu_l4_rule.test_rule.resource_id
  name		= tencentcloud_dayu_l4_rule.test_rule.name
}
`

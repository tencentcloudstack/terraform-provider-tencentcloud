package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataDayuL7RulesName = "data.tencentcloud_dayu_l7_rules.id_test"

func TestAccTencentCloudDataDayuL7Rules(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuL7RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTencentCloudDataDayuL7RulesBaic, defaultDayuBgpIp, defaultSshCertificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuL7RuleExists("tencentcloud_dayu_l7_rule.test_rule"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.name", "rule_test"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.switch", "true"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.domain", "zhaoshaona.com"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.source_type", "2"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.protocol", "https"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.ssl_id", defaultSshCertificate),
					resource.TestCheckResourceAttrSet(testDataDayuL7RulesName, "list.0.status"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.source_list.#", "2"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.health_check_switch", "true"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.health_check_interval", "30"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.health_check_code", "31"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.health_check_health_num", "5"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.health_check_unhealth_num", "10"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.health_check_path", "/"),
					resource.TestCheckResourceAttr(testDataDayuL7RulesName, "list.0.health_check_method", "GET"),
				),
			},
		},
	})
}

const testAccTencentCloudDataDayuL7RulesBaic = `
resource "tencentcloud_dayu_l7_rule" "test_rule" {
  resource_type         = "bgpip"
  resource_id 			= "%s"
  name					= "rule_test"
  domain				= "zhaoshaona.com"
  protocol				= "https"
  switch				= true
  source_type			= 2
  source_list 			= ["1.1.1.1:80","2.2.2.2"]
  ssl_id				= "%s"
  health_check_switch	= true
  health_check_code		= 31
  health_check_interval = 30
  health_check_method	= "GET"
  health_check_path		= "/"
  health_check_health_num = 5
  health_check_unhealth_num = 10
}

data "tencentcloud_dayu_l7_rules" "id_test" {
  resource_type = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id      = tencentcloud_dayu_l7_rule.test_rule.resource_id
  domain		= tencentcloud_dayu_l7_rule.test_rule.domain

}
`

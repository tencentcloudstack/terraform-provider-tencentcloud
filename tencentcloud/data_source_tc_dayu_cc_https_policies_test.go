package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataDayuCCHttpsPoliciesName = "data.tencentcloud_dayu_cc_https_policies.id_test"

func TestAccTencentCloudDataDayuCCHttpsPolicies(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpsPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTencentCloudDataDayuCCHttpsPoliciesBasic, defaultDayuBgpIp, defaultSshCertificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuCCHttpsPolicyExists("tencentcloud_dayu_cc_https_policy.test_policy"),
					resource.TestCheckResourceAttr(testDataDayuCCHttpsPoliciesName, "list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataDayuCCHttpsPoliciesName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataDayuCCHttpsPoliciesName, "list.0.policy_id"),
					resource.TestCheckResourceAttr(testDataDayuCCHttpsPoliciesName, "list.0.name", "policy_test"),
					resource.TestCheckResourceAttr(testDataDayuCCHttpsPoliciesName, "list.0.exe_mode", "drop"),
					resource.TestCheckResourceAttr(testDataDayuCCHttpsPoliciesName, "list.0.rule_list.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataDayuCCHttpsPoliciesBasic = `
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

resource "tencentcloud_dayu_cc_https_policy" "test_policy" {
  resource_type         = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id 			= tencentcloud_dayu_l7_rule.test_rule.resource_id
  rule_id				= tencentcloud_dayu_l7_rule.test_rule.rule_id
  domain				= tencentcloud_dayu_l7_rule.test_rule.domain
  name					= "policy_test"
  exe_mode				= "drop"
  switch				= true

  rule_list {
	skey 				= "cgi"
	operator			= "include"
	value				= "123"
	}
}

data "tencentcloud_dayu_cc_https_policies" "id_test" {
  resource_type = tencentcloud_dayu_cc_https_policy.test_policy.resource_type
  resource_id = tencentcloud_dayu_cc_https_policy.test_policy.resource_id
  name = tencentcloud_dayu_cc_https_policy.test_policy.name
}
`

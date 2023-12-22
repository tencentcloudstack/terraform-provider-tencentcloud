package dayu_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataDayuCCHttpPoliciesName = "data.tencentcloud_dayu_cc_http_policies.id_test"

func TestAccTencentCloudDataDayuCCHttpPolicies(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDayuCCHttpPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTencentCloudDataDayuCCHttpPoliciesBaic, tcacctest.DefaultDayuBgpIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists("tencentcloud_dayu_cc_http_policy.test_policy"),
					resource.TestCheckResourceAttr(testDataDayuCCHttpPoliciesName, "list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataDayuCCHttpPoliciesName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataDayuCCHttpPoliciesName, "list.0.policy_id"),
					resource.TestCheckResourceAttr(testDataDayuCCHttpPoliciesName, "list.0.name", "policy_match"),
					resource.TestCheckResourceAttr(testDataDayuCCHttpPoliciesName, "list.0.smode", "matching"),
					resource.TestCheckResourceAttr(testDataDayuCCHttpPoliciesName, "list.0.action", "drop"),
					resource.TestCheckResourceAttr(testDataDayuCCHttpPoliciesName, "list.0.rule_list.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataDayuCCHttpPoliciesBaic = `
resource "tencentcloud_dayu_cc_http_policy" "test_policy" {
  resource_type         = "bgpip"
  resource_id 			= "%s"
  name					= "policy_match"
  smode					= "matching"
  action				= "drop"
  switch				= true

  rule_list {
	skey 				= "host"
	operator			= "include"
	value				= "123"
	}
}

data "tencentcloud_dayu_cc_http_policies" "id_test" {
  resource_type = tencentcloud_dayu_cc_http_policy.test_policy.resource_type
  resource_id = tencentcloud_dayu_cc_http_policy.test_policy.resource_id
  policy_id      = tencentcloud_dayu_cc_http_policy.test_policy.policy_id
}
`

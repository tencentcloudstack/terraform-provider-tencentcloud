package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosCcPrecisionPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosCcPrecisionPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "domain", "t.baidu.com"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "instance_id", "bgpip-0000078h"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "ip", "212.64.62.191"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "policy_action", "alg"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "policy_list.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "protocol", "http"),
				),
			},
			{
				Config: testAccAntiddosCcPrecisionPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "policy_action", "drop"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_precision_policy.cc_precision_policy", "policy_list.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_cc_precision_policy.cc_precision_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosCcPrecisionPolicy = `

resource "tencentcloud_antiddos_cc_precision_policy" "cc_precision_policy" {
	instance_id   = "bgpip-0000078h"
	ip            = "212.64.62.191"
	protocol      = "http"
	domain        = "t.baidu.com"
	policy_action = "alg"
	policy_list {
	  field_type     = "value"
	  field_name     = "cgi"
	  value          = "a.com"
	  value_operator = "equal"
	}
  
	policy_list {
	  field_type     = "value"
	  field_name     = "ua"
	  value          = "test"
	  value_operator = "equal"
	}
  }

`

const testAccAntiddosCcPrecisionPolicyUpdate = `

resource "tencentcloud_antiddos_cc_precision_policy" "cc_precision_policy" {
	instance_id   = "bgpip-0000078h"
	ip            = "212.64.62.191"
	protocol      = "http"
	domain        = "t.baidu.com"
	policy_action = "drop"
	policy_list {
	  field_type     = "value"
	  field_name     = "cgi"
	  value          = "a.com"
	  value_operator = "equal"
	}
  }

`

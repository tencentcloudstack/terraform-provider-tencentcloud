package cfw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCfwNatPolicyOrderConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwNatPolicyOrderConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy_order_config.example", "id"),
				),
			},
			{
				Config: testAccCfwNatPolicyOrderConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy_order_config.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_nat_policy_order_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwNatPolicyOrderConfig = `
resource "tencentcloud_cfw_nat_policy" "in_example1" {
  source_content = "1.1.1.1/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "in_example2" {
  source_content = "3.3.3.3/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "ANY"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "in_example3" {
  source_content = "6.6.6.6/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "UDP"
  rule_action    = "accept"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "out_example1" {
  source_content = "1.1.1.1/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "out_example2" {
  source_content = "3.3.3.3/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "ANY"
  rule_action    = "accept"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}


resource "tencentcloud_cfw_nat_policy_order_config" "example" {
  inbound_rule_uuid_list = [
    tencentcloud_cfw_nat_policy.in_example3.uuid,
    tencentcloud_cfw_nat_policy.in_example1.uuid,
    tencentcloud_cfw_nat_policy.in_example2.uuid,
  ]

  outbound_rule_uuid_list = [
    tencentcloud_cfw_nat_policy.out_example2.uuid,
    tencentcloud_cfw_nat_policy.out_example1.uuid,
  ]
}
`

const testAccCfwNatPolicyOrderConfigUpdate = `
resource "tencentcloud_cfw_nat_policy" "in_example1" {
  source_content = "1.1.1.1/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "in_example2" {
  source_content = "3.3.3.3/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "ANY"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "in_example3" {
  source_content = "6.6.6.6/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "UDP"
  rule_action    = "accept"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "out_example1" {
  source_content = "1.1.1.1/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy" "out_example2" {
  source_content = "3.3.3.3/24"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "ANY"
  rule_action    = "accept"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "remark."
  scope          = "ALL"
}


resource "tencentcloud_cfw_nat_policy_order_config" "example" {
  inbound_rule_uuid_list = [
    tencentcloud_cfw_nat_policy.in_example2.uuid,
    tencentcloud_cfw_nat_policy.in_example3.uuid,
    tencentcloud_cfw_nat_policy.in_example1.uuid,
  ]

  outbound_rule_uuid_list = [
  	tencentcloud_cfw_nat_policy.out_example1.uuid,
    tencentcloud_cfw_nat_policy.out_example2.uuid,
  ]
}
`

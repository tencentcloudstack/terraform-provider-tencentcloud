package cfw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCfwEdgePolicyOrderConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwEdgePolicyOrderConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy_order_config.example", "id"),
				),
			},
			{
				Config: testAccCfwEdgePolicyOrderConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy_order_config.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_edge_policy_order_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwEdgePolicyOrderConfig = `
resource "tencentcloud_cfw_edge_policy" "in_example1" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy" "in_example2" {
  source_content = "2.2.2.2/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy" "in_example3" {
  source_content = "3.3.3.3/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy" "out_example1" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy" "out_example2" {
  source_content = "2.2.2.2/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy_order_config" "example" {
  inbound_rule_uuid_list = [
    tencentcloud_cfw_edge_policy.in_example3.uuid,
    tencentcloud_cfw_edge_policy.in_example1.uuid,
    tencentcloud_cfw_edge_policy.in_example2.uuid,
  ]

  outbound_rule_uuid_list = [
    tencentcloud_cfw_edge_policy.out_example2.uuid,
    tencentcloud_cfw_edge_policy.out_example1.uuid,
  ]
}
`

const testAccCfwEdgePolicyOrderConfigUpdate = `
resource "tencentcloud_cfw_edge_policy" "in_example1" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy" "in_example2" {
  source_content = "2.2.2.2/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy" "in_example3" {
  source_content = "3.3.3.3/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy" "out_example1" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy" "out_example2" {
  source_content = "2.2.2.2/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 0
  enable         = "true"
  description    = "policy description."
  scope          = "all"
}

resource "tencentcloud_cfw_edge_policy_order_config" "example" {
  inbound_rule_uuid_list = [
    tencentcloud_cfw_edge_policy.in_example2.uuid,
    tencentcloud_cfw_edge_policy.in_example3.uuid,
    tencentcloud_cfw_edge_policy.in_example1.uuid,
  ]

  outbound_rule_uuid_list = [
    tencentcloud_cfw_edge_policy.out_example2.uuid,
    tencentcloud_cfw_edge_policy.out_example1.uuid,
  ]
}
`

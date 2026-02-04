package cfw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCfwVpcPolicyOrderConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwVpcPolicyOrderConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy_order_config.example", "id"),
				),
			},
			{
				Config: testAccCfwVpcPolicyOrderConfigUpdate,
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

const testAccCfwVpcPolicyOrderConfig = `
resource "tencentcloud_cfw_vpc_policy" "example1" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  dest_content   = "192.168.0.2"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}

resource "tencentcloud_cfw_vpc_policy" "example2" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  dest_content   = "192.168.0.2"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}

resource "tencentcloud_cfw_vpc_policy" "example3" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  dest_content   = "192.168.0.2"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}

resource "tencentcloud_cfw_nat_policy_order_config" "example" {
  rule_uuid_list = [
    tencentcloud_cfw_vpc_policy.in_example3.uuid,
    tencentcloud_cfw_vpc_policy.in_example1.uuid,
    tencentcloud_cfw_vpc_policy.in_example2.uuid,
  ]
}
`

const testAccCfwVpcPolicyOrderConfigUpdate = `
resource "tencentcloud_cfw_vpc_policy" "example1" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  dest_content   = "192.168.0.2"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}

resource "tencentcloud_cfw_vpc_policy" "example2" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  dest_content   = "192.168.0.2"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}

resource "tencentcloud_cfw_vpc_policy" "example3" {
  source_content = "0.0.0.0/0"
  source_type    = "net"
  dest_content   = "192.168.0.2"
  dest_type      = "net"
  protocol       = "ANY"
  rule_action    = "log"
  port           = "-1/-1"
  description    = "description."
  enable         = "true"
  fw_group_id    = "ALL"
}

resource "tencentcloud_cfw_nat_policy_order_config" "example" {
  rule_uuid_list = [
    tencentcloud_cfw_vpc_policy.in_example2.uuid,
    tencentcloud_cfw_vpc_policy.in_example3.uuid,
    tencentcloud_cfw_vpc_policy.in_example1.uuid,
  ]
}
`

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
resource "tencentcloud_cfw_nat_policy" "example" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "111"
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy_order_config" "example" {
  uuid        = tencentcloud_cfw_nat_policy.example.id
  order_index = 1
}
`

const testAccCfwNatPolicyOrderConfigUpdate = `
resource "tencentcloud_cfw_nat_policy" "example" {
  source_content = "1.1.1.1/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "111"
  scope          = "ALL"
}

resource "tencentcloud_cfw_nat_policy_order_config" "example" {
  uuid        = tencentcloud_cfw_nat_policy.example.id
  order_index = 2
}
`

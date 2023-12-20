package cfw_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwNatPolicyResource_basic -v
func TestAccTencentCloudNeedFixCfwNatPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwNatPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "source_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "source_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "target_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "target_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "protocol"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "rule_action"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "port"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "direction"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "enable"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "description"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_nat_policy.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfwNatPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "source_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "source_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "target_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "target_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "protocol"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "rule_action"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "port"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "direction"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "enable"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_policy.example", "description"),
				),
			},
		},
	})
}

const testAccCfwNatPolicy = `
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
  description    = "policy description."
}
`

const testAccCfwNatPolicyUpdate = `
resource "tencentcloud_cfw_nat_policy" "example" {
  source_content = "2.2.2.2/0"
  source_type    = "net"
  target_content = "3.3.3.3/0"
  target_type    = "net"
  protocol       = "TCP"
  rule_action    = "drop"
  port           = "-1/-1"
  direction      = 1
  enable         = "true"
  description    = "policy description update."
}
`

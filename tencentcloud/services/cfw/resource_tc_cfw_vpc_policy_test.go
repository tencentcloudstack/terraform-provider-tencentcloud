package cfw_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwVpcPolicyResource_basic -v
func TestAccTencentCloudNeedFixCfwVpcPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwVpcPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_policy.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_vpc_policy.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfwVpcPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_policy.example", "id"),
				),
			},
		},
	})
}

const testAccCfwVpcPolicy = `
resource "tencentcloud_cfw_vpc_policy" "example" {
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
`

const testAccCfwVpcPolicyUpdate = `
resource "tencentcloud_cfw_vpc_policy" "example" {
  source_content = "1.1.1.1"
  source_type    = "net"
  dest_content   = "2.2.2.2"
  dest_type      = "net"
  protocol       = "TCP"
  rule_action    = "accept"
  port           = "-1/-1"
  description    = "description update."
  enable         = "false"
  fw_group_id    = "ALL"
}
`

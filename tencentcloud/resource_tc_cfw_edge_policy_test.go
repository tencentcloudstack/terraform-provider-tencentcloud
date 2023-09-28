package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwEdgePolicyResource_basic -v
func TestAccTencentCloudNeedFixCfwEdgePolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwEdgePolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "source_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "source_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "target_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "target_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "protocol"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "rule_action"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "port"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "direction"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "enable"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "scope"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_edge_policy.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfwEdgePolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "source_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "source_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "target_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "target_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "protocol"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "rule_action"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "port"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "direction"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "enable"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_policy.example", "scope"),
				),
			},
		},
	})
}

const testAccCfwEdgePolicy = `
resource "tencentcloud_cfw_edge_policy" "example" {
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
`

const testAccCfwEdgePolicyUpdate = `
resource "tencentcloud_cfw_edge_policy" "example" {
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
  scope          = "all"
}
`

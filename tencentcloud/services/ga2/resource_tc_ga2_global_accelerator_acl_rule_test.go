package ga2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudGa2GlobalAcceleratorAclRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGa2GlobalAcceleratorAclRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_global_accelerator_acl_rule.example", "id"),
				),
			},
			{
				Config: testAccGa2GlobalAcceleratorAclRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_global_accelerator_acl_rule.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_ga2_global_accelerator_acl_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGa2GlobalAcceleratorAclRule = `
resource "tencentcloud_ga2_global_accelerator_acl_rule" "example" {
  global_accelerator_id           = "ga-xxxxxxxx"
  global_accelerator_acl_policy_id = "aclpol-xxxxxxxx"
  protocol                        = "TCP"
  port                            = "80"
  source_cidr_block               = "10.0.0.0/24"
  policy                          = "ACCEPT"
  description                     = "tf example acl rule"
}
`

const testAccGa2GlobalAcceleratorAclRuleUpdate = `
resource "tencentcloud_ga2_global_accelerator_acl_rule" "example" {
  global_accelerator_id           = "ga-xxxxxxxx"
  global_accelerator_acl_policy_id = "aclpol-xxxxxxxx"
  protocol                        = "TCP"
  port                            = "443"
  source_cidr_block               = "10.0.1.0/24"
  policy                          = "DROP"
  description                     = "tf example acl rule updated"
}
`

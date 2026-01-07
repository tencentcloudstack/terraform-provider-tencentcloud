package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVpcPrivateNatGatewayTranslationAclRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPrivateNatGatewayTranslationAclRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway_translation_acl_rule.example", "id"),
				),
			},
			{
				Config: testAccVpcPrivateNatGatewayTranslationAclRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway_translation_acl_rule.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_private_nat_gateway_translation_acl_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcPrivateNatGatewayTranslationAclRule = `
resource "tencentcloud_vpc_private_nat_gateway_translation_acl_rule" "example" {
  nat_gateway_id        = "intranat-bw389ya1"
  translation_direction = "LOCAL"
  translation_type      = "NETWORK_LAYER"
  translation_ip        = "2.2.2.2"
  original_ip           = "1.1.1.1"
  translation_acl_rules {
    protocol         = "TCP"
    source_port      = "80"
    destination_port = "8080"
    destination_cidr = "8.8.8.8"
    description      = "remark."
  }
}
`

const testAccVpcPrivateNatGatewayTranslationAclRuleUpdate = `
resource "tencentcloud_vpc_private_nat_gateway_translation_acl_rule" "example" {
  nat_gateway_id        = "intranat-bw389ya1"
  translation_direction = "LOCAL"
  translation_type      = "NETWORK_LAYER"
  translation_ip        = "2.2.2.2"
  original_ip           = "1.1.1.1"
  translation_acl_rules {
    protocol         = "TCP"
    source_port      = "90"
    destination_port = "9090"
    destination_cidr = "9.9.9.9"
    description      = "remark update."
  }
}
`

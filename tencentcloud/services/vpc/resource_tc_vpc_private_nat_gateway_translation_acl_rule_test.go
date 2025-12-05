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
  nat_gateway_id        = ""
  translation_direction = ""
  translation_type      = ""
  translation_ip        = ""
  original_ip           = ""
  translation_acl_rule {
    protocol         = ""
    source_port      = ""
    source_cidr      = ""
    destination_port = ""
    destination_cidr = ""
    action           = ""
    description      = ""
  }
}
`

const testAccVpcPrivateNatGatewayTranslationAclRuleUpdate = `
resource "tencentcloud_vpc_private_nat_gateway_translation_acl_rule" "example" {
  nat_gateway_id        = ""
  translation_direction = ""
  translation_type      = ""
  translation_ip        = ""
  original_ip           = ""
  translation_acl_rule {
    protocol         = ""
    source_port      = ""
    source_cidr      = ""
    destination_port = ""
    destination_cidr = ""
    action           = ""
    description      = ""
  }
}
`

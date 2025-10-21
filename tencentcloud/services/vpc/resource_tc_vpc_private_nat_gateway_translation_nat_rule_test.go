package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVpcPrivateNatGatewayTranslationNatRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPrivateNatGatewayTranslationNatRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example", "nat_gateway_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example", "translation_nat_rules.#"),
				),
			},
			{
				Config: testAccVpcPrivateNatGatewayTranslationNatRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example", "nat_gateway_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example", "translation_nat_rules.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_private_nat_gateway_translation_nat_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcPrivateNatGatewayTranslationNatRule = `
resource "tencentcloud_vpc_private_nat_gateway_translation_nat_rule" "example" {
  nat_gateway_id = "intranat-r46f6pxl"
  translation_nat_rules {
    translation_direction = "LOCAL"
    translation_type      = "NETWORK_LAYER"
    translation_ip        = "2.2.2.2"
    description           = "remark."
    original_ip           = "1.1.1.1"
  }

  translation_nat_rules {
    translation_direction = "LOCAL"
    translation_type      = "TRANSPORT_LAYER"
    translation_ip        = "3.3.3.3"
    description           = "remark."
  }
}
`

const testAccVpcPrivateNatGatewayTranslationNatRuleUpdate = `
resource "tencentcloud_vpc_private_nat_gateway_translation_nat_rule" "example" {
  nat_gateway_id = "intranat-r46f6pxl"
  translation_nat_rules {
    translation_direction = "LOCAL"
    translation_type      = "NETWORK_LAYER"
    translation_ip        = "2.2.2.2"
    description           = "remark."
    original_ip           = "1.1.1.1"
  }
}
`

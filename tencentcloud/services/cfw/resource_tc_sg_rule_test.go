package cfw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudSgRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSgRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sg_rule.sg_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "enable", "1"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.description", "111111"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.dest_content", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.dest_type", "net"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.port", "-1/-1"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.protocol", "ANY"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.rule_action", "accept"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.source_content", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.source_type", "net"),
				),
			},
			{
				ResourceName:      "tencentcloud_sg_rule.sg_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSgRuleUpEnable,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sg_rule.sg_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "enable", "0"),
				),
			},
			{
				Config: testAccSgRuleUpData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sg_rule.sg_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "enable", "0"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.description", "11111122"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.dest_content", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.dest_type", "net"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.port", "-1/-1"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.protocol", "ANY"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.rule_action", "accept"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.source_content", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("tencentcloud_sg_rule.sg_rule", "data.0.source_type", "net"),
				),
			},
		},
	})
}

const testAccSgRule = `

resource "tencentcloud_sg_rule" "sg_rule" {
  enable = 1

  data {
    description         = "111111"
    dest_content        = "0.0.0.0/0"
    dest_type           = "net"
    port                = "-1/-1"
    protocol            = "ANY"
    rule_action         = "accept"
    source_content      = "0.0.0.0/0"
    source_type         = "net"
  }
}
`
const testAccSgRuleUpEnable = `

resource "tencentcloud_sg_rule" "sg_rule" {
  enable = 0

  data {
    description         = "111111"
    dest_content        = "0.0.0.0/0"
    dest_type           = "net"
    port                = "-1/-1"
    protocol            = "ANY"
    rule_action         = "accept"
    source_content      = "0.0.0.0/0"
    source_type         = "net"
  }
}
`

const testAccSgRuleUpData = `

resource "tencentcloud_sg_rule" "sg_rule" {
  enable = 0

  data {
    description         = "11111122"
    dest_content        = "0.0.0.0/0"
    dest_type           = "net"
    port                = "-1/-1"
    protocol            = "ANY"
    rule_action         = "accept"
    source_content      = "0.0.0.0/0"
    source_type         = "net"
  }
}
`

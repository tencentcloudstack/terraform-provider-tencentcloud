package privatedns_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudPrivateDnsForwardRuleResource_basic -v
func TestAccTencentCloudPrivateDnsForwardRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsForwardRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "rule_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "rule_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "end_point_id"),
				),
			},
			{
				Config: testAccPrivateDnsForwardRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "rule_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "rule_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_forward_rule.example", "end_point_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_private_dns_forward_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPrivateDnsForwardRule = `
resource "tencentcloud_private_dns_forward_rule" "example" {
  rule_name    = "tf-example"
  rule_type    = "DOWN"
  zone_id      = "zone-cmmbvaq8"
  end_point_id = "eid-72dbbc79fb"
}
`

const testAccPrivateDnsForwardRuleUpdate = `
resource "tencentcloud_private_dns_forward_rule" "example" {
  rule_name    = "tf-example-update"
  rule_type    = "DOWN"
  zone_id      = "zone-cmmbvaq8"
  end_point_id = "eid-72dbbc79fb"
}
`

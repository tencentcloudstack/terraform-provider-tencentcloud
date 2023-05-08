package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoWafRuleGroupsDataSource -v
func TestAccTencentCloudTeoWafRuleGroupsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoWafRuleGroups,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_teo_waf_rule_groups.waf_rule_groups"),
				),
			},
		},
	})
}

const testAccDataSourceTeoWafRuleGroupsVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}

variable "entity" {
  default = "` + defaultZoneName + `"
}
`

const testAccDataSourceTeoWafRuleGroups = testAccDataSourceTeoWafRuleGroupsVar + `

data "tencentcloud_teo_waf_rule_groups" "waf_rule_groups" {
  zone_id = var.zone_id
  entity  = var.entity
}

`

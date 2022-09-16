package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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

const testAccDataSourceTeoWafRuleGroups = `

data "tencentcloud_teo_waf_rule_groups" "waf_rule_groups" {
  }

`

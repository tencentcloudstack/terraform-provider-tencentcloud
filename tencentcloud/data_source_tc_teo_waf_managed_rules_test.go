
package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoWafManagedRulesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoWafManagedRules,
				Check: resource.ComposeTestCheckFunc(
				  testAccCheckTencentCloudDataSourceID("data.tencentcloud_teo_waf_managed_rules.waf_managed_rules"),
				),
			},
		},
	})
}

const testAccDataSourceTeoWafManagedRules = `

data "tencentcloud_teo_waf_managed_rules" "waf_managed_rules" {
  zone_id = ""
  entity = ""
    }

`

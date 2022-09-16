package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoBotManagedRulesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoBotManagedRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_teo_bot_managed_rules.bot_managed_rules"),
				),
			},
		},
	})
}

const testAccDataSourceTeoBotManagedRules = `

data "tencentcloud_teo_bot_managed_rules" "bot_managed_rules" {
  zone_id = ""
  entity = ""
    }

`

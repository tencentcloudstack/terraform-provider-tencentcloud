package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoBotManagedRulesDataSource -v
func TestAccTencentCloudNeedFixTeoBotManagedRulesDataSource(t *testing.T) {
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

const testAccDataSourceTeoBotManagedRulesVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}

variable "entity" {
  default = "` + defaultZoneName + `"
}
`

const testAccDataSourceTeoBotManagedRules = testAccDataSourceTeoBotManagedRulesVar + `

data "tencentcloud_teo_bot_managed_rules" "bot_managed_rules" {
  zone_id = var.zone_id
  entity = var.entity
}

`

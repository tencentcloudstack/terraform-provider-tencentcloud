package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoBotPortraitRulesDataSource -v
func TestAccTencentCloudTeoBotPortraitRulesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoBotPortraitRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_teo_bot_portrait_rules.bot_portrait_rules"),
				),
			},
		},
	})
}

const testAccDataSourceTeoBotPortraitRulesVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}

variable "entity" {
  default = "` + defaultZoneName + `"
}
`

const testAccDataSourceTeoBotPortraitRules = testAccDataSourceTeoBotPortraitRulesVar + `

data "tencentcloud_teo_bot_portrait_rules" "bot_portrait_rules" {
  zone_id = var.zone_id
  entity = var.entity
}

`

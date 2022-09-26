package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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

const testAccDataSourceTeoBotPortraitRules = `

data "tencentcloud_teo_bot_portrait_rules" "bot_portrait_rules" {
  zone_id = ""
  entity = ""
    }

`

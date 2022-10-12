package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoRuleEngineSettingsDataSource -v
func TestAccTencentCloudTeoRuleEngineSettingsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoRuleEngineSettings,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_teo_rule_engine_settings.rule_engine_settings"),
				),
			},
		},
	})
}

const testAccDataSourceTeoRuleEngineSettings = `

data "tencentcloud_teo_rule_engine_settings" "rule_engine_settings" {
}

`

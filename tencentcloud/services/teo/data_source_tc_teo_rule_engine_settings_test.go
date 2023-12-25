package teo_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoRuleEngineSettingsDataSource -v
func TestAccTencentCloudTeoRuleEngineSettingsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoRuleEngineSettings,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_rule_engine_settings.rule_engine_settings"),
				),
			},
		},
	})
}

const testAccDataSourceTeoRuleEngineSettings = `

data "tencentcloud_teo_rule_engine_settings" "rule_engine_settings" {
}

`

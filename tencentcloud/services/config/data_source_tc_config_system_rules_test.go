package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigSystemRulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigSystemRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_system_rules.example"),
				),
			},
		},
	})
}

func TestAccTencentCloudConfigSystemRulesDataSource_withFilters(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigSystemRulesDataSourceWithFilters,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_system_rules.example_with_filters"),
				),
			},
		},
	})
}

const testAccConfigSystemRulesDataSource = `
data "tencentcloud_config_system_rules" "example" {}
`

const testAccConfigSystemRulesDataSourceWithFilters = `
data "tencentcloud_config_system_rules" "example_with_filters" {
  risk_level = 1
}
`

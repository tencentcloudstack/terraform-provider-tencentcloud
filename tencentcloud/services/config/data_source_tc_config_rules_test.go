package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigRulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_rules.example"),
				),
			},
		},
	})
}

func TestAccTencentCloudConfigRulesDataSource_withFilters(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRulesDataSourceWithFilters,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_rules.example_with_filters"),
				),
			},
		},
	})
}

const testAccConfigRulesDataSource = `
data "tencentcloud_config_rules" "example" {}
`

const testAccConfigRulesDataSourceWithFilters = `
data "tencentcloud_config_rules" "example_with_filters" {
  state      = "ACTIVE"
  risk_level = [1, 2]
  order_type = "desc"
}
`

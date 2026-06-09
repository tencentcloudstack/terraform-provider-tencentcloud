package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigCompliancePacksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigCompliancePacksDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_compliance_packs.example"),
				),
			},
		},
	})
}

func TestAccTencentCloudConfigCompliancePacksDataSource_withFilters(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigCompliancePacksDataSourceWithFilters,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_compliance_packs.example_with_filters"),
				),
			},
		},
	})
}

const testAccConfigCompliancePacksDataSource = `
data "tencentcloud_config_compliance_packs" "example" {}
`

const testAccConfigCompliancePacksDataSourceWithFilters = `
data "tencentcloud_config_compliance_packs" "example_with_filters" {
  risk_level  = [1, 2]
  status      = "ACTIVE"
  order_type  = "desc"
}
`

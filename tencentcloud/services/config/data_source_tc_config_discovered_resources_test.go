package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigDiscoveredResourcesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDiscoveredResourcesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_discovered_resources.example"),
				),
			},
		},
	})
}

func TestAccTencentCloudConfigDiscoveredResourcesDataSource_withFilters(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDiscoveredResourcesDataSourceWithFilters,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_discovered_resources.example_with_filters"),
				),
			},
		},
	})
}

const testAccConfigDiscoveredResourcesDataSource = `
data "tencentcloud_config_discovered_resources" "example" {}
`

const testAccConfigDiscoveredResourcesDataSourceWithFilters = `
data "tencentcloud_config_discovered_resources" "example_with_filters" {
  order_type = "desc"
}
`

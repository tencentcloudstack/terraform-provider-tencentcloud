package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMonitorProductNamesapce(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorProductNamespace(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_product_namespace.instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_product_namespace.instances",
						"list.#"),
				),
			},
		},
	})
}

func testAccDataSourceMonitorProductNamespace() string {
	return `
data "tencentcloud_monitor_product_namespace" "instances" {
  name = "Redis"
}`
}

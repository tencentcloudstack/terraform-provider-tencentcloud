package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonitorProductNamesapce(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorProductNamespace(),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_product_namespace.instances"),
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

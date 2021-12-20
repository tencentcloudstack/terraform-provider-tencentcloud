package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudKubernetesChartsDataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.tencentcloud_kubernetes_charts.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubernetesCharts,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "chart_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceKubernetesCharts = `
data "tencentcloud_kubernetes_charts" "test" {
}
`

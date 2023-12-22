package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesChartsDataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.tencentcloud_kubernetes_charts.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
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

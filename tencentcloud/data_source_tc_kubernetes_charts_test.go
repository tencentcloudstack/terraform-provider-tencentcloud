package tencentcloud

import (
	"regexp"
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
					resource.TestCheckResourceAttrSet(dataSourceName, "locked_versions.cbs"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locked_versions.cos"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locked_versions.tcr"),
				),
			},
			{
				Config: testAccDataSourceKubernetesChartsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "locked_versions.cbs"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locked_versions.cos"),
					resource.TestCheckResourceAttr(dataSourceName, "locked_versions.tcr", "__fake_version__"),
				),
			},
			{
				Config: testAccDataSourceKubernetesChartsUpdate2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "locked_versions.cbs"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locked_versions.cos"),
					resource.TestMatchResourceAttr(dataSourceName, "locked_versions.tcr", regexp.MustCompile(`^\d+\.\d+\.\d+$`)),
				),
			},
		},
	})
}

const testAccDataSourceKubernetesCharts = `
data "tencentcloud_kubernetes_charts" "test" {
}
`

const testAccDataSourceKubernetesChartsUpdate = `
data "tencentcloud_kubernetes_charts" "test" {
  locked_versions = { tcr: "__fake_version__" }
}
`

const testAccDataSourceKubernetesChartsUpdate2 = `
data "tencentcloud_kubernetes_charts" "test" {
  update_locked_versions = true
  locked_versions = { tcr: "__fake_version__" }
}
`

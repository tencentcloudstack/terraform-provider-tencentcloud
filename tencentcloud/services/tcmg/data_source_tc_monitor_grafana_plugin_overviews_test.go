package tcmg_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorGrafanaPluginOverviewsDataSource_basic -v
func TestAccTencentCloudMonitorGrafanaPluginOverviewsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaPluginOverviewsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_grafana_plugin_overviews.plugin_overviews"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_grafana_plugin_overviews.plugin_overviews", "plugin_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_grafana_plugin_overviews.plugin_overviews", "plugin_set.0.plugin_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_grafana_plugin_overviews.plugin_overviews", "plugin_set.0.version"),
				),
			},
		},
	})
}

const testAccMonitorGrafanaPluginOverviewsDataSource = `

data "tencentcloud_monitor_grafana_plugin_overviews" "plugin_overviews" {
}

`

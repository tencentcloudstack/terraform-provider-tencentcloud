package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorPluginOverviewsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorPluginOverviewsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_plugin_overviews.plugin_overviews")),
			},
		},
	})
}

const testAccMonitorPluginOverviewsDataSource = `

data "tencentcloud_monitor_plugin_overviews" "plugin_overviews" {
  }

`

package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorStatisticDataDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorStatisticDataDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_statistic_data.statistic_data")),
			},
		},
	})
}

const testAccMonitorStatisticDataDataSource = `

data "tencentcloud_monitor_statistic_data" "statistic_data" {
  module = ""
  namespace = ""
  metric_names = 
  conditions {
		key = ""
		operator = ""
		value = 

  }
        group_bys = 
  }

`

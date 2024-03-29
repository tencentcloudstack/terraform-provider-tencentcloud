package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMonitorStatisticDataDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorStatisticDataDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_statistic_data.statistic_data"),
				),
			},
		},
	})
}

const testAccMonitorStatisticDataDataSource = `

data "tencentcloud_monitor_statistic_data" "statistic_data" {
  module       = "monitor"
  namespace    = "QCE/TKE2"
  metric_names = ["cpu_usage"]
  conditions {
    key      = "tke_cluster_instance_id"
    operator = "="
    value    = ["cls-mw2w40s7"]
  }
}

`

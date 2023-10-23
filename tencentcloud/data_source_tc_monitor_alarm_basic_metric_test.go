package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlarmBasicMetricDataSource_basic -v
func TestAccTencentCloudMonitorAlarmBasicMetricDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmBasicMetricDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.dimensions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.dimensions.0.dimensions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.meaning.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.meaning.0.zh"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.metric_c_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.metric_e_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.period.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.periods.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.periods.0.period"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.periods.0.stat_type.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_metric.alarm_metric", "metric_set.0.unit"),
				),
			},
		},
	})
}

const testAccMonitorAlarmBasicMetricDataSource = `

data "tencentcloud_monitor_alarm_basic_metric" "alarm_metric" {
  namespace   = "qce/cvm"
  metric_name = "WanOuttraffic"
  dimensions  = ["uuid"]
}

`

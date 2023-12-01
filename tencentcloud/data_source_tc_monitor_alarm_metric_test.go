package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlarmMetricDataSource_basic -v
func TestAccTencentCloudMonitorAlarmMetricDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmMetricDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_metric.alarm_metric"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.dimensions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.is_advanced"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.is_open"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.metric_config.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.metric_config.0.continue_period.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.metric_config.0.operator.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.metric_config.0.period.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.min"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.product_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_metric.alarm_metric", "metrics.0.unit"),
				),
			},
		},
	})
}

const testAccMonitorAlarmMetricDataSource = `

data "tencentcloud_monitor_alarm_metric" "alarm_metric" {
  module       = "monitor"
  monitor_type = "Monitoring"
  namespace    = "cvm_device"
}

`

package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlarmHistoryDataSource_basic -v
func TestAccTencentCloudMonitorAlarmHistoryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmHistoryDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_history.alarm_history"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.alarm_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.alarm_object"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.alarm_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.alarm_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.content"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.dimensions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.first_occur_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.last_occur_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.metrics_info.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.metrics_info.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.metrics_info.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.metrics_info.0.period"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.metrics_info.0.qce_namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.metrics_info.0.value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.monitor_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.origin_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.policy_exists"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.policy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.policy_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.project_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_history.alarm_history", "histories.0.vpc"),
				),
			},
		},
	})
}

const testAccMonitorAlarmHistoryDataSource = `

data "tencentcloud_monitor_alarm_history" "alarm_history" {
  module        = "monitor"
  order         = "DESC"
  start_time    = 1696608000
  end_time      = 1697212799
  monitor_types = ["MT_QCE"]
  project_ids   = [0]
  namespaces {
    monitor_type = "CpuUsage"
    namespace    = "cvm_device"
  }
  policy_name = "terraform_test"
  content     = "CPU利用率 > 3%"
  policy_ids  = ["policy-iejtp4ue"]
}

`

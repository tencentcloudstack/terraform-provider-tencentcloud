package monitor_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlamBasicAlarmsDataSource_basic -v
func TestAccTencentCloudMonitorAlamBasicAlarmsDataSource_basic(t *testing.T) {
	t.Parallel()

	startTime := time.Now().AddDate(0, 0, -1).Unix()
	endTime := time.Now().Unix()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMonitorAlarmBasciAlarmsDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_basic_alarms.alarms"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.alarm_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.alarm_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.content"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.dimensions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.duration"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.first_occur_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.last_occur_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.metric_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.obj_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.obj_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.project_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.view_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_basic_alarms.alarms", "alarms.0.vpc"),
				),
			},
		},
	})
}

const testAccMonitorAlarmBasciAlarmsDataSource = `

data "tencentcloud_monitor_alarm_basic_alarms" "alarms" {
  module             = "monitor"
  start_time         = %v
  end_time           = %v
  occur_time_order   = "DESC"
  project_ids        = [0]
  view_names         = ["cvm_device"]
  alarm_status       = [1]
  instance_group_ids = [5497073]
  metric_names       = ["cpu_usage"]
}

`

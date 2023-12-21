package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlarmAllNamespacesDataSource_basic -v
func TestAccTencentCloudMonitorAlarmAllNamespacesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmAllNamespacesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces", "common_namespaces.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces", "qce_namespaces_new.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces", "qce_namespaces_new.0.config"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces", "qce_namespaces_new.0.dashboard_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces", "qce_namespaces_new.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces", "qce_namespaces_new.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces", "qce_namespaces_new.0.product_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces", "qce_namespaces_new.0.sort_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_all_namespaces.alarm_all_namespaces", "qce_namespaces_new.0.value"),
				),
			},
		},
	})
}

const testAccMonitorAlarmAllNamespacesDataSource = `

data "tencentcloud_monitor_alarm_all_namespaces" "alarm_all_namespaces" {
  scene_type    = "ST_ALARM"
  module        = "monitor"
  monitor_types = ["MT_QCE"]
  ids           = ["qaap_tunnel_l4_listeners"]
}

`

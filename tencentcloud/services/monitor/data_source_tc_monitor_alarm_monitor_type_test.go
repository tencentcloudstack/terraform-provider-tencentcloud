package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlarmMonitorTypeDataSource_basic -v
func TestAccTencentCloudMonitorAlarmMonitorTypeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmMonitorTypeDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_monitor_type.alarm_monitor_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_monitor_type.alarm_monitor_type", "monitor_type_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_monitor_type.alarm_monitor_type", "monitor_type_infos.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_monitor_type.alarm_monitor_type", "monitor_type_infos.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_monitor_type.alarm_monitor_type", "monitor_type_infos.0.sort_id"),
				),
			},
		},
	})
}

const testAccMonitorAlarmMonitorTypeDataSource = `

data "tencentcloud_monitor_alarm_monitor_type" "alarm_monitor_type" {
}

`

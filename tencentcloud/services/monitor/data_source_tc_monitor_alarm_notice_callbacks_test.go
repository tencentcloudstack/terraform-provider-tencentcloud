package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlarmNoticeCallbacksDataSource_basic -v
func TestAccTencentCloudMonitorAlarmNoticeCallbacksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmNoticeCallbacksDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_notice_callbacks.alarm_notice_callbacks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notice_callbacks.alarm_notice_callbacks", "url_notices.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notice_callbacks.alarm_notice_callbacks", "url_notices.0.url"),
				),
			},
		},
	})
}

const testAccMonitorAlarmNoticeCallbacksDataSource = `

data "tencentcloud_monitor_alarm_notice_callbacks" "alarm_notice_callbacks" {
}

`

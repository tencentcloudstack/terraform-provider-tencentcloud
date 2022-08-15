package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMonitorAlarmNoticeResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmNotice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_monitor_alarm_notice.example", "name", "test_alarm_notice_123"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_alarm_notice.example", "notice_type", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_alarm_notice.example", "notice_language", "zh-CN"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_alarm_notice.example", "is_preset", "0"),
				),
			},
		},
	})
}

const testAccMonitorAlarmNotice string = `
resource "tencentcloud_monitor_alarm_notice" "example" {
  name                  = "test_alarm_notice_123"
  notice_type           = "ALL"
  notice_language       = "zh-CN"

}
`

package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMonitorAlarmNoticeResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmNotice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_monitor_alarm_notice.example", "name", "test_alarm_notice"),
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
  name            = "test_alarm_notice"
  notice_language = "zh-CN"
  notice_type     = "ALL"

  url_notices {
    end_time   = 86399
    is_valid = 0
    start_time = 0
    url        = "https://www.mytest.com/validate"
    weekday    = [
      1,
      2,
      3,
      4,
      5,
      6,
      7,
    ]
  }

  user_notices {
    end_time                 = 86399
    group_ids                = []
    need_phone_arrive_notice = 1
    notice_way               = [
      "EMAIL",
      "SMS",
    ]
    phone_call_type       = "CIRCLE"
    phone_circle_interval = 180
    phone_circle_times    = 2
    phone_inner_interval  = 180
    phone_order           = []
    receiver_type         = "USER"
    start_time            = 0
    user_ids              = [
      11082189,
      11082190,
    ]
    weekday = [
      1,
      2,
      3,
      4,
      5,
      6,
      7,
    ]
  }
}
`

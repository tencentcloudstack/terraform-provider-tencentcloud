package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccAlarmNoticesDatasourceBasic -v
func TestAccAlarmNoticesDatasourceBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlarmNotices(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_notices.notices"),
					resource.TestCheckResourceAttr("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.0.amp_consumer_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.0.is_preset"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.0.notice_language"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.0.notice_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.0.policy_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.0.updated_by"),
					resource.TestCheckResourceAttr("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.0.user_notices.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAlarmNotices() string {
	return `
data "tencentcloud_monitor_alarm_notices" "notices" {
  order      = "DESC"
  notice_ids = ["notice-f2svbu3w"]
}`
}

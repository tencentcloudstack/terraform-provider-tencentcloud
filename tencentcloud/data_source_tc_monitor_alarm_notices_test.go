package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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
					//resource.TestCheckResourceAttr("data.tencentcloud_monitor_alarm_notices.notices", "alarm_notice.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_monitor_alarm_notices.notices", "order", "DESC"),
					resource.TestCheckResourceAttr("data.tencentcloud_monitor_alarm_notices.notices", "owner_uid", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAlarmNotices() string {
	return `data "tencentcloud_monitor_alarm_notices" "notices" {
    order = "DESC"
    owner_uid = 1
    name = ""
    receiver_type = ""
    user_ids = []
    group_ids = []
    notice_ids = []
}`
}

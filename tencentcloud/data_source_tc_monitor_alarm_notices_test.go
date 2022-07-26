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
					resource.TestCheckResourceAttr("data.tencentcloud_monitor_alarm_notices.notices", "notices.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAlarmNotices() string {
	return `data "tencentcloud_monitor_alarm_notices" "notices" {
  module     = "monitor"
  pagenumber = 1
  pagesize   = 20
  order      = "DESC"
}`
}

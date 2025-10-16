package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsAlarmMessagesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsAlarmMessagesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_alarm_messages.wedata_ops_alarm_messages"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_alarm_messages.wedata_ops_alarm_messages", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_alarm_messages.wedata_ops_alarm_messages", "data.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_alarm_messages.wedata_ops_alarm_messages", "data.0.items.#", "1"),
				),
			},
		},
	})
}

const testAccWedataOpsAlarmMessagesDataSource = `

data "tencentcloud_wedata_ops_alarm_messages" "wedata_ops_alarm_messages" {
  project_id  = "1859317240494305280"
  start_time  = "2025-10-14 21:09:26"
  end_time    = "2025-10-14 21:10:26"
  alarm_level = 1
  time_zone   = "UTC+8"
}
`

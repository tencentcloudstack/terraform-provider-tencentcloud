package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsAlarmMessageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsAlarmMessageDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_alarm_message.wedata_ops_alarm_message"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_alarm_message.wedata_ops_alarm_message", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_alarm_message.wedata_ops_alarm_message", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataOpsAlarmMessageDataSource = `

data "tencentcloud_wedata_ops_alarm_message" "wedata_ops_alarm_message" {
  project_id  = "1859317240494305280"
  alarm_message_id  = 263840
}
`

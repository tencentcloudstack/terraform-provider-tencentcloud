package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClsAlarmNoticeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsAlarmNotice,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm_notice.alarm_notice", "id")),
			},
			{
				ResourceName:      "tencentcloud_cls_alarm_notice.alarm_notice",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsAlarmNotice = `

resource "tencentcloud_cls_alarm_notice" "alarm_notice" {
  name = "notice"
  type = "Trigger"
  notice_receivers {
		receiver_type = "Uin"
		receiver_ids = 
		receiver_channels = 
		start_time = "00:00:00"
		end_time = "23:59:59"
		index = 1

  }
  web_callbacks {
		url = "http://www.testnotice.com/callback"
		callback_type = "WeCom"
		method = "POST"
		headers = 
		body = "null"
		index = 10

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`

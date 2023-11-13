package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClsAlarmResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsAlarm,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm", "id")),
			},
			{
				ResourceName:      "tencentcloud_cls_alarm.alarm",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsAlarm = `

resource "tencentcloud_cls_alarm" "alarm" {
  name = "alarm"
  alarm_targets {
		topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
		query = "* | select count(*) as count"
		number = 1
		start_time_offset = 0
		end_time_offset = 0
		logset_id = "5cd3a17e-1111-418c-afd7-77b365397426"

  }
  monitor_time {
		type = "Period"
		time = 1

  }
  condition = "$1&gt;100"
  trigger_count = 5
  alarm_period = 5
  alarm_notice_ids = 
  status = true
  message_template = "test"
  call_back {
		body = "test"
		headers = 

  }
  analysis {
		name = "analysis"
		type = "query"
		content = "content"
		config_info {
			key = "key"
			value = "value"
		}

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`

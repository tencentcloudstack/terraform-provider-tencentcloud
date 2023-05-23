package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: testAccClsAlarmUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "name", "terraform-alarm-for-test"),
				),
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
  name             = "terraform-alarm-test"
  alarm_notice_ids = [
    "notice-0850756b-245d-4bc7-bb27-2a58fffc780b",
  ]
  alarm_period     = 15
  condition        = "test"
  message_template = "{{.Label}}"
  status           = true
  tags             = {
    "createdBy" = "terraform"
  }
  trigger_count = 1

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "33aaf0ae-6163-411b-a415-9f27450f68db"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "88735a07-bea4-4985-8763-e9deb6da4fad"
  }

  analysis {
    content = "__FILENAME__"
    name    = "terraform"
    type    = "field"

    config_info {
      key   = "QueryIndex"
      value = "1"
    }
  }

  monitor_time {
    time = 1
    type = "Period"
  }
}

`

const testAccClsAlarmUpdate = `

resource "tencentcloud_cls_alarm" "alarm" {
  name             = "terraform-alarm-for-test"
  alarm_notice_ids = [
    "notice-0850756b-245d-4bc7-bb27-2a58fffc780b",
  ]
  alarm_period     = 15
  condition        = "test"
  message_template = "{{.Label}}"
  status           = true
  tags             = {
    "createdBy" = "terraform"
  }
  trigger_count = 1

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "33aaf0ae-6163-411b-a415-9f27450f68db"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "88735a07-bea4-4985-8763-e9deb6da4fad"
  }

  analysis {
    content = "__FILENAME__"
    name    = "terraform"
    type    = "field"

    config_info {
      key   = "QueryIndex"
      value = "1"
    }
  }

  monitor_time {
    time = 1
    type = "Period"
  }
}

`

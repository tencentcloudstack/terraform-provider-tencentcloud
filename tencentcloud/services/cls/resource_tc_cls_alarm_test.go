package cls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudClsAlarmResource_basic -v
func TestAccTencentCloudClsAlarmResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsAlarm,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "condition", "test"),
				),
			},
			{
				Config: testAccClsAlarmUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "condition", "test update"),
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
  name             = "tf-example"
  alarm_notice_ids = [
    "notice-d365c616-1ae2-4a77-863a-9777453ab9d5",
  ]
  alarm_period     = 15
  condition        = "test"
  alarm_level      = 0
  message_template = "{{.Label}}"
  status           = true
  tags = {
    "createdBy" = "terraform"
  }
  trigger_count = 1

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "dac3e1a9-d22c-403b-a129-f94f666a33af"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "775c0bc2-2246-43a0-8eb2-f5bc248be183"
	syntax_rule       = 0
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
  name             = "tf-example-update"
  alarm_notice_ids = [
    "notice-d365c616-1ae2-4a77-863a-9777453ab9d5",
  ]
  alarm_period     = 15
  condition        = "test update"
  alarm_level      = 1
  message_template = "{{.Label}}"
  status           = true
  tags = {
    "createdBy" = "terraform"
  }
  trigger_count = 1

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "dac3e1a9-d22c-403b-a129-f94f666a33af"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "775c0bc2-2246-43a0-8eb2-f5bc248be183"
	syntax_rule       = 1
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

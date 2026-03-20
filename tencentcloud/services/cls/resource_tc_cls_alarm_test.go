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

// go test -i; go test -test.run TestAccTencentCloudClsAlarmResource_classifications -v
func TestAccTencentCloudClsAlarmResource_classifications(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsAlarmWithClassifications,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "name", "tf-example-classifications"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "classifications.%", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "classifications.category", "application"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "classifications.severity", "high"),
				),
			},
			{
				Config: testAccClsAlarmWithClassificationsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "name", "tf-example-classifications"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "classifications.%", "3"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "classifications.category", "application"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "classifications.severity", "critical"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "classifications.environment", "production"),
				),
			},
			{
				Config: testAccClsAlarmWithoutClassifications,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm", "name", "tf-example-no-classifications"),
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

// go test -i; go test -test.run TestAccTencentCloudClsAlarmResource_monitorNotice -v
func TestAccTencentCloudClsAlarmResource_monitorNotice(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsAlarmWithMonitorNotice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm_monitor", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm_monitor", "name", "tf-example-monitor-notice"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm_monitor", "monitor_notice.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm_monitor", "monitor_notice.0.notices.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm_monitor", "monitor_notice.0.notices.0.notice_id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm_monitor", "monitor_notice.0.notices.0.alarm_levels.#", "2"),
				),
			},
			{
				Config: testAccClsAlarmUpdateMonitorNotice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm_monitor", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm_monitor", "name", "tf-example-monitor-notice-update"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm_monitor", "monitor_notice.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm_monitor", "monitor_notice.0.notices.#", "2"),
				),
			},
			{
				Config: testAccClsAlarmSwitchToMonitorNotice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_alarm.alarm_switch", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_alarm.alarm_switch", "monitor_notice.#", "1"),
					resource.TestCheckNoResourceAttr("tencentcloud_cls_alarm.alarm_switch", "alarm_notice_ids"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_alarm.alarm_monitor",
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

const testAccClsAlarmWithClassifications = `
resource "tencentcloud_cls_alarm" "alarm" {
  name             = "tf-example-classifications"
  alarm_notice_ids = [
    "notice-d365c616-1ae2-4a77-863a-9777453ab9d5",
  ]
  alarm_period     = 15
  condition        = "test"
  alarm_level      = 0
  message_template = "{{.Label}}"
  status           = true
  trigger_count    = 1
  classifications = {
    category = "application"
    severity = "high"
  }

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "dac3e1a9-d22c-403b-a129-f94f666a33af"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "775c0bc2-2246-43a0-8eb2-f5bc248be183"
    syntax_rule       = 0
  }

  monitor_time {
    time = 1
    type = "Period"
  }
}
`

const testAccClsAlarmWithClassificationsUpdate = `
resource "tencentcloud_cls_alarm" "alarm" {
  name             = "tf-example-classifications"
  alarm_notice_ids = [
    "notice-d365c616-1ae2-4a77-863a-9777453ab9d5",
  ]
  alarm_period     = 15
  condition        = "test"
  alarm_level      = 0
  message_template = "{{.Label}}"
  status           = true
  trigger_count    = 1
  classifications = {
    category    = "application"
    severity    = "critical"
    environment = "production"
  }

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "dac3e1a9-d22c-403b-a129-f94f666a33af"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "775c0bc2-2246-43a0-8eb2-f5bc248be183"
    syntax_rule       = 0
  }

  monitor_time {
    time = 1
    type = "Period"
  }
}
`

const testAccClsAlarmWithoutClassifications = `
resource "tencentcloud_cls_alarm" "alarm" {
  name             = "tf-example-no-classifications"
  alarm_notice_ids = [
    "notice-d365c616-1ae2-4a77-863a-9777453ab9d5",
  ]
  alarm_period     = 15
  condition        = "test"
  alarm_level      = 0
  message_template = "{{.Label}}"
  status           = true
  trigger_count    = 1

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "dac3e1a9-d22c-403b-a129-f94f666a33af"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "775c0bc2-2246-43a0-8eb2-f5bc248be183"
    syntax_rule       = 0
  }

  monitor_time {
    time = 1
    type = "Period"
  }
}
`

const testAccClsAlarmWithMonitorNotice = `
resource "tencentcloud_cls_alarm" "alarm_monitor" {
  name          = "tf-example-monitor-notice"
  alarm_period  = 15
  condition     = "test"
  alarm_level   = 0
  status        = true
  trigger_count = 1

  monitor_notice {
    notices {
      notice_id    = "notice-12345678-1234-1234-1234-123456789012"
      alarm_levels = [0, 2]
    }
  }

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "dac3e1a9-d22c-403b-a129-f94f666a33af"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "775c0bc2-2246-43a0-8eb2-f5bc248be183"
    syntax_rule       = 0
  }

  monitor_time {
    time = 1
    type = "Period"
  }
}
`

const testAccClsAlarmUpdateMonitorNotice = `
resource "tencentcloud_cls_alarm" "alarm_monitor" {
  name          = "tf-example-monitor-notice-update"
  alarm_period  = 15
  condition     = "test update"
  alarm_level   = 1
  status        = true
  trigger_count = 1

  monitor_notice {
    notices {
      notice_id       = "notice-12345678-1234-1234-1234-123456789012"
      content_tmpl_id = "tmpl-87654321-4321-4321-4321-210987654321"
      alarm_levels    = [0, 1]
    }
    notices {
      notice_id    = "notice-abcdefgh-abcd-abcd-abcd-abcdefghijkl"
      alarm_levels = [2]
    }
  }

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "dac3e1a9-d22c-403b-a129-f94f666a33af"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "775c0bc2-2246-43a0-8eb2-f5bc248be183"
    syntax_rule       = 1
  }

  monitor_time {
    time = 1
    type = "Period"
  }
}
`

const testAccClsAlarmSwitchToMonitorNotice = `
resource "tencentcloud_cls_alarm" "alarm_switch" {
  name          = "tf-example-switch-notice"
  alarm_period  = 15
  condition     = "test"
  alarm_level   = 0
  status        = true
  trigger_count = 1

  monitor_notice {
    notices {
      notice_id    = "notice-12345678-1234-1234-1234-123456789012"
      alarm_levels = [0, 2]
    }
  }

  alarm_targets {
    end_time_offset   = 0
    logset_id         = "dac3e1a9-d22c-403b-a129-f94f666a33af"
    number            = 1
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    topic_id          = "775c0bc2-2246-43a0-8eb2-f5bc248be183"
    syntax_rule       = 0
  }

  monitor_time {
    time = 1
    type = "Period"
  }
}
`

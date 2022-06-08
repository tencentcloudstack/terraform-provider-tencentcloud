package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMonitorAlarmPolicyResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmPolicy,
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrSet("tencentcloud_monitor_alarm_policy.policy", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_alarm_policy.policy", "policy_name", "terraform"),
				),
			},
		},
	})
}

const testAccMonitorAlarmPolicy string = `
resource "tencentcloud_monitor_alarm_policy" "policy" {
  enable       = 1
  monitor_type = "MT_QCE"
  namespace    = "cvm_device"
  notice_ids   = [
    "notice-f2svbu3w",
  ]
  policy_name  = "terraform"
  project_id   = 0

  conditions {
    is_union_rule = 0

    rules {
      continue_period  = 5
      description      = "CPUUtilization"
      is_power_notice  = 0
      metric_name      = "CpuUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "PublicBandwidthUtilization"
      is_power_notice  = 0
      metric_name      = "Outratio"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "MemoryUtilization"
      is_power_notice  = 0
      metric_name      = "MemUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "DiskUtilization"
      is_power_notice  = 0
      metric_name      = "CvmDiskUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
  }

  event_conditions {
    continue_period  = 0
    description      = "DiskReadonly"
    is_power_notice  = 0
    metric_name      = "disk_readonly"
    notice_frequency = 0
    period           = 0
  }

  policy_tag {
    key   = "test-tag"
    value = "unit-test"
  }
}
`

package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsAlarmRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsAlarmRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_level", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_name", "tf_test"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_types.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_types.0", "failure"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "description", "ccc"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "monitor_object_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "monitor_object_ids.0", "20230906105118824"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "monitor_object_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "project_id", "1859317240494305280"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_escalation_interval", "15"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_recipient_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_recipient_ids.0", "100029411056"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_recipient_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_ways.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_ways.0", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.notify_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.notify_interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.0.days_of_week.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.0.days_of_week.0", "6"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.0.days_of_week.1", "7"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.0.end_time", "21:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.0.start_time", "10:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.0.data_backfill_or_rerun_trigger", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.0.trigger", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataOpsAlarmRuleUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_level", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_name", "tf_test_up"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_types.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_types.0", "failure"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_types.1", "reconciliationFailure"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "description", "qqq"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "monitor_object_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "monitor_object_ids.0", "3bec54e4-cd0a-4163-9318-65f0fe115ee9"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "monitor_object_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "project_id", "1859317240494305280"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_escalation_interval", "15"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_recipient_ids.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_recipient_ids.0", "100028448903"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_recipient_ids.1", "100029411056"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_recipient_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_ways.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_ways.0", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_ways.1", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_ways.2", "3"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.alarm_ways.3", "4"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.notify_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.notify_interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.0.days_of_week.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.0.days_of_week.0", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.0.end_time", "01:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_groups.0.notification_fatigue.0.quiet_intervals.0.start_time", "00:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.0.data_backfill_or_rerun_trigger", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.0.trigger", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.0.reconciliation_ext_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.0.reconciliation_ext_info.0.hour", "0"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.0.reconciliation_ext_info.0.min", "0"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.0.reconciliation_ext_info.0.mismatch_count", "0"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "alarm_rule_detail.0.reconciliation_ext_info.0.rule_type", "reconciliationFailure"),
				),
			},
		},
	})
}

const testAccWedataOpsAlarmRule = `

resource "tencentcloud_wedata_ops_alarm_rule" "wedata_ops_alarm_rule" {
    alarm_level         = 1
    alarm_rule_name     = "tf_test"
    alarm_types         = [
        "failure",
    ]
    description         = "ccc"
    monitor_object_ids  = [
        "20230906105118824",
    ]
    monitor_object_type = 1
    project_id          = "1859317240494305280"

    alarm_groups {
        alarm_escalation_interval      = 15
        alarm_escalation_recipient_ids = []
        alarm_recipient_ids            = [
            "100029411056",
        ]
        alarm_recipient_type           = 1
        alarm_ways                     = [
            "1",
        ]

        notification_fatigue {
            notify_count    = 1
            notify_interval = 5

            quiet_intervals {
                days_of_week = [
                    6,
                    7,
                ]
                end_time     = "21:00:00"
                start_time   = "10:00:00"
            }
        }
    }

    alarm_rule_detail {
        data_backfill_or_rerun_trigger = 1
        trigger                        = 2
    }
}
`

const testAccWedataOpsAlarmRuleUp = `
resource "tencentcloud_wedata_ops_alarm_rule" "wedata_ops_alarm_rule" {
    alarm_level         = 2
    alarm_rule_name     = "tf_test_up"
    alarm_types         = [
        "failure",
        "reconciliationFailure",
    ]
    description         = "qqq"
    monitor_object_ids  = [
        "3bec54e4-cd0a-4163-9318-65f0fe115ee9",
    ]
    monitor_object_type = 2
    project_id          = "1859317240494305280"

    alarm_groups {
        alarm_escalation_interval      = 15
        alarm_escalation_recipient_ids = []
        alarm_recipient_ids            = [
            "100028448903",
            "100029411056",
        ]
        alarm_recipient_type           = 1
        alarm_ways                     = [
            "1",
            "2",
            "3",
            "4",
        ]

        notification_fatigue {
            notify_count    = 1
            notify_interval = 5

            quiet_intervals {
                days_of_week = [
                    1,
                ]
                end_time     = "01:00:00"
                start_time   = "00:00:00"
            }
        }
    }

    alarm_rule_detail {
        data_backfill_or_rerun_trigger = 1
        trigger                        = 2

        reconciliation_ext_info {
            hour           = 0
            min            = 0
            mismatch_count = 0
            rule_type      = "reconciliationFailure"
        }
    }
}
`

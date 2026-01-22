---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_alarm_rule"
sidebar_current: "docs-tencentcloud-resource-wedata_ops_alarm_rule"
description: |-
  Provides a resource to create a wedata ops alarm rule
---

# tencentcloud_wedata_ops_alarm_rule

Provides a resource to create a wedata ops alarm rule

## Example Usage

```hcl
resource "tencentcloud_wedata_ops_alarm_rule" "wedata_ops_alarm_rule" {
  alarm_level     = 1
  alarm_rule_name = "tf_test"
  alarm_types = [
    "failure",
  ]
  description = "ccc"
  monitor_object_ids = [
    "20230906105118824",
  ]
  monitor_object_type = 1
  project_id          = "1859317240494305280"

  alarm_groups {
    alarm_escalation_interval      = 15
    alarm_escalation_recipient_ids = []
    alarm_recipient_ids = [
      "100029411056",
    ]
    alarm_recipient_type = 1
    alarm_ways = [
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
        end_time   = "21:00:00"
        start_time = "10:00:00"
      }
    }
  }

  alarm_rule_detail {
    data_backfill_or_rerun_trigger = 1
    trigger                        = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `alarm_groups` - (Required, List) Alarm receiver configuration information.
* `alarm_rule_name` - (Required, String) Alert rule name.
* `alarm_types` - (Required, Set: [`String`]) Alarm Rule Monitoring Types: failure: failure alarm; overtime: timeout alarm; success: success alarm; backTrackingOrRerunSuccess: backTrackingOrRerunSuccess: backTrackingOrRerunFailure: backTrackingOrRerunFailure. Project Fluctuation Alarms: projectFailureInstanceUpwardFluctuationAlarm: alarm if the upward fluctuation rate of failed instances exceeds the threshold. projectSuccessInstanceDownwardFluctuationAlarm: alarm if the downward fluctuation rate of successful instances exceeds the threshold. Offline Integration Task Reconciliation Alarms: reconciliationFailure: offline reconciliation task failure alarm; reconciliationOvertime: offline reconciliation task timeout alarm; reconciliationMismatch: alarm if the number of inconsistent entries in a data reconciliation task exceeds the threshold. Example value: ["failure"].
* `monitor_object_ids` - (Required, Set: [`String`]) A list of monitored object business IDs. Different business IDs are passed in based on the MonitorType setting. For example, 1 (Task) - MonitorObjectIds is a list of task IDs; 2 (Workflow) - MonitorObjectIds is a list of workflow IDs (workflow IDs can be obtained from the ListWorkflows interface); 3 (Project) - MonitorObjectIds is a list of project IDs. Example value: ["ddc"].
* `project_id` - (Required, String, ForceNew) Project id.
* `alarm_level` - (Optional, Int) Alarm level: 1. Normal, 2. Major, 3. Urgent (default 1. Normal).
* `alarm_rule_detail` - (Optional, List) Alarm rule configuration information: Success alarms do not require configuration. Failure alarms can be configured as either first-failure alarms or all retry failure alarms. Timeout configuration requires the timeout type and timeout threshold. Project fluctuation alarms require the fluctuation rate and anti-shake period.
* `description` - (Optional, String) Alarm rule description.
* `monitor_object_type` - (Optional, Int) Monitoring object type, Task-based monitoring: Configurable by task/workflow/project: 1. Task, 2. Workflow, 3. Project (default is 1. Task). Project-based monitoring: Alerts for overall project task fluctuations, 7: Project fluctuation monitoring alerts.

The `alarm_groups` object supports the following:

* `alarm_escalation_interval` - (Optional, Int) Alarm escalation interval.
* `alarm_escalation_recipient_ids` - (Optional, Set) Alarm escalator ID list. If the alarm receiver or the upper escalator does not confirm the alarm within the alarm interval, the alarm will be sent to the next level escalator.
* `alarm_recipient_ids` - (Optional, Set) Depending on the type of AlarmRecipientType, this list has different business IDs: 1 (Specified Person): Alarm Recipient ID List; 2 (Task Responsible Person): No configuration required; 3 (Duty Roster): Duty Roster ID List.
* `alarm_recipient_type` - (Optional, Int) Alarm Recipient Type: 1. Designated Personnel, 2. Task Responsible Personnel, 3. Duty Roster (Default: 1. Designated Personnel).
* `alarm_ways` - (Optional, Set) Alert Channels: 1: Email, 2: SMS, 3: WeChat, 4: Voice, 5: WeChat Enterprise, 6: Http, 7: WeChat Enterprise Group, 8: Lark Group, 9: DingTalk Group, 10: Slack Group, 11: Teams Group (Default: Email), Only one channel can be selected.
* `notification_fatigue` - (Optional, List) Alarm notification fatigue configuration.
* `web_hooks` - (Optional, List) List of webhook addresses for corporate WeChat groups, Feishu groups, DingTalk groups, Slack groups, and Teams groups.

The `alarm_rule_detail` object supports the following:

* `data_backfill_or_rerun_time_out_ext_info` - (Optional, List) Detailed configuration of re-running and re-recording instance timeout.
* `data_backfill_or_rerun_trigger` - (Optional, Int) Re-recording trigger timing: 1 - Triggered by the first failure; 2 - Triggered by completion of all retries.
* `project_instance_statistics_alarm_info_list` - (Optional, List) Project fluctuation alarm configuration details.
* `reconciliation_ext_info` - (Optional, List) Offline integrated reconciliation alarm configuration information.
* `time_out_ext_info` - (Optional, List) Periodic instance timeout configuration details.
* `trigger` - (Optional, Int) Failure trigger timing: 1 - Triggered on first failure; 2 -- Triggered when all retries complete (default).

The `data_backfill_or_rerun_time_out_ext_info` object of `alarm_rule_detail` supports the following:

* `hour` - (Optional, Int) Specify the timeout value in hours. The default value is 0.
* `min` - (Optional, Int) The timeout value is specified in minutes. The default value is 1.
* `rule_type` - (Optional, Int) Timeout alarm configuration: 1. Estimated running time exceeded, 2. Estimated completion time exceeded, 3. Estimated waiting time for scheduling exceeded, 4. Estimated completion within the period but not completed.
* `schedule_time_zone` - (Optional, String) The time zone configuration corresponding to the timeout period, such as UTC+7, the default is UTC+8.
* `type` - (Optional, Int) Timeout value configuration type: 1-Specified value; 2-Average value.

The `notification_fatigue` object of `alarm_groups` supports the following:

* `notify_count` - (Optional, Int) Number of alarms.
* `notify_interval` - (Optional, Int) Alarm interval, in minutes.
* `quiet_intervals` - (Optional, List) Do not disturb time, for example, the example value [{DaysOfWeek: [1, 2], StartTime: "00:00:00", EndTime: "09:00:00"}] means do not disturb from 00:00 to 09:00 every Monday and Tuesday.

The `project_instance_statistics_alarm_info_list` object of `alarm_rule_detail` supports the following:

* `alarm_type` - (Required, String) Alarm type: projectFailureInstanceUpwardFluctuationAlarm: Failure instance upward fluctuation alarm; projectSuccessInstanceDownwardFluctuationAlarm: Success instance downward fluctuation alarm.
* `instance_count` - (Optional, Int) The cumulative number of instances on the day; the downward fluctuation of the number of failed instances on the day.
* `instance_threshold_count_percent` - (Optional, Int) The alarm threshold for the proportion of instance successes fluctuating downwards; the alarm threshold for the proportion of instance failures fluctuating upwards.
* `instance_threshold_count` - (Optional, Int) The cumulative instance number fluctuation threshold.
* `is_cumulant` - (Optional, Bool) Whether to calculate cumulatively, false: continuous, true: cumulative.
* `stabilize_statistics_cycle` - (Optional, Int) Stability statistics period (number of anti-shake configuration statistics periods).
* `stabilize_threshold` - (Optional, Int) Stability threshold (number of statistical cycles for anti-shake configuration).

The `quiet_intervals` object of `notification_fatigue` supports the following:

* `days_of_week` - (Optional, Set) According to the ISO standard, 1 represents Monday and 7 represents Sunday.
* `end_time` - (Optional, String) End time, with precision of hours, minutes, and seconds, in the format of HH:mm:ss.
* `start_time` - (Optional, String) Start time, with precision of hours, minutes, and seconds, in the format of HH:mm:ss.

The `reconciliation_ext_info` object of `alarm_rule_detail` supports the following:

* `hour` - (Optional, Int) Reconciliation task timeout threshold: hours, default is 0.
* `min` - (Optional, Int) Reconciliation task timeout threshold: minutes, default is 1.
* `mismatch_count` - (Optional, Int) Reconciliation inconsistency threshold, RuleType=reconciliationMismatch. This field needs to be configured and has no default value.
* `rule_type` - (Optional, String) Offline alarm rule types: reconciliationFailure: Offline reconciliation failure alarm; reconciliationOvertime: Offline reconciliation task timeout alarm (timeout must be configured); reconciliationMismatch: Offline reconciliation mismatch alarm (mismatch threshold must be configured).

The `time_out_ext_info` object of `alarm_rule_detail` supports the following:

* `hour` - (Optional, Int) Specify the timeout value in hours. The default value is 0.
* `min` - (Optional, Int) The timeout value is specified in minutes. The default value is 1.
* `rule_type` - (Optional, Int) Timeout alarm configuration: 1. Estimated running time exceeded, 2. Estimated completion time exceeded, 3. Estimated waiting time for scheduling exceeded, 4. Estimated completion within the period but not completed.
* `schedule_time_zone` - (Optional, String) The time zone configuration corresponding to the timeout period, such as UTC+7, the default is UTC+8.
* `type` - (Optional, Int) Timeout value configuration type: 1-Specified value; 2-Average value.

The `web_hooks` object of `alarm_groups` supports the following:

* `alarm_way` - (Optional, String) Alert channel value: 7. Enterprise WeChat group, 8. Feishu group, 9. DingTalk group, 10. Slack group, 11. Teams group.
* `web_hooks` - (Optional, Set) List of webhook addresses for the alarm group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

wedata ops alarm rule can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule projectId#askId
```


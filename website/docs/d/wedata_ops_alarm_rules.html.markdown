---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_alarm_rules"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_alarm_rules"
description: |-
  Use this data source to query detailed information of wedata ops alarm rules
---

# tencentcloud_wedata_ops_alarm_rules

Use this data source to query detailed information of wedata ops alarm rules

## Example Usage

```hcl
data "tencentcloud_wedata_ops_alarm_rules" "wedata_ops_alarm_rules" {
  project_id = "1859317240494305280"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project id.
* `alarm_level` - (Optional, Int) Alarm level: 1. Normal, 2. Major, 3. Urgent.
* `alarm_recipient_id` - (Optional, String) Query the alarm rules configured for the corresponding alarm recipient.
* `alarm_type` - (Optional, String) Alarm Rule Monitoring Types: failure: failure alarm; overtime: timeout alarm; success: success alarm; backTrackingOrRerunSuccess: backTrackingOrRerunSuccess: backTrackingOrRerunFailure: backTrackingOrRerunFailure. Project Fluctuation Alarms: projectFailureInstanceUpwardFluctuationAlarm: alarm if the upward fluctuation rate of failed instances exceeds the threshold. projectSuccessInstanceDownwardFluctuationAlarm: alarm if the downward fluctuation rate of successful instances exceeds the threshold. Offline Integration Task Reconciliation Alarms: reconciliationFailure: offline reconciliation task failure alarm; reconciliationOvertime: offline reconciliation task timeout alarm; reconciliationMismatch: alarm if the number of inconsistent entries in a data reconciliation task exceeds the threshold. Example value: ["failure"].
* `create_time_from` - (Optional, String) The start time of the alarm rule creation time range, in the format of 2025-08-17 00:00:00.
* `create_time_to` - (Optional, String) The end time of the alarm rule creation time range, in the format of "2025-08-26 23:59:59".
* `create_user_uin` - (Optional, String) Alarm rule creator filtering.
* `keyword` - (Optional, String) Query the corresponding alarm rule based on the alarm rule ID/rule name.
* `monitor_object_type` - (Optional, Int) Monitoring object type, Task dimension monitoring: can be configured according to task/workflow/project: 1.Task, 2.Workflow, 3.Project (default is 1.Task) Project dimension monitoring: Project overall task fluctuation alarm, 7: Project fluctuation monitoring alarm.
* `result_output_file` - (Optional, String) Used to save results.
* `task_id` - (Optional, String) Query alarm rules based on task ID.
* `update_time_from` - (Optional, String) Last updated time filter alarm rules, format such as "2025-08-26 00:00:00".
* `update_time_to` - (Optional, String) Last updated time filter alarm rule format such as: "2025-08-26 23:59:59".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Alarm information response.



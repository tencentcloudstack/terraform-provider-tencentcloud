---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_baseline"
sidebar_current: "docs-tencentcloud-resource-wedata_baseline"
description: |-
  Provides a resource to create a wedata baseline
---

# tencentcloud_wedata_baseline

Provides a resource to create a wedata baseline

## Example Usage

```hcl
resource "tencentcloud_wedata_baseline" "example" {
  project_id     = "1927766435649077248"
  baseline_name  = "tf_example"
  baseline_type  = "D"
  create_uin     = "100028439226"
  create_name    = "tf_user"
  in_charge_uin  = "tf_user"
  in_charge_name = "100028439226"
  promise_tasks {
    project_id          = "1927766435649077248"
    task_name           = "tf_demo_task"
    task_id             = "20231030145334153"
    task_cycle          = "D"
    workflow_name       = "dataflow_mpp"
    workflow_id         = "e4dafb2e-76eb-11ee-bfeb-b8cef68a6637"
    task_in_charge_name = ";tf_user;"
  }
  promise_time   = "00:00:00"
  warning_margin = 30
  is_new_alarm   = true
  baseline_create_alarm_rule_request {
    alarm_types = [
      "baseLineBroken",
      "baseLineWarning",
      "baseLineTaskFailure"
    ]
    alarm_level = 2
    alarm_ways = [
      "email",
      "sms"
    ]
    alarm_recipient_type = 1
    alarm_recipients = [
      "tf_user"
    ]
    alarm_recipient_ids = [
      "100028439226"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `baseline_name` - (Required, String) Baseline Name.
* `baseline_type` - (Required, String) D or H; representing daily baseline and hourly baseline respectively.
* `create_name` - (Required, String) Creator Name.
* `create_uin` - (Required, String) Creator ID.
* `in_charge_name` - (Required, String) Baseline Owner Name.
* `in_charge_uin` - (Required, String) Baseline Owner ID.
* `is_new_alarm` - (Required, Bool) Is it a newly created alarm rule.
* `project_id` - (Required, String) Project ID.
* `promise_tasks` - (Required, List) Promise Tasks.
* `promise_time` - (Required, String) Service Assurance Time.
* `warning_margin` - (Required, Int) Warning Margin in minutes.
* `alarm_rule_dto` - (Optional, List) Existing Alarm Rule Information.
* `baseline_create_alarm_rule_request` - (Optional, List) Description of the New Alarm Rule.

The `alarm_rule_dto` object supports the following:

* `alarm_level_type` - (Optional, String) Important;Urgent;Normal.
* `alarm_rule_id` - (Optional, String) Alarm Rule ID.

The `baseline_create_alarm_rule_request` object supports the following:

* `alarm_level` - (Optional, Int) Alarm Level, 1. Normal, 2. Important, 3. Urgent (default is 1. Normal)Note: This field may return null, indicating no valid value.
* `alarm_recipient_ids` - (Optional, Set) Alarm Recipient IDsNote: This field may return null, indicating no valid value.
* `alarm_recipient_type` - (Optional, Int) Alarm Recipient Type: 1. Specified Personnel, 2. Task Owner, 3. Duty Roster (default is 1. Specified Personnel)Note: This field may return null, indicating no valid value.
* `alarm_recipients` - (Optional, Set) Alarm RecipientsNote: This field may return null, indicating no valid value.
* `alarm_types` - (Optional, Set) Alarm Types, 1. Failure Alarm, 2. Timeout Alarm, 3. Success Alarm, 4. Baseline Violation, 5. Baseline Warning, 6. Baseline Task Failure (default is 1. Failure Alarm)Note: This field may return null, indicating no valid value.
* `alarm_ways` - (Optional, Set) Alarm Methods, 1. Email, 2. SMS, 3. WeChat, 4. Voice, 5. Enterprise WeChat, 6. HTTP, 7. Enterprise WeChat Group; Alarm method code list (default is 1. Email)Note: This field may return null, indicating no valid value.
* `creator_id` - (Optional, String) Creator NameNote: This field may return null, indicating no valid value.
* `creator` - (Optional, String) Creator UINNote: This field may return null, indicating no valid value.
* `ext_info` - (Optional, String) Extended Information, 1. Estimated Runtime (default), 2. Estimated Completion Time, 3. Estimated Scheduling Time, 4. Incomplete within the Cycle; Value Types: 1. Specified Value, 2. Historical Average (default is 1. Specified Value)Note: This field may return null, indicating no valid value.
* `monitor_object_ids` - (Optional, Set) Monitoring ObjectsNote: This field may return null, indicating no valid value.
* `monitor_type` - (Optional, Int) Monitoring Type, 1. Task, 2. Workflow, 3. Project, 4. Baseline (default is 1. Task)Note: This field may return null, indicating no valid value.
* `project_id` - (Optional, String) Project NameNote: This field may return null, indicating no valid value.
* `rule_name` - (Optional, String) Rule NameNote: This field may return null, indicating no valid value.

The `promise_tasks` object supports the following:

* `project_id` - (Optional, String) Project ID.
* `task_cycle` - (Optional, String) Task Scheduling Cycle.
* `task_id` - (Optional, String) Task ID.
* `task_in_charge_name` - (Optional, String) Task Owner Name.
* `task_in_charge_uin` - (Optional, String) Task Owner ID.
* `task_name` - (Optional, String) Task Name.
* `workflow_id` - (Optional, String) Workflow ID.
* `workflow_name` - (Optional, String) Workflow Name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `baseline_id` - Baseline ID.


## Import

wedata baseline can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_baseline.example 1927766435649077248#2
```


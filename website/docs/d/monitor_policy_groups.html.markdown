---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_policy_groups"
sidebar_current: "docs-tencentcloud-datasource-monitor_policy_groups"
description: |-
  Use this data source to query monitor policy groups (There is a lot of data and it is recommended to output to a file)
---

# tencentcloud_monitor_policy_groups

Use this data source to query monitor policy groups (There is a lot of data and it is recommended to output to a file)

## Example Usage

```hcl
data "tencentcloud_monitor_policy_groups" "groups" {
  policy_view_names = ["REDIS-CLUSTER", "cvm_device"]
}

data "tencentcloud_monitor_policy_groups" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Policy group name for query.
* `policy_view_names` - (Optional) The policy view for query.
* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list policy groups. Each element contains the following attributes:
  * `can_set_default` - Whether it can be set as the default policy.
  * `conditions` - A list of threshold rules. Each element contains the following attributes:
    * `alarm_notify_period` - Alarm sending cycle per second. `<0` does not fire, `0` only fires once, and `>0` fires every triggerTime second.
    * `alarm_notify_type` - Alarm sending convergence type. `0` continuous alarm, `1` index alarm.
    * `calc_type` - Compare type, `1` means more than, `2`  means greater than or equal, `3` means less than, `4` means less than or equal to, `5` means equal, `6` means not equal, `7` means days rose, `8` means days fell, `9` means weeks rose, `10` means weeks fell, `11` means period rise, `12` means period fell.
    * `calc_value` - Threshold value.
    * `continue_time` - How long does the triggering rule last (per second).
    * `metric_id` - The ID of this metric.
    * `metric_show_name` - The name of this metric.
    * `metric_unit` - The unit of this metric.
    * `period` - Data aggregation cycle (unit second).
    * `rule_id` - Threshold rule ID.
  * `event_conditions` - A list of event rules. Each element contains the following attributes:
    * `alarm_notify_period` - Alarm sending cycle per second. `<0` does not fire, `0` only fires once, and `>0` fires every triggerTime second.
    * `alarm_notify_type` - Alarm sending convergence type. `0` continuous alarm, `1` index alarm.
    * `event_id` - The ID of this event metric.
    * `event_show_name` - The name of this event metric.
    * `rule_id` - Threshold rule ID.
  * `group_id` - The policy group id.
  * `group_name` - The policy group name.
  * `insert_time` - The policy group create timestamp.
  * `is_default` - If is default policy group or not, `0` represents the non-default policy, and `1` represents the default policy.
  * `is_open` - Whether open or not.
  * `last_edit_uin` - Recently edited user uin.
  * `no_shielded_sum` - Number of unmasked instances of policy group bindings.
  * `parent_group_id` - Parent policy group ID.
  * `policy_view_name` - The policy group view name.
  * `project_id` - The project ID to which the policy group belongs.
  * `receivers` - A list of receivers. Each element contains the following attributes:
    * `end_time` - End of alarm period. Meaning with `start_time`.
    * `need_send_notice` - Do need a telephone alarm contact prompt.You don't need 0, you need 1.
    * `notify_way` - Method of warning notification.Optional `CALL`,`EMAIL`,`SITE`,`SMS`,`WECHAT`.
    * `person_interval` - Telephone warning to individual interval (seconds).
    * `receive_language` - Alert sending language.
    * `receiver_group_list` - Alarm receive group ID list.
    * `receiver_type` - Receive type. Optional 'group' or 'user'.
    * `receiver_user_list` - Alarm receiver ID list.
    * `recover_notify` - Restore notification mode. Optional "SMS".
    * `round_interval` - Telephone alarm interval per round (seconds).
    * `round_number` - Telephone alarm number.
    * `send_for` - Telephone warning time.Option "OCCUR", "RECOVER".
    * `start_time` - Alarm period start time.Range [0,86399], which removes the date after it is converted to Beijing time as a Unix timestamp, for example 7200 means '10:0:0'.
    * `uid_list` - The phone alerts the receiver uid.
  * `remark` - Policy group remarks.
  * `update_time` - The policy group update timestamp.
  * `use_sum` - Number of instances of policy group bindings.



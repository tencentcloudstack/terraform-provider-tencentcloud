---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_policy_group"
sidebar_current: "docs-tencentcloud-resource-monitor_policy_group"
description: |-
  Provides a policy group resource for monitor.
---

# tencentcloud_monitor_policy_group

Provides a policy group resource for monitor.

## Example Usage

```hcl
resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "nice_group"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  is_union_rule    = 1
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
  conditions {
    metric_id           = 30
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 2
    calc_value          = 30
    calc_period         = 300
    continue_period     = 2
  }
  event_conditions {
    event_id            = 39
    alarm_notify_type   = 0
    alarm_notify_period = 300
  }
  event_conditions {
    event_id            = 40
    alarm_notify_type   = 0
    alarm_notify_period = 300
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required) Policy group name, length should between 1 and 20.
* `policy_view_name` - (Required, ForceNew) Policy view name, eg:`cvm_device`,`BANDWIDTHPACKAGE`, refer to `data.tencentcloud_monitor_policy_conditions(policy_view_name)`.
* `remark` - (Required, ForceNew) Policy group's remark information.
* `conditions` - (Optional) A list of threshold rules. Each element contains the following attributes:
* `event_conditions` - (Optional) A list of event rules. Each element contains the following attributes:
* `is_union_rule` - (Optional) The and or relation of indicator alarm rule. Valid values: `0`, `1`. `0` represents or rule (if any rule is met, the alarm will be raised), `1` represents and rule (if all rules are met, the alarm will be raised).The default is 0.
* `project_id` - (Optional, ForceNew) The project id to which the policy group belongs, default is `0`.

The `conditions` object supports the following:

* `alarm_notify_period` - (Required) Alarm sending cycle per second. <0 does not fire, `0` only fires once, and >0 fires every triggerTime second.
* `alarm_notify_type` - (Required) Alarm sending convergence type. `0` continuous alarm, `1` index alarm.
* `metric_id` - (Required) Id of the metric, refer to `data.tencentcloud_monitor_policy_conditions(metric_id)`.
* `calc_period` - (Optional) Data aggregation cycle (unit of second), if the metric has a default value can not be filled, refer to `data.tencentcloud_monitor_policy_conditions(period_keys)`.
* `calc_type` - (Optional) Compare type. Valid value ranges: [1~12]. `1` means more than, `2` means greater than or equal, `3` means less than, `4` means less than or equal to, `5` means equal, `6` means not equal, `7` means days rose, `8` means days fell, `9` means weeks rose, `10` means weeks fell, `11` means period rise, `12` means period fell, refer to `data.tencentcloud_monitor_policy_conditions(calc_type_keys)`.
* `calc_value` - (Optional) Threshold value, refer to `data.tencentcloud_monitor_policy_conditions(calc_value_*)`.
* `continue_period` - (Optional) The rule triggers an alert that lasts for several detection cycles, refer to `data.tencentcloud_monitor_policy_conditions(period_num_keys)`.

The `event_conditions` object supports the following:

* `alarm_notify_period` - (Required) Alarm sending cycle per second. <0 does not fire, `0` only fires once, and >0 fires every triggerTime second.
* `alarm_notify_type` - (Required) Alarm sending convergence type. `0` continuous alarm, `1` index alarm.
* `event_id` - (Required) The ID of this event metric, refer to `data.tencentcloud_monitor_policy_conditions(event_id).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `binding_objects` - A list binding objects(list only those in the `provider.region`). Each element contains the following attributes:
  * `dimensions_json` - Represents a collection of dimensions of an object instance, json format.
  * `is_shielded` - Whether the object is shielded or not, 0 means unshielded and 1 means shielded.
  * `region` - The region where the object is located.
  * `unique_id` - Object unique id.
* `dimension_group` - A list of dimensions for this policy group.
* `last_edit_uin` - Recently edited user uin.
* `receivers` - A list of receivers. Each element contains the following attributes:
  * `end_time` - End of alarm period. Meaning with `start_time`.
  * `need_send_notice` - Do need a telephone alarm contact prompt. You don't need `0`, you need `1`.
  * `notify_way` - Method of warning notification. Valid values: "SMS", "SITE", "EMAIL", "CALL", "WECHAT".
  * `person_interval` - Telephone warning to individual interval (seconds).
  * `receive_language` - Alert sending language.
  * `receiver_group_list` - Alarm receive group ID list.
  * `receiver_type` - Receive type. Valid values: group, user. 'group' (receiving group) or 'user' (receiver).
  * `receiver_user_list` - Alarm receiver id list.
  * `recover_notify` - Restore notification mode. Optional "SMS".
  * `round_interval` - Telephone alarm interval per round (seconds).
  * `round_number` - Telephone alarm number.
  * `send_for` - Telephone warning time. Valid values: "OCCUR","RECOVER".
  * `start_time` - Alarm period start time.Range [0,86400], which removes the date after it is converted to Beijing time as a Unix timestamp, for example 7200 means '10:0:0'.
  * `uid_list` - The phone alerts the receiver uid.
* `support_regions` - Support regions this policy group.
* `update_time` - The policy group update time.


## Import

Policy group instance can be imported, e.g.

```
$ terraform import tencentcloud_monitor_policy_group.group group-id
```


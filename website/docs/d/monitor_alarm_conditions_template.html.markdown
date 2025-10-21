---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_conditions_template"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_conditions_template"
description: |-
  Use this data source to query detailed information of monitor alarm_conditions_template
---

# tencentcloud_monitor_alarm_conditions_template

Use this data source to query detailed information of monitor alarm_conditions_template

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_conditions_template" "alarm_conditions_template" {
  module             = "monitor"
  view_name          = "cvm_device"
  group_name         = "keep-template"
  group_id           = "7803070"
  update_time_order  = "desc=descending"
  policy_count_order = "asc=ascending"
}
```

## Argument Reference

The following arguments are supported:

* `module` - (Required, String) Fixed value, as&amp;amp;#39; monitor &amp;amp;#39;.
* `group_id` - (Optional, String) Filter queries based on trigger condition template ID.
* `group_name` - (Optional, String) Filter queries based on trigger condition template names.
* `policy_count_order` - (Optional, String) Specify the sorting method based on the number of binding policies, asc=ascending, desc=descending.
* `result_output_file` - (Optional, String) Used to save results.
* `update_time_order` - (Optional, String) Specify the sorting method by update time, asc=ascending, desc=descending.
* `view_name` - (Optional, String) View name, composed of [DescribeAllNamespaces]( https://cloud.tencent.com/document/product/248/48683 )Obtain. For cloud product monitoring, retrieve the QceNamespacesNew. N.ID parameter from the interface, such as cvm_ Device.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_group_list` - Template List.
  * `conditions` - Indicator alarm rules.
    * `alarm_notify_period` - Alarm notification frequency.
    * `alarm_notify_type` - Predefined repeated notification strategy (0- alarm only once, 1- exponential alarm, 2- connection alarm).
    * `calc_type` - Detection method.
    * `calc_value` - Detection value.
    * `continue_time` - Duration in seconds.
    * `is_advanced` - Whether it is an advanced indicator, 0: No; 1: Yes.
    * `is_open` - Whether to activate advanced indicators, 0: No; 1: Yes.
    * `metric_display_name` - Indicator display name (external).
    * `metric_id` - Indicator ID.
    * `period` - Cycle.
    * `product_id` - Product ID.
    * `rule_id` - Rule ID.
    * `unit` - Indicator unit.
  * `event_conditions` - Event alarm rules.
    * `alarm_notify_period` - Alarm notification frequency.
    * `alarm_notify_type` - Predefined repeated notification strategy (0- alarm only once, 1- exponential alarm, 2- connection alarm).
    * `event_display_name` - Event Display Name (External).
    * `event_id` - Event ID.
    * `rule_id` - Rule ID.
  * `group_id` - Template Policy Group ID.
  * `group_name` - Template Policy Group Name.
  * `insert_time` - Creation time.
  * `is_union_rule` - Is it a relationship with.
  * `last_edit_uin` - Last modified by UIN.
  * `policy_groups` - Associate Alert Policy Group.
    * `can_set_default` - Can it be set as the default alarm strategy.
    * `enable` - Alarm Policy Enable Status.
    * `group_id` - Alarm Policy Group ID.
    * `group_name` - Alarm Policy Group Name.
    * `insert_time` - Creation time.
    * `is_default` - Is it the default alarm policy.
    * `is_union_rule` - Is it a relationship rule with.
    * `last_edit_uin` - Last modified by UIN.
    * `no_shielded_instance_count` - Number of unshielded instances.
    * `parent_group_id` - Parent Policy Group ID.
    * `project_id` - Project ID.
    * `receiver_infos` - Alarm receiving object information.
      * `end_time` - Effective period end time.
      * `need_send_notice` - Do you need to send a notification.
      * `notify_way` - Alarm reception channel.
      * `person_interval` - Telephone alarm to personal interval (seconds).
      * `receiver_group_list` - Message receiving group list.
      * `receiver_type` - Receiver type.
      * `receiver_user_list` - Recipient list. List of recipient IDs queried through the platform interface.
      * `recover_notify` - Alarm recovery notification method.
      * `round_interval` - Telephone alarm interval per round (seconds).
      * `round_number` - Number of phone alarm rounds.
      * `send_for` - Timing of telephone alarm notification. Optional OCCUR (notification during alarm), RECOVER (notification during recovery).
      * `start_time` - Effective period start time.
      * `uid_list` - Telephone alarm receiver uid.
    * `remark` - Remarks.
    * `total_instance_count` - Total number of bound instances.
    * `update_time` - Modification time.
    * `view_name` - View.
  * `remark` - Remarks.
  * `update_time` - Update time.
  * `view_name` - View.



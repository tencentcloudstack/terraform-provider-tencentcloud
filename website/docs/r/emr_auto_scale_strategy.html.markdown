---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_auto_scale_strategy"
sidebar_current: "docs-tencentcloud-resource-emr_auto_scale_strategy"
description: |-
  Provides a resource to create a emr emr_auto_scale_strategy
---

# tencentcloud_emr_auto_scale_strategy

Provides a resource to create a emr emr_auto_scale_strategy

## Example Usage

```hcl
resource "tencentcloud_emr_auto_scale_strategy" "emr_auto_scale_strategy" {
  instance_id   = "emr-rzrochgp"
  strategy_type = 2
  time_auto_scale_strategy {
    strategy_name    = "tf-test1"
    interval_time    = 100
    scale_action     = 1
    scale_num        = 1
    strategy_status  = 1
    retry_valid_time = 60
    repeat_strategy {
      repeat_type = "DAY"
      day_repeat {
        execute_at_time_of_day = "16:30:00"
        step                   = 1
      }
      expire = "2026-02-20 23:59:59"
    }
    grace_down_flag = false
    tags {
      tag_key   = "createBy"
      tag_value = "terraform"
    }
    config_group_assigned = "{\"HDFS-2.8.5\":-1,\"YARN-2.8.5\":-1}"
    measure_method        = "INSTANCE"
    terminate_policy      = "DEFAULT"
    soft_deploy_info      = [1, 2]
    service_node_info     = [7]
    priority              = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `strategy_type` - (Required, Int) 1 means expansion and contraction according to load rules, 2 means expansion and contraction according to time rules. Must be filled in and match the following rule policy.
* `load_auto_scale_strategy` - (Optional, List) Expansion rules based on load.
* `time_auto_scale_strategy` - (Optional, List) Rules for scaling up and down over time.

The `conditions` object of `load_metrics` supports the following:

* `compare_method` - (Required, Int) Conditional comparison method, 1 means greater than, 2 means less than, 3 means greater than or equal to, 4 means less than or equal to.
* `threshold` - (Optional, Float64) Conditional threshold.

The `day_repeat` object of `repeat_strategy` supports the following:

* `execute_at_time_of_day` - (Required, String) Repeat the specific time when the task is executed, such as "01:02:00".
* `step` - (Required, Int) Executed every Step day.

The `load_auto_scale_strategy` object supports the following:

* `calm_down_time` - (Optional, Int) Cooling time for rules to take effect.
* `config_group_assigned` - (Optional, String) Default configuration group.
* `grace_down_flag` - (Optional, Bool) Elegant shrink switch.
* `grace_down_time` - (Optional, Int) Graceful downsizing waiting time.
* `load_metrics_conditions` - (Optional, List) Multiple indicator trigger conditions.
* `measure_method` - (Optional, String) Expansion resource calculation methods, "DEFAULT", "INSTANCE", "CPU", "MEMORYGB".
"DEFAULT" means the default mode, which has the same meaning as "INSTANCE".
"INSTANCE" means calculation based on nodes, the default method.
"CPU" means calculated based on the number of cores of the machine.
"MEMORYGB" means calculated based on the number of machine memory.
* `period_valid` - (Optional, String) Effective time for the rule to take effect.
* `priority` - (Optional, Int) Rule priority, invalid when added, defaults to auto-increment.
* `process_method` - (Optional, Int) Indicator processing method, 1 represents MAX, 2 represents MIN, and 3 represents AVG.
* `scale_action` - (Optional, Int) Expansion and contraction actions, 1 means expansion, 2 means shrinkage.
* `scale_num` - (Optional, Int) The amount of expansion and contraction each time the rule takes effect.
* `strategy_id` - (Optional, Int) Rule ID.
* `strategy_name` - (Optional, String) Rule name.
* `strategy_status` - (Optional, Int) Rule status, 1 means enabled, 3 means disabled.
* `tags` - (Optional, List) Binding tag list.
* `yarn_node_label` - (Optional, String) Rule expansion specifies yarn node label.

The `load_metrics_conditions` object of `load_auto_scale_strategy` supports the following:

* `load_metrics` - (Optional, List) Trigger rule conditions.

The `load_metrics` object of `load_metrics_conditions` supports the following:

* `conditions` - (Optional, List) Trigger condition.
* `load_metrics` - (Optional, String) Expansion and contraction load indicators.
* `metric_id` - (Optional, Int) Rule metadata record ID.
* `statistic_period` - (Optional, Int) The regular statistical period provides 1min, 3min, and 5min.
* `trigger_threshold` - (Optional, Int) The number of triggers. When the number of consecutive triggers exceeds TriggerThreshold, the expansion and contraction will begin.

The `month_repeat` object of `repeat_strategy` supports the following:

* `days_of_month_range` - (Required, Set) The description of the day period in each month, the length can only be 2, for example, [2,10] means the 2-10th of each month.
* `execute_at_time_of_day` - (Required, String) Repeat the specific time when the task is executed, such as "01:02:00".

The `not_repeat` object of `repeat_strategy` supports the following:

* `execute_at` - (Required, String) The specific and complete time of the task execution, the format is "2020-07-13 00:00:00".

The `repeat_strategy` object of `time_auto_scale_strategy` supports the following:

* `repeat_type` - (Required, String) The value range is "DAY", "DOW", "DOM", and "NONE", which respectively represent daily repetition, weekly repetition, monthly repetition and one-time execution. Required.
* `day_repeat` - (Optional, List) Repeat rules by day, valid when RepeatType is "DAY".
* `expire` - (Optional, String) Rule expiration time. After this time, the rule will automatically be placed in a suspended state, in the form of "2020-07-23 00:00:00". Required.
* `month_repeat` - (Optional, List) Repeat rules by month, valid when RepeatType is "DOM".
* `not_repeat` - (Optional, List) Execute the rule once, effective when RepeatType is "NONE".
* `week_repeat` - (Optional, List) Repeat rules by week, valid when RepeatType is "DOW".

The `tags` object of `load_auto_scale_strategy` supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

The `tags` object of `time_auto_scale_strategy` supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

The `time_auto_scale_strategy` object supports the following:

* `interval_time` - (Required, Int) The cooling time after the policy is triggered. During this period, elastic expansion and contraction will not be triggered.
* `priority` - (Required, Int) Rule priority, the smaller it is, the higher it is.
* `repeat_strategy` - (Required, List) Time expansion and contraction repetition strategy.
* `retry_valid_time` - (Required, Int) When multiple rules are triggered at the same time and some of them are not actually executed, retries will be made within this time range.
* `scale_action` - (Required, Int) Expansion and contraction actions, 1 means expansion, 2 means shrinkage.
* `scale_num` - (Required, Int) The number of expansions and contractions.
* `strategy_name` - (Required, String) Policy name, unique within the cluster.
* `strategy_status` - (Required, Int) Rule status, 1 means valid, 2 means invalid, and 3 means suspended. Required.
* `compensate_flag` - (Optional, Int) Compensation expansion, 0 means not enabled, 1 means enabled.
* `config_group_assigned` - (Optional, String) Default configuration group.
* `grace_down_flag` - (Optional, Bool) Elegant shrink switch.
* `grace_down_time` - (Optional, Int) Graceful downsizing waiting time.
* `group_id` - (Optional, Int) scaling group id.
* `max_use` - (Optional, Int) Maximum usage time, seconds, minimum 1 hour, maximum 24 hours.
* `measure_method` - (Optional, String) Expansion resource calculation methods, "DEFAULT", "INSTANCE", "CPU", "MEMORYGB".
"DEFAULT" means the default mode, which has the same meaning as "INSTANCE".
"INSTANCE" means calculation based on nodes, the default method.
"CPU" means calculated based on the number of cores of the machine.
"MEMORYGB" means calculated based on the number of machine memory.
* `service_node_info` - (Optional, Set) Start process list.
* `soft_deploy_info` - (Optional, Set) Node deployment service list. Only fill in HDFS and YARN for deployment services. [Mapping relationship table corresponding to component names](https://cloud.tencent.com/document/product/589/98760).
* `tags` - (Optional, List) Binding tag list.
* `terminate_policy` - (Optional, String) Destruction strategy, "DEFAULT", the default destruction strategy, shrinkage is triggered by shrinkage rules, "TIMING" means scheduled destruction.

The `week_repeat` object of `repeat_strategy` supports the following:

* `days_of_week` - (Required, Set) The numerical description of the days of the week, for example, [1,3,4] means Monday, Wednesday, and Thursday every week.
* `execute_at_time_of_day` - (Required, String) Repeat the specific time when the task is executed, such as "01:02:00".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

emr emr_auto_scale_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_emr_auto_scale_strategy.emr_auto_scale_strategy emr_auto_scale_strategy_id
```


---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scaling_policy"
sidebar_current: "docs-tencentcloud-resource-as_scaling_policy"
description: |-
  Provides a resource for an AS (Auto scaling) policy.
---

# tencentcloud_as_scaling_policy

Provides a resource for an AS (Auto scaling) policy.

## Example Usage

```hcl
resource "tencentcloud_as_scaling_policy" "scaling_policy" {
  scaling_group_id    = "asg-n32ymck2"
  policy_name         = "tf-as-scaling-policy"
  adjustment_type     = "EXACT_CAPACITY"
  adjustment_value    = 0
  comparison_operator = "GREATER_THAN"
  metric_name         = "CPU_UTILIZATION"
  threshold           = 80
  period              = 300
  continuous_time     = 10
  statistic           = "AVERAGE"
  cooldown            = 360
}
```

## Argument Reference

The following arguments are supported:

* `adjustment_type` - (Required) Specifies whether the adjustment is an absolute number or a percentage of the current capacity. Valid values: `CHANGE_IN_CAPACITY`, `EXACT_CAPACITY` and `PERCENT_CHANGE_IN_CAPACITY`.
* `adjustment_value` - (Required) Define the number of instances by which to scale.For `CHANGE_IN_CAPACITY` type or PERCENT_CHANGE_IN_CAPACITY, a positive increment adds to the current capacity and a negative value removes from the current capacity. For `EXACT_CAPACITY` type, it defines an absolute number of the existing Auto Scaling group size.
* `comparison_operator` - (Required) Comparison operator. Valid values: `GREATER_THAN`, `GREATER_THAN_OR_EQUAL_TO`, `LESS_THAN`, `LESS_THAN_OR_EQUAL_TO`, `EQUAL_TO` and `NOT_EQUAL_TO`.
* `continuous_time` - (Required) Retry times. Valid value ranges: (1~10).
* `metric_name` - (Required) Name of an indicator. Valid values: `CPU_UTILIZATION`, `MEM_UTILIZATION`, `LAN_TRAFFIC_OUT`, `LAN_TRAFFIC_IN`, `WAN_TRAFFIC_OUT` and `WAN_TRAFFIC_IN`.
* `period` - (Required) Time period in second. Valid values: `60` and `300`.
* `policy_name` - (Required) Name of a policy used to define a reaction when an alarm is triggered.
* `scaling_group_id` - (Required, ForceNew) ID of a scaling group.
* `threshold` - (Required) Alarm threshold.
* `cooldown` - (Optional) Cooldwon time in second. Default is `30`0.
* `notification_user_group_ids` - (Optional) An ID group of users to be notified when an alarm is triggered.
* `statistic` - (Optional) Statistic types. Valid values: `AVERAGE`, `MAXIMUM` and `MINIMUM`. Default is `AVERAGE`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.




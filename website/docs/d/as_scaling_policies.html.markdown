---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scaling_policies"
sidebar_current: "docs-tencentcloud-datasource-as_scaling_policies"
description: |-
  Use this data source to query detailed information of scaling policy.
---

# tencentcloud_as_scaling_policies

Use this data source to query detailed information of scaling policy.

## Example Usage

```hcl
data "tencentcloud_as_scaling_policies" "as_scaling_policies" {
  scaling_policy_id  = "asg-mvyghxu7"
  result_output_file = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `policy_name` - (Optional) Scaling policy name.
* `result_output_file` - (Optional) Used to save results.
* `scaling_group_id` - (Optional) Scaling group ID.
* `scaling_policy_id` - (Optional) Scaling policy ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_policy_list` - A list of scaling policy. Each element contains the following attributes:
  * `adjustment_type` - Adjustment type of the scaling rule.
  * `adjustment_value` - Adjustment value of the scaling rule.
  * `comparison_operator` - Comparison operator.
  * `continuous_time` - Retry times.
  * `cooldown` - Cooldown time of the scaling rule.
  * `metric_name` - Name of an indicator.
  * `notification_user_group_ids` - Users need to be notified when an alarm is triggered.
  * `period` - Time period in second.
  * `policy_name` - Scaling policy name.
  * `scaling_group_id` - Scaling policy ID.
  * `statistic` - Statistic types.
  * `threshold` - Alarm threshold.



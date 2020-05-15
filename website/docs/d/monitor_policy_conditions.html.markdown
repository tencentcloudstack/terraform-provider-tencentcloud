---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_policy_conditions"
sidebar_current: "docs-tencentcloud-datasource-monitor_policy_conditions"
description: |-
  Use this data source to query monitor policy conditions(There is a lot of data and it is recommended to output to a file)
---

# tencentcloud_monitor_policy_conditions

Use this data source to query monitor policy conditions(There is a lot of data and it is recommended to output to a file)

## Example Usage

```hcl
data "tencentcloud_monitor_policy_conditions" "monitor_policy_conditions" {
  name               = "Cloud Virtual Machine"
  result_output_file = "./tencentcloud_monitor_policy_conditions.txt"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the policy name, support partial matching, eg:`Cloud Virtual Machine`,`Virtual`,`Cloud Load Banlancer-Private CLB Listener`.
* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list policy condition. Each element contains the following attributes:
  * `event_metrics` - A list of event condition metrics. Each element contains the following attributes:
    * `event_id` - The id of this event metric.
    * `event_show_name` - The name of this event metric.
    * `need_recovered` - Whether to recover.
  * `is_support_multi_region` - Whether to support multi region.
  * `metrics` - A list of event condition metrics. Each element contains the following attributes:
    * `calc_type_keys` - Calculate type of this metric.
    * `calc_type_need` - Whether `calc_type` required in the configuration.
    * `calc_value_default` - The default calculate value of this metric.
    * `calc_value_fixed` - The fixed calculate value of this metric.
    * `calc_value_max` - The max calculate value of this metric.
    * `calc_value_min` - The min calculate value of this metric.
    * `calc_value_need` - Whether `calc_value` required in the configuration.
    * `continue_time_default` - The default continue time(seconds) config for this metric.
    * `continue_time_keys` - The continue time(seconds) keys for this metric.
    * `continue_time_need` - Whether `continue_time` required in the configuration.
    * `metric_id` - The id of this metric.
    * `metric_show_name` - The name of this metric.
    * `metric_unit` - The unit of this metric.
    * `period_default` - The default data time(seconds) config for this metric.
    * `period_keys` - The data time(seconds) keys for this metric.
    * `period_need` - Whether `period` required in the configuration.
    * `period_num_default` - The default period number config for this metric.
    * `period_num_keys` - The period number keys for this metric.
    * `period_num_need` - Whether `period_num` required in the configuration.
    * `stat_type_p10` - Data aggregation mode, cycle of 10 seconds.
    * `stat_type_p1800` - Data aggregation mode, cycle of 1800 seconds.
    * `stat_type_p300` - Data aggregation mode, cycle of 300 seconds.
    * `stat_type_p3600` - Data aggregation mode, cycle of 3600 seconds.
    * `stat_type_p5` - Data aggregation mode, cycle of 5 seconds.
    * `stat_type_p600` - Data aggregation mode, cycle of 600 seconds.
    * `stat_type_p60` - Data aggregation mode, cycle of 60 seconds.
    * `stat_type_p86400` - Data aggregation mode, cycle of 86400 seconds.
  * `name` - Name of this policy name.
  * `policy_view_name` - Policy view name, eg:`cvm_device`,`BANDWIDTHPACKAGE`, refer to `data.tencentcloud_monitor_policy_conditions(policy_view_name)`.
  * `support_regions` - Support regions of this policy view.



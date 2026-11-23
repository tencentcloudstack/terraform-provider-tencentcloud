---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_statistic_data"
sidebar_current: "docs-tencentcloud-datasource-monitor_statistic_data"
description: |-
  Use this data source to query detailed information of monitor statistic_data
---

# tencentcloud_monitor_statistic_data

Use this data source to query detailed information of monitor statistic_data

## Example Usage

```hcl
data "tencentcloud_monitor_statistic_data" "statistic_data" {
  module       = "monitor"
  namespace    = "QCE/TKE2"
  metric_names = ["cpu_usage"]
  conditions {
    key      = "tke_cluster_instance_id"
    operator = "="
    value    = ["cls-mw2w40s7"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `metric_names` - (Required, Set: [`String`]) Metric name list.
* `module` - (Required, String) Module, whose value is fixed at monitor.
* `namespace` - (Required, String) Namespace. Valid values: QCE, TKE2.
* `conditions` - (Optional, List) Dimension condition. The = and in operators are supported.
* `group_bys` - (Optional, Set: [`String`]) GroupBy by the specified dimension.
* `result_output_file` - (Optional, String) Used to save results.

The `conditions` object supports the following:

* `key` - (Required, String) Dimension.
* `operator` - (Required, String) Operator. Valid values: eq (equal to), ne (not equal to), in.
* `value` - (Required, Set) Dimension value. If Operator is eq or ne, only the first element will be used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Monitoring data.
  * `metric_name` - Metric name.
  * `points` - Monitoring data point.
    * `dimensions` - Combination of instance object dimensions.
      * `name` - Instance dimension name.
      * `value` - Instance dimension value.
    * `values` - Data point list.
      * `timestamp` - Time point when this monitoring data point is generated.
      * `value` - Monitoring data point valueNote: this field may return null, indicating that no valid values can be obtained.
* `end_time` - End time, which is the current time by default, such as 2020-12-08T19:51:23+08:00.
* `period` - Statistical period in seconds. Default value: 300. Optional values: 60, 300, 3,600, and 86,400.Due to the storage period limit, the statistical period is subject to the time range of statistics:60s: The time range is less than 12 hours, and the timespan between StartTime and the current time cannot exceed 15 days.300s: The time range is less than three days, and the timespan between StartTime and the current time cannot exceed 31 days.3,600s: The time range is less than 30 days, and the timespan between StartTime and the current time cannot exceed 93 days.86,400s: The time range is less than 186 days, and the timespan between StartTime and the current time cannot exceed 186 days.
* `start_time` - Start time, which is the current time by default, such as 2020-12-08T19:51:23+08:00.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeStatisticData](https://www.tencentcloud.com/document/api/248/39481)


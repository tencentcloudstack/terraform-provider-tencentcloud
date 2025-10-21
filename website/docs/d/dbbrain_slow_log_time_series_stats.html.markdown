---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_slow_log_time_series_stats"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_slow_log_time_series_stats"
description: |-
  Use this data source to query detailed information of dbbrain slow_log_time_series_stats
---

# tencentcloud_dbbrain_slow_log_time_series_stats

Use this data source to query detailed information of dbbrain slow_log_time_series_stats

## Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_time_series_stats" "test" {
  instance_id = "%s"
  start_time  = "%s"
  end_time    = "%s"
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time, such as `2019-09-10 12:13:14`, the interval between the end time and the start time can be up to 7 days.
* `instance_id` - (Required, String) Instance ID.
* `start_time` - (Required, String) Start time, such as `2019-09-10 12:13:14`.
* `product` - (Optional, String) Service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database CynosDB for MySQL, the default is `mysql`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `period` - The unit time interval between bars, in seconds.
* `series_data` - Instan1ce cpu utilization monitoring data within a unit time interval.
  * `series` - Monitor metrics.
    * `metric` - Indicator name.
    * `unit` - Indicator unit.
    * `values` - Index value. Note: This field may return null, indicating that no valid value can be obtained.
  * `timestamp` - The timestamp corresponding to the monitoring indicator.
* `time_series` - Statistics on the number of slow logs in a unit time interval.
  * `count` - total.
  * `timestamp` - Statistics start time.



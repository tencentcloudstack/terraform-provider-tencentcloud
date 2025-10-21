---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_listener_statistics"
sidebar_current: "docs-tencentcloud-datasource-gaap_listener_statistics"
description: |-
  Use this data source to query detailed information of gaap listener statistics
---

# tencentcloud_gaap_listener_statistics

Use this data source to query detailed information of gaap listener statistics

## Example Usage

```hcl
data "tencentcloud_gaap_listener_statistics" "listener_statistics" {
  listener_id  = "listener-xxxxxx"
  start_time   = "2023-10-19 00:00:00"
  end_time     = "2023-10-19 23:59:59"
  metric_names = ["InBandwidth", "OutBandwidth", "InPackets", "OutPackets", "Concurrent"]
  granularity  = 300
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End Time.
* `granularity` - (Required, Int) Monitoring granularity, currently supporting 300 3600 86400, in seconds.The query time range does not exceed 1 day and supports a minimum granularity of 300 seconds;The query interval should not exceed 7 days and support a minimum granularity of 3600 seconds;The query interval exceeds 7 days and supports a minimum granularity of 86400 seconds.
* `listener_id` - (Required, String) Listener Id.
* `metric_names` - (Required, Set: [`String`]) List of statistical indicator names. Supporting: InBandwidth, OutBandwidth, Concurrent, InPackets, OutPackets.
* `start_time` - (Required, String) Start Time.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `statistics_data` - Channel Group Statistics.
  * `metric_data` - Metric Data.
    * `data` - Statistical data valueNote: This field may return null, indicating that a valid value cannot be obtained.
    * `time` - Time.
  * `metric_name` - Metric Name.



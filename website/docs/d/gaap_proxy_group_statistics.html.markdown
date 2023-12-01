---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_proxy_group_statistics"
sidebar_current: "docs-tencentcloud-datasource-gaap_proxy_group_statistics"
description: |-
  Use this data source to query detailed information of gaap proxy group statistics
---

# tencentcloud_gaap_proxy_group_statistics

Use this data source to query detailed information of gaap proxy group statistics

## Example Usage

```hcl
data "tencentcloud_gaap_proxy_group_statistics" "proxy_group_statistics" {
  group_id     = "link-8lpyo88p"
  start_time   = "2023-10-09 00:00:00"
  end_time     = "2023-10-09 23:59:59"
  metric_names = ["InBandwidth", "OutBandwidth", "InFlow", "OutFlow"]
  granularity  = 300
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End Time.
* `granularity` - (Required, Int) Monitoring granularity, currently supporting 60 300 3600 86400, in seconds.When the time range does not exceed 1 day, support a minimum granularity of 60 seconds;When the time range does not exceed 7 days, support a minimum granularity of 3600 seconds;When the time range does not exceed 30 days, the minimum granularity supported is 86400 seconds.
* `group_id` - (Required, String) Group Id.
* `metric_names` - (Required, Set: [`String`]) Metric Names. support, InBandwidth, OutBandwidth, Concurrent, InPackets, OutPackets.
* `start_time` - (Required, String) Start Time.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `statistics_data` - proxy Group Statistics.
  * `metric_data` - Metric Data.
    * `data` - DataNote: This field may return null, indicating that a valid value cannot be obtained.
    * `time` - Time.
  * `metric_name` - Metric Name.



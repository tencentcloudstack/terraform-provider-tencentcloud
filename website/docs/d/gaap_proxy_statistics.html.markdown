---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_proxy_statistics"
sidebar_current: "docs-tencentcloud-datasource-gaap_proxy_statistics"
description: |-
  Use this data source to query detailed information of gaap proxy statistics
---

# tencentcloud_gaap_proxy_statistics

Use this data source to query detailed information of gaap proxy statistics

## Example Usage

```hcl
data "tencentcloud_gaap_proxy_statistics" "proxy_statistics" {
  proxy_id     = "link-8lpyo88p"
  start_time   = "2023-10-09 00:00:00"
  end_time     = "2023-10-09 23:59:59"
  metric_names = ["InBandwidth", "OutBandwidth", "InFlow", "OutFlow", "InPackets", "OutPackets", "Concurrent", "HttpQPS", "HttpsQPS", "Latency", "PacketLoss"]
  granularity  = 300
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End Time(2019-03-25 12:00:00).
* `granularity` - (Required, Int) Monitoring granularity, currently supporting 60 300 3600 86400, in seconds.When the time range does not exceed 3 days, support a minimum granularity of 60 seconds;When the time range does not exceed 7 days, support a minimum granularity of 300 seconds;When the time range does not exceed 30 days, the minimum granularity supported is 3600 seconds.
* `metric_names` - (Required, Set: [`String`]) Metric Names. Valid values: InBandwidth,OutBandwidth, Concurrent, InPackets, OutPackets, PacketLoss, Latency, HttpQPS, HttpsQPS.
* `proxy_id` - (Required, String) Proxy Id.
* `start_time` - (Required, String) Start Time(2019-03-25 12:00:00).
* `isp` - (Optional, String) Operator (valid when the proxy is a three network proxy), supports CMCC, CUCC, CTCC, and merges data from the three operators if null values are passed or not passed.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `statistics_data` - proxy Statistics.
  * `metric_data` - Metric Data.
    * `data` - DataNote: This field may return null, indicating that a valid value cannot be obtained.
    * `time` - Time.
  * `metric_name` - Metric Name.



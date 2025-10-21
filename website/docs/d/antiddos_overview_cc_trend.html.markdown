---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_overview_cc_trend"
sidebar_current: "docs-tencentcloud-datasource-antiddos_overview_cc_trend"
description: |-
  Use this data source to query detailed information of antiddos overview_cc_trend
---

# tencentcloud_antiddos_overview_cc_trend

Use this data source to query detailed information of antiddos overview_cc_trend

## Example Usage

```hcl
data "tencentcloud_antiddos_overview_cc_trend" "overview_cc_trend" {
  period      = 300
  start_time  = "2023-11-20 00:00:00"
  end_time    = "2023-11-21 00:00:00"
  metric_name = "inqps"
  business    = "bgpip"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) EndTime.
* `metric_name` - (Required, String) Indicator, values [inqps (peak total requests, dropqps (peak attack requests)), incount (number of requests), dropcount (number of attacks)].
* `period` - (Required, Int) Statistical granularity, values [300 (5 minutes), 3600 (hours), 86400 (days)].
* `start_time` - (Required, String) StartTime.
* `business` - (Optional, String) Dayu sub product code (bgpip represents advanced defense IP; net represents professional version of advanced defense IP).
* `ip_list` - (Optional, Set: [`String`]) resource id list.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Data.



---
subcategory: "Anti-DDoS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_bgp_biz_trend"
sidebar_current: "docs-tencentcloud-datasource-antiddos_bgp_biz_trend"
description: |-
  Use this data source to query detailed information of antiddos bgp_biz_trend
---

# tencentcloud_antiddos_bgp_biz_trend

Use this data source to query detailed information of antiddos bgp_biz_trend

## Example Usage

```hcl
data "tencentcloud_antiddos_bgp_biz_trend" "bgp_biz_trend" {
  business    = "bgp-multip"
  start_time  = "2023-11-22 09:25:00"
  end_time    = "2023-11-22 10:25:00"
  metric_name = "intraffic"
  instance_id = "bgp-00000ry7"
  flag        = 0
}
```

## Argument Reference

The following arguments are supported:

* `business` - (Required, String) Dayu sub product code (bgpip represents advanced defense IP; net represents professional version of advanced defense IP).
* `end_time` - (Required, String) Statistic end time.
* `flag` - (Required, Int) 0 represents fixed time, 1 represents custom time.
* `instance_id` - (Required, String) Antiddos InstanceId.
* `metric_name` - (Required, String) Statistic metric name, for example: intraffic, outtraffic, inpkg, outpkg.
* `start_time` - (Required, String) Statistic start time.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data_list` - Values at various time points on the graph.
* `max_data` - Returns the maximum value of an array.
* `total` - Number of values in the curve graph.



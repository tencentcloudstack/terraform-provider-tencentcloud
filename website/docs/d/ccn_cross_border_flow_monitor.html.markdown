---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_cross_border_flow_monitor"
sidebar_current: "docs-tencentcloud-datasource-ccn_cross_border_flow_monitor"
description: |-
  Use this data source to query detailed information of vpc cross_border_flow_monitor
---

# tencentcloud_ccn_cross_border_flow_monitor

Use this data source to query detailed information of vpc cross_border_flow_monitor

## Example Usage

```hcl
data "tencentcloud_ccn_cross_border_flow_monitor" "cross_border_flow_monitor" {
  source_region      = "ap-guangzhou"
  destination_region = "ap-singapore"
  ccn_id             = "ccn-39lqkygf"
  ccn_uin            = "979137"
  period             = 60
  start_time         = "2023-01-01 00:00:00"
  end_time           = "2023-01-01 01:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String) CcnId.
* `ccn_uin` - (Required, String) CcnUin.
* `destination_region` - (Required, String) DestinationRegion.
* `end_time` - (Required, String) EndTime.
* `period` - (Required, Int) TimePeriod.
* `source_region` - (Required, String) SourceRegion.
* `start_time` - (Required, String) StartTime.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cross_border_flow_monitor_data` - monitor data of cross border.
  * `in_bandwidth` - in bandwidth, `bps`.
  * `in_pkg` - in pkg, `pps`.
  * `out_bandwidth` - out bandwidth, `bps`.
  * `out_pkg` - out pkg, `pps`.



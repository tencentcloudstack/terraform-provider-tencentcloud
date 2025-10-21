---
subcategory: "Anti-DDoS(antiddos)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_overview_attack_trend"
sidebar_current: "docs-tencentcloud-datasource-antiddos_overview_attack_trend"
description: |-
  Use this data source to query detailed information of antiddos overview_attack_trend
---

# tencentcloud_antiddos_overview_attack_trend

Use this data source to query detailed information of antiddos overview_attack_trend

## Example Usage

```hcl
data "tencentcloud_antiddos_overview_attack_trend" "overview_attack_trend" {
  type       = "ddos"
  dimension  = "attackcount"
  period     = 86400
  start_time = "2023-11-21 10:28:31"
  end_time   = "2023-11-22 10:28:31"
}
```

## Argument Reference

The following arguments are supported:

* `dimension` - (Required, String) Latitude, currently only attackcount is supported.
* `end_time` - (Required, String) Protection Overview Attack Trend End Time.
* `period` - (Required, Int) Period, currently only 86400 is supported.
* `start_time` - (Required, String) Protection Overview Attack Trend Start Time.
* `type` - (Required, String) Attack type: cc, ddos.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Number of attacks per cycle point.
* `period_point_count` - Number of period points included.



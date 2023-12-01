---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_regions"
sidebar_current: "docs-tencentcloud-datasource-monitor_tmp_regions"
description: |-
  Use this data source to query detailed information of monitor tmp_regions
---

# tencentcloud_monitor_tmp_regions

Use this data source to query detailed information of monitor tmp_regions

## Example Usage

```hcl
data "tencentcloud_monitor_tmp_regions" "tmp_regions" {
  pay_mode = 1
}
```

## Argument Reference

The following arguments are supported:

* `pay_mode` - (Optional, Int) Pay mode. `1`-Prepaid, `2`-Postpaid, `3`-All regions (default is all regions if not filled in).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region_set` - Region set.
  * `area` - Region area.
  * `region_id` - Region ID.
  * `region_name` - Region name.
  * `region_pay_mode` - Region pay mode.
  * `region_short_name` - Region short name.
  * `region_state` - Region status (0-unavailable; 1-available).
  * `region` - Region.



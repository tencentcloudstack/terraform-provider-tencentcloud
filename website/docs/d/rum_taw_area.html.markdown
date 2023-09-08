---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_taw_area"
sidebar_current: "docs-tencentcloud-datasource-rum_taw_area"
description: |-
  Use this data source to query detailed information of rum taw_area
---

# tencentcloud_rum_taw_area

Use this data source to query detailed information of rum taw_area

## Example Usage

```hcl
data "tencentcloud_rum_taw_area" "taw_area" {
  area_ids      =
  area_keys     =
  area_statuses =
}
```

## Argument Reference

The following arguments are supported:

* `area_ids` - (Optional, Set: [`Int`]) Area id.
* `area_keys` - (Optional, Set: [`String`]) Area key.
* `area_statuses` - (Optional, Set: [`Int`]) Area status `1`:valid; `2`:invalid.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `area_set` - Area list.
  * `area_abbr` - Region abbreviation.
  * `area_id` - Area id.
  * `area_key` - Area key.
  * `area_name` - Area name.
  * `area_region_code` - Area code.
  * `area_region_id` - Area code id.
  * `area_status` - Area status `1`:&amp;#39;valid&amp;#39;; `2`:&amp;#39;invalid&amp;#39;.



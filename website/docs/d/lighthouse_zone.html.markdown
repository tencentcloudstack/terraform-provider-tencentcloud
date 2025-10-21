---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_zone"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_zone"
description: |-
  Use this data source to query detailed information of lighthouse zone
---

# tencentcloud_lighthouse_zone

Use this data source to query detailed information of lighthouse zone

## Example Usage

```hcl
data "tencentcloud_lighthouse_zone" "zone" {
  order_field = "ZONE"
  order       = "ASC"
}
```

## Argument Reference

The following arguments are supported:

* `order_field` - (Optional, String) Sorting field. Valid values:
- ZONE: Sort by the availability zone.
- INSTANCE_DISPLAY_LABEL: Sort by visibility labels (HIDDEN, NORMAL and SELECTED). Default: [HIDDEN, NORMAL, SELECTED].
Sort by availability zone by default.
* `order` - (Optional, String) Specifies how availability zones are listed. Valid values:
- ASC: Ascending sort.
- DESC: Descending sort.
The default value is ASC.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zone_info_set` - List of zone info.
  * `instance_display_label` - Instance purchase page availability zone display label.
  * `zone_name` - Chinese name of availability zone.
  * `zone` - Availability zone.



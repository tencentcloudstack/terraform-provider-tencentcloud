---
subcategory: "Regional Management(region)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_zones"
sidebar_current: "docs-tencentcloud-datasource-zones"
description: |-
  Use this data source to query availability zones supported by a cloud product.
---

# tencentcloud_zones

Use this data source to query availability zones supported by a cloud product.

## Example Usage

```hcl
data "tencentcloud_zones" "example" {
  product = "cvm"
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) Product name to query, e.g. `cvm`. Use `tencentcloud_products` to get available product names.
* `result_output_file` - (Optional, String) Used to save results.
* `scene` - (Optional, Int) Scene control parameter. `0` or not set means do not query optional business whitelist; `1` means query optional business whitelist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zone_list` - Zone list.
  * `parent_zone_id` - Parent zone ID.
  * `parent_zone_name` - Parent zone description.
  * `parent_zone` - Parent zone identifier.
  * `zone_id` - Zone ID.
  * `zone_name` - Zone description, e.g. `Guangzhou Zone 3`.
  * `zone_state` - Zone status, `AVAILABLE` or `UNAVAILABLE`.
  * `zone_type` - Zone type.
  * `zone` - Zone name, e.g. `ap-guangzhou-3`.



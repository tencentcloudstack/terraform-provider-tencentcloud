---
subcategory: "Regional Management(region)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_regions"
sidebar_current: "docs-tencentcloud-datasource-regions"
description: |-
  Use this data source to query regions supported by a cloud product.
---

# tencentcloud_regions

Use this data source to query regions supported by a cloud product.

## Example Usage

```hcl
data "tencentcloud_regions" "example" {
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

* `region_list` - Region list.
  * `location_m_c` - Region description in different languages.
  * `region_id_m_c` - Region ID for console.
  * `region_name_m_c` - Region description displayed in console.
  * `region_name` - Region name, e.g. `South China (Guangzhou)`.
  * `region_state` - Region availability status, e.g. `AVAILABLE`.
  * `region_type_m_c` - Console type, null when called via API.
  * `region` - Region identifier, e.g. `ap-guangzhou`.



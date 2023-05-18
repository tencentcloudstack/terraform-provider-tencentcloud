---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_regions"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_regions"
description: |-
  Use this data source to query detailed information of sqlserver datasource_regions
---

# tencentcloud_sqlserver_regions

Use this data source to query detailed information of sqlserver datasource_regions

## Example Usage

```hcl
data "tencentcloud_sqlserver_regions" "datasource_regions" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region_set` - Region information array.
  * `region_id` - Numeric ID of region.
  * `region_name` - Region name.
  * `region_state` - Current purchasability of this region. UNAVAILABLE: not purchasable, AVAILABLE: purchasable.
  * `region` - Region ID in the format of ap-guangzhou.



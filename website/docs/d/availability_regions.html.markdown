---
subcategory: "Provider Data Sources"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_availability_regions"
sidebar_current: "docs-tencentcloud-datasource-availability_regions"
description: |-
  Use this data source to get the available regions. By default only `AVAILABLE` regions will be returned, but `UNAVAILABLE` regions can also be fetched when `include_unavailable` is specified.
---

# tencentcloud_availability_regions

Use this data source to get the available regions. By default only `AVAILABLE` regions will be returned, but `UNAVAILABLE` regions can also be fetched when `include_unavailable` is specified.

## Example Usage

```hcl
data "tencentcloud_availability_regions" "my_favourite_region" {
  name = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `include_unavailable` - (Optional) A bool variable indicates that the query will include `UNAVAILABLE` regions.
* `name` - (Optional) When specified, only the region with the exactly name match will be returned. `default` value means it consistent with the TIC resource stack region.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `regions` - A list of regions will be exported and its every element contains the following attributes:
  * `description` - The description of the region, like `Guangzhou Region`.
  * `name` - The name of the region, like `ap-guangzhou`.
  * `state` - The state of the region, indicate availability using `AVAILABLE` and `UNAVAILABLE` values.



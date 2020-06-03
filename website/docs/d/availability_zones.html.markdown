---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_availability_zones"
sidebar_current: "docs-tencentcloud-datasource-availability_zones"
description: |-
  Use this data source to get the available zones in current region. By default only `AVAILABLE` zones will be returned, but `UNAVAILABLE` zones can also be fetched when `include_unavailable` is specified.
---

# tencentcloud_availability_zones

Use this data source to get the available zones in current region. By default only `AVAILABLE` zones will be returned, but `UNAVAILABLE` zones can also be fetched when `include_unavailable` is specified.

## Example Usage

```hcl
data "tencentcloud_availability_zones" "my_favourite_zone" {
  name = "ap-guangzhou-3"
}
```

## Argument Reference

The following arguments are supported:

* `include_unavailable` - (Optional) A bool variable indicates that the query will include `UNAVAILABLE` zones.
* `name` - (Optional) When specified, only the zone with the exactly name match will be returned.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zones` - A list of zones will be exported and its every element contains the following attributes:
  * `description` - The description of the zone, like `Guangzhou Zone 3`.
  * `id` - An internal id for the zone, like `200003`, usually not so useful.
  * `name` - The name of the zone, like `ap-guangzhou-3`.
  * `state` - The state of the zone, indicate availability using `AVAILABLE` and `UNAVAILABLE` values.



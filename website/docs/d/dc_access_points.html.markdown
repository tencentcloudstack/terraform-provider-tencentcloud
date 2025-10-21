---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_access_points"
sidebar_current: "docs-tencentcloud-datasource-dc_access_points"
description: |-
  Use this data source to query detailed information of dc access_points
---

# tencentcloud_dc_access_points

Use this data source to query detailed information of dc access_points

## Example Usage

```hcl
data "tencentcloud_dc_access_points" "access_points" {
  region_id = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Optional, String) Access point region, which can be queried through `DescribeRegions`.You can call `DescribeRegions` to get the region ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_point_set` - Access point information.
  * `access_point_id` - Unique access point ID.
  * `access_point_name` - Access point name.
  * `access_point_type` - Access point type. Valid values: `VXLAN`, `QCPL`, and `QCAR`.Note: this field may return `null`, indicating that no valid values can be obtained.
  * `area` - Access point regionNote: this field may return `null`, indicating that no valid values can be obtained.
  * `available_port_type` - Available port type at the access point. Valid values: 1000BASE-T: gigabit electrical port; 1000BASE-LX: 10 km gigabit single-mode optical port; 1000BASE-ZX: 80 km gigabit single-mode optical port; 10GBASE-LR: 10 km 10-gigabit single-mode optical port; 10GBASE-ZR: 80 km 10-gigabit single-mode optical port; 10GBASE-LH: 40 km 10-gigabit single-mode optical port; 100GBASE-LR4: 10 km 100-gigabit single-mode optical portfiber optic port.Note: this field may return `null`, indicating that no valid value is obtained.
  * `city` - City where the access point is locatedNote: this field may return `null`, indicating that no valid values can be obtained.
  * `coordinate` - Latitude and longitude of the access pointNote: this field may return `null`, indicating that no valid values can be obtained.
    * `lat` - Latitude.
    * `lng` - Longitude.
  * `line_operator` - List of ISPs supported by access point.
  * `location` - Access point location.
  * `region_id` - ID of the region that manages the access point.
  * `state` - Access point status. Valid values: available, unavailable.



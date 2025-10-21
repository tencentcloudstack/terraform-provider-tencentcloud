---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_access_regions_by_dest_region"
sidebar_current: "docs-tencentcloud-datasource-gaap_access_regions_by_dest_region"
description: |-
  Use this data source to query detailed information of gaap access regions by dest region
---

# tencentcloud_gaap_access_regions_by_dest_region

Use this data source to query detailed information of gaap access regions by dest region

## Example Usage

```hcl
data "tencentcloud_gaap_access_regions_by_dest_region" "access_regions_by_dest_region" {
  dest_region = "SouthChina"
}
```

## Argument Reference

The following arguments are supported:

* `dest_region` - (Required, String) Origin region.
* `ip_address_version` - (Optional, String) IP version, can be taken as IPv4 or IPv6, with a default value of IPv4.
* `package_type` - (Optional, String) Channel package type, where Thunder represents a standard proxy group, Accelerator represents a game accelerator proxy, and CrossBorder represents a cross-border proxy.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_region_set` - List of available acceleration zone information.
  * `bandwidth_list` - Optional bandwidth value array.
  * `concurrent_list` - Optional concurrency value array.
  * `feature_bitmap` - The type of computer room, where dc represents the DataCenter data center, ec represents the feature bitmap, and each bit represents a feature, where:0, indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained. Edge nodes.
  * `idc_type` - The type of computer room, where dc represents the DataCenter data center and ec represents the EdgeComputing edge node.
  * `region_area_name` - Region name of the computer room.
  * `region_area` - Region of the computer room.
  * `region_id` - Region id.
  * `region_name` - Chinese or English name of the region.



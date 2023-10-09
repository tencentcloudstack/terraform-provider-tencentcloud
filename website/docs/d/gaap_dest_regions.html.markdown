---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_dest_regions"
sidebar_current: "docs-tencentcloud-datasource-gaap_dest_regions"
description: |-
  Use this data source to query detailed information of gaap dest regions
---

# tencentcloud_gaap_dest_regions

Use this data source to query detailed information of gaap dest regions

## Example Usage

```hcl
data "tencentcloud_gaap_dest_regions" "dest_regions" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dest_region_set` - Source Site Area Details List.
  * `feature_bitmap` - Property bitmap, where each bit represents a property, where:0, indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `idc_type` - The type of computer room, where dc represents the DataCenter data center and ec represents the EdgeComputing edge node.
  * `region_area_name` - Region name of the computer room.
  * `region_area` - Region of the computer room.
  * `region_id` - region ID.
  * `region_name` - region name.
  * `support_feature` - Ability to access regional supportNote: This field may return null, indicating that a valid value cannot be obtained.
    * `network_type` - A list of network types supported by the access area, with normal indicating support for regular BGP, cn2 indicating premium BGP, triple indicating three networks, and secure_EIP represents a custom secure EIP.



---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_region_and_price"
sidebar_current: "docs-tencentcloud-datasource-gaap_region_and_price"
description: |-
  Use this data source to query detailed information of gaap region and price
---

# tencentcloud_gaap_region_and_price

Use this data source to query detailed information of gaap region and price

## Example Usage

```hcl
data "tencentcloud_gaap_region_and_price" "region_and_price" {
}
```

## Argument Reference

The following arguments are supported:

* `ip_address_version` - (Optional, String) IP version. Available values: IPv4, IPv6. Default is IPv4.
* `package_type` - (Optional, String) Type of channel package. `Thunder` represents standard channel group, `Accelerator` represents game accelerator channel, and `CrossBorder` represents cross-border channel.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bandwidth_unit_price` - proxy bandwidth cost gradient price.
  * `bandwidth_range` - Band width Range.
  * `bandwidth_unit_price` - Band width Unit Price, Unit:yuan/Mbps/day.
  * `discount_bandwidth_unit_price` - Bandwidth discount price, unit:yuan/Mbps/day.
* `currency` - Bandwidth Price Currency Type:CNYUSD.
* `dest_region_set` - Source Site Area Details List.
  * `feature_bitmap` - Property bitmap, where each bit represents a property, where:0, indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `idc_type` - Type of computer room, dc represents DataCenter data center, ec represents EdgeComputing edge node.
  * `region_area_name` - Region name of the computer room.
  * `region_area` - Region of the computer room.
  * `region_id` - Region Id.
  * `region_name` - Region Name.
  * `support_feature` - Ability to access regional supportNote: This field may return null, indicating that a valid value cannot be obtained.
    * `network_type` - A list of network types supported by the access area, with `normal` indicating support for regular BGP, `cn2` indicating premium BGP, `triple` indicating three networks, and `secure_eip` represents a custom secure EIP.



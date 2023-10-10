---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_proxy_groups"
sidebar_current: "docs-tencentcloud-datasource-gaap_proxy_groups"
description: |-
  Use this data source to query detailed information of gaap proxy groups
---

# tencentcloud_gaap_proxy_groups

Use this data source to query detailed information of gaap proxy groups

## Example Usage

```hcl
data "tencentcloud_gaap_proxy_groups" "proxy_groups" {
  project_id = 0
  filters {
    name   = "GroupId"
    values = ["lg-5anbbou5"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, Int) Project ID. Value range:-1, All projects under this user0, default projectOther values, specified items.
* `filters` - (Optional, List) Filter conditions,The upper limit of Filter.Values per request is 5.RealServerRegion - String - Required: No - (filtering criteria) Filter by real server region, refer to the RegionId in the returned results of the DescribeDestRegions interface.PackageType - String - Required: No - (Filter condition) proxy group type, where &amp;#39;Thunder&amp;#39; represents the standard proxy group and &amp;#39;Accelerator&amp;#39; represents the silver acceleration proxy group.
* `result_output_file` - (Optional, String) Used to save results.
* `tag_set` - (Optional, List) Tag list, when this field exists, pulls the resource list under the corresponding tag.Supports a maximum of 5 labels. When there are two or more labels and any one of them is met, the proxy group will be pulled out.

The `filters` object supports the following:

* `name` - (Required, String) Filter conditions.
* `values` - (Required, Set) filtering value.

The `tag_set` object supports the following:

* `tag_key` - (Required, String) Tag Key.
* `tag_value` - (Required, String) Tag Value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `proxy_group_list` - List of proxy groups.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `create_time` - Create TimeNote: This field may return null, indicating that a valid value cannot be obtained.
  * `domain` - proxy group domain nameNote: This field may return null, indicating that a valid value cannot be obtained.
  * `feature_bitmap` - Property bitmap, where each bit represents a property, where:0, indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `group_id` - proxy group Id.
  * `group_name` - proxy Group NameNote: This field may return null, indicating that a valid value cannot be obtained.
  * `http3_supported` - Supports the identification of Http3 features, where:0 indicates shutdown;1 indicates enabled.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `project_id` - Project Id.
  * `proxy_type` - Does the proxy group include Microsoft proxysNote: This field may return null, indicating that a valid value cannot be obtained.
  * `real_server_region_info` - Real Server Region Info.
    * `feature_bitmap` - Property bitmap, where each bit represents a property, where:0, indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.
    * `idc_type` - The type of computer room, where &#39;dc&#39; represents the DataCenter data center and &#39;ec&#39; represents the EdgeComputing edge node.
    * `region_area_name` - Region name of the computer room.
    * `region_area` - Region of the computer room.
    * `region_id` - Region Id.
    * `region_name` - Region Name.
    * `support_feature` - Ability to access regional supportNote: This field may return null, indicating that a valid value cannot be obtained.
      * `network_type` - A list of network types supported by the access area, with &#39;normal&#39; indicating support for regular BGP, &#39;cn2&#39; indicating premium BGP, &#39;triple&#39; indicating three networks, and &#39;secure_EIP&#39; represents a custom secure EIP.
  * `status` - proxy group status.Among them,&#39;RUNNING&#39; indicates running;&#39;CREATING&#39; indicates being created;&#39;DESTROYING&#39; indicates being destroyed;&#39;MOVING&#39; indicates that the proxy is being migrated;&#39;CHANGING&#39; indicates partial deployment.
  * `tag_set` - Tag Set.
    * `tag_key` - Tag Key.
    * `tag_value` - Tag Value.
  * `version` - proxy Group VersionNote: This field may return null, indicating that a valid value cannot be obtained.



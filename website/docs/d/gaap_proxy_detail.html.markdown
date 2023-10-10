---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_proxy_detail"
sidebar_current: "docs-tencentcloud-datasource-gaap_proxy_detail"
description: |-
  Use this data source to query detailed information of gaap proxy detail
---

# tencentcloud_gaap_proxy_detail

Use this data source to query detailed information of gaap proxy detail

## Example Usage

```hcl
data "tencentcloud_gaap_proxy_detail" "proxy_detail" {
  proxy_id = "link-8lpyo88p"
}
```

## Argument Reference

The following arguments are supported:

* `proxy_id` - (Required, String) Proxy Id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `proxy_detail` - Proxy Detail.
  * `access_region_info` - Detailed information about the access region, including the region ID and domain name.Note: This field may return null, indicating that a valid value cannot be obtained.
    * `feature_bitmap` - Property bitmap, where each bit represents a property, where:0 indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.
    * `idc_type` - The type of computer room, where dc represents the DataCenter data center and ec represents the EdgeComputing edge node.
    * `region_area_name` - Region name of the computer room.
    * `region_area` - Region of the computer room.
    * `region_id` - Region Id.
    * `region_name` - Region Name.
    * `support_feature` - Ability to access regional supportNote: This field may return null, indicating that a valid value cannot be obtained.
      * `network_type` - A list of network types supported by the access area, with normal indicating support for regular BGP, cn2 indicating premium BGP, triple indicating three networks, and secure_EIP represents a custom secure EIP.
  * `access_region` - Access Region.
  * `ban_status` - Blocking and Unblocking Status: BANNED indicates that the ban has been lifted, RECOVER indicates that the ban has been lifted or not, BANNING indicates that the ban is in progress, RECOVERING indicates that the ban is being lifted, BAN_FAILED indicates that the ban has failed, RECOVER_FAILED indicates that the unblocking has failed.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `bandwidth` - Band width, in Mbps.
  * `billing_type` - Billing type: 0 represents bandwidth based billing, and 1 represents traffic based billing.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `client_ip_method` - The method of obtaining client IP through proxys, where 0 represents TOA and 1 represents Proxy ProtocolNote: This field may return null, indicating that a valid value cannot be obtained.
  * `concurrent` - Concurrent, in 10000 pieces/second.
  * `create_time` - The creation time, using a Unix timestamp, represents the number of seconds that have passed since January 1, 1970 (midnight UTC/GMT).
  * `domain` - Domain.
  * `feature_bitmap` - Property bitmap, where each bit represents a property, where:0 indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `forward_ip` - proxy forwarding IP.
  * `group_id` - proxy group ID, which exists when a proxy belongs to a certain proxy group.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `http3_supported` - Identification that supports the Http3 protocol, where:0 indicates shutdown;1 indicates enabled.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `in_ban_blacklist` - Is it on the banned blacklist? 0 indicates not on the blacklist, and 1 indicates on the blacklist.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `instance_id` - (Old parameter, please use ProxyId) Proxy instance ID.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `ip_address_version` - IP version: IPv4, IPv6Note: This field may return null, indicating that a valid value cannot be obtained.
  * `ip_list` - IP ListNote: This field may return null, indicating that a valid value cannot be obtained.
    * `bandwidth` - Band width.
    * `ip` - IP.
    * `provider` - Supplier, BGP represents default, CMCC represents China Mobile, CUCC represents China Unicom, and CTCC represents China Telecom.
  * `ip` - IP.
  * `modify_config_time` - Configuration change timeNote: This field may return null, indicating that a valid value cannot be obtained.
  * `network_type` - Network type: normal represents regular BGP, cn2 represents premium BGP, triple represents triple network, secure_EIP represents customized security EIPNote: This field may return null, indicating that a valid value cannot be obtained.
  * `package_type` - proxy package type: Thunder represents standard proxy, Accelerator represents silver acceleration proxy,CrossBorder represents a cross-border proxy.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `policy_id` - Security policy ID, which exists when a security policy is set.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `project_id` - Project Id.
  * `proxy_id` - (New parameter) proxy instance ID.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `proxy_name` - Proxy Name.
  * `proxy_type` - proxy type, 100 represents THUNDER proxy, 103 represents Microsoft cooperation proxyNote: This field may return null, indicating that a valid value cannot be obtained.
  * `real_server_region_info` - Detailed information of the real server region, including the region ID and domain name.Note: This field may return null, indicating that a valid value cannot be obtained.
    * `feature_bitmap` - Property bitmap, where each bit represents a property, where:0 indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.
    * `idc_type` - The type of computer room, where dc represents the DataCenter data center and ec represents the EdgeComputing edge node.
    * `region_area_name` - Region name of the computer room.
    * `region_area` - Region of the computer room.
    * `region_id` - Region Id.
    * `region_name` - Region Name.
    * `support_feature` - Ability to access regional supportNote: This field may return null, indicating that a valid value cannot be obtained.
      * `network_type` - A list of network types supported by the access area, with normal indicating support for regular BGP, cn2 indicating premium BGP, triple indicating three networks, and secure_EIP represents a custom secure EIP.
  * `real_server_region` - Real Server Region.
  * `related_global_domains` - List of domain names associated with resolutionNote: This field may return null, indicating that a valid value cannot be obtained.
  * `scalarable` - 1. This proxy can be scaled and expanded; 0, this proxy cannot be scaled or expanded.
  * `status` - proxy status. Among them:RUNNING indicates running;CREATING indicates being created;DESTROYING indicates being destroyed;OPENING indicates being opened;CLOSING indicates being closed;Closed indicates that it has been closed;ADJUSTING represents a configuration change in progress;ISOLATING indicates being isolated;ISOLATED indicates that it has been isolated;CLONING indicates copying;RECOVERING indicates that the proxy is being maintained;MOVING indicates that migration is in progress.
  * `support_protocols` - Supported protocol types.
  * `support_security` - Does it support security group configurationNote: This field may return null, indicating that a valid value cannot be obtained.
  * `tag_set` - tag list, when there are no labels, this field is an empty list.Note: This field may return null, indicating that a valid value cannot be obtained.
    * `tag_key` - Tag Key.
    * `tag_value` - Tag Value.
  * `version` - Version 1.0, 2.0, 3.0.



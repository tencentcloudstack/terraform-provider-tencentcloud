---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_origin_acl"
sidebar_current: "docs-tencentcloud-datasource-teo_origin_acl"
description: |-
  Use this data source to query detailed information of TEO origin acl
---

# tencentcloud_teo_origin_acl

Use this data source to query detailed information of TEO origin acl

## Example Usage

### Query origin acl by zone Id

```hcl
data "tencentcloud_teo_origin_acl" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Specifies the site ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `origin_acl_info` - Describes the binding relationship between the l7 acceleration domain/l4 proxy instance and the origin server IP range.
  * `current_origin_acl` - Currently effective origin ACLs. This field is empty when origin protection is not enabled.
Note: This field may return null, which indicates a failure to obtain a valid value.
    * `active_time` - Version effective time in UTC+8, following the date and time format of the ISO 8601 standard.
Note: This field may return null, which indicates a failure to obtain a valid value.
    * `entire_addresses` - IP range details.
Note: This field may return null, which indicates a failure to obtain a valid value.
      * `i_pv4` - (**Deprecated**) Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead. IPv4 subnet.
      * `i_pv6` - (**Deprecated**) Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead. IPv6 subnet.
      * `ipv4` - IPv4 subnet.
      * `ipv6` - IPv6 subnet.
    * `is_planed` - This parameter is used to record whether "I've upgraded to the lastest version" is completed before the origin ACLs version is effective. valid values:.
- true: specifies that the version is effective and the update to the latest version is confirmed.
- false: when the version takes effect, the confirmation of updating to the latest origin ACLs are not completed. The IP range is forcibly updated to the latest version in the backend. When this parameter returns false, please confirm in time whether your origin server firewall configuration has been updated to the latest version to avoid origin-pull failure.
Note: This field may return null, which indicates a failure to obtain a valid value.
    * `version` - Version number.
Note: This field may return null, which indicates a failure to obtain a valid value.
  * `l4_proxy_ids` - The list of L4 proxy instances that enable the origin ACLs. This field is empty when origin protection is not enabled.
  * `l7_hosts` - The list of L7 accelerated domains that enable the origin ACLs. This field is empty when origin protection is not enabled.
  * `next_origin_acl` - When the origin ACLs are updated, this field will be returned with the next version's origin IP range to take effect, including a comparison with the current origin IP range. This field is empty if not updated or origin protection is not enabled.
Note: This field may return null, which indicates a failure to obtain a valid value.
    * `added_addresses` - The latest origin IP range newly-added compared with the origin IP range in CurrentOrginACL.
      * `i_pv4` - (**Deprecated**) Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead. IPv4 subnet.
      * `i_pv6` - (**Deprecated**) Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead. IPv6 subnet.
      * `ipv4` - IPv4 subnet.
      * `ipv6` - IPv6 subnet.
    * `entire_addresses` - IP range details.
      * `i_pv4` - (**Deprecated**) Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead. IPv4 subnet.
      * `i_pv6` - (**Deprecated**) Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead. IPv6 subnet.
      * `ipv4` - IPv4 subnet.
      * `ipv6` - IPv6 subnet.
    * `no_change_addresses` - The latest origin IP range is unchanged compared with the origin IP range in CurrentOrginACL.
      * `i_pv4` - (**Deprecated**) Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead. IPv4 subnet.
      * `i_pv6` - (**Deprecated**) Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead. IPv6 subnet.
      * `ipv4` - IPv4 subnet.
      * `ipv6` - IPv6 subnet.
    * `planned_active_time` - Version effective time, which adopts UTC+8 and follows the date and time format of the ISO 8601 standard.
    * `removed_addresses` - The latest origin IP range deleted compared with the origin IP range in CurrentOrginACL.
      * `i_pv4` - (**Deprecated**) Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead. IPv4 subnet.
      * `i_pv6` - (**Deprecated**) Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead. IPv6 subnet.
      * `ipv4` - IPv4 subnet.
      * `ipv6` - IPv6 subnet.
    * `version` - Version number.
  * `status` - Origin protection status. Vaild values:
- online: in effect;
- offline: disabled;
- updating: configuration deployment in progress.



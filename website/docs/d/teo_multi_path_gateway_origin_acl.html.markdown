---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateway_origin_acl"
sidebar_current: "docs-tencentcloud-datasource-teo_multi_path_gateway_origin_acl"
description: |-
  Use this data source to query detailed information of TEO multi-path gateway origin acl
---

# tencentcloud_teo_multi_path_gateway_origin_acl

Use this data source to query detailed information of TEO multi-path gateway origin acl

## Example Usage

### Query multi-path gateway origin acl by zone_id and gateway_id

```hcl
data "tencentcloud_teo_multi_path_gateway_origin_acl" "example" {
  zone_id    = "zone-2noqxz9b6klw"
  gateway_id = "gw-12345678"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) Gateway ID.
* `zone_id` - (Required, String) Zone ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `multi_path_gateway_origin_acl_info` - Multi-path gateway origin ACL info.
  * `multi_path_gateway_current_origin_acl` - Currently effective origin ACLs.
    * `entire_addresses` - IP CIDR details.
      * `ipv4` - IPv4 CIDR list.
      * `ipv6` - IPv6 CIDR list.
    * `is_planed` - Whether the update to the latest origin IP CIDR has been confirmed.
    * `version` - Version number.
  * `multi_path_gateway_next_origin_acl` - Next version origin ACLs.
    * `added_addresses` - Added IP CIDRs compared to current.
      * `ipv4` - IPv4 CIDR list.
      * `ipv6` - IPv6 CIDR list.
    * `entire_addresses` - IP CIDR details.
      * `ipv4` - IPv4 CIDR list.
      * `ipv6` - IPv6 CIDR list.
    * `no_change_addresses` - Unchanged IP CIDRs compared to current.
      * `ipv4` - IPv4 CIDR list.
      * `ipv6` - IPv6 CIDR list.
    * `removed_addresses` - Removed IP CIDRs compared to current.
      * `ipv4` - IPv4 CIDR list.
      * `ipv6` - IPv6 CIDR list.
    * `version` - Version number.



---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_confirm_multi_path_gateway_origin_acl"
sidebar_current: "docs-tencentcloud-resource-teo_confirm_multi_path_gateway_origin_acl"
description: |-
  Provides a resource to manage the confirmation of TEO multi-path gateway origin ACL updates.
---

# tencentcloud_teo_confirm_multi_path_gateway_origin_acl

Provides a resource to manage the confirmation of TEO multi-path gateway origin ACL updates.

## Example Usage

### Confirm origin ACL version

```hcl
resource "tencentcloud_teo_confirm_multi_path_gateway_origin_acl" "example" {
  zone_id            = "zone-3edjdliiw3he"
  gateway_id         = "gw-abc12345"
  origin_acl_version = 2
}
```

### Read-only without confirming

```hcl
resource "tencentcloud_teo_confirm_multi_path_gateway_origin_acl" "example" {
  zone_id    = "zone-3edjdliiw3he"
  gateway_id = "gw-abc12345"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Gateway ID.
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `origin_acl_version` - (Optional, Int) Origin ACL version number to confirm.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `multi_path_gateway_origin_acl_info` - Multi-path gateway origin ACL info.
  * `multi_path_gateway_current_origin_acl` - Current active origin ACL info.
    * `entire_addresses` - Current origin IP segment details.
      * `ipv4` - IPv4 segment list.
      * `ipv6` - IPv6 segment list.
    * `is_planed` - Whether the update confirmation is completed.
    * `version` - Current version number.
  * `multi_path_gateway_next_origin_acl` - Next version origin ACL info when there is a pending update.
    * `added_addresses` - Added IP segments compared to current ACL.
      * `ipv4` - IPv4 segment list.
      * `ipv6` - IPv6 segment list.
    * `entire_addresses` - Next version origin IP segment details.
      * `ipv4` - IPv4 segment list.
      * `ipv6` - IPv6 segment list.
    * `no_change_addresses` - Unchanged IP segments compared to current ACL.
      * `ipv4` - IPv4 segment list.
      * `ipv6` - IPv6 segment list.
    * `removed_addresses` - Removed IP segments compared to current ACL.
      * `ipv4` - IPv4 segment list.
      * `ipv6` - IPv6 segment list.
    * `version` - Next version number.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.

## Import

TEO confirm multi-path gateway origin ACL can be imported using the id, e.g.

```
terraform import tencentcloud_teo_confirm_multi_path_gateway_origin_acl.example zone-3edjdliiw3he#gw-abc12345
```


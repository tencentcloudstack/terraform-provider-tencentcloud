---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_origin_acl"
sidebar_current: "docs-tencentcloud-resource-teo_origin_acl"
description: |-
  Provides a resource to create a TEO origin acl
---

# tencentcloud_teo_origin_acl

Provides a resource to create a TEO origin acl

~> **NOTE:** This resource must exclusive in one origin acl, do not declare additional rule resources of this origin acl elsewhere.

## Example Usage

```hcl
resource "tencentcloud_teo_origin_acl" "example" {
  zone_id = "zone-39quuimqg8r6"
  l7_hosts = [
    "example1.com",
    "example2.com",
    "example3.com",
  ]

  l4_proxy_ids = [
    "sid-3dwf5252ravl",
    "sid-3dwfxzt8ed3l",
    "sid-3dwfy5mpwnk4",
    "sid-3dwfyaj6qeys",
  ]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Specifies the site ID.
* `l4_proxy_ids` - (Optional, Set: [`String`]) he list of L4 proxy Instances that require enabling origin ACLs. This list must be empty when the request parameter L4EnableMode is set to 'all'.
* `l7_hosts` - (Optional, Set: [`String`]) The list of L7 acceleration domains that require enabling the origin ACLs. This list must be empty when the request parameter L7EnableMode is set to 'all'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `30m`) Used when creating the resource.
* `update` - (Defaults to `30m`) Used when updating the resource.
* `delete` - (Defaults to `30m`) Used when deleting the resource.

## Import

TEO origin acl can be imported using the zone_id, e.g.

````
terraform import tencentcloud_teo_origin_acl.example zone-39quuimqg8r6
````


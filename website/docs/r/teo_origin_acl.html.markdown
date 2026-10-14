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
  zone_id           = "zone-3fkff38fyw8s"
  origin_acl_family = "gaz"
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
* `origin_acl_family` - (Optional, String) Origin ACL control domain. Valid values: gaz, mlc, emc, plat-gaz, plat-mlc, plat-emc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `30m`) Used when creating the resource.
* `update` - (Defaults to `30m`) Used when updating the resource.
* `delete` - (Defaults to `30m`) Used when deleting the resource.
## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeOriginACL](https://cloud.tencent.com/document/api/1552/120408) - This API is used to query the binding relationship between L7 acceleration domains/L4 proxy instances and origin ACLs under a site, as well as IP range details.
* Create: [EnableOriginACL](https://cloud.tencent.com/document/api/1552/120406) - This API is used to enable origin protection for a site for the first time. Enabled, EdgeOne will use specific origin-pull IP ranges for L7 acceleration domains/L4 proxy instances.
* Create: [ModifyOriginACL](https://cloud.tencent.com/document/api/1552/120405) - This API is used to enable or disable specific origin IP ranges for L7 acceleration domains or L4 proxy instances. The maximum quantity for single submission is 200 L7 acceleration domains or 100 L4 proxy instances, with mixed submissions
* Read: [DescribeOriginACL](https://cloud.tencent.com/document/api/1552/120408) - This API is used to query the binding relationship between L7 acceleration domains/L4 proxy instances and origin ACLs under a site, as well as IP range details.
* Update: [DescribeOriginACL](https://cloud.tencent.com/document/api/1552/120408) - This API is used to query the binding relationship between L7 acceleration domains/L4 proxy instances and origin ACLs under a site, as well as IP range details.
* Update: [ModifyOriginACL](https://cloud.tencent.com/document/api/1552/120405) - This API is used to enable or disable specific origin IP ranges for L7 acceleration domains or L4 proxy instances. The maximum quantity for single submission is 200 L7 acceleration domains or 100 L4 proxy instances, with mixed submissions
* Delete: [DescribeOriginACL](https://cloud.tencent.com/document/api/1552/120408) - This API is used to query the binding relationship between L7 acceleration domains/L4 proxy instances and origin ACLs under a site, as well as IP range details.
* Delete: [DisableOriginACL](https://cloud.tencent.com/document/api/1552/120407) - This API is used to disable 'Origin Protection' of a site. Once disabled, resources related to it will no longer use only the origin ACLs provided by "origin protection" to request your origin, and stops sending update notifications on the

## Import

TEO origin acl can be imported using the id, e.g.

````
terraform import tencentcloud_teo_origin_acl.example zone-3fkff38fyw8s
````


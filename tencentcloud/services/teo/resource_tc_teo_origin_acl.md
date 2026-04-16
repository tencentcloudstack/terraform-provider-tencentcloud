Provides a resource to create a TEO origin acl

~> **NOTE:** This resource must exclusive in one origin acl, do not declare additional rule resources of this origin acl elsewhere.

Example Usage

```hcl
resource "tencentcloud_teo_origin_acl" "example" {
  zone_id        = "zone-39quuimqg8r6"
  origin_acl_family = "gaz"
  l7_hosts       = [
    "example1.com",
    "example2.com",
    "example3.com",
  ]

  l4_proxy_ids   = [
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

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Specifies the site ID.
* `origin_acl_family` - (Optional, Computed) The control domain for origin ACL. Available values: 'gaz' (Global Availability Zone), 'mlc' (Mainland China), 'emc' (Excluding Mainland China), 'plat-gaz' (Platform Global Availability Zone), 'plat-mlc' (Platform Mainland China), 'plat-emc' (Platform Excluding Mainland China). If not specified, the default global control domain will be used.
* `l7_hosts` - (Optional, Computed) The list of L7 acceleration domains that require enabling the origin ACLs. This list must be empty when the request parameter L7EnableMode is set to 'all'.
* `l4_proxy_ids` - (Optional, Computed) The list of L4 proxy Instances that require enabling origin ACLs. This list must be empty when the request parameter L4EnableMode is set to 'all'.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as `zone_id`.

Import

TEO origin acl can be imported using the zone_id, e.g.

````
terraform import tencentcloud_teo_origin_acl.example zone-39quuimqg8r6
````

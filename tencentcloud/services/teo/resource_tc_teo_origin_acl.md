Provides a resource to create a TEO origin acl

~> **NOTE:** This resource must exclusive in one origin acl, do not declare additional rule resources of this origin acl elsewhere.

Example Usage

```hcl
resource "tencentcloud_teo_origin_acl" "example" {
  zone_id        = "zone-39quuimqg8r6"
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

Import

TEO origin acl can be imported using the zone_id, e.g.

````
terraform import tencentcloud_teo_origin_acl.example zone-39quuimqg8r6
````

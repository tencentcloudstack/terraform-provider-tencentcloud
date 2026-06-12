---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_domain_shared_cname_attachment"
sidebar_current: "docs-tencentcloud-resource-teo_domain_shared_cname_attachment"
description: |-
  Provides a resource to manage TEO (EdgeOne) acceleration domain shared CNAME binding attachment
---

# tencentcloud_teo_domain_shared_cname_attachment

Provides a resource to manage TEO (EdgeOne) acceleration domain shared CNAME binding attachment

## Example Usage

```hcl
resource "tencentcloud_teo_domain_shared_cname_attachment" "example" {
  zone_id = "zone-2qtuhspy7cr6"

  bind_shared_cname_maps {
    shared_cname = "shared.example.com"
    domain_names = ["domain1.example.com", "domain2.example.com"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `bind_shared_cname_maps` - (Required, List, ForceNew) The binding relationships between acceleration domains and shared CNAMEs.
* `zone_id` - (Required, String, ForceNew) The zone ID that the acceleration domain belongs to.

The `bind_shared_cname_maps` object supports the following:

* `domain_names` - (Required, List, ForceNew) The acceleration domain names to bind, up to 20.
* `shared_cname` - (Required, String, ForceNew) The shared CNAME to bind to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo domain_shared_cname_attachment can be imported using the zone_id#shared_cname#domain_names (domain names joined by comma), e.g.
```
terraform import tencentcloud_teo_domain_shared_cname_attachment.example zone-2qtuhspy7cr6#shared.example.com#domain1.example.com,domain2.example.com
```


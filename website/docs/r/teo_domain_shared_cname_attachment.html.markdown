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
  zone_id      = "zone-2qtuhspy7cr6"
  shared_cname = "shared.example.com"
  domain_names = ["domain1.example.com", "domain2.example.com"]
}
```

## Argument Reference

The following arguments are supported:

* `domain_names` - (Required, Set: [`String`]) The acceleration domain names to bind.
* `shared_cname` - (Required, String, ForceNew) The shared CNAME to bind to.
* `zone_id` - (Required, String, ForceNew) The zone ID that the acceleration domain belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo domain_shared_cname_attachment can be imported using the zone_id#shared_cname, e.g.
```
terraform import tencentcloud_teo_domain_shared_cname_attachment.example zone-2qtuhspy7cr6#shared.example.com
```


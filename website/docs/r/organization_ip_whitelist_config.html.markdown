---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_ip_whitelist_config"
sidebar_current: "docs-tencentcloud-resource-organization_ip_whitelist_config"
description: |-
  Provides a resource to create a Organization IP whitelist config
---

# tencentcloud_organization_ip_whitelist_config

Provides a resource to create a Organization IP whitelist config

## Example Usage

```hcl
resource "tencentcloud_organization_ip_whitelist_config" "example" {
  zone_id = "z-1os7c9znogct"
  ip_whitelist = [
    "10.0.0.0/24",
    "192.168.1.0/24",
    "172.16.10.0/24",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `ip_whitelist` - (Required, List: [`String`]) IP whitelist entries.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Organization IP whitelist config can be imported using the zoneId, e.g.

```
terraform import tencentcloud_organization_ip_whitelist_config.example z-1os7c9znogct
```


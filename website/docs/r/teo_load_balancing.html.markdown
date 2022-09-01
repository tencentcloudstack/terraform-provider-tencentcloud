---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_load_balancing"
sidebar_current: "docs-tencentcloud-resource-teo_load_balancing"
description: |-
  Provides a resource to create a teo loadBalancing
---

# tencentcloud_teo_load_balancing

Provides a resource to create a teo loadBalancing

## Example Usage

```hcl
resource "tencentcloud_teo_load_balancing" "load_balancing" {
  zone_id = tencentcloud_teo_zone.zone.id

  host = "sfurnace.work"
  origin_id = [
    split("#", tencentcloud_teo_origin_group.group0.id)[1]
  ]
  ttl  = 600
  type = "proxied"
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) Subdomain name. You can use @ to represent the root domain.
* `origin_id` - (Required, Set: [`String`]) ID of the origin group used.
* `type` - (Required, String) Proxy mode. Valid values: dns_only: Only DNS, proxied: Enable proxy.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `ttl` - (Optional, Int) Indicates DNS TTL time when Type=dns_only.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - Schedules domain names, Note: This field may return null, indicating that no valid value can be obtained.
* `load_balancing_id` - CLB instance ID.
* `update_time` - Update time.


## Import

teo loadBalancing can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_load_balancing.loadBalancing loadBalancing_id
```


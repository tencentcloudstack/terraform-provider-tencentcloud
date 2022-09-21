---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_load_balancing"
sidebar_current: "docs-tencentcloud-resource-teo_load_balancing"
description: |-
  Provides a resource to create a teo load_balancing
---

# tencentcloud_teo_load_balancing

Provides a resource to create a teo load_balancing

## Example Usage

```hcl
resource "tencentcloud_teo_load_balancing" "load_balancing" {
  #  backup_origin_group_id = "origin-a499ca4b-3721-11ed-b9c1-5254005a52aa"
  host            = "www.toutiao2.com"
  origin_group_id = "origin-4f8a30b2-3720-11ed-b66b-525400dceb86"
  status          = "online"
  tags            = {}
  ttl             = 600
  type            = "proxied"
  zone_id         = "zone-297z8rf93cfw"
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) Subdomain name. You can use @ to represent the root domain.
* `origin_group_id` - (Required, String) ID of the origin group to use.
* `type` - (Required, String) Proxy mode.- `dns_only`: Only DNS.- `proxied`: Enable proxy.
* `zone_id` - (Required, String) Site ID.
* `backup_origin_group_id` - (Optional, String) ID of the backup origin group to use.
* `status` - (Optional, String) Status of the task. Valid values to set: `online`, `offline`. During status change, the status is `process`.
* `ttl` - (Optional, Int) Indicates DNS TTL time when `Type` is dns_only.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - Schedules domain names. Note: This field may return null, indicating that no valid value can be obtained.
* `load_balancing_id` - Load balancer instance ID.
* `update_time` - Last modification date.


## Import

teo load_balancing can be imported using the zone_id#loadBalancing_id, e.g.
```
$ terraform import tencentcloud_teo_load_balancing.load_balancing zone-297z8rf93cfw#lb-2a93c649-3719-11ed-b9c1-5254005a52aa
```


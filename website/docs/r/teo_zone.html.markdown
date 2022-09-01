---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone"
sidebar_current: "docs-tencentcloud-resource-teo_zone"
description: |-
  Provides a resource to create a teo zone
---

# tencentcloud_teo_zone

Provides a resource to create a teo zone

## Example Usage

```hcl
resource "tencentcloud_teo_zone" "zone" {
  name           = "sfurnace.work"
  plan_type      = "ent_cm_with_bot"
  type           = "full"
  paused         = false
  cname_speed_up = "enabled"

  #  vanity_name_servers {
  #    switch  = "on"
  #    servers = ["2.2.2.2"]
  #  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Site name.
* `plan_type` - (Required, String) Plan type of the zone. See details in data source `zone_available_plans`.
* `cname_speed_up` - (Optional, String) Specifies whether to enable CNAME acceleration, enabled: Enable; disabled: Disable.
* `paused` - (Optional, Bool) Indicates whether the site is disabled.
* `tags` - (Optional, Map) Tag description list.
* `type` - (Optional, String) Specifies how the site is connected to EdgeOne.
* `vanity_name_servers` - (Optional, List) User-defined name server information.

The `vanity_name_servers` object supports the following:

* `switch` - (Required, String) Whether to enable the custom name server.
* `servers` - (Optional, Set) List of custom name servers.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `area` - Acceleration area of the zone. Valid values: `mainland`, `overseas`.
* `cname_status` - Ownership verification status of the site when it accesses via CNAME.
* `created_on` - Site creation date.
* `modified_on` - Site modification date.
* `name_servers` - List of name servers assigned to users by Tencent Cloud.
* `original_name_servers` - List of name servers used.
* `status` - Site status.
* `vanity_name_servers_ips` - User-defined name server IP information.
  * `ipv4` - IPv4 address of the custom name server.
  * `name` - Name of the custom name server.


## Import

teo zone can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_zone.zone zone_id
```


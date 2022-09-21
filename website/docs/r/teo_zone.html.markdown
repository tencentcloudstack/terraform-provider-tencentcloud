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
  zone_name = "toutiao2.com"
  plan_type = "sta"
  type      = "full"
  paused    = false
  #  vanity_name_servers {
  #    switch = ""
  #    servers = ""
  #
  #  }
  cname_speed_up = "enabled"
  #  tags {
  #    tag_key = ""
  #    tag_value = ""
  #
  #  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `plan_type` - (Required, String) Plan type of the zone. See details in data source `zone_available_plans`.
* `zone_name` - (Required, String) Site name.
* `cname_speed_up` - (Optional, String) Specifies whether CNAME acceleration is enabled. Valid values: `enabled`, `disabled`.
* `paused` - (Optional, Bool) Indicates whether the site is disabled.
* `tags` - (Optional, Map) Tag description list.
* `type` - (Optional, String) Specifies how the site is connected to EdgeOne.- `full`: The site is connected via NS.- `partial`: The site is connected via CNAME.
* `vanity_name_servers` - (Optional, List) User-defined name server information. Note: This field may return null, indicating that no valid value can be obtained.

The `vanity_name_servers` object supports the following:

* `switch` - (Required, String) Whether to enable the custom name server.- `on`: Enable.- `off`: Disable.
* `servers` - (Optional, Set) List of custom name servers.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `area` - Acceleration area of the zone. Valid values: `mainland`, `overseas`.
* `cname_status` - Ownership verification status of the site when it accesses via CNAME.- `finished`: The site is verified.- `pending`: The site is waiting for verification.
* `created_on` - Site creation date.
* `modified_on` - Site modification date.
* `name_servers` - List of name servers assigned by Tencent Cloud.
* `original_name_servers` - Name server used by the site.
* `resources` - Billing resources of the zone.
  * `area` - Valid values: `mainland`, `overseas`.
  * `auto_renew_flag` - Whether to automatically renew. Valid values:- `0`: Default.- `1`: Enable automatic renewal.- `2`: Disable automatic renewal.
  * `create_time` - Resource creation date.
  * `enable_time` - Enable time of the resource.
  * `expire_time` - Expire time of the resource.
  * `id` - Resource ID.
  * `pay_mode` - Resource pay mode. Valid values:- `0`: post pay mode.
  * `plan_id` - Associated plan ID.
  * `status` - Status of the resource. Valid values: `normal`, `isolated`, `destroyed`.
  * `sv` - Price inquiry parameters.
    * `key` - Parameter Key.
    * `value` - Parameter Value.
* `status` - Site status. Valid values:- `active`: NS is switched.- `pending`: NS is not switched.- `moved`: NS is moved.- `deactivated`: this site is blocked.
* `vanity_name_servers_ips` - User-defined name server IP information. Note: This field may return null, indicating that no valid value can be obtained.
  * `ipv4` - IPv4 address of the custom name server.
  * `name` - Name of the custom name server.
* `zone_id` - Site ID.


## Import

teo zone can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_zone.zone zone_id
```


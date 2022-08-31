---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_application_proxy"
sidebar_current: "docs-tencentcloud-resource-teo_application_proxy"
description: |-
  Provides a resource to create a teo applicationProxy
---

# tencentcloud_teo_application_proxy

Provides a resource to create a teo applicationProxy

## Example Usage

```hcl
resource "tencentcloud_teo_application_proxy" "applicationProxy" {
  zone_id              = ""
  zone_name            = ""
  proxy_name           = ""
  plat_type            = ""
  security_type        = ""
  accelerate_type      = ""
  session_persist_time = ""
  proxy_type           = ""
}
```

## Argument Reference

The following arguments are supported:

* `accelerate_type` - (Required, Int) - 0: Disable acceleration.- 1: Enable acceleration.
* `plat_type` - (Required, String) Scheduling mode.- ip: Anycast IP.- domain: CNAME.
* `proxy_name` - (Required, String) Layer-4 proxy name.
* `security_type` - (Required, Int) - 0: Disable security protection.- 1: Enable security protection.
* `zone_id` - (Required, String) Site ID.
* `zone_name` - (Required, String) Site name.
* `proxy_type` - (Optional, String) Specifies how a layer-4 proxy is created.- hostname: Subdomain name.- instance: Instance.
* `session_persist_time` - (Optional, Int) Session persistence duration. Value range: 30-3600 (in seconds).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `host_id` - ID of the layer-7 domain name.
* `proxy_id` - Proxy ID.
* `schedule_value` - Scheduling information.
* `update_time` - Last modification date.


## Import

teo applicationProxy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_application_proxy.applicationProxy zone_id#applicationProxy_id
```


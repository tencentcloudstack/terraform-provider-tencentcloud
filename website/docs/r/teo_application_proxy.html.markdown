---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_application_proxy"
sidebar_current: "docs-tencentcloud-resource-teo_application_proxy"
description: |-
  Provides a resource to create a teo application proxy
---

# tencentcloud_teo_application_proxy

Provides a resource to create a teo application proxy

## Example Usage

```hcl
resource "tencentcloud_teo_application_proxy" "example" {
  accelerate_type      = 0
  plat_type            = "domain"
  proxy_name           = "test"
  proxy_type           = "instance"
  security_type        = 1
  session_persist_time = 0
  status               = "online"
  zone_id              = "zone-2o0l8g7zisgt"

  ipv6 {
    switch = "off"
  }
}
```

## Argument Reference

The following arguments are supported:

* `accelerate_type` - (Required, Int) - `0`: Disable acceleration.- `1`: Enable acceleration.
* `plat_type` - (Required, String) Scheduling mode.- `ip`: Anycast IP.- `domain`: CNAME.
* `proxy_name` - (Required, String) When `ProxyType` is hostname, `ProxyName` is the domain or subdomain name.When `ProxyType` is instance, `ProxyName` is the name of proxy application.
* `security_type` - (Required, Int) - `0`: Disable security protection.- `1`: Enable security protection.
* `zone_id` - (Required, String) Site ID.
* `ipv6` - (Optional, List) IPv6 access configuration.
* `proxy_type` - (Optional, String) Layer 4 proxy mode. Valid values:- `hostname`: subdomain mode.- `instance`: instance mode.
* `session_persist_time` - (Optional, Int) Session persistence duration. Value range: 30-3600 (in seconds), default value is 600.
* `status` - (Optional, String) Status of this application proxy. Valid values to set is `online` and `offline`.- `online`: Enable.- `offline`: Disable.- `progress`: Deploying.- `stopping`: Deactivating.- `fail`: Deploy or deactivate failed.

The `ipv6` object supports the following:

* `switch` - (Required, String) - `on`: Enable.- `off`: Disable.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `area` - Acceleration area. Valid values: `mainland`, `overseas`.
* `ban_status` - Application proxy block status. Valid values: `banned`, `banning`, `recover`, `recovering`.
* `host_id` - When `ProxyType` is hostname, this field is the ID of the subdomain.
* `proxy_id` - Proxy ID.
* `schedule_value` - Scheduling information.
* `update_time` - Last modification date.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateApplicationProxy](https://www.tencentcloud.com/document/api/1145/50985)
  - [DescribeApplicationProxies](https://www.tencentcloud.com/document/api/1145/50981)
  - [DescribeZones](https://www.tencentcloud.com/document/api/1145/50481)
* Read: [DescribeApplicationProxies](https://www.tencentcloud.com/document/api/1145/50981)
* Update:
  - [ModifyApplicationProxy](https://www.tencentcloud.com/document/api/1145/50980)
  - [ModifyApplicationProxyStatus](https://www.tencentcloud.com/document/api/1145/50977)
* Delete:
  - [DeleteApplicationProxy](https://www.tencentcloud.com/document/api/1145/50983)
  - [ModifyApplicationProxyStatus](https://www.tencentcloud.com/document/api/1145/50977)

## Import

teo application_proxy can be imported using the zoneId#proxyId, e.g.
```
terraform import tencentcloud_teo_application_proxy.example zone-2983wizgxqvm#proxy-6972528a-373a-11ed-afca-52540044a456
```


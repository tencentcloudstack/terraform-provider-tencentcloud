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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [UpdateIPWhitelist](https://cloud.tencent.com/document/api/850/130525) - 更新新建ip白名单
* Read: [GetIPWhitelist](https://cloud.tencent.com/document/api/850/131130) - 获取ip白名单
* Update: [UpdateIPWhitelist](https://cloud.tencent.com/document/api/850/130525) - 更新新建ip白名单
* Delete: [UpdateIPWhitelist](https://cloud.tencent.com/document/api/850/130525) - 更新新建ip白名单

## Import

Organization IP whitelist config can be imported using the zoneId, e.g.

```
terraform import tencentcloud_organization_ip_whitelist_config.example z-1os7c9znogct
```


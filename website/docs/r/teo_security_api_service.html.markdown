---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_api_service"
sidebar_current: "docs-tencentcloud-resource-teo_security_api_service"
description: |-
  Provides a resource to create a TEO API security service.
---

# tencentcloud_teo_security_api_service

Provides a resource to create a TEO API security service.

## Example Usage

```hcl
resource "tencentcloud_teo_security_api_service" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_services {
    name      = "tf-example"
    base_path = "/api/v1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_services` - (Required, List) API service configuration. Only one service is allowed per request.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `api_services` object supports the following:

* `base_path` - (Required, String) API service base path, e.g. `/tt`.
* `name` - (Required, String) API service name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO security API service can be imported using the zoneId#apiServiceId, e.g.

```
terraform import tencentcloud_teo_security_api_service.example zone-3fkff38fyw8s#apisrv-0000038524
```


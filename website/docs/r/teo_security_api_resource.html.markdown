---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_api_resource"
sidebar_current: "docs-tencentcloud-resource-teo_security_api_resource"
description: |-
  Provides a resource to create a TEO API security resource.
---

# tencentcloud_teo_security_api_resource

Provides a resource to create a TEO API security resource.

## Example Usage

```hcl
resource "tencentcloud_teo_security_api_service" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_services {
    name      = "tf-example"
    base_path = "/api/v1"
  }
}

resource "tencentcloud_teo_security_api_resource" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_resources {
    name               = "tf-example"
    path               = "/api/v1/orders"
    api_service_ids    = [tencentcloud_teo_security_api_service.example.api_services[0].id]
    methods            = ["GET", "POST"]
    request_constraint = "$${http.request.body.form['operationType']} in ['query', 'create']"
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_resources` - (Required, List) API resource configuration. Only one resource is allowed per request.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `api_resources` object supports the following:

* `name` - (Required, String) API resource name.
* `path` - (Required, String) API resource path, e.g. `/ava`.
* `api_service_ids` - (Optional, List) Associated API service ID list.
* `methods` - (Optional, List) Allowed HTTP methods. Valid values: `GET`, `POST`, `PUT`, `HEAD`, `PATCH`, `OPTIONS`, `DELETE`.
* `request_constraint` - (Optional, String) Request content matching rule expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO security API resource can be imported using the zoneId#apiResourceId, e.g.

```
terraform import tencentcloud_teo_security_api_resource.example zone-3fkff38fyw8s#apires-0000039154
```


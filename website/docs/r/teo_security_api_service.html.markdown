---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_api_service"
sidebar_current: "docs-tencentcloud-resource-teo_security_api_service"
description: |-
  Provides a resource to create a TEO security API service.
---

# tencentcloud_teo_security_api_service

Provides a resource to create a TEO security API service.

## Example Usage

```hcl
resource "tencentcloud_teo_security_api_service" "example" {
  zone_id = "zone-2qtuhspy7cr6"

  api_services {
    name      = "my-api-service"
    base_path = "/api/v1"
  }

  api_resources {
    name            = "my-api-resource"
    api_service_ids = ["svc-id-12345"]
    path            = "/api/v1/users"
    methods         = ["GET", "POST"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_services` - (Required, List, ForceNew) API service list.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `api_resources` - (Optional, List) API resource list.

The `api_resources` object supports the following:

* `api_service_ids` - (Optional, List) API service IDs associated with the API resource.
* `id` - (Optional, String) Resource ID.
* `methods` - (Optional, List) HTTP methods. Supported values: GET, POST, PUT, HEAD, PATCH, OPTIONS, DELETE.
* `name` - (Optional, String) Resource name.
* `path` - (Optional, String) Resource path.
* `request_constraint` - (Optional, String) Request content matching rule.

The `api_services` object supports the following:

* `base_path` - (Required, String) Base path of the API service.
* `name` - (Required, String) API service name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_service_ids` - API service ID list.


## Import

TEO security API service can be imported using the joint id "zone_id#api_service_ids", e.g.

```
terraform import tencentcloud_teo_security_api_service.example zone-2qtuhspy7cr6#svc-id-12345,svc-id-67890
```


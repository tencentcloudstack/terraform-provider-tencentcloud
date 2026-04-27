---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_api_resource"
sidebar_current: "docs-tencentcloud-resource-teo_security_api_resource"
description: |-
  Provides a resource to create a TEO (EdgeOne) security API resource, which is used to define API endpoints and their associated API services for security protection.
---

# tencentcloud_teo_security_api_resource

Provides a resource to create a TEO (EdgeOne) security API resource, which is used to define API endpoints and their associated API services for security protection.

## Example Usage

```hcl
resource "tencentcloud_teo_security_api_resource" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  api_resources {
    name               = "test-api-resource"
    api_service_ids    = ["svc-123"]
    path               = "/api/v1/test"
    methods            = ["GET", "POST"]
    request_constraint = jsonencode({ "key" : "value" })
  }
  api_resources {
    name    = "test-api-resource-2"
    path    = "/api/v2/test"
    methods = ["GET"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_resources` - (Required, List) API resource list.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `api_resources` object supports the following:

* `name` - (Required, String) API resource name.
* `api_service_ids` - (Optional, List) API service IDs associated with the API resource.
* `methods` - (Optional, List) Request method list. Supported values: GET, POST, PUT, HEAD, PATCH, OPTIONS, DELETE.
* `path` - (Optional, String) Resource path.
* `request_constraint` - (Optional, String) Request content matching rule, must conform to expression syntax.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_resource_ids` - API resource IDs returned by server after creation.


## Import

teo security_api_resource can be imported using the id, e.g.

```
terraform import tencentcloud_teo_security_api_resource.example zone_id
```


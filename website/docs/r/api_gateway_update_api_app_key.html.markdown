---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_update_api_app_key"
sidebar_current: "docs-tencentcloud-resource-api_gateway_update_api_app_key"
description: |-
  Provides a resource to create a apiGateway update_api_app_key
---

# tencentcloud_api_gateway_update_api_app_key

Provides a resource to create a apiGateway update_api_app_key

## Example Usage

```hcl
resource "tencentcloud_api_gateway_update_api_app_key" "example" {
  api_app_id  = "app-krljp4wn"
  api_app_key = "APID6JmG21yRCc03h4z16hlsTqj1wpO3dB3ZQcUP"
}
```

## Argument Reference

The following arguments are supported:

* `api_app_id` - (Required, String, ForceNew) Application unique ID.
* `api_app_key` - (Required, String, ForceNew) Key of the application.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.




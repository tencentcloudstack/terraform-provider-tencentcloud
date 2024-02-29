---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_update_service"
sidebar_current: "docs-tencentcloud-resource-api_gateway_update_service"
description: |-
  Provides a resource to create a apigateway update_service
---

# tencentcloud_api_gateway_update_service

Provides a resource to create a apigateway update_service

## Example Usage

```hcl
resource "tencentcloud_api_gateway_update_service" "example" {
  service_id       = "service-oczq2nyk"
  environment_name = "test"
  version_name     = "20240204142759-b5a4f741-adc0-4964-b01b-2a4a04ff6964"
}
```

## Argument Reference

The following arguments are supported:

* `environment_name` - (Required, String, ForceNew) The name of the environment to be switched, currently supporting three environments: test (test environment), prepub (pre release environment), and release (release environment).
* `service_id` - (Required, String, ForceNew) Service ID.
* `version_name` - (Required, String, ForceNew) The version number of the switch.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.




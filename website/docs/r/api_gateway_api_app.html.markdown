---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_app"
sidebar_current: "docs-tencentcloud-resource-api_gateway_api_app"
description: |-
  Provides a resource to create a APIGateway ApiApp
---

# tencentcloud_api_gateway_api_app

Provides a resource to create a APIGateway ApiApp

## Example Usage

```hcl
resource "tencentcloud_api_gateway_api_app" "my_api_app" {
  api_app_name = "app_test1"
  api_app_desc = "app desc."
}
```

## Argument Reference

The following arguments are supported:

* `api_app_name` - (Required, String) Api app name.
* `api_app_desc` - (Optional, String) App description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_app_id` - Api app ID.
* `api_app_key` - Api app key.
* `api_app_secret` - Api app secret.
* `created_time` - Api app created time.
* `modified_time` - Api app modified time.



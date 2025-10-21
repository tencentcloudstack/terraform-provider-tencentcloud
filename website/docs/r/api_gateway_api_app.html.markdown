---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_app"
sidebar_current: "docs-tencentcloud-resource-api_gateway_api_app"
description: |-
  Provides a resource to create a APIGateway ApiApp
---

# tencentcloud_api_gateway_api_app

Provides a resource to create a APIGateway ApiApp

## Example Usage

### Create a basic apigateway api_app

```hcl
resource "tencentcloud_api_gateway_api_app" "example" {
  api_app_name = "tf_example"
  api_app_desc = "app desc."
}
```

### Bind Tag

```hcl
resource "tencentcloud_api_gateway_api_app" "example" {
  api_app_name = "tf_example"
  api_app_desc = "app desc."

  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_app_name` - (Required, String) Api app name.
* `api_app_desc` - (Optional, String) App description.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_app_id` - Api app ID.
* `api_app_key` - Api app key.
* `api_app_secret` - Api app secret.
* `created_time` - Api app created time.
* `modified_time` - Api app modified time.


## Import

apigateway api_app can be imported using the id, e.g.

```
terraform import tencentcloud_api_gateway_api_app.example app-poe0pyex
```


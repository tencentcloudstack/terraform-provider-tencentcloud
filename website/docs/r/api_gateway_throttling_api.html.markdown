---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_throttling_api"
sidebar_current: "docs-tencentcloud-resource-api_gateway_throttling_api"
description: |-
  Use this resource to create API gateway throttling API.
---

# tencentcloud_api_gateway_throttling_api

Use this resource to create API gateway throttling API.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id            = tencentcloud_api_gateway_service.service.id
  api_name              = "hello_update"
  api_desc              = "my hello api update"
  auth_type             = "SECRET"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "POST"
  request_parameters {
    name          = "email"
    position      = "QUERY"
    type          = "string"
    desc          = "your email please?"
    default_value = "tom@qq.com"
    required      = true
  }
  service_config_type      = "HTTP"
  service_config_timeout   = 10
  service_config_url       = "http://www.tencent.com"
  service_config_path      = "/user"
  service_config_method    = "POST"
  response_type            = "XML"
  response_success_example = "<note>success</note>"
  response_fail_example    = "<note>fail</note>"
  response_error_codes {
    code           = 10
    msg            = "system error"
    desc           = "system error code"
    converted_code = -10
    need_convert   = true
  }
}

resource "tencentcloud_api_gateway_throttling_api" "foo" {
  service_id       = tencentcloud_api_gateway_service.service.id
  strategy         = "400"
  environment_name = "test"
  api_ids          = [tencentcloud_api_gateway_api.api.id]
}
```

## Argument Reference

The following arguments are supported:

* `api_ids` - (Required) List of API ID.
* `environment_name` - (Required) List of Environment names.
* `service_id` - (Required, ForceNew) Service ID for query.
* `strategy` - (Required) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_environment_strategies` - List of throttling policies bound to API.
  * `api_id` - Unique API ID.
  * `api_name` - Custom API name.
  * `method` - API method.
  * `path` - API path.
  * `strategy_list` - Environment throttling information.
    * `environment_name` - Environment name.
    * `quota` - Throttling value.



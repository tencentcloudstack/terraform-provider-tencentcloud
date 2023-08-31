---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_throttling_apis"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_throttling_apis"
description: |-
  Use this data source to query API gateway throttling APIs.
---

# tencentcloud_api_gateway_throttling_apis

Use this data source to query API gateway throttling APIs.

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
  api_name              = "tf_example"
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

  release_limit = 100
  pre_limit     = 100
  test_limit    = 100
}

data "tencentcloud_api_gateway_throttling_apis" "id" {
  service_id = tencentcloud_api_gateway_api.service_id
}

data "tencentcloud_api_gateway_throttling_apis" "foo" {
  service_id        = tencentcloud_api_gateway_api.service.service_id
  environment_names = ["release", "test"]
}
```

## Argument Reference

The following arguments are supported:

* `environment_names` - (Optional, List: [`String`]) Environment list.
* `result_output_file` - (Optional, String) Used to save results.
* `service_id` - (Optional, String) Unique service ID of API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of policies bound to API.
  * `api_environment_strategies` - List of throttling policies bound to API.
    * `api_id` - Unique API ID.
    * `api_name` - Custom API name.
    * `method` - API method.
    * `path` - API path.
    * `strategy_list` - Environment throttling information.
      * `environment_name` - Environment name.
      * `quota` - Throttling value.
  * `service_id` - Unique service ID of API.



---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_plugins"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_plugins"
description: |-
  Use this data source to query detailed information of apigateway plugin
---

# tencentcloud_api_gateway_plugins

Use this data source to query detailed information of apigateway plugin

## Example Usage

```hcl
data "tencentcloud_api_gateway_plugins" "example" {
  service_id       = tencentcloud_api_gateway_service_release.example.service_id
  plugin_id        = tencentcloud_api_gateway_plugin.example.id
  environment_name = "release"
}

resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf_example"
  protocol     = "http&https"
  service_desc = "desc."
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
  tags = {
    testKey = "testValue"
  }
  release_limit = 500
  pre_limit     = 500
  test_limit    = 500
}

resource "tencentcloud_api_gateway_api" "example" {
  service_id            = tencentcloud_api_gateway_service.example.id
  api_name              = "hello"
  api_desc              = "my hello api"
  auth_type             = "NONE"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "GET"

  request_parameters {
    name          = "name"
    position      = "QUERY"
    type          = "string"
    desc          = "who are you?"
    default_value = "tom"
    required      = true
  }
  service_config_type      = "HTTP"
  service_config_timeout   = 15
  service_config_url       = "http://www.qq.com"
  service_config_path      = "/user"
  service_config_method    = "GET"
  response_type            = "HTML"
  response_success_example = "success"
  response_fail_example    = "fail"
  response_error_codes {
    code           = 500
    msg            = "system error"
    desc           = "system error code"
    converted_code = 5000
    need_convert   = true
  }
}

resource "tencentcloud_api_gateway_service_release" "example" {
  service_id       = tencentcloud_api_gateway_api.example.service_id
  environment_name = "release"
  release_desc     = "desc."
}

resource "tencentcloud_api_gateway_plugin" "example" {
  plugin_name = "tf-example"
  plugin_type = "IPControl"
  plugin_data = jsonencode({
    "type" : "white_list",
    "blocks" : "1.1.1.1",
  })
  description = "desc."
}
```

## Argument Reference

The following arguments are supported:

* `environment_name` - (Required, String) Environmental information.
* `plugin_id` - (Required, String) The plugin ID to query.
* `service_id` - (Required, String) The service ID to query.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - List of plugin related APIs.
  * `api_id` - API ID.
  * `api_name` - API name.
  * `api_type` - API type.
  * `attached_other_plugin` - Whether the API is bound to other plugins.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `is_attached` - Whether the API is bound to the current plugin.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `method` - API method.
  * `path` - API path.



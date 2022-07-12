---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api"
sidebar_current: "docs-tencentcloud-resource-api_gateway_api"
description: |-
  Use this resource to create API of API gateway.
---

# tencentcloud_api_gateway_api

Use this resource to create API of API gateway.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id            = tencentcloud_api_gateway_service.service.id
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
    code           = 100
    msg            = "system error"
    desc           = "system error code"
    converted_code = -100
    need_convert   = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_name` - (Required, String) Custom API name.
* `request_config_path` - (Required, String) Request frontend path configuration. Like `/user/getinfo`.
* `service_id` - (Required, String, ForceNew) Which service this API belongs. Refer to resource `tencentcloud_api_gateway_service`.
* `api_desc` - (Optional, String) Custom API description.
* `auth_type` - (Optional, String) API authentication type. Valid values: `SECRET` (key pair authentication),`NONE` (no authentication). Default value: `NONE`.
* `enable_cors` - (Optional, Bool) Whether to enable CORS. Default value: `true`.
* `pre_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `protocol` - (Optional, String, ForceNew) API frontend request type. Valid values: `HTTP`, `WEBSOCKET`. Default value: `HTTP`.
* `release_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `request_config_method` - (Optional, String) Request frontend method configuration. Valid values: `GET`,`POST`,`PUT`,`DELETE`,`HEAD`,`ANY`. Default value: `GET`.
* `request_parameters` - (Optional, Set) Frontend request parameters.
* `response_error_codes` - (Optional, Set) Custom error code configuration. Must keep at least one after set.
* `response_fail_example` - (Optional, String) Response failure sample of custom response configuration.
* `response_success_example` - (Optional, String) Successful response sample of custom response configuration.
* `response_type` - (Optional, String) Return type. Valid values: `HTML`, `JSON`, `TEXT`, `BINARY`, `XML`. Default value: `HTML`.
* `service_config_method` - (Optional, String) API backend service request method, such as `GET`. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_method` and backend method `service_config_method` can be different.
* `service_config_mock_return_message` - (Optional, String) Returned information of API backend mocking. This parameter is required when `service_config_type` is `MOCK`.
* `service_config_path` - (Optional, String) API backend service path, such as /path. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_path` and backend path `service_config_path` can be different.
* `service_config_product` - (Optional, String) Backend type. This parameter takes effect when VPC is enabled. Currently, only `clb` is supported.
* `service_config_scf_function_name` - (Optional, String) SCF function name. This parameter takes effect when `service_config_type` is `SCF`.
* `service_config_scf_function_namespace` - (Optional, String) SCF function namespace. This parameter takes effect when `service_config_type` is `SCF`.
* `service_config_scf_function_qualifier` - (Optional, String) SCF function version. This parameter takes effect when `service_config_type` is `SCF`.
* `service_config_timeout` - (Optional, Int) API backend service timeout period in seconds. Default value: `5`.
* `service_config_type` - (Optional, String) API backend service type. Valid values: `WEBSOCKET`, `HTTP`, `SCF`, `MOCK`. Default value: `HTTP`.
* `service_config_url` - (Optional, String) API backend service url. This parameter is required when `service_config_type` is `HTTP`.
* `service_config_vpc_id` - (Optional, String) Unique VPC ID.
* `test_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.

The `request_parameters` object supports the following:

* `name` - (Required, String) Parameter name.
* `position` - (Required, String) Parameter location.
* `type` - (Required, String) Parameter type.
* `default_value` - (Optional, String) Parameter default value.
* `desc` - (Optional, String) Parameter description.
* `required` - (Optional, Bool) If this parameter required. Default value: `false`.

The `response_error_codes` object supports the following:

* `code` - (Required, Int) Custom response configuration error code.
* `msg` - (Required, String) Custom response configuration error message.
* `converted_code` - (Optional, Int) Custom error code conversion.
* `desc` - (Optional, String) Parameter description.
* `need_convert` - (Optional, Bool) Whether to enable error code conversion. Default value: `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `update_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.



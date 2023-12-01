---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_usage_plans"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_api_usage_plans"
description: |-
  Use this data source to query detailed information of apigateway api_usage_plan
---

# tencentcloud_api_gateway_api_usage_plans

Use this data source to query detailed information of apigateway api_usage_plan

## Example Usage

```hcl
data "tencentcloud_api_gateway_api_usage_plans" "example" {
  service_id = tencentcloud_api_gateway_usage_plan_attachment.example.service_id
}

resource "tencentcloud_api_gateway_usage_plan" "example" {
  usage_plan_name         = "tf_example"
  usage_plan_desc         = "desc."
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf_example"
  protocol     = "http&https"
  service_desc = "desc."
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "example" {
  service_id            = tencentcloud_api_gateway_service.example.id
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
    desc          = "desc."
    default_value = "test@qq.com"
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
    code           = 500
    msg            = "system error"
    desc           = "system error code"
    converted_code = 5000
    need_convert   = true
  }
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "example" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.example.id
  service_id    = tencentcloud_api_gateway_service.example.id
  environment   = "release"
  bind_type     = "API"
  api_id        = tencentcloud_api_gateway_api.example.id
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String) The unique ID of the service to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - API binding usage plan list.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `api_id` - API unique ID.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `api_name` - API name.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `created_time` - Create a time using a schedule.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `environment` - Use the service environment bound by the plan.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `in_use_request_num` - The quota that has already been used.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `max_request_num_pre_sec` - Request QPS upper limit, -1 indicates no limit.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `max_request_num` - Request total quota, -1 indicates no limit.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `method` - API method.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `modified_time` - Use the last modification time of the plan.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `path` - API path.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `service_id` - Service unique ID.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `service_name` - Service name.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `usage_plan_desc` - Description of the usage plan.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `usage_plan_id` - Use the unique ID of the plan.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `usage_plan_name` - Use the name of the plan.Note: This field may return null, indicating that a valid value cannot be obtained.



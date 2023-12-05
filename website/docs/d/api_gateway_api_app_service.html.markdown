---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_app_service"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_api_app_service"
description: |-
  Use this data source to query detailed information of apigateway api_app_services
---

# tencentcloud_api_gateway_api_app_service

Use this data source to query detailed information of apigateway api_app_services

## Example Usage

```hcl
data "tencentcloud_api_gateway_api_app_service" "example" {
  service_id = tencentcloud_api_gateway_api.example.service_id
  api_region = "ap-guangzhou"
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
  auth_type             = "APP"
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
```

## Argument Reference

The following arguments are supported:

* `api_region` - (Required, String) Territory to which the service belongs.
* `service_id` - (Required, String) The unique ID of the service to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_id_status_set` - API list.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `api_business_type` - API business type.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `api_desc` - API DescriptionNote: This field may return null, indicating that a valid value cannot be obtained.
  * `api_id` - API unique ID.
  * `api_name` - API name.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `api_type` - API type.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `auth_relation_api_id` - Unique ID of the association authorization API.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `auth_type` - Authorization type.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `created_time` - Service creation time.
  * `is_debug_after_charge` - Whether to debug after purchase.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `method` - API METHOD.
  * `modified_time` - Service modification time.
  * `oauth_config` - OAuth configuration information.Note: This field may return null, indicating that a valid value cannot be obtained.
    * `login_redirect_url` - Redirect address, used to guide users in login operations.
    * `public_key` - Public key, used to verify user tokens.
    * `token_location` - Token passes the position.
  * `path` - API PATH.
  * `protocol` - API protocol.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `service_id` - Service unique ID.
  * `token_location` - OAuth2.0 API request, token storage location.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `uniq_vpc_id` - VPC unique ID.Note: This field may return null, indicating that a valid value cannot be obtained.
* `api_total_count` - Total number of APIs.Note: This field may return null, indicating that a valid value cannot be obtained.
* `available_environments` - List of service environments.Note: This field may return null, indicating that a valid value cannot be obtained.
* `created_time` - Service creation time.Note: This field may return null, indicating that a valid value cannot be obtained.
* `inner_http_port` - Internal network access HTTP service port number.
* `inner_https_port` - Internal network access https port number.
* `internal_sub_domain` - Intranet access sub domain name.
* `ip_version` - IP version.Note: This field may return null, indicating that a valid value cannot be obtained.
* `modified_time` - Service modification time.Note: This field may return null, indicating that a valid value cannot be obtained.
* `net_types` - A list of network types, where INNER represents internal network access and OUTER represents external network access.
* `outer_sub_domain` - External network access sub domain name.
* `protocol` - Service support protocol, optional values are http, https, and http&amp;amp;https.
* `service_desc` - Service description.Note: This field may return null, indicating that a valid value cannot be obtained.
* `service_name` - Service name.Note: This field may return null, indicating that a valid value cannot be obtained.
* `set_id` - Reserved fields.Note: This field may return null, indicating that a valid value cannot be obtained.
* `usage_plan_list` - Use a plan array.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `created_time` - Use planned time.
  * `environment` - Environment name.
  * `max_request_num_pre_sec` - Use plan qps, -1 indicates no restrictions.
  * `modified_time` - Use the schedule to modify the time.
  * `usage_plan_desc` - Use plan description.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `usage_plan_id` - Use a unique ID for the plan.
  * `usage_plan_name` - Use the plan name.
* `usage_plan_total_count` - Total number of usage plans.Note: This field may return null, indicating that a valid value cannot be obtained.
* `user_type` - The user type of this service.Note: This field may return null, indicating that a valid value cannot be obtained.



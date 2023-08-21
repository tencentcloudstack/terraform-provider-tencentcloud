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
* `service_id` - (Required, String, ForceNew) The unique ID of the service where the API is located. Refer to resource `tencentcloud_api_gateway_service`.
* `api_business_type` - (Optional, String) When `auth_type` is OAUTH, this field is valid, NORMAL: Business API, OAUTH: Authorization API.
* `api_desc` - (Optional, String) Custom API description.
* `api_type` - (Optional, String) API type, supports NORMAL (regular API) and TSF (microservice API), defaults to NORMAL.
* `auth_relation_api_id` - (Optional, String) The unique ID of the associated authorization API takes effect when AuthType is OAUTH and ApiBusinessType is NORMAL. The unique ID of the oauth2.0 authorized API that identifies the business API binding.
* `auth_type` - (Optional, String) API authentication type. Support SECRET (Key Pair Authentication), NONE (Authentication Exemption), OAUTH, APP (Application Authentication). The default is NONE.
* `constant_parameters` - (Optional, Set) Constant parameter.
* `eiam_app_id` - (Optional, String) EIAM application ID.
* `eiam_app_type` - (Optional, String) EIAM application type.
* `eiam_auth_type` - (Optional, String) The EIAM application authentication type supports AuthenticationOnly, Authentication, and Authorization.
* `enable_cors` - (Optional, Bool) Whether to enable CORS. Default value: `true`.
* `event_bus_id` - (Optional, String) Event bus ID.
* `is_base64_encoded` - (Optional, Bool) Whether to enable Base64 encoding will only take effect when the backend is scf.
* `is_debug_after_charge` - (Optional, Bool) Charge after starting debugging. (Cloud Market Reserved Fields).
* `is_delete_response_error_codes` - (Optional, Bool) Do you want to delete the custom response configuration error code? If it is not passed or False is passed, it will not be deleted. If True is passed, all custom response configuration error codes for this API will be deleted.
* `micro_services` - (Optional, Set) API bound microservice list.
* `oauth_config` - (Optional, List) OAuth configuration. Effective when AuthType is OAUTH.
* `owner` - (Optional, String) Owner of resources.
* `pre_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `protocol` - (Optional, String, ForceNew) API frontend request type. Valid values: `HTTP`, `WEBSOCKET`. Default value: `HTTP`.
* `release_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `request_config_method` - (Optional, String) Request frontend method configuration. Valid values: `GET`,`POST`,`PUT`,`DELETE`,`HEAD`,`ANY`. Default value: `GET`.
* `request_parameters` - (Optional, Set) Frontend request parameters.
* `response_error_codes` - (Optional, Set) Custom error code configuration. Must keep at least one after set.
* `response_fail_example` - (Optional, String) Response failure sample of custom response configuration.
* `response_success_example` - (Optional, String) Successful response sample of custom response configuration.
* `response_type` - (Optional, String) Return type. Valid values: `HTML`, `JSON`, `TEXT`, `BINARY`, `XML`. Default value: `HTML`.
* `service_config_cos_config` - (Optional, List) API backend COS configuration. If ServiceType is COS, then this parameter must be passed.Note: This field may return null, indicating that a valid value cannot be obtained.
* `service_config_method` - (Optional, String) API backend service request method, such as `GET`. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_method` and backend method `service_config_method` can be different.
* `service_config_mock_return_message` - (Optional, String) Returned information of API backend mocking. This parameter is required when `service_config_type` is `MOCK`.
* `service_config_path` - (Optional, String) API backend service path, such as /path. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_path` and backend path `service_config_path` can be different.
* `service_config_product` - (Optional, String) Backend type. Effective when enabling vpc, currently supported types are clb, cvm, and upstream.
* `service_config_scf_function_name` - (Optional, String) SCF function name. This parameter takes effect when `service_config_type` is `SCF`.
* `service_config_scf_function_namespace` - (Optional, String) SCF function namespace. This parameter takes effect when `service_config_type` is `SCF`.
* `service_config_scf_function_qualifier` - (Optional, String) SCF function version. This parameter takes effect when `service_config_type` is `SCF`.
* `service_config_scf_function_type` - (Optional, String) Scf function type. Effective when the backend type is SCF. Support Event Triggering (EVENT) and HTTP Direct Cloud Function (HTTP).
* `service_config_scf_is_integrated_response` - (Optional, Bool) Whether to enable response integration. Effective when the backend type is SCF.
* `service_config_timeout` - (Optional, Int) API backend service timeout period in seconds. Default value: `5`.
* `service_config_type` - (Optional, String) The backend service type of the API. Supports HTTP, MOCK, TSF, SCF, WEBSOCKET, COS, TARGET (internal testing).
* `service_config_upstream_id` - (Optional, String) Only required when binding to VPC channelsNote: This field may return null, indicating that a valid value cannot be obtained.
* `service_config_url` - (Optional, String) The backend service URL of the API. If the ServiceType is HTTP, this parameter must be passed.
* `service_config_vpc_id` - (Optional, String) Unique VPC ID.
* `service_config_websocket_cleanup_function_name` - (Optional, String) Scf websocket cleaning function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_cleanup_function_namespace` - (Optional, String) Scf websocket cleans up the function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_cleanup_function_qualifier` - (Optional, String) Scf websocket cleaning function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_register_function_name` - (Optional, String) Scf websocket registration function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_register_function_namespace` - (Optional, String) Scf websocket registers function namespaces. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_register_function_qualifier` - (Optional, String) Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_transport_function_name` - (Optional, String) Scf websocket transfer function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_transport_function_namespace` - (Optional, String) Scf websocket transfer function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_transport_function_qualifier` - (Optional, String) Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_parameters` - (Optional, List) The backend service parameters of the API.
* `service_tsf_health_check_conf` - (Optional, List) Health check configuration for microservices.
* `service_tsf_load_balance_conf` - (Optional, List) Load balancing configuration for microservices.
* `tags` - (Optional, Map) Tag description list.
* `target_namespace_id` - (Optional, String) Tsf serverless namespace ID. (In internal testing).
* `target_services_health_check_conf` - (Optional, List) Target health check configuration. (Internal testing stage).
* `target_services_load_balance_conf` - (Optional, Int) Target type load balancing configuration. (Internal testing stage).
* `target_services` - (Optional, List) Target type backend resource information. (Internal testing stage).
* `test_limit` - (Optional, Int) API QPS value. Enter a positive number to limit the API query rate per second `QPS`.
* `token_timeout` - (Optional, Int) The effective time of the EIAM application token, measured in seconds, defaults to 7200 seconds.
* `user_type` - (Optional, String) User type.

The `constant_parameters` object supports the following:

* `default_value` - (Optional, String) Default value for constant parameters. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
* `desc` - (Optional, String) Constant parameter description. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
* `name` - (Optional, String) Constant parameter name. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
* `position` - (Optional, String) Constant parameter position. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.

The `micro_services` object supports the following:

* `cluster_id` - (Required, String) Micro service cluster.
* `micro_service_name` - (Required, String) Microservice name.
* `namespace_id` - (Required, String) Microservice namespace.

The `oauth_config` object supports the following:

* `public_key` - (Required, String) Public key, used to verify user tokens.
* `token_location` - (Required, String) Token passes the position.
* `login_redirect_url` - (Optional, String) Redirect address, used to guide users in login operations.

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

The `service_config_cos_config` object supports the following:

* `action` - (Required, String) The API calls the backend COS method, and the optional values for the front-end request method and Action are:GET: GetObjectPUT: PutObjectPOST: PostObject, AppendObjectHEAD: HeadObjectDELETE: DeleteObject.Note: This field may return null, indicating that a valid value cannot be obtained.
* `bucket_name` - (Required, String) The bucket name of the API backend COS.Note: This field may return null, indicating that a valid value cannot be obtained.
* `authorization` - (Optional, Bool) The API calls the signature switch of the backend COS, which defaults to false.Note: This field may return null, indicating that a valid value cannot be obtained.
* `path_match_mode` - (Optional, String) Path matching mode for API backend COS, optional values:BackEndPath: Backend path matchingFullPath: Full Path MatchingThe default value is: BackEndPathNote: This field may return null, indicating that a valid value cannot be obtained.

The `service_parameters` object supports the following:

* `default_value` - (Optional, String) The default value for the backend service parameters of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
* `name` - (Optional, String) The backend service parameter name of the API. This parameter is only used when ServiceType is HTTP. The front and rear parameter names can be different.Note: This field may return null, indicating that a valid value cannot be obtained.
* `position` - (Optional, String) The backend service parameter location of the API, such as head. This parameter is only used when ServiceType is HTTP. The parameter positions at the front and rear ends can be configured differently.Note: This field may return null, indicating that a valid value cannot be obtained.
* `relevant_request_parameter_desc` - (Optional, String) Remarks on the backend service parameters of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
* `relevant_request_parameter_name` - (Optional, String) The name of the front-end parameter corresponding to the backend service parameter of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
* `relevant_request_parameter_position` - (Optional, String) The location of the front-end parameters corresponding to the backend service parameters of the API, such as head. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
* `relevant_request_parameter_type` - (Optional, String) The backend service parameter type of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.

The `service_tsf_health_check_conf` object supports the following:

* `error_threshold_percentage` - (Optional, Int) Threshold percentage.Note: This field may return null, indicating that a valid value cannot be obtained.
* `is_health_check` - (Optional, Bool) Whether to initiate a health check.Note: This field may return null, indicating that a valid value cannot be obtained.
* `request_volume_threshold` - (Optional, Int) Health check threshold.Note: This field may return null, indicating that a valid value cannot be obtained.
* `sleep_window_in_milliseconds` - (Optional, Int) Window size.Note: This field may return null, indicating that a valid value cannot be obtained.

The `service_tsf_load_balance_conf` object supports the following:

* `is_load_balance` - (Optional, Bool) Is load balancing enabled.Note: This field may return null, indicating that a valid value cannot be obtained.
* `method` - (Optional, String) Load balancing method.Note: This field may return null, indicating that a valid value cannot be obtained.
* `session_stick_required` - (Optional, Bool) Whether to enable session persistence.Note: This field may return null, indicating that a valid value cannot be obtained.
* `session_stick_timeout` - (Optional, Int) Session hold timeout.Note: This field may return null, indicating that a valid value cannot be obtained.

The `target_services_health_check_conf` object supports the following:

* `error_threshold_percentage` - (Optional, Int) Threshold percentage.Note: This field may return null, indicating that a valid value cannot be obtained.
* `is_health_check` - (Optional, Bool) Whether to initiate a health check.Note: This field may return null, indicating that a valid value cannot be obtained.
* `request_volume_threshold` - (Optional, Int) Health check threshold.Note: This field may return null, indicating that a valid value cannot be obtained.
* `sleep_window_in_milliseconds` - (Optional, Int) Window size.Note: This field may return null, indicating that a valid value cannot be obtained.

The `target_services` object supports the following:

* `host_ip` - (Required, String) Host IP of the CVM.
* `vm_ip` - (Required, String) vm ip.
* `vm_port` - (Required, Int) vm port.
* `vpc_id` - (Required, String) vpc id.
* `docker_ip` - (Optional, String) docker ip.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `update_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.



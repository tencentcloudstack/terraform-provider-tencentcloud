---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_app_api"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_api_app_api"
description: |-
  Use this data source to query detailed information of apiGateway api_app_api
---

# tencentcloud_api_gateway_api_app_api

Use this data source to query detailed information of apiGateway api_app_api

## Example Usage

```hcl
data "tencentcloud_api_gateway_api_app_api" "example" {
  service_id = "service-nxz6yync"
  api_id     = "api-0cvmf4x4"
  api_region = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required, String) API interface unique ID.
* `api_region` - (Required, String) Api region.
* `service_id` - (Required, String) The unique ID of the service where the API resides.
* `result_output_file` - (Optional, String) Used to save apiAppApis.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - API details.
  * `api_business_type` - Type of OAUTH API. Possible values are NORMAL (Business API), OAUTH (Authorization API).
  * `api_desc` - Description of the API interface.
  * `api_id` - API interface unique ID.
  * `api_name` - The name of the API interface.
  * `api_type` - API type. Possible values are NORMAL (normal API) and TSF (microservice API).
  * `auth_relation_api_id` - OAUTH The unique ID of the authorization API associated with the business API.
  * `auth_type` - API authentication type. Possible values are SECRET (key pair authentication), NONE (authentication-free), and OAUTH.
  * `base64_encoded_trigger_rules` - Header triggers rules, and the total number of rules does not exceed 10.
    * `name` - Header for encoding triggering, optional values Accept and Content_Type correspond to Accept and Content-Type in the actual data flow request header.
    * `value` - An array of optional values for the header triggered by encoding. The maximum string length of the array element is 40. The elements can include numbers, English letters and special characters. The optional values for special characters are: `.` `+` ` *` `-` `/` `_` For example [ application/x-vpeg005, application/xhtml+xml, application/vnd.ms -project, application/vnd.rn-rn_music_package] etc. are all legal.
  * `constant_parameters` - Constant parameters.
    * `default_value` - Constant parameter default value. This parameter is only used if the ServiceType is HTTP.
    * `desc` - Constant parameter description. This parameter is only used if the ServiceType is HTTP.
    * `name` - Constant parameter name. This parameter is only used if the ServiceType is HTTP.
    * `position` - Constant parameter position. This parameter is only used if the ServiceType is HTTP.
  * `created_time` - Creation time, expressed in accordance with the ISO8601 standard and using UTC time. The format is: YYYY-MM-DDThh:mm:ssZ.
  * `enable_cors` - Whether to enable cross-domain.
  * `environments` - API published environment information.
  * `internal_domain` - WEBSOCKET pushback address.
  * `is_base64_encoded` - Whether to enable Base64 encoding will only take effect when the backend is scf.
  * `is_base64_trigger` - Whether to enable Base64-encoded header triggering will only take effect when the backend is scf.
  * `is_debug_after_charge` - Whether to debug after purchase (parameters reserved in the cloud market).
  * `micro_services_info` - Microservice information details.
  * `micro_services` - API binding microservice list.
    * `cluster_id` - Microservice cluster ID.
    * `micro_service_name` - Microservice name.
    * `namespace_id` - Microservice namespace ID.
  * `modified_time` - Last modification time, expressed in accordance with the ISO8601 standard and using UTC time. The format is: YYYY-MM-DDThh:mm:ssZ.
  * `oauth_config` - OAUTH configuration.
    * `login_redirect_url` - Redirect address, used to guide users to log in.
    * `public_key` - Public key, used to verify user token.
    * `token_location` - Token delivery position.
  * `protocol` - The front-end request type of the API, such as HTTP or HTTPS or HTTP and HTTPS.
  * `request_config` - The requested frontend configuration.
    * `method` - API request method, such as GET.
    * `path` - API path, such as /path.
  * `request_parameters` - Front-end request parameters.
    * `default_value` - API front-end parameter default value.
    * `desc` - API front-end parameter remarks.
    * `name` - API front-end parameter name.
    * `position` - The front-end parameter position of the API, such as header. Currently supports header, query, path.
    * `required` - .
    * `type` - API front-end parameter type, such as String, int.
  * `response_error_codes` - User-defined error code configuration.
    * `code` - Custom response configuration error code.
    * `converted_code` - Custom error code conversion.
    * `desc` - Custom response configuration error code remarks.
    * `msg` - Custom response configuration error message.
    * `need_convert` - Whether it is necessary to enable error code conversion.
  * `response_fail_example` - Custom response configuration failure response example.
  * `response_success_example` - Custom response configuration successful response example.
  * `response_type` - Return type.
  * `service_config` - Backend service configuration for the API.
    * `method` - API backend service request method, such as GET. If ServiceType is HTTP, this parameter is required. The front-end and back-end methods can be different.
    * `path` - API backend service path, such as /path. If ServiceType is HTTP, this parameter is required. The front-end and back-end paths can be different.
    * `product` - Backend type. It takes effect when vpc is enabled. Currently supported types are clb, cvm and upstream.
    * `uniq_vpc_id` - The unique ID of the vpc.
    * `upstream_id` - Only required when binding vpc channel.
    * `url` - API&amp;#39;s backend service url. If ServiceType is HTTP, this parameter must be passed.
  * `service_desc` - A description of the service where the API resides.
  * `service_id` - The unique ID of the service where the API resides.
  * `service_mock_return_message` - APIs backend Mock returns information. If ServiceType is Mock, this parameter must be passed.
  * `service_name` - The name of the service where the API resides.
  * `service_parameters` - API backend service parameters.
    * `default_value` - Default values for the APIs backend service parameters. This parameter is only used if the ServiceType is HTTP.
    * `name` - The backend service parameter name of the API. This parameter will be used only if the ServiceType is HTTP. The front-end and back-end parameter names can be different.
    * `position` - The backend service parameter location of the API, such as head. This parameter is only used if the ServiceType is HTTP. The front-end and back-end parameter positions can be configured differently.
    * `relevant_request_parameter_desc` - Remarks on the backend service parameters of the API. This parameter is only used if the ServiceType is HTTP.
    * `relevant_request_parameter_name` - The front-end parameter name corresponding to the back-end service parameter of the API. This parameter is only used if the ServiceType is HTTP.
    * `relevant_request_parameter_position` - The front-end parameter position corresponding to the back-end service parameter of the API, such as head. This parameter is only used if the ServiceType is HTTP.
  * `service_scf_function_name` - Scf function name. Effective when the backend type is SCF.
  * `service_scf_function_namespace` - Scf function namespace. Effective when the backend type is SCF.
  * `service_scf_function_qualifier` - Scf function version. Effective when the backend type is SCF.
  * `service_scf_is_integrated_response` - Whether to enable integrated response.
  * `service_timeout` - The backend service timeout of the API, in seconds.
  * `service_tsf_health_check_conf` - Health check configuration for microservices.
    * `error_threshold_percentage` - Threshold percentage.
    * `is_health_check` - Whether to enable health check.
    * `request_volume_threshold` - Health check threshold.
    * `sleep_window_in_milliseconds` - Window size.
  * `service_tsf_load_balance_conf` - Load balancing configuration for microservices.
    * `is_load_balance` - Whether to enable load balancing.
    * `method` - Load balancing method.
    * `session_stick_required` - Whether to enable session persistence.
    * `session_stick_timeout` - Session retention timeout.
  * `service_type` - The backend service type of the API. Possible values are HTTP, MOCK, TSF, CLB, SCF, WEBSOCKET, and TARGET (internal testing).
  * `service_websocket_cleanup_function_name` - Scf websocket cleaning function. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.
  * `service_websocket_cleanup_function_namespace` - Scf websocket cleanup function namespace. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.
  * `service_websocket_cleanup_function_qualifier` - Scf websocket cleanup function version. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.
  * `service_websocket_register_function_name` - Scf websocket registration function namespace. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.
  * `service_websocket_register_function_namespace` - Scf websocket registration function namespace. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.
  * `service_websocket_register_function_qualifier` - Scf websocket transfer function version. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.
  * `service_websocket_transport_function_name` - Scf websocket transfer function. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.
  * `service_websocket_transport_function_namespace` - Scf websocket transfer function namespace. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.
  * `service_websocket_transport_function_qualifier` - Scf websocket transfer function version. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.
  * `tags` - API binding tag information.
    * `key` - Key of the label.
    * `value` - The value of the note.



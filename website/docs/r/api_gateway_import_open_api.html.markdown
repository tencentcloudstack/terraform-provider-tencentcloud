---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_import_open_api"
sidebar_current: "docs-tencentcloud-resource-api_gateway_import_open_api"
description: |-
  Provides a resource to create a apiGateway import_open_api
---

# tencentcloud_api_gateway_import_open_api

Provides a resource to create a apiGateway import_open_api

## Example Usage

### Import open Api by YAML

```hcl
resource "tencentcloud_api_gateway_import_open_api" "example" {
  service_id      = "service-nxz6yync"
  content         = "info:\n  title: keep-service\n  version: 1.0.1\nopenapi: 3.0.0\npaths:\n  /api/test:\n    get:\n      description: desc\n      operationId: test\n      responses:\n        '200':\n          content:\n            text/html:\n              example: '200'\n          description: '200'\n        default:\n          content:\n            text/html:\n              example: '400'\n          description: '400'\n      x-apigw-api-business-type: NORMAL\n      x-apigw-api-type: NORMAL\n      x-apigw-backend:\n        ServiceConfig:\n          Method: GET\n          Path: /test\n          Url: http://domain.com\n        ServiceType: HTTP\n      x-apigw-cors: false\n      x-apigw-protocol: HTTP\n      x-apigw-service-timeout: 15\n"
  encode_type     = "YAML"
  content_version = "openAPI"
}
```

### Import open Api by JSON

```hcl
resource "tencentcloud_api_gateway_import_open_api" "example" {
  service_id      = "service-nxz6yync"
  content         = "{\"openapi\": \"3.0.0\", \"info\": {\"title\": \"keep-service\", \"version\": \"1.0.1\"}, \"paths\": {\"/api/test\": {\"get\": {\"operationId\": \"test\", \"description\": \"desc\", \"responses\": {\"200\": {\"description\": \"200\", \"content\": {\"text/html\": {\"example\": \"200\"}}}, \"default\": {\"content\": {\"text/html\": {\"example\": \"400\"}}, \"description\": \"400\"}}, \"x-apigw-api-type\": \"NORMAL\", \"x-apigw-api-business-type\": \"NORMAL\", \"x-apigw-protocol\": \"HTTP\", \"x-apigw-cors\": false, \"x-apigw-service-timeout\": 15, \"x-apigw-backend\": {\"ServiceType\": \"HTTP\", \"ServiceConfig\": {\"Url\": \"http://domain.com\", \"Path\": \"/test\", \"Method\": \"GET\"}}}}}}"
  encode_type     = "JSON"
  content_version = "openAPI"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String, ForceNew) OpenAPI body content.
* `service_id` - (Required, String, ForceNew) The unique ID of the service where the API is located.
* `content_version` - (Optional, String, ForceNew) The Content version defaults to OpenAPI and currently only supports OpenAPI.
* `encode_type` - (Optional, String, ForceNew) The Content format can only be YAML or JSON, and the default is YAML.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_business_type` - When `auth_type` is OAUTH, this field is valid, NORMAL: Business API, OAUTH: Authorization API.
* `api_desc` - Custom API description.
* `api_id` - Custom Api Id.
* `api_name` - Custom API name.
* `api_type` - API type, supports NORMAL (regular API) and TSF (microservice API), defaults to NORMAL.
* `auth_relation_api_id` - The unique ID of the associated authorization API takes effect when AuthType is OAUTH and ApiBusinessType is NORMAL. The unique ID of the oauth2.0 authorized API that identifies the business API binding.
* `auth_type` - API authentication type. Support SECRET (Key Pair Authentication), NONE (Authentication Exemption), OAUTH, APP (Application Authentication). The default is NONE.
* `constant_parameters` - Constant parameter.
  * `default_value` - Default value for constant parameters. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `desc` - Constant parameter description. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `name` - Constant parameter name. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `position` - Constant parameter position. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `enable_cors` - Whether to enable CORS. Default value: `true`.
* `is_base64_encoded` - Whether to enable Base64 encoding will only take effect when the backend is scf.
* `is_debug_after_charge` - Charge after starting debugging. (Cloud Market Reserved Fields).
* `is_delete_response_error_codes` - Do you want to delete the custom response configuration error code? If it is not passed or False is passed, it will not be deleted. If True is passed, all custom response configuration error codes for this API will be deleted.
* `micro_services` - API bound microservice list.
  * `cluster_id` - Micro service cluster.
  * `micro_service_name` - Microservice name.
  * `namespace_id` - Microservice namespace.
* `oauth_config` - OAuth configuration. Effective when AuthType is OAUTH.
  * `login_redirect_url` - Redirect address, used to guide users in login operations.
  * `public_key` - Public key, used to verify user tokens.
  * `token_location` - Token passes the position.
* `protocol` - API frontend request type. Valid values: `HTTP`, `WEBSOCKET`. Default value: `HTTP`.
* `request_config_method` - Request frontend method configuration. Valid values: `GET`,`POST`,`PUT`,`DELETE`,`HEAD`,`ANY`. Default value: `GET`.
* `request_config_path` - Request frontend path configuration. Like `/user/getinfo`.
* `request_parameters` - Frontend request parameters.
  * `default_value` - Parameter default value.
  * `desc` - Parameter description.
  * `name` - Parameter name.
  * `position` - Parameter location.
  * `required` - If this parameter required. Default value: `false`.
  * `type` - Parameter type.
* `response_error_codes` - Custom error code configuration. Must keep at least one after set.
  * `code` - Custom response configuration error code.
  * `converted_code` - Custom error code conversion.
  * `desc` - Parameter description.
  * `msg` - Custom response configuration error message.
  * `need_convert` - Whether to enable error code conversion. Default value: `false`.
* `response_fail_example` - Response failure sample of custom response configuration.
* `response_success_example` - Successful response sample of custom response configuration.
* `response_type` - Return type. Valid values: `HTML`, `JSON`, `TEXT`, `BINARY`, `XML`. Default value: `HTML`.
* `service_config_cos_config` - API backend COS configuration. If ServiceType is COS, then this parameter must be passed.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `action` - The API calls the backend COS method, and the optional values for the front-end request method and Action are:GET: GetObjectPUT: PutObjectPOST: PostObject, AppendObjectHEAD: HeadObjectDELETE: DeleteObject.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `authorization` - The API calls the signature switch of the backend COS, which defaults to false.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `bucket_name` - The bucket name of the API backend COS.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `path_match_mode` - Path matching mode for API backend COS, optional values:BackEndPath: Backend path matchingFullPath: Full Path MatchingThe default value is: BackEndPathNote: This field may return null, indicating that a valid value cannot be obtained.
* `service_config_method` - API backend service request method, such as `GET`. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_method` and backend method `service_config_method` can be different.
* `service_config_mock_return_message` - Returned information of API backend mocking. This parameter is required when `service_config_type` is `MOCK`.
* `service_config_path` - API backend service path, such as /path. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_path` and backend path `service_config_path` can be different.
* `service_config_product` - Backend type. Effective when enabling vpc, currently supported types are clb, cvm, and upstream.
* `service_config_scf_function_name` - SCF function name. This parameter takes effect when `service_config_type` is `SCF`.
* `service_config_scf_function_namespace` - SCF function namespace. This parameter takes effect when `service_config_type` is `SCF`.
* `service_config_scf_function_qualifier` - SCF function version. This parameter takes effect when `service_config_type` is `SCF`.
* `service_config_scf_function_type` - Scf function type. Effective when the backend type is SCF. Support Event Triggering (EVENT) and HTTP Direct Cloud Function (HTTP).
* `service_config_scf_is_integrated_response` - Whether to enable response integration. Effective when the backend type is SCF.
* `service_config_timeout` - API backend service timeout period in seconds. Default value: `5`.
* `service_config_type` - The backend service type of the API. Supports HTTP, MOCK, TSF, SCF, WEBSOCKET, COS, TARGET (internal testing).
* `service_config_upstream_id` - Only required when binding to VPC channelsNote: This field may return null, indicating that a valid value cannot be obtained.
* `service_config_url` - The backend service URL of the API. If the ServiceType is HTTP, this parameter must be passed.
* `service_config_vpc_id` - Unique VPC ID.
* `service_config_websocket_cleanup_function_name` - Scf websocket cleaning function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_cleanup_function_namespace` - Scf websocket cleans up the function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_cleanup_function_qualifier` - Scf websocket cleaning function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_register_function_name` - Scf websocket registration function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_register_function_namespace` - Scf websocket registers function namespaces. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_register_function_qualifier` - Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_transport_function_name` - Scf websocket transfer function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_transport_function_namespace` - Scf websocket transfer function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_config_websocket_transport_function_qualifier` - Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.
* `service_parameters` - The backend service parameters of the API.
  * `default_value` - The default value for the backend service parameters of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `name` - The backend service parameter name of the API. This parameter is only used when ServiceType is HTTP. The front and rear parameter names can be different.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `position` - The backend service parameter location of the API, such as head. This parameter is only used when ServiceType is HTTP. The parameter positions at the front and rear ends can be configured differently.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `relevant_request_parameter_desc` - Remarks on the backend service parameters of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `relevant_request_parameter_name` - The name of the front-end parameter corresponding to the backend service parameter of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `relevant_request_parameter_position` - The location of the front-end parameters corresponding to the backend service parameters of the API, such as head. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `relevant_request_parameter_type` - The backend service parameter type of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.
* `service_tsf_health_check_conf` - Health check configuration for microservices.
  * `error_threshold_percentage` - Threshold percentage.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `is_health_check` - Whether to initiate a health check.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `request_volume_threshold` - Health check threshold.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `sleep_window_in_milliseconds` - Window size.Note: This field may return null, indicating that a valid value cannot be obtained.
* `service_tsf_load_balance_conf` - Load balancing configuration for microservices.
  * `is_load_balance` - Is load balancing enabled.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `method` - Load balancing method.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `session_stick_required` - Whether to enable session persistence.Note: This field may return null, indicating that a valid value cannot be obtained.
  * `session_stick_timeout` - Session hold timeout.Note: This field may return null, indicating that a valid value cannot be obtained.
* `update_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.



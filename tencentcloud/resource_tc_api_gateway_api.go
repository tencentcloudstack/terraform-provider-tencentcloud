/*
Use this resource to create API of API gateway.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf-example"
  protocol     = "http&https"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id            = tencentcloud_api_gateway_service.example.id
  api_name              = "tf-example"
  api_desc              = "desc."
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

  release_limit    = 500
  pre_limit        = 500
  test_limit       = 500
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudAPIGatewayAPI() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayAPICreate,
		Read:   resourceTencentCloudAPIGatewayAPIRead,
		Update: resourceTencentCloudAPIGatewayAPIUpdate,
		Delete: resourceTencentCloudAPIGatewayAPIDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique ID of the service where the API is located. Refer to resource `tencentcloud_api_gateway_service`.",
			},
			"api_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom API name.",
			},
			"api_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom API description.",
			},
			"api_type": {
				Optional:     true,
				Type:         schema.TypeString,
				Default:      API_GATEWAY_API_TYPE_NORMAL,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_APT_TYPES),
				Description:  "API type, supports NORMAL (regular API) and TSF (microservice API), defaults to NORMAL.",
			},
			"auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      API_GATEWAY_AUTH_TYPE_NONE,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_AUTH_TYPES),
				Description:  "API authentication type. Support SECRET (Key Pair Authentication), NONE (Authentication Exemption), OAUTH, APP (Application Authentication). The default is NONE.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      API_GATEWAY_API_PROTOCOL_HTTP,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_API_PROTOCOLS),
				Description:  "API frontend request type. Valid values: `HTTP`, `WEBSOCKET`. Default value: `HTTP`.",
			},
			"enable_cors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable CORS. Default value: `true`.",
			},
			"request_config_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Request frontend path configuration. Like `/user/getinfo`.",
			},
			"request_config_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "GET",
				ValidateFunc: validateAllowedStringValue([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "ANY"}),
				Description:  "Request frontend method configuration. Valid values: `GET`,`POST`,`PUT`,`DELETE`,`HEAD`,`ANY`. Default value: `GET`.",
			},
			"constant_parameters": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "Constant parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Constant parameter name. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Constant parameter description. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"position": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Constant parameter position. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default value for constant parameters. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"request_parameters": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Frontend request parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter name.",
						},
						"position": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter location.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter type.",
						},
						"desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter description.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter default value.",
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "If this parameter required. Default value: `false`.",
						},
					},
				},
			},
			"micro_services": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "API bound microservice list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Micro service cluster.",
						},
						"namespace_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Microservice namespace.",
						},
						"micro_service_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Microservice name.",
						},
					},
				},
			},
			"service_tsf_load_balance_conf": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Load balancing configuration for microservices.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_load_balance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Is load balancing enabled.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Load balancing method.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"session_stick_required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable session persistence.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"session_stick_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Session hold timeout.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"service_tsf_health_check_conf": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Health check configuration for microservices.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_health_check": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to initiate a health check.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"request_volume_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Health check threshold.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"sleep_window_in_milliseconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Window size.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"error_threshold_percentage": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Threshold percentage.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"target_services": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Target type backend resource information. (Internal testing stage).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "vm ip.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "vpc id.",
						},
						"vm_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "vm port.",
						},
						"host_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Host IP of the CVM.",
						},
						"docker_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "docker ip.",
						},
					},
				},
			},
			"target_services_load_balance_conf": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Target type load balancing configuration. (Internal testing stage).",
			},
			"target_services_health_check_conf": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Target health check configuration. (Internal testing stage).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_health_check": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to initiate a health check.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"request_volume_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Health check threshold.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"sleep_window_in_milliseconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Window size.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"error_threshold_percentage": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Threshold percentage.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"api_business_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "When `auth_type` is OAUTH, this field is valid, NORMAL: Business API, OAUTH: Authorization API.",
			},
			"service_config_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      API_GATEWAY_SERVICE_TYPE_HTTP,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_SERVICE_TYPES),
				Description:  "The backend service type of the API. Supports HTTP, MOCK, TSF, SCF, WEBSOCKET, COS, TARGET (internal testing).",
			},
			"service_config_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "API backend service timeout period in seconds. Default value: `5`.",
			},
			"service_config_product": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backend type. Effective when enabling vpc, currently supported types are clb, cvm, and upstream.",
			},
			"service_config_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique VPC ID.",
			},
			"service_config_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The backend service URL of the API. If the ServiceType is HTTP, this parameter must be passed.",
			},
			"service_config_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API backend service path, such as /path. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_path` and backend path `service_config_path` can be different.",
			},
			"service_config_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API backend service request method, such as `GET`. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_method` and backend method `service_config_method` can be different.",
			},
			"service_config_upstream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Only required when binding to VPC channelsNote: This field may return null, indicating that a valid value cannot be obtained.",
			},
			"service_config_cos_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "API backend COS configuration. If ServiceType is COS, then this parameter must be passed.Note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The API calls the backend COS method, and the optional values for the front-end request method and Action are:GET: GetObjectPUT: PutObjectPOST: PostObject, AppendObjectHEAD: HeadObjectDELETE: DeleteObject.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"bucket_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The bucket name of the API backend COS.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"authorization": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The API calls the signature switch of the backend COS, which defaults to false.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"path_match_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path matching mode for API backend COS, optional values:BackEndPath: Backend path matchingFullPath: Full Path MatchingThe default value is: BackEndPathNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"service_config_scf_function_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SCF function name. This parameter takes effect when `service_config_type` is `SCF`.",
			},
			"service_config_scf_function_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SCF function namespace. This parameter takes effect when `service_config_type` is `SCF`.",
			},
			"service_config_scf_function_qualifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SCF function version. This parameter takes effect when `service_config_type` is `SCF`.",
			},
			"service_config_scf_function_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf function type. Effective when the backend type is SCF. Support Event Triggering (EVENT) and HTTP Direct Cloud Function (HTTP).",
			},
			"service_config_scf_is_integrated_response": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable response integration. Effective when the backend type is SCF.",
			},
			"service_config_mock_return_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returned information of API backend mocking. This parameter is required when `service_config_type` is `MOCK`.",
			},
			"service_config_websocket_register_function_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket registration function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_cleanup_function_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleaning function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_transport_function_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_register_function_namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket registers function namespaces. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_register_function_qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_transport_function_namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_transport_function_qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_cleanup_function_namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleans up the function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_cleanup_function_qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleaning function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"is_debug_after_charge": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Charge after starting debugging. (Cloud Market Reserved Fields).",
			},
			"is_delete_response_error_codes": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Do you want to delete the custom response configuration error code? If it is not passed or False is passed, it will not be deleted. If True is passed, all custom response configuration error codes for this API will be deleted.",
			},
			"response_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_API_RESPONSE_TYPES),
				Description:  "Return type. Valid values: `HTML`, `JSON`, `TEXT`, `BINARY`, `XML`. Default value: `HTML`.",
			},
			"response_success_example": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Successful response sample of custom response configuration.",
			},
			"response_fail_example": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Response failure sample of custom response configuration.",
			},
			"response_error_codes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Custom error code configuration. Must keep at least one after set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Custom response configuration error code.",
						},
						"msg": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Custom response configuration error message.",
						},
						"desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter description.",
						},
						"converted_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Custom error code conversion.",
						},
						"need_convert": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable error code conversion. Default value: `false`.",
						},
					},
				},
			},
			"auth_relation_api_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the associated authorization API takes effect when AuthType is OAUTH and ApiBusinessType is NORMAL. The unique ID of the oauth2.0 authorized API that identifies the business API binding.",
			},
			"service_parameters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The backend service parameters of the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The backend service parameter name of the API. This parameter is only used when ServiceType is HTTP. The front and rear parameter names can be different.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"position": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The backend service parameter location of the API, such as head. This parameter is only used when ServiceType is HTTP. The parameter positions at the front and rear ends can be configured differently.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"relevant_request_parameter_position": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The location of the front-end parameters corresponding to the backend service parameters of the API, such as head. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"relevant_request_parameter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the front-end parameter corresponding to the backend service parameter of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The default value for the backend service parameters of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"relevant_request_parameter_desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remarks on the backend service parameters of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"relevant_request_parameter_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The backend service parameter type of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"oauth_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "OAuth configuration. Effective when AuthType is OAUTH.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Public key, used to verify user tokens.",
						},
						"token_location": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Token passes the position.",
						},
						"login_redirect_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Redirect address, used to guide users in login operations.",
						},
					},
				},
			},
			"target_namespace_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Tsf serverless namespace ID. (In internal testing).",
			},
			"user_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User type.",
			},
			"is_base64_encoded": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable Base64 encoding will only take effect when the backend is scf.",
			},
			"event_bus_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Event bus ID.",
			},
			"eiam_app_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "EIAM application type.",
			},
			"eiam_auth_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The EIAM application authentication type supports AuthenticationOnly, Authentication, and Authorization.",
			},
			"eiam_app_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "EIAM application ID.",
			},
			"token_timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The effective time of the EIAM application token, measured in seconds, defaults to 7200 seconds.",
			},
			"owner": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Owner of resources.",
			},
			"release_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "API QPS value. Enter a positive number to limit the API query rate per second `QPS`.",
			},
			"pre_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "API QPS value. Enter a positive number to limit the API query rate per second `QPS`.",
			},
			"test_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "API QPS value. Enter a positive number to limit the API query rate per second `QPS`.",
			},
			// Computed values.
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayAPICreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api.create")()

	var (
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		err               error
		response          = apigateway.NewCreateApiResponse()
		request           = apigateway.NewCreateApiRequest()
		serviceId         = d.Get("service_id").(string)
		has               bool

		releaseLimit int
		preLimit     int
		testLimit    int
	)

	request.ServiceId = &serviceId
	request.ApiName = helper.String(d.Get("api_name").(string))
	if object, ok := d.GetOk("api_desc"); ok {
		request.ApiDesc = helper.String(object.(string))
	}

	request.ApiType = helper.String(d.Get("api_type").(string))
	request.AuthType = helper.String(d.Get("auth_type").(string))
	request.Protocol = helper.String(d.Get("protocol").(string))
	request.EnableCORS = helper.Bool(d.Get("enable_cors").(bool))
	request.RequestConfig =
		&apigateway.ApiRequestConfig{Path: helper.String(d.Get("request_config_path").(string)),
			Method: helper.String(d.Get("request_config_method").(string))}

	if v, ok := d.GetOk("constant_parameters"); ok {
		constantParameters := v.(*schema.Set).List()
		request.ConstantParameters = make([]*apigateway.ConstantParameter, 0, len(constantParameters))
		for _, item := range constantParameters {
			dMap := item.(map[string]interface{})
			constantParameter := apigateway.ConstantParameter{}
			if v, ok := dMap["name"]; ok {
				constantParameter.Name = helper.String(v.(string))
			}
			if v, ok := dMap["desc"]; ok {
				constantParameter.Desc = helper.String(v.(string))
			}
			if v, ok := dMap["position"]; ok {
				constantParameter.Position = helper.String(v.(string))
			}
			if v, ok := dMap["default_value"]; ok {
				constantParameter.DefaultValue = helper.String(v.(string))
			}
			request.ConstantParameters = append(request.ConstantParameters, &constantParameter)
		}
	}

	if object, ok := d.GetOk("request_parameters"); ok {
		parameters := object.(*schema.Set).List()
		request.RequestParameters = make([]*apigateway.RequestParameter, 0, len(parameters))
		for _, parameter := range parameters {
			parameterMap := parameter.(map[string]interface{})
			requestParameter := &apigateway.RequestParameter{
				Name:     helper.String(parameterMap["name"].(string)),
				Position: helper.String(parameterMap["position"].(string)),
				Type:     helper.String(parameterMap["type"].(string)),
				Required: helper.Bool(parameterMap["required"].(bool)),
			}
			if parameterMap["desc"] != nil {
				requestParameter.Desc = helper.String(parameterMap["desc"].(string))
			}
			if parameterMap["default_value"] != nil {
				requestParameter.DefaultValue = helper.String(parameterMap["default_value"].(string))
			}
			request.RequestParameters = append(request.RequestParameters, requestParameter)
		}
	}

	if v, ok := d.GetOk("micro_services"); ok {
		microServices := v.(*schema.Set).List()
		request.MicroServices = make([]*apigateway.MicroServiceReq, 0, len(microServices))
		for _, item := range microServices {
			dMap := item.(map[string]interface{})
			microServiceReq := apigateway.MicroServiceReq{}
			if v, ok := dMap["cluster_id"]; ok {
				microServiceReq.ClusterId = helper.String(v.(string))
			}
			if v, ok := dMap["namespace_id"]; ok {
				microServiceReq.NamespaceId = helper.String(v.(string))
			}
			if v, ok := dMap["micro_service_name"]; ok {
				microServiceReq.MicroServiceName = helper.String(v.(string))
			}
			request.MicroServices = append(request.MicroServices, &microServiceReq)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "service_tsf_load_balance_conf"); ok {
		tsfLoadBalanceConfResp := apigateway.TsfLoadBalanceConfResp{}
		if v, ok := dMap["is_load_balance"]; ok {
			tsfLoadBalanceConfResp.IsLoadBalance = helper.Bool(v.(bool))
		}
		if v, ok := dMap["method"]; ok {
			tsfLoadBalanceConfResp.Method = helper.String(v.(string))
		}
		if v, ok := dMap["session_stick_required"]; ok {
			tsfLoadBalanceConfResp.SessionStickRequired = helper.Bool(v.(bool))
		}
		if v, ok := dMap["session_stick_timeout"]; ok {
			tsfLoadBalanceConfResp.SessionStickTimeout = helper.IntInt64(v.(int))
		}
		request.ServiceTsfLoadBalanceConf = &tsfLoadBalanceConfResp
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "service_tsf_health_check_conf"); ok {
		healthCheckConf := apigateway.HealthCheckConf{}
		if v, ok := dMap["is_health_check"]; ok {
			healthCheckConf.IsHealthCheck = helper.Bool(v.(bool))
		}
		if v, ok := dMap["request_volume_threshold"]; ok {
			healthCheckConf.RequestVolumeThreshold = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["sleep_window_in_milliseconds"]; ok {
			healthCheckConf.SleepWindowInMilliseconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["error_threshold_percentage"]; ok {
			healthCheckConf.ErrorThresholdPercentage = helper.IntInt64(v.(int))
		}
		request.ServiceTsfHealthCheckConf = &healthCheckConf
	}

	if v, ok := d.GetOk("target_services"); ok {
		targetServices := v.(*schema.Set).List()
		request.TargetServices = make([]*apigateway.TargetServicesReq, 0, len(targetServices))
		for _, item := range targetServices {
			dMap := item.(map[string]interface{})
			targetServicesReq := apigateway.TargetServicesReq{}
			if v, ok := dMap["vm_ip"]; ok {
				targetServicesReq.VmIp = helper.String(v.(string))
			}
			if v, ok := dMap["vpc_id"]; ok {
				targetServicesReq.VpcId = helper.String(v.(string))
			}
			if v, ok := dMap["vm_port"]; ok {
				targetServicesReq.VmPort = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["host_ip"]; ok {
				targetServicesReq.HostIp = helper.String(v.(string))
			}
			if v, ok := dMap["docker_ip"]; ok {
				targetServicesReq.DockerIp = helper.String(v.(string))
			}
			request.TargetServices = append(request.TargetServices, &targetServicesReq)
		}
	}

	if v, ok := d.GetOkExists("target_services_load_balance_conf"); ok {
		request.TargetServicesLoadBalanceConf = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "target_services_health_check_conf"); ok {
		healthCheckConf := apigateway.HealthCheckConf{}
		if v, ok := dMap["is_health_check"]; ok {
			healthCheckConf.IsHealthCheck = helper.Bool(v.(bool))
		}
		if v, ok := dMap["request_volume_threshold"]; ok {
			healthCheckConf.RequestVolumeThreshold = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["sleep_window_in_milliseconds"]; ok {
			healthCheckConf.SleepWindowInMilliseconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["error_threshold_percentage"]; ok {
			healthCheckConf.ErrorThresholdPercentage = helper.IntInt64(v.(int))
		}
		request.TargetServicesHealthCheckConf = &healthCheckConf
	}

	if *request.AuthType == "OAUTH" {
		if v, ok := d.GetOk("api_business_type"); ok {
			request.ApiBusinessType = helper.String(v.(string))
		}
	}

	var serviceType = d.Get("service_config_type").(string)
	request.ServiceType = &serviceType
	request.ServiceTimeout = helper.IntInt64(d.Get("service_config_timeout").(int))

	switch serviceType {

	case API_GATEWAY_SERVICE_TYPE_WEBSOCKET, API_GATEWAY_SERVICE_TYPE_HTTP:
		serviceConfigProduct := d.Get("service_config_product").(string)
		serviceConfigVpcId := d.Get("service_config_vpc_id").(string)
		serviceConfigUrl := d.Get("service_config_url").(string)
		serviceConfigPath := d.Get("service_config_path").(string)
		serviceConfigMethod := d.Get("service_config_method").(string)
		serviceConfigUpstreamId := d.Get("service_config_upstream_id").(string)
		if serviceConfigProduct != "" {
			if serviceConfigVpcId == "" {
				return fmt.Errorf("`service_config_product` need param `service_config_vpc_id`")
			}
		}
		if serviceConfigUrl == "" || serviceConfigPath == "" || serviceConfigMethod == "" {
			return fmt.Errorf("`service_config_url`,`service_config_path`,`service_config_method` is needed if `service_config_type` is `WEBSOCKET` or `HTTP`")
		}
		request.ServiceConfig = &apigateway.ServiceConfig{}
		if serviceConfigProduct != "" {
			request.ServiceConfig.Product = &serviceConfigProduct
		}
		if serviceConfigVpcId != "" {
			request.ServiceConfig.UniqVpcId = &serviceConfigVpcId
		}
		if serviceConfigUpstreamId != "" {
			request.ServiceConfig.UpstreamId = &serviceConfigUpstreamId
		}

		request.ServiceConfig.Url = &serviceConfigUrl
		request.ServiceConfig.Path = &serviceConfigPath
		request.ServiceConfig.Method = &serviceConfigMethod

	case API_GATEWAY_SERVICE_TYPE_MOCK:
		serviceConfigMockReturnMessage := d.Get("service_config_mock_return_message").(string)
		if serviceConfigMockReturnMessage == "" {
			return fmt.Errorf("`service_config_mock_return_message` is needed if `service_config_type` is `MOCK`")
		}
		request.ServiceMockReturnMessage = &serviceConfigMockReturnMessage

	case API_GATEWAY_SERVICE_TYPE_SCF:
		scfFunctionName := d.Get("service_config_scf_function_name").(string)
		scfFunctionNamespace := d.Get("service_config_scf_function_namespace").(string)
		scfFunctionQualifier := d.Get("service_config_scf_function_qualifier").(string)
		scfFunctionType := d.Get("service_config_scf_function_type").(string)
		scfFunctionIntegratedResponse := d.Get("service_config_scf_is_integrated_response").(bool)
		if scfFunctionName == "" || scfFunctionNamespace == "" || scfFunctionQualifier == "" || scfFunctionType == "" {
			return fmt.Errorf("`service_config_scf_function_name`,`service_config_scf_function_namespace`,`service_config_scf_function_qualifier`, `service_config_scf_function_type` is needed if `service_config_type` is `SCF`")
		}
		request.ServiceScfFunctionName = &scfFunctionName
		request.ServiceScfFunctionNamespace = &scfFunctionNamespace
		request.ServiceScfFunctionQualifier = &scfFunctionQualifier
		request.ServiceScfFunctionType = &scfFunctionType
		request.ServiceScfIsIntegratedResponse = &scfFunctionIntegratedResponse

	case API_GATEWAY_SERVICE_TYPE_COS:
		if dMap, ok := helper.InterfacesHeadMap(d, "service_config_cos_config"); ok {
			cosConfig := apigateway.ServiceConfig{}.CosConfig
			if v, ok := dMap["action"]; ok {
				cosConfig.Action = helper.String(v.(string))
			}
			if v, ok := dMap["bucket_name"]; ok {
				cosConfig.BucketName = helper.String(v.(string))
			}
			if v, ok := dMap["authorization"]; ok {
				cosConfig.Authorization = helper.Bool(v.(bool))
			}
			if v, ok := dMap["path_match_mode"]; ok {
				cosConfig.PathMatchMode = helper.String(v.(string))
			}
			request.ServiceConfig.CosConfig = cosConfig
		}
	case API_GATEWAY_SERVICE_TYPE_TSF:
		serviceWebsocketRegisterFunctionName := d.Get("service_config_websocket_register_function_name").(string)
		serviceWebsocketRegisterFunctionNamespace := d.Get("service_config_websocket_register_function_namespace").(string)
		serviceWebsocketRegisterFunctionQualifier := d.Get("service_config_websocket_register_function_qualifier").(string)
		serviceWebsocketCleanupFunctionName := d.Get("service_config_websocket_cleanup_function_name").(string)
		serviceWebsocketCleanupFunctionNamespace := d.Get("service_config_websocket_cleanup_function_namespace").(string)
		serviceWebsocketCleanupFunctionQualifier := d.Get("service_config_websocket_cleanup_function_qualifier").(string)
		serviceWebsocketTransportFunctionName := d.Get("service_config_websocket_transport_function_name").(string)
		serviceWebsocketTransportFunctionNamespace := d.Get("service_config_websocket_transport_function_namespace").(string)
		serviceWebsocketTransportFunctionQualifier := d.Get("service_config_websocket_transport_function_qualifier").(string)

		request.ServiceWebsocketRegisterFunctionName = &serviceWebsocketRegisterFunctionName
		request.ServiceWebsocketRegisterFunctionNamespace = &serviceWebsocketRegisterFunctionNamespace
		request.ServiceWebsocketRegisterFunctionQualifier = &serviceWebsocketRegisterFunctionQualifier
		request.ServiceWebsocketCleanupFunctionName = &serviceWebsocketCleanupFunctionName
		request.ServiceWebsocketCleanupFunctionNamespace = &serviceWebsocketCleanupFunctionNamespace
		request.ServiceWebsocketCleanupFunctionQualifier = &serviceWebsocketCleanupFunctionQualifier
		request.ServiceWebsocketTransportFunctionName = &serviceWebsocketTransportFunctionName
		request.ServiceWebsocketTransportFunctionNamespace = &serviceWebsocketTransportFunctionNamespace
		request.ServiceWebsocketTransportFunctionQualifier = &serviceWebsocketTransportFunctionQualifier
	}

	if v, ok := d.GetOkExists("is_debug_after_charge"); ok {
		request.IsDebugAfterCharge = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_delete_response_error_codes"); ok {
		request.IsDeleteResponseErrorCodes = helper.Bool(v.(bool))
	}

	request.ResponseType = helper.String(d.Get("response_type").(string))

	if object, ok := d.GetOk("response_success_example"); ok {
		request.ResponseSuccessExample = helper.String(object.(string))
	}

	if object, ok := d.GetOk("response_fail_example"); ok {
		request.ResponseFailExample = helper.String(object.(string))
	}

	if v, ok := d.GetOk("auth_relation_api_id"); ok {
		request.AuthRelationApiId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_parameters"); ok {
		serviceParameters := v.(*schema.Set).List()
		request.ServiceParameters = make([]*apigateway.ServiceParameter, 0, len(serviceParameters))
		for _, item := range serviceParameters {
			dMap := item.(map[string]interface{})
			serviceParameter := apigateway.ServiceParameter{}
			if v, ok := dMap["name"]; ok {
				serviceParameter.Name = helper.String(v.(string))
			}
			if v, ok := dMap["position"]; ok {
				serviceParameter.Position = helper.String(v.(string))
			}
			if v, ok := dMap["relevant_request_parameter_position"]; ok {
				serviceParameter.RelevantRequestParameterPosition = helper.String(v.(string))
			}
			if v, ok := dMap["relevant_request_parameter_name"]; ok {
				serviceParameter.RelevantRequestParameterName = helper.String(v.(string))
			}
			if v, ok := dMap["default_value"]; ok {
				serviceParameter.DefaultValue = helper.String(v.(string))
			}
			if v, ok := dMap["relevant_request_parameter_desc"]; ok {
				serviceParameter.RelevantRequestParameterDesc = helper.String(v.(string))
			}
			if v, ok := dMap["relevant_request_parameter_type"]; ok {
				serviceParameter.RelevantRequestParameterType = helper.String(v.(string))
			}
			request.ServiceParameters = append(request.ServiceParameters, &serviceParameter)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "oauth_config"); ok {
		oauthConfig := apigateway.OauthConfig{}
		if v, ok := dMap["public_key"]; ok {
			oauthConfig.PublicKey = helper.String(v.(string))
		}
		if v, ok := dMap["token_location"]; ok {
			oauthConfig.TokenLocation = helper.String(v.(string))
		}
		if v, ok := dMap["login_redirect_url"]; ok {
			oauthConfig.LoginRedirectUrl = helper.String(v.(string))
		}
		request.OauthConfig = &oauthConfig
	}

	if object, ok := d.GetOk("response_error_codes"); ok {
		codes := object.(*schema.Set).List()
		request.ResponseErrorCodes = make([]*apigateway.ResponseErrorCodeReq, 0, len(codes))
		for _, code := range codes {
			codeMap := code.(map[string]interface{})
			codeReq := &apigateway.ResponseErrorCodeReq{}
			codeReq.Code = helper.IntInt64(codeMap["code"].(int))
			codeReq.Msg = helper.String(codeMap["msg"].(string))

			if codeMap["desc"] != nil {
				codeReq.Desc = helper.String(codeMap["desc"].(string))
			}
			if codeMap["converted_code"] != nil {
				codeReq.ConvertedCode = helper.IntInt64(codeMap["converted_code"].(int))
			}
			if codeMap["need_convert"] != nil {
				codeReq.NeedConvert = helper.Bool(codeMap["need_convert"].(bool))
			}
			if *codeReq.NeedConvert && codeReq.ConvertedCode == nil {
				return fmt.Errorf("`need_convert` need `converted_code`setted")
			}
			request.ResponseErrorCodes = append(request.ResponseErrorCodes, codeReq)
		}
	}

	if v, ok := d.GetOk("target_namespace_id"); ok {
		request.TargetNamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_type"); ok {
		request.UserType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_base64_encoded"); ok {
		request.IsBase64Encoded = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("event_bus_id"); ok {
		request.EventBusId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_scf_function_type"); ok {
		request.ServiceScfFunctionType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("eiam_app_type"); ok {
		request.EIAMAppType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("eiam_auth_type"); ok {
		request.EIAMAuthType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("eiam_app_id"); ok {
		request.EIAMAppId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("token_timeout"); ok {
		request.TokenTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("owner"); ok {
		request.Owner = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pre_limit"); ok {
		preLimit = v.(int)
	}

	if v, ok := d.GetOk("release_limit"); ok {
		releaseLimit = v.(int)
	}

	if v, ok := d.GetOk("test_limit"); ok {
		testLimit = v.(int)
	}

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, err = apiGatewayService.DescribeService(ctx, serviceId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("service %s not exist on server", serviceId)
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = apiGatewayService.client.UseAPIGatewayClient().CreateApi(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if response == nil || response.Response.Result == nil || response.Response.Result.ApiId == nil {
		return fmt.Errorf("create API fail, return nil response")
	}

	if preLimit != 0 {
		_, err = apiGatewayService.ModifyApiEnvironmentStrategy(ctx, serviceId, int64(preLimit), "prepub", []string{*response.Response.Result.ApiId})
		if err != nil {
			return err
		}
	}

	if releaseLimit != 0 {
		_, err = apiGatewayService.ModifyApiEnvironmentStrategy(ctx, serviceId, int64(releaseLimit), "release", []string{*response.Response.Result.ApiId})
		if err != nil {
			return err
		}
	}

	if testLimit != 0 {
		_, err = apiGatewayService.ModifyApiEnvironmentStrategy(ctx, serviceId, int64(testLimit), "test", []string{*response.Response.Result.ApiId})
		if err != nil {
			return err
		}
	}

	d.SetId(*response.Response.Result.ApiId)

	return resourceTencentCloudAPIGatewayAPIRead(d, meta)
}

func resourceTencentCloudAPIGatewayAPIRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api.read")()
	defer inconsistentCheck(d, meta)()

	var (
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiId             = d.Id()
		serviceId         = d.Get("service_id").(string)
		info              apigateway.ApiInfo
		has               bool
		err               error
		releaseLimit      int64
		preLimit          int64
		testLimit         int64
	)

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, err = apiGatewayService.DescribeApi(ctx, serviceId, apiId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("service_id", info.ServiceId)
	_ = d.Set("api_name", info.ApiName)
	_ = d.Set("api_desc", info.ApiDesc)
	_ = d.Set("api_type", info.ApiType)
	_ = d.Set("auth_type", info.AuthType)
	_ = d.Set("protocol", info.Protocol)
	_ = d.Set("request_config_path", info.RequestConfig.Path)
	_ = d.Set("request_config_method", info.RequestConfig.Method)
	_ = d.Set("enable_cors", info.EnableCORS)

	if info.ConstantParameters != nil {
		constantParametersList := []interface{}{}
		for _, constantParameters := range info.ConstantParameters {
			constantParametersMap := map[string]interface{}{}

			if constantParameters.Name != nil {
				constantParametersMap["name"] = constantParameters.Name
			}

			if constantParameters.Desc != nil {
				constantParametersMap["desc"] = constantParameters.Desc
			}

			if constantParameters.Position != nil {
				constantParametersMap["position"] = constantParameters.Position
			}

			if constantParameters.DefaultValue != nil {
				constantParametersMap["default_value"] = constantParameters.DefaultValue
			}

			constantParametersList = append(constantParametersList, constantParametersMap)
		}

		_ = d.Set("constant_parameters", constantParametersList)
	}

	if info.RequestParameters != nil {
		list := make([]map[string]interface{}, 0, len(info.RequestParameters))
		for _, param := range info.RequestParameters {
			list = append(list, map[string]interface{}{
				"name":          param.Name,
				"position":      param.Position,
				"type":          param.Type,
				"desc":          param.Desc,
				"default_value": param.DefaultValue,
				"required":      param.Required,
			})
		}
		_ = d.Set("request_parameters", list)
	}

	if info.MicroServices != nil {
		microServicesList := []interface{}{}
		for _, microServices := range info.MicroServices {
			microServicesMap := map[string]interface{}{}

			if microServices.ClusterId != nil {
				microServicesMap["cluster_id"] = microServices.ClusterId
			}

			if microServices.NamespaceId != nil {
				microServicesMap["namespace_id"] = microServices.NamespaceId
			}

			if microServices.MicroServiceName != nil {
				microServicesMap["micro_service_name"] = microServices.MicroServiceName
			}

			microServicesList = append(microServicesList, microServicesMap)
		}

		_ = d.Set("micro_services", microServicesList)

	}

	if info.ServiceTsfLoadBalanceConf != nil {
		ServiceTsfLoadBalanceConfMap := map[string]interface{}{}

		if info.ServiceTsfLoadBalanceConf.IsLoadBalance != nil {
			ServiceTsfLoadBalanceConfMap["is_load_balance"] = info.ServiceTsfLoadBalanceConf.IsLoadBalance
		}

		if info.ServiceTsfLoadBalanceConf.Method != nil {
			ServiceTsfLoadBalanceConfMap["method"] = info.ServiceTsfLoadBalanceConf.Method
		}

		if info.ServiceTsfLoadBalanceConf.SessionStickRequired != nil {
			ServiceTsfLoadBalanceConfMap["session_stick_required"] = info.ServiceTsfLoadBalanceConf.SessionStickRequired
		}

		if info.ServiceTsfLoadBalanceConf.SessionStickTimeout != nil {
			ServiceTsfLoadBalanceConfMap["session_stick_timeout"] = info.ServiceTsfLoadBalanceConf.SessionStickTimeout
		}

		_ = d.Set("service_tsf_load_balance_conf", []interface{}{ServiceTsfLoadBalanceConfMap})
	}

	if info.ServiceTsfHealthCheckConf != nil {
		serviceTsfHealthCheckConfMap := map[string]interface{}{}

		if info.ServiceTsfHealthCheckConf.IsHealthCheck != nil {
			serviceTsfHealthCheckConfMap["is_health_check"] = info.ServiceTsfHealthCheckConf.IsHealthCheck
		}

		if info.ServiceTsfHealthCheckConf.RequestVolumeThreshold != nil {
			serviceTsfHealthCheckConfMap["request_volume_threshold"] = info.ServiceTsfHealthCheckConf.RequestVolumeThreshold
		}

		if info.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds != nil {
			serviceTsfHealthCheckConfMap["sleep_window_in_milliseconds"] = info.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds
		}

		if info.ServiceTsfHealthCheckConf.ErrorThresholdPercentage != nil {
			serviceTsfHealthCheckConfMap["error_threshold_percentage"] = info.ServiceTsfHealthCheckConf.ErrorThresholdPercentage
		}

		_ = d.Set("service_tsf_health_check_conf", []interface{}{serviceTsfHealthCheckConfMap})
	}

	if info.ApiBusinessType != nil {
		_ = d.Set("api_business_type", info.ApiBusinessType)
	}

	_ = d.Set("service_config_type", info.ServiceType)
	_ = d.Set("service_config_timeout", info.ServiceTimeout)
	if info.ServiceConfig != nil {
		if info.ServiceConfig.Product != nil {
			_ = d.Set("service_config_product", info.ServiceConfig.Product)
		}

		if info.ServiceConfig.UniqVpcId != nil {
			_ = d.Set("service_config_vpc_id", info.ServiceConfig.UniqVpcId)
		}

		if info.ServiceConfig.Url != nil {
			_ = d.Set("service_config_url", info.ServiceConfig.Url)
		}

		if info.ServiceConfig.Path != nil {
			_ = d.Set("service_config_path", info.ServiceConfig.Path)
		}

		if info.ServiceConfig.Method != nil {
			_ = d.Set("service_config_method", info.ServiceConfig.Method)
		}

		if info.ServiceConfig.UpstreamId != nil {
			_ = d.Set("service_config_upstream_id", info.ServiceConfig.UpstreamId)
		}

		if info.ServiceConfig.CosConfig != nil {
			cosConfigMap := map[string]interface{}{}

			if info.ServiceConfig.CosConfig.Action != nil {
				cosConfigMap["action"] = info.ServiceConfig.CosConfig.Action
			}

			if info.ServiceConfig.CosConfig.BucketName != nil {
				cosConfigMap["bucket_name"] = info.ServiceConfig.CosConfig.BucketName
			}

			if info.ServiceConfig.CosConfig.Authorization != nil {
				cosConfigMap["authorization"] = info.ServiceConfig.CosConfig.Authorization
			}

			if info.ServiceConfig.CosConfig.PathMatchMode != nil {
				cosConfigMap["path_match_mode"] = info.ServiceConfig.CosConfig.PathMatchMode
			}

			_ = d.Set("service_config_cos_config", []interface{}{cosConfigMap})
		}
	}

	_ = d.Set("service_config_scf_function_name", info.ServiceScfFunctionName)
	_ = d.Set("service_config_scf_function_namespace", info.ServiceScfFunctionNamespace)
	_ = d.Set("service_config_scf_function_qualifier", info.ServiceScfFunctionQualifier)
	_ = d.Set("service_config_scf_is_integrated_response", info.ServiceScfIsIntegratedResponse)
	_ = d.Set("service_config_mock_return_message", info.ServiceMockReturnMessage)

	_ = d.Set("service_config_websocket_register_function_name", info.ServiceWebsocketRegisterFunctionName)
	_ = d.Set("service_config_websocket_cleanup_function_name", info.ServiceWebsocketCleanupFunctionName)
	_ = d.Set("service_config_websocket_transport_function_name", info.ServiceWebsocketTransportFunctionName)
	_ = d.Set("service_config_websocket_register_function_namespace", info.ServiceWebsocketRegisterFunctionNamespace)
	_ = d.Set("service_config_websocket_register_function_qualifier", info.ServiceWebsocketRegisterFunctionQualifier)
	_ = d.Set("service_config_websocket_transport_function_namespace", info.ServiceWebsocketTransportFunctionNamespace)
	_ = d.Set("service_config_websocket_transport_function_qualifier", info.ServiceWebsocketTransportFunctionQualifier)
	_ = d.Set("service_config_websocket_cleanup_function_namespace", info.ServiceWebsocketCleanupFunctionNamespace)
	_ = d.Set("service_config_websocket_cleanup_function_qualifier", info.ServiceWebsocketCleanupFunctionQualifier)

	_ = d.Set("is_debug_after_charge", info.IsDebugAfterCharge)
	_ = d.Set("response_type", info.ResponseType)
	_ = d.Set("response_success_example", info.ResponseSuccessExample)
	_ = d.Set("response_fail_example", info.ResponseFailExample)
	_ = d.Set("auth_relation_api_id", info.AuthRelationApiId)
	_ = d.Set("update_time", info.ModifiedTime)
	_ = d.Set("create_time", info.CreatedTime)

	if info.ServiceParameters != nil {
		serviceParametersList := []interface{}{}
		for _, serviceParameters := range info.ServiceParameters {
			serviceParametersMap := map[string]interface{}{}

			if serviceParameters.Name != nil {
				serviceParametersMap["name"] = serviceParameters.Name
			}

			if serviceParameters.Position != nil {
				serviceParametersMap["position"] = serviceParameters.Position
			}

			if serviceParameters.RelevantRequestParameterPosition != nil {
				serviceParametersMap["relevant_request_parameter_position"] = serviceParameters.RelevantRequestParameterPosition
			}

			if serviceParameters.RelevantRequestParameterName != nil {
				serviceParametersMap["relevant_request_parameter_name"] = serviceParameters.RelevantRequestParameterName
			}

			if serviceParameters.DefaultValue != nil {
				serviceParametersMap["default_value"] = serviceParameters.DefaultValue
			}

			if serviceParameters.RelevantRequestParameterDesc != nil {
				serviceParametersMap["relevant_request_parameter_desc"] = serviceParameters.RelevantRequestParameterDesc
			}

			serviceParametersList = append(serviceParametersList, serviceParametersMap)
		}

		_ = d.Set("service_parameters", serviceParametersList)

	}

	if info.OauthConfig != nil {
		oauthConfigMap := map[string]interface{}{}

		if info.OauthConfig.PublicKey != nil {
			oauthConfigMap["public_key"] = info.OauthConfig.PublicKey
		}

		if info.OauthConfig.TokenLocation != nil {
			oauthConfigMap["token_location"] = info.OauthConfig.TokenLocation
		}

		if info.OauthConfig.LoginRedirectUrl != nil {
			oauthConfigMap["login_redirect_url"] = info.OauthConfig.LoginRedirectUrl
		}

		_ = d.Set("oauth_config", []interface{}{oauthConfigMap})
	}

	if info.ResponseErrorCodes != nil {
		list := make([]map[string]interface{}, 0, len(info.ResponseErrorCodes))
		for _, code := range info.ResponseErrorCodes {
			list = append(list, map[string]interface{}{
				"code":           code.Code,
				"msg":            code.Msg,
				"desc":           code.Desc,
				"converted_code": code.ConvertedCode,
				"need_convert":   code.NeedConvert,
			})
		}
		_ = d.Set("response_error_codes", list)
	}

	_ = d.Set("is_base64_encoded", info.IsBase64Encoded)

	environmentList, err := apiGatewayService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{}, apiId)
	if len(environmentList) == 0 {
		return nil
	}

	environmentSet := environmentList[0].EnvironmentStrategySet
	for _, envSet := range environmentSet {
		if envSet == nil {
			continue
		}
		if *envSet.EnvironmentName == API_GATEWAY_SERVICE_ENV_PREPUB {
			if *envSet.Quota == -1 {
				preLimit = QUOTA_MAX
				continue
			}
			preLimit = *envSet.Quota
			continue
		}
		if *envSet.EnvironmentName == API_GATEWAY_SERVICE_ENV_TEST {
			if *envSet.Quota == -1 {
				testLimit = QUOTA_MAX
				continue
			}
			testLimit = *envSet.Quota
			continue
		}
		if *envSet.EnvironmentName == API_GATEWAY_SERVICE_ENV_RELEASE {
			if *envSet.Quota == -1 {
				releaseLimit = QUOTA_MAX
				continue
			}
			releaseLimit = *envSet.Quota
			continue
		}
	}

	_ = d.Set("pre_limit", preLimit)
	_ = d.Set("test_limit", testLimit)
	_ = d.Set("release_limit", releaseLimit)

	return nil
}

func resourceTencentCloudAPIGatewayAPIUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api.update")()

	var (
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		response          = apigateway.NewModifyApiResponse()
		request           = apigateway.NewModifyApiRequest()
		apiId             = d.Id()
		serviceId         = d.Get("service_id").(string)
		err               error
		releaseLimit      int
		preLimit          int
		testLimit         int
	)

	immutableArgs := []string{"target_services"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	d.Partial(true)
	request.ServiceId = &serviceId
	request.ApiId = &apiId
	request.ApiName = helper.String(d.Get("api_name").(string))
	if object, ok := d.GetOk("api_desc"); ok {
		request.ApiDesc = helper.String(object.(string))
	}

	request.ApiType = helper.String(d.Get("api_type").(string))
	request.AuthType = helper.String(d.Get("auth_type").(string))
	request.Protocol = helper.String(d.Get("protocol").(string))
	request.EnableCORS = helper.Bool(d.Get("enable_cors").(bool))
	request.RequestConfig = &apigateway.RequestConfig{
		Path:   helper.String(d.Get("request_config_path").(string)),
		Method: helper.String(d.Get("request_config_method").(string)),
	}

	if v, ok := d.GetOk("constant_parameters"); ok {
		constantParameters := v.(*schema.Set).List()
		request.ConstantParameters = make([]*apigateway.ConstantParameter, 0, len(constantParameters))
		for _, item := range constantParameters {
			dMap := item.(map[string]interface{})
			constantParameter := apigateway.ConstantParameter{}
			if v, ok := dMap["name"]; ok {
				constantParameter.Name = helper.String(v.(string))
			}
			if v, ok := dMap["desc"]; ok {
				constantParameter.Desc = helper.String(v.(string))
			}
			if v, ok := dMap["position"]; ok {
				constantParameter.Position = helper.String(v.(string))
			}
			if v, ok := dMap["default_value"]; ok {
				constantParameter.DefaultValue = helper.String(v.(string))
			}
			request.ConstantParameters = append(request.ConstantParameters, &constantParameter)
		}
	}

	if object, ok := d.GetOk("request_parameters"); ok {
		parameters := object.(*schema.Set).List()
		request.RequestParameters = make([]*apigateway.ReqParameter, 0, len(parameters))
		for _, parameter := range parameters {
			parameterMap := parameter.(map[string]interface{})
			requestParameter := &apigateway.ReqParameter{}
			requestParameter.Name = helper.String(parameterMap["name"].(string))
			requestParameter.Position = helper.String(parameterMap["position"].(string))
			requestParameter.Type = helper.String(parameterMap["type"].(string))
			requestParameter.Required = helper.Bool(parameterMap["required"].(bool))
			if parameterMap["desc"] != nil {
				requestParameter.Desc = helper.String(parameterMap["desc"].(string))
			}
			if parameterMap["default_value"] != nil {
				requestParameter.DefaultValue = helper.String(parameterMap["default_value"].(string))
			}
			request.RequestParameters = append(request.RequestParameters, requestParameter)
		}
	}

	if v, ok := d.GetOk("micro_services"); ok {
		microServices := v.(*schema.Set).List()
		request.MicroServices = make([]*apigateway.MicroServiceReq, 0, len(microServices))
		for _, item := range microServices {
			dMap := item.(map[string]interface{})
			microServiceReq := apigateway.MicroServiceReq{}
			if v, ok := dMap["cluster_id"]; ok {
				microServiceReq.ClusterId = helper.String(v.(string))
			}
			if v, ok := dMap["namespace_id"]; ok {
				microServiceReq.NamespaceId = helper.String(v.(string))
			}
			if v, ok := dMap["micro_service_name"]; ok {
				microServiceReq.MicroServiceName = helper.String(v.(string))
			}
			request.MicroServices = append(request.MicroServices, &microServiceReq)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "service_tsf_load_balance_conf"); ok {
		tsfLoadBalanceConfResp := apigateway.TsfLoadBalanceConfResp{}
		if v, ok := dMap["is_load_balance"]; ok {
			tsfLoadBalanceConfResp.IsLoadBalance = helper.Bool(v.(bool))
		}
		if v, ok := dMap["method"]; ok {
			tsfLoadBalanceConfResp.Method = helper.String(v.(string))
		}
		if v, ok := dMap["session_stick_required"]; ok {
			tsfLoadBalanceConfResp.SessionStickRequired = helper.Bool(v.(bool))
		}
		if v, ok := dMap["session_stick_timeout"]; ok {
			tsfLoadBalanceConfResp.SessionStickTimeout = helper.IntInt64(v.(int))
		}
		request.ServiceTsfLoadBalanceConf = &tsfLoadBalanceConfResp
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "service_tsf_health_check_conf"); ok {
		healthCheckConf := apigateway.HealthCheckConf{}
		if v, ok := dMap["is_health_check"]; ok {
			healthCheckConf.IsHealthCheck = helper.Bool(v.(bool))
		}
		if v, ok := dMap["request_volume_threshold"]; ok {
			healthCheckConf.RequestVolumeThreshold = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["sleep_window_in_milliseconds"]; ok {
			healthCheckConf.SleepWindowInMilliseconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["error_threshold_percentage"]; ok {
			healthCheckConf.ErrorThresholdPercentage = helper.IntInt64(v.(int))
		}
		request.ServiceTsfHealthCheckConf = &healthCheckConf
	}

	if v, ok := d.GetOkExists("target_services_load_balance_conf"); ok {
		request.TargetServicesLoadBalanceConf = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "target_services_health_check_conf"); ok {
		healthCheckConf := apigateway.HealthCheckConf{}
		if v, ok := dMap["is_health_check"]; ok {
			healthCheckConf.IsHealthCheck = helper.Bool(v.(bool))
		}
		if v, ok := dMap["request_volume_threshold"]; ok {
			healthCheckConf.RequestVolumeThreshold = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["sleep_window_in_milliseconds"]; ok {
			healthCheckConf.SleepWindowInMilliseconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["error_threshold_percentage"]; ok {
			healthCheckConf.ErrorThresholdPercentage = helper.IntInt64(v.(int))
		}
		request.TargetServicesHealthCheckConf = &healthCheckConf
	}

	if *request.AuthType == "OAUTH" {
		if v, ok := d.GetOk("api_business_type"); ok {
			request.ApiBusinessType = helper.String(v.(string))
		}
	}

	var serviceType = d.Get("service_config_type").(string)
	request.ServiceType = &serviceType
	request.ServiceTimeout = helper.IntInt64(d.Get("service_config_timeout").(int))

	switch serviceType {
	case API_GATEWAY_SERVICE_TYPE_WEBSOCKET, API_GATEWAY_SERVICE_TYPE_HTTP:
		serviceConfigProduct := d.Get("service_config_product").(string)
		serviceConfigVpcId := d.Get("service_config_vpc_id").(string)
		serviceConfigUrl := d.Get("service_config_url").(string)
		serviceConfigPath := d.Get("service_config_path").(string)
		serviceConfigMethod := d.Get("service_config_method").(string)
		serviceConfigUpstreamId := d.Get("service_config_upstream_id").(string)
		if serviceConfigProduct != "" {
			if serviceConfigProduct != "clb" {
				return fmt.Errorf("`service_config_product` only support `clb` now")
			}
			if serviceConfigVpcId == "" {
				return fmt.Errorf("`service_config_product` need param `service_config_vpc_id`")
			}
		}
		if serviceConfigUrl == "" || serviceConfigPath == "" || serviceConfigMethod == "" {
			return fmt.Errorf("`service_config_url`,`service_config_path`,`service_config_method` is needed if `service_config_type` is `WEBSOCKET` or `HTTP`")
		}
		request.ServiceConfig = &apigateway.ServiceConfig{}
		if serviceConfigProduct != "" {
			request.ServiceConfig.Product = &serviceConfigProduct
		}
		if serviceConfigVpcId != "" {
			request.ServiceConfig.UniqVpcId = &serviceConfigVpcId
		}
		if serviceConfigUpstreamId != "" {
			request.ServiceConfig.UpstreamId = &serviceConfigUpstreamId
		}
		request.ServiceConfig.Url = &serviceConfigUrl
		request.ServiceConfig.Path = &serviceConfigPath
		request.ServiceConfig.Method = &serviceConfigMethod

	case API_GATEWAY_SERVICE_TYPE_MOCK:
		serviceConfigMockReturnMessage := d.Get("service_config_mock_return_message").(string)
		if serviceConfigMockReturnMessage == "" {
			return fmt.Errorf("`service_config_mock_return_message` is needed if `service_config_type` is `MOCK`")
		}
		request.ServiceMockReturnMessage = &serviceConfigMockReturnMessage

	case API_GATEWAY_SERVICE_TYPE_SCF:
		scfFunctionName := d.Get("service_config_scf_function_name").(string)
		scfFunctionNamespace := d.Get("service_config_scf_function_namespace").(string)
		scfFunctionQualifier := d.Get("service_config_scf_function_qualifier").(string)
		scfFunctionType := d.Get("service_config_scf_function_type").(string)
		scfFunctionIntegratedResponse := d.Get("service_config_scf_is_integrated_response").(bool)
		if scfFunctionName == "" || scfFunctionNamespace == "" || scfFunctionQualifier == "" || scfFunctionType == "" {
			return fmt.Errorf("`service_config_scf_function_name`,`service_config_scf_function_namespace`,`service_config_scf_function_qualifier`, `service_config_scf_function_type` is needed if `service_config_type` is `SCF`")
		}
		request.ServiceScfFunctionName = &scfFunctionName
		request.ServiceScfFunctionNamespace = &scfFunctionNamespace
		request.ServiceScfFunctionQualifier = &scfFunctionQualifier
		request.ServiceScfFunctionType = &scfFunctionType
		request.ServiceScfIsIntegratedResponse = &scfFunctionIntegratedResponse

	case API_GATEWAY_SERVICE_TYPE_COS:
		if dMap, ok := helper.InterfacesHeadMap(d, "service_config_cos_config"); ok {
			cosConfig := apigateway.ServiceConfig{}.CosConfig
			if v, ok := dMap["action"]; ok {
				cosConfig.Action = helper.String(v.(string))
			}
			if v, ok := dMap["bucket_name"]; ok {
				cosConfig.BucketName = helper.String(v.(string))
			}
			if v, ok := dMap["authorization"]; ok {
				cosConfig.Authorization = helper.Bool(v.(bool))
			}
			if v, ok := dMap["path_match_mode"]; ok {
				cosConfig.PathMatchMode = helper.String(v.(string))
			}
			request.ServiceConfig.CosConfig = cosConfig
		}
	case API_GATEWAY_SERVICE_TYPE_TSF:
		serviceWebsocketRegisterFunctionName := d.Get("service_config_websocket_register_function_name").(string)
		serviceWebsocketRegisterFunctionNamespace := d.Get("service_config_websocket_register_function_namespace").(string)
		serviceWebsocketRegisterFunctionQualifier := d.Get("service_config_websocket_register_function_qualifier").(string)
		serviceWebsocketCleanupFunctionName := d.Get("service_config_websocket_cleanup_function_name").(string)
		serviceWebsocketCleanupFunctionNamespace := d.Get("service_config_websocket_cleanup_function_namespace").(string)
		serviceWebsocketCleanupFunctionQualifier := d.Get("service_config_websocket_cleanup_function_qualifier").(string)
		serviceWebsocketTransportFunctionName := d.Get("service_config_websocket_transport_function_name").(string)
		serviceWebsocketTransportFunctionNamespace := d.Get("service_config_websocket_transport_function_namespace").(string)
		serviceWebsocketTransportFunctionQualifier := d.Get("service_config_websocket_transport_function_qualifier").(string)

		request.ServiceWebsocketRegisterFunctionName = &serviceWebsocketRegisterFunctionName
		request.ServiceWebsocketRegisterFunctionNamespace = &serviceWebsocketRegisterFunctionNamespace
		request.ServiceWebsocketRegisterFunctionQualifier = &serviceWebsocketRegisterFunctionQualifier
		request.ServiceWebsocketCleanupFunctionName = &serviceWebsocketCleanupFunctionName
		request.ServiceWebsocketCleanupFunctionNamespace = &serviceWebsocketCleanupFunctionNamespace
		request.ServiceWebsocketCleanupFunctionQualifier = &serviceWebsocketCleanupFunctionQualifier
		request.ServiceWebsocketTransportFunctionName = &serviceWebsocketTransportFunctionName
		request.ServiceWebsocketTransportFunctionNamespace = &serviceWebsocketTransportFunctionNamespace
		request.ServiceWebsocketTransportFunctionQualifier = &serviceWebsocketTransportFunctionQualifier
	}

	if v, ok := d.GetOkExists("is_debug_after_charge"); ok {
		request.IsDebugAfterCharge = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_delete_response_error_codes"); ok {
		request.IsDeleteResponseErrorCodes = helper.Bool(v.(bool))
	}

	request.ResponseType = helper.String(d.Get("response_type").(string))

	if object, ok := d.GetOk("response_success_example"); ok {
		request.ResponseSuccessExample = helper.String(object.(string))
	}

	if object, ok := d.GetOk("response_fail_example"); ok {
		request.ResponseFailExample = helper.String(object.(string))
	}

	if v, ok := d.GetOk("auth_relation_api_id"); ok {
		request.AuthRelationApiId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_parameters"); ok {
		serviceParameters := v.(*schema.Set).List()
		request.ServiceParameters = make([]*apigateway.ServiceParameter, 0, len(serviceParameters))
		for _, item := range serviceParameters {
			dMap := item.(map[string]interface{})
			serviceParameter := apigateway.ServiceParameter{}
			if v, ok := dMap["name"]; ok {
				serviceParameter.Name = helper.String(v.(string))
			}
			if v, ok := dMap["position"]; ok {
				serviceParameter.Position = helper.String(v.(string))
			}
			if v, ok := dMap["relevant_request_parameter_position"]; ok {
				serviceParameter.RelevantRequestParameterPosition = helper.String(v.(string))
			}
			if v, ok := dMap["relevant_request_parameter_name"]; ok {
				serviceParameter.RelevantRequestParameterName = helper.String(v.(string))
			}
			if v, ok := dMap["default_value"]; ok {
				serviceParameter.DefaultValue = helper.String(v.(string))
			}
			if v, ok := dMap["relevant_request_parameter_desc"]; ok {
				serviceParameter.RelevantRequestParameterDesc = helper.String(v.(string))
			}
			if v, ok := dMap["relevant_request_parameter_type"]; ok {
				serviceParameter.RelevantRequestParameterType = helper.String(v.(string))
			}
			request.ServiceParameters = append(request.ServiceParameters, &serviceParameter)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "oauth_config"); ok {
		oauthConfig := apigateway.OauthConfig{}
		if v, ok := dMap["public_key"]; ok {
			oauthConfig.PublicKey = helper.String(v.(string))
		}
		if v, ok := dMap["token_location"]; ok {
			oauthConfig.TokenLocation = helper.String(v.(string))
		}
		if v, ok := dMap["login_redirect_url"]; ok {
			oauthConfig.LoginRedirectUrl = helper.String(v.(string))
		}
		request.OauthConfig = &oauthConfig
	}

	if v, ok := d.GetOkExists("is_base64_encoded"); ok {
		request.IsBase64Encoded = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("event_bus_id"); ok {
		request.EventBusId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_scf_function_type"); ok {
		request.ServiceScfFunctionType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("eiam_app_type"); ok {
		request.EIAMAppType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("eiam_auth_type"); ok {
		request.EIAMAuthType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("eiam_app_id"); ok {
		request.EIAMAppId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("token_timeout"); ok {
		request.TokenTimeout = helper.IntInt64(v.(int))
	}

	oldInterface, newInterface := d.GetChange("response_error_codes")

	if oldInterface.(*schema.Set).Len() > 0 && newInterface.(*schema.Set).Len() == 0 {
		return fmt.Errorf("`response_error_codes` must keep at least one after set")
	}

	if object, ok := d.GetOk("response_error_codes"); ok {
		codes := object.(*schema.Set).List()
		request.ResponseErrorCodes = make([]*apigateway.ResponseErrorCodeReq, 0, len(codes))
		for _, code := range codes {
			codeMap := code.(map[string]interface{})
			codeReq := &apigateway.ResponseErrorCodeReq{}
			codeReq.Code = helper.IntInt64(codeMap["code"].(int))
			codeReq.Msg = helper.String(codeMap["msg"].(string))

			if codeMap["desc"] != nil {
				codeReq.Desc = helper.String(codeMap["desc"].(string))
			}
			if codeMap["converted_code"] != nil {
				codeReq.ConvertedCode = helper.IntInt64(codeMap["converted_code"].(int))
			}
			if codeMap["need_convert"] != nil {
				codeReq.NeedConvert = helper.Bool(codeMap["need_convert"].(bool))
			}
			if *codeReq.NeedConvert && codeReq.ConvertedCode == nil {
				return fmt.Errorf("`need_convert` need `converted_code`setted")
			}
			request.ResponseErrorCodes = append(request.ResponseErrorCodes, codeReq)
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = apiGatewayService.client.UseAPIGatewayClient().ModifyApi(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if response == nil {
		return fmt.Errorf("modify API fail, return nil response")
	}

	if d.HasChange("pre_limit") {
		if v, ok := d.GetOk("pre_limit"); ok {
			preLimit = v.(int)
		}
		if preLimit != 0 {
			_, err = apiGatewayService.ModifyApiEnvironmentStrategy(ctx, serviceId, int64(preLimit), "prepub", []string{apiId})
			if err != nil {
				return err
			}
		}

	}

	if d.HasChange("release_limit") {
		if v, ok := d.GetOk("release_limit"); ok {
			releaseLimit = v.(int)
		}
		if releaseLimit != 0 {
			_, err = apiGatewayService.ModifyApiEnvironmentStrategy(ctx, serviceId, int64(preLimit), "release", []string{apiId})
			if err != nil {
				return err
			}
		}

	}

	if d.HasChange("test_limit") {
		if v, ok := d.GetOk("test_limit"); ok {
			testLimit = v.(int)
		}
		if testLimit != 0 {
			_, err = apiGatewayService.ModifyApiEnvironmentStrategy(ctx, serviceId, int64(preLimit), "test", []string{apiId})
			if err != nil {
				return err
			}
		}

	}

	d.Partial(false)
	return resourceTencentCloudAPIGatewayAPIRead(d, meta)
}

func resourceTencentCloudAPIGatewayAPIDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api.delete")()

	var (
		apiGatewayService       = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		logId                   = getLogId(contextNil)
		ctx                     = context.WithValue(context.TODO(), logIdKey, logId)
		apiId                   = d.Id()
		serviceId               = d.Get("service_id").(string)
		limitNumber       int64 = QUOTA
		err               error
	)

	for _, v := range API_GATEWAY_SERVICE_ENVS {
		_, err = apiGatewayService.ModifyApiEnvironmentStrategy(ctx, serviceId, limitNumber, v, []string{apiId})
		if err != nil {
			return err
		}
	}

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err = apiGatewayService.DeleteApi(ctx, serviceId, apiId)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
}

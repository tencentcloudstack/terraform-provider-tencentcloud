package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudApiGatewayImportOpenApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApiGatewayImportOpenApiCreate,
		Read:   resourceTencentCloudApiGatewayImportOpenApiRead,
		Delete: resourceTencentCloudApiGatewayImportOpenApiDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the service where the API is located.",
			},
			"content": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "OpenAPI body content.",
			},
			"encode_type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				Default:      IMPORT_OPEN_API_ENCODE_TYPE_YAML,
				ValidateFunc: validateAllowedStringValue(IMPORT_OPEN_API_ENCODE_TYPE),
				Description:  "The Content format can only be YAML or JSON, and the default is YAML.",
			},
			"content_version": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Default:     "OpenAPI",
				Description: "The Content version defaults to OpenAPI and currently only supports OpenAPI.",
			},
			// Computed
			"api_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom Api Id.",
			},
			"api_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom API name.",
			},
			"api_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom API description.",
			},
			"api_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "API type, supports NORMAL (regular API) and TSF (microservice API), defaults to NORMAL.",
			},
			"auth_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API authentication type. Support SECRET (Key Pair Authentication), NONE (Authentication Exemption), OAUTH, APP (Application Authentication). The default is NONE.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API frontend request type. Valid values: `HTTP`, `WEBSOCKET`. Default value: `HTTP`.",
			},
			"enable_cors": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable CORS. Default value: `true`.",
			},
			"request_config_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Request frontend path configuration. Like `/user/getinfo`.",
			},
			"request_config_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Request frontend method configuration. Valid values: `GET`,`POST`,`PUT`,`DELETE`,`HEAD`,`ANY`. Default value: `GET`.",
			},
			"constant_parameters": {
				Computed:    true,
				Type:        schema.TypeSet,
				Description: "Constant parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Constant parameter name. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Constant parameter description. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Constant parameter position. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default value for constant parameters. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"request_parameters": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Frontend request parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter name.",
						},
						"position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter location.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter type.",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter description.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter default value.",
						},
						"required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If this parameter required. Default value: `false`.",
						},
					},
				},
			},
			"micro_services": {
				Computed:    true,
				Type:        schema.TypeSet,
				Description: "API bound microservice list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Micro service cluster.",
						},
						"namespace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microservice namespace.",
						},
						"micro_service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microservice name.",
						},
					},
				},
			},
			"service_tsf_load_balance_conf": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Load balancing configuration for microservices.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_load_balance": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is load balancing enabled.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load balancing method.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"session_stick_required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable session persistence.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"session_stick_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session hold timeout.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"service_tsf_health_check_conf": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Health check configuration for microservices.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_health_check": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to initiate a health check.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"request_volume_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health check threshold.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"sleep_window_in_milliseconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Window size.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"error_threshold_percentage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Threshold percentage.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"api_business_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "When `auth_type` is OAUTH, this field is valid, NORMAL: Business API, OAUTH: Authorization API.",
			},
			"service_config_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The backend service type of the API. Supports HTTP, MOCK, TSF, SCF, WEBSOCKET, COS, TARGET (internal testing).",
			},
			"service_config_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "API backend service timeout period in seconds. Default value: `5`.",
			},
			"service_config_product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backend type. Effective when enabling vpc, currently supported types are clb, cvm, and upstream.",
			},
			"service_config_vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique VPC ID.",
			},
			"service_config_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The backend service URL of the API. If the ServiceType is HTTP, this parameter must be passed.",
			},
			"service_config_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API backend service path, such as /path. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_path` and backend path `service_config_path` can be different.",
			},
			"service_config_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API backend service request method, such as `GET`. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_method` and backend method `service_config_method` can be different.",
			},
			"service_config_upstream_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only required when binding to VPC channelsNote: This field may return null, indicating that a valid value cannot be obtained.",
			},
			"service_config_cos_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "API backend COS configuration. If ServiceType is COS, then this parameter must be passed.Note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API calls the backend COS method, and the optional values for the front-end request method and Action are:GET: GetObjectPUT: PutObjectPOST: PostObject, AppendObjectHEAD: HeadObjectDELETE: DeleteObject.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"bucket_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The bucket name of the API backend COS.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"authorization": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The API calls the signature switch of the backend COS, which defaults to false.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"path_match_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path matching mode for API backend COS, optional values:BackEndPath: Backend path matchingFullPath: Full Path MatchingThe default value is: BackEndPathNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"service_config_scf_function_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF function name. This parameter takes effect when `service_config_type` is `SCF`.",
			},
			"service_config_scf_function_namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF function namespace. This parameter takes effect when `service_config_type` is `SCF`.",
			},
			"service_config_scf_function_qualifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF function version. This parameter takes effect when `service_config_type` is `SCF`.",
			},
			"service_config_scf_function_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf function type. Effective when the backend type is SCF. Support Event Triggering (EVENT) and HTTP Direct Cloud Function (HTTP).",
			},
			"service_config_scf_is_integrated_response": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable response integration. Effective when the backend type is SCF.",
			},
			"service_config_mock_return_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returned information of API backend mocking. This parameter is required when `service_config_type` is `MOCK`.",
			},
			"service_config_websocket_register_function_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket registration function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_cleanup_function_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleaning function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_transport_function_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_register_function_namespace": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket registers function namespaces. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_register_function_qualifier": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_transport_function_namespace": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_transport_function_qualifier": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_cleanup_function_namespace": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleans up the function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"service_config_websocket_cleanup_function_qualifier": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleaning function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},
			"is_debug_after_charge": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Charge after starting debugging. (Cloud Market Reserved Fields).",
			},
			"is_delete_response_error_codes": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Do you want to delete the custom response configuration error code? If it is not passed or False is passed, it will not be deleted. If True is passed, all custom response configuration error codes for this API will be deleted.",
			},
			"response_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Return type. Valid values: `HTML`, `JSON`, `TEXT`, `BINARY`, `XML`. Default value: `HTML`.",
			},
			"response_success_example": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Successful response sample of custom response configuration.",
			},
			"response_fail_example": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Response failure sample of custom response configuration.",
			},
			"response_error_codes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Custom error code configuration. Must keep at least one after set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Custom response configuration error code.",
						},
						"msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom response configuration error message.",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter description.",
						},
						"converted_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Custom error code conversion.",
						},
						"need_convert": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable error code conversion. Default value: `false`.",
						},
					},
				},
			},
			"auth_relation_api_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the associated authorization API takes effect when AuthType is OAUTH and ApiBusinessType is NORMAL. The unique ID of the oauth2.0 authorized API that identifies the business API binding.",
			},
			"service_parameters": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The backend service parameters of the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend service parameter name of the API. This parameter is only used when ServiceType is HTTP. The front and rear parameter names can be different.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend service parameter location of the API, such as head. This parameter is only used when ServiceType is HTTP. The parameter positions at the front and rear ends can be configured differently.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"relevant_request_parameter_position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The location of the front-end parameters corresponding to the backend service parameters of the API, such as head. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"relevant_request_parameter_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the front-end parameter corresponding to the backend service parameter of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default value for the backend service parameters of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"relevant_request_parameter_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remarks on the backend service parameters of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"relevant_request_parameter_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend service parameter type of the API. This parameter is only used when ServiceType is HTTP.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"oauth_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "OAuth configuration. Effective when AuthType is OAUTH.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public key, used to verify user tokens.",
						},
						"token_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Token passes the position.",
						},
						"login_redirect_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Redirect address, used to guide users in login operations.",
						},
					},
				},
			},
			"is_base64_encoded": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable Base64 encoding will only take effect when the backend is scf.",
			},
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

func resourceTencentCloudApiGatewayImportOpenApiCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_import_open_api.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		request   = apiGateway.NewImportOpenApiRequest()
		response  = apiGateway.NewImportOpenApiResponse()
		serviceId string
		apiId     string
	)

	if v, ok := d.GetOk("service_id"); ok {
		request.ServiceId = helper.String(v.(string))
		serviceId = v.(string)
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("encode_type"); ok {
		request.EncodeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content_version"); ok {
		request.ContentVersion = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().ImportOpenApi(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("apiGateway importOpenApi not exists")
			return resource.NonRetryableError(e)
		}

		if *result.Response.Result.ApiSet[0].Status == "success" {
			response = result
			return nil
		}

		return resource.RetryableError(fmt.Errorf("create apiGateway importOpenApi is running, status: %s", *result.Response.Result.ApiSet[0].Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create apiGateway importOpenApi failed, reason:%+v", logId, err)
		return err
	}

	apiId = *response.Response.Result.ApiSet[0].ApiId
	d.SetId(strings.Join([]string{serviceId, apiId}, FILED_SP))
	return resourceTencentCloudApiGatewayImportOpenApiRead(d, meta)
}

func resourceTencentCloudApiGatewayImportOpenApiRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_import_open_api.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	serviceId := idSplit[0]
	apiId := idSplit[1]

	info, err := service.DescribeApiGatewayImportOpenApiById(ctx, serviceId, apiId)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApiGatewayImportOpenApi` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("service_id", info.ServiceId)
	_ = d.Set("api_id", info.ApiId)
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
	} else {
		_ = d.Set("service_tsf_load_balance_conf", []interface{}{})
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
	} else {
		_ = d.Set("service_tsf_health_check_conf", []interface{}{})
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
		} else {
			_ = d.Set("service_config_cos_config", []interface{}{})
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
	} else {
		_ = d.Set("oauth_config", []interface{}{})
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

	return nil
}

func resourceTencentCloudApiGatewayImportOpenApiDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_import_open_api.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	serviceId := idSplit[0]
	apiId := idSplit[1]

	if err := service.DeleteApiGatewayImportOpenApiById(ctx, serviceId, apiId); err != nil {
		return err
	}

	return nil
}

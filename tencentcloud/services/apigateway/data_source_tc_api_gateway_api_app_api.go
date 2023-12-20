package apigateway

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudApiGatewayApiAppApi() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApiGatewayApiAppApiRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the service where the API resides.",
			},
			"api_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "API interface unique ID.",
			},
			"api_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Api region.",
			},
			// computed
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "API details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the service where the API resides.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the service where the API resides.",
						},
						"service_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description of the service where the API resides.",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API interface unique ID.",
						},
						"api_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the API interface.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time, expressed in accordance with the ISO8601 standard and using UTC time. The format is: YYYY-MM-DDThh:mm:ssZ.",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modification time, expressed in accordance with the ISO8601 standard and using UTC time. The format is: YYYY-MM-DDThh:mm:ssZ.",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the API interface.",
						},
						"api_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API type. Possible values are NORMAL (normal API) and TSF (microservice API).",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The front-end request type of the API, such as HTTP or HTTPS or HTTP and HTTPS.",
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API authentication type. Possible values are SECRET (key pair authentication), NONE (authentication-free), and OAUTH.",
						},
						"api_business_type": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "Type of OAUTH API. Possible values are NORMAL (Business API), OAUTH (Authorization API).",
						},
						"auth_relation_api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OAUTH The unique ID of the authorization API associated with the business API.",
						},
						"oauth_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "OAUTH configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Public key, used to verify user token.",
									},
									"token_location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Token delivery position.",
									},
									"login_redirect_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Redirect address, used to guide users to log in.",
									},
								},
							},
						},
						"is_debug_after_charge": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to debug after purchase (parameters reserved in the cloud market).",
						},
						"request_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The requested frontend configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API path, such as /path.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API request method, such as GET.",
									},
								},
							},
						},
						"response_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Return type.",
						},
						"response_success_example": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom response configuration successful response example.",
						},
						"response_fail_example": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom response configuration failure response example.",
						},
						"response_error_codes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User-defined error code configuration.",
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
										Description: "Custom response configuration error code remarks.",
									},
									"converted_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Custom error code conversion.",
									},
									"need_convert": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether it is necessary to enable error code conversion.",
									},
								},
							},
						},
						"request_parameters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Front-end request parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API front-end parameter name.",
									},
									"position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The front-end parameter position of the API, such as header. Currently supports header, query, path.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API front-end parameter type, such as String, int.",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API front-end parameter default value.",
									},
									"required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: ".",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API front-end parameter remarks.",
									},
								},
							},
						},
						"service_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The backend service timeout of the API, in seconds.",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend service type of the API. Possible values are HTTP, MOCK, TSF, CLB, SCF, WEBSOCKET, and TARGET (internal testing).",
						},
						"service_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Backend service configuration for the API.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"product": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Backend type. It takes effect when vpc is enabled. Currently supported types are clb, cvm and upstream.",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique ID of the vpc.",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API&amp;#39;s backend service url. If ServiceType is HTTP, this parameter must be passed.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API backend service path, such as /path. If ServiceType is HTTP, this parameter is required. The front-end and back-end paths can be different.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API backend service request method, such as GET. If ServiceType is HTTP, this parameter is required. The front-end and back-end methods can be different.",
									},
									"upstream_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Only required when binding vpc channel.",
									},
								},
							},
						},
						"service_parameters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "API backend service parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The backend service parameter name of the API. This parameter will be used only if the ServiceType is HTTP. The front-end and back-end parameter names can be different.",
									},
									"position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The backend service parameter location of the API, such as head. This parameter is only used if the ServiceType is HTTP. The front-end and back-end parameter positions can be configured differently.",
									},
									"relevant_request_parameter_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The front-end parameter position corresponding to the back-end service parameter of the API, such as head. This parameter is only used if the ServiceType is HTTP.",
									},
									"relevant_request_parameter_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The front-end parameter name corresponding to the back-end service parameter of the API. This parameter is only used if the ServiceType is HTTP.",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default values for the APIs backend service parameters. This parameter is only used if the ServiceType is HTTP.",
									},
									"relevant_request_parameter_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Remarks on the backend service parameters of the API. This parameter is only used if the ServiceType is HTTP.",
									},
								},
							},
						},
						"constant_parameters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Constant parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Constant parameter name. This parameter is only used if the ServiceType is HTTP.",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Constant parameter description. This parameter is only used if the ServiceType is HTTP.",
									},
									"position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Constant parameter position. This parameter is only used if the ServiceType is HTTP.",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Constant parameter default value. This parameter is only used if the ServiceType is HTTP.",
									},
								},
							},
						},
						"service_mock_return_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "APIs backend Mock returns information. If ServiceType is Mock, this parameter must be passed.",
						},
						"service_scf_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf function name. Effective when the backend type is SCF.",
						},
						"service_scf_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf function namespace. Effective when the backend type is SCF.",
						},
						"service_scf_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf function version. Effective when the backend type is SCF.",
						},
						"service_scf_is_integrated_response": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable integrated response.",
						},
						"service_websocket_register_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket registration function namespace. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.",
						},
						"service_websocket_register_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket registration function namespace. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.",
						},
						"service_websocket_register_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket transfer function version. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.",
						},
						"service_websocket_cleanup_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket cleaning function. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.",
						},
						"service_websocket_cleanup_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket cleanup function namespace. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.",
						},
						"service_websocket_cleanup_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket cleanup function version. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.",
						},
						"internal_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WEBSOCKET pushback address.",
						},
						"service_websocket_transport_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket transfer function. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.",
						},
						"service_websocket_transport_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket transfer function namespace. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.",
						},
						"service_websocket_transport_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket transfer function version. Valid when the front-end type is WEBSOCKET and the back-end type is SCF.",
						},
						"micro_services": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "API binding microservice list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice cluster ID.",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice namespace ID.",
									},
									"micro_service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice name.",
									},
								},
							},
						},
						"micro_services_info": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Microservice information details.",
						},
						"service_tsf_load_balance_conf": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Load balancing configuration for microservices.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_load_balance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable load balancing.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Load balancing method.",
									},
									"session_stick_required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable session persistence.",
									},
									"session_stick_timeout": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Session retention timeout.",
									},
								},
							},
						},
						"service_tsf_health_check_conf": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Health check configuration for microservices.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_health_check": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable health check.",
									},
									"request_volume_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Health check threshold.",
									},
									"sleep_window_in_milliseconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Window size.",
									},
									"error_threshold_percentage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Threshold percentage.",
									},
								},
							},
						},
						"enable_cors": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable cross-domain.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "API binding tag information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key of the label.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the note.",
									},
								},
							},
						},
						"environments": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "API published environment information.",
						},
						"is_base64_encoded": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable Base64 encoding will only take effect when the backend is scf.",
						},
						"is_base64_trigger": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable Base64-encoded header triggering will only take effect when the backend is scf.",
						},
						"base64_encoded_trigger_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Header triggers rules, and the total number of rules does not exceed 10.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Header for encoding triggering, optional values Accept and Content_Type correspond to Accept and Content-Type in the actual data flow request header.",
									},
									"value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "An array of optional values for the header triggered by encoding. The maximum string length of the array element is 40. The elements can include numbers, English letters and special characters. The optional values for special characters are: `.` `+` ` *` `-` `/` `_` For example [ application/x-vpeg005, application/xhtml+xml, application/vnd.ms -project, application/vnd.rn-rn_music_package] etc. are all legal.",
									},
								},
							},
						},
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save apiAppApis.",
			},
		},
	}
}

func dataSourceTencentCloudApiGatewayApiAppApiRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_api_app_api.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		apiAppApi  *apigateway.ApiInfo
		service_id string
		api_id     string
		api_region string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
		service_id = v.(string)
	}

	if v, ok := d.GetOk("api_id"); ok {
		paramMap["APIId"] = helper.String(v.(string))
		api_id = v.(string)
	}

	if v, ok := d.GetOk("api_region"); ok {
		paramMap["ApiRegion"] = helper.String(v.(string))
		api_region = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApiGatewayApiAppApiByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		apiAppApi = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)
	if apiAppApi != nil {
		apiInfoMap := map[string]interface{}{}

		if apiAppApi.ServiceId != nil {
			apiInfoMap["service_id"] = apiAppApi.ServiceId
		}

		if apiAppApi.ServiceName != nil {
			apiInfoMap["service_name"] = apiAppApi.ServiceName
		}

		if apiAppApi.ServiceDesc != nil {
			apiInfoMap["service_desc"] = apiAppApi.ServiceDesc
		}

		if apiAppApi.ApiId != nil {
			apiInfoMap["api_id"] = apiAppApi.ApiId
		}

		if apiAppApi.ApiDesc != nil {
			apiInfoMap["api_desc"] = apiAppApi.ApiDesc
		}

		if apiAppApi.CreatedTime != nil {
			apiInfoMap["created_time"] = apiAppApi.CreatedTime
		}

		if apiAppApi.ModifiedTime != nil {
			apiInfoMap["modified_time"] = apiAppApi.ModifiedTime
		}

		if apiAppApi.ApiName != nil {
			apiInfoMap["api_name"] = apiAppApi.ApiName
		}

		if apiAppApi.ApiType != nil {
			apiInfoMap["api_type"] = apiAppApi.ApiType
		}

		if apiAppApi.Protocol != nil {
			apiInfoMap["protocol"] = apiAppApi.Protocol
		}

		if apiAppApi.AuthType != nil {
			apiInfoMap["auth_type"] = apiAppApi.AuthType
		}

		if apiAppApi.ApiBusinessType != nil {
			apiInfoMap["api_business_type"] = apiAppApi.ApiBusinessType
		}

		if apiAppApi.AuthRelationApiId != nil {
			apiInfoMap["auth_relation_api_id"] = apiAppApi.AuthRelationApiId
		}

		if apiAppApi.OauthConfig != nil {
			oauthConfigMap := map[string]interface{}{}

			if apiAppApi.OauthConfig.PublicKey != nil {
				oauthConfigMap["public_key"] = apiAppApi.OauthConfig.PublicKey
			}

			if apiAppApi.OauthConfig.TokenLocation != nil {
				oauthConfigMap["token_location"] = apiAppApi.OauthConfig.TokenLocation
			}

			if apiAppApi.OauthConfig.LoginRedirectUrl != nil {
				oauthConfigMap["login_redirect_url"] = apiAppApi.OauthConfig.LoginRedirectUrl
			}

			apiInfoMap["oauth_config"] = []interface{}{oauthConfigMap}
		}

		if apiAppApi.IsDebugAfterCharge != nil {
			apiInfoMap["is_debug_after_charge"] = apiAppApi.IsDebugAfterCharge
		}

		if apiAppApi.RequestConfig != nil {
			requestConfigMap := map[string]interface{}{}

			if apiAppApi.RequestConfig.Path != nil {
				requestConfigMap["path"] = apiAppApi.RequestConfig.Path
			}

			if apiAppApi.RequestConfig.Method != nil {
				requestConfigMap["method"] = apiAppApi.RequestConfig.Method
			}

			apiInfoMap["request_config"] = []interface{}{requestConfigMap}
		}

		if apiAppApi.ResponseType != nil {
			apiInfoMap["response_type"] = apiAppApi.ResponseType
		}

		if apiAppApi.ResponseSuccessExample != nil {
			apiInfoMap["response_success_example"] = apiAppApi.ResponseSuccessExample
		}

		if apiAppApi.ResponseFailExample != nil {
			apiInfoMap["response_fail_example"] = apiAppApi.ResponseFailExample
		}

		if apiAppApi.ResponseErrorCodes != nil {
			responseErrorCodesList := []interface{}{}
			for _, responseErrorCodes := range apiAppApi.ResponseErrorCodes {
				responseErrorCodesMap := map[string]interface{}{}

				if responseErrorCodes.Code != nil {
					responseErrorCodesMap["code"] = responseErrorCodes.Code
				}

				if responseErrorCodes.Msg != nil {
					responseErrorCodesMap["msg"] = responseErrorCodes.Msg
				}

				if responseErrorCodes.Desc != nil {
					responseErrorCodesMap["desc"] = responseErrorCodes.Desc
				}

				if responseErrorCodes.ConvertedCode != nil {
					responseErrorCodesMap["converted_code"] = responseErrorCodes.ConvertedCode
				}

				if responseErrorCodes.NeedConvert != nil {
					responseErrorCodesMap["need_convert"] = responseErrorCodes.NeedConvert
				}

				responseErrorCodesList = append(responseErrorCodesList, responseErrorCodesMap)
			}

			apiInfoMap["response_error_codes"] = responseErrorCodesList
		}

		if apiAppApi.RequestParameters != nil {
			requestParametersList := []interface{}{}
			for _, requestParameters := range apiAppApi.RequestParameters {
				requestParametersMap := map[string]interface{}{}

				if requestParameters.Name != nil {
					requestParametersMap["name"] = requestParameters.Name
				}

				if requestParameters.Position != nil {
					requestParametersMap["position"] = requestParameters.Position
				}

				if requestParameters.Type != nil {
					requestParametersMap["type"] = requestParameters.Type
				}

				if requestParameters.DefaultValue != nil {
					requestParametersMap["default_value"] = requestParameters.DefaultValue
				}

				if requestParameters.Required != nil {
					requestParametersMap["required"] = requestParameters.Required
				}

				if requestParameters.Desc != nil {
					requestParametersMap["desc"] = requestParameters.Desc
				}

				requestParametersList = append(requestParametersList, requestParametersMap)
			}

			apiInfoMap["request_parameters"] = requestParametersList
		}

		if apiAppApi.ServiceTimeout != nil {
			apiInfoMap["service_timeout"] = apiAppApi.ServiceTimeout
		}

		if apiAppApi.ServiceType != nil {
			apiInfoMap["service_type"] = apiAppApi.ServiceType
		}

		if apiAppApi.ServiceConfig != nil {
			serviceConfigMap := map[string]interface{}{}

			if apiAppApi.ServiceConfig.Product != nil {
				serviceConfigMap["product"] = apiAppApi.ServiceConfig.Product
			}

			if apiAppApi.ServiceConfig.UniqVpcId != nil {
				serviceConfigMap["uniq_vpc_id"] = apiAppApi.ServiceConfig.UniqVpcId
			}

			if apiAppApi.ServiceConfig.Url != nil {
				serviceConfigMap["url"] = apiAppApi.ServiceConfig.Url
			}

			if apiAppApi.ServiceConfig.Path != nil {
				serviceConfigMap["path"] = apiAppApi.ServiceConfig.Path
			}

			if apiAppApi.ServiceConfig.Method != nil {
				serviceConfigMap["method"] = apiAppApi.ServiceConfig.Method
			}

			if apiAppApi.ServiceConfig.UpstreamId != nil {
				serviceConfigMap["upstream_id"] = apiAppApi.ServiceConfig.UpstreamId
			}

			apiInfoMap["service_config"] = []interface{}{serviceConfigMap}
		}

		if apiAppApi.ServiceParameters != nil {
			serviceParametersList := []interface{}{}
			for _, serviceParameters := range apiAppApi.ServiceParameters {
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

			apiInfoMap["service_parameters"] = serviceParametersList
		}

		if apiAppApi.ConstantParameters != nil {
			constantParametersList := []interface{}{}
			for _, constantParameters := range apiAppApi.ConstantParameters {
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

			apiInfoMap["constant_parameters"] = constantParametersList
		}

		if apiAppApi.ServiceMockReturnMessage != nil {
			apiInfoMap["service_mock_return_message"] = apiAppApi.ServiceMockReturnMessage
		}

		if apiAppApi.ServiceScfFunctionName != nil {
			apiInfoMap["service_scf_function_name"] = apiAppApi.ServiceScfFunctionName
		}

		if apiAppApi.ServiceScfFunctionNamespace != nil {
			apiInfoMap["service_scf_function_namespace"] = apiAppApi.ServiceScfFunctionNamespace
		}

		if apiAppApi.ServiceScfFunctionQualifier != nil {
			apiInfoMap["service_scf_function_qualifier"] = apiAppApi.ServiceScfFunctionQualifier
		}

		if apiAppApi.ServiceScfIsIntegratedResponse != nil {
			apiInfoMap["service_scf_is_integrated_response"] = apiAppApi.ServiceScfIsIntegratedResponse
		}

		if apiAppApi.ServiceWebsocketRegisterFunctionName != nil {
			apiInfoMap["service_websocket_register_function_name"] = apiAppApi.ServiceWebsocketRegisterFunctionName
		}

		if apiAppApi.ServiceWebsocketRegisterFunctionNamespace != nil {
			apiInfoMap["service_websocket_register_function_namespace"] = apiAppApi.ServiceWebsocketRegisterFunctionNamespace
		}

		if apiAppApi.ServiceWebsocketRegisterFunctionQualifier != nil {
			apiInfoMap["service_websocket_register_function_qualifier"] = apiAppApi.ServiceWebsocketRegisterFunctionQualifier
		}

		if apiAppApi.ServiceWebsocketCleanupFunctionName != nil {
			apiInfoMap["service_websocket_cleanup_function_name"] = apiAppApi.ServiceWebsocketCleanupFunctionName
		}

		if apiAppApi.ServiceWebsocketCleanupFunctionNamespace != nil {
			apiInfoMap["service_websocket_cleanup_function_namespace"] = apiAppApi.ServiceWebsocketCleanupFunctionNamespace
		}

		if apiAppApi.ServiceWebsocketCleanupFunctionQualifier != nil {
			apiInfoMap["service_websocket_cleanup_function_qualifier"] = apiAppApi.ServiceWebsocketCleanupFunctionQualifier
		}

		if apiAppApi.InternalDomain != nil {
			apiInfoMap["internal_domain"] = apiAppApi.InternalDomain
		}

		if apiAppApi.ServiceWebsocketTransportFunctionName != nil {
			apiInfoMap["service_websocket_transport_function_name"] = apiAppApi.ServiceWebsocketTransportFunctionName
		}

		if apiAppApi.ServiceWebsocketTransportFunctionNamespace != nil {
			apiInfoMap["service_websocket_transport_function_namespace"] = apiAppApi.ServiceWebsocketTransportFunctionNamespace
		}

		if apiAppApi.ServiceWebsocketTransportFunctionQualifier != nil {
			apiInfoMap["service_websocket_transport_function_qualifier"] = apiAppApi.ServiceWebsocketTransportFunctionQualifier
		}

		if apiAppApi.MicroServices != nil {
			microServicesList := []interface{}{}
			for _, microServices := range apiAppApi.MicroServices {
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

			apiInfoMap["micro_services"] = microServicesList
		}

		if apiAppApi.MicroServicesInfo != nil {
			apiInfoMap["micro_services_info"] = apiAppApi.MicroServicesInfo
		}

		if apiAppApi.ServiceTsfLoadBalanceConf != nil {
			serviceTsfLoadBalanceConfMap := map[string]interface{}{}

			if apiAppApi.ServiceTsfLoadBalanceConf.IsLoadBalance != nil {
				serviceTsfLoadBalanceConfMap["is_load_balance"] = apiAppApi.ServiceTsfLoadBalanceConf.IsLoadBalance
			}

			if apiAppApi.ServiceTsfLoadBalanceConf.Method != nil {
				serviceTsfLoadBalanceConfMap["method"] = apiAppApi.ServiceTsfLoadBalanceConf.Method
			}

			if apiAppApi.ServiceTsfLoadBalanceConf.SessionStickRequired != nil {
				serviceTsfLoadBalanceConfMap["session_stick_required"] = apiAppApi.ServiceTsfLoadBalanceConf.SessionStickRequired
			}

			if apiAppApi.ServiceTsfLoadBalanceConf.SessionStickTimeout != nil {
				serviceTsfLoadBalanceConfMap["session_stick_timeout"] = apiAppApi.ServiceTsfLoadBalanceConf.SessionStickTimeout
			}

			apiInfoMap["service_tsf_load_balance_conf"] = []interface{}{serviceTsfLoadBalanceConfMap}
		}

		if apiAppApi.ServiceTsfHealthCheckConf != nil {
			serviceTsfHealthCheckConfMap := map[string]interface{}{}

			if apiAppApi.ServiceTsfHealthCheckConf.IsHealthCheck != nil {
				serviceTsfHealthCheckConfMap["is_health_check"] = apiAppApi.ServiceTsfHealthCheckConf.IsHealthCheck
			}

			if apiAppApi.ServiceTsfHealthCheckConf.RequestVolumeThreshold != nil {
				serviceTsfHealthCheckConfMap["request_volume_threshold"] = apiAppApi.ServiceTsfHealthCheckConf.RequestVolumeThreshold
			}

			if apiAppApi.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds != nil {
				serviceTsfHealthCheckConfMap["sleep_window_in_milliseconds"] = apiAppApi.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds
			}

			if apiAppApi.ServiceTsfHealthCheckConf.ErrorThresholdPercentage != nil {
				serviceTsfHealthCheckConfMap["error_threshold_percentage"] = apiAppApi.ServiceTsfHealthCheckConf.ErrorThresholdPercentage
			}

			apiInfoMap["service_tsf_health_check_conf"] = []interface{}{serviceTsfHealthCheckConfMap}
		}

		if apiAppApi.EnableCORS != nil {
			apiInfoMap["enable_cors"] = apiAppApi.EnableCORS
		}

		if apiAppApi.Tags != nil {
			tagsList := []interface{}{}
			for _, tags := range apiAppApi.Tags {
				tagsMap := map[string]interface{}{}

				if tags.Key != nil {
					tagsMap["key"] = tags.Key
				}

				if tags.Value != nil {
					tagsMap["value"] = tags.Value
				}

				tagsList = append(tagsList, tagsMap)
			}

			apiInfoMap["tags"] = tagsList
		}

		if apiAppApi.Environments != nil {
			apiInfoMap["environments"] = apiAppApi.Environments
		}

		if apiAppApi.IsBase64Encoded != nil {
			apiInfoMap["is_base64_encoded"] = apiAppApi.IsBase64Encoded
		}

		if apiAppApi.IsBase64Trigger != nil {
			apiInfoMap["is_base64_trigger"] = apiAppApi.IsBase64Trigger
		}

		if apiAppApi.Base64EncodedTriggerRules != nil {
			base64EncodedTriggerRulesList := []interface{}{}
			for _, base64EncodedTriggerRules := range apiAppApi.Base64EncodedTriggerRules {
				base64EncodedTriggerRulesMap := map[string]interface{}{}

				if base64EncodedTriggerRules.Name != nil {
					base64EncodedTriggerRulesMap["name"] = base64EncodedTriggerRules.Name
				}

				if base64EncodedTriggerRules.Value != nil {
					base64EncodedTriggerRulesMap["value"] = base64EncodedTriggerRules.Value
				}

				base64EncodedTriggerRulesList = append(base64EncodedTriggerRulesList, base64EncodedTriggerRulesMap)
			}

			apiInfoMap["base64_encoded_trigger_rules"] = base64EncodedTriggerRulesList
		}

		tmpList = append(tmpList, apiInfoMap)
		_ = d.Set("result", tmpList)
	}

	d.SetId(strings.Join([]string{service_id, api_id, api_region}, tccommon.FILED_SP))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

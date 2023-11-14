/*
Provides a resource to create a apigateway a_p_i

Example Usage

```hcl
resource "tencentcloud_apigateway_a_p_i" "a_p_i" {
  service_id = ""
  service_type = ""
  service_timeout =
  protocol = ""
  request_config {
		path = ""
		method = ""

  }
  api_name = ""
  api_desc = ""
  api_type = ""
  auth_type = ""
  enable_c_o_r_s =
  constant_parameters {
		name = ""
		desc = ""
		position = ""
		default_value = ""

  }
  request_parameters {
		name = ""
		desc = ""
		position = ""
		type = ""
		default_value = ""
		required =

  }
  api_business_type = ""
  service_mock_return_message = ""
  micro_services {
		cluster_id = ""
		namespace_id = ""
		micro_service_name = ""

  }
  service_tsf_load_balance_conf {
		is_load_balance =
		method = ""
		session_stick_required =
		session_stick_timeout =

  }
  service_tsf_health_check_conf {
		is_health_check =
		request_volume_threshold =
		sleep_window_in_milliseconds =
		error_threshold_percentage =

  }
  target_services {
		vm_ip = ""
		vpc_id = ""
		vm_port =
		host_ip = ""
		docker_ip = ""

  }
  target_services_load_balance_conf =
  target_services_health_check_conf {
		is_health_check =
		request_volume_threshold =
		sleep_window_in_milliseconds =
		error_threshold_percentage =

  }
  service_scf_function_name = ""
  service_websocket_register_function_name = ""
  service_websocket_cleanup_function_name = ""
  service_websocket_transport_function_name = ""
  service_scf_function_namespace = ""
  service_scf_function_qualifier = ""
  service_websocket_register_function_namespace = ""
  service_websocket_register_function_qualifier = ""
  service_websocket_transport_function_namespace = ""
  service_websocket_transport_function_qualifier = ""
  service_websocket_cleanup_function_namespace = ""
  service_websocket_cleanup_function_qualifier = ""
  service_scf_is_integrated_response =
  is_debug_after_charge =
  is_delete_response_error_codes =
  response_type = ""
  response_success_example = ""
  response_fail_example = ""
  service_config {
		product = ""
		uniq_vpc_id = ""
		url = ""
		path = ""
		method = ""
		upstream_id = ""
		cos_config {
			action = ""
			bucket_name = ""
			authorization =
			path_match_mode = ""
		}

  }
  auth_relation_api_id = ""
  service_parameters {
		name = ""
		position = ""
		relevant_request_parameter_position = ""
		relevant_request_parameter_name = ""
		default_value = ""
		relevant_request_parameter_desc = ""
		relevant_request_parameter_type = ""

  }
  oauth_config {
		public_key = ""
		token_location = ""
		login_redirect_url = ""

  }
  response_error_codes {
		code =
		msg = ""
		desc = ""
		converted_code =
		need_convert =

  }
  target_namespace_id = ""
  user_type = ""
  is_base64_encoded =
  event_bus_id = ""
  service_scf_function_type = ""
  e_i_a_m_app_type = ""
  e_i_a_m_auth_type = ""
  token_timeout =
  e_i_a_m_app_id = ""
  owner = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

apigateway a_p_i can be imported using the id, e.g.

```
terraform import tencentcloud_apigateway_a_p_i.a_p_i a_p_i_id
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
	"log"
)

func resourceTencentCloudApigatewayAPI() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApigatewayAPICreate,
		Read:   resourceTencentCloudApigatewayAPIRead,
		Update: resourceTencentCloudApigatewayAPIUpdate,
		Delete: resourceTencentCloudApigatewayAPIDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the service where the API is located.",
			},

			"service_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The backend service type of the API. Supports HTTP、MOCK、TSF、SCF、EB、TARGET、VPC、UPSTREAM、GRPC、COS、WEBSOCKET.",
			},

			"service_timeout": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The timeout time of the API&amp;amp;#39;s backend service, in seconds.",
			},

			"protocol": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The timeout time of the API&amp;amp;#39;s backend service, in seconds.",
			},

			"request_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The front-end configuration of the request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path.",
						},
						"method": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Method.",
						},
					},
				},
			},

			"api_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User defined API name.",
			},

			"api_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User defined API interface description.",
			},

			"api_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "API type, supports NORMAL (regular API) and TSF (microservice API), defaults to NORMAL.",
			},

			"auth_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "API authentication type. Support SECRET (Key Pair Authentication), NONE (Authentication Exemption), OAUTH, APP (Application Authentication). The default is NONE.",
			},

			"enable_c_o_r_s": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable cross domain.",
			},

			"constant_parameters": {
				Optional:    true,
				Type:        schema.TypeList,
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
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Front end request parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request Parameter Name.",
						},
						"desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Describe.",
						},
						"position": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter position.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default value.",
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Required.",
						},
					},
				},
			},

			"api_business_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "When AuthType is OAUTH, this field is valid, NORMAL: Business API OAUTH: Authorization API.",
			},

			"service_mock_return_message": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The backend mock of the API returns information. If ServiceType is Mock, this parameter must be passed.",
			},

			"micro_services": {
				Optional:    true,
				Type:        schema.TypeList,
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
							Description: "Vm ip.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vpc id.",
						},
						"vm_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Vm port.",
						},
						"host_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Host IP of the CVM.",
						},
						"docker_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Docker ip.",
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

			"service_scf_function_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf function name. Effective when the backend type is SCF.",
			},

			"service_websocket_register_function_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket registration function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},

			"service_websocket_cleanup_function_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleaning function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},

			"service_websocket_transport_function_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},

			"service_scf_function_namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "SCF function namespace. Effective when the backend type is SCF.",
			},

			"service_scf_function_qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf function version. Effective when the backend type is SCF.",
			},

			"service_websocket_register_function_namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket registers function namespaces. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},

			"service_websocket_register_function_qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},

			"service_websocket_transport_function_namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},

			"service_websocket_transport_function_qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},

			"service_websocket_cleanup_function_namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleans up the function namespace. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},

			"service_websocket_cleanup_function_qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleaning function version. It takes effect when the current end type is WEBSOCKET and the backend type is SCF.",
			},

			"service_scf_is_integrated_response": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable response integration. Effective when the backend type is SCF.",
			},

			"is_debug_after_charge": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Charge after starting debugging. (Cloud Market Reserved Fields).",
			},

			"is_delete_response_error_codes": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Do you want to delete the custom response configuration error code? If it is not passed or False is passed, it will not be deleted. If True is passed, all custom response configuration error codes for this API will be deleted.",
			},

			"response_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Return type.",
			},

			"response_success_example": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Custom response configuration successful response example.",
			},

			"response_fail_example": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Custom response configuration failed response example.",
			},

			"service_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "API backend service configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Backend type. Effective when enabling vpc, currently supported types are clb, cvm, and upstream.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The unique ID of the vpc.",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The backend service URL of the API. If the ServiceType is HTTP, this parameter must be passed.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The backend service path of the API, such as/path. If the ServiceType is HTTP, this parameter must be passed. The front and rear paths can be different.",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The backend service request method for APIs, such as GET. If the ServiceType is HTTP, this parameter must be passed. The front and rear methods can be different.",
						},
						"upstream_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Only required when binding to VPC channelsNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"cos_config": {
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
					},
				},
			},

			"auth_relation_api_id": {
				Optional:    true,
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

			"response_error_codes": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "User defined error code configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Custom response configuration error code.",
						},
						"msg": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom response configuration error information.",
						},
						"desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom response configuration error code comments.",
						},
						"converted_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Custom error code conversion.",
						},
						"need_convert": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Do you need to enable error code conversion.",
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
				Type:        schema.TypeBool,
				Description: "Whether to enable Base64 encoding will only take effect when the backend is scf.",
			},

			"event_bus_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Event bus ID.",
			},

			"service_scf_function_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Scf function type. Effective when the backend type is SCF. Support Event Triggering (EVENT) and HTTP Direct Cloud Function (HTTP).",
			},

			"e_i_a_m_app_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "EIAM application type.",
			},

			"e_i_a_m_auth_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The EIAM application authentication type supports AuthenticationOnly, Authentication, and Authorization.",
			},

			"token_timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The effective time of the EIAM application token, measured in seconds, defaults to 7200 seconds.",
			},

			"e_i_a_m_app_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "EIAM application ID.",
			},

			"owner": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Owner of resources.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudApigatewayAPICreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_a_p_i.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = apigateway.NewCreateApiRequest()
		response = apigateway.NewCreateApiResponse()
		apiId    string
	)
	if v, ok := d.GetOk("service_id"); ok {
		request.ServiceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_type"); ok {
		request.ServiceType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("service_timeout"); ok {
		request.ServiceTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "request_config"); ok {
		apiRequestConfig := apigateway.ApiRequestConfig{}
		if v, ok := dMap["path"]; ok {
			apiRequestConfig.Path = helper.String(v.(string))
		}
		if v, ok := dMap["method"]; ok {
			apiRequestConfig.Method = helper.String(v.(string))
		}
		request.RequestConfig = &apiRequestConfig
	}

	if v, ok := d.GetOk("api_name"); ok {
		request.ApiName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("api_desc"); ok {
		request.ApiDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("api_type"); ok {
		request.ApiType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_type"); ok {
		request.AuthType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_c_o_r_s"); ok {
		request.EnableCORS = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("constant_parameters"); ok {
		for _, item := range v.([]interface{}) {
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

	if v, ok := d.GetOk("request_parameters"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			requestParameter := apigateway.RequestParameter{}
			if v, ok := dMap["name"]; ok {
				requestParameter.Name = helper.String(v.(string))
			}
			if v, ok := dMap["desc"]; ok {
				requestParameter.Desc = helper.String(v.(string))
			}
			if v, ok := dMap["position"]; ok {
				requestParameter.Position = helper.String(v.(string))
			}
			if v, ok := dMap["type"]; ok {
				requestParameter.Type = helper.String(v.(string))
			}
			if v, ok := dMap["default_value"]; ok {
				requestParameter.DefaultValue = helper.String(v.(string))
			}
			if v, ok := dMap["required"]; ok {
				requestParameter.Required = helper.Bool(v.(bool))
			}
			request.RequestParameters = append(request.RequestParameters, &requestParameter)
		}
	}

	if v, ok := d.GetOk("api_business_type"); ok {
		request.ApiBusinessType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_mock_return_message"); ok {
		request.ServiceMockReturnMessage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("micro_services"); ok {
		for _, item := range v.([]interface{}) {
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
		for _, item := range v.([]interface{}) {
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

	if v, ok := d.GetOk("service_scf_function_name"); ok {
		request.ServiceScfFunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_websocket_register_function_name"); ok {
		request.ServiceWebsocketRegisterFunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_websocket_cleanup_function_name"); ok {
		request.ServiceWebsocketCleanupFunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_websocket_transport_function_name"); ok {
		request.ServiceWebsocketTransportFunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_scf_function_namespace"); ok {
		request.ServiceScfFunctionNamespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_scf_function_qualifier"); ok {
		request.ServiceScfFunctionQualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_websocket_register_function_namespace"); ok {
		request.ServiceWebsocketRegisterFunctionNamespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_websocket_register_function_qualifier"); ok {
		request.ServiceWebsocketRegisterFunctionQualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_websocket_transport_function_namespace"); ok {
		request.ServiceWebsocketTransportFunctionNamespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_websocket_transport_function_qualifier"); ok {
		request.ServiceWebsocketTransportFunctionQualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_websocket_cleanup_function_namespace"); ok {
		request.ServiceWebsocketCleanupFunctionNamespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_websocket_cleanup_function_qualifier"); ok {
		request.ServiceWebsocketCleanupFunctionQualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("service_scf_is_integrated_response"); ok {
		request.ServiceScfIsIntegratedResponse = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_debug_after_charge"); ok {
		request.IsDebugAfterCharge = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_delete_response_error_codes"); ok {
		request.IsDeleteResponseErrorCodes = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("response_type"); ok {
		request.ResponseType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("response_success_example"); ok {
		request.ResponseSuccessExample = helper.String(v.(string))
	}

	if v, ok := d.GetOk("response_fail_example"); ok {
		request.ResponseFailExample = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "service_config"); ok {
		serviceConfig := apigateway.ServiceConfig{}
		if v, ok := dMap["product"]; ok {
			serviceConfig.Product = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpc_id"]; ok {
			serviceConfig.UniqVpcId = helper.String(v.(string))
		}
		if v, ok := dMap["url"]; ok {
			serviceConfig.Url = helper.String(v.(string))
		}
		if v, ok := dMap["path"]; ok {
			serviceConfig.Path = helper.String(v.(string))
		}
		if v, ok := dMap["method"]; ok {
			serviceConfig.Method = helper.String(v.(string))
		}
		if v, ok := dMap["upstream_id"]; ok {
			serviceConfig.UpstreamId = helper.String(v.(string))
		}
		if cosConfigMap, ok := helper.InterfaceToMap(dMap, "cos_config"); ok {
			cosConfig := apigateway.CosConfig{}
			if v, ok := cosConfigMap["action"]; ok {
				cosConfig.Action = helper.String(v.(string))
			}
			if v, ok := cosConfigMap["bucket_name"]; ok {
				cosConfig.BucketName = helper.String(v.(string))
			}
			if v, ok := cosConfigMap["authorization"]; ok {
				cosConfig.Authorization = helper.Bool(v.(bool))
			}
			if v, ok := cosConfigMap["path_match_mode"]; ok {
				cosConfig.PathMatchMode = helper.String(v.(string))
			}
			serviceConfig.CosConfig = &cosConfig
		}
		request.ServiceConfig = &serviceConfig
	}

	if v, ok := d.GetOk("auth_relation_api_id"); ok {
		request.AuthRelationApiId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_parameters"); ok {
		for _, item := range v.([]interface{}) {
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

	if v, ok := d.GetOk("response_error_codes"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			responseErrorCodeReq := apigateway.ResponseErrorCodeReq{}
			if v, ok := dMap["code"]; ok {
				responseErrorCodeReq.Code = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["msg"]; ok {
				responseErrorCodeReq.Msg = helper.String(v.(string))
			}
			if v, ok := dMap["desc"]; ok {
				responseErrorCodeReq.Desc = helper.String(v.(string))
			}
			if v, ok := dMap["converted_code"]; ok {
				responseErrorCodeReq.ConvertedCode = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["need_convert"]; ok {
				responseErrorCodeReq.NeedConvert = helper.Bool(v.(bool))
			}
			request.ResponseErrorCodes = append(request.ResponseErrorCodes, &responseErrorCodeReq)
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

	if v, ok := d.GetOk("e_i_a_m_app_type"); ok {
		request.EIAMAppType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("e_i_a_m_auth_type"); ok {
		request.EIAMAuthType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("token_timeout"); ok {
		request.TokenTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("e_i_a_m_app_id"); ok {
		request.EIAMAppId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner"); ok {
		request.Owner = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().CreateApi(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create apigateway API failed, reason:%+v", logId, err)
		return err
	}

	apiId = *response.Response.ApiId
	d.SetId(apiId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::apigw:%s:uin/:apiId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayAPIRead(d, meta)
}

func resourceTencentCloudApigatewayAPIRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_a_p_i.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	aPIId := d.Id()

	API, err := service.DescribeApigatewayAPIById(ctx, apiId)
	if err != nil {
		return err
	}

	if API == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApigatewayAPI` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if API.ServiceId != nil {
		_ = d.Set("service_id", API.ServiceId)
	}

	if API.ServiceType != nil {
		_ = d.Set("service_type", API.ServiceType)
	}

	if API.ServiceTimeout != nil {
		_ = d.Set("service_timeout", API.ServiceTimeout)
	}

	if API.Protocol != nil {
		_ = d.Set("protocol", API.Protocol)
	}

	if API.RequestConfig != nil {
		requestConfigMap := map[string]interface{}{}

		if API.RequestConfig.Path != nil {
			requestConfigMap["path"] = API.RequestConfig.Path
		}

		if API.RequestConfig.Method != nil {
			requestConfigMap["method"] = API.RequestConfig.Method
		}

		_ = d.Set("request_config", []interface{}{requestConfigMap})
	}

	if API.ApiName != nil {
		_ = d.Set("api_name", API.ApiName)
	}

	if API.ApiDesc != nil {
		_ = d.Set("api_desc", API.ApiDesc)
	}

	if API.ApiType != nil {
		_ = d.Set("api_type", API.ApiType)
	}

	if API.AuthType != nil {
		_ = d.Set("auth_type", API.AuthType)
	}

	if API.EnableCORS != nil {
		_ = d.Set("enable_c_o_r_s", API.EnableCORS)
	}

	if API.ConstantParameters != nil {
		constantParametersList := []interface{}{}
		for _, constantParameters := range API.ConstantParameters {
			constantParametersMap := map[string]interface{}{}

			if API.ConstantParameters.Name != nil {
				constantParametersMap["name"] = API.ConstantParameters.Name
			}

			if API.ConstantParameters.Desc != nil {
				constantParametersMap["desc"] = API.ConstantParameters.Desc
			}

			if API.ConstantParameters.Position != nil {
				constantParametersMap["position"] = API.ConstantParameters.Position
			}

			if API.ConstantParameters.DefaultValue != nil {
				constantParametersMap["default_value"] = API.ConstantParameters.DefaultValue
			}

			constantParametersList = append(constantParametersList, constantParametersMap)
		}

		_ = d.Set("constant_parameters", constantParametersList)

	}

	if API.RequestParameters != nil {
		requestParametersList := []interface{}{}
		for _, requestParameters := range API.RequestParameters {
			requestParametersMap := map[string]interface{}{}

			if API.RequestParameters.Name != nil {
				requestParametersMap["name"] = API.RequestParameters.Name
			}

			if API.RequestParameters.Desc != nil {
				requestParametersMap["desc"] = API.RequestParameters.Desc
			}

			if API.RequestParameters.Position != nil {
				requestParametersMap["position"] = API.RequestParameters.Position
			}

			if API.RequestParameters.Type != nil {
				requestParametersMap["type"] = API.RequestParameters.Type
			}

			if API.RequestParameters.DefaultValue != nil {
				requestParametersMap["default_value"] = API.RequestParameters.DefaultValue
			}

			if API.RequestParameters.Required != nil {
				requestParametersMap["required"] = API.RequestParameters.Required
			}

			requestParametersList = append(requestParametersList, requestParametersMap)
		}

		_ = d.Set("request_parameters", requestParametersList)

	}

	if API.ApiBusinessType != nil {
		_ = d.Set("api_business_type", API.ApiBusinessType)
	}

	if API.ServiceMockReturnMessage != nil {
		_ = d.Set("service_mock_return_message", API.ServiceMockReturnMessage)
	}

	if API.MicroServices != nil {
		microServicesList := []interface{}{}
		for _, microServices := range API.MicroServices {
			microServicesMap := map[string]interface{}{}

			if API.MicroServices.ClusterId != nil {
				microServicesMap["cluster_id"] = API.MicroServices.ClusterId
			}

			if API.MicroServices.NamespaceId != nil {
				microServicesMap["namespace_id"] = API.MicroServices.NamespaceId
			}

			if API.MicroServices.MicroServiceName != nil {
				microServicesMap["micro_service_name"] = API.MicroServices.MicroServiceName
			}

			microServicesList = append(microServicesList, microServicesMap)
		}

		_ = d.Set("micro_services", microServicesList)

	}

	if API.ServiceTsfLoadBalanceConf != nil {
		serviceTsfLoadBalanceConfMap := map[string]interface{}{}

		if API.ServiceTsfLoadBalanceConf.IsLoadBalance != nil {
			serviceTsfLoadBalanceConfMap["is_load_balance"] = API.ServiceTsfLoadBalanceConf.IsLoadBalance
		}

		if API.ServiceTsfLoadBalanceConf.Method != nil {
			serviceTsfLoadBalanceConfMap["method"] = API.ServiceTsfLoadBalanceConf.Method
		}

		if API.ServiceTsfLoadBalanceConf.SessionStickRequired != nil {
			serviceTsfLoadBalanceConfMap["session_stick_required"] = API.ServiceTsfLoadBalanceConf.SessionStickRequired
		}

		if API.ServiceTsfLoadBalanceConf.SessionStickTimeout != nil {
			serviceTsfLoadBalanceConfMap["session_stick_timeout"] = API.ServiceTsfLoadBalanceConf.SessionStickTimeout
		}

		_ = d.Set("service_tsf_load_balance_conf", []interface{}{serviceTsfLoadBalanceConfMap})
	}

	if API.ServiceTsfHealthCheckConf != nil {
		serviceTsfHealthCheckConfMap := map[string]interface{}{}

		if API.ServiceTsfHealthCheckConf.IsHealthCheck != nil {
			serviceTsfHealthCheckConfMap["is_health_check"] = API.ServiceTsfHealthCheckConf.IsHealthCheck
		}

		if API.ServiceTsfHealthCheckConf.RequestVolumeThreshold != nil {
			serviceTsfHealthCheckConfMap["request_volume_threshold"] = API.ServiceTsfHealthCheckConf.RequestVolumeThreshold
		}

		if API.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds != nil {
			serviceTsfHealthCheckConfMap["sleep_window_in_milliseconds"] = API.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds
		}

		if API.ServiceTsfHealthCheckConf.ErrorThresholdPercentage != nil {
			serviceTsfHealthCheckConfMap["error_threshold_percentage"] = API.ServiceTsfHealthCheckConf.ErrorThresholdPercentage
		}

		_ = d.Set("service_tsf_health_check_conf", []interface{}{serviceTsfHealthCheckConfMap})
	}

	if API.TargetServices != nil {
		targetServicesList := []interface{}{}
		for _, targetServices := range API.TargetServices {
			targetServicesMap := map[string]interface{}{}

			if API.TargetServices.VmIp != nil {
				targetServicesMap["vm_ip"] = API.TargetServices.VmIp
			}

			if API.TargetServices.VpcId != nil {
				targetServicesMap["vpc_id"] = API.TargetServices.VpcId
			}

			if API.TargetServices.VmPort != nil {
				targetServicesMap["vm_port"] = API.TargetServices.VmPort
			}

			if API.TargetServices.HostIp != nil {
				targetServicesMap["host_ip"] = API.TargetServices.HostIp
			}

			if API.TargetServices.DockerIp != nil {
				targetServicesMap["docker_ip"] = API.TargetServices.DockerIp
			}

			targetServicesList = append(targetServicesList, targetServicesMap)
		}

		_ = d.Set("target_services", targetServicesList)

	}

	if API.TargetServicesLoadBalanceConf != nil {
		_ = d.Set("target_services_load_balance_conf", API.TargetServicesLoadBalanceConf)
	}

	if API.TargetServicesHealthCheckConf != nil {
		targetServicesHealthCheckConfMap := map[string]interface{}{}

		if API.TargetServicesHealthCheckConf.IsHealthCheck != nil {
			targetServicesHealthCheckConfMap["is_health_check"] = API.TargetServicesHealthCheckConf.IsHealthCheck
		}

		if API.TargetServicesHealthCheckConf.RequestVolumeThreshold != nil {
			targetServicesHealthCheckConfMap["request_volume_threshold"] = API.TargetServicesHealthCheckConf.RequestVolumeThreshold
		}

		if API.TargetServicesHealthCheckConf.SleepWindowInMilliseconds != nil {
			targetServicesHealthCheckConfMap["sleep_window_in_milliseconds"] = API.TargetServicesHealthCheckConf.SleepWindowInMilliseconds
		}

		if API.TargetServicesHealthCheckConf.ErrorThresholdPercentage != nil {
			targetServicesHealthCheckConfMap["error_threshold_percentage"] = API.TargetServicesHealthCheckConf.ErrorThresholdPercentage
		}

		_ = d.Set("target_services_health_check_conf", []interface{}{targetServicesHealthCheckConfMap})
	}

	if API.ServiceScfFunctionName != nil {
		_ = d.Set("service_scf_function_name", API.ServiceScfFunctionName)
	}

	if API.ServiceWebsocketRegisterFunctionName != nil {
		_ = d.Set("service_websocket_register_function_name", API.ServiceWebsocketRegisterFunctionName)
	}

	if API.ServiceWebsocketCleanupFunctionName != nil {
		_ = d.Set("service_websocket_cleanup_function_name", API.ServiceWebsocketCleanupFunctionName)
	}

	if API.ServiceWebsocketTransportFunctionName != nil {
		_ = d.Set("service_websocket_transport_function_name", API.ServiceWebsocketTransportFunctionName)
	}

	if API.ServiceScfFunctionNamespace != nil {
		_ = d.Set("service_scf_function_namespace", API.ServiceScfFunctionNamespace)
	}

	if API.ServiceScfFunctionQualifier != nil {
		_ = d.Set("service_scf_function_qualifier", API.ServiceScfFunctionQualifier)
	}

	if API.ServiceWebsocketRegisterFunctionNamespace != nil {
		_ = d.Set("service_websocket_register_function_namespace", API.ServiceWebsocketRegisterFunctionNamespace)
	}

	if API.ServiceWebsocketRegisterFunctionQualifier != nil {
		_ = d.Set("service_websocket_register_function_qualifier", API.ServiceWebsocketRegisterFunctionQualifier)
	}

	if API.ServiceWebsocketTransportFunctionNamespace != nil {
		_ = d.Set("service_websocket_transport_function_namespace", API.ServiceWebsocketTransportFunctionNamespace)
	}

	if API.ServiceWebsocketTransportFunctionQualifier != nil {
		_ = d.Set("service_websocket_transport_function_qualifier", API.ServiceWebsocketTransportFunctionQualifier)
	}

	if API.ServiceWebsocketCleanupFunctionNamespace != nil {
		_ = d.Set("service_websocket_cleanup_function_namespace", API.ServiceWebsocketCleanupFunctionNamespace)
	}

	if API.ServiceWebsocketCleanupFunctionQualifier != nil {
		_ = d.Set("service_websocket_cleanup_function_qualifier", API.ServiceWebsocketCleanupFunctionQualifier)
	}

	if API.ServiceScfIsIntegratedResponse != nil {
		_ = d.Set("service_scf_is_integrated_response", API.ServiceScfIsIntegratedResponse)
	}

	if API.IsDebugAfterCharge != nil {
		_ = d.Set("is_debug_after_charge", API.IsDebugAfterCharge)
	}

	if API.IsDeleteResponseErrorCodes != nil {
		_ = d.Set("is_delete_response_error_codes", API.IsDeleteResponseErrorCodes)
	}

	if API.ResponseType != nil {
		_ = d.Set("response_type", API.ResponseType)
	}

	if API.ResponseSuccessExample != nil {
		_ = d.Set("response_success_example", API.ResponseSuccessExample)
	}

	if API.ResponseFailExample != nil {
		_ = d.Set("response_fail_example", API.ResponseFailExample)
	}

	if API.ServiceConfig != nil {
		serviceConfigMap := map[string]interface{}{}

		if API.ServiceConfig.Product != nil {
			serviceConfigMap["product"] = API.ServiceConfig.Product
		}

		if API.ServiceConfig.UniqVpcId != nil {
			serviceConfigMap["uniq_vpc_id"] = API.ServiceConfig.UniqVpcId
		}

		if API.ServiceConfig.Url != nil {
			serviceConfigMap["url"] = API.ServiceConfig.Url
		}

		if API.ServiceConfig.Path != nil {
			serviceConfigMap["path"] = API.ServiceConfig.Path
		}

		if API.ServiceConfig.Method != nil {
			serviceConfigMap["method"] = API.ServiceConfig.Method
		}

		if API.ServiceConfig.UpstreamId != nil {
			serviceConfigMap["upstream_id"] = API.ServiceConfig.UpstreamId
		}

		if API.ServiceConfig.CosConfig != nil {
			cosConfigMap := map[string]interface{}{}

			if API.ServiceConfig.CosConfig.Action != nil {
				cosConfigMap["action"] = API.ServiceConfig.CosConfig.Action
			}

			if API.ServiceConfig.CosConfig.BucketName != nil {
				cosConfigMap["bucket_name"] = API.ServiceConfig.CosConfig.BucketName
			}

			if API.ServiceConfig.CosConfig.Authorization != nil {
				cosConfigMap["authorization"] = API.ServiceConfig.CosConfig.Authorization
			}

			if API.ServiceConfig.CosConfig.PathMatchMode != nil {
				cosConfigMap["path_match_mode"] = API.ServiceConfig.CosConfig.PathMatchMode
			}

			serviceConfigMap["cos_config"] = []interface{}{cosConfigMap}
		}

		_ = d.Set("service_config", []interface{}{serviceConfigMap})
	}

	if API.AuthRelationApiId != nil {
		_ = d.Set("auth_relation_api_id", API.AuthRelationApiId)
	}

	if API.ServiceParameters != nil {
		serviceParametersList := []interface{}{}
		for _, serviceParameters := range API.ServiceParameters {
			serviceParametersMap := map[string]interface{}{}

			if API.ServiceParameters.Name != nil {
				serviceParametersMap["name"] = API.ServiceParameters.Name
			}

			if API.ServiceParameters.Position != nil {
				serviceParametersMap["position"] = API.ServiceParameters.Position
			}

			if API.ServiceParameters.RelevantRequestParameterPosition != nil {
				serviceParametersMap["relevant_request_parameter_position"] = API.ServiceParameters.RelevantRequestParameterPosition
			}

			if API.ServiceParameters.RelevantRequestParameterName != nil {
				serviceParametersMap["relevant_request_parameter_name"] = API.ServiceParameters.RelevantRequestParameterName
			}

			if API.ServiceParameters.DefaultValue != nil {
				serviceParametersMap["default_value"] = API.ServiceParameters.DefaultValue
			}

			if API.ServiceParameters.RelevantRequestParameterDesc != nil {
				serviceParametersMap["relevant_request_parameter_desc"] = API.ServiceParameters.RelevantRequestParameterDesc
			}

			if API.ServiceParameters.RelevantRequestParameterType != nil {
				serviceParametersMap["relevant_request_parameter_type"] = API.ServiceParameters.RelevantRequestParameterType
			}

			serviceParametersList = append(serviceParametersList, serviceParametersMap)
		}

		_ = d.Set("service_parameters", serviceParametersList)

	}

	if API.OauthConfig != nil {
		oauthConfigMap := map[string]interface{}{}

		if API.OauthConfig.PublicKey != nil {
			oauthConfigMap["public_key"] = API.OauthConfig.PublicKey
		}

		if API.OauthConfig.TokenLocation != nil {
			oauthConfigMap["token_location"] = API.OauthConfig.TokenLocation
		}

		if API.OauthConfig.LoginRedirectUrl != nil {
			oauthConfigMap["login_redirect_url"] = API.OauthConfig.LoginRedirectUrl
		}

		_ = d.Set("oauth_config", []interface{}{oauthConfigMap})
	}

	if API.ResponseErrorCodes != nil {
		responseErrorCodesList := []interface{}{}
		for _, responseErrorCodes := range API.ResponseErrorCodes {
			responseErrorCodesMap := map[string]interface{}{}

			if API.ResponseErrorCodes.Code != nil {
				responseErrorCodesMap["code"] = API.ResponseErrorCodes.Code
			}

			if API.ResponseErrorCodes.Msg != nil {
				responseErrorCodesMap["msg"] = API.ResponseErrorCodes.Msg
			}

			if API.ResponseErrorCodes.Desc != nil {
				responseErrorCodesMap["desc"] = API.ResponseErrorCodes.Desc
			}

			if API.ResponseErrorCodes.ConvertedCode != nil {
				responseErrorCodesMap["converted_code"] = API.ResponseErrorCodes.ConvertedCode
			}

			if API.ResponseErrorCodes.NeedConvert != nil {
				responseErrorCodesMap["need_convert"] = API.ResponseErrorCodes.NeedConvert
			}

			responseErrorCodesList = append(responseErrorCodesList, responseErrorCodesMap)
		}

		_ = d.Set("response_error_codes", responseErrorCodesList)

	}

	if API.TargetNamespaceId != nil {
		_ = d.Set("target_namespace_id", API.TargetNamespaceId)
	}

	if API.UserType != nil {
		_ = d.Set("user_type", API.UserType)
	}

	if API.IsBase64Encoded != nil {
		_ = d.Set("is_base64_encoded", API.IsBase64Encoded)
	}

	if API.EventBusId != nil {
		_ = d.Set("event_bus_id", API.EventBusId)
	}

	if API.ServiceScfFunctionType != nil {
		_ = d.Set("service_scf_function_type", API.ServiceScfFunctionType)
	}

	if API.EIAMAppType != nil {
		_ = d.Set("e_i_a_m_app_type", API.EIAMAppType)
	}

	if API.EIAMAuthType != nil {
		_ = d.Set("e_i_a_m_auth_type", API.EIAMAuthType)
	}

	if API.TokenTimeout != nil {
		_ = d.Set("token_timeout", API.TokenTimeout)
	}

	if API.EIAMAppId != nil {
		_ = d.Set("e_i_a_m_app_id", API.EIAMAppId)
	}

	if API.Owner != nil {
		_ = d.Set("owner", API.Owner)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "apigw", "apiId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudApigatewayAPIUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_a_p_i.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifyApiIncrementRequest  = apigateway.NewModifyApiIncrementRequest()
		modifyApiIncrementResponse = apigateway.NewModifyApiIncrementResponse()
	)

	aPIId := d.Id()

	request.ApiId = &apiId

	immutableArgs := []string{"service_id", "service_type", "service_timeout", "protocol", "request_config", "api_name", "api_desc", "api_type", "auth_type", "enable_c_o_r_s", "constant_parameters", "request_parameters", "api_business_type", "service_mock_return_message", "micro_services", "service_tsf_load_balance_conf", "service_tsf_health_check_conf", "target_services", "target_services_load_balance_conf", "target_services_health_check_conf", "service_scf_function_name", "service_websocket_register_function_name", "service_websocket_cleanup_function_name", "service_websocket_transport_function_name", "service_scf_function_namespace", "service_scf_function_qualifier", "service_websocket_register_function_namespace", "service_websocket_register_function_qualifier", "service_websocket_transport_function_namespace", "service_websocket_transport_function_qualifier", "service_websocket_cleanup_function_namespace", "service_websocket_cleanup_function_qualifier", "service_scf_is_integrated_response", "is_debug_after_charge", "is_delete_response_error_codes", "response_type", "response_success_example", "response_fail_example", "service_config", "auth_relation_api_id", "service_parameters", "oauth_config", "response_error_codes", "target_namespace_id", "user_type", "is_base64_encoded", "event_bus_id", "service_scf_function_type", "e_i_a_m_app_type", "e_i_a_m_auth_type", "token_timeout", "e_i_a_m_app_id", "owner"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("service_id") {
		if v, ok := d.GetOk("service_id"); ok {
			request.ServiceId = helper.String(v.(string))
		}
	}

	if d.HasChange("service_type") {
		if v, ok := d.GetOk("service_type"); ok {
			request.ServiceType = helper.String(v.(string))
		}
	}

	if d.HasChange("service_timeout") {
		if v, ok := d.GetOkExists("service_timeout"); ok {
			request.ServiceTimeout = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("protocol") {
		if v, ok := d.GetOk("protocol"); ok {
			request.Protocol = helper.String(v.(string))
		}
	}

	if d.HasChange("request_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "request_config"); ok {
			apiRequestConfig := apigateway.ApiRequestConfig{}
			if v, ok := dMap["path"]; ok {
				apiRequestConfig.Path = helper.String(v.(string))
			}
			if v, ok := dMap["method"]; ok {
				apiRequestConfig.Method = helper.String(v.(string))
			}
			request.RequestConfig = &apiRequestConfig
		}
	}

	if d.HasChange("api_name") {
		if v, ok := d.GetOk("api_name"); ok {
			request.ApiName = helper.String(v.(string))
		}
	}

	if d.HasChange("api_desc") {
		if v, ok := d.GetOk("api_desc"); ok {
			request.ApiDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("api_type") {
		if v, ok := d.GetOk("api_type"); ok {
			request.ApiType = helper.String(v.(string))
		}
	}

	if d.HasChange("auth_type") {
		if v, ok := d.GetOk("auth_type"); ok {
			request.AuthType = helper.String(v.(string))
		}
	}

	if d.HasChange("enable_c_o_r_s") {
		if v, ok := d.GetOkExists("enable_c_o_r_s"); ok {
			request.EnableCORS = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("constant_parameters") {
		if v, ok := d.GetOk("constant_parameters"); ok {
			for _, item := range v.([]interface{}) {
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
	}

	if d.HasChange("request_parameters") {
		if v, ok := d.GetOk("request_parameters"); ok {
			for _, item := range v.([]interface{}) {
				requestParameter := apigateway.RequestParameter{}
				if v, ok := dMap["name"]; ok {
					requestParameter.Name = helper.String(v.(string))
				}
				if v, ok := dMap["desc"]; ok {
					requestParameter.Desc = helper.String(v.(string))
				}
				if v, ok := dMap["position"]; ok {
					requestParameter.Position = helper.String(v.(string))
				}
				if v, ok := dMap["type"]; ok {
					requestParameter.Type = helper.String(v.(string))
				}
				if v, ok := dMap["default_value"]; ok {
					requestParameter.DefaultValue = helper.String(v.(string))
				}
				if v, ok := dMap["required"]; ok {
					requestParameter.Required = helper.Bool(v.(bool))
				}
				request.RequestParameters = append(request.RequestParameters, &requestParameter)
			}
		}
	}

	if d.HasChange("api_business_type") {
		if v, ok := d.GetOk("api_business_type"); ok {
			request.ApiBusinessType = helper.String(v.(string))
		}
	}

	if d.HasChange("service_mock_return_message") {
		if v, ok := d.GetOk("service_mock_return_message"); ok {
			request.ServiceMockReturnMessage = helper.String(v.(string))
		}
	}

	if d.HasChange("micro_services") {
		if v, ok := d.GetOk("micro_services"); ok {
			for _, item := range v.([]interface{}) {
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
	}

	if d.HasChange("service_tsf_load_balance_conf") {
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
	}

	if d.HasChange("service_tsf_health_check_conf") {
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
	}

	if d.HasChange("target_services_load_balance_conf") {
		if v, ok := d.GetOkExists("target_services_load_balance_conf"); ok {
			request.TargetServicesLoadBalanceConf = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("target_services_health_check_conf") {
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
	}

	if d.HasChange("service_scf_function_name") {
		if v, ok := d.GetOk("service_scf_function_name"); ok {
			request.ServiceScfFunctionName = helper.String(v.(string))
		}
	}

	if d.HasChange("service_websocket_register_function_name") {
		if v, ok := d.GetOk("service_websocket_register_function_name"); ok {
			request.ServiceWebsocketRegisterFunctionName = helper.String(v.(string))
		}
	}

	if d.HasChange("service_websocket_cleanup_function_name") {
		if v, ok := d.GetOk("service_websocket_cleanup_function_name"); ok {
			request.ServiceWebsocketCleanupFunctionName = helper.String(v.(string))
		}
	}

	if d.HasChange("service_websocket_transport_function_name") {
		if v, ok := d.GetOk("service_websocket_transport_function_name"); ok {
			request.ServiceWebsocketTransportFunctionName = helper.String(v.(string))
		}
	}

	if d.HasChange("service_scf_function_namespace") {
		if v, ok := d.GetOk("service_scf_function_namespace"); ok {
			request.ServiceScfFunctionNamespace = helper.String(v.(string))
		}
	}

	if d.HasChange("service_scf_function_qualifier") {
		if v, ok := d.GetOk("service_scf_function_qualifier"); ok {
			request.ServiceScfFunctionQualifier = helper.String(v.(string))
		}
	}

	if d.HasChange("service_websocket_register_function_namespace") {
		if v, ok := d.GetOk("service_websocket_register_function_namespace"); ok {
			request.ServiceWebsocketRegisterFunctionNamespace = helper.String(v.(string))
		}
	}

	if d.HasChange("service_websocket_register_function_qualifier") {
		if v, ok := d.GetOk("service_websocket_register_function_qualifier"); ok {
			request.ServiceWebsocketRegisterFunctionQualifier = helper.String(v.(string))
		}
	}

	if d.HasChange("service_websocket_transport_function_namespace") {
		if v, ok := d.GetOk("service_websocket_transport_function_namespace"); ok {
			request.ServiceWebsocketTransportFunctionNamespace = helper.String(v.(string))
		}
	}

	if d.HasChange("service_websocket_transport_function_qualifier") {
		if v, ok := d.GetOk("service_websocket_transport_function_qualifier"); ok {
			request.ServiceWebsocketTransportFunctionQualifier = helper.String(v.(string))
		}
	}

	if d.HasChange("service_websocket_cleanup_function_namespace") {
		if v, ok := d.GetOk("service_websocket_cleanup_function_namespace"); ok {
			request.ServiceWebsocketCleanupFunctionNamespace = helper.String(v.(string))
		}
	}

	if d.HasChange("service_websocket_cleanup_function_qualifier") {
		if v, ok := d.GetOk("service_websocket_cleanup_function_qualifier"); ok {
			request.ServiceWebsocketCleanupFunctionQualifier = helper.String(v.(string))
		}
	}

	if d.HasChange("service_scf_is_integrated_response") {
		if v, ok := d.GetOkExists("service_scf_is_integrated_response"); ok {
			request.ServiceScfIsIntegratedResponse = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("is_debug_after_charge") {
		if v, ok := d.GetOkExists("is_debug_after_charge"); ok {
			request.IsDebugAfterCharge = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("is_delete_response_error_codes") {
		if v, ok := d.GetOkExists("is_delete_response_error_codes"); ok {
			request.IsDeleteResponseErrorCodes = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("response_type") {
		if v, ok := d.GetOk("response_type"); ok {
			request.ResponseType = helper.String(v.(string))
		}
	}

	if d.HasChange("response_success_example") {
		if v, ok := d.GetOk("response_success_example"); ok {
			request.ResponseSuccessExample = helper.String(v.(string))
		}
	}

	if d.HasChange("response_fail_example") {
		if v, ok := d.GetOk("response_fail_example"); ok {
			request.ResponseFailExample = helper.String(v.(string))
		}
	}

	if d.HasChange("service_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "service_config"); ok {
			serviceConfig := apigateway.ServiceConfig{}
			if v, ok := dMap["product"]; ok {
				serviceConfig.Product = helper.String(v.(string))
			}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				serviceConfig.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["url"]; ok {
				serviceConfig.Url = helper.String(v.(string))
			}
			if v, ok := dMap["path"]; ok {
				serviceConfig.Path = helper.String(v.(string))
			}
			if v, ok := dMap["method"]; ok {
				serviceConfig.Method = helper.String(v.(string))
			}
			if v, ok := dMap["upstream_id"]; ok {
				serviceConfig.UpstreamId = helper.String(v.(string))
			}
			if cosConfigMap, ok := helper.InterfaceToMap(dMap, "cos_config"); ok {
				cosConfig := apigateway.CosConfig{}
				if v, ok := cosConfigMap["action"]; ok {
					cosConfig.Action = helper.String(v.(string))
				}
				if v, ok := cosConfigMap["bucket_name"]; ok {
					cosConfig.BucketName = helper.String(v.(string))
				}
				if v, ok := cosConfigMap["authorization"]; ok {
					cosConfig.Authorization = helper.Bool(v.(bool))
				}
				if v, ok := cosConfigMap["path_match_mode"]; ok {
					cosConfig.PathMatchMode = helper.String(v.(string))
				}
				serviceConfig.CosConfig = &cosConfig
			}
			request.ServiceConfig = &serviceConfig
		}
	}

	if d.HasChange("auth_relation_api_id") {
		if v, ok := d.GetOk("auth_relation_api_id"); ok {
			request.AuthRelationApiId = helper.String(v.(string))
		}
	}

	if d.HasChange("service_parameters") {
		if v, ok := d.GetOk("service_parameters"); ok {
			for _, item := range v.([]interface{}) {
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
	}

	if d.HasChange("oauth_config") {
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
	}

	if d.HasChange("response_error_codes") {
		if v, ok := d.GetOk("response_error_codes"); ok {
			for _, item := range v.([]interface{}) {
				responseErrorCodeReq := apigateway.ResponseErrorCodeReq{}
				if v, ok := dMap["code"]; ok {
					responseErrorCodeReq.Code = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["msg"]; ok {
					responseErrorCodeReq.Msg = helper.String(v.(string))
				}
				if v, ok := dMap["desc"]; ok {
					responseErrorCodeReq.Desc = helper.String(v.(string))
				}
				if v, ok := dMap["converted_code"]; ok {
					responseErrorCodeReq.ConvertedCode = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["need_convert"]; ok {
					responseErrorCodeReq.NeedConvert = helper.Bool(v.(bool))
				}
				request.ResponseErrorCodes = append(request.ResponseErrorCodes, &responseErrorCodeReq)
			}
		}
	}

	if d.HasChange("is_base64_encoded") {
		if v, ok := d.GetOkExists("is_base64_encoded"); ok {
			request.IsBase64Encoded = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("event_bus_id") {
		if v, ok := d.GetOk("event_bus_id"); ok {
			request.EventBusId = helper.String(v.(string))
		}
	}

	if d.HasChange("service_scf_function_type") {
		if v, ok := d.GetOk("service_scf_function_type"); ok {
			request.ServiceScfFunctionType = helper.String(v.(string))
		}
	}

	if d.HasChange("e_i_a_m_app_type") {
		if v, ok := d.GetOk("e_i_a_m_app_type"); ok {
			request.EIAMAppType = helper.String(v.(string))
		}
	}

	if d.HasChange("e_i_a_m_auth_type") {
		if v, ok := d.GetOk("e_i_a_m_auth_type"); ok {
			request.EIAMAuthType = helper.String(v.(string))
		}
	}

	if d.HasChange("token_timeout") {
		if v, ok := d.GetOkExists("token_timeout"); ok {
			request.TokenTimeout = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("e_i_a_m_app_id") {
		if v, ok := d.GetOk("e_i_a_m_app_id"); ok {
			request.EIAMAppId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().ModifyApiIncrement(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update apigateway API failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("apigw", "apiId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayAPIRead(d, meta)
}

func resourceTencentCloudApigatewayAPIDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_a_p_i.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	aPIId := d.Id()

	if err := service.DeleteApigatewayAPIById(ctx, apiId); err != nil {
		return err
	}

	return nil
}

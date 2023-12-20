package apigateway

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func DataSourceTencentCloudAPIGatewayAPIs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayAPIsRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service ID for query.",
			},
			"api_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom API name.",
			},
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Created API ID.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			// Computed values.
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of APIs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Which service this API belongs. Refer to resource `tencentcloud_api_gateway_service`.",
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
						"auth_type": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "API authentication type. Valid values: `SECRET`, `NONE`. " +
								"`SECRET` means key pair authentication, `NONE` means no authentication.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API frontend request type, such as `HTTP`,`WEBSOCKET`.",
						},
						"enable_cors": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable CORS.",
						},
						"request_config_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request frontend path configuration. Like `/user/getinfo`.",
						},
						"request_config_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request frontend method configuration. Like `GET`,`POST`,`PUT`,`DELETE`,`HEAD`,`ANY`.",
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
										Description: "If this parameter required.",
									},
								},
							},
						},
						"service_config_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API backend service type.",
						},
						"service_config_timeout": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "API backend service timeout period in seconds.",
						},
						"service_config_product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backend type. This parameter takes effect when VPC is enabled. Currently, only `clb` is supported.",
						},
						"service_config_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique VPC ID.",
						},
						"service_config_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API backend service url. This parameter is required when `service_config_type` is `HTTP`.",
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
						"service_config_scf_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SCF function name. This parameter takes effect when `service_config_type` is `SCF`.",
						},
						"service_config_scf_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SCF function namespace. This parameter takes effect when  `service_config_type` is `SCF`.",
						},
						"service_config_scf_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SCF function version. This parameter takes effect when `service_config_type`  is `SCF`.",
						},
						"service_config_mock_return_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Returned information of API backend mocking.",
						},
						"response_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Return type.",
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
						"modify_time": {
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
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayAPIsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_apis.read")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		apiName           = d.Get("api_name").(string)
		apiId             = d.Get("api_id").(string)
		serviceId         = d.Get("service_id").(string)
		apiSet            []*apigateway.DescribeApisStatusResultApiIdStatusSetInfo
		err               error
	)

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		apiSet, err = apiGatewayService.DescribeApisStatus(ctx, serviceId, apiName, apiId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(apiSet))
	for _, apiKey := range apiSet {
		var (
			info apigateway.ApiInfo
			has  bool
			item = make(map[string]interface{})
		)
		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			info, has, err = apiGatewayService.DescribeApi(ctx, *apiKey.ServiceId, *apiKey.ApiId)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
		if !has {
			continue
		}

		item["service_id"] = info.ServiceId
		item["api_name"] = info.ApiName
		item["api_desc"] = info.ApiDesc
		item["auth_type"] = info.AuthType
		item["protocol"] = info.Protocol
		item["enable_cors"] = info.EnableCORS
		item["response_type"] = info.ResponseType
		item["response_success_example"] = info.ResponseSuccessExample
		item["response_fail_example"] = info.ResponseFailExample
		item["service_config_type"] = info.ServiceType
		item["service_config_timeout"] = info.ServiceTimeout
		item["service_config_scf_function_name"] = info.ServiceScfFunctionName
		item["service_config_scf_function_namespace"] = info.ServiceScfFunctionNamespace
		item["service_config_scf_function_qualifier"] = info.ServiceScfFunctionQualifier
		item["service_config_mock_return_message"] = info.ServiceMockReturnMessage
		item["modify_time"] = info.ModifiedTime
		item["create_time"] = info.CreatedTime

		if info.RequestConfig != nil {
			item["request_config_path"] = info.RequestConfig.Path
			item["request_config_method"] = info.RequestConfig.Method
		} else {
			item["request_config_path"] = ""
			item["request_config_method"] = ""
		}

		paramList := make([]map[string]interface{}, 0, len(info.RequestParameters))
		if info.RequestParameters != nil {
			for _, param := range info.RequestParameters {
				paramList = append(paramList, map[string]interface{}{
					"name":          param.Name,
					"position":      param.Position,
					"type":          param.Type,
					"desc":          param.Desc,
					"default_value": param.DefaultValue,
					"required":      param.Required,
				})
			}
		}
		item["request_parameters"] = paramList

		if info.ServiceConfig != nil {
			item["service_config_product"] = info.ServiceConfig.Product
			item["service_config_vpc_id"] = info.ServiceConfig.UniqVpcId
			item["service_config_url"] = info.ServiceConfig.Url
			item["service_config_path"] = info.ServiceConfig.Path
			item["service_config_method"] = info.ServiceConfig.Method
		} else {
			item["service_config_product"] = ""
			item["service_config_vpc_id"] = ""
			item["service_config_url"] = ""
			item["service_config_path"] = ""
			item["service_config_method"] = ""
		}

		codeList := make([]map[string]interface{}, 0, len(info.ResponseErrorCodes))
		if info.ResponseErrorCodes != nil {
			for _, code := range info.ResponseErrorCodes {
				codeList = append(codeList, map[string]interface{}{
					"code":           code.Code,
					"msg":            code.Msg,
					"desc":           code.Desc,
					"converted_code": code.ConvertedCode,
					"need_convert":   code.NeedConvert,
				})
			}
		}
		item["response_error_codes"] = codeList

		list = append(list, item)
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{apiName, apiId}, tccommon.FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}

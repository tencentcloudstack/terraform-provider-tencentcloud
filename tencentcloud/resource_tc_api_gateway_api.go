/*
Use this resource to create API of API gateway.

Example Usage

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
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Description: "Which service this API belongs. Refer to resource `tencentcloud_api_gateway_service`.",
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
			"auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      API_GATEWAY_AUTH_TYPE_NONE,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_AUTH_TYPES),
				Description:  "API authentication type. Valid values: `SECRET` (key pair authentication),`NONE` (no authentication). Default value: `NONE`.",
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
			"service_config_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      API_GATEWAY_SERVICE_TYPE_HTTP,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_SERVICE_TYPES),
				Description:  "API backend service type. Valid values: `WEBSOCKET`, `HTTP`, `SCF`, `MOCK`. Default value: `HTTP`.",
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
				Description: "Backend type. This parameter takes effect when VPC is enabled. Currently, only `clb` is supported.",
			},
			"service_config_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique VPC ID.",
			},
			"service_config_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API backend service url. This parameter is required when `service_config_type` is `HTTP`.",
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
			"service_config_mock_return_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returned information of API backend mocking. This parameter is required when `service_config_type` is `MOCK`.",
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
				Description: "Successful response sample of custom response configuration.",
			},
			"response_fail_example": {
				Type:        schema.TypeString,
				Optional:    true,
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

	request.AuthType = helper.String(d.Get("auth_type").(string))
	request.Protocol = helper.String(d.Get("protocol").(string))
	request.EnableCORS = helper.Bool(d.Get("enable_cors").(bool))
	request.RequestConfig =
		&apigateway.ApiRequestConfig{Path: helper.String(d.Get("request_config_path").(string)),
			Method: helper.String(d.Get("request_config_method").(string))}

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
		if scfFunctionName == "" || scfFunctionNamespace == "" || scfFunctionQualifier == "" {
			return fmt.Errorf("`service_config_scf_function_name`,`service_config_scf_function_namespace`,`service_config_scf_function_qualifier` is needed if `service_config_type` is `SCF`")
		}
		request.ServiceScfFunctionName = &scfFunctionName
		request.ServiceScfFunctionNamespace = &scfFunctionNamespace
		request.ServiceScfFunctionQualifier = &scfFunctionQualifier
	}

	request.ResponseType = helper.String(d.Get("response_type").(string))

	if object, ok := d.GetOk("response_success_example"); ok {
		request.ResponseSuccessExample = helper.String(object.(string))
	}

	if object, ok := d.GetOk("response_fail_example"); ok {
		request.ResponseFailExample = helper.String(object.(string))
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
	_ = d.Set("auth_type", info.AuthType)
	_ = d.Set("protocol", info.Protocol)
	_ = d.Set("enable_cors", info.EnableCORS)
	_ = d.Set("response_type", info.ResponseType)
	_ = d.Set("response_success_example", info.ResponseSuccessExample)
	_ = d.Set("response_fail_example", info.ResponseFailExample)
	_ = d.Set("service_config_type", info.ServiceType)
	_ = d.Set("service_config_timeout", info.ServiceTimeout)
	_ = d.Set("service_config_scf_function_name", info.ServiceScfFunctionName)
	_ = d.Set("service_config_scf_function_namespace", info.ServiceScfFunctionNamespace)
	_ = d.Set("service_config_scf_function_qualifier", info.ServiceScfFunctionQualifier)
	_ = d.Set("service_config_mock_return_message", info.ServiceMockReturnMessage)
	_ = d.Set("update_time", info.ModifiedTime)
	_ = d.Set("create_time", info.CreatedTime)

	if info.RequestConfig != nil {
		_ = d.Set("request_config_path", info.RequestConfig.Path)
		_ = d.Set("request_config_method", info.RequestConfig.Method)
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

	if info.ServiceConfig != nil {
		_ = d.Set("service_config_product", info.ServiceConfig.Product)
		_ = d.Set("service_config_vpc_id", info.ServiceConfig.UniqVpcId)
		_ = d.Set("service_config_url", info.ServiceConfig.Url)
		_ = d.Set("service_config_path", info.ServiceConfig.Path)
		_ = d.Set("service_config_method", info.ServiceConfig.Method)
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

		releaseLimit int
		preLimit     int
		testLimit    int
	)
	d.Partial(true)
	request.ServiceId = &serviceId
	request.ApiId = &apiId
	request.ApiName = helper.String(d.Get("api_name").(string))
	if object, ok := d.GetOk("api_desc"); ok {
		request.ApiDesc = helper.String(object.(string))
	}

	request.AuthType = helper.String(d.Get("auth_type").(string))
	request.Protocol = helper.String(d.Get("protocol").(string))
	request.EnableCORS = helper.Bool(d.Get("enable_cors").(bool))
	request.RequestConfig = &apigateway.RequestConfig{
		Path:   helper.String(d.Get("request_config_path").(string)),
		Method: helper.String(d.Get("request_config_method").(string)),
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
		if scfFunctionName == "" || scfFunctionNamespace == "" || scfFunctionQualifier == "" {
			return fmt.Errorf("`service_config_scf_function_name`,`service_config_scf_function_namespace`,`service_config_scf_function_qualifier` is needed if `service_config_type` is `SCF`")
		}
		request.ServiceScfFunctionName = &scfFunctionName
		request.ServiceScfFunctionNamespace = &scfFunctionNamespace
		request.ServiceScfFunctionQualifier = &scfFunctionQualifier
	}

	request.ResponseType = helper.String(d.Get("response_type").(string))

	if object, ok := d.GetOk("response_success_example"); ok {
		request.ResponseSuccessExample = helper.String(object.(string))
	}

	if object, ok := d.GetOk("response_fail_example"); ok {
		request.ResponseFailExample = helper.String(object.(string))
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

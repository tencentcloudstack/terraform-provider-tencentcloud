/*
Use this resource to create API gateway throttling API.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
    service_id            = tencentcloud_api_gateway_service.service.id
    api_name              = "hello_update"
    api_desc              = "my hello api update"
    auth_type             = "SECRET"
    protocol              = "HTTP"
    enable_cors           = true
    request_config_path   = "/user/info"
    request_config_method = "POST"
    request_parameters {
    	name          = "email"
        position      = "QUERY"
        type          = "string"
        desc          = "your email please?"
        default_value = "tom@qq.com"
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
    	code           = 10
        msg            = "system error"
       	desc           = "system error code"
       	converted_code = -10
        need_convert   = true
    }
}

resource "tencentcloud_api_gateway_throttling_api" "foo" {
	service_id       = tencentcloud_api_gateway_service.service.id
	strategy         = "400"
	environment_name = "test"
	api_ids          = [tencentcloud_api_gateway_api.api.id]
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAPIGatewayThrottlingAPI() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayThrottlingAPICreate,
		Read:   resourceTencentCloudAPIGatewayThrottlingAPIRead,
		Update: resourceTencentCloudAPIGatewayThrottlingAPIUpdate,
		Delete: resourceTencentCloudAPIGatewayThrottlingAPIDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				ForceNew:     true,
				Description:  "Service ID for query.",
			},
			"strategy": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "API QPS value. Enter a positive number to limit the API query rate per second `QPS`.",
			},
			"environment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "List of Environment names.",
			},
			"api_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of API ID.",
			},
			//compute
			"api_environment_strategies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of throttling policies bound to API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique API ID.",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom API name.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API path.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API method.",
						},
						"strategy_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Environment throttling information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Environment name.",
									},
									"quota": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Throttling value.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAPIGatewayThrottlingAPICreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_api.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Get("service_id").(string)
		strategy          = d.Get("strategy").(int)
		environmentName   = d.Get("environment_name").(string)
		apiIds            = d.Get("api_ids").([]interface{})
		apiIdResults      []string
		err               error
	)

	for _, v := range apiIds {
		apiIdResults = append(apiIdResults, v.(string))
	}

	_, err = throttlingService.ModifyApiEnvironmentStrategy(ctx, serviceId, int64(strategy), environmentName, apiIdResults)
	if err != nil {
		return err
	}
	d.SetId(strings.Join([]string{serviceId, environmentName}, FILED_SP))

	return resourceTencentCloudAPIGatewayThrottlingAPIRead(d, meta)
}

func resourceTencentCloudAPIGatewayThrottlingAPIRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_api.read")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		err               error
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]
	environmentName := results[1]
	environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName})
	if err != nil {
		return err
	}
	if len(environmentList) == 0 {
		d.SetId("")
		return nil
	}

	environmentResults := make([]map[string]interface{}, 0, len(environmentList))
	for _, envList := range environmentList {
		environmentSet := envList.EnvironmentStrategySet
		strategy_list := make([]map[string]interface{}, 0, len(environmentSet))
		for _, envSet := range environmentSet {
			if envSet == nil {
				continue
			}
			strategy_list = append(strategy_list, map[string]interface{}{
				"environment_name": envSet.EnvironmentName,
				"quota":            envSet.Quota,
			})
		}

		item := map[string]interface{}{
			"api_id":        envList.ApiId,
			"api_name":      envList.ApiName,
			"path":          envList.Path,
			"method":        envList.Method,
			"strategy_list": strategy_list,
		}
		environmentResults = append(environmentResults, item)
	}

	_ = d.Set("service_id", serviceId)
	_ = d.Set("api_environment_strategies", environmentResults)

	return nil
}

func resourceTencentCloudAPIGatewayThrottlingAPIUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_api.update")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		err               error
		strategy          int64
		environmentName   string
		apiIds            []string
		hasChange         bool
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]

	strategy = int64(d.Get("strategy").(int))
	if d.HasChange("strategy") {
		hasChange = true
	}

	environmentName = d.Get("environment_name").(string)
	if d.HasChange("environment_name") {
		hasChange = true
	}

	apiIds = helper.InterfacesStrings(d.Get("api_ids").([]interface{}))
	if d.HasChange("api_ids") {
		hasChange = true
	}

	if hasChange {
		_, err = throttlingService.ModifyApiEnvironmentStrategy(ctx, serviceId, strategy, environmentName, apiIds)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayThrottlingAPIRead(d, meta)
}

func resourceTencentCloudAPIGatewayThrottlingAPIDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_api.delete")()

	var (
		logId                   = getLogId(contextNil)
		ctx                     = context.WithValue(context.TODO(), logIdKey, logId)
		id                      = d.Id()
		throttlingService       = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		strategy          int64 = QUOTA_MAX
		apiList           []string
		err               error
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]
	environmentName := results[1]

	environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName})
	if err != nil {
		return err
	}
	for _, envList := range environmentList {
		if envList == nil || envList.ApiId == nil {
			continue
		}
		apiList = append(apiList, *envList.ApiId)
	}

	if len(apiList) == 0 {
		return nil
	}

	_, err = throttlingService.ModifyApiEnvironmentStrategy(ctx, serviceId, strategy, environmentName, apiList)
	if err != nil {
		return err
	}

	return nil
}

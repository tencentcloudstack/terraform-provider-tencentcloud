/*
Use this resource to create API gateway throttling server.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_throttling_service" "service" {
	service_id        = tencentcloud_api_gateway_service.service.id
	strategy          = "400"
	environment_names = ["release"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAPIGatewayThrottlingService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayThrottlingServiceCreate,
		Read:   resourceTencentCloudAPIGatewayThrottlingServiceRead,
		Update: resourceTencentCloudAPIGatewayThrottlingServiceUpdate,
		Delete: resourceTencentCloudAPIGatewayThrottlingServiceDelete,

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
				Description: "Server QPS value. The service throttling value. Enter a positive number to limit the server query rate per second `QPS`.",
			},
			"environment_names": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of Environment names.",
			},
			//compute
			"environments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of Throttling policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Environment name.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access service environment URL.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Release status.",
						},
						"version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Published version number.",
						},
						"strategy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Throttling value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAPIGatewayThrottlingServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_service.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Get("service_id").(string)
		strategy          = d.Get("strategy").(int)
		environmentName   = d.Get("environment_names").([]interface{})
		nameResults       []string
		err               error
	)

	for _, v := range environmentName {
		nameResults = append(nameResults, v.(string))
	}

	_, err = throttlingService.ModifyServiceEnvironmentStrategy(ctx, serviceId, int64(strategy), nameResults)
	if err != nil {
		return err
	}
	d.SetId(serviceId)
	return resourceTencentCloudAPIGatewayThrottlingServiceRead(d, meta)
}

func resourceTencentCloudAPIGatewayThrottlingServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_service.read")()
	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Id()
		err               error
	)

	environmentList, err := throttlingService.DescribeServiceEnvironmentStrategyList(ctx, serviceId)
	if err != nil {
		return err
	}

	environmentResults := make([]map[string]interface{}, 0, len(environmentList))
	for _, value := range environmentList {
		if value == nil {
			continue
		}
		item := map[string]interface{}{
			"environment_name": value.EnvironmentName,
			"url":              value.Url,
			"status":           value.Status,
			"version_name":     value.VersionName,
			"strategy":         value.Strategy,
		}
		environmentResults = append(environmentResults, item)
	}
	_ = d.Set("service_id", serviceId)
	_ = d.Set("environments", environmentResults)
	return nil
}

func resourceTencentCloudAPIGatewayThrottlingServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_service.update")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Id()
		err               error
		strategy          int64
		environmentNames  []string
	)

	if d.HasChange("strategy") || d.HasChange("environment_names") {
		environmentNames = helper.InterfacesStrings(d.Get("environment_names").([]interface{}))
		strategy = int64(d.Get("strategy").(int))

		_, err = throttlingService.ModifyServiceEnvironmentStrategy(ctx, serviceId, strategy, environmentNames)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayThrottlingServiceRead(d, meta)
}

func resourceTencentCloudAPIGatewayThrottlingServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_service.delete")()

	var (
		logId                   = getLogId(contextNil)
		ctx                     = context.WithValue(context.TODO(), logIdKey, logId)
		serviceId               = d.Id()
		throttlingService       = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		strategy          int64 = STRATEGY_MAX
		environmentNames  []string
		err               error
	)

	environmentList, err := throttlingService.DescribeServiceEnvironmentStrategyList(ctx, serviceId)
	if err != nil {
		return err
	}
	for _, envList := range environmentList {
		if envList == nil || envList.EnvironmentName == nil {
			continue
		}
		environmentNames = append(environmentNames, *envList.EnvironmentName)
	}

	if len(environmentNames) == 0 {
		return nil
	}

	_, err = throttlingService.ModifyServiceEnvironmentStrategy(ctx, serviceId, strategy, environmentNames)
	if err != nil {
		return err
	}

	return nil
}

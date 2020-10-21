/*
Use this resource to create IP strategy of API gateway.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_ip_strategy" "test"{
    service_id    = tencentcloud_api_gateway_service.service.id
    strategy_name = "tf_test"
    strategy_type = "BLACK"
    strategy_data = "9.9.9.9"
}
```

Import

IP strategy of API gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_ip_strategy.test service-ohxqslqe#IPStrategy-q1lk8ud2
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func resourceTencentCloudAPIGatewayIPStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayIPStrategyCreate,
		Read:   resourceTencentCloudAPIGatewayIPStrategyRead,
		Update: resourceTencentCloudAPIGatewayIPStrategyUpdate,
		Delete: resourceTencentCloudAPIGatewayIPStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "The ID of the API gateway service.",
			},
			"strategy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "User defined strategy name.",
			},
			"strategy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "Blacklist or whitelist.",
			},
			"strategy_data": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIp,
				Description:  "IP address data.",
			},
			// Computed values.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
			},
			"strategy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP policy ID.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayIPStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_ip_strategy.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Get("service_id").(string)
		strategyName      = d.Get("strategy_name").(string)
		strategyType      = d.Get("strategy_type").(string)
		strategyData      = d.Get("strategy_data").(string)
		strategyId        string
		err               error
	)
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		strategyId, err = apiGatewayService.CreateIPStrategy(ctx, serviceId, strategyName, strategyType, strategyData)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{serviceId, strategyId}, FILED_SP))

	//wait ip strategy create ok
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, err := apiGatewayService.DescribeIPStrategyStatus(ctx, serviceId, strategyId)
		if err != nil {
			return retryError(err, InternalError)
		}
		if has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("strategyID %s not found on server", strategyId))

	}); err != nil {
		return err
	}

	return resourceTencentCloudAPIGatewayIPStrategyRead(d, meta)
}

func resourceTencentCloudAPIGatewayIPStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_ip_strategy.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		IpStatus          *apigateway.IPStrategy
		err               error
		has               bool
	)

	idSplit := strings.Split(id, FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("IP strategy is can't read, id is borken, id is %s", d.Id())
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		IpStatus, has, err = apiGatewayService.DescribeIPStrategyStatus(ctx, serviceId, strategyId)
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

	_ = d.Set("service_id", *IpStatus.ServiceId)
	_ = d.Set("strategy_name", *IpStatus.StrategyName)
	_ = d.Set("strategy_type", *IpStatus.StrategyType)
	_ = d.Set("strategy_data", *IpStatus.StrategyData)
	_ = d.Set("strategy_id", *IpStatus.StrategyId)
	_ = d.Set("create_time", *IpStatus.CreatedTime)

	return nil
}

func resourceTencentCloudAPIGatewayIPStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_ip_strategy.update")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
	)

	idSplit := strings.Split(id, FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("IP strategy is can't update, id is borken, id is %s", d.Id())
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]

	if d.HasChange("strategy_data") {
		var (
			strategyData = d.Get("strategy_data").(string)
			err          error
		)

		if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			err = apiGatewayService.UpdateIPStrategy(ctx, serviceId, strategyId, strategyData)

			if err != nil {
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayIPStrategyRead(d, meta)
}

func resourceTencentCloudAPIGatewayIPStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_ip_strategy.delete")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		err               error
	)

	idSplit := strings.Split(id, FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("IP strategy is can't delete, id is borken, id is %s", d.Id())
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err = apiGatewayService.DeleteIPStrategy(ctx, serviceId, strategyId)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

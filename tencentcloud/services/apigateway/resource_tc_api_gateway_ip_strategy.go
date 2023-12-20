package apigateway

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func ResourceTencentCloudAPIGatewayIPStrategy() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "The ID of the API gateway service.",
			},
			"strategy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "User defined strategy name.",
			},
			"strategy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "Blacklist or whitelist.",
			},
			"strategy_data": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateIp,
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
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_ip_strategy.create")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId         = d.Get("service_id").(string)
		strategyName      = d.Get("strategy_name").(string)
		strategyType      = d.Get("strategy_type").(string)
		strategyData      = d.Get("strategy_data").(string)
		strategyId        string
		err               error
	)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		strategyId, err = apiGatewayService.CreateIPStrategy(ctx, serviceId, strategyName, strategyType, strategyData)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{serviceId, strategyId}, tccommon.FILED_SP))

	//wait ip strategy create ok
	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, err := apiGatewayService.DescribeIPStrategyStatus(ctx, serviceId, strategyId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
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
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_ip_strategy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id                = d.Id()
		IpStatus          *apigateway.IPStrategy
		err               error
		has               bool
	)

	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("IP strategy is can't read, id is borken, id is %s", d.Id())
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		IpStatus, has, err = apiGatewayService.DescribeIPStrategyStatus(ctx, serviceId, strategyId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("service_id", IpStatus.ServiceId)
	_ = d.Set("strategy_name", IpStatus.StrategyName)
	_ = d.Set("strategy_type", IpStatus.StrategyType)
	_ = d.Set("strategy_data", IpStatus.StrategyData)
	_ = d.Set("strategy_id", IpStatus.StrategyId)
	_ = d.Set("create_time", IpStatus.CreatedTime)

	return nil
}

func resourceTencentCloudAPIGatewayIPStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_ip_strategy.update")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id                = d.Id()
	)

	idSplit := strings.Split(id, tccommon.FILED_SP)
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

		if err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err = apiGatewayService.UpdateIPStrategy(ctx, serviceId, strategyId, strategyData)

			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayIPStrategyRead(d, meta)
}

func resourceTencentCloudAPIGatewayIPStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_ip_strategy.delete")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id                = d.Id()
		err               error
	)

	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("IP strategy is can't delete, id is borken, id is %s", d.Id())
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		err = apiGatewayService.DeleteIPStrategy(ctx, serviceId, strategyId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

package apigateway

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudAPIGatewayStrategyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayStrategyAttachmentCreate,
		Read:   resourceTencentCloudAPIGatewayStrategyAttachmentRead,
		Delete: resourceTencentCloudAPIGatewayStrategyAttachmentDelete,
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
			"strategy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "The ID of the API gateway strategy.",
			},
			"environment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(API_GATEWAY_SERVICE_ENVS),
				Description:  "The environment of the strategy association. Valid values: `test`, `release`, `prepub`.",
			},
			"bind_api_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "The API that needs to be bound.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayStrategyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_strategy_attachment.create")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId         = d.Get("service_id").(string)
		strategyId        = d.Get("strategy_id").(string)
		envName           = d.Get("environment_name").(string)
		bindApiId         = d.Get("bind_api_id").(string)
		err               error
		has               bool
	)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err = apiGatewayService.CreateStrategyAttachment(ctx, serviceId, strategyId, envName, bindApiId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{serviceId, strategyId, bindApiId, envName}, tccommon.FILED_SP))

	//wait IP strategy create ok
	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		has, err = apiGatewayService.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		if has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("IP strategy attachment %s not found on server", strings.Join([]string{strategyId, bindApiId}, tccommon.FILED_SP)))
	}); err != nil {
		return err
	}

	return resourceTencentCloudAPIGatewayStrategyAttachmentRead(d, meta)
}

func resourceTencentCloudAPIGatewayStrategyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_strategy_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		attachmentId      = d.Id()
		err               error
		has               bool
	)

	idSplit := strings.Split(attachmentId, tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("IP strategy attachment id is broken, id is %s", attachmentId)
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]
	bindApiId := idSplit[2]
	envname := idSplit[3]

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		has, err = apiGatewayService.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
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

	_ = d.Set("service_id", serviceId)
	_ = d.Set("strategy_id", strategyId)
	_ = d.Set("bind_api_id", bindApiId)
	_ = d.Set("environment_name", envname)

	return nil
}

func resourceTencentCloudAPIGatewayStrategyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_strategy_attachment.delete")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId         = d.Get("service_id").(string)
		strategyId        = d.Get("strategy_id").(string)
		envName           = d.Get("environment_name").(string)
		bindApiId         = d.Get("bind_api_id").(string)
	)

	has, err := apiGatewayService.DeleteStrategyAttachment(ctx, serviceId, strategyId, envName, bindApiId)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("IP strategy is still exist")
	}

	return nil
}

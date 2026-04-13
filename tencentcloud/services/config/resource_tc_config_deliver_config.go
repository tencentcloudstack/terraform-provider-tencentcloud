package config

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudConfigDeliverConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigDeliverConfigCreate,
		Read:   resourceTencentCloudConfigDeliverConfigRead,
		Update: resourceTencentCloudConfigDeliverConfigUpdate,
		Delete: resourceTencentCloudConfigDeliverConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Delivery switch. Valid values: 0 (disabled), 1 (enabled).",
			},

			"deliver_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Delivery service name.",
			},

			"target_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource ARN. COS format: qcs::cos:$region:$account:prefix/$appid/$BucketName. CLS format: qcs::cls:$region:$account:cls/topicId.",
			},

			"deliver_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log prefix for stored delivery content.",
			},

			"deliver_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Delivery type. Valid values: COS, CLS.",
			},

			"deliver_content_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Delivery content type. Valid values: 1 (configuration change), 2 (resource list), 3 (all).",
			},

			// Computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the delivery configuration.",
			},
		},
	}
}

func resourceTencentCloudConfigDeliverConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_deliver_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudConfigDeliverConfigUpdate(d, meta)
}

func resourceTencentCloudConfigDeliverConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_deliver_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeConfigDeliver(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_config_deliver_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Status != nil {
		_ = d.Set("status", int(*respData.Status))
	}

	if respData.DeliverName != nil {
		_ = d.Set("deliver_name", respData.DeliverName)
	}

	if respData.TargetArn != nil {
		_ = d.Set("target_arn", respData.TargetArn)
	}

	if respData.DeliverPrefix != nil {
		_ = d.Set("deliver_prefix", respData.DeliverPrefix)
	}

	if respData.DeliverType != nil {
		_ = d.Set("deliver_type", respData.DeliverType)
	}

	if respData.DeliverContentType != nil {
		_ = d.Set("deliver_content_type", int(*respData.DeliverContentType))
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	return nil
}

func resourceTencentCloudConfigDeliverConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_deliver_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = configv20220802.NewUpdateConfigDeliverRequest()
	)

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("deliver_name"); ok {
		request.DeliverName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_arn"); ok {
		request.TargetArn = helper.String(v.(string))
	}

	if v, ok := d.GetOk("deliver_prefix"); ok {
		request.DeliverPrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("deliver_type"); ok {
		request.DeliverType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("deliver_content_type"); ok {
		request.DeliverContentType = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateConfigDeliverWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update config deliver config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudConfigDeliverConfigRead(d, meta)
}

func resourceTencentCloudConfigDeliverConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_deliver_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

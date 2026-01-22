package bh

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhReconnectionSettingConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhReconnectionSettingConfigCreate,
		Read:   resourceTencentCloudBhReconnectionSettingConfigRead,
		Update: resourceTencentCloudBhReconnectionSettingConfigUpdate,
		Delete: resourceTencentCloudBhReconnectionSettingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"reconnection_max_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Retry count, value range: 0-20.",
			},

			"enable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "true: limit reconnection count, false: do not limit reconnection count.",
			},
		},
	}
}

func resourceTencentCloudBhReconnectionSettingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_reconnection_setting_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudBhReconnectionSettingConfigUpdate(d, meta)
}

func resourceTencentCloudBhReconnectionSettingConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_reconnection_setting_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeBhReconnectionSettingConfigById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_reconnection_setting_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Reconnection != nil {
		if respData.Reconnection.ReconnectionMaxCount != nil {
			_ = d.Set("reconnection_max_count", respData.Reconnection.ReconnectionMaxCount)
		}

		if respData.Reconnection.Enable != nil {
			_ = d.Set("enable", respData.Reconnection.Enable)
		}
	}

	return nil
}

func resourceTencentCloudBhReconnectionSettingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_reconnection_setting_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewModifyReconnectionSettingRequest()
	)

	if v, ok := d.GetOkExists("reconnection_max_count"); ok {
		request.ReconnectionMaxCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("enable"); ok {
		request.Enable = helper.Bool(v.(bool))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyReconnectionSettingWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update bh reconnection setting config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudBhReconnectionSettingConfigRead(d, meta)
}

func resourceTencentCloudBhReconnectionSettingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_reconnection_setting_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

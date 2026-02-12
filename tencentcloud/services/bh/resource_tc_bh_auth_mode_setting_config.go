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

func ResourceTencentCloudBhAuthModeSettingConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhAuthModeSettingConfigCreate,
		Read:   resourceTencentCloudBhAuthModeSettingConfigRead,
		Update: resourceTencentCloudBhAuthModeSettingConfigUpdate,
		Delete: resourceTencentCloudBhAuthModeSettingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auth_mode": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Two-factor authentication, 0-disabled, 1-OTP, 2-SMS, 3-USB Key.",
			},

			"resource_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Resource type, 0: normal 1: national cryptography.",
			},
		},
	}
}

func resourceTencentCloudBhAuthModeSettingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_setting_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudBhAuthModeSettingConfigUpdate(d, meta)
}

func resourceTencentCloudBhAuthModeSettingConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_setting_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeBhAuthModeSettingConfigById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_auth_mode_setting_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AuthModeGM != nil {
		if respData.AuthModeGM.AuthMode != nil {
			_ = d.Set("auth_mode", respData.AuthModeGM.AuthMode)
		}
	}

	return nil
}

func resourceTencentCloudBhAuthModeSettingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_setting_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewModifyAuthModeSettingRequest()
	)

	if v, ok := d.GetOkExists("auth_mode"); ok {
		request.AuthMode = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("resource_type"); ok {
		request.ResourceType = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyAuthModeSettingWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update bh auth mode setting config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudBhAuthModeSettingConfigRead(d, meta)
}

func resourceTencentCloudBhAuthModeSettingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_setting_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

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

func ResourceTencentCloudBhAuthModeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhAuthModeConfigCreate,
		Read:   resourceTencentCloudBhAuthModeConfigRead,
		Update: resourceTencentCloudBhAuthModeConfigUpdate,
		Delete: resourceTencentCloudBhAuthModeConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auth_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Double factor authentication configuration. Valid values: `0` (disabled), `1` (OTP), `2` (SMS), `3` (USB Key, only valid when ResourceType=1 and AuthModeGM is not passed). Note: At least one of AuthMode and AuthModeGM must be passed.",
			},
			"auth_mode_gm": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "National secret double factor authentication configuration. Valid values: `0` (disabled), `1` (OTP), `2` (SMS), `3` (USB Key). Note: At least one of AuthMode and AuthModeGM must be passed. AuthModeGM has higher priority than ResourceType.",
			},
		},
	}
}

func resourceTencentCloudBhAuthModeConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudBhAuthModeConfigUpdate(d, meta)
}

func resourceTencentCloudBhAuthModeConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeBhAuthModeConfigById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_auth_mode_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	// Set auth_mode if returned
	if respData.AuthMode != nil {
		if respData.AuthMode.AuthMode != nil {
			_ = d.Set("auth_mode", respData.AuthMode.AuthMode)
		}
	}

	// Set auth_mode_gm if returned
	if respData.AuthModeGM != nil {
		if respData.AuthModeGM.AuthMode != nil {
			_ = d.Set("auth_mode_gm", respData.AuthModeGM.AuthMode)
		}
	}

	return nil
}

func resourceTencentCloudBhAuthModeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewModifyAuthModeSettingRequest()
	)

	// Detect changes and build request
	if d.HasChange("auth_mode") || d.HasChange("auth_mode_gm") {
		if v, ok := d.GetOkExists("auth_mode"); ok {
			request.AuthMode = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("auth_mode_gm"); ok {
			request.AuthModeGM = helper.IntUint64(v.(int))
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
			log.Printf("[CRITAL]%s update bh auth mode config failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBhAuthModeConfigRead(d, meta)
}

func resourceTencentCloudBhAuthModeConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

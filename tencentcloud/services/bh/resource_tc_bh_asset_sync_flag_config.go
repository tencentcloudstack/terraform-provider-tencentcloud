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

func ResourceTencentCloudBhAssetSyncFlagConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhAssetSyncFlagConfigCreate,
		Read:   resourceTencentCloudBhAssetSyncFlagConfigRead,
		Update: resourceTencentCloudBhAssetSyncFlagConfigUpdate,
		Delete: resourceTencentCloudBhAssetSyncFlagConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_sync": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable asset auto-sync, false - disabled, true - enabled.",
			},

			// computed
			"role_granted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the role has been authorized, false - not authorized, true - authorized.",
			},
		},
	}
}

func resourceTencentCloudBhAssetSyncFlagConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_asset_sync_flag_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudBhAssetSyncFlagConfigUpdate(d, meta)
}

func resourceTencentCloudBhAssetSyncFlagConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_asset_sync_flag_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeBhAssetSyncFlagConfigById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_asset_sync_flag_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AutoSync != nil {
		_ = d.Set("auto_sync", respData.AutoSync)
	}

	if respData.RoleGranted != nil {
		_ = d.Set("role_granted", respData.RoleGranted)
	}

	return nil
}

func resourceTencentCloudBhAssetSyncFlagConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_asset_sync_flag_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewModifyAssetSyncFlagRequest()
	)

	if v, ok := d.GetOkExists("auto_sync"); ok {
		request.AutoSync = helper.Bool(v.(bool))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyAssetSyncFlagWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update bh asset sync flag config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudBhAssetSyncFlagConfigRead(d, meta)
}

func resourceTencentCloudBhAssetSyncFlagConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_asset_sync_flag_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

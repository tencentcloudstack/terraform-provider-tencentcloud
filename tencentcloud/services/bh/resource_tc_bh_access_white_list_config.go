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

func ResourceTencentCloudBhAccessWhiteListConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhAccessWhiteListConfigCreate,
		Read:   resourceTencentCloudBhAccessWhiteListConfigRead,
		Update: resourceTencentCloudBhAccessWhiteListConfigUpdate,
		Delete: resourceTencentCloudBhAccessWhiteListConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"allow_any": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "true: allow all source IPs; false: do not allow all source IPs.",
			},

			"allow_auto": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "true: allow automatically added IPs; false: do not allow automatically added IPs.",
			},
		},
	}
}

func resourceTencentCloudBhAccessWhiteListConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_access_white_list_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudBhAccessWhiteListConfigUpdate(d, meta)
}

func resourceTencentCloudBhAccessWhiteListConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_access_white_list_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeBhAccessWhiteListConfigById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `bh_access_white_list_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AllowAny != nil {
		_ = d.Set("allow_any", respData.AllowAny)
	}

	if respData.AllowAuto != nil {
		_ = d.Set("allow_auto", respData.AllowAuto)
	}

	return nil
}

func resourceTencentCloudBhAccessWhiteListConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_access_white_list_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	if d.HasChange("allow_any") {
		request := bhv20230418.NewModifyAccessWhiteListStatusRequest()
		if v, ok := d.GetOkExists("allow_any"); ok {
			request.AllowAny = helper.Bool(v.(bool))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyAccessWhiteListStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bh access white list config failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("allow_auto") {
		request := bhv20230418.NewModifyAccessWhiteListAutoStatusRequest()
		if v, ok := d.GetOkExists("allow_auto"); ok {
			request.AllowAuto = helper.Bool(v.(bool))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyAccessWhiteListAutoStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bh access white list config failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBhAccessWhiteListConfigRead(d, meta)
}

func resourceTencentCloudBhAccessWhiteListConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_access_white_list_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

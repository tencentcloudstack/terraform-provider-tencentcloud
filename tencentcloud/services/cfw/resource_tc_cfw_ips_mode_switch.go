package cfw

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwIpsModeSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwIpsModeSwitchCreate,
		Read:   resourceTencentCloudCfwIpsModeSwitchRead,
		Update: resourceTencentCloudCfwIpsModeSwitchUpdate,
		Delete: resourceTencentCloudCfwIpsModeSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Protection mode: 0-observation mode, 1-interception mode, 2-strict mode.",
			},
		},
	}
}

func resourceTencentCloudCfwIpsModeSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_ips_mode_switch.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudCfwIpsModeSwitchUpdate(d, meta)
}

func resourceTencentCloudCfwIpsModeSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_ips_mode_switch.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeCfwIpsModeSwitchById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_ips_mode_switch` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Data != nil {
		if respData.Data.Mode != nil {
			_ = d.Set("mode", respData.Data.Mode)
		}
	}

	return nil
}

func resourceTencentCloudCfwIpsModeSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_ips_mode_switch.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cfw.NewModifyIpsModeSwitchRequest()
	)

	if v, ok := d.GetOkExists("mode"); ok {
		request.Mode = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyIpsModeSwitchWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update cfw ips mode switch failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudCfwIpsModeSwitchRead(d, meta)
}

func resourceTencentCloudCfwIpsModeSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_ips_mode_switch.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

package teo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoConfirmOriginAclUpdateOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoConfirmOriginAclUpdateOperationCreate,
		Read:   resourceTencentCloudTeoConfirmOriginAclUpdateOperationRead,
		Delete: resourceTencentCloudTeoConfirmOriginAclUpdateOperationDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
		},
	}
}

func resourceTencentCloudTeoConfirmOriginAclUpdateOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_confirm_origin_acl_update_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewConfirmOriginACLUpdateRequest()
	)

	zoneId := d.Get("zone_id").(string)
	request.ZoneId = helper.String(zoneId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ConfirmOriginACLUpdateWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s confirm origin acl update failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.BuildToken())
	return resourceTencentCloudTeoConfirmOriginAclUpdateOperationRead(d, meta)
}

func resourceTencentCloudTeoConfirmOriginAclUpdateOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_confirm_origin_acl_update_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoConfirmOriginAclUpdateOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_confirm_origin_acl_update_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

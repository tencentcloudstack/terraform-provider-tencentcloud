package vcube

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vcubev20220410 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vcube/v20220410"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVcubeRenewVideoOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVcubeRenewVideoOperationCreate,
		Read:   resourceTencentCloudVcubeRenewVideoOperationRead,
		Delete: resourceTencentCloudVcubeRenewVideoOperationDelete,
		Schema: map[string]*schema.Schema{
			"license_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "License ID for video playback renewal.",
			},
		},
	}
}

func resourceTencentCloudVcubeRenewVideoOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vcube_renew_video_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = vcubev20220410.NewRenewVideoRequest()
		licenseId string
	)

	if v, ok := d.GetOkExists("license_id"); ok {
		request.LicenseId = helper.IntUint64(v.(int))
		licenseId = helper.IntToStr(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVcubeV20220410Client().RenewVideoWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vcube renew video operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(licenseId)
	return resourceTencentCloudVcubeRenewVideoOperationRead(d, meta)
}

func resourceTencentCloudVcubeRenewVideoOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vcube_renew_video_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVcubeRenewVideoOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vcube_renew_video_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

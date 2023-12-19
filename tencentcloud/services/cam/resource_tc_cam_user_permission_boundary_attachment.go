package cam

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCamUserPermissionBoundaryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamUserPermissionBoundaryAttachmentCreate,
		Read:   resourceTencentCloudCamUserPermissionBoundaryAttachmentRead,
		Delete: resourceTencentCloudCamUserPermissionBoundaryAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_uin": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Sub account Uin.",
			},

			"policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Policy ID.",
			},
		},
	}
}

func resourceTencentCloudCamUserPermissionBoundaryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user_permission_boundary_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = cam.NewPutUserPermissionsBoundaryRequest()
		targetUin string
		policyId  string
	)
	if v, ok := d.GetOkExists("target_uin"); ok {
		targetUin = helper.IntToStr(v.(int))
		request.TargetUin = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = helper.IntToStr(v.(int))
		request.PolicyId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().PutUserPermissionsBoundary(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam UserPermissionBoundary failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(targetUin + tccommon.FILED_SP + policyId)

	return resourceTencentCloudCamUserPermissionBoundaryAttachmentRead(d, meta)
}

func resourceTencentCloudCamUserPermissionBoundaryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user_permission_boundary_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	targetUin := idSplit[0]

	UserPermissionBoundary, err := service.DescribeCamUserPermissionBoundaryById(ctx, targetUin)
	if err != nil {
		return err
	}

	if UserPermissionBoundary == nil || UserPermissionBoundary.Response == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamUserPermissionBoundary` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if UserPermissionBoundary.Response.PolicyId != nil {
		_ = d.Set("policy_id", UserPermissionBoundary.Response.PolicyId)
	}
	return nil
}

func resourceTencentCloudCamUserPermissionBoundaryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user_permission_boundary_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	targetUin := idSplit[0]

	if err := service.DeleteCamUserPermissionBoundaryById(ctx, targetUin); err != nil {
		return err
	}

	return nil
}

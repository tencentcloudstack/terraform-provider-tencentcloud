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

func ResourceTencentCloudCamRolePermissionBoundaryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamRolePermissionBoundaryAttachmentCreate,
		Read:   resourceTencentCloudCamRolePermissionBoundaryAttachmentRead,
		Delete: resourceTencentCloudCamRolePermissionBoundaryAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Role ID.",
			},

			"role_id": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Role ID (at least one should be filled in with the role name).",
			},

			"role_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Role name (at least one should be filled in with the role ID).",
			},
		},
	}
}

func resourceTencentCloudCamRolePermissionBoundaryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role_permission_boundary_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cam.NewPutRolePermissionsBoundaryRequest()
		policyId string
		roleId   string
		roleName string
	)
	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = helper.IntToStr(v.(int))
		request.PolicyId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("role_id"); ok {
		roleId = v.(string)
		request.RoleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role_name"); ok {
		roleName = v.(string)
		request.RoleName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().PutRolePermissionsBoundary(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam RolePermissionBoundaryAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(policyId + tccommon.FILED_SP + roleId + tccommon.FILED_SP + roleName)

	return resourceTencentCloudCamRolePermissionBoundaryAttachmentRead(d, meta)
}

func resourceTencentCloudCamRolePermissionBoundaryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role_permission_boundary_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	roleId := idSplit[1]
	roleName := idSplit[2]

	if roleId == "" {
		roleInfo, err := service.DescribeCamTagRoleById(ctx, roleName, roleId)
		if err != nil {
			return err
		}
		if roleInfo == nil {
			return fmt.Errorf("role info is null")
		}
		roleId = *roleInfo.RoleId
	}

	RolePermissionBoundaryAttachment, err := service.DescribeCamRolePermissionBoundaryAttachmentById(ctx, roleId, policyId)
	if err != nil {
		return err
	}

	if RolePermissionBoundaryAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamRolePermissionBoundaryAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if RolePermissionBoundaryAttachment.PolicyId != nil {
		_ = d.Set("policy_id", RolePermissionBoundaryAttachment.PolicyId)
	}

	_ = d.Set("role_id", roleId)
	_ = d.Set("role_name", roleName)

	return nil
}

func resourceTencentCloudCamRolePermissionBoundaryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role_permission_boundary_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	roleId := idSplit[1]
	roleName := idSplit[2]

	if err := service.DeleteCamRolePermissionBoundaryAttachmentById(ctx, roleId, roleName); err != nil {
		return err
	}

	return nil
}

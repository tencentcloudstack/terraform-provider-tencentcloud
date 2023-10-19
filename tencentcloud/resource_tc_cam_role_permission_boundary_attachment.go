/*
Provides a resource to create a cam role_permission_boundary_attachment

Example Usage

```hcl
resource "tencentcloud_cam_role_permission_boundary_attachment" "role_permission_boundary_attachment" {
  policy_id = 1
  role_name = "test-cam-tag"
}
```

Import

cam role_permission_boundary_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cam_role_permission_boundary_attachment.role_permission_boundary_attachment role_permission_boundary_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamRolePermissionBoundaryAttachment() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_cam_role_permission_boundary_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().PutRolePermissionsBoundary(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam RolePermissionBoundaryAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(policyId + FILED_SP + roleId + FILED_SP + roleName)

	return resourceTencentCloudCamRolePermissionBoundaryAttachmentRead(d, meta)
}

func resourceTencentCloudCamRolePermissionBoundaryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_permission_boundary_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
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
	defer logElapsed("resource.tencentcloud_cam_role_permission_boundary_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
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

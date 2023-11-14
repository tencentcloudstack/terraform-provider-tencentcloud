/*
Provides a resource to create a cam role_permission_boundary

Example Usage

```hcl
resource "tencentcloud_cam_role_permission_boundary" "role_permission_boundary" {
  policy_id =
  role_id = ""
  role_name = ""
}
```

Import

cam role_permission_boundary can be imported using the id, e.g.

```
terraform import tencentcloud_cam_role_permission_boundary.role_permission_boundary role_permission_boundary_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCamRolePermissionBoundary() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamRolePermissionBoundaryCreate,
		Read:   resourceTencentCloudCamRolePermissionBoundaryRead,
		Delete: resourceTencentCloudCamRolePermissionBoundaryDelete,
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

func resourceTencentCloudCamRolePermissionBoundaryCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_permission_boundary.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cam.NewPutRolePermissionsBoundaryRequest()
		response = cam.NewPutRolePermissionsBoundaryResponse()
		roleId   string
		policyId int
	)
	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = v.(int64)
		request.PolicyId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("role_id"); ok {
		roleId = v.(string)
		request.RoleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role_name"); ok {
		request.RoleName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().PutRolePermissionsBoundary(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam RolePermissionBoundary failed, reason:%+v", logId, err)
		return err
	}

	roleId = *response.Response.RoleId
	d.SetId(strings.Join([]string{roleId, helper.Int64ToStr(policyId)}, FILED_SP))

	return resourceTencentCloudCamRolePermissionBoundaryRead(d, meta)
}

func resourceTencentCloudCamRolePermissionBoundaryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_permission_boundary.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	roleId := idSplit[0]
	policyId := idSplit[1]

	RolePermissionBoundary, err := service.DescribeCamRolePermissionBoundaryById(ctx, roleId, policyId)
	if err != nil {
		return err
	}

	if RolePermissionBoundary == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamRolePermissionBoundary` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if RolePermissionBoundary.PolicyId != nil {
		_ = d.Set("policy_id", RolePermissionBoundary.PolicyId)
	}

	if RolePermissionBoundary.RoleId != nil {
		_ = d.Set("role_id", RolePermissionBoundary.RoleId)
	}

	if RolePermissionBoundary.RoleName != nil {
		_ = d.Set("role_name", RolePermissionBoundary.RoleName)
	}

	return nil
}

func resourceTencentCloudCamRolePermissionBoundaryDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_permission_boundary.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	roleId := idSplit[0]
	policyId := idSplit[1]

	if err := service.DeleteCamRolePermissionBoundaryById(ctx, roleId, policyId); err != nil {
		return err
	}

	return nil
}

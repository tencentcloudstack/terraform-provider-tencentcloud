/*
Provides a resource to create a cam tag_role

Example Usage

```hcl
resource "tencentcloud_cam_tag_role" "tag_role" {
  role_name = ""
  role_id = ""
}
```

Import

cam tag_role can be imported using the id, e.g.

```
terraform import tencentcloud_cam_tag_role.tag_role tag_role_id
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
)

func resourceTencentCloudCamTagRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamTagRoleCreate,
		Read:   resourceTencentCloudCamTagRoleRead,
		Delete: resourceTencentCloudCamTagRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"role_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Enter at least one role name and role ID.",
			},

			"role_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Enter at least one role ID and role name.",
			},
		},
	}
}

func resourceTencentCloudCamTagRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_tag_role.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cam.NewTagRoleRequest()
		response = cam.NewTagRoleResponse()
		roleId   string
	)
	if v, ok := d.GetOk("role_name"); ok {
		request.RoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role_id"); ok {
		roleId = v.(string)
		request.RoleId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().TagRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam TagRole failed, reason:%+v", logId, err)
		return err
	}

	roleId = *response.Response.RoleId
	d.SetId(roleId)

	return resourceTencentCloudCamTagRoleRead(d, meta)
}

func resourceTencentCloudCamTagRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_tag_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	tagRoleId := d.Id()

	TagRole, err := service.DescribeCamTagRoleById(ctx, roleId)
	if err != nil {
		return err
	}

	if TagRole == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamTagRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if TagRole.RoleName != nil {
		_ = d.Set("role_name", TagRole.RoleName)
	}

	if TagRole.RoleId != nil {
		_ = d.Set("role_id", TagRole.RoleId)
	}

	return nil
}

func resourceTencentCloudCamTagRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_tag_role.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	tagRoleId := d.Id()

	if err := service.DeleteCamTagRoleById(ctx, roleId); err != nil {
		return err
	}

	return nil
}

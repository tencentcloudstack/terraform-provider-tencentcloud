/*
Provides a resource to create a cam role

Example Usage

```hcl
resource "tencentcloud_cam_role" "role" {
  role_name = ""
  policy_document = ""
  description = ""
  console_login =
  session_duration =
}
```

Import

cam role can be imported using the id, e.g.

```
terraform import tencentcloud_cam_role.role role_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamRoleCreate,
		Read:   resourceTencentCloudCamRoleRead,
		Update: resourceTencentCloudCamRoleUpdate,
		Delete: resourceTencentCloudCamRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"role_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Role name. The length is 1-128 characters and can contain English letters, numbers, and+=,. @ - _.",
			},

			"policy_document": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Policy document, example:{version:2.0,statement:[{action:name/sts:AssumeRole,effect:allow,principal:{service:[cloudaudit.cloud.tencent.com,cls.cloud.tencent.com]}}]}，Principal is used to specify the authorization object for a role. To obtain this parameter, please refer to Obtaining Role Details（https://cloud.tencent.com/document/product/598/36221） Output parameter RoleInfo.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Role Description.",
			},

			"console_login": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Allow login if 1 is allowed and 0 is not allowed.",
			},

			"session_duration": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum validity limit for applying for role temporary keys (range: 0-43200).",
			},
		},
	}
}

func resourceTencentCloudCamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cam.NewCreateRoleRequest()
		response = cam.NewCreateRoleResponse()
		roleId   string
	)
	if v, ok := d.GetOk("role_name"); ok {
		request.RoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_document"); ok {
		request.PolicyDocument = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("console_login"); ok {
		request.ConsoleLogin = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("session_duration"); ok {
		request.SessionDuration = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam Role failed, reason:%+v", logId, err)
		return err
	}

	roleId = *response.Response.RoleId
	d.SetId(roleId)

	return resourceTencentCloudCamRoleRead(d, meta)
}

func resourceTencentCloudCamRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	roleId := d.Id()

	Role, err := service.DescribeCamRoleById(ctx, roleId)
	if err != nil {
		return err
	}

	if Role == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if Role.RoleName != nil {
		_ = d.Set("role_name", Role.RoleName)
	}

	if Role.PolicyDocument != nil {
		_ = d.Set("policy_document", Role.PolicyDocument)
	}

	if Role.Description != nil {
		_ = d.Set("description", Role.Description)
	}

	if Role.ConsoleLogin != nil {
		_ = d.Set("console_login", Role.ConsoleLogin)
	}

	if Role.SessionDuration != nil {
		_ = d.Set("session_duration", Role.SessionDuration)
	}

	return nil
}

func resourceTencentCloudCamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		updateRoleDescriptionRequest  = cam.NewUpdateRoleDescriptionRequest()
		updateRoleDescriptionResponse = cam.NewUpdateRoleDescriptionResponse()
	)

	roleId := d.Id()

	request.RoleId = &roleId

	immutableArgs := []string{"role_name", "policy_document", "description", "console_login", "session_duration"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("role_name") {
		if v, ok := d.GetOk("role_name"); ok {
			request.RoleName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("console_login") {
		if v, ok := d.GetOkExists("console_login"); ok {
			request.ConsoleLogin = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateRoleDescription(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cam Role failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCamRoleRead(d, meta)
}

func resourceTencentCloudCamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	roleId := d.Id()

	if err := service.DeleteCamRoleById(ctx, roleId); err != nil {
		return err
	}

	return nil
}

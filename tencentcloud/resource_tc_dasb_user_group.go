/*
Provides a resource to create a dasb user_group

Example Usage

```hcl
resource "tencentcloud_dasb_user_group" "example" {
  name = "tf_example_update"
}
```

Or

```hcl
resource "tencentcloud_dasb_user_group" "example" {
  name          = "tf_example_update"
  department_id = "1.2"
}
```

Import

dasb user_group can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_user_group.example 16
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDasbUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbUserGroupCreate,
		Read:   resourceTencentCloudDasbUserGroupRead,
		Update: resourceTencentCloudDasbUserGroupUpdate,
		Delete: resourceTencentCloudDasbUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User group name, maximum length 32 characters.",
			},
			"department_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the department to which the user group belongs, such as: 1.2.3.",
			},
		},
	}
}

func resourceTencentCloudDasbUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_user_group.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = dasb.NewCreateUserGroupRequest()
		response    = dasb.NewCreateUserGroupResponse()
		userGroupId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		request.DepartmentId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().CreateUserGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb UserGroup failed, reason:%+v", logId, err)
		return err
	}

	userGroupIdInt := *response.Response.Id
	userGroupId = strconv.FormatUint(userGroupIdInt, 10)
	d.SetId(userGroupId)

	return resourceTencentCloudDasbUserGroupRead(d, meta)
}

func resourceTencentCloudDasbUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_user_group.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		service     = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
		userGroupId = d.Id()
	)

	UserGroup, err := service.DescribeDasbUserGroupById(ctx, userGroupId)
	if err != nil {
		return err
	}

	if UserGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbUserGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if UserGroup.Name != nil {
		_ = d.Set("name", UserGroup.Name)
	}

	if UserGroup.Department != nil {
		departmentId := *UserGroup.Department.Id
		if departmentId != "1" {
			_ = d.Set("department_id", departmentId)
		}
	}

	return nil
}

func resourceTencentCloudDasbUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_user_group.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = dasb.NewModifyUserGroupRequest()
		userGroupId = d.Id()
	)

	userGroupIdInt, _ := strconv.ParseUint(userGroupId, 10, 64)
	request.Id = &userGroupIdInt

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		request.DepartmentId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().ModifyUserGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dasb UserGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDasbUserGroupRead(d, meta)
}

func resourceTencentCloudDasbUserGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_user_group.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		service     = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
		userGroupId = d.Id()
	)

	if err := service.DeleteDasbUserGroupById(ctx, userGroupId); err != nil {
		return err
	}

	return nil
}

/*
Provides a resource to create a bi project_user_role

Example Usage

```hcl
resource "tencentcloud_bi_project_user_role" "project_user_role" {
  project_id = 123
  user_list {
		user_id = "abc"
		user_name = "abc"
		corp_id = "abc"
		email = "abc@tencent.com"
		last_login = "2023-05-11 16:59:16"
		status = 1
		first_modify = 1
		phone_number = "12345678910"
		area_code = "86"
		created_user = "abc"
		created_at = "2023-05-11 16:59:16"
		updated_user = "abc"
		updated_at = "2023-05-11 16:59:16"
		global_user_name = "abc"
		mobile = "12345678910"

  }
  user_info_list {
		user_id = "abc"
		user_name = "abc"
		email = "abc@tencent.com"
		phone_number = "12345678910"
		area_code = "86"

  }
}
```

Import

bi project_user_role can be imported using the id, e.g.

```
terraform import tencentcloud_bi_project_user_role.project_user_role project_user_role_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudBiProjectUserRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiProjectUserRoleCreate,
		Read:   resourceTencentCloudBiProjectUserRoleRead,
		Update: resourceTencentCloudBiProjectUserRoleUpdate,
		Delete: resourceTencentCloudBiProjectUserRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project id.",
			},

			"user_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "User list (deprecated).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User id.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Username.",
						},
						"corp_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enterprise id(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "E-mail(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"last_login": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last login time, public cloud unrelated fields(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Disabled state(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"first_modify": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "First login to change password, public cloud unrelated fields(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"phone_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Phone number(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"area_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Mobile area code(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"created_user": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Created by(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"created_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Created at(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"updated_user": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Updated by(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Updated at(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"global_user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Global role name(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"mobile": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Mobile number, public cloud unrelated fields(Note: This field may return null, indicating that no valid value can be obtained).",
						},
					},
				},
			},

			"user_info_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "User List (New).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User id.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Username.",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "E-mail(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"phone_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Phone number(Note: This field may return null, indicating that no valid value can be obtained).",
						},
						"area_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Mobile area code(Note: This field may return null, indicating that no valid value can be obtained).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudBiProjectUserRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_project_user_role.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = bi.NewCreateUserRoleProjectRequest()
		response = bi.NewCreateUserRoleProjectResponse()
		userId   string
	)
	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("user_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			userIdAndUserName := bi.UserIdAndUserName{}
			if v, ok := dMap["user_id"]; ok {
				userIdAndUserName.UserId = helper.String(v.(string))
			}
			if v, ok := dMap["user_name"]; ok {
				userIdAndUserName.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["corp_id"]; ok {
				userIdAndUserName.CorpId = helper.String(v.(string))
			}
			if v, ok := dMap["email"]; ok {
				userIdAndUserName.Email = helper.String(v.(string))
			}
			if v, ok := dMap["last_login"]; ok {
				userIdAndUserName.LastLogin = helper.String(v.(string))
			}
			if v, ok := dMap["status"]; ok {
				userIdAndUserName.Status = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["first_modify"]; ok {
				userIdAndUserName.FirstModify = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["phone_number"]; ok {
				userIdAndUserName.PhoneNumber = helper.String(v.(string))
			}
			if v, ok := dMap["area_code"]; ok {
				userIdAndUserName.AreaCode = helper.String(v.(string))
			}
			if v, ok := dMap["created_user"]; ok {
				userIdAndUserName.CreatedUser = helper.String(v.(string))
			}
			if v, ok := dMap["created_at"]; ok {
				userIdAndUserName.CreatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["updated_user"]; ok {
				userIdAndUserName.UpdatedUser = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				userIdAndUserName.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["global_user_name"]; ok {
				userIdAndUserName.GlobalUserName = helper.String(v.(string))
			}
			if v, ok := dMap["mobile"]; ok {
				userIdAndUserName.Mobile = helper.String(v.(string))
			}
			request.UserList = append(request.UserList, &userIdAndUserName)
		}
	}

	if v, ok := d.GetOk("user_info_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			userInfo := bi.UserInfo{}
			if v, ok := dMap["user_id"]; ok {
				userInfo.UserId = helper.String(v.(string))
			}
			if v, ok := dMap["user_name"]; ok {
				userInfo.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["email"]; ok {
				userInfo.Email = helper.String(v.(string))
			}
			if v, ok := dMap["phone_number"]; ok {
				userInfo.PhoneNumber = helper.String(v.(string))
			}
			if v, ok := dMap["area_code"]; ok {
				userInfo.AreaCode = helper.String(v.(string))
			}
			request.UserInfoList = append(request.UserInfoList, &userInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().CreateUserRoleProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create bi projectUserRole failed, reason:%+v", logId, err)
		return err
	}

	userId = *response.Response.UserId
	d.SetId(userId)

	return resourceTencentCloudBiProjectUserRoleRead(d, meta)
}

func resourceTencentCloudBiProjectUserRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_project_user_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}

	projectUserRoleId := d.Id()

	projectUserRole, err := service.DescribeBiProjectUserRoleById(ctx, userId)
	if err != nil {
		return err
	}

	if projectUserRole == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `BiProjectUserRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if projectUserRole.ProjectId != nil {
		_ = d.Set("project_id", projectUserRole.ProjectId)
	}

	if projectUserRole.UserList != nil {
		userListList := []interface{}{}
		for _, userList := range projectUserRole.UserList {
			userListMap := map[string]interface{}{}

			if projectUserRole.UserList.UserId != nil {
				userListMap["user_id"] = projectUserRole.UserList.UserId
			}

			if projectUserRole.UserList.UserName != nil {
				userListMap["user_name"] = projectUserRole.UserList.UserName
			}

			if projectUserRole.UserList.CorpId != nil {
				userListMap["corp_id"] = projectUserRole.UserList.CorpId
			}

			if projectUserRole.UserList.Email != nil {
				userListMap["email"] = projectUserRole.UserList.Email
			}

			if projectUserRole.UserList.LastLogin != nil {
				userListMap["last_login"] = projectUserRole.UserList.LastLogin
			}

			if projectUserRole.UserList.Status != nil {
				userListMap["status"] = projectUserRole.UserList.Status
			}

			if projectUserRole.UserList.FirstModify != nil {
				userListMap["first_modify"] = projectUserRole.UserList.FirstModify
			}

			if projectUserRole.UserList.PhoneNumber != nil {
				userListMap["phone_number"] = projectUserRole.UserList.PhoneNumber
			}

			if projectUserRole.UserList.AreaCode != nil {
				userListMap["area_code"] = projectUserRole.UserList.AreaCode
			}

			if projectUserRole.UserList.CreatedUser != nil {
				userListMap["created_user"] = projectUserRole.UserList.CreatedUser
			}

			if projectUserRole.UserList.CreatedAt != nil {
				userListMap["created_at"] = projectUserRole.UserList.CreatedAt
			}

			if projectUserRole.UserList.UpdatedUser != nil {
				userListMap["updated_user"] = projectUserRole.UserList.UpdatedUser
			}

			if projectUserRole.UserList.UpdatedAt != nil {
				userListMap["updated_at"] = projectUserRole.UserList.UpdatedAt
			}

			if projectUserRole.UserList.GlobalUserName != nil {
				userListMap["global_user_name"] = projectUserRole.UserList.GlobalUserName
			}

			if projectUserRole.UserList.Mobile != nil {
				userListMap["mobile"] = projectUserRole.UserList.Mobile
			}

			userListList = append(userListList, userListMap)
		}

		_ = d.Set("user_list", userListList)

	}

	if projectUserRole.UserInfoList != nil {
		userInfoListList := []interface{}{}
		for _, userInfoList := range projectUserRole.UserInfoList {
			userInfoListMap := map[string]interface{}{}

			if projectUserRole.UserInfoList.UserId != nil {
				userInfoListMap["user_id"] = projectUserRole.UserInfoList.UserId
			}

			if projectUserRole.UserInfoList.UserName != nil {
				userInfoListMap["user_name"] = projectUserRole.UserInfoList.UserName
			}

			if projectUserRole.UserInfoList.Email != nil {
				userInfoListMap["email"] = projectUserRole.UserInfoList.Email
			}

			if projectUserRole.UserInfoList.PhoneNumber != nil {
				userInfoListMap["phone_number"] = projectUserRole.UserInfoList.PhoneNumber
			}

			if projectUserRole.UserInfoList.AreaCode != nil {
				userInfoListMap["area_code"] = projectUserRole.UserInfoList.AreaCode
			}

			userInfoListList = append(userInfoListList, userInfoListMap)
		}

		_ = d.Set("user_info_list", userInfoListList)

	}

	return nil
}

func resourceTencentCloudBiProjectUserRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_project_user_role.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := bi.NewModifyUserRoleProjectRequest()

	projectUserRoleId := d.Id()

	request.UserId = &userId

	immutableArgs := []string{"project_id", "user_list", "user_info_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOkExists("project_id"); ok {
			request.ProjectId = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().ModifyUserRoleProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update bi projectUserRole failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudBiProjectUserRoleRead(d, meta)
}

func resourceTencentCloudBiProjectUserRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_project_user_role.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}
	projectUserRoleId := d.Id()

	if err := service.DeleteBiProjectUserRoleById(ctx, userId); err != nil {
		return err
	}

	return nil
}

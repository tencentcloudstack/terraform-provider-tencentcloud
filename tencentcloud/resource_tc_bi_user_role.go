/*
Provides a resource to create a bi user_role

Example Usage

```hcl
resource "tencentcloud_bi_user_role" "user_role" {
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

bi user_role can be imported using the id, e.g.

```
terraform import tencentcloud_bi_user_role.user_role user_role_id
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

func resourceTencentCloudBiUserRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiUserRoleCreate,
		Read:   resourceTencentCloudBiUserRoleRead,
		Update: resourceTencentCloudBiUserRoleUpdate,
		Delete: resourceTencentCloudBiUserRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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

func resourceTencentCloudBiUserRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_user_role.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = bi.NewCreateUserRoleRequest()
		response = bi.NewCreateUserRoleResponse()
		userId   string
	)
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().CreateUserRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create bi userRole failed, reason:%+v", logId, err)
		return err
	}

	userId = *response.Response.UserId
	d.SetId(userId)

	return resourceTencentCloudBiUserRoleRead(d, meta)
}

func resourceTencentCloudBiUserRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_user_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}

	userRoleId := d.Id()

	userRole, err := service.DescribeBiUserRoleById(ctx, userId)
	if err != nil {
		return err
	}

	if userRole == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `BiUserRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if userRole.UserList != nil {
		userListList := []interface{}{}
		for _, userList := range userRole.UserList {
			userListMap := map[string]interface{}{}

			if userRole.UserList.UserId != nil {
				userListMap["user_id"] = userRole.UserList.UserId
			}

			if userRole.UserList.UserName != nil {
				userListMap["user_name"] = userRole.UserList.UserName
			}

			if userRole.UserList.CorpId != nil {
				userListMap["corp_id"] = userRole.UserList.CorpId
			}

			if userRole.UserList.Email != nil {
				userListMap["email"] = userRole.UserList.Email
			}

			if userRole.UserList.LastLogin != nil {
				userListMap["last_login"] = userRole.UserList.LastLogin
			}

			if userRole.UserList.Status != nil {
				userListMap["status"] = userRole.UserList.Status
			}

			if userRole.UserList.FirstModify != nil {
				userListMap["first_modify"] = userRole.UserList.FirstModify
			}

			if userRole.UserList.PhoneNumber != nil {
				userListMap["phone_number"] = userRole.UserList.PhoneNumber
			}

			if userRole.UserList.AreaCode != nil {
				userListMap["area_code"] = userRole.UserList.AreaCode
			}

			if userRole.UserList.CreatedUser != nil {
				userListMap["created_user"] = userRole.UserList.CreatedUser
			}

			if userRole.UserList.CreatedAt != nil {
				userListMap["created_at"] = userRole.UserList.CreatedAt
			}

			if userRole.UserList.UpdatedUser != nil {
				userListMap["updated_user"] = userRole.UserList.UpdatedUser
			}

			if userRole.UserList.UpdatedAt != nil {
				userListMap["updated_at"] = userRole.UserList.UpdatedAt
			}

			if userRole.UserList.GlobalUserName != nil {
				userListMap["global_user_name"] = userRole.UserList.GlobalUserName
			}

			if userRole.UserList.Mobile != nil {
				userListMap["mobile"] = userRole.UserList.Mobile
			}

			userListList = append(userListList, userListMap)
		}

		_ = d.Set("user_list", userListList)

	}

	if userRole.UserInfoList != nil {
		userInfoListList := []interface{}{}
		for _, userInfoList := range userRole.UserInfoList {
			userInfoListMap := map[string]interface{}{}

			if userRole.UserInfoList.UserId != nil {
				userInfoListMap["user_id"] = userRole.UserInfoList.UserId
			}

			if userRole.UserInfoList.UserName != nil {
				userInfoListMap["user_name"] = userRole.UserInfoList.UserName
			}

			if userRole.UserInfoList.Email != nil {
				userInfoListMap["email"] = userRole.UserInfoList.Email
			}

			if userRole.UserInfoList.PhoneNumber != nil {
				userInfoListMap["phone_number"] = userRole.UserInfoList.PhoneNumber
			}

			if userRole.UserInfoList.AreaCode != nil {
				userInfoListMap["area_code"] = userRole.UserInfoList.AreaCode
			}

			userInfoListList = append(userInfoListList, userInfoListMap)
		}

		_ = d.Set("user_info_list", userInfoListList)

	}

	return nil
}

func resourceTencentCloudBiUserRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_user_role.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := bi.NewModifyUserRoleRequest()

	userRoleId := d.Id()

	request.UserId = &userId

	immutableArgs := []string{"user_list", "user_info_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().ModifyUserRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update bi userRole failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudBiUserRoleRead(d, meta)
}

func resourceTencentCloudBiUserRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_user_role.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}
	userRoleId := d.Id()

	if err := service.DeleteBiUserRoleById(ctx, userId); err != nil {
		return err
	}

	return nil
}

package bi

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBiUserRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiUserRoleCreate,
		Read:   resourceTencentCloudBiUserRoleRead,
		Update: resourceTencentCloudBiUserRoleUpdate,
		Delete: resourceTencentCloudBiUserRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"role_id_list": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Role id list.",
			},

			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User id.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username.",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "E-mail(Note: This field may return null, indicating that no valid value can be obtained).",
			},
			"phone_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Phone number(Note: This field may return null, indicating that no valid value can be obtained).",
			},
			"area_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mobile area code(Note: This field may return null, indicating that no valid value can be obtained).",
			},
		},
	}
}

func resourceTencentCloudBiUserRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_user_role.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = bi.NewCreateUserRoleRequest()
		userId  string
	)
	if v, ok := d.GetOk("role_id_list"); ok {
		roleIdListSet := v.(*schema.Set).List()
		for i := range roleIdListSet {
			roleIdList := roleIdListSet[i].(int)
			request.RoleIdList = append(request.RoleIdList, helper.IntInt64(roleIdList))
		}
	}

	var userInfo bi.UserInfo
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		userInfo.UserId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("user_name"); ok {
		userInfo.UserName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("email"); ok {
		userInfo.Email = helper.String(v.(string))
	}
	if v, ok := d.GetOk("phone_number"); ok {
		userInfo.PhoneNumber = helper.String(v.(string))
	}
	if v, ok := d.GetOk("area_code"); ok {
		userInfo.AreaCode = helper.String(v.(string))
	}
	request.UserInfoList = append(request.UserInfoList, &userInfo)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBiClient().CreateUserRole(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create bi userRole failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(userId)

	return resourceTencentCloudBiUserRoleRead(d, meta)
}

func resourceTencentCloudBiUserRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_user_role.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := BiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	userId := d.Id()
	userRole, err := service.DescribeBiUserRoleById(ctx, userId)
	if err != nil {
		return err
	}

	if userRole == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `BiUserRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if userRole.RoleIdList != nil {
		_ = d.Set("role_id_list", userRole.RoleIdList)
	}

	if userRole.UserId != nil {
		_ = d.Set("user_id", userRole.UserId)
	}

	if userRole.UserName != nil {
		_ = d.Set("user_name", userRole.UserName)
	}

	if userRole.Email != nil {
		_ = d.Set("email", userRole.Email)
	}

	if userRole.PhoneNumber != nil {
		_ = d.Set("phone_number", userRole.PhoneNumber)
	}

	if userRole.AreaCode != nil {
		_ = d.Set("area_code", userRole.AreaCode)
	}

	return nil
}

func resourceTencentCloudBiUserRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_user_role.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := bi.NewModifyUserRoleRequest()

	userId := d.Id()
	request.UserId = &userId

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = helper.String(v.(string))
	}

	if d.HasChange("role_id_list") {
		if v, ok := d.GetOk("role_id_list"); ok {
			roleIdListSet := v.(*schema.Set).List()
			for i := range roleIdListSet {
				roleIdList := roleIdListSet[i].(int)
				request.RoleIdList = append(request.RoleIdList, helper.IntInt64(roleIdList))
			}
		}
	}

	if d.HasChange("email") {
		if v, ok := d.GetOk("email"); ok {
			request.Email = helper.String(v.(string))
		}
	}

	if d.HasChange("phone_number") {
		if v, ok := d.GetOk("phone_number"); ok {
			request.PhoneNumber = helper.String(v.(string))
		}
	}

	if d.HasChange("area_code") {
		if v, ok := d.GetOk("area_code"); ok {
			request.AreaCode = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBiClient().ModifyUserRole(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_bi_user_role.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := BiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	userId := d.Id()

	if err := service.DeleteBiUserRoleById(ctx, userId); err != nil {
		return err
	}

	return nil
}

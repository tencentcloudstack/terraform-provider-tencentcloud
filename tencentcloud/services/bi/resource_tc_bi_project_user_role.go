package bi

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBiProjectUserRole() *schema.Resource {
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

			"role_id_list": {
				Optional: true,
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

func resourceTencentCloudBiProjectUserRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_project_user_role.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = bi.NewCreateUserRoleProjectRequest()
		projectId int
		userId    string
	)
	if v, ok := d.GetOkExists("project_id"); ok {
		projectId = v.(int)
		request.ProjectId = helper.IntInt64(v.(int))
	}

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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBiClient().CreateUserRoleProject(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create bi projectUserRole failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{strconv.Itoa(projectId), userId}, tccommon.FILED_SP))

	return resourceTencentCloudBiProjectUserRoleRead(d, meta)
}

func resourceTencentCloudBiProjectUserRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_project_user_role.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := BiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	projectIdInt, _ := strconv.ParseInt(projectId, 10, 64)
	userId := idSplit[1]

	userRole, err := service.DescribeBiProjectUserRoleById(ctx, projectIdInt, userId)
	if err != nil {
		return err
	}

	if userRole == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `BiProjectUserRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectIdInt)

	if userRole.RoleList != nil {
		var roles []*int64
		for _, v := range userRole.RoleList {
			roles = append(roles, v.RoleId)
		}
		_ = d.Set("role_id_list", roles)
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

func resourceTencentCloudBiProjectUserRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_project_user_role.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := bi.NewModifyUserRoleProjectRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	projectIdInt, _ := strconv.ParseInt(projectId, 10, 64)
	userId := idSplit[1]

	request.UserId = &userId
	request.ProjectId = &projectIdInt

	immutableArgs := []string{"area_code", "phone_number"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBiClient().ModifyUserRoleProject(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_bi_project_user_role.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := BiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	projectIdInt, _ := strconv.ParseInt(projectId, 10, 64)
	userId := idSplit[1]

	if err := service.DeleteBiProjectUserRoleById(ctx, projectIdInt, userId); err != nil {
		return err
	}

	return nil
}

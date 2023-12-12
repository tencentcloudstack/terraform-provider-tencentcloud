package bh

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDasbUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbUserCreate,
		Read:   resourceTencentCloudDasbUserRead,
		Update: resourceTencentCloudDasbUserUpdate,
		Delete: resourceTencentCloudDasbUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Username, 3-20 characters, must start with an English letter and cannot contain characters other than letters, numbers, '.', '_', '-'.",
			},
			"real_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Real name, maximum length 20 characters, cannot contain blank characters.",
			},
			"phone": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fill in the mainland mobile phone number directly. If it is a number from other countries or regions, enter it in the format of country area code|mobile phone number. For example: +852|xxxxxxxx.",
			},
			"email": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Email.",
			},
			"validate_from": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User effective time, such as: 2021-09-22T00:00:00+00:00If the effective and expiry time are not filled in, the user will be valid for a long time.",
			},
			"validate_to": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User expiration time, such as: 2021-09-23T00:00:00+00:00If the effective and expiry time are not filled in, the user will be valid for a long time.",
			},
			"group_id_set": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "The set of user group IDs to which it belongs.",
			},
			"auth_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Authentication method, 0 - local, 1 - LDAP, 2 - OAuth. If not passed, the default is 0.",
			},
			"validate_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Access time period limit, a string composed of 0 and 1, length 168 (7 * 24), representing the time period the user is allowed to access in a week. The Nth character in the string represents the Nth hour of the week, 0 - means access is not allowed, 1 - means access is allowed.",
			},
			"department_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Department ID, such as: 1.2.3.",
			},
		},
	}
}

func resourceTencentCloudDasbUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = dasb.NewCreateUserRequest()
		response = dasb.NewCreateUserResponse()
		userId   string
	)

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("real_name"); ok {
		request.RealName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("phone"); ok {
		request.Phone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("email"); ok {
		request.Email = helper.String(v.(string))
	}

	if v, ok := d.GetOk("validate_from"); ok {
		request.ValidateFrom = helper.String(v.(string))
	}

	if v, ok := d.GetOk("validate_to"); ok {
		request.ValidateTo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id_set"); ok {
		groupIdSetSet := v.(*schema.Set).List()
		for i := range groupIdSetSet {
			groupIdSet := groupIdSetSet[i].(int)
			request.GroupIdSet = append(request.GroupIdSet, helper.IntUint64(groupIdSet))
		}
	}

	if v, ok := d.GetOkExists("auth_type"); ok {
		request.AuthType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("validate_time"); ok {
		request.ValidateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		request.DepartmentId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().CreateUser(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response.Id == nil {
			e = fmt.Errorf("dasb user not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb user failed, reason:%+v", logId, err)
		return err
	}

	userIdInt := *response.Response.Id
	userId = strconv.FormatUint(userIdInt, 10)
	d.SetId(userId)

	return resourceTencentCloudDasbUserRead(d, meta)
}

func resourceTencentCloudDasbUserRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userId  = d.Id()
	)

	user, err := service.DescribeDasbUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbUser` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if user.UserName != nil {
		_ = d.Set("user_name", user.UserName)
	}

	if user.RealName != nil {
		_ = d.Set("real_name", user.RealName)
	}

	if user.Phone != nil {
		_ = d.Set("phone", user.Phone)
	}

	if user.Email != nil {
		_ = d.Set("email", user.Email)
	}

	if user.ValidateFrom != nil {
		_ = d.Set("validate_from", user.ValidateFrom)
	}

	if user.ValidateTo != nil {
		_ = d.Set("validate_to", user.ValidateTo)
	}

	if user.GroupSet != nil {
		tmpList := make([]*uint64, 0)
		for _, item := range user.GroupSet {
			if item.Id != nil {
				tmpList = append(tmpList, item.Id)
			}
		}

		_ = d.Set("group_id_set", tmpList)
	}

	if user.AuthType != nil {
		_ = d.Set("auth_type", user.AuthType)
	}

	if user.ValidateTime != nil {
		_ = d.Set("validate_time", user.ValidateTime)
	}

	if user.DepartmentId != nil {
		if *user.DepartmentId != "1" {
			_ = d.Set("department_id", user.DepartmentId)
		}
	}

	return nil
}

func resourceTencentCloudDasbUserUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = dasb.NewModifyUserRequest()
		userId  = d.Id()
	)

	immutableArgs := []string{"user_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.Id = helper.StrToUint64Point(userId)

	if v, ok := d.GetOk("real_name"); ok {
		request.RealName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("phone"); ok {
		request.Phone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("email"); ok {
		request.Email = helper.String(v.(string))
	}

	if v, ok := d.GetOk("validate_from"); ok {
		request.ValidateFrom = helper.String(v.(string))
	}

	if v, ok := d.GetOk("validate_to"); ok {
		request.ValidateTo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id_set"); ok {
		groupIdSetSet := v.(*schema.Set).List()
		for i := range groupIdSetSet {
			groupIdSet := groupIdSetSet[i].(int)
			request.GroupIdSet = append(request.GroupIdSet, helper.IntUint64(groupIdSet))
		}
	}

	if v, ok := d.GetOkExists("auth_type"); ok {
		request.AuthType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("validate_time"); ok {
		request.ValidateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		request.DepartmentId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().ModifyUser(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dasb user failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDasbUserRead(d, meta)
}

func resourceTencentCloudDasbUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userId  = d.Id()
	)

	if err := service.DeleteDasbUserById(ctx, userId); err != nil {
		return err
	}

	return nil
}

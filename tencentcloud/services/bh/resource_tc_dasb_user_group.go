package bh

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDasbUserGroup() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().CreateUserGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
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
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().ModifyUserGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userGroupId = d.Id()
	)

	if err := service.DeleteDasbUserGroupById(ctx, userGroupId); err != nil {
		return err
	}

	return nil
}

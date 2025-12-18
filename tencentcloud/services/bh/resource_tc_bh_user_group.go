package bh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhUserGroupCreate,
		Read:   resourceTencentCloudBhUserGroupRead,
		Update: resourceTencentCloudBhUserGroupUpdate,
		Delete: resourceTencentCloudBhUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User group name, maximum length 32 characters.",
			},

			"department_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Department ID to which the user group belongs, e.g.: 1.2.3.",
			},

			// computed
			"user_group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "User group ID.",
			},
		},
	}
}

func resourceTencentCloudBhUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request     = bhv20230418.NewCreateUserGroupRequest()
		response    = bhv20230418.NewCreateUserGroupResponse()
		userGroupId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		request.DepartmentId = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().CreateUserGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bh user group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bh user group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Id == nil {
		return fmt.Errorf("Id is nil.")
	}

	userGroupId = helper.UInt64ToStr(*response.Response.Id)
	d.SetId(userGroupId)
	return resourceTencentCloudBhUserGroupRead(d, meta)
}

func resourceTencentCloudBhUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service     = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userGroupId = d.Id()
	)

	respData, err := service.DescribeBhUserGroupById(ctx, userGroupId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_user_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Department != nil {
		if respData.Department.Id != nil {
			dResp, err := service.DescribeBhDepartments(ctx)
			if err != nil {
				return err
			}

			if dResp == nil {
				return fmt.Errorf("Departments is nil")
			}

			if dResp.Enabled != nil && *dResp.Enabled {
				_ = d.Set("department_id", respData.Department.Id)
			}
		}
	}

	if respData.Id != nil {
		_ = d.Set("user_group_id", respData.Id)
	}

	return nil
}

func resourceTencentCloudBhUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		userGroupId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"name", "department_id"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := bhv20230418.NewModifyUserGroupRequest()
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("department_id"); ok {
			request.DepartmentId = helper.String(v.(string))
		}

		request.Id = helper.StrToUint64Point(userGroupId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyUserGroupWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bh user group failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBhUserGroupRead(d, meta)
}

func resourceTencentCloudBhUserGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request     = bhv20230418.NewDeleteUserGroupsRequest()
		userGroupId = d.Id()
	)

	request.IdSet = append(request.IdSet, helper.StrToUint64Point(userGroupId))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().DeleteUserGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bh user group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

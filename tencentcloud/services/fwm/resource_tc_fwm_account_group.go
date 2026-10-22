package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudFwmAccountGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmAccountGroupCreate,
		Read:   resourceTencentCloudFwmAccountGroupRead,
		Update: resourceTencentCloudFwmAccountGroupUpdate,
		Delete: resourceTencentCloudFwmAccountGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account group name, up to 50 characters, unique under the same admin AppId.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the account group, only settable on create.",
			},

			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account group ID, e.g. acg-xxxxxxxx.",
			},

			"member_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Member count of the account group.",
			},

			"can_delete": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the account group can be deleted.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the account group.",
			},
		},
	}
}

func resourceTencentCloudFwmAccountGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_account_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = fwmv20250611.NewCreateAccountGroupRequest()
		response = fwmv20250611.NewCreateAccountGroupResponse()
	)

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateAccountGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm account group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm account group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.GroupId == nil {
		log.Printf("[CRITAL]%s create fwm account group return empty GroupId, logId=%s, d.Id()=%s", logId, logId, d.Id())
		return fmt.Errorf("GroupId is nil, create fwm account group failed.")
	}

	groupId := *response.Response.GroupId
	d.SetId(groupId)
	return resourceTencentCloudFwmAccountGroupRead(d, meta)
}

func resourceTencentCloudFwmAccountGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_account_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	groupId := d.Id()

	respData, err := service.DescribeFwmAccountGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] fwm_account_group id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.GroupId != nil {
		_ = d.Set("group_id", respData.GroupId)
	}

	if respData.GroupName != nil {
		_ = d.Set("group_name", respData.GroupName)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.MemberCount != nil {
		_ = d.Set("member_count", respData.MemberCount)
	}

	if respData.CanDelete != nil {
		_ = d.Set("can_delete", respData.CanDelete)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	return nil
}

func resourceTencentCloudFwmAccountGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_account_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	needChange := false
	mutableArgs := []string{"group_name"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		groupId := d.Id()
		request := fwmv20250611.NewModifyAccountGroupRequest()
		request.GroupId = helper.String(groupId)
		if v, ok := d.GetOk("group_name"); ok {
			request.GroupName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyAccountGroupWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify fwm account group failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update fwm account group failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudFwmAccountGroupRead(d, meta)
}

func resourceTencentCloudFwmAccountGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_account_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	groupId := d.Id()

	err := service.DeleteFwmAccountGroupById(ctx, groupId)
	if err != nil {
		log.Printf("[CRITAL]%s delete fwm account group failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

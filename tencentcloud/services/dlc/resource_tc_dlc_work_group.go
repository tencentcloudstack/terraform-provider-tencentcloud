package dlc

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcWorkGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcWorkGroupCreate,
		Read:   resourceTencentCloudDlcWorkGroupRead,
		Update: resourceTencentCloudDlcWorkGroupUpdate,
		Delete: resourceTencentCloudDlcWorkGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"work_group_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Working group name.",
			},

			"work_group_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Working group description.",
			},

			"user_ids": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Collection of IDs of users to be bound to working groups.",
			},

			// computed
			"work_group_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Working group ID.",
			},
		},
	}
}

func resourceTencentCloudDlcWorkGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_work_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = dlc.NewCreateWorkGroupRequest()
		response    = dlc.NewCreateWorkGroupResponse()
		workGroupId int64
	)

	if v, ok := d.GetOk("work_group_name"); ok {
		request.WorkGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_group_description"); ok {
		request.WorkGroupDescription = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateWorkGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc workGroup failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dlc workGroup failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.WorkGroupId == nil {
		return fmt.Errorf("WorkGroupId is nil.")
	}

	workGroupId = *response.Response.WorkGroupId
	d.SetId(helper.Int64ToStr(workGroupId))

	return resourceTencentCloudDlcWorkGroupRead(d, meta)
}

func resourceTencentCloudDlcWorkGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_work_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		workGroupId = d.Id()
	)

	workGroup, err := service.DescribeDlcWorkGroupById(ctx, workGroupId)
	if err != nil {
		return err
	}

	if workGroup == nil {
		log.Printf("[WARN]%s resource `DlcWorkGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if workGroup.WorkGroupName != nil {
		_ = d.Set("work_group_name", workGroup.WorkGroupName)
	}

	if workGroup.WorkGroupDescription != nil {
		_ = d.Set("work_group_description", workGroup.WorkGroupDescription)
	}

	if workGroup.UserSet != nil {
		userIds := make([]*string, 0, len(workGroup.UserSet))
		for _, user := range workGroup.UserSet {
			userIds = append(userIds, user.UserId)
		}

		_ = d.Set("user_ids", userIds)
	}

	if workGroup.WorkGroupId != nil {
		_ = d.Set("work_group_id", workGroup.WorkGroupId)
	}

	return nil
}

func resourceTencentCloudDlcWorkGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_work_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = dlc.NewModifyWorkGroupRequest()
		workGroupId = d.Id()
	)

	if d.HasChange("work_group_description") {
		if v, ok := d.GetOk("work_group_description"); ok {
			request.WorkGroupDescription = helper.String(v.(string))
		}

		request.WorkGroupId = helper.Int64(helper.StrToInt64(workGroupId))
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().ModifyWorkGroup(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update dlc workGroup failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudDlcWorkGroupRead(d, meta)
}

func resourceTencentCloudDlcWorkGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_work_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		workGroupId = d.Id()
	)

	if err := service.DeleteDlcWorkGroupById(ctx, workGroupId); err != nil {
		return err
	}

	return nil
}

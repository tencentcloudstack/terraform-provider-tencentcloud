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

func ResourceTencentCloudDlcBindWorkGroupsToUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcBindWorkGroupsToUserCreateAttachment,
		Read:   resourceTencentCloudDlcBindWorkGroupsToUserReadAttachment,
		Delete: resourceTencentCloudDlcBindWorkGroupsToUserDeleteAttachment,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"add_info": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Information about bound working groups and users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "User ID, which matches Uin on the CAM side.",
						},
						"work_group_ids": {
							Type:        schema.TypeSet,
							Required:    true,
							ForceNew:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "Collections of IDs of working groups.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcBindWorkGroupsToUserCreateAttachment(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_bind_work_groups_to_user_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = dlc.NewBindWorkGroupsToUserRequest()
		userId  string
	)

	if dMap, ok := helper.InterfacesHeadMap(d, "add_info"); ok {
		workGroupIdSetOfUserId := dlc.WorkGroupIdSetOfUserId{}
		if v, ok := dMap["user_id"]; ok {
			workGroupIdSetOfUserId.UserId = helper.String(v.(string))
			userId = v.(string)
		}

		if v, ok := dMap["work_group_ids"]; ok {
			workGroupIdsSet := v.(*schema.Set).List()
			for i := range workGroupIdsSet {
				workGroupIds := workGroupIdsSet[i].(int)
				workGroupIdSetOfUserId.WorkGroupIds = append(workGroupIdSetOfUserId.WorkGroupIds, helper.IntInt64(workGroupIds))
			}
		}

		request.AddInfo = &workGroupIdSetOfUserId
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().BindWorkGroupsToUser(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc bind work groups to user failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dlc bind work groups to user failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(userId)
	return resourceTencentCloudDlcBindWorkGroupsToUserReadAttachment(d, meta)
}

func resourceTencentCloudDlcBindWorkGroupsToUserReadAttachment(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_bind_work_groups_to_user_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	parm := make(map[string]interface{})
	parm["UserId"] = helper.String(d.Id())
	parm["Type"] = helper.String("Group")
	bindWorkGroupsToUser, err := service.DescribeDlcDescribeUserInfoByFilter(ctx, parm)
	if err != nil {
		return err
	}

	if bindWorkGroupsToUser == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dlc_bind_work_groups_to_user_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	var workId []int64
	if bindWorkGroupsToUser.WorkGroupInfo != nil {
		addInfoMap := map[string]interface{}{}
		if len(bindWorkGroupsToUser.WorkGroupInfo.WorkGroupSet) > 1 {
			for _, v := range bindWorkGroupsToUser.WorkGroupInfo.WorkGroupSet {
				if v.WorkGroupId != nil {
					workId = append(workId, *v.WorkGroupId)
				}
			}

			addInfoMap["work_group_ids"] = workId
		}

		addInfoMap["user_id"] = parm["UserId"]
		_ = d.Set("add_info", []interface{}{addInfoMap})
	}

	return nil
}

func resourceTencentCloudDlcBindWorkGroupsToUserDeleteAttachment(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_bind_work_groups_to_user_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userId         = d.Id()
		workGroupIdSet []*int64
	)

	if dMap, ok := helper.InterfacesHeadMap(d, "add_info"); ok {
		if v, ok := dMap["work_group_ids"]; ok {
			workGroupIdsSet := v.(*schema.Set).List()
			for i := range workGroupIdsSet {
				workGroupIds := workGroupIdsSet[i].(int)
				workGroupIdSet = append(workGroupIdSet, helper.IntInt64(workGroupIds))
			}
		}
	}

	if err := service.DeleteDlcBindWorkGroupsToUserById(ctx, userId, workGroupIdSet); err != nil {
		return err
	}

	return nil
}

package dlc

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcAddUsersToWorkGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcAddUsersToWorkGroupAttachmentCreate,
		Read:   resourceTencentCloudDlcAddUsersToWorkGroupAttachmentRead,
		Delete: resourceTencentCloudDlcAddUsersToWorkGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"add_info": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Work group and user information to operate on.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"work_group_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Work group id.",
						},
						"user_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "User id set, matched with CAM side uin.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcAddUsersToWorkGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_add_users_to_work_group_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request     = dlc.NewAddUsersToWorkGroupRequest()
		workGroupId string
		ids         []string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "add_info"); ok {
		userIdSetOfWorkGroupId := dlc.UserIdSetOfWorkGroupId{}
		if v, ok := dMap["work_group_id"]; ok {
			workGroupId = helper.IntToStr(v.(int))
			userIdSetOfWorkGroupId.WorkGroupId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["user_ids"]; ok {
			userIdsSet := v.(*schema.Set).List()
			for i := range userIdsSet {
				userIds := userIdsSet[i].(string)
				ids = append(ids, userIds)
				userIdSetOfWorkGroupId.UserIds = append(userIdSetOfWorkGroupId.UserIds, &userIds)
			}
		}
		request.AddInfo = &userIdSetOfWorkGroupId
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().AddUsersToWorkGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dlc addUsersToWorkGroupAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(workGroupId + tccommon.FILED_SP + strings.Join(ids, "|"))

	return resourceTencentCloudDlcAddUsersToWorkGroupAttachmentRead(d, meta)
}

func resourceTencentCloudDlcAddUsersToWorkGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_add_users_to_work_group_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	workGroupId := idSplit[0]
	userIds := idSplit[1]

	ids := strings.Split(userIds, "|")
	if len(ids) < 1 {
		return fmt.Errorf("ids is null,%s", d.Id())
	}

	addUsersToWorkGroupAttachment, err := service.DescribeDlcWorkGroupById(ctx, workGroupId)
	if err != nil {
		return err
	}

	if addUsersToWorkGroupAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DlcAddUsersToWorkGroupAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if addUsersToWorkGroupAttachment.UserSet != nil {
		userMap := make(map[string]struct{}, len(addUsersToWorkGroupAttachment.UserSet))
		for _, user := range addUsersToWorkGroupAttachment.UserSet {
			userMap[*user.UserId] = struct{}{}
		}

		for _, id := range ids {
			if _, ok := userMap[id]; !ok {
				return fmt.Errorf("AddUsersToWorkGroup fail, id %s,workGroupId %s", id, workGroupId)
			}
		}
	}

	return nil
}

func resourceTencentCloudDlcAddUsersToWorkGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_add_users_to_work_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	workGroupId := idSplit[0]
	userIds := idSplit[1]

	ids := strings.Split(userIds, "|")
	if len(ids) < 1 {
		return fmt.Errorf("ids is null,%s", d.Id())
	}

	if err := service.DeleteDlcUsersToWorkGroupAttachmentById(ctx, workGroupId, ids); err != nil {
		return err
	}

	return nil
}

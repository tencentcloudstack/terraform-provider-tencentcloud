package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcBindWorkGroupsToUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcBindWorkGroupsToUserCreateAttachment,
		Read:   resourceTencentCloudDlcBindWorkGroupsToUserReadAttachment,
		Delete: resourceTencentCloudDlcBindWorkGroupsToUserDeleteAttachment,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"add_info": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Bind user and workgroup information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User id, matched with CAM side uin.",
						},
						"work_group_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Required:    true,
							Description: "Work group id set.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcBindWorkGroupsToUserCreateAttachment(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_bind_work_groups_to_user_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = dlc.NewBindWorkGroupsToUserRequest()
		userId  string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "add_info"); ok {
		workGroupIdSetOfUserId := dlc.WorkGroupIdSetOfUserId{}
		if v, ok := dMap["user_id"]; ok {
			userId = v.(string)
			workGroupIdSetOfUserId.UserId = helper.String(v.(string))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().BindWorkGroupsToUser(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dlc bindWorkGroupsToUser failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(userId)

	return resourceTencentCloudDlcBindWorkGroupsToUserReadAttachment(d, meta)
}

func resourceTencentCloudDlcBindWorkGroupsToUserReadAttachment(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_bind_work_groups_to_user_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	parm := make(map[string]interface{})
	parm["UserId"] = helper.String(d.Id())
	parm["Type"] = helper.String("Group")
	bindWorkGroupsToUser, err := service.DescribeDlcDescribeUserInfoByFilter(ctx, parm)
	if err != nil {
		return err
	}

	if bindWorkGroupsToUser == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DlcBindWorkGroupsToUser` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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
	defer logElapsed("resource.tencentcloud_dlc_bind_work_groups_to_user_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	userId := d.Id()
	var workGroupIdSet []*int64
	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}
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

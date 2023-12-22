package tsf

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
)

func ResourceTencentCloudTsfOperateGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfOperateGroupCreate,
		Read:   resourceTencentCloudTsfOperateGroupRead,
		Update: resourceTencentCloudTsfOperateGroupUpdate,
		Delete: resourceTencentCloudTsfOperateGroupDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "group id.",
			},

			"operate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Operation, `start`- start the group, `stop`- stop the group.",
			},
		},
	}
}

func resourceTencentCloudTsfOperateGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_operate_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var groupId string
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}

	d.SetId(groupId)

	return resourceTencentCloudTsfOperateGroupUpdate(d, meta)
}

func resourceTencentCloudTsfOperateGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_operate_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	groupId := d.Id()
	startGroup, err := service.DescribeTsfStartGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if startGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfOperateGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if startGroup.GroupId != nil {
		_ = d.Set("group_id", startGroup.GroupId)
	}

	return nil
}

func resourceTencentCloudTsfOperateGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_operate_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	groupId := d.Id()
	if v, ok := d.GetOk("operate"); ok {
		operate := v.(string)
		if operate == "start" {
			request := tsf.NewStartGroupRequest()
			request.GroupId = &groupId
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().StartGroup(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update tsf startGroup failed, reason:%+v", logId, err)
				return err
			}
		}
		if operate == "stop" {
			request := tsf.NewStopGroupRequest()
			request.GroupId = &groupId
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().StopGroup(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update tsf stopGroup failed, reason:%+v", logId, err)
				return err
			}
		}

		service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			groupInfo, err := service.DescribeTsfStartGroupById(ctx, groupId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			if groupInfo == nil {
				err = fmt.Errorf("group %s not exists", groupId)
				return resource.NonRetryableError(err)
			}
			if operate == "start" && *groupInfo.GroupStatus == "Running" {
				return nil
			}
			if operate == "stop" && *groupInfo.GroupStatus == "Paused" {
				return nil
			}
			if operate == "start" && *groupInfo.GroupStatus == "Paused" {
				return resource.RetryableError(fmt.Errorf("start or stop operation status is %s", *groupInfo.GroupStatus))
			}
			if operate == "stop" && *groupInfo.GroupStatus == "Running" {
				return resource.RetryableError(fmt.Errorf("start or stop operation status is %s", *groupInfo.GroupStatus))
			}
			if *groupInfo.GroupStatus == "Waiting" || *groupInfo.GroupStatus == "Updating" {
				return resource.RetryableError(fmt.Errorf("start or stop operation status is %s", *groupInfo.GroupStatus))
			}
			err = fmt.Errorf("start or stop operation status is %v, we won't wait for it finish", *groupInfo.GroupStatus)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s start or stop operation, reason:%s\n ", logId, err.Error())
			return err
		}

	}

	return resourceTencentCloudTsfOperateGroupRead(d, meta)
}

func resourceTencentCloudTsfOperateGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_operate_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

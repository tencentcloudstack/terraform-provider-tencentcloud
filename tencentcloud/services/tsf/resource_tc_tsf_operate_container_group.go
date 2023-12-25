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

func ResourceTencentCloudTsfOperateContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfOperateContainerGroupCreate,
		Read:   resourceTencentCloudTsfOperateContainerGroupRead,
		Update: resourceTencentCloudTsfOperateContainerGroupUpdate,
		Delete: resourceTencentCloudTsfOperateContainerGroupDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "group Id.",
			},

			"operate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Operation, `start`- start the container, `stop`- stop the container.",
			},
		},
	}
}

func resourceTencentCloudTsfOperateContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_operate_container_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var groupId string
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}

	d.SetId(groupId)

	return resourceTencentCloudTsfOperateContainerGroupUpdate(d, meta)
}

func resourceTencentCloudTsfOperateContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_operate_container_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	groupId := d.Id()
	startContainerGroup, err := service.DescribeTsfStartContainerGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if startContainerGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfOperateContainerGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("group_id", groupId)

	return nil
}

func resourceTencentCloudTsfOperateContainerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_operate_container_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	groupId := d.Id()

	if v, ok := d.GetOk("operate"); ok {
		var status bool
		operate := v.(string)
		if operate == "start" {
			request := tsf.NewStartContainerGroupRequest()
			request.GroupId = &groupId
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().StartContainerGroup(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				status = *result.Response.Result
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update tsf startContainerGroup failed, reason:%+v", logId, err)
				return err
			}

			if !status {
				return fmt.Errorf("[CRITAL]%s start tsf containerGroup failed", logId)
			}
		} else if operate == "stop" {
			request := tsf.NewStopContainerGroupRequest()
			request.GroupId = &groupId
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().StopContainerGroup(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				status = *result.Response.Result
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update tsf stopContainerGroup failed, reason:%+v", logId, err)
				return err
			}

			if !status {
				return fmt.Errorf("[CRITAL]%s stop tsf containerGroup failed", logId)
			}
		} else {
			return fmt.Errorf("[CRITAL]%s operate type error, %s", logId, operate)
		}

		service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			groupInfo, err := service.DescribeTsfStartContainerGroupById(ctx, groupId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			if groupInfo == nil {
				err = fmt.Errorf("group %s not exists", groupId)
				return resource.NonRetryableError(err)
			}
			if operate == "start" && *groupInfo.Status == "Running" {
				return nil
			}
			if operate == "stop" && *groupInfo.Status == "Paused" {
				return nil
			}
			if operate == "start" && *groupInfo.Status == "Paused" {
				return resource.RetryableError(fmt.Errorf("start or stop operation status is %s", *groupInfo.Status))
			}
			if operate == "stop" && *groupInfo.Status == "Running" {
				return resource.RetryableError(fmt.Errorf("start or stop operation status is %s", *groupInfo.Status))
			}
			if *groupInfo.Status == "Waiting" || *groupInfo.Status == "Updating" {
				return resource.RetryableError(fmt.Errorf("start or stop operation status is %s", *groupInfo.Status))
			}
			err = fmt.Errorf("start or stop operation status is %v, we won't wait for it finish", *groupInfo.Status)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s start or stop operation, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudTsfOperateContainerGroupRead(d, meta)
}

func resourceTencentCloudTsfOperateContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_Operate_container_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

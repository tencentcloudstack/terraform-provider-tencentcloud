/*
Provides a resource to create a tsf operate_group

Example Usage

```hcl
resource "tencentcloud_tsf_operate_group" "operate_group" {
  group_id = "group-ynd95rea"
  operate  = "start"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
)

func resourceTencentCloudTsfOperateGroup() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_tsf_operate_group.create")()
	defer inconsistentCheck(d, meta)()

	var groupId string
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}

	d.SetId(groupId)

	return resourceTencentCloudTsfOperateGroupUpdate(d, meta)
}

func resourceTencentCloudTsfOperateGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_operate_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

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
	defer logElapsed("resource.tencentcloud_tsf_operate_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	groupId := d.Id()
	if v, ok := d.GetOk("operate"); ok {
		operate := v.(string)
		if operate == "start" {
			request := tsf.NewStartGroupRequest()
			request.GroupId = &groupId
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().StartGroup(request)
				if e != nil {
					return retryError(e)
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
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().StopGroup(request)
				if e != nil {
					return retryError(e)
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

		service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			groupInfo, err := service.DescribeTsfStartGroupById(ctx, groupId)
			if err != nil {
				return retryError(err)
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
	defer logElapsed("resource.tencentcloud_tsf_operate_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

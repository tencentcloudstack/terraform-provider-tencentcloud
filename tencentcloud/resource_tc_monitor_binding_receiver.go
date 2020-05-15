/*
Provides a resource for bind receivers to a policy group resource.

Example Usage

```hcl
data "tencentcloud_cam_groups" "groups" {
  //You should first create a user group with CAM
}

resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "nice_group"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
}

resource "tencentcloud_monitor_binding_receiver" "receiver" {
  group_id = tencentcloud_monitor_policy_group.group.id
  receivers {
    start_time          = 0
    end_time            = 86399
    notify_way          = ["SMS"]
    receiver_type       = "group"
    receiver_group_list = [data.tencentcloud_cam_groups.groups.group_list[0].group_id]
    receive_language    = "en-US"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"time"
)

func resourceTencentMonitorBindingAlarmReceiver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentMonitorBindingAlarmReceiverCreate,
		Read:   resourceTencentMonitorBindingAlarmReceiverRead,
		Update: resourceTencentMonitorBindingAlarmReceiverUpdate,
		Delete: resourceTencentMonitorBindingAlarmReceiverDelete,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Policy group id for binding receivers.",
			},
			"receivers": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A list of receivers(will overwrite the configuration of the server or other resources). Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validateIntegerInRange(0, 86399),
							Description:  "Alarm period start time.Range [0,86399], which removes the date after it is converted to Beijing time as a Unix timestamp, for example 7200 means '10:0:0'.",
						},
						"end_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      86400,
							ValidateFunc: validateIntegerInRange(0, 86399),
							Description:  "End of alarm period. Meaning with `start_time`.",
						},
						"notify_way": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							MinItems:    1,
							Description: `Method of warning notification.Optional ` + helper.SliceFieldSerialize(monitorNotifyWays) + `.`,
						},
						"receiver_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue(monitorReceiverTypes),
							Description:  "Receive type. Optional " + helper.SliceFieldSerialize(monitorReceiverTypes) + ".",
						},
						"receiver_group_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
							Description: "Alarm receive group id list.",
						},
						"receiver_user_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
							Description: "Alarm receiver id list.",
						},
						"receive_language": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(monitorReceiveLanguages),
							Default:      monitorReceiveLanguageCN,
							Description:  "Alert sending language. Optional " + helper.SliceFieldSerialize(monitorReceiveLanguages) + ".",
						},
					},
				},
			},
		},
	}
}

func resourceTencentMonitorBindingAlarmReceiverCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_binding_receiver.create")()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), "logId", logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewModifyAlarmReceiversRequest()
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}

	if info == nil {
		return fmt.Errorf("policy group %d not exist", groupId)
	}

	request.GroupId = &groupId
	request.Module = helper.String("monitor")
	request.ReceiverInfos = make([]*monitor.ReceiverInfo, 0, 10)

	for _, iface := range d.Get("receivers").([]interface{}) {
		var receiverInfo monitor.ReceiverInfo
		ifaceMap := iface.(map[string]interface{})

		receiverInfo.StartTime = helper.IntInt64(ifaceMap["start_time"].(int))
		receiverInfo.EndTime = helper.IntInt64(ifaceMap["end_time"].(int))
		receiverInfo.NotifyWay = helper.InterfacesStringsPoint(ifaceMap["notify_way"].([]interface{}))
		receiverInfo.ReceiverType = helper.String(ifaceMap["receiver_type"].(string))

		if ifaceMap["receiver_group_list"] != nil {
			receiverInfo.ReceiverGroupList = helper.InterfacesIntInt64Point(ifaceMap["receiver_group_list"].([]interface{}))
		}
		if ifaceMap["receiver_user_list"] != nil {
			receiverInfo.ReceiverUserList = helper.InterfacesIntInt64Point(ifaceMap["receiver_user_list"].([]interface{}))
		}

		if *receiverInfo.ReceiverType == monitorReceiverTypeGroup {
			if len(receiverInfo.ReceiverGroupList) < 1 {
				return fmt.Errorf("miss field receiver_group_list, this array at least  has one element when you choose `group` receiver_type")
			}
		}
		if *receiverInfo.ReceiverType == monitorReceiverTypeUser {
			if len(receiverInfo.ReceiverUserList) < 1 {
				return fmt.Errorf("miss field receiver_user_list, this array at least  has one element when you choose `user` receiver_type")
			}
		}

		receiverInfo.ReceiveLanguage = helper.String(ifaceMap["receive_language"].(string))
		request.ReceiverInfos = append(request.ReceiverInfos, &receiverInfo)
	}

	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().ModifyAlarmReceivers(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", groupId))
	time.Sleep(3 * time.Second)

	return resourceTencentMonitorBindingAlarmReceiverRead(d, meta)
}

func resourceTencentMonitorBindingAlarmReceiverRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_binding_receiver.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), "logId", logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return nil
	}

	list := make([]interface{}, 0, len(info.ReceiverInfos))

	for _, receiver := range info.ReceiverInfos {
		var receiverMap = make(map[string]interface{})
		receiverMap["start_time"] = receiver.StartTime
		receiverMap["end_time"] = receiver.EndTime
		receiverMap["notify_way"] = receiver.NotifyWay
		receiverMap["receiver_type"] = receiver.ReceiverType
		receiverMap["receiver_group_list"] = receiver.ReceiverGroupList
		receiverMap["receiver_user_list"] = receiver.ReceiverUserList
		receiverMap["receive_language"] = receiver.ReceiveLanguage
		list = append(list, receiverMap)

	}
	return d.Set("receivers", list)

}
func resourceTencentMonitorBindingAlarmReceiverUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_binding_receiver.update")()
	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), "logId", logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewModifyAlarmReceiversRequest()
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return nil
	}

	request.GroupId = &groupId
	request.Module = helper.String("monitor")
	request.ReceiverInfos = make([]*monitor.ReceiverInfo, 0, 10)

	for _, iface := range d.Get("receivers").([]interface{}) {
		var receiverInfo monitor.ReceiverInfo
		ifaceMap := iface.(map[string]interface{})

		receiverInfo.StartTime = helper.IntInt64(ifaceMap["start_time"].(int))
		receiverInfo.EndTime = helper.IntInt64(ifaceMap["end_time"].(int))
		receiverInfo.NotifyWay = helper.InterfacesStringsPoint(ifaceMap["notify_way"].([]interface{}))
		receiverInfo.ReceiverType = helper.String(ifaceMap["receiver_type"].(string))

		if ifaceMap["receiver_group_list"] != nil {
			receiverInfo.ReceiverGroupList = helper.InterfacesIntInt64Point(ifaceMap["receiver_group_list"].([]interface{}))
		}
		if ifaceMap["receiver_user_list"] != nil {
			receiverInfo.ReceiverUserList = helper.InterfacesIntInt64Point(ifaceMap["receiver_user_list"].([]interface{}))
		}
		if *receiverInfo.ReceiverType == monitorReceiverTypeGroup {
			if len(receiverInfo.ReceiverGroupList) < 1 {
				return fmt.Errorf("miss field receiver_group_list, this array at least  has one element when you choose `group` receiver_type")
			}
		}
		if *receiverInfo.ReceiverType == monitorReceiverTypeUser {
			if len(receiverInfo.ReceiverUserList) < 1 {
				return fmt.Errorf("miss field receiver_user_list, this array at least  has one element when you choose `user` receiver_type")
			}
		}
		receiverInfo.ReceiveLanguage = helper.String(ifaceMap["receive_language"].(string))
		request.ReceiverInfos = append(request.ReceiverInfos, &receiverInfo)
	}

	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().ModifyAlarmReceivers(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	return resourceTencentMonitorBindingAlarmReceiverRead(d, meta)
}

func resourceTencentMonitorBindingAlarmReceiverDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_binding_receiver.delete")()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), "logId", logId)
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewModifyAlarmReceiversRequest()
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return nil
	}

	request.GroupId = &groupId
	request.Module = helper.String("monitor")
	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().ModifyAlarmReceivers(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

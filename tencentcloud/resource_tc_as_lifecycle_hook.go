package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
)

func resourceTencentCloudAsLifecycleHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsLifecycleHookCreate,
		Read:   resourceTencentCloudAsLifecycleHookRead,
		Update: resourceTencentCloudAsLifecycleHookUpdate,
		Delete: resourceTencentCloudAsLifecycleHookDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lifecycle_hook_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lifecycle_transition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"INSTANCE_LAUNCHING", "INSTANCE_TERMINATING"}),
			},
			"default_result": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "CONTINUE",
				ValidateFunc: validateAllowedStringValue([]string{"CONTINUE", "ABANDON"}),
			},
			"heartbeat_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: validateIntegerInRange(30, 3600),
			},
			"notification_metadata": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"notification_target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"CMQ_QUEUE", "CMQ_TOPIC"}),
			},
			"notification_queue_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_topic_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTencentCloudAsLifecycleHookCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)

	request := as.NewCreateLifecycleHookRequest()
	request.AutoScalingGroupId = stringToPointer(d.Get("scaling_group_id").(string))
	request.LifecycleHookName = stringToPointer(d.Get("lifecycle_hook_name").(string))
	request.LifecycleTransition = stringToPointer(d.Get("lifecycle_transition").(string))

	if v, ok := d.GetOk("default_result"); ok {
		request.DefaultResult = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("heartbeat_timeout"); ok {
		heartbeatTimeout := int64(v.(int))
		request.HeartbeatTimeout = &heartbeatTimeout
	}
	if v, ok := d.GetOk("notification_metadata"); ok {
		request.NotificationMetadata = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("notification_target_type"); ok {
		request.NotificationTarget = &as.NotificationTarget{}
		request.NotificationTarget.TargetType = stringToPointer(v.(string))
		if v.(string) == "CMQ_QUEUE" {
			if vv, ok := d.GetOk("notification_queue_name"); ok {
				request.NotificationTarget.QueueName = stringToPointer(vv.(string))
			} else {
				return fmt.Errorf("queue_name must not be null when target_type is CMQ_QUEUE")
			}
		} else if v.(string) == "CMQ_TOPIC" {
			if vv, ok := d.GetOk("notification_topic_name"); ok {
				request.NotificationTarget.TopicName = stringToPointer(vv.(string))
			} else {
				return fmt.Errorf("topic_name must ot be null when target_type is CMQ_TOPIC")
			}
		}
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().CreateLifecycleHook(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.LifecycleHookId == nil {
		return fmt.Errorf("lifecycle hook id is nil")
	}
	d.SetId(*response.Response.LifecycleHookId)

	return resourceTencentCloudAsLifecycleHookRead(d, meta)
}

func resourceTencentCloudAsLifecycleHookRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	lifecycleHookId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	lifecycleHook, err := asService.DescribeLifecycleHookById(ctx, lifecycleHookId)
	if err != nil {
		return err
	}

	d.Set("scaling_group_id", *lifecycleHook.AutoScalingGroupId)
	d.Set("lifecycle_hook_name", *lifecycleHook.LifecycleHookName)
	d.Set("lifecycle_transition", *lifecycleHook.LifecycleTransition)
	if lifecycleHook.DefaultResult != nil {
		d.Set("default_result", *lifecycleHook.DefaultResult)
	}
	if lifecycleHook.HeartbeatTimeout != nil {
		d.Set("heartbeat_timeout", *lifecycleHook.HeartbeatTimeout)
	}
	if lifecycleHook.NotificationMetadata != nil {
		d.Set("notification_metadata", *lifecycleHook.NotificationMetadata)
	}
	if lifecycleHook.NotificationTarget != nil {
		d.Set("notification_target_type", *lifecycleHook.NotificationTarget.TargetType)
		if lifecycleHook.NotificationTarget.QueueName != nil {
			d.Set("notification_queue_name", *lifecycleHook.NotificationTarget.QueueName)
		}
		if lifecycleHook.NotificationTarget.TopicName != nil {
			d.Set("notification_topic_name", *lifecycleHook.NotificationTarget.TopicName)
		}
	}

	return nil
}

func resourceTencentCloudAsLifecycleHookUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)

	request := as.NewUpgradeLifecycleHookRequest()
	lifecycleHookId := d.Id()
	request.LifecycleHookId = &lifecycleHookId
	request.LifecycleHookName = stringToPointer(d.Get("lifecycle_hook_name").(string))
	request.LifecycleTransition = stringToPointer(d.Get("lifecycle_transition").(string))
	if v, ok := d.GetOk("default_result"); ok {
		request.DefaultResult = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("heartbeat_timeout"); ok {
		heartbeatTimeout := int64(v.(int))
		request.HeartbeatTimeout = &heartbeatTimeout
	}
	if v, ok := d.GetOk("notification_metadata"); ok {
		request.NotificationMetadata = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("notification_target_type"); ok {
		request.NotificationTarget = &as.NotificationTarget{}
		request.NotificationTarget.TargetType = stringToPointer(v.(string))
		if v.(string) == "CMQ_QUEUE" {
			if vv, ok := d.GetOk("notification_queue_name"); ok {
				request.NotificationTarget.QueueName = stringToPointer(vv.(string))
			} else {
				return fmt.Errorf("queue_name must not be null when target_type is CMQ_QUEUE")
			}
		} else if v.(string) == "CMQ_TOPIC" {
			if vv, ok := d.GetOk("notification_topic_name"); ok {
				request.NotificationTarget.TopicName = stringToPointer(vv.(string))
			} else {
				return fmt.Errorf("topic_name must ot be null when target_type is CMQ_TOPIC")
			}
		}
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().UpgradeLifecycleHook(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func resourceTencentCloudAsLifecycleHookDelete(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	lifecycleHookId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := asService.DeleteLifecycleHook(ctx, lifecycleHookId)
	if err != nil {
		return err
	}

	return nil
}

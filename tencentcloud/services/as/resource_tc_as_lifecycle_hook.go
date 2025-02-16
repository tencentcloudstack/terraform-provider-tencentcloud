package as

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsLifecycleHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsLifecycleHookCreate,
		Read:   resourceTencentCloudAsLifecycleHookRead,
		Update: resourceTencentCloudAsLifecycleHookUpdate,
		Delete: resourceTencentCloudAsLifecycleHookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a scaling group.",
			},
			"lifecycle_hook_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the lifecycle hook.",
			},
			"lifecycle_transition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"INSTANCE_LAUNCHING", "INSTANCE_TERMINATING"}),
				Description:  "The instance state to which you want to attach the lifecycle hook. Valid values: `INSTANCE_LAUNCHING` and `INSTANCE_TERMINATING`.",
			},
			"default_result": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "CONTINUE",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"CONTINUE", "ABANDON"}),
				Description:  "Defines the action the AS group should take when the lifecycle hook timeout elapses or if an unexpected failure occurs. Valid values: `CONTINUE` and `ABANDON`. The default value is `CONTINUE`.",
			},
			"heartbeat_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: tccommon.ValidateIntegerInRange(30, 7200),
				Description:  "Defines the amount of time, in seconds, that can elapse before the lifecycle hook times out. Valid value ranges: (30~7200). and default value is `300`.",
			},
			"notification_metadata": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Contains additional information that you want to include any time AS sends a message to the notification target.",
			},
			"notification_target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"CMQ_QUEUE", "CMQ_TOPIC", "TDMQ_CMQ_QUEUE", "TDMQ_CMQ_TOPIC"}),
				Description:  "Target type. Valid values: `CMQ_QUEUE`, `CMQ_TOPIC`, `TDMQ_CMQ_QUEUE`, `TDMQ_CMQ_TOPIC`.",
			},
			"notification_queue_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "For CMQ_QUEUE type, a name of queue must be set.",
			},
			"notification_topic_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "For CMQ_TOPIC type, a name of topic must be set.",
			},
			"lifecycle_transition_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The scenario where the lifecycle hook is applied. `EXTENSION`: the lifecycle hook will be triggered when AttachInstances, DetachInstances or RemoveInstaces is called. `NORMAL`: the lifecycle hook is not triggered by the above APIs.",
			},
			"lifecycle_command": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Optional:    true,
				Description: "Remote command execution object. `NotificationTarget` and `LifecycleCommand` cannot be specified at the same time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Remote command ID. It is required to execute a command.",
						},
						"parameters": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom parameter. The field type is JSON encoded string. For example, {\"varA\": \"222\"}.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAsLifecycleHookCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_lifecycle_hook.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewCreateLifecycleHookRequest()
	request.AutoScalingGroupId = helper.String(d.Get("scaling_group_id").(string))
	request.LifecycleHookName = helper.String(d.Get("lifecycle_hook_name").(string))
	request.LifecycleTransition = helper.String(d.Get("lifecycle_transition").(string))

	if v, ok := d.GetOk("default_result"); ok {
		request.DefaultResult = helper.String(v.(string))
	}
	if v, ok := d.GetOk("heartbeat_timeout"); ok {
		heartbeatTimeout := int64(v.(int))
		request.HeartbeatTimeout = &heartbeatTimeout
	}
	if v, ok := d.GetOk("notification_metadata"); ok {
		request.NotificationMetadata = helper.String(v.(string))
	}
	if v, ok := d.GetOk("lifecycle_transition_type"); ok {
		request.LifecycleTransitionType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("notification_target_type"); ok {
		request.NotificationTarget = &as.NotificationTarget{}
		request.NotificationTarget.TargetType = helper.String(v.(string))
		if v.(string) == "CMQ_QUEUE" {
			if vv, ok := d.GetOk("notification_queue_name"); ok {
				request.NotificationTarget.QueueName = helper.String(vv.(string))
			} else {
				return fmt.Errorf("queue_name must not be null when target_type is CMQ_QUEUE")
			}
		} else if v.(string) == "CMQ_TOPIC" {
			if vv, ok := d.GetOk("notification_topic_name"); ok {
				request.NotificationTarget.TopicName = helper.String(vv.(string))
			} else {
				return fmt.Errorf("topic_name must ot be null when target_type is CMQ_TOPIC")
			}
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "lifecycle_command"); ok {
		lifecycleCommand := as.LifecycleCommand{}
		if v, ok := dMap["command_id"]; ok {
			lifecycleCommand.CommandId = helper.String(v.(string))
		}
		if v, ok := dMap["parameters"]; ok {
			lifecycleCommand.Parameters = helper.String(v.(string))
		}
		request.LifecycleCommand = &lifecycleCommand
	}

	var lifecycleHookId string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateLifecycleHook(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
			if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if strings.Contains(e.GetCode(), "LimitExceeded.QuotaNotEnough") {
					return resource.RetryableError(err)
				}
			}

			return tccommon.RetryError(err)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if response == nil || response.Response == nil || response.Response.LifecycleHookId == nil {
			return resource.NonRetryableError(fmt.Errorf("AS LifecycleHook not exists"))
		}

		lifecycleHookId = *response.Response.LifecycleHookId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create AS LifecycleHook failed, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(lifecycleHookId)

	return resourceTencentCloudAsLifecycleHookRead(d, meta)
}

func resourceTencentCloudAsLifecycleHookRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_lifecycle_hook.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	lifecycleHookId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		lifecycleHook, has, e := asService.DescribeLifecycleHookById(ctx, lifecycleHookId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}
		_ = d.Set("scaling_group_id", *lifecycleHook.AutoScalingGroupId)
		_ = d.Set("lifecycle_hook_name", *lifecycleHook.LifecycleHookName)
		_ = d.Set("lifecycle_transition", *lifecycleHook.LifecycleTransition)
		if lifecycleHook.DefaultResult != nil {
			_ = d.Set("default_result", *lifecycleHook.DefaultResult)
		}
		if lifecycleHook.HeartbeatTimeout != nil {
			_ = d.Set("heartbeat_timeout", *lifecycleHook.HeartbeatTimeout)
		}
		if lifecycleHook.NotificationMetadata != nil {
			_ = d.Set("notification_metadata", *lifecycleHook.NotificationMetadata)
		}
		if lifecycleHook.LifecycleTransitionType != nil {
			_ = d.Set("lifecycle_transition_type", *lifecycleHook.LifecycleTransitionType)
		}
		if lifecycleHook.NotificationTarget != nil {
			_ = d.Set("notification_target_type", *lifecycleHook.NotificationTarget.TargetType)
			if lifecycleHook.NotificationTarget.QueueName != nil {
				_ = d.Set("notification_queue_name", *lifecycleHook.NotificationTarget.QueueName)
			}
			if lifecycleHook.NotificationTarget.TopicName != nil {
				_ = d.Set("notification_topic_name", *lifecycleHook.NotificationTarget.TopicName)
			}
		}
		if lifecycleHook.LifecycleCommand != nil {
			commandMap := map[string]interface{}{}
			if lifecycleHook.LifecycleCommand.CommandId != nil {
				commandMap["command_id"] = lifecycleHook.LifecycleCommand.CommandId
			}
			if lifecycleHook.LifecycleCommand.Parameters != nil {
				commandMap["parameters"] = lifecycleHook.LifecycleCommand.Parameters
			}
			_ = d.Set("lifecycle_command", []interface{}{commandMap})
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAsLifecycleHookUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_lifecycle_hook.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewUpgradeLifecycleHookRequest()
	lifecycleHookId := d.Id()
	request.LifecycleHookId = &lifecycleHookId
	request.LifecycleHookName = helper.String(d.Get("lifecycle_hook_name").(string))
	request.LifecycleTransition = helper.String(d.Get("lifecycle_transition").(string))
	if v, ok := d.GetOk("default_result"); ok {
		request.DefaultResult = helper.String(v.(string))
	}
	if v, ok := d.GetOk("heartbeat_timeout"); ok {
		heartbeatTimeout := int64(v.(int))
		request.HeartbeatTimeout = &heartbeatTimeout
	}
	if v, ok := d.GetOk("notification_metadata"); ok {
		request.NotificationMetadata = helper.String(v.(string))
	}
	if v, ok := d.GetOk("lifecycle_transition_type"); ok {
		request.LifecycleTransitionType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("notification_target_type"); ok {
		request.NotificationTarget = &as.NotificationTarget{}
		request.NotificationTarget.TargetType = helper.String(v.(string))
		if v.(string) == "CMQ_QUEUE" {
			if vv, ok := d.GetOk("notification_queue_name"); ok {
				request.NotificationTarget.QueueName = helper.String(vv.(string))
			} else {
				return fmt.Errorf("queue_name must not be null when target_type is CMQ_QUEUE")
			}
		} else if v.(string) == "CMQ_TOPIC" {
			if vv, ok := d.GetOk("notification_topic_name"); ok {
				request.NotificationTarget.TopicName = helper.String(vv.(string))
			} else {
				return fmt.Errorf("topic_name must ot be null when target_type is CMQ_TOPIC")
			}
		}
	}

	if d.HasChange("lifecycle_command") {
		if dMap, ok := helper.InterfacesHeadMap(d, "lifecycle_command"); ok {
			lifecycleCommand := as.LifecycleCommand{}
			if v, ok := dMap["command_id"]; ok && v != "" {
				lifecycleCommand.CommandId = helper.String(v.(string))
			}
			if v, ok := dMap["parameters"]; ok && v != "" {
				lifecycleCommand.Parameters = helper.String(v.(string))
			}
			request.LifecycleCommand = &lifecycleCommand
		}
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().UpgradeLifecycleHook(request)
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
	defer tccommon.LogElapsed("resource.tencentcloud_as_lifecycle_hook.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	lifecycleHookId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := asService.DeleteLifecycleHook(ctx, lifecycleHookId)
	if err != nil {
		return err
	}

	return nil
}

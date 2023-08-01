/*
Provides a resource for an AS (Auto scaling) lifecycle hook.

Example Usage

Create a basic LifecycleHook

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.image.images.0.image_id
  instance_types     = ["SA1.SMALL1", "SA2.SMALL1", "SA2.SMALL2", "SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
  }
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name = "tf-example"
  configuration_id   = tencentcloud_as_scaling_config.example.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_as_lifecycle_hook" "example" {
  scaling_group_id         = tencentcloud_as_scaling_group.example.id
  lifecycle_hook_name      = "tf-as-lifecycle-hook"
  lifecycle_transition     = "INSTANCE_LAUNCHING"
  default_result           = "CONTINUE"
  heartbeat_timeout        = 500
  notification_metadata    = "tf test"
}
```

If `notification_target_type` is `CMQ_QUEUE`

```hcl
resource "tencentcloud_as_lifecycle_hook" "example" {
  scaling_group_id         = tencentcloud_as_scaling_group.example.id
  lifecycle_hook_name      = "tf-as-lifecycle-hook"
  lifecycle_transition     = "INSTANCE_LAUNCHING"
  default_result           = "CONTINUE"
  heartbeat_timeout        = 500
  notification_metadata    = "tf test"
  notification_target_type = "CMQ_QUEUE"
  notification_queue_name  = "lifcyclehook"
}
```

Or `notification_target_type` is `CMQ_TOPIC`

```hcl
resource "tencentcloud_as_lifecycle_hook" "example" {
  scaling_group_id         = tencentcloud_as_scaling_group.example.id
  lifecycle_hook_name      = "tf-as-lifecycle-hook"
  lifecycle_transition     = "INSTANCE_LAUNCHING"
  default_result           = "CONTINUE"
  heartbeat_timeout        = 500
  notification_metadata    = "tf test"
  notification_target_type = "CMQ_TOPIC"
  notification_topic_name  = "lifcyclehook"
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
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAsLifecycleHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsLifecycleHookCreate,
		Read:   resourceTencentCloudAsLifecycleHookRead,
		Update: resourceTencentCloudAsLifecycleHookUpdate,
		Delete: resourceTencentCloudAsLifecycleHookDelete,

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
				ValidateFunc: validateAllowedStringValue([]string{"INSTANCE_LAUNCHING", "INSTANCE_TERMINATING"}),
				Description:  "The instance state to which you want to attach the lifecycle hook. Valid values: `INSTANCE_LAUNCHING` and `INSTANCE_TERMINATING`.",
			},
			"default_result": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "CONTINUE",
				ValidateFunc: validateAllowedStringValue([]string{"CONTINUE", "ABANDON"}),
				Description:  "Defines the action the AS group should take when the lifecycle hook timeout elapses or if an unexpected failure occurs. Valid values: `CONTINUE` and `ABANDON`. The default value is `CONTINUE`.",
			},
			"heartbeat_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: validateIntegerInRange(30, 7200),
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
				ValidateFunc: validateAllowedStringValue([]string{"CMQ_QUEUE", "CMQ_TOPIC"}),
				Description:  "Target type. Valid values: `CMQ_QUEUE`, `CMQ_TOPIC`.",
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
		},
	}
}

func resourceTencentCloudAsLifecycleHookCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_lifecycle_hook.create")()

	logId := getLogId(contextNil)

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
	defer logElapsed("resource.tencentcloud_as_lifecycle_hook.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	lifecycleHookId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		lifecycleHook, has, e := asService.DescribeLifecycleHookById(ctx, lifecycleHookId)
		if e != nil {
			return retryError(e)
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
		if lifecycleHook.NotificationTarget != nil {
			_ = d.Set("notification_target_type", *lifecycleHook.NotificationTarget.TargetType)
			if lifecycleHook.NotificationTarget.QueueName != nil {
				_ = d.Set("notification_queue_name", *lifecycleHook.NotificationTarget.QueueName)
			}
			if lifecycleHook.NotificationTarget.TopicName != nil {
				_ = d.Set("notification_topic_name", *lifecycleHook.NotificationTarget.TopicName)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAsLifecycleHookUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_lifecycle_hook.update")()

	logId := getLogId(contextNil)

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
	defer logElapsed("resource.tencentcloud_as_lifecycle_hook.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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

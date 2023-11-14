/*
Provides a resource to create a tdmq subscription_attachment

Example Usage

```hcl
resource "tencentcloud_tdmq_subscription_attachment" "subscription_attachment" {
  environment_id = ""
  topic_name = ""
  subscription_name = ""
  is_idempotent =
  remark = ""
  cluster_id = ""
  auto_create_policy_topic =
  post_fix_pattern = ""
}
```

Import

tdmq subscription_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_subscription_attachment.subscription_attachment subscription_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTdmqSubscriptionAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqSubscriptionAttachmentCreate,
		Read:   resourceTencentCloudTdmqSubscriptionAttachmentRead,
		Delete: resourceTencentCloudTdmqSubscriptionAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Environment (namespace) name.",
			},

			"topic_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Topic name.",
			},

			"subscription_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Subscriber name, no more than 128 characters.",
			},

			"is_idempotent": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether it is idempotent to create, if not, it is not allowed to create a subscription relationship with the same name.",
			},

			"remark": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Remarks, within 128 characters.",
			},

			"cluster_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the Pulsar cluster.",
			},

			"auto_create_policy_topic": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically create dead letters and retry topics, True means to create, False means not to create, the default is to automatically create dead letters and retry topics.",
			},

			"post_fix_pattern": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Specifies the dead letter and retry topic name specification, LEGACY indicates the historical naming convention, COMMUNITY indicates the Pulsar community naming convention.",
			},
		},
	}
}

func resourceTencentCloudTdmqSubscriptionAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscription_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateSubscriptionRequest()
		response  = tdmq.NewCreateSubscriptionResponse()
		clusterId string
	)
	if v, ok := d.GetOk("environment_id"); ok {
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subscription_name"); ok {
		request.SubscriptionName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_idempotent"); ok {
		request.IsIdempotent = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_create_policy_topic"); ok {
		request.AutoCreatePolicyTopic = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("post_fix_pattern"); ok {
		request.PostFixPattern = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateSubscription(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmq subscriptionAttachment failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	return resourceTencentCloudTdmqSubscriptionAttachmentRead(d, meta)
}

func resourceTencentCloudTdmqSubscriptionAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscription_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	subscriptionAttachmentId := d.Id()

	subscriptionAttachment, err := service.DescribeTdmqSubscriptionAttachmentById(ctx, clusterId)
	if err != nil {
		return err
	}

	if subscriptionAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqSubscriptionAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if subscriptionAttachment.EnvironmentId != nil {
		_ = d.Set("environment_id", subscriptionAttachment.EnvironmentId)
	}

	if subscriptionAttachment.TopicName != nil {
		_ = d.Set("topic_name", subscriptionAttachment.TopicName)
	}

	if subscriptionAttachment.SubscriptionName != nil {
		_ = d.Set("subscription_name", subscriptionAttachment.SubscriptionName)
	}

	if subscriptionAttachment.IsIdempotent != nil {
		_ = d.Set("is_idempotent", subscriptionAttachment.IsIdempotent)
	}

	if subscriptionAttachment.Remark != nil {
		_ = d.Set("remark", subscriptionAttachment.Remark)
	}

	if subscriptionAttachment.ClusterId != nil {
		_ = d.Set("cluster_id", subscriptionAttachment.ClusterId)
	}

	if subscriptionAttachment.AutoCreatePolicyTopic != nil {
		_ = d.Set("auto_create_policy_topic", subscriptionAttachment.AutoCreatePolicyTopic)
	}

	if subscriptionAttachment.PostFixPattern != nil {
		_ = d.Set("post_fix_pattern", subscriptionAttachment.PostFixPattern)
	}

	return nil
}

func resourceTencentCloudTdmqSubscriptionAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscription_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	subscriptionAttachmentId := d.Id()

	if err := service.DeleteTdmqSubscriptionAttachmentById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}

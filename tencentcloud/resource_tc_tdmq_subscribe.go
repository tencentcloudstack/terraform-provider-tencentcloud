/*
Provides a resource to create a tdmq subscribe

Example Usage

```hcl
resource "tencentcloud_tdmq_subscribe" "subscribe" {
  topic_name = "topic_name"
  subscription_name = "subscription_name"
  protocol = "HTTP"
  endpoint = &lt;nil&gt;
  notify_strategy = "EXPONENTIAL_DECAY_RETRY"
  filter_tag = &lt;nil&gt;
  binding_key = &lt;nil&gt;
  notify_content_format = "JSON"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tdmq subscribe can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_subscribe.subscribe subscribe_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTdmqSubscribe() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqSubscribeCreate,
		Read:   resourceTencentCloudTdmqSubscribeRead,
		Update: resourceTencentCloudTdmqSubscribeUpdate,
		Delete: resourceTencentCloudTdmqSubscribeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.",
			},

			"subscription_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subscription name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.",
			},

			"protocol": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Ubscription protocol. Currently, two protocols are supported: HTTP and queue. To use the HTTP protocol, you need to build your own web server to receive messages. With the queue protocol, messages are automatically pushed to a CMQ queue and you can pull them concurrently.",
			},

			"endpoint": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`Endpoint` for notification receipt, which is distinguished by `Protocol`. For `http`, `Endpoint` must begin with `http://` and `host` can be a domain name or IP. For `Queue`, enter `QueueName`. Note that currently the push service cannot push messages to a VPC; therefore, if a VPC domain name or address is entered for `Endpoint`, pushed messages will not be received. Currently, messages can be pushed only to the public network and classic network.",
			},

			"notify_strategy": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CMQ push server retry policy in case an error occurs while pushing a message to `Endpoint`. Valid values: 1. `BACKOFF_RETRY`: backoff retry, which is to retry at a fixed interval, discard the message after a certain number of retries, and continue to push the next message; 2. `EXPONENTIAL_DECAY_RETRY`: exponential decay retry, which is to retry at an exponentially increasing interval, such as 1s, 2s, 4s, 8s, and so on. As a message can be retained in a topic for one day, failed messages will be discarded at most after one day of retry. Default value: `EXPONENTIAL_DECAY_RETRY`.",
			},

			"filter_tag": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Message body tag (used for message filtering). The number of tags cannot exceed 5, and each tag can contain up to 16 characters. It is used in conjunction with the `MsgTag` parameter of `(Batch)PublishMessage`. Rules: 1. If `FilterTag` is not configured, no matter whether `MsgTag` is configured, the subscription will receive all messages published to the topic; 2. If the array of `FilterTag` values has a value, only when at least one of the values in the array also exists in the array of `MsgTag` values (i.e., `FilterTag` and `MsgTag` have an intersection) can the subscription receive messages published to the topic; 3. If the array of `FilterTag` values has a value, but `MsgTag` is not configured, then no message published to the topic will be received, which can be considered as a special case of rule 2 as `FilterTag` and `MsgTag` do not intersect in this case. The overall design idea of rules is based on the intention of the subscriber.",
			},

			"binding_key": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The number of `BindingKey` cannot exceed 5, and the length of each `BindingKey` cannot exceed 64 bytes. This field indicates the filtering policy for subscribing to and receiving messages. Each `BindingKey` includes up to 15 dots (namely up to 16 segments).",
			},

			"notify_content_format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Push content format. Valid values: 1. JSON; 2. SIMPLIFIED, i.e., the raw format. If `Protocol` is `queue`, this value must be `SIMPLIFIED`. If `Protocol` is `http`, both options are acceptable, and the default value is `JSON`.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTdmqSubscribeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscribe.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = tdmq.NewCreateCmqSubscribeRequest()
		response         = tdmq.NewCreateCmqSubscribeResponse()
		topicName        string
		subscriptionName string
	)
	if v, ok := d.GetOk("topic_name"); ok {
		topicName = v.(string)
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subscription_name"); ok {
		subscriptionName = v.(string)
		request.SubscriptionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("endpoint"); ok {
		request.Endpoint = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notify_strategy"); ok {
		request.NotifyStrategy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_tag"); ok {
		filterTagSet := v.(*schema.Set).List()
		for i := range filterTagSet {
			filterTag := filterTagSet[i].(string)
			request.FilterTag = append(request.FilterTag, &filterTag)
		}
	}

	if v, ok := d.GetOk("binding_key"); ok {
		bindingKeySet := v.(*schema.Set).List()
		for i := range bindingKeySet {
			bindingKey := bindingKeySet[i].(string)
			request.BindingKey = append(request.BindingKey, &bindingKey)
		}
	}

	if v, ok := d.GetOk("notify_content_format"); ok {
		request.NotifyContentFormat = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateCmqSubscribe(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmq subscribe failed, reason:%+v", logId, err)
		return err
	}

	topicName = *response.Response.TopicName
	d.SetId(strings.Join([]string{topicName, subscriptionName}, FILED_SP))

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tdmq:%s:uin/:topic/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTdmqSubscribeRead(d, meta)
}

func resourceTencentCloudTdmqSubscribeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscribe.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicName := idSplit[0]
	subscriptionName := idSplit[1]

	subscribe, err := service.DescribeTdmqSubscribeById(ctx, topicName, subscriptionName)
	if err != nil {
		return err
	}

	if subscribe == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqSubscribe` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if subscribe.TopicName != nil {
		_ = d.Set("topic_name", subscribe.TopicName)
	}

	if subscribe.SubscriptionName != nil {
		_ = d.Set("subscription_name", subscribe.SubscriptionName)
	}

	if subscribe.Protocol != nil {
		_ = d.Set("protocol", subscribe.Protocol)
	}

	if subscribe.Endpoint != nil {
		_ = d.Set("endpoint", subscribe.Endpoint)
	}

	if subscribe.NotifyStrategy != nil {
		_ = d.Set("notify_strategy", subscribe.NotifyStrategy)
	}

	if subscribe.FilterTag != nil {
		_ = d.Set("filter_tag", subscribe.FilterTag)
	}

	if subscribe.BindingKey != nil {
		_ = d.Set("binding_key", subscribe.BindingKey)
	}

	if subscribe.NotifyContentFormat != nil {
		_ = d.Set("notify_content_format", subscribe.NotifyContentFormat)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tdmq", "topic", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTdmqSubscribeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscribe.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyCmqSubscriptionAttributeRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicName := idSplit[0]
	subscriptionName := idSplit[1]

	request.TopicName = &topicName
	request.SubscriptionName = &subscriptionName

	immutableArgs := []string{"topic_name", "subscription_name", "protocol", "endpoint", "notify_strategy", "filter_tag", "binding_key", "notify_content_format"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("topic_name") {
		if v, ok := d.GetOk("topic_name"); ok {
			request.TopicName = helper.String(v.(string))
		}
	}

	if d.HasChange("subscription_name") {
		if v, ok := d.GetOk("subscription_name"); ok {
			request.SubscriptionName = helper.String(v.(string))
		}
	}

	if d.HasChange("notify_strategy") {
		if v, ok := d.GetOk("notify_strategy"); ok {
			request.NotifyStrategy = helper.String(v.(string))
		}
	}

	if d.HasChange("binding_key") {
		if v, ok := d.GetOk("binding_key"); ok {
			bindingKeySet := v.(*schema.Set).List()
			for i := range bindingKeySet {
				bindingKey := bindingKeySet[i].(string)
				request.BindingKey = append(request.BindingKey, &bindingKey)
			}
		}
	}

	if d.HasChange("notify_content_format") {
		if v, ok := d.GetOk("notify_content_format"); ok {
			request.NotifyContentFormat = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyCmqSubscriptionAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmq subscribe failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tdmq", "topic", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTdmqSubscribeRead(d, meta)
}

func resourceTencentCloudTdmqSubscribeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscribe.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicName := idSplit[0]
	subscriptionName := idSplit[1]

	if err := service.DeleteTdmqSubscribeById(ctx, topicName, subscriptionName); err != nil {
		return err
	}

	return nil
}

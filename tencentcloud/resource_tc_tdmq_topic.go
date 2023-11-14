/*
Provides a resource to create a tdmq topic

Example Usage

```hcl
resource "tencentcloud_tdmq_topic" "topic" {
  topic_name = "topic_name"
  max_msg_size = 65536
  filter_type = 1
  msg_retention_seconds = 86400
  trace = true
}
```

Import

tdmq topic can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_topic.topic topic_id
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
)

func resourceTencentCloudTdmqTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqTopicCreate,
		Read:   resourceTencentCloudTdmqTopicRead,
		Update: resourceTencentCloudTdmqTopicUpdate,
		Delete: resourceTencentCloudTdmqTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.",
			},

			"max_msg_size": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum message length. Value range: 1024-65536 bytes (i.e., 1-64 KB). Default value: 65536.",
			},

			"filter_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Used to specify the message match policy for the topic. `1`: tag match policy (default value); `2`: routing match policy.",
			},

			"msg_retention_seconds": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Message retention period. Value range: 60-86400 seconds (i.e., 1 minute-1 day). Default value: 86400.",
			},

			"trace": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable message trace. true: yes; false: no. If this field is left empty, the feature will not be enabled.",
			},
		},
	}
}

func resourceTencentCloudTdmqTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_topic.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateCmqTopicRequest()
		response  = tdmq.NewCreateCmqTopicResponse()
		topicName string
	)
	if v, ok := d.GetOk("topic_name"); ok {
		topicName = v.(string)
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_msg_size"); ok {
		request.MaxMsgSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("filter_type"); ok {
		request.FilterType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("msg_retention_seconds"); ok {
		request.MsgRetentionSeconds = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("trace"); ok {
		request.Trace = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateCmqTopic(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmq topic failed, reason:%+v", logId, err)
		return err
	}

	topicName = *response.Response.TopicName
	d.SetId(topicName)

	return resourceTencentCloudTdmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_topic.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	topicId := d.Id()

	topic, err := service.DescribeTdmqTopicById(ctx, topicName)
	if err != nil {
		return err
	}

	if topic == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqTopic` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if topic.TopicName != nil {
		_ = d.Set("topic_name", topic.TopicName)
	}

	if topic.MaxMsgSize != nil {
		_ = d.Set("max_msg_size", topic.MaxMsgSize)
	}

	if topic.FilterType != nil {
		_ = d.Set("filter_type", topic.FilterType)
	}

	if topic.MsgRetentionSeconds != nil {
		_ = d.Set("msg_retention_seconds", topic.MsgRetentionSeconds)
	}

	if topic.Trace != nil {
		_ = d.Set("trace", topic.Trace)
	}

	return nil
}

func resourceTencentCloudTdmqTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_topic.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyCmqTopicAttributeRequest()

	topicId := d.Id()

	request.TopicName = &topicName

	immutableArgs := []string{"topic_name", "max_msg_size", "filter_type", "msg_retention_seconds", "trace"}

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

	if d.HasChange("max_msg_size") {
		if v, ok := d.GetOkExists("max_msg_size"); ok {
			request.MaxMsgSize = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("msg_retention_seconds") {
		if v, ok := d.GetOkExists("msg_retention_seconds"); ok {
			request.MsgRetentionSeconds = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("trace") {
		if v, ok := d.GetOkExists("trace"); ok {
			request.Trace = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyCmqTopicAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmq topic failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_topic.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	topicId := d.Id()

	if err := service.DeleteTdmqTopicById(ctx, topicName); err != nil {
		return err
	}

	return nil
}

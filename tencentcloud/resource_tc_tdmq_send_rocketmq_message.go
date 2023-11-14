/*
Provides a resource to create a tdmq send_rocketmq_message

Example Usage

```hcl
resource "tencentcloud_tdmq_send_rocketmq_message" "send_rocketmq_message" {
  cluster_id = ""
  namespace_id = ""
  topic_name = ""
  msg_body = ""
  msg_key = ""
  msg_tag = ""
}
```

Import

tdmq send_rocketmq_message can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message send_rocketmq_message_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTdmqSendRocketmqMessage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqSendRocketmqMessageCreate,
		Read:   resourceTencentCloudTdmqSendRocketmqMessageRead,
		Delete: resourceTencentCloudTdmqSendRocketmqMessageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster id.",
			},

			"namespace_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespaces.",
			},

			"topic_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Topic name.",
			},

			"msg_body": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Information.",
			},

			"msg_key": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Message key information.",
			},

			"msg_tag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Message tag information.",
			},
		},
	}
}

func resourceTencentCloudTdmqSendRocketmqMessageCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_send_rocketmq_message.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = tdmq.NewSendRocketMQMessageRequest()
		response    = tdmq.NewSendRocketMQMessageResponse()
		clusterId   string
		namespaceId string
		topicName   string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		namespaceId = v.(string)
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		topicName = v.(string)
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("msg_body"); ok {
		request.MsgBody = helper.String(v.(string))
	}

	if v, ok := d.GetOk("msg_key"); ok {
		request.MsgKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("msg_tag"); ok {
		request.MsgTag = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().SendRocketMQMessage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tdmq sendRocketmqMessage failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(strings.Join([]string{clusterId, namespaceId, topicName}, FILED_SP))

	return resourceTencentCloudTdmqSendRocketmqMessageRead(d, meta)
}

func resourceTencentCloudTdmqSendRocketmqMessageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_send_rocketmq_message.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTdmqSendRocketmqMessageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_send_rocketmq_message.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

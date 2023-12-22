package trocket

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdmqSendRocketmqMessage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqSendRocketmqMessageCreate,
		Read:   resourceTencentCloudTdmqSendRocketmqMessageRead,
		Delete: resourceTencentCloudTdmqSendRocketmqMessageDelete,

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
				Description: "topic name.",
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
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_send_rocketmq_message.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = tdmq.NewSendRocketMQMessageRequest()
		clusterId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().SendRocketMQMessage(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("tdmq sendRocketmqMessage not exists")
			return resource.NonRetryableError(e)
		}

		if !*result.Response.Result {
			e = fmt.Errorf("send tdmq sendRocketmqMessage status is false, requestId: %s, MsgId: %s", *result.Response.RequestId, *result.Response.MsgId)
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate tdmq sendRocketmqMessage failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId)

	return resourceTencentCloudTdmqSendRocketmqMessageRead(d, meta)
}

func resourceTencentCloudTdmqSendRocketmqMessageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_send_rocketmq_message.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTdmqSendRocketmqMessageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_send_rocketmq_message.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

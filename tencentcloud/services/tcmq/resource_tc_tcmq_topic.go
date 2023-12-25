package tcmq

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcmqTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcmqTopicCreate,
		Read:   resourceTencentCloudTcmqTopicRead,
		Update: resourceTencentCloudTcmqTopicUpdate,
		Delete: resourceTencentCloudTcmqTopicDelete,
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
				Default:     65536,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum message length. Value range: 1024-65536 bytes (i.e., 1-64 KB). Default value: 65536.",
			},

			"filter_type": {
				Default:     1,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Used to specify the message match policy for the topic. `1`: tag match policy (default value); `2`: routing match policy.",
			},

			"msg_retention_seconds": {
				Default:     86400,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Message retention period. Value range: 60-86400 seconds (i.e., 1 minute-1 day). Default value: 86400.",
			},

			"trace": {
				Default:     true,
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable message trace. true: yes; false: no. If this field is left empty, the feature will not be enabled.",
			},
		},
	}
}

func resourceTencentCloudTcmqTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_topic.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = tcmq.NewCreateCmqTopicRequest()
		topicName string
	)
	if v, ok := d.GetOk("topic_name"); ok {
		topicName = v.(string)
		request.TopicName = helper.String(topicName)
	}

	if v, _ := d.GetOk("max_msg_size"); v != nil {
		request.MaxMsgSize = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("filter_type"); v != nil {
		request.FilterType = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("msg_retention_seconds"); v != nil {
		request.MsgRetentionSeconds = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("trace"); v != nil {
		request.Trace = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateCmqTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcmq topic failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(topicName)

	return resourceTencentCloudTcmqTopicRead(d, meta)
}

func resourceTencentCloudTcmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_topic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	topicName := d.Id()

	topic, err := service.DescribeTcmqTopicById(ctx, topicName)
	if err != nil {
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "ResourceNotFound" {
				return nil
			}
		}
		return err
	}

	if topic == nil {
		d.SetId("")
		return fmt.Errorf("resource %s does not exist", topicName)
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

func resourceTencentCloudTcmqTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_topic.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tcmq.NewModifyCmqTopicAttributeRequest()

	topicName := d.Id()

	request.TopicName = &topicName
	if d.HasChange("topic_name") {
		if v, ok := d.GetOk("topic_name"); ok {
			request.TopicName = helper.String(v.(string))
		}
	}

	if d.HasChange("max_msg_size") {
		if v, _ := d.GetOk("max_msg_size"); v != nil {
			request.MaxMsgSize = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("msg_retention_seconds") {
		if v, _ := d.GetOk("msg_retention_seconds"); v != nil {
			request.MsgRetentionSeconds = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("trace") {
		if v, _ := d.GetOk("trace"); v != nil {
			request.Trace = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyCmqTopicAttribute(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcmq topic failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmqTopicRead(d, meta)
}

func resourceTencentCloudTcmqTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_topic.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	topicName := d.Id()

	if err := service.DeleteTcmqTopicById(ctx, topicName); err != nil {
		return err
	}

	return nil
}

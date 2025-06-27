package trocket

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTrocketRocketmqTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTrocketRocketmqTopicCreate,
		Read:   resourceTencentCloudTrocketRocketmqTopicRead,
		Update: resourceTencentCloudTrocketRocketmqTopicUpdate,
		Delete: resourceTencentCloudTrocketRocketmqTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance Id.",
			},

			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "topic.",
			},

			"topic_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Topic type. `UNSPECIFIED`: not specified, `NORMAL`: normal message, `FIFO`: sequential message, `DELAY`: delayed message.",
			},

			"queue_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of queue. Must be greater than or equal to 3.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "remark.",
			},
		},
	}
}

func resourceTencentCloudTrocketRocketmqTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_topic.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = trocket.NewCreateTopicRequest()
		response   = trocket.NewCreateTopicResponse()
		instanceId string
		topic      string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic"); ok {
		request.Topic = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_type"); ok {
		request.TopicType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("queue_num"); ok {
		request.QueueNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTrocketClient().CreateTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create trocket rocketmqTopic failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create trocket rocketmqTopic failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.InstanceId == nil || response.Response.Topic == nil {
		return fmt.Errorf("InstanceId or Topic is nil.")
	}

	instanceId = *response.Response.InstanceId
	topic = *response.Response.Topic
	d.SetId(strings.Join([]string{instanceId, topic}, tccommon.FILED_SP))
	return resourceTencentCloudTrocketRocketmqTopicRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_topic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	topic := idSplit[1]

	rocketmqTopic, err := service.DescribeTrocketRocketmqTopicById(ctx, instanceId, topic)
	if err != nil {
		return err
	}

	if rocketmqTopic == nil {
		log.Printf("[WARN]%s resource `TrocketRocketmqTopic` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if rocketmqTopic.InstanceId != nil {
		_ = d.Set("instance_id", rocketmqTopic.InstanceId)
	}

	if rocketmqTopic.Topic != nil {
		_ = d.Set("topic", rocketmqTopic.Topic)
	}

	if rocketmqTopic.TopicType != nil {
		_ = d.Set("topic_type", rocketmqTopic.TopicType)
	}

	if rocketmqTopic.QueueNum != nil {
		_ = d.Set("queue_num", rocketmqTopic.QueueNum)
	}

	if rocketmqTopic.Remark != nil {
		_ = d.Set("remark", rocketmqTopic.Remark)
	}

	return nil
}

func resourceTencentCloudTrocketRocketmqTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_topic.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = trocket.NewModifyTopicRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	topic := idSplit[1]

	if d.HasChange("queue_num") || d.HasChange("remark") {
		request.InstanceId = &instanceId
		request.Topic = &topic

		if v, ok := d.GetOkExists("queue_num"); ok {
			request.QueueNum = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTrocketClient().ModifyTopic(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update trocket rocketmqTopic failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTrocketRocketmqTopicRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_topic.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	topic := idSplit[1]

	if err := service.DeleteTrocketRocketmqTopicById(ctx, instanceId, topic); err != nil {
		return err
	}

	return nil
}

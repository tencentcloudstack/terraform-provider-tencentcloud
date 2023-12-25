package trocket

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdmqRocketmqTopic() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRocketmqTopicRead,
		Create: resourceTencentCloudTdmqRocketmqTopicCreate,
		Update: resourceTencentCloudTdmqRocketmqTopicUpdate,
		Delete: resourceTencentCloudTdmqRocketmqTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Topic name, which can contain 3-64 letters, digits, hyphens, and underscores.",
			},

			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Topic namespace. Currently, you can create topics only in one single namespace.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Topic type. Valid values: Normal, GlobalOrder, PartitionedOrder.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Topic remarks (up to 128 characters).",
			},

			"partition_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "Number of partitions.",
			},

			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation time in milliseconds.",
			},

			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Update time in milliseconds.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmqRocketmq_topic.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = tdmqRocketmq.NewCreateRocketMQTopicRequest()
		clusterId     string
		namespaceName string
		topicName     string
		topicType     string
	)

	if v, ok := d.GetOk("topic_name"); ok {
		topicName = v.(string)
		request.Topic = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		namespaceName = v.(string)
		request.Namespaces = []*string{&namespaceName}
	}

	if v, ok := d.GetOk("type"); ok {
		topicType = v.(string)
		request.Type = helper.String(topicType)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {

		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("partition_num"); ok {
		request.PartitionNum = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateRocketMQTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq topic failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + tccommon.FILED_SP + namespaceName + tccommon.FILED_SP + topicType + tccommon.FILED_SP + topicName)
	return resourceTencentCloudTdmqRocketmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmqRocketmq_topic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TdmqRocketmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceName := idSplit[1]
	topicType := idSplit[2]
	topicName := idSplit[3]

	topicList, err := service.DescribeTdmqRocketmqTopic(ctx, clusterId, namespaceName, topicName)
	if err != nil {
		return err
	}

	if len(topicList) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `topic` %s does not exist", topicName)
	}

	topic := topicList[0]

	_ = d.Set("topic_name", topic.Name)
	_ = d.Set("namespace_name", namespaceName)
	_ = d.Set("type", topicType)
	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("remark", topic.Remark)
	_ = d.Set("partition_num", topic.PartitionNum)
	_ = d.Set("create_time", topic.CreateTime)
	_ = d.Set("update_time", topic.UpdateTime)

	return nil
}

func resourceTencentCloudTdmqRocketmqTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmqRocketmq_topic.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tdmqRocketmq.NewModifyRocketMQTopicRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceName := idSplit[1]
	topicName := idSplit[3]

	request.ClusterId = &clusterId
	request.NamespaceId = &namespaceName
	request.Topic = &topicName

	if d.HasChange("topic") {

		return fmt.Errorf("`topic` do not support change now.")

	}

	if d.HasChange("namespace_name") {

		return fmt.Errorf("`namespace_name` do not support change now.")

	}

	if d.HasChange("type") {

		return fmt.Errorf("`type` do not support change now.")

	}

	if d.HasChange("cluster_id") {

		return fmt.Errorf("`cluster_id` do not support change now.")

	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

	}

	if d.HasChange("partition_num") {
		if v, ok := d.GetOk("partition_num"); ok {
			request.PartitionNum = helper.IntInt64(v.(int))
		}

	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRocketMQTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq topic failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmqRocketmq_topic.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TdmqRocketmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceName := idSplit[1]
	topicName := idSplit[3]

	if err := service.DeleteTdmqRocketmqTopicById(ctx, clusterId, namespaceName, topicName); err != nil {
		return err
	}

	return nil
}

/*
Provides a resource to create a tdmqRocketmq topic

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq"
	remark = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  namespace_name = "test_namespace"
  ttl = 65000
  retention_time = 65000
  remark = "test namespace"
}

resource "tencentcloud_tdmq_rocketmq_topic" "topic" {
  topic_name = "test_rocketmq_topic"
  namespace_name = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
  type = "Normal"
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  remark = "test rocketmq topic"
}
```
Import

tdmqRocketmq topic can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_topic.topic topic_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqRocketmqTopic() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_topic.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateRocketMQTopic(request)
		if e != nil {
			return retryError(e)
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

	d.SetId(clusterId + FILED_SP + namespaceName + FILED_SP + topicType + FILED_SP + topicName)
	return resourceTencentCloudTdmqRocketmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_topic.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
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
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_topic.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmqRocketmq.NewModifyRocketMQTopicRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyRocketMQTopic(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_topic.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
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

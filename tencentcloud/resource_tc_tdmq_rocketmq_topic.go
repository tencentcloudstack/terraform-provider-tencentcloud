/*
Provides a resource to create a tdmqRocketmq topic

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_topic" "topic" {
  topic = &lt;nil&gt;
  namespaces = &lt;nil&gt;
  type = &lt;nil&gt;
  cluster_id = &lt;nil&gt;
  remark = &lt;nil&gt;
  partition_num = &lt;nil&gt;
      }
```

Import

tdmqRocketmq topic can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rocketmq_topic.topic topic_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTdmqRocketmqTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRocketmqTopicCreate,
		Read:   resourceTencentCloudTdmqRocketmqTopicRead,
		Update: resourceTencentCloudTdmqRocketmqTopicUpdate,
		Delete: resourceTencentCloudTdmqRocketmqTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic name, which can contain 3-64 letters, digits, hyphens, and underscores.",
			},

			"namespaces": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Topic namespace. Currently, you can create topics only in one single namespace.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic type. Valid values: Normal, GlobalOrder, PartitionedOrder.",
			},

			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Topic remarks (up to 128 characters).",
			},

			"partition_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of partitions, which doesn&amp;amp;#39;t take effect for globally sequential messages.",
			},

			"name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Topic name.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Creation time in milliseconds.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Update time in milliseconds.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_topic.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = tdmqRocketmq.NewCreateRocketMQTopicRequest()
		response    = tdmqRocketmq.NewCreateRocketMQTopicResponse()
		clusterId   string
		namespaceId string
		topic       string
	)
	if v, ok := d.GetOk("topic"); ok {
		topic = v.(string)
		request.Topic = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesSet := v.(*schema.Set).List()
		for i := range namespacesSet {
			namespaces := namespacesSet[i].(string)
			request.Namespaces = append(request.Namespaces, &namespaces)
		}
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("partition_num"); ok {
		request.PartitionNum = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().CreateRocketMQTopic(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq topic failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(strings.Join([]string{clusterId, namespaceId, topic}, FILED_SP))

	return resourceTencentCloudTdmqRocketmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_topic.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceId := idSplit[1]
	topic := idSplit[2]

	topic, err := service.DescribeTdmqRocketmqTopicById(ctx, clusterId, namespaceId, topic)
	if err != nil {
		return err
	}

	if topic == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRocketmqTopic` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if topic.Topic != nil {
		_ = d.Set("topic", topic.Topic)
	}

	if topic.Namespaces != nil {
		_ = d.Set("namespaces", topic.Namespaces)
	}

	if topic.Type != nil {
		_ = d.Set("type", topic.Type)
	}

	if topic.ClusterId != nil {
		_ = d.Set("cluster_id", topic.ClusterId)
	}

	if topic.Remark != nil {
		_ = d.Set("remark", topic.Remark)
	}

	if topic.PartitionNum != nil {
		_ = d.Set("partition_num", topic.PartitionNum)
	}

	if topic.Name != nil {
		_ = d.Set("name", topic.Name)
	}

	if topic.CreateTime != nil {
		_ = d.Set("create_time", topic.CreateTime)
	}

	if topic.UpdateTime != nil {
		_ = d.Set("update_time", topic.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTdmqRocketmqTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_topic.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmqRocketmq.NewModifyRocketMQTopicRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceId := idSplit[1]
	topic := idSplit[2]

	request.ClusterId = &clusterId
	request.NamespaceId = &namespaceId
	request.Topic = &topic

	immutableArgs := []string{"topic", "namespaces", "type", "cluster_id", "remark", "partition_num", "name", "create_time", "update_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	if d.HasChange("partition_num") {
		if v, ok := d.GetOkExists("partition_num"); ok {
			request.PartitionNum = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().ModifyRocketMQTopic(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmqRocketmq topic failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_topic.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceId := idSplit[1]
	topic := idSplit[2]

	if err := service.DeleteTdmqRocketmqTopicById(ctx, clusterId, namespaceId, topic); err != nil {
		return err
	}

	return nil
}

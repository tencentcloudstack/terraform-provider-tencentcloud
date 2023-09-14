/*
Provides a resource to create a trocket rocketmq_topic

Example Usage

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-xxxxx"
  subnet_id     = "subnet-xxxxx"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_topic" "rocketmq_topic" {
  instance_id = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  topic       = "test_topic"
  topic_type  = "NORMAL"
  queue_num   = 4
  remark      = "test for terraform"
}
```

Import

trocket rocketmq_topic can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_topic.rocketmq_topic instanceId#topic
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
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTrocketRocketmqTopic() *schema.Resource {
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
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Instance Id.",
			},

			"topic": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "topic.",
			},

			"topic_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Topic type. `UNSPECIFIED`: not specified, `NORMAL`: normal message, `FIFO`: sequential message, `DELAY`: delayed message.",
			},

			"queue_num": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of queue. Must be greater than or equal to 3.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "remark.",
			},
		},
	}
}

func resourceTencentCloudTrocketRocketmqTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_topic.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTrocketClient().CreateTopic(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create trocket rocketmqTopic failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	topic = *response.Response.Topic
	d.SetId(instanceId + FILED_SP + topic)

	return resourceTencentCloudTrocketRocketmqTopicRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_topic.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TrocketService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
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
		d.SetId("")
		log.Printf("[WARN]%s resource `TrocketRocketmqTopic` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_topic.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := trocket.NewModifyTopicRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	topic := idSplit[1]

	request.InstanceId = &instanceId
	request.Topic = &topic

	if d.HasChange("queue_num") {
		if v, ok := d.GetOkExists("queue_num"); ok {
			request.QueueNum = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTrocketClient().ModifyTopic(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update trocket rocketmqTopic failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTrocketRocketmqTopicRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_topic.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TrocketService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
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

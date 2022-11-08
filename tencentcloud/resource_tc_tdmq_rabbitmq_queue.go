/*
Provides a resource to create a tdmq rabbitmq_queue

Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_queue" "rabbitmq_queue" {
  queue = ""
  cluster_id = ""
  vhost_id = ""
  auto_delete = ""
  remark = ""
  dead_letter_exchange = ""
  dead_letter_routing_key = ""
}

```
Import

tdmq rabbitmq_queue can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_queue.rabbitmq_queue rabbitmqQueue_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqRabbitmqQueue() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRabbitmqQueueRead,
		Create: resourceTencentCloudTdmqRabbitmqQueueCreate,
		Update: resourceTencentCloudTdmqRabbitmqQueueUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"queue": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "queue name, 3~64 characters.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cluster id.",
			},

			"vhost_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vhost name.",
			},

			"auto_delete": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "auto delete.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "queue description, 128 characters or less.",
			},

			"dead_letter_exchange": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "dead letter exchange.",
			},

			"dead_letter_routing_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "dead letter routing key.",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqQueueCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_queue.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateAMQPQueueRequest()
		clusterId string
		vHostId   string
		queue     string
	)

	if v, ok := d.GetOk("queue"); ok {
		queue = v.(string)
		request.Queue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vhost_id"); ok {
		vHostId = v.(string)
		request.VHostId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("auto_delete"); v != nil {
		request.AutoDelete = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("remark"); ok {

		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dead_letter_exchange"); ok {

		request.DeadLetterExchange = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dead_letter_routing_key"); ok {

		request.DeadLetterRoutingKey = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateAMQPQueue(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqQueue failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + FILED_SP + vHostId + FILED_SP + queue)
	return resourceTencentCloudTdmqRabbitmqQueueRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqQueueRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_queue.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	vHostId := idSplit[1]
	queue := idSplit[2]

	rabbitmqQueue, err := service.DescribeTdmqRabbitmqQueue(ctx, clusterId, vHostId, queue)

	if err != nil {
		return err
	}

	if rabbitmqQueue == nil {
		d.SetId("")
		return fmt.Errorf("resource `rabbitmqQueue` %s does not exist", queue)
	}

	_ = d.Set("queue", queue)
	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("vhost_id", vHostId)

	if rabbitmqQueue.AutoDelete != nil {
		_ = d.Set("auto_delete", rabbitmqQueue.AutoDelete)
	}

	if rabbitmqQueue.Remark != nil {
		_ = d.Set("remark", rabbitmqQueue.Remark)
	}

	if rabbitmqQueue.DeadLetterExchange != nil {
		_ = d.Set("dead_letter_exchange", rabbitmqQueue.DeadLetterExchange)
	}

	if rabbitmqQueue.DeadLetterRoutingKey != nil {
		_ = d.Set("dead_letter_routing_key", rabbitmqQueue.DeadLetterRoutingKey)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_queue.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyAMQPQueueRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	vHostId := idSplit[1]
	queue := idSplit[2]

	request.ClusterId = &clusterId
	request.VHostId = &vHostId
	request.Queue = &queue

	if d.HasChange("queue") {
		return fmt.Errorf("`queue` do not support change now.")
	}

	if d.HasChange("cluster_id") {
		return fmt.Errorf("`cluster_id` do not support change now.")
	}

	if d.HasChange("vhost_id") {
		return fmt.Errorf("`vhost_id` do not support change now.")
	}

	if d.HasChange("auto_delete") {
		if v, ok := d.GetOk("auto_delete"); ok {
			request.AutoDelete = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	if d.HasChange("dead_letter_exchange") {
		if v, ok := d.GetOk("dead_letter_exchange"); ok {
			request.DeadLetterExchange = helper.String(v.(string))
		}
	}

	if d.HasChange("dead_letter_routing_key") {
		if v, ok := d.GetOk("dead_letter_routing_key"); ok {
			request.DeadLetterRoutingKey = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyAMQPQueue(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqQueue failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRabbitmqQueueRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqQueueDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_queue.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	vHostId := idSplit[1]
	queue := idSplit[2]

	if err := service.DeleteTdmqRabbitmqQueueById(ctx, clusterId, vHostId, queue); err != nil {
		return err
	}

	return nil
}

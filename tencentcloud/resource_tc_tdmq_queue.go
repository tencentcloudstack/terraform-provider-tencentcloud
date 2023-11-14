/*
Provides a resource to create a tdmq queue

Example Usage

```hcl
resource "tencentcloud_tdmq_queue" "queue" {
  queue_name = "queue_name"
  max_msg_heap_num = 10000000
  polling_wait_seconds = 0
  visibility_timeout = 30
  max_msg_size = 65536
  msg_retention_seconds = 3600
  rewind_seconds = 0
  transaction = 1
  first_query_interval = 1
  max_query_count = 1
  dead_letter_queue_name = "dead_letter_queue_name"
  policy = 0
  max_receive_count = 50
  max_time_to_live = 300
  trace = false
  retention_size_in_m_b = 0
}
```

Import

tdmq queue can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_queue.queue queue_id
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

func resourceTencentCloudTdmqQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqQueueCreate,
		Read:   resourceTencentCloudTdmqQueueRead,
		Update: resourceTencentCloudTdmqQueueUpdate,
		Delete: resourceTencentCloudTdmqQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"queue_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Queue name, which must be unique under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.",
			},

			"max_msg_heap_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum number of heaped messages. The value range is 1,000,000-10,000,000 during the beta test and can be 1,000,000-1,000,000,000 after the product is officially released. The default value is 10,000,000 during the beta test and will be 100,000,000 after the product is officially released.",
			},

			"polling_wait_seconds": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Long polling wait time for message reception. Value range: 0-30 seconds. Default value: 0.",
			},

			"visibility_timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Message visibility timeout period. Value range: 1-43200 seconds (i.e., 12 hours). Default value: 30.",
			},

			"max_msg_size": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum message length. Value range: 1024-65536 bytes (i.e., 1-64 KB). Default value: 65536.",
			},

			"msg_retention_seconds": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The max period during which a message is retained before it is automatically acknowledged. Value range: 30-43,200 seconds (30 seconds to 12 hours). Default value: 3600 seconds (1 hour).",
			},

			"rewind_seconds": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Rewindable time of messages in the queue. Value range: 0-1,296,000s (if message rewind is enabled). The value “0” indicates that message rewind is not enabled.",
			},

			"transaction": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "1: transaction queue; 0: general queue.",
			},

			"first_query_interval": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "First lookback interval.",
			},

			"max_query_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum number of lookbacks.",
			},

			"dead_letter_queue_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Dead letter queue name.",
			},

			"policy": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Dead letter policy. 0: message has been consumed multiple times but not deleted; 1: `Time-To-Live` has elapsed.",
			},

			"max_receive_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum receipt times. Value range: 1-1000.",
			},

			"max_time_to_live": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum period in seconds before an unconsumed message expires, which is required if `policy` is 1. Value range: 300-43200. This value should be smaller than `msgRetentionSeconds` (maximum message retention period).",
			},

			"trace": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable message trace. true: yes; false: no. If this field is not configured, the feature will not be enabled.",
			},

			"retention_size_in_m_b": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Queue storage space configured for message rewind. Value range: 10,240-512,000 MB (if message rewind is enabled). The value “0” indicates that message rewind is not enabled.",
			},
		},
	}
}

func resourceTencentCloudTdmqQueueCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_queue.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateCmqQueueRequest()
		response  = tdmq.NewCreateCmqQueueResponse()
		queueName string
	)
	if v, ok := d.GetOk("queue_name"); ok {
		queueName = v.(string)
		request.QueueName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_msg_heap_num"); ok {
		request.MaxMsgHeapNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("polling_wait_seconds"); ok {
		request.PollingWaitSeconds = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("visibility_timeout"); ok {
		request.VisibilityTimeout = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_msg_size"); ok {
		request.MaxMsgSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("msg_retention_seconds"); ok {
		request.MsgRetentionSeconds = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("rewind_seconds"); ok {
		request.RewindSeconds = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("transaction"); ok {
		request.Transaction = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("first_query_interval"); ok {
		request.FirstQueryInterval = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_query_count"); ok {
		request.MaxQueryCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("dead_letter_queue_name"); ok {
		request.DeadLetterQueueName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("policy"); ok {
		request.Policy = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_receive_count"); ok {
		request.MaxReceiveCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_time_to_live"); ok {
		request.MaxTimeToLive = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("trace"); ok {
		request.Trace = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("retention_size_in_m_b"); ok {
		request.RetentionSizeInMB = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateCmqQueue(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmq queue failed, reason:%+v", logId, err)
		return err
	}

	queueName = *response.Response.QueueName
	d.SetId(queueName)

	return resourceTencentCloudTdmqQueueRead(d, meta)
}

func resourceTencentCloudTdmqQueueRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_queue.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	queueId := d.Id()

	queue, err := service.DescribeTdmqQueueById(ctx, queueName)
	if err != nil {
		return err
	}

	if queue == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqQueue` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if queue.QueueName != nil {
		_ = d.Set("queue_name", queue.QueueName)
	}

	if queue.MaxMsgHeapNum != nil {
		_ = d.Set("max_msg_heap_num", queue.MaxMsgHeapNum)
	}

	if queue.PollingWaitSeconds != nil {
		_ = d.Set("polling_wait_seconds", queue.PollingWaitSeconds)
	}

	if queue.VisibilityTimeout != nil {
		_ = d.Set("visibility_timeout", queue.VisibilityTimeout)
	}

	if queue.MaxMsgSize != nil {
		_ = d.Set("max_msg_size", queue.MaxMsgSize)
	}

	if queue.MsgRetentionSeconds != nil {
		_ = d.Set("msg_retention_seconds", queue.MsgRetentionSeconds)
	}

	if queue.RewindSeconds != nil {
		_ = d.Set("rewind_seconds", queue.RewindSeconds)
	}

	if queue.Transaction != nil {
		_ = d.Set("transaction", queue.Transaction)
	}

	if queue.FirstQueryInterval != nil {
		_ = d.Set("first_query_interval", queue.FirstQueryInterval)
	}

	if queue.MaxQueryCount != nil {
		_ = d.Set("max_query_count", queue.MaxQueryCount)
	}

	if queue.DeadLetterQueueName != nil {
		_ = d.Set("dead_letter_queue_name", queue.DeadLetterQueueName)
	}

	if queue.Policy != nil {
		_ = d.Set("policy", queue.Policy)
	}

	if queue.MaxReceiveCount != nil {
		_ = d.Set("max_receive_count", queue.MaxReceiveCount)
	}

	if queue.MaxTimeToLive != nil {
		_ = d.Set("max_time_to_live", queue.MaxTimeToLive)
	}

	if queue.Trace != nil {
		_ = d.Set("trace", queue.Trace)
	}

	if queue.RetentionSizeInMB != nil {
		_ = d.Set("retention_size_in_m_b", queue.RetentionSizeInMB)
	}

	return nil
}

func resourceTencentCloudTdmqQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_queue.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyCmqQueueAttributeRequest()

	queueId := d.Id()

	request.QueueName = &queueName

	immutableArgs := []string{"queue_name", "max_msg_heap_num", "polling_wait_seconds", "visibility_timeout", "max_msg_size", "msg_retention_seconds", "rewind_seconds", "transaction", "first_query_interval", "max_query_count", "dead_letter_queue_name", "policy", "max_receive_count", "max_time_to_live", "trace", "retention_size_in_m_b"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("queue_name") {
		if v, ok := d.GetOk("queue_name"); ok {
			request.QueueName = helper.String(v.(string))
		}
	}

	if d.HasChange("max_msg_heap_num") {
		if v, ok := d.GetOkExists("max_msg_heap_num"); ok {
			request.MaxMsgHeapNum = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("polling_wait_seconds") {
		if v, ok := d.GetOkExists("polling_wait_seconds"); ok {
			request.PollingWaitSeconds = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("visibility_timeout") {
		if v, ok := d.GetOkExists("visibility_timeout"); ok {
			request.VisibilityTimeout = helper.IntUint64(v.(int))
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

	if d.HasChange("rewind_seconds") {
		if v, ok := d.GetOkExists("rewind_seconds"); ok {
			request.RewindSeconds = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("transaction") {
		if v, ok := d.GetOkExists("transaction"); ok {
			request.Transaction = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("first_query_interval") {
		if v, ok := d.GetOkExists("first_query_interval"); ok {
			request.FirstQueryInterval = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_query_count") {
		if v, ok := d.GetOkExists("max_query_count"); ok {
			request.MaxQueryCount = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("dead_letter_queue_name") {
		if v, ok := d.GetOk("dead_letter_queue_name"); ok {
			request.DeadLetterQueueName = helper.String(v.(string))
		}
	}

	if d.HasChange("policy") {
		if v, ok := d.GetOkExists("policy"); ok {
			request.Policy = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_receive_count") {
		if v, ok := d.GetOkExists("max_receive_count"); ok {
			request.MaxReceiveCount = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_time_to_live") {
		if v, ok := d.GetOkExists("max_time_to_live"); ok {
			request.MaxTimeToLive = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("trace") {
		if v, ok := d.GetOkExists("trace"); ok {
			request.Trace = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("retention_size_in_m_b") {
		if v, ok := d.GetOkExists("retention_size_in_m_b"); ok {
			request.RetentionSizeInMB = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyCmqQueueAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmq queue failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqQueueRead(d, meta)
}

func resourceTencentCloudTdmqQueueDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_queue.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	queueId := d.Id()

	if err := service.DeleteTdmqQueueById(ctx, queueName); err != nil {
		return err
	}

	return nil
}

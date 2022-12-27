/*
Provides a resource to create a tcmq queue

Example Usage

```hcl
resource "tencentcloud_tcmq_queue" "queue" {
  queue_name = "queue_name"
}
```

Import

tcmq queue can be imported using the id, e.g.

```
terraform import tencentcloud_tcmq_queue.queue queue_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcmqQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcmqQueueCreate,
		Read:   resourceTencentCloudTcmqQueueRead,
		Update: resourceTencentCloudTcmqQueueUpdate,
		Delete: resourceTencentCloudTcmqQueueDelete,
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
				Default:     10000000,
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
				Default:     30,
				Type:        schema.TypeInt,
				Description: "Message visibility timeout period. Value range: 1-43200 seconds (i.e., 12 hours). Default value: 30.",
			},

			"max_msg_size": {
				Optional:    true,
				Default:     65536,
				Type:        schema.TypeInt,
				Description: "Maximum message length. Value range: 1024-65536 bytes (i.e., 1-64 KB). Default value: 65536.",
			},

			"msg_retention_seconds": {
				Optional:    true,
				Default:     3600,
				Type:        schema.TypeInt,
				Description: "The max period during which a message is retained before it is automatically acknowledged. Value range: 30-43,200 seconds (30 seconds to 12 hours). Default value: 3600 seconds (1 hour).",
			},

			"rewind_seconds": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Rewindable time of messages in the queue. Value range: 0-1,296,000s (if message rewind is enabled). The value `0` indicates that message rewind is not enabled.",
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
				Default:     50,
				Type:        schema.TypeInt,
				Description: "Maximum receipt times. Value range: 1-1000.",
			},

			"max_time_to_live": {
				Optional:    true,
				Default:     300,
				Type:        schema.TypeInt,
				Description: "Maximum period in seconds before an unconsumed message expires, which is required if `policy` is 1. Value range: 300-43200. This value should be smaller than `msgRetentionSeconds` (maximum message retention period).",
			},

			"trace": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable message trace. true: yes; false: no. If this field is not configured, the feature will not be enabled.",
			},

			"retention_size_in_mb": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Queue storage space configured for message rewind. Value range: 10,240-512,000 MB (if message rewind is enabled). The value `0` indicates that message rewind is not enabled.",
			},
		},
	}
}

func resourceTencentCloudTcmqQueueCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcmq_queue.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tcmq.NewCreateCmqQueueRequest()
		queueName string
	)
	if v, ok := d.GetOk("queue_name"); ok {
		queueName = v.(string)
		request.QueueName = helper.String(queueName)
	}

	if v, _ := d.GetOk("max_msg_heap_num"); v != nil {
		request.MaxMsgHeapNum = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("polling_wait_seconds"); v != nil {
		request.PollingWaitSeconds = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("visibility_timeout"); v != nil {
		request.VisibilityTimeout = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("max_msg_size"); v != nil {
		request.MaxMsgSize = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("msg_retention_seconds"); v != nil {
		request.MsgRetentionSeconds = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("rewind_seconds"); v != nil {
		request.RewindSeconds = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("transaction"); v != nil {
		request.Transaction = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("first_query_interval"); v != nil {
		request.FirstQueryInterval = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("max_query_count"); v != nil {
		request.MaxQueryCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("dead_letter_queue_name"); ok {
		request.DeadLetterQueueName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("policy"); v != nil {
		request.Policy = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("max_receive_count"); v != nil {
		request.MaxReceiveCount = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("max_time_to_live"); v != nil {
		request.MaxTimeToLive = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("trace"); v != nil {
		request.Trace = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("retention_size_in_mb"); v != nil {
		request.RetentionSizeInMB = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateCmqQueue(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcmq queue failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(queueName)

	return resourceTencentCloudTcmqQueueRead(d, meta)
}

func resourceTencentCloudTcmqQueueRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcmq_queue.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	queueName := d.Id()

	queue, err := service.DescribeTcmqQueueById(ctx, queueName)
	if err != nil {
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "ResourceNotFound" {
				return nil
			}
		}
		return err
	}

	if queue == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
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

	if queue.TransactionPolicy != nil {
		if queue.TransactionPolicy.FirstQueryInterval != nil {
			_ = d.Set("first_query_interval", queue.TransactionPolicy.FirstQueryInterval)
		}

		if queue.TransactionPolicy.MaxQueryCount != nil {
			_ = d.Set("max_query_count", queue.TransactionPolicy.MaxQueryCount)
		}
	}

	if len(queue.DeadLetterSource) > 0 && queue.DeadLetterSource[0].QueueName != nil {
		_ = d.Set("dead_letter_queue_name", queue.DeadLetterSource[0].QueueName)
	}

	if queue.DeadLetterPolicy != nil {
		if queue.DeadLetterPolicy.Policy != nil {
			_ = d.Set("policy", queue.DeadLetterPolicy.Policy)
		}

		if queue.DeadLetterPolicy.MaxReceiveCount != nil {
			_ = d.Set("max_receive_count", queue.DeadLetterPolicy.MaxReceiveCount)
		}

		if queue.DeadLetterPolicy.MaxTimeToLive != nil {
			_ = d.Set("max_time_to_live", queue.DeadLetterPolicy.MaxTimeToLive)
		}
	}

	if queue.Trace != nil {
		_ = d.Set("trace", queue.Trace)
	}

	if queue.RetentionSizeInMB != nil {
		_ = d.Set("retention_size_in_mb", queue.RetentionSizeInMB)
	}

	return nil
}

func resourceTencentCloudTcmqQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcmq_queue.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcmq.NewModifyCmqQueueAttributeRequest()

	queueName := d.Id()

	request.QueueName = &queueName
	if d.HasChange("queue_name") {
		if v, ok := d.GetOk("queue_name"); ok {
			request.QueueName = helper.String(v.(string))
		}
	}

	if d.HasChange("max_msg_heap_num") {
		if v, _ := d.GetOk("max_msg_heap_num"); v != nil {
			request.MaxMsgHeapNum = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("polling_wait_seconds") {
		if v, _ := d.GetOk("polling_wait_seconds"); v != nil {
			request.PollingWaitSeconds = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("visibility_timeout") {
		if v, _ := d.GetOk("visibility_timeout"); v != nil {
			request.VisibilityTimeout = helper.IntUint64(v.(int))
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

	if d.HasChange("rewind_seconds") {
		if v, _ := d.GetOk("rewind_seconds"); v != nil {
			request.RewindSeconds = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("transaction") {
		if v, _ := d.GetOk("transaction"); v != nil {
			request.Transaction = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("first_query_interval") {
		if v, _ := d.GetOk("first_query_interval"); v != nil {
			request.FirstQueryInterval = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_query_count") {
		if v, _ := d.GetOk("max_query_count"); v != nil {
			request.MaxQueryCount = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("dead_letter_queue_name") {
		if v, ok := d.GetOk("dead_letter_queue_name"); ok {
			request.DeadLetterQueueName = helper.String(v.(string))
		}
	}

	if d.HasChange("policy") {
		if v, _ := d.GetOk("policy"); v != nil {
			request.Policy = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_receive_count") {
		if v, _ := d.GetOk("max_receive_count"); v != nil {
			request.MaxReceiveCount = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_time_to_live") {
		if v, _ := d.GetOk("max_time_to_live"); v != nil {
			request.MaxTimeToLive = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("trace") {
		if v, _ := d.GetOk("trace"); v != nil {
			request.Trace = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("retention_size_in_mb") {
		if v, _ := d.GetOk("retention_size_in_mb"); v != nil {
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
		log.Printf("[CRITAL]%s create tcmq queue failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmqQueueRead(d, meta)
}

func resourceTencentCloudTcmqQueueDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcmq_queue.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	queueName := d.Id()

	if err := service.DeleteTcmqQueueById(ctx, queueName); err != nil {
		return err
	}

	return nil
}

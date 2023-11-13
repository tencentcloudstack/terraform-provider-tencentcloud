/*
Use this data source to query detailed information of tdmq dead_letter_source_queue

Example Usage

```hcl
data "tencentcloud_tdmq_dead_letter_source_queue" "dead_letter_source_queue" {
  dead_letter_queue_name = ""
  source_queue_name = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdmqDeadLetterSourceQueue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqDeadLetterSourceQueueRead,
		Schema: map[string]*schema.Schema{
			"dead_letter_queue_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Dead letter queue name.",
			},

			"source_queue_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter by SourceQueueName.",
			},

			"queue_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Dead letter queue source queue.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message queue ID.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"queue_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message queue name.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqDeadLetterSourceQueueRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_dead_letter_source_queue.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("dead_letter_queue_name"); ok {
		paramMap["DeadLetterQueueName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_queue_name"); ok {
		paramMap["SourceQueueName"] = helper.String(v.(string))
	}

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var queueSet []*tdmq.CmqDeadLetterSource

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqDeadLetterSourceQueueByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		queueSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(queueSet))
	tmpList := make([]map[string]interface{}, 0, len(queueSet))

	if queueSet != nil {
		for _, cmqDeadLetterSource := range queueSet {
			cmqDeadLetterSourceMap := map[string]interface{}{}

			if cmqDeadLetterSource.QueueId != nil {
				cmqDeadLetterSourceMap["queue_id"] = cmqDeadLetterSource.QueueId
			}

			if cmqDeadLetterSource.QueueName != nil {
				cmqDeadLetterSourceMap["queue_name"] = cmqDeadLetterSource.QueueName
			}

			ids = append(ids, *cmqDeadLetterSource.QueueId)
			tmpList = append(tmpList, cmqDeadLetterSourceMap)
		}

		_ = d.Set("queue_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

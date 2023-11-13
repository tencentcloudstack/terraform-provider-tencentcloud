/*
Use this data source to query detailed information of tdmq queue

Example Usage

```hcl
data "tencentcloud_tdmq_queue" "queue" {
  offset = 0
  limit = 20
  queue_name = "queue_name"
  queue_name_list = &lt;nil&gt;
  is_tag_filter = true
  filters {
		name = "tag"
		values = &lt;nil&gt;

  }
  total_count = &lt;nil&gt;
  queue_list {
		queue_id = &lt;nil&gt;
		queue_name = &lt;nil&gt;
		qps = &lt;nil&gt;
		bps = &lt;nil&gt;
		max_delay_seconds = &lt;nil&gt;
		max_msg_heap_num = &lt;nil&gt;
		polling_wait_seconds = &lt;nil&gt;
		msg_retention_seconds = &lt;nil&gt;
		visibility_timeout = &lt;nil&gt;
		max_msg_size = &lt;nil&gt;
		rewind_seconds = &lt;nil&gt;
		create_time = &lt;nil&gt;
		last_modify_time = &lt;nil&gt;
		active_msg_num = &lt;nil&gt;
		inactive_msg_num = &lt;nil&gt;
		delay_msg_num = &lt;nil&gt;
		rewind_msg_num = &lt;nil&gt;
		min_msg_time = &lt;nil&gt;
		transaction = &lt;nil&gt;
		dead_letter_source {
			queue_id = &lt;nil&gt;
			queue_name = &lt;nil&gt;
		}
		dead_letter_policy {
			dead_letter_queue = &lt;nil&gt;
			policy = &lt;nil&gt;
			max_time_to_live = &lt;nil&gt;
			max_receive_count = &lt;nil&gt;
		}
		transaction_policy {
			first_query_interval = &lt;nil&gt;
			max_query_count = &lt;nil&gt;
		}
		create_uin = &lt;nil&gt;
		tags {
			tag_key = &lt;nil&gt;
			tag_value = &lt;nil&gt;
		}
		trace = &lt;nil&gt;
		tenant_id = &lt;nil&gt;
		namespace_name = &lt;nil&gt;
		status = &lt;nil&gt;
		max_unacked_msg_num = &lt;nil&gt;
		max_msg_backlog_size = &lt;nil&gt;
		retention_size_in_m_b = &lt;nil&gt;

  }
  request_id = &lt;nil&gt;
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

func dataSourceTencentCloudTdmqQueue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqQueueRead,
		Schema: map[string]*schema.Schema{
			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Starting position of a queue list to be returned on the current page in case of paginated return. If a value is entered, limit must be specified. If this parameter is left empty, 0 will be used by default.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of queues to be returned per page in case of paginated return. If this parameter is not passed in, 20 will be used by default. Maximum value: 50.",
			},

			"queue_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter by QueueName.",
			},

			"queue_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by CMQ queue name.",
			},

			"is_tag_filter": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "For filtering by tag, this parameter must be set to `true`.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter. Currently, you can filter by tag. The tag name must be prefixed with “tag:”, such as “tag: owner”, “tag: environment”, or “tag: business”.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter parameter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Value.",
						},
					},
				},
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "The number of queues.",
			},

			"queue_list": {
				Type:        schema.TypeList,
				Description: "Queue list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_id": {
							Type:        schema.TypeString,
							Description: "Message queue ID.",
						},
						"queue_name": {
							Type:        schema.TypeString,
							Description: "Message queue name.",
						},
						"qps": {
							Type:        schema.TypeInt,
							Description: "Limit of the number of messages produced per second. The value for consumed messages is 1.1 times this value.",
						},
						"bps": {
							Type:        schema.TypeInt,
							Description: "Bandwidth limit.",
						},
						"max_delay_seconds": {
							Type:        schema.TypeInt,
							Description: "Maximum retention period for inflight messages.",
						},
						"max_msg_heap_num": {
							Type:        schema.TypeInt,
							Description: "Maximum number of heaped messages. The value range is 1,000,000-10,000,000 during the beta test and can be 1,000,000-1,000,000,000 after the product is officially released. The default value is 10,000,000 during the beta test and will be 100,000,000 after the product is officially released.",
						},
						"polling_wait_seconds": {
							Type:        schema.TypeInt,
							Description: "Long polling wait time for message reception. Value range: 0-30 seconds. Default value: 0.",
						},
						"msg_retention_seconds": {
							Type:        schema.TypeInt,
							Description: "The max period during which a message is retained before it is automatically acknowledged. Value range: 30-43,200 seconds (30 seconds to 12 hours). Default value: 3600 seconds (1 hour).",
						},
						"visibility_timeout": {
							Type:        schema.TypeInt,
							Description: "Message visibility timeout period. Value range: 1-43200 seconds (i.e., 12 hours). Default value: 30.",
						},
						"max_msg_size": {
							Type:        schema.TypeInt,
							Description: "Maximum message length. Value range: 1024-65536 bytes (i.e., 1-64 KB). Default value: 65536.",
						},
						"rewind_seconds": {
							Type:        schema.TypeInt,
							Description: "Rewindable time of messages in the queue. Value range: 0-1,296,000s (if message rewind is enabled). The value “0” indicates that message rewind is not enabled.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Description: "Queue creation time. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"last_modify_time": {
							Type:        schema.TypeInt,
							Description: "Time when the queue attribute is last modified. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"active_msg_num": {
							Type:        schema.TypeInt,
							Description: "Total number of messages in `Active` status (i.e., unconsumed) in the queue, which is an approximate value.",
						},
						"inactive_msg_num": {
							Type:        schema.TypeInt,
							Description: "Total number of messages in `Inactive` status (i.e., being consumed) in the queue, which is an approximate value.",
						},
						"delay_msg_num": {
							Type:        schema.TypeInt,
							Description: "Number of delayed messages.",
						},
						"rewind_msg_num": {
							Type:        schema.TypeInt,
							Description: "Number of retained messages which have been deleted by the `DelMsg` API but are still within their rewind time range.",
						},
						"min_msg_time": {
							Type:        schema.TypeInt,
							Description: "Minimum unconsumed time of message in seconds.",
						},
						"transaction": {
							Type:        schema.TypeBool,
							Description: "1: transaction queue; 0: general queue.",
						},
						"dead_letter_source": {
							Type:        schema.TypeList,
							Description: "Dead letter queue.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"queue_id": {
										Type:        schema.TypeString,
										Description: "Message queue ID.",
									},
									"queue_name": {
										Type:        schema.TypeString,
										Description: "Message queue name.",
									},
								},
							},
						},
						"dead_letter_policy": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Dead letter queue policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dead_letter_queue": {
										Type:        schema.TypeString,
										Description: "Dead letter queue.",
									},
									"policy": {
										Type:        schema.TypeInt,
										Description: "Dead letter queue policy.",
									},
									"max_time_to_live": {
										Type:        schema.TypeInt,
										Description: "Maximum period in seconds before an unconsumed message expires, which is required if `Policy` is 1. Value range: 300-43200. This value should be smaller than `MsgRetentionSeconds` (maximum message retention period).",
									},
									"max_receive_count": {
										Type:        schema.TypeInt,
										Description: "Maximum number of receipts.",
									},
								},
							},
						},
						"transaction_policy": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Transaction message policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"first_query_interval": {
										Type:        schema.TypeInt,
										Description: "First lookback time.",
									},
									"max_query_count": {
										Type:        schema.TypeInt,
										Description: "Maximum number of queries.",
									},
								},
							},
						},
						"create_uin": {
							Type:        schema.TypeInt,
							Description: "Creator `Uin`.",
						},
						"tags": {
							Type:        schema.TypeList,
							Description: "Associated tag.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Description: "Value of the tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Description: "Value of the tag value.",
									},
								},
							},
						},
						"trace": {
							Type:        schema.TypeBool,
							Description: "Message trace. true: enabled; false: not enabled.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Description: "Tenant ID.",
						},
						"namespace_name": {
							Type:        schema.TypeString,
							Description: "Namespace name.",
						},
						"status": {
							Type:        schema.TypeInt,
							Description: "Cluster status. `0`: creating; `1`: normal; `2`: terminating; `3`: deleted; `4`: isolated; `5`: creation failed; `6`: deletion failed.",
						},
						"max_unacked_msg_num": {
							Type:        schema.TypeInt,
							Description: "The maximum number of unacknowledged messages.",
						},
						"max_msg_backlog_size": {
							Type:        schema.TypeInt,
							Description: "Maximum size of heaped messages in bytes.",
						},
						"retention_size_in_m_b": {
							Type:        schema.TypeInt,
							Description: "Queue storage space configured for message rewind. Value range: 1,024-10,240 MB (if message rewind is enabled). The value “0” indicates that message rewind is not enabled.",
						},
					},
				},
			},

			"request_id": {
				Type:        schema.TypeString,
				Description: "The unique request ID, which is returned for each request. RequestId is required for locating a problem.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqQueueRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_queue.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("queue_name"); ok {
		paramMap["QueueName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("queue_name_list"); ok {
		queueNameListSet := v.(*schema.Set).List()
		paramMap["QueueNameList"] = helper.InterfacesStringsPoint(queueNameListSet)
	}

	if v, _ := d.GetOk("is_tag_filter"); v != nil {
		paramMap["IsTagFilter"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tdmq.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tdmq.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("queue_list"); ok {
		queueListSet := v.([]interface{})
		tmpSet := make([]*tdmq.CmqQueue, 0, len(queueListSet))

		for _, item := range queueListSet {
			cmqQueue := tdmq.CmqQueue{}
			cmqQueueMap := item.(map[string]interface{})

			if v, ok := cmqQueueMap["queue_id"]; ok {
				cmqQueue.QueueId = helper.String(v.(string))
			}
			if v, ok := cmqQueueMap["queue_name"]; ok {
				cmqQueue.QueueName = helper.String(v.(string))
			}
			if v, ok := cmqQueueMap["qps"]; ok {
				cmqQueue.Qps = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["bps"]; ok {
				cmqQueue.Bps = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["max_delay_seconds"]; ok {
				cmqQueue.MaxDelaySeconds = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["max_msg_heap_num"]; ok {
				cmqQueue.MaxMsgHeapNum = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["polling_wait_seconds"]; ok {
				cmqQueue.PollingWaitSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["msg_retention_seconds"]; ok {
				cmqQueue.MsgRetentionSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["visibility_timeout"]; ok {
				cmqQueue.VisibilityTimeout = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["max_msg_size"]; ok {
				cmqQueue.MaxMsgSize = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["rewind_seconds"]; ok {
				cmqQueue.RewindSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["create_time"]; ok {
				cmqQueue.CreateTime = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["last_modify_time"]; ok {
				cmqQueue.LastModifyTime = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["active_msg_num"]; ok {
				cmqQueue.ActiveMsgNum = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["inactive_msg_num"]; ok {
				cmqQueue.InactiveMsgNum = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["delay_msg_num"]; ok {
				cmqQueue.DelayMsgNum = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["rewind_msg_num"]; ok {
				cmqQueue.RewindMsgNum = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["min_msg_time"]; ok {
				cmqQueue.MinMsgTime = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["transaction"]; ok {
				cmqQueue.Transaction = helper.Bool(v.(bool))
			}
			if v, ok := cmqQueueMap["dead_letter_source"]; ok {
				for _, item := range v.([]interface{}) {
					deadLetterSourceMap := item.(map[string]interface{})
					cmqDeadLetterSource := tdmq.CmqDeadLetterSource{}
					if v, ok := deadLetterSourceMap["queue_id"]; ok {
						cmqDeadLetterSource.QueueId = helper.String(v.(string))
					}
					if v, ok := deadLetterSourceMap["queue_name"]; ok {
						cmqDeadLetterSource.QueueName = helper.String(v.(string))
					}
					cmqQueue.DeadLetterSource = append(cmqQueue.DeadLetterSource, &cmqDeadLetterSource)
				}
			}
			if deadLetterPolicyMap, ok := helper.InterfaceToMap(cmqQueueMap, "dead_letter_policy"); ok {
				cmqDeadLetterPolicy := tdmq.CmqDeadLetterPolicy{}
				if v, ok := deadLetterPolicyMap["dead_letter_queue"]; ok {
					cmqDeadLetterPolicy.DeadLetterQueue = helper.String(v.(string))
				}
				if v, ok := deadLetterPolicyMap["policy"]; ok {
					cmqDeadLetterPolicy.Policy = helper.IntUint64(v.(int))
				}
				if v, ok := deadLetterPolicyMap["max_time_to_live"]; ok {
					cmqDeadLetterPolicy.MaxTimeToLive = helper.IntUint64(v.(int))
				}
				if v, ok := deadLetterPolicyMap["max_receive_count"]; ok {
					cmqDeadLetterPolicy.MaxReceiveCount = helper.IntUint64(v.(int))
				}
				cmqQueue.DeadLetterPolicy = &cmqDeadLetterPolicy
			}
			if transactionPolicyMap, ok := helper.InterfaceToMap(cmqQueueMap, "transaction_policy"); ok {
				cmqTransactionPolicy := tdmq.CmqTransactionPolicy{}
				if v, ok := transactionPolicyMap["first_query_interval"]; ok {
					cmqTransactionPolicy.FirstQueryInterval = helper.IntUint64(v.(int))
				}
				if v, ok := transactionPolicyMap["max_query_count"]; ok {
					cmqTransactionPolicy.MaxQueryCount = helper.IntUint64(v.(int))
				}
				cmqQueue.TransactionPolicy = &cmqTransactionPolicy
			}
			if v, ok := cmqQueueMap["create_uin"]; ok {
				cmqQueue.CreateUin = helper.IntUint64(v.(int))
			}
			if v, ok := cmqQueueMap["tags"]; ok {
				for _, item := range v.([]interface{}) {
					tagsMap := item.(map[string]interface{})
					tag := tdmq.Tag{}
					if v, ok := tagsMap["tag_key"]; ok {
						tag.TagKey = helper.String(v.(string))
					}
					if v, ok := tagsMap["tag_value"]; ok {
						tag.TagValue = helper.String(v.(string))
					}
					cmqQueue.Tags = append(cmqQueue.Tags, &tag)
				}
			}
			if v, ok := cmqQueueMap["trace"]; ok {
				cmqQueue.Trace = helper.Bool(v.(bool))
			}
			if v, ok := cmqQueueMap["tenant_id"]; ok {
				cmqQueue.TenantId = helper.String(v.(string))
			}
			if v, ok := cmqQueueMap["namespace_name"]; ok {
				cmqQueue.NamespaceName = helper.String(v.(string))
			}
			if v, ok := cmqQueueMap["status"]; ok {
				cmqQueue.Status = helper.IntInt64(v.(int))
			}
			if v, ok := cmqQueueMap["max_unacked_msg_num"]; ok {
				cmqQueue.MaxUnackedMsgNum = helper.IntInt64(v.(int))
			}
			if v, ok := cmqQueueMap["max_msg_backlog_size"]; ok {
				cmqQueue.MaxMsgBacklogSize = helper.IntInt64(v.(int))
			}
			if v, ok := cmqQueueMap["retention_size_in_m_b"]; ok {
				cmqQueue.RetentionSizeInMB = helper.IntUint64(v.(int))
			}
			tmpSet = append(tmpSet, &cmqQueue)
		}
		paramMap["queue_list"] = tmpSet
	}

	if v, ok := d.GetOk("request_id"); ok {
		paramMap["RequestId"] = helper.String(v.(string))
	}

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var queueList []*tdmq.CmqQueue

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqQueueByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		queueList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(queueList))
	tmpList := make([]map[string]interface{}, 0, len(queueList))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

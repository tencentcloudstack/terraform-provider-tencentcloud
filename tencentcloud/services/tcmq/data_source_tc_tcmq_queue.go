package tcmq

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcmqQueue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcmqQueueRead,
		Schema: map[string]*schema.Schema{
			"offset": {
				Default:     0,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Starting position of a queue list to be returned on the current page in case of paginated return. If a value is entered, limit must be specified. If this parameter is left empty, 0 will be used by default.",
			},

			"limit": {
				Default:     20,
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
				Description: "Filter. Currently, you can filter by tag. The tag name must be prefixed with `tag:`, such as `tag: owner`, `tag: environment`, or `tag: business`.",
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

			"queue_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Queue list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Message queue ID.",
						},
						"queue_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Message queue name.",
						},
						"qps": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Limit of the number of messages produced per second. The value for consumed messages is 1.1 times this value.",
						},
						"bps": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Bandwidth limit.",
						},
						"max_delay_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum retention period for inflight messages.",
						},
						"max_msg_heap_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum number of heaped messages. The value range is 1,000,000-10,000,000 during the beta test and can be 1,000,000-1,000,000,000 after the product is officially released. The default value is 10,000,000 during the beta test and will be 100,000,000 after the product is officially released.",
						},
						"polling_wait_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Long polling wait time for message reception. Value range: 0-30 seconds. Default value: 0.",
						},
						"msg_retention_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The max period during which a message is retained before it is automatically acknowledged. Value range: 30-43,200 seconds (30 seconds to 12 hours). Default value: 3600 seconds (1 hour).",
						},
						"visibility_timeout": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Message visibility timeout period. Value range: 1-43200 seconds (i.e., 12 hours). Default value: 30.",
						},
						"max_msg_size": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum message length. Value range: 1024-65536 bytes (i.e., 1-64 KB). Default value: 65536.",
						},
						"rewind_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Rewindable time of messages in the queue. Value range: 0-1,296,000s (if message rewind is enabled). The value `0` indicates that message rewind is not enabled.",
						},
						"create_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Queue creation time. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"last_modify_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Time when the queue attribute is last modified. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"active_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Total number of messages in `Active` status (i.e., unconsumed) in the queue, which is an approximate value.",
						},
						"inactive_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Total number of messages in `Inactive` status (i.e., being consumed) in the queue, which is an approximate value.",
						},
						"delay_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Number of delayed messages.",
						},
						"rewind_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Number of retained messages which have been deleted by the `DelMsg` API but are still within their rewind time range.",
						},
						"min_msg_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Minimum unconsumed time of message in seconds.",
						},
						"transaction": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "1: transaction queue; 0: general queue.",
						},
						"dead_letter_source": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Dead letter queue.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"queue_id": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Message queue ID.",
									},
									"queue_name": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Message queue name.",
									},
								},
							},
						},
						"dead_letter_policy": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Dead letter queue policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dead_letter_queue": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Dead letter queue.",
									},
									"policy": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "Dead letter queue policy.",
									},
									"max_time_to_live": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "Maximum period in seconds before an unconsumed message expires, which is required if `Policy` is 1. Value range: 300-43200. This value should be smaller than `MsgRetentionSeconds` (maximum message retention period).",
									},
									"max_receive_count": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "Maximum number of receipts.",
									},
								},
							},
						},
						"transaction_policy": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Transaction message policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"first_query_interval": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "First lookback time.",
									},
									"max_query_count": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "Maximum number of queries.",
									},
								},
							},
						},
						"create_uin": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Creator `Uin`.",
						},
						"tags": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Associated tag.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Value of the tag key.",
									},
									"tag_value": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Value of the tag value.",
									},
								},
							},
						},
						"trace": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "Message trace. true: enabled; false: not enabled.",
						},
						"tenant_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Tenant ID.",
						},
						"namespace_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Namespace name.",
						},
						"status": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Cluster status. `0`: creating; `1`: normal; `2`: terminating; `3`: deleted; `4`: isolated; `5`: creation failed; `6`: deletion failed.",
						},
						"max_unacked_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The maximum number of unacknowledged messages.",
						},
						"max_msg_backlog_size": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum size of heaped messages in bytes.",
						},
						"retention_size_in_mb": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Queue storage space configured for message rewind. Value range: 1,024-10,240 MB (if message rewind is enabled). The value `0` indicates that message rewind is not enabled.",
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

func dataSourceTencentCloudTcmqQueueRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcmq_queue.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	if v, ok := d.GetOk("queue_name"); ok {
		paramMap["queue_name"] = v.(string)
	}

	if v, ok := d.GetOk("queue_name_list"); ok {
		queueNameListSet := v.(*schema.Set).List()
		queueNameList := make([]string, 0)
		for i := range queueNameListSet {
			queueName := queueNameListSet[i].(string)
			queueNameList = append(queueNameList, queueName)
		}
		paramMap["queue_name_list"] = queueNameList
	}

	if v, _ := d.GetOk("is_tag_filter"); v != nil {
		paramMap["is_tag_filter"] = v.(bool)
	}

	if v, ok := d.GetOk("filters"); ok {
		filters := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			filter := item.(map[string]interface{})
			name := filter["name"].(string)
			values := make([]string, 0)
			values = append(values, filter["values"].([]string)...)
			filters = append(filters, map[string]interface{}{
				"name":   name,
				"values": values,
			})
		}
		paramMap["fileters"] = filters

	}

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var queueList []*tcmq.CmqQueue
	queueNames := make([]string, 0)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcmqQueueByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		queueList = result
		return nil
	})
	if err != nil {
		return err
	}
	result := make([]map[string]interface{}, 0)
	for _, queue := range queueList {
		queueNames = append(queueNames, *queue.QueueName)
		queueItem := make(map[string]interface{})
		if queue.QueueId != nil {
			queueItem["queue_id"] = *queue.QueueId
		}
		if queue.QueueName != nil {
			queueItem["queue_name"] = *queue.QueueName
		}
		if queue.Qps != nil {
			queueItem["qps"] = *queue.Qps
		}
		if queue.Bps != nil {
			queueItem["bps"] = *queue.Bps
		}
		if queue.MaxDelaySeconds != nil {
			queueItem["max_delay_seconds"] = *queue.MaxDelaySeconds
		}
		if queue.MaxMsgHeapNum != nil {
			queueItem["max_msg_heap_num"] = *queue.MaxMsgHeapNum
		}
		if queue.PollingWaitSeconds != nil {
			queueItem["polling_wait_seconds"] = *queue.PollingWaitSeconds
		}
		if queue.MsgRetentionSeconds != nil {
			queueItem["msg_retention_seconds"] = *queue.MsgRetentionSeconds
		}
		if queue.VisibilityTimeout != nil {
			queueItem["visibility_timeout"] = *queue.VisibilityTimeout
		}
		if queue.MaxMsgSize != nil {
			queueItem["max_msg_size"] = *queue.MaxMsgSize
		}
		if queue.RewindSeconds != nil {
			queueItem["rewind_seconds"] = *queue.RewindSeconds
		}
		if queue.CreateTime != nil {
			queueItem["create_time"] = *queue.CreateTime
		}
		if queue.LastModifyTime != nil {
			queueItem["last_modify_time"] = *queue.LastModifyTime
		}
		if queue.ActiveMsgNum != nil {
			queueItem["active_msg_num"] = *queue.ActiveMsgNum
		}
		if queue.InactiveMsgNum != nil {
			queueItem["inactive_msg_num"] = *queue.InactiveMsgNum
		}
		if queue.DelayMsgNum != nil {
			queueItem["delay_msg_num"] = *queue.DelayMsgNum
		}
		if queue.RewindMsgNum != nil {
			queueItem["rewind_msg_num"] = *queue.RewindMsgNum
		}
		if queue.MinMsgTime != nil {
			queueItem["min_msg_time"] = *queue.MinMsgTime
		}
		if queue.Transaction != nil {
			queueItem["transaction"] = *queue.Transaction
		}
		if queue.CreateUin != nil {
			queueItem["create_uin"] = *queue.CreateUin
		}
		if queue.Trace != nil {
			queueItem["trace"] = *queue.Trace
		}
		if queue.TenantId != nil {
			queueItem["tenant_id"] = *queue.TenantId
		}
		if queue.NamespaceName != nil {
			queueItem["namespace_name"] = *queue.NamespaceName
		}
		if queue.Status != nil {
			queueItem["status"] = *queue.Status
		}
		if queue.MaxUnackedMsgNum != nil {
			queueItem["max_unacked_msg_num"] = *queue.MaxUnackedMsgNum
		}
		if queue.MaxMsgBacklogSize != nil {
			queueItem["max_msg_backlog_size"] = *queue.MaxMsgBacklogSize
		}
		if queue.RetentionSizeInMB != nil {
			queueItem["retention_size_in_mb"] = *queue.RetentionSizeInMB
		}
		if queue.DeadLetterSource != nil {
			deadLetterSource := make([]map[string]interface{}, 0)
			for _, item := range queue.DeadLetterSource {
				deadLetterSource = append(deadLetterSource, map[string]interface{}{
					"queue_id":   item.QueueId,
					"queue_name": item.QueueName,
				})
			}
			queueItem["dead_letter_source"] = deadLetterSource
		}
		if queue.Tags != nil {
			tags := make([]map[string]interface{}, 0)
			for _, item := range queue.Tags {
				tags = append(tags, map[string]interface{}{
					"tag_key":   item.TagKey,
					"tag_value": item.TagValue,
				})
			}
			queueItem["tags"] = tags
		}
		if queue.DeadLetterPolicy != nil {
			deadLetterPolicy := make(map[string]interface{})
			deadLetterPolicy["dead_letter_queue"] = queue.DeadLetterPolicy.DeadLetterQueue
			deadLetterPolicy["policy"] = queue.DeadLetterPolicy.Policy
			deadLetterPolicy["max_time_to_live"] = queue.DeadLetterPolicy.MaxTimeToLive
			deadLetterPolicy["max_receive_count"] = queue.DeadLetterPolicy.MaxReceiveCount
			queueItem["dead_letter_policy"] = deadLetterPolicy
		}
		if queue.TransactionPolicy != nil {
			transactionPolicy := make(map[string]interface{})
			transactionPolicy["first_query_interval"] = queue.TransactionPolicy.FirstQueryInterval
			transactionPolicy["max_query_count"] = queue.TransactionPolicy.MaxQueryCount
			queueItem["transaction_policy"] = transactionPolicy
		}

		result = append(result, queueItem)
	}
	d.SetId(helper.DataResourceIdsHash(queueNames))
	_ = d.Set("queue_list", result)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}

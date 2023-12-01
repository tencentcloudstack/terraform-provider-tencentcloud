/*
Use this data source to query detailed information of tdmq message

Example Usage

```hcl
data "tencentcloud_tdmq_rocketmq_messages" "message" {
  cluster_id     = "rocketmq-rkrbm52djmro"
  environment_id = "keep_ns"
  topic_name     = "keep-topic"
  msg_id         = "A9FE8D0567FE15DB97425FC08EEF0000"
  query_dlq_msg  = false
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

func dataSourceTencentCloudTdmqRocketmqMessages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqMessageRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster id.",
			},
			"environment_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment.",
			},
			"topic_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic, groupId is passed when querying dead letters.",
			},
			"msg_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Message ID.",
			},
			"query_dlq_msg": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "The value is true when querying dead letters, only valid for Rocketmq.",
			},
			"body": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Message body.",
			},
			"properties": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Detailed parameters.",
			},
			"produce_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Production time.",
			},
			"producer_addr": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Producer address.",
			},
			"message_tracks": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Consumer Group ConsumptionNote: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "consumer group.",
						},
						"consume_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "consumption status.",
						},
						"track_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "message track type.",
						},
						"exception_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Exception informationNote: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},
			"show_topic_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The topic name displayed on the details pageNote: This field may return null, indicating that no valid value can be obtained.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqMessageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_rocketmq_messages.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
		msgDetail *tdmq.DescribeRocketMQMsgResponseParams
		clusterId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("environment_id"); ok {
		paramMap["EnvironmentId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		paramMap["TopicName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("msg_id"); ok {
		paramMap["MsgId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("query_dlq_msg"); v != nil {
		paramMap["QueryDlqMsg"] = helper.Bool(v.(bool))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqMessageByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		msgDetail = result
		return nil
	})

	if err != nil {
		return err
	}

	if msgDetail.Body != nil {
		_ = d.Set("body", msgDetail.Body)
	}

	if msgDetail.Properties != nil {
		_ = d.Set("properties", msgDetail.Properties)
	}

	if msgDetail.ProduceTime != nil {
		_ = d.Set("produce_time", msgDetail.ProduceTime)
	}

	if msgDetail.ProducerAddr != nil {
		_ = d.Set("producer_addr", msgDetail.ProducerAddr)
	}

	if msgDetail.MessageTracks != nil {
		tmpList := []interface{}{}
		for _, rocketMQMessageTrack := range msgDetail.MessageTracks {
			rocketMQMessageTrackMap := map[string]interface{}{}

			if rocketMQMessageTrack.Group != nil {
				rocketMQMessageTrackMap["group"] = rocketMQMessageTrack.Group
			}

			if rocketMQMessageTrack.ConsumeStatus != nil {
				rocketMQMessageTrackMap["consume_status"] = rocketMQMessageTrack.ConsumeStatus
			}

			if rocketMQMessageTrack.TrackType != nil {
				rocketMQMessageTrackMap["track_type"] = rocketMQMessageTrack.TrackType
			}

			if rocketMQMessageTrack.ExceptionDesc != nil {
				rocketMQMessageTrackMap["exception_desc"] = rocketMQMessageTrack.ExceptionDesc
			}

			tmpList = append(tmpList, rocketMQMessageTrackMap)
		}

		_ = d.Set("message_tracks", tmpList)
	}

	if msgDetail.ShowTopicName != nil {
		_ = d.Set("show_topic_name", msgDetail.ShowTopicName)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

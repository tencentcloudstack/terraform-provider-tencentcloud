/*
Use this data source to query detailed information of ckafka group_offsets

Example Usage

```hcl
data "tencentcloud_ckafka_group_offsets" "group_offsets" {
  instance_id = "ckafka-xxxxxx"
  group = "xxxxxx"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCkafkaGroupOffsets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaGroupOffsetsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"group": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Kafka consumer group name.",
			},

			"topics": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "An array of topic names subscribed by the group, if there is no such array, it means all topic information under the specified group.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "fuzzy match topicName.",
			},

			"topic_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The topic array, where each element is a json object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "topicName.",
						},
						"partitions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "he topic partition array, where each element is a json object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"partition": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "topic partitionId.",
									},
									"offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The offset of the position.",
									},
									"metadata": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "When consumers submit messages, they can pass in metadata for other purposes. Currently, it is usually an empty string.",
									},
									"error_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ErrorCode.",
									},
									"log_end_offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The latest offset of the current partition.",
									},
									"lag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of unconsumed messages.",
									},
								},
							},
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

func dataSourceTencentCloudCkafkaGroupOffsetsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_group_offsets.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group"); ok {
		paramMap["group"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topics"); ok {
		topicsSet := v.(*schema.Set).List()
		paramMap["topics"] = helper.InterfacesStringsPoint(topicsSet)
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["search_word"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var groupOffsetTopics []*ckafka.GroupOffsetTopic

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaGroupOffsetsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		groupOffsetTopics = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groupOffsetTopics))
	groupOffsetResponseMap := map[string]interface{}{}

	topicList := []interface{}{}
	for _, topic := range groupOffsetTopics {
		topicMap := map[string]interface{}{}

		if topic.Topic != nil {
			topicMap["topic"] = topic.Topic
			ids = append(ids, *topic.Topic)

		}

		if topic.Partitions != nil {
			partitionsList := []interface{}{}
			for _, partitions := range topic.Partitions {
				partitionsMap := map[string]interface{}{}

				if partitions.Partition != nil {
					partitionsMap["partition"] = partitions.Partition
				}

				if partitions.Offset != nil {
					partitionsMap["offset"] = partitions.Offset
				}

				if partitions.Metadata != nil {
					partitionsMap["metadata"] = partitions.Metadata
				}

				if partitions.ErrorCode != nil {
					partitionsMap["error_code"] = partitions.ErrorCode
				}

				if partitions.LogEndOffset != nil {
					partitionsMap["log_end_offset"] = partitions.LogEndOffset
				}

				if partitions.Lag != nil {
					partitionsMap["lag"] = partitions.Lag
				}

				partitionsList = append(partitionsList, partitionsMap)
			}

			topicMap["partitions"] = partitionsList
		}

		topicList = append(topicList, topicMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("topic_list", topicList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), groupOffsetResponseMap); e != nil {
			return e
		}
	}
	return nil
}

package ckafka

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaDatahubGroupOffsets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaDatahubGroupOffsetsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "topic name that the task subscribe.",
			},

			"group": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Kafka consumer group.",
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
							Description: "topic name.",
						},
						"partitions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The topic partition array, where each element is a json object.",
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
										Description: "consumer offset.",
									},
									"metadata": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Usually an empty string.",
									},
									"error_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Error Code.",
									},
									"log_end_offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "partition Log End Offset.",
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

func dataSourceTencentCloudCkafkaDatahubGroupOffsetsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_datahub_group_offsets.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		paramMap["name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group"); ok {
		paramMap["group"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["search_word"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result []*ckafka.GroupOffsetTopic

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		groupOffsetTopics, e := service.DescribeCkafkaDatahubGroupOffsetsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = groupOffsetTopics
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	topicList := make([]map[string]interface{}, 0, len(result))
	for _, topic := range result {
		topicMap := make(map[string]interface{})

		if topic.Topic != nil {
			topicMap["topic"] = topic.Topic
			ids = append(ids, *topic.Topic)

		}

		if topic.Partitions != nil {
			partitionsList := make([]map[string]interface{}, 0)
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
		if e := tccommon.WriteToFile(output.(string), topicList); e != nil {
			return e
		}
	}
	return nil
}

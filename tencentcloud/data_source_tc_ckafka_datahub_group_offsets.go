/*
Use this data source to query detailed information of ckafka datahub_group_offsets

Example Usage

```hcl
data "tencentcloud_ckafka_datahub_group_offsets" "datahub_group_offsets" {
  name = "1300xxxx-topicName"
  group = "datahub-task-lzp7qb7e"
  search_word = ""
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

func dataSourceTencentCloudCkafkaDatahubGroupOffsets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaDatahubGroupOffsetsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic name that the task subscribe.",
			},

			"group": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Kafka consumer group.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy match topicName.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of matching results.",
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
										Description: "Topic name.",
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
													Description: "Topic partitionId.",
												},
												"offset": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Consumer offset.",
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
													Description: "Partition Log End Offset.",
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
	defer logElapsed("data_source.tencentcloud_ckafka_datahub_group_offsets.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		paramMap["Name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group"); ok {
		paramMap["Group"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.GroupOffsetResponse

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaDatahubGroupOffsetsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		groupOffsetResponseMap := map[string]interface{}{}

		if result.TotalCount != nil {
			groupOffsetResponseMap["total_count"] = result.TotalCount
		}

		if result.TopicList != nil {
			topicListList := []interface{}{}
			for _, topicList := range result.TopicList {
				topicListMap := map[string]interface{}{}

				if topicList.Topic != nil {
					topicListMap["topic"] = topicList.Topic
				}

				if topicList.Partitions != nil {
					partitionsList := []interface{}{}
					for _, partitions := range topicList.Partitions {
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

					topicListMap["partitions"] = []interface{}{partitionsList}
				}

				topicListList = append(topicListList, topicListMap)
			}

			groupOffsetResponseMap["topic_list"] = []interface{}{topicListList}
		}

		ids = append(ids, *result.Name)
		_ = d.Set("result", groupOffsetResponseMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), groupOffsetResponseMap); e != nil {
			return e
		}
	}
	return nil
}

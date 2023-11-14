/*
Use this data source to query detailed information of ckafka datahub_topic

Example Usage

```hcl
data "tencentcloud_ckafka_datahub_topic" "datahub_topic" {
  search_word = "topicName"
  offset = 0
  limit = 20
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

func dataSourceTencentCloudCkafkaDatahubTopic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaDatahubTopicRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query key word.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The offset position of this query, the default is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of results returned this time, the default is 50, and the maximum value is 50.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Topic list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total Count.",
						},
						"topic_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Topic list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name.",
									},
									"topic_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Topic name.",
									},
									"topic_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Topic Id.",
									},
									"partition_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of partitions.",
									},
									"retention_ms": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Expiration.",
									},
									"note": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Remark.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Status, 1 in use, 2 in deletion.",
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

func dataSourceTencentCloudCkafkaDatahubTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_datahub_topic.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.DescribeDatahubTopicsResp

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaDatahubTopicByFilter(ctx, paramMap)
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
		describeDatahubTopicsRespMap := map[string]interface{}{}

		if result.TotalCount != nil {
			describeDatahubTopicsRespMap["total_count"] = result.TotalCount
		}

		if result.TopicList != nil {
			topicListList := []interface{}{}
			for _, topicList := range result.TopicList {
				topicListMap := map[string]interface{}{}

				if topicList.Name != nil {
					topicListMap["name"] = topicList.Name
				}

				if topicList.TopicName != nil {
					topicListMap["topic_name"] = topicList.TopicName
				}

				if topicList.TopicId != nil {
					topicListMap["topic_id"] = topicList.TopicId
				}

				if topicList.PartitionNum != nil {
					topicListMap["partition_num"] = topicList.PartitionNum
				}

				if topicList.RetentionMs != nil {
					topicListMap["retention_ms"] = topicList.RetentionMs
				}

				if topicList.Note != nil {
					topicListMap["note"] = topicList.Note
				}

				if topicList.Status != nil {
					topicListMap["status"] = topicList.Status
				}

				topicListList = append(topicListList, topicListMap)
			}

			describeDatahubTopicsRespMap["topic_list"] = []interface{}{topicListList}
		}

		ids = append(ids, *result.TopicName)
		_ = d.Set("result", describeDatahubTopicsRespMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), describeDatahubTopicsRespMap); e != nil {
			return e
		}
	}
	return nil
}

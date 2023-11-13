/*
Use this data source to query detailed information of ckafka topic_sync_replica

Example Usage

```hcl
data "tencentcloud_ckafka_topic_sync_replica" "topic_sync_replica" {
  instance_id = "InstanceId"
  topic_name = "TopicName"
  out_of_sync_replica_only = true
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

func dataSourceTencentCloudCkafkaTopicSyncReplica() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaTopicSyncReplicaRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"topic_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "TopicName.",
			},

			"out_of_sync_replica_only": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Filter only unsynced replicas.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Return topic copy details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_in_sync_replica_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Topic details and copy collection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"partition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Partition name.",
									},
									"leader": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Leader Id.",
									},
									"replica": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Replica set.",
									},
									"in_sync_replica": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ISR.",
									},
									"begin_offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "BeginOffset.",
									},
									"end_offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "EndOffset.",
									},
									"message_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Message Count.",
									},
									"out_of_sync_replica": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Out Of Sync Replica.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number.",
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

func dataSourceTencentCloudCkafkaTopicSyncReplicaRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_topic_sync_replica.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		paramMap["TopicName"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("out_of_sync_replica_only"); v != nil {
		paramMap["OutOfSyncReplicaOnly"] = helper.Bool(v.(bool))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.TopicInSyncReplicaResult

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaTopicSyncReplicaByFilter(ctx, paramMap)
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
		topicInSyncReplicaResultMap := map[string]interface{}{}

		if result.TopicInSyncReplicaList != nil {
			topicInSyncReplicaListList := []interface{}{}
			for _, topicInSyncReplicaList := range result.TopicInSyncReplicaList {
				topicInSyncReplicaListMap := map[string]interface{}{}

				if topicInSyncReplicaList.Partition != nil {
					topicInSyncReplicaListMap["partition"] = topicInSyncReplicaList.Partition
				}

				if topicInSyncReplicaList.Leader != nil {
					topicInSyncReplicaListMap["leader"] = topicInSyncReplicaList.Leader
				}

				if topicInSyncReplicaList.Replica != nil {
					topicInSyncReplicaListMap["replica"] = topicInSyncReplicaList.Replica
				}

				if topicInSyncReplicaList.InSyncReplica != nil {
					topicInSyncReplicaListMap["in_sync_replica"] = topicInSyncReplicaList.InSyncReplica
				}

				if topicInSyncReplicaList.BeginOffset != nil {
					topicInSyncReplicaListMap["begin_offset"] = topicInSyncReplicaList.BeginOffset
				}

				if topicInSyncReplicaList.EndOffset != nil {
					topicInSyncReplicaListMap["end_offset"] = topicInSyncReplicaList.EndOffset
				}

				if topicInSyncReplicaList.MessageCount != nil {
					topicInSyncReplicaListMap["message_count"] = topicInSyncReplicaList.MessageCount
				}

				if topicInSyncReplicaList.OutOfSyncReplica != nil {
					topicInSyncReplicaListMap["out_of_sync_replica"] = topicInSyncReplicaList.OutOfSyncReplica
				}

				topicInSyncReplicaListList = append(topicInSyncReplicaListList, topicInSyncReplicaListMap)
			}

			topicInSyncReplicaResultMap["topic_in_sync_replica_list"] = []interface{}{topicInSyncReplicaListList}
		}

		if result.TotalCount != nil {
			topicInSyncReplicaResultMap["total_count"] = result.TotalCount
		}

		ids = append(ids, *result.InstanceId)
		_ = d.Set("result", topicInSyncReplicaResultMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), topicInSyncReplicaResultMap); e != nil {
			return e
		}
	}
	return nil
}

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

			"topic_in_sync_replica_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Topic details and copy collection.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "partition name.",
						},
						"leader": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Leader Id.",
						},
						"replica": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "replica set.",
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
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		paramMap["topic_name"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("out_of_sync_replica_only"); v != nil {
		paramMap["out_of_sync_replica_only"] = helper.Bool(v.(bool))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.TopicInSyncReplicaInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		topicInSyncReplicaInfos, e := service.DescribeCkafkaTopicSyncReplicaByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = topicInSyncReplicaInfos
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	topicInSyncReplicaList := make([]interface{}, 0, len(result))

	for _, topicInSyncReplica := range result {
		topicInSyncReplicaMap := map[string]interface{}{}

		if topicInSyncReplica.Partition != nil {
			ids = append(ids, *topicInSyncReplica.Partition)
			topicInSyncReplicaMap["partition"] = topicInSyncReplica.Partition
		}

		if topicInSyncReplica.Leader != nil {
			topicInSyncReplicaMap["leader"] = topicInSyncReplica.Leader
		}

		if topicInSyncReplica.Replica != nil {
			topicInSyncReplicaMap["replica"] = topicInSyncReplica.Replica
		}

		if topicInSyncReplica.InSyncReplica != nil {
			topicInSyncReplicaMap["in_sync_replica"] = topicInSyncReplica.InSyncReplica
		}

		if topicInSyncReplica.BeginOffset != nil {
			topicInSyncReplicaMap["begin_offset"] = topicInSyncReplica.BeginOffset
		}

		if topicInSyncReplica.EndOffset != nil {
			topicInSyncReplicaMap["end_offset"] = topicInSyncReplica.EndOffset
		}

		if topicInSyncReplica.MessageCount != nil {
			topicInSyncReplicaMap["message_count"] = topicInSyncReplica.MessageCount
		}

		if topicInSyncReplica.OutOfSyncReplica != nil {
			topicInSyncReplicaMap["out_of_sync_replica"] = topicInSyncReplica.OutOfSyncReplica
		}

		topicInSyncReplicaList = append(topicInSyncReplicaList, topicInSyncReplicaMap)
	}

	_ = d.Set("topic_in_sync_replica_list", topicInSyncReplicaList)

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), topicInSyncReplicaList); e != nil {
			return e
		}
	}
	return nil
}

/*
Use this data source to query detailed information of ckafka topic_subscribe_group

Example Usage

```hcl
data "tencentcloud_ckafka_topic_subscribe_group" "topic_subscribe_group" {
  instance_id = "ckafka-xxxxxx"
  topic_name = "xxxxxx"
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

func dataSourceTencentCloudCkafkaTopicSubscribeGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaTopicSubscribeGroupRead,
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

			"groups_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Consumer group information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error code, normally 0.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group state description (commonly Empty, Stable, and Dead states): Dead: The consumption group does not exist Empty: The consumption group does not currently have any consumer subscriptions PreparingRebalance: The consumption group is in the rebalance state CompletingRebalance: The consumption group is in the rebalance state Stable: Each consumer in the consumption group has joined and is in a stable state.",
						},
						"protocol_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol type selected by the consumption group is normally the consumer, but some systems use their own protocol, such as kafka-connect, which uses connect. Only the standard consumer protocol, this interface knows the format of the specific allocation method, and can analyze the specific partition allocation.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Common consumer partition allocation algorithms are as follows (the default option for Kafka consumer SDK is range) range|roundrobin| sticky.",
						},
						"members": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "This array contains information only if state is Stable and protocol_type is consumer.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"member_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID that the coordinator generated for consumer.",
									},
									"client_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The client.id information set by the client consumer SDK itself.",
									},
									"client_host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Generally store the customer&#39;s IP address.",
									},
									"assignment": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Stores the partition information assigned to the consumer.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "assignment version information.",
												},
												"topics": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "topic list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "topic name.",
															},
															"partitions": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeInt,
																},
																Computed:    true,
																Description: "partition list.",
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
						"group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kafka consumer group.",
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

func dataSourceTencentCloudCkafkaTopicSubscribeGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_topic_subscribe_group.read")()
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

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.GroupInfoResponse

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		groupInfos, e := service.DescribeCkafkaTopicSubscribeGroupByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = groupInfos
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))

	groupsInfoList := []interface{}{}
	for _, groupsInfo := range result {
		groupsInfoMap := map[string]interface{}{}

		if groupsInfo.ErrorCode != nil {
			groupsInfoMap["error_code"] = groupsInfo.ErrorCode
		}

		if groupsInfo.State != nil {
			groupsInfoMap["state"] = groupsInfo.State
		}

		if groupsInfo.ProtocolType != nil {
			groupsInfoMap["protocol_type"] = groupsInfo.ProtocolType
		}

		if groupsInfo.Protocol != nil {
			groupsInfoMap["protocol"] = groupsInfo.Protocol
		}

		if groupsInfo.Members != nil {
			membersList := []interface{}{}
			for _, members := range groupsInfo.Members {
				membersMap := map[string]interface{}{}

				if members.MemberId != nil {
					membersMap["member_id"] = members.MemberId
				}

				if members.ClientId != nil {
					membersMap["client_id"] = members.ClientId
				}

				if members.ClientHost != nil {
					membersMap["client_host"] = members.ClientHost
				}

				if members.Assignment != nil {
					assignmentMap := map[string]interface{}{}

					if members.Assignment.Version != nil {
						assignmentMap["version"] = members.Assignment.Version
					}

					if members.Assignment.Topics != nil {
						topicsList := []interface{}{}
						for _, topics := range members.Assignment.Topics {
							topicsMap := map[string]interface{}{}

							if topics.Topic != nil {
								topicsMap["topic"] = topics.Topic
							}

							if topics.Partitions != nil {
								topicsMap["partitions"] = topics.Partitions
							}

							topicsList = append(topicsList, topicsMap)
						}

						assignmentMap["topics"] = []interface{}{topicsList}
					}

					membersMap["assignment"] = []interface{}{assignmentMap}
				}

				membersList = append(membersList, membersMap)
			}

			groupsInfoMap["members"] = membersList
		}

		if groupsInfo.Group != nil {
			ids = append(ids, *groupsInfo.Group)

			groupsInfoMap["group"] = groupsInfo.Group
		}

		groupsInfoList = append(groupsInfoList, groupsInfoMap)
	}

	_ = d.Set("groups_info", groupsInfoList)

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), groupsInfoList); e != nil {
			return e
		}
	}
	return nil
}

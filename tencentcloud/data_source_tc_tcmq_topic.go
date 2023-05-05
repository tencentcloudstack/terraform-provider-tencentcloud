/*
Use this data source to query detailed information of tcmq topic

Example Usage

```hcl
data "tencentcloud_tcmq_topic" "topic" {
  topic_name = "topic_name"
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

func dataSourceTencentCloudTcmqTopic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcmqTopicRead,
		Schema: map[string]*schema.Schema{
			"offset": {
				Default:     0,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Starting position of the list of topics to be returned on the current page in case of paginated return. If a value is entered, limit is required. If this parameter is left empty, 0 will be used by default.",
			},

			"limit": {
				Default:     20,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of topics to be returned per page in case of paginated return. If this parameter is not passed in, 20 will be used by default. Maximum value: 50.",
			},

			"topic_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search by TopicName.",
			},

			"topic_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by CMQ topic name.",
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

			"topic_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Topic list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Topic ID.",
						},
						"topic_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Topic name.",
						},
						"msg_retention_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum lifecycle of message in topic. After the period specified by this parameter has elapsed since a message is sent to the topic, the message will be deleted no matter whether it has been successfully pushed to the user. This parameter is measured in seconds and defaulted to one day (86,400 seconds), which cannot be modified.",
						},
						"max_msg_size": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum message size, which ranges from 1,024 to 1,048,576 bytes (i.e., 1-1,024 KB). The default value is 65,536.",
						},
						"qps": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Number of messages published per second.",
						},
						"filter_type": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Filtering policy selected when a subscription is created: If `filterType` is 1, `FilterTag` will be used for filtering. If `filterType` is 2, `BindingKey` will be used for filtering.",
						},
						"create_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Topic creation time. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"last_modify_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Time when the topic attribute is last modified. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"msg_count": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Number of current messages in the topic (number of retained messages).",
						},
						"create_uin": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Creator `Uin`. The `resource` field for CAM authentication is composed of this field.",
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
						"broker_type": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Valid values: `0` (Pulsar), `1` (RocketMQ).",
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

func dataSourceTencentCloudTcmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcmq_topic.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	if v, ok := d.GetOk("topic_name"); ok {
		paramMap["topic_name"] = v.(string)
	}

	if v, ok := d.GetOk("topic_name_list"); ok {
		topicNameListSet := v.(*schema.Set).List()
		topicNameList := make([]string, 0)
		for i := range topicNameListSet {
			topicName := topicNameListSet[i].(string)
			topicNameList = append(topicNameList, topicName)
		}
		paramMap["topic_name_list"] = topicNameList

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
		paramMap["filters"] = filters

	}

	service := TcmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var topicList []*tdmq.CmqTopic
	topicNames := make([]string, 0)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcmqTopicByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		topicList = result
		return nil
	})
	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)
	for _, topic := range topicList {
		topicNames = append(topicNames, *topic.TopicName)
		topicItem := make(map[string]interface{})
		if topic.TenantId != nil {
			topicItem["topic_id"] = *topic.TenantId
		}
		if topic.TopicName != nil {
			topicItem["topic_name"] = *topic.TopicName
		}
		if topic.MsgRetentionSeconds != nil {
			topicItem["msg_retention_seconds"] = *topic.MsgRetentionSeconds
		}
		if topic.MaxMsgSize != nil {
			topicItem["max_msg_size"] = *topic.MaxMsgSize
		}
		if topic.Qps != nil {
			topicItem["qps"] = *topic.Qps
		}
		if topic.FilterType != nil {
			topicItem["filter_type"] = *topic.FilterType
		}
		if topic.CreateTime != nil {
			topicItem["create_time"] = *topic.CreateTime
		}
		if topic.LastModifyTime != nil {
			topicItem["last_modify_time"] = *topic.LastModifyTime
		}
		if topic.MsgCount != nil {
			topicItem["msg_count"] = *topic.MsgCount
		}
		if topic.CreateUin != nil {
			topicItem["create_uin"] = *topic.CreateUin
		}
		if topic.Trace != nil {
			topicItem["trace"] = *topic.Trace
		}
		if topic.TenantId != nil {
			topicItem["tenant_id"] = *topic.TenantId
		}
		if topic.NamespaceName != nil {
			topicItem["namespace_name"] = *topic.NamespaceName
		}
		if topic.Status != nil {
			topicItem["status"] = *topic.Status
		}
		if topic.BrokerType != nil {
			topicItem["broker_type"] = *topic.BrokerType
		}

		if topic.Tags != nil {
			tags := make([]map[string]interface{}, 0)
			for _, item := range topic.Tags {
				tags = append(tags, map[string]interface{}{
					"tag_key":   *item.TagKey,
					"tag_value": *item.TagValue,
				})
			}
			topicItem["tags"] = tags
		}

		result = append(result, topicItem)
	}
	d.SetId(helper.DataResourceIdsHash(topicNames))
	_ = d.Set("topic_list", result)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}

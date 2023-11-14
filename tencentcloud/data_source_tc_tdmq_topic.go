/*
Use this data source to query detailed information of tdmq topic

Example Usage

```hcl
data "tencentcloud_tdmq_topic" "topic" {
  offset = 0
  limit = 20
  topic_name = "topic_name"
  topic_name_list = &lt;nil&gt;
  is_tag_filter = &lt;nil&gt;
  filters {
		name = "tag"
		values = &lt;nil&gt;

  }
  topic_list {
		topic_id = &lt;nil&gt;
		topic_name = &lt;nil&gt;
		msg_retention_seconds = &lt;nil&gt;
		max_msg_size = &lt;nil&gt;
		qps = &lt;nil&gt;
		filter_type = &lt;nil&gt;
		create_time = &lt;nil&gt;
		last_modify_time = &lt;nil&gt;
		msg_count = &lt;nil&gt;
		create_uin = &lt;nil&gt;
		tags {
			tag_key = &lt;nil&gt;
			tag_value = &lt;nil&gt;
		}
		trace = &lt;nil&gt;
		tenant_id = &lt;nil&gt;
		namespace_name = &lt;nil&gt;
		status = &lt;nil&gt;
		broker_type = &lt;nil&gt;

  }
  total_count = &lt;nil&gt;
  request_id = &lt;nil&gt;
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

func dataSourceTencentCloudTdmqTopic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqTopicRead,
		Schema: map[string]*schema.Schema{
			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Starting position of the list of topics to be returned on the current page in case of paginated return. If a value is entered, limit is required. If this parameter is left empty, 0 will be used by default.",
			},

			"limit": {
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
				Description: "Filter. Currently, you can filter by tag. The tag name must be prefixed with “tag:”, such as “tag: owner”, “tag: environment”, or “tag: business”.",
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
				Type:        schema.TypeList,
				Description: "Topic list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Type:        schema.TypeString,
							Description: "Topic ID.",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Description: "Topic name.",
						},
						"msg_retention_seconds": {
							Type:        schema.TypeInt,
							Description: "Maximum lifecycle of message in topic. After the period specified by this parameter has elapsed since a message is sent to the topic, the message will be deleted no matter whether it has been successfully pushed to the user. This parameter is measured in seconds and defaulted to one day (86,400 seconds), which cannot be modified.",
						},
						"max_msg_size": {
							Type:        schema.TypeInt,
							Description: "Maximum message size, which ranges from 1,024 to 1,048,576 bytes (i.e., 1-1,024 KB). The default value is 65,536.",
						},
						"qps": {
							Type:        schema.TypeInt,
							Description: "Number of messages published per second.",
						},
						"filter_type": {
							Type:        schema.TypeInt,
							Description: "Filtering policy selected when a subscription is created: If `filterType` is 1, `FilterTag` will be used for filtering. If `filterType` is 2, `BindingKey` will be used for filtering.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Description: "Topic creation time. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"last_modify_time": {
							Type:        schema.TypeInt,
							Description: "Time when the topic attribute is last modified. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"msg_count": {
							Type:        schema.TypeInt,
							Description: "Number of current messages in the topic (number of retained messages).",
						},
						"create_uin": {
							Type:        schema.TypeInt,
							Description: "Creator `Uin`. The `resource` field for CAM authentication is composed of this field.",
						},
						"tags": {
							Type:        schema.TypeList,
							Description: "Associated tag.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Description: "Value of the tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Description: "Value of the tag value.",
									},
								},
							},
						},
						"trace": {
							Type:        schema.TypeBool,
							Description: "Message trace. true: enabled; false: not enabled.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Description: "Tenant ID.",
						},
						"namespace_name": {
							Type:        schema.TypeString,
							Description: "Namespace name.",
						},
						"status": {
							Type:        schema.TypeInt,
							Description: "Cluster status. `0`: creating; `1`: normal; `2`: terminating; `3`: deleted; `4`: isolated; `5`: creation failed; `6`: deletion failed.",
						},
						"broker_type": {
							Type:        schema.TypeInt,
							Description: "Valid values: `0` (Pulsar), `1` (RocketMQ).",
						},
					},
				},
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "The total number of topics.",
			},

			"request_id": {
				Type:        schema.TypeString,
				Description: "The unique request ID, which is returned for each request. RequestId is required for locating a problem.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_topic.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		paramMap["TopicName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name_list"); ok {
		topicNameListSet := v.(*schema.Set).List()
		paramMap["TopicNameList"] = helper.InterfacesStringsPoint(topicNameListSet)
	}

	if v, _ := d.GetOk("is_tag_filter"); v != nil {
		paramMap["IsTagFilter"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tdmq.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tdmq.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("topic_list"); ok {
		topicListSet := v.([]interface{})
		tmpSet := make([]*tdmq.CmqTopic, 0, len(topicListSet))

		for _, item := range topicListSet {
			cmqTopic := tdmq.CmqTopic{}
			cmqTopicMap := item.(map[string]interface{})

			if v, ok := cmqTopicMap["topic_id"]; ok {
				cmqTopic.TopicId = helper.String(v.(string))
			}
			if v, ok := cmqTopicMap["topic_name"]; ok {
				cmqTopic.TopicName = helper.String(v.(string))
			}
			if v, ok := cmqTopicMap["msg_retention_seconds"]; ok {
				cmqTopic.MsgRetentionSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := cmqTopicMap["max_msg_size"]; ok {
				cmqTopic.MaxMsgSize = helper.IntUint64(v.(int))
			}
			if v, ok := cmqTopicMap["qps"]; ok {
				cmqTopic.Qps = helper.IntUint64(v.(int))
			}
			if v, ok := cmqTopicMap["filter_type"]; ok {
				cmqTopic.FilterType = helper.IntUint64(v.(int))
			}
			if v, ok := cmqTopicMap["create_time"]; ok {
				cmqTopic.CreateTime = helper.IntUint64(v.(int))
			}
			if v, ok := cmqTopicMap["last_modify_time"]; ok {
				cmqTopic.LastModifyTime = helper.IntUint64(v.(int))
			}
			if v, ok := cmqTopicMap["msg_count"]; ok {
				cmqTopic.MsgCount = helper.IntUint64(v.(int))
			}
			if v, ok := cmqTopicMap["create_uin"]; ok {
				cmqTopic.CreateUin = helper.IntUint64(v.(int))
			}
			if v, ok := cmqTopicMap["tags"]; ok {
				for _, item := range v.([]interface{}) {
					tagsMap := item.(map[string]interface{})
					tag := tdmq.Tag{}
					if v, ok := tagsMap["tag_key"]; ok {
						tag.TagKey = helper.String(v.(string))
					}
					if v, ok := tagsMap["tag_value"]; ok {
						tag.TagValue = helper.String(v.(string))
					}
					cmqTopic.Tags = append(cmqTopic.Tags, &tag)
				}
			}
			if v, ok := cmqTopicMap["trace"]; ok {
				cmqTopic.Trace = helper.Bool(v.(bool))
			}
			if v, ok := cmqTopicMap["tenant_id"]; ok {
				cmqTopic.TenantId = helper.String(v.(string))
			}
			if v, ok := cmqTopicMap["namespace_name"]; ok {
				cmqTopic.NamespaceName = helper.String(v.(string))
			}
			if v, ok := cmqTopicMap["status"]; ok {
				cmqTopic.Status = helper.IntInt64(v.(int))
			}
			if v, ok := cmqTopicMap["broker_type"]; ok {
				cmqTopic.BrokerType = helper.IntInt64(v.(int))
			}
			tmpSet = append(tmpSet, &cmqTopic)
		}
		paramMap["topic_list"] = tmpSet
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("request_id"); ok {
		paramMap["RequestId"] = helper.String(v.(string))
	}

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var topicList []*tdmq.CmqTopic

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqTopicByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		topicList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(topicList))
	tmpList := make([]map[string]interface{}, 0, len(topicList))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

package cls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsTopicsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "<li>topicName: Filter by **log topic name**. Fuzzy match is implemented by default. You can use the `PreciseSearch` parameter to set exact match. Type: String. Required. No. <br><li>logsetName: Filter by **logset name**. Fuzzy match is implemented by default. You can use the `PreciseSearch` parameter to set exact match. Type: String. Required: No. <br><li>topicId: Filter by **log topic ID**. Type: String. Required: No. <br><li>logsetId: Filter by **logset ID**. You can call `DescribeLogsets` to query the list of created logsets or log in to the console to view them. You can also call `CreateLogset` to create a logset. Type: String. Required: No. <br><li>tagKey: Filter by **tag key**. Type: String. Required: No. <br><li>tag:tagKey: Filter by **tag key-value pair**. The `tagKey` should be replaced with a specified tag key, such as `tag:exampleKey`. Type: String. Required: No. <br><li>storageType: Filter by **log topic storage type**. Valid values: `hot` (standard storage) and `cold` (IA storage). Type: String. Required: No. Each request can have up to 10 `Filters` and 100 `Filter.Values`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Value to be filtered.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"precise_search": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Match mode for `Filters` fields.\n- 0: Fuzzy match for `topicName` and `logsetName`. This is the default value.\n- 1: Exact match for `topicName`.\n- 2: Exact match for `logsetName`.\n- 3: Exact match for `topicName` and `logsetName`.",
			},

			"biz_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Topic type\n- 0 (default): Log topic.\n- 1: Metric topic.",
			},

			"topics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Log topic list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logset_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Logset ID.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Topic ID.",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Topic Name.",
						},
						"partition_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of topic partitions.",
						},
						"index": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the topic has indexing enabled (the topic type must be log topic).",
						},
						"assumer_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cloud product identifier. When the topic is created by other cloud products, this field displays the name of the cloud product, such as CDN, TKE.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Creation time.",
						},
						"status": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the topic has log collection enabled. true: collection enabled; false: collection disabled.Log collection is enabled by default when creating a log topic, and this field can be modified by calling ModifyTopic through the SDK.The console currently does not support modifying this parameter.",
						},
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Tag information bound to the topicNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The tag key.\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The tag value.\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"auto_split": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether automatic split is enabled for this topic\nNote: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"max_split_partitions": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximum number of partitions to split into for this topic if automatic split is enabled\nNote: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage type of the topicNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Lifecycle in days. Value range: 1-3600 (3640 indicates permanent retention)\nNote: This field may return `null`, indicating that no valid value was found.",
						},
						"sub_assumer_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cloud product sub-identifier. If the log topic is created by another cloud product, this field returns the name of the cloud product and its log type, such as `TKE-Audit` or `TKE-Event`. Some products only return the cloud product identifier (`AssumerName`), without this field.\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"describes": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Topic description\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"hot_period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Enable log sinking, with the lifecycle of standard storage, where hotPeriod < Period.For standard storage, hotPeriod is used, and for infrequent access storage, it is Period-hotPeriod. (The topic type must be a log topic)HotPeriod=0 indicates that log sinking is not enabled.\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"biz_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Topic type.\n- 0:  log  Topic  \n- 1: Metric Topic\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"is_web_tracking": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Free authentication switch. false: disabled; true: enabled.After enabling, anonymous access to the log topic will be supported for specified operations. For details, please refer to Log Topic (https://intl.cloud.tencent.com/document/product/614/41035?from_cn_redirect=1).Note: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudClsTopicsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cls_topics.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clsv20201016.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := clsv20201016.Filter{}
			if v, ok := filtersMap["key"].(string); ok && v != "" {
				filter.Key = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOkExists("precise_search"); ok {
		paramMap["PreciseSearch"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("biz_type"); ok {
		paramMap["BizType"] = helper.IntUint64(v.(int))
	}

	var respData []*clsv20201016.TopicInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsTopicsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	ids := make([]string, 0, len(respData))
	topicsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, topics := range respData {
			topicsMap := map[string]interface{}{}
			if topics.LogsetId != nil {
				topicsMap["logset_id"] = topics.LogsetId
			}

			if topics.TopicId != nil {
				topicsMap["topic_id"] = topics.TopicId
				ids = append(ids, *topics.TopicId)
			}

			if topics.TopicName != nil {
				topicsMap["topic_name"] = topics.TopicName
			}

			if topics.PartitionCount != nil {
				topicsMap["partition_count"] = topics.PartitionCount
			}

			if topics.Index != nil {
				topicsMap["index"] = topics.Index
			}

			if topics.AssumerName != nil {
				topicsMap["assumer_name"] = topics.AssumerName
			}

			if topics.CreateTime != nil {
				topicsMap["create_time"] = topics.CreateTime
			}

			if topics.Status != nil {
				topicsMap["status"] = topics.Status
			}

			tagsList := make([]map[string]interface{}, 0, len(topics.Tags))
			if topics.Tags != nil {
				for _, tags := range topics.Tags {
					tagsMap := map[string]interface{}{}
					if tags.Key != nil {
						tagsMap["key"] = tags.Key
					}

					if tags.Value != nil {
						tagsMap["value"] = tags.Value
					}

					tagsList = append(tagsList, tagsMap)
				}

				topicsMap["tags"] = tagsList
			}

			if topics.AutoSplit != nil {
				topicsMap["auto_split"] = topics.AutoSplit
			}

			if topics.MaxSplitPartitions != nil {
				topicsMap["max_split_partitions"] = topics.MaxSplitPartitions
			}

			if topics.StorageType != nil {
				topicsMap["storage_type"] = topics.StorageType
			}

			if topics.Period != nil {
				topicsMap["period"] = topics.Period
			}

			if topics.SubAssumerName != nil {
				topicsMap["sub_assumer_name"] = topics.SubAssumerName
			}

			if topics.Describes != nil {
				topicsMap["describes"] = topics.Describes
			}

			if topics.HotPeriod != nil {
				topicsMap["hot_period"] = topics.HotPeriod
			}

			if topics.BizType != nil {
				topicsMap["biz_type"] = topics.BizType
			}

			if topics.IsWebTracking != nil {
				topicsMap["is_web_tracking"] = topics.IsWebTracking
			}

			topicsList = append(topicsList, topicsMap)
		}

		_ = d.Set("topics", topicsList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), topicsList); e != nil {
			return e
		}
	}

	return nil
}

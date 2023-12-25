package trocket

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdmqRocketmqTopic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqRocketmqTopicRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID.",
			},

			"namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace.",
			},

			"filter_type": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Filter by topic type. Valid values: `Normal`, `GlobalOrder`, `PartitionedOrder`, `Transaction`.",
			},

			"filter_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search by topic name. Fuzzy query is supported.",
			},

			"topics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of topic information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Topic name.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Topic name.",
						},
						"partition_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of read/write partitions.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation time in milliseconds.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Update time in milliseconds.",
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

func dataSourceTencentCloudTdmqRocketmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmqRocketmq_topic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["cluster_id"] = v.(string)
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		paramMap["namespace_id"] = v.(string)
	}

	if v, ok := d.GetOk("filter_type"); ok {
		filterTypes := v.(*schema.Set).List()
		filterTypeList := make([]string, 0)
		for _, filterType := range filterTypes {
			filter_type := filterType.(string)
			filterTypeList = append(filterTypeList, filter_type)
		}
		paramMap["filter_type"] = filterTypeList
	}

	if v, ok := d.GetOk("filter_name"); ok {
		paramMap["filter_name"] = v.(string)
	}

	tdmqRocketmqService := TdmqRocketmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var topics []*tdmqRocketmq.RocketMQTopic
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := tdmqRocketmqService.DescribeTdmqRocketmqTopicByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		topics = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read TdmqRocketmq topics failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0)
	topicList := []interface{}{}
	for _, topic := range topics {
		topicMap := map[string]interface{}{}
		ids = append(ids, *topic.Name)
		topicMap["name"] = topic.Name
		if topic.Remark != nil {
			topicMap["remark"] = topic.Remark
		}
		if topic.PartitionNum != nil {
			topicMap["partition_num"] = topic.PartitionNum
		}
		if topic.CreateTime != nil {
			topicMap["create_time"] = topic.CreateTime
		}
		if topic.UpdateTime != nil {
			topicMap["update_time"] = topic.UpdateTime
		}

		topicList = append(topicList, topicMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("topics", topicList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), topicList); e != nil {
			return e
		}
	}

	return nil
}

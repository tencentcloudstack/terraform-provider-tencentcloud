package trocket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	trocketv20230308 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTrocketRocketmqInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTrocketRocketmqInstancesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter query criteria list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"tag_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Tag values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance name.",
						},
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "version.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance type. EXPERIMENT: trial version; BASIC: Basic Edition; PRO: Professional Edition; PLATINUM: Platinum Edition.",
						},
						"instance_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance status. RUNNING, Running; MAINTAINING: Under maintenance; ABNORMAL: abnormal; OVERDUE: arrears; DESTROYED: Deleted; CREATING: Creating; MODIFYING: In the process of transformation; CREATE_FAILURE: Creation failed; MODIFY_FAILURE: Transformation failed; DELETING: deleting.",
						},
						"topic_num_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximum number of instance topics.",
						},
						"group_num_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximum number of instance consumer groups.",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pay mode. - POSTPAID: postpaid; - PREPAID: prepaid.",
						},
						"expiry_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Expiration timestamp, **Unix timestamp (in milliseconds)**.",
						},
						"remark": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Remark.",
						},
						"topic_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Topic nums.",
						},
						"group_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Group nums.",
						},
						"tag_list": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"sku_code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Product Specifications.",
						},
						"tps_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "TPS current limit value.",
						},
						"scaled_tps_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Elastic TPS current limit value.",
						},
						"message_retention": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Message retention time, in hours.",
						},
						"max_message_delay": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximum delay message duration in hours.",
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Whether to renew automatically, only for prepaid clusters (0: not renew automatically; 1: renew automatically).",
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

func dataSourceTencentCloudTrocketRocketmqInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_trocket_rocketmq_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*trocketv20230308.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := trocketv20230308.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
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

	if v, ok := d.GetOk("tag_filters"); ok {
		tagFiltersSet := v.([]interface{})
		tmpSet := make([]*trocketv20230308.TagFilter, 0, len(tagFiltersSet))
		for _, item := range tagFiltersSet {
			tagFiltersMap := item.(map[string]interface{})
			tagFilter := trocketv20230308.TagFilter{}
			if v, ok := tagFiltersMap["tag_key"].(string); ok && v != "" {
				tagFilter.TagKey = helper.String(v)
			}

			if v, ok := tagFiltersMap["tag_values"]; ok {
				tagValuesSet := v.(*schema.Set).List()
				for i := range tagValuesSet {
					tagValues := tagValuesSet[i].(string)
					tagFilter.TagValues = append(tagFilter.TagValues, helper.String(tagValues))
				}
			}

			tmpSet = append(tmpSet, &tagFilter)
		}

		paramMap["TagFilters"] = tmpSet
	}

	var respData []*trocketv20230308.InstanceItem
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTrocketRocketmqInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, data := range respData {
			dataMap := map[string]interface{}{}
			if data.InstanceId != nil {
				dataMap["instance_id"] = data.InstanceId
			}

			if data.InstanceName != nil {
				dataMap["instance_name"] = data.InstanceName
			}

			if data.Version != nil {
				dataMap["version"] = data.Version
			}

			if data.InstanceType != nil {
				dataMap["instance_type"] = data.InstanceType
			}

			if data.InstanceStatus != nil {
				dataMap["instance_status"] = data.InstanceStatus
			}

			if data.TopicNumLimit != nil {
				dataMap["topic_num_limit"] = data.TopicNumLimit
			}

			if data.GroupNumLimit != nil {
				dataMap["group_num_limit"] = data.GroupNumLimit
			}

			if data.PayMode != nil {
				dataMap["pay_mode"] = data.PayMode
			}

			if data.ExpiryTime != nil {
				dataMap["expiry_time"] = data.ExpiryTime
			}

			if data.Remark != nil {
				dataMap["remark"] = data.Remark
			}

			if data.TopicNum != nil {
				dataMap["topic_num"] = data.TopicNum
			}

			if data.GroupNum != nil {
				dataMap["group_num"] = data.GroupNum
			}

			tagListList := make([]map[string]interface{}, 0, len(data.TagList))
			if data.TagList != nil {
				for _, tagList := range data.TagList {
					tagListMap := map[string]interface{}{}

					if tagList.TagKey != nil {
						tagListMap["tag_key"] = tagList.TagKey
					}

					if tagList.TagValue != nil {
						tagListMap["tag_value"] = tagList.TagValue
					}

					tagListList = append(tagListList, tagListMap)
				}

				dataMap["tag_list"] = tagListList
			}
			if data.SkuCode != nil {
				dataMap["sku_code"] = data.SkuCode
			}

			if data.TpsLimit != nil {
				dataMap["tps_limit"] = data.TpsLimit
			}

			if data.ScaledTpsLimit != nil {
				dataMap["scaled_tps_limit"] = data.ScaledTpsLimit
			}

			if data.MessageRetention != nil {
				dataMap["message_retention"] = data.MessageRetention
			}

			if data.MaxMessageDelay != nil {
				dataMap["max_message_delay"] = data.MaxMessageDelay
			}

			if data.RenewFlag != nil {
				dataMap["renew_flag"] = data.RenewFlag
			}

			dataList = append(dataList, dataMap)
		}

		_ = d.Set("data", dataList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataList); e != nil {
			return e
		}
	}

	return nil
}

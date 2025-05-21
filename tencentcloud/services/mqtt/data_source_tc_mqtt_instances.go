package mqtt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMqttInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMqttInstancesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Query criteria list, supporting the following fields: InstanceName: cluster name, fuzzy search, InstanceId: cluster ID, precise search, InstanceStatus: cluster status search (RUNNING - Running, CREATING - Creating, MODIFYING - Changing, DELETING - Deleting).",
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
							Description: "Values.",
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
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_values": {
							Type:        schema.TypeSet,
							Optional:    true,
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
							Computed:    true,
							Description: "Instacen ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instacen name.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instacen version.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type. BASIC- Basic Edition; PRO- professional edition; PLATINUM- Platinum version.",
						},
						"instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status. RUNNING- In operation; MAINTAINING- Under Maintenance; ABNORMAL- abnormal; OVERDUE- Arrears of fees; DESTROYED- Deleted; CREATING- Creating in progress; MODIFYING- In the process of transformation; CREATE_FAILURE- Creation failed; MODIFY_FAILURE- Transformation failed; DELETING- deleting.",
						},
						"topic_num_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of instance topics.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark.",
						},
						"topic_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Topic num.",
						},
						"sku_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product specifications.",
						},
						"tps_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Elastic TPS current limit value.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation time, millisecond timestamp.",
						},
						"max_subscription_per_client": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of subscriptions per client.",
						},
						"client_num_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of client connections online.",
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to renew automatically. Only the annual and monthly package cluster is effective. 1: Automatic renewal; 0: Non automatic renewal.",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing mode, POSTPAID, pay as you go PREPAID, annual and monthly package.",
						},
						"expiry_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Expiration time, millisecond level timestamp.",
						},
						"destroy_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Pre destruction time, millisecond timestamp.",
						},
						"authorization_policy_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Limit on the number of authorization rules.",
						},
						"max_ca_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum CA quota.",
						},
						"max_subscription": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of subscriptions.",
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

func dataSourceTencentCloudMqttInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mqtt_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*mqttv20240516.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := mqttv20240516.Filter{}
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
		tmpSet := make([]*mqttv20240516.TagFilter, 0, len(tagFiltersSet))
		for _, item := range tagFiltersSet {
			tagFiltersMap := item.(map[string]interface{})
			tagFilter := mqttv20240516.TagFilter{}
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

	var respData []*mqttv20240516.MQTTInstanceItem
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMqttInstancesByFilter(ctx, paramMap)
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

			if data.Remark != nil {
				dataMap["remark"] = data.Remark
			}

			if data.TopicNum != nil {
				dataMap["topic_num"] = data.TopicNum
			}

			if data.SkuCode != nil {
				dataMap["sku_code"] = data.SkuCode
			}

			if data.TpsLimit != nil {
				dataMap["tps_limit"] = data.TpsLimit
			}

			if data.CreateTime != nil {
				dataMap["create_time"] = data.CreateTime
			}

			if data.MaxSubscriptionPerClient != nil {
				dataMap["max_subscription_per_client"] = data.MaxSubscriptionPerClient
			}

			if data.ClientNumLimit != nil {
				dataMap["client_num_limit"] = data.ClientNumLimit
			}

			if data.RenewFlag != nil {
				dataMap["renew_flag"] = data.RenewFlag
			}

			if data.PayMode != nil {
				dataMap["pay_mode"] = data.PayMode
			}

			if data.ExpiryTime != nil {
				dataMap["expiry_time"] = data.ExpiryTime
			}

			if data.DestroyTime != nil {
				dataMap["destroy_time"] = data.DestroyTime
			}

			if data.AuthorizationPolicyLimit != nil {
				dataMap["authorization_policy_limit"] = data.AuthorizationPolicyLimit
			}

			if data.MaxCaNum != nil {
				dataMap["max_ca_num"] = data.MaxCaNum
			}

			if data.MaxSubscription != nil {
				dataMap["max_subscription"] = data.MaxSubscription
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

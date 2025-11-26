package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmInstancePackageList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmInstancePackageListRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name, supported list as follows:\n- InstanceId: instance ID.\n- InstanceName: instance name.\n- ResourceId: package ID.\n- PackageType: package type. This is a required parameter, not passing it will cause interface query failure.",
						},
						"value": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field value.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable fuzzy query, only supports filter field name as domain.\nWhen fuzzy query is enabled, maximum Value length is 1, otherwise maximum Value length is 5. (Reserved field, not currently used).",
						},
					},
				},
			},

			"is_used": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether used: 0 not used 1 used.",
			},

			"instance_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance package list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance package resource ID.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package type\nFREE: Free version\nSTANDARD: Standard version\nULTIMATE: Ultimate version.",
						},
						"current_deadline": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package expiration time.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package creation time.",
						},
						"is_expire": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether expired 0 no 1 yes.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status\nENABLED: Normal\nDISABLED: Disabled.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether auto-renew 0 no 1 yes.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark.",
						},
						"cost_item_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Billing item.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cost_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Billing item name.",
									},
									"cost_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Billing item value.",
									},
								},
							},
						},
						"min_check_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum check interval time s.",
						},
						"min_global_ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum TTL s.",
						},
						"traffic_strategy": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Traffic strategy type: ALL return all, WEIGHT weight.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"schedule_strategy": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Strategy type: LOCATION schedule by geographic location, DELAY schedule by delay.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
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

func dataSourceTencentCloudIgtmInstancePackageListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_instance_package_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*igtmv20231024.ResourceFilter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			resourceFilter := igtmv20231024.ResourceFilter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				resourceFilter.Name = helper.String(v)
			}

			if v, ok := filtersMap["value"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					resourceFilter.Value = append(resourceFilter.Value, helper.String(value))
				}
			}

			if v, ok := filtersMap["fuzzy"].(bool); ok {
				resourceFilter.Fuzzy = helper.Bool(v)
			}

			tmpSet = append(tmpSet, &resourceFilter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOkExists("is_used"); ok {
		paramMap["IsUsed"] = helper.IntUint64(v.(int))
	}

	var respData []*igtmv20231024.InstancePackage
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIgtmInstancePackageListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	instanceSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, instanceSet := range respData {
			instanceSetMap := map[string]interface{}{}
			if instanceSet.ResourceId != nil {
				instanceSetMap["resource_id"] = instanceSet.ResourceId
			}

			if instanceSet.InstanceId != nil {
				instanceSetMap["instance_id"] = instanceSet.InstanceId
			}

			if instanceSet.InstanceName != nil {
				instanceSetMap["instance_name"] = instanceSet.InstanceName
			}

			if instanceSet.PackageType != nil {
				instanceSetMap["package_type"] = instanceSet.PackageType
			}

			if instanceSet.CurrentDeadline != nil {
				instanceSetMap["current_deadline"] = instanceSet.CurrentDeadline
			}

			if instanceSet.CreateTime != nil {
				instanceSetMap["create_time"] = instanceSet.CreateTime
			}

			if instanceSet.IsExpire != nil {
				instanceSetMap["is_expire"] = instanceSet.IsExpire
			}

			if instanceSet.Status != nil {
				instanceSetMap["status"] = instanceSet.Status
			}

			if instanceSet.AutoRenewFlag != nil {
				instanceSetMap["auto_renew_flag"] = instanceSet.AutoRenewFlag
			}

			if instanceSet.Remark != nil {
				instanceSetMap["remark"] = instanceSet.Remark
			}

			costItemListList := make([]map[string]interface{}, 0, len(instanceSet.CostItemList))
			if instanceSet.CostItemList != nil {
				for _, costItemList := range instanceSet.CostItemList {
					costItemListMap := map[string]interface{}{}
					if costItemList.CostName != nil {
						costItemListMap["cost_name"] = costItemList.CostName
					}

					if costItemList.CostValue != nil {
						costItemListMap["cost_value"] = costItemList.CostValue
					}

					costItemListList = append(costItemListList, costItemListMap)
				}

				instanceSetMap["cost_item_list"] = costItemListList
			}

			if instanceSet.MinCheckInterval != nil {
				instanceSetMap["min_check_interval"] = instanceSet.MinCheckInterval
			}

			if instanceSet.MinGlobalTtl != nil {
				instanceSetMap["min_global_ttl"] = instanceSet.MinGlobalTtl
			}

			if instanceSet.TrafficStrategy != nil {
				instanceSetMap["traffic_strategy"] = instanceSet.TrafficStrategy
			}

			if instanceSet.ScheduleStrategy != nil {
				instanceSetMap["schedule_strategy"] = instanceSet.ScheduleStrategy
			}

			instanceSetList = append(instanceSetList, instanceSetMap)
		}

		_ = d.Set("instance_set", instanceSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceSetList); e != nil {
			return e
		}
	}

	return nil
}

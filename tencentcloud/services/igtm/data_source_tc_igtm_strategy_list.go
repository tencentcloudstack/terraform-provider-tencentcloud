package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmStrategyList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmStrategyListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Strategy filter conditions: StrategyName: strategy name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name, supported list as follows:\n- type: main resource type, CDN.\n- instanceId: IGTM instance ID. This is a required parameter, failure to pass will cause interface query failure.",
						},
						"value": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable fuzzy query, only supports filter field name as domain.\nWhen fuzzy query is enabled, Value maximum length is 1, otherwise Value maximum length is 5. (Reserved field, currently unused).",
						},
					},
				},
			},

			"strategy_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Strategy list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Strategy name.",
						},
						"source": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Address source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_line_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resolution request source line ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resolution request source line name.",
									},
								},
							},
						},
						"strategy_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategy ID.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health status: ok healthy, warn risk, down failure.",
						},
						"activate_main_pool_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Activated main pool ID, null means unknown.",
						},
						"activate_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current activated address pool level, 0 means fallback activated, null means unknown.",
						},
						"active_pool_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current activated address pool set type: main main pool; fallback fallback pool.",
						},
						"active_traffic_strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current activated address pool traffic strategy: all resolve all; weight load balancing.",
						},
						"monitor_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Monitor count.",
						},
						"is_enabled": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether enabled: ENABLED enabled; DISABLED disabled.",
						},
						"keep_domain_records": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to retain lines: enabled retain, disabled not retain, only retain default lines.",
						},
						"switch_pool_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheduling mode: AUTO default; PAUSE only pause without switching.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
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

func dataSourceTencentCloudIgtmStrategyListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_strategy_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(nil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

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

	var respData []*igtmv20231024.Strategy
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIgtmStrategyListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	strategySetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, strategySet := range respData {
			strategySetMap := map[string]interface{}{}
			if strategySet.InstanceId != nil {
				strategySetMap["instance_id"] = strategySet.InstanceId
			}

			if strategySet.Name != nil {
				strategySetMap["name"] = strategySet.Name
			}

			sourceList := make([]map[string]interface{}, 0, len(strategySet.Source))
			if strategySet.Source != nil {
				for _, source := range strategySet.Source {
					sourceMap := map[string]interface{}{}
					if source.DnsLineId != nil {
						sourceMap["dns_line_id"] = source.DnsLineId
					}

					if source.Name != nil {
						sourceMap["name"] = source.Name
					}

					sourceList = append(sourceList, sourceMap)
				}

				strategySetMap["source"] = sourceList
			}

			if strategySet.StrategyId != nil {
				strategySetMap["strategy_id"] = strategySet.StrategyId
			}

			if strategySet.Status != nil {
				strategySetMap["status"] = strategySet.Status
			}

			if strategySet.ActivateMainPoolId != nil {
				strategySetMap["activate_main_pool_id"] = strategySet.ActivateMainPoolId
			}

			if strategySet.ActivateLevel != nil {
				strategySetMap["activate_level"] = strategySet.ActivateLevel
			}

			if strategySet.ActivePoolType != nil {
				strategySetMap["active_pool_type"] = strategySet.ActivePoolType
			}

			if strategySet.ActiveTrafficStrategy != nil {
				strategySetMap["active_traffic_strategy"] = strategySet.ActiveTrafficStrategy
			}

			if strategySet.MonitorNum != nil {
				strategySetMap["monitor_num"] = strategySet.MonitorNum
			}

			if strategySet.IsEnabled != nil {
				strategySetMap["is_enabled"] = strategySet.IsEnabled
			}

			if strategySet.KeepDomainRecords != nil {
				strategySetMap["keep_domain_records"] = strategySet.KeepDomainRecords
			}

			if strategySet.SwitchPoolType != nil {
				strategySetMap["switch_pool_type"] = strategySet.SwitchPoolType
			}

			if strategySet.CreatedOn != nil {
				strategySetMap["created_on"] = strategySet.CreatedOn
			}

			if strategySet.UpdatedOn != nil {
				strategySetMap["updated_on"] = strategySet.UpdatedOn
			}

			strategySetList = append(strategySetList, strategySetMap)
		}

		_ = d.Set("strategy_set", strategySetList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), strategySetList); e != nil {
			return e
		}
	}

	return nil
}

/*
Use this data source to query detailed information of tse group_list

Example Usage

```hcl
data "tencentcloud_tse_group_list" "group_list" {
  gateway_id = ""
  filters {
		name = "GroupId"
		values =

  }
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseGroupList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGroupListRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions, valid value:Name,GroupId.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter values.",
						},
					},
				},
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Groups infomation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total count.",
						},
						"gateway_group_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Group list of gateway.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group Id.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group name.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group description.",
									},
									"node_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Group node configration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"specification": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Group specification, 1c2g|2c4g|4c8g|8c16g.",
												},
												"number": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Group node number, 2-50.",
												},
											},
										},
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group status.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group create time.",
									},
									"is_first_group": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether it is the default group- 0: false.- 1: yes.",
									},
									"binding_strategy": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Associated strategy informationNote: This field may return null, indicating that a valid value is not available.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"strategy_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Strategy ID.",
												},
												"strategy_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Strategy nameNote: This field may return null, indicating that a valid value is not available.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Create timeNote: This field may return null, indicating that a valid value is not available.",
												},
												"modify_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Modify timeNote: This field may return null, indicating that a valid value is not available.",
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Description of strategyNote: This field may return null, indicating that a valid value is not available.",
												},
												"config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Auto scaling configurationNote: This field may return null, indicating that a valid value is not available.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"max_replicas": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum number of replicasNote: This field may return null, indicating that a valid value is not available.",
															},
															"metrics": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Metric listNote: This field may return null, indicating that a valid value is not available.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Metric typeNote: This field may return null, indicating that a valid value is not available.",
																		},
																		"resource_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Metric resource nameNote: This field may return null, indicating that a valid value is not available.",
																		},
																		"target_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Metric target typeNote: This field may return null, indicating that a valid value is not available.",
																		},
																		"target_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Metric target valueNote: This field may return null, indicating that a valid value is not available.",
																		},
																	},
																},
															},
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to enable metric auto scalingNote: This field may return null, indicating that a valid value is not available.",
															},
															"create_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Create timeNote: This field may return null, indicating that a valid value is not available.",
															},
															"modify_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Modify timeNote: This field may return null, indicating that a valid value is not available.",
															},
															"strategy_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Strategy IDNote: This field may return null, indicating that a valid value is not available.",
															},
															"auto_scaler_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Auto scaler IDNote: This field may return null, indicating that a valid value is not available.",
															},
														},
													},
												},
												"gateway_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway IDNote: This field may return null, indicating that a valid value is not available.",
												},
												"cron_config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Timing scaling configurationNote: This field may return null, indicating that a valid value is not available.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to enable timing auto scaling.",
															},
															"params": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Params of timing auto scaling.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"period": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Period of timing auto scaling.",
																		},
																		"start_at": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Start time.",
																		},
																		"target_replicas": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Target replicas.",
																		},
																		"crontab": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cron expression.",
																		},
																	},
																},
															},
															"create_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Create time.",
															},
															"modify_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Modify time.",
															},
															"strategy_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Strategy ID.",
															},
														},
													},
												},
												"max_replicas": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum number of replicas.",
												},
											},
										},
									},
									"gateway_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway ID.",
									},
									"internet_max_bandwidth_out": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Public network outbound traffic bandwidth.",
									},
									"modify_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Modify time.",
									},
									"subnet_ids": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet IDs.",
									},
								},
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

func dataSourceTencentCloudTseGroupListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_group_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tse.Filter{}
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

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tse.NativeGatewayServerGroups

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseGroupListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		nativeGatewayServerGroupsMap := map[string]interface{}{}

		if result.TotalCount != nil {
			nativeGatewayServerGroupsMap["total_count"] = result.TotalCount
		}

		if result.GatewayGroupList != nil {
			gatewayGroupListList := []interface{}{}
			for _, gatewayGroupList := range result.GatewayGroupList {
				gatewayGroupListMap := map[string]interface{}{}

				if gatewayGroupList.GroupId != nil {
					gatewayGroupListMap["group_id"] = gatewayGroupList.GroupId
				}

				if gatewayGroupList.Name != nil {
					gatewayGroupListMap["name"] = gatewayGroupList.Name
				}

				if gatewayGroupList.Description != nil {
					gatewayGroupListMap["description"] = gatewayGroupList.Description
				}

				if gatewayGroupList.NodeConfig != nil {
					nodeConfigMap := map[string]interface{}{}

					if gatewayGroupList.NodeConfig.Specification != nil {
						nodeConfigMap["specification"] = gatewayGroupList.NodeConfig.Specification
					}

					if gatewayGroupList.NodeConfig.Number != nil {
						nodeConfigMap["number"] = gatewayGroupList.NodeConfig.Number
					}

					gatewayGroupListMap["node_config"] = []interface{}{nodeConfigMap}
				}

				if gatewayGroupList.Status != nil {
					gatewayGroupListMap["status"] = gatewayGroupList.Status
				}

				if gatewayGroupList.CreateTime != nil {
					gatewayGroupListMap["create_time"] = gatewayGroupList.CreateTime
				}

				if gatewayGroupList.IsFirstGroup != nil {
					gatewayGroupListMap["is_first_group"] = gatewayGroupList.IsFirstGroup
				}

				if gatewayGroupList.BindingStrategy != nil {
					bindingStrategyMap := map[string]interface{}{}

					if gatewayGroupList.BindingStrategy.StrategyId != nil {
						bindingStrategyMap["strategy_id"] = gatewayGroupList.BindingStrategy.StrategyId
					}

					if gatewayGroupList.BindingStrategy.StrategyName != nil {
						bindingStrategyMap["strategy_name"] = gatewayGroupList.BindingStrategy.StrategyName
					}

					if gatewayGroupList.BindingStrategy.CreateTime != nil {
						bindingStrategyMap["create_time"] = gatewayGroupList.BindingStrategy.CreateTime
					}

					if gatewayGroupList.BindingStrategy.ModifyTime != nil {
						bindingStrategyMap["modify_time"] = gatewayGroupList.BindingStrategy.ModifyTime
					}

					if gatewayGroupList.BindingStrategy.Description != nil {
						bindingStrategyMap["description"] = gatewayGroupList.BindingStrategy.Description
					}

					if gatewayGroupList.BindingStrategy.Config != nil {
						configMap := map[string]interface{}{}

						if gatewayGroupList.BindingStrategy.Config.MaxReplicas != nil {
							configMap["max_replicas"] = gatewayGroupList.BindingStrategy.Config.MaxReplicas
						}

						if gatewayGroupList.BindingStrategy.Config.Metrics != nil {
							metricsList := []interface{}{}
							for _, metrics := range gatewayGroupList.BindingStrategy.Config.Metrics {
								metricsMap := map[string]interface{}{}

								if metrics.Type != nil {
									metricsMap["type"] = metrics.Type
								}

								if metrics.ResourceName != nil {
									metricsMap["resource_name"] = metrics.ResourceName
								}

								if metrics.TargetType != nil {
									metricsMap["target_type"] = metrics.TargetType
								}

								if metrics.TargetValue != nil {
									metricsMap["target_value"] = metrics.TargetValue
								}

								metricsList = append(metricsList, metricsMap)
							}

							configMap["metrics"] = []interface{}{metricsList}
						}

						if gatewayGroupList.BindingStrategy.Config.Enabled != nil {
							configMap["enabled"] = gatewayGroupList.BindingStrategy.Config.Enabled
						}

						if gatewayGroupList.BindingStrategy.Config.CreateTime != nil {
							configMap["create_time"] = gatewayGroupList.BindingStrategy.Config.CreateTime
						}

						if gatewayGroupList.BindingStrategy.Config.ModifyTime != nil {
							configMap["modify_time"] = gatewayGroupList.BindingStrategy.Config.ModifyTime
						}

						if gatewayGroupList.BindingStrategy.Config.StrategyId != nil {
							configMap["strategy_id"] = gatewayGroupList.BindingStrategy.Config.StrategyId
						}

						if gatewayGroupList.BindingStrategy.Config.AutoScalerId != nil {
							configMap["auto_scaler_id"] = gatewayGroupList.BindingStrategy.Config.AutoScalerId
						}

						bindingStrategyMap["config"] = []interface{}{configMap}
					}

					if gatewayGroupList.BindingStrategy.GatewayId != nil {
						bindingStrategyMap["gateway_id"] = gatewayGroupList.BindingStrategy.GatewayId
					}

					if gatewayGroupList.BindingStrategy.CronConfig != nil {
						cronConfigMap := map[string]interface{}{}

						if gatewayGroupList.BindingStrategy.CronConfig.Enabled != nil {
							cronConfigMap["enabled"] = gatewayGroupList.BindingStrategy.CronConfig.Enabled
						}

						if gatewayGroupList.BindingStrategy.CronConfig.Params != nil {
							paramsList := []interface{}{}
							for _, params := range gatewayGroupList.BindingStrategy.CronConfig.Params {
								paramsMap := map[string]interface{}{}

								if params.Period != nil {
									paramsMap["period"] = params.Period
								}

								if params.StartAt != nil {
									paramsMap["start_at"] = params.StartAt
								}

								if params.TargetReplicas != nil {
									paramsMap["target_replicas"] = params.TargetReplicas
								}

								if params.Crontab != nil {
									paramsMap["crontab"] = params.Crontab
								}

								paramsList = append(paramsList, paramsMap)
							}

							cronConfigMap["params"] = []interface{}{paramsList}
						}

						if gatewayGroupList.BindingStrategy.CronConfig.CreateTime != nil {
							cronConfigMap["create_time"] = gatewayGroupList.BindingStrategy.CronConfig.CreateTime
						}

						if gatewayGroupList.BindingStrategy.CronConfig.ModifyTime != nil {
							cronConfigMap["modify_time"] = gatewayGroupList.BindingStrategy.CronConfig.ModifyTime
						}

						if gatewayGroupList.BindingStrategy.CronConfig.StrategyId != nil {
							cronConfigMap["strategy_id"] = gatewayGroupList.BindingStrategy.CronConfig.StrategyId
						}

						bindingStrategyMap["cron_config"] = []interface{}{cronConfigMap}
					}

					if gatewayGroupList.BindingStrategy.MaxReplicas != nil {
						bindingStrategyMap["max_replicas"] = gatewayGroupList.BindingStrategy.MaxReplicas
					}

					gatewayGroupListMap["binding_strategy"] = []interface{}{bindingStrategyMap}
				}

				if gatewayGroupList.GatewayId != nil {
					gatewayGroupListMap["gateway_id"] = gatewayGroupList.GatewayId
				}

				if gatewayGroupList.InternetMaxBandwidthOut != nil {
					gatewayGroupListMap["internet_max_bandwidth_out"] = gatewayGroupList.InternetMaxBandwidthOut
				}

				if gatewayGroupList.ModifyTime != nil {
					gatewayGroupListMap["modify_time"] = gatewayGroupList.ModifyTime
				}

				if gatewayGroupList.SubnetIds != nil {
					gatewayGroupListMap["subnet_ids"] = gatewayGroupList.SubnetIds
				}

				gatewayGroupListList = append(gatewayGroupListList, gatewayGroupListMap)
			}

			nativeGatewayServerGroupsMap["gateway_group_list"] = []interface{}{gatewayGroupListList}
		}

		ids = append(ids, *result.GroupId)
		_ = d.Set("result", nativeGatewayServerGroupsMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), nativeGatewayServerGroupsMap); e != nil {
			return e
		}
	}
	return nil
}

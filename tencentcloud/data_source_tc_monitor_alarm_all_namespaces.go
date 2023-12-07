package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMonitorAlarmAllNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmAllNamespacesRead,
		Schema: map[string]*schema.Schema{
			"scene_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Currently, only ST_ALARM=alarm type is filtered based on usage scenarios.",
			},

			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Fixed value, as `monitor`.",
			},

			"monitor_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter based on monitoring type, do not fill in default, check all types MT_QCE=cloud product monitoring.",
			},

			"ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter based on the Id of the namespace without filling in the default query for all.",
			},

			"qce_namespaces_new": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Types of alarm strategies for cloud products.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace labeling.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace name.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace value.",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product Name.",
						},
						"config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration information.",
						},
						"available_regions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of supported regions.",
						},
						"sort_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sort Id.",
						},
						"dashboard_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique representation in dashboard.",
						},
					},
				},
			},

			"custom_namespaces_new": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Other alarm strategy types are currently not supported.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace labeling.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace name.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace value.",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product Name.",
						},
						"config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration information.",
						},
						"available_regions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of supported regions.",
						},
						"sort_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sort Id.",
						},
						"dashboard_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique representation in dashboard.",
						},
					},
				},
			},

			"common_namespaces": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "General alarm strategy types (including: application performance monitoring, front-end performance monitoring, cloud dial testing).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace labeling.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace name.",
						},
						"monitor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring type.",
						},
						"dimensions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Dimension Information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dimension key identifier, backend English name.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dimension key name, Chinese and English frontend display name.",
									},
									"is_required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Required or not.",
									},
									"operators": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of supported operators.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Operator identification.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Operator Display Name.",
												},
											},
										},
									},
									"is_multiple": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Do you support multiple selections.",
									},
									"is_mutable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can I modify it after creation.",
									},
									"is_visible": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to display to users.",
									},
									"can_filter_policy": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can it be used to filter the policy list.",
									},
									"can_filter_history": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can it be used to filter alarm history.",
									},
									"can_group_by": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can it be used as an aggregation dimension.",
									},
									"must_group_by": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Must it be used as an aggregation dimension.",
									},
									"show_value_replace": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key to replace in front-end translation.",
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

func dataSourceTencentCloudMonitorAlarmAllNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_monitor_alarm_all_namespaces.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("scene_type"); ok {
		paramMap["SceneType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("monitor_types"); ok {
		monitorTypesSet := v.(*schema.Set).List()
		paramMap["MonitorTypes"] = helper.InterfacesStringsPoint(monitorTypesSet)
	}

	if v, ok := d.GetOk("ids"); ok {
		idsSet := v.(*schema.Set).List()
		paramMap["Ids"] = helper.InterfacesStringsPoint(idsSet)
	}

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	var qceNamespacesNew []*monitor.CommonNamespace
	var customNamespacesNew []*monitor.CommonNamespace
	var commonNamespaces []*monitor.CommonNamespaceNew
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		qce, custom, common, e := service.DescribeMonitorAlarmAllNamespacesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		qceNamespacesNew = qce
		customNamespacesNew = custom
		commonNamespaces = common
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	if qceNamespacesNew != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, commonNamespace := range qceNamespacesNew {
			commonNamespaceMap := map[string]interface{}{}

			if commonNamespace.Id != nil {
				commonNamespaceMap["id"] = commonNamespace.Id
			}

			if commonNamespace.Name != nil {
				commonNamespaceMap["name"] = commonNamespace.Name
			}

			if commonNamespace.Value != nil {
				commonNamespaceMap["value"] = commonNamespace.Value
			}

			if commonNamespace.ProductName != nil {
				commonNamespaceMap["product_name"] = commonNamespace.ProductName
			}

			if commonNamespace.Config != nil {
				commonNamespaceMap["config"] = commonNamespace.Config
			}

			if commonNamespace.AvailableRegions != nil {
				commonNamespaceMap["available_regions"] = commonNamespace.AvailableRegions
			}

			if commonNamespace.SortId != nil {
				commonNamespaceMap["sort_id"] = commonNamespace.SortId
			}

			if commonNamespace.DashboardId != nil {
				commonNamespaceMap["dashboard_id"] = commonNamespace.DashboardId
			}

			ids = append(ids, *commonNamespace.Id)
			tmpList = append(tmpList, commonNamespaceMap)
		}

		_ = d.Set("qce_namespaces_new", tmpList)
	}

	if customNamespacesNew != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, commonNamespace := range customNamespacesNew {
			commonNamespaceMap := map[string]interface{}{}

			if commonNamespace.Id != nil {
				commonNamespaceMap["id"] = commonNamespace.Id
			}

			if commonNamespace.Name != nil {
				commonNamespaceMap["name"] = commonNamespace.Name
			}

			if commonNamespace.Value != nil {
				commonNamespaceMap["value"] = commonNamespace.Value
			}

			if commonNamespace.ProductName != nil {
				commonNamespaceMap["product_name"] = commonNamespace.ProductName
			}

			if commonNamespace.Config != nil {
				commonNamespaceMap["config"] = commonNamespace.Config
			}

			if commonNamespace.AvailableRegions != nil {
				commonNamespaceMap["available_regions"] = commonNamespace.AvailableRegions
			}

			if commonNamespace.SortId != nil {
				commonNamespaceMap["sort_id"] = commonNamespace.SortId
			}

			if commonNamespace.DashboardId != nil {
				commonNamespaceMap["dashboard_id"] = commonNamespace.DashboardId
			}

			ids = append(ids, *commonNamespace.Id)
			tmpList = append(tmpList, commonNamespaceMap)
		}

		_ = d.Set("custom_namespaces_new", tmpList)
	}

	if commonNamespaces != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, commonNamespaceNew := range commonNamespaces {
			commonNamespaceNewMap := map[string]interface{}{}

			if commonNamespaceNew.Id != nil {
				commonNamespaceNewMap["id"] = commonNamespaceNew.Id
			}

			if commonNamespaceNew.Name != nil {
				commonNamespaceNewMap["name"] = commonNamespaceNew.Name
			}

			if commonNamespaceNew.MonitorType != nil {
				commonNamespaceNewMap["monitor_type"] = commonNamespaceNew.MonitorType
			}

			if commonNamespaceNew.Dimensions != nil {
				dimensionsList := []interface{}{}
				for _, dimensions := range commonNamespaceNew.Dimensions {
					dimensionsMap := map[string]interface{}{}

					if dimensions.Key != nil {
						dimensionsMap["key"] = dimensions.Key
					}

					if dimensions.Name != nil {
						dimensionsMap["name"] = dimensions.Name
					}

					if dimensions.IsRequired != nil {
						dimensionsMap["is_required"] = dimensions.IsRequired
					}

					if dimensions.Operators != nil {
						operatorsList := []interface{}{}
						for _, operators := range dimensions.Operators {
							operatorsMap := map[string]interface{}{}

							if operators.Id != nil {
								operatorsMap["id"] = operators.Id
							}

							if operators.Name != nil {
								operatorsMap["name"] = operators.Name
							}

							operatorsList = append(operatorsList, operatorsMap)
						}

						dimensionsMap["operators"] = operatorsList
					}

					if dimensions.IsMultiple != nil {
						dimensionsMap["is_multiple"] = dimensions.IsMultiple
					}

					if dimensions.IsMutable != nil {
						dimensionsMap["is_mutable"] = dimensions.IsMutable
					}

					if dimensions.IsVisible != nil {
						dimensionsMap["is_visible"] = dimensions.IsVisible
					}

					if dimensions.CanFilterPolicy != nil {
						dimensionsMap["can_filter_policy"] = dimensions.CanFilterPolicy
					}

					if dimensions.CanFilterHistory != nil {
						dimensionsMap["can_filter_history"] = dimensions.CanFilterHistory
					}

					if dimensions.CanGroupBy != nil {
						dimensionsMap["can_group_by"] = dimensions.CanGroupBy
					}

					if dimensions.MustGroupBy != nil {
						dimensionsMap["must_group_by"] = dimensions.MustGroupBy
					}

					if dimensions.ShowValueReplace != nil {
						dimensionsMap["show_value_replace"] = dimensions.ShowValueReplace
					}

					dimensionsList = append(dimensionsList, dimensionsMap)
				}

				commonNamespaceNewMap["dimensions"] = dimensionsList
			}

			ids = append(ids, *commonNamespaceNew.Id)
			tmpList = append(tmpList, commonNamespaceNewMap)
		}

		_ = d.Set("common_namespaces", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}

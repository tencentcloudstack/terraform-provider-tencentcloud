package emr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEmrServiceNodeInfos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEmrServiceNodeInfosRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EMR Instance ID.",
			},

			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page Number.",
			},

			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of Items per Page.",
			},

			"search_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search Field.",
			},

			"conf_status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Configuration Status, -2: Configuration Failed, -1: Configuration Expired, 1: Synchronized, -99 All.",
			},

			"maintain_state_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter Condition: Maintenance Status - 0 represents all statuses, 1 represents normal mode, 2 represents maintenance mode.",
			},

			"operator_state_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter Condition: Operation Status - 0 represents all statuses, 1 represents started, 2 represents stopped.",
			},

			"health_state_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter Conditions: Health Status, 0 represents unavailable, 1 represents good, -2 represents unknown, -99 represents all, -3 represents potential risks, -4 represents not detected.",
			},

			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service Component Name, all in uppercase, e.g., YARN.",
			},

			"node_type_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Node Names: master, core, task, common, router, all.",
			},

			"data_node_maintenance_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter Condition: Whether DN is in Maintenance Mode - 0 represents all statuses, 1 represents in maintenance mode.",
			},

			"search_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search Fields.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Types Supported for Search.",
						},
						"search_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Values Supported for Search.",
						},
					},
				},
			},

			"total_cnt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Count.",
			},

			"alias_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Serialization of Aliases for All Nodes in the Cluster.",
			},

			"support_node_flag_filter_list": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Supported FlagNode List.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"service_node_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Service Node Detail Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the node where the process resides.",
						},
						"node_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node Type.",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node Name.",
						},
						"service_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Service Status.",
						},
						"monitor_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Monitor Status.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status.",
						},
						"ports_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Process Port Information.",
						},
						"last_restart_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Most Recent Restart Time.",
						},
						"flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Flag.",
						},
						"conf_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Configuration Group ID.",
						},
						"conf_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration Group Name.",
						},
						"conf_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Configuration Status.",
						},
						"service_detection_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Process Detection Information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"detect_alert": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detection Alert Level.",
									},
									"detect_function_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detection Function Description.",
									},
									"detect_function_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detection Function Result.",
									},
									"detect_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detection Time.",
									},
								},
							},
						},
						"node_flag_filter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node Flag Filter.",
						},
						"health_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Process Health Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Health Status Code.",
									},
									"text": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Health Status Description.",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Health Status Description.",
									},
								},
							},
						},
						"is_support_role_monitor": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Monitoring is Supported.",
						},
						"stop_policies": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Stop Policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy Name.",
									},
									"display_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy Display Name.",
									},
									"describe": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy Description.",
									},
									"batch_size_range": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Batch  Node Count Optional Range.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"is_default": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether it is the Default Policy.",
									},
								},
							},
						},
						"ha_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HA State.",
						},
						"name_service": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name Service.",
						},
						"is_federation": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Federation is Supported.",
						},
						"data_node_maintenance_state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data Node Maintenance State.",
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

func dataSourceTencentCloudEmrServiceNodeInfosRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_emr_service_node_infos.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(instanceId)
	}

	if v, ok := d.GetOk("offset"); ok {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("limit"); ok {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("search_text"); ok {
		paramMap["SearchText"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("conf_status"); ok {
		paramMap["ConfStatus"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("maintain_state_id"); ok {
		paramMap["MaintainStateId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("operator_state_id"); ok {
		paramMap["OperatorStateId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("health_state_id"); ok {
		paramMap["HealthStateId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_name"); ok {
		paramMap["ServiceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("node_type_name"); ok {
		paramMap["NodeTypeName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("data_node_maintenance_id"); ok {
		paramMap["DataNodeMaintenanceId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("search_fields"); ok {
		searchFields := v.([]interface{})
		searchFieldList := make([]interface{}, 0, len(searchFields))
		for _, searchField := range searchFields {
			searchItem := &emr.SearchItem{}
			searchFieldMap := searchField.(map[string]interface{})
			if v, ok := searchFieldMap["search_type"]; ok {
				searchItem.SearchType = helper.String(v.(string))
			}
			if v, ok := searchFieldMap["search_value"]; ok {
				searchItem.SearchValue = helper.String(v.(string))
			}
			searchFieldList = append(searchFieldList, searchItem)
		}

		paramMap["SearchFields"] = searchFieldList
	}

	var respData *emr.DescribeServiceNodeInfosResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEmrServiceNodeInfosByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if respData.TotalCnt != nil {
		_ = d.Set("total_cnt", respData.TotalCnt)
	}

	if respData.AliasInfo != nil {
		_ = d.Set("alias_info", respData.AliasInfo)
	}

	if respData.SupportNodeFlagFilterList != nil {
		_ = d.Set("support_node_flag_filter_list", respData.SupportNodeFlagFilterList)
	}

	serviceNodeListList := make([]map[string]interface{}, 0, len(respData.ServiceNodeList))
	if respData.ServiceNodeList != nil {
		for _, serviceNodeList := range respData.ServiceNodeList {
			serviceNodeListMap := map[string]interface{}{}

			if serviceNodeList.Ip != nil {
				serviceNodeListMap["ip"] = serviceNodeList.Ip
			}

			if serviceNodeList.NodeType != nil {
				serviceNodeListMap["node_type"] = serviceNodeList.NodeType
			}

			if serviceNodeList.NodeName != nil {
				serviceNodeListMap["node_name"] = serviceNodeList.NodeName
			}

			if serviceNodeList.ServiceStatus != nil {
				serviceNodeListMap["service_status"] = serviceNodeList.ServiceStatus
			}

			if serviceNodeList.MonitorStatus != nil {
				serviceNodeListMap["monitor_status"] = serviceNodeList.MonitorStatus
			}

			if serviceNodeList.Status != nil {
				serviceNodeListMap["status"] = serviceNodeList.Status
			}

			if serviceNodeList.PortsInfo != nil {
				serviceNodeListMap["port_info"] = serviceNodeList.PortsInfo
			}

			if serviceNodeList.LastRestartTime != nil {
				serviceNodeListMap["last_restart_time"] = serviceNodeList.LastRestartTime
			}

			if serviceNodeList.Flag != nil {
				serviceNodeListMap["flag"] = serviceNodeList.Flag
			}

			if serviceNodeList.ConfGroupId != nil {
				serviceNodeListMap["conf_group_id"] = serviceNodeList.ConfGroupId
			}

			if serviceNodeList.ConfGroupName != nil {
				serviceNodeListMap["conf_group_name"] = serviceNodeList.ConfGroupName
			}

			if serviceNodeList.ConfStatus != nil {
				serviceNodeListMap["conf_status"] = serviceNodeList.ConfStatus
			}

			serviceDetectionInfoList := make([]map[string]interface{}, 0, len(serviceNodeList.ServiceDetectionInfo))
			if serviceNodeList.ServiceDetectionInfo != nil {
				for _, serviceDetectionInfo := range serviceNodeList.ServiceDetectionInfo {
					serviceDetectionInfoMap := map[string]interface{}{}

					if serviceDetectionInfo.DetectAlert != nil {
						serviceDetectionInfoMap["detect_alert"] = serviceDetectionInfo.DetectAlert
					}

					if serviceDetectionInfo.DetectFunctionKey != nil {
						serviceDetectionInfoMap["detect_function_key"] = serviceDetectionInfo.DetectFunctionKey
					}

					if serviceDetectionInfo.DetectFunctionValue != nil {
						serviceDetectionInfoMap["detect_function_value"] = serviceDetectionInfo.DetectFunctionValue
					}

					if serviceDetectionInfo.DetectTime != nil {
						serviceDetectionInfoMap["detect_time"] = serviceDetectionInfo.DetectTime
					}

					serviceDetectionInfoList = append(serviceDetectionInfoList, serviceDetectionInfoMap)
				}

				serviceNodeListMap["service_detection_info"] = serviceDetectionInfoList
			}
			if serviceNodeList.NodeFlagFilter != nil {
				serviceNodeListMap["node_flag_filter"] = serviceNodeList.NodeFlagFilter
			}

			healthStatusMap := map[string]interface{}{}

			if serviceNodeList.HealthStatus != nil {
				if serviceNodeList.HealthStatus.Code != nil {
					healthStatusMap["code"] = serviceNodeList.HealthStatus.Code
				}

				if serviceNodeList.HealthStatus.Text != nil {
					healthStatusMap["text"] = serviceNodeList.HealthStatus.Text
				}

				if serviceNodeList.HealthStatus.Desc != nil {
					healthStatusMap["desc"] = serviceNodeList.HealthStatus.Desc
				}

				serviceNodeListMap["health_status"] = []interface{}{healthStatusMap}
			}

			if serviceNodeList.IsSupportRoleMonitor != nil {
				serviceNodeListMap["is_support_role_monitor"] = serviceNodeList.IsSupportRoleMonitor
			}

			stopPolicyList := make([]interface{}, 0, len(serviceNodeList.StopPolicies))

			if serviceNodeList.StopPolicies != nil {
				for _, stopPolicie := range serviceNodeList.StopPolicies {
					stopPoliciesMap := map[string]interface{}{}
					if stopPolicie.Name != nil {
						stopPoliciesMap["name"] = stopPolicie.Name
					}

					if stopPolicie.DisplayName != nil {
						stopPoliciesMap["display_name"] = stopPolicie.DisplayName
					}

					if stopPolicie.Describe != nil {
						stopPoliciesMap["describe"] = stopPolicie.Describe
					}

					if stopPolicie.BatchSizeRange != nil {
						stopPoliciesMap["bath_size_range"] = stopPolicie.BatchSizeRange
					}

					if stopPolicie.IsDefault != nil {
						stopPoliciesMap["is_default"] = stopPolicie.IsDefault
					}
					stopPolicyList = append(stopPolicyList, stopPoliciesMap)
				}

				_ = d.Set("stop_policies", stopPolicyList)
			}

			if serviceNodeList.HAState != nil {
				serviceNodeListMap["ha_state"] = serviceNodeList.HAState
			}

			if serviceNodeList.NameService != nil {
				serviceNodeListMap["name_service"] = serviceNodeList.NameService
			}

			if serviceNodeList.IsFederation != nil {
				serviceNodeListMap["is_federation"] = serviceNodeList.IsFederation
			}

			if serviceNodeList.DataNodeMaintenanceState != nil {
				serviceNodeListMap["data_node_maintenance_state"] = serviceNodeList.DataNodeMaintenanceState
			}

			serviceNodeListList = append(serviceNodeListList, serviceNodeListMap)
		}

		_ = d.Set("service_node_list", serviceNodeListList)
	}

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

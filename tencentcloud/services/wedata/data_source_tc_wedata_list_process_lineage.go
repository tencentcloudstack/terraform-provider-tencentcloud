package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataListProcessLineage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataListProcessLineageRead,
		Schema: map[string]*schema.Schema{
			"process_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task unique ID.",
			},

			"process_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task type: SCHEDULE_TASK, INTEGRATION_TASK, THIRD_REPORT, TABLE_MODEL, MODEL_METRIC, METRIC_METRIC, DATA_SERVICE.",
			},

			"platform": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source: WEDATA|THIRD, default WEDATA.",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Lineage pair list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Entity original unique ID.\n\nNote: When lineage is for table columns, the unique ID should be TableResourceUniqueId::FieldName.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Entity type.\nTABLE|METRIC|MODEL|SERVICE|COLUMN.",
									},
									"platform": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source: WEDATA|THIRD.\nDefault wedata.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Business name: database.table|metric name|model name|field name.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description: table type|metric description|model description|field description.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update time.",
									},
									"resource_properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Resource additional extension parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property name.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property value.",
												},
											},
										},
									},
									"lineage_node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Lineage node unique identifier.",
									},
								},
							},
						},
						"target": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Entity original unique ID.\n\nNote: When lineage is for table columns, the unique ID should be TableResourceUniqueId::FieldName.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Entity type.\nTABLE|METRIC|MODEL|SERVICE|COLUMN.",
									},
									"platform": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source: WEDATA|THIRD.\nDefault wedata.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Business name: database.table|metric name|model name|field name.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description: table type|metric description|model description|field description.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update time.",
									},
									"resource_properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Resource additional extension parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property name.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property value.",
												},
											},
										},
									},
									"lineage_node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Lineage node unique identifier.",
									},
								},
							},
						},
						"processes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Lineage processing procedures.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"process_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Original unique ID.",
									},
									"process_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Task type.\nSCHEDULE_TASK,\nINTEGRATION_TASK,\nTHIRD_REPORT,\nTABLE_MODEL,\nMODEL_METRIC,\nMETRIC_METRIC,\nDATA_SERVICE.",
									},
									"platform": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "WEDATA, THIRD.",
									},
									"process_sub_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Task subtype.\nSQL_TASK,\nINTEGRATED_STREAM,\nINTEGRATED_OFFLINE.",
									},
									"process_properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Additional extension parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property name.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property value.",
												},
											},
										},
									},
									"lineage_node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Lineage task unique node ID.",
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

func dataSourceTencentCloudWedataListProcessLineageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_list_process_lineage.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("process_id"); ok {
		paramMap["ProcessId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("process_type"); ok {
		paramMap["ProcessType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("platform"); ok {
		paramMap["Platform"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.LineagePair
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataListProcessLineageByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	itemsList := make([]map[string]interface{}, 0, len(respData))
	for _, items := range respData {
		itemsMap := map[string]interface{}{}
		sourceMap := map[string]interface{}{}
		if items.Source != nil {
			if items.Source.ResourceUniqueId != nil {
				sourceMap["resource_unique_id"] = items.Source.ResourceUniqueId
			}

			if items.Source.ResourceType != nil {
				sourceMap["resource_type"] = items.Source.ResourceType
			}

			if items.Source.Platform != nil {
				sourceMap["platform"] = items.Source.Platform
			}

			if items.Source.ResourceName != nil {
				sourceMap["resource_name"] = items.Source.ResourceName
			}

			if items.Source.Description != nil {
				sourceMap["description"] = items.Source.Description
			}

			if items.Source.CreateTime != nil {
				sourceMap["create_time"] = items.Source.CreateTime
			}

			if items.Source.UpdateTime != nil {
				sourceMap["update_time"] = items.Source.UpdateTime
			}

			resourcePropertiesList := make([]map[string]interface{}, 0, len(items.Source.ResourceProperties))
			if items.Source.ResourceProperties != nil {
				for _, resourceProperties := range items.Source.ResourceProperties {
					resourcePropertiesMap := map[string]interface{}{}

					if resourceProperties.Name != nil {
						resourcePropertiesMap["name"] = resourceProperties.Name
					}

					if resourceProperties.Value != nil {
						resourcePropertiesMap["value"] = resourceProperties.Value
					}

					resourcePropertiesList = append(resourcePropertiesList, resourcePropertiesMap)
				}

				sourceMap["resource_properties"] = resourcePropertiesList
			}

			if items.Source.LineageNodeId != nil {
				sourceMap["lineage_node_id"] = items.Source.LineageNodeId
			}

			itemsMap["source"] = []interface{}{sourceMap}
		}

		targetMap := map[string]interface{}{}
		if items.Target != nil {
			if items.Target.ResourceUniqueId != nil {
				targetMap["resource_unique_id"] = items.Target.ResourceUniqueId
			}

			if items.Target.ResourceType != nil {
				targetMap["resource_type"] = items.Target.ResourceType
			}

			if items.Target.Platform != nil {
				targetMap["platform"] = items.Target.Platform
			}

			if items.Target.ResourceName != nil {
				targetMap["resource_name"] = items.Target.ResourceName
			}

			if items.Target.Description != nil {
				targetMap["description"] = items.Target.Description
			}

			if items.Target.CreateTime != nil {
				targetMap["create_time"] = items.Target.CreateTime
			}

			if items.Target.UpdateTime != nil {
				targetMap["update_time"] = items.Target.UpdateTime
			}

			resourcePropertiesList := make([]map[string]interface{}, 0, len(items.Target.ResourceProperties))
			if items.Target.ResourceProperties != nil {
				for _, resourceProperties := range items.Target.ResourceProperties {
					resourcePropertiesMap := map[string]interface{}{}

					if resourceProperties.Name != nil {
						resourcePropertiesMap["name"] = resourceProperties.Name
					}

					if resourceProperties.Value != nil {
						resourcePropertiesMap["value"] = resourceProperties.Value
					}

					resourcePropertiesList = append(resourcePropertiesList, resourcePropertiesMap)
				}

				targetMap["resource_properties"] = resourcePropertiesList
			}

			if items.Target.LineageNodeId != nil {
				targetMap["lineage_node_id"] = items.Target.LineageNodeId
			}

			itemsMap["target"] = []interface{}{targetMap}
		}

		processesList := make([]map[string]interface{}, 0, len(items.Processes))
		if items.Processes != nil {
			for _, processes := range items.Processes {
				processesMap := map[string]interface{}{}
				if processes.ProcessId != nil {
					processesMap["process_id"] = processes.ProcessId
				}

				if processes.ProcessType != nil {
					processesMap["process_type"] = processes.ProcessType
				}

				if processes.Platform != nil {
					processesMap["platform"] = processes.Platform
				}

				if processes.ProcessSubType != nil {
					processesMap["process_sub_type"] = processes.ProcessSubType
				}

				processPropertiesList := make([]map[string]interface{}, 0, len(processes.ProcessProperties))
				if processes.ProcessProperties != nil {
					for _, processProperties := range processes.ProcessProperties {
						processPropertiesMap := map[string]interface{}{}
						if processProperties.Name != nil {
							processPropertiesMap["name"] = processProperties.Name
						}

						if processProperties.Value != nil {
							processPropertiesMap["value"] = processProperties.Value
						}

						processPropertiesList = append(processPropertiesList, processPropertiesMap)
					}

					processesMap["process_properties"] = processPropertiesList
				}

				if processes.LineageNodeId != nil {
					processesMap["lineage_node_id"] = processes.LineageNodeId
				}

				processesList = append(processesList, processesMap)
			}

			itemsMap["processes"] = processesList
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("items", itemsList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}

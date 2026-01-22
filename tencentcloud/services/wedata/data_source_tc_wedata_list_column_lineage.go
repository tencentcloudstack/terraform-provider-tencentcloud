package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataListColumnLineage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataListColumnLineageRead,
		Schema: map[string]*schema.Schema{
			"table_unique_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Table unique ID.",
			},

			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Lineage direction INPUT|OUTPUT.",
			},

			"column_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Column name.",
			},

			"platform": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source: WEDATA|THIRD, default WEDATA.",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Lineage record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Current resource.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Entity original unique ID.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Business name: database.table|metric name|model name|field name.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Entity type\nTABLE|METRIC|MODEL|SERVICE|COLUMN.",
									},
									"lineage_node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Lineage node unique identifier.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description: table type|metric description|model description|field description.",
									},
									"platform": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source: WEDATA|THIRD\ndefault wedata.",
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
								},
							},
						},
						"relation": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Relation.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"relation_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Relation ID.",
									},
									"source_unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source unique lineage ID.",
									},
									"target_unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target unique lineage ID.",
									},
									"processes": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Lineage processing process.",
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
													Description: "Task type\n    //Scheduling task\n    SCHEDULE_TASK,\n    //Integration task\n    INTEGRATION_TASK,\n    //Third-party reporting\n    THIRD_REPORT,\n    //Data modeling\n    TABLE_MODEL,\n    //Model creates metrics\n    MODEL_METRIC,\n    //Atomic metric creates derived metric\n    METRIC_METRIC,\n    //Data service\n    DATA_SERVICE.",
												},
												"platform": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "WEDATA, THIRD.",
												},
												"process_sub_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Task subtype\n SQL_TASK,\n    //Integration real-time task lineage\n    INTEGRATED_STREAM,\n    //Integration offline task lineage\n    INTEGRATED_OFFLINE.",
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

func dataSourceTencentCloudWedataListColumnLineageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_list_column_lineage.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("table_unique_id"); ok {
		paramMap["TableUniqueId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("direction"); ok {
		paramMap["Direction"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("column_name"); ok {
		paramMap["ColumnName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("platform"); ok {
		paramMap["Platform"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.LineageNodeInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataListColumnLineageByFilter(ctx, paramMap)
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
		resourceMap := map[string]interface{}{}
		if items.Resource != nil {
			if items.Resource.ResourceUniqueId != nil {
				resourceMap["resource_unique_id"] = items.Resource.ResourceUniqueId
			}

			if items.Resource.ResourceName != nil {
				resourceMap["resource_name"] = items.Resource.ResourceName
			}

			if items.Resource.ResourceType != nil {
				resourceMap["resource_type"] = items.Resource.ResourceType
			}

			if items.Resource.LineageNodeId != nil {
				resourceMap["lineage_node_id"] = items.Resource.LineageNodeId
			}

			if items.Resource.Description != nil {
				resourceMap["description"] = items.Resource.Description
			}

			if items.Resource.Platform != nil {
				resourceMap["platform"] = items.Resource.Platform
			}

			if items.Resource.CreateTime != nil {
				resourceMap["create_time"] = items.Resource.CreateTime
			}

			if items.Resource.UpdateTime != nil {
				resourceMap["update_time"] = items.Resource.UpdateTime
			}

			resourcePropertiesList := make([]map[string]interface{}, 0, len(items.Resource.ResourceProperties))
			if items.Resource.ResourceProperties != nil {
				for _, resourceProperties := range items.Resource.ResourceProperties {
					resourcePropertiesMap := map[string]interface{}{}

					if resourceProperties.Name != nil {
						resourcePropertiesMap["name"] = resourceProperties.Name
					}

					if resourceProperties.Value != nil {
						resourcePropertiesMap["value"] = resourceProperties.Value
					}

					resourcePropertiesList = append(resourcePropertiesList, resourcePropertiesMap)
				}

				resourceMap["resource_properties"] = resourcePropertiesList
			}

			itemsMap["resource"] = []interface{}{resourceMap}
		}

		relationMap := map[string]interface{}{}
		if items.Relation != nil {
			if items.Relation.RelationId != nil {
				relationMap["relation_id"] = items.Relation.RelationId
			}

			if items.Relation.SourceUniqueId != nil {
				relationMap["source_unique_id"] = items.Relation.SourceUniqueId
			}

			if items.Relation.TargetUniqueId != nil {
				relationMap["target_unique_id"] = items.Relation.TargetUniqueId
			}

			processesList := make([]map[string]interface{}, 0, len(items.Relation.Processes))
			if items.Relation.Processes != nil {
				for _, processes := range items.Relation.Processes {
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

				relationMap["processes"] = processesList
			}

			itemsMap["relation"] = []interface{}{relationMap}
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

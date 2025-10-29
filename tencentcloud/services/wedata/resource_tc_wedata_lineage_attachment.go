package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataLineageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataLineageAttachmentCreate,
		Read:   resourceTencentCloudWedataLineageAttachmentRead,
		Delete: resourceTencentCloudWedataLineageAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"relations": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "List of lineage relationships to be registered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_unique_id": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Entity original unique ID.\\n\nNote: When lineage is for table columns, the unique ID should be passed as TableResourceUniqueId::FieldName.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Entity type.\nTABLE|METRIC|MODEL|SERVICE|COLUMN.",
									},
									"platform": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Source: WEDATA|THIRD.\nDefault is wedata.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Business name: database.table | metric name | model name | field name.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Description: table type | metric description | model description | field description.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Creation time.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Update time.",
									},
									"resource_properties": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										Description: "Resource additional extension parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Property name.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Property value.",
												},
											},
										},
									},
									"lineage_node_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Lineage node unique identifier.",
									},
								},
							},
						},
						"target": {
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_unique_id": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Entity original unique ID.\\n\nNote: When lineage is for table columns, the unique ID should be passed as TableResourceUniqueId::FieldName.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Entity type.\nTABLE|METRIC|MODEL|SERVICE|COLUMN.",
									},
									"platform": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Source: WEDATA|THIRD.\nDefault is wedata.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Business name: database.table | metric name | model name | field name.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Description: table type | metric description | model description | field description.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Creation time.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Update time.",
									},
									"resource_properties": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										Description: "Resource additional extension parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Property name.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Property value.",
												},
											},
										},
									},
									"lineage_node_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Lineage node unique identifier.",
									},
								},
							},
						},
						"processes": {
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Lineage processing process.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"process_id": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Original unique ID.",
									},
									"process_type": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Task type.\n    //Scheduled task\n    SCHEDULE_TASK,\n    //Integration task\n    INTEGRATION_TASK,\n    //Third-party reporting\n    THIRD_REPORT,\n    //Data modeling\n    TABLE_MODEL,\n    //Model creates metric\n    MODEL_METRIC,\n    //Atomic metric creates derived metric\n    METRIC_METRIC,\n    //Data service\n    DATA_SERVICE.",
									},
									"platform": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "WEDATA, THIRD.",
									},
									"process_sub_type": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Task subtype.\n SQL_TASK,\n    //Integrated real-time task lineage\n    INTEGRATED_STREAM,\n    //Integrated offline task lineage\n    INTEGRATED_OFFLINE.",
									},
									"process_properties": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										Description: "Additional extension parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Property name.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Property value.",
												},
											},
										},
									},
									"lineage_node_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Lineage task unique node ID.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudWedataLineageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_lineage_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                  = tccommon.GetLogId(tccommon.ContextNil)
		ctx                    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request                = wedatav20250806.NewRegisterLineageRequest()
		sourceResourceUniqueId string
		sourceResourceType     string
		sourcePlatform         string
		targetResourceUniqueId string
		targetResourceType     string
		targetPlatform         string
		processId              string
		processType            string
		processPlatform        string
	)

	if v, ok := d.GetOk("relations"); ok {
		for _, item := range v.([]interface{}) {
			relationsMap := item.(map[string]interface{})
			lineagePair := wedatav20250806.LineagePair{}
			if sourceMap, ok := helper.ConvertInterfacesHeadToMap(relationsMap["source"]); ok {
				lineageResouce := wedatav20250806.LineageResouce{}
				if v, ok := sourceMap["resource_unique_id"].(string); ok && v != "" {
					lineageResouce.ResourceUniqueId = helper.String(v)
					sourceResourceUniqueId = v
				}

				if v, ok := sourceMap["resource_type"].(string); ok && v != "" {
					lineageResouce.ResourceType = helper.String(v)
					sourceResourceType = v
				}

				if v, ok := sourceMap["platform"].(string); ok && v != "" {
					lineageResouce.Platform = helper.String(v)
					sourcePlatform = v
				}

				if v, ok := sourceMap["resource_name"].(string); ok && v != "" {
					lineageResouce.ResourceName = helper.String(v)
				}

				if v, ok := sourceMap["description"].(string); ok && v != "" {
					lineageResouce.Description = helper.String(v)
				}

				if v, ok := sourceMap["create_time"].(string); ok && v != "" {
					lineageResouce.CreateTime = helper.String(v)
				}

				if v, ok := sourceMap["update_time"].(string); ok && v != "" {
					lineageResouce.UpdateTime = helper.String(v)
				}

				if v, ok := sourceMap["resource_properties"]; ok {
					for _, item := range v.([]interface{}) {
						resourcePropertiesMap := item.(map[string]interface{})
						lineageProperty := wedatav20250806.LineageProperty{}
						if v, ok := resourcePropertiesMap["name"].(string); ok && v != "" {
							lineageProperty.Name = helper.String(v)
						}

						if v, ok := resourcePropertiesMap["value"].(string); ok && v != "" {
							lineageProperty.Value = helper.String(v)
						}

						lineageResouce.ResourceProperties = append(lineageResouce.ResourceProperties, &lineageProperty)
					}
				}

				if v, ok := sourceMap["lineage_node_id"].(string); ok && v != "" {
					lineageResouce.LineageNodeId = helper.String(v)
				}

				lineagePair.Source = &lineageResouce
			}

			if targetMap, ok := helper.ConvertInterfacesHeadToMap(relationsMap["target"]); ok {
				lineageResouce2 := wedatav20250806.LineageResouce{}
				if v, ok := targetMap["resource_unique_id"].(string); ok && v != "" {
					lineageResouce2.ResourceUniqueId = helper.String(v)
					targetResourceUniqueId = v
				}

				if v, ok := targetMap["resource_type"].(string); ok && v != "" {
					lineageResouce2.ResourceType = helper.String(v)
					targetResourceType = v
				}

				if v, ok := targetMap["platform"].(string); ok && v != "" {
					lineageResouce2.Platform = helper.String(v)
					targetPlatform = v
				}

				if v, ok := targetMap["resource_name"].(string); ok && v != "" {
					lineageResouce2.ResourceName = helper.String(v)
				}

				if v, ok := targetMap["description"].(string); ok && v != "" {
					lineageResouce2.Description = helper.String(v)
				}

				if v, ok := targetMap["create_time"].(string); ok && v != "" {
					lineageResouce2.CreateTime = helper.String(v)
				}

				if v, ok := targetMap["update_time"].(string); ok && v != "" {
					lineageResouce2.UpdateTime = helper.String(v)
				}

				if v, ok := targetMap["resource_properties"]; ok {
					for _, item := range v.([]interface{}) {
						resourcePropertiesMap := item.(map[string]interface{})
						lineageProperty := wedatav20250806.LineageProperty{}
						if v, ok := resourcePropertiesMap["name"].(string); ok && v != "" {
							lineageProperty.Name = helper.String(v)
						}

						if v, ok := resourcePropertiesMap["value"].(string); ok && v != "" {
							lineageProperty.Value = helper.String(v)
						}

						lineageResouce2.ResourceProperties = append(lineageResouce2.ResourceProperties, &lineageProperty)
					}
				}

				if v, ok := targetMap["lineage_node_id"].(string); ok && v != "" {
					lineageResouce2.LineageNodeId = helper.String(v)
				}

				lineagePair.Target = &lineageResouce2
			}

			if v, ok := relationsMap["processes"]; ok {
				for _, item := range v.([]interface{}) {
					processesMap := item.(map[string]interface{})
					lineageProcess := wedatav20250806.LineageProcess{}
					if v, ok := processesMap["process_id"].(string); ok && v != "" {
						lineageProcess.ProcessId = helper.String(v)
						processId = v
					}

					if v, ok := processesMap["process_type"].(string); ok && v != "" {
						lineageProcess.ProcessType = helper.String(v)
						processType = v
					}

					if v, ok := processesMap["platform"].(string); ok && v != "" {
						lineageProcess.Platform = helper.String(v)
						processPlatform = v
					}

					if v, ok := processesMap["process_sub_type"].(string); ok && v != "" {
						lineageProcess.ProcessSubType = helper.String(v)
					}

					if v, ok := processesMap["process_properties"]; ok {
						for _, item := range v.([]interface{}) {
							processPropertiesMap := item.(map[string]interface{})
							lineageProperty := wedatav20250806.LineageProperty{}
							if v, ok := processPropertiesMap["name"].(string); ok && v != "" {
								lineageProperty.Name = helper.String(v)
							}

							if v, ok := processPropertiesMap["value"].(string); ok && v != "" {
								lineageProperty.Value = helper.String(v)
							}

							lineageProcess.ProcessProperties = append(lineageProcess.ProcessProperties, &lineageProperty)
						}
					}

					if v, ok := processesMap["lineage_node_id"].(string); ok && v != "" {
						lineageProcess.LineageNodeId = helper.String(v)
					}

					lineagePair.Processes = append(lineagePair.Processes, &lineageProcess)
				}
			}

			request.Relations = append(request.Relations, &lineagePair)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().RegisterLineageWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata lineage attachment failed, Response is nil."))
		}

		if *result.Response.Data.Status != 1 {
			return resource.NonRetryableError(fmt.Errorf("Status is not 1."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata lineage attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	sourceStr := strings.Join([]string{sourceResourceUniqueId, sourceResourceType, sourcePlatform}, tccommon.COMMA_SP)
	targetStr := strings.Join([]string{targetResourceUniqueId, targetResourceType, targetPlatform}, tccommon.COMMA_SP)
	processStr := strings.Join([]string{processId, processType, processPlatform}, tccommon.COMMA_SP)
	d.SetId(strings.Join([]string{sourceStr, targetStr, processStr}, tccommon.FILED_SP))
	return resourceTencentCloudWedataLineageAttachmentRead(d, meta)
}

func resourceTencentCloudWedataLineageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_lineage_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	sourceObj := idSplit[0]
	targetObj := idSplit[1]
	processObj := idSplit[2]

	sourceObjSplit := strings.Split(sourceObj, tccommon.COMMA_SP)
	if len(sourceObjSplit) != 3 {
		return fmt.Errorf("source ID is broken,%s", d.Id())
	}

	targetObjSplit := strings.Split(targetObj, tccommon.COMMA_SP)
	if len(targetObjSplit) != 3 {
		return fmt.Errorf("target ID is broken,%s", d.Id())
	}

	processObjSplit := strings.Split(processObj, tccommon.COMMA_SP)
	if len(processObjSplit) != 3 {
		return fmt.Errorf("process ID is broken,%s", d.Id())
	}

	sourceResourceUniqueId := sourceObjSplit[0]
	sourceResourceType := sourceObjSplit[1]
	sourcePlatform := sourceObjSplit[2]

	targetResourceUniqueId := targetObjSplit[0]
	targetResourceType := targetObjSplit[1]
	targetPlatform := targetObjSplit[2]

	processId := processObjSplit[0]
	processType := processObjSplit[1]
	processPlatform := processObjSplit[2]

	has, err := service.DescribeWedataLineageAttachmentById(ctx, sourceResourceUniqueId, sourceResourceType, sourcePlatform, targetResourceUniqueId, targetResourceType, targetPlatform, processId, processType, processPlatform)
	if err != nil {
		return err
	}

	if !has {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_lineage_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTencentCloudWedataLineageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_lineage_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wedatav20250806.NewDeleteLineageRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	sourceObj := idSplit[0]
	targetObj := idSplit[1]
	processObj := idSplit[2]

	sourceObjSplit := strings.Split(sourceObj, tccommon.COMMA_SP)
	if len(sourceObjSplit) != 3 {
		return fmt.Errorf("source ID is broken,%s", d.Id())
	}

	targetObjSplit := strings.Split(targetObj, tccommon.COMMA_SP)
	if len(targetObjSplit) != 3 {
		return fmt.Errorf("target ID is broken,%s", d.Id())
	}

	processObjSplit := strings.Split(processObj, tccommon.COMMA_SP)
	if len(processObjSplit) != 3 {
		return fmt.Errorf("process ID is broken,%s", d.Id())
	}

	sourceResourceUniqueId := sourceObjSplit[0]
	sourceResourceType := sourceObjSplit[1]
	sourcePlatform := sourceObjSplit[2]

	targetResourceUniqueId := targetObjSplit[0]
	targetResourceType := targetObjSplit[1]
	targetPlatform := targetObjSplit[2]

	processId := processObjSplit[0]
	processType := processObjSplit[1]
	processPlatform := processObjSplit[2]

	request.Relations = []*wedatav20250806.LineagePair{
		{
			Source: &wedatav20250806.LineageResouce{
				ResourceUniqueId: &sourceResourceUniqueId,
				ResourceType:     &sourceResourceType,
				Platform:         &sourcePlatform,
			},
			Target: &wedatav20250806.LineageResouce{
				ResourceUniqueId: &targetResourceUniqueId,
				ResourceType:     &targetResourceType,
				Platform:         &targetPlatform,
			},
			Processes: []*wedatav20250806.LineageProcess{
				{
					ProcessId:   &processId,
					ProcessType: &processType,
					Platform:    &processPlatform,
				},
			},
		},
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteLineageWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata lineage attachment failed, Response is nil."))
		}

		if *result.Response.Data.Status != 1 {
			return resource.NonRetryableError(fmt.Errorf("Status is not 1."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata lineage attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

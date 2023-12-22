package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataIntegrationTaskNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataIntegrationTaskNodeCreate,
		Read:   resourceTencentCloudWedataIntegrationTaskNodeRead,
		Update: resourceTencentCloudWedataIntegrationTaskNodeUpdate,
		Delete: resourceTencentCloudWedataIntegrationTaskNodeDelete,

		Schema: map[string]*schema.Schema{
			// create
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The task id to which the node belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Node Name.",
			},
			"node_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Node type: INPUT, OUTPUT, JOIN, FILTER, TRANSFORM.",
			},
			"data_source_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data source type: MYSQL, POSTGRE, ORACLE, SQLSERVER, FTP, HIVE, HDFS, ICEBERG, KAFKA, HBASE, SPARK, TBASE, DB2, DM, GAUSSDB, GBASE, IMPALA, ES, S3_DATAINSIGHT, GREENPLUM, PHOENIX, SAP_HANA, SFTP, OCEANBASE, CLICKHOUSE, KUDU, VERTICA, REDIS, COS, DLC, DORIS, CKAFKA, DTS_KAFKA, S3, CDW, TDSQLC, TDSQL, MONGODB, SYBASE, REST_API, StarRocks, TCHOUSE_X.",
			},
			"task_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Task type, 201: real-time task, 202: offline task.",
			},
			"task_mode": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Task display mode, 0: canvas mode, 1: form mode.",
			},
			//"config": {
			//	Type:        schema.TypeList,
			//	Required:    true,
			//	Description: "Node configuration information.",
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"name": {
			//				Type:        schema.TypeString,
			//				Optional:    true,
			//				Description: "Configuration name.",
			//			},
			//			"value": {
			//				Type:        schema.TypeString,
			//				Optional:    true,
			//				Description: "Configuration value.",
			//			},
			//		},
			//	},
			//},
			// modify
			"node_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Node information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datasource_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Datasource ID.",
						},
						"config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Node configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration name.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration value.",
									},
								},
							},
						},
						"ext_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Node extension configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration name.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration value.",
									},
								},
							},
						},
						"schema": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Schema information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Schema ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Schema name.",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Schema type.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Schema value.",
									},
									"properties": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Schema extended attributes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Attributes name.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Attributes value.",
												},
											},
										},
									},
									"alias": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Schema alias.",
									},
									"comment": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Schema comment.",
									},
								},
							},
						},
						"node_mapping": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Node mapping.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Source node ID.",
									},
									"sink_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Sink node ID.",
									},
									"source_schema": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Source node schema information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema ID.",
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema name.",
												},
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema type.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema value.",
												},
												"properties": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Schema extended attributes.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Attributes name.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Attributes value.",
															},
														},
													},
												},
												"alias": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema alias.",
												},
												"comment": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema comment.",
												},
											},
										},
									},
									"schema_mappings": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Schema mapping information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source_schema_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema ID from source node.",
												},
												"sink_schema_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema ID from sink node.",
												},
											},
										},
									},
									"ext_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Node extension configuration information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Configuration name.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Configuration value.",
												},
											},
										},
									},
								},
							},
						},
						"app_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "User App Id.",
						},
						"creator_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Creator User ID.",
						},
						"operator_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Operator User ID.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Owner User ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Create time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Update time.",
						},
					},
				},
			},
			"node_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node ID.",
			},
		},
	}
}

func resourceTencentCloudWedataIntegrationTaskNodeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_integration_task_node.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = wedata.NewCreateIntegrationNodeRequest()
		response      = wedata.NewCreateIntegrationNodeResponse()
		modifyRequest = wedata.NewModifyIntegrationNodeRequest()
		projectId     string
		taskId        string
		nodeId        string
		taskType      int
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOkExists("task_type"); ok {
		request.TaskType = helper.IntUint64(v.(int))
		taskType = v.(int)
	}

	createIntegrationNodeInfo := wedata.IntegrationNodeInfo{}
	if v, ok := d.GetOk("task_id"); ok {
		createIntegrationNodeInfo.TaskId = helper.String(v.(string))
		taskId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		createIntegrationNodeInfo.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("node_type"); ok {
		createIntegrationNodeInfo.NodeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_source_type"); ok {
		createIntegrationNodeInfo.DataSourceType = helper.String(v.(string))
	}

	//if v, ok := d.GetOk("config"); ok {
	//	for _, item := range v.([]interface{}) {
	//		configMap := item.(map[string]interface{})
	//		recordField := wedata.RecordField{}
	//		if v, ok := configMap["name"]; ok {
	//			recordField.Name = helper.String(v.(string))
	//		}
	//
	//		if v, ok := configMap["value"]; ok {
	//			recordField.Value = helper.String(v.(string))
	//		}
	//
	//		createIntegrationNodeInfo.Config = append(createIntegrationNodeInfo.Config, &recordField)
	//	}
	//}

	// create
	request.NodeInfo = &createIntegrationNodeInfo
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().CreateIntegrationNode(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata integrationTaskNode failed, reason:%+v", logId, err)
		return err
	}

	nodeId = *response.Response.Id
	d.SetId(strings.Join([]string{projectId, nodeId}, tccommon.FILED_SP))

	// modify
	modifyRequest.ProjectId = &projectId
	modifyRequest.TaskType = helper.IntUint64(taskType)
	if v, ok := d.GetOkExists("task_mode"); ok {
		modifyRequest.TaskMode = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "node_info"); ok {
		integrationNodeInfo := wedata.IntegrationNodeInfo{}
		integrationNodeInfo.Id = &nodeId
		integrationNodeInfo.TaskId = &taskId

		if v, ok := dMap["datasource_id"]; ok {
			integrationNodeInfo.DatasourceId = helper.String(v.(string))
		}

		if v, ok := dMap["config"]; ok {
			for _, item := range v.([]interface{}) {
				configMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := configMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := configMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationNodeInfo.Config = append(integrationNodeInfo.Config, &recordField)
			}
		}

		if v, ok := dMap["ext_config"]; ok {
			for _, item := range v.([]interface{}) {
				extConfigMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := extConfigMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := extConfigMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationNodeInfo.ExtConfig = append(integrationNodeInfo.ExtConfig, &recordField)
			}
		}

		if v, ok := dMap["schema"]; ok {
			for _, item := range v.([]interface{}) {
				schemaMap := item.(map[string]interface{})
				integrationNodeSchema := wedata.IntegrationNodeSchema{}
				if v, ok := schemaMap["id"]; ok {
					integrationNodeSchema.Id = helper.String(v.(string))
				}

				if v, ok := schemaMap["name"]; ok {
					integrationNodeSchema.Name = helper.String(v.(string))
				}

				if v, ok := schemaMap["type"]; ok {
					integrationNodeSchema.Type = helper.String(v.(string))
				}

				if v, ok := schemaMap["value"]; ok {
					integrationNodeSchema.Value = helper.String(v.(string))
				}

				if v, ok := schemaMap["properties"]; ok {
					for _, item := range v.([]interface{}) {
						propertiesMap := item.(map[string]interface{})
						recordField := wedata.RecordField{}
						if v, ok := propertiesMap["name"]; ok {
							recordField.Name = helper.String(v.(string))
						}

						if v, ok := propertiesMap["value"]; ok {
							recordField.Value = helper.String(v.(string))
						}

						integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
					}
				}

				if v, ok := schemaMap["alias"]; ok {
					integrationNodeSchema.Alias = helper.String(v.(string))
				}

				if v, ok := schemaMap["comment"]; ok {
					integrationNodeSchema.Comment = helper.String(v.(string))
				}

				integrationNodeInfo.Schema = append(integrationNodeInfo.Schema, &integrationNodeSchema)
			}
		}

		if nodeMappingMap, ok := helper.InterfaceToMap(dMap, "node_mapping"); ok {
			integrationNodeMapping := wedata.IntegrationNodeMapping{}
			if v, ok := nodeMappingMap["source_id"]; ok {
				integrationNodeMapping.SourceId = helper.String(v.(string))
			}

			if v, ok := nodeMappingMap["sink_id"]; ok {
				integrationNodeMapping.SinkId = helper.String(v.(string))
			}

			if v, ok := nodeMappingMap["source_schema"]; ok {
				for _, item := range v.([]interface{}) {
					sourceSchemaMap := item.(map[string]interface{})
					integrationNodeSchema := wedata.IntegrationNodeSchema{}
					if v, ok := sourceSchemaMap["id"]; ok {
						integrationNodeSchema.Id = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["name"]; ok {
						integrationNodeSchema.Name = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["type"]; ok {
						integrationNodeSchema.Type = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["value"]; ok {
						integrationNodeSchema.Value = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["properties"]; ok {
						for _, item := range v.([]interface{}) {
							propertiesMap := item.(map[string]interface{})
							recordField := wedata.RecordField{}
							if v, ok := propertiesMap["name"]; ok {
								recordField.Name = helper.String(v.(string))
							}

							if v, ok := propertiesMap["value"]; ok {
								recordField.Value = helper.String(v.(string))
							}

							integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
						}
					}

					if v, ok := sourceSchemaMap["alias"]; ok {
						integrationNodeSchema.Alias = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["comment"]; ok {
						integrationNodeSchema.Comment = helper.String(v.(string))
					}

					integrationNodeMapping.SourceSchema = append(integrationNodeMapping.SourceSchema, &integrationNodeSchema)
				}
			}

			if v, ok := nodeMappingMap["schema_mappings"]; ok {
				for _, item := range v.([]interface{}) {
					schemaMappingsMap := item.(map[string]interface{})
					integrationNodeSchemaMapping := wedata.IntegrationNodeSchemaMapping{}
					if v, ok := schemaMappingsMap["source_schema_id"]; ok {
						integrationNodeSchemaMapping.SourceSchemaId = helper.String(v.(string))
					}

					if v, ok := schemaMappingsMap["sink_schema_id"]; ok {
						integrationNodeSchemaMapping.SinkSchemaId = helper.String(v.(string))
					}

					integrationNodeMapping.SchemaMappings = append(integrationNodeMapping.SchemaMappings, &integrationNodeSchemaMapping)
				}
			}

			if v, ok := nodeMappingMap["ext_config"]; ok {
				for _, item := range v.([]interface{}) {
					extConfigMap := item.(map[string]interface{})
					recordField := wedata.RecordField{}
					if v, ok := extConfigMap["name"]; ok {
						recordField.Name = helper.String(v.(string))
					}

					if v, ok := extConfigMap["value"]; ok {
						recordField.Value = helper.String(v.(string))
					}

					integrationNodeMapping.ExtConfig = append(integrationNodeMapping.ExtConfig, &recordField)
				}
			}

			integrationNodeInfo.NodeMapping = &integrationNodeMapping
		}

		if v, ok := dMap["app_id"]; ok {
			integrationNodeInfo.AppId = helper.String(v.(string))
		}

		if v, ok := dMap["creator_uin"]; ok {
			integrationNodeInfo.CreatorUin = helper.String(v.(string))
		}

		if v, ok := dMap["operator_uin"]; ok {
			integrationNodeInfo.OperatorUin = helper.String(v.(string))
		}

		if v, ok := dMap["owner_uin"]; ok {
			integrationNodeInfo.OwnerUin = helper.String(v.(string))
		}

		if v, ok := dMap["create_time"]; ok {
			integrationNodeInfo.CreateTime = helper.String(v.(string))
		}

		if v, ok := dMap["update_time"]; ok {
			integrationNodeInfo.UpdateTime = helper.String(v.(string))
		}

		modifyRequest.NodeInfo = &integrationNodeInfo
	}

	// create
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().ModifyIntegrationNode(modifyRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyRequest.GetAction(), modifyRequest.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata integrationTaskNode failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegrationTaskNodeRead(d, meta)
}

func resourceTencentCloudWedataIntegrationTaskNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_integration_task_node.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	nodeId := idSplit[1]

	integrationTaskNode, err := service.DescribeWedataIntegrationTaskNodeById(ctx, projectId, nodeId)
	if err != nil {
		return err
	}

	if integrationTaskNode == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataIntegrationTaskNode` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("node_id", nodeId)

	if integrationTaskNode.NodeInfo != nil {
		nodeInfoMap := map[string]interface{}{}

		if integrationTaskNode.NodeInfo.TaskId != nil {
			_ = d.Set("task_id", integrationTaskNode.NodeInfo.TaskId)
		}

		if integrationTaskNode.NodeInfo.Name != nil {
			_ = d.Set("name", integrationTaskNode.NodeInfo.Name)
		}

		if integrationTaskNode.NodeInfo.NodeType != nil {
			_ = d.Set("node_type", integrationTaskNode.NodeInfo.NodeType)
		}

		if integrationTaskNode.NodeInfo.DataSourceType != nil {
			_ = d.Set("data_source_type", integrationTaskNode.NodeInfo.DataSourceType)
		}

		if integrationTaskNode.NodeInfo.DatasourceId != nil {
			nodeInfoMap["datasource_id"] = integrationTaskNode.NodeInfo.DatasourceId
		}

		if integrationTaskNode.NodeInfo.Config != nil {
			configList := []interface{}{}
			for _, config := range integrationTaskNode.NodeInfo.Config {
				configMap := map[string]interface{}{}

				if config.Name != nil {
					configMap["name"] = config.Name
				}

				if config.Value != nil {
					configMap["value"] = config.Value
				}

				configList = append(configList, configMap)
			}

			nodeInfoMap["config"] = configList
		}

		if integrationTaskNode.NodeInfo.ExtConfig != nil {
			extConfigList := []interface{}{}
			for _, extConfig := range integrationTaskNode.NodeInfo.ExtConfig {
				extConfigMap := map[string]interface{}{}

				if extConfig.Name != nil {
					extConfigMap["name"] = extConfig.Name
				}

				if extConfig.Value != nil {
					extConfigMap["value"] = extConfig.Value
				}

				extConfigList = append(extConfigList, extConfigMap)
			}

			nodeInfoMap["ext_config"] = extConfigList
		}

		if integrationTaskNode.NodeInfo.Schema != nil {
			schemaList := []interface{}{}
			for _, Nschema := range integrationTaskNode.NodeInfo.Schema {
				schemaMap := map[string]interface{}{}

				if Nschema.Id != nil {
					schemaMap["id"] = Nschema.Id
				}

				if Nschema.Name != nil {
					schemaMap["name"] = Nschema.Name
				}

				if Nschema.Type != nil {
					schemaMap["type"] = Nschema.Type
				}

				if Nschema.Value != nil {
					schemaMap["value"] = Nschema.Value
				}

				if Nschema.Properties != nil {
					propertiesList := []interface{}{}
					for _, properties := range Nschema.Properties {
						propertiesMap := map[string]interface{}{}

						if properties.Name != nil {
							propertiesMap["name"] = properties.Name
						}

						if properties.Value != nil {
							propertiesMap["value"] = properties.Value
						}

						propertiesList = append(propertiesList, propertiesMap)
					}

					schemaMap["properties"] = propertiesList
				}

				if Nschema.Alias != nil {
					schemaMap["alias"] = Nschema.Alias
				}

				if Nschema.Comment != nil {
					schemaMap["comment"] = Nschema.Comment
				}

				schemaList = append(schemaList, schemaMap)
			}

			nodeInfoMap["schema"] = schemaList
		}

		if integrationTaskNode.NodeInfo.NodeMapping != nil {
			nodeMappingMap := map[string]interface{}{}

			if integrationTaskNode.NodeInfo.NodeMapping.SourceId != nil {
				nodeMappingMap["source_id"] = integrationTaskNode.NodeInfo.NodeMapping.SourceId
			}

			if integrationTaskNode.NodeInfo.NodeMapping.SinkId != nil {
				nodeMappingMap["sink_id"] = integrationTaskNode.NodeInfo.NodeMapping.SinkId
			}

			if integrationTaskNode.NodeInfo.NodeMapping.SourceSchema != nil {
				sourceSchemaList := []interface{}{}
				for _, sourceSchema := range integrationTaskNode.NodeInfo.NodeMapping.SourceSchema {
					sourceSchemaMap := map[string]interface{}{}

					if sourceSchema.Id != nil {
						sourceSchemaMap["id"] = sourceSchema.Id
					}

					if sourceSchema.Name != nil {
						sourceSchemaMap["name"] = sourceSchema.Name
					}

					if sourceSchema.Type != nil {
						sourceSchemaMap["type"] = sourceSchema.Type
					}

					if sourceSchema.Value != nil {
						sourceSchemaMap["value"] = sourceSchema.Value
					}

					if sourceSchema.Properties != nil {
						propertiesList := []interface{}{}
						for _, properties := range sourceSchema.Properties {
							propertiesMap := map[string]interface{}{}

							if properties.Name != nil {
								propertiesMap["name"] = properties.Name
							}

							if properties.Value != nil {
								propertiesMap["value"] = properties.Value
							}

							propertiesList = append(propertiesList, propertiesMap)
						}

						sourceSchemaMap["properties"] = propertiesList
					}

					if sourceSchema.Alias != nil {
						sourceSchemaMap["alias"] = sourceSchema.Alias
					}

					if sourceSchema.Comment != nil {
						sourceSchemaMap["comment"] = sourceSchema.Comment
					}

					sourceSchemaList = append(sourceSchemaList, sourceSchemaMap)
				}

				nodeMappingMap["source_schema"] = sourceSchemaList
			}

			if integrationTaskNode.NodeInfo.NodeMapping.SchemaMappings != nil {
				schemaMappingsList := []interface{}{}
				for _, schemaMappings := range integrationTaskNode.NodeInfo.NodeMapping.SchemaMappings {
					schemaMappingsMap := map[string]interface{}{}

					if schemaMappings.SourceSchemaId != nil {
						schemaMappingsMap["source_schema_id"] = schemaMappings.SourceSchemaId
					}

					if schemaMappings.SinkSchemaId != nil {
						schemaMappingsMap["sink_schema_id"] = schemaMappings.SinkSchemaId
					}

					schemaMappingsList = append(schemaMappingsList, schemaMappingsMap)
				}

				nodeMappingMap["schema_mappings"] = schemaMappingsList
			}

			if integrationTaskNode.NodeInfo.NodeMapping.ExtConfig != nil {
				extConfigList := []interface{}{}
				for _, extConfig := range integrationTaskNode.NodeInfo.NodeMapping.ExtConfig {
					extConfigMap := map[string]interface{}{}

					if extConfig.Name != nil {
						extConfigMap["name"] = extConfig.Name
					}

					if extConfig.Value != nil {
						extConfigMap["value"] = extConfig.Value
					}

					extConfigList = append(extConfigList, extConfigMap)
				}

				nodeMappingMap["ext_config"] = extConfigList
			}

			nodeInfoMap["node_mapping"] = []interface{}{nodeMappingMap}
		}

		if integrationTaskNode.NodeInfo.AppId != nil {
			nodeInfoMap["app_id"] = integrationTaskNode.NodeInfo.AppId
		}

		if integrationTaskNode.NodeInfo.CreatorUin != nil {
			nodeInfoMap["creator_uin"] = integrationTaskNode.NodeInfo.CreatorUin
		}

		if integrationTaskNode.NodeInfo.OperatorUin != nil {
			nodeInfoMap["operator_uin"] = integrationTaskNode.NodeInfo.OperatorUin
		}

		if integrationTaskNode.NodeInfo.OwnerUin != nil {
			nodeInfoMap["owner_uin"] = integrationTaskNode.NodeInfo.OwnerUin
		}

		if integrationTaskNode.NodeInfo.CreateTime != nil {
			nodeInfoMap["create_time"] = integrationTaskNode.NodeInfo.CreateTime
		}

		if integrationTaskNode.NodeInfo.UpdateTime != nil {
			nodeInfoMap["update_time"] = integrationTaskNode.NodeInfo.UpdateTime
		}

		_ = d.Set("node_info", []interface{}{nodeInfoMap})
	}

	return nil
}

func resourceTencentCloudWedataIntegrationTaskNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_integration_task_node.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = wedata.NewModifyIntegrationNodeRequest()
	)

	immutableArgs := []string{"project_id", "task_id", "task_type", "task_mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	nodeId := idSplit[1]

	request.ProjectId = &projectId
	if v, ok := d.GetOkExists("task_type"); ok {
		request.TaskType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("task_mode"); ok {
		request.TaskMode = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "node_info"); ok {
		integrationNodeInfo := wedata.IntegrationNodeInfo{}
		integrationNodeInfo.Id = &nodeId

		if v, ok := d.GetOk("task_id"); ok {
			integrationNodeInfo.TaskId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("name"); ok {
			integrationNodeInfo.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("node_type"); ok {
			integrationNodeInfo.NodeType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("data_source_type"); ok {
			integrationNodeInfo.DataSourceType = helper.String(v.(string))
		}

		if v, ok := dMap["datasource_id"]; ok {
			integrationNodeInfo.DatasourceId = helper.String(v.(string))
		}

		if v, ok := dMap["config"]; ok {
			for _, item := range v.([]interface{}) {
				configMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := configMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := configMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationNodeInfo.Config = append(integrationNodeInfo.Config, &recordField)
			}
		}

		if v, ok := dMap["ext_config"]; ok {
			for _, item := range v.([]interface{}) {
				extConfigMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := extConfigMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := extConfigMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationNodeInfo.ExtConfig = append(integrationNodeInfo.ExtConfig, &recordField)
			}
		}

		if v, ok := dMap["schema"]; ok {
			for _, item := range v.([]interface{}) {
				schemaMap := item.(map[string]interface{})
				integrationNodeSchema := wedata.IntegrationNodeSchema{}
				if v, ok := schemaMap["id"]; ok {
					integrationNodeSchema.Id = helper.String(v.(string))
				}

				if v, ok := schemaMap["name"]; ok {
					integrationNodeSchema.Name = helper.String(v.(string))
				}

				if v, ok := schemaMap["type"]; ok {
					integrationNodeSchema.Type = helper.String(v.(string))
				}

				if v, ok := schemaMap["value"]; ok {
					integrationNodeSchema.Value = helper.String(v.(string))
				}

				if v, ok := schemaMap["properties"]; ok {
					for _, item := range v.([]interface{}) {
						propertiesMap := item.(map[string]interface{})
						recordField := wedata.RecordField{}
						if v, ok := propertiesMap["name"]; ok {
							recordField.Name = helper.String(v.(string))
						}

						if v, ok := propertiesMap["value"]; ok {
							recordField.Value = helper.String(v.(string))
						}

						integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
					}
				}

				if v, ok := schemaMap["alias"]; ok {
					integrationNodeSchema.Alias = helper.String(v.(string))
				}

				if v, ok := schemaMap["comment"]; ok {
					integrationNodeSchema.Comment = helper.String(v.(string))
				}

				integrationNodeInfo.Schema = append(integrationNodeInfo.Schema, &integrationNodeSchema)
			}
		}

		if nodeMappingMap, ok := helper.InterfaceToMap(dMap, "node_mapping"); ok {
			integrationNodeMapping := wedata.IntegrationNodeMapping{}
			if v, ok := nodeMappingMap["source_id"]; ok {
				integrationNodeMapping.SourceId = helper.String(v.(string))
			}

			if v, ok := nodeMappingMap["sink_id"]; ok {
				integrationNodeMapping.SinkId = helper.String(v.(string))
			}

			if v, ok := nodeMappingMap["source_schema"]; ok {
				for _, item := range v.([]interface{}) {
					sourceSchemaMap := item.(map[string]interface{})
					integrationNodeSchema := wedata.IntegrationNodeSchema{}
					if v, ok := sourceSchemaMap["id"]; ok {
						integrationNodeSchema.Id = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["name"]; ok {
						integrationNodeSchema.Name = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["type"]; ok {
						integrationNodeSchema.Type = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["value"]; ok {
						integrationNodeSchema.Value = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["properties"]; ok {
						for _, item := range v.([]interface{}) {
							propertiesMap := item.(map[string]interface{})
							recordField := wedata.RecordField{}
							if v, ok := propertiesMap["name"]; ok {
								recordField.Name = helper.String(v.(string))
							}

							if v, ok := propertiesMap["value"]; ok {
								recordField.Value = helper.String(v.(string))
							}

							integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
						}
					}

					if v, ok := sourceSchemaMap["alias"]; ok {
						integrationNodeSchema.Alias = helper.String(v.(string))
					}

					if v, ok := sourceSchemaMap["comment"]; ok {
						integrationNodeSchema.Comment = helper.String(v.(string))
					}

					integrationNodeMapping.SourceSchema = append(integrationNodeMapping.SourceSchema, &integrationNodeSchema)
				}
			}

			if v, ok := nodeMappingMap["schema_mappings"]; ok {
				for _, item := range v.([]interface{}) {
					schemaMappingsMap := item.(map[string]interface{})
					integrationNodeSchemaMapping := wedata.IntegrationNodeSchemaMapping{}
					if v, ok := schemaMappingsMap["source_schema_id"]; ok {
						integrationNodeSchemaMapping.SourceSchemaId = helper.String(v.(string))
					}

					if v, ok := schemaMappingsMap["sink_schema_id"]; ok {
						integrationNodeSchemaMapping.SinkSchemaId = helper.String(v.(string))
					}

					integrationNodeMapping.SchemaMappings = append(integrationNodeMapping.SchemaMappings, &integrationNodeSchemaMapping)
				}
			}

			if v, ok := nodeMappingMap["ext_config"]; ok {
				for _, item := range v.([]interface{}) {
					extConfigMap := item.(map[string]interface{})
					recordField := wedata.RecordField{}
					if v, ok := extConfigMap["name"]; ok {
						recordField.Name = helper.String(v.(string))
					}

					if v, ok := extConfigMap["value"]; ok {
						recordField.Value = helper.String(v.(string))
					}

					integrationNodeMapping.ExtConfig = append(integrationNodeMapping.ExtConfig, &recordField)
				}
			}

			integrationNodeInfo.NodeMapping = &integrationNodeMapping
		}

		if v, ok := dMap["app_id"]; ok {
			integrationNodeInfo.AppId = helper.String(v.(string))
		}

		if v, ok := dMap["creator_uin"]; ok {
			integrationNodeInfo.CreatorUin = helper.String(v.(string))
		}

		if v, ok := dMap["operator_uin"]; ok {
			integrationNodeInfo.OperatorUin = helper.String(v.(string))
		}

		if v, ok := dMap["owner_uin"]; ok {
			integrationNodeInfo.OwnerUin = helper.String(v.(string))
		}

		if v, ok := dMap["create_time"]; ok {
			integrationNodeInfo.CreateTime = helper.String(v.(string))
		}

		if v, ok := dMap["update_time"]; ok {
			integrationNodeInfo.UpdateTime = helper.String(v.(string))
		}

		request.NodeInfo = &integrationNodeInfo
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().ModifyIntegrationNode(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update wedata integrationTaskNode failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegrationTaskNodeRead(d, meta)
}

func resourceTencentCloudWedataIntegrationTaskNodeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_integration_task_node.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	nodeId := idSplit[1]

	if err := service.DeleteWedataIntegrationTaskNodeById(ctx, projectId, nodeId); err != nil {
		return err
	}

	return nil
}

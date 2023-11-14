/*
Provides a resource to create a wedata integration_task_node

Example Usage

```hcl
resource "tencentcloud_wedata_integration_task_node" "integration_task_node" {
  node_info {
		id = ""
		task_id = "j84cc717e-215b-4960-9575-898586bae37f"
		name = "input_name"
		node_type = "INPUT"
		data_source_type = "MYSQL"
		description = "Node for test"
		datasource_id = "100"
		config {
			name = "Database"
			value = "db"
		}
		ext_config {
			name = "x"
			value = "320"
		}
		schema {
			id = "796598528"
			name = "col_name"
			type = "string"
			value = "1"
			properties {
				name = "name"
				value = "value"
			}
			alias = "name"
			comment = "comment"
		}
		node_mapping {
			source_id = "10"
			sink_id = "11"
			source_schema {
				id = "796598528"
				name = "col_name"
				type = "string"
				value = "1"
				properties {
					name = "name"
					value = "value"
				}
				alias = "name"
				comment = "comment"
			}
			schema_mappings {
				source_schema_id = "200"
				sink_schema_id = "300"
			}
			ext_config {
				name = "x"
				value = "320"
			}
		}
		app_id = "1315000000"
		project_id = "1455251608631480391"
		creator_uin = "100028448000"
		operator_uin = "100028448000"
		owner_uin = "100028448000"
		create_time = "2023-10-17 18:02:46"
		update_time = "2023-10-17 18:02:46"

  }
  project_id = "1455251608631480391"
  task_type = 201
}
```

Import

wedata integration_task_node can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_integration_task_node.integration_task_node integration_task_node_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWedataIntegration_task_node() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataIntegration_task_nodeCreate,
		Read:   resourceTencentCloudWedataIntegration_task_nodeRead,
		Update: resourceTencentCloudWedataIntegration_task_nodeUpdate,
		Delete: resourceTencentCloudWedataIntegration_task_nodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"node_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Node information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Node ID.",
						},
						"task_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The task id to which the node belongs.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Node Name.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Node type: INPUT,OUTPUT,JOIN,FILTER,TRANSFORM.",
						},
						"data_source_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data source type: MYSQL, POSTGRE, ORACLE, SQLSERVER, FTP, HIVE, HDFS, ICEBERG, KAFKA, HBASE, SPARK, TBASE, DB2, DM, GAUSSDB, GBASE, IMPALA, ES, S3_DATAINSIGHT, GREENPLUM, PHOENIX, SAP_HANA, SFTP, OCEANBASE, CLICKHOUSE, KUDU, VERTICA, REDIS, COS, DLC, DORIS, CKAFKA, DTS_KAFKA, S3, CDW, TDSQLC, TDSQL, MONGODB, SYBASE, REST_API, StarRocks, TCHOUSE_X.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Node Description.",
						},
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
							Description: "User App Id.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Project ID.",
						},
						"creator_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Creator User ID.",
						},
						"operator_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Operator User ID.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Owner User ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Create time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
					},
				},
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"task_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Task type, 201: real-time task, 202: offline task.",
			},
		},
	}
}

func resourceTencentCloudWedataIntegration_task_nodeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_task_node.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = wedata.NewCreateIntegrationNodeRequest()
		response = wedata.NewCreateIntegrationNodeResponse()
		id       string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "node_info"); ok {
		integrationNodeInfo := wedata.IntegrationNodeInfo{}
		if v, ok := dMap["id"]; ok {
			integrationNodeInfo.Id = helper.String(v.(string))
		}
		if v, ok := dMap["task_id"]; ok {
			integrationNodeInfo.TaskId = helper.String(v.(string))
		}
		if v, ok := dMap["name"]; ok {
			integrationNodeInfo.Name = helper.String(v.(string))
		}
		if v, ok := dMap["node_type"]; ok {
			integrationNodeInfo.NodeType = helper.String(v.(string))
		}
		if v, ok := dMap["data_source_type"]; ok {
			integrationNodeInfo.DataSourceType = helper.String(v.(string))
		}
		if v, ok := dMap["description"]; ok {
			integrationNodeInfo.Description = helper.String(v.(string))
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
		if v, ok := dMap["project_id"]; ok {
			integrationNodeInfo.ProjectId = helper.String(v.(string))
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

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("task_type"); ok {
		request.TaskType = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateIntegrationNode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata integration_task_node failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.Id
	d.SetId(id)

	return resourceTencentCloudWedataIntegration_task_nodeRead(d, meta)
}

func resourceTencentCloudWedataIntegration_task_nodeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_task_node.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	integration_task_nodeId := d.Id()

	integration_task_node, err := service.DescribeWedataIntegration_task_nodeById(ctx, id)
	if err != nil {
		return err
	}

	if integration_task_node == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataIntegration_task_node` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if integration_task_node.NodeInfo != nil {
		nodeInfoMap := map[string]interface{}{}

		if integration_task_node.NodeInfo.Id != nil {
			nodeInfoMap["id"] = integration_task_node.NodeInfo.Id
		}

		if integration_task_node.NodeInfo.TaskId != nil {
			nodeInfoMap["task_id"] = integration_task_node.NodeInfo.TaskId
		}

		if integration_task_node.NodeInfo.Name != nil {
			nodeInfoMap["name"] = integration_task_node.NodeInfo.Name
		}

		if integration_task_node.NodeInfo.NodeType != nil {
			nodeInfoMap["node_type"] = integration_task_node.NodeInfo.NodeType
		}

		if integration_task_node.NodeInfo.DataSourceType != nil {
			nodeInfoMap["data_source_type"] = integration_task_node.NodeInfo.DataSourceType
		}

		if integration_task_node.NodeInfo.Description != nil {
			nodeInfoMap["description"] = integration_task_node.NodeInfo.Description
		}

		if integration_task_node.NodeInfo.DatasourceId != nil {
			nodeInfoMap["datasource_id"] = integration_task_node.NodeInfo.DatasourceId
		}

		if integration_task_node.NodeInfo.Config != nil {
			configList := []interface{}{}
			for _, config := range integration_task_node.NodeInfo.Config {
				configMap := map[string]interface{}{}

				if config.Name != nil {
					configMap["name"] = config.Name
				}

				if config.Value != nil {
					configMap["value"] = config.Value
				}

				configList = append(configList, configMap)
			}

			nodeInfoMap["config"] = []interface{}{configList}
		}

		if integration_task_node.NodeInfo.ExtConfig != nil {
			extConfigList := []interface{}{}
			for _, extConfig := range integration_task_node.NodeInfo.ExtConfig {
				extConfigMap := map[string]interface{}{}

				if extConfig.Name != nil {
					extConfigMap["name"] = extConfig.Name
				}

				if extConfig.Value != nil {
					extConfigMap["value"] = extConfig.Value
				}

				extConfigList = append(extConfigList, extConfigMap)
			}

			nodeInfoMap["ext_config"] = []interface{}{extConfigList}
		}

		if integration_task_node.NodeInfo.Schema != nil {
			schemaList := []interface{}{}
			for _, schema := range integration_task_node.NodeInfo.Schema {
				schemaMap := map[string]interface{}{}

				if schema.Id != nil {
					schemaMap["id"] = schema.Id
				}

				if schema.Name != nil {
					schemaMap["name"] = schema.Name
				}

				if schema.Type != nil {
					schemaMap["type"] = schema.Type
				}

				if schema.Value != nil {
					schemaMap["value"] = schema.Value
				}

				if schema.Properties != nil {
					propertiesList := []interface{}{}
					for _, properties := range schema.Properties {
						propertiesMap := map[string]interface{}{}

						if properties.Name != nil {
							propertiesMap["name"] = properties.Name
						}

						if properties.Value != nil {
							propertiesMap["value"] = properties.Value
						}

						propertiesList = append(propertiesList, propertiesMap)
					}

					schemaMap["properties"] = []interface{}{propertiesList}
				}

				if schema.Alias != nil {
					schemaMap["alias"] = schema.Alias
				}

				if schema.Comment != nil {
					schemaMap["comment"] = schema.Comment
				}

				schemaList = append(schemaList, schemaMap)
			}

			nodeInfoMap["schema"] = []interface{}{schemaList}
		}

		if integration_task_node.NodeInfo.NodeMapping != nil {
			nodeMappingMap := map[string]interface{}{}

			if integration_task_node.NodeInfo.NodeMapping.SourceId != nil {
				nodeMappingMap["source_id"] = integration_task_node.NodeInfo.NodeMapping.SourceId
			}

			if integration_task_node.NodeInfo.NodeMapping.SinkId != nil {
				nodeMappingMap["sink_id"] = integration_task_node.NodeInfo.NodeMapping.SinkId
			}

			if integration_task_node.NodeInfo.NodeMapping.SourceSchema != nil {
				sourceSchemaList := []interface{}{}
				for _, sourceSchema := range integration_task_node.NodeInfo.NodeMapping.SourceSchema {
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

						sourceSchemaMap["properties"] = []interface{}{propertiesList}
					}

					if sourceSchema.Alias != nil {
						sourceSchemaMap["alias"] = sourceSchema.Alias
					}

					if sourceSchema.Comment != nil {
						sourceSchemaMap["comment"] = sourceSchema.Comment
					}

					sourceSchemaList = append(sourceSchemaList, sourceSchemaMap)
				}

				nodeMappingMap["source_schema"] = []interface{}{sourceSchemaList}
			}

			if integration_task_node.NodeInfo.NodeMapping.SchemaMappings != nil {
				schemaMappingsList := []interface{}{}
				for _, schemaMappings := range integration_task_node.NodeInfo.NodeMapping.SchemaMappings {
					schemaMappingsMap := map[string]interface{}{}

					if schemaMappings.SourceSchemaId != nil {
						schemaMappingsMap["source_schema_id"] = schemaMappings.SourceSchemaId
					}

					if schemaMappings.SinkSchemaId != nil {
						schemaMappingsMap["sink_schema_id"] = schemaMappings.SinkSchemaId
					}

					schemaMappingsList = append(schemaMappingsList, schemaMappingsMap)
				}

				nodeMappingMap["schema_mappings"] = []interface{}{schemaMappingsList}
			}

			if integration_task_node.NodeInfo.NodeMapping.ExtConfig != nil {
				extConfigList := []interface{}{}
				for _, extConfig := range integration_task_node.NodeInfo.NodeMapping.ExtConfig {
					extConfigMap := map[string]interface{}{}

					if extConfig.Name != nil {
						extConfigMap["name"] = extConfig.Name
					}

					if extConfig.Value != nil {
						extConfigMap["value"] = extConfig.Value
					}

					extConfigList = append(extConfigList, extConfigMap)
				}

				nodeMappingMap["ext_config"] = []interface{}{extConfigList}
			}

			nodeInfoMap["node_mapping"] = []interface{}{nodeMappingMap}
		}

		if integration_task_node.NodeInfo.AppId != nil {
			nodeInfoMap["app_id"] = integration_task_node.NodeInfo.AppId
		}

		if integration_task_node.NodeInfo.ProjectId != nil {
			nodeInfoMap["project_id"] = integration_task_node.NodeInfo.ProjectId
		}

		if integration_task_node.NodeInfo.CreatorUin != nil {
			nodeInfoMap["creator_uin"] = integration_task_node.NodeInfo.CreatorUin
		}

		if integration_task_node.NodeInfo.OperatorUin != nil {
			nodeInfoMap["operator_uin"] = integration_task_node.NodeInfo.OperatorUin
		}

		if integration_task_node.NodeInfo.OwnerUin != nil {
			nodeInfoMap["owner_uin"] = integration_task_node.NodeInfo.OwnerUin
		}

		if integration_task_node.NodeInfo.CreateTime != nil {
			nodeInfoMap["create_time"] = integration_task_node.NodeInfo.CreateTime
		}

		if integration_task_node.NodeInfo.UpdateTime != nil {
			nodeInfoMap["update_time"] = integration_task_node.NodeInfo.UpdateTime
		}

		_ = d.Set("node_info", []interface{}{nodeInfoMap})
	}

	if integration_task_node.ProjectId != nil {
		_ = d.Set("project_id", integration_task_node.ProjectId)
	}

	if integration_task_node.TaskType != nil {
		_ = d.Set("task_type", integration_task_node.TaskType)
	}

	return nil
}

func resourceTencentCloudWedataIntegration_task_nodeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_task_node.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := wedata.NewModifyIntegrationNodeRequest()

	integration_task_nodeId := d.Id()

	request.Id = &id

	immutableArgs := []string{"node_info", "project_id", "task_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("node_info") {
		if dMap, ok := helper.InterfacesHeadMap(d, "node_info"); ok {
			integrationNodeInfo := wedata.IntegrationNodeInfo{}
			if v, ok := dMap["id"]; ok {
				integrationNodeInfo.Id = helper.String(v.(string))
			}
			if v, ok := dMap["task_id"]; ok {
				integrationNodeInfo.TaskId = helper.String(v.(string))
			}
			if v, ok := dMap["name"]; ok {
				integrationNodeInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["node_type"]; ok {
				integrationNodeInfo.NodeType = helper.String(v.(string))
			}
			if v, ok := dMap["data_source_type"]; ok {
				integrationNodeInfo.DataSourceType = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				integrationNodeInfo.Description = helper.String(v.(string))
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
			if v, ok := dMap["project_id"]; ok {
				integrationNodeInfo.ProjectId = helper.String(v.(string))
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
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("task_type") {
		if v, ok := d.GetOkExists("task_type"); ok {
			request.TaskType = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().ModifyIntegrationNode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update wedata integration_task_node failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegration_task_nodeRead(d, meta)
}

func resourceTencentCloudWedataIntegration_task_nodeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_task_node.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	integration_task_nodeId := d.Id()

	if err := service.DeleteWedataIntegration_task_nodeById(ctx, id); err != nil {
		return err
	}

	return nil
}

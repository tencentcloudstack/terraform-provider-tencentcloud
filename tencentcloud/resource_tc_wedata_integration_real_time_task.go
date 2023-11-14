/*
Provides a resource to create a wedata integration_real_time_task

Example Usage

```hcl
resource "tencentcloud_wedata_integration_real_time_task" "integration_real_time_task" {
  task_info {
		task_name = "TaskTest_10"
		description = "Task for test"
		sync_type = 2
		task_type = 201
		workflow_id = "1"
		task_id = "j84cc717e-215b-4960-9575-898586bae37f"
		schedule_task_id = "1"
		task_group_id = "1"
		project_id = "1455251608631480391"
		creator_uin = "100028448000"
		operator_uin = "100028448000"
		owner_uin = "100028448000"
		app_id = "1315000000"
		status = 1
		nodes {
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
		executor_id = "2000"
		config {
			name = "Database"
			value = "db"
		}
		ext_config {
			name = "x"
			value = "320"
		}
		execute_context {
			name = ""
			value = ""
		}
		mappings {
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
		task_mode = "1"
		incharge = ""
		offline_task_add_entity {
			workflow_name = "workflow_test"
			dependency_workflow = "no"
			start_time = "2023-12-31 00:00:00"
			end_time = "2099-12-31 00:00:00"
			cycle_type = 0
			cycle_step = 1
			delay_time = 0
			crontab_expression = "0 0 1 * * ?"
			retry_wait = 5
			retriable = 1
			try_limit = 1
			run_priority = 6
			product_name = "DATA_INTEGRATION"
			self_depend = 3
			task_action = "1"
			execution_end_time = "16:59"
			execution_start_time = "02:00"
			task_auto_submit = false
			instance_init_strategy = &lt;nil&gt;
		}
		executor_group_name = "executor1"
		in_long_manager_url = "172.16.0.3:8083"
		in_long_stream_id = "b_q3b502073-1cac-4a7b-a67f-d30314833a32"
		in_long_manager_version = "v16"
		data_proxy_url =
		submit = false
		input_datasource_type = "MYSQL"
		output_datasource_type = "MYSQL"
		num_records_in = 1000
		num_records_out = 1000
		reader_delay =
		num_restarts = 1
		create_time = "2023-10-12 17:17:14"
		update_time = "2023-10-12 17:17:14"
		last_run_time = "2023-10-12 17:17:14"
		stop_time = "2023-10-12 17:17:14"
		has_version = false
		locked = false
		locker = "100028578868"
		running_cu =
		task_alarm_regular_list = &lt;nil&gt;
		switch_resource = 0
		read_phase = 0
		instance_version = 1

  }
  project_id = "1455251608631480391"
}
```

Import

wedata integration_real_time_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_integration_real_time_task.integration_real_time_task integration_real_time_task_id
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

func resourceTencentCloudWedataIntegration_real_time_task() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataIntegration_real_time_taskCreate,
		Read:   resourceTencentCloudWedataIntegration_real_time_taskRead,
		Update: resourceTencentCloudWedataIntegration_real_time_taskUpdate,
		Delete: resourceTencentCloudWedataIntegration_real_time_taskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"task_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Task Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description information.",
						},
						"sync_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Synchronization type: 1. Whole database synchronization, 2. Single table synchronization.",
						},
						"task_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Task type: 201. Real-time, 202. Offline.",
						},
						"workflow_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The workflow id to which the task belongs.",
						},
						"task_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task ID.",
						},
						"schedule_task_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task scheduling id (job id such as oceanus or us).",
						},
						"task_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Inlong Task Group ID.",
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
						"app_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User App Id.",
						},
						"status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Task status 1. Not started | Task initialization, 2. Task starting, 3. Running, 4. Paused, 5. Task stopping, 6. Stopped, 7. Execution failed, 8. deleted, 9. Locked, 404 .unknown status.",
						},
						"nodes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Task Node Information.",
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
						"executor_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Executor resource ID.",
						},
						"config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Task configuration.",
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
						"execute_context": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Execute context.",
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
						"mappings": {
							Type:        schema.TypeList,
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
						"task_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task display mode, 0: canvas mode, 1: form mode.",
						},
						"incharge": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Incharge user.",
						},
						"offline_task_add_entity": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Offline task scheduling configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workflow_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the workflow to which the task belongs.",
									},
									"dependency_workflow": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to support workflow dependencies: yes / no, default value: no.",
									},
									"start_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Effective start time, the format is yyyy-MM-dd HH:mm:ss.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Effective end time, the format is yyyy-MM-dd HH:mm:ss.",
									},
									"cycle_type": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Scheduling type, 0: crontab type, 1: minutes, 2: hours, 3: days, 4: weeks, 5: months, 6: one-time, 7: user-driven, 10: elastic period (week), 11: elastic period (month), 12: year, 13: instant trigger.",
									},
									"cycle_step": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Interval time of scheduling, the minimum value: 1.",
									},
									"delay_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Execution time, unit is minutes, only available for day/week/month/year scheduling. For example, daily scheduling is executed once every day at 02:00, and the delayTime is 120 minutes.",
									},
									"crontab_expression": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Crontab expression.",
									},
									"retry_wait": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Retry waiting time, unit is minutes.",
									},
									"retriable": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Whether to retry.",
									},
									"try_limit": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Number of retries.",
									},
									"run_priority": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Task running priority.",
									},
									"product_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Product name: DATA_INTEGRATION.",
									},
									"self_depend": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Self-dependent rules, 1: Ordered serial one at a time, queued execution, 2: Unordered serial one at a time, not queued execution, 3: Parallel, multiple at once.",
									},
									"task_action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Flexible cycle configuration, if it is a weekly task: 1 is Sunday, 2 is Monday, 3 is Tuesday, and so on. If it is a monthly task: &amp;#39;1,3&amp;#39; represents the 1st and 3rd; &amp;#39;L&amp;#39; represents the end of the month.",
									},
									"execution_end_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Scheduling execution end time.",
									},
									"execution_start_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Scheduling execution start time.",
									},
									"task_auto_submit": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to automatically submit.",
									},
									"instance_init_strategy": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Instance initialization strategy.",
									},
								},
							},
						},
						"executor_group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Executor group name.",
						},
						"in_long_manager_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "InLong manager url.",
						},
						"in_long_stream_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "InLong stream id.",
						},
						"in_long_manager_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "InLong manager version.",
						},
						"data_proxy_url": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Data proxy url.",
						},
						"submit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the task version has been submitted for operation and maintenance.",
						},
						"input_datasource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Input datasource type.",
						},
						"output_datasource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Output datasource type.",
						},
						"num_records_in": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of reads.",
						},
						"num_records_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of writes.",
						},
						"reader_delay": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "Read latency.",
						},
						"num_restarts": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Times of restarts.",
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
						"last_run_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The last time the task was run.",
						},
						"stop_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The time the task was stopped.",
						},
						"has_version": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the task been submitted.",
						},
						"locked": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the task been locked.",
						},
						"locker": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User locked task.",
						},
						"running_cu": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "The amount of resources consumed by real-time task.",
						},
						"task_alarm_regular_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Task alarm regular.",
						},
						"switch_resource": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Resource tiering status, 0: in progress, 1: successful, 2: failed.",
						},
						"read_phase": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Reading stage, 0: full amount, 1: partial full amount, 2: all incremental.",
						},
						"instance_version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Instance version.",
						},
					},
				},
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},
		},
	}
}

func resourceTencentCloudWedataIntegration_real_time_taskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_real_time_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = wedata.NewCreateIntegrationTaskRequest()
		response = wedata.NewCreateIntegrationTaskResponse()
		taskId   string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "task_info"); ok {
		integrationTaskInfo := wedata.IntegrationTaskInfo{}
		if v, ok := dMap["task_name"]; ok {
			integrationTaskInfo.TaskName = helper.String(v.(string))
		}
		if v, ok := dMap["description"]; ok {
			integrationTaskInfo.Description = helper.String(v.(string))
		}
		if v, ok := dMap["sync_type"]; ok {
			integrationTaskInfo.SyncType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["task_type"]; ok {
			integrationTaskInfo.TaskType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["workflow_id"]; ok {
			integrationTaskInfo.WorkflowId = helper.String(v.(string))
		}
		if v, ok := dMap["task_id"]; ok {
			integrationTaskInfo.TaskId = helper.String(v.(string))
		}
		if v, ok := dMap["schedule_task_id"]; ok {
			integrationTaskInfo.ScheduleTaskId = helper.String(v.(string))
		}
		if v, ok := dMap["task_group_id"]; ok {
			integrationTaskInfo.TaskGroupId = helper.String(v.(string))
		}
		if v, ok := dMap["project_id"]; ok {
			integrationTaskInfo.ProjectId = helper.String(v.(string))
		}
		if v, ok := dMap["creator_uin"]; ok {
			integrationTaskInfo.CreatorUin = helper.String(v.(string))
		}
		if v, ok := dMap["operator_uin"]; ok {
			integrationTaskInfo.OperatorUin = helper.String(v.(string))
		}
		if v, ok := dMap["owner_uin"]; ok {
			integrationTaskInfo.OwnerUin = helper.String(v.(string))
		}
		if v, ok := dMap["app_id"]; ok {
			integrationTaskInfo.AppId = helper.String(v.(string))
		}
		if v, ok := dMap["status"]; ok {
			integrationTaskInfo.Status = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["nodes"]; ok {
			for _, item := range v.([]interface{}) {
				nodesMap := item.(map[string]interface{})
				integrationNodeInfo := wedata.IntegrationNodeInfo{}
				if v, ok := nodesMap["id"]; ok {
					integrationNodeInfo.Id = helper.String(v.(string))
				}
				if v, ok := nodesMap["task_id"]; ok {
					integrationNodeInfo.TaskId = helper.String(v.(string))
				}
				if v, ok := nodesMap["name"]; ok {
					integrationNodeInfo.Name = helper.String(v.(string))
				}
				if v, ok := nodesMap["node_type"]; ok {
					integrationNodeInfo.NodeType = helper.String(v.(string))
				}
				if v, ok := nodesMap["data_source_type"]; ok {
					integrationNodeInfo.DataSourceType = helper.String(v.(string))
				}
				if v, ok := nodesMap["description"]; ok {
					integrationNodeInfo.Description = helper.String(v.(string))
				}
				if v, ok := nodesMap["datasource_id"]; ok {
					integrationNodeInfo.DatasourceId = helper.String(v.(string))
				}
				if v, ok := nodesMap["config"]; ok {
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
				if v, ok := nodesMap["ext_config"]; ok {
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
				if v, ok := nodesMap["schema"]; ok {
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
				if nodeMappingMap, ok := helper.InterfaceToMap(nodesMap, "node_mapping"); ok {
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
				if v, ok := nodesMap["app_id"]; ok {
					integrationNodeInfo.AppId = helper.String(v.(string))
				}
				if v, ok := nodesMap["project_id"]; ok {
					integrationNodeInfo.ProjectId = helper.String(v.(string))
				}
				if v, ok := nodesMap["creator_uin"]; ok {
					integrationNodeInfo.CreatorUin = helper.String(v.(string))
				}
				if v, ok := nodesMap["operator_uin"]; ok {
					integrationNodeInfo.OperatorUin = helper.String(v.(string))
				}
				if v, ok := nodesMap["owner_uin"]; ok {
					integrationNodeInfo.OwnerUin = helper.String(v.(string))
				}
				if v, ok := nodesMap["create_time"]; ok {
					integrationNodeInfo.CreateTime = helper.String(v.(string))
				}
				if v, ok := nodesMap["update_time"]; ok {
					integrationNodeInfo.UpdateTime = helper.String(v.(string))
				}
				integrationTaskInfo.Nodes = append(integrationTaskInfo.Nodes, &integrationNodeInfo)
			}
		}
		if v, ok := dMap["executor_id"]; ok {
			integrationTaskInfo.ExecutorId = helper.String(v.(string))
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
				integrationTaskInfo.Config = append(integrationTaskInfo.Config, &recordField)
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
				integrationTaskInfo.ExtConfig = append(integrationTaskInfo.ExtConfig, &recordField)
			}
		}
		if v, ok := dMap["execute_context"]; ok {
			for _, item := range v.([]interface{}) {
				executeContextMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := executeContextMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}
				if v, ok := executeContextMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}
				integrationTaskInfo.ExecuteContext = append(integrationTaskInfo.ExecuteContext, &recordField)
			}
		}
		if v, ok := dMap["mappings"]; ok {
			for _, item := range v.([]interface{}) {
				mappingsMap := item.(map[string]interface{})
				integrationNodeMapping := wedata.IntegrationNodeMapping{}
				if v, ok := mappingsMap["source_id"]; ok {
					integrationNodeMapping.SourceId = helper.String(v.(string))
				}
				if v, ok := mappingsMap["sink_id"]; ok {
					integrationNodeMapping.SinkId = helper.String(v.(string))
				}
				if v, ok := mappingsMap["source_schema"]; ok {
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
				if v, ok := mappingsMap["schema_mappings"]; ok {
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
				if v, ok := mappingsMap["ext_config"]; ok {
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
				integrationTaskInfo.Mappings = append(integrationTaskInfo.Mappings, &integrationNodeMapping)
			}
		}
		if v, ok := dMap["task_mode"]; ok {
			integrationTaskInfo.TaskMode = helper.String(v.(string))
		}
		if v, ok := dMap["incharge"]; ok {
			integrationTaskInfo.Incharge = helper.String(v.(string))
		}
		if offlineTaskAddEntityMap, ok := helper.InterfaceToMap(dMap, "offline_task_add_entity"); ok {
			offlineTaskAddParam := wedata.OfflineTaskAddParam{}
			if v, ok := offlineTaskAddEntityMap["workflow_name"]; ok {
				offlineTaskAddParam.WorkflowName = helper.String(v.(string))
			}
			if v, ok := offlineTaskAddEntityMap["dependency_workflow"]; ok {
				offlineTaskAddParam.DependencyWorkflow = helper.String(v.(string))
			}
			if v, ok := offlineTaskAddEntityMap["start_time"]; ok {
				offlineTaskAddParam.StartTime = helper.String(v.(string))
			}
			if v, ok := offlineTaskAddEntityMap["end_time"]; ok {
				offlineTaskAddParam.EndTime = helper.String(v.(string))
			}
			if v, ok := offlineTaskAddEntityMap["cycle_type"]; ok {
				offlineTaskAddParam.CycleType = helper.IntUint64(v.(int))
			}
			if v, ok := offlineTaskAddEntityMap["cycle_step"]; ok {
				offlineTaskAddParam.CycleStep = helper.IntUint64(v.(int))
			}
			if v, ok := offlineTaskAddEntityMap["delay_time"]; ok {
				offlineTaskAddParam.DelayTime = helper.IntUint64(v.(int))
			}
			if v, ok := offlineTaskAddEntityMap["crontab_expression"]; ok {
				offlineTaskAddParam.CrontabExpression = helper.String(v.(string))
			}
			if v, ok := offlineTaskAddEntityMap["retry_wait"]; ok {
				offlineTaskAddParam.RetryWait = helper.IntUint64(v.(int))
			}
			if v, ok := offlineTaskAddEntityMap["retriable"]; ok {
				offlineTaskAddParam.Retriable = helper.IntUint64(v.(int))
			}
			if v, ok := offlineTaskAddEntityMap["try_limit"]; ok {
				offlineTaskAddParam.TryLimit = helper.IntUint64(v.(int))
			}
			if v, ok := offlineTaskAddEntityMap["run_priority"]; ok {
				offlineTaskAddParam.RunPriority = helper.IntUint64(v.(int))
			}
			if v, ok := offlineTaskAddEntityMap["product_name"]; ok {
				offlineTaskAddParam.ProductName = helper.String(v.(string))
			}
			if v, ok := offlineTaskAddEntityMap["self_depend"]; ok {
				offlineTaskAddParam.SelfDepend = helper.IntUint64(v.(int))
			}
			if v, ok := offlineTaskAddEntityMap["task_action"]; ok {
				offlineTaskAddParam.TaskAction = helper.String(v.(string))
			}
			if v, ok := offlineTaskAddEntityMap["execution_end_time"]; ok {
				offlineTaskAddParam.ExecutionEndTime = helper.String(v.(string))
			}
			if v, ok := offlineTaskAddEntityMap["execution_start_time"]; ok {
				offlineTaskAddParam.ExecutionStartTime = helper.String(v.(string))
			}
			if v, ok := offlineTaskAddEntityMap["task_auto_submit"]; ok {
				offlineTaskAddParam.TaskAutoSubmit = helper.Bool(v.(bool))
			}
			if v, ok := offlineTaskAddEntityMap["instance_init_strategy"]; ok {
				offlineTaskAddParam.InstanceInitStrategy = helper.String(v.(string))
			}
			integrationTaskInfo.OfflineTaskAddEntity = &offlineTaskAddParam
		}
		if v, ok := dMap["executor_group_name"]; ok {
			integrationTaskInfo.ExecutorGroupName = helper.String(v.(string))
		}
		if v, ok := dMap["in_long_manager_url"]; ok {
			integrationTaskInfo.InLongManagerUrl = helper.String(v.(string))
		}
		if v, ok := dMap["in_long_stream_id"]; ok {
			integrationTaskInfo.InLongStreamId = helper.String(v.(string))
		}
		if v, ok := dMap["in_long_manager_version"]; ok {
			integrationTaskInfo.InLongManagerVersion = helper.String(v.(string))
		}
		if v, ok := dMap["data_proxy_url"]; ok {
			dataProxyUrlSet := v.(*schema.Set).List()
			for i := range dataProxyUrlSet {
				dataProxyUrl := dataProxyUrlSet[i].(string)
				integrationTaskInfo.DataProxyUrl = append(integrationTaskInfo.DataProxyUrl, &dataProxyUrl)
			}
		}
		if v, ok := dMap["submit"]; ok {
			integrationTaskInfo.Submit = helper.Bool(v.(bool))
		}
		if v, ok := dMap["input_datasource_type"]; ok {
			integrationTaskInfo.InputDatasourceType = helper.String(v.(string))
		}
		if v, ok := dMap["output_datasource_type"]; ok {
			integrationTaskInfo.OutputDatasourceType = helper.String(v.(string))
		}
		if v, ok := dMap["num_records_in"]; ok {
			integrationTaskInfo.NumRecordsIn = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["num_records_out"]; ok {
			integrationTaskInfo.NumRecordsOut = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["reader_delay"]; ok {
			integrationTaskInfo.ReaderDelay = helper.Float64(v.(float64))
		}
		if v, ok := dMap["num_restarts"]; ok {
			integrationTaskInfo.NumRestarts = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["create_time"]; ok {
			integrationTaskInfo.CreateTime = helper.String(v.(string))
		}
		if v, ok := dMap["update_time"]; ok {
			integrationTaskInfo.UpdateTime = helper.String(v.(string))
		}
		if v, ok := dMap["last_run_time"]; ok {
			integrationTaskInfo.LastRunTime = helper.String(v.(string))
		}
		if v, ok := dMap["stop_time"]; ok {
			integrationTaskInfo.StopTime = helper.String(v.(string))
		}
		if v, ok := dMap["has_version"]; ok {
			integrationTaskInfo.HasVersion = helper.Bool(v.(bool))
		}
		if v, ok := dMap["locked"]; ok {
			integrationTaskInfo.Locked = helper.Bool(v.(bool))
		}
		if v, ok := dMap["locker"]; ok {
			integrationTaskInfo.Locker = helper.String(v.(string))
		}
		if v, ok := dMap["running_cu"]; ok {
			integrationTaskInfo.RunningCu = helper.Float64(v.(float64))
		}
		if v, ok := dMap["task_alarm_regular_list"]; ok {
			taskAlarmRegularListSet := v.(*schema.Set).List()
			for i := range taskAlarmRegularListSet {
				taskAlarmRegularList := taskAlarmRegularListSet[i].(string)
				integrationTaskInfo.TaskAlarmRegularList = append(integrationTaskInfo.TaskAlarmRegularList, &taskAlarmRegularList)
			}
		}
		if v, ok := dMap["switch_resource"]; ok {
			integrationTaskInfo.SwitchResource = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["read_phase"]; ok {
			integrationTaskInfo.ReadPhase = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["instance_version"]; ok {
			integrationTaskInfo.InstanceVersion = helper.IntInt64(v.(int))
		}
		request.TaskInfo = &integrationTaskInfo
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateIntegrationTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata integration_real_time_task failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudWedataIntegration_real_time_taskRead(d, meta)
}

func resourceTencentCloudWedataIntegration_real_time_taskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_real_time_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	integration_real_time_taskId := d.Id()

	integration_real_time_task, err := service.DescribeWedataIntegration_real_time_taskById(ctx, taskId)
	if err != nil {
		return err
	}

	if integration_real_time_task == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataIntegration_real_time_task` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if integration_real_time_task.TaskInfo != nil {
		taskInfoMap := map[string]interface{}{}

		if integration_real_time_task.TaskInfo.TaskName != nil {
			taskInfoMap["task_name"] = integration_real_time_task.TaskInfo.TaskName
		}

		if integration_real_time_task.TaskInfo.Description != nil {
			taskInfoMap["description"] = integration_real_time_task.TaskInfo.Description
		}

		if integration_real_time_task.TaskInfo.SyncType != nil {
			taskInfoMap["sync_type"] = integration_real_time_task.TaskInfo.SyncType
		}

		if integration_real_time_task.TaskInfo.TaskType != nil {
			taskInfoMap["task_type"] = integration_real_time_task.TaskInfo.TaskType
		}

		if integration_real_time_task.TaskInfo.WorkflowId != nil {
			taskInfoMap["workflow_id"] = integration_real_time_task.TaskInfo.WorkflowId
		}

		if integration_real_time_task.TaskInfo.TaskId != nil {
			taskInfoMap["task_id"] = integration_real_time_task.TaskInfo.TaskId
		}

		if integration_real_time_task.TaskInfo.ScheduleTaskId != nil {
			taskInfoMap["schedule_task_id"] = integration_real_time_task.TaskInfo.ScheduleTaskId
		}

		if integration_real_time_task.TaskInfo.TaskGroupId != nil {
			taskInfoMap["task_group_id"] = integration_real_time_task.TaskInfo.TaskGroupId
		}

		if integration_real_time_task.TaskInfo.ProjectId != nil {
			taskInfoMap["project_id"] = integration_real_time_task.TaskInfo.ProjectId
		}

		if integration_real_time_task.TaskInfo.CreatorUin != nil {
			taskInfoMap["creator_uin"] = integration_real_time_task.TaskInfo.CreatorUin
		}

		if integration_real_time_task.TaskInfo.OperatorUin != nil {
			taskInfoMap["operator_uin"] = integration_real_time_task.TaskInfo.OperatorUin
		}

		if integration_real_time_task.TaskInfo.OwnerUin != nil {
			taskInfoMap["owner_uin"] = integration_real_time_task.TaskInfo.OwnerUin
		}

		if integration_real_time_task.TaskInfo.AppId != nil {
			taskInfoMap["app_id"] = integration_real_time_task.TaskInfo.AppId
		}

		if integration_real_time_task.TaskInfo.Status != nil {
			taskInfoMap["status"] = integration_real_time_task.TaskInfo.Status
		}

		if integration_real_time_task.TaskInfo.Nodes != nil {
			nodesList := []interface{}{}
			for _, nodes := range integration_real_time_task.TaskInfo.Nodes {
				nodesMap := map[string]interface{}{}

				if nodes.Id != nil {
					nodesMap["id"] = nodes.Id
				}

				if nodes.TaskId != nil {
					nodesMap["task_id"] = nodes.TaskId
				}

				if nodes.Name != nil {
					nodesMap["name"] = nodes.Name
				}

				if nodes.NodeType != nil {
					nodesMap["node_type"] = nodes.NodeType
				}

				if nodes.DataSourceType != nil {
					nodesMap["data_source_type"] = nodes.DataSourceType
				}

				if nodes.Description != nil {
					nodesMap["description"] = nodes.Description
				}

				if nodes.DatasourceId != nil {
					nodesMap["datasource_id"] = nodes.DatasourceId
				}

				if nodes.Config != nil {
					configList := []interface{}{}
					for _, config := range nodes.Config {
						configMap := map[string]interface{}{}

						if config.Name != nil {
							configMap["name"] = config.Name
						}

						if config.Value != nil {
							configMap["value"] = config.Value
						}

						configList = append(configList, configMap)
					}

					nodesMap["config"] = []interface{}{configList}
				}

				if nodes.ExtConfig != nil {
					extConfigList := []interface{}{}
					for _, extConfig := range nodes.ExtConfig {
						extConfigMap := map[string]interface{}{}

						if extConfig.Name != nil {
							extConfigMap["name"] = extConfig.Name
						}

						if extConfig.Value != nil {
							extConfigMap["value"] = extConfig.Value
						}

						extConfigList = append(extConfigList, extConfigMap)
					}

					nodesMap["ext_config"] = []interface{}{extConfigList}
				}

				if nodes.Schema != nil {
					schemaList := []interface{}{}
					for _, schema := range nodes.Schema {
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

					nodesMap["schema"] = []interface{}{schemaList}
				}

				if nodes.NodeMapping != nil {
					nodeMappingMap := map[string]interface{}{}

					if nodes.NodeMapping.SourceId != nil {
						nodeMappingMap["source_id"] = nodes.NodeMapping.SourceId
					}

					if nodes.NodeMapping.SinkId != nil {
						nodeMappingMap["sink_id"] = nodes.NodeMapping.SinkId
					}

					if nodes.NodeMapping.SourceSchema != nil {
						sourceSchemaList := []interface{}{}
						for _, sourceSchema := range nodes.NodeMapping.SourceSchema {
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

					if nodes.NodeMapping.SchemaMappings != nil {
						schemaMappingsList := []interface{}{}
						for _, schemaMappings := range nodes.NodeMapping.SchemaMappings {
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

					if nodes.NodeMapping.ExtConfig != nil {
						extConfigList := []interface{}{}
						for _, extConfig := range nodes.NodeMapping.ExtConfig {
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

					nodesMap["node_mapping"] = []interface{}{nodeMappingMap}
				}

				if nodes.AppId != nil {
					nodesMap["app_id"] = nodes.AppId
				}

				if nodes.ProjectId != nil {
					nodesMap["project_id"] = nodes.ProjectId
				}

				if nodes.CreatorUin != nil {
					nodesMap["creator_uin"] = nodes.CreatorUin
				}

				if nodes.OperatorUin != nil {
					nodesMap["operator_uin"] = nodes.OperatorUin
				}

				if nodes.OwnerUin != nil {
					nodesMap["owner_uin"] = nodes.OwnerUin
				}

				if nodes.CreateTime != nil {
					nodesMap["create_time"] = nodes.CreateTime
				}

				if nodes.UpdateTime != nil {
					nodesMap["update_time"] = nodes.UpdateTime
				}

				nodesList = append(nodesList, nodesMap)
			}

			taskInfoMap["nodes"] = []interface{}{nodesList}
		}

		if integration_real_time_task.TaskInfo.ExecutorId != nil {
			taskInfoMap["executor_id"] = integration_real_time_task.TaskInfo.ExecutorId
		}

		if integration_real_time_task.TaskInfo.Config != nil {
			configList := []interface{}{}
			for _, config := range integration_real_time_task.TaskInfo.Config {
				configMap := map[string]interface{}{}

				if config.Name != nil {
					configMap["name"] = config.Name
				}

				if config.Value != nil {
					configMap["value"] = config.Value
				}

				configList = append(configList, configMap)
			}

			taskInfoMap["config"] = []interface{}{configList}
		}

		if integration_real_time_task.TaskInfo.ExtConfig != nil {
			extConfigList := []interface{}{}
			for _, extConfig := range integration_real_time_task.TaskInfo.ExtConfig {
				extConfigMap := map[string]interface{}{}

				if extConfig.Name != nil {
					extConfigMap["name"] = extConfig.Name
				}

				if extConfig.Value != nil {
					extConfigMap["value"] = extConfig.Value
				}

				extConfigList = append(extConfigList, extConfigMap)
			}

			taskInfoMap["ext_config"] = []interface{}{extConfigList}
		}

		if integration_real_time_task.TaskInfo.ExecuteContext != nil {
			executeContextList := []interface{}{}
			for _, executeContext := range integration_real_time_task.TaskInfo.ExecuteContext {
				executeContextMap := map[string]interface{}{}

				if executeContext.Name != nil {
					executeContextMap["name"] = executeContext.Name
				}

				if executeContext.Value != nil {
					executeContextMap["value"] = executeContext.Value
				}

				executeContextList = append(executeContextList, executeContextMap)
			}

			taskInfoMap["execute_context"] = []interface{}{executeContextList}
		}

		if integration_real_time_task.TaskInfo.Mappings != nil {
			mappingsList := []interface{}{}
			for _, mappings := range integration_real_time_task.TaskInfo.Mappings {
				mappingsMap := map[string]interface{}{}

				if mappings.SourceId != nil {
					mappingsMap["source_id"] = mappings.SourceId
				}

				if mappings.SinkId != nil {
					mappingsMap["sink_id"] = mappings.SinkId
				}

				if mappings.SourceSchema != nil {
					sourceSchemaList := []interface{}{}
					for _, sourceSchema := range mappings.SourceSchema {
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

					mappingsMap["source_schema"] = []interface{}{sourceSchemaList}
				}

				if mappings.SchemaMappings != nil {
					schemaMappingsList := []interface{}{}
					for _, schemaMappings := range mappings.SchemaMappings {
						schemaMappingsMap := map[string]interface{}{}

						if schemaMappings.SourceSchemaId != nil {
							schemaMappingsMap["source_schema_id"] = schemaMappings.SourceSchemaId
						}

						if schemaMappings.SinkSchemaId != nil {
							schemaMappingsMap["sink_schema_id"] = schemaMappings.SinkSchemaId
						}

						schemaMappingsList = append(schemaMappingsList, schemaMappingsMap)
					}

					mappingsMap["schema_mappings"] = []interface{}{schemaMappingsList}
				}

				if mappings.ExtConfig != nil {
					extConfigList := []interface{}{}
					for _, extConfig := range mappings.ExtConfig {
						extConfigMap := map[string]interface{}{}

						if extConfig.Name != nil {
							extConfigMap["name"] = extConfig.Name
						}

						if extConfig.Value != nil {
							extConfigMap["value"] = extConfig.Value
						}

						extConfigList = append(extConfigList, extConfigMap)
					}

					mappingsMap["ext_config"] = []interface{}{extConfigList}
				}

				mappingsList = append(mappingsList, mappingsMap)
			}

			taskInfoMap["mappings"] = []interface{}{mappingsList}
		}

		if integration_real_time_task.TaskInfo.TaskMode != nil {
			taskInfoMap["task_mode"] = integration_real_time_task.TaskInfo.TaskMode
		}

		if integration_real_time_task.TaskInfo.Incharge != nil {
			taskInfoMap["incharge"] = integration_real_time_task.TaskInfo.Incharge
		}

		if integration_real_time_task.TaskInfo.OfflineTaskAddEntity != nil {
			offlineTaskAddEntityMap := map[string]interface{}{}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.WorkflowName != nil {
				offlineTaskAddEntityMap["workflow_name"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.WorkflowName
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.DependencyWorkflow != nil {
				offlineTaskAddEntityMap["dependency_workflow"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.DependencyWorkflow
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.StartTime != nil {
				offlineTaskAddEntityMap["start_time"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.StartTime
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.EndTime != nil {
				offlineTaskAddEntityMap["end_time"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.EndTime
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.CycleType != nil {
				offlineTaskAddEntityMap["cycle_type"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.CycleType
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.CycleStep != nil {
				offlineTaskAddEntityMap["cycle_step"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.CycleStep
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.DelayTime != nil {
				offlineTaskAddEntityMap["delay_time"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.DelayTime
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.CrontabExpression != nil {
				offlineTaskAddEntityMap["crontab_expression"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.CrontabExpression
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.RetryWait != nil {
				offlineTaskAddEntityMap["retry_wait"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.RetryWait
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.Retriable != nil {
				offlineTaskAddEntityMap["retriable"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.Retriable
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.TryLimit != nil {
				offlineTaskAddEntityMap["try_limit"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.TryLimit
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.RunPriority != nil {
				offlineTaskAddEntityMap["run_priority"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.RunPriority
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.ProductName != nil {
				offlineTaskAddEntityMap["product_name"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.ProductName
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.SelfDepend != nil {
				offlineTaskAddEntityMap["self_depend"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.SelfDepend
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.TaskAction != nil {
				offlineTaskAddEntityMap["task_action"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.TaskAction
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.ExecutionEndTime != nil {
				offlineTaskAddEntityMap["execution_end_time"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.ExecutionEndTime
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.ExecutionStartTime != nil {
				offlineTaskAddEntityMap["execution_start_time"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.ExecutionStartTime
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.TaskAutoSubmit != nil {
				offlineTaskAddEntityMap["task_auto_submit"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.TaskAutoSubmit
			}

			if integration_real_time_task.TaskInfo.OfflineTaskAddEntity.InstanceInitStrategy != nil {
				offlineTaskAddEntityMap["instance_init_strategy"] = integration_real_time_task.TaskInfo.OfflineTaskAddEntity.InstanceInitStrategy
			}

			taskInfoMap["offline_task_add_entity"] = []interface{}{offlineTaskAddEntityMap}
		}

		if integration_real_time_task.TaskInfo.ExecutorGroupName != nil {
			taskInfoMap["executor_group_name"] = integration_real_time_task.TaskInfo.ExecutorGroupName
		}

		if integration_real_time_task.TaskInfo.InLongManagerUrl != nil {
			taskInfoMap["in_long_manager_url"] = integration_real_time_task.TaskInfo.InLongManagerUrl
		}

		if integration_real_time_task.TaskInfo.InLongStreamId != nil {
			taskInfoMap["in_long_stream_id"] = integration_real_time_task.TaskInfo.InLongStreamId
		}

		if integration_real_time_task.TaskInfo.InLongManagerVersion != nil {
			taskInfoMap["in_long_manager_version"] = integration_real_time_task.TaskInfo.InLongManagerVersion
		}

		if integration_real_time_task.TaskInfo.DataProxyUrl != nil {
			taskInfoMap["data_proxy_url"] = integration_real_time_task.TaskInfo.DataProxyUrl
		}

		if integration_real_time_task.TaskInfo.Submit != nil {
			taskInfoMap["submit"] = integration_real_time_task.TaskInfo.Submit
		}

		if integration_real_time_task.TaskInfo.InputDatasourceType != nil {
			taskInfoMap["input_datasource_type"] = integration_real_time_task.TaskInfo.InputDatasourceType
		}

		if integration_real_time_task.TaskInfo.OutputDatasourceType != nil {
			taskInfoMap["output_datasource_type"] = integration_real_time_task.TaskInfo.OutputDatasourceType
		}

		if integration_real_time_task.TaskInfo.NumRecordsIn != nil {
			taskInfoMap["num_records_in"] = integration_real_time_task.TaskInfo.NumRecordsIn
		}

		if integration_real_time_task.TaskInfo.NumRecordsOut != nil {
			taskInfoMap["num_records_out"] = integration_real_time_task.TaskInfo.NumRecordsOut
		}

		if integration_real_time_task.TaskInfo.ReaderDelay != nil {
			taskInfoMap["reader_delay"] = integration_real_time_task.TaskInfo.ReaderDelay
		}

		if integration_real_time_task.TaskInfo.NumRestarts != nil {
			taskInfoMap["num_restarts"] = integration_real_time_task.TaskInfo.NumRestarts
		}

		if integration_real_time_task.TaskInfo.CreateTime != nil {
			taskInfoMap["create_time"] = integration_real_time_task.TaskInfo.CreateTime
		}

		if integration_real_time_task.TaskInfo.UpdateTime != nil {
			taskInfoMap["update_time"] = integration_real_time_task.TaskInfo.UpdateTime
		}

		if integration_real_time_task.TaskInfo.LastRunTime != nil {
			taskInfoMap["last_run_time"] = integration_real_time_task.TaskInfo.LastRunTime
		}

		if integration_real_time_task.TaskInfo.StopTime != nil {
			taskInfoMap["stop_time"] = integration_real_time_task.TaskInfo.StopTime
		}

		if integration_real_time_task.TaskInfo.HasVersion != nil {
			taskInfoMap["has_version"] = integration_real_time_task.TaskInfo.HasVersion
		}

		if integration_real_time_task.TaskInfo.Locked != nil {
			taskInfoMap["locked"] = integration_real_time_task.TaskInfo.Locked
		}

		if integration_real_time_task.TaskInfo.Locker != nil {
			taskInfoMap["locker"] = integration_real_time_task.TaskInfo.Locker
		}

		if integration_real_time_task.TaskInfo.RunningCu != nil {
			taskInfoMap["running_cu"] = integration_real_time_task.TaskInfo.RunningCu
		}

		if integration_real_time_task.TaskInfo.TaskAlarmRegularList != nil {
			taskInfoMap["task_alarm_regular_list"] = integration_real_time_task.TaskInfo.TaskAlarmRegularList
		}

		if integration_real_time_task.TaskInfo.SwitchResource != nil {
			taskInfoMap["switch_resource"] = integration_real_time_task.TaskInfo.SwitchResource
		}

		if integration_real_time_task.TaskInfo.ReadPhase != nil {
			taskInfoMap["read_phase"] = integration_real_time_task.TaskInfo.ReadPhase
		}

		if integration_real_time_task.TaskInfo.InstanceVersion != nil {
			taskInfoMap["instance_version"] = integration_real_time_task.TaskInfo.InstanceVersion
		}

		_ = d.Set("task_info", []interface{}{taskInfoMap})
	}

	if integration_real_time_task.ProjectId != nil {
		_ = d.Set("project_id", integration_real_time_task.ProjectId)
	}

	return nil
}

func resourceTencentCloudWedataIntegration_real_time_taskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_real_time_task.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := wedata.NewModifyIntegrationTaskRequest()

	integration_real_time_taskId := d.Id()

	request.TaskId = &taskId

	immutableArgs := []string{"task_info", "project_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("task_info") {
		if dMap, ok := helper.InterfacesHeadMap(d, "task_info"); ok {
			integrationTaskInfo := wedata.IntegrationTaskInfo{}
			if v, ok := dMap["task_name"]; ok {
				integrationTaskInfo.TaskName = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				integrationTaskInfo.Description = helper.String(v.(string))
			}
			if v, ok := dMap["sync_type"]; ok {
				integrationTaskInfo.SyncType = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["task_type"]; ok {
				integrationTaskInfo.TaskType = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["workflow_id"]; ok {
				integrationTaskInfo.WorkflowId = helper.String(v.(string))
			}
			if v, ok := dMap["task_id"]; ok {
				integrationTaskInfo.TaskId = helper.String(v.(string))
			}
			if v, ok := dMap["schedule_task_id"]; ok {
				integrationTaskInfo.ScheduleTaskId = helper.String(v.(string))
			}
			if v, ok := dMap["task_group_id"]; ok {
				integrationTaskInfo.TaskGroupId = helper.String(v.(string))
			}
			if v, ok := dMap["project_id"]; ok {
				integrationTaskInfo.ProjectId = helper.String(v.(string))
			}
			if v, ok := dMap["creator_uin"]; ok {
				integrationTaskInfo.CreatorUin = helper.String(v.(string))
			}
			if v, ok := dMap["operator_uin"]; ok {
				integrationTaskInfo.OperatorUin = helper.String(v.(string))
			}
			if v, ok := dMap["owner_uin"]; ok {
				integrationTaskInfo.OwnerUin = helper.String(v.(string))
			}
			if v, ok := dMap["app_id"]; ok {
				integrationTaskInfo.AppId = helper.String(v.(string))
			}
			if v, ok := dMap["status"]; ok {
				integrationTaskInfo.Status = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["nodes"]; ok {
				for _, item := range v.([]interface{}) {
					nodesMap := item.(map[string]interface{})
					integrationNodeInfo := wedata.IntegrationNodeInfo{}
					if v, ok := nodesMap["id"]; ok {
						integrationNodeInfo.Id = helper.String(v.(string))
					}
					if v, ok := nodesMap["task_id"]; ok {
						integrationNodeInfo.TaskId = helper.String(v.(string))
					}
					if v, ok := nodesMap["name"]; ok {
						integrationNodeInfo.Name = helper.String(v.(string))
					}
					if v, ok := nodesMap["node_type"]; ok {
						integrationNodeInfo.NodeType = helper.String(v.(string))
					}
					if v, ok := nodesMap["data_source_type"]; ok {
						integrationNodeInfo.DataSourceType = helper.String(v.(string))
					}
					if v, ok := nodesMap["description"]; ok {
						integrationNodeInfo.Description = helper.String(v.(string))
					}
					if v, ok := nodesMap["datasource_id"]; ok {
						integrationNodeInfo.DatasourceId = helper.String(v.(string))
					}
					if v, ok := nodesMap["config"]; ok {
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
					if v, ok := nodesMap["ext_config"]; ok {
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
					if v, ok := nodesMap["schema"]; ok {
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
					if nodeMappingMap, ok := helper.InterfaceToMap(nodesMap, "node_mapping"); ok {
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
					if v, ok := nodesMap["app_id"]; ok {
						integrationNodeInfo.AppId = helper.String(v.(string))
					}
					if v, ok := nodesMap["project_id"]; ok {
						integrationNodeInfo.ProjectId = helper.String(v.(string))
					}
					if v, ok := nodesMap["creator_uin"]; ok {
						integrationNodeInfo.CreatorUin = helper.String(v.(string))
					}
					if v, ok := nodesMap["operator_uin"]; ok {
						integrationNodeInfo.OperatorUin = helper.String(v.(string))
					}
					if v, ok := nodesMap["owner_uin"]; ok {
						integrationNodeInfo.OwnerUin = helper.String(v.(string))
					}
					if v, ok := nodesMap["create_time"]; ok {
						integrationNodeInfo.CreateTime = helper.String(v.(string))
					}
					if v, ok := nodesMap["update_time"]; ok {
						integrationNodeInfo.UpdateTime = helper.String(v.(string))
					}
					integrationTaskInfo.Nodes = append(integrationTaskInfo.Nodes, &integrationNodeInfo)
				}
			}
			if v, ok := dMap["executor_id"]; ok {
				integrationTaskInfo.ExecutorId = helper.String(v.(string))
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
					integrationTaskInfo.Config = append(integrationTaskInfo.Config, &recordField)
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
					integrationTaskInfo.ExtConfig = append(integrationTaskInfo.ExtConfig, &recordField)
				}
			}
			if v, ok := dMap["execute_context"]; ok {
				for _, item := range v.([]interface{}) {
					executeContextMap := item.(map[string]interface{})
					recordField := wedata.RecordField{}
					if v, ok := executeContextMap["name"]; ok {
						recordField.Name = helper.String(v.(string))
					}
					if v, ok := executeContextMap["value"]; ok {
						recordField.Value = helper.String(v.(string))
					}
					integrationTaskInfo.ExecuteContext = append(integrationTaskInfo.ExecuteContext, &recordField)
				}
			}
			if v, ok := dMap["mappings"]; ok {
				for _, item := range v.([]interface{}) {
					mappingsMap := item.(map[string]interface{})
					integrationNodeMapping := wedata.IntegrationNodeMapping{}
					if v, ok := mappingsMap["source_id"]; ok {
						integrationNodeMapping.SourceId = helper.String(v.(string))
					}
					if v, ok := mappingsMap["sink_id"]; ok {
						integrationNodeMapping.SinkId = helper.String(v.(string))
					}
					if v, ok := mappingsMap["source_schema"]; ok {
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
					if v, ok := mappingsMap["schema_mappings"]; ok {
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
					if v, ok := mappingsMap["ext_config"]; ok {
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
					integrationTaskInfo.Mappings = append(integrationTaskInfo.Mappings, &integrationNodeMapping)
				}
			}
			if v, ok := dMap["task_mode"]; ok {
				integrationTaskInfo.TaskMode = helper.String(v.(string))
			}
			if v, ok := dMap["incharge"]; ok {
				integrationTaskInfo.Incharge = helper.String(v.(string))
			}
			if offlineTaskAddEntityMap, ok := helper.InterfaceToMap(dMap, "offline_task_add_entity"); ok {
				offlineTaskAddParam := wedata.OfflineTaskAddParam{}
				if v, ok := offlineTaskAddEntityMap["workflow_name"]; ok {
					offlineTaskAddParam.WorkflowName = helper.String(v.(string))
				}
				if v, ok := offlineTaskAddEntityMap["dependency_workflow"]; ok {
					offlineTaskAddParam.DependencyWorkflow = helper.String(v.(string))
				}
				if v, ok := offlineTaskAddEntityMap["start_time"]; ok {
					offlineTaskAddParam.StartTime = helper.String(v.(string))
				}
				if v, ok := offlineTaskAddEntityMap["end_time"]; ok {
					offlineTaskAddParam.EndTime = helper.String(v.(string))
				}
				if v, ok := offlineTaskAddEntityMap["cycle_type"]; ok {
					offlineTaskAddParam.CycleType = helper.IntUint64(v.(int))
				}
				if v, ok := offlineTaskAddEntityMap["cycle_step"]; ok {
					offlineTaskAddParam.CycleStep = helper.IntUint64(v.(int))
				}
				if v, ok := offlineTaskAddEntityMap["delay_time"]; ok {
					offlineTaskAddParam.DelayTime = helper.IntUint64(v.(int))
				}
				if v, ok := offlineTaskAddEntityMap["crontab_expression"]; ok {
					offlineTaskAddParam.CrontabExpression = helper.String(v.(string))
				}
				if v, ok := offlineTaskAddEntityMap["retry_wait"]; ok {
					offlineTaskAddParam.RetryWait = helper.IntUint64(v.(int))
				}
				if v, ok := offlineTaskAddEntityMap["retriable"]; ok {
					offlineTaskAddParam.Retriable = helper.IntUint64(v.(int))
				}
				if v, ok := offlineTaskAddEntityMap["try_limit"]; ok {
					offlineTaskAddParam.TryLimit = helper.IntUint64(v.(int))
				}
				if v, ok := offlineTaskAddEntityMap["run_priority"]; ok {
					offlineTaskAddParam.RunPriority = helper.IntUint64(v.(int))
				}
				if v, ok := offlineTaskAddEntityMap["product_name"]; ok {
					offlineTaskAddParam.ProductName = helper.String(v.(string))
				}
				if v, ok := offlineTaskAddEntityMap["self_depend"]; ok {
					offlineTaskAddParam.SelfDepend = helper.IntUint64(v.(int))
				}
				if v, ok := offlineTaskAddEntityMap["task_action"]; ok {
					offlineTaskAddParam.TaskAction = helper.String(v.(string))
				}
				if v, ok := offlineTaskAddEntityMap["execution_end_time"]; ok {
					offlineTaskAddParam.ExecutionEndTime = helper.String(v.(string))
				}
				if v, ok := offlineTaskAddEntityMap["execution_start_time"]; ok {
					offlineTaskAddParam.ExecutionStartTime = helper.String(v.(string))
				}
				if v, ok := offlineTaskAddEntityMap["task_auto_submit"]; ok {
					offlineTaskAddParam.TaskAutoSubmit = helper.Bool(v.(bool))
				}
				if v, ok := offlineTaskAddEntityMap["instance_init_strategy"]; ok {
					offlineTaskAddParam.InstanceInitStrategy = helper.String(v.(string))
				}
				integrationTaskInfo.OfflineTaskAddEntity = &offlineTaskAddParam
			}
			if v, ok := dMap["executor_group_name"]; ok {
				integrationTaskInfo.ExecutorGroupName = helper.String(v.(string))
			}
			if v, ok := dMap["in_long_manager_url"]; ok {
				integrationTaskInfo.InLongManagerUrl = helper.String(v.(string))
			}
			if v, ok := dMap["in_long_stream_id"]; ok {
				integrationTaskInfo.InLongStreamId = helper.String(v.(string))
			}
			if v, ok := dMap["in_long_manager_version"]; ok {
				integrationTaskInfo.InLongManagerVersion = helper.String(v.(string))
			}
			if v, ok := dMap["data_proxy_url"]; ok {
				dataProxyUrlSet := v.(*schema.Set).List()
				for i := range dataProxyUrlSet {
					dataProxyUrl := dataProxyUrlSet[i].(string)
					integrationTaskInfo.DataProxyUrl = append(integrationTaskInfo.DataProxyUrl, &dataProxyUrl)
				}
			}
			if v, ok := dMap["submit"]; ok {
				integrationTaskInfo.Submit = helper.Bool(v.(bool))
			}
			if v, ok := dMap["input_datasource_type"]; ok {
				integrationTaskInfo.InputDatasourceType = helper.String(v.(string))
			}
			if v, ok := dMap["output_datasource_type"]; ok {
				integrationTaskInfo.OutputDatasourceType = helper.String(v.(string))
			}
			if v, ok := dMap["num_records_in"]; ok {
				integrationTaskInfo.NumRecordsIn = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["num_records_out"]; ok {
				integrationTaskInfo.NumRecordsOut = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["reader_delay"]; ok {
				integrationTaskInfo.ReaderDelay = helper.Float64(v.(float64))
			}
			if v, ok := dMap["num_restarts"]; ok {
				integrationTaskInfo.NumRestarts = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["create_time"]; ok {
				integrationTaskInfo.CreateTime = helper.String(v.(string))
			}
			if v, ok := dMap["update_time"]; ok {
				integrationTaskInfo.UpdateTime = helper.String(v.(string))
			}
			if v, ok := dMap["last_run_time"]; ok {
				integrationTaskInfo.LastRunTime = helper.String(v.(string))
			}
			if v, ok := dMap["stop_time"]; ok {
				integrationTaskInfo.StopTime = helper.String(v.(string))
			}
			if v, ok := dMap["has_version"]; ok {
				integrationTaskInfo.HasVersion = helper.Bool(v.(bool))
			}
			if v, ok := dMap["locked"]; ok {
				integrationTaskInfo.Locked = helper.Bool(v.(bool))
			}
			if v, ok := dMap["locker"]; ok {
				integrationTaskInfo.Locker = helper.String(v.(string))
			}
			if v, ok := dMap["running_cu"]; ok {
				integrationTaskInfo.RunningCu = helper.Float64(v.(float64))
			}
			if v, ok := dMap["task_alarm_regular_list"]; ok {
				taskAlarmRegularListSet := v.(*schema.Set).List()
				for i := range taskAlarmRegularListSet {
					taskAlarmRegularList := taskAlarmRegularListSet[i].(string)
					integrationTaskInfo.TaskAlarmRegularList = append(integrationTaskInfo.TaskAlarmRegularList, &taskAlarmRegularList)
				}
			}
			if v, ok := dMap["switch_resource"]; ok {
				integrationTaskInfo.SwitchResource = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["read_phase"]; ok {
				integrationTaskInfo.ReadPhase = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["instance_version"]; ok {
				integrationTaskInfo.InstanceVersion = helper.IntInt64(v.(int))
			}
			request.TaskInfo = &integrationTaskInfo
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().ModifyIntegrationTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update wedata integration_real_time_task failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegration_real_time_taskRead(d, meta)
}

func resourceTencentCloudWedataIntegration_real_time_taskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_real_time_task.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	integration_real_time_taskId := d.Id()

	if err := service.DeleteWedataIntegration_real_time_taskById(ctx, taskId); err != nil {
		return err
	}

	return nil
}

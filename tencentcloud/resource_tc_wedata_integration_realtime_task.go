/*
Provides a resource to create a wedata integration_realtime_task

Example Usage

```hcl
resource "tencentcloud_wedata_integration_realtime_task" "example" {
  project_id  = "1612982498218618880"
  task_name   = "tf_example"
  task_mode   = "1"
  description = "description."
  sync_type   = 1
  task_info {
    incharge    = "100028439226"
    executor_id = "20230313175748567418"
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "TaskManager"
      value = "1"
    }
    config {
      name  = "JobManager"
      value = "1"
    }
    config {
      name  = "TolerateDirtyData"
      value = "0"
    }
    config {
      name  = "CheckpointingInterval"
      value = "1"
    }
    config {
      name  = "CheckpointingIntervalUnit"
      value = "min"
    }
    config {
      name  = "RestartStrategyFixedDelayAttempts"
      value = "-1"
    }
    config {
      name  = "ResourceAllocationType"
      value = "0"
    }
    config {
      name  = "TaskAlarmRegularList"
      value = "35"
    }
  }
}
```

Import

wedata integration_realtime_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_integration_realtime_task.example 1776563389209296896#h9d39630a-ae45-4460-90b2-0b093cbfef5d
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWedataIntegrationRealtimeTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataIntegrationRealtimeTaskCreate,
		Read:   resourceTencentCloudWedataIntegrationRealtimeTaskRead,
		Update: resourceTencentCloudWedataIntegrationRealtimeTaskUpdate,
		Delete: resourceTencentCloudWedataIntegrationRealtimeTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// createTask
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},
			"task_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task name.",
			},
			"sync_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Synchronization type: 1. Whole database synchronization, 2. Single table synchronization.",
			},
			"task_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task display mode, 0: canvas mode, 1: form mode.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description information.",
			},
			// modifyTask
			"task_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Task Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workflow_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The workflow id to which the task belongs.",
						},
						"schedule_task_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Task scheduling id (job id such as oceanus or us).",
						},
						"task_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Inlong Task Group ID.",
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
						"app_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "User App Id.",
						},
						"status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Task status 1. Not started | Task initialization, 2. Task starting, 3. Running, 4. Paused, 5. Task stopping, 6. Stopped, 7. Execution failed, 8. deleted, 9. Locked, 404. unknown status.",
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
						//"task_mode": {
						//	Type:        schema.TypeString,
						//	Optional:    true,
						//	Description: "Task display mode, 0: canvas mode, 1: form mode.",
						//},
						"incharge": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Incharge user.",
						},
						//"offline_task_add_entity": {
						//	Type:        schema.TypeList,
						//	MaxItems:    1,
						//	Optional:    true,
						//	Description: "Offline task scheduling configuration.",
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"workflow_name": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "The name of the workflow to which the task belongs.",
						//			},
						//			"dependency_workflow": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Whether to support workflow dependencies: yes / no, default value: no.",
						//			},
						//			"start_time": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Effective start time, the format is yyyy-MM-dd HH:mm:ss.",
						//			},
						//			"end_time": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Effective end time, the format is yyyy-MM-dd HH:mm:ss.",
						//			},
						//			"cycle_type": {
						//				Type:        schema.TypeInt,
						//				Optional:    true,
						//				Description: "Scheduling type, 0: crontab type, 1: minutes, 2: hours, 3: days, 4: weeks, 5: months, 6: one-time, 7: user-driven, 10: elastic period (week), 11: elastic period (month), 12: year, 13: instant trigger.",
						//			},
						//			"cycle_step": {
						//				Type:        schema.TypeInt,
						//				Optional:    true,
						//				Description: "Interval time of scheduling, the minimum value: 1.",
						//			},
						//			"delay_time": {
						//				Type:        schema.TypeInt,
						//				Optional:    true,
						//				Description: "Execution time, unit is minutes, only available for day/week/month/year scheduling. For example, daily scheduling is executed once every day at 02:00, and the delayTime is 120 minutes.",
						//			},
						//			"crontab_expression": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Crontab expression.",
						//			},
						//			"retry_wait": {
						//				Type:        schema.TypeInt,
						//				Optional:    true,
						//				Description: "Retry waiting time, unit is minutes.",
						//			},
						//			"retriable": {
						//				Type:        schema.TypeInt,
						//				Optional:    true,
						//				Description: "Whether to retry.",
						//			},
						//			"try_limit": {
						//				Type:        schema.TypeInt,
						//				Optional:    true,
						//				Description: "Number of retries.",
						//			},
						//			"run_priority": {
						//				Type:        schema.TypeInt,
						//				Optional:    true,
						//				Description: "Task running priority.",
						//			},
						//			"product_name": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Product name: DATA_INTEGRATION.",
						//			},
						//			"self_depend": {
						//				Type:        schema.TypeInt,
						//				Optional:    true,
						//				Description: "Self-dependent rules, 1: Ordered serial one at a time, queued execution, 2: Unordered serial one at a time, not queued execution, 3: Parallel, multiple at once.",
						//			},
						//			"task_action": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Flexible cycle configuration, if it is a weekly task: 1 is Sunday, 2 is Monday, 3 is Tuesday, and so on. If it is a monthly task: &amp;#39;1,3&amp;#39; represents the 1st and 3rd; &amp;#39;L&amp;#39; represents the end of the month.",
						//			},
						//			"execution_end_time": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Scheduling execution end time.",
						//			},
						//			"execution_start_time": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Scheduling execution start time.",
						//			},
						//			"task_auto_submit": {
						//				Type:        schema.TypeBool,
						//				Optional:    true,
						//				Description: "Whether to automatically submit.",
						//			},
						//			"instance_init_strategy": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Instance initialization strategy.",
						//			},
						//		},
						//	},
						//},
						"executor_group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Executor group name.",
						},
						"in_long_manager_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "InLong manager url.",
						},
						"in_long_stream_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "InLong stream id.",
						},
						"in_long_manager_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "InLong manager version.",
						},
						"data_proxy_url": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Data proxy url.",
						},
						"submit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether the task version has been submitted for operation and maintenance.",
						},
						"input_datasource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Input datasource type.",
						},
						"output_datasource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Output datasource type.",
						},
						"num_records_in": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Number of reads.",
						},
						"num_records_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Number of writes.",
						},
						"reader_delay": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Computed:    true,
							Description: "Read latency.",
						},
						"num_restarts": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Times of restarts.",
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
						"last_run_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The last time the task was run.",
						},
						"stop_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The time the task was stopped.",
						},
						"has_version": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether the task been submitted.",
						},
						"locked": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether the task been locked.",
						},
						"locker": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "User locked task.",
						},
						"running_cu": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Computed:    true,
							Description: "The amount of resources consumed by real-time task.",
						},
						"task_alarm_regular_list": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Task alarm regular.",
						},
						"switch_resource": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Resource tiering status, 0: in progress, 1: successful, 2: failed.",
						},
						"read_phase": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Reading stage, 0: full amount, 1: partial full amount, 2: all incremental.",
						},
						"instance_version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Instance version.",
						},
					},
				},
			},
			"task_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Task ID.",
			},
		},
	}
}

func resourceTencentCloudWedataIntegrationRealtimeTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_realtime_task.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		request       = wedata.NewCreateIntegrationTaskRequest()
		response      = wedata.NewCreateIntegrationTaskResponse()
		modifyRequest = wedata.NewModifyIntegrationTaskRequest()
		projectId     string
		taskId        string
		taskName      string
		taskMode      string
		syncType      int
		description   string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	createIntegrationTaskInfo := wedata.IntegrationTaskInfo{}
	createIntegrationTaskInfo.ProjectId = &projectId
	createIntegrationTaskInfo.TaskType = helper.IntInt64(201)

	if v, ok := d.GetOk("task_name"); ok {
		createIntegrationTaskInfo.TaskName = helper.String(v.(string))
		taskName = v.(string)
	}

	if v, ok := d.GetOkExists("sync_type"); ok {
		createIntegrationTaskInfo.SyncType = helper.IntInt64(v.(int))
		syncType = v.(int)
	}

	if v, ok := d.GetOk("task_mode"); ok {
		createIntegrationTaskInfo.TaskMode = helper.String(v.(string))
		taskMode = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		createIntegrationTaskInfo.Description = helper.String(v.(string))
		description = v.(string)
	}

	request.TaskInfo = &createIntegrationTaskInfo
	// create
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
		log.Printf("[CRITAL]%s create wedata integrationRealtimeTask failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(strings.Join([]string{projectId, taskId}, FILED_SP))

	// modify
	if dMap, ok := helper.InterfacesHeadMap(d, "task_info"); ok {
		integrationTaskInfo := wedata.IntegrationTaskInfo{}
		integrationTaskInfo.TaskName = &taskName
		integrationTaskInfo.Description = &description
		integrationTaskInfo.SyncType = helper.IntInt64(syncType)
		integrationTaskInfo.TaskType = helper.IntInt64(201)
		integrationTaskInfo.TaskId = &taskId
		integrationTaskInfo.TaskMode = &taskMode

		//if v, ok := dMap["workflow_id"]; ok {
		//	integrationTaskInfo.WorkflowId = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["schedule_task_id"]; ok {
		//	integrationTaskInfo.ScheduleTaskId = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["task_group_id"]; ok {
		//	integrationTaskInfo.TaskGroupId = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["creator_uin"]; ok {
		//	integrationTaskInfo.CreatorUin = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["operator_uin"]; ok {
		//	integrationTaskInfo.OperatorUin = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["owner_uin"]; ok {
		//	integrationTaskInfo.OwnerUin = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["app_id"]; ok {
		//	integrationTaskInfo.AppId = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["status"]; ok {
		//	integrationTaskInfo.Status = helper.IntInt64(v.(int))
		//}

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

		if v, ok := dMap["incharge"]; ok {
			integrationTaskInfo.Incharge = helper.String(v.(string))
		}

		//if offlineTaskAddEntityMap, ok := helper.InterfaceToMap(dMap, "offline_task_add_entity"); ok {
		//	offlineTaskAddParam := wedata.OfflineTaskAddParam{}
		//	if v, ok := offlineTaskAddEntityMap["workflow_name"]; ok {
		//		offlineTaskAddParam.WorkflowName = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["dependency_workflow"]; ok {
		//		offlineTaskAddParam.DependencyWorkflow = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["start_time"]; ok {
		//		offlineTaskAddParam.StartTime = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["end_time"]; ok {
		//		offlineTaskAddParam.EndTime = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["cycle_type"]; ok {
		//		offlineTaskAddParam.CycleType = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["cycle_step"]; ok {
		//		offlineTaskAddParam.CycleStep = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["delay_time"]; ok {
		//		offlineTaskAddParam.DelayTime = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["crontab_expression"]; ok {
		//		offlineTaskAddParam.CrontabExpression = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["retry_wait"]; ok {
		//		offlineTaskAddParam.RetryWait = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["retriable"]; ok {
		//		offlineTaskAddParam.Retriable = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["try_limit"]; ok {
		//		offlineTaskAddParam.TryLimit = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["run_priority"]; ok {
		//		offlineTaskAddParam.RunPriority = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["product_name"]; ok {
		//		offlineTaskAddParam.ProductName = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["self_depend"]; ok {
		//		offlineTaskAddParam.SelfDepend = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["task_action"]; ok {
		//		offlineTaskAddParam.TaskAction = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["execution_end_time"]; ok {
		//		offlineTaskAddParam.ExecutionEndTime = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["execution_start_time"]; ok {
		//		offlineTaskAddParam.ExecutionStartTime = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["task_auto_submit"]; ok {
		//		offlineTaskAddParam.TaskAutoSubmit = helper.Bool(v.(bool))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["instance_init_strategy"]; ok {
		//		offlineTaskAddParam.InstanceInitStrategy = helper.String(v.(string))
		//	}
		//
		//	integrationTaskInfo.OfflineTaskAddEntity = &offlineTaskAddParam
		//}

		//if v, ok := dMap["executor_group_name"]; ok {
		//	integrationTaskInfo.ExecutorGroupName = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["in_long_manager_url"]; ok {
		//	integrationTaskInfo.InLongManagerUrl = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["in_long_stream_id"]; ok {
		//	integrationTaskInfo.InLongStreamId = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["in_long_manager_version"]; ok {
		//	integrationTaskInfo.InLongManagerVersion = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["data_proxy_url"]; ok {
		//	dataProxyUrlSet := v.(*schema.Set).List()
		//	for i := range dataProxyUrlSet {
		//		dataProxyUrl := dataProxyUrlSet[i].(string)
		//		integrationTaskInfo.DataProxyUrl = append(integrationTaskInfo.DataProxyUrl, &dataProxyUrl)
		//	}
		//}
		//
		//if v, ok := dMap["submit"]; ok {
		//	integrationTaskInfo.Submit = helper.Bool(v.(bool))
		//}
		//
		//if v, ok := dMap["input_datasource_type"]; ok {
		//	integrationTaskInfo.InputDatasourceType = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["output_datasource_type"]; ok {
		//	integrationTaskInfo.OutputDatasourceType = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["num_records_in"]; ok {
		//	integrationTaskInfo.NumRecordsIn = helper.IntInt64(v.(int))
		//}
		//
		//if v, ok := dMap["num_records_out"]; ok {
		//	integrationTaskInfo.NumRecordsOut = helper.IntInt64(v.(int))
		//}
		//
		//if v, ok := dMap["reader_delay"]; ok {
		//	integrationTaskInfo.ReaderDelay = helper.Float64(v.(float64))
		//}
		//
		//if v, ok := dMap["num_restarts"]; ok {
		//	integrationTaskInfo.NumRestarts = helper.IntInt64(v.(int))
		//}
		//
		//if v, ok := dMap["create_time"]; ok {
		//	integrationTaskInfo.CreateTime = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["update_time"]; ok {
		//	integrationTaskInfo.UpdateTime = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["last_run_time"]; ok {
		//	integrationTaskInfo.LastRunTime = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["stop_time"]; ok {
		//	integrationTaskInfo.StopTime = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["has_version"]; ok {
		//	integrationTaskInfo.HasVersion = helper.Bool(v.(bool))
		//}
		//
		//if v, ok := dMap["locked"]; ok {
		//	integrationTaskInfo.Locked = helper.Bool(v.(bool))
		//}
		//
		//if v, ok := dMap["locker"]; ok {
		//	integrationTaskInfo.Locker = helper.String(v.(string))
		//}
		//
		//if v, ok := dMap["running_cu"]; ok {
		//	integrationTaskInfo.RunningCu = helper.Float64(v.(float64))
		//}
		//
		//if v, ok := dMap["task_alarm_regular_list"]; ok {
		//	taskAlarmRegularListSet := v.(*schema.Set).List()
		//	for i := range taskAlarmRegularListSet {
		//		taskAlarmRegularList := taskAlarmRegularListSet[i].(string)
		//		integrationTaskInfo.TaskAlarmRegularList = append(integrationTaskInfo.TaskAlarmRegularList, &taskAlarmRegularList)
		//	}
		//}
		//
		//if v, ok := dMap["switch_resource"]; ok {
		//	integrationTaskInfo.SwitchResource = helper.IntInt64(v.(int))
		//}
		//
		//if v, ok := dMap["read_phase"]; ok {
		//	integrationTaskInfo.ReadPhase = helper.IntInt64(v.(int))
		//}
		//
		//if v, ok := dMap["instance_version"]; ok {
		//	integrationTaskInfo.InstanceVersion = helper.IntInt64(v.(int))
		//}

		modifyRequest.TaskInfo = &integrationTaskInfo
	}

	modifyRequest.ProjectId = &projectId
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().ModifyIntegrationTask(modifyRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyRequest.GetAction(), modifyRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("wedata integrationRealtimeTask not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify wedata integrationRealtimeTask failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegrationRealtimeTaskRead(d, meta)
}

func resourceTencentCloudWedataIntegrationRealtimeTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_realtime_task.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	taskId := idSplit[1]

	integrationRealtimeTask, err := service.DescribeWedataIntegrationRealtimeTaskById(ctx, projectId, taskId)
	if err != nil {
		return err
	}

	if integrationRealtimeTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataIntegrationRealtimeTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("task_id", taskId)

	if integrationRealtimeTask.TaskInfo != nil {
		taskInfoMap := map[string]interface{}{}

		if integrationRealtimeTask.TaskInfo.TaskName != nil {
			_ = d.Set("task_name", integrationRealtimeTask.TaskInfo.TaskName)
		}

		if integrationRealtimeTask.TaskInfo.Description != nil {
			_ = d.Set("description", integrationRealtimeTask.TaskInfo.Description)
		}

		if integrationRealtimeTask.TaskInfo.SyncType != nil {
			_ = d.Set("sync_type", integrationRealtimeTask.TaskInfo.SyncType)
		}

		if integrationRealtimeTask.TaskInfo.WorkflowId != nil {
			taskInfoMap["workflow_id"] = integrationRealtimeTask.TaskInfo.WorkflowId
		}

		if integrationRealtimeTask.TaskInfo.ScheduleTaskId != nil {
			taskInfoMap["schedule_task_id"] = integrationRealtimeTask.TaskInfo.ScheduleTaskId
		}

		if integrationRealtimeTask.TaskInfo.TaskGroupId != nil {
			taskInfoMap["task_group_id"] = integrationRealtimeTask.TaskInfo.TaskGroupId
		}

		if integrationRealtimeTask.TaskInfo.CreatorUin != nil {
			taskInfoMap["creator_uin"] = integrationRealtimeTask.TaskInfo.CreatorUin
		}

		if integrationRealtimeTask.TaskInfo.OperatorUin != nil {
			taskInfoMap["operator_uin"] = integrationRealtimeTask.TaskInfo.OperatorUin
		}

		if integrationRealtimeTask.TaskInfo.OwnerUin != nil {
			taskInfoMap["owner_uin"] = integrationRealtimeTask.TaskInfo.OwnerUin
		}

		if integrationRealtimeTask.TaskInfo.AppId != nil {
			taskInfoMap["app_id"] = integrationRealtimeTask.TaskInfo.AppId
		}

		if integrationRealtimeTask.TaskInfo.Status != nil {
			taskInfoMap["status"] = integrationRealtimeTask.TaskInfo.Status
		}

		if integrationRealtimeTask.TaskInfo.Nodes != nil {
			nodesList := []interface{}{}
			for _, nodes := range integrationRealtimeTask.TaskInfo.Nodes {
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

					nodesMap["config"] = configList
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

					nodesMap["ext_config"] = extConfigList
				}

				if nodes.Schema != nil {
					schemaList := []interface{}{}
					for _, nSchema := range nodes.Schema {
						schemaMap := map[string]interface{}{}

						if nSchema.Id != nil {
							schemaMap["id"] = nSchema.Id
						}

						if nSchema.Name != nil {
							schemaMap["name"] = nSchema.Name
						}

						if nSchema.Type != nil {
							schemaMap["type"] = nSchema.Type
						}

						if nSchema.Value != nil {
							schemaMap["value"] = nSchema.Value
						}

						if nSchema.Properties != nil {
							propertiesList := []interface{}{}
							for _, properties := range nSchema.Properties {
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

						if nSchema.Alias != nil {
							schemaMap["alias"] = nSchema.Alias
						}

						if nSchema.Comment != nil {
							schemaMap["comment"] = nSchema.Comment
						}

						schemaList = append(schemaList, schemaMap)
					}

					nodesMap["schema"] = schemaList
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

						nodeMappingMap["schema_mappings"] = schemaMappingsList
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

						nodeMappingMap["ext_config"] = extConfigList
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

			taskInfoMap["nodes"] = nodesList
		}

		if integrationRealtimeTask.TaskInfo.ExecutorId != nil {
			taskInfoMap["executor_id"] = integrationRealtimeTask.TaskInfo.ExecutorId
		}

		if integrationRealtimeTask.TaskInfo.Config != nil {
			configList := []interface{}{}
			for _, config := range integrationRealtimeTask.TaskInfo.Config {
				configMap := map[string]interface{}{}

				if config.Name != nil {
					configMap["name"] = config.Name
				}

				if config.Value != nil {
					configMap["value"] = config.Value
				}

				configList = append(configList, configMap)
			}

			taskInfoMap["config"] = configList
		}

		if integrationRealtimeTask.TaskInfo.ExtConfig != nil {
			extConfigList := []interface{}{}
			for _, extConfig := range integrationRealtimeTask.TaskInfo.ExtConfig {
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

		if integrationRealtimeTask.TaskInfo.ExecuteContext != nil {
			executeContextList := []interface{}{}
			for _, executeContext := range integrationRealtimeTask.TaskInfo.ExecuteContext {
				executeContextMap := map[string]interface{}{}

				if executeContext.Name != nil {
					executeContextMap["name"] = executeContext.Name
				}

				if executeContext.Value != nil {
					executeContextMap["value"] = executeContext.Value
				}

				executeContextList = append(executeContextList, executeContextMap)
			}

			taskInfoMap["execute_context"] = executeContextList
		}

		if integrationRealtimeTask.TaskInfo.Mappings != nil {
			mappingsList := []interface{}{}
			for _, mappings := range integrationRealtimeTask.TaskInfo.Mappings {
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

					mappingsMap["source_schema"] = sourceSchemaList
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

					mappingsMap["schema_mappings"] = schemaMappingsList
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

					mappingsMap["ext_config"] = extConfigList
				}

				mappingsList = append(mappingsList, mappingsMap)
			}

			taskInfoMap["mappings"] = mappingsList
		}

		if integrationRealtimeTask.TaskInfo.TaskMode != nil {
			_ = d.Set("task_mode", integrationRealtimeTask.TaskInfo.TaskMode)
		}

		if integrationRealtimeTask.TaskInfo.Incharge != nil {
			taskInfoMap["incharge"] = integrationRealtimeTask.TaskInfo.Incharge
		}

		if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity != nil {
			offlineTaskAddEntityMap := map[string]interface{}{}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.WorkflowName != nil {
				offlineTaskAddEntityMap["workflow_name"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.WorkflowName
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.DependencyWorkflow != nil {
				offlineTaskAddEntityMap["dependency_workflow"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.DependencyWorkflow
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.StartTime != nil {
				offlineTaskAddEntityMap["start_time"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.StartTime
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.EndTime != nil {
				offlineTaskAddEntityMap["end_time"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.EndTime
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.CycleType != nil {
				offlineTaskAddEntityMap["cycle_type"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.CycleType
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.CycleStep != nil {
				offlineTaskAddEntityMap["cycle_step"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.CycleStep
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.DelayTime != nil {
				offlineTaskAddEntityMap["delay_time"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.DelayTime
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.CrontabExpression != nil {
				offlineTaskAddEntityMap["crontab_expression"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.CrontabExpression
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.RetryWait != nil {
				offlineTaskAddEntityMap["retry_wait"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.RetryWait
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.Retriable != nil {
				offlineTaskAddEntityMap["retriable"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.Retriable
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.TryLimit != nil {
				offlineTaskAddEntityMap["try_limit"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.TryLimit
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.RunPriority != nil {
				offlineTaskAddEntityMap["run_priority"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.RunPriority
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.ProductName != nil {
				offlineTaskAddEntityMap["product_name"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.ProductName
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.SelfDepend != nil {
				offlineTaskAddEntityMap["self_depend"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.SelfDepend
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.TaskAction != nil {
				offlineTaskAddEntityMap["task_action"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.TaskAction
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.ExecutionEndTime != nil {
				offlineTaskAddEntityMap["execution_end_time"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.ExecutionEndTime
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.ExecutionStartTime != nil {
				offlineTaskAddEntityMap["execution_start_time"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.ExecutionStartTime
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.TaskAutoSubmit != nil {
				offlineTaskAddEntityMap["task_auto_submit"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.TaskAutoSubmit
			}

			if integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.InstanceInitStrategy != nil {
				offlineTaskAddEntityMap["instance_init_strategy"] = integrationRealtimeTask.TaskInfo.OfflineTaskAddEntity.InstanceInitStrategy
			}

			taskInfoMap["offline_task_add_entity"] = []interface{}{offlineTaskAddEntityMap}
		}

		if integrationRealtimeTask.TaskInfo.ExecutorGroupName != nil {
			taskInfoMap["executor_group_name"] = integrationRealtimeTask.TaskInfo.ExecutorGroupName
		}

		if integrationRealtimeTask.TaskInfo.InLongManagerUrl != nil {
			taskInfoMap["in_long_manager_url"] = integrationRealtimeTask.TaskInfo.InLongManagerUrl
		}

		if integrationRealtimeTask.TaskInfo.InLongStreamId != nil {
			taskInfoMap["in_long_stream_id"] = integrationRealtimeTask.TaskInfo.InLongStreamId
		}

		if integrationRealtimeTask.TaskInfo.InLongManagerVersion != nil {
			taskInfoMap["in_long_manager_version"] = integrationRealtimeTask.TaskInfo.InLongManagerVersion
		}

		if integrationRealtimeTask.TaskInfo.DataProxyUrl != nil {
			taskInfoMap["data_proxy_url"] = integrationRealtimeTask.TaskInfo.DataProxyUrl
		}

		if integrationRealtimeTask.TaskInfo.Submit != nil {
			taskInfoMap["submit"] = integrationRealtimeTask.TaskInfo.Submit
		}

		if integrationRealtimeTask.TaskInfo.InputDatasourceType != nil {
			taskInfoMap["input_datasource_type"] = integrationRealtimeTask.TaskInfo.InputDatasourceType
		}

		if integrationRealtimeTask.TaskInfo.OutputDatasourceType != nil {
			taskInfoMap["output_datasource_type"] = integrationRealtimeTask.TaskInfo.OutputDatasourceType
		}

		if integrationRealtimeTask.TaskInfo.NumRecordsIn != nil {
			taskInfoMap["num_records_in"] = integrationRealtimeTask.TaskInfo.NumRecordsIn
		}

		if integrationRealtimeTask.TaskInfo.NumRecordsOut != nil {
			taskInfoMap["num_records_out"] = integrationRealtimeTask.TaskInfo.NumRecordsOut
		}

		if integrationRealtimeTask.TaskInfo.ReaderDelay != nil {
			taskInfoMap["reader_delay"] = integrationRealtimeTask.TaskInfo.ReaderDelay
		}

		if integrationRealtimeTask.TaskInfo.NumRestarts != nil {
			taskInfoMap["num_restarts"] = integrationRealtimeTask.TaskInfo.NumRestarts
		}

		if integrationRealtimeTask.TaskInfo.CreateTime != nil {
			taskInfoMap["create_time"] = integrationRealtimeTask.TaskInfo.CreateTime
		}

		if integrationRealtimeTask.TaskInfo.UpdateTime != nil {
			taskInfoMap["update_time"] = integrationRealtimeTask.TaskInfo.UpdateTime
		}

		if integrationRealtimeTask.TaskInfo.LastRunTime != nil {
			taskInfoMap["last_run_time"] = integrationRealtimeTask.TaskInfo.LastRunTime
		}

		if integrationRealtimeTask.TaskInfo.StopTime != nil {
			taskInfoMap["stop_time"] = integrationRealtimeTask.TaskInfo.StopTime
		}

		if integrationRealtimeTask.TaskInfo.HasVersion != nil {
			taskInfoMap["has_version"] = integrationRealtimeTask.TaskInfo.HasVersion
		}

		if integrationRealtimeTask.TaskInfo.Locked != nil {
			taskInfoMap["locked"] = integrationRealtimeTask.TaskInfo.Locked
		}

		if integrationRealtimeTask.TaskInfo.Locker != nil {
			taskInfoMap["locker"] = integrationRealtimeTask.TaskInfo.Locker
		}

		if integrationRealtimeTask.TaskInfo.RunningCu != nil {
			taskInfoMap["running_cu"] = integrationRealtimeTask.TaskInfo.RunningCu
		}

		if integrationRealtimeTask.TaskInfo.TaskAlarmRegularList != nil {
			taskInfoMap["task_alarm_regular_list"] = integrationRealtimeTask.TaskInfo.TaskAlarmRegularList
		}

		if integrationRealtimeTask.TaskInfo.SwitchResource != nil {
			taskInfoMap["switch_resource"] = integrationRealtimeTask.TaskInfo.SwitchResource
		}

		if integrationRealtimeTask.TaskInfo.ReadPhase != nil {
			taskInfoMap["read_phase"] = integrationRealtimeTask.TaskInfo.ReadPhase
		}

		if integrationRealtimeTask.TaskInfo.InstanceVersion != nil {
			taskInfoMap["instance_version"] = integrationRealtimeTask.TaskInfo.InstanceVersion
		}

		_ = d.Set("task_info", []interface{}{taskInfoMap})
	}

	return nil
}

func resourceTencentCloudWedataIntegrationRealtimeTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_realtime_task.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = wedata.NewModifyIntegrationTaskRequest()
	)

	immutableArgs := []string{"project_id", "task_mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	taskId := idSplit[1]

	request.ProjectId = &projectId
	if dMap, ok := helper.InterfacesHeadMap(d, "task_info"); ok {
		integrationTaskInfo := wedata.IntegrationTaskInfo{}
		integrationTaskInfo.TaskId = &taskId
		integrationTaskInfo.TaskType = helper.IntInt64(201)

		if v, ok := d.GetOk("task_name"); ok {
			integrationTaskInfo.TaskName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			integrationTaskInfo.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("sync_type"); ok {
			integrationTaskInfo.SyncType = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["workflow_id"]; ok {
			integrationTaskInfo.WorkflowId = helper.String(v.(string))
		}

		if v, ok := dMap["schedule_task_id"]; ok {
			integrationTaskInfo.ScheduleTaskId = helper.String(v.(string))
		}

		if v, ok := dMap["task_group_id"]; ok {
			integrationTaskInfo.TaskGroupId = helper.String(v.(string))
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

		if v, ok := dMap["incharge"]; ok {
			integrationTaskInfo.Incharge = helper.String(v.(string))
		}

		//if offlineTaskAddEntityMap, ok := helper.InterfaceToMap(dMap, "offline_task_add_entity"); ok {
		//	offlineTaskAddParam := wedata.OfflineTaskAddParam{}
		//	if v, ok := offlineTaskAddEntityMap["workflow_name"]; ok {
		//		offlineTaskAddParam.WorkflowName = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["dependency_workflow"]; ok {
		//		offlineTaskAddParam.DependencyWorkflow = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["start_time"]; ok {
		//		offlineTaskAddParam.StartTime = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["end_time"]; ok {
		//		offlineTaskAddParam.EndTime = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["cycle_type"]; ok {
		//		offlineTaskAddParam.CycleType = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["cycle_step"]; ok {
		//		offlineTaskAddParam.CycleStep = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["delay_time"]; ok {
		//		offlineTaskAddParam.DelayTime = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["crontab_expression"]; ok {
		//		offlineTaskAddParam.CrontabExpression = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["retry_wait"]; ok {
		//		offlineTaskAddParam.RetryWait = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["retriable"]; ok {
		//		offlineTaskAddParam.Retriable = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["try_limit"]; ok {
		//		offlineTaskAddParam.TryLimit = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["run_priority"]; ok {
		//		offlineTaskAddParam.RunPriority = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["product_name"]; ok {
		//		offlineTaskAddParam.ProductName = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["self_depend"]; ok {
		//		offlineTaskAddParam.SelfDepend = helper.IntUint64(v.(int))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["task_action"]; ok {
		//		offlineTaskAddParam.TaskAction = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["execution_end_time"]; ok {
		//		offlineTaskAddParam.ExecutionEndTime = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["execution_start_time"]; ok {
		//		offlineTaskAddParam.ExecutionStartTime = helper.String(v.(string))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["task_auto_submit"]; ok {
		//		offlineTaskAddParam.TaskAutoSubmit = helper.Bool(v.(bool))
		//	}
		//
		//	if v, ok := offlineTaskAddEntityMap["instance_init_strategy"]; ok {
		//		offlineTaskAddParam.InstanceInitStrategy = helper.String(v.(string))
		//	}
		//
		//	integrationTaskInfo.OfflineTaskAddEntity = &offlineTaskAddParam
		//}

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

	request.ProjectId = &projectId
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
		log.Printf("[CRITAL]%s update wedata integrationRealtimeTask failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegrationRealtimeTaskRead(d, meta)
}

func resourceTencentCloudWedataIntegrationRealtimeTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_realtime_task.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	taskId := idSplit[1]

	if err := service.DeleteWedataIntegrationRealtimeTaskById(ctx, projectId, taskId); err != nil {
		return err
	}

	return nil
}

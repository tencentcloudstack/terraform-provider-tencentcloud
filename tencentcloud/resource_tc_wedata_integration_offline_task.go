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

func resourceTencentCloudWedataIntegrationOfflineTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataIntegrationOfflineTaskCreate,
		Read:   resourceTencentCloudWedataIntegrationOfflineTaskRead,
		Update: resourceTencentCloudWedataIntegrationOfflineTaskUpdate,
		Delete: resourceTencentCloudWedataIntegrationOfflineTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// OfflineTask
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},
			"cycle_step": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Interval time of scheduling, the minimum value: 1.",
			},
			"delay_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Execution time, unit is minutes, only available for day/week/month/year scheduling. For example, daily scheduling is executed once every day at 02:00, and the delayTime is 120 minutes.",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Effective end time, the format is yyyy-MM-dd HH:mm:ss.",
			},
			"notes": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Description information.",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Effective start time, the format is yyyy-MM-dd HH:mm:ss.",
			},
			"task_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task name.",
			},
			"task_action": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Scheduling configuration: flexible period configuration, only available for hourly/weekly/monthly/yearly scheduling. If the hourly task is specified to run at 0:00, 3:00 and 4:00 every day, it is 0,3,4.",
			},
			"task_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task display mode, 0: canvas mode, 1: form mode.",
			},
			// IntegrationTask
			"task_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Task Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sync_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Synchronization type: 1. Whole database synchronization, 2. Single table synchronization.",
						},
						"workflow_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The workflow id to which the task belongs.",
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
							Description: "Task status 1. Not started | Task initialization, 2. Task starting, 3. Running, 4. Paused, 5. Task stopping, 6. Stopped, 7. Execution failed, 8. deleted, 9. Locked, 404. unknown status.",
						},
						//"nodes": {
						//	Type:        schema.TypeList,
						//	Optional:    true,
						//	Description: "Task Node Information.",
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Node ID.",
						//			},
						//			"task_id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "The task id to which the node belongs.",
						//			},
						//			"name": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Node Name.",
						//			},
						//			"node_type": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Node type: INPUT,OUTPUT,JOIN,FILTER,TRANSFORM.",
						//			},
						//			"data_source_type": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Data source type: MYSQL, POSTGRE, ORACLE, SQLSERVER, FTP, HIVE, HDFS, ICEBERG, KAFKA, HBASE, SPARK, TBASE, DB2, DM, GAUSSDB, GBASE, IMPALA, ES, S3_DATAINSIGHT, GREENPLUM, PHOENIX, SAP_HANA, SFTP, OCEANBASE, CLICKHOUSE, KUDU, VERTICA, REDIS, COS, DLC, DORIS, CKAFKA, DTS_KAFKA, S3, CDW, TDSQLC, TDSQL, MONGODB, SYBASE, REST_API, StarRocks, TCHOUSE_X.",
						//			},
						//			"description": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Node Description.",
						//			},
						//			"datasource_id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Datasource ID.",
						//			},
						//			"config": {
						//				Type:        schema.TypeList,
						//				Optional:    true,
						//				Description: "Node configuration information.",
						//				Elem: &schema.Resource{
						//					Schema: map[string]*schema.Schema{
						//						"name": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Configuration name.",
						//						},
						//						"value": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Configuration value.",
						//						},
						//					},
						//				},
						//			},
						//			"ext_config": {
						//				Type:        schema.TypeList,
						//				Optional:    true,
						//				Description: "Node extension configuration information.",
						//				Elem: &schema.Resource{
						//					Schema: map[string]*schema.Schema{
						//						"name": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Configuration name.",
						//						},
						//						"value": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Configuration value.",
						//						},
						//					},
						//				},
						//			},
						//			"schema": {
						//				Type:        schema.TypeList,
						//				Optional:    true,
						//				Description: "Schema information.",
						//				Elem: &schema.Resource{
						//					Schema: map[string]*schema.Schema{
						//						"id": {
						//							Type:        schema.TypeString,
						//							Required:    true,
						//							Description: "Schema ID.",
						//						},
						//						"name": {
						//							Type:        schema.TypeString,
						//							Required:    true,
						//							Description: "Schema name.",
						//						},
						//						"type": {
						//							Type:        schema.TypeString,
						//							Required:    true,
						//							Description: "Schema type.",
						//						},
						//						"value": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Schema value.",
						//						},
						//						"properties": {
						//							Type:        schema.TypeList,
						//							Optional:    true,
						//							Description: "Schema extended attributes.",
						//							Elem: &schema.Resource{
						//								Schema: map[string]*schema.Schema{
						//									"name": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Attributes name.",
						//									},
						//									"value": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Attributes value.",
						//									},
						//								},
						//							},
						//						},
						//						"alias": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Schema alias.",
						//						},
						//						"comment": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Schema comment.",
						//						},
						//					},
						//				},
						//			},
						//			"node_mapping": {
						//				Type:        schema.TypeList,
						//				MaxItems:    1,
						//				Optional:    true,
						//				Description: "Node mapping.",
						//				Elem: &schema.Resource{
						//					Schema: map[string]*schema.Schema{
						//						"source_id": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Source node ID.",
						//						},
						//						"sink_id": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Sink node ID.",
						//						},
						//						"source_schema": {
						//							Type:        schema.TypeList,
						//							Optional:    true,
						//							Description: "Source node schema information.",
						//							Elem: &schema.Resource{
						//								Schema: map[string]*schema.Schema{
						//									"id": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema ID.",
						//									},
						//									"name": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema name.",
						//									},
						//									"type": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema type.",
						//									},
						//									"value": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Schema value.",
						//									},
						//									"properties": {
						//										Type:        schema.TypeList,
						//										Optional:    true,
						//										Description: "Schema extended attributes.",
						//										Elem: &schema.Resource{
						//											Schema: map[string]*schema.Schema{
						//												"name": {
						//													Type:        schema.TypeString,
						//													Optional:    true,
						//													Description: "Attributes name.",
						//												},
						//												"value": {
						//													Type:        schema.TypeString,
						//													Optional:    true,
						//													Description: "Attributes value.",
						//												},
						//											},
						//										},
						//									},
						//									"alias": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Schema alias.",
						//									},
						//									"comment": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Schema comment.",
						//									},
						//								},
						//							},
						//						},
						//						"schema_mappings": {
						//							Type:        schema.TypeList,
						//							Optional:    true,
						//							Description: "Schema mapping information.",
						//							Elem: &schema.Resource{
						//								Schema: map[string]*schema.Schema{
						//									"source_schema_id": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema ID from source node.",
						//									},
						//									"sink_schema_id": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema ID from sink node.",
						//									},
						//								},
						//							},
						//						},
						//						"ext_config": {
						//							Type:        schema.TypeList,
						//							Optional:    true,
						//							Description: "Node extension configuration information.",
						//							Elem: &schema.Resource{
						//								Schema: map[string]*schema.Schema{
						//									"name": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Configuration name.",
						//									},
						//									"value": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Configuration value.",
						//									},
						//								},
						//							},
						//						},
						//					},
						//				},
						//			},
						//			"app_id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "User App Id.",
						//			},
						//			"project_id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Project ID.",
						//			},
						//			"creator_uin": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Creator User ID.",
						//			},
						//			"operator_uin": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Operator User ID.",
						//			},
						//			"owner_uin": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Owner User ID.",
						//			},
						//			"create_time": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Create time.",
						//			},
						//			"update_time": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Update time.",
						//			},
						//		},
						//	},
						//},
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
									//"workflow_name": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "The name of the workflow to which the task belongs.",
									//},
									//"dependency_workflow": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Whether to support workflow dependencies: yes / no, default value: no.",
									//},
									//"start_time": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Effective start time, the format is yyyy-MM-dd HH:mm:ss.",
									//},
									//"end_time": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Effective end time, the format is yyyy-MM-dd HH:mm:ss.",
									//},
									"cycle_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Scheduling type, 0: crontab type, 1: minutes, 2: hours, 3: days, 4: weeks, 5: months, 6: one-time, 7: user-driven, 10: elastic period (week), 11: elastic period (month), 12: year, 13: instant trigger.",
									},
									//"cycle_step": {
									//	Type:        schema.TypeInt,
									//	Optional:    true,
									//	Description: "Interval time of scheduling, the minimum value: 1.",
									//},
									//"delay_time": {
									//	Type:        schema.TypeInt,
									//	Optional:    true,
									//	Description: "Execution time, unit is minutes, only available for day/week/month/year scheduling. For example, daily scheduling is executed once every day at 02:00, and the delayTime is 120 minutes.",
									//},
									"crontab_expression": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Crontab expression.",
									},
									"retry_wait": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Retry waiting time, unit is minutes.",
									},
									"retriable": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Whether to retry.",
									},
									"try_limit": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of retries.",
									},
									//"run_priority": {
									//	Type:        schema.TypeInt,
									//	Optional:    true,
									//	Description: "Task running priority.",
									//},
									//"product_name": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Product name: DATA_INTEGRATION.",
									//},
									"self_depend": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Self-dependent rules, 1: Ordered serial one at a time, queued execution, 2: Unordered serial one at a time, not queued execution, 3: Parallel, multiple at once.",
									},
									//"task_action": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Flexible cycle configuration, if it is a weekly task: 1 is Sunday, 2 is Monday, 3 is Tuesday, and so on. If it is a monthly task: 1, represents the 1st and 3rd; L represents the end of the month.",
									//},
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
									//"task_auto_submit": {
									//	Type:        schema.TypeBool,
									//	Optional:    true,
									//	Description: "Whether to automatically submit.",
									//},
									//"instance_init_strategy": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Instance initialization strategy.",
									//},
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
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
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
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
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
			// computed
			"task_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Task ID.",
			},
		},
	}
}

func resourceTencentCloudWedataIntegrationOfflineTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_offline_task.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                        = getLogId(contextNil)
		request                      = wedata.NewCreateOfflineTaskRequest()
		response                     = wedata.NewCreateOfflineTaskResponse()
		modifyIntegrationTaskRequest = wedata.NewModifyIntegrationTaskRequest()
		projectId                    string
		taskId                       string
		taskName                     string
		notes                        string
		taskAction                   string
		startTime                    string
		endTime                      string
		cycleStep                    int
		delayTime                    int
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOkExists("cycle_step"); ok {
		request.CycleStep = helper.IntInt64(v.(int))
		cycleStep = v.(int)
	}

	if v, ok := d.GetOkExists("delay_time"); ok {
		request.DelayTime = helper.IntInt64(v.(int))
		delayTime = v.(int)
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
		endTime = v.(string)
	}

	if v, ok := d.GetOk("notes"); ok {
		request.Notes = helper.String(v.(string))
		notes = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
		startTime = v.(string)
	}

	if v, ok := d.GetOk("task_name"); ok {
		request.TaskName = helper.String(v.(string))
		taskName = v.(string)
	}

	request.TypeId = helper.IntInt64(27)

	if v, ok := d.GetOk("task_action"); ok {
		request.TaskAction = helper.String(v.(string))
		taskAction = v.(string)
	}

	if v, ok := d.GetOk("task_mode"); ok {
		request.TaskMode = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateOfflineTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("wedata integrationOfflineTask not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata integrationOfflineTask failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(strings.Join([]string{projectId, taskId}, FILED_SP))

	// Create IntegrationTask

	if dMap, ok := helper.InterfacesHeadMap(d, "task_info"); ok {
		integrationTaskInfo := wedata.IntegrationTaskInfo{}
		integrationTaskInfo.ProjectId = &projectId
		integrationTaskInfo.TaskId = &taskId
		integrationTaskInfo.TaskName = &taskName
		integrationTaskInfo.Description = &notes
		integrationTaskInfo.TaskType = helper.IntInt64(202)

		if v, ok := dMap["sync_type"]; ok {
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

		//if v, ok := dMap["nodes"]; ok {
		//	for _, item := range v.([]interface{}) {
		//		nodesMap := item.(map[string]interface{})
		//		integrationNodeInfo := wedata.IntegrationNodeInfo{}
		//		if v, ok := nodesMap["id"]; ok {
		//			integrationNodeInfo.Id = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["task_id"]; ok {
		//			integrationNodeInfo.TaskId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["name"]; ok {
		//			integrationNodeInfo.Name = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["node_type"]; ok {
		//			integrationNodeInfo.NodeType = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["data_source_type"]; ok {
		//			integrationNodeInfo.DataSourceType = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["description"]; ok {
		//			integrationNodeInfo.Description = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["datasource_id"]; ok {
		//			integrationNodeInfo.DatasourceId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["config"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				configMap := item.(map[string]interface{})
		//				recordField := wedata.RecordField{}
		//				if v, ok := configMap["name"]; ok {
		//					recordField.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := configMap["value"]; ok {
		//					recordField.Value = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.Config = append(integrationNodeInfo.Config, &recordField)
		//			}
		//		}
		//
		//		if v, ok := nodesMap["ext_config"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				extConfigMap := item.(map[string]interface{})
		//				recordField := wedata.RecordField{}
		//				if v, ok := extConfigMap["name"]; ok {
		//					recordField.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := extConfigMap["value"]; ok {
		//					recordField.Value = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.ExtConfig = append(integrationNodeInfo.ExtConfig, &recordField)
		//			}
		//		}
		//
		//		if v, ok := nodesMap["schema"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				schemaMap := item.(map[string]interface{})
		//				integrationNodeSchema := wedata.IntegrationNodeSchema{}
		//				if v, ok := schemaMap["id"]; ok {
		//					integrationNodeSchema.Id = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["name"]; ok {
		//					integrationNodeSchema.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["type"]; ok {
		//					integrationNodeSchema.Type = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["value"]; ok {
		//					integrationNodeSchema.Value = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["properties"]; ok {
		//					for _, item := range v.([]interface{}) {
		//						propertiesMap := item.(map[string]interface{})
		//						recordField := wedata.RecordField{}
		//						if v, ok := propertiesMap["name"]; ok {
		//							recordField.Name = helper.String(v.(string))
		//						}
		//
		//						if v, ok := propertiesMap["value"]; ok {
		//							recordField.Value = helper.String(v.(string))
		//						}
		//
		//						integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
		//					}
		//				}
		//
		//				if v, ok := schemaMap["alias"]; ok {
		//					integrationNodeSchema.Alias = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["comment"]; ok {
		//					integrationNodeSchema.Comment = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.Schema = append(integrationNodeInfo.Schema, &integrationNodeSchema)
		//			}
		//		}
		//		if nodeMappingMap, ok := helper.InterfaceToMap(nodesMap, "node_mapping"); ok {
		//			integrationNodeMapping := wedata.IntegrationNodeMapping{}
		//			if v, ok := nodeMappingMap["source_id"]; ok {
		//				integrationNodeMapping.SourceId = helper.String(v.(string))
		//			}
		//
		//			if v, ok := nodeMappingMap["sink_id"]; ok {
		//				integrationNodeMapping.SinkId = helper.String(v.(string))
		//			}
		//
		//			if v, ok := nodeMappingMap["source_schema"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					sourceSchemaMap := item.(map[string]interface{})
		//					integrationNodeSchema := wedata.IntegrationNodeSchema{}
		//					if v, ok := sourceSchemaMap["id"]; ok {
		//						integrationNodeSchema.Id = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["name"]; ok {
		//						integrationNodeSchema.Name = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["type"]; ok {
		//						integrationNodeSchema.Type = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["value"]; ok {
		//						integrationNodeSchema.Value = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["properties"]; ok {
		//						for _, item := range v.([]interface{}) {
		//							propertiesMap := item.(map[string]interface{})
		//							recordField := wedata.RecordField{}
		//							if v, ok := propertiesMap["name"]; ok {
		//								recordField.Name = helper.String(v.(string))
		//							}
		//
		//							if v, ok := propertiesMap["value"]; ok {
		//								recordField.Value = helper.String(v.(string))
		//							}
		//
		//							integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
		//						}
		//					}
		//
		//					if v, ok := sourceSchemaMap["alias"]; ok {
		//						integrationNodeSchema.Alias = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["comment"]; ok {
		//						integrationNodeSchema.Comment = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.SourceSchema = append(integrationNodeMapping.SourceSchema, &integrationNodeSchema)
		//				}
		//			}
		//
		//			if v, ok := nodeMappingMap["schema_mappings"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					schemaMappingsMap := item.(map[string]interface{})
		//					integrationNodeSchemaMapping := wedata.IntegrationNodeSchemaMapping{}
		//					if v, ok := schemaMappingsMap["source_schema_id"]; ok {
		//						integrationNodeSchemaMapping.SourceSchemaId = helper.String(v.(string))
		//					}
		//
		//					if v, ok := schemaMappingsMap["sink_schema_id"]; ok {
		//						integrationNodeSchemaMapping.SinkSchemaId = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.SchemaMappings = append(integrationNodeMapping.SchemaMappings, &integrationNodeSchemaMapping)
		//				}
		//			}
		//
		//			if v, ok := nodeMappingMap["ext_config"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					extConfigMap := item.(map[string]interface{})
		//					recordField := wedata.RecordField{}
		//					if v, ok := extConfigMap["name"]; ok {
		//						recordField.Name = helper.String(v.(string))
		//					}
		//
		//					if v, ok := extConfigMap["value"]; ok {
		//						recordField.Value = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.ExtConfig = append(integrationNodeMapping.ExtConfig, &recordField)
		//				}
		//			}
		//
		//			integrationNodeInfo.NodeMapping = &integrationNodeMapping
		//		}
		//
		//		if v, ok := nodesMap["app_id"]; ok {
		//			integrationNodeInfo.AppId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["project_id"]; ok {
		//			integrationNodeInfo.ProjectId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["creator_uin"]; ok {
		//			integrationNodeInfo.CreatorUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["operator_uin"]; ok {
		//			integrationNodeInfo.OperatorUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["owner_uin"]; ok {
		//			integrationNodeInfo.OwnerUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["create_time"]; ok {
		//			integrationNodeInfo.CreateTime = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["update_time"]; ok {
		//			integrationNodeInfo.UpdateTime = helper.String(v.(string))
		//		}
		//		integrationTaskInfo.Nodes = append(integrationTaskInfo.Nodes, &integrationNodeInfo)
		//	}
		//}

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
			//if v, ok := offlineTaskAddEntityMap["workflow_name"]; ok {
			//	offlineTaskAddParam.WorkflowName = helper.String(v.(string))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["dependency_workflow"]; ok {
			//	offlineTaskAddParam.DependencyWorkflow = helper.String(v.(string))
			//}

			offlineTaskAddParam.StartTime = &startTime
			offlineTaskAddParam.EndTime = &endTime
			offlineTaskAddParam.CycleStep = helper.IntUint64(cycleStep)
			offlineTaskAddParam.DelayTime = helper.IntUint64(delayTime)
			offlineTaskAddParam.TaskAction = &taskAction

			if v, ok := offlineTaskAddEntityMap["cycle_type"]; ok {
				offlineTaskAddParam.CycleType = helper.IntUint64(v.(int))
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

			//if v, ok := offlineTaskAddEntityMap["run_priority"]; ok {
			//	offlineTaskAddParam.RunPriority = helper.IntUint64(v.(int))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["product_name"]; ok {
			//	offlineTaskAddParam.ProductName = helper.String(v.(string))
			//}

			if v, ok := offlineTaskAddEntityMap["self_depend"]; ok {
				offlineTaskAddParam.SelfDepend = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["execution_end_time"]; ok {
				offlineTaskAddParam.ExecutionEndTime = helper.String(v.(string))
			}

			if v, ok := offlineTaskAddEntityMap["execution_start_time"]; ok {
				offlineTaskAddParam.ExecutionStartTime = helper.String(v.(string))
			}

			//if v, ok := offlineTaskAddEntityMap["task_auto_submit"]; ok {
			//	offlineTaskAddParam.TaskAutoSubmit = helper.Bool(v.(bool))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["instance_init_strategy"]; ok {
			//	offlineTaskAddParam.InstanceInitStrategy = helper.String(v.(string))
			//}

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

		modifyIntegrationTaskRequest.TaskInfo = &integrationTaskInfo
	}

	modifyIntegrationTaskRequest.ProjectId = &projectId
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().ModifyIntegrationTask(modifyIntegrationTaskRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyIntegrationTaskRequest.GetAction(), modifyIntegrationTaskRequest.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata integration_real_time_task failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegrationOfflineTaskRead(d, meta)
}

func resourceTencentCloudWedataIntegrationOfflineTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_offline_task.read")()
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

	integrationOfflineTask, err := service.DescribeWedataIntegrationOfflineTaskById(ctx, projectId, taskId)
	if err != nil {
		return err
	}

	if integrationOfflineTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataIntegrationOfflineTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("task_id", taskId)

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.CycleStep != nil {
		_ = d.Set("cycle_step", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.CycleStep)
	}

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.DelayTime != nil {
		_ = d.Set("delay_time", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.DelayTime)
	}

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.EndTime != nil {
		_ = d.Set("end_time", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.EndTime)
	}

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.StartTime != nil {
		_ = d.Set("start_time", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.StartTime)
	}

	if integrationOfflineTask.TaskInfo.TaskName != nil {
		_ = d.Set("task_name", integrationOfflineTask.TaskInfo.TaskName)
	}

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.TaskAction != nil {
		_ = d.Set("task_action", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.TaskAction)
	}

	if integrationOfflineTask.TaskInfo.TaskMode != nil {
		_ = d.Set("task_mode", integrationOfflineTask.TaskInfo.TaskMode)
	}

	return nil
}

func resourceTencentCloudWedataIntegrationOfflineTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_offline_task.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = wedata.NewModifyIntegrationTaskRequest()
		taskName   string
		notes      string
		taskAction string
		startTime  string
		endTime    string
		cycleStep  int
		delayTime  int
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

	if v, ok := d.GetOkExists("cycle_step"); ok {
		cycleStep = v.(int)
	}

	if v, ok := d.GetOkExists("delay_time"); ok {
		delayTime = v.(int)
	}

	if v, ok := d.GetOk("end_time"); ok {
		endTime = v.(string)
	}

	if v, ok := d.GetOk("notes"); ok {
		notes = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		startTime = v.(string)
	}

	if v, ok := d.GetOk("task_name"); ok {
		taskName = v.(string)
	}

	if v, ok := d.GetOk("task_action"); ok {
		taskAction = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "task_info"); ok {
		integrationTaskInfo := wedata.IntegrationTaskInfo{}
		integrationTaskInfo.ProjectId = &projectId
		integrationTaskInfo.TaskId = &taskId
		integrationTaskInfo.TaskName = &taskName
		integrationTaskInfo.Description = &notes
		integrationTaskInfo.TaskType = helper.IntInt64(202)

		if v, ok := dMap["sync_type"]; ok {
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

		//if v, ok := dMap["nodes"]; ok {
		//	for _, item := range v.([]interface{}) {
		//		nodesMap := item.(map[string]interface{})
		//		integrationNodeInfo := wedata.IntegrationNodeInfo{}
		//		if v, ok := nodesMap["id"]; ok {
		//			integrationNodeInfo.Id = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["task_id"]; ok {
		//			integrationNodeInfo.TaskId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["name"]; ok {
		//			integrationNodeInfo.Name = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["node_type"]; ok {
		//			integrationNodeInfo.NodeType = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["data_source_type"]; ok {
		//			integrationNodeInfo.DataSourceType = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["description"]; ok {
		//			integrationNodeInfo.Description = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["datasource_id"]; ok {
		//			integrationNodeInfo.DatasourceId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["config"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				configMap := item.(map[string]interface{})
		//				recordField := wedata.RecordField{}
		//				if v, ok := configMap["name"]; ok {
		//					recordField.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := configMap["value"]; ok {
		//					recordField.Value = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.Config = append(integrationNodeInfo.Config, &recordField)
		//			}
		//		}
		//
		//		if v, ok := nodesMap["ext_config"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				extConfigMap := item.(map[string]interface{})
		//				recordField := wedata.RecordField{}
		//				if v, ok := extConfigMap["name"]; ok {
		//					recordField.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := extConfigMap["value"]; ok {
		//					recordField.Value = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.ExtConfig = append(integrationNodeInfo.ExtConfig, &recordField)
		//			}
		//		}
		//
		//		if v, ok := nodesMap["schema"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				schemaMap := item.(map[string]interface{})
		//				integrationNodeSchema := wedata.IntegrationNodeSchema{}
		//				if v, ok := schemaMap["id"]; ok {
		//					integrationNodeSchema.Id = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["name"]; ok {
		//					integrationNodeSchema.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["type"]; ok {
		//					integrationNodeSchema.Type = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["value"]; ok {
		//					integrationNodeSchema.Value = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["properties"]; ok {
		//					for _, item := range v.([]interface{}) {
		//						propertiesMap := item.(map[string]interface{})
		//						recordField := wedata.RecordField{}
		//						if v, ok := propertiesMap["name"]; ok {
		//							recordField.Name = helper.String(v.(string))
		//						}
		//
		//						if v, ok := propertiesMap["value"]; ok {
		//							recordField.Value = helper.String(v.(string))
		//						}
		//
		//						integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
		//					}
		//				}
		//
		//				if v, ok := schemaMap["alias"]; ok {
		//					integrationNodeSchema.Alias = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["comment"]; ok {
		//					integrationNodeSchema.Comment = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.Schema = append(integrationNodeInfo.Schema, &integrationNodeSchema)
		//			}
		//		}
		//		if nodeMappingMap, ok := helper.InterfaceToMap(nodesMap, "node_mapping"); ok {
		//			integrationNodeMapping := wedata.IntegrationNodeMapping{}
		//			if v, ok := nodeMappingMap["source_id"]; ok {
		//				integrationNodeMapping.SourceId = helper.String(v.(string))
		//			}
		//
		//			if v, ok := nodeMappingMap["sink_id"]; ok {
		//				integrationNodeMapping.SinkId = helper.String(v.(string))
		//			}
		//
		//			if v, ok := nodeMappingMap["source_schema"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					sourceSchemaMap := item.(map[string]interface{})
		//					integrationNodeSchema := wedata.IntegrationNodeSchema{}
		//					if v, ok := sourceSchemaMap["id"]; ok {
		//						integrationNodeSchema.Id = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["name"]; ok {
		//						integrationNodeSchema.Name = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["type"]; ok {
		//						integrationNodeSchema.Type = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["value"]; ok {
		//						integrationNodeSchema.Value = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["properties"]; ok {
		//						for _, item := range v.([]interface{}) {
		//							propertiesMap := item.(map[string]interface{})
		//							recordField := wedata.RecordField{}
		//							if v, ok := propertiesMap["name"]; ok {
		//								recordField.Name = helper.String(v.(string))
		//							}
		//
		//							if v, ok := propertiesMap["value"]; ok {
		//								recordField.Value = helper.String(v.(string))
		//							}
		//
		//							integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
		//						}
		//					}
		//
		//					if v, ok := sourceSchemaMap["alias"]; ok {
		//						integrationNodeSchema.Alias = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["comment"]; ok {
		//						integrationNodeSchema.Comment = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.SourceSchema = append(integrationNodeMapping.SourceSchema, &integrationNodeSchema)
		//				}
		//			}
		//
		//			if v, ok := nodeMappingMap["schema_mappings"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					schemaMappingsMap := item.(map[string]interface{})
		//					integrationNodeSchemaMapping := wedata.IntegrationNodeSchemaMapping{}
		//					if v, ok := schemaMappingsMap["source_schema_id"]; ok {
		//						integrationNodeSchemaMapping.SourceSchemaId = helper.String(v.(string))
		//					}
		//
		//					if v, ok := schemaMappingsMap["sink_schema_id"]; ok {
		//						integrationNodeSchemaMapping.SinkSchemaId = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.SchemaMappings = append(integrationNodeMapping.SchemaMappings, &integrationNodeSchemaMapping)
		//				}
		//			}
		//
		//			if v, ok := nodeMappingMap["ext_config"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					extConfigMap := item.(map[string]interface{})
		//					recordField := wedata.RecordField{}
		//					if v, ok := extConfigMap["name"]; ok {
		//						recordField.Name = helper.String(v.(string))
		//					}
		//
		//					if v, ok := extConfigMap["value"]; ok {
		//						recordField.Value = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.ExtConfig = append(integrationNodeMapping.ExtConfig, &recordField)
		//				}
		//			}
		//
		//			integrationNodeInfo.NodeMapping = &integrationNodeMapping
		//		}
		//
		//		if v, ok := nodesMap["app_id"]; ok {
		//			integrationNodeInfo.AppId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["project_id"]; ok {
		//			integrationNodeInfo.ProjectId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["creator_uin"]; ok {
		//			integrationNodeInfo.CreatorUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["operator_uin"]; ok {
		//			integrationNodeInfo.OperatorUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["owner_uin"]; ok {
		//			integrationNodeInfo.OwnerUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["create_time"]; ok {
		//			integrationNodeInfo.CreateTime = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["update_time"]; ok {
		//			integrationNodeInfo.UpdateTime = helper.String(v.(string))
		//		}
		//		integrationTaskInfo.Nodes = append(integrationTaskInfo.Nodes, &integrationNodeInfo)
		//	}
		//}

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
			//if v, ok := offlineTaskAddEntityMap["workflow_name"]; ok {
			//	offlineTaskAddParam.WorkflowName = helper.String(v.(string))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["dependency_workflow"]; ok {
			//	offlineTaskAddParam.DependencyWorkflow = helper.String(v.(string))
			//}

			offlineTaskAddParam.StartTime = &startTime
			offlineTaskAddParam.EndTime = &endTime
			offlineTaskAddParam.CycleStep = helper.IntUint64(cycleStep)
			offlineTaskAddParam.DelayTime = helper.IntUint64(delayTime)
			offlineTaskAddParam.TaskAction = &taskAction

			if v, ok := offlineTaskAddEntityMap["cycle_type"]; ok {
				offlineTaskAddParam.CycleType = helper.IntUint64(v.(int))
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

			//if v, ok := offlineTaskAddEntityMap["run_priority"]; ok {
			//	offlineTaskAddParam.RunPriority = helper.IntUint64(v.(int))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["product_name"]; ok {
			//	offlineTaskAddParam.ProductName = helper.String(v.(string))
			//}

			if v, ok := offlineTaskAddEntityMap["self_depend"]; ok {
				offlineTaskAddParam.SelfDepend = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["execution_end_time"]; ok {
				offlineTaskAddParam.ExecutionEndTime = helper.String(v.(string))
			}

			if v, ok := offlineTaskAddEntityMap["execution_start_time"]; ok {
				offlineTaskAddParam.ExecutionStartTime = helper.String(v.(string))
			}

			//if v, ok := offlineTaskAddEntityMap["task_auto_submit"]; ok {
			//	offlineTaskAddParam.TaskAutoSubmit = helper.Bool(v.(bool))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["instance_init_strategy"]; ok {
			//	offlineTaskAddParam.InstanceInitStrategy = helper.String(v.(string))
			//}

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
		log.Printf("[CRITAL]%s modify wedata integration_real_time_task failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegrationOfflineTaskRead(d, meta)
}

func resourceTencentCloudWedataIntegrationOfflineTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_offline_task.delete")()
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

	if err := service.DeleteWedataIntegrationOfflineTaskById(ctx, projectId, taskId); err != nil {
		return err
	}

	return nil
}

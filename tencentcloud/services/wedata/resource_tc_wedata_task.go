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

func ResourceTencentCloudWedataTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataTaskCreate,
		Read:   resourceTencentCloudWedataTaskRead,
		Update: resourceTencentCloudWedataTaskUpdate,
		Delete: resourceTencentCloudWedataTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID.",
			},

			"task_base_attribute": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Basic task attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task name.",
						},
						"task_type_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task type ID:\n\n* 21:JDBC SQL\n* 23:TDSQL-PostgreSQL\n* 26:OfflineSynchronization\n* 30:Python\n* 31:PySpark\n* 32:DLC SQL\n* 33:Impala\n* 34:Hive SQL\n* 35:Shell\n* 36:Spark SQL\n* 38:Shell Form Mode\n* 39:Spark\n* 40:TCHouse-P\n* 41:Kettle\n* 42:Tchouse-X\n* 43:TCHouse-X SQL\n* 46:DLC Spark\n* 47:TiOne\n* 48:Trino\n* 50:DLC PySpark\n* 92:MapReduce\n* 130:Branch Node\n* 131:Merged Node\n* 132:Notebook\n* 133:SSH\n* 134:StarRocks\n* 137:For-each\n* 138:Setats SQL.",
						},
						"workflow_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Workflow ID.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Task owner ID, defaults to current user.",
						},
						"task_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Task description.",
						},
					},
				},
			},

			"task_configuration": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Task configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_group": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Resource group ID: Need to obtain ExecutorGroupId via DescribeNormalSchedulerExecutorGroups.",
						},
						"code_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Base64 encoded code content.",
						},
						"task_ext_configuration_list": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Task extended attribute configuration list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"param_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter name.",
									},
									"param_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter value.",
									},
								},
							},
						},
						"data_cluster": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"broker_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specified running node.",
						},
						"yarn_queue": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Resource pool queue name, need to obtain via DescribeProjectClusterQueues.",
						},
						"source_service_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Source data source ID, separated by `;`, need to obtain via DescribeDataSourceWithoutInfo.",
						},
						"target_service_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Target data source ID, separated by `;`, need to obtain via DescribeDataSourceWithoutInfo.",
						},
						"task_scheduling_parameter_list": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Scheduling parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"param_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter name.",
									},
									"param_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter value.",
									},
								},
							},
						},
						"bundle_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "ID used by Bundle.",
						},
						"bundle_info": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Bundle information.",
						},
					},
				},
			},

			"task_scheduler_configuration": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Task scheduling configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cycle_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Cycle type: Defaults to DAY_CYCLE.\n\nSupported types are\n\n* ONEOFF_CYCLE: One-time\n* YEAR_CYCLE: Yearly\n* MONTH_CYCLE: Monthly\n* WEEK_CYCLE: Weekly\n* DAY_CYCLE: Daily\n* HOUR_CYCLE: Hourly\n* MINUTE_CYCLE: Minutely\n* CRONTAB_CYCLE: Crontab expression type.",
						},
						"schedule_time_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Time zone, defaults to UTC+8.",
						},
						"crontab_expression": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Cron expression, defaults to 0 0 0 * * `?` *.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Effective date, defaults to 00:00:00 of current date.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "End date, defaults to 2099-12-31 23:59:59.",
						},
						"execution_start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Execution time left-closed interval, default 00:00.",
						},
						"execution_end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Execution time right-closed interval, default 23:59.",
						},
						"schedule_run_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Scheduling type: 0 Normal scheduling 1 Empty run scheduling, defaults to 0.",
						},
						"calendar_open": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Calendar scheduling: Values are 0 and 1, 1 for enabled, 0 for disabled, defaults to 0.",
						},
						"calendar_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Calendar scheduling calendar ID.",
						},
						"self_depend": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Self-dependency, default value serial, values: parallel, serial, orderly.",
						},
						"upstream_dependency_config_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Upstream dependency array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Task ID.",
									},
									"main_cyclic_config": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Main dependency configuration, values:\n\n* CRONTAB\n* DAY\n* HOUR\n* LIST_DAY\n* LIST_HOUR\n* LIST_MINUTE\n* MINUTE\n* MONTH\n* RANGE_DAY\n* RANGE_HOUR\n* RANGE_MINUTE\n* WEEK\n* YEAR.",
									},
									"subordinate_cyclic_config": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Secondary dependency configuration, values:\n* ALL_DAY_OF_YEAR\n* ALL_MONTH_OF_YEAR\n* CURRENT\n* CURRENT_DAY\n* CURRENT_HOUR\n* CURRENT_MINUTE\n* CURRENT_MONTH\n* CURRENT_WEEK\n* CURRENT_YEAR\n* PREVIOUS_BEGIN_OF_MONTH\n* PREVIOUS_DAY\n* PREVIOUS_DAY_LATER_OFFSET_HOUR\n* PREVIOUS_DAY_LATER_OFFSET_MINUTE\n* PREVIOUS_END_OF_MONTH\n* PREVIOUS_FRIDAY\n* PREVIOUS_HOUR\n* PREVIOUS_HOUR_CYCLE\n* PREVIOUS_HOUR_LATER_OFFSET_MINUTE\n* PREVIOUS_MINUTE_CYCLE\n* PREVIOUS_MONTH\n* PREVIOUS_WEEK\n* PREVIOUS_WEEKEND\n* RECENT_DATE.",
									},
									"offset": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Offset in interval and list modes.",
									},
									"dependency_strategy": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Dependency execution strategy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"polling_null_strategy": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Strategy for waiting for upstream task instances: EXECUTING; WAITING.",
												},
												"task_dependency_executing_strategies": {
													Type:        schema.TypeSet,
													Optional:    true,
													Computed:    true,
													Description: "This field is required only when PollingNullStrategy is EXECUTING, List type: NOT_EXIST (default, when minute depends on minute/hour depends on hour, parent instance is not within the downstream instance scheduling time range); PARENT_EXPIRED (parent instance failed); PARENT_TIMEOUT (parent instance timed out). If any of the above scenarios is met, the parent task instance dependency judgment can be passed; otherwise, it is necessary to wait for the parent instance.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"task_dependency_executing_timeout_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "This field is required only when TaskDependencyExecutingStrategies contains PARENT_TIMEOUT, the timeout time for downstream tasks depending on parent instance execution, unit: minutes.",
												},
											},
										},
									},
								},
							},
						},
						"event_listener_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Event array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"event_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Event name.",
									},
									"event_sub_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Event cycle: SECOND, MIN, HOUR, DAY.",
									},
									"event_broadcast_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Event broadcast type: SINGLE, BROADCAST.",
									},
									"properties_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Extended information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"param_key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Parameter name.",
												},
												"param_value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Parameter value.",
												},
											},
										},
									},
								},
							},
						},
						"run_priority": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Task scheduling priority: 4 for high, 5 for medium, 6 for low, default: 6.",
						},
						"retry_wait": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Retry strategy - retry waiting time, unit: minutes: default: 5.",
						},
						"max_retry_attempts": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Retry strategy - maximum number of attempts, default: 4.",
						},
						"execution_ttl": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Timeout handling strategy - running time timeout (unit: minutes) defaults to -1.",
						},
						"wait_execution_total_ttl": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Timeout handling strategy - total waiting time timeout (unit: minutes) defaults to -1.",
						},
						"allow_redo_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Rerun & backfill configuration, defaults to ALL; ALL: can rerun or backfill after success or failure; FAILURE: cannot rerun or backfill after success, can rerun or backfill after failure; NONE: cannot rerun or backfill after success or failure.",
						},
						"param_task_out_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Output parameter array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"param_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter name.",
									},
									"param_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter definition.",
									},
								},
							},
						},
						"param_task_in_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Input parameter array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"param_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter name.",
									},
									"param_desc": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter description: format is project_identifier.task_name.parameter_name; example: project_wedata_1.sh_250820_104107.pp_out.",
									},
									"from_task_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parent task ID.",
									},
									"from_param_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parent task parameter key.",
									},
								},
							},
						},
						"task_output_registry_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Output registration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datasource_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Data source ID.",
									},
									"database_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Database name.",
									},
									"table_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Table name.",
									},
									"partition_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Partition name.",
									},
									"data_flow_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Input/output table type\n      Input stream\n UPSTREAM,\n      Output stream\n  DOWNSTREAM.",
									},
									"table_physical_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Table physical unique ID.",
									},
									"db_guid": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Database unique identifier.",
									},
									"table_guid": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Table unique identifier.",
									},
								},
							},
						},
						"init_strategy": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "**Instance generation strategy**\n* T_PLUS_0: T+0 generation, default strategy\n* T_PLUS_1: T+1 generation.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudWedataTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		projectId string
		taskId    string
	)
	var (
		request  = wedatav20250806.NewCreateTaskRequest()
		response = wedatav20250806.NewCreateTaskResponse()
	)

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(projectId)
	}

	if taskBaseAttributeMap, ok := helper.InterfacesHeadMap(d, "task_base_attribute"); ok {
		createTaskBaseAttribute := wedatav20250806.CreateTaskBaseAttribute{}
		if v, ok := taskBaseAttributeMap["task_name"]; ok {
			createTaskBaseAttribute.TaskName = helper.String(v.(string))
		}
		if v, ok := taskBaseAttributeMap["task_type_id"]; ok {
			createTaskBaseAttribute.TaskTypeId = helper.String(v.(string))
		}
		if v, ok := taskBaseAttributeMap["workflow_id"]; ok {
			createTaskBaseAttribute.WorkflowId = helper.String(v.(string))
		}
		if v, ok := taskBaseAttributeMap["owner_uin"]; ok {
			createTaskBaseAttribute.OwnerUin = helper.String(v.(string))
		}
		if v, ok := taskBaseAttributeMap["task_description"]; ok {
			createTaskBaseAttribute.TaskDescription = helper.String(v.(string))
		}
		request.TaskBaseAttribute = &createTaskBaseAttribute
	}

	if taskConfigurationMap, ok := helper.InterfacesHeadMap(d, "task_configuration"); ok {
		createTaskConfiguration := wedatav20250806.CreateTaskConfiguration{}
		if v, ok := taskConfigurationMap["resource_group"]; ok {
			createTaskConfiguration.ResourceGroup = helper.String(v.(string))
		}
		if v, ok := taskConfigurationMap["code_content"]; ok {
			createTaskConfiguration.CodeContent = helper.String(v.(string))
		}
		if v, ok := taskConfigurationMap["task_ext_configuration_list"]; ok {
			for _, item := range v.(*schema.Set).List() {
				taskExtConfigurationListMap := item.(map[string]interface{})
				taskExtParameter := wedatav20250806.TaskExtParameter{}
				if v, ok := taskExtConfigurationListMap["param_key"]; ok {
					taskExtParameter.ParamKey = helper.String(v.(string))
				}
				if v, ok := taskExtConfigurationListMap["param_value"]; ok {
					taskExtParameter.ParamValue = helper.String(v.(string))
				}
				createTaskConfiguration.TaskExtConfigurationList = append(createTaskConfiguration.TaskExtConfigurationList, &taskExtParameter)
			}
		}
		if v, ok := taskConfigurationMap["data_cluster"]; ok {
			createTaskConfiguration.DataCluster = helper.String(v.(string))
		}
		if v, ok := taskConfigurationMap["broker_ip"]; ok {
			createTaskConfiguration.BrokerIp = helper.String(v.(string))
		}
		if v, ok := taskConfigurationMap["yarn_queue"]; ok {
			createTaskConfiguration.YarnQueue = helper.String(v.(string))
		}
		if v, ok := taskConfigurationMap["source_service_id"]; ok {
			createTaskConfiguration.SourceServiceId = helper.String(v.(string))
		}
		if v, ok := taskConfigurationMap["target_service_id"]; ok {
			createTaskConfiguration.TargetServiceId = helper.String(v.(string))
		}
		if v, ok := taskConfigurationMap["task_scheduling_parameter_list"]; ok {
			for _, item := range v.(*schema.Set).List() {
				taskSchedulingParameterListMap := item.(map[string]interface{})
				taskSchedulingParameter := wedatav20250806.TaskSchedulingParameter{}
				if v, ok := taskSchedulingParameterListMap["param_key"]; ok {
					taskSchedulingParameter.ParamKey = helper.String(v.(string))
				}
				if v, ok := taskSchedulingParameterListMap["param_value"]; ok {
					taskSchedulingParameter.ParamValue = helper.String(v.(string))
				}
				createTaskConfiguration.TaskSchedulingParameterList = append(createTaskConfiguration.TaskSchedulingParameterList, &taskSchedulingParameter)
			}
		}
		if v, ok := taskConfigurationMap["bundle_id"]; ok {
			createTaskConfiguration.BundleId = helper.String(v.(string))
		}
		if v, ok := taskConfigurationMap["bundle_info"]; ok {
			createTaskConfiguration.BundleInfo = helper.String(v.(string))
		}
		request.TaskConfiguration = &createTaskConfiguration
	}

	if taskSchedulerConfigurationMap, ok := helper.InterfacesHeadMap(d, "task_scheduler_configuration"); ok {
		createTaskSchedulerConfiguration := wedatav20250806.CreateTaskSchedulerConfiguration{}
		if v, ok := taskSchedulerConfigurationMap["cycle_type"]; ok {
			createTaskSchedulerConfiguration.CycleType = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["schedule_time_zone"]; ok {
			createTaskSchedulerConfiguration.ScheduleTimeZone = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["crontab_expression"]; ok {
			createTaskSchedulerConfiguration.CrontabExpression = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["start_time"]; ok {
			createTaskSchedulerConfiguration.StartTime = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["end_time"]; ok {
			createTaskSchedulerConfiguration.EndTime = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["execution_start_time"]; ok {
			createTaskSchedulerConfiguration.ExecutionStartTime = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["execution_end_time"]; ok {
			createTaskSchedulerConfiguration.ExecutionEndTime = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["schedule_run_type"]; ok {
			createTaskSchedulerConfiguration.ScheduleRunType = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["calendar_open"]; ok {
			createTaskSchedulerConfiguration.CalendarOpen = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["calendar_id"]; ok {
			createTaskSchedulerConfiguration.CalendarId = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["self_depend"]; ok {
			createTaskSchedulerConfiguration.SelfDepend = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["upstream_dependency_config_list"]; ok {
			for _, item := range v.([]interface{}) {
				upstreamDependencyConfigListMap := item.(map[string]interface{})
				dependencyTaskBrief := wedatav20250806.DependencyTaskBrief{}
				if v, ok := upstreamDependencyConfigListMap["task_id"]; ok {
					dependencyTaskBrief.TaskId = helper.String(v.(string))
				}
				if v, ok := upstreamDependencyConfigListMap["main_cyclic_config"]; ok {
					dependencyTaskBrief.MainCyclicConfig = helper.String(v.(string))
				}
				if v, ok := upstreamDependencyConfigListMap["subordinate_cyclic_config"]; ok {
					dependencyTaskBrief.SubordinateCyclicConfig = helper.String(v.(string))
				}
				if v, ok := upstreamDependencyConfigListMap["offset"]; ok {
					dependencyTaskBrief.Offset = helper.String(v.(string))
				}
				if dependencyStrategyMap, ok := helper.ConvertInterfacesHeadToMap(upstreamDependencyConfigListMap["dependency_strategy"]); ok {
					dependencyStrategyTask := wedatav20250806.DependencyStrategyTask{}
					if v, ok := dependencyStrategyMap["polling_null_strategy"]; ok {
						dependencyStrategyTask.PollingNullStrategy = helper.String(v.(string))
					}
					if v, ok := dependencyStrategyMap["task_dependency_executing_strategies"]; ok {
						taskDependencyExecutingStrategiesSet := v.(*schema.Set).List()
						for i := range taskDependencyExecutingStrategiesSet {
							taskDependencyExecutingStrategies := taskDependencyExecutingStrategiesSet[i].(string)
							dependencyStrategyTask.TaskDependencyExecutingStrategies = append(dependencyStrategyTask.TaskDependencyExecutingStrategies, helper.String(taskDependencyExecutingStrategies))
						}
					}
					if v, ok := dependencyStrategyMap["task_dependency_executing_timeout_value"]; ok {
						dependencyStrategyTask.TaskDependencyExecutingTimeoutValue = helper.IntInt64(v.(int))
					}
					dependencyTaskBrief.DependencyStrategy = &dependencyStrategyTask
				}
				createTaskSchedulerConfiguration.UpstreamDependencyConfigList = append(createTaskSchedulerConfiguration.UpstreamDependencyConfigList, &dependencyTaskBrief)
			}
		}
		if v, ok := taskSchedulerConfigurationMap["event_listener_list"]; ok {
			for _, item := range v.([]interface{}) {
				eventListenerListMap := item.(map[string]interface{})
				eventListener := wedatav20250806.EventListener{}
				if v, ok := eventListenerListMap["event_name"]; ok {
					eventListener.EventName = helper.String(v.(string))
				}
				if v, ok := eventListenerListMap["event_sub_type"]; ok {
					eventListener.EventSubType = helper.String(v.(string))
				}
				if v, ok := eventListenerListMap["event_broadcast_type"]; ok {
					eventListener.EventBroadcastType = helper.String(v.(string))
				}
				if v, ok := eventListenerListMap["properties_list"]; ok {
					for _, item := range v.([]interface{}) {
						propertiesListMap := item.(map[string]interface{})
						paramInfo := wedatav20250806.ParamInfo{}
						if v, ok := propertiesListMap["param_key"]; ok {
							paramInfo.ParamKey = helper.String(v.(string))
						}
						if v, ok := propertiesListMap["param_value"]; ok {
							paramInfo.ParamValue = helper.String(v.(string))
						}
						eventListener.PropertiesList = append(eventListener.PropertiesList, &paramInfo)
					}
				}
				createTaskSchedulerConfiguration.EventListenerList = append(createTaskSchedulerConfiguration.EventListenerList, &eventListener)
			}
		}
		if v, ok := taskSchedulerConfigurationMap["run_priority"]; ok {
			createTaskSchedulerConfiguration.RunPriority = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["retry_wait"]; ok {
			createTaskSchedulerConfiguration.RetryWait = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["max_retry_attempts"]; ok {
			createTaskSchedulerConfiguration.MaxRetryAttempts = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["execution_ttl"]; ok {
			createTaskSchedulerConfiguration.ExecutionTTL = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["wait_execution_total_ttl"]; ok {
			createTaskSchedulerConfiguration.WaitExecutionTotalTTL = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["allow_redo_type"]; ok {
			createTaskSchedulerConfiguration.AllowRedoType = helper.String(v.(string))
		}
		if v, ok := taskSchedulerConfigurationMap["param_task_out_list"]; ok {
			for _, item := range v.([]interface{}) {
				paramTaskOutListMap := item.(map[string]interface{})
				outTaskParameter := wedatav20250806.OutTaskParameter{}
				if v, ok := paramTaskOutListMap["param_key"]; ok {
					outTaskParameter.ParamKey = helper.String(v.(string))
				}
				if v, ok := paramTaskOutListMap["param_value"]; ok {
					outTaskParameter.ParamValue = helper.String(v.(string))
				}
				createTaskSchedulerConfiguration.ParamTaskOutList = append(createTaskSchedulerConfiguration.ParamTaskOutList, &outTaskParameter)
			}
		}
		if v, ok := taskSchedulerConfigurationMap["param_task_in_list"]; ok {
			for _, item := range v.([]interface{}) {
				paramTaskInListMap := item.(map[string]interface{})
				inTaskParameter := wedatav20250806.InTaskParameter{}
				if v, ok := paramTaskInListMap["param_key"]; ok {
					inTaskParameter.ParamKey = helper.String(v.(string))
				}
				if v, ok := paramTaskInListMap["param_desc"]; ok {
					inTaskParameter.ParamDesc = helper.String(v.(string))
				}
				if v, ok := paramTaskInListMap["from_task_id"]; ok {
					inTaskParameter.FromTaskId = helper.String(v.(string))
				}
				if v, ok := paramTaskInListMap["from_param_key"]; ok {
					inTaskParameter.FromParamKey = helper.String(v.(string))
				}
				createTaskSchedulerConfiguration.ParamTaskInList = append(createTaskSchedulerConfiguration.ParamTaskInList, &inTaskParameter)
			}
		}
		if v, ok := taskSchedulerConfigurationMap["task_output_registry_list"]; ok {
			for _, item := range v.([]interface{}) {
				taskOutputRegistryListMap := item.(map[string]interface{})
				taskDataRegistry := wedatav20250806.TaskDataRegistry{}
				if v, ok := taskOutputRegistryListMap["datasource_id"]; ok {
					taskDataRegistry.DatasourceId = helper.String(v.(string))
				}
				if v, ok := taskOutputRegistryListMap["database_name"]; ok {
					taskDataRegistry.DatabaseName = helper.String(v.(string))
				}
				if v, ok := taskOutputRegistryListMap["table_name"]; ok {
					taskDataRegistry.TableName = helper.String(v.(string))
				}
				if v, ok := taskOutputRegistryListMap["partition_name"]; ok {
					taskDataRegistry.PartitionName = helper.String(v.(string))
				}
				if v, ok := taskOutputRegistryListMap["data_flow_type"]; ok {
					taskDataRegistry.DataFlowType = helper.String(v.(string))
				}
				if v, ok := taskOutputRegistryListMap["table_physical_id"]; ok {
					taskDataRegistry.TablePhysicalId = helper.String(v.(string))
				}
				if v, ok := taskOutputRegistryListMap["db_guid"]; ok {
					taskDataRegistry.DbGuid = helper.String(v.(string))
				}
				if v, ok := taskOutputRegistryListMap["table_guid"]; ok {
					taskDataRegistry.TableGuid = helper.String(v.(string))
				}
				createTaskSchedulerConfiguration.TaskOutputRegistryList = append(createTaskSchedulerConfiguration.TaskOutputRegistryList, &taskDataRegistry)
			}
		}
		if v, ok := taskSchedulerConfigurationMap["init_strategy"]; ok {
			createTaskSchedulerConfiguration.InitStrategy = helper.String(v.(string))
		}
		request.TaskSchedulerConfiguration = &createTaskSchedulerConfiguration
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata task failed, reason:%+v", logId, err)
		return err
	}

	if response.Response != nil && response.Response.Data != nil && response.Response.Data.TaskId != nil {
		taskId = *response.Response.Data.TaskId
	} else {
		return fmt.Errorf("taskId is nil")
	}
	_ = response

	d.SetId(strings.Join([]string{projectId, taskId}, tccommon.FILED_SP))

	return resourceTencentCloudWedataTaskRead(d, meta)
}

func resourceTencentCloudWedataTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	taskId := idSplit[1]

	var response *wedatav20250806.GetTaskResponseParams
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataTaskById(ctx, projectId, taskId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete wedata task failed, reason:%+v", logId, err)
		return err
	}

	if response == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `wedata_task` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)

	respData := response.Data
	taskBaseAttributeMap := map[string]interface{}{}
	if respData.TaskBaseAttribute != nil {

		if respData.TaskBaseAttribute.TaskName != nil {
			taskBaseAttributeMap["task_name"] = respData.TaskBaseAttribute.TaskName
		}

		if respData.TaskBaseAttribute.TaskTypeId != nil {
			taskBaseAttributeMap["task_type_id"] = helper.UInt64ToStr(*respData.TaskBaseAttribute.TaskTypeId)
		}

		if respData.TaskBaseAttribute.WorkflowId != nil {
			taskBaseAttributeMap["workflow_id"] = respData.TaskBaseAttribute.WorkflowId
		}

		if respData.TaskBaseAttribute.OwnerUin != nil {
			taskBaseAttributeMap["owner_uin"] = respData.TaskBaseAttribute.OwnerUin
		}

		if respData.TaskBaseAttribute.TaskDescription != nil {
			taskBaseAttributeMap["task_description"] = respData.TaskBaseAttribute.TaskDescription
		}

		_ = d.Set("task_base_attribute", []interface{}{taskBaseAttributeMap})
	}

	taskConfigurationMap := map[string]interface{}{}

	if respData.TaskConfiguration != nil {
		if respData.TaskConfiguration.ResourceGroup != nil {
			taskConfigurationMap["resource_group"] = respData.TaskConfiguration.ResourceGroup
		}

		if respData.TaskConfiguration.CodeContent != nil {
			taskConfigurationMap["code_content"] = respData.TaskConfiguration.CodeContent
		}

		taskExtConfigurationListList := make([]map[string]interface{}, 0, len(respData.TaskConfiguration.TaskExtConfigurationList))
		if respData.TaskConfiguration.TaskExtConfigurationList != nil {
			for _, taskExtConfigurationList := range respData.TaskConfiguration.TaskExtConfigurationList {
				taskExtConfigurationListMap := map[string]interface{}{}

				if taskExtConfigurationList.ParamKey != nil {
					taskExtConfigurationListMap["param_key"] = taskExtConfigurationList.ParamKey
				}

				if taskExtConfigurationList.ParamValue != nil {
					taskExtConfigurationListMap["param_value"] = taskExtConfigurationList.ParamValue
				}

				taskExtConfigurationListList = append(taskExtConfigurationListList, taskExtConfigurationListMap)
			}

			taskConfigurationMap["task_ext_configuration_list"] = taskExtConfigurationListList
		}
		if respData.TaskConfiguration.DataCluster != nil {
			taskConfigurationMap["data_cluster"] = respData.TaskConfiguration.DataCluster
		}

		if respData.TaskConfiguration.BrokerIp != nil {
			taskConfigurationMap["broker_ip"] = respData.TaskConfiguration.BrokerIp
		}

		if respData.TaskConfiguration.YarnQueue != nil {
			taskConfigurationMap["yarn_queue"] = respData.TaskConfiguration.YarnQueue
		}

		if respData.TaskConfiguration.SourceServiceId != nil {
			taskConfigurationMap["source_service_id"] = respData.TaskConfiguration.SourceServiceId
		}

		if respData.TaskConfiguration.TargetServiceId != nil {
			taskConfigurationMap["target_service_id"] = respData.TaskConfiguration.TargetServiceId
		}

		taskSchedulingParameterListList := make([]map[string]interface{}, 0, len(respData.TaskConfiguration.TaskSchedulingParameterList))
		if respData.TaskConfiguration.TaskSchedulingParameterList != nil {
			for _, taskSchedulingParameterList := range respData.TaskConfiguration.TaskSchedulingParameterList {
				taskSchedulingParameterListMap := map[string]interface{}{}

				if taskSchedulingParameterList.ParamKey != nil {
					taskSchedulingParameterListMap["param_key"] = taskSchedulingParameterList.ParamKey
				}

				if taskSchedulingParameterList.ParamValue != nil {
					taskSchedulingParameterListMap["param_value"] = taskSchedulingParameterList.ParamValue
				}

				taskSchedulingParameterListList = append(taskSchedulingParameterListList, taskSchedulingParameterListMap)
			}

			taskConfigurationMap["task_scheduling_parameter_list"] = taskSchedulingParameterListList
		}
		if respData.TaskConfiguration.BundleId != nil {
			taskConfigurationMap["bundle_id"] = respData.TaskConfiguration.BundleId
		}

		if respData.TaskConfiguration.BundleInfo != nil {
			taskConfigurationMap["bundle_info"] = respData.TaskConfiguration.BundleInfo
		}

		_ = d.Set("task_configuration", []interface{}{taskConfigurationMap})
	}

	taskSchedulerConfigurationMap := map[string]interface{}{}

	if respData.TaskSchedulerConfiguration != nil {
		if respData.TaskSchedulerConfiguration.CycleType != nil {
			taskSchedulerConfigurationMap["cycle_type"] = respData.TaskSchedulerConfiguration.CycleType
		}

		if respData.TaskSchedulerConfiguration.ScheduleTimeZone != nil {
			taskSchedulerConfigurationMap["schedule_time_zone"] = respData.TaskSchedulerConfiguration.ScheduleTimeZone
		}

		if respData.TaskSchedulerConfiguration.CrontabExpression != nil {
			taskSchedulerConfigurationMap["crontab_expression"] = respData.TaskSchedulerConfiguration.CrontabExpression
		}

		if respData.TaskSchedulerConfiguration.StartTime != nil {
			taskSchedulerConfigurationMap["start_time"] = respData.TaskSchedulerConfiguration.StartTime
		}

		if respData.TaskSchedulerConfiguration.EndTime != nil {
			taskSchedulerConfigurationMap["end_time"] = respData.TaskSchedulerConfiguration.EndTime
		}

		if respData.TaskSchedulerConfiguration.ExecutionStartTime != nil {
			taskSchedulerConfigurationMap["execution_start_time"] = respData.TaskSchedulerConfiguration.ExecutionStartTime
		}

		if respData.TaskSchedulerConfiguration.ExecutionEndTime != nil {
			taskSchedulerConfigurationMap["execution_end_time"] = respData.TaskSchedulerConfiguration.ExecutionEndTime
		}

		if respData.TaskSchedulerConfiguration.ScheduleRunType != nil {
			taskSchedulerConfigurationMap["schedule_run_type"] = helper.Int64ToStr(*respData.TaskSchedulerConfiguration.ScheduleRunType)
		}

		if respData.TaskSchedulerConfiguration.CalendarOpen != nil {
			taskSchedulerConfigurationMap["calendar_open"] = respData.TaskSchedulerConfiguration.CalendarOpen
		}

		if respData.TaskSchedulerConfiguration.CalendarId != nil {
			taskSchedulerConfigurationMap["calendar_id"] = respData.TaskSchedulerConfiguration.CalendarId
		}

		if respData.TaskSchedulerConfiguration.SelfDepend != nil {
			taskSchedulerConfigurationMap["self_depend"] = respData.TaskSchedulerConfiguration.SelfDepend
		}

		upstreamDependencyConfigListList := make([]map[string]interface{}, 0, len(respData.TaskSchedulerConfiguration.UpstreamDependencyConfigList))
		if respData.TaskSchedulerConfiguration.UpstreamDependencyConfigList != nil {
			for _, upstreamDependencyConfigList := range respData.TaskSchedulerConfiguration.UpstreamDependencyConfigList {
				upstreamDependencyConfigListMap := map[string]interface{}{}

				if upstreamDependencyConfigList.TaskId != nil {
					upstreamDependencyConfigListMap["task_id"] = upstreamDependencyConfigList.TaskId
				}

				if upstreamDependencyConfigList.MainCyclicConfig != nil {
					upstreamDependencyConfigListMap["main_cyclic_config"] = upstreamDependencyConfigList.MainCyclicConfig
				}

				if upstreamDependencyConfigList.SubordinateCyclicConfig != nil {
					upstreamDependencyConfigListMap["subordinate_cyclic_config"] = upstreamDependencyConfigList.SubordinateCyclicConfig
				}

				if upstreamDependencyConfigList.Offset != nil {
					upstreamDependencyConfigListMap["offset"] = upstreamDependencyConfigList.Offset
				}

				dependencyStrategyMap := map[string]interface{}{}

				if upstreamDependencyConfigList.DependencyStrategy != nil {
					if upstreamDependencyConfigList.DependencyStrategy.PollingNullStrategy != nil {
						dependencyStrategyMap["polling_null_strategy"] = upstreamDependencyConfigList.DependencyStrategy.PollingNullStrategy
					}

					if upstreamDependencyConfigList.DependencyStrategy.TaskDependencyExecutingStrategies != nil {
						dependencyStrategyMap["task_dependency_executing_strategies"] = upstreamDependencyConfigList.DependencyStrategy.TaskDependencyExecutingStrategies
					}

					if upstreamDependencyConfigList.DependencyStrategy.TaskDependencyExecutingTimeoutValue != nil {
						dependencyStrategyMap["task_dependency_executing_timeout_value"] = upstreamDependencyConfigList.DependencyStrategy.TaskDependencyExecutingTimeoutValue
					}

					upstreamDependencyConfigListMap["dependency_strategy"] = []interface{}{dependencyStrategyMap}
				}

				upstreamDependencyConfigListList = append(upstreamDependencyConfigListList, upstreamDependencyConfigListMap)
			}

			taskSchedulerConfigurationMap["upstream_dependency_config_list"] = upstreamDependencyConfigListList
		}

		eventListenerListList := make([]map[string]interface{}, 0, len(respData.TaskSchedulerConfiguration.EventListenerList))
		if respData.TaskSchedulerConfiguration.EventListenerList != nil {
			for _, eventListenerList := range respData.TaskSchedulerConfiguration.EventListenerList {
				eventListenerListMap := map[string]interface{}{}

				if eventListenerList.EventName != nil {
					eventListenerListMap["event_name"] = eventListenerList.EventName
				}

				if eventListenerList.EventSubType != nil {
					eventListenerListMap["event_sub_type"] = eventListenerList.EventSubType
				}

				if eventListenerList.EventBroadcastType != nil {
					eventListenerListMap["event_broadcast_type"] = eventListenerList.EventBroadcastType
				}

				propertiesListList := make([]map[string]interface{}, 0, len(eventListenerList.PropertiesList))
				if eventListenerList.PropertiesList != nil {
					for _, propertiesList := range eventListenerList.PropertiesList {
						propertiesListMap := map[string]interface{}{}

						if propertiesList.ParamKey != nil {
							propertiesListMap["param_key"] = propertiesList.ParamKey
						}

						if propertiesList.ParamValue != nil {
							propertiesListMap["param_value"] = propertiesList.ParamValue
						}

						propertiesListList = append(propertiesListList, propertiesListMap)
					}

					eventListenerListMap["properties_list"] = propertiesListList
				}
				eventListenerListList = append(eventListenerListList, eventListenerListMap)
			}

			taskSchedulerConfigurationMap["event_listener_list"] = eventListenerListList
		}
		if respData.TaskSchedulerConfiguration.RunPriority != nil {
			taskSchedulerConfigurationMap["run_priority"] = helper.UInt64ToStr(*respData.TaskSchedulerConfiguration.RunPriority)
		}

		if respData.TaskSchedulerConfiguration.RetryWait != nil {
			taskSchedulerConfigurationMap["retry_wait"] = helper.Int64ToStr(*respData.TaskSchedulerConfiguration.RetryWait)
		}

		if respData.TaskSchedulerConfiguration.MaxRetryAttempts != nil {
			taskSchedulerConfigurationMap["max_retry_attempts"] = helper.Int64ToStr(*respData.TaskSchedulerConfiguration.MaxRetryAttempts)
		}

		if respData.TaskSchedulerConfiguration.ExecutionTTL != nil {
			taskSchedulerConfigurationMap["execution_ttl"] = helper.Int64ToStr(*respData.TaskSchedulerConfiguration.ExecutionTTL)
		}

		if respData.TaskSchedulerConfiguration.WaitExecutionTotalTTL != nil {
			taskSchedulerConfigurationMap["wait_execution_total_ttl"] = respData.TaskSchedulerConfiguration.WaitExecutionTotalTTL
		}

		if respData.TaskSchedulerConfiguration.AllowRedoType != nil {
			taskSchedulerConfigurationMap["allow_redo_type"] = respData.TaskSchedulerConfiguration.AllowRedoType
		}

		paramTaskOutListList := make([]map[string]interface{}, 0, len(respData.TaskSchedulerConfiguration.ParamTaskOutList))
		if respData.TaskSchedulerConfiguration.ParamTaskOutList != nil {
			for _, paramTaskOutList := range respData.TaskSchedulerConfiguration.ParamTaskOutList {
				paramTaskOutListMap := map[string]interface{}{}

				if paramTaskOutList.ParamKey != nil {
					paramTaskOutListMap["param_key"] = paramTaskOutList.ParamKey
				}

				if paramTaskOutList.ParamValue != nil {
					paramTaskOutListMap["param_value"] = paramTaskOutList.ParamValue
				}

				paramTaskOutListList = append(paramTaskOutListList, paramTaskOutListMap)
			}

			taskSchedulerConfigurationMap["param_task_out_list"] = paramTaskOutListList
		}
		paramTaskInListList := make([]map[string]interface{}, 0, len(respData.TaskSchedulerConfiguration.ParamTaskInList))
		if respData.TaskSchedulerConfiguration.ParamTaskInList != nil {
			for _, paramTaskInList := range respData.TaskSchedulerConfiguration.ParamTaskInList {
				paramTaskInListMap := map[string]interface{}{}

				if paramTaskInList.ParamKey != nil {
					paramTaskInListMap["param_key"] = paramTaskInList.ParamKey
				}

				if paramTaskInList.ParamDesc != nil {
					paramTaskInListMap["param_desc"] = paramTaskInList.ParamDesc
				}

				if paramTaskInList.FromTaskId != nil {
					paramTaskInListMap["from_task_id"] = paramTaskInList.FromTaskId
				}

				if paramTaskInList.FromParamKey != nil {
					paramTaskInListMap["from_param_key"] = paramTaskInList.FromParamKey
				}

				paramTaskInListList = append(paramTaskInListList, paramTaskInListMap)
			}

			taskSchedulerConfigurationMap["param_task_in_list"] = paramTaskInListList
		}
		taskOutputRegistryListList := make([]map[string]interface{}, 0, len(respData.TaskSchedulerConfiguration.TaskOutputRegistryList))
		if respData.TaskSchedulerConfiguration.TaskOutputRegistryList != nil {
			for _, taskOutputRegistryList := range respData.TaskSchedulerConfiguration.TaskOutputRegistryList {
				taskOutputRegistryListMap := map[string]interface{}{}

				if taskOutputRegistryList.DatasourceId != nil {
					taskOutputRegistryListMap["datasource_id"] = taskOutputRegistryList.DatasourceId
				}

				if taskOutputRegistryList.DatabaseName != nil {
					taskOutputRegistryListMap["database_name"] = taskOutputRegistryList.DatabaseName
				}

				if taskOutputRegistryList.TableName != nil {
					taskOutputRegistryListMap["table_name"] = taskOutputRegistryList.TableName
				}

				if taskOutputRegistryList.PartitionName != nil {
					taskOutputRegistryListMap["partition_name"] = taskOutputRegistryList.PartitionName
				}

				if taskOutputRegistryList.DataFlowType != nil {
					taskOutputRegistryListMap["data_flow_type"] = taskOutputRegistryList.DataFlowType
				}

				if taskOutputRegistryList.TablePhysicalId != nil {
					taskOutputRegistryListMap["table_physical_id"] = taskOutputRegistryList.TablePhysicalId
				}

				if taskOutputRegistryList.DbGuid != nil {
					taskOutputRegistryListMap["db_guid"] = taskOutputRegistryList.DbGuid
				}

				if taskOutputRegistryList.TableGuid != nil {
					taskOutputRegistryListMap["table_guid"] = taskOutputRegistryList.TableGuid
				}

				taskOutputRegistryListList = append(taskOutputRegistryListList, taskOutputRegistryListMap)
			}

			taskSchedulerConfigurationMap["task_output_registry_list"] = taskOutputRegistryListList
		}
		if respData.TaskSchedulerConfiguration.InitStrategy != nil {
			taskSchedulerConfigurationMap["init_strategy"] = respData.TaskSchedulerConfiguration.InitStrategy
		}

		_ = d.Set("task_scheduler_configuration", []interface{}{taskSchedulerConfigurationMap})
	}
	_ = projectId
	_ = taskId
	return nil
}

func resourceTencentCloudWedataTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_task.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	taskId := idSplit[1]

	needChange := false
	mutableArgs := []string{"task_base_attribute", "task_configuration", "task_scheduler_configuration"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wedatav20250806.NewUpdateTaskRequest()
		request.ProjectId = helper.String(projectId)
		request.TaskId = helper.String(taskId)
		request.Task = &wedatav20250806.UpdateTaskBrief{}
		if taskBaseAttributeMap, ok := helper.InterfacesHeadMap(d, "task_base_attribute"); ok {
			updateTaskBaseAttribute := wedatav20250806.UpdateTaskBaseAttribute{}
			if v, ok := taskBaseAttributeMap["task_name"]; ok {
				updateTaskBaseAttribute.TaskName = helper.String(v.(string))
			}
			if v, ok := taskBaseAttributeMap["owner_uin"]; ok {
				updateTaskBaseAttribute.OwnerUin = helper.String(v.(string))
			}
			if v, ok := taskBaseAttributeMap["task_description"]; ok {
				updateTaskBaseAttribute.TaskDescription = helper.String(v.(string))
			}
			request.Task.TaskBaseAttribute = &updateTaskBaseAttribute
		}

		if taskConfigurationMap, ok := helper.InterfacesHeadMap(d, "task_configuration"); ok {
			taskConfiguration := wedatav20250806.TaskConfiguration{}
			if v, ok := taskConfigurationMap["resource_group"]; ok {
				taskConfiguration.ResourceGroup = helper.String(v.(string))
			}
			if v, ok := taskConfigurationMap["code_content"]; ok {
				taskConfiguration.CodeContent = helper.String(v.(string))
			}
			if v, ok := taskConfigurationMap["task_ext_configuration_list"]; ok {
				for _, item := range v.(*schema.Set).List() {
					taskExtConfigurationListMap := item.(map[string]interface{})
					taskExtParameter := wedatav20250806.TaskExtParameter{}
					if v, ok := taskExtConfigurationListMap["param_key"]; ok {
						taskExtParameter.ParamKey = helper.String(v.(string))
					}
					if v, ok := taskExtConfigurationListMap["param_value"]; ok {
						taskExtParameter.ParamValue = helper.String(v.(string))
					}
					taskConfiguration.TaskExtConfigurationList = append(taskConfiguration.TaskExtConfigurationList, &taskExtParameter)
				}
			}
			if v, ok := taskConfigurationMap["data_cluster"]; ok {
				taskConfiguration.DataCluster = helper.String(v.(string))
			}
			if v, ok := taskConfigurationMap["broker_ip"]; ok {
				taskConfiguration.BrokerIp = helper.String(v.(string))
			}
			if v, ok := taskConfigurationMap["yarn_queue"]; ok {
				taskConfiguration.YarnQueue = helper.String(v.(string))
			}
			if v, ok := taskConfigurationMap["source_service_id"]; ok {
				taskConfiguration.SourceServiceId = helper.String(v.(string))
			}
			if v, ok := taskConfigurationMap["target_service_id"]; ok {
				taskConfiguration.TargetServiceId = helper.String(v.(string))
			}
			if v, ok := taskConfigurationMap["task_scheduling_parameter_list"]; ok {
				for _, item := range v.(*schema.Set).List() {
					taskSchedulingParameterListMap := item.(map[string]interface{})
					taskSchedulingParameter := wedatav20250806.TaskSchedulingParameter{}
					if v, ok := taskSchedulingParameterListMap["param_key"]; ok {
						taskSchedulingParameter.ParamKey = helper.String(v.(string))
					}
					if v, ok := taskSchedulingParameterListMap["param_value"]; ok {
						taskSchedulingParameter.ParamValue = helper.String(v.(string))
					}
					taskConfiguration.TaskSchedulingParameterList = append(taskConfiguration.TaskSchedulingParameterList, &taskSchedulingParameter)
				}
			}
			if v, ok := taskConfigurationMap["bundle_id"]; ok {
				taskConfiguration.BundleId = helper.String(v.(string))
			}
			if v, ok := taskConfigurationMap["bundle_info"]; ok {
				taskConfiguration.BundleInfo = helper.String(v.(string))
			}
			request.Task.TaskConfiguration = &taskConfiguration
		}

		if taskSchedulerConfigurationMap, ok := helper.InterfacesHeadMap(d, "task_scheduler_configuration"); ok {
			taskSchedulerConfiguration := wedatav20250806.TaskSchedulerConfiguration{}
			if v, ok := taskSchedulerConfigurationMap["cycle_type"]; ok {
				taskSchedulerConfiguration.CycleType = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["schedule_time_zone"]; ok {
				taskSchedulerConfiguration.ScheduleTimeZone = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["crontab_expression"]; ok {
				taskSchedulerConfiguration.CrontabExpression = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["start_time"]; ok {
				taskSchedulerConfiguration.StartTime = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["end_time"]; ok {
				taskSchedulerConfiguration.EndTime = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["execution_start_time"]; ok {
				taskSchedulerConfiguration.ExecutionStartTime = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["execution_end_time"]; ok {
				taskSchedulerConfiguration.ExecutionEndTime = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["schedule_run_type"]; ok {
				taskSchedulerConfiguration.ScheduleRunType = helper.StrToInt64Point(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["calendar_open"]; ok {
				taskSchedulerConfiguration.CalendarOpen = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["calendar_id"]; ok {
				taskSchedulerConfiguration.CalendarId = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["self_depend"]; ok {
				taskSchedulerConfiguration.SelfDepend = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["upstream_dependency_config_list"]; ok {
				for _, item := range v.([]interface{}) {
					upstreamDependencyConfigListMap := item.(map[string]interface{})
					dependencyTaskBrief := wedatav20250806.DependencyTaskBrief{}
					if v, ok := upstreamDependencyConfigListMap["task_id"]; ok {
						dependencyTaskBrief.TaskId = helper.String(v.(string))
					}
					if v, ok := upstreamDependencyConfigListMap["main_cyclic_config"]; ok {
						dependencyTaskBrief.MainCyclicConfig = helper.String(v.(string))
					}
					if v, ok := upstreamDependencyConfigListMap["subordinate_cyclic_config"]; ok {
						dependencyTaskBrief.SubordinateCyclicConfig = helper.String(v.(string))
					}
					if v, ok := upstreamDependencyConfigListMap["offset"]; ok {
						dependencyTaskBrief.Offset = helper.String(v.(string))
					}
					if dependencyStrategyMap, ok := helper.ConvertInterfacesHeadToMap(upstreamDependencyConfigListMap["dependency_strategy"]); ok {
						dependencyStrategyTask := wedatav20250806.DependencyStrategyTask{}
						if v, ok := dependencyStrategyMap["polling_null_strategy"]; ok {
							dependencyStrategyTask.PollingNullStrategy = helper.String(v.(string))
						}
						if v, ok := dependencyStrategyMap["task_dependency_executing_strategies"]; ok {
							taskDependencyExecutingStrategiesSet := v.(*schema.Set).List()
							for i := range taskDependencyExecutingStrategiesSet {
								taskDependencyExecutingStrategies := taskDependencyExecutingStrategiesSet[i].(string)
								dependencyStrategyTask.TaskDependencyExecutingStrategies = append(dependencyStrategyTask.TaskDependencyExecutingStrategies, helper.String(taskDependencyExecutingStrategies))
							}
						}
						if v, ok := dependencyStrategyMap["task_dependency_executing_timeout_value"]; ok {
							dependencyStrategyTask.TaskDependencyExecutingTimeoutValue = helper.IntInt64(v.(int))
						}
						dependencyTaskBrief.DependencyStrategy = &dependencyStrategyTask
					}
					taskSchedulerConfiguration.UpstreamDependencyConfigList = append(taskSchedulerConfiguration.UpstreamDependencyConfigList, &dependencyTaskBrief)
				}
			}
			if v, ok := taskSchedulerConfigurationMap["event_listener_list"]; ok {
				for _, item := range v.([]interface{}) {
					eventListenerListMap := item.(map[string]interface{})
					eventListener := wedatav20250806.EventListener{}
					if v, ok := eventListenerListMap["event_name"]; ok {
						eventListener.EventName = helper.String(v.(string))
					}
					if v, ok := eventListenerListMap["event_sub_type"]; ok {
						eventListener.EventSubType = helper.String(v.(string))
					}
					if v, ok := eventListenerListMap["event_broadcast_type"]; ok {
						eventListener.EventBroadcastType = helper.String(v.(string))
					}
					if v, ok := eventListenerListMap["properties_list"]; ok {
						for _, item := range v.([]interface{}) {
							propertiesListMap := item.(map[string]interface{})
							paramInfo := wedatav20250806.ParamInfo{}
							if v, ok := propertiesListMap["param_key"]; ok {
								paramInfo.ParamKey = helper.String(v.(string))
							}
							if v, ok := propertiesListMap["param_value"]; ok {
								paramInfo.ParamValue = helper.String(v.(string))
							}
							eventListener.PropertiesList = append(eventListener.PropertiesList, &paramInfo)
						}
					}
					taskSchedulerConfiguration.EventListenerList = append(taskSchedulerConfiguration.EventListenerList, &eventListener)
				}
			}
			if v, ok := taskSchedulerConfigurationMap["run_priority"]; ok {
				taskSchedulerConfiguration.RunPriority = helper.StrToUint64Point(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["retry_wait"]; ok {
				taskSchedulerConfiguration.RetryWait = helper.StrToInt64Point(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["max_retry_attempts"]; ok {
				taskSchedulerConfiguration.MaxRetryAttempts = helper.StrToInt64Point(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["execution_ttl"]; ok {
				taskSchedulerConfiguration.ExecutionTTL = helper.StrToInt64Point(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["wait_execution_total_ttl"]; ok {
				taskSchedulerConfiguration.WaitExecutionTotalTTL = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["allow_redo_type"]; ok {
				taskSchedulerConfiguration.AllowRedoType = helper.String(v.(string))
			}
			if v, ok := taskSchedulerConfigurationMap["param_task_out_list"]; ok {
				for _, item := range v.([]interface{}) {
					paramTaskOutListMap := item.(map[string]interface{})
					outTaskParameter := wedatav20250806.OutTaskParameter{}
					if v, ok := paramTaskOutListMap["param_key"]; ok {
						outTaskParameter.ParamKey = helper.String(v.(string))
					}
					if v, ok := paramTaskOutListMap["param_value"]; ok {
						outTaskParameter.ParamValue = helper.String(v.(string))
					}
					taskSchedulerConfiguration.ParamTaskOutList = append(taskSchedulerConfiguration.ParamTaskOutList, &outTaskParameter)
				}
			}
			if v, ok := taskSchedulerConfigurationMap["param_task_in_list"]; ok {
				for _, item := range v.([]interface{}) {
					paramTaskInListMap := item.(map[string]interface{})
					inTaskParameter := wedatav20250806.InTaskParameter{}
					if v, ok := paramTaskInListMap["param_key"]; ok {
						inTaskParameter.ParamKey = helper.String(v.(string))
					}
					if v, ok := paramTaskInListMap["param_desc"]; ok {
						inTaskParameter.ParamDesc = helper.String(v.(string))
					}
					if v, ok := paramTaskInListMap["from_task_id"]; ok {
						inTaskParameter.FromTaskId = helper.String(v.(string))
					}
					if v, ok := paramTaskInListMap["from_param_key"]; ok {
						inTaskParameter.FromParamKey = helper.String(v.(string))
					}
					taskSchedulerConfiguration.ParamTaskInList = append(taskSchedulerConfiguration.ParamTaskInList, &inTaskParameter)
				}
			}
			if v, ok := taskSchedulerConfigurationMap["task_output_registry_list"]; ok {
				for _, item := range v.([]interface{}) {
					taskOutputRegistryListMap := item.(map[string]interface{})
					taskDataRegistry := wedatav20250806.TaskDataRegistry{}
					if v, ok := taskOutputRegistryListMap["datasource_id"]; ok {
						taskDataRegistry.DatasourceId = helper.String(v.(string))
					}
					if v, ok := taskOutputRegistryListMap["database_name"]; ok {
						taskDataRegistry.DatabaseName = helper.String(v.(string))
					}
					if v, ok := taskOutputRegistryListMap["table_name"]; ok {
						taskDataRegistry.TableName = helper.String(v.(string))
					}
					if v, ok := taskOutputRegistryListMap["partition_name"]; ok {
						taskDataRegistry.PartitionName = helper.String(v.(string))
					}
					if v, ok := taskOutputRegistryListMap["data_flow_type"]; ok {
						taskDataRegistry.DataFlowType = helper.String(v.(string))
					}
					if v, ok := taskOutputRegistryListMap["table_physical_id"]; ok {
						taskDataRegistry.TablePhysicalId = helper.String(v.(string))
					}
					if v, ok := taskOutputRegistryListMap["db_guid"]; ok {
						taskDataRegistry.DbGuid = helper.String(v.(string))
					}
					if v, ok := taskOutputRegistryListMap["table_guid"]; ok {
						taskDataRegistry.TableGuid = helper.String(v.(string))
					}
					taskSchedulerConfiguration.TaskOutputRegistryList = append(taskSchedulerConfiguration.TaskOutputRegistryList, &taskDataRegistry)
				}
			}
			if v, ok := taskSchedulerConfigurationMap["init_strategy"]; ok {
				taskSchedulerConfiguration.InitStrategy = helper.String(v.(string))
			}
			request.Task.TaskSchedulerConfiguration = &taskSchedulerConfiguration
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateTaskWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update wedata task failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = projectId
	_ = taskId
	return resourceTencentCloudWedataTaskRead(d, meta)
}

func resourceTencentCloudWedataTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	taskId := idSplit[1]

	var (
		request  = wedatav20250806.NewDeleteTaskRequest()
		response = wedatav20250806.NewDeleteTaskResponse()
	)

	request.ProjectId = helper.String(projectId)
	request.TaskId = helper.String(taskId)

	if v, ok := d.GetOkExists("operate_inform"); ok {
		request.OperateInform = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("delete_mode"); ok {
		request.DeleteMode = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete wedata task failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = projectId
	_ = taskId
	return nil
}

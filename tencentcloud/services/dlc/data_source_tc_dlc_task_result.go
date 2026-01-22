package dlc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcTaskResult() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcTaskResultRead,
		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique task ID.",
			},

			"next_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The pagination information returned by the last response. This parameter can be omitted for the first response, where the data will be returned from the beginning. The data with a volume set by the `MaxResults` field is returned each time.",
			},

			"max_results": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of returned rows. Value range: 0-1,000. Default value: 1,000.",
			},

			"is_transform_data_type": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to convert the data type.",
			},

			"task_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The queried task information. If the returned value is empty, the task with the entered task ID does not exist. The task result will be returned only if the task status is `2` (succeeded).\nNote: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique task ID.",
						},
						"datasource_connection_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the default selected data source when the current job is executed\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"database_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the default selected database when the current job is executed\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"sql": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The currently executed SQL statement. Each task contains one SQL statement.",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of the executed task. Valid values: `DDL`, `DML`, `DQL`.",
						},
						"state": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "u200cThe current task status. Valid values: `0` (initializing), `1` (executing), `2` (executed), `3` (writing data), `4` (queuing), u200c`-1` (failed), and `-3` (canceled). Only when the task is successfully executed, a task execution result will be returned.",
						},
						"data_amount": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Amount of the data scanned in bytes.",
						},
						"used_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The compute time in ms.",
						},
						"output_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Address of the COS bucket for storing the task result.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task creation timestamp.",
						},
						"output_message": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task execution information. `success` will be returned if the task succeeds; otherwise, the failure cause will be returned.",
						},
						"row_affect_info": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Number of affected rows.",
						},
						"result_schema": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Schema information of the result\nNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Column name, which is case-insensitive and can contain up to 25 characters.",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Column type. Valid values:\nstring|tinyint|smallint|int|bigint|boolean|float|double|decimal|timestamp|date|binary|array<data_type>|map<primitive_type, data_type>|struct<col_name : data_type [COMMENT col_comment], ...>|uniontype<data_type, data_type, ...>.",
									},
									"comment": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Class comment.\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"precision": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length of the entire numeric value\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"scale": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length of the decimal part\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"nullable": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether the column is null.\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"position": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Field position\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field creation time\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"modified_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field modification time\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"is_partition": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the column is the partition field.\nNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"result_set": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Result information. After it is unescaped, each element of the outer array is a data row.\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"next_token": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pagination information. If there is no more result data, `nextToken` will be empty.",
						},
						"percentage": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Task progress (%).",
						},
						"progress_detail": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task progress details.",
						},
						"display_format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Console display format. Valid values: `table`, `text`.",
						},
						"total_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The task time in ms.",
						},
						"query_result_time": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "Time consumed to get results\nNote: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudDlcTaskResultRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_task_result.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		taskId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("task_id"); ok {
		paramMap["TaskId"] = helper.String(v.(string))
		taskId = v.(string)
	}

	if v, ok := d.GetOk("next_token"); ok {
		paramMap["NextToken"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_results"); ok {
		paramMap["MaxResults"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_transform_data_type"); ok {
		paramMap["IsTransformDataType"] = helper.Bool(v.(bool))
	}

	var respData *dlcv20210125.DescribeTaskResultResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcTaskResultByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	taskInfoMap := map[string]interface{}{}
	if respData.TaskInfo != nil {
		if respData.TaskInfo.TaskId != nil {
			taskInfoMap["task_id"] = respData.TaskInfo.TaskId
		}

		if respData.TaskInfo.DatasourceConnectionName != nil {
			taskInfoMap["datasource_connection_name"] = respData.TaskInfo.DatasourceConnectionName
		}

		if respData.TaskInfo.DatabaseName != nil {
			taskInfoMap["database_name"] = respData.TaskInfo.DatabaseName
		}

		if respData.TaskInfo.SQL != nil {
			taskInfoMap["sql"] = respData.TaskInfo.SQL
		}

		if respData.TaskInfo.SQLType != nil {
			taskInfoMap["sql_type"] = respData.TaskInfo.SQLType
		}

		if respData.TaskInfo.State != nil {
			taskInfoMap["state"] = respData.TaskInfo.State
		}

		if respData.TaskInfo.DataAmount != nil {
			taskInfoMap["data_amount"] = respData.TaskInfo.DataAmount
		}

		if respData.TaskInfo.UsedTime != nil {
			taskInfoMap["used_time"] = respData.TaskInfo.UsedTime
		}

		if respData.TaskInfo.OutputPath != nil {
			taskInfoMap["output_path"] = respData.TaskInfo.OutputPath
		}

		if respData.TaskInfo.CreateTime != nil {
			taskInfoMap["create_time"] = respData.TaskInfo.CreateTime
		}

		if respData.TaskInfo.OutputMessage != nil {
			taskInfoMap["output_message"] = respData.TaskInfo.OutputMessage
		}

		if respData.TaskInfo.RowAffectInfo != nil {
			taskInfoMap["row_affect_info"] = respData.TaskInfo.RowAffectInfo
		}

		resultSchemaList := make([]map[string]interface{}, 0, len(respData.TaskInfo.ResultSchema))
		if respData.TaskInfo.ResultSchema != nil {
			for _, resultSchema := range respData.TaskInfo.ResultSchema {
				resultSchemaMap := map[string]interface{}{}

				if resultSchema.Name != nil {
					resultSchemaMap["name"] = resultSchema.Name
				}

				if resultSchema.Type != nil {
					resultSchemaMap["type"] = resultSchema.Type
				}

				if resultSchema.Comment != nil {
					resultSchemaMap["comment"] = resultSchema.Comment
				}

				if resultSchema.Precision != nil {
					resultSchemaMap["precision"] = resultSchema.Precision
				}

				if resultSchema.Scale != nil {
					resultSchemaMap["scale"] = resultSchema.Scale
				}

				if resultSchema.Nullable != nil {
					resultSchemaMap["nullable"] = resultSchema.Nullable
				}

				if resultSchema.Position != nil {
					resultSchemaMap["position"] = resultSchema.Position
				}

				if resultSchema.CreateTime != nil {
					resultSchemaMap["create_time"] = resultSchema.CreateTime
				}

				if resultSchema.ModifiedTime != nil {
					resultSchemaMap["modified_time"] = resultSchema.ModifiedTime
				}

				if resultSchema.IsPartition != nil {
					resultSchemaMap["is_partition"] = resultSchema.IsPartition
				}

				resultSchemaList = append(resultSchemaList, resultSchemaMap)
			}

			taskInfoMap["result_schema"] = resultSchemaList
		}
		if respData.TaskInfo.ResultSet != nil {
			taskInfoMap["result_set"] = respData.TaskInfo.ResultSet
		}

		if respData.TaskInfo.NextToken != nil {
			taskInfoMap["next_token"] = respData.TaskInfo.NextToken
		}

		if respData.TaskInfo.Percentage != nil {
			taskInfoMap["percentage"] = respData.TaskInfo.Percentage
		}

		if respData.TaskInfo.ProgressDetail != nil {
			taskInfoMap["progress_detail"] = respData.TaskInfo.ProgressDetail
		}

		if respData.TaskInfo.DisplayFormat != nil {
			taskInfoMap["display_format"] = respData.TaskInfo.DisplayFormat
		}

		if respData.TaskInfo.TotalTime != nil {
			taskInfoMap["total_time"] = respData.TaskInfo.TotalTime
		}

		if respData.TaskInfo.QueryResultTime != nil {
			taskInfoMap["query_result_time"] = respData.TaskInfo.QueryResultTime
		}

		_ = d.Set("task_info", []interface{}{taskInfoMap})
	}

	d.SetId(taskId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), taskInfoMap); e != nil {
			return e
		}
	}

	return nil
}

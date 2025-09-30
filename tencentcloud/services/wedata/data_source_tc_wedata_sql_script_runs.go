package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataSqlScriptRuns() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataSqlScriptRunsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"script_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Script ID.",
			},

			"job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Job ID.",
			},

			"search_word": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search keyword.",
			},

			"execute_user_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Execute user UIN.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Start time.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "End time.",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data exploration tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data exploration task ID.",
						},
						"job_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data exploration task name.",
						},
						"job_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job type.",
						},
						"script_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Script ID.",
						},
						"job_execution_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subtask list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"job_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data exploration task ID.",
									},
									"job_execution_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subquery task ID.",
									},
									"job_execution_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subquery name.",
									},
									"script_content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subquery SQL content.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subquery status.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create time.",
									},
									"execute_stage_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution phase.",
									},
									"log_file_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log file path.",
									},
									"result_file_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Result file path.",
									},
									"result_preview_file_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Preview result file path.",
									},
									"result_total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total number of rows in the task execution result.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update time.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "End time.",
									},
									"time_cost": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Time consumed.",
									},
									"context_script_content": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Context SQL content.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"result_preview_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of rows for previewing the task execution results.",
									},
									"result_effect_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of rows affected by the task execution result.",
									},
									"collecting_total_result": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether collecting full results: default false, true means collecting full results, used for frontend polling.",
									},
									"script_content_truncate": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the script content is truncated.",
									},
								},
							},
						},
						"script_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Script content.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task status.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud owner account UIN.",
						},
						"user_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account UIN.",
						},
						"time_cost": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time consumed.",
						},
						"script_content_truncate": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the script content is truncated.",
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

func dataSourceTencentCloudWedataSqlScriptRunsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_sql_script_runs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("script_id"); ok {
		paramMap["ScriptId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("execute_user_uin"); ok {
		paramMap["ExecuteUserUin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.JobDto
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataSqlScriptRunsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, data := range respData {
			dataMap := map[string]interface{}{}
			if data.JobId != nil {
				dataMap["job_id"] = data.JobId
			}

			if data.JobName != nil {
				dataMap["job_name"] = data.JobName
			}

			if data.JobType != nil {
				dataMap["job_type"] = data.JobType
			}

			if data.ScriptId != nil {
				dataMap["script_id"] = data.ScriptId
			}

			jobExecutionListList := make([]map[string]interface{}, 0, len(data.JobExecutionList))
			if data.JobExecutionList != nil {
				for _, jobExecutionList := range data.JobExecutionList {
					jobExecutionListMap := map[string]interface{}{}

					if jobExecutionList.JobId != nil {
						jobExecutionListMap["job_id"] = jobExecutionList.JobId
					}

					if jobExecutionList.JobExecutionId != nil {
						jobExecutionListMap["job_execution_id"] = jobExecutionList.JobExecutionId
					}

					if jobExecutionList.JobExecutionName != nil {
						jobExecutionListMap["job_execution_name"] = jobExecutionList.JobExecutionName
					}

					if jobExecutionList.ScriptContent != nil {
						jobExecutionListMap["script_content"] = jobExecutionList.ScriptContent
					}

					if jobExecutionList.Status != nil {
						jobExecutionListMap["status"] = jobExecutionList.Status
					}

					if jobExecutionList.CreateTime != nil {
						jobExecutionListMap["create_time"] = jobExecutionList.CreateTime
					}

					if jobExecutionList.ExecuteStageInfo != nil {
						jobExecutionListMap["execute_stage_info"] = jobExecutionList.ExecuteStageInfo
					}

					if jobExecutionList.LogFilePath != nil {
						jobExecutionListMap["log_file_path"] = jobExecutionList.LogFilePath
					}

					if jobExecutionList.ResultFilePath != nil {
						jobExecutionListMap["result_file_path"] = jobExecutionList.ResultFilePath
					}

					if jobExecutionList.ResultPreviewFilePath != nil {
						jobExecutionListMap["result_preview_file_path"] = jobExecutionList.ResultPreviewFilePath
					}

					if jobExecutionList.ResultTotalCount != nil {
						jobExecutionListMap["result_total_count"] = jobExecutionList.ResultTotalCount
					}

					if jobExecutionList.UpdateTime != nil {
						jobExecutionListMap["update_time"] = jobExecutionList.UpdateTime
					}

					if jobExecutionList.EndTime != nil {
						jobExecutionListMap["end_time"] = jobExecutionList.EndTime
					}

					if jobExecutionList.TimeCost != nil {
						jobExecutionListMap["time_cost"] = jobExecutionList.TimeCost
					}

					if jobExecutionList.ContextScriptContent != nil {
						jobExecutionListMap["context_script_content"] = jobExecutionList.ContextScriptContent
					}

					if jobExecutionList.ResultPreviewCount != nil {
						jobExecutionListMap["result_preview_count"] = jobExecutionList.ResultPreviewCount
					}

					if jobExecutionList.ResultEffectCount != nil {
						jobExecutionListMap["result_effect_count"] = jobExecutionList.ResultEffectCount
					}

					if jobExecutionList.CollectingTotalResult != nil {
						jobExecutionListMap["collecting_total_result"] = jobExecutionList.CollectingTotalResult
					}

					if jobExecutionList.ScriptContentTruncate != nil {
						jobExecutionListMap["script_content_truncate"] = jobExecutionList.ScriptContentTruncate
					}

					jobExecutionListList = append(jobExecutionListList, jobExecutionListMap)
				}

				dataMap["job_execution_list"] = jobExecutionListList
			}
			if data.ScriptContent != nil {
				dataMap["script_content"] = data.ScriptContent
			}

			if data.Status != nil {
				dataMap["status"] = data.Status
			}

			if data.CreateTime != nil {
				dataMap["create_time"] = data.CreateTime
			}

			if data.UpdateTime != nil {
				dataMap["update_time"] = data.UpdateTime
			}

			if data.EndTime != nil {
				dataMap["end_time"] = data.EndTime
			}

			if data.OwnerUin != nil {
				dataMap["owner_uin"] = data.OwnerUin
			}

			if data.UserUin != nil {
				dataMap["user_uin"] = data.UserUin
			}

			if data.TimeCost != nil {
				dataMap["time_cost"] = data.TimeCost
			}

			if data.ScriptContentTruncate != nil {
				dataMap["script_content_truncate"] = data.ScriptContentTruncate
			}

			dataList = append(dataList, dataMap)
		}

		_ = d.Set("data", dataList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataList); e != nil {
			return e
		}
	}

	return nil
}

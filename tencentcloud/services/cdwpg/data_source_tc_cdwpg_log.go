package cdwpg

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwpgv20201230 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdwpgLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdwpgLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time.",
			},

			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Database.",
			},

			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort by.",
			},

			"order_by_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Ascending/Descending.",
			},

			"duration": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Filter duration.",
			},

			"slow_log_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Slow sql log details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_time": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "Total time spent.",
						},
						"total_call_times": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Total call count.",
						},
						"normal_querys": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Slow sql.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"call_times": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Call count.",
									},
									"shared_read_blocks": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Shared read blocks.",
									},
									"shared_write_blocks": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Shared write blocks.",
									},
									"database_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Database.",
									},
									"normal_query": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Desensitized query.",
									},
									"max_elapsed_query": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Longest execution time query.",
									},
									"cost_time": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Total time spent.",
									},
									"client_ip": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Client ip.",
									},
									"user_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Username.",
									},
									"total_call_times_percent": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Total call count percentage.",
									},
									"total_cost_time_percent": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Total time spent percentage.",
									},
									"min_cost_time": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Minimum time spent.",
									},
									"max_cost_time": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Maximum time spent.",
									},
									"first_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Earliest timestamp.",
									},
									"last_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Latest timestamp.",
									},
									"read_cost_time": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Total read I/O time.",
									},
									"write_cost_time": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Total write I/O time.",
									},
								},
							},
						},
					},
				},
			},

			"error_log_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Error log details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Username.",
						},
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database.",
						},
						"error_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Error time.",
						},
						"error_message": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Error message.",
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

func dataSourceTencentCloudCdwpgLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cdwpg_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CdwpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database"); ok {
		paramMap["Database"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("duration"); ok {
		paramMap["Duration"] = helper.Float64(v.(float64))
	}

	var respData *cdwpgv20201230.DescribeSlowLogResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdwpgLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	var instanceId string
	slowLogDetailsMap := map[string]interface{}{}

	if respData.SlowLogDetails != nil {
		if respData.SlowLogDetails.TotalTime != nil {
			slowLogDetailsMap["total_time"] = respData.SlowLogDetails.TotalTime
		}

		if respData.SlowLogDetails.TotalCallTimes != nil {
			slowLogDetailsMap["total_call_times"] = respData.SlowLogDetails.TotalCallTimes
		}

		normalQuerysList := make([]map[string]interface{}, 0, len(respData.SlowLogDetails.NormalQuerys))
		if respData.SlowLogDetails.NormalQuerys != nil {
			for _, normalQuerys := range respData.SlowLogDetails.NormalQuerys {
				normalQuerysMap := map[string]interface{}{}

				if normalQuerys.CallTimes != nil {
					normalQuerysMap["call_times"] = normalQuerys.CallTimes
				}

				if normalQuerys.SharedReadBlocks != nil {
					normalQuerysMap["shared_read_blocks"] = normalQuerys.SharedReadBlocks
				}

				if normalQuerys.SharedWriteBlocks != nil {
					normalQuerysMap["shared_write_blocks"] = normalQuerys.SharedWriteBlocks
				}

				if normalQuerys.DatabaseName != nil {
					normalQuerysMap["database_name"] = normalQuerys.DatabaseName
				}

				if normalQuerys.NormalQuery != nil {
					normalQuerysMap["normal_query"] = normalQuerys.NormalQuery
				}

				if normalQuerys.MaxElapsedQuery != nil {
					normalQuerysMap["max_elapsed_query"] = normalQuerys.MaxElapsedQuery
				}

				if normalQuerys.CostTime != nil {
					normalQuerysMap["cost_time"] = normalQuerys.CostTime
				}

				if normalQuerys.ClientIp != nil {
					normalQuerysMap["client_ip"] = normalQuerys.ClientIp
				}

				if normalQuerys.UserName != nil {
					normalQuerysMap["user_name"] = normalQuerys.UserName
				}

				if normalQuerys.TotalCallTimesPercent != nil {
					normalQuerysMap["total_call_times_percent"] = normalQuerys.TotalCallTimesPercent
				}

				if normalQuerys.TotalCostTimePercent != nil {
					normalQuerysMap["total_cost_time_percent"] = normalQuerys.TotalCostTimePercent
				}

				if normalQuerys.MinCostTime != nil {
					normalQuerysMap["min_cost_time"] = normalQuerys.MinCostTime
				}

				if normalQuerys.MaxCostTime != nil {
					normalQuerysMap["max_cost_time"] = normalQuerys.MaxCostTime
				}

				if normalQuerys.FirstTime != nil {
					normalQuerysMap["first_time"] = normalQuerys.FirstTime
				}

				if normalQuerys.LastTime != nil {
					normalQuerysMap["last_time"] = normalQuerys.LastTime
				}

				if normalQuerys.ReadCostTime != nil {
					normalQuerysMap["read_cost_time"] = normalQuerys.ReadCostTime
				}

				if normalQuerys.WriteCostTime != nil {
					normalQuerysMap["write_cost_time"] = normalQuerys.WriteCostTime
				}

				normalQuerysList = append(normalQuerysList, normalQuerysMap)
			}

			slowLogDetailsMap["normal_querys"] = normalQuerysList
		}
		_ = d.Set("slow_log_details", []interface{}{slowLogDetailsMap})
	}

	paramMap1 := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap1["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap1["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap1["EndTime"] = helper.String(v.(string))
	}

	var respData1 []*cdwpgv20201230.ErrorLogDetail
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdwpgLogByFilter1(ctx, paramMap1)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData1 = result
		return nil
	})
	if err != nil {
		return err
	}

	errorLogDetailsList := make([]map[string]interface{}, 0, len(respData1))
	if respData1 != nil {
		for _, errorLogDetails := range respData1 {
			errorLogDetailsMap := map[string]interface{}{}

			if errorLogDetails.UserName != nil {
				errorLogDetailsMap["user_name"] = errorLogDetails.UserName
			}

			if errorLogDetails.Database != nil {
				errorLogDetailsMap["database"] = errorLogDetails.Database
			}

			if errorLogDetails.ErrorTime != nil {
				errorLogDetailsMap["error_time"] = errorLogDetails.ErrorTime
			}

			if errorLogDetails.ErrorMessage != nil {
				errorLogDetailsMap["error_message"] = errorLogDetails.ErrorMessage
			}

			errorLogDetailsList = append(errorLogDetailsList, errorLogDetailsMap)
		}

		_ = d.Set("error_log_details", errorLogDetailsList)
	}

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

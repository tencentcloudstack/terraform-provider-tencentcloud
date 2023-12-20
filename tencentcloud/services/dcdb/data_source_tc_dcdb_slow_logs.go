package dcdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcdbSlowLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbSlowLogsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of `tdsqlshard-ow728lmc`.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query start time in the format of 2016-07-23 14:55:20.",
			},

			"shard_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance shard ID in the format of `shard-rc754ljk`.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query end time in the format of 2016-08-22 14:55:20.",
			},

			"db": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specific name of the database to be queried.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting metric. Valid values: query_time_sum, query_count.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting order. Valid values: desc, asc.",
			},

			"slave": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Query slow queries from either the primary or the replica. Valid values: 0 (primary), 1 (replica).",
			},

			"lock_time_sum": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "Total statement lock time.",
			},

			"query_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total number of statement queries.",
			},

			"query_time_sum": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "Total statement query time.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Slow query log data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_sum": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Statement checksum for querying details.",
						},
						"db": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
						"finger_print": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Abstracted SQL statement.",
						},
						"lock_time_avg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Average lock time.",
						},
						"lock_time_max": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Maximum lock time.",
						},
						"lock_time_min": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Minimum lock time.",
						},
						"lock_time_sum": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Total lock time.",
						},
						"query_count": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Number of queries.",
						},
						"query_time_avg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Average query time.",
						},
						"query_time_max": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Maximum query time.",
						},
						"query_time_min": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Minimum query time.",
						},
						"query_time_sum": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Total query time.",
						},
						"rows_examined_sum": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Number of scanned rows.",
						},
						"rows_sent_sum": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Number of sent rows.",
						},
						"ts_max": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last execution time.",
						},
						"ts_min": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "First execution time.",
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account.",
						},
						"example_sql": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sample SQLNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host address of account.",
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

func dataSourceTencentCloudDcdbSlowLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_slow_logs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var (
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("shard_id"); ok {
		paramMap["ShardId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db"); ok {
		paramMap["Db"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("slave"); v != nil {
		paramMap["Slave"] = helper.IntInt64(v.(int))
	}

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		resp     *dcdb.DescribeDBSlowLogsResponseParams
		slowLogs []*dcdb.SlowLogData
		e        error
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		slowLogs, resp, e = service.DescribeDcdbSlowLogsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG]%s quey dcdb slow log success, slowLogs.len:%v, resp:[%v], \n ", //result.LockTimeSum:[%v], result.QueryTimeSum:[%v]
		logId, len(slowLogs), resp)
	// logId, len(slowLogs), result, result.LockTimeSum, result.QueryTimeSum)

	if resp != nil {
		if resp.LockTimeSum != nil {
			_ = d.Set("lock_time_sum", resp.LockTimeSum)
		}

		if resp.QueryCount != nil {
			_ = d.Set("query_count", resp.QueryCount)
		}

		if resp.QueryTimeSum != nil {
			_ = d.Set("query_time_sum", resp.QueryTimeSum)
		}
	}

	if slowLogs != nil {
		slowLogDataList := []interface{}{}
		for _, slowLogData := range slowLogs {
			slowLogDataMap := map[string]interface{}{}

			if slowLogData.CheckSum != nil {
				slowLogDataMap["check_sum"] = slowLogData.CheckSum
			}

			if slowLogData.Db != nil {
				slowLogDataMap["db"] = slowLogData.Db
			}

			if slowLogData.FingerPrint != nil {
				slowLogDataMap["finger_print"] = slowLogData.FingerPrint
			}

			if slowLogData.LockTimeAvg != nil {
				slowLogDataMap["lock_time_avg"] = slowLogData.LockTimeAvg
			}

			if slowLogData.LockTimeMax != nil {
				slowLogDataMap["lock_time_max"] = slowLogData.LockTimeMax
			}

			if slowLogData.LockTimeMin != nil {
				slowLogDataMap["lock_time_min"] = slowLogData.LockTimeMin
			}

			if slowLogData.LockTimeSum != nil {
				slowLogDataMap["lock_time_sum"] = slowLogData.LockTimeSum
			}

			if slowLogData.QueryCount != nil {
				slowLogDataMap["query_count"] = slowLogData.QueryCount
			}

			if slowLogData.QueryTimeAvg != nil {
				slowLogDataMap["query_time_avg"] = slowLogData.QueryTimeAvg
			}

			if slowLogData.QueryTimeMax != nil {
				slowLogDataMap["query_time_max"] = slowLogData.QueryTimeMax
			}

			if slowLogData.QueryTimeMin != nil {
				slowLogDataMap["query_time_min"] = slowLogData.QueryTimeMin
			}

			if slowLogData.QueryTimeSum != nil {
				slowLogDataMap["query_time_sum"] = slowLogData.QueryTimeSum
			}

			if slowLogData.RowsExaminedSum != nil {
				slowLogDataMap["rows_examined_sum"] = slowLogData.RowsExaminedSum
			}

			if slowLogData.RowsSentSum != nil {
				slowLogDataMap["rows_sent_sum"] = slowLogData.RowsSentSum
			}

			if slowLogData.TsMax != nil {
				slowLogDataMap["ts_max"] = slowLogData.TsMax
			}

			if slowLogData.TsMin != nil {
				slowLogDataMap["ts_min"] = slowLogData.TsMin
			}

			if slowLogData.User != nil {
				slowLogDataMap["user"] = slowLogData.User
			}

			if slowLogData.ExampleSql != nil {
				slowLogDataMap["example_sql"] = slowLogData.ExampleSql
			}

			if slowLogData.Host != nil {
				slowLogDataMap["host"] = slowLogData.Host
			}

			slowLogDataList = append(slowLogDataList, slowLogDataMap)
		}

		_ = d.Set("data", slowLogDataList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), slowLogs); e != nil {
			return e
		}
	}
	return nil
}

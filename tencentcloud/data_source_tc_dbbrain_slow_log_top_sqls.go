/*
Use this data source to query detailed information of dbbrain slow_log_top_sqls

Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_top_sqls" "test" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  sort_by = "QueryTimeMax"
  order_by = "ASC"
  product = "mysql"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainSlowLogTopSqls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSlowLogTopSqlsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time, such as `2019-09-10 12:13:14`.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The deadline, such as `2019-09-11 10:13:14`, the interval between the deadline and the start time is less than 7 days.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort key, currently supports sort keys such as QueryTime, ExecTimes, RowsSent, LockTime and RowsExamined, the default is QueryTime.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The sorting method supports ASC (ascending) and DESC (descending). The default is DESC.",
			},

			"schema_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Array of database names.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schema": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DB name.",
						},
					},
				},
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database CynosDB for MySQL, the default is `mysql`.",
			},

			"rows": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Slow log top sql list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lock_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "SQL total lock waiting time, in seconds.",
						},
						"lock_time_max": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Maximum lock waiting time, in seconds.",
						},
						"lock_time_min": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Minimum lock waiting time, in seconds.",
						},
						"rows_examined": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total scan lines.",
						},
						"rows_examined_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of scan lines.",
						},
						"rows_examined_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum number of scan lines.",
						},
						"query_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total time, in seconds.",
						},
						"query_time_max": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The maximum execution time, in seconds.",
						},
						"query_time_min": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The minimum execution time, in seconds.",
						},
						"rows_sent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total number of rows returned.",
						},
						"rows_sent_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of rows returned.",
						},
						"rows_sent_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum number of rows returned.",
						},
						"exec_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Execution times.",
						},
						"sql_template": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "sql template.",
						},
						"sql_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL with parameters (random).",
						},
						"schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
						"query_time_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total time-consuming ratio, unit %.",
						},
						"lock_time_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The ratio of the total lock waiting time of SQL, in %.",
						},
						"rows_examined_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The proportion of the total number of scanned lines, unit %.",
						},
						"rows_sent_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The proportion of the total number of rows returned, in %.",
						},
						"query_time_avg": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average execution time, in seconds.",
						},
						"rows_sent_avg": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "average number of rows returned.",
						},
						"lock_time_avg": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average lock waiting time, in seconds.",
						},
						"rows_examined_avg": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "average number of lines scanned.",
						},
						"md5": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MD5 value of SOL template.",
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

func dataSourceTencentCloudDbbrainSlowLogTopSqlsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_slow_log_top_sqls.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var id string
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
		id = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["start_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["sort_by"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["order_by"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("schema_list"); ok {
		schemaListSet := v.([]interface{})
		tmpSet := make([]*dbbrain.SchemaItem, 0, len(schemaListSet))

		for _, item := range schemaListSet {
			schemaItem := dbbrain.SchemaItem{}
			schemaItemMap := item.(map[string]interface{})

			if v, ok := schemaItemMap["schema"]; ok {
				schemaItem.Schema = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &schemaItem)
		}
		paramMap["schema_list"] = tmpSet
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var rows []*dbbrain.SlowLogTopSqlItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainSlowLogTopSqlsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		rows = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(rows))

	if rows != nil {
		for _, slowLogTopSqlItem := range rows {
			slowLogTopSqlItemMap := map[string]interface{}{}

			if slowLogTopSqlItem.LockTime != nil {
				slowLogTopSqlItemMap["lock_time"] = slowLogTopSqlItem.LockTime
			}

			if slowLogTopSqlItem.LockTimeMax != nil {
				slowLogTopSqlItemMap["lock_time_max"] = slowLogTopSqlItem.LockTimeMax
			}

			if slowLogTopSqlItem.LockTimeMin != nil {
				slowLogTopSqlItemMap["lock_time_min"] = slowLogTopSqlItem.LockTimeMin
			}

			if slowLogTopSqlItem.RowsExamined != nil {
				slowLogTopSqlItemMap["rows_examined"] = slowLogTopSqlItem.RowsExamined
			}

			if slowLogTopSqlItem.RowsExaminedMax != nil {
				slowLogTopSqlItemMap["rows_examined_max"] = slowLogTopSqlItem.RowsExaminedMax
			}

			if slowLogTopSqlItem.RowsExaminedMin != nil {
				slowLogTopSqlItemMap["rows_examined_min"] = slowLogTopSqlItem.RowsExaminedMin
			}

			if slowLogTopSqlItem.QueryTime != nil {
				slowLogTopSqlItemMap["query_time"] = slowLogTopSqlItem.QueryTime
			}

			if slowLogTopSqlItem.QueryTimeMax != nil {
				slowLogTopSqlItemMap["query_time_max"] = slowLogTopSqlItem.QueryTimeMax
			}

			if slowLogTopSqlItem.QueryTimeMin != nil {
				slowLogTopSqlItemMap["query_time_min"] = slowLogTopSqlItem.QueryTimeMin
			}

			if slowLogTopSqlItem.RowsSent != nil {
				slowLogTopSqlItemMap["rows_sent"] = slowLogTopSqlItem.RowsSent
			}

			if slowLogTopSqlItem.RowsSentMax != nil {
				slowLogTopSqlItemMap["rows_sent_max"] = slowLogTopSqlItem.RowsSentMax
			}

			if slowLogTopSqlItem.RowsSentMin != nil {
				slowLogTopSqlItemMap["rows_sent_min"] = slowLogTopSqlItem.RowsSentMin
			}

			if slowLogTopSqlItem.ExecTimes != nil {
				slowLogTopSqlItemMap["exec_times"] = slowLogTopSqlItem.ExecTimes
			}

			if slowLogTopSqlItem.SqlTemplate != nil {
				slowLogTopSqlItemMap["sql_template"] = slowLogTopSqlItem.SqlTemplate
			}

			if slowLogTopSqlItem.SqlText != nil {
				slowLogTopSqlItemMap["sql_text"] = slowLogTopSqlItem.SqlText
			}

			if slowLogTopSqlItem.Schema != nil {
				slowLogTopSqlItemMap["schema"] = slowLogTopSqlItem.Schema
			}

			if slowLogTopSqlItem.QueryTimeRatio != nil {
				slowLogTopSqlItemMap["query_time_ratio"] = slowLogTopSqlItem.QueryTimeRatio
			}

			if slowLogTopSqlItem.LockTimeRatio != nil {
				slowLogTopSqlItemMap["lock_time_ratio"] = slowLogTopSqlItem.LockTimeRatio
			}

			if slowLogTopSqlItem.RowsExaminedRatio != nil {
				slowLogTopSqlItemMap["rows_examined_ratio"] = slowLogTopSqlItem.RowsExaminedRatio
			}

			if slowLogTopSqlItem.RowsSentRatio != nil {
				slowLogTopSqlItemMap["rows_sent_ratio"] = slowLogTopSqlItem.RowsSentRatio
			}

			if slowLogTopSqlItem.QueryTimeAvg != nil {
				slowLogTopSqlItemMap["query_time_avg"] = slowLogTopSqlItem.QueryTimeAvg
			}

			if slowLogTopSqlItem.RowsSentAvg != nil {
				slowLogTopSqlItemMap["rows_sent_avg"] = slowLogTopSqlItem.RowsSentAvg
			}

			if slowLogTopSqlItem.LockTimeAvg != nil {
				slowLogTopSqlItemMap["lock_time_avg"] = slowLogTopSqlItem.LockTimeAvg
			}

			if slowLogTopSqlItem.RowsExaminedAvg != nil {
				slowLogTopSqlItemMap["rows_examined_avg"] = slowLogTopSqlItem.RowsExaminedAvg
			}

			if slowLogTopSqlItem.Md5 != nil {
				slowLogTopSqlItemMap["md5"] = slowLogTopSqlItem.Md5
			}

			tmpList = append(tmpList, slowLogTopSqlItemMap)
		}

		_ = d.Set("rows", tmpList)
	}

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

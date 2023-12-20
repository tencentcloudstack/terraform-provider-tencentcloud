package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbInstanceSlowQueries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbInstanceSlowQueriesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Earliest transaction start time.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Latest transaction start time.",
			},

			"username": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "user name.",
			},

			"host": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Client host.",
			},

			"database": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort field, optional values: QueryTime, LockTime, RowsExamined, RowsSent.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort type, optional values: asc, desc.",
			},

			"slow_queries": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Slow query records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Execution timestamp.",
						},
						"query_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Execution time in seconds.",
						},
						"sql_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL statement.",
						},
						"user_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Client host.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "user name.",
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
						"lock_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Lock duration in seconds.",
						},
						"rows_examined": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Scan Rows.",
						},
						"rows_sent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Return the number of rows.",
						},
						"sql_template": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL template.",
						},
						"sql_md5": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL statement md5.",
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

func dataSourceTencentCloudCynosdbInstanceSlowQueriesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_instance_slow_queries.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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

	if v, ok := d.GetOk("username"); ok {
		paramMap["Username"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		paramMap["Host"] = helper.String(v.(string))
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

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var slowQueries []*cynosdb.SlowQueriesItem

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbInstanceSlowQueriesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		slowQueries = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(slowQueries))
	tmpList := make([]map[string]interface{}, 0, len(slowQueries))

	if slowQueries != nil {
		for _, slowQueriesItem := range slowQueries {
			slowQueriesItemMap := map[string]interface{}{}

			if slowQueriesItem.Timestamp != nil {
				slowQueriesItemMap["timestamp"] = slowQueriesItem.Timestamp
			}

			if slowQueriesItem.QueryTime != nil {
				slowQueriesItemMap["query_time"] = slowQueriesItem.QueryTime
			}

			if slowQueriesItem.SqlText != nil {
				slowQueriesItemMap["sql_text"] = slowQueriesItem.SqlText
			}

			if slowQueriesItem.UserHost != nil {
				slowQueriesItemMap["user_host"] = slowQueriesItem.UserHost
			}

			if slowQueriesItem.UserName != nil {
				slowQueriesItemMap["user_name"] = slowQueriesItem.UserName
			}

			if slowQueriesItem.Database != nil {
				slowQueriesItemMap["database"] = slowQueriesItem.Database
			}

			if slowQueriesItem.LockTime != nil {
				slowQueriesItemMap["lock_time"] = slowQueriesItem.LockTime
			}

			if slowQueriesItem.RowsExamined != nil {
				slowQueriesItemMap["rows_examined"] = slowQueriesItem.RowsExamined
			}

			if slowQueriesItem.RowsSent != nil {
				slowQueriesItemMap["rows_sent"] = slowQueriesItem.RowsSent
			}

			if slowQueriesItem.SqlTemplate != nil {
				slowQueriesItemMap["sql_template"] = slowQueriesItem.SqlTemplate
			}

			if slowQueriesItem.SqlMd5 != nil {
				slowQueriesItemMap["sql_md5"] = slowQueriesItem.SqlMd5
			}

			ids = append(ids, *slowQueriesItem.SqlMd5)
			tmpList = append(tmpList, slowQueriesItemMap)
		}

		_ = d.Set("slow_queries", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

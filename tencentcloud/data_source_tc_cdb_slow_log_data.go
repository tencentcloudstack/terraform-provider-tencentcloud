/*
Use this data source to query detailed information of cdb slow_log_data

Example Usage

```hcl
data "tencentcloud_cdb_slow_log_data" "slow_log_data" {
  instance_id = ""
  start_time =
  end_time =
  user_hosts =
  user_names =
  data_bases =
  sort_by = ""
  order_by = ""
  inst_type = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdbSlowLogData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbSlowLogDataRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start timestamp. For example 1585142640 .",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End timestamp. For example 1585142640 .",
			},

			"user_hosts": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of client hosts.",
			},

			"user_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of client usernames.",
			},

			"data_bases": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of databases accessed.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort field. Currently supported: Timestamp, QueryTime, LockTime, RowsExamined, RowsSent.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort in ascending or descending order. Currently supported: ASC,DESC.",
			},

			"inst_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Only valid when the instance is the master instance or disaster recovery instance, the optional value: slave, which means to pull the log of the slave machine.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Query records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sql execution time.",
						},
						"query_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Sql execution time (seconds).",
						},
						"sql_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sql statement.",
						},
						"user_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Client address.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User name.",
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
						"lock_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Lock duration (seconds).",
						},
						"rows_examined": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of rows to scan.",
						},
						"rows_sent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of rows in the result set.",
						},
						"sql_template": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sql template.",
						},
						"md5": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The md5 of the Sql statement.",
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

func dataSourceTencentCloudCdbSlowLogDataRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_slow_log_data.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("user_hosts"); ok {
		userHostsSet := v.(*schema.Set).List()
		paramMap["UserHosts"] = helper.InterfacesStringsPoint(userHostsSet)
	}

	if v, ok := d.GetOk("user_names"); ok {
		userNamesSet := v.(*schema.Set).List()
		paramMap["UserNames"] = helper.InterfacesStringsPoint(userNamesSet)
	}

	if v, ok := d.GetOk("data_bases"); ok {
		dataBasesSet := v.(*schema.Set).List()
		paramMap["DataBases"] = helper.InterfacesStringsPoint(dataBasesSet)
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("inst_type"); ok {
		paramMap["InstType"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbSlowLogDataByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		totalCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(totalCount))
	if items != nil {
		for _, slowLogItem := range items {
			slowLogItemMap := map[string]interface{}{}

			if slowLogItem.Timestamp != nil {
				slowLogItemMap["timestamp"] = slowLogItem.Timestamp
			}

			if slowLogItem.QueryTime != nil {
				slowLogItemMap["query_time"] = slowLogItem.QueryTime
			}

			if slowLogItem.SqlText != nil {
				slowLogItemMap["sql_text"] = slowLogItem.SqlText
			}

			if slowLogItem.UserHost != nil {
				slowLogItemMap["user_host"] = slowLogItem.UserHost
			}

			if slowLogItem.UserName != nil {
				slowLogItemMap["user_name"] = slowLogItem.UserName
			}

			if slowLogItem.Database != nil {
				slowLogItemMap["database"] = slowLogItem.Database
			}

			if slowLogItem.LockTime != nil {
				slowLogItemMap["lock_time"] = slowLogItem.LockTime
			}

			if slowLogItem.RowsExamined != nil {
				slowLogItemMap["rows_examined"] = slowLogItem.RowsExamined
			}

			if slowLogItem.RowsSent != nil {
				slowLogItemMap["rows_sent"] = slowLogItem.RowsSent
			}

			if slowLogItem.SqlTemplate != nil {
				slowLogItemMap["sql_template"] = slowLogItem.SqlTemplate
			}

			if slowLogItem.Md5 != nil {
				slowLogItemMap["md5"] = slowLogItem.Md5
			}

			ids = append(ids, *slowLogItem.InstanceId)
			tmpList = append(tmpList, slowLogItemMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}

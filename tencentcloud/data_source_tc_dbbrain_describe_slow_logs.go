/*
Use this data source to query detailed information of dbbrain describe_slow_logs

Example Usage

```hcl
data "tencentcloud_dbbrain_describe_slow_logs" "describe_slow_logs" {
  product = ""
  instance_id = ""
  md5 = ""
  start_time = ""
  end_time = ""
  d_b =
  key =
  user =
  ip =
  time =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainDescribeSlowLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDescribeSlowLogsRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include： mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
			},

			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Id.",
			},

			"md5": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Md5 value of sql template.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time, such as 2019-09-10 12:13:14.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The deadline, such as 2019-09-11 10:13:14, the interval between the deadline and the start time is less than 7 days.",
			},

			"d_b": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Database list.",
			},

			"key": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Keywords.",
			},

			"user": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "User.",
			},

			"ip": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Ip.",
			},

			"time": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Time-consuming interval, the left and right boundaries of the time-consuming interval correspond to the 0th element and the first element of the array respectively.",
			},

			"rows": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Slow log details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Slow log start time.",
						},
						"sql_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sql statement.",
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User sourceNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"user_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ip sourceNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"query_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Execution time, in seconds.",
						},
						"lock_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Lock time, in secondsNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"rows_examined": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Scan linesNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"rows_sent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Return the number of rowsNote: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudDbbrainDescribeSlowLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_describe_slow_logs.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("md5"); ok {
		paramMap["Md5"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b"); ok {
		dBSet := v.(*schema.Set).List()
		paramMap["DB"] = helper.InterfacesStringsPoint(dBSet)
	}

	if v, ok := d.GetOk("key"); ok {
		keySet := v.(*schema.Set).List()
		paramMap["Key"] = helper.InterfacesStringsPoint(keySet)
	}

	if v, ok := d.GetOk("user"); ok {
		userSet := v.(*schema.Set).List()
		paramMap["User"] = helper.InterfacesStringsPoint(userSet)
	}

	if v, ok := d.GetOk("ip"); ok {
		ipSet := v.(*schema.Set).List()
		paramMap["Ip"] = helper.InterfacesStringsPoint(ipSet)
	}

	if v, ok := d.GetOk("time"); ok {
		timeSet := v.(*schema.Set).List()
		for i := range timeSet {
			time := timeSet[i].(int)
			paramMap["Time"] = append(paramMap["Time"], helper.IntInt64(time))
		}
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var rows []*dbbrain.SlowLogInfoItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainDescribeSlowLogsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		rows = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(rows))
	tmpList := make([]map[string]interface{}, 0, len(rows))

	if rows != nil {
		for _, slowLogInfoItem := range rows {
			slowLogInfoItemMap := map[string]interface{}{}

			if slowLogInfoItem.Timestamp != nil {
				slowLogInfoItemMap["timestamp"] = slowLogInfoItem.Timestamp
			}

			if slowLogInfoItem.SqlText != nil {
				slowLogInfoItemMap["sql_text"] = slowLogInfoItem.SqlText
			}

			if slowLogInfoItem.Database != nil {
				slowLogInfoItemMap["database"] = slowLogInfoItem.Database
			}

			if slowLogInfoItem.UserName != nil {
				slowLogInfoItemMap["user_name"] = slowLogInfoItem.UserName
			}

			if slowLogInfoItem.UserHost != nil {
				slowLogInfoItemMap["user_host"] = slowLogInfoItem.UserHost
			}

			if slowLogInfoItem.QueryTime != nil {
				slowLogInfoItemMap["query_time"] = slowLogInfoItem.QueryTime
			}

			if slowLogInfoItem.LockTime != nil {
				slowLogInfoItemMap["lock_time"] = slowLogInfoItem.LockTime
			}

			if slowLogInfoItem.RowsExamined != nil {
				slowLogInfoItemMap["rows_examined"] = slowLogInfoItem.RowsExamined
			}

			if slowLogInfoItem.RowsSent != nil {
				slowLogInfoItemMap["rows_sent"] = slowLogInfoItem.RowsSent
			}

			ids = append(ids, *slowLogInfoItem.Md5)
			tmpList = append(tmpList, slowLogInfoItemMap)
		}

		_ = d.Set("rows", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

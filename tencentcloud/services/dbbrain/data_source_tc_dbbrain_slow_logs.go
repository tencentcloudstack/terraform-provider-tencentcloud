package dbbrain

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainSlowLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSlowLogsRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
			},

			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance Id.",
			},

			"md5": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "md5 value of sql template.",
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

			"db": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "database list.",
			},

			"key": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "keywords.",
			},

			"user": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "user.",
			},

			"ip": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ip.",
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
							Description: "sql statement.",
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "database.",
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
							Description: "lock time, in secondsNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"rows_examined": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "scan linesNote: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudDbbrainSlowLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_slow_logs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	var (
		md5        string
		instanceId string
		product    string
	)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
		product = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("md5"); ok {
		paramMap["Md5"] = helper.String(v.(string))
		md5 = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db"); ok {
		dbSet := v.(*schema.Set).List()
		paramMap["Db"] = helper.InterfacesStringsPoint(dbSet)
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
		tmpList := []interface{}{}
		for _, v := range timeSet {
			if v != nil {
				time := v.(int)
				tmpList = append(tmpList, helper.IntInt64(time))
			}
		}
		paramMap["Time"] = tmpList
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var rows []*dbbrain.SlowLogInfoItem

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainSlowLogsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
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

			ids = append(ids, strings.Join([]string{md5, instanceId, product}, tccommon.FILED_SP))
			tmpList = append(tmpList, slowLogInfoItemMap)
		}

		_ = d.Set("rows", tmpList)
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

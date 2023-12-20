package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlSlowLogData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlSlowLogDataRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start timestamp. For example 1585142640.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End timestamp. For example 1585142640.",
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
							Description: "client address.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "user name.",
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "database name.",
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

func dataSourceTencentCloudMysqlSlowLogDataRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_slow_log_data.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
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

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var items []*cdb.SlowLogItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlSlowLogDataByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(items))
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

			tmpList = append(tmpList, slowLogItemMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

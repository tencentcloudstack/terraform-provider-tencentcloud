package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbAuditLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbAuditLogsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time, format: 2017-07-12 10:29:20.",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The end time is in the format of 2017-07-12 10:29:20.",
			},
			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort by. The supported values include: ASC - ascending order, DESC - descending order.",
			},
			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort fields. The supported values include: timestamp - timestamp; &amp;#39;effectRows&amp;#39; - affects the number of rows; &amp;#39;execTime&amp;#39; - Execution time.",
			},
			"filter": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Filter conditions. You can filter logs according to the set filtering criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Client address.",
						},
						"user": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "User name.",
						},
						"db_name": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Database name.",
						},
						"table_name": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Table name.",
						},
						"policy_name": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Audit policy name.",
						},
						"sql": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL statement. Supports fuzzy matching.",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL type. Currently supported: SELECT, Insert, UPDATE, DELETE, CREATE, DROP, ALT, SET, REPLACE, EXECUTE.",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Execution time. Unit: ms. Indicates audit logs with a filter execution time greater than this value.",
						},
						"affect_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Affects the number of rows. Indicates that filtering affects audit logs with rows greater than this value.",
						},
						"sql_types": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "SQL type. Supports simultaneous querying of multiple types. Currently supported: SELECT, Insert, UPDATE, DELETE, CREATE, DROP, ALT, SET, REPLACE, EXECUTE.",
						},
						"sqls": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "SQL statement. Supports passing multiple SQL statements.",
						},
						"sent_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Returns the number of rows.",
						},
						"thread_id": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Thread ID.",
						},
					},
				},
			},
			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Audit log details. Note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"affect_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Affects the number of rows.",
						},
						"err_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Error code.",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL type.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Table name.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Audit policy name.",
						},
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
						"sql": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL statement.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Client address.",
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User name.",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Execution time.",
						},
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp.",
						},
						"sent_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of rows sent.",
						},
						"thread_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Execution thread ID.",
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

func dataSourceTencentCloudCynosdbAuditLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_audit_logs.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		items      []*cynosdb.AuditLog
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

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "filter"); ok {
		auditLogFilter := cynosdb.AuditLogFilter{}
		if v, ok := dMap["host"]; ok {
			hostSet := v.(*schema.Set).List()
			auditLogFilter.Host = helper.InterfacesStringsPoint(hostSet)
		}
		if v, ok := dMap["user"]; ok {
			userSet := v.(*schema.Set).List()
			auditLogFilter.User = helper.InterfacesStringsPoint(userSet)
		}
		if v, ok := dMap["db_name"]; ok {
			dBNameSet := v.(*schema.Set).List()
			auditLogFilter.DBName = helper.InterfacesStringsPoint(dBNameSet)
		}
		if v, ok := dMap["table_name"]; ok {
			tableNameSet := v.(*schema.Set).List()
			auditLogFilter.TableName = helper.InterfacesStringsPoint(tableNameSet)
		}
		if v, ok := dMap["policy_name"]; ok {
			policyNameSet := v.(*schema.Set).List()
			auditLogFilter.PolicyName = helper.InterfacesStringsPoint(policyNameSet)
		}
		if v, ok := dMap["sql"]; ok {
			auditLogFilter.Sql = helper.String(v.(string))
		}
		if v, ok := dMap["sql_type"]; ok {
			auditLogFilter.SqlType = helper.String(v.(string))
		}
		if v, ok := dMap["exec_time"]; ok {
			auditLogFilter.ExecTime = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["affect_rows"]; ok {
			auditLogFilter.AffectRows = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["sql_types"]; ok {
			sqlTypesSet := v.(*schema.Set).List()
			auditLogFilter.SqlTypes = helper.InterfacesStringsPoint(sqlTypesSet)
		}
		if v, ok := dMap["sqls"]; ok {
			sqlsSet := v.(*schema.Set).List()
			auditLogFilter.Sqls = helper.InterfacesStringsPoint(sqlsSet)
		}
		if v, ok := dMap["sent_rows"]; ok {
			auditLogFilter.SentRows = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["thread_id"]; ok {
			threadIdSet := v.(*schema.Set).List()
			auditLogFilter.ThreadId = helper.InterfacesStringsPoint(threadIdSet)
		}
		paramMap["filter"] = &auditLogFilter
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbAuditLogsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		items = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, auditLog := range items {
			auditLogMap := map[string]interface{}{}

			if auditLog.AffectRows != nil {
				auditLogMap["affect_rows"] = auditLog.AffectRows
			}

			if auditLog.ErrCode != nil {
				auditLogMap["err_code"] = auditLog.ErrCode
			}

			if auditLog.SqlType != nil {
				auditLogMap["sql_type"] = auditLog.SqlType
			}

			if auditLog.TableName != nil {
				auditLogMap["table_name"] = auditLog.TableName
			}

			if auditLog.InstanceName != nil {
				auditLogMap["instance_name"] = auditLog.InstanceName
			}

			if auditLog.PolicyName != nil {
				auditLogMap["policy_name"] = auditLog.PolicyName
			}

			if auditLog.DBName != nil {
				auditLogMap["db_name"] = auditLog.DBName
			}

			if auditLog.Sql != nil {
				auditLogMap["sql"] = auditLog.Sql
			}

			if auditLog.Host != nil {
				auditLogMap["host"] = auditLog.Host
			}

			if auditLog.User != nil {
				auditLogMap["user"] = auditLog.User
			}

			if auditLog.ExecTime != nil {
				auditLogMap["exec_time"] = auditLog.ExecTime
			}

			if auditLog.Timestamp != nil {
				auditLogMap["timestamp"] = auditLog.Timestamp
			}

			if auditLog.SentRows != nil {
				auditLogMap["sent_rows"] = auditLog.SentRows
			}

			if auditLog.ThreadId != nil {
				auditLogMap["thread_id"] = auditLog.ThreadId
			}

			tmpList = append(tmpList, auditLogMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}

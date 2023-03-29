/*
Provides a resource to create a mysql audit_log_file

Example Usage

```hcl
resource "tencentcloud_mysql_audit_log_file" "audit_log_file" {
  instance_id = "cdb-fitq5t9h"
  start_time  = "2023-03-28 20:14:00"
  end_time    = "2023-03-29 20:14:00"
  order       = "ASC"
  order_by    = "timestamp"
  filter {
    host = ["30.50.207.46"]
    user = ["keep_dbbrain"]
  }
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlAuditLogFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlAuditLogFileCreate,
		Read:   resourceTencentCloudMysqlAuditLogFileRead,
		Delete: resourceTencentCloudMysqlAuditLogFileDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "The ID of instance.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "start time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "end time.",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Sort by. supported values are: `ASC`- ascending order, `DESC`- descending order.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Sort field. supported values include:`timestamp` - timestamp; `affectRows` - affected rows; `execTime` - execution time.",
			},

			"filter": {
				Optional:    true,
				Type:        schema.TypeList,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Filter condition. Logs can be filtered according to the filter conditions set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Client address.",
						},
						"user": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "User name.",
						},
						"db_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Database name.",
						},
						"table_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Table name.",
						},
						"policy_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "The name of policy.",
						},
						"sql": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL statement. support fuzzy matching.",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL type. Currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Execution time. The unit is: ms. Indicates to filter audit logs whose execution time is greater than this value.",
						},
						"affect_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Affects the number of rows. Indicates to filter audit logs whose number of affected rows is greater than this value.",
						},
						"sql_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "SQL type. Supports simultaneous query of multiple types. Currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.",
						},
						"sqls": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "SQL statement. Support passing multiple sql statements.",
						},
					},
				},
			},

			"file_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "size of file(KB).",
			},

			"download_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "download url.",
			},
		},
	}
}

func resourceTencentCloudMysqlAuditLogFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_audit_log_file.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mysql.NewCreateAuditLogFileRequest()
		response   = mysql.NewCreateAuditLogFileResponse()
		instanceId string
		fileName   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		request.Order = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		request.OrderBy = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "filter"); ok {
		auditLogFilter := mysql.AuditLogFilter{}
		if v, ok := dMap["host"]; ok {
			hostSet := v.(*schema.Set).List()
			for i := range hostSet {
				host := hostSet[i].(string)
				auditLogFilter.Host = append(auditLogFilter.Host, &host)
			}
		}
		if v, ok := dMap["user"]; ok {
			userSet := v.(*schema.Set).List()
			for i := range userSet {
				user := userSet[i].(string)
				auditLogFilter.User = append(auditLogFilter.User, &user)
			}
		}
		if v, ok := dMap["db_name"]; ok {
			dBNameSet := v.(*schema.Set).List()
			for i := range dBNameSet {
				dBName := dBNameSet[i].(string)
				auditLogFilter.DBName = append(auditLogFilter.DBName, &dBName)
			}
		}
		if v, ok := dMap["table_name"]; ok {
			tableNameSet := v.(*schema.Set).List()
			for i := range tableNameSet {
				tableName := tableNameSet[i].(string)
				auditLogFilter.TableName = append(auditLogFilter.TableName, &tableName)
			}
		}
		if v, ok := dMap["policy_name"]; ok {
			policyNameSet := v.(*schema.Set).List()
			for i := range policyNameSet {
				policyName := policyNameSet[i].(string)
				auditLogFilter.PolicyName = append(auditLogFilter.PolicyName, &policyName)
			}
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
			for i := range sqlTypesSet {
				sqlTypes := sqlTypesSet[i].(string)
				auditLogFilter.SqlTypes = append(auditLogFilter.SqlTypes, &sqlTypes)
			}
		}
		if v, ok := dMap["sqls"]; ok {
			sqlsSet := v.(*schema.Set).List()
			for i := range sqlsSet {
				sqls := sqlsSet[i].(string)
				auditLogFilter.Sqls = append(auditLogFilter.Sqls, &sqls)
			}
		}
		request.Filter = &auditLogFilter
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().CreateAuditLogFile(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql auditLogFile failed, reason:%+v", logId, err)
		return err
	}

	fileName = *response.Response.FileName
	d.SetId(instanceId + FILED_SP + fileName)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf(
		[]string{"creating"},
		[]string{"success", "failed"},
		1*readRetryTimeout,
		time.Second,
		service.MysqlAuditLogFileStateRefreshFunc(instanceId, fileName, []string{}),
	)

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudMysqlAuditLogFileRead(d, meta)
}

func resourceTencentCloudMysqlAuditLogFileRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_audit_log_file.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	auditLogFile, err := service.DescribeMysqlAuditLogFileById(ctx, instanceId, fileName)
	if err != nil {
		return err
	}

	if auditLogFile == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlAuditLogFile` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if auditLogFile.FileSize != nil {
		_ = d.Set("file_size", auditLogFile.FileSize)
	}

	if auditLogFile.DownloadUrl != nil {
		_ = d.Set("download_url", auditLogFile.DownloadUrl)
	}

	return nil
}

func resourceTencentCloudMysqlAuditLogFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_audit_log_file.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	if err := service.DeleteMysqlAuditLogFileById(ctx, instanceId, fileName); err != nil {
		return err
	}

	return nil
}

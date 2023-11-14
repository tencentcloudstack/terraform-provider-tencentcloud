/*
Provides a resource to create a cdb audit_log_file

Example Usage

```hcl
resource "tencentcloud_cdb_audit_log_file" "audit_log_file" {
  instance_id = "cdb-c1nl9rpv"
  start_time = "2022-07-12 10:29:20"
  end_time = "2022-08-12 10:29:20"
  order = "ASC"
  order_by = ""
  filter {
		host = &lt;nil&gt;
		user = &lt;nil&gt;
		d_b_name = &lt;nil&gt;
		table_name = &lt;nil&gt;
		policy_name = &lt;nil&gt;
		sql = &lt;nil&gt;
		sql_type = &lt;nil&gt;
		exec_time = &lt;nil&gt;
		affect_rows = &lt;nil&gt;
		sql_types = &lt;nil&gt;
		sqls = &lt;nil&gt;

  }
    }
```

Import

cdb audit_log_file can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_audit_log_file.audit_log_file audit_log_file_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudCdbAuditLogFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbAuditLogFileCreate,
		Read:   resourceTencentCloudCdbAuditLogFileRead,
		Update: resourceTencentCloudCdbAuditLogFileUpdate,
		Delete: resourceTencentCloudCdbAuditLogFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort by. supported values are: ASC-ascending order, DESC-descending order.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort field. supported values include:timestamp - timestampaffectRows - affected rowsexecTime - execution time.",
			},

			"filter": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Filter condition. Logs can be filterd according to the filter conditions set.",
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
						"d_b_name": {
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
				Description: "Size of file(KB).",
			},

			"download_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Download url.",
			},
		},
	}
}

func resourceTencentCloudCdbAuditLogFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_audit_log_file.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewCreateAuditLogFileRequest()
		response   = cdb.NewCreateAuditLogFileResponse()
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
		auditLogFilter := cdb.AuditLogFilter{}
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
		if v, ok := dMap["d_b_name"]; ok {
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().CreateAuditLogFile(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdb auditLogFile failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, fileName}, FILED_SP))

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 1*readRetryTimeout, time.Second, service.CdbAuditLogFileStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbAuditLogFileRead(d, meta)
}

func resourceTencentCloudCdbAuditLogFileRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_audit_log_file.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	auditLogFile, err := service.DescribeCdbAuditLogFileById(ctx, instanceId, fileName)
	if err != nil {
		return err
	}

	if auditLogFile == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbAuditLogFile` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if auditLogFile.InstanceId != nil {
		_ = d.Set("instance_id", auditLogFile.InstanceId)
	}

	if auditLogFile.StartTime != nil {
		_ = d.Set("start_time", auditLogFile.StartTime)
	}

	if auditLogFile.EndTime != nil {
		_ = d.Set("end_time", auditLogFile.EndTime)
	}

	if auditLogFile.Order != nil {
		_ = d.Set("order", auditLogFile.Order)
	}

	if auditLogFile.OrderBy != nil {
		_ = d.Set("order_by", auditLogFile.OrderBy)
	}

	if auditLogFile.Filter != nil {
		filterMap := map[string]interface{}{}

		if auditLogFile.Filter.Host != nil {
			filterMap["host"] = auditLogFile.Filter.Host
		}

		if auditLogFile.Filter.User != nil {
			filterMap["user"] = auditLogFile.Filter.User
		}

		if auditLogFile.Filter.DBName != nil {
			filterMap["d_b_name"] = auditLogFile.Filter.DBName
		}

		if auditLogFile.Filter.TableName != nil {
			filterMap["table_name"] = auditLogFile.Filter.TableName
		}

		if auditLogFile.Filter.PolicyName != nil {
			filterMap["policy_name"] = auditLogFile.Filter.PolicyName
		}

		if auditLogFile.Filter.Sql != nil {
			filterMap["sql"] = auditLogFile.Filter.Sql
		}

		if auditLogFile.Filter.SqlType != nil {
			filterMap["sql_type"] = auditLogFile.Filter.SqlType
		}

		if auditLogFile.Filter.ExecTime != nil {
			filterMap["exec_time"] = auditLogFile.Filter.ExecTime
		}

		if auditLogFile.Filter.AffectRows != nil {
			filterMap["affect_rows"] = auditLogFile.Filter.AffectRows
		}

		if auditLogFile.Filter.SqlTypes != nil {
			filterMap["sql_types"] = auditLogFile.Filter.SqlTypes
		}

		if auditLogFile.Filter.Sqls != nil {
			filterMap["sqls"] = auditLogFile.Filter.Sqls
		}

		_ = d.Set("filter", []interface{}{filterMap})
	}

	if auditLogFile.FileSize != nil {
		_ = d.Set("file_size", auditLogFile.FileSize)
	}

	if auditLogFile.DownloadUrl != nil {
		_ = d.Set("download_url", auditLogFile.DownloadUrl)
	}

	return nil
}

func resourceTencentCloudCdbAuditLogFileUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_audit_log_file.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"instance_id", "start_time", "end_time", "order", "order_by", "filter", "file_size", "download_url"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudCdbAuditLogFileRead(d, meta)
}

func resourceTencentCloudCdbAuditLogFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_audit_log_file.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	if err := service.DeleteCdbAuditLogFileById(ctx, instanceId, fileName); err != nil {
		return err
	}

	return nil
}

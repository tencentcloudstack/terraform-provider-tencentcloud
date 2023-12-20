package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudCynosdbAuditLogFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbAuditLogFileCreate,
		Read:   resourceTencentCloudCynosdbAuditLogFileRead,
		Delete: resourceTencentCloudCynosdbAuditLogFileDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"start_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Start time.",
			},

			"end_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"order": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Sort by. Supported values are: `ASC` - ascending, `DESC` - descending.",
			},

			"order_by": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Sort field. supported values are:\n`timestamp` - timestamp\n`affectRows` - affected rows\n`execTime` - execution time.",
			},

			"filter": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
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
							Description: "Client host.",
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
							Description: "The name of database.",
						},
						"table_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "The name of table.",
						},
						"policy_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "The name of audit policy.",
						},
						"sql": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL statement. Support fuzzy matching.",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL type. currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Execution time. The unit is: ms. Indicates to filter audit logs whose execution time is greater than this value.",
						},
						"affect_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Affects the number of rows. Indicates that the audit log whose number of affected rows is greater than this value is filtered.",
						},
						"sql_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "SQL type. Supports simultaneous query of multiple types. currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.",
						},
						"sqls": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "SQL statement. Support passing multiple sql statements.",
						},
						"sent_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Return the number of rows.",
						},
						"thread_id": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "The ID of thread.",
						},
					},
				},
			},
			// computed
			"file_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Audit log file name.",
			},
			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Audit log file creation time. The format is 2019-03-20 17:09:13.",
			},
			"file_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "File size, The unit is KB.",
			},
			"download_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The download address of the audit logs.",
			},
			"err_msg": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Error message.",
			},
		},
	}
}

func resourceTencentCloudCynosdbAuditLogFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_log_file.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = cynosdb.NewCreateAuditLogFileRequest()
		response   = cynosdb.NewCreateAuditLogFileResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
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
		auditLogFilter := cynosdb.AuditLogFilter{}
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
		if v, ok := dMap["sent_rows"]; ok {
			auditLogFilter.SentRows = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["thread_id"]; ok {
			threadIdSet := v.(*schema.Set).List()
			for i := range threadIdSet {
				threadId := threadIdSet[i].(string)
				auditLogFilter.ThreadId = append(auditLogFilter.ThreadId, &threadId)
			}
		}
		request.Filter = &auditLogFilter
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().CreateAuditLogFile(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb auditLogFile failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		request := cynosdb.NewDescribeAuditLogFilesRequest()
		request.InstanceId = helper.String(instanceId)
		request.FileName = response.Response.FileName
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeAuditLogFiles(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if len(result.Response.Items) > 0 && *result.Response.Items[0].Status == "success" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("%s not ready", *response.Response.FileName))
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb auditLogFile failed, reason:%+v", logId, err)
		return err
	}

	auditLogFileId := strings.Join([]string{instanceId, *response.Response.FileName}, tccommon.FILED_SP)
	d.SetId(auditLogFileId)

	return resourceTencentCloudCynosdbAuditLogFileRead(d, meta)
}

func resourceTencentCloudCynosdbAuditLogFileRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_log_file.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	auditLogFile, err := service.DescribeCynosdbAuditLogFileById(ctx, instanceId, fileName)
	if err != nil {
		return err
	}

	if auditLogFile == nil {
		d.SetId("")
		return fmt.Errorf("resource `CynosdbAuditLogFile` %s does not exist", d.Id())
	}

	_ = d.Set("file_name", *auditLogFile.FileName)
	_ = d.Set("create_time", *auditLogFile.CreateTime)
	_ = d.Set("file_size", *auditLogFile.FileSize)
	_ = d.Set("download_url", *auditLogFile.DownloadUrl)
	_ = d.Set("err_msg", *auditLogFile.ErrMsg)

	return nil
}

func resourceTencentCloudCynosdbAuditLogFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_log_file.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	if err := service.DeleteCynosdbAuditLogFileById(ctx, instanceId, fileName); err != nil {
		return err
	}

	return nil
}

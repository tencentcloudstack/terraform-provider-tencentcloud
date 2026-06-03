package postgresql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresAuditLogFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresAuditLogFileCreate,
		Read:   resourceTencentCloudPostgresAuditLogFileRead,
		Delete: resourceTencentCloudPostgresAuditLogFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Start time, format: `2026-03-25 00:00:00`.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "End time, format: `2026-03-25 01:00:00`.",
			},
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Product name, fixed value: `postgres`.",
			},
			"filter": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Filter conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"affect_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Affect rows.",
						},
						"db_name": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Database name list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Execution time.",
						},
						"host": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Host list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"sql": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "SQL statement.",
						},
						"user": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "User name list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"sql_type": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "SQL type list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			// computed
			"file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Audit log file name.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task status. Values: `success`, `running`, `failed`.",
			},
			"file_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "File size in MB.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
			"download_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Download URL.",
			},
			"err_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Error message.",
			},
			"progress": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Download progress.",
			},
			"finish_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Finish time.",
			},
		},
	}
}

func resourceTencentCloudPostgresAuditLogFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgres_audit_log_file.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		_          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request    = postgresql.NewCreateAuditLogFileRequest()
		instanceId string
		product    string
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

	if v, ok := d.GetOk("product"); ok {
		product = v.(string)
		request.Product = helper.String(product)
	}

	if v, ok := d.GetOk("filter"); ok {
		filterList := v.([]interface{})
		if len(filterList) > 0 {
			filterMap := filterList[0].(map[string]interface{})
			auditLogFilter := &postgresql.AuditLogFilter{}

			if v, ok := filterMap["affect_rows"].(int); ok && v != 0 {
				auditLogFilter.AffectRows = helper.IntUint64(v)
			}

			if v, ok := filterMap["db_name"]; ok {
				dbNameList := v.([]interface{})
				for _, item := range dbNameList {
					auditLogFilter.DBName = append(auditLogFilter.DBName, helper.String(item.(string)))
				}
			}

			if v, ok := filterMap["exec_time"].(int); ok && v != 0 {
				auditLogFilter.ExecTime = helper.IntUint64(v)
			}

			if v, ok := filterMap["host"]; ok {
				hostList := v.([]interface{})
				for _, item := range hostList {
					auditLogFilter.Host = append(auditLogFilter.Host, helper.String(item.(string)))
				}
			}

			if v, ok := filterMap["sql"].(string); ok && v != "" {
				auditLogFilter.Sql = helper.String(v)
			}

			if v, ok := filterMap["user"]; ok {
				userList := v.([]interface{})
				for _, item := range userList {
					auditLogFilter.User = append(auditLogFilter.User, helper.String(item.(string)))
				}
			}

			if v, ok := filterMap["sql_type"]; ok {
				sqlTypeList := v.([]interface{})
				for _, item := range sqlTypeList {
					auditLogFilter.SqlType = append(auditLogFilter.SqlType, helper.String(item.(string)))
				}
			}

			request.Filter = auditLogFilter
		}
	}

	// Get existing file list before creation to identify the new file
	var existingFileNames map[string]bool
	describeRequest := postgresql.NewDescribeAuditLogFilesRequest()
	describeRequest.InstanceId = helper.String(instanceId)
	describeRequest.Product = helper.String(product)
	describeRequest.Limit = helper.Uint64(300)
	describeRequest.Offset = helper.Uint64(0)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DescribeAuditLogFiles(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		existingFileNames = make(map[string]bool)
		if result != nil && result.Response != nil && result.Response.Items != nil {
			for _, item := range result.Response.Items {
				if item.FileName != nil {
					existingFileNames[*item.FileName] = true
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s describe audit log files before creation failed, reason:%+v", logId, err)
		return err
	}

	// Call CreateAuditLogFile
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateAuditLogFile(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create postgres audit log file failed, response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create postgres audit log file failed, reason:%+v", logId, err)
		return err
	}

	// Poll DescribeAuditLogFiles to find the new file and wait for success
	var fileName string
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		descReq := postgresql.NewDescribeAuditLogFilesRequest()
		descReq.InstanceId = helper.String(instanceId)
		descReq.Product = helper.String(product)
		descReq.Limit = helper.Uint64(300)
		descReq.Offset = helper.Uint64(0)

		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DescribeAuditLogFiles(descReq)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Items == nil {
			return resource.RetryableError(fmt.Errorf("waiting for audit log file to be created"))
		}

		// Find the new file that wasn't in the existing list
		for _, item := range result.Response.Items {
			if item.FileName == nil {
				continue
			}

			if existingFileNames[*item.FileName] {
				continue
			}

			// Found the new file
			if item.Status != nil && *item.Status == "success" {
				fileName = *item.FileName
				return nil
			} else if item.Status != nil && *item.Status == "failed" {
				errMsg := ""
				if item.ErrMsg != nil {
					errMsg = *item.ErrMsg
				}

				return resource.NonRetryableError(fmt.Errorf("audit log file creation failed: %s", errMsg))
			}

			// Still running
			return resource.RetryableError(fmt.Errorf("audit log file is still being created, status: %s", *item.Status))
		}

		return resource.RetryableError(fmt.Errorf("waiting for audit log file to appear"))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create postgres audit log file failed during polling, reason:%+v", logId, err)
		return err
	}

	if fileName == "" {
		return fmt.Errorf("audit log file name is empty after creation")
	}

	d.SetId(strings.Join([]string{instanceId, fileName}, tccommon.FILED_SP))

	return resourceTencentCloudPostgresAuditLogFileRead(d, meta)
}

func resourceTencentCloudPostgresAuditLogFileRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgres_audit_log_file.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		_     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	instanceId := idSplit[0]
	fileName := idSplit[1]

	// Read product from state; if not set (e.g. import), default to "postgres"
	product := "postgres"
	if v, ok := d.GetOk("product"); ok && v.(string) != "" {
		product = v.(string)
	}

	request := postgresql.NewDescribeAuditLogFilesRequest()
	request.InstanceId = helper.String(instanceId)
	request.Product = helper.String(product)
	request.FileName = helper.String(fileName)
	request.Limit = helper.Uint64(300)
	request.Offset = helper.Uint64(0)

	var fileInfo *postgresql.AuditLogFile
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DescribeAuditLogFiles(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Items == nil || len(result.Response.Items) == 0 {
			return nil
		}

		fileInfo = result.Response.Items[0]
		return nil
	})

	if err != nil {
		return err
	}

	if fileInfo == nil {
		log.Printf("[WARN]%s resource `tencentcloud_postgres_audit_log_file` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("product", product)

	if fileInfo.FileName != nil {
		_ = d.Set("file_name", fileInfo.FileName)
	}

	if fileInfo.Status != nil {
		_ = d.Set("status", fileInfo.Status)
	}

	if fileInfo.FileSize != nil {
		_ = d.Set("file_size", fileInfo.FileSize)
	}

	if fileInfo.CreateTime != nil {
		_ = d.Set("create_time", fileInfo.CreateTime)
	}

	if fileInfo.DownloadUrl != nil {
		_ = d.Set("download_url", fileInfo.DownloadUrl)
	}

	if fileInfo.ErrMsg != nil {
		_ = d.Set("err_msg", fileInfo.ErrMsg)
	}

	if fileInfo.Progress != nil {
		_ = d.Set("progress", fileInfo.Progress)
	}

	if fileInfo.FinishTime != nil {
		_ = d.Set("finish_time", fileInfo.FinishTime)
	}

	return nil
}

func resourceTencentCloudPostgresAuditLogFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgres_audit_log_file.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		_     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	instanceId := idSplit[0]
	fileName := idSplit[1]

	product := "postgres"
	if v, ok := d.GetOk("product"); ok && v.(string) != "" {
		product = v.(string)
	}

	request := postgresql.NewDeleteAuditLogFileRequest()
	request.InstanceId = helper.String(instanceId)
	request.Product = helper.String(product)
	request.FileName = helper.String(fileName)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DeleteAuditLogFile(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete postgres audit log file failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

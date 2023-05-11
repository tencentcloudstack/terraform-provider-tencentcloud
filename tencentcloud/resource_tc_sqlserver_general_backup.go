/*
Provides a resource to create a sqlserver general_backup

Example Usage

```hcl
resource "tencentcloud_sqlserver_general_backup" "general_backup" {
  strategy = 0
  db_names = ["db1", "db2"]
  instance_id = "mssql-i1z41iwd"
  backup_name = "bk_name"
}
```

Import

sqlserver general_backups can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_backups.general_backups general_backups_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverGeneralBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralBackupCreate,
		Read:   resourceTencentCloudSqlserverGeneralBackupRead,
		Update: resourceTencentCloudSqlserverGeneralBackupUpdate,
		Delete: resourceTencentCloudSqlserverGeneralBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"strategy": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "Backup policy (0: instance backup, 1: multi-database backup).",
			},
			"db_names": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of names of databases to be backed up (required only for multi-database backup).",
			},
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of mssql-i1z41iwd.",
			},
			"backup_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup name. If this parameter is left empty, a backup name in the format of [Instance ID]_[Backup start timestamp] will be automatically generated.",
			},
			"backup_files": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Details of backup file list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup file format (pkg - packaged backup file, single - single library backup file).",
						},
						"backup_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "backup name.",
						},
						"backup_way": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup mode, 0-scheduled backup; 1-manual temporary backup; 2-regular backup.",
						},
						"cross_backup_addr": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Destination domain download link for cross-region backup.",
						},
						"cross_backup_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Target region and backup status of cross-region backup.",
						},
						"dbs": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The name of the library for backing up files.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "end time.",
						},
						"external_addr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "External network download address, this value is not returned for single-database backup files; the download address of single-database backup files is obtained through the DescribeBackupFiles interface.",
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "backup file name.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Aggregate Id, this value is not returned for packaged backup files. Use this value to call the DescribeBackupFiles interface to obtain the detailed information of a single database backup file.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "backup id.",
						},
						"internal_addr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intranet download address, this value is not returned for a single database backup file; the download address of a single database backup file is obtained through the DescribeBackupFiles interface.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "file size(k).",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start time.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup file status (0-creating; 1-success; 2-failure).",
						},
						"strategy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup strategy (0-instance backup; 1-multi-database backup).",
						},
					},
				},
			},
			"backup_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup ID.",
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_backup.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = sqlserver.NewCreateBackupRequest()
		instanceId string
		flowId     string
		backupId   uint64
		startStr   string
		endStr     string
		err        error
	)

	if v, ok := d.GetOk("strategy"); ok {
		request.Strategy = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("db_names"); ok {
		dBNamesSet := v.(*schema.Set).List()
		for i := range dBNamesSet {
			dBNames := dBNamesSet[i].(string)
			request.DBNames = append(request.DBNames, &dBNames)
		}
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = *helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_name"); ok {
		request.BackupName = helper.String(v.(string))
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CreateBackup(request)
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			err = fmt.Errorf("sqlserver Backup %s not exists", instanceId)
			return resource.NonRetryableError(err)
		}

		flowId = strconv.FormatInt(*result.Response.FlowId, 10)
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver Backup failed, reason:%+v", logId, err)
		return err
	}

	// waiting for backup done.
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBackupByFlowId(ctx, instanceId, flowId)
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			err = fmt.Errorf("sqlserver Backup %s not exists", instanceId)
			return resource.NonRetryableError(err)
		}

		if *result.Response.Status == SQLSERVER_BACKUP_RUNNING {
			return resource.RetryableError(fmt.Errorf("create sqlserver Backup task status is running"))
		}

		if *result.Response.Status == SQLSERVER_BACKUP_SUCCESS {
			backupId = *result.Response.Id
			startStr = *result.Response.StartTime
			endStr = *result.Response.EndTime
			return nil
		}

		if *result.Response.Status == SQLSERVER_BACKUP_FAIL {
			return resource.NonRetryableError(fmt.Errorf("create sqlserver Backup task status is failed"))
		}

		err = fmt.Errorf("create sqlserver Backup task status is %v, we won't wait for it finish", *result.Response.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver Backup task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{instanceId, strconv.Itoa(int(backupId)), flowId, startStr, endStr}, FILED_SP))
	return resourceTencentCloudSqlserverGeneralBackupRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_backup.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId string
		startStr   string
		endStr     string
		backupId   uint64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId = idSplit[0]
	tempD, _ := strconv.Atoi(idSplit[1])
	backupId = uint64(tempD)
	startStr = idSplit[3]
	endStr = idSplit[4]

	backupList, err := service.DescribeSqlserverBackupByBackupId(ctx, instanceId, startStr, endStr, backupId)
	if err != nil {
		return err
	}

	if backupList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralBackups` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	list := make([]map[string]interface{}, 0, len(backupList))
	backupInfo := backupList[0]
	var infoMap = map[string]interface{}{
		"backup_format":       backupInfo.BackupFormat,
		"backup_name":         backupInfo.BackupName,
		"backup_way":          backupInfo.BackupWay,
		"cross_backup_addr":   backupInfo.CrossBackupAddr,
		"cross_backup_status": backupInfo.CrossBackupStatus,
		"dbs":                 backupInfo.DBs,
		"end_time":            backupInfo.EndTime,
		"external_addr":       backupInfo.ExternalAddr,
		"file_name":           backupInfo.FileName,
		"group_id":            backupInfo.GroupId,
		"id":                  backupInfo.Id,
		"internal_addr":       backupInfo.InternalAddr,
		"region":              backupInfo.Region,
		"size":                backupInfo.Size,
		"start_time":          backupInfo.StartTime,
		"status":              backupInfo.Status,
		"strategy":            backupInfo.Strategy,
	}
	list = append(list, infoMap)
	_ = d.Set("backup_files", list)
	_ = d.Set("backup_name", backupInfo.BackupName)
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("strategy", backupInfo.Strategy)
	return nil
}

func resourceTencentCloudSqlserverGeneralBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_backup.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = sqlserver.NewModifyBackupNameRequest()
		instanceId string
		backupId   uint64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId = idSplit[0]
	tempD, _ := strconv.Atoi(idSplit[1])
	backupId = uint64(tempD)

	backupList := d.Get("backup_files").([]interface{})
	backupInfo := backupList[0].(map[string]interface{})
	groupID := backupInfo["group_id"].(string)

	if d.HasChange("backup_name") {
		if v, ok := d.GetOk("backup_name"); ok {
			request.BackupName = helper.String(v.(string))
		}
	}

	request.InstanceId = &instanceId
	request.BackupId = &backupId
	request.GroupId = &groupID

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyBackupName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver generalBackups failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverGeneralBackupRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_backup.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId string
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId = idSplit[0]

	backupList := d.Get("backup_files").([]interface{})
	backupInfo := backupList[0].(map[string]interface{})
	fileName := backupInfo["file_name"].(string)
	if err := service.DeleteSqlserverGeneralBackupsById(ctx, instanceId, fileName); err != nil {
		return err
	}

	return nil
}

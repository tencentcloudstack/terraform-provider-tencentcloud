package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlBackupCreate,
		Read:   resourceTencentCloudMysqlBackupRead,
		Delete: resourceTencentCloudMysqlBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, such as `cdb-c1nl9rpv`. It is identical to the instance ID displayed in the database console page.",
			},
			"backup_method": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target backup method. Supported values include: `logical` - logical cold backup, `physical` - physical cold backup, `snapshot` - snapshot backup. Basic edition instances only support snapshot backup.",
			},
			"backup_db_table_list": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "List of databases and tables to backup. Only valid when `backup_method` is `logical`. The specified databases and tables must exist, otherwise backup may fail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"table": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Table name. If specified, backup this table in the database. If not specified, backup the entire database.",
						},
					},
				},
			},
			"manual_backup_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Manual backup alias. Maximum length is 60 characters.",
			},
			"encryption_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Whether to encrypt physical backup. Supported values include: `on` - yes, `off` - no. Only valid when `backup_method` is `physical`. If not specified, the instance's default backup encryption policy is used.",
			},
			"backup_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "ID of the backup task.",
			},
			"intranet_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Intranet download URL of the backup file.",
			},
			"internet_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Internet download URL of the backup file.",
			},
		},
	}
}

func resourceTencentCloudMysqlBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		instanceId string
	)
	request := mysql.NewCreateBackupRequest()

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_method"); ok {
		request.BackupMethod = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_db_table_list"); ok {
		dbTableList := v.([]interface{})
		for _, item := range dbTableList {
			dbTableMap := item.(map[string]interface{})
			backupItem := &mysql.BackupItem{}
			if db, exist := dbTableMap["database"]; exist {
				backupItem.Db = helper.String(db.(string))
			}
			if table, exist := dbTableMap["table"]; exist {
				if tableStr := table.(string); tableStr != "" {
					backupItem.Table = helper.String(tableStr)
				}
			}
			request.BackupDBTableList = append(request.BackupDBTableList, backupItem)
		}
	}

	if v, ok := d.GetOk("manual_backup_name"); ok {
		request.ManualBackupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("encryption_flag"); ok {
		request.EncryptionFlag = helper.String(v.(string))
	}

	var response *mysql.CreateBackupResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateBackup(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql_backup failed, reason:%+v", logId, err)
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.BackupId == nil {
		return fmt.Errorf("mysql_backup create failed, BackupId is empty")
	}
	backupId := helper.UInt64ToStr(*response.Response.BackupId)

	d.SetId(backupId + tccommon.FILED_SP + instanceId)

	// Wait for the backup task to finish since CreateBackup is asynchronous.
	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		backupInfos, e := service.DescribeBackupsByMysqlId(ctx, instanceId, 1000)
		if e != nil {
			return tccommon.RetryError(e)
		}

		for _, item := range backupInfos {
			if item.BackupId == nil || helper.Int64ToStr(*item.BackupId) != backupId {
				continue
			}
			if item.Status == nil {
				return resource.RetryableError(fmt.Errorf("mysql_backup [%s] status is nil, retrying", backupId))
			}
			switch *item.Status {
			case "SUCCESS":
				return nil
			case "FAILED":
				return resource.NonRetryableError(fmt.Errorf("mysql_backup [%s] create failed, status: FAILED", backupId))
			default:
				return resource.RetryableError(fmt.Errorf("mysql_backup [%s] is in unknown status [%s], retrying", backupId, *item.Status))
			}
		}

		return resource.RetryableError(fmt.Errorf("mysql_backup [%s] not found yet, retrying", backupId))
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait for mysql_backup [%s] ready failed, reason:%+v", logId, backupId, err)
		return err
	}

	return resourceTencentCloudMysqlBackupRead(d, meta)
}

func resourceTencentCloudMysqlBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	backupId := idSplit[0]
	instanceId := idSplit[1]

	var backupInfo *mysql.BackupInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		backupInfos, e := service.DescribeBackupsByMysqlId(ctx, instanceId, 1000)
		if e != nil {
			return tccommon.RetryError(e)
		}
		for _, item := range backupInfos {
			if item.BackupId != nil && helper.Int64ToStr(*item.BackupId) == backupId {
				backupInfo = item
				break
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if backupInfo == nil {
		log.Printf("[CRUD] mysql_backup id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	if backupInfo.BackupId != nil {
		_ = d.Set("backup_id", backupInfo.BackupId)
	}
	if backupInfo.Type != nil {
		_ = d.Set("backup_method", backupInfo.Type)
	}
	if backupInfo.ManualBackupName != nil {
		_ = d.Set("manual_backup_name", backupInfo.ManualBackupName)
	}
	if backupInfo.EncryptionFlag != nil && backupInfo.Type != nil && *backupInfo.Type == "physical" {
		_ = d.Set("encryption_flag", backupInfo.EncryptionFlag)
	}
	if backupInfo.IntranetUrl != nil {
		_ = d.Set("intranet_url", backupInfo.IntranetUrl)
	}
	if backupInfo.InternetUrl != nil {
		_ = d.Set("internet_url", backupInfo.InternetUrl)
	}

	return nil
}

func resourceTencentCloudMysqlBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	backupId := idSplit[0]
	instanceId := idSplit[1]

	request := mysql.NewDeleteBackupRequest()
	request.InstanceId = helper.String(instanceId)
	request.BackupId = helper.StrToInt64Point(backupId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().DeleteBackup(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete mysql_backup failed, reason:%+v", logId, err)
		return err
	}

	// Wait for the backup to be removed since DeleteBackup is asynchronous.
	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		backupInfos, e := service.DescribeBackupsByMysqlId(ctx, instanceId, 1000)
		if e != nil {
			return tccommon.RetryError(e)
		}

		for _, item := range backupInfos {
			if item.BackupId != nil && helper.Int64ToStr(*item.BackupId) == backupId {
				return resource.RetryableError(fmt.Errorf("mysql_backup [%s] is still being deleted, retrying", backupId))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait for mysql_backup [%s] delete failed, reason:%+v", logId, backupId, err)
		return err
	}

	return nil
}

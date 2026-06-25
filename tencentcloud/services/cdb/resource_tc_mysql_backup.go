package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

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
		},
	}
}

func resourceTencentCloudMysqlBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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
	if backupInfo.EncryptionFlag != nil {
		_ = d.Set("encryption_flag", backupInfo.EncryptionFlag)
	}

	return nil
}

func resourceTencentCloudMysqlBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	return nil
}

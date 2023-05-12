/*
Provides a resource to create a sqlserver full_backup_migration

Example Usage

```hcl
resource "tencentcloud_sqlserver_full_backup_migration" "full_backup_migration" {
  instance_id = "mssql-i1z41iwd"
  recovery_type = "FULL"
  upload_type = "COS_URL"
  migration_name = "test_migration"
  backup_files = []
}
```

Import

sqlserver full_backup_migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_full_backup_migration.full_backup_migration full_backup_migration_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverFullBackupMigration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverFullBackupMigrationCreate,
		Read:   resourceTencentCloudSqlserverFullBackupMigrationRead,
		Update: resourceTencentCloudSqlserverFullBackupMigrationUpdate,
		Delete: resourceTencentCloudSqlserverFullBackupMigrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of imported target instance.",
			},
			"recovery_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Migration task restoration type. FULL: full backup restoration, FULL_LOG: full backup and transaction log restoration, FULL_DIFF: full backup and differential backup restoration.",
			},
			"upload_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup upload type. COS_URL: the backup is stored in users Cloud Object Storage, with URL provided. COS_UPLOAD: the backup is stored in the application Cloud Object Storage and needs to be uploaded by the user.",
			},
			"migration_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task name.",
			},
			"backup_files": {
				Optional:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "If the UploadType is COS_URL, fill in the URL here. If the UploadType is COS_UPLOAD, fill in the name of the backup file here. Only 1 backup file is supported, but a backup file can involve multiple databases.",
			},
			"backup_migration_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "migration id.",
			},
			"backup_migration_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "backup migration set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "app id.",
						},
						"backup_files": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "backup files list.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "end time.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"is_recovery": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether it is the final recovery, the field of the full import task is empty.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "msg.",
						},
						"migration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "migration id.",
						},
						"migration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "migration name.",
						},
						"recovery_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "recovery type.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start time.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Migration task status, 2-created, 7-full import, 8-waiting for increment, 9-import successful, 10-import failed, 12-incremental import.",
						},
						"upload_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup user upload type, COS_URL - the backup is placed on the user's object storage, and the URL is provided. COS_UPLOAD-The backup is placed on the object storage of the business, and the user uploads it.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverFullBackupMigrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_full_backup_migration.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		request           = sqlserver.NewCreateBackupMigrationRequest()
		response          = sqlserver.NewCreateBackupMigrationResponse()
		instanceId        string
		backupMigrationId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("recovery_type"); ok {
		request.RecoveryType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upload_type"); ok {
		request.UploadType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("migration_name"); ok {
		request.MigrationName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_files"); ok {
		backupFilesSet := v.(*schema.Set).List()
		for i := range backupFilesSet {
			backupFiles := backupFilesSet[i].(string)
			request.BackupFiles = append(request.BackupFiles, &backupFiles)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CreateBackupMigration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver full backup migration %s not exists", instanceId)
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver fullBackupMigration failed, reason:%+v", logId, err)
		return err
	}

	backupMigrationId = *response.Response.BackupMigrationId
	d.SetId(instanceId)
	_ = d.Set("backup_migration_id", backupMigrationId)

	return resourceTencentCloudSqlserverFullBackupMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverFullBackupMigrationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_full_backup_migration.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		service           = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId        = d.Id()
		backupMigrationId = d.Get("backup_migration_id").(string)
	)

	fullBackupMigration, err := service.DescribeSqlserverFullBackupMigrationById(ctx, instanceId, backupMigrationId)
	if err != nil {
		return err
	}

	if fullBackupMigration == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverFullBackupMigration` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	list := make([]map[string]interface{}, 0)
	var infoMap = map[string]interface{}{}
	if fullBackupMigration.AppId != nil {
		infoMap["app_id"] = fullBackupMigration.AppId
	}
	if fullBackupMigration.BackupFiles != nil {
		infoMap["backup_files"] = fullBackupMigration.BackupFiles
	}
	if fullBackupMigration.CreateTime != nil {
		infoMap["create_time"] = fullBackupMigration.CreateTime
	}
	if fullBackupMigration.EndTime != nil {
		infoMap["end_time"] = fullBackupMigration.EndTime
	}
	if fullBackupMigration.InstanceId != nil {
		infoMap["instance_id"] = fullBackupMigration.InstanceId
	}
	if fullBackupMigration.IsRecovery != nil {
		infoMap["is_recovery"] = fullBackupMigration.IsRecovery
	}
	if fullBackupMigration.Message != nil {
		infoMap["message"] = fullBackupMigration.Message
	}
	if fullBackupMigration.MigrationId != nil {
		infoMap["migration_id"] = fullBackupMigration.MigrationId
	}
	if fullBackupMigration.MigrationName != nil {
		infoMap["migration_name"] = fullBackupMigration.MigrationName
	}
	if fullBackupMigration.RecoveryType != nil {
		infoMap["recovery_type"] = fullBackupMigration.RecoveryType
	}
	if fullBackupMigration.Region != nil {
		infoMap["region"] = fullBackupMigration.Region
	}
	if fullBackupMigration.StartTime != nil {
		infoMap["start_time"] = fullBackupMigration.StartTime
	}
	if fullBackupMigration.Status != nil {
		infoMap["status"] = fullBackupMigration.Status
	}
	if fullBackupMigration.UploadType != nil {
		infoMap["upload_type"] = fullBackupMigration.UploadType
	}

	list = append(list, infoMap)
	_ = d.Set("backup_migration_set", list)

	return nil
}

func resourceTencentCloudSqlserverFullBackupMigrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_full_backup_migration.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		request           = sqlserver.NewModifyBackupMigrationRequest()
		instanceId        = d.Id()
		backupMigrationId = d.Get("backup_migration_id").(string)
	)

	request.InstanceId = &instanceId
	request.BackupMigrationId = &backupMigrationId

	if d.HasChange("migration_name") {
		if v, ok := d.GetOk("migration_name"); ok {
			request.MigrationName = helper.String(v.(string))
		}
	}

	if d.HasChange("recovery_type") {
		if v, ok := d.GetOk("recovery_type"); ok {
			request.RecoveryType = helper.String(v.(string))
		}
	}

	if d.HasChange("upload_type") {
		if v, ok := d.GetOk("upload_type"); ok {
			request.UploadType = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyBackupMigration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("update sqlserver full backup migration %s not exists", instanceId)
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver fullBackupMigration failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverFullBackupMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverFullBackupMigrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_full_backup_migration.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		service           = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId        = d.Id()
		backupMigrationId = d.Get("backup_migration_id").(string)
	)

	if err := service.DeleteSqlserverFullBackupMigrationById(ctx, instanceId, backupMigrationId); err != nil {
		return err
	}

	return nil
}

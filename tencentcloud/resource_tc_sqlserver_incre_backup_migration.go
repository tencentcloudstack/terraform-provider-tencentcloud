/*
Provides a resource to create a sqlserver incre_backup_migration

Example Usage

```hcl
resource "tencentcloud_sqlserver_incre_backup_migration" "incre_backup_migration" {
  instance_id = "mssql-i1z41iwd"
  backup_migration_id = "migration_00001"
  backup_files =
  is_recovery = "No"
}
```

Import

sqlserver incre_backup_migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration incre_backup_migration_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSqlserverIncreBackupMigration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverIncreBackupMigrationCreate,
		Read:   resourceTencentCloudSqlserverIncreBackupMigrationRead,
		Update: resourceTencentCloudSqlserverIncreBackupMigrationUpdate,
		Delete: resourceTencentCloudSqlserverIncreBackupMigrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of imported target instance.",
			},

			"backup_migration_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup import task ID, which is returned through the API CreateBackupMigration.",
			},

			"backup_files": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Incremental backup file. If the UploadType of a full backup file is COS_URL, fill in URL here. If the UploadType is COS_UPLOAD, fill in the name of the backup file here. Only 1 backup file is supported, but a backup file can involve multiple databases.",
			},

			"is_recovery": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether restoration is required. No: not required. Yes: required. Not required by default.",
			},
		},
	}
}

func resourceTencentCloudSqlserverIncreBackupMigrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_incre_backup_migration.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewCreateIncrementalMigrationRequest()
		response   = sqlserver.NewCreateIncrementalMigrationResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_migration_id"); ok {
		request.BackupMigrationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_files"); ok {
		backupFilesSet := v.(*schema.Set).List()
		for i := range backupFilesSet {
			backupFiles := backupFilesSet[i].(string)
			request.BackupFiles = append(request.BackupFiles, &backupFiles)
		}
	}

	if v, ok := d.GetOk("is_recovery"); ok {
		request.IsRecovery = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CreateIncrementalMigration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver increBackupMigration failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudSqlserverIncreBackupMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverIncreBackupMigrationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_incre_backup_migration.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	increBackupMigrationId := d.Id()

	increBackupMigration, err := service.DescribeSqlserverIncreBackupMigrationById(ctx, instanceId)
	if err != nil {
		return err
	}

	if increBackupMigration == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverIncreBackupMigration` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if increBackupMigration.InstanceId != nil {
		_ = d.Set("instance_id", increBackupMigration.InstanceId)
	}

	if increBackupMigration.BackupMigrationId != nil {
		_ = d.Set("backup_migration_id", increBackupMigration.BackupMigrationId)
	}

	if increBackupMigration.BackupFiles != nil {
		_ = d.Set("backup_files", increBackupMigration.BackupFiles)
	}

	if increBackupMigration.IsRecovery != nil {
		_ = d.Set("is_recovery", increBackupMigration.IsRecovery)
	}

	return nil
}

func resourceTencentCloudSqlserverIncreBackupMigrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_incre_backup_migration.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyIncrementalMigrationRequest()

	increBackupMigrationId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "backup_migration_id", "backup_files", "is_recovery"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("backup_migration_id") {
		if v, ok := d.GetOk("backup_migration_id"); ok {
			request.BackupMigrationId = helper.String(v.(string))
		}
	}

	if d.HasChange("backup_files") {
		if v, ok := d.GetOk("backup_files"); ok {
			backupFilesSet := v.(*schema.Set).List()
			for i := range backupFilesSet {
				backupFiles := backupFilesSet[i].(string)
				request.BackupFiles = append(request.BackupFiles, &backupFiles)
			}
		}
	}

	if d.HasChange("is_recovery") {
		if v, ok := d.GetOk("is_recovery"); ok {
			request.IsRecovery = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyIncrementalMigration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver increBackupMigration failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverIncreBackupMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverIncreBackupMigrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_incre_backup_migration.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	increBackupMigrationId := d.Id()

	if err := service.DeleteSqlserverIncreBackupMigrationById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}

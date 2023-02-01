/*
Provides a resource to create a sqlserver backup_migration

Example Usage

```hcl
resource "tencentcloud_sqlserver_backup_migration" "backup_migration" {
  instance_id = ""
  recovery_type = ""
  upload_type = ""
  migration_name = ""
  backup_files =
}
```

Import

sqlserver backup_migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_backup_migration.backup_migration backup_migration_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSqlserverBackupMigration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverBackupMigrationCreate,
		Read:   resourceTencentCloudSqlserverBackupMigrationRead,
		Update: resourceTencentCloudSqlserverBackupMigrationUpdate,
		Delete: resourceTencentCloudSqlserverBackupMigrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Import target instance ID.",
			},

			"recovery_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Migration task recovery type, FULL - full backup recovery, FULL_ LOG - full backup+transaction log recovery, FULL_ DIFF - full backup+differential backup recovery.",
			},

			"upload_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup upload type, COS_ The URL - backup is placed on the user&amp;#39;s object store and provides the URL. COS_ UPLOAD - The backup is placed on the object storage of the business and needs to be uploaded by the user.",
			},

			"migration_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task name.",
			},

			"backup_files": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "UploadType is COS_ Fill in the URL here, COS_ Fill in the name of the backup file here. Only one backup file is supported, but one backup file can contain multiple libraries.",
			},
		},
	}
}

func resourceTencentCloudSqlserverBackupMigrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_backup_migration.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request           = sqlserver.NewCreateBackupMigrationRequest()
		response          = sqlserver.NewCreateBackupMigrationResponse()
		backupMigrationId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver backupMigration failed, reason:%+v", logId, err)
		return err
	}

	backupMigrationId = *response.Response.BackupMigrationId
	d.SetId(helper.String(backupMigrationId))

	return resourceTencentCloudSqlserverBackupMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverBackupMigrationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_backup_migration.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupMigrationId := d.Id()

	backupMigration, err := service.DescribeSqlserverBackupMigrationById(ctx, backupMigrationId)
	if err != nil {
		return err
	}

	if backupMigration == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverBackupMigration` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backupMigration.InstanceId != nil {
		_ = d.Set("instance_id", backupMigration.InstanceId)
	}

	if backupMigration.RecoveryType != nil {
		_ = d.Set("recovery_type", backupMigration.RecoveryType)
	}

	if backupMigration.UploadType != nil {
		_ = d.Set("upload_type", backupMigration.UploadType)
	}

	if backupMigration.MigrationName != nil {
		_ = d.Set("migration_name", backupMigration.MigrationName)
	}

	if backupMigration.BackupFiles != nil {
		_ = d.Set("backup_files", backupMigration.BackupFiles)
	}

	return nil
}

func resourceTencentCloudSqlserverBackupMigrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_backup_migration.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyBackupMigrationRequest()

	backupMigrationId := d.Id()

	request.BackupMigrationId = &backupMigrationId

	immutableArgs := []string{"instance_id", "recovery_type", "upload_type", "migration_name", "backup_files"}

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

	if d.HasChange("migration_name") {
		if v, ok := d.GetOk("migration_name"); ok {
			request.MigrationName = helper.String(v.(string))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyBackupMigration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver backupMigration failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverBackupMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverBackupMigrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_backup_migration.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	backupMigrationId := d.Id()

	if err := service.DeleteSqlserverBackupMigrationById(ctx, backupMigrationId); err != nil {
		return err
	}

	return nil
}

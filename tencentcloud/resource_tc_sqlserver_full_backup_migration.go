/*
Provides a resource to create a sqlserver full_backup_migration

Example Usage

```hcl
resource "tencentcloud_sqlserver_full_backup_migration" "my_migration" {
  instance_id = "mssql-qelbzgwf"
  recovery_type = "FULL"
  upload_type = "COS_URL"
  migration_name = "migration_test"
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
	"strings"

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
				Description: "Backup upload type. COS_URL: the backup is stored in users Cloud Object Storage, with URL provided. COS_UPLOAD: the backup is stored in the applications Cloud Object Storage and needs to be uploaded by the user.",
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
		request.BackupFiles = helper.InterfacesStringsPoint(v.([]interface{}))
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
	d.SetId(strings.Join([]string{instanceId, backupMigrationId}, FILED_SP))

	return resourceTencentCloudSqlserverFullBackupMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverFullBackupMigrationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_full_backup_migration.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	backupMigrationId := idSplit[1]

	fullBackupMigration, err := service.DescribeSqlserverFullBackupMigrationById(ctx, instanceId, backupMigrationId)
	if err != nil {
		return err
	}

	if fullBackupMigration == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverFullBackupMigration` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if fullBackupMigration.InstanceId != nil {
		_ = d.Set("instance_id", fullBackupMigration.InstanceId)
	}

	if fullBackupMigration.RecoveryType != nil {
		_ = d.Set("recovery_type", fullBackupMigration.RecoveryType)
	}

	if fullBackupMigration.UploadType != nil {
		_ = d.Set("upload_type", fullBackupMigration.UploadType)
	}

	if fullBackupMigration.MigrationName != nil {
		_ = d.Set("migration_name", fullBackupMigration.MigrationName)
	}

	if fullBackupMigration.BackupFiles != nil {
		_ = d.Set("backup_files", fullBackupMigration.BackupFiles)
	}

	return nil
}

func resourceTencentCloudSqlserverFullBackupMigrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_full_backup_migration.update")()
	defer inconsistentCheck(d, meta)()

	immutableArgs := []string{"instance_id", "backup_files"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	var (
		logId   = getLogId(contextNil)
		request = sqlserver.NewModifyBackupMigrationRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	backupMigrationId := idSplit[1]

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
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	backupMigrationId := idSplit[1]

	if err := service.DeleteSqlserverFullBackupMigrationById(ctx, instanceId, backupMigrationId); err != nil {
		return err
	}

	return nil
}

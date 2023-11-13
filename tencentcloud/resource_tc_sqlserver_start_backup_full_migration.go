/*
Provides a resource to create a sqlserver start_backup_full_migration

Example Usage

```hcl
resource "tencentcloud_sqlserver_start_backup_full_migration" "start_backup_full_migration" {
  instance_id = "mssql-i1z41iwd"
  backup_migration_id = ""
}
```

Import

sqlserver start_backup_full_migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_start_backup_full_migration.start_backup_full_migration start_backup_full_migration_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSqlserverStartBackupFullMigration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverStartBackupFullMigrationCreate,
		Read:   resourceTencentCloudSqlserverStartBackupFullMigrationRead,
		Delete: resourceTencentCloudSqlserverStartBackupFullMigrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of imported target instance.",
			},

			"backup_migration_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Backup import task ID, returned by the CreateBackupMigration interface.",
			},
		},
	}
}

func resourceTencentCloudSqlserverStartBackupFullMigrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_start_backup_full_migration.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewStartBackupMigrationRequest()
		response   = sqlserver.NewStartBackupMigrationResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_migration_id"); ok {
		request.BackupMigrationId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().StartBackupMigration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate sqlserver startBackupFullMigration failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudSqlserverStartBackupFullMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverStartBackupFullMigrationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_start_backup_full_migration.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSqlserverStartBackupFullMigrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_start_backup_full_migration.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

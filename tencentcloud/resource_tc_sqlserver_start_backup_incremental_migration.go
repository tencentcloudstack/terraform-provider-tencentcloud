/*
Provides a resource to create a sqlserver start_backup_incremental_migration

Example Usage

```hcl
resource "tencentcloud_sqlserver_start_backup_incremental_migration" "start_backup_incremental_migration" {
  instance_id              = "mssql-i1z41iwd"
  backup_migration_id      = "mssql-backup-migration-cg0ffgqt"
  incremental_migration_id = "mssql-incremental-migration-kp7bgv8p"
}
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

func resourceTencentCloudSqlserverStartBackupIncrementalMigration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverStartBackupIncrementalMigrationCreate,
		Read:   resourceTencentCloudSqlserverStartBackupIncrementalMigrationRead,
		Delete: resourceTencentCloudSqlserverStartBackupIncrementalMigrationDelete,

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
			"incremental_migration_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Incremental backup import task ID.",
			},
		},
	}
}

func resourceTencentCloudSqlserverStartBackupIncrementalMigrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_start_backup_incremental_migration.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = sqlserver.NewStartIncrementalMigrationRequest()
		instanceId string
		flowId     uint64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_migration_id"); ok {
		request.BackupMigrationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("incremental_migration_id"); ok {
		request.IncrementalMigrationId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().StartIncrementalMigration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate sqlserver startBackupIncrementalMigration failed, reason:%+v", logId, err)
		return err
	}

	// wait
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, int64(flowId))
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			e = fmt.Errorf("sqlserver startBackupIncrementalMigration instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *result.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver startBackupIncrementalMigration task status is running"))
		}

		if *result.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		}

		if *result.Status == SQLSERVER_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("sqlserver startBackupIncrementalMigration task status is failed"))
		}

		e = fmt.Errorf("sqlserver startBackupIncrementalMigration task status is %v, we won't wait for it finish", *result.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s sqlserver startBackupIncrementalMigration task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverStartBackupIncrementalMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverStartBackupIncrementalMigrationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_start_backup_incremental_migration.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSqlserverStartBackupIncrementalMigrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_start_backup_incremental_migration.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

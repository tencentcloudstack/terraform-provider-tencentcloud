package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverStartBackupFullMigration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverStartBackupFullMigrationCreate,
		Read:   resourceTencentCloudSqlserverStartBackupFullMigrationRead,
		Delete: resourceTencentCloudSqlserverStartBackupFullMigrationDelete,

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
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_start_backup_full_migration.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = sqlserver.NewStartBackupMigrationRequest()
		instanceId string
		flowId     uint64
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("backup_migration_id"); ok {
		request.BackupMigrationId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().StartBackupMigration(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate sqlserver startBackupFullMigration failed, reason:%+v", logId, err)
		return err
	}

	// wait
	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, int64(flowId))
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			e = fmt.Errorf("sqlserver startBackupFullMigration instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *result.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver startBackupFullMigration task status is running"))
		}

		if *result.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		}

		if *result.Status == SQLSERVER_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("sqlserver startBackupFullMigration task status is failed"))
		}

		e = fmt.Errorf("sqlserver startBackupFullMigration task status is %v, we won't wait for it finish", *result.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s sqlserver startBackupFullMigration task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverStartBackupFullMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverStartBackupFullMigrationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_start_backup_full_migration.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSqlserverStartBackupFullMigrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_start_backup_full_migration.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

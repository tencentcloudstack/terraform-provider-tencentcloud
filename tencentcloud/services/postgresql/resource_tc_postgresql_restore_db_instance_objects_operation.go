package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresqlv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlRestoreDbInstanceObjectsOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlRestoreDbInstanceObjectsOperationCreate,
		Read:   resourceTencentCloudPostgresqlRestoreDbInstanceObjectsOperationRead,
		Delete: resourceTencentCloudPostgresqlRestoreDbInstanceObjectsOperationDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "PostgreSQL instance ID, e.g. `postgres-6bwgamo3`.",
			},

			"restore_objects": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of database objects to restore. The restored object name format will be `${original}_bak_${timestamp}`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"backup_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup set ID used for restoration. Exactly one of `backup_set_id` or `restore_target_time` must be specified.",
			},

			"restore_target_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Point-in-time target for restoration (Beijing time), e.g. `2024-04-30 00:20:27`. Exactly one of `backup_set_id` or `restore_target_time` must be specified.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlRestoreDbInstanceObjectsOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_restore_db_instance_objects_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = postgresqlv20170312.NewRestoreDBInstanceObjectsRequest()
		response     = postgresqlv20170312.NewRestoreDBInstanceObjectsResponse()
		dbInstanceId string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dbInstanceId = v.(string)
	}

	if v, ok := d.GetOk("restore_objects"); ok {
		for _, item := range v.([]interface{}) {
			request.RestoreObjects = append(request.RestoreObjects, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("backup_set_id"); ok {
		request.BackupSetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("restore_target_time"); ok {
		request.RestoreTargetTime = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlV20170312Client().RestoreDBInstanceObjectsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("RestoreDBInstanceObjects failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create postgresql restore db instance objects operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}

	taskId := *response.Response.TaskId

	// Poll DescribeTasks until Status == "Success"
	flowRequest := postgresqlv20170312.NewDescribeTasksRequest()
	flowRequest.TaskId = helper.Int64Uint64(taskId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlV20170312Client().DescribeTasksWithContext(ctx, flowRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, flowRequest.GetAction(), flowRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskSet == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeTasks response is nil."))
		}

		if len(result.Response.TaskSet) == 0 {
			return resource.RetryableError(fmt.Errorf("waiting for task to initialize."))
		}

		if result.Response.TaskSet[0].Status == nil {
			return resource.RetryableError(fmt.Errorf("task status is nil, waiting."))
		}

		status := *result.Response.TaskSet[0].Status
		if status == "Success" {
			return nil
		}

		if status == "Fail" || status == "Pause" {
			return resource.NonRetryableError(fmt.Errorf("restore db instance objects task failed with status: %s.", status))
		}

		return resource.RetryableError(fmt.Errorf("restore db instance objects task still running, current status: %s.", status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s wait for postgresql restore db instance objects task failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dbInstanceId)
	return resourceTencentCloudPostgresqlRestoreDbInstanceObjectsOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlRestoreDbInstanceObjectsOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_restore_db_instance_objects_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlRestoreDbInstanceObjectsOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_restore_db_instance_objects_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

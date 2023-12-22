package postgresql

import (
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlDeleteLogBackupOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlDeleteLogBackupOperationCreate,
		Read:   resourceTencentCloudPostgresqlDeleteLogBackupOperationRead,
		Delete: resourceTencentCloudPostgresqlDeleteLogBackupOperationDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"log_backup_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Log backup ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlDeleteLogBackupOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_delete_log_backup_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = postgresql.NewDeleteLogBackupRequest()
		dBInstanceId string
		logBackupId  string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dBInstanceId = v.(string)
	}

	if v, ok := d.GetOk("log_backup_id"); ok {
		request.LogBackupId = helper.String(v.(string))
		logBackupId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DeleteLogBackup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql DeleteLogBackupOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{dBInstanceId, logBackupId}, tccommon.FILED_SP))

	return resourceTencentCloudPostgresqlDeleteLogBackupOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlDeleteLogBackupOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_delete_log_backup_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlDeleteLogBackupOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_delete_log_backup_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

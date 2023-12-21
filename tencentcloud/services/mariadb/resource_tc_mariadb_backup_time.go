package mariadb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMariadbBackupTime() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbBackupTimeCreate,
		Read:   resourceTencentCloudMariadbBackupTimeRead,
		Update: resourceTencentCloudMariadbBackupTimeUpdate,
		Delete: resourceTencentCloudMariadbBackupTimeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},
			"start_backup_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time of daily backup window in the format of `mm:ss`, such as 22:00.",
			},
			"end_backup_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time of daily backup window in the format of `mm:ss`, such as 23:59.",
			},
		},
	}
}

func resourceTencentCloudMariadbBackupTimeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_backup_time.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbBackupTimeUpdate(d, meta)
}

func resourceTencentCloudMariadbBackupTimeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_backup_time.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	backupTime, err := service.DescribeMariadbBackupTimeById(ctx, instanceId)
	if err != nil {
		return err
	}

	if backupTime == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbBackupTime` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backupTime.InstanceId != nil {
		_ = d.Set("instance_id", backupTime.InstanceId)
	}

	if backupTime.StartBackupTime != nil {
		_ = d.Set("start_backup_time", backupTime.StartBackupTime)
	}

	if backupTime.EndBackupTime != nil {
		_ = d.Set("end_backup_time", backupTime.EndBackupTime)
	}

	return nil
}

func resourceTencentCloudMariadbBackupTimeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_backup_time.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mariadb.NewModifyBackupTimeRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId
	if v, ok := d.GetOk("start_backup_time"); ok {
		request.StartBackupTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_backup_time"); ok {
		request.EndBackupTime = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().ModifyBackupTime(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if *result.Response.Status != MODIFY_BACKUPTIME_SUCCESS {
			return resource.NonRetryableError(fmt.Errorf("update mariadb backupTime status is fail"))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update mariadb backupTime failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbBackupTimeRead(d, meta)
}

func resourceTencentCloudMariadbBackupTimeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_backup_time.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

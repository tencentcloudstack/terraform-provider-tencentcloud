package cdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlBackupEncryptionStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlBackupEncryptionStatusCreate,
		Read:   resourceTencentCloudMysqlBackupEncryptionStatusRead,
		Update: resourceTencentCloudMysqlBackupEncryptionStatusUpdate,
		Delete: resourceTencentCloudMysqlBackupEncryptionStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, in the format: cdb-XXXX. Same instance ID as displayed in the ApsaraDB for Console page.",
			},

			"encryption_status": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Whether physical backup encryption is enabled for the instance. Possible values are `on`, `off`.",
			},
		},
	}
}

func resourceTencentCloudMysqlBackupEncryptionStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup_encryption_status.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlBackupEncryptionStatusUpdate(d, meta)
}

func resourceTencentCloudMysqlBackupEncryptionStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup_encryption_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	backupEncryptionStatus, err := service.DescribeMysqlBackupEncryptionStatusById(ctx, instanceId)
	if err != nil {
		return err
	}

	if backupEncryptionStatus == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlBackupEncryptionStatus` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if backupEncryptionStatus.EncryptionStatus != nil {
		_ = d.Set("encryption_status", backupEncryptionStatus.EncryptionStatus)
	}

	return nil
}

func resourceTencentCloudMysqlBackupEncryptionStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup_encryption_status.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := mysql.NewModifyBackupEncryptionStatusRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("encryption_status"); ok {
		request.EncryptionStatus = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ModifyBackupEncryptionStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql backupEncryptionStatus failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlBackupEncryptionStatusRead(d, meta)
}

func resourceTencentCloudMysqlBackupEncryptionStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup_encryption_status.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

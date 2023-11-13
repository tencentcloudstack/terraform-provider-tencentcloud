/*
Provides a resource to create a cdb backup_encryption_status

Example Usage

```hcl
resource "tencentcloud_cdb_backup_encryption_status" "backup_encryption_status" {
  instance_id = "cdb-c1nl9rpv"
  encryption_status = "on"
}
```

Import

cdb backup_encryption_status can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_backup_encryption_status.backup_encryption_status backup_encryption_status_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCdbBackupEncryptionStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbBackupEncryptionStatusCreate,
		Read:   resourceTencentCloudCdbBackupEncryptionStatusRead,
		Update: resourceTencentCloudCdbBackupEncryptionStatusUpdate,
		Delete: resourceTencentCloudCdbBackupEncryptionStatusDelete,
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
				Description: "Whether physical backup encryption is enabled for the instance. Possible values are on, off.",
			},
		},
	}
}

func resourceTencentCloudCdbBackupEncryptionStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_backup_encryption_status.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCdbBackupEncryptionStatusUpdate(d, meta)
}

func resourceTencentCloudCdbBackupEncryptionStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_backup_encryption_status.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupEncryptionStatusId := d.Id()

	backupEncryptionStatus, err := service.DescribeCdbBackupEncryptionStatusById(ctx, instanceId)
	if err != nil {
		return err
	}

	if backupEncryptionStatus == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbBackupEncryptionStatus` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backupEncryptionStatus.InstanceId != nil {
		_ = d.Set("instance_id", backupEncryptionStatus.InstanceId)
	}

	if backupEncryptionStatus.EncryptionStatus != nil {
		_ = d.Set("encryption_status", backupEncryptionStatus.EncryptionStatus)
	}

	return nil
}

func resourceTencentCloudCdbBackupEncryptionStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_backup_encryption_status.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyBackupEncryptionStatusRequest()

	backupEncryptionStatusId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "encryption_status"}

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

	if d.HasChange("encryption_status") {
		if v, ok := d.GetOk("encryption_status"); ok {
			request.EncryptionStatus = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyBackupEncryptionStatus(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb backupEncryptionStatus failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCdbBackupEncryptionStatusRead(d, meta)
}

func resourceTencentCloudCdbBackupEncryptionStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_backup_encryption_status.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

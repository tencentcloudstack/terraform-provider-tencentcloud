/*
Provides a resource to create a cdb remote_backup_config

Example Usage

```hcl
resource "tencentcloud_cdb_remote_backup_config" "remote_backup_config" {
  instance_id = "cdb-c1nl9rpv"
  remote_backup_save = "on"
  remote_binlog_save = "on"
  remote_region =
  expire_days = 7
}
```

Import

cdb remote_backup_config can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_remote_backup_config.remote_backup_config remote_backup_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"log"
)

func resourceTencentCloudCdbRemoteBackupConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbRemoteBackupConfigCreate,
		Read:   resourceTencentCloudCdbRemoteBackupConfigRead,
		Update: resourceTencentCloudCdbRemoteBackupConfigUpdate,
		Delete: resourceTencentCloudCdbRemoteBackupConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.",
			},

			"remote_backup_save": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Remote data backup switch, off - disable remote backup, on - enable remote backup.",
			},

			"remote_binlog_save": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Off-site log backup switch, off - off off-site backup, on-on off-site backup, only when the parameter RemoteBackupSave is on, the RemoteBinlogSave parameter can be set to on.",
			},

			"remote_region": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "User settings off-site backup region list.",
			},

			"expire_days": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Remote backup retention time, in days.",
			},
		},
	}
}

func resourceTencentCloudCdbRemoteBackupConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_remote_backup_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCdbRemoteBackupConfigUpdate(d, meta)
}

func resourceTencentCloudCdbRemoteBackupConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_remote_backup_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	remoteBackupConfigId := d.Id()

	remoteBackupConfig, err := service.DescribeCdbRemoteBackupConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if remoteBackupConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbRemoteBackupConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if remoteBackupConfig.InstanceId != nil {
		_ = d.Set("instance_id", remoteBackupConfig.InstanceId)
	}

	if remoteBackupConfig.RemoteBackupSave != nil {
		_ = d.Set("remote_backup_save", remoteBackupConfig.RemoteBackupSave)
	}

	if remoteBackupConfig.RemoteBinlogSave != nil {
		_ = d.Set("remote_binlog_save", remoteBackupConfig.RemoteBinlogSave)
	}

	if remoteBackupConfig.RemoteRegion != nil {
		_ = d.Set("remote_region", remoteBackupConfig.RemoteRegion)
	}

	if remoteBackupConfig.ExpireDays != nil {
		_ = d.Set("expire_days", remoteBackupConfig.ExpireDays)
	}

	return nil
}

func resourceTencentCloudCdbRemoteBackupConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_remote_backup_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyRemoteBackupConfigRequest()

	remoteBackupConfigId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "remote_backup_save", "remote_binlog_save", "remote_region", "expire_days"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyRemoteBackupConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb remoteBackupConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCdbRemoteBackupConfigRead(d, meta)
}

func resourceTencentCloudCdbRemoteBackupConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_remote_backup_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

/*
Provides a resource to create a cynosdb backup_config

Example Usage

```hcl
resource "tencentcloud_cynosdb_backup_config" "backup_config" {
  cluster_id = &lt;nil&gt;
  backup_time_beg = 3600
  backup_time_end = &lt;nil&gt;
  reserve_duration = &lt;nil&gt;
  backup_freq = &lt;nil&gt;
  backup_type = &lt;nil&gt;
}
```

Import

cynosdb backup_config can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_backup_config.backup_config backup_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCynosdbBackupConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbBackupConfigCreate,
		Read:   resourceTencentCloudCynosdbBackupConfigRead,
		Update: resourceTencentCloudCynosdbBackupConfigUpdate,
		Delete: resourceTencentCloudCynosdbBackupConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of cluster.",
			},

			"backup_time_beg": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Indicates the full backup start time, [0-24*3600], such as 0:00, 1:00, 2:00 are 0, 3600, 7200 respectively.",
			},

			"backup_time_end": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Indicates the full backup end time, [0-24*3600], such as 0:00, 1:00, 2:00 are 0, 3600, 7200 respectively.",
			},

			"reserve_duration": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Indicates the length of time to keep the backup, in seconds, beyond which time it will be cleared, seven days is expressed as 3600*24*7=604800, the maximum is 158112000.",
			},

			"backup_freq": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "This parameter does not currently support modification and does not need to be filled in. Backup frequency, an array of length 7, corresponding to the backup methods from Monday to Sunday, full-full backup, increment-incremental backup.",
			},

			"backup_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "This parameter does not currently support modification and does not need to be filled in. Backup method, logic-logic backup, snapshot-snapshot backup.",
			},
		},
	}
}

func resourceTencentCloudCynosdbBackupConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_backup_config.create")()
	defer inconsistentCheck(d, meta)()

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	d.SetId(clusterId)

	return resourceTencentCloudCynosdbBackupConfigUpdate(d, meta)
}

func resourceTencentCloudCynosdbBackupConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_backup_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupConfigId := d.Id()

	backupConfig, err := service.DescribeCynosdbBackupConfigById(ctx, clusterId)
	if err != nil {
		return err
	}

	if backupConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbBackupConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backupConfig.ClusterId != nil {
		_ = d.Set("cluster_id", backupConfig.ClusterId)
	}

	if backupConfig.BackupTimeBeg != nil {
		_ = d.Set("backup_time_beg", backupConfig.BackupTimeBeg)
	}

	if backupConfig.BackupTimeEnd != nil {
		_ = d.Set("backup_time_end", backupConfig.BackupTimeEnd)
	}

	if backupConfig.ReserveDuration != nil {
		_ = d.Set("reserve_duration", backupConfig.ReserveDuration)
	}

	if backupConfig.BackupFreq != nil {
		_ = d.Set("backup_freq", backupConfig.BackupFreq)
	}

	if backupConfig.BackupType != nil {
		_ = d.Set("backup_type", backupConfig.BackupType)
	}

	return nil
}

func resourceTencentCloudCynosdbBackupConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_backup_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyBackupConfigRequest()

	backupConfigId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_id", "backup_time_beg", "backup_time_end", "reserve_duration", "backup_freq", "backup_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("backup_time_beg") {
		if v, ok := d.GetOkExists("backup_time_beg"); ok {
			request.BackupTimeBeg = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("backup_time_end") {
		if v, ok := d.GetOkExists("backup_time_end"); ok {
			request.BackupTimeEnd = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("reserve_duration") {
		if v, ok := d.GetOkExists("reserve_duration"); ok {
			request.ReserveDuration = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("backup_freq") {
		if v, ok := d.GetOk("backup_freq"); ok {
			backupFreqSet := v.(*schema.Set).List()
			for i := range backupFreqSet {
				backupFreq := backupFreqSet[i].(string)
				request.BackupFreq = append(request.BackupFreq, &backupFreq)
			}
		}
	}

	if d.HasChange("backup_type") {
		if v, ok := d.GetOk("backup_type"); ok {
			request.BackupType = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyBackupConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb backupConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbBackupConfigRead(d, meta)
}

func resourceTencentCloudCynosdbBackupConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_backup_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

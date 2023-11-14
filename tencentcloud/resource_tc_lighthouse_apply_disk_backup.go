/*
Provides a resource to create a lighthouse apply_disk_backup

Example Usage

```hcl
resource "tencentcloud_lighthouse_apply_disk_backup" "apply_disk_backup" {
  disk_id = "lhdisk-123456"
  disk_backup_id = "lhbak-1234556"
}
```

Import

lighthouse apply_disk_backup can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_apply_disk_backup.apply_disk_backup apply_disk_backup_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudLighthouseApplyDiskBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseApplyDiskBackupCreate,
		Read:   resourceTencentCloudLighthouseApplyDiskBackupRead,
		Delete: resourceTencentCloudLighthouseApplyDiskBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Disk ID.",
			},

			"disk_backup_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Disk backup ID.",
			},
		},
	}
}

func resourceTencentCloudLighthouseApplyDiskBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_apply_disk_backup.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = lighthouse.NewApplyDiskBackupRequest()
		response     = lighthouse.NewApplyDiskBackupResponse()
		diskBackupId string
	)
	if v, ok := d.GetOk("disk_id"); ok {
		request.DiskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("disk_backup_id"); ok {
		diskBackupId = v.(string)
		request.DiskBackupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ApplyDiskBackup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse applyDiskBackup failed, reason:%+v", logId, err)
		return err
	}

	diskBackupId = *response.Response.DiskBackupId
	d.SetId(diskBackupId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseApplyDiskBackupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseApplyDiskBackupRead(d, meta)
}

func resourceTencentCloudLighthouseApplyDiskBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_apply_disk_backup.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseApplyDiskBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_apply_disk_backup.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

/*
Provides a resource to create a cbs disk_backup

Example Usage

```hcl
resource "tencentcloud_cbs_disk_backup" "disk_backup" {
  disk_id = "disk-xxx"
  disk_backup_name = "xxx"
}
```

Import

cbs disk_backup can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_disk_backup.disk_backup disk_backup_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCbsDiskBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsDiskBackupCreate,
		Read:   resourceTencentCloudCbsDiskBackupRead,
		Delete: resourceTencentCloudCbsDiskBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the original cloud disk of the backup point, which can be queried through the DescribeDisks API.",
			},

			"disk_backup_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Backup point name.",
			},
		},
	}
}

func resourceTencentCloudCbsDiskBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_disk_backup.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = cbs.NewCreateDiskBackupRequest()
		response     = cbs.NewCreateDiskBackupResponse()
		diskBackupId string
	)
	if v, ok := d.GetOk("disk_id"); ok {
		request.DiskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("disk_backup_name"); ok {
		request.DiskBackupName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().CreateDiskBackup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs DiskBackup failed, reason:%+v", logId, err)
		return err
	}

	diskBackupId = *response.Response.DiskBackupId
	d.SetId(diskBackupId)

	service := CbsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"CREATING"}, 10*readRetryTimeout, time.Second, service.CbsDiskBackupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCbsDiskBackupRead(d, meta)
}

func resourceTencentCloudCbsDiskBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_disk_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CbsService{client: meta.(*TencentCloudClient).apiV3Conn}

	diskBackupId := d.Id()

	DiskBackup, err := service.DescribeCbsDiskBackupById(ctx, diskBackupId)
	if err != nil {
		return err
	}

	if DiskBackup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CbsDiskBackup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if DiskBackup.DiskId != nil {
		_ = d.Set("disk_id", DiskBackup.DiskId)
	}

	if DiskBackup.DiskBackupName != nil {
		_ = d.Set("disk_backup_name", DiskBackup.DiskBackupName)
	}

	return nil
}

func resourceTencentCloudCbsDiskBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_disk_backup.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CbsService{client: meta.(*TencentCloudClient).apiV3Conn}
	diskBackupId := d.Id()

	if err := service.DeleteCbsDiskBackupById(ctx, diskBackupId); err != nil {
		return err
	}

	return nil
}

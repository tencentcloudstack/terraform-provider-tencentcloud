package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudLighthouseDiskBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseDiskBackupCreate,
		Read:   resourceTencentCloudLighthouseDiskBackupRead,
		Update: resourceTencentCloudLighthouseDiskBackupUpdate,
		Delete: resourceTencentCloudLighthouseDiskBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Disk ID. Only data disks are supported to create disk backup.",
			},

			"disk_backup_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Disk backup name. The maximum length is 90 characters.",
			},
		},
	}
}

func resourceTencentCloudLighthouseDiskBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_disk_backup.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = lighthouse.NewCreateDiskBackupRequest()
		response     = lighthouse.NewCreateDiskBackupResponse()
		diskBackupId string
	)
	if v, ok := d.GetOk("disk_id"); ok {
		request.DiskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("disk_backup_name"); ok {
		request.DiskBackupName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().CreateDiskBackup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse diskBackup failed, reason:%+v", logId, err)
		return err
	}

	diskBackupId = *response.Response.DiskBackupId
	d.SetId(diskBackupId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"NORMAL"}, 20*readRetryTimeout, time.Second, service.LighthouseDiskBackupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseDiskBackupRead(d, meta)
}

func resourceTencentCloudLighthouseDiskBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_disk_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	diskBackupId := d.Id()

	diskBackup, err := service.DescribeLighthouseDiskBackupById(ctx, diskBackupId)
	if err != nil {
		return err
	}

	if diskBackup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseDiskBackup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if diskBackup.DiskId != nil {
		_ = d.Set("disk_id", diskBackup.DiskId)
	}

	if diskBackup.DiskBackupName != nil {
		_ = d.Set("disk_backup_name", diskBackup.DiskBackupName)
	}

	return nil
}

func resourceTencentCloudLighthouseDiskBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_disk_backup.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := lighthouse.NewModifyDiskBackupsAttributeRequest()

	diskBackupId := d.Id()

	request.DiskBackupIds = []*string{&diskBackupId}

	if d.HasChange("disk_backup_name") {
		if v, ok := d.GetOk("disk_backup_name"); ok {
			request.DiskBackupName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyDiskBackupsAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update lighthouse diskBackup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLighthouseDiskBackupRead(d, meta)
}

func resourceTencentCloudLighthouseDiskBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_disk_backup.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	diskBackupId := d.Id()

	if err := service.DeleteLighthouseDiskBackupById(ctx, diskBackupId); err != nil {
		return err
	}

	return nil
}

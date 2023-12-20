package cbs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func ResourceTencentCloudCbsDiskBackup() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_disk_backup.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		diskId         string
		diskBackupName string
	)
	if v, ok := d.GetOk("disk_id"); ok {
		diskId = v.(string)
	}

	if v, ok := d.GetOk("disk_backup_name"); ok {
		diskBackupName = v.(string)
	}

	service := CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	diskBackupId, err := service.CreateDiskBackup(ctx, diskId, diskBackupName)
	if err != nil {
		return nil
	}
	d.SetId(diskBackupId)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		diskBackup, e := service.DescribeCbsDiskBackupById(ctx, diskBackupId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if *diskBackup.DiskBackupState != "NORMAL" {
			return resource.RetryableError(fmt.Errorf("DiskBackupState not ready"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs DiskBackup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCbsDiskBackupRead(d, meta)
}

func resourceTencentCloudCbsDiskBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_disk_backup.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_disk_backup.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	diskBackupId := d.Id()

	if err := service.DeleteCbsDiskBackupById(ctx, diskBackupId); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		diskBackup, e := service.DescribeCbsDiskBackupById(ctx, diskBackupId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if diskBackup == nil {
			return nil
		}
		return resource.RetryableError(errors.New("Disk backup still deleting"))
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cbs DiskBackup failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

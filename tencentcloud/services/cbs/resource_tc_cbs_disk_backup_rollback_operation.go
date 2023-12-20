package cbs

import (
	"context"
	"fmt"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func ResourceTencentCloudCbsDiskBackupRollbackOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsDiskBackupRollbackOperationCreate,
		Read:   resourceTencentCloudCbsDiskBackupRollbackOperationRead,
		Delete: resourceTencentCloudCbsDiskBackupRollbackOperationDelete,

		Schema: map[string]*schema.Schema{
			"disk_backup_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud disk backup point ID.",
			},
			"disk_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud disk backup point original cloud disk ID.",
			},
			"is_rollback_completed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the rollback is completed. `true` meaing rollback completed, `false` meaning still rollbacking.",
			},
		},
	}
}

func resourceTencentCloudCbsDiskBackupRollbackOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_disk_backup_rollback_operation.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	diskBackupId := d.Get("disk_backup_id").(string)
	diskId := d.Get("disk_id").(string)

	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := cbsService.ApplyDiskBackup(ctx, diskBackupId, diskId); err != nil {
		return err
	}
	// deal with state sync delay
	time.Sleep(time.Second * 1)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		disk, e := cbsService.DescribeDiskById(ctx, diskId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if *disk.Rollbacking {
			return resource.RetryableError(errors.New("Disk still rollbacking"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(diskBackupId + tccommon.FILED_SP + diskId)

	return resourceTencentCloudCbsDiskBackupRollbackOperationRead(d, meta)
}

func resourceTencentCloudCbsDiskBackupRollbackOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_disk_backup_rollback_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	diskId := idSplit[1]

	disk, err := cbsService.DescribeDiskById(ctx, diskId)
	if err != nil {
		return err
	}
	_ = d.Set("is_rollback_completed", !*disk.Rollbacking)
	return nil
}

func resourceTencentCloudCbsDiskBackupRollbackOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_disk_backup_rollback_operation.delete")()

	return nil
}

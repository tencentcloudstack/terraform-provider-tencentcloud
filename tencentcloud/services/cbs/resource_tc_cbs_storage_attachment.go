package cbs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
)

func ResourceTencentCloudCbsStorageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsStorageAttachmentCreate,
		Read:   resourceTencentCloudCbsStorageAttachmentRead,
		Delete: resourceTencentCloudCbsStorageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the mounted CBS.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CVM instance.",
			},
		},
	}
}

func resourceTencentCloudCbsStorageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_attachment.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Get("storage_id").(string)
	instanceId := d.Get("instance_id").(string)

	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := cbsService.AttachDisk(ctx, storageId, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage attach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(storageId)

	err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		storage, e := cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if storage == nil {
			return resource.NonRetryableError(fmt.Errorf("cbs storage is nil"))
		}
		if *storage.DiskState == CBS_STORAGE_STATUS_ATTACHING || *storage.DiskState == CBS_STORAGE_STATUS_UNATTACHED {
			return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
		}
		if *storage.DiskState == CBS_STORAGE_STATUS_ATTACHED {
			return nil
		}
		e = fmt.Errorf("cbs storage status is %s, we won't wait for it finish.", *storage.DiskState)
		return resource.NonRetryableError(e)
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage attach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudCbsStorageAttachmentRead(d, meta)
}

func resourceTencentCloudCbsStorageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var storage *cbs.Disk
	var e error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		storage, e = cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s describe cbs storage attach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if storage == nil || !*storage.Attached {
		d.SetId("")
		return nil
	}
	_ = d.Set("storage_id", storage.DiskId)
	_ = d.Set("instance_id", storage.InstanceId)

	return nil
}

func resourceTencentCloudCbsStorageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_attachment.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	instanceId := ""
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		storage, e := cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if storage == nil || !*storage.Attached {
			log.Printf("[DEBUG]%s disk id %s is not attached", logId, storageId)
			return nil
		}
		instanceId = *storage.InstanceId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s describe cbs storage attach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if instanceId == "" {
		return nil
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := cbsService.DetachDisk(ctx, storageId, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage detach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		storage, e := cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if storage == nil || !*storage.Attached {
			return nil
		}
		if *storage.Attached {
			return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
		}
		return resource.NonRetryableError(fmt.Errorf("cbs storage status is %s, we won't wait for it finish.", *storage.DiskState))
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage detach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}

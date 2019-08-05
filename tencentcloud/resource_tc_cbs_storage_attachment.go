/*
Provides a CBS storage attachment resource.

Example Usage

```hcl
resource "tencentcloud_cbs_storage_attachment" "attachment" {
  storage_id  = "disk-kdt0sq6m"
  instance_id = "ins-jqlegd42"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudCbsStorageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsStorageAttachmentCreate,
		Read:   resourceTencentCloudCbsStorageAttachmentRead,
		Delete: resourceTencentCloudCbsStorageAttachmentDelete,

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
	defer LogElapsed("resource.tencentcloud_cbs_storage_attachment.create")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	storageId := d.Get("storage_id").(string)
	instanceId := d.Get("instance_id").(string)

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		e := cbsService.AttachDisk(ctx, storageId, instanceId)
		if e != nil {
			if e.(*errors.TencentCloudSDKError).GetCode() == "ResourceBusy" || e.(*errors.TencentCloudSDKError).GetCode() == "InvalidInstance.NotSupported" {
				return resource.RetryableError(fmt.Errorf("cbs storage busy: %s", e.Error()))
			}
			return resource.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage attach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(storageId)

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		storage, e := cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return resource.NonRetryableError(e)
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
	defer LogElapsed("resource.tencentcloud_cbs_storage_attachment.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	storage, err := cbsService.DescribeDiskById(ctx, storageId)
	if err != nil {
		return err
	}

	if *storage.Attached == false {
		log.Printf("[DEBUG]%s, disk id %s is not attached", logId, storageId)
		d.SetId("")
	}

	d.Set("storage_id", storage.DiskId)
	d.Set("instance_id", storage.InstanceId)

	return nil
}

func resourceTencentCloudCbsStorageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("resource.tencentcloud_cbs_storage_attachment.delete")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	storage, err := cbsService.DescribeDiskById(ctx, storageId)
	if err != nil {
		return err
	}
	if *storage.Attached == false {
		log.Printf("[DEBUG]%s, disk id %s is not attached", logId, storageId)
		return nil
	}
	instanceId := *storage.InstanceId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		e := cbsService.DetachDisk(ctx, storageId, instanceId)
		if e != nil {
			if e.(*errors.TencentCloudSDKError).GetCode() == "ResourceBusy" || e.(*errors.TencentCloudSDKError).GetCode() == "InvalidInstance.NotSupported" {
				return resource.RetryableError(fmt.Errorf("cbs storage busy: %s", e.Error()))
			}
			return resource.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage detach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		storage, e := cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if *storage.Attached == true {
			return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
		}
		if *storage.Attached == false {
			return nil
		}
		e = fmt.Errorf("cbs storage status is %s, we won't wait for it finish.", *storage.DiskState)
		return resource.NonRetryableError(e)
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage detach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}

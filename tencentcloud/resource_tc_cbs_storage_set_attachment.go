/*
Provides a CBS storage set attachment resource.

Example Usage

```hcl
resource "tencentcloud_cbs_storage_set_attachment" "attachment" {
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

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
)

func resourceTencentCloudCbsStorageSetAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsStorageSetAttachmentCreate,
		Read:   resourceTencentCloudCbsStorageSetAttachmentRead,
		Delete: resourceTencentCloudCbsStorageSetAttachmentDelete,
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

func resourceTencentCloudCbsStorageSetAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage_set_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	storageId := d.Get("storage_id").(string)
	instanceId := d.Get("instance_id").(string)

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := cbsService.AttachDisk(ctx, storageId, instanceId)
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && IsContains(CVM_RETRYABLE_ERROR, ee.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("cbs attach error: %s, retrying", ee.Error()))
			}
			return resource.NonRetryableError(ee)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage attach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(storageId)

	return nil
}

func resourceTencentCloudCbsStorageSetAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage_set_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var storage *cbs.Disk
	var errRet error
	storage, errRet = cbsService.DescribeDiskById(ctx, storageId)
	if errRet != nil {
		log.Printf("[CRITAL]%s describe cbs storage attach failed, reason:%s\n ", logId, errRet.Error())
		return errRet
	}

	if storage == nil || !*storage.Attached {
		d.SetId("")
		return nil
	}
	_ = d.Set("storage_id", storage.DiskId)
	_ = d.Set("instance_id", storage.InstanceId)

	return nil
}

func resourceTencentCloudCbsStorageSetAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage_set_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	storageId := d.Id()
	instanceId := d.Get("instance_id").(string)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := cbsService.DetachDisk(ctx, storageId, instanceId)
		if errRet != nil {
			log.Printf("[CRITAL][detach disk]%s api[%s] fail, reason[%s]\n",
				logId, "detach", errRet.Error())
			e, ok := errRet.(*sdkErrors.TencentCloudSDKError)
			if ok && IsContains(CVM_RETRYABLE_ERROR, e.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("[detach]disk detach error: %s, retrying", e.Error()))
			}
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage detach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}

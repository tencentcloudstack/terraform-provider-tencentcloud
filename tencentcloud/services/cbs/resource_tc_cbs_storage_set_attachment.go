package cbs

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
)

func ResourceTencentCloudCbsStorageSetAttachment() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set_attachment.create")()

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
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, ee.Code) {
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
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
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
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set_attachment.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	storageId := d.Id()
	instanceId := d.Get("instance_id").(string)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cbsService.DetachDisk(ctx, storageId, instanceId)
		if errRet != nil {
			log.Printf("[CRITAL][detach disk]%s api[%s] fail, reason[%s]\n",
				logId, "detach", errRet.Error())
			e, ok := errRet.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, e.Code) {
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

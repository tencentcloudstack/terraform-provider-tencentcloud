package cos

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudCosObjectAbortMultipartUploadOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosObjectAbortMultipartUploadOperationCreate,
		Read:   resourceTencentCloudCosObjectAbortMultipartUploadOperationRead,
		Delete: resourceTencentCloudCosObjectAbortMultipartUploadOperationDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Bucket.",
			},
			"key": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Object key.",
			},
			"upload_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Multipart uploaded id.",
			},
		},
	}
}

func resourceTencentCloudCosObjectAbortMultipartUploadOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_abort_multipart_upload_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	uploadId := d.Get("upload_id").(string)
	_, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Object.AbortMultipartUpload(ctx, key, uploadId)
	if err != nil {
		log.Printf("[CRITAL]%s AbortMultipartUpload failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(uploadId)

	return resourceTencentCloudCosObjectAbortMultipartUploadOperationRead(d, meta)
}

func resourceTencentCloudCosObjectAbortMultipartUploadOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_abort_multipart_upload_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCosObjectAbortMultipartUploadOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_abort_multipart_upload_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

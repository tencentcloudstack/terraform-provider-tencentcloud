package cos

import (
	"context"
	"io"
	"log"
	"os"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudCosObjectDownloadOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosObjectDownloadOperationCreate,
		Read:   resourceTencentCloudCosObjectDownloadOperationRead,
		Delete: resourceTencentCloudCosObjectDownloadOperationDelete,

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
			"download_path": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Download path.",
			},
		},
	}
}

func resourceTencentCloudCosObjectDownloadOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_download_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	downloadPath := d.Get("download_path").(string)
	resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Object.Get(ctx, key, nil)
	if err != nil {
		log.Printf("[CRITAL]%s object download failed, reason:%+v", logId, err)
		return err
	}
	defer resp.Body.Close()

	fd, err := os.OpenFile(downloadPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}

	_, err = io.Copy(fd, resp.Body)
	fd.Close()
	if err != nil {
		return err
	}

	d.SetId(bucket + tccommon.FILED_SP + key)

	return resourceTencentCloudCosObjectDownloadOperationRead(d, meta)
}

func resourceTencentCloudCosObjectDownloadOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_download_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCosObjectDownloadOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_download_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

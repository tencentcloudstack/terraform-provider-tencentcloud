package tencentcloud

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudCosObjectDownloadOperation() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_cos_object_download_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	downloadPath := d.Get("download_path").(string)
	resp, err := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Object.Get(ctx, key, nil)
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

	d.SetId(bucket + FILED_SP + key)

	return resourceTencentCloudCosObjectDownloadOperationRead(d, meta)
}

func resourceTencentCloudCosObjectDownloadOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_object_download_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCosObjectDownloadOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_object_download_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

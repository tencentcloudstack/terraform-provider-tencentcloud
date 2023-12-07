package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudCosObjectCopyOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosObjectCopyOperationCreate,
		Read:   resourceTencentCloudCosObjectCopyOperationRead,
		Delete: resourceTencentCloudCosObjectCopyOperationDelete,

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
			"source_url": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Object key.",
			},
		},
	}
}

func resourceTencentCloudCosObjectCopyOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_object_copy_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	sourceURL := d.Get("source_url").(string)
	_, _, err := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Object.Copy(ctx, key, sourceURL, nil)
	if err != nil {
		log.Printf("[CRITAL]%s Restore failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(bucket + FILED_SP + key)

	return resourceTencentCloudCosObjectCopyOperationRead(d, meta)
}

func resourceTencentCloudCosObjectCopyOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_object_copy_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCosObjectCopyOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_object_copy_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

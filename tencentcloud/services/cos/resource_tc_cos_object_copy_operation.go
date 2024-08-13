package cos

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudCosObjectCopyOperation() *schema.Resource {
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
			"cdc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "CDC cluster ID.",
			},
		},
	}
}

func resourceTencentCloudCosObjectCopyOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_copy_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	sourceURL := d.Get("source_url").(string)
	cdcId := d.Get("cdc_id").(string)
	_, _, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket, cdcId).Object.Copy(ctx, key, sourceURL, nil)
	if err != nil {
		log.Printf("[CRITAL]%s Restore failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{bucket, key}, tccommon.FILED_SP))

	return resourceTencentCloudCosObjectCopyOperationRead(d, meta)
}

func resourceTencentCloudCosObjectCopyOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_copy_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCosObjectCopyOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_copy_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

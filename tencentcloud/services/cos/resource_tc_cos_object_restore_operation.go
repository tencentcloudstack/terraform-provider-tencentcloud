package cos

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func ResourceTencentCloudCosObjectRestoreOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosObjectRestoreOperationCreate,
		Read:   resourceTencentCloudCosObjectRestoreOperationRead,
		Delete: resourceTencentCloudCosObjectRestoreOperationDelete,

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
			"tier": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
				Description: "when restoring, Tier can be specified as the supported recovery model.\n" +
					"There are three recovery models for recovering archived storage type data, which are:\n" +
					"- Expedited: quick retrieval mode, and the recovery task can be completed in 1-5 minutes.\n" +
					"- Standard: standard retrieval mode. Recovery task is completed within 3-5 hours.\n" +
					"- Bulk: batch retrieval mode, and the recovery task is completed within 5-12 hours.\n" +
					"For deep recovery archive storage type data, there are two recovery models, which are:\n" +
					"- Standard: standard retrieval mode, recovery time is 12-24 hours.\n" +
					"- Bulk: batch retrieval mode, recovery time is 24-48 hours.",
			},
			"days": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Specifies the valid duration of the restored temporary copy in days.",
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

func resourceTencentCloudCosObjectRestoreOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_restore_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	tier := d.Get("tier").(string)
	days := d.Get("days").(int)
	cdcId := d.Get("cdc_id").(string)
	opt := &cos.ObjectRestoreOptions{
		Days: days,
		Tier: &cos.CASJobParameters{
			Tier: tier,
		},
	}

	_, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket, cdcId).Object.PostRestore(ctx, key, opt)
	if err != nil {
		log.Printf("[CRITAL]%s Restore failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{bucket, key}, tccommon.FILED_SP))

	return resourceTencentCloudCosObjectRestoreOperationRead(d, meta)
}

func resourceTencentCloudCosObjectRestoreOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_restore_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCosObjectRestoreOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_object_restore_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

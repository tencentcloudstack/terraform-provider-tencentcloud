package ci

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTencentCloudCIOriginalImageProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCIOriginalImageProtectionCreate,
		Read:   resourceTencentCloudCIOriginalImageProtectionRead,
		Update: resourceTencentCloudCIOriginalImageProtectionUpdate,
		Delete: resourceTencentCloudCIOriginalImageProtectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The name of a bucket, the format should be [custom name]-[appid], for example `mycos-1258798060`.",
				ValidateFunc: tccommon.ValidateCosBucketName,
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Whether original image protection is set, options: on/off.",
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
		},
	}
}

func resourceTencentCloudCIOriginalImageProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_original_image_protection.create")()

	d.SetId(d.Get("bucket").(string))
	return resourceTencentCloudCIOriginalImageProtectionUpdate(d, meta)
}

func resourceTencentCloudCIOriginalImageProtectionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_original_image_protection.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()
	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	res, err := service.GetCiOriginalImageProtectionById(ctx, bucket)

	if err != nil {
		return err
	}
	log.Printf("[DEBUG] bucket=%s status=%s", bucket, res.OriginProtectStatus)
	_ = d.Set("bucket", bucket)
	_ = d.Set("status", res.OriginProtectStatus)
	return nil
}

func resourceTencentCloudCIOriginalImageProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_original_image_protection.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()
	if d.HasChange("status") {
		var err error

		service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		newStatus := d.Get("status")
		if newStatus == "on" {
			err = service.OpenCiOriginalImageProtectionById(ctx, bucket)
		} else {
			err = service.CloseCiOriginalImageProtectionById(ctx, bucket)
		}

		if err != nil {
			return err
		}
		return resourceTencentCloudCIOriginalImageProtectionRead(d, meta)
	}
	return nil
}

func resourceTencentCloudCIOriginalImageProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_original_image_protection.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Get("bucket").(string)
	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := service.CloseCiOriginalImageProtectionById(ctx, bucket)

	if err != nil {
		return err
	}
	return nil
}

/*
Manage original image protection functionality

Example Usage

```hcl

resource "tencentcloud_ci_original_image_protection" "foo" {
	bucket = "examplebucket-1250000000"
	status = "on"
}

```

Import

Resource original image protection can be imported using the id, e.g.

```
$ terraform import tencentcloud_ci_original_image_protection.example examplebucket-1250000000
```

*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTencentCloudCIOriginalImageProtection() *schema.Resource {
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
				ValidateFunc: validateCosBucketName,
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
	defer logElapsed("resource.tencentcloud_ci_original_image_protection.create")()

	d.SetId(d.Get("bucket").(string))
	return resourceTencentCloudCIOriginalImageProtectionUpdate(d, meta)
}

func resourceTencentCloudCIOriginalImageProtectionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_original_image_protection.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()
	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
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
	defer logElapsed("resource.tencentcloud_ci_original_image_protection.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()
	if d.HasChange("status") {
		var err error

		service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
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
	defer logElapsed("resource.tencentcloud_ci_original_image_protection.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Get("bucket").(string)
	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := service.CloseCiOriginalImageProtectionById(ctx, bucket)

	if err != nil {
		return err
	}
	return nil
}

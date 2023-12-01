/*
Manage Guetzli compression functionality

Example Usage

```hcl

resource "tencentcloud_ci_guetzli" "foo" {
	bucket = "examplebucket-1250000000"
	status = "on"
}

```

Import

Resource guetzli can be imported using the id, e.g.

```
$ terraform import tencentcloud_ci_guetzli.example examplebucket-1250000000
```

*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTencentCloudCIGuetzli() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCIGuetzliCreate,
		Read:   resourceTencentCloudCIGuetzliRead,
		Update: resourceTencentCloudCIGuetzliUpdate,
		Delete: resourceTencentCloudCIGuetzliDelete,
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
				Description:  "Whether Guetzli is set, options: on/off.",
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
		},
	}
}

func resourceTencentCloudCIGuetzliCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_guetzli.create")()

	d.SetId(d.Get("bucket").(string))
	return resourceTencentCloudCIGuetzliUpdate(d, meta)
}

func resourceTencentCloudCIGuetzliRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_guetzli.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()
	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	res, err := service.GetCiGuetzliById(ctx, bucket)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] bucket=%s status=%s", bucket, res.GuetzliStatus)
	_ = d.Set("bucket", bucket)
	_ = d.Set("status", res.GuetzliStatus)
	return nil
}

func resourceTencentCloudCIGuetzliUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_guetzli.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()
	if d.HasChange("status") {
		var err error
		service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

		newStatus := d.Get("status")
		if newStatus == "on" {
			err = service.OpenCiGuetzliById(ctx, bucket)
		} else {
			err = service.CloseCiGuetzliById(ctx, bucket)
		}
		if err != nil {
			return err
		}
		return resourceTencentCloudCIGuetzliRead(d, meta)

	}
	return nil
}

func resourceTencentCloudCIGuetzliDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_guetzli.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Get("bucket").(string)
	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := service.CloseCiGuetzliById(ctx, bucket)

	if err != nil {
		return err
	}
	return nil
}

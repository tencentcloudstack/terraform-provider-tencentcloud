/*
Provides a resource to create a ci bucket_pic_style

Example Usage

```hcl
resource "tencentcloud_ci_bucket_pic_style" "bucket_pic_style" {
  bucket     = "terraform-ci-xxxxxx"
  style_name = "rayscale_2"
  style_body = "imageMogr2/thumbnail/20x/crop/20x20/gravity/center/interlace/0/quality/100"
}
```

Import

ci bucket_pic_style can be imported using the bucket#styleName, e.g.

```
terraform import tencentcloud_ci_bucket_pic_style.bucket_pic_style terraform-ci-xxxxxx#rayscale_2
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	ci "github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCiBucketPicStyle() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiBucketPicStyleCreate,
		Read:   resourceTencentCloudCiBucketPicStyleRead,
		Delete: resourceTencentCloudCiBucketPicStyleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateCosBucketName,
				Description:  "bucket name.",
			},
			"style_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "style name, style names are case-sensitive, and a combination of uppercase and lowercase letters, numbers, and `$ + _ ( )` is supported.",
			},

			"style_body": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "style details, example: mageMogr2/grayscale/1.",
			},
		},
	}
}

func resourceTencentCloudCiBucketPicStyleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_bucket_pic_style.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		styleName string
		styleBody string
		bucket    string
	)
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, _ := d.GetOk("style_name"); v != nil {
		styleName = v.(string)
	}

	if v, ok := d.GetOk("style_body"); ok {
		styleBody = v.(string)
	}

	ciClient := meta.(*TencentCloudClient).apiV3Conn.UsePicClient(bucket)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := ciClient.CI.AddStyle(ctx, &ci.AddStyleOptions{
			StyleName: styleName,
			StyleBody: styleBody,
		})
		if e != nil {
			time.Sleep(5 * time.Second)
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, response status [%s]\n", logId, "AddStyle", result.Status)

		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci bucketPicStyle failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(bucket + FILED_SP + styleName)

	return resourceTencentCloudCiBucketPicStyleRead(d, meta)
}

func resourceTencentCloudCiBucketPicStyleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_bucket_pic_style.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	styleName := idSplit[1]

	bucketPicStyle, err := service.DescribeCiBucketPicStyleById(ctx, bucket, styleName)
	if err != nil {
		return err
	}

	if bucketPicStyle == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)
	_ = d.Set("style_name", bucketPicStyle.StyleName)
	_ = d.Set("style_body", bucketPicStyle.StyleBody)

	return nil
}

func resourceTencentCloudCiBucketPicStyleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_bucket_pic_style.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	styleName := idSplit[1]

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	if err := service.DeleteCiBucketPicStyleById(ctx, bucket, styleName); err != nil {
		return err
	}

	return nil
}

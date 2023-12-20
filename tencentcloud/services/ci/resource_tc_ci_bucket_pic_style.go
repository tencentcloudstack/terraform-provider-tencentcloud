package ci

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	ci "github.com/tencentyun/cos-go-sdk-v5"
)

func ResourceTencentCloudCiBucketPicStyle() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateCosBucketName,
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
	defer tccommon.LogElapsed("resource.tencentcloud_ci_bucket_pic_style.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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

	ciClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePicClient(bucket)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := ciClient.CI.AddStyle(ctx, &ci.AddStyleOptions{
			StyleName: styleName,
			StyleBody: styleBody,
		})
		if e != nil {
			time.Sleep(5 * time.Second)
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, response status [%s]\n", logId, "AddStyle", result.Status)

		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci bucketPicStyle failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(bucket + tccommon.FILED_SP + styleName)

	return resourceTencentCloudCiBucketPicStyleRead(d, meta)
}

func resourceTencentCloudCiBucketPicStyleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_bucket_pic_style.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
	defer tccommon.LogElapsed("resource.tencentcloud_ci_bucket_pic_style.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	styleName := idSplit[1]

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.DeleteCiBucketPicStyleById(ctx, bucket, styleName); err != nil {
		return err
	}

	return nil
}

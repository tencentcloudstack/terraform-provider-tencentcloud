/*
Provides a resource to create a ci bucket_attachment

Example Usage

```hcl
resource "tencentcloud_ci_bucket_attachment" "bucket_attachment" {
  bucket = "terraform-ci-xxxxxx"
  ci_status = &lt;nil&gt;
}
```

Import

ci bucket_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_ci_bucket_attachment.bucket_attachment bucket_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	ci "github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"time"
)

func resourceTencentCloudCiBucketAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiBucketAttachmentCreate,
		Read:   resourceTencentCloudCiBucketAttachmentRead,
		Delete: resourceTencentCloudCiBucketAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Bucket name.",
			},

			"ci_status": {
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Binding object storage state, `on`: bound, `off`: unbound, `unbinding`: unbinding.",
			},
		},
	}
}

func resourceTencentCloudCiBucketAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_bucket_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = ci.NewOpenCIServiceRequest()
		response = ci.NewOpenCIServiceResponse()
		bucket   string
	)
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
		request.bucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ci_status"); ok {
		request.ciStatus = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().OpenCIService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci bucketAttachment failed, reason:%+v", logId, err)
		return err
	}

	bucket = *response.Response.Bucket
	d.SetId(bucket)

	return resourceTencentCloudCiBucketAttachmentRead(d, meta)
}

func resourceTencentCloudCiBucketAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_bucket_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	bucketAttachmentId := d.Id()

	bucketAttachment, err := service.DescribeCiBucketAttachmentById(ctx, bucket)
	if err != nil {
		return err
	}

	if bucketAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiBucketAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if bucketAttachment.bucket != nil {
		_ = d.Set("bucket", bucketAttachment.bucket)
	}

	if bucketAttachment.ciStatus != nil {
		_ = d.Set("ci_status", bucketAttachment.ciStatus)
	}

	return nil
}

func resourceTencentCloudCiBucketAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_bucket_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	bucketAttachmentId := d.Id()

	if err := service.DeleteCiBucketAttachmentById(ctx, bucket); err != nil {
		return err
	}

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"off"}, 0*readRetryTimeout, time.Second, service.CiBucketAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

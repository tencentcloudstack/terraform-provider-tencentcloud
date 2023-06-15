/*
Provides a resource to create a cos bucket_version

Example Usage

```hcl
resource "tencentcloud_cos_bucket_version" "bucket_version" {
  bucket = "mycos-1258798060"
  status = "Enabled"
}
```

Import

cos bucket_version can be imported using the id, e.g.

```
terraform import tencentcloud_cos_bucket_version.bucket_version bucket_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCosBucketVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketVersionCreate,
		Read:   resourceTencentCloudCosBucketVersionRead,
		Update: resourceTencentCloudCosBucketVersionUpdate,
		Delete: resourceTencentCloudCosBucketVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.",
			},

			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Whether to enable versioning. Valid values: `Suspended`, `Enabled`.",
			},
		},
	}
}

func resourceTencentCloudCosBucketVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_version.create")()
	defer inconsistentCheck(d, meta)()

	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	}

	d.SetId(bucket)

	return resourceTencentCloudCosBucketVersionUpdate(d, meta)
}

func resourceTencentCloudCosBucketVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CosService{client: meta.(*TencentCloudClient).apiV3Conn}

	bucket := d.Id()

	bucketVersion, err := service.DescribeCosBucketVersionById(ctx, bucket)
	if err != nil {
		return err
	}

	if bucketVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CosBucketVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("bucket", bucket)

	if bucketVersion.Status != "" {
		_ = d.Set("status", bucketVersion.Status)
	}

	return nil
}

func resourceTencentCloudCosBucketVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_version.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()

	request := cos.BucketPutVersionOptions{}
	if v, ok := d.GetOk("status"); ok {
		request.Status = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Bucket.PutVersioning(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%+v], response status [%s]\n", logId, "PutVersioning", request, result.Status)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cos versioning failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCosBucketVersionRead(d, meta)
}

func resourceTencentCloudCosBucketVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

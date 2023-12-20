package cos

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func ResourceTencentCloudCosBucketVersion() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_version.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	}

	d.SetId(bucket)

	return resourceTencentCloudCosBucketVersionUpdate(d, meta)
}

func resourceTencentCloudCosBucketVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_version.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_version.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()

	request := cos.BucketPutVersionOptions{}
	if v, ok := d.GetOk("status"); ok {
		request.Status = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Bucket.PutVersioning(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_version.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

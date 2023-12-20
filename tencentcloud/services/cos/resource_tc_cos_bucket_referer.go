package cos

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func ResourceTencentCloudCosBucketReferer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketRefererCreate,
		Read:   resourceTencentCloudCosBucketRefererRead,
		Update: resourceTencentCloudCosBucketRefererUpdate,
		Delete: resourceTencentCloudCosBucketRefererDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.",
			},
			"referer_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hotlink protection type. Enumerated values: `Black-List`, `White-List`.",
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Whether to enable hotlink protection. Enumerated values: `Enabled`, `Disabled`.",
			},
			"domain_list": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "A list of domain names in the blocklist/allowlist.",
			},
			"empty_refer_configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether to allow access with an empty referer. Enumerated values: `Allow`, `Deny` (default).",
			},
		},
	}
}

func resourceTencentCloudCosBucketRefererCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_referer.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	}

	d.SetId(bucket)

	return resourceTencentCloudCosBucketRefererUpdate(d, meta)
}

func resourceTencentCloudCosBucketRefererRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_referer.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	bucket := d.Id()

	bucketReferer, err := service.DescribeCosBucketRefererById(ctx, bucket)
	if err != nil {
		return err
	}

	if bucketReferer == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CosBucketReferer` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("bucket", bucket)

	if bucketReferer.Status != "" {
		_ = d.Set("status", bucketReferer.Status)
	}

	if bucketReferer.RefererType != "" {
		_ = d.Set("referer_type", bucketReferer.RefererType)
	}

	if len(bucketReferer.DomainList) > 0 {
		_ = d.Set("domain_list", bucketReferer.DomainList)
	}

	if bucketReferer.EmptyReferConfiguration != "" {
		_ = d.Set("empty_refer_configuration", bucketReferer.EmptyReferConfiguration)
	}

	return nil
}

func resourceTencentCloudCosBucketRefererUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_referer.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()

	request := cos.BucketPutRefererOptions{}
	if v, ok := d.GetOk("status"); ok {
		request.Status = v.(string)
	}
	if v, ok := d.GetOk("referer_type"); ok {
		request.RefererType = v.(string)
	}
	if v, ok := d.GetOk("domain_list"); ok {
		domainListSet := v.(*schema.Set).List()
		for i := range domainListSet {
			domainList := domainListSet[i].(string)
			request.DomainList = append(request.DomainList, domainList)
		}
	}
	if v, ok := d.GetOk("empty_refer_configuration"); ok {
		request.EmptyReferConfiguration = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Bucket.PutReferer(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%+v], response status [%s]\n", logId, "PutReferer", request, result.Status)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cos bucketReferer failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCosBucketRefererRead(d, meta)
}

func resourceTencentCloudCosBucketRefererDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_referer.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

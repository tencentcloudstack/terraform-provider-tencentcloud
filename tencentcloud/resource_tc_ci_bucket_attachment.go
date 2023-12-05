package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
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
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateCosBucketName,
				Description:  "bucket name.",
			},
			"ci_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Binding object storage state, `on`: bound, `off`: unbound, `unbinding`: unbinding.",
			},
		},
	}
}

func resourceTencentCloudCiBucketAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_bucket_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	ciClient := meta.(*TencentCloudClient).apiV3Conn.UsePicClient(bucket)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := ciClient.CI.OpenCIService(ctx)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response status [%s]\n", logId, "OpenCIService", bucket, result.Status)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci bucket failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(bucket)

	return resourceTencentCloudCiBucketAttachmentRead(d, meta)
}

func resourceTencentCloudCiBucketAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_bucket_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	bucket := d.Id()

	result, err := service.DescribeCiBucketById(ctx, bucket)
	if err != nil {
		return err
	}

	if result == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)
	_ = d.Set("ci_status", result.CIStatus)

	return nil
}

func resourceTencentCloudCiBucketAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_bucket_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	bucket := d.Id()

	if err := service.DeleteCiBucketById(ctx, bucket); err != nil {
		return err
	}

	retryErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCiBucketById(ctx, bucket)
		if e != nil {
			return retryError(e)
		}
		if result.CIStatus == "unbinding" {
			return resource.RetryableError(fmt.Errorf("Binding bucket status is %s , retry...", result.CIStatus))
		}
		if result.CIStatus == "off" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Binding bucket status is %s , retry...", result.CIStatus))
	})
	if retryErr != nil {
		return retryErr
	}

	return nil
}

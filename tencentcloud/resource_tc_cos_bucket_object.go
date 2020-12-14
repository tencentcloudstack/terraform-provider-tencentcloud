/*
Provides a COS object resource to put an object(content or file) to the bucket.

Example Usage

Uploading a file to a bucket

```hcl
resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket = "mycos-1258798060"
  key    = "new_object_key"
  source = "path/to/file"
}
```

Uploading a content to a bucket

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket  = tencentcloud_cos_bucket.mycos.bucket
  key     = "new_object_key"
  content = "the content that you want to upload."
}
```
*/
package tencentcloud

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/go-homedir"
)

func resourceTencentCloudCosBucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketObjectCreate,
		Read:   resourceTencentCloudCosBucketObjectRead,
		Update: resourceTencentCloudCosBucketObjectUpdate,
		Delete: resourceTencentCloudCosBucketObjectDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of a bucket to use. Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the object once it is in the bucket.",
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
				Description:   "The path to the source file being uploaded to the bucket.",
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
				Description:   "Literal string value to use as the object content, which will be uploaded as UTF-8-encoded text.",
			},
			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  s3.ObjectCannedACLPrivate,
				ValidateFunc: validateAllowedStringValue([]string{
					s3.ObjectCannedACLPrivate,
					s3.ObjectCannedACLPublicRead,
					s3.ObjectCannedACLPublicReadWrite,
				}),
				Description: "The canned ACL to apply. Available values include `private`, `public-read`, and `public-read-write`. Defaults to `private`.",
			},
			"cache_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies caching behavior along the request/reply chain. For further details, RFC2616 can be referred.",
			},
			"content_disposition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies presentational information for the object.",
			},
			"content_encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A standard MIME type describing the format of the object data.",
			},
			"storage_class": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(availableCosStorageClass),
				Description:  "Object storage type, Available values include `STANDARD`, `STANDARD_IA` and `ARCHIVE`.",
			},
			"etag": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ETag generated for the object (an MD5 sum of the object content).",
			},
		},
	}
}

func resourceTencentCloudCosBucketObjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_object.create")()

	logId := getLogId(contextNil)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	var body io.ReadSeeker
	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("cos object source (%s) homedir expand error: %s", source, err.Error())
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("cos object source (%s) open error: %s", source, err.Error())
		}
		body = file
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("closing cos object source (%s) error: %s", path, err.Error())
			}
		}()
	} else if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		body = bytes.NewReader([]byte(content))
	} else {
		return fmt.Errorf("must specify \"source\" or \"content\" field")
	}

	request := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	if v, ok := d.GetOk("acl"); ok {
		request.ACL = aws.String(v.(string))
	}
	if v, ok := d.GetOk("cache_control"); ok {
		request.CacheControl = aws.String(v.(string))
	}
	if v, ok := d.GetOk("content_disposition"); ok {
		request.ContentDisposition = aws.String(v.(string))
	}
	if v, ok := d.GetOk("content_encoding"); ok {
		request.ContentEncoding = aws.String(v.(string))
	}
	if v, ok := d.GetOk("content_type"); ok {
		request.ContentType = aws.String(v.(string))
	}
	if v, ok := d.GetOk("storage_class"); ok {
		request.StorageClass = aws.String(v.(string))
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseCosClient().PutObject(request)
	if err != nil {
		return fmt.Errorf("putting object (%s) in cos bucket (%s) error: %s", key, bucket, err.Error())
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put object", request.String(), response.String())

	d.SetId(bucket + key)
	return resourceTencentCloudCosBucketObjectRead(d, meta)
}

func resourceTencentCloudCosBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_object.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	cosService := CosService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	response, err := cosService.HeadObject(ctx, bucket, key)
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			log.Printf("[WARN]%s object (%s) in bucket (%s) not found, error code (404)", logId, key, bucket)
			d.SetId("")
			return nil
		}
		return err
	}

	_ = d.Set("cache_control", response.CacheControl)
	_ = d.Set("content_disposition", response.ContentDisposition)
	_ = d.Set("content_encoding", response.ContentEncoding)
	_ = d.Set("content_type", response.ContentType)
	_ = d.Set("etag", strings.Trim(*response.ETag, `"`))
	_ = d.Set("storage_class", s3.StorageClassStandard)
	if response.StorageClass != nil {
		_ = d.Set("storage_class", response.StorageClass)
	}

	return nil
}

func resourceTencentCloudCosBucketObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_object.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	fields := []string{
		"cache_control",
		"content_disposition",
		"content_encoding",
		"content_type",
		"source",
		"content",
		"storage_class",
		"etag",
	}
	for _, key := range fields {
		if d.HasChange(key) {
			return resourceTencentCloudCosBucketObjectCreate(d, meta)
		}
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	cosService := CosService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	if d.HasChange("acl") {
		acl := d.Get("acl").(string)
		err := cosService.PutObjectAcl(ctx, bucket, key, acl)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudCosBucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_object.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	cosService := CosService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := cosService.DeleteObject(ctx, bucket, key)
	if err != nil {
		return err
	}

	return nil
}

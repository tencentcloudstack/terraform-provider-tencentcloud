package tencentcloud

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudCosBucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketObjectCreate,
		Read:   resourceTencentCloudCosBucketObjectRead,
		Update: resourceTencentCloudCosBucketObjectUpdate,
		Delete: resourceTencentCloudCosBucketObjectDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
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
			},
			"cache_control": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_disposition": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validateAllowedStringValue([]string{
					// todo: check enum value
					s3.ObjectStorageClassStandard,
					s3.ObjectStorageClassStandardIa,
				}),
			},
			"etag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudCosBucketObjectCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)

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
		return fmt.Errorf("puting object (%s) in cos bucket (%s) error: %s", key, bucket, err.Error())
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put object", request.String(), response.String())

	d.SetId(bucket + key)
	return resourceTencentCloudCosBucketObjectRead(d, meta)
}

func resourceTencentCloudCosBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	cosService := CosService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	response, err := cosService.HeadObject(ctx, bucket, key)
	if err != nil {
		return err
	}

	d.Set("cache_control", response.CacheControl)
	d.Set("content_disposition", response.ContentDisposition)
	d.Set("content_encoding", response.ContentEncoding)
	d.Set("content_type", response.ContentType)
	d.Set("etag", strings.Trim(*response.ETag, `"`))
	d.Set("storage_class", s3.StorageClassStandard)
	if response.StorageClass != nil {
		d.Set("storage_class", response.StorageClass)
	}

	return nil
}

func resourceTencentCloudCosBucketObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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

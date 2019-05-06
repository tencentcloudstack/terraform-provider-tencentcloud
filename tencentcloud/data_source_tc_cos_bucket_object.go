package tencentcloud

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudCosBucketObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCosBucketObjectsRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cache_control": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_disposition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_encoding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_language": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"server_side_encryption": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"website_redirect_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudCosBucketObjectsRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	versionId := ""
	if v, ok := d.GetOk("version_id"); ok {
		versionId = v.(string)
	}
	cosService := CosService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	info, err := cosService.HeadObject(ctx, bucket, key)
	if err != nil {
		return err
	}

	ids := []string{bucket, key, versionId}
	d.SetId(dataResourceIdsHash(ids))
	d.Set("cache_control", info.CacheControl)
	d.Set("content_disposition", info.ContentDisposition)
	d.Set("content_encoding", info.ContentEncoding)
	d.Set("content_language", info.ContentLanguage)
	d.Set("content_length", info.ContentLength)
	d.Set("content_type", info.ContentType)
	d.Set("etag", strings.Trim(*info.ETag, `"`))
	d.Set("expiration", info.Expiration)
	d.Set("expires", info.Expires)
	d.Set("last_modified", info.LastModified.Format(time.RFC1123))
	d.Set("metadata", pointersMapToStringMap(info.Metadata))
	d.Set("server_side_encryption", info.ServerSideEncryption)
	d.Set("version_id", info.VersionId)
	d.Set("website_redirect_location", info.WebsiteRedirectLocation)
	d.Set("storage_class", s3.StorageClassStandard)
	if info.StorageClass != nil {
		d.Set("storage_class", info.StorageClass)
	}

	return nil
}

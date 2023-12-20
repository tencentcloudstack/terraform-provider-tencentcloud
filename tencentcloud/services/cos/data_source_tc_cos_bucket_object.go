package cos

import (
	"context"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCosBucketObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCosBucketObjectsRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the bucket that contains the objects to query.",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The full path to the object inside the bucket.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"cache_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies caching behavior along the request/reply chain.",
			},
			"content_disposition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies presentational information for the object.",
			},
			"content_encoding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A standard MIME type describing the format of the object data.",
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ETag generated for the object, which is may not equal to MD5 value.",
			},
			"last_modified": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified date of the object.",
			},
			"storage_class": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object storage type such as STANDARD.",
			},
		},
	}
}

func dataSourceTencentCloudCosBucketObjectsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cos_bucket_object.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	outputMap := make(map[string]string)
	outputMap["bucket"] = bucket
	outputMap["key"] = key
	cosService := CosService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	info, err := cosService.HeadObject(ctx, bucket, key)
	if err != nil {
		return err
	}

	ids := []string{bucket, key}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("cache_control", info.CacheControl)
	outputMap["cache_control"] = getStringValue(info.CacheControl)
	_ = d.Set("content_disposition", info.ContentDisposition)
	outputMap["content_disposition"] = getStringValue(info.ContentDisposition)
	_ = d.Set("content_encoding", info.ContentEncoding)
	outputMap["content_encoding"] = getStringValue(info.ContentEncoding)
	_ = d.Set("content_type", info.ContentType)
	outputMap["content_type"] = getStringValue(info.ContentType)
	etag := getStringValue(info.ETag)
	_ = d.Set("etag", strings.Trim(etag, `"`))
	outputMap["etag"] = strings.Trim(etag, `"`)
	_ = d.Set("last_modified", info.LastModified.Format(time.RFC1123))
	outputMap["last_modified"] = info.LastModified.Format(time.RFC1123)
	_ = d.Set("storage_class", s3.StorageClassStandard)
	outputMap["storage_class"] = s3.StorageClassStandard
	if info.StorageClass != nil {
		_ = d.Set("storage_class", info.StorageClass)
		outputMap["storage_class"] = getStringValue(info.StorageClass)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = tccommon.WriteToFile(output.(string), outputMap); err != nil {
			return err
		}
	}

	return nil
}

func getStringValue(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

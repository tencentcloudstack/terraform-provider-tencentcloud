package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudCosBucketObjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCosBucketObjectsRead,

		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"etag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCosBucketObjectsRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	bucketName := d.Get("bucket_name").(string)
	keyPrefixString := ""
	if keyPrefix, ok := d.GetOk("key_prefix"); ok {
		keyPrefixString = keyPrefix.(string)
	}

	cosService := CosService{client: meta.(*TencentCloudClient).apiV3Conn}
	result, err := cosService.GetObjectList(ctx, bucketName, keyPrefixString)
	if err != nil {
		return err
	}

	objectList := make([]map[string]interface{}, 0, len(result.Contents))
	for _, item := range result.Contents {
		object := map[string]interface{}{
			"key":           item.Key,
			"last_modified": item.LastModified,
			"etag":          item.ETag,
			"size":          item.Size,
			"storage_class": item.StorageClass,
		}
		objectList = append(objectList, object)
	}
	ids := make([]string, 2)
	ids[0] = bucketName
	ids[1] = keyPrefixString
	d.SetId(dataResourceIdsHash(ids))
	err = d.Set("object_list", objectList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set object list fail, reason:%s\n ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		writeToFile(output.(string), objectList)
	}
	return nil
}

/*
Use this data source to query the COS buckets of the current Tencent Cloud user.

Example Usage

```hcl
data "tencentcloud_cos_buckets" "cos_buckets" {
	bucket_prefix = "tf-bucket-"
    result_output_file = "mytestpath"
}
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudCosBuckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCosBucketsRead,

		Schema: map[string]*schema.Schema{
			"bucket_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A prefix string to filter results by bucket name",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"bucket_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of bucket. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bucket name, the format likes <bucket>-<appid>.",
						},
						"cors_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of CORS rule configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_origins": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Specifies which origins are allowed.",
									},
									"allowed_methods": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Specifies which methods are allowed. Can be GET, PUT, POST, DELETE or HEAD.",
									},
									"allowed_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Specifies which headers are allowed.",
									},
									"max_age_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies time in seconds that browser can cache the response for a preflight request.",
									},
									"expose_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Specifies expose header in the response.",
									},
								},
							},
						},
						"lifecycle_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The lifecycle configuration of a bucket.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter_prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Object key prefix identifying one or more objects to which the rule applies.",
									},
									"transition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a period in the object's transitions (documented below).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the date after which you want the corresponding action to take effect.",
												},
												"days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of days after object creation when the specific rule action takes effect.",
												},
												"storage_class": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the storage class to which you want the object to transition. Available values include STANDARD, STANDARD_IA and ARCHIVE.",
												},
											},
										},
									},
									"expiration": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a period in the object's expire (documented below).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the date after which you want the corresponding action to take effect.",
												},
												"days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of days after object creation when the specific rule action takes effect.",
												},
											},
										},
									},
								},
							},
						},
						"website": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of one element containing configuration parameters used when the bucket is used as a website.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index_document": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "COS returns this index document when requests are made to the root domain or any of the subfolders.",
									},
									"error_document": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An absolute path to the document to return in case of a 4XX error.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCosBucketsRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cosService := CosService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	buckets, err := cosService.ListBuckets(ctx)
	if err != nil {
		return err
	}

	prefix := d.Get("bucket_prefix").(string)
	bucketList := make([]map[string]interface{}, 0, len(buckets))
	for _, v := range buckets {
		bucket := make(map[string]interface{})
		if prefix != "" && !strings.HasPrefix(*v.Name, prefix) {
			continue
		}

		bucket["bucket"] = *v.Name
		corsRules, err := cosService.GetBucketCors(ctx, *v.Name)
		if err != nil {
			return err
		}
		bucket["cors_rules"] = corsRules
		lifecycleRules, err := cosService.GetDataSourceBucketLifecycle(ctx, *v.Name)
		if err != nil {
			return err
		}
		bucket["lifecycle_rules"] = lifecycleRules
		website, err := cosService.GetBucketWebsite(ctx, *v.Name)
		if err != nil {
			return err
		}
		bucket["website"] = website

		bucketList = append(bucketList, bucket)
	}

	ids := make([]string, 2)
	ids[0] = "bucketlist"
	ids[1] = prefix
	d.SetId(dataResourceIdsHash(ids))
	if err := d.Set("bucket_list", bucketList); err != nil {
		return fmt.Errorf("setting bucket list error: %s", err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = writeToFile(output.(string), bucketList); err != nil {
			return err
		}
	}

	return nil
}

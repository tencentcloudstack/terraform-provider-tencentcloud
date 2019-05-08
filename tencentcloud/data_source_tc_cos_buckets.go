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
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cors_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_origins": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allowed_methods": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allowed_headers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"max_age_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"expose_headers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"lifecycle_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"transition": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"days": {
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
									"expiration": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"days": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"website": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index_document": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"error_document": {
										Type:     schema.TypeString,
										Computed: true,
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

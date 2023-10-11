/*
Use this data source to query detailed information of ssl describe_host_cos_instance_list

Example Usage

```hcl
data "tencentcloud_ssl_describe_host_cos_instance_list" "describe_host_cos_instance_list" {
  certificate_id = "8u8DII0l"
  resource_type = "cos"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeHostCosInstanceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostCosInstanceListRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID to be deployed.",
			},

			"resource_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy resource type cos.",
			},

			"is_cache": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to query the cache, 1: Yes; 0: No, the default is the query cache, the cache is half an hour.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of filter parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter parameter key.",
						},
						"filter_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter parameter value.",
						},
					},
				},
			},

			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "COS instance listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "domain name.",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Binded certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enabled: domain name online statusDisabled: Domain name offline status.",
						},
						"bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reserve bucket nameNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Barrel areaNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
					},
				},
			},

			"async_total_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of asynchronous refreshNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"async_offset": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Asynchronous refresh current execution numberNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"async_cache_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current cache read timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeHostCosInstanceListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_host_cos_instance_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("certificate_id"); ok {
		paramMap["CertificateId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		paramMap["ResourceType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("is_cache"); v != nil {
		paramMap["IsCache"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*ssl.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := ssl.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["filter_key"]; ok {
				filter.FilterKey = helper.String(v.(string))
			}
			if v, ok := filterMap["filter_value"]; ok {
				filter.FilterValue = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceList *ssl.DescribeHostCosInstanceListResponseParams

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeHostCosInstanceListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	tmpList := make([]map[string]interface{}, 0)

	if instanceList != nil && instanceList.InstanceList != nil {

		for _, cosInstanceDetail := range instanceList.InstanceList {
			cosInstanceDetailMap := map[string]interface{}{}

			if cosInstanceDetail.Domain != nil {
				cosInstanceDetailMap["domain"] = cosInstanceDetail.Domain
			}

			if cosInstanceDetail.CertId != nil {
				cosInstanceDetailMap["cert_id"] = cosInstanceDetail.CertId
			}

			if cosInstanceDetail.Status != nil {
				cosInstanceDetailMap["status"] = cosInstanceDetail.Status
			}

			if cosInstanceDetail.Bucket != nil {
				cosInstanceDetailMap["bucket"] = cosInstanceDetail.Bucket
			}

			if cosInstanceDetail.Region != nil {
				cosInstanceDetailMap["region"] = cosInstanceDetail.Region
			}

			ids = append(ids, *cosInstanceDetail.CertId)
			tmpList = append(tmpList, cosInstanceDetailMap)
		}

		_ = d.Set("instance_list", tmpList)
	}

	if instanceList != nil && instanceList.AsyncTotalNum != nil {
		_ = d.Set("async_total_num", instanceList.AsyncTotalNum)
	}

	if instanceList != nil && instanceList.AsyncOffset != nil {
		_ = d.Set("async_offset", instanceList.AsyncOffset)
	}

	if instanceList != nil && instanceList.AsyncCacheTime != nil {
		_ = d.Set("async_cache_time", instanceList.AsyncCacheTime)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

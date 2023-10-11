/*
Use this data source to query detailed information of ssl describe_host_live_instance_list

Example Usage

```hcl
data "tencentcloud_ssl_describe_host_live_instance_list" "describe_host_live_instance_list" {
  certificate_id = "8u8DII0l"
  resource_type = "live"
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

func dataSourceTencentCloudSslDescribeHostLiveInstanceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostLiveInstanceListRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID to be deployed.",
			},

			"resource_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy resource type.",
			},

			"is_cache": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to query the cache, 1: Yes; 0: No, the default is the query cache, the cache is half an hour.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of filtering parameters; Filterkey: domainmatch.",
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

			"old_certificate_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Deployed certificate ID.",
			},

			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Live instance listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "-1: Unrelated certificate of domain name.1: The domain name HTTPS has been opened.0: The domain name HTTPS has been closed.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeHostLiveInstanceListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_host_live_instance_list.read")()
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

	if v, ok := d.GetOk("old_certificate_id"); ok {
		paramMap["OldCertificateId"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceList []*ssl.LiveInstanceDetail

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeHostLiveInstanceListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	tmpList := make([]map[string]interface{}, 0, len(instanceList))

	if instanceList != nil {
		for _, liveInstanceDetail := range instanceList {
			liveInstanceDetailMap := map[string]interface{}{}

			if liveInstanceDetail.Domain != nil {
				liveInstanceDetailMap["domain"] = liveInstanceDetail.Domain
			}

			if liveInstanceDetail.CertId != nil {
				liveInstanceDetailMap["cert_id"] = liveInstanceDetail.CertId
			}

			if liveInstanceDetail.Status != nil {
				liveInstanceDetailMap["status"] = liveInstanceDetail.Status
			}

			ids = append(ids, *liveInstanceDetail.CertId)
			tmpList = append(tmpList, liveInstanceDetailMap)
		}

		_ = d.Set("instance_list", tmpList)
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

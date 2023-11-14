/*
Use this data source to query detailed information of tsf microservice

Example Usage

```hcl
data "tencentcloud_tsf_microservice" "microservice" {
  namespace_id = "ns-123456"
  search_word = ""
  order_by = ""
  order_type = 0
  status =
  microservice_id_list =
  microservice_name_list =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfMicroservice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfMicroserviceRead,
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace id.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search word.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field.",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Sorting type field. 0 or 1.",
			},

			"status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Status filter，online、offline、single_online.",
			},

			"microservice_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Microservice id list.",
			},

			"microservice_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of service names for search.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Microservice paging list information. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Microservice paging list information. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Microservice list information. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"microservice_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice Id. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"microservice_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice name. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"microservice_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice description. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"create_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CreationTime. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"update_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Last update time.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace Id.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"run_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Run instance count in namespace.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"critical_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Offline instance count.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
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

func dataSourceTencentCloudTsfMicroserviceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_microservice.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("namespace_id"); ok {
		paramMap["NamespaceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		paramMap["Status"] = helper.InterfacesStringsPoint(statusSet)
	}

	if v, ok := d.GetOk("microservice_id_list"); ok {
		microserviceIdListSet := v.(*schema.Set).List()
		paramMap["MicroserviceIdList"] = helper.InterfacesStringsPoint(microserviceIdListSet)
	}

	if v, ok := d.GetOk("microservice_name_list"); ok {
		microserviceNameListSet := v.(*schema.Set).List()
		paramMap["MicroserviceNameList"] = helper.InterfacesStringsPoint(microserviceNameListSet)
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tsf.TsfPageMicroservice

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfMicroserviceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		tsfPageMicroserviceMap := map[string]interface{}{}

		if result.TotalCount != nil {
			tsfPageMicroserviceMap["total_count"] = result.TotalCount
		}

		if result.Content != nil {
			contentList := []interface{}{}
			for _, content := range result.Content {
				contentMap := map[string]interface{}{}

				if content.MicroserviceId != nil {
					contentMap["microservice_id"] = content.MicroserviceId
				}

				if content.MicroserviceName != nil {
					contentMap["microservice_name"] = content.MicroserviceName
				}

				if content.MicroserviceDesc != nil {
					contentMap["microservice_desc"] = content.MicroserviceDesc
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.NamespaceId != nil {
					contentMap["namespace_id"] = content.NamespaceId
				}

				if content.RunInstanceCount != nil {
					contentMap["run_instance_count"] = content.RunInstanceCount
				}

				if content.CriticalInstanceCount != nil {
					contentMap["critical_instance_count"] = content.CriticalInstanceCount
				}

				contentList = append(contentList, contentMap)
			}

			tsfPageMicroserviceMap["content"] = []interface{}{contentList}
		}

		ids = append(ids, *result.NamespaceId)
		_ = d.Set("result", tsfPageMicroserviceMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageMicroserviceMap); e != nil {
			return e
		}
	}
	return nil
}

/*
Use this data source to query detailed information of apigateway a_p_i

Example Usage

```hcl
data "tencentcloud_apigateway_a_p_i" "a_p_i" {
  service_id = ""
    tags = {
    "createdBy" = "terraform"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudApigatewayAPI() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApigatewayAPIRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the service to be queried.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "API binding usage plan list.            Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of usage plans for the API binding.                Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"api_usage_plan_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of API binding usage plans.                Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service unique ID.                    Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API unique ID.                    Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"api_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API name.                     Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API path.                     Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API methods.                     Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"usage_plan_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Use the plan&amp;#39;s unique ID.                     Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"usage_plan_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Use the name of the plan.                     Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"usage_plan_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the usage plan.                     Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"environment": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Use a plan-bound service environment.                     Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"in_use_request_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Already used quota.                     Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"max_request_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total request quota, -1 means no limit.                     Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"max_request_num_pre_sec": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Request QPS cap, -1 means no limit.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Use scheduled creation time.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"modified_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Use the time the plan was last modified.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service name.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudApigatewayAPIRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_apigateway_a_p_i.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*apigateway.ApiUsagePlanSet

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApigatewayAPIByFilter(ctx, paramMap)
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
		apiUsagePlanSetMap := map[string]interface{}{}

		if result.TotalCount != nil {
			apiUsagePlanSetMap["total_count"] = result.TotalCount
		}

		if result.ApiUsagePlanList != nil {
			apiUsagePlanListList := []interface{}{}
			for _, apiUsagePlanList := range result.ApiUsagePlanList {
				apiUsagePlanListMap := map[string]interface{}{}

				if apiUsagePlanList.ServiceId != nil {
					apiUsagePlanListMap["service_id"] = apiUsagePlanList.ServiceId
				}

				if apiUsagePlanList.ApiId != nil {
					apiUsagePlanListMap["api_id"] = apiUsagePlanList.ApiId
				}

				if apiUsagePlanList.ApiName != nil {
					apiUsagePlanListMap["api_name"] = apiUsagePlanList.ApiName
				}

				if apiUsagePlanList.Path != nil {
					apiUsagePlanListMap["path"] = apiUsagePlanList.Path
				}

				if apiUsagePlanList.Method != nil {
					apiUsagePlanListMap["method"] = apiUsagePlanList.Method
				}

				if apiUsagePlanList.UsagePlanId != nil {
					apiUsagePlanListMap["usage_plan_id"] = apiUsagePlanList.UsagePlanId
				}

				if apiUsagePlanList.UsagePlanName != nil {
					apiUsagePlanListMap["usage_plan_name"] = apiUsagePlanList.UsagePlanName
				}

				if apiUsagePlanList.UsagePlanDesc != nil {
					apiUsagePlanListMap["usage_plan_desc"] = apiUsagePlanList.UsagePlanDesc
				}

				if apiUsagePlanList.Environment != nil {
					apiUsagePlanListMap["environment"] = apiUsagePlanList.Environment
				}

				if apiUsagePlanList.InUseRequestNum != nil {
					apiUsagePlanListMap["in_use_request_num"] = apiUsagePlanList.InUseRequestNum
				}

				if apiUsagePlanList.MaxRequestNum != nil {
					apiUsagePlanListMap["max_request_num"] = apiUsagePlanList.MaxRequestNum
				}

				if apiUsagePlanList.MaxRequestNumPreSec != nil {
					apiUsagePlanListMap["max_request_num_pre_sec"] = apiUsagePlanList.MaxRequestNumPreSec
				}

				if apiUsagePlanList.CreatedTime != nil {
					apiUsagePlanListMap["created_time"] = apiUsagePlanList.CreatedTime
				}

				if apiUsagePlanList.ModifiedTime != nil {
					apiUsagePlanListMap["modified_time"] = apiUsagePlanList.ModifiedTime
				}

				if apiUsagePlanList.ServiceName != nil {
					apiUsagePlanListMap["service_name"] = apiUsagePlanList.ServiceName
				}

				apiUsagePlanListList = append(apiUsagePlanListList, apiUsagePlanListMap)
			}

			apiUsagePlanSetMap["api_usage_plan_list"] = []interface{}{apiUsagePlanListList}
		}

		ids = append(ids, *result.ApiId)
		_ = d.Set("result", apiUsagePlanSetMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), apiUsagePlanSetMap); e != nil {
			return e
		}
	}
	return nil
}

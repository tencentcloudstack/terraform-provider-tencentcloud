package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAPIGatewayApiUsagePlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayApiUsagePlanRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the service to be queried.",
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "API binding usage plan list.Note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service unique ID.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service name.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API unique ID.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API name.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API path.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API method.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"usage_plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the unique ID of the plan.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"usage_plan_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the name of the plan.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"usage_plan_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the usage plan.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the service environment bound by the plan.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"in_use_request_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The quota that has already been used.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"max_request_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Request total quota, -1 indicates no limit.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"max_request_num_pre_sec": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Request QPS upper limit, -1 indicates no limit.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create a time using a schedule.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the last modification time of the plan.Note: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudAPIGatewayApiUsagePlanRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_api_usage_plans.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		result  []*apigateway.ApiUsagePlan
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeAPIGatewayApiUsagePlanByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		result = response
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		apiUsagePlanListList := []interface{}{}
		for _, apiUsagePlanList := range result {
			apiUsagePlanListMap := map[string]interface{}{}

			if apiUsagePlanList.ServiceId != nil {
				apiUsagePlanListMap["service_id"] = apiUsagePlanList.ServiceId
			}

			if apiUsagePlanList.ServiceName != nil {
				apiUsagePlanListMap["service_name"] = apiUsagePlanList.ServiceName
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

			ids = append(ids, *apiUsagePlanList.UsagePlanId)
			apiUsagePlanListList = append(apiUsagePlanListList, apiUsagePlanListMap)
		}

		_ = d.Set("result", apiUsagePlanListList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

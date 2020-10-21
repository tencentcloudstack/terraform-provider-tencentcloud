/*
Use this data source to query API gateway usage plans.

Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
  usage_plan_name         = "my_plan"
  usage_plan_desc         = "nice plan"
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

data "tencentcloud_api_gateway_usage_plans" "name" {
  usage_plan_name = tencentcloud_api_gateway_usage_plan.plan.usage_plan_name
}

data "tencentcloud_api_gateway_usage_plans" "id" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func dataSourceTencentCloudAPIGatewayUsagePlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayUsagePlansRead,

		Schema: map[string]*schema.Schema{
			"usage_plan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the usage plan.",
			},
			"usage_plan_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the usage plan.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			// Computed values.
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of usage plans.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"usage_plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the usage plan.",
						},
						"usage_plan_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the usage plan.",
						},
						"usage_plan_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom usage plan description.",
						},
						"max_request_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of requests allowed. Valid values: -1, [1,99999999]. The default value is -1, which indicates no limit.",
						},
						"max_request_num_pre_sec": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Limit of requests per second. Valid values: -1, [1,2000]. The default value is -1, which indicates no limit.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
						},
						"attach_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Attach service and API list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service ID.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service name.",
									},
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The API ID, this value is empty if attach service.",
									},
									"api_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The API name, this value is empty if attach service.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The API path, this value is empty if attach service.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The API method, this value is empty if attach service.",
									},
									"environment": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The environment name.",
									},
									"modify_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
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

func dataSourceTencentCloudAPIGatewayUsagePlansRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_usage_plans.read")

	var (
		logId                      = getLogId(contextNil)
		ctx                        = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService          = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		infos                      []*apigateway.UsagePlanStatusInfo
		list                       []map[string]interface{}
		usagePlanId, usagePlanName string
		err                        error
	)

	if v, ok := d.GetOk("usage_plan_id"); ok {
		usagePlanId = v.(string)
	}
	if v, ok := d.GetOk("usage_plan_name"); ok {
		usagePlanName = v.(string)
	}

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeUsagePlansStatus(ctx, usagePlanId, usagePlanName)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		var (
			infoMap    = make(map[string]interface{}, 10)
			attachList []*apigateway.UsagePlanEnvironment
		)
		for _, bindType := range API_GATEWAY_TYPES {
			var tyepeList []*apigateway.UsagePlanEnvironment
			if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				tyepeList, err = apiGatewayService.DescribeUsagePlanEnvironments(ctx, *info.UsagePlanId, bindType)
				if err != nil {
					return retryError(err, InternalError)
				}
				attachList = append(attachList, tyepeList...)
				return nil
			}); err != nil {
				return err
			}
		}
		infoAttachList := make([]map[string]interface{}, 0, len(attachList))
		for _, v := range attachList {
			infoAttachList = append(infoAttachList, map[string]interface{}{
				"service_id":   v.ServiceId,
				"service_name": v.ServiceName,
				"api_id":       v.ApiId,
				"api_name":     v.ApiName,
				"path":         v.Path,
				"method":       v.Method,
				"environment":  v.Environment,
				"modify_time":  v.ModifiedTime,
				"create_time":  v.CreatedTime,
			})
		}
		infoMap["usage_plan_id"] = info.UsagePlanId
		infoMap["usage_plan_name"] = info.UsagePlanName
		infoMap["usage_plan_desc"] = info.UsagePlanDesc
		infoMap["max_request_num"] = info.MaxRequestNum
		infoMap["max_request_num_pre_sec"] = info.MaxRequestNumPreSec
		infoMap["modify_time"] = info.ModifiedTime
		infoMap["create_time"] = info.CreatedTime
		infoMap["attach_list"] = infoAttachList

		list = append(list, infoMap)
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
	}

	d.SetId(strings.Join([]string{usagePlanId, usagePlanName}, FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil
}

/*
Use this data source to query detailed information of apigateway api_app

Example Usage

```hcl
data "tencentcloud_apigateway_api_app" "api_app" {
  service_id = ""
  a_p_i_ids =
  filters {
		name = ""
		values =

  }
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

func dataSourceTencentCloudApigatewayApiApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApigatewayApiAppRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service ID.",
			},

			"a_p_i_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of API IDs.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions. Supports ApiAppId, Environment, KeyWord (can match name or ID).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter value of the field.",
						},
					},
				},
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of APIs bound by the application.Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Quantity.",
						},
						"api_app_api_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Application bound API information arrayNote: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_app_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application NameNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"api_app_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application ID.",
									},
									"a_p_i_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API IDNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"a_p_i_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API nameNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service IDNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"authorized_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authorization binding time, expressed in accordance with the ISO8601 standard and using UTC time. The format is: YYYY-MM-DDThh:mm:ssZ.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"api_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api&amp;#39;s regionNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"environment_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authorization binding environmentNote: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudApigatewayApiAppRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_apigateway_api_app.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("a_p_i_ids"); ok {
		aPIIdsSet := v.(*schema.Set).List()
		paramMap["APIIds"] = helper.InterfacesStringsPoint(aPIIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*apigateway.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := apigateway.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*apigateway.ApiAppApiInfos

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApigatewayApiAppByFilter(ctx, paramMap)
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
		apiAppApiInfosMap := map[string]interface{}{}

		if result.TotalCount != nil {
			apiAppApiInfosMap["total_count"] = result.TotalCount
		}

		if result.ApiAppApiSet != nil {
			apiAppApiSetList := []interface{}{}
			for _, apiAppApiSet := range result.ApiAppApiSet {
				apiAppApiSetMap := map[string]interface{}{}

				if apiAppApiSet.ApiAppName != nil {
					apiAppApiSetMap["api_app_name"] = apiAppApiSet.ApiAppName
				}

				if apiAppApiSet.ApiAppId != nil {
					apiAppApiSetMap["api_app_id"] = apiAppApiSet.ApiAppId
				}

				if apiAppApiSet.APIId != nil {
					apiAppApiSetMap["a_p_i_id"] = apiAppApiSet.APIId
				}

				if apiAppApiSet.APIName != nil {
					apiAppApiSetMap["a_p_i_name"] = apiAppApiSet.APIName
				}

				if apiAppApiSet.ServiceId != nil {
					apiAppApiSetMap["service_id"] = apiAppApiSet.ServiceId
				}

				if apiAppApiSet.AuthorizedTime != nil {
					apiAppApiSetMap["authorized_time"] = apiAppApiSet.AuthorizedTime
				}

				if apiAppApiSet.ApiRegion != nil {
					apiAppApiSetMap["api_region"] = apiAppApiSet.ApiRegion
				}

				if apiAppApiSet.EnvironmentName != nil {
					apiAppApiSetMap["environment_name"] = apiAppApiSet.EnvironmentName
				}

				apiAppApiSetList = append(apiAppApiSetList, apiAppApiSetMap)
			}

			apiAppApiInfosMap["api_app_api_set"] = []interface{}{apiAppApiSetList}
		}

		ids = append(ids, *result.ApiAppId)
		_ = d.Set("result", apiAppApiInfosMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), apiAppApiInfosMap); e != nil {
			return e
		}
	}
	return nil
}

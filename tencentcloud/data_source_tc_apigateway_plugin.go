/*
Use this data source to query detailed information of apigateway plugin

Example Usage

```hcl
data "tencentcloud_apigateway_plugin" "plugin" {
  service_id = ""
  plugin_id = ""
  environment_name = ""
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

func dataSourceTencentCloudApigatewayPlugin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApigatewayPluginRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The service ID to be queried.",
			},

			"plugin_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The plug-in ID to be queried.",
			},

			"environment_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment information.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of plug-in related APIs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of plug-in related APIs.",
						},
						"api_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "API information related to the plug-in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"a_p_i_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API ID.",
									},
									"a_p_i_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API name.",
									},
									"api_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API type.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API path.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API method.",
									},
									"attached_other_plugin": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the API is bound to other plug-ins.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"is_attached": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the API is bound to the current plug-in.Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudApigatewayPluginRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_apigateway_plugin.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plugin_id"); ok {
		paramMap["PluginId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("environment_name"); ok {
		paramMap["EnvironmentName"] = helper.String(v.(string))
	}

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*apigateway.ApiInfoSummary

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApigatewayPluginByFilter(ctx, paramMap)
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
		apiInfoSummaryMap := map[string]interface{}{}

		if result.TotalCount != nil {
			apiInfoSummaryMap["total_count"] = result.TotalCount
		}

		if result.ApiSet != nil {
			apiSetList := []interface{}{}
			for _, apiSet := range result.ApiSet {
				apiSetMap := map[string]interface{}{}

				if apiSet.APIId != nil {
					apiSetMap["a_p_i_id"] = apiSet.APIId
				}

				if apiSet.APIName != nil {
					apiSetMap["a_p_i_name"] = apiSet.APIName
				}

				if apiSet.ApiType != nil {
					apiSetMap["api_type"] = apiSet.ApiType
				}

				if apiSet.Path != nil {
					apiSetMap["path"] = apiSet.Path
				}

				if apiSet.Method != nil {
					apiSetMap["method"] = apiSet.Method
				}

				if apiSet.AttachedOtherPlugin != nil {
					apiSetMap["attached_other_plugin"] = apiSet.AttachedOtherPlugin
				}

				if apiSet.IsAttached != nil {
					apiSetMap["is_attached"] = apiSet.IsAttached
				}

				apiSetList = append(apiSetList, apiSetMap)
			}

			apiInfoSummaryMap["api_set"] = []interface{}{apiSetList}
		}

		ids = append(ids, *result.PluginId)
		_ = d.Set("result", apiInfoSummaryMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), apiInfoSummaryMap); e != nil {
			return e
		}
	}
	return nil
}

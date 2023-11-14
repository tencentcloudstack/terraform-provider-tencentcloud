/*
Use this data source to query detailed information of apigateway upstream

Example Usage

```hcl
data "tencentcloud_apigateway_upstream" "upstream" {
  upstream_id = ""
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

func dataSourceTencentCloudApigatewayUpstream() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApigatewayUpstreamRead,
		Schema: map[string]*schema.Schema{
			"upstream_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backend channel ID.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "ServiceId and ApiId filter query.",
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
				Description: "Query results.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total.",
						},
						"bind_api_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Bound API information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"a_p_i_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api unique id.",
									},
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service unique id.",
									},
									"a_p_i_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api nameNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service nameNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"bind_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Binding time.",
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

func dataSourceTencentCloudApigatewayUpstreamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_apigateway_upstream.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("upstream_id"); ok {
		paramMap["UpstreamId"] = helper.String(v.(string))
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

	var result []*apigateway.DescribeUpstreamBindApis

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApigatewayUpstreamByFilter(ctx, paramMap)
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
		describeUpstreamBindApisMap := map[string]interface{}{}

		if result.TotalCount != nil {
			describeUpstreamBindApisMap["total_count"] = result.TotalCount
		}

		if result.BindApiSet != nil {
			bindApiSetList := []interface{}{}
			for _, bindApiSet := range result.BindApiSet {
				bindApiSetMap := map[string]interface{}{}

				if bindApiSet.APIId != nil {
					bindApiSetMap["a_p_i_id"] = bindApiSet.APIId
				}

				if bindApiSet.ServiceId != nil {
					bindApiSetMap["service_id"] = bindApiSet.ServiceId
				}

				if bindApiSet.APIName != nil {
					bindApiSetMap["a_p_i_name"] = bindApiSet.APIName
				}

				if bindApiSet.ServiceName != nil {
					bindApiSetMap["service_name"] = bindApiSet.ServiceName
				}

				if bindApiSet.BindTime != nil {
					bindApiSetMap["bind_time"] = bindApiSet.BindTime
				}

				bindApiSetList = append(bindApiSetList, bindApiSetMap)
			}

			describeUpstreamBindApisMap["bind_api_set"] = []interface{}{bindApiSetList}
		}

		ids = append(ids, *result.UpstreamId)
		_ = d.Set("result", describeUpstreamBindApisMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), describeUpstreamBindApisMap); e != nil {
			return e
		}
	}
	return nil
}

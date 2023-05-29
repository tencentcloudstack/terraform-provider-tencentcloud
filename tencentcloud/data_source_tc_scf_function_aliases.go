/*
Use this data source to query detailed information of scf function_aliases

Example Usage

```hcl
data "tencentcloud_scf_function_aliases" "function_aliases" {
  function_name = "keep-1676351130"
  namespace     = "default"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudScfFunctionAliases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfFunctionAliasesRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace.",
			},

			"function_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "If this parameter is provided, only aliases associated with this function version will be returned.",
			},

			"aliases": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Alias list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Master version pointed to by the alias.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alias name.",
						},
						"routing_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Routing information of aliasNote: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"additional_version_weights": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Additional version with random weight-based routing.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Function version name.",
												},
												"weight": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Version weight.",
												},
											},
										},
									},
									"addition_version_matchs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Additional version with rule-based routing.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Function version name.",
												},
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Matching rule key. When the API is called, pass in the `key` to route the request to the specified version based on the matching ruleHeader method:Enter invoke.headers.User for `key` and pass in `RoutingKey:{User:value}` when invoking a function through `invoke` for invocation based on rule matching.",
												},
												"method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match method. Valid values:range: range matchexact: exact string match.",
												},
												"expression": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Rule requirements for range match:It should be described in an open or closed range, i.e., `(a,b)` or `[a,b]`, where both a and b are integersRule requirements for exact match:Exact string match.",
												},
											},
										},
									},
								},
							},
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DescriptionNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timeNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"mod_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update timeNote: this field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudScfFunctionAliasesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_function_aliases.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_version"); ok {
		paramMap["FunctionVersion"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var aliases []*scf.Alias

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfFunctionAliasesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		aliases = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(aliases))
	tmpList := make([]map[string]interface{}, 0, len(aliases))

	if aliases != nil {
		for _, alias := range aliases {
			aliasMap := map[string]interface{}{}

			if alias.FunctionVersion != nil {
				aliasMap["function_version"] = alias.FunctionVersion
			}

			if alias.Name != nil {
				aliasMap["name"] = alias.Name
			}

			if alias.RoutingConfig != nil {
				routingConfigMap := map[string]interface{}{}

				if alias.RoutingConfig.AdditionalVersionWeights != nil {
					additionalVersionWeightsList := []interface{}{}
					for _, additionalVersionWeights := range alias.RoutingConfig.AdditionalVersionWeights {
						additionalVersionWeightsMap := map[string]interface{}{}

						if additionalVersionWeights.Version != nil {
							additionalVersionWeightsMap["version"] = additionalVersionWeights.Version
						}

						if additionalVersionWeights.Weight != nil {
							additionalVersionWeightsMap["weight"] = additionalVersionWeights.Weight
						}

						additionalVersionWeightsList = append(additionalVersionWeightsList, additionalVersionWeightsMap)
					}

					routingConfigMap["additional_version_weights"] = additionalVersionWeightsList
				}

				if alias.RoutingConfig.AddtionVersionMatchs != nil {
					addtionVersionMatchsList := []interface{}{}
					for _, addtionVersionMatchs := range alias.RoutingConfig.AddtionVersionMatchs {
						addtionVersionMatchsMap := map[string]interface{}{}

						if addtionVersionMatchs.Version != nil {
							addtionVersionMatchsMap["version"] = addtionVersionMatchs.Version
						}

						if addtionVersionMatchs.Key != nil {
							addtionVersionMatchsMap["key"] = addtionVersionMatchs.Key
						}

						if addtionVersionMatchs.Method != nil {
							addtionVersionMatchsMap["method"] = addtionVersionMatchs.Method
						}

						if addtionVersionMatchs.Expression != nil {
							addtionVersionMatchsMap["expression"] = addtionVersionMatchs.Expression
						}

						addtionVersionMatchsList = append(addtionVersionMatchsList, addtionVersionMatchsMap)
					}

					routingConfigMap["addition_version_matchs"] = addtionVersionMatchsList
				}

				aliasMap["routing_config"] = []interface{}{routingConfigMap}
			}

			if alias.Description != nil {
				aliasMap["description"] = alias.Description
			}

			if alias.AddTime != nil {
				aliasMap["add_time"] = alias.AddTime
			}

			if alias.ModTime != nil {
				aliasMap["mod_time"] = alias.ModTime
			}

			ids = append(ids, *alias.Name)
			tmpList = append(tmpList, aliasMap)
		}

		_ = d.Set("aliases", tmpList)
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

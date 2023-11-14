/*
Use this data source to query detailed information of tse gateway_canary_rules

Example Usage

```hcl
data "tencentcloud_tse_gateway_canary_rules" "gateway_canary_rules" {
  gateway_id = "gateway-xxxxxx"
  service_id = "451a9920-e67a-4519-af41-fccac0e72005"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseGatewayCanaryRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayCanaryRulesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service ID.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Canary rule configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"canary_rule_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Canary rule list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"priority": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Priority. The value ranges from 0 to 100; the larger the value, the higher the priority; the priority cannot be repeated between different rules.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The status of canary rule.",
									},
									"condition_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Parameter matching condition list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Type.Reference value:- path- method- query- header- cookie- body- system.",
												},
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Parameter name.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Operator.Reference value:- le- eq- lt- ne- ge- gt- regex- exists- in- not in- prefix- exact- regex.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Parameter value.",
												},
												"delimiter": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Delimiter. valid when operator is in or not in, reference value:- ,- ;- n.",
												},
												"global_config_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Global configuration ID.",
												},
												"global_config_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Global configuration name.",
												},
											},
										},
									},
									"balanced_service_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Service weight configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"service_i_d": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Service ID.",
												},
												"service_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Service name.",
												},
												"upstream_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Upstream name.",
												},
												"percent": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Percent, 10 is 10%，valid values:0 to 100.",
												},
											},
										},
									},
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service ID.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service name.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total count.",
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

func dataSourceTencentCloudTseGatewayCanaryRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_gateway_canary_rules.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tse.CloudAPIGatewayCanaryRuleList

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseGatewayCanaryRulesByFilter(ctx, paramMap)
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
		cloudAPIGatewayCanaryRuleListMap := map[string]interface{}{}

		if result.CanaryRuleList != nil {
			canaryRuleListList := []interface{}{}
			for _, canaryRuleList := range result.CanaryRuleList {
				canaryRuleListMap := map[string]interface{}{}

				if canaryRuleList.Priority != nil {
					canaryRuleListMap["priority"] = canaryRuleList.Priority
				}

				if canaryRuleList.Enabled != nil {
					canaryRuleListMap["enabled"] = canaryRuleList.Enabled
				}

				if canaryRuleList.ConditionList != nil {
					conditionListList := []interface{}{}
					for _, conditionList := range canaryRuleList.ConditionList {
						conditionListMap := map[string]interface{}{}

						if conditionList.Type != nil {
							conditionListMap["type"] = conditionList.Type
						}

						if conditionList.Key != nil {
							conditionListMap["key"] = conditionList.Key
						}

						if conditionList.Operator != nil {
							conditionListMap["operator"] = conditionList.Operator
						}

						if conditionList.Value != nil {
							conditionListMap["value"] = conditionList.Value
						}

						if conditionList.Delimiter != nil {
							conditionListMap["delimiter"] = conditionList.Delimiter
						}

						if conditionList.GlobalConfigId != nil {
							conditionListMap["global_config_id"] = conditionList.GlobalConfigId
						}

						if conditionList.GlobalConfigName != nil {
							conditionListMap["global_config_name"] = conditionList.GlobalConfigName
						}

						conditionListList = append(conditionListList, conditionListMap)
					}

					canaryRuleListMap["condition_list"] = []interface{}{conditionListList}
				}

				if canaryRuleList.BalancedServiceList != nil {
					balancedServiceListList := []interface{}{}
					for _, balancedServiceList := range canaryRuleList.BalancedServiceList {
						balancedServiceListMap := map[string]interface{}{}

						if balancedServiceList.ServiceID != nil {
							balancedServiceListMap["service_i_d"] = balancedServiceList.ServiceID
						}

						if balancedServiceList.ServiceName != nil {
							balancedServiceListMap["service_name"] = balancedServiceList.ServiceName
						}

						if balancedServiceList.UpstreamName != nil {
							balancedServiceListMap["upstream_name"] = balancedServiceList.UpstreamName
						}

						if balancedServiceList.Percent != nil {
							balancedServiceListMap["percent"] = balancedServiceList.Percent
						}

						balancedServiceListList = append(balancedServiceListList, balancedServiceListMap)
					}

					canaryRuleListMap["balanced_service_list"] = []interface{}{balancedServiceListList}
				}

				if canaryRuleList.ServiceId != nil {
					canaryRuleListMap["service_id"] = canaryRuleList.ServiceId
				}

				if canaryRuleList.ServiceName != nil {
					canaryRuleListMap["service_name"] = canaryRuleList.ServiceName
				}

				canaryRuleListList = append(canaryRuleListList, canaryRuleListMap)
			}

			cloudAPIGatewayCanaryRuleListMap["canary_rule_list"] = []interface{}{canaryRuleListList}
		}

		if result.TotalCount != nil {
			cloudAPIGatewayCanaryRuleListMap["total_count"] = result.TotalCount
		}

		ids = append(ids, *result.GatewayId)
		_ = d.Set("result", cloudAPIGatewayCanaryRuleListMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), cloudAPIGatewayCanaryRuleListMap); e != nil {
			return e
		}
	}
	return nil
}

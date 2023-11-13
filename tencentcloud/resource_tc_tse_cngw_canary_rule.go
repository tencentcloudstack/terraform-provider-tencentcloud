/*
Provides a resource to create a tse cngw_canary_rule

Example Usage

```hcl
resource "tencentcloud_tse_cngw_canary_rule" "cngw_canary_rule" {
  gateway_id = "gateway-xxxxxx"
  service_id = "451a9920-e67a-4519-af41-fccac0e72005"
  canary_rule {
		priority = 10
		enabled = true
		condition_list {
			type = ""
			key = ""
			operator = ""
			value = ""
			delimiter = ""
			global_config_id = ""
			global_config_name = ""
		}
		balanced_service_list {
			service_i_d = ""
			service_name = ""
			upstream_name = ""
			percent =
		}
		service_id = ""
		service_name = ""

  }
}
```

Import

tse cngw_canary_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_canary_rule.cngw_canary_rule cngw_canary_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTseCngwCanaryRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwCanaryRuleCreate,
		Read:   resourceTencentCloudTseCngwCanaryRuleRead,
		Update: resourceTencentCloudTseCngwCanaryRuleUpdate,
		Delete: resourceTencentCloudTseCngwCanaryRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

			"canary_rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Canary rule configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Priority. The value ranges from 0 to 100; the larger the value, the higher the priority; the priority cannot be repeated between different rules.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "The status of canary rule.",
						},
						"condition_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Parameter matching condition list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type.Reference value:- path- method- query- header- cookie- body- system.",
									},
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Parameter name.",
									},
									"operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Operator.Reference value:- le- eq- lt- ne- ge- gt- regex- exists- in- not in- prefix- exact- regex.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Parameter value.",
									},
									"delimiter": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Delimiter. valid when operator is in or not in, reference value:- ,- ;- n.",
									},
									"global_config_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Global configuration ID.",
									},
									"global_config_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Global configuration name.",
									},
								},
							},
						},
						"balanced_service_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Service weight configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_i_d": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Service ID, required when used as an input parameter.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Service name, meaningless when used as an input parameter.",
									},
									"upstream_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Upstream name, meaningless when used as an input parameter.",
									},
									"percent": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Percent, 10 is 10%，valid values:0 to 100.",
									},
								},
							},
						},
						"service_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service ID.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service name.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTseCngwCanaryRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_canary_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayCanaryRuleRequest()
		response  = tse.NewCreateCloudNativeAPIGatewayCanaryRuleResponse()
		gatewayId string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		request.ServiceId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "canary_rule"); ok {
		cloudNativeAPIGatewayCanaryRule := tse.CloudNativeAPIGatewayCanaryRule{}
		if v, ok := dMap["priority"]; ok {
			cloudNativeAPIGatewayCanaryRule.Priority = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["enabled"]; ok {
			cloudNativeAPIGatewayCanaryRule.Enabled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["condition_list"]; ok {
			for _, item := range v.([]interface{}) {
				conditionListMap := item.(map[string]interface{})
				cloudNativeAPIGatewayCanaryRuleCondition := tse.CloudNativeAPIGatewayCanaryRuleCondition{}
				if v, ok := conditionListMap["type"]; ok {
					cloudNativeAPIGatewayCanaryRuleCondition.Type = helper.String(v.(string))
				}
				if v, ok := conditionListMap["key"]; ok {
					cloudNativeAPIGatewayCanaryRuleCondition.Key = helper.String(v.(string))
				}
				if v, ok := conditionListMap["operator"]; ok {
					cloudNativeAPIGatewayCanaryRuleCondition.Operator = helper.String(v.(string))
				}
				if v, ok := conditionListMap["value"]; ok {
					cloudNativeAPIGatewayCanaryRuleCondition.Value = helper.String(v.(string))
				}
				if v, ok := conditionListMap["delimiter"]; ok {
					cloudNativeAPIGatewayCanaryRuleCondition.Delimiter = helper.String(v.(string))
				}
				if v, ok := conditionListMap["global_config_id"]; ok {
					cloudNativeAPIGatewayCanaryRuleCondition.GlobalConfigId = helper.String(v.(string))
				}
				if v, ok := conditionListMap["global_config_name"]; ok {
					cloudNativeAPIGatewayCanaryRuleCondition.GlobalConfigName = helper.String(v.(string))
				}
				cloudNativeAPIGatewayCanaryRule.ConditionList = append(cloudNativeAPIGatewayCanaryRule.ConditionList, &cloudNativeAPIGatewayCanaryRuleCondition)
			}
		}
		if v, ok := dMap["balanced_service_list"]; ok {
			for _, item := range v.([]interface{}) {
				balancedServiceListMap := item.(map[string]interface{})
				cloudNativeAPIGatewayBalancedService := tse.CloudNativeAPIGatewayBalancedService{}
				if v, ok := balancedServiceListMap["service_i_d"]; ok {
					cloudNativeAPIGatewayBalancedService.ServiceID = helper.String(v.(string))
				}
				if v, ok := balancedServiceListMap["service_name"]; ok {
					cloudNativeAPIGatewayBalancedService.ServiceName = helper.String(v.(string))
				}
				if v, ok := balancedServiceListMap["upstream_name"]; ok {
					cloudNativeAPIGatewayBalancedService.UpstreamName = helper.String(v.(string))
				}
				if v, ok := balancedServiceListMap["percent"]; ok {
					cloudNativeAPIGatewayBalancedService.Percent = helper.Float64(v.(float64))
				}
				cloudNativeAPIGatewayCanaryRule.BalancedServiceList = append(cloudNativeAPIGatewayCanaryRule.BalancedServiceList, &cloudNativeAPIGatewayBalancedService)
			}
		}
		if v, ok := dMap["service_id"]; ok {
			cloudNativeAPIGatewayCanaryRule.ServiceId = helper.String(v.(string))
		}
		if v, ok := dMap["service_name"]; ok {
			cloudNativeAPIGatewayCanaryRule.ServiceName = helper.String(v.(string))
		}
		request.CanaryRule = &cloudNativeAPIGatewayCanaryRule
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateCloudNativeAPIGatewayCanaryRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwCanaryRule failed, reason:%+v", logId, err)
		return err
	}

	gatewayId = *response.Response.GatewayId
	d.SetId(gatewayId)

	return resourceTencentCloudTseCngwCanaryRuleRead(d, meta)
}

func resourceTencentCloudTseCngwCanaryRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_canary_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	cngwCanaryRuleId := d.Id()

	cngwCanaryRule, err := service.DescribeTseCngwCanaryRuleById(ctx, gatewayId)
	if err != nil {
		return err
	}

	if cngwCanaryRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwCanaryRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cngwCanaryRule.GatewayId != nil {
		_ = d.Set("gateway_id", cngwCanaryRule.GatewayId)
	}

	if cngwCanaryRule.ServiceId != nil {
		_ = d.Set("service_id", cngwCanaryRule.ServiceId)
	}

	if cngwCanaryRule.CanaryRule != nil {
		canaryRuleMap := map[string]interface{}{}

		if cngwCanaryRule.CanaryRule.Priority != nil {
			canaryRuleMap["priority"] = cngwCanaryRule.CanaryRule.Priority
		}

		if cngwCanaryRule.CanaryRule.Enabled != nil {
			canaryRuleMap["enabled"] = cngwCanaryRule.CanaryRule.Enabled
		}

		if cngwCanaryRule.CanaryRule.ConditionList != nil {
			conditionListList := []interface{}{}
			for _, conditionList := range cngwCanaryRule.CanaryRule.ConditionList {
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

			canaryRuleMap["condition_list"] = []interface{}{conditionListList}
		}

		if cngwCanaryRule.CanaryRule.BalancedServiceList != nil {
			balancedServiceListList := []interface{}{}
			for _, balancedServiceList := range cngwCanaryRule.CanaryRule.BalancedServiceList {
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

			canaryRuleMap["balanced_service_list"] = []interface{}{balancedServiceListList}
		}

		if cngwCanaryRule.CanaryRule.ServiceId != nil {
			canaryRuleMap["service_id"] = cngwCanaryRule.CanaryRule.ServiceId
		}

		if cngwCanaryRule.CanaryRule.ServiceName != nil {
			canaryRuleMap["service_name"] = cngwCanaryRule.CanaryRule.ServiceName
		}

		_ = d.Set("canary_rule", []interface{}{canaryRuleMap})
	}

	return nil
}

func resourceTencentCloudTseCngwCanaryRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_canary_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tse.NewModifyCloudNativeAPIGatewayCanaryRuleRequest()

	cngwCanaryRuleId := d.Id()

	request.GatewayId = &gatewayId

	immutableArgs := []string{"gateway_id", "service_id", "canary_rule"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("gateway_id") {
		if v, ok := d.GetOk("gateway_id"); ok {
			request.GatewayId = helper.String(v.(string))
		}
	}

	if d.HasChange("service_id") {
		if v, ok := d.GetOk("service_id"); ok {
			request.ServiceId = helper.String(v.(string))
		}
	}

	if d.HasChange("canary_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "canary_rule"); ok {
			cloudNativeAPIGatewayCanaryRule := tse.CloudNativeAPIGatewayCanaryRule{}
			if v, ok := dMap["priority"]; ok {
				cloudNativeAPIGatewayCanaryRule.Priority = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["enabled"]; ok {
				cloudNativeAPIGatewayCanaryRule.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := dMap["condition_list"]; ok {
				for _, item := range v.([]interface{}) {
					conditionListMap := item.(map[string]interface{})
					cloudNativeAPIGatewayCanaryRuleCondition := tse.CloudNativeAPIGatewayCanaryRuleCondition{}
					if v, ok := conditionListMap["type"]; ok {
						cloudNativeAPIGatewayCanaryRuleCondition.Type = helper.String(v.(string))
					}
					if v, ok := conditionListMap["key"]; ok {
						cloudNativeAPIGatewayCanaryRuleCondition.Key = helper.String(v.(string))
					}
					if v, ok := conditionListMap["operator"]; ok {
						cloudNativeAPIGatewayCanaryRuleCondition.Operator = helper.String(v.(string))
					}
					if v, ok := conditionListMap["value"]; ok {
						cloudNativeAPIGatewayCanaryRuleCondition.Value = helper.String(v.(string))
					}
					if v, ok := conditionListMap["delimiter"]; ok {
						cloudNativeAPIGatewayCanaryRuleCondition.Delimiter = helper.String(v.(string))
					}
					if v, ok := conditionListMap["global_config_id"]; ok {
						cloudNativeAPIGatewayCanaryRuleCondition.GlobalConfigId = helper.String(v.(string))
					}
					if v, ok := conditionListMap["global_config_name"]; ok {
						cloudNativeAPIGatewayCanaryRuleCondition.GlobalConfigName = helper.String(v.(string))
					}
					cloudNativeAPIGatewayCanaryRule.ConditionList = append(cloudNativeAPIGatewayCanaryRule.ConditionList, &cloudNativeAPIGatewayCanaryRuleCondition)
				}
			}
			if v, ok := dMap["balanced_service_list"]; ok {
				for _, item := range v.([]interface{}) {
					balancedServiceListMap := item.(map[string]interface{})
					cloudNativeAPIGatewayBalancedService := tse.CloudNativeAPIGatewayBalancedService{}
					if v, ok := balancedServiceListMap["service_i_d"]; ok {
						cloudNativeAPIGatewayBalancedService.ServiceID = helper.String(v.(string))
					}
					if v, ok := balancedServiceListMap["service_name"]; ok {
						cloudNativeAPIGatewayBalancedService.ServiceName = helper.String(v.(string))
					}
					if v, ok := balancedServiceListMap["upstream_name"]; ok {
						cloudNativeAPIGatewayBalancedService.UpstreamName = helper.String(v.(string))
					}
					if v, ok := balancedServiceListMap["percent"]; ok {
						cloudNativeAPIGatewayBalancedService.Percent = helper.Float64(v.(float64))
					}
					cloudNativeAPIGatewayCanaryRule.BalancedServiceList = append(cloudNativeAPIGatewayCanaryRule.BalancedServiceList, &cloudNativeAPIGatewayBalancedService)
				}
			}
			if v, ok := dMap["service_id"]; ok {
				cloudNativeAPIGatewayCanaryRule.ServiceId = helper.String(v.(string))
			}
			if v, ok := dMap["service_name"]; ok {
				cloudNativeAPIGatewayCanaryRule.ServiceName = helper.String(v.(string))
			}
			request.CanaryRule = &cloudNativeAPIGatewayCanaryRule
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().ModifyCloudNativeAPIGatewayCanaryRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwCanaryRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwCanaryRuleRead(d, meta)
}

func resourceTencentCloudTseCngwCanaryRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_canary_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	cngwCanaryRuleId := d.Id()

	if err := service.DeleteTseCngwCanaryRuleById(ctx, gatewayId); err != nil {
		return err
	}

	return nil
}

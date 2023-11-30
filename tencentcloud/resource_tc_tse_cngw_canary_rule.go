package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Description: "gateway ID.",
			},

			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "service ID.",
			},

			"canary_rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "canary rule configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "priority. The value ranges from 0 to 100; the larger the value, the higher the priority; the priority cannot be repeated between different rules.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "the status of canary rule.",
						},
						"condition_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "parameter matching condition list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "type.Reference value:`path`,`method`,`query`,`header`,`cookie`,`body`,`system`.",
									},
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "parameter name.",
									},
									"operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "operator.Reference value:`le`,`eq`,`lt`,`ne`,`ge`,`gt`,`regex`,`exists`,`in`,`not in`,`prefix`,`exact`,`regex`.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "parameter value.",
									},
									"delimiter": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "delimiter. valid when operator is in or not in, reference value:`,`, `;`,`\\n`.",
									},
									"global_config_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "global configuration ID.",
									},
									"global_config_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "global configuration name.",
									},
								},
							},
						},
						"balanced_service_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "service weight configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "service ID, required when used as an input parameter.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "service name, meaningless when used as an input parameter.",
									},
									"upstream_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "upstream name, meaningless when used as an input parameter.",
									},
									"percent": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "percent, 10 is 10%, valid values:0 to 100.",
									},
								},
							},
						},
						"service_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "service ID.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "service name.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
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
		gatewayId string
		serviceId string
		priority  int
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		serviceId = v.(string)
		request.ServiceId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "canary_rule"); ok {
		cloudNativeAPIGatewayCanaryRule := tse.CloudNativeAPIGatewayCanaryRule{}
		if v, ok := dMap["priority"]; ok {
			priority = v.(int)
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
				if v, ok := balancedServiceListMap["service_id"]; ok {
					cloudNativeAPIGatewayBalancedService.ServiceID = helper.String(v.(string))
				}
				if v, ok := balancedServiceListMap["service_name"]; ok {
					cloudNativeAPIGatewayBalancedService.ServiceName = helper.String(v.(string))
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
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwCanaryRule failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(gatewayId + FILED_SP + serviceId + FILED_SP + strconv.Itoa(priority))

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tse:%s:uin/:cngw_canary_rule/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTseCngwCanaryRuleRead(d, meta)
}

func resourceTencentCloudTseCngwCanaryRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_canary_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	serviceId := idSplit[1]
	priority := idSplit[2]

	cngwCanaryRule, err := service.DescribeTseCngwCanaryRuleById(ctx, gatewayId, serviceId, priority)
	if err != nil {
		return err
	}

	if cngwCanaryRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwCanaryRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)
	_ = d.Set("service_id", serviceId)

	if cngwCanaryRule != nil {
		canaryRuleMap := map[string]interface{}{}

		if cngwCanaryRule.Priority != nil {
			canaryRuleMap["priority"] = cngwCanaryRule.Priority
		}

		if cngwCanaryRule.Enabled != nil {
			canaryRuleMap["enabled"] = cngwCanaryRule.Enabled
		}

		if cngwCanaryRule.ConditionList != nil {
			conditionListList := []interface{}{}
			for _, conditionList := range cngwCanaryRule.ConditionList {
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

			canaryRuleMap["condition_list"] = conditionListList
		}

		if cngwCanaryRule.BalancedServiceList != nil {
			balancedServiceListList := []interface{}{}
			for _, balancedServiceList := range cngwCanaryRule.BalancedServiceList {
				balancedServiceListMap := map[string]interface{}{}

				if balancedServiceList.ServiceID != nil {
					balancedServiceListMap["service_id"] = balancedServiceList.ServiceID
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

			canaryRuleMap["balanced_service_list"] = balancedServiceListList
		}

		if cngwCanaryRule.ServiceId != nil {
			canaryRuleMap["service_id"] = cngwCanaryRule.ServiceId
		}

		if cngwCanaryRule.ServiceName != nil {
			canaryRuleMap["service_name"] = cngwCanaryRule.ServiceName
		}

		_ = d.Set("canary_rule", []interface{}{canaryRuleMap})
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tse", "cngw_canary_rule", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTseCngwCanaryRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_canary_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tse.NewModifyCloudNativeAPIGatewayCanaryRuleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	serviceId := idSplit[1]
	priority := idSplit[2]

	priorityInt64, err := strconv.ParseInt(priority, 10, 64)
	if err != nil {
		return err
	}

	request.GatewayId = &gatewayId
	request.ServiceId = &serviceId
	request.Priority = &priorityInt64

	immutableArgs := []string{"gateway_id", "service_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
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
				if v, ok := balancedServiceListMap["service_id"]; ok {
					cloudNativeAPIGatewayBalancedService.ServiceID = helper.String(v.(string))
				}
				if v, ok := balancedServiceListMap["service_name"]; ok {
					cloudNativeAPIGatewayBalancedService.ServiceName = helper.String(v.(string))
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

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
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

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tse", "cngw_canary_rule", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTseCngwCanaryRuleRead(d, meta)
}

func resourceTencentCloudTseCngwCanaryRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_canary_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	serviceId := idSplit[1]
	priority := idSplit[2]

	if err := service.DeleteTseCngwCanaryRuleById(ctx, gatewayId, serviceId, priority); err != nil {
		return err
	}

	return nil
}

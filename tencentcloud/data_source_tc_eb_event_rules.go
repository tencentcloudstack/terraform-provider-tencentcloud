package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEbEventRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbEventRulesRead,
		Schema: map[string]*schema.Schema{
			"event_bus_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "event bus Id.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "According to which field to sort the returned results, the following fields are supported: AddTime (creation time), ModTime (modification time).",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Return results in ascending or descending order, optional values ASC (ascending) and DESC (descending).",
			},

			"rules": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Event rule information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status.",
						},
						"mod_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "modify time.",
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "enable switch.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "description.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "rule Id.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"event_bus_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "event bus Id.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "rule name.",
						},
						"targets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Target brief information, note: this field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "target Id.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "target type.",
									},
								},
							},
						},
						"dead_letter_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The dlq rule set by rule. It may be null. Note: this field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dispose_method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Support three modes of dlq, discarding, ignoring errors and continuing to pass, corresponding to: DLQ, DROP, IGNORE_ERROR.",
									},
									"ckafka_delivery_params": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "After setting the DLQ mode, this option is required. The error message will be delivered to the corresponding kafka topic Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"topic_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ckafka topic name.",
												},
												"resource_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ckafka resource qcs six-segment.",
												},
											},
										},
									},
								},
							},
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

func dataSourceTencentCloudEbEventRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eb_event_rules.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("event_bus_id"); ok {
		paramMap["EventBusId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var rules []*eb.Rule

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbEventRulesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		rules = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(rules))
	tmpList := make([]map[string]interface{}, 0, len(rules))

	if rules != nil {
		for _, rule := range rules {
			ruleMap := map[string]interface{}{}

			if rule.Status != nil {
				ruleMap["status"] = rule.Status
			}

			if rule.ModTime != nil {
				ruleMap["mod_time"] = rule.ModTime
			}

			if rule.Enable != nil {
				ruleMap["enable"] = rule.Enable
			}

			if rule.Description != nil {
				ruleMap["description"] = rule.Description
			}

			if rule.RuleId != nil {
				ruleMap["rule_id"] = rule.RuleId
			}

			if rule.AddTime != nil {
				ruleMap["add_time"] = rule.AddTime
			}

			if rule.EventBusId != nil {
				ruleMap["event_bus_id"] = rule.EventBusId
			}

			if rule.RuleName != nil {
				ruleMap["rule_name"] = rule.RuleName
			}

			if rule.Targets != nil {
				targetsList := []interface{}{}
				for _, targets := range rule.Targets {
					targetsMap := map[string]interface{}{}

					if targets.TargetId != nil {
						targetsMap["target_id"] = targets.TargetId
					}

					if targets.Type != nil {
						targetsMap["type"] = targets.Type
					}

					targetsList = append(targetsList, targetsMap)
				}

				ruleMap["targets"] = []interface{}{targetsList}
			}

			if rule.DeadLetterConfig != nil {
				deadLetterConfigMap := map[string]interface{}{}

				if rule.DeadLetterConfig.DisposeMethod != nil {
					deadLetterConfigMap["dispose_method"] = rule.DeadLetterConfig.DisposeMethod
				}

				if rule.DeadLetterConfig.CkafkaDeliveryParams != nil {
					ckafkaDeliveryParamsMap := map[string]interface{}{}

					if rule.DeadLetterConfig.CkafkaDeliveryParams.TopicName != nil {
						ckafkaDeliveryParamsMap["topic_name"] = rule.DeadLetterConfig.CkafkaDeliveryParams.TopicName
					}

					if rule.DeadLetterConfig.CkafkaDeliveryParams.ResourceDescription != nil {
						ckafkaDeliveryParamsMap["resource_description"] = rule.DeadLetterConfig.CkafkaDeliveryParams.ResourceDescription
					}

					deadLetterConfigMap["ckafka_delivery_params"] = []interface{}{ckafkaDeliveryParamsMap}
				}

				ruleMap["dead_letter_config"] = []interface{}{deadLetterConfigMap}
			}

			ids = append(ids, *rule.EventBusId)
			tmpList = append(tmpList, ruleMap)
		}

		_ = d.Set("rules", tmpList)
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

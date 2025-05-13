package waf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafBotSceneUCBRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafBotSceneUCBRuleCreate,
		Read:   resourceTencentCloudWafBotSceneUCBRuleRead,
		Update: resourceTencentCloudWafBotSceneUCBRuleUpdate,
		Delete: resourceTencentCloudWafBotSceneUCBRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"scene_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "When calling at the BOT global whitelist, pass `global`; When configuring BOT scenarios, transmit the specific scenario ID.",
			},

			"rule": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Rule content, add encoding SceneId information. When calling at the BOT global whitelist, SceneId is set to `global` and RuleType is passed as 10, Action is `permit`; When configuring BOT scenarios, SceneId is the scenario ID.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain.",
						},

						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule name.",
						},

						"rule": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Specific rule items of UCB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key.",
									},

									"op": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Operator.",
									},

									"value": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Value.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"basic_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "String type value.",
												},

												"logic_value": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Bool type value.",
												},

												"belong_value": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "String array type value.",
												},

												"valid_key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Indicate valid fields.",
												},

												"multi_value": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "String array type value.",
												},
											},
										},
									},

									"op_op": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional Supplementary Operators.",
									},

									"op_arg": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Optional supplementary parameters.",
									},

									"op_value": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Optional supplementary values.",
									},

									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When using header parameter values.",
									},

									"areas": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Regional selection.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"country": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "In addition to standard countries, the country also supports two special identifiers: domestic and foreign.",
												},

												"region": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Province.",
												},

												"city": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "City.",
												},
											},
										},
									},

									"lang": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Language environment.",
									},
								},
							},
						},

						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Disposal action.",
						},

						"on_off": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule switch.",
						},

						"rule_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule type.",
						},

						"prior": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule priority.",
						},

						"timestamp": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Modifying timestamps.",
						},

						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Label.",
						},

						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Entry ID.",
						},

						"scene_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Scene ID.",
						},

						"valid_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Valid time.",
						},

						"appid": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Appid.",
						},
						"addition_arg": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Additional parameters.",
						},

						"desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule description.",
						},

						"rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule ID.",
						},

						"pre_define": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "True - System preset rules False - Custom rules.",
						},

						"job_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Scheduled task type.",
						},

						"job_date_time": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Scheduled task configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timed": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Time parameter for timed execution.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_date_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Start timestamp, in seconds.",
												},

												"end_date_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "End timestamp, in seconds.",
												},
											},
										},
									},

									"cron": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Time parameter for cycle execution.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "On what day of each month is it executed.",
												},

												"w_days": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "What day of the week is executed each week.",
												},

												"start_time": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Start time.",
												},
												"end_time": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "End time.",
												},
											},
										},
									},

									"time_t_zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Time zone.",
									},
								},
							},
						},

						"expire_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Effective deadline.",
						},

						"valid_status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Effective -1, Invalid -0.",
						},

						"block_page_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Customize interception page ID.",
						},

						"action_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "When Action=intercept, this field is mandatory.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action.",
									},

									"proportion": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Proportion.",
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

func resourceTencentCloudWafBotSceneUCBRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_scene_ucb_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = waf.NewModifyBotSceneUCBRuleRequest()
		response = waf.NewModifyBotSceneUCBRuleResponse()
		domain   string
		sceneId  string
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("scene_id"); ok {
		request.SceneId = helper.String(v.(string))
		sceneId = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
		inOutputBotUCBRule := waf.InOutputBotUCBRule{}
		if v, ok := dMap["domain"]; ok {
			inOutputBotUCBRule.Domain = helper.String(v.(string))
		}

		if v, ok := dMap["name"]; ok {
			inOutputBotUCBRule.Name = helper.String(v.(string))
		}

		if v, ok := dMap["rule"]; ok {
			for _, item := range v.([]interface{}) {
				if ruleMap, ok := item.(map[string]interface{}); ok && ruleMap != nil {
					inOutputUCBRuleEntry := waf.InOutputUCBRuleEntry{}
					if v, ok := ruleMap["key"]; ok {
						inOutputUCBRuleEntry.Key = helper.String(v.(string))
					}

					if v, ok := ruleMap["op"]; ok {
						inOutputUCBRuleEntry.Op = helper.String(v.(string))
					}

					if valueMap, ok := helper.InterfaceToMap(ruleMap, "value"); ok {
						uCBEntryValue := waf.UCBEntryValue{}
						if v, ok := valueMap["basic_value"]; ok {
							uCBEntryValue.BasicValue = helper.String(v.(string))
						}

						if v, ok := valueMap["logic_value"]; ok {
							uCBEntryValue.LogicValue = helper.Bool(v.(bool))
						}

						if v, ok := valueMap["belong_value"]; ok {
							belongValueSet := v.(*schema.Set).List()
							for i := range belongValueSet {
								if belongValueSet[i] != nil {
									belongValue := belongValueSet[i].(string)
									uCBEntryValue.BelongValue = append(uCBEntryValue.BelongValue, &belongValue)
								}
							}
						}

						if v, ok := valueMap["valid_key"]; ok {
							uCBEntryValue.ValidKey = helper.String(v.(string))
						}

						if v, ok := valueMap["multi_value"]; ok {
							multiValueSet := v.(*schema.Set).List()
							for i := range multiValueSet {
								if multiValueSet[i] != nil {
									multiValue := multiValueSet[i].(string)
									uCBEntryValue.MultiValue = append(uCBEntryValue.MultiValue, &multiValue)
								}
							}
						}

						inOutputUCBRuleEntry.Value = &uCBEntryValue
					}

					if v, ok := ruleMap["op_op"]; ok {
						inOutputUCBRuleEntry.OpOp = helper.String(v.(string))
					}

					if v, ok := ruleMap["op_arg"]; ok {
						opArgSet := v.(*schema.Set).List()
						for i := range opArgSet {
							if opArgSet[i] != nil {
								opArg := opArgSet[i].(string)
								inOutputUCBRuleEntry.OpArg = append(inOutputUCBRuleEntry.OpArg, &opArg)
							}
						}
					}

					if v, ok := ruleMap["op_value"]; ok {
						inOutputUCBRuleEntry.OpValue = helper.Float64(v.(float64))
					}

					if v, ok := ruleMap["name"]; ok {
						inOutputUCBRuleEntry.Name = helper.String(v.(string))
					}

					if v, ok := ruleMap["areas"]; ok {
						for _, item := range v.([]interface{}) {
							areasMap := item.(map[string]interface{})
							area := waf.Area{}
							if v, ok := areasMap["country"]; ok {
								area.Country = helper.String(v.(string))
							}

							if v, ok := areasMap["region"]; ok {
								area.Region = helper.String(v.(string))
							}

							if v, ok := areasMap["city"]; ok {
								area.City = helper.String(v.(string))
							}

							inOutputUCBRuleEntry.Areas = append(inOutputUCBRuleEntry.Areas, &area)
						}
					}

					if v, ok := ruleMap["lang"]; ok {
						inOutputUCBRuleEntry.Lang = helper.String(v.(string))
					}

					inOutputBotUCBRule.Rule = append(inOutputBotUCBRule.Rule, &inOutputUCBRuleEntry)
				}
			}
		}

		if v, ok := dMap["action"]; ok {
			inOutputBotUCBRule.Action = helper.String(v.(string))
		}

		if v, ok := dMap["on_off"]; ok {
			inOutputBotUCBRule.OnOff = helper.String(v.(string))
		}

		if v, ok := dMap["rule_type"]; ok {
			inOutputBotUCBRule.RuleType = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["prior"]; ok {
			inOutputBotUCBRule.Prior = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["timestamp"]; ok {
			inOutputBotUCBRule.Timestamp = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["label"]; ok {
			inOutputBotUCBRule.Label = helper.String(v.(string))
		}

		if v, ok := dMap["id"]; ok {
			inOutputBotUCBRule.Id = helper.String(v.(string))
		}

		if v, ok := dMap["scene_id"]; ok {
			inOutputBotUCBRule.SceneId = helper.String(v.(string))
		}

		if v, ok := dMap["valid_time"]; ok {
			inOutputBotUCBRule.ValidTime = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["appid"]; ok {
			inOutputBotUCBRule.Appid = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["addition_arg"]; ok {
			inOutputBotUCBRule.AdditionArg = helper.String(v.(string))
		}

		if v, ok := dMap["desc"]; ok {
			inOutputBotUCBRule.Desc = helper.String(v.(string))
		}

		if v, ok := dMap["rule_id"]; ok {
			inOutputBotUCBRule.RuleId = helper.String(v.(string))
		}

		if v, ok := dMap["pre_define"]; ok {
			inOutputBotUCBRule.PreDefine = helper.Bool(v.(bool))
		}

		if v, ok := dMap["job_type"]; ok {
			inOutputBotUCBRule.JobType = helper.String(v.(string))
		}

		if jobDateTimeMap, ok := helper.InterfaceToMap(dMap, "job_date_time"); ok {
			jobDateTime := waf.JobDateTime{}
			if v, ok := jobDateTimeMap["timed"]; ok {
				for _, item := range v.([]interface{}) {
					if timedMap, ok := item.(map[string]interface{}); ok && timedMap != nil {
						timedJob := waf.TimedJob{}
						if v, ok := timedMap["start_date_time"]; ok {
							timedJob.StartDateTime = helper.IntUint64(v.(int))
						}

						if v, ok := timedMap["end_date_time"]; ok {
							timedJob.EndDateTime = helper.IntUint64(v.(int))
						}

						jobDateTime.Timed = append(jobDateTime.Timed, &timedJob)
					}
				}
			}

			if v, ok := jobDateTimeMap["cron"]; ok {
				for _, item := range v.([]interface{}) {
					if cronMap, ok := item.(map[string]interface{}); ok && cronMap != nil {
						cronJob := waf.CronJob{}
						if v, ok := cronMap["days"]; ok {
							daysSet := v.(*schema.Set).List()
							for i := range daysSet {
								days := daysSet[i].(int)
								cronJob.Days = append(cronJob.Days, helper.IntUint64(days))
							}
						}

						if v, ok := cronMap["w_days"]; ok {
							wDaysSet := v.(*schema.Set).List()
							for i := range wDaysSet {
								wDays := wDaysSet[i].(int)
								cronJob.WDays = append(cronJob.WDays, helper.IntUint64(wDays))
							}
						}

						if v, ok := cronMap["start_time"]; ok {
							cronJob.StartTime = helper.String(v.(string))
						}

						if v, ok := cronMap["end_time"]; ok {
							cronJob.EndTime = helper.String(v.(string))
						}

						jobDateTime.Cron = append(jobDateTime.Cron, &cronJob)
					}
				}
			}

			if v, ok := jobDateTimeMap["time_t_zone"]; ok {
				jobDateTime.TimeTZone = helper.String(v.(string))
			}

			inOutputBotUCBRule.JobDateTime = &jobDateTime
		}

		if v, ok := dMap["expire_time"]; ok {
			inOutputBotUCBRule.ExpireTime = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["valid_status"]; ok {
			inOutputBotUCBRule.ValidStatus = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["block_page_id"]; ok {
			inOutputBotUCBRule.BlockPageId = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["action_list"]; ok {
			for _, item := range v.([]interface{}) {
				if actionListMap, ok := item.(map[string]interface{}); ok && actionListMap != nil {
					uCBActionProportion := waf.UCBActionProportion{}
					if v, ok := actionListMap["action"]; ok {
						uCBActionProportion.Action = helper.String(v.(string))
					}

					if v, ok := actionListMap["proportion"]; ok {
						uCBActionProportion.Proportion = helper.Float64(v.(float64))
					}

					inOutputBotUCBRule.ActionList = append(inOutputBotUCBRule.ActionList, &uCBActionProportion)
				}
			}
		}

		request.Rule = &inOutputBotUCBRule
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyBotSceneUCBRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create waf bot scene ucb rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf bot scene ucb rule failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	ruleId := ""
	d.SetId(strings.Join([]string{domain, sceneId, ruleId}, tccommon.FILED_SP))

	return resourceTencentCloudWafBotSceneUCBRuleRead(d, meta)
}

func resourceTencentCloudWafBotSceneUCBRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_scene_ucb_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	sceneId := idSplit[1]
	ruleId := idSplit[2]

	respData, err := service.DescribeWafBotSceneUCBRuleById(ctx, domain, sceneId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafBotSceneUCBRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	if respData.SceneId != nil {
		_ = d.Set("scene_id", respData.SceneId)
	}

	// if waf_bot_scene_ucb_rule.Rule != nil {
	// 	ruleMap := map[string]interface{}{}

	// 	if waf_bot_scene_ucb_rule.Rule.Domain != nil {
	// 		ruleMap["domain"] = waf_bot_scene_ucb_rule.Rule.Domain
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.Name != nil {
	// 		ruleMap["name"] = waf_bot_scene_ucb_rule.Rule.Name
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.Rule != nil {
	// 		ruleList := []interface{}{}
	// 		for _, rule := range waf_bot_scene_ucb_rule.Rule.Rule {
	// 			ruleMap := map[string]interface{}{}

	// 			if rule.Key != nil {
	// 				ruleMap["key"] = rule.Key
	// 			}

	// 			if rule.Op != nil {
	// 				ruleMap["op"] = rule.Op
	// 			}

	// 			if rule.Value != nil {
	// 				valueMap := map[string]interface{}{}

	// 				if rule.Value.BasicValue != nil {
	// 					valueMap["basic_value"] = rule.Value.BasicValue
	// 				}

	// 				if rule.Value.LogicValue != nil {
	// 					valueMap["logic_value"] = rule.Value.LogicValue
	// 				}

	// 				if rule.Value.BelongValue != nil {
	// 					valueMap["belong_value"] = rule.Value.BelongValue
	// 				}

	// 				if rule.Value.ValidKey != nil {
	// 					valueMap["valid_key"] = rule.Value.ValidKey
	// 				}

	// 				if rule.Value.MultiValue != nil {
	// 					valueMap["multi_value"] = rule.Value.MultiValue
	// 				}

	// 				ruleMap["value"] = []interface{}{valueMap}
	// 			}

	// 			if rule.OpOp != nil {
	// 				ruleMap["op_op"] = rule.OpOp
	// 			}

	// 			if rule.OpArg != nil {
	// 				ruleMap["op_arg"] = rule.OpArg
	// 			}

	// 			if rule.OpValue != nil {
	// 				ruleMap["op_value"] = rule.OpValue
	// 			}

	// 			if rule.Name != nil {
	// 				ruleMap["name"] = rule.Name
	// 			}

	// 			if rule.Areas != nil {
	// 				areasList := []interface{}{}
	// 				for _, areas := range rule.Areas {
	// 					areasMap := map[string]interface{}{}

	// 					if areas.Country != nil {
	// 						areasMap["country"] = areas.Country
	// 					}

	// 					if areas.Region != nil {
	// 						areasMap["region"] = areas.Region
	// 					}

	// 					if areas.City != nil {
	// 						areasMap["city"] = areas.City
	// 					}

	// 					areasList = append(areasList, areasMap)
	// 				}

	// 				ruleMap["areas"] = areasList
	// 			}

	// 			if rule.Lang != nil {
	// 				ruleMap["lang"] = rule.Lang
	// 			}

	// 			ruleList = append(ruleList, ruleMap)
	// 		}

	// 		ruleMap["rule"] = ruleList
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.Action != nil {
	// 		ruleMap["action"] = waf_bot_scene_ucb_rule.Rule.Action
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.OnOff != nil {
	// 		ruleMap["on_off"] = waf_bot_scene_ucb_rule.Rule.OnOff
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.RuleType != nil {
	// 		ruleMap["rule_type"] = waf_bot_scene_ucb_rule.Rule.RuleType
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.Prior != nil {
	// 		ruleMap["prior"] = waf_bot_scene_ucb_rule.Rule.Prior
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.Timestamp != nil {
	// 		ruleMap["timestamp"] = waf_bot_scene_ucb_rule.Rule.Timestamp
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.Label != nil {
	// 		ruleMap["label"] = waf_bot_scene_ucb_rule.Rule.Label
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.Id != nil {
	// 		ruleMap["id"] = waf_bot_scene_ucb_rule.Rule.Id
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.SceneId != nil {
	// 		ruleMap["scene_id"] = waf_bot_scene_ucb_rule.Rule.SceneId
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.ValidTime != nil {
	// 		ruleMap["valid_time"] = waf_bot_scene_ucb_rule.Rule.ValidTime
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.Appid != nil {
	// 		ruleMap["appid"] = waf_bot_scene_ucb_rule.Rule.Appid
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.AdditionArg != nil {
	// 		ruleMap["addition_arg"] = waf_bot_scene_ucb_rule.Rule.AdditionArg
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.Desc != nil {
	// 		ruleMap["desc"] = waf_bot_scene_ucb_rule.Rule.Desc
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.RuleId != nil {
	// 		ruleMap["rule_id"] = waf_bot_scene_ucb_rule.Rule.RuleId
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.PreDefine != nil {
	// 		ruleMap["pre_define"] = waf_bot_scene_ucb_rule.Rule.PreDefine
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.JobType != nil {
	// 		ruleMap["job_type"] = waf_bot_scene_ucb_rule.Rule.JobType
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.JobDateTime != nil {
	// 		jobDateTimeMap := map[string]interface{}{}

	// 		if waf_bot_scene_ucb_rule.Rule.JobDateTime.Timed != nil {
	// 			timedList := []interface{}{}
	// 			for _, timed := range waf_bot_scene_ucb_rule.Rule.JobDateTime.Timed {
	// 				timedMap := map[string]interface{}{}

	// 				if timed.StartDateTime != nil {
	// 					timedMap["start_date_time"] = timed.StartDateTime
	// 				}

	// 				if timed.EndDateTime != nil {
	// 					timedMap["end_date_time"] = timed.EndDateTime
	// 				}

	// 				timedList = append(timedList, timedMap)
	// 			}

	// 			jobDateTimeMap["timed"] = timedList
	// 		}

	// 		if waf_bot_scene_ucb_rule.Rule.JobDateTime.Cron != nil {
	// 			cronList := []interface{}{}
	// 			for _, cron := range waf_bot_scene_ucb_rule.Rule.JobDateTime.Cron {
	// 				cronMap := map[string]interface{}{}

	// 				if cron.Days != nil {
	// 					cronMap["days"] = cron.Days
	// 				}

	// 				if cron.WDays != nil {
	// 					cronMap["w_days"] = cron.WDays
	// 				}

	// 				if cron.StartTime != nil {
	// 					cronMap["start_time"] = cron.StartTime
	// 				}

	// 				if cron.EndTime != nil {
	// 					cronMap["end_time"] = cron.EndTime
	// 				}

	// 				cronList = append(cronList, cronMap)
	// 			}

	// 			jobDateTimeMap["cron"] = cronList
	// 		}

	// 		if waf_bot_scene_ucb_rule.Rule.JobDateTime.TimeTZone != nil {
	// 			jobDateTimeMap["time_t_zone"] = waf_bot_scene_ucb_rule.Rule.JobDateTime.TimeTZone
	// 		}

	// 		ruleMap["job_date_time"] = []interface{}{jobDateTimeMap}
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.ExpireTime != nil {
	// 		ruleMap["expire_time"] = waf_bot_scene_ucb_rule.Rule.ExpireTime
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.ValidStatus != nil {
	// 		ruleMap["valid_status"] = waf_bot_scene_ucb_rule.Rule.ValidStatus
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.BlockPageId != nil {
	// 		ruleMap["block_page_id"] = waf_bot_scene_ucb_rule.Rule.BlockPageId
	// 	}

	// 	if waf_bot_scene_ucb_rule.Rule.ActionList != nil {
	// 		actionListList := []interface{}{}
	// 		for _, actionList := range waf_bot_scene_ucb_rule.Rule.ActionList {
	// 			actionListMap := map[string]interface{}{}

	// 			if actionList.Action != nil {
	// 				actionListMap["action"] = actionList.Action
	// 			}

	// 			if actionList.Proportion != nil {
	// 				actionListMap["proportion"] = actionList.Proportion
	// 			}

	// 			actionListList = append(actionListList, actionListMap)
	// 		}

	// 		ruleMap["action_list"] = actionListList
	// 	}

	// 	_ = d.Set("rule", []interface{}{ruleMap})
	// }

	return nil
}

func resourceTencentCloudWafBotSceneUCBRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_scene_ucb_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = waf.NewModifyBotSceneUCBRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	sceneId := idSplit[1]
	ruleId := idSplit[2]

	request.Domain = &domain
	request.SceneId = &sceneId

	if dMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
		inOutputBotUCBRule := waf.InOutputBotUCBRule{}
		if v, ok := dMap["domain"]; ok {
			inOutputBotUCBRule.Domain = helper.String(v.(string))
		}

		if v, ok := dMap["name"]; ok {
			inOutputBotUCBRule.Name = helper.String(v.(string))
		}

		if v, ok := dMap["rule"]; ok {
			for _, item := range v.([]interface{}) {
				ruleMap := item.(map[string]interface{})
				inOutputUCBRuleEntry := waf.InOutputUCBRuleEntry{}
				if v, ok := ruleMap["key"]; ok {
					inOutputUCBRuleEntry.Key = helper.String(v.(string))
				}

				if v, ok := ruleMap["op"]; ok {
					inOutputUCBRuleEntry.Op = helper.String(v.(string))
				}

				if valueMap, ok := helper.InterfaceToMap(ruleMap, "value"); ok {
					uCBEntryValue := waf.UCBEntryValue{}
					if v, ok := valueMap["basic_value"]; ok {
						uCBEntryValue.BasicValue = helper.String(v.(string))
					}

					if v, ok := valueMap["logic_value"]; ok {
						uCBEntryValue.LogicValue = helper.Bool(v.(bool))
					}

					if v, ok := valueMap["belong_value"]; ok {
						belongValueSet := v.(*schema.Set).List()
						for i := range belongValueSet {
							if belongValueSet[i] != nil {
								belongValue := belongValueSet[i].(string)
								uCBEntryValue.BelongValue = append(uCBEntryValue.BelongValue, &belongValue)
							}
						}
					}

					if v, ok := valueMap["valid_key"]; ok {
						uCBEntryValue.ValidKey = helper.String(v.(string))
					}

					if v, ok := valueMap["multi_value"]; ok {
						multiValueSet := v.(*schema.Set).List()
						for i := range multiValueSet {
							if multiValueSet[i] != nil {
								multiValue := multiValueSet[i].(string)
								uCBEntryValue.MultiValue = append(uCBEntryValue.MultiValue, &multiValue)
							}
						}
					}

					inOutputUCBRuleEntry.Value = &uCBEntryValue
				}

				if v, ok := ruleMap["op_op"]; ok {
					inOutputUCBRuleEntry.OpOp = helper.String(v.(string))
				}

				if v, ok := ruleMap["op_arg"]; ok {
					opArgSet := v.(*schema.Set).List()
					for i := range opArgSet {
						if opArgSet[i] != nil {
							opArg := opArgSet[i].(string)
							inOutputUCBRuleEntry.OpArg = append(inOutputUCBRuleEntry.OpArg, &opArg)
						}
					}
				}

				if v, ok := ruleMap["op_value"]; ok {
					inOutputUCBRuleEntry.OpValue = helper.Float64(v.(float64))
				}

				if v, ok := ruleMap["name"]; ok {
					inOutputUCBRuleEntry.Name = helper.String(v.(string))
				}

				if v, ok := ruleMap["areas"]; ok {
					for _, item := range v.([]interface{}) {
						areasMap := item.(map[string]interface{})
						area := waf.Area{}
						if v, ok := areasMap["country"]; ok {
							area.Country = helper.String(v.(string))
						}

						if v, ok := areasMap["region"]; ok {
							area.Region = helper.String(v.(string))
						}

						if v, ok := areasMap["city"]; ok {
							area.City = helper.String(v.(string))
						}

						inOutputUCBRuleEntry.Areas = append(inOutputUCBRuleEntry.Areas, &area)
					}
				}

				if v, ok := ruleMap["lang"]; ok {
					inOutputUCBRuleEntry.Lang = helper.String(v.(string))
				}

				inOutputBotUCBRule.Rule = append(inOutputBotUCBRule.Rule, &inOutputUCBRuleEntry)
			}
		}

		if v, ok := dMap["action"]; ok {
			inOutputBotUCBRule.Action = helper.String(v.(string))
		}

		if v, ok := dMap["on_off"]; ok {
			inOutputBotUCBRule.OnOff = helper.String(v.(string))
		}

		if v, ok := dMap["rule_type"]; ok {
			inOutputBotUCBRule.RuleType = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["prior"]; ok {
			inOutputBotUCBRule.Prior = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["timestamp"]; ok {
			inOutputBotUCBRule.Timestamp = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["label"]; ok {
			inOutputBotUCBRule.Label = helper.String(v.(string))
		}

		if v, ok := dMap["id"]; ok {
			inOutputBotUCBRule.Id = helper.String(v.(string))
		}

		if v, ok := dMap["scene_id"]; ok {
			inOutputBotUCBRule.SceneId = helper.String(v.(string))
		}

		if v, ok := dMap["valid_time"]; ok {
			inOutputBotUCBRule.ValidTime = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["appid"]; ok {
			inOutputBotUCBRule.Appid = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["addition_arg"]; ok {
			inOutputBotUCBRule.AdditionArg = helper.String(v.(string))
		}

		if v, ok := dMap["desc"]; ok {
			inOutputBotUCBRule.Desc = helper.String(v.(string))
		}

		inOutputBotUCBRule.RuleId = &ruleId

		if v, ok := dMap["pre_define"]; ok {
			inOutputBotUCBRule.PreDefine = helper.Bool(v.(bool))
		}

		if v, ok := dMap["job_type"]; ok {
			inOutputBotUCBRule.JobType = helper.String(v.(string))
		}

		if jobDateTimeMap, ok := helper.InterfaceToMap(dMap, "job_date_time"); ok {
			jobDateTime := waf.JobDateTime{}
			if v, ok := jobDateTimeMap["timed"]; ok {
				for _, item := range v.([]interface{}) {
					timedMap := item.(map[string]interface{})
					timedJob := waf.TimedJob{}
					if v, ok := timedMap["start_date_time"]; ok {
						timedJob.StartDateTime = helper.IntUint64(v.(int))
					}

					if v, ok := timedMap["end_date_time"]; ok {
						timedJob.EndDateTime = helper.IntUint64(v.(int))
					}

					jobDateTime.Timed = append(jobDateTime.Timed, &timedJob)
				}
			}

			if v, ok := jobDateTimeMap["cron"]; ok {
				for _, item := range v.([]interface{}) {
					cronMap := item.(map[string]interface{})
					cronJob := waf.CronJob{}
					if v, ok := cronMap["days"]; ok {
						daysSet := v.(*schema.Set).List()
						for i := range daysSet {
							days := daysSet[i].(int)
							cronJob.Days = append(cronJob.Days, helper.IntUint64(days))
						}
					}

					if v, ok := cronMap["w_days"]; ok {
						wDaysSet := v.(*schema.Set).List()
						for i := range wDaysSet {
							wDays := wDaysSet[i].(int)
							cronJob.WDays = append(cronJob.WDays, helper.IntUint64(wDays))
						}
					}

					if v, ok := cronMap["start_time"]; ok {
						cronJob.StartTime = helper.String(v.(string))
					}

					if v, ok := cronMap["end_time"]; ok {
						cronJob.EndTime = helper.String(v.(string))
					}

					jobDateTime.Cron = append(jobDateTime.Cron, &cronJob)
				}
			}

			if v, ok := jobDateTimeMap["time_t_zone"]; ok {
				jobDateTime.TimeTZone = helper.String(v.(string))
			}

			inOutputBotUCBRule.JobDateTime = &jobDateTime
		}

		if v, ok := dMap["expire_time"]; ok {
			inOutputBotUCBRule.ExpireTime = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["valid_status"]; ok {
			inOutputBotUCBRule.ValidStatus = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["block_page_id"]; ok {
			inOutputBotUCBRule.BlockPageId = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["action_list"]; ok {
			for _, item := range v.([]interface{}) {
				actionListMap := item.(map[string]interface{})
				uCBActionProportion := waf.UCBActionProportion{}
				if v, ok := actionListMap["action"]; ok {
					uCBActionProportion.Action = helper.String(v.(string))
				}

				if v, ok := actionListMap["proportion"]; ok {
					uCBActionProportion.Proportion = helper.Float64(v.(float64))
				}

				inOutputBotUCBRule.ActionList = append(inOutputBotUCBRule.ActionList, &uCBActionProportion)
			}
		}

		request.Rule = &inOutputBotUCBRule
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyBotSceneUCBRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update waf bot scene ucb rule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWafBotSceneUCBRuleRead(d, meta)
}

func resourceTencentCloudWafBotSceneUCBRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_scene_ucb_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	sceneId := idSplit[1]
	ruleId := idSplit[2]

	if err := service.DeleteWafBotSceneUCBRuleById(ctx, domain, sceneId, ruleId); err != nil {
		return err
	}

	return nil
}

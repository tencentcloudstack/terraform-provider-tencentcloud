package waf

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

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
													Computed:    true,
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
							Computed:    true,
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
										Computed:    true,
										Description: "Time parameter for timed execution.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_date_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Start timestamp, in seconds.",
												},

												"end_date_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "End timestamp, in seconds.",
												},
											},
										},
									},

									"cron": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Time parameter for cycle execution.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Computed:    true,
													Description: "On what day of each month is it executed.",
												},

												"w_days": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Computed:    true,
													Description: "What day of the week is executed each week.",
												},

												"start_time": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Start time.",
												},

												"end_time": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
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
							Computed:    true,
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

			// computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID.",
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
					var base46Flag bool
					if v, ok := ruleMap["key"]; ok {
						inOutputUCBRuleEntry.Key = helper.String(v.(string))
					}

					if v, ok := ruleMap["op"]; ok {
						inOutputUCBRuleEntry.Op = helper.String(v.(string))
						if v.(string) == "rematch" {
							base46Flag = true
						}
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
							if base46Flag {
								for i := range multiValueSet {
									if multiValueSet[i] != nil {
										multiValue := multiValueSet[i].(string)
										bs64Str := helper.String(base64.URLEncoding.EncodeToString([]byte(multiValue)))
										uCBEntryValue.MultiValue = append(uCBEntryValue.MultiValue, bs64Str)
									}
								}
							} else {
								for i := range multiValueSet {
									if multiValueSet[i] != nil {
										multiValue := multiValueSet[i].(string)
										uCBEntryValue.MultiValue = append(uCBEntryValue.MultiValue, &multiValue)
									}
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
							if areasMap, ok := item.(map[string]interface{}); ok && areasMap != nil {
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

		tmpTime := time.Now().UnixMilli()
		inOutputBotUCBRule.Timestamp = &tmpTime

		if v, ok := dMap["label"]; ok {
			inOutputBotUCBRule.Label = helper.String(v.(string))
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
								if daysSet[i] != nil {
									days := daysSet[i].(int)
									cronJob.Days = append(cronJob.Days, helper.IntUint64(days))
								}
							}
						}

						if v, ok := cronMap["w_days"]; ok {
							wDaysSet := v.(*schema.Set).List()
							for i := range wDaysSet {
								if wDaysSet[i] != nil {
									wDays := wDaysSet[i].(int)
									cronJob.WDays = append(cronJob.WDays, helper.IntUint64(wDays))
								}
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

	if response.Response.RuleIdList == nil || len(response.Response.RuleIdList) < 1 {
		return fmt.Errorf("RuleIdList is nil.")
	}

	ruleId := *response.Response.RuleIdList[0]
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

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	ruleList := make([]map[string]interface{}, 0, 1)
	ruleMap := make(map[string]interface{})
	if respData.Domain != nil {
		ruleMap["domain"] = respData.Domain
	}

	if respData.Name != nil {
		ruleMap["name"] = respData.Name
	}

	if respData.Rule != nil && len(respData.Rule) > 0 {
		tmpList := make([]map[string]interface{}, 0, len(respData.Rule))
		for _, item := range respData.Rule {
			dMap := make(map[string]interface{})
			var base46Flag bool
			if item.Key != nil {
				dMap["key"] = item.Key
			}

			if item.Op != nil {
				dMap["op"] = item.Op
				if *item.Op == "rematch" {
					base46Flag = true
				}
			}

			if item.Value != nil {
				valueList := make([]map[string]interface{}, 0, 1)
				valueMap := make(map[string]interface{})
				if item.Value.BasicValue != nil {
					valueMap["basic_value"] = item.Value.BasicValue
				}

				if item.Value.LogicValue != nil {
					valueMap["logic_value"] = item.Value.LogicValue
				}

				if item.Value.BelongValue != nil {
					valueMap["belong_value"] = item.Value.BelongValue
				}

				if item.Value.ValidKey != nil {
					valueMap["valid_key"] = item.Value.ValidKey
				}

				if item.Value.MultiValue != nil {
					if base46Flag {
						tmpMvList := make([]string, 0, len(item.Value.MultiValue))
						for _, item := range item.Value.MultiValue {
							decoded, e := base64.StdEncoding.DecodeString(*item)
							if e != nil {
								return fmt.Errorf("[%s] base64 decode error: %s", *item, e.Error())
							}

							tmpMvList = append(tmpMvList, string(decoded))
						}

						valueMap["multi_value"] = tmpMvList
					} else {
						valueMap["multi_value"] = item.Value.MultiValue
					}
				}

				valueList = append(valueList, valueMap)
				dMap["value"] = valueList
			}

			if item.OpOp != nil {
				dMap["op_op"] = item.OpOp
			}

			if item.OpArg != nil {
				dMap["op_arg"] = item.OpArg
			}

			if item.OpValue != nil {
				dMap["op_value"] = item.OpValue
			}

			if item.Name != nil {
				dMap["name"] = item.Name
			}

			if item.Areas != nil && len(item.Areas) > 0 {
				areasList := make([]map[string]interface{}, 0, 1)
				areasMap := make(map[string]interface{})
				for _, item := range item.Areas {
					if item.Country != nil {
						areasMap["country"] = item.Country
					}

					if item.Region != nil {
						areasMap["region"] = item.Region
					}

					if item.City != nil {
						areasMap["city"] = item.City
					}

					areasList = append(areasList, areasMap)
				}

				dMap["areas"] = areasList
			}

			if item.Lang != nil {
				dMap["lang"] = item.Lang
			}

			tmpList = append(tmpList, dMap)
		}

		ruleMap["rule"] = tmpList
	}

	if respData.Action != nil {
		ruleMap["action"] = respData.Action
	}

	if respData.OnOff != nil {
		ruleMap["on_off"] = respData.OnOff
	}

	if respData.RuleType != nil {
		ruleMap["rule_type"] = respData.RuleType
	}

	if respData.Prior != nil {
		ruleMap["prior"] = respData.Prior
	}

	if respData.Label != nil {
		ruleMap["label"] = respData.Label
	}

	if respData.Id != nil {
		ruleMap["id"] = respData.Id
	}

	if respData.SceneId != nil {
		ruleMap["scene_id"] = respData.SceneId
	}

	if respData.ValidTime != nil {
		ruleMap["valid_time"] = respData.ValidTime
	}

	if respData.Appid != nil {
		ruleMap["appid"] = respData.Appid
	}

	if respData.AdditionArg != nil {
		ruleMap["addition_arg"] = respData.AdditionArg
	}

	if respData.Desc != nil {
		ruleMap["desc"] = respData.Desc
	}

	if respData.PreDefine != nil {
		ruleMap["pre_define"] = respData.PreDefine
	}

	if respData.JobType != nil {
		ruleMap["job_type"] = respData.JobType
	}

	if respData.JobDateTime != nil {
		jdtList := make([]map[string]interface{}, 0, 1)
		jdtMap := make(map[string]interface{})
		if respData.JobDateTime.Timed != nil && len(respData.JobDateTime.Timed) > 0 {
			tList := make([]map[string]interface{}, 0, len(respData.JobDateTime.Timed))
			for _, item := range respData.JobDateTime.Timed {
				tMap := make(map[string]interface{})
				if item.StartDateTime != nil {
					tMap["start_date_time"] = item.StartDateTime
				}

				if item.EndDateTime != nil {
					tMap["end_date_time"] = item.EndDateTime
				}

				tList = append(tList, tMap)
			}

			jdtMap["timed"] = tList
		}

		if respData.JobDateTime.Cron != nil && len(respData.JobDateTime.Cron) > 0 {
			cList := make([]map[string]interface{}, 0, len(respData.JobDateTime.Cron))
			for _, item := range respData.JobDateTime.Cron {
				cMap := make(map[string]interface{})
				if item.Days != nil {
					cMap["days"] = item.Days
				}

				if item.WDays != nil {
					cMap["w_days"] = item.WDays
				}

				if item.StartTime != nil {
					cMap["start_time"] = item.StartTime
				}

				if item.EndTime != nil {
					cMap["end_time"] = item.EndTime
				}

				cList = append(cList, cMap)
			}

			jdtMap["cron"] = cList
		}

		if respData.JobDateTime.TimeTZone != nil {
			jdtMap["time_t_zone"] = respData.JobDateTime.TimeTZone
		}

		jdtList = append(jdtList, jdtMap)
		ruleMap["job_date_time"] = jdtList
	}

	if respData.ExpireTime != nil {
		ruleMap["expire_time"] = respData.ExpireTime
	}

	if respData.ValidStatus != nil {
		ruleMap["valid_status"] = respData.ValidStatus
	}

	if respData.BlockPageId != nil {
		ruleMap["block_page_id"] = respData.BlockPageId
	}

	if respData.ActionList != nil && len(respData.ActionList) > 0 {
		alList := make([]map[string]interface{}, 0, len(respData.ActionList))
		for _, item := range respData.ActionList {
			alMap := make(map[string]interface{})
			if item.Action != nil {
				alMap["action"] = item.Action
			}

			if item.Proportion != nil {
				alMap["proportion"] = item.Proportion
			}

			alList = append(alList, alMap)
		}

		ruleMap["action_list"] = alList
	}

	ruleList = append(ruleList, ruleMap)
	_ = d.Set("rule", ruleList)

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
				if ruleMap, ok := item.(map[string]interface{}); ok && ruleMap != nil {
					inOutputUCBRuleEntry := waf.InOutputUCBRuleEntry{}
					var base46Flag bool
					if v, ok := ruleMap["key"]; ok {
						inOutputUCBRuleEntry.Key = helper.String(v.(string))
					}

					if v, ok := ruleMap["op"]; ok {
						inOutputUCBRuleEntry.Op = helper.String(v.(string))
						if v.(string) == "rematch" {
							base46Flag = true
						}
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
							if base46Flag {
								for i := range multiValueSet {
									if multiValueSet[i] != nil {
										multiValue := multiValueSet[i].(string)
										bs64Str := helper.String(base64.URLEncoding.EncodeToString([]byte(multiValue)))
										uCBEntryValue.MultiValue = append(uCBEntryValue.MultiValue, bs64Str)
									}
								}
							} else {
								for i := range multiValueSet {
									if multiValueSet[i] != nil {
										multiValue := multiValueSet[i].(string)
										uCBEntryValue.MultiValue = append(uCBEntryValue.MultiValue, &multiValue)
									}
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
							if areasMap, ok := item.(map[string]interface{}); ok && areasMap != nil {
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

		tmpTime := time.Now().UnixMilli()
		inOutputBotUCBRule.Timestamp = &tmpTime

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
								if daysSet[i] != nil {
									days := daysSet[i].(int)
									cronJob.Days = append(cronJob.Days, helper.IntUint64(days))
								}
							}
						}

						if v, ok := cronMap["w_days"]; ok {
							wDaysSet := v.(*schema.Set).List()
							for i := range wDaysSet {
								if wDaysSet[i] != nil {
									wDays := wDaysSet[i].(int)
									cronJob.WDays = append(cronJob.WDays, helper.IntUint64(wDays))
								}
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

package tsf

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTsfLaneRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfLaneRuleCreate,
		Read:   resourceTencentCloudTsfLaneRuleRead,
		Update: resourceTencentCloudTsfLaneRuleUpdate,
		Delete: resourceTencentCloudTsfLaneRuleDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Rule id.",
			},

			"rule_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "lane rule name.",
			},

			"remark": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Lane rule notes.",
			},

			"rule_tag_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "list of swimlane rule labels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "label ID.",
						},
						"tag_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "label name.",
						},
						"tag_operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "label operator.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "tag value.",
						},
						"lane_rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "lane rule ID.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "creation time.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "update time.",
						},
					},
				},
			},

			"rule_tag_relationship": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "lane rule label relationship.",
			},

			"lane_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "lane ID.",
			},

			"priority": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Priority.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "open state, true/false, default: false.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "creation time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "update time.",
			},

			"program_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program id list.",
			},
		},
	}
}

func resourceTencentCloudTsfLaneRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_lane_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = tsf.NewCreateLaneRuleRequest()
		response = tsf.NewCreateLaneRuleResponse()
		ruleId   string
	)
	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_tag_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			laneRuleTag := tsf.LaneRuleTag{}
			if v, ok := dMap["tag_id"]; ok {
				laneRuleTag.TagId = helper.String(v.(string))
			}
			if v, ok := dMap["tag_name"]; ok {
				laneRuleTag.TagName = helper.String(v.(string))
			}
			if v, ok := dMap["tag_operator"]; ok {
				laneRuleTag.TagOperator = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				laneRuleTag.TagValue = helper.String(v.(string))
			}
			if v, ok := dMap["lane_rule_id"]; ok {
				laneRuleTag.LaneRuleId = helper.String(v.(string))
			}
			if v, ok := dMap["create_time"]; ok {
				laneRuleTag.CreateTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["update_time"]; ok {
				laneRuleTag.UpdateTime = helper.IntInt64(v.(int))
			}
			request.RuleTagList = append(request.RuleTagList, &laneRuleTag)
		}
	}

	if v, ok := d.GetOk("rule_tag_relationship"); ok {
		request.RuleTagRelationship = helper.String(v.(string))
	}

	if v, ok := d.GetOk("lane_id"); ok {
		request.LaneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().CreateLaneRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf laneRule failed, reason:%+v", logId, err)
		return err
	}

	ruleId = *response.Response.Result
	d.SetId(ruleId)

	return resourceTencentCloudTsfLaneRuleUpdate(d, meta)
}

func resourceTencentCloudTsfLaneRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_lane_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ruleId := d.Id()

	laneRule, err := service.DescribeTsfLaneRuleById(ctx, ruleId)
	if err != nil {
		return err
	}

	if laneRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfLaneRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if laneRule.RuleId != nil {
		_ = d.Set("rule_id", laneRule.RuleId)
	}

	if laneRule.RuleName != nil {
		_ = d.Set("rule_name", laneRule.RuleName)
	}

	if laneRule.Remark != nil {
		_ = d.Set("remark", laneRule.Remark)
	}

	if laneRule.RuleTagList != nil {
		ruleTagListList := []interface{}{}
		for _, ruleTagList := range laneRule.RuleTagList {
			ruleTagListMap := map[string]interface{}{}

			if ruleTagList.TagId != nil {
				ruleTagListMap["tag_id"] = ruleTagList.TagId
			}

			if ruleTagList.TagName != nil {
				ruleTagListMap["tag_name"] = ruleTagList.TagName
			}

			if ruleTagList.TagOperator != nil {
				ruleTagListMap["tag_operator"] = ruleTagList.TagOperator
			}

			if ruleTagList.TagValue != nil {
				ruleTagListMap["tag_value"] = ruleTagList.TagValue
			}

			if ruleTagList.LaneRuleId != nil {
				ruleTagListMap["lane_rule_id"] = ruleTagList.LaneRuleId
			}

			if ruleTagList.CreateTime != nil {
				ruleTagListMap["create_time"] = ruleTagList.CreateTime
			}

			if ruleTagList.UpdateTime != nil {
				ruleTagListMap["update_time"] = ruleTagList.UpdateTime
			}

			ruleTagListList = append(ruleTagListList, ruleTagListMap)
		}

		_ = d.Set("rule_tag_list", ruleTagListList)

	}

	if laneRule.RuleTagRelationship != nil {
		_ = d.Set("rule_tag_relationship", laneRule.RuleTagRelationship)
	}

	if laneRule.LaneId != nil {
		_ = d.Set("lane_id", laneRule.LaneId)
	}

	if laneRule.Priority != nil {
		_ = d.Set("priority", laneRule.Priority)
	}

	if laneRule.Enable != nil {
		_ = d.Set("enable", laneRule.Enable)
	}

	if laneRule.CreateTime != nil {
		_ = d.Set("create_time", laneRule.CreateTime)
	}

	if laneRule.UpdateTime != nil {
		_ = d.Set("update_time", laneRule.UpdateTime)
	}

	// if laneRule.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", laneRule.ProgramIdList)
	// }

	return nil
}

func resourceTencentCloudTsfLaneRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_lane_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tsf.NewModifyLaneRuleRequest()

	ruleId := d.Id()

	request.RuleId = &ruleId

	immutableArgs := []string{"rule_id", "priority", "create_time", "update_time", "program_id_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	v, _ := d.GetOk("enable")
	request.Enable = helper.Bool(v.(bool))

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_tag_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			laneRuleTag := tsf.LaneRuleTag{}
			if v, ok := dMap["tag_id"]; ok {
				laneRuleTag.TagId = helper.String(v.(string))
			}
			if v, ok := dMap["tag_name"]; ok {
				laneRuleTag.TagName = helper.String(v.(string))
			}
			if v, ok := dMap["tag_operator"]; ok {
				laneRuleTag.TagOperator = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				laneRuleTag.TagValue = helper.String(v.(string))
			}
			if v, ok := dMap["lane_rule_id"]; ok {
				laneRuleTag.LaneRuleId = helper.String(v.(string))
			}
			if v, ok := dMap["create_time"]; ok {
				laneRuleTag.CreateTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["update_time"]; ok {
				laneRuleTag.UpdateTime = helper.IntInt64(v.(int))
			}
			request.RuleTagList = append(request.RuleTagList, &laneRuleTag)
		}
	}

	if v, ok := d.GetOk("rule_tag_relationship"); ok {
		request.RuleTagRelationship = helper.String(v.(string))
	}

	if v, ok := d.GetOk("lane_id"); ok {
		request.LaneId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().ModifyLaneRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf laneRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfLaneRuleRead(d, meta)
}

func resourceTencentCloudTsfLaneRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_lane_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	ruleId := d.Id()

	if err := service.DeleteTsfLaneRuleById(ctx, ruleId); err != nil {
		return err
	}

	return nil
}

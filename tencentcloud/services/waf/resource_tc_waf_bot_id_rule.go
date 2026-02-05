package waf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafBotIdRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafBotIdRuleCreate,
		Read:   resourceTencentCloudWafBotIdRuleRead,
		Update: resourceTencentCloudWafBotIdRuleUpdate,
		Delete: resourceTencentCloudWafBotIdRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name.",
			},

			"scene_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Scene ID.",
			},

			"data": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Configuration information, supports batch processing.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule ID.",
						},
						"status": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Rule switch.",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action configuration.",
						},
						"bot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Rule name.",
						},
						"redirect": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Redirect path.",
						},
					},
				},
			},

			"global_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "0-global settings do not take effect 1-global switch configuration field takes effect 2-global action configuration field takes effect 3-both global switch and action fields take effect 4-only modify global redirect path 5-only modify global protection level.",
			},

			"protect_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Protection level: normal-normal; strict-strict.",
			},
		},
	}
}

func resourceTencentCloudWafBotIdRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_id_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		domain  string
		sceneId string
	)

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	if v, ok := d.GetOk("scene_id"); ok {
		sceneId = v.(string)
	}

	d.SetId(strings.Join([]string{domain, sceneId}, tccommon.FILED_SP))
	return resourceTencentCloudWafBotIdRuleUpdate(d, meta)
}

func resourceTencentCloudWafBotIdRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_id_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	sceneId := idSplit[1]

	respData, err := service.DescribeWafBotIdRuleById(ctx, domain, sceneId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_bot_id_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("scene_id", sceneId)

	if respData.Data != nil {
		dataList := make([]map[string]interface{}, 0, len(respData.Data))
		for _, data := range respData.Data {
			dataMap := map[string]interface{}{}
			if data.RuleId != nil {
				dataMap["rule_id"] = data.RuleId
			}

			if data.Status != nil {
				dataMap["status"] = data.Status
			}

			if data.Action != nil {
				dataMap["action"] = data.Action
			}

			if data.BotId != nil {
				dataMap["bot_id"] = data.BotId
			}

			if data.Redirect != nil {
				dataMap["redirect"] = data.Redirect
			}

			dataList = append(dataList, dataMap)
		}

		_ = d.Set("data", dataList)
	}

	if respData.StatInfo != nil {
		if respData.StatInfo.ProtectLevel != nil {
			_ = d.Set("protect_level", respData.StatInfo.ProtectLevel)
		}
	}

	return nil
}

func resourceTencentCloudWafBotIdRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_id_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	sceneId := idSplit[1]

	needChange := false
	mutableArgs := []string{"data", "global_switch", "status", "rule_action", "global_redirect", "protect_level"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wafv20180125.NewModifyBotIdRuleRequest()
		if v, ok := d.GetOk("data"); ok {
			for _, item := range v.([]interface{}) {
				dataMap := item.(map[string]interface{})
				botIdConfig := wafv20180125.BotIdConfig{}
				if v, ok := dataMap["rule_id"].(string); ok && v != "" {
					botIdConfig.RuleId = helper.String(v)
				}

				if v, ok := dataMap["status"].(bool); ok {
					botIdConfig.Status = helper.Bool(v)
				}

				if v, ok := dataMap["action"].(string); ok && v != "" {
					botIdConfig.Action = helper.String(v)
				}

				if v, ok := dataMap["bot_id"].(string); ok && v != "" {
					botIdConfig.BotId = helper.String(v)
				}

				if v, ok := dataMap["redirect"].(string); ok && v != "" {
					botIdConfig.Redirect = helper.String(v)
				}

				request.Data = append(request.Data, &botIdConfig)
			}
		}

		if v, ok := d.GetOkExists("global_switch"); ok {
			request.GlobalSwitch = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("protect_level"); ok {
			request.ProtectLevel = helper.String(v.(string))
		}

		request.Domain = &domain
		request.SceneId = &sceneId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyBotIdRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf bot id rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWafBotIdRuleRead(d, meta)
}

func resourceTencentCloudWafBotIdRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_bot_id_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

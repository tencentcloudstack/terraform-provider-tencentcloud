package ga2

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGa2ForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGa2ForwardingRuleCreate,
		Read:   resourceTencentCloudGa2ForwardingRuleRead,
		Update: resourceTencentCloudGa2ForwardingRuleUpdate,
		Delete: resourceTencentCloudGa2ForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Global accelerator instance ID.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Listener ID.",
			},
			"forwarding_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Forwarding policy ID.",
			},
			"rule_conditions": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Layer-7 forwarding rule conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_condition_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule condition type.",
						},
						"rule_condition_value": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Rule condition values.",
						},
					},
				},
			},
			"rule_actions": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Layer-7 forwarding rule actions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_action_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule action type.",
						},
						"rule_action_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule action value.",
						},
					},
				},
			},
			"origin_headers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Origin headers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Header key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Header value.",
						},
					},
				},
			},
			"enable_origin_sni": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable origin SNI.",
			},
			"origin_sni": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Origin SNI value.",
			},
			"origin_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Origin host value.",
			},
			"forwarding_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Forwarding rule ID.",
			},
		},
	}
}

func resourceTencentCloudGa2ForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request            = ga2v20250115.NewCreateForwardingRuleRequest()
		response           = ga2v20250115.NewCreateForwardingRuleResponse()
		gaId               string
		listenerId         string
		forwardingPolicyId string
		forwardingRuleId   string
		taskId             string
	)

	if v, ok := d.GetOk("global_accelerator_id"); ok {
		gaId = v.(string)
		request.GlobalAcceleratorId = helper.String(gaId)
	}

	if v, ok := d.GetOk("listener_id"); ok {
		listenerId = v.(string)
		request.ListenerId = helper.String(listenerId)
	}

	if v, ok := d.GetOk("forwarding_policy_id"); ok {
		forwardingPolicyId = v.(string)
		request.ForwardingPolicyId = helper.String(forwardingPolicyId)
	}

	if v, ok := d.GetOk("rule_conditions"); ok {
		request.RuleConditions = buildRuleConditions(v.([]interface{}))
	}

	if v, ok := d.GetOk("rule_actions"); ok {
		request.RuleActions = buildRuleActions(v.([]interface{}))
	}

	if v, ok := d.GetOk("origin_headers"); ok {
		request.OriginHeaders = buildOriginHeaders(v.([]interface{}))
	}

	if v, ok := d.GetOkExists("enable_origin_sni"); ok {
		request.EnableOriginSni = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("origin_sni"); ok {
		request.OriginSni = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_host"); ok {
		request.OriginHost = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateForwardingRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 forwarding rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 forwarding rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ForwardingRuleId == nil {
		return fmt.Errorf("ForwardingRuleId is nil.")
	}
	forwardingRuleId = *response.Response.ForwardingRuleId

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}
	taskId = *response.Response.TaskId

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{gaId, listenerId, forwardingPolicyId, forwardingRuleId}, tccommon.FILED_SP))
	return resourceTencentCloudGa2ForwardingRuleRead(d, meta)
}

func resourceTencentCloudGa2ForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	gaId, listenerId, forwardingPolicyId, forwardingRuleId, err := parseGa2ForwardingRuleId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeGa2ForwardingRuleById(ctx, gaId, listenerId, forwardingPolicyId, forwardingRuleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_forwarding_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("global_accelerator_id", gaId)
	_ = d.Set("listener_id", listenerId)
	_ = d.Set("forwarding_policy_id", forwardingPolicyId)

	if respData.ForwardingRuleId != nil {
		_ = d.Set("forwarding_rule_id", respData.ForwardingRuleId)
	}

	if respData.RuleCondition != nil {
		_ = d.Set("rule_conditions", flattenRuleConditions(respData.RuleCondition))
	}

	if respData.RuleAction != nil {
		_ = d.Set("rule_actions", flattenRuleActions(respData.RuleAction))
	}

	if respData.OriginHeaders != nil {
		_ = d.Set("origin_headers", flattenOriginHeaders(respData.OriginHeaders))
	}

	if respData.EnableOriginSni != nil {
		_ = d.Set("enable_origin_sni", respData.EnableOriginSni)
	}

	if respData.OriginSni != nil {
		_ = d.Set("origin_sni", respData.OriginSni)
	}

	if respData.OriginHost != nil {
		_ = d.Set("origin_host", respData.OriginHost)
	}

	return nil
}

func resourceTencentCloudGa2ForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	gaId, listenerId, forwardingPolicyId, forwardingRuleId, err := parseGa2ForwardingRuleId(d.Id())
	if err != nil {
		return err
	}

	needChange := d.HasChange("rule_conditions") || d.HasChange("rule_actions") || d.HasChange("origin_headers") ||
		d.HasChange("enable_origin_sni") || d.HasChange("origin_sni") || d.HasChange("origin_host")

	if !needChange {
		return resourceTencentCloudGa2ForwardingRuleRead(d, meta)
	}

	request := ga2v20250115.NewModifyForwardingRuleRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)
	request.ForwardingPolicyId = helper.String(forwardingPolicyId)
	request.ForwardingRuleId = helper.String(forwardingRuleId)

	if v, ok := d.GetOk("rule_conditions"); ok {
		request.RuleConditions = buildRuleConditions(v.([]interface{}))
	}

	if v, ok := d.GetOk("rule_actions"); ok {
		request.RuleActions = buildRuleActions(v.([]interface{}))
	}

	if v, ok := d.GetOk("origin_headers"); ok {
		request.OriginHeaders = buildOriginHeaders(v.([]interface{}))
	}

	if v, ok := d.GetOkExists("enable_origin_sni"); ok {
		request.EnableOriginSni = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("origin_sni"); ok {
		request.OriginSni = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_host"); ok {
		request.OriginHost = helper.String(v.(string))
	}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyForwardingRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify ga2 forwarding rule failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update ga2 forwarding rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceTencentCloudGa2ForwardingRuleRead(d, meta)
}

func resourceTencentCloudGa2ForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = ga2v20250115.NewDeleteForwardingRuleRequest()
	)

	gaId, listenerId, forwardingPolicyId, forwardingRuleId, err := parseGa2ForwardingRuleId(d.Id())
	if err != nil {
		return err
	}

	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)
	request.ForwardingPolicyId = helper.String(forwardingPolicyId)
	request.ForwardingRuleId = helper.String(forwardingRuleId)

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteForwardingRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 forwarding rule failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 forwarding rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}

func parseGa2ForwardingRuleId(id string) (gaId, listenerId, forwardingPolicyId, forwardingRuleId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 4 {
		err = fmt.Errorf("invalid resource id %q, expected format <global_accelerator_id>%s<listener_id>%s<forwarding_policy_id>%s<forwarding_rule_id>", id, tccommon.FILED_SP, tccommon.FILED_SP, tccommon.FILED_SP)
		return
	}
	gaId, listenerId, forwardingPolicyId, forwardingRuleId = parts[0], parts[1], parts[2], parts[3]
	if gaId == "" || listenerId == "" || forwardingPolicyId == "" || forwardingRuleId == "" {
		err = fmt.Errorf("invalid resource id %q, components must all be non-empty", id)
		return
	}
	return
}

func buildRuleConditions(rawList []interface{}) []*ga2v20250115.RuleCondition {
	result := make([]*ga2v20250115.RuleCondition, 0, len(rawList))
	for _, item := range rawList {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		rc := &ga2v20250115.RuleCondition{}

		if v, ok := m["rule_condition_type"].(string); ok && v != "" {
			rc.RuleConditionType = helper.String(v)
		}
		if v, ok := m["rule_condition_value"]; ok {
			rc.RuleConditionValue = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		result = append(result, rc)
	}
	return result
}

func flattenRuleConditions(items []*ga2v20250115.RuleCondition) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		m := map[string]interface{}{}
		if item.RuleConditionType != nil {
			m["rule_condition_type"] = *item.RuleConditionType
		}
		if item.RuleConditionValue != nil {
			m["rule_condition_value"] = helper.PStrings(item.RuleConditionValue)
		}
		result = append(result, m)
	}
	return result
}

func buildRuleActions(rawList []interface{}) []*ga2v20250115.RuleAction {
	result := make([]*ga2v20250115.RuleAction, 0, len(rawList))
	for _, item := range rawList {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		ra := &ga2v20250115.RuleAction{}

		if v, ok := m["rule_action_type"].(string); ok && v != "" {
			ra.RuleActionType = helper.String(v)
		}
		if v, ok := m["rule_action_value"].(string); ok && v != "" {
			ra.RuleActionValue = helper.String(v)
		}

		result = append(result, ra)
	}
	return result
}

func flattenRuleActions(items []*ga2v20250115.RuleAction) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		m := map[string]interface{}{}
		if item.RuleActionType != nil {
			m["rule_action_type"] = *item.RuleActionType
		}
		if item.RuleActionValue != nil {
			m["rule_action_value"] = *item.RuleActionValue
		}
		result = append(result, m)
	}
	return result
}

func buildOriginHeaders(rawList []interface{}) []*ga2v20250115.OriginHeader {
	result := make([]*ga2v20250115.OriginHeader, 0, len(rawList))
	for _, item := range rawList {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		oh := &ga2v20250115.OriginHeader{}

		if v, ok := m["key"].(string); ok && v != "" {
			oh.Key = helper.String(v)
		}
		if v, ok := m["value"].(string); ok && v != "" {
			oh.Value = helper.String(v)
		}

		result = append(result, oh)
	}
	return result
}

func flattenOriginHeaders(items []*ga2v20250115.OriginHeader) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		m := map[string]interface{}{}
		if item.Key != nil {
			m["key"] = *item.Key
		}
		if item.Value != nil {
			m["value"] = *item.Value
		}
		result = append(result, m)
	}
	return result
}

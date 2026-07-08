package ga2

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// ResourceTencentCloudGa2ForwardingRule manages a Tencent Cloud GA2 layer-7 forwarding rule.
//
// All write APIs (CreateForwardingRule / ModifyForwardingRule / DeleteForwardingRule) are
// asynchronous: each returns a TaskId that must be polled via DescribeTaskResult until
// Status == "SUCCESS" before this resource considers the operation complete.
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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Global accelerator instance ID this forwarding rule belongs to.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Listener ID this forwarding rule belongs to.",
			},
			"forwarding_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Forwarding policy ID this forwarding rule belongs to.",
			},
			"rule_conditions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Layer-7 forwarding rule condition list. Maximum of 1 element. Treated as an unordered set; HCL element order has no semantic meaning.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_condition_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Layer-7 forwarding rule condition type. Valid values: `Path`.",
						},
						"rule_condition_value": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: "Layer-7 forwarding rule condition values. Each value must match the regular expression " +
								"`^[a-zA-Z0-9_.-/]{1,80}$`. Maximum of 1 element. Treated as an unordered set.",
						},
					},
				},
			},
			"rule_actions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Layer-7 forwarding rule action list. Treated as an unordered set; HCL element order has no semantic meaning.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_action_type": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Layer-7 forwarding rule action type. Valid values: `ForwardGroup` (forward to an endpoint group), " +
								"`Drop` (drop the request).",
						},
						"rule_action_value": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Layer-7 forwarding rule action value. Not required when `rule_action_type` is `Drop`. " +
								"Required when `rule_action_type` is `ForwardGroup`, in which case it must be a custom endpoint group ID " +
								"(the default endpoint group is not supported).",
						},
					},
				},
			},
			"origin_headers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Description: "Origin request header list. Maximum of 5 elements. Required when `rule_actions.rule_action_type` " +
					"is `ForwardGroup`. Treated as an unordered set; HCL element order has no semantic meaning.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin request header key. Must contain only printable ASCII characters and must not contain `()<>@,;:\\\"/[ ]?={}`. Length must be between 1 and 40 characters.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin request header value. Maximum length is 128 characters. If the value contains `$`, only `$remote_addr` or `$remote_port` are supported.",
						},
					},
				},
			},
			"enable_origin_sni": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Whether to enable origin SNI. Default: `false`. Required when `rule_actions.rule_action_type` " +
					"is `ForwardGroup`.",
			},
			"origin_sni": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Origin SNI value. Maximum length is 80 characters. Required when `enable_origin_sni` is `true`, " +
					"and also required when `rule_actions.rule_action_type` is `ForwardGroup`.",
			},
			"origin_host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Origin host value. Maximum length is 80 characters. Required when `rule_actions.rule_action_type` " +
					"is `ForwardGroup`.",
			},

			// Computed
			"forwarding_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Layer-7 forwarding rule ID.",
			},
		},
	}
}

func resourceTencentCloudGa2ForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = ga2v20250115.NewCreateForwardingRuleRequest()
		response   = ga2v20250115.NewCreateForwardingRuleResponse()
		gaId       string
		listenerId string
		policyId   string
		ruleId     string
		taskId     string
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
		policyId = v.(string)
		request.ForwardingPolicyId = helper.String(policyId)
	}

	if v, ok := d.GetOk("rule_conditions"); ok {
		request.RuleConditions = buildGa2ForwardingRuleConditions(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("rule_actions"); ok {
		request.RuleActions = buildGa2ForwardingRuleActions(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("origin_headers"); ok {
		request.OriginHeaders = buildGa2ForwardingRuleOriginHeaders(v.(*schema.Set).List())
	}

	//nolint:staticcheck
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
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

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
	ruleId = *response.Response.ForwardingRuleId

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}
	taskId = *response.Response.TaskId

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{gaId, listenerId, policyId, ruleId}, tccommon.FILED_SP))
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

	gaId, listenerId, policyId, ruleId, err := parseGa2ForwardingRuleId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeGa2ForwardingRuleById(ctx, gaId, listenerId, policyId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_forwarding_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.GlobalAcceleratorId != nil {
		_ = d.Set("global_accelerator_id", respData.GlobalAcceleratorId)
	}

	if respData.ListenerId != nil {
		_ = d.Set("listener_id", respData.ListenerId)
	}

	if respData.ForwardingPolicyId != nil {
		_ = d.Set("forwarding_policy_id", respData.ForwardingPolicyId)
	}

	if respData.ForwardingRuleId != nil {
		_ = d.Set("forwarding_rule_id", respData.ForwardingRuleId)
	}

	// Note: SDK response struct uses singular field names (RuleCondition / RuleAction),
	// while the schema and request structs use the plural form (rule_conditions / rule_actions).
	if len(respData.RuleCondition) > 0 {
		_ = d.Set("rule_conditions", flattenGa2ForwardingRuleConditions(respData.RuleCondition))
	}

	if len(respData.RuleAction) > 0 {
		_ = d.Set("rule_actions", flattenGa2ForwardingRuleActions(respData.RuleAction))
	}

	if len(respData.OriginHeaders) > 0 {
		_ = d.Set("origin_headers", flattenGa2ForwardingRuleOriginHeaders(respData.OriginHeaders))
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

	gaId, listenerId, policyId, ruleId, err := parseGa2ForwardingRuleId(d.Id())
	if err != nil {
		return err
	}

	// Body fields supported by ModifyForwardingRule (the 4 ID fields are ForceNew, so they
	// cannot trigger Update directly).
	bodyFields := []string{
		"rule_conditions", "rule_actions", "origin_headers",
		"enable_origin_sni", "origin_sni", "origin_host",
	}
	needModify := false
	for _, f := range bodyFields {
		if d.HasChange(f) {
			needModify = true
			break
		}
	}

	if !needModify {
		return resourceTencentCloudGa2ForwardingRuleRead(d, meta)
	}

	request := ga2v20250115.NewModifyForwardingRuleRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)
	request.ForwardingPolicyId = helper.String(policyId)
	request.ForwardingRuleId = helper.String(ruleId)

	if v, ok := d.GetOk("rule_conditions"); ok {
		request.RuleConditions = buildGa2ForwardingRuleConditions(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("rule_actions"); ok {
		request.RuleActions = buildGa2ForwardingRuleActions(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("origin_headers"); ok {
		request.OriginHeaders = buildGa2ForwardingRuleOriginHeaders(v.(*schema.Set).List())
	}

	//nolint:staticcheck
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
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

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

	gaId, listenerId, policyId, ruleId, err := parseGa2ForwardingRuleId(d.Id())
	if err != nil {
		return err
	}

	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)
	request.ForwardingPolicyId = helper.String(policyId)
	request.ForwardingRuleId = helper.String(ruleId)

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteForwardingRuleWithContext(ctx, request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound" {
					return nil
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

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

// parseGa2ForwardingRuleId splits the composite resource ID into its four components.
func parseGa2ForwardingRuleId(id string) (gaId, listenerId, policyId, ruleId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 4 {
		err = fmt.Errorf("invalid resource id %q, expected format <global_accelerator_id>%s<listener_id>%s<forwarding_policy_id>%s<forwarding_rule_id>", id, tccommon.FILED_SP, tccommon.FILED_SP, tccommon.FILED_SP)
		return
	}
	gaId, listenerId, policyId, ruleId = parts[0], parts[1], parts[2], parts[3]
	if gaId == "" || listenerId == "" || policyId == "" || ruleId == "" {
		err = fmt.Errorf("invalid resource id %q, components must all be non-empty", id)
		return
	}
	return
}

// buildGa2ForwardingRuleConditions converts the schema set into the SDK RuleCondition slice.
func buildGa2ForwardingRuleConditions(rawList []interface{}) []*ga2v20250115.RuleCondition {
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
		if v, ok := m["rule_condition_value"]; ok && v != nil {
			rc.RuleConditionValue = expandGa2ForwardingRuleStringSet(v.(*schema.Set))
		}

		result = append(result, rc)
	}
	return result
}

// buildGa2ForwardingRuleActions converts the schema set into the SDK RuleAction slice.
func buildGa2ForwardingRuleActions(rawList []interface{}) []*ga2v20250115.RuleAction {
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

// buildGa2ForwardingRuleOriginHeaders converts the schema set into the SDK OriginHeader slice.
func buildGa2ForwardingRuleOriginHeaders(rawList []interface{}) []*ga2v20250115.OriginHeader {
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

// expandGa2ForwardingRuleStringSet converts a TypeSet of strings into a []*string suitable for the SDK.
func expandGa2ForwardingRuleStringSet(set *schema.Set) []*string {
	if set == nil {
		return nil
	}
	raw := set.List()
	if len(raw) == 0 {
		return nil
	}
	result := make([]*string, 0, len(raw))
	for _, item := range raw {
		s, ok := item.(string)
		if !ok || s == "" {
			continue
		}
		v := s
		result = append(result, &v)
	}
	return result
}

// flattenGa2ForwardingRuleConditions maps SDK RuleCondition slice back into the schema set payload.
func flattenGa2ForwardingRuleConditions(items []*ga2v20250115.RuleCondition) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		m := map[string]interface{}{}
		if item.RuleConditionType != nil {
			m["rule_condition_type"] = *item.RuleConditionType
		}
		if len(item.RuleConditionValue) > 0 {
			m["rule_condition_value"] = helper.PStrings(item.RuleConditionValue)
		}
		result = append(result, m)
	}
	return result
}

// flattenGa2ForwardingRuleActions maps SDK RuleAction slice back into the schema set payload.
func flattenGa2ForwardingRuleActions(items []*ga2v20250115.RuleAction) []map[string]interface{} {
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

// flattenGa2ForwardingRuleOriginHeaders maps SDK OriginHeader slice back into the schema set payload.
func flattenGa2ForwardingRuleOriginHeaders(items []*ga2v20250115.OriginHeader) []map[string]interface{} {
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

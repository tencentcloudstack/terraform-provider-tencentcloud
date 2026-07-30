package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var MAX_CREATE_RULES_LEN = 20

func ResourceTencentCloudVpcPrivateNatGatewayTranslationNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleCreate,
		Read:   resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleRead,
		Update: resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleUpdate,
		Delete: resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Private NAT gateway unique ID, such as: `intranat-xxxxxxxx`.",
			},

			"translation_nat_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"local_network_layer_rules",
					"local_transport_layer_rules",
					"peer_network_layer_rules",
				},
				Deprecated:  "It has been deprecated from version 1.82.98, please use local_network_layer_rules / local_transport_layer_rules / peer_network_layer_rules instead. Cannot be used together with the new fields.",
				Description: "(Deprecated) Translation rule object array. Use the typed list fields `local_network_layer_rules`, `local_transport_layer_rules`, and `peer_network_layer_rules` instead. The legacy field continues to work but does not benefit from in-place ModifyPrivateNatGatewayTranslationNatRule support.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"translation_direction": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translation rule target, optional values \"LOCAL\",\"PEER\".",
						},
						"translation_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translation rule type, optional values \"NETWORK_LAYER\",\"TRANSPORT_LAYER\".",
						},
						"translation_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translation IP, when translation rule type is transport layer, it is an IP pool.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translation rule description.",
						},
						"original_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source IP, valid when translation rule type is network layer.",
						},
					},
				},
			},

			"local_network_layer_rules": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"translation_nat_rules"},
				Description:   "Translation rules for the LOCAL direction at the NETWORK_LAYER (three-layer). Identity is keyed by `original_ip` (unique within this bucket). Editing `description` or `translation_ip` is applied in place via ModifyPrivateNatGatewayTranslationNatRule; editing `original_ip` is treated as deleting the old rule and creating a new one.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"translation_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translated (mapped-after) IP for this rule. Can be modified in place.",
						},
						"original_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Original (mapped-before) IP for this rule. Acts as the rule identity within this bucket.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Translation rule description.",
						},
					},
				},
			},

			"local_transport_layer_rules": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"translation_nat_rules"},
				Description:   "Translation rules for the LOCAL direction at the TRANSPORT_LAYER (four-layer). Identity is keyed by `translation_ip` (unique within this bucket). `original_ip` is not applicable for transport-layer rules and is intentionally not exposed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"translation_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translated IP pool for this rule (transport-layer rules use an IP pool). Acts as the rule identity within this bucket.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Translation rule description.",
						},
					},
				},
			},

			"peer_network_layer_rules": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"translation_nat_rules"},
				Description:   "Translation rules for the PEER direction at the NETWORK_LAYER (three-layer). Identity is keyed by `original_ip` (unique within this bucket). Editing `description` or `translation_ip` is applied in place via ModifyPrivateNatGatewayTranslationNatRule; editing `original_ip` is treated as deleting the old rule and creating a new one.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"translation_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translated (mapped-after) IP for this rule. Can be modified in place.",
						},
						"original_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Original (mapped-before) IP for this rule. Acts as the rule identity within this bucket.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Translation rule description.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_nat_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = vpcv20170312.NewCreatePrivateNatGatewayTranslationNatRuleRequest()
		natGatewayId string
	)

	if v, ok := d.GetOk("nat_gateway_id"); ok {
		request.NatGatewayId = helper.String(v.(string))
		natGatewayId = v.(string)
	}

	var allRules []*vpcv20170312.TranslationNatRuleInput
	if v, ok := d.GetOk("translation_nat_rules"); ok {
		for _, item := range v.(*schema.Set).List() {
			allRules = append(allRules, buildLegacyRuleInput(item.(map[string]interface{})))
		}
	}
	allRules = append(allRules, buildTypedRuleInputs(d)...)

	for i := 0; i < len(allRules); i += MAX_CREATE_RULES_LEN {
		end := i + MAX_CREATE_RULES_LEN
		if end > len(allRules) {
			end = len(allRules)
		}

		batchRules := allRules[i:end]
		request.TranslationNatRules = batchRules
		request.NatGatewayId = helper.String(natGatewayId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreatePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil || result.Response.NatGatewayId == nil {
				return resource.NonRetryableError(fmt.Errorf("Create vpc private nat gateway translation nat rule failed, Response is nil."))
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s create vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	d.SetId(natGatewayId)
	return resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleRead(d, meta)
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_nat_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		natGatewayId = d.Id()
	)

	respData, err := service.DescribeVpcPrivateNatGatewayTranslationNatRuleById(ctx, natGatewayId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_private_nat_gateway_translation_nat_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("nat_gateway_id", natGatewayId)

	// Bucket the API response in its authoritative order.
	legacyList := make([]map[string]interface{}, 0, len(respData))
	localNetworkList := make([]map[string]interface{}, 0)
	localTransportList := make([]map[string]interface{}, 0)
	peerNetworkList := make([]map[string]interface{}, 0)

	for _, item := range respData {
		legacyList = append(legacyList, flattenLegacyRule(item))

		if item.TranslationDirection == nil || item.TranslationType == nil {
			continue
		}
		direction := *item.TranslationDirection
		typ := *item.TranslationType

		switch {
		case direction == "LOCAL" && typ == "NETWORK_LAYER":
			localNetworkList = append(localNetworkList, flattenNetworkLayerRule(item))
		case direction == "LOCAL" && typ == "TRANSPORT_LAYER":
			localTransportList = append(localTransportList, flattenTransportLayerRule(item))
		case direction == "PEER" && typ == "NETWORK_LAYER":
			peerNetworkList = append(peerNetworkList, flattenNetworkLayerRule(item))
		}
	}

	_ = d.Set("translation_nat_rules", legacyList)
	_ = d.Set("local_network_layer_rules", localNetworkList)
	_ = d.Set("local_transport_layer_rules", localTransportList)
	_ = d.Set("peer_network_layer_rules", peerNetworkList)

	return nil
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_nat_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		natGatewayId = d.Id()
	)

	// Legacy field path (Set semantics: add/remove via Set.Difference).
	if d.HasChange("translation_nat_rules") {
		oldInterface, newInterface := d.GetChange("translation_nat_rules")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := olds.Difference(news).List()
		add := news.Difference(olds).List()

		if len(remove) > 0 {
			request := vpcv20170312.NewDeletePrivateNatGatewayTranslationNatRuleRequest()
			for _, item := range remove {
				request.TranslationNatRules = append(request.TranslationNatRules, buildLegacyRule(item.(map[string]interface{})))
			}
			request.NatGatewayId = &natGatewayId
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeletePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				}
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				return nil
			})
			if reqErr != nil {
				log.Printf("[CRITAL]%s delete vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}

		if len(add) > 0 {
			var allRules []*vpcv20170312.TranslationNatRuleInput
			for _, item := range add {
				allRules = append(allRules, buildLegacyRuleInput(item.(map[string]interface{})))
			}
			for i := 0; i < len(allRules); i += MAX_CREATE_RULES_LEN {
				end := i + MAX_CREATE_RULES_LEN
				if end > len(allRules) {
					end = len(allRules)
				}
				request := vpcv20170312.NewCreatePrivateNatGatewayTranslationNatRuleRequest()
				request.NatGatewayId = helper.String(natGatewayId)
				request.TranslationNatRules = allRules[i:end]
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreatePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					}
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					if result == nil || result.Response == nil || result.Response.NatGatewayId == nil {
						return resource.NonRetryableError(fmt.Errorf("Create vpc private nat gateway translation nat rule failed, Response is nil."))
					}
					return nil
				})
				if reqErr != nil {
					log.Printf("[CRITAL]%s create vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			}
		}
	}

	// Typed list fields — diffed by index per (direction, type) bucket.
	if err := updateTypedRuleField(ctx, logId, meta, d, natGatewayId,
		"local_network_layer_rules", "LOCAL", "NETWORK_LAYER", true); err != nil {
		return err
	}
	if err := updateTypedRuleField(ctx, logId, meta, d, natGatewayId,
		"local_transport_layer_rules", "LOCAL", "TRANSPORT_LAYER", false); err != nil {
		return err
	}
	if err := updateTypedRuleField(ctx, logId, meta, d, natGatewayId,
		"peer_network_layer_rules", "PEER", "NETWORK_LAYER", true); err != nil {
		return err
	}

	return resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleRead(d, meta)
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_nat_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request      = vpcv20170312.NewDeletePrivateNatGatewayTranslationNatRuleRequest()
		natGatewayId = d.Id()
	)

	respData, err := service.DescribeVpcPrivateNatGatewayTranslationNatRuleById(ctx, natGatewayId)
	if err != nil {
		return err
	}

	for _, item := range respData {
		translationNatRule := vpcv20170312.TranslationNatRule{}
		if item.TranslationDirection != nil {
			translationNatRule.TranslationDirection = item.TranslationDirection
		}
		if item.TranslationType != nil {
			translationNatRule.TranslationType = item.TranslationType
		}
		if item.TranslationIp != nil {
			translationNatRule.TranslationIp = item.TranslationIp
		}
		if item.Description != nil {
			translationNatRule.Description = item.Description
		}
		if item.OriginalIp != nil {
			translationNatRule.OriginalIp = item.OriginalIp
		}
		request.TranslationNatRules = append(request.TranslationNatRules, &translationNatRule)
	}

	request.NatGatewayId = &natGatewayId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeletePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// Helpers.
//
// Per API uniqueness contract: three-layer rules key on `original_ip`;
// four-layer rules key on `translation_ip`. The `includeOriginalIp` flag
// selects bucket shape across all helpers below.

// listItemsToMaps casts a TypeList raw value to []map[string]interface{}.
func listItemsToMaps(raw interface{}) []map[string]interface{} {
	if raw == nil {
		return nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

// equalRuleMap compares two rule maps across all bucket-relevant fields.
func equalRuleMap(a, b map[string]interface{}, includeOriginalIp bool) bool {
	keys := []string{"translation_ip", "description"}
	if includeOriginalIp {
		keys = append(keys, "original_ip")
	}
	for _, k := range keys {
		av, _ := a[k].(string)
		bv, _ := b[k].(string)
		if av != bv {
			return false
		}
	}
	return true
}

type modifyPair struct {
	oldItem map[string]interface{}
	newItem map[string]interface{}
}

// diffByIndex partitions an old vs new typed-list pair by index position:
//   - Equal-length slots with field differences → MODIFY (covers in-place key
//     IP rename, since identity is the slot index, not the IP value).
//   - Tail of new beyond min length → CREATE.
//   - Tail of old beyond min length → DELETE.
func diffByIndex(oldList, newList []map[string]interface{}, includeOriginalIp bool) (
	toCreate []map[string]interface{},
	toDelete []map[string]interface{},
	toModify []modifyPair,
) {
	minLen := len(oldList)
	if len(newList) < minLen {
		minLen = len(newList)
	}
	for i := 0; i < minLen; i++ {
		if !equalRuleMap(oldList[i], newList[i], includeOriginalIp) {
			toModify = append(toModify, modifyPair{oldItem: oldList[i], newItem: newList[i]})
		}
	}
	if len(newList) > minLen {
		toCreate = append(toCreate, newList[minLen:]...)
	}
	if len(oldList) > minLen {
		toDelete = append(toDelete, oldList[minLen:]...)
	}
	return
}

// buildLegacyRuleInput builds a Create-shape input from a legacy field map.
func buildLegacyRuleInput(m map[string]interface{}) *vpcv20170312.TranslationNatRuleInput {
	r := &vpcv20170312.TranslationNatRuleInput{}
	if v, ok := m["translation_direction"].(string); ok && v != "" {
		r.TranslationDirection = helper.String(v)
	}
	if v, ok := m["translation_type"].(string); ok && v != "" {
		r.TranslationType = helper.String(v)
	}
	if v, ok := m["translation_ip"].(string); ok && v != "" {
		r.TranslationIp = helper.String(v)
	}
	if v, ok := m["description"].(string); ok && v != "" {
		r.Description = helper.String(v)
	}
	if v, ok := m["original_ip"].(string); ok && v != "" {
		r.OriginalIp = helper.String(v)
	}
	return r
}

// buildLegacyRule builds a Delete-shape rule from a legacy field map.
func buildLegacyRule(m map[string]interface{}) *vpcv20170312.TranslationNatRule {
	r := &vpcv20170312.TranslationNatRule{}
	if v, ok := m["translation_direction"].(string); ok && v != "" {
		r.TranslationDirection = helper.String(v)
	}
	if v, ok := m["translation_type"].(string); ok && v != "" {
		r.TranslationType = helper.String(v)
	}
	if v, ok := m["translation_ip"].(string); ok && v != "" {
		r.TranslationIp = helper.String(v)
	}
	if v, ok := m["description"].(string); ok && v != "" {
		r.Description = helper.String(v)
	}
	if v, ok := m["original_ip"].(string); ok && v != "" {
		r.OriginalIp = helper.String(v)
	}
	return r
}

// buildTypedRuleInput tags a typed-list map with (direction, type) for Create.
func buildTypedRuleInput(m map[string]interface{}, direction, typ string, includeOriginalIp bool) *vpcv20170312.TranslationNatRuleInput {
	r := &vpcv20170312.TranslationNatRuleInput{
		TranslationDirection: helper.String(direction),
		TranslationType:      helper.String(typ),
	}
	if v, ok := m["translation_ip"].(string); ok && v != "" {
		r.TranslationIp = helper.String(v)
	}
	if v, ok := m["description"].(string); ok && v != "" {
		r.Description = helper.String(v)
	}
	if includeOriginalIp {
		if v, ok := m["original_ip"].(string); ok && v != "" {
			r.OriginalIp = helper.String(v)
		}
	}
	return r
}

// buildTypedRule mirrors buildTypedRuleInput for the Delete shape.
func buildTypedRule(m map[string]interface{}, direction, typ string, includeOriginalIp bool) *vpcv20170312.TranslationNatRule {
	r := &vpcv20170312.TranslationNatRule{
		TranslationDirection: helper.String(direction),
		TranslationType:      helper.String(typ),
	}
	if v, ok := m["translation_ip"].(string); ok && v != "" {
		r.TranslationIp = helper.String(v)
	}
	if v, ok := m["description"].(string); ok && v != "" {
		r.Description = helper.String(v)
	}
	if includeOriginalIp {
		if v, ok := m["original_ip"].(string); ok && v != "" {
			r.OriginalIp = helper.String(v)
		}
	}
	return r
}

// buildModifyDiff packs (oldItem, newItem) into a TranslationNatRuleDiff.
// Three-layer rules populate Old/New OriginalIp; four-layer rules omit them.
// `description` is sent unconditionally so a user-cleared description (empty
// string) actually reaches the API instead of being silently swallowed.
func buildModifyDiff(oldItem, newItem map[string]interface{}, direction, typ string, includeOriginalIp bool) *vpcv20170312.TranslationNatRuleDiff {
	d := &vpcv20170312.TranslationNatRuleDiff{
		TranslationDirection: helper.String(direction),
		TranslationType:      helper.String(typ),
	}
	if v, ok := oldItem["translation_ip"].(string); ok && v != "" {
		d.OldTranslationIp = helper.String(v)
	}
	if v, ok := newItem["translation_ip"].(string); ok && v != "" {
		d.TranslationIp = helper.String(v)
	}
	if v, ok := newItem["description"].(string); ok {
		d.Description = helper.String(v)
	}
	if includeOriginalIp {
		if v, ok := oldItem["original_ip"].(string); ok && v != "" {
			d.OldOriginalIp = helper.String(v)
		}
		if v, ok := newItem["original_ip"].(string); ok && v != "" {
			d.OriginalIp = helper.String(v)
		}
	}
	return d
}

// buildTypedRuleInputs concatenates Create inputs from the three typed-list
// fields in canonical order: LOCAL/NETWORK → LOCAL/TRANSPORT → PEER/NETWORK.
func buildTypedRuleInputs(d *schema.ResourceData) []*vpcv20170312.TranslationNatRuleInput {
	var out []*vpcv20170312.TranslationNatRuleInput
	for _, m := range listItemsToMaps(d.Get("local_network_layer_rules")) {
		out = append(out, buildTypedRuleInput(m, "LOCAL", "NETWORK_LAYER", true))
	}
	for _, m := range listItemsToMaps(d.Get("local_transport_layer_rules")) {
		out = append(out, buildTypedRuleInput(m, "LOCAL", "TRANSPORT_LAYER", false))
	}
	for _, m := range listItemsToMaps(d.Get("peer_network_layer_rules")) {
		out = append(out, buildTypedRuleInput(m, "PEER", "NETWORK_LAYER", true))
	}
	return out
}

// flattenLegacyRule produces a state map for the legacy `translation_nat_rules` Set.
func flattenLegacyRule(item *vpcv20170312.TranslationNatRule) map[string]interface{} {
	out := map[string]interface{}{}
	if item.TranslationDirection != nil {
		out["translation_direction"] = *item.TranslationDirection
	}
	if item.TranslationType != nil {
		out["translation_type"] = *item.TranslationType
	}
	if item.TranslationIp != nil {
		out["translation_ip"] = *item.TranslationIp
	}
	if item.Description != nil {
		out["description"] = *item.Description
	}
	if item.OriginalIp != nil {
		out["original_ip"] = *item.OriginalIp
	}
	return out
}

// flattenNetworkLayerRule produces a state map for three-layer typed lists.
func flattenNetworkLayerRule(item *vpcv20170312.TranslationNatRule) map[string]interface{} {
	out := map[string]interface{}{}
	if item.TranslationIp != nil {
		out["translation_ip"] = *item.TranslationIp
	}
	if item.OriginalIp != nil {
		out["original_ip"] = *item.OriginalIp
	}
	if item.Description != nil {
		out["description"] = *item.Description
	}
	return out
}

// flattenTransportLayerRule produces a state map for the four-layer typed list.
// `original_ip` is intentionally omitted: the API rejects it on four-layer rules.
func flattenTransportLayerRule(item *vpcv20170312.TranslationNatRule) map[string]interface{} {
	out := map[string]interface{}{}
	if item.TranslationIp != nil {
		out["translation_ip"] = *item.TranslationIp
	}
	if item.Description != nil {
		out["description"] = *item.Description
	}
	return out
}

// updateTypedRuleField runs the positional diff for one typed-list field and
// dispatches the work through applyRulesDiff (DELETE → MODIFY → CREATE).
func updateTypedRuleField(
	ctx context.Context,
	logId string,
	meta interface{},
	d *schema.ResourceData,
	natGatewayId string,
	field string,
	direction string,
	typ string,
	includeOriginalIp bool,
) error {
	if !d.HasChange(field) {
		return nil
	}

	oldRaw, newRaw := d.GetChange(field)
	oldMaps := listItemsToMaps(oldRaw)
	newMaps := listItemsToMaps(newRaw)
	createMaps, deleteMaps, modifyPairs := diffByIndex(oldMaps, newMaps, includeOriginalIp)

	var deleteRules []*vpcv20170312.TranslationNatRule
	for _, m := range deleteMaps {
		deleteRules = append(deleteRules, buildTypedRule(m, direction, typ, includeOriginalIp))
	}
	var createRules []*vpcv20170312.TranslationNatRuleInput
	for _, m := range createMaps {
		createRules = append(createRules, buildTypedRuleInput(m, direction, typ, includeOriginalIp))
	}
	var modifyDiffs []*vpcv20170312.TranslationNatRuleDiff
	for _, p := range modifyPairs {
		modifyDiffs = append(modifyDiffs, buildModifyDiff(p.oldItem, p.newItem, direction, typ, includeOriginalIp))
	}

	return applyRulesDiff(ctx, logId, meta, natGatewayId, deleteRules, modifyDiffs, createRules)
}

// applyRulesDiff issues DELETE → MODIFY → CREATE for one typed-list field.
// MODIFY is one request per pair (SDK constraint: single rule per modify).
// CREATE is batched by MAX_CREATE_RULES_LEN. Each call is retry-wrapped.
func applyRulesDiff(
	ctx context.Context,
	logId string,
	meta interface{},
	natGatewayId string,
	deleteRules []*vpcv20170312.TranslationNatRule,
	modifyDiffs []*vpcv20170312.TranslationNatRuleDiff,
	createRules []*vpcv20170312.TranslationNatRuleInput,
) error {
	// 1. DELETE.
	if len(deleteRules) > 0 {
		request := vpcv20170312.NewDeletePrivateNatGatewayTranslationNatRuleRequest()
		request.NatGatewayId = &natGatewayId
		request.TranslationNatRules = deleteRules

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeletePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})
		if reqErr != nil {
			log.Printf("[CRITAL]%s delete vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	// 2. MODIFY.
	for _, diff := range modifyDiffs {
		request := vpcv20170312.NewModifyPrivateNatGatewayTranslationNatRuleRequest()
		request.NatGatewayId = &natGatewayId
		request.TranslationNatRules = []*vpcv20170312.TranslationNatRuleDiff{diff}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyPrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})
		if reqErr != nil {
			log.Printf("[CRITAL]%s modify vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	// 3. CREATE.
	for i := 0; i < len(createRules); i += MAX_CREATE_RULES_LEN {
		end := i + MAX_CREATE_RULES_LEN
		if end > len(createRules) {
			end = len(createRules)
		}

		request := vpcv20170312.NewCreatePrivateNatGatewayTranslationNatRuleRequest()
		request.NatGatewayId = helper.String(natGatewayId)
		request.TranslationNatRules = createRules[i:end]

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreatePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil || result.Response.NatGatewayId == nil {
				return resource.NonRetryableError(fmt.Errorf("Create vpc private nat gateway translation nat rule failed, Response is nil."))
			}
			return nil
		})
		if reqErr != nil {
			log.Printf("[CRITAL]%s create vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return nil
}

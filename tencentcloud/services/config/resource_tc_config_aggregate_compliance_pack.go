package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions (repo convention: only assert the
// extension interfaces; the base interface is checked implicitly by the
// factory return value).
var (
	_ resource.ResourceWithConfigure   = &ConfigAggregateCompliancePackResource{}
	_ resource.ResourceWithImportState = &ConfigAggregateCompliancePackResource{}
)

// NewConfigAggregateCompliancePackResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewConfigAggregateCompliancePackResource() resource.Resource {
	return &ConfigAggregateCompliancePackResource{}
}

// ConfigAggregateCompliancePackResource implements resource.Resource for
// tencentcloud_config_aggregate_compliance_pack.
//
// The cloud side locates a resource through the combination of an aggregator
// account group id and the compliance pack id, so the Terraform id is the two
// joined with `#` (tccommon.FILED_SP) and split for every CRUD call.
type ConfigAggregateCompliancePackResource struct {
	fw.ResourceWithConfigure
}

// ConfigAggregateCompliancePackResourceModel maps the schema attributes. Every
// schema attribute appears here exactly once.
type ConfigAggregateCompliancePackResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	CompliancePackName types.String `tfsdk:"compliance_pack_name"`
	AccountGroupId     types.String `tfsdk:"account_group_id"`
	RiskLevel          types.Int64  `tfsdk:"risk_level"`
	Description        types.String `tfsdk:"description"`
	ConfigRules        types.List   `tfsdk:"config_rules"`
	Status             types.String `tfsdk:"status"`
	CompliancePackId  types.String `tfsdk:"compliance_pack_id"`
	CreateTime         types.String `tfsdk:"create_time"`
}

// aggregateCompliancePackRuleModel maps the config_rules nested object.
type aggregateCompliancePackRuleModel struct {
	RuleName              types.String `tfsdk:"rule_name"`
	RiskLevel             types.Int64  `tfsdk:"risk_level"`
	Identifier            types.String `tfsdk:"identifier"`
	Description           types.String `tfsdk:"description"`
	ManagedRuleIdentifier types.String `tfsdk:"managed_rule_identifier"`
	ConfigRuleId          types.String `tfsdk:"config_rule_id"`
	CompliancePackId      types.String `tfsdk:"compliance_pack_id"`
	InputParameter        types.List   `tfsdk:"input_parameter"`
}

// aggregateCompliancePackInputParameterModel maps the input_parameter nested object.
type aggregateCompliancePackInputParameterModel struct {
	ParameterKey types.String `tfsdk:"parameter_key"`
	Type          types.String `tfsdk:"type"`
	Value         types.String `tfsdk:"value"`
}

func (r *ConfigAggregateCompliancePackResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_config_aggregate_compliance_pack"
}

func (r *ConfigAggregateCompliancePackResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create an aggregate (account group) compliance pack of Config.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<account_group_id>#<compliance_pack_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"compliance_pack_name": schema.StringAttribute{
				Required:    true,
				Description: "Compliance pack name.",
			},
			"account_group_id": schema.StringAttribute{
				Required:    true,
				Description:  "Account group id. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"risk_level": schema.Int64Attribute{
				Required:    true,
				Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the compliance pack.",
			},
			"config_rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of compliance pack rules.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_name": schema.StringAttribute{
							Required:    true,
							Description: "Rule name.",
						},
						"risk_level": schema.Int64Attribute{
							Required:    true,
							Description: "Rule risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
						},
						"input_parameter": schema.ListNestedAttribute{
							Required:    true,
							Description: "Rule input parameters.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"parameter_key": schema.StringAttribute{
										Required:    true,
										Description: "Parameter key.",
									},
									"type": schema.StringAttribute{
										Optional:    true,
										Description: "Parameter type. Required: `Require`, optional: `Optional`.",
									},
									"value": schema.StringAttribute{
										Optional:    true,
										Description: "Parameter value.",
									},
								},
							},
						},
						"identifier": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Rule identifier (managed rule name or custom rule cloud function ARN).",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Rule description.",
						},
						"managed_rule_identifier": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Managed rule identifier (preset rule identity).",
						},
						"config_rule_id": schema.StringAttribute{
							Optional:    true,
							Computed: true,
							Description: "Config rule id.",
						},
						"compliance_pack_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Compliance pack id that this rule belongs to.",
						},
					},
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Compliance pack status.",
			},
			"compliance_pack_id": schema.StringAttribute{
				Computed:    true,
				Description: "Compliance pack id.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the compliance pack.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// splitAggregateCompliancePackId parses `<account_group_id>#<compliance_pack_id>`.
func splitAggregateCompliancePackId(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// Create handles `terraform apply` for new resources.
func (r *ConfigAggregateCompliancePackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConfigAggregateCompliancePackResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId := plan.AccountGroupId.ValueString()

	request := configv20220802.NewAddAggregateCompliancePackRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.CompliancePackName = plan.CompliancePackName.ValueStringPointer()
	request.RiskLevel = helper.Uint64(uint64(plan.RiskLevel.ValueInt64()))
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}

	rules, diags := expandAggregateCompliancePackRules(ctx, plan.ConfigRules)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ConfigRules = rules

	service := NewConfigService(r.Client())
	compliancePackId, err := service.AddAggregateCompliancePack(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Config aggregate compliance pack",
			fmt.Sprintf("AddAggregateCompliancePack failed: %v", err),
		)
		return
	}

	plan.CompliancePackId = types.StringValue(compliancePackId)
	plan.Id = types.StringValue(strings.Join([]string{accountGroupId, compliancePackId}, tccommon.FILED_SP))

	// Persist the identity before any read-back so a partial failure still
	// leaves a manageable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := service.DescribeAggregateCompliancePackById(ctx, accountGroupId, compliancePackId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Config aggregate compliance pack after create", err.Error())
		return
	}
	if item == nil || item.CompliancePackId == nil {
		resp.Diagnostics.AddError(
			"Config aggregate compliance pack not readable after create",
			fmt.Sprintf("resource %s was created but cannot be read back", plan.Id.ValueString()),
		)
		return
	}
	refreshAggregateCompliancePackModel(ctx, &plan, accountGroupId, compliancePackId, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *ConfigAggregateCompliancePackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConfigAggregateCompliancePackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, compliancePackId, ok := splitAggregateCompliancePackId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("aggregate compliance pack id is broken, id is %q (expected `<account_group_id>%s<compliance_pack_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	service := NewConfigService(r.Client())
	item, err := service.DescribeAggregateCompliancePackById(ctx, accountGroupId, compliancePackId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Config aggregate compliance pack", err.Error())
		return
	}

	// Confirmed "not found" on the cloud side => drop it from state. A nil
	// response is an API failure and must not land here.
	if item == nil || item.CompliancePackId == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("Config aggregate compliance pack %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	refreshAggregateCompliancePackModel(ctx, &state, accountGroupId, compliancePackId, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *ConfigAggregateCompliancePackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConfigAggregateCompliancePackResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ConfigAggregateCompliancePackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, compliancePackId, ok := splitAggregateCompliancePackId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("aggregate compliance pack id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	request := configv20220802.NewUpdateAggregateCompliancePackRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.CompliancePackId = helper.String(compliancePackId)
	request.CompliancePackName = plan.CompliancePackName.ValueStringPointer()
	request.RiskLevel = helper.Uint64(uint64(plan.RiskLevel.ValueInt64()))
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}

	rules, diags := expandAggregateCompliancePackRules(ctx, plan.ConfigRules)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ConfigRules = rules

	service := NewConfigService(r.Client())
	if err := service.UpdateAggregateCompliancePack(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating Config aggregate compliance pack",
			fmt.Sprintf("UpdateAggregateCompliancePack failed: %v", err),
		)
		return
	}

	item, err := service.DescribeAggregateCompliancePackById(ctx, accountGroupId, compliancePackId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Config aggregate compliance pack after update", err.Error())
		return
	}
	if item == nil || item.CompliancePackId == nil {
		resp.Diagnostics.AddError(
			"Config aggregate compliance pack not readable after update",
			fmt.Sprintf("resource %s was updated but cannot be read back", state.Id.ValueString()),
		)
		return
	}
	plan.Id = state.Id
	refreshAggregateCompliancePackModel(ctx, &plan, accountGroupId, compliancePackId, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *ConfigAggregateCompliancePackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConfigAggregateCompliancePackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, compliancePackId, ok := splitAggregateCompliancePackId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("aggregate compliance pack id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	request := configv20220802.NewDeleteAggregateCompliancePackRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.CompliancePackId = helper.String(compliancePackId)

	service := NewConfigService(r.Client())
	if err := service.DeleteAggregateCompliancePack(ctx, request); err != nil {
		if tccommon.IsExpectError(err, []string{"ResourceNotFound", "ResourceInexistent"}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting Config aggregate compliance pack",
			fmt.Sprintf("DeleteAggregateCompliancePack failed: %v", err),
		)
		return
	}
}

// ImportState supports `terraform import` using the
// `<account_group_id>#<compliance_pack_id>` form.
func (r *ConfigAggregateCompliancePackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := splitAggregateCompliancePackId(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<account_group_id>%s<compliance_pack_id>`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// expandAggregateCompliancePackRules converts the plan list of rules into the
// SDK CompliancePackRule slice.
func expandAggregateCompliancePackRules(ctx context.Context, list types.List) ([]*configv20220802.CompliancePackRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var rules []aggregateCompliancePackRuleModel
	diags.Append(list.ElementsAs(ctx, &rules, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*configv20220802.CompliancePackRule, 0, len(rules))
	for _, rule := range rules {
		r := &configv20220802.CompliancePackRule{
			RuleName:  rule.RuleName.ValueStringPointer(),
			RiskLevel: helper.Uint64(uint64(rule.RiskLevel.ValueInt64())),
		}
		if !rule.Identifier.IsNull() && !rule.Identifier.IsUnknown() {
			r.Identifier = rule.Identifier.ValueStringPointer()
		}
		if !rule.Description.IsNull() && !rule.Description.IsUnknown() {
			r.Description = rule.Description.ValueStringPointer()
		}
		if !rule.ManagedRuleIdentifier.IsNull() && !rule.ManagedRuleIdentifier.IsUnknown() {
			r.ManagedRuleIdentifier = rule.ManagedRuleIdentifier.ValueStringPointer()
		}
		if !rule.ConfigRuleId.IsNull() && !rule.ConfigRuleId.IsUnknown() {
			r.ConfigRuleId = rule.ConfigRuleId.ValueStringPointer()
		}
		if !rule.CompliancePackId.IsNull() && !rule.CompliancePackId.IsUnknown() {
			r.CompliancePackId = rule.CompliancePackId.ValueStringPointer()
		}
		r.InputParameter = expandAggregateInputParameters(rule.InputParameter)
		out = append(out, r)
	}

	return out, diags
}

// expandAggregateInputParameters converts the plan parameter list into the SDK
// InputParameter slice.
func expandAggregateInputParameters(list types.List) []*configv20220802.InputParameter {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var params []aggregateCompliancePackInputParameterModel
	// Errors here are ignored: the list is structurally validated at plan time.
	_ = list.ElementsAs(context.Background(), &params, false)

	out := make([]*configv20220802.InputParameter, 0, len(params))
	for _, p := range params {
		param := &configv20220802.InputParameter{
			ParameterKey: p.ParameterKey.ValueStringPointer(),
		}
		if !p.Type.IsNull() && !p.Type.IsUnknown() {
			param.Type = p.Type.ValueStringPointer()
		}
		if !p.Value.IsNull() && !p.Value.IsUnknown() {
			param.Value = p.Value.ValueStringPointer()
		}
		out = append(out, param)
	}
	return out
}

// refreshAggregateCompliancePackModel copies cloud-owned fields into the
// model. For fields that the Describe response does not echo back
// (managed_rule_identifier and per-rule compliance_pack_id), the previous state
// value is retained instead of being cleared, so drift detection for those
// fields is limited (documented in the resource .md).
func refreshAggregateCompliancePackModel(ctx context.Context, m *ConfigAggregateCompliancePackResourceModel, accountGroupId, compliancePackId string, item *configv20220802.DescribeAggregateCompliancePackResponseParams) {
	m.AccountGroupId = types.StringValue(accountGroupId)
	m.CompliancePackId = types.StringValue(compliancePackId)
	m.Id = types.StringValue(strings.Join([]string{accountGroupId, compliancePackId}, tccommon.FILED_SP))

	m.CompliancePackName = types.StringPointerValue(item.CompliancePackName)
	m.Description = types.StringPointerValue(item.Description)
	m.CreateTime = types.StringPointerValue(item.CreateTime)
	m.Status = types.StringPointerValue(item.Status)
	if item.RiskLevel != nil {
		m.RiskLevel = types.Int64Value(int64(*item.RiskLevel))
	}

	m.ConfigRules = flattenAggregateComplianceConfigRules(ctx, item.ConfigRules, m.ConfigRules)
}

// configRulesAttrTypes is the object type of the config_rules nested object,
// used to build types.List values.
func configRulesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"rule_name":               types.StringType,
		"risk_level":              types.Int64Type,
		"identifier":              types.StringType,
		"description":             types.StringType,
		"managed_rule_identifier": types.StringType,
		"config_rule_id":          types.StringType,
		"compliance_pack_id":      types.StringType,
		"input_parameter":        types.ListType{ElemType: types.ObjectType{AttrTypes: inputParameterAttrTypes()}},
	}
}

func inputParameterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"parameter_key": types.StringType,
		"type":           types.StringType,
		"value":          types.StringType,
	}
}

// flattenAggregateComplianceConfigRules builds the config_rules list from the
// Describe response. The previous plan/state list is passed in so that fields
// the cloud does not echo back (managed_rule_identifier, compliance_pack_id)
// keep their prior value instead of being nulled. Rules are correlated by
// identifier (when available) falling back to config_rule_id then rule_name.
func flattenAggregateComplianceConfigRules(ctx context.Context, rules []*configv20220802.ComplianceConfigRule, prevList types.List) types.List {
	if len(rules) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: configRulesAttrTypes()})
	}

	var prev []aggregateCompliancePackRuleModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		_ = prevList.ElementsAs(ctx, &prev, false)
	}

	out := make([]aggregateCompliancePackRuleModel, 0, len(rules))
	for _, rule := range rules {
		rm := aggregateCompliancePackRuleModel{
			RuleName: types.StringPointerValue(rule.RuleName),
		}
		if rule.RiskLevel != nil {
			rm.RiskLevel = types.Int64Value(int64(*rule.RiskLevel))
		}
		rm.Identifier = types.StringPointerValue(rule.Identifier)
		rm.Description = types.StringPointerValue(rule.Description)
		rm.ConfigRuleId = types.StringPointerValue(rule.ConfigRuleId)

		// managed_rule_identifier: map it from the cloud Identifier if the prior
		// state had it set (it is the documented alias). Otherwise: if prior
		// state carried managed_rule_identifier and identifier both, then for the rule,
		// we should restore managed_rule_identifier from prior state matching.
		rm.ManagedRuleIdentifier = types.StringNull()

		// compliance_pack_id: dropped by cloud at rule level. Restore from prior state matched by identifier/rule_id/name.
		rm.CompliancePackId = types.StringNull()

		// input_parameter
		rm.InputParameter = flattenAggregateInputParameters(ctx, rule.InputParameter)

		// Try to carry over write-only fields from the previous rule state.
		key := ""
		if rule.Identifier != nil && *rule.Identifier != "" {
			key = *rule.Identifier
		} else if rule.ConfigRuleId != nil && *rule.ConfigRuleId != "" {
			key = *rule.ConfigRuleId
		} else if rule.RuleName != nil {
			key = *rule.RuleName
		}
		for _, p := range prev {
			pKey := ""
			if !p.Identifier.IsNull() && !p.Identifier.IsUnknown() {
				pKey = p.Identifier.ValueString()
			} else if !p.ConfigRuleId.IsNull() && !p.ConfigRuleId.IsUnknown() {
				pKey = p.ConfigRuleId.ValueString()
			} else if !p.RuleName.IsNull() && !p.RuleName.IsUnknown() {
				pKey = p.RuleName.ValueString()
			}
			if pKey != "" && pKey == key {
				rm.ManagedRuleIdentifier = p.ManagedRuleIdentifier
				rm.CompliancePackId = p.CompliancePackId
				break
			}
		}

		out = append(out, rm)
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: configRulesAttrTypes()}, out)
	// On diagnostic failure prefer a null list over writing a malformed value.
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: configRulesAttrTypes()})
	}
	return list
}

// flattenAggregateInputParameters builds the input_parameter list from the
// Describe response.
func flattenAggregateInputParameters(ctx context.Context, params []*configv20220802.InputParameter) types.List {
	if len(params) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: inputParameterAttrTypes()})
	}

	out := make([]aggregateCompliancePackInputParameterModel, 0, len(params))
	for _, p := range params {
		pm := aggregateCompliancePackInputParameterModel{
			ParameterKey: types.StringPointerValue(p.ParameterKey),
			Type:           types.StringPointerValue(p.Type),
			Value:          types.StringPointerValue(p.Value),
		}
		out = append(out, pm)
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: inputParameterAttrTypes()}, out)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: inputParameterAttrTypes()})
	}
	return list
}
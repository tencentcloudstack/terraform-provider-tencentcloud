package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	sdkerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions (repo convention: only assert the extension
// interfaces; the base interface is checked implicitly by the factory return value).
var (
	_ resource.ResourceWithConfigure   = &ConfigAggregateConfigRuleResource{}
	_ resource.ResourceWithImportState = &ConfigAggregateConfigRuleResource{}
)

// NewConfigAggregateConfigRuleResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewConfigAggregateConfigRuleResource() resource.Resource {
	return &ConfigAggregateConfigRuleResource{}
}

// ConfigAggregateConfigRuleResource implements resource.Resource for
// tencentcloud_config_aggregate_config_rule.
type ConfigAggregateConfigRuleResource struct {
	fw.ResourceWithConfigure
}

// ConfigAggregateConfigRuleResourceModel maps the schema attributes.
type ConfigAggregateConfigRuleResourceModel struct {
	Identifier               types.String `tfsdk:"identifier"`
	IdentifierType           types.String `tfsdk:"identifier_type"`
	RuleName                types.String `tfsdk:"rule_name"`
	ResourceType            types.Set    `tfsdk:"resource_type"`
	TriggerType             types.List    `tfsdk:"trigger_type"`
	RiskLevel               types.Int64   `tfsdk:"risk_level"`
	AccountGroupId           types.String `tfsdk:"account_group_id"`
	InputParameter           types.List    `tfsdk:"input_parameter"`
	Description              types.String `tfsdk:"description"`
	RegionScope            types.Set    `tfsdk:"region_scope"`
	TagsScope               types.List    `tfsdk:"tags_scope"`
	ExcludeResourceIdsScope types.Set    `tfsdk:"exclude_resource_ids_scope"`
	Id                       types.String `tfsdk:"id"`
	RuleId                  types.String `tfsdk:"rule_id"`
	ConfigRule              types.Object `tfsdk:"config_rule"`
}

// TriggerTypeModel maps the trigger_type nested object.
type TriggerTypeModel struct {
	MessageType               types.String `tfsdk:"message_type"`
	MaximumExecutionFrequency types.String `tfsdk:"maximum_execution_frequency"`
}

// InputParameterModel maps the input_parameter nested object.
type InputParameterModel struct {
	ParameterKey types.String `tfsdk:"parameter_key"`
	Type         types.String `tfsdk:"type"`
	Value        types.String `tfsdk:"value"`
}

// TagsScopeModel maps the tags_scope nested object.
type TagsScopeModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// ConfigRuleModel maps the computed config_rule output object.
type ConfigRuleModel struct {
	Identifier              types.String `tfsdk:"identifier"`
	RuleName                types.String `tfsdk:"rule_name"`
	InputParameter          types.List   `tfsdk:"input_parameter"`
	SourceCondition         types.List   `tfsdk:"source_condition"`
	ResourceType            types.Set    `tfsdk:"resource_type"`
	Labels                  types.Set    `tfsdk:"labels"`
	RiskLevel               types.Int64  `tfsdk:"risk_level"`
	ServiceFunction         types.String `tfsdk:"service_function"`
	CreateTime              types.String `tfsdk:"create_time"`
	Description            types.String `tfsdk:"description"`
	Status                  types.String `tfsdk:"status"`
	ComplianceResult        types.String `tfsdk:"compliance_result"`
	Annotation              types.List   `tfsdk:"annotation"`
	ConfigRuleInvokedTime   types.String `tfsdk:"config_rule_invoked_time"`
	ConfigRuleId            types.String `tfsdk:"config_rule_id"`
	IdentifierType          types.String `tfsdk:"identifier_type"`
	CompliancePackId        types.String `tfsdk:"compliance_pack_id"`
	TriggerType            types.List   `tfsdk:"trigger_type"`
	ManageInputParameter    types.List   `tfsdk:"manage_input_parameter"`
	CompliancePackName      types.String `tfsdk:"compliance_pack_name"`
	RegionsScope           types.Set    `tfsdk:"regions_scope"`
	TagsScope               types.List   `tfsdk:"tags_scope"`
	ExcludeResourceIdsScope types.Set   `tfsdk:"exclude_resource_ids_scope"`
	AccountGroupId          types.String `tfsdk:"account_group_id"`
	AccountGroupName        types.String `tfsdk:"account_group_name"`
	RuleOwnerId             types.Int64  `tfsdk:"rule_owner_id"`
	ManageTriggerType       types.Set    `tfsdk:"manage_trigger_type"`
}

// SourceConditionModel maps the source_condition nested object (read-only output).
type SourceConditionModel struct {
	EmptyAs      types.String `tfsdk:"empty_as"`
	SelectPath   types.String `tfsdk:"select_path"`
	Operator     types.String `tfsdk:"operator"`
	Required     types.Bool   `tfsdk:"required"`
	DesiredValue types.String `tfsdk:"desired_value"`
}

// AnnotationModel maps the annotation nested object (read-only output).
type AnnotationModel struct {
	Configuration types.String `tfsdk:"configuration"`
	DesiredValue  types.String `tfsdk:"desired_value"`
	Operator      types.String `tfsdk:"operator"`
	Property      types.String `tfsdk:"property"`
}

// ManageInputParameterModel maps the manage_input_parameter nested object (read-only output).
type ManageInputParameterModel struct {
	ValueType    types.String `tfsdk:"value_type"`
	ParameterKey types.String `tfsdk:"parameter_key"`
	Type         types.String `tfsdk:"type"`
	DefaultValue types.String `tfsdk:"default_value"`
	Description  types.String `tfsdk:"description"`
}

func configRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"identifier":                  types.StringType,
		"rule_name":                   types.StringType,
		"input_parameter":             types.ListType{ElemType: types.ObjectType{AttrTypes: inputParameterAttrTypes()}},
		"source_condition":            types.ListType{ElemType: types.ObjectType{AttrTypes: sourceConditionAttrTypes()}},
		"resource_type":               types.SetType{ElemType: types.StringType},
		"labels":                      types.SetType{ElemType: types.StringType},
		"risk_level":                  types.Int64Type,
		"service_function":            types.StringType,
		"create_time":                 types.StringType,
		"description":                 types.StringType,
		"status":                      types.StringType,
		"compliance_result":           types.StringType,
		"annotation":                  types.ListType{ElemType: types.ObjectType{AttrTypes: annotationAttrTypes()}},
		"config_rule_invoked_time":    types.StringType,
		"config_rule_id":              types.StringType,
		"identifier_type":             types.StringType,
		"compliance_pack_id":          types.StringType,
		"trigger_type":                types.ListType{ElemType: types.ObjectType{AttrTypes: triggerTypeAttrTypes()}},
		"manage_input_parameter":      types.ListType{ElemType: types.ObjectType{AttrTypes: manageInputParameterAttrTypes()}},
		"compliance_pack_name":        types.StringType,
		"regions_scope":               types.SetType{ElemType: types.StringType},
		"tags_scope":                  types.ListType{ElemType: types.ObjectType{AttrTypes: tagsScopeAttrTypes()}},
		"exclude_resource_ids_scope":  types.SetType{ElemType: types.StringType},
		"account_group_id":            types.StringType,
		"account_group_name":          types.StringType,
		"rule_owner_id":               types.Int64Type,
		"manage_trigger_type":         types.SetType{ElemType: types.StringType},
	}
}

func triggerTypeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"message_type":                types.StringType,
		"maximum_execution_frequency": types.StringType,
	}
}

func inputParameterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"parameter_key": types.StringType,
		"type":          types.StringType,
		"value":         types.StringType,
	}
}

func tagsScopeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}
}

func sourceConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"empty_as":      types.StringType,
		"select_path":   types.StringType,
		"operator":      types.StringType,
		"required":      types.BoolType,
		"desired_value": types.StringType,
	}
}

func annotationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"configuration": types.StringType,
		"desired_value": types.StringType,
		"operator":      types.StringType,
		"property":      types.StringType,
	}
}

func manageInputParameterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value_type":    types.StringType,
		"parameter_key": types.StringType,
		"type":          types.StringType,
		"default_value": types.StringType,
		"description":   types.StringType,
	}
}

func (r *ConfigAggregateConfigRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_config_aggregate_config_rule"
}

func (r *ConfigAggregateConfigRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a Config aggregate config rule under an account group.",
		Attributes: map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				Required:    true,
				Description: "Rule template identifier. For system preset rules use the identifier name; for custom rules use the cloud function ARN (region:functionName).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"identifier_type": schema.StringAttribute{
				Required:    true,
				Description: "Rule template type. Valid values: `SYSTEM` (system preset), `CUSTOMIZE` (custom).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_name": schema.StringAttribute{
				Required:    true,
				Description: "Rule name.",
			},
			"resource_type": schema.SetAttribute{
				Required:    true,
				Description: "List of supported resource types (e.g. QCS::CAM::User). The rule only applies to resources of the specified types.",
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"trigger_type": schema.ListNestedAttribute{
				Required:    true,
				Description: "Trigger type list, up to 2 entries.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"message_type": schema.StringAttribute{
							Required:    true,
							Description: "Trigger message type. Valid values: `ScheduledNotification` (periodic), `ConfigurationItemChangeNotification` (change).",
						},
						"maximum_execution_frequency": schema.StringAttribute{
							Optional:    true,
							Description: "Maximum execution frequency, only for `ScheduledNotification`. e.g. `TwentyFour_Hours`.",
						},
					},
				},
			},
			"risk_level": schema.Int64Attribute{
				Required:    true,
				Description: "Risk level. Valid values: `1` (high risk), `2` (medium risk), `3` (low risk).",
			},
			"account_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Account group ID. The aggregate rule is created under this account group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"input_parameter": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Rule input parameter list.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"parameter_key": schema.StringAttribute{
							Required:    true,
							Description: "Parameter key.",
						},
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "Parameter type: `Require` or `Optional`.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "Parameter value.",
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Rule description.",
			},
			"region_scope": schema.SetAttribute{
				Optional:    true,
				Description: "Region scope list; the rule only applies to resources in the specified regions.",
				ElementType: types.StringType,
			},
			"tags_scope": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Tag scope list; the rule only applies to resources with the specified tags.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag_key": schema.StringAttribute{
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": schema.StringAttribute{
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"exclude_resource_ids_scope": schema.SetAttribute{
				Optional:    true,
				Description: "Resource ID list excluded from rule evaluation.",
				ElementType: types.StringType,
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<account_group_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rule_id": schema.StringAttribute{
				Computed:    true,
				Description: "Config rule ID. Returned by the cloud after creation and refreshed by Read.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"config_rule": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Rule detail returned by the cloud.",
				Attributes:  configRuleAttributes(),
			},
		},
	}
}

// configRuleAttributes returns the schema attributes for the computed config_rule block.
func configRuleAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"identifier":                schema.StringAttribute{Computed: true, Description: "Rule identifier."},
		"rule_name":                 schema.StringAttribute{Computed: true, Description: "Rule name."},
		"input_parameter":           schema.ListNestedAttribute{Computed: true, Description: "Rule parameters.", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{"parameter_key": schema.StringAttribute{Computed: true, Description: "Parameter key."}, "type": schema.StringAttribute{Computed: true, Description: "Parameter type."}, "value": schema.StringAttribute{Computed: true, Description: "Parameter value."}}}},
		"source_condition":          schema.ListNestedAttribute{Computed: true, Description: "Rule trigger conditions.", NestedObject: schema.NestedAttributeObject{Attributes: sourceConditionSchemaAttrs()}},
		"resource_type":             schema.SetAttribute{Computed: true, ElementType: types.StringType, Description: "Supported resource types."},
		"labels":                    schema.SetAttribute{Computed: true, ElementType: types.StringType, Description: "Rule labels."},
		"risk_level":                schema.Int64Attribute{Computed: true, Description: "Risk level. 1:low, 2:medium, 3:high."},
		"service_function":          schema.StringAttribute{Computed: true, Description: "Corresponding cloud function."},
		"create_time":               schema.StringAttribute{Computed: true, Description: "Creation time."},
		"description":               schema.StringAttribute{Computed: true, Description: "Rule description."},
		"status":                    schema.StringAttribute{Computed: true, Description: "Rule status. ACTIVE: enabled, NO_ACTIVE: stopped."},
		"compliance_result":         schema.StringAttribute{Computed: true, Description: "Compliance result."},
		"annotation":                schema.ListNestedAttribute{Computed: true, Description: "Non-compliance annotation.", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{"configuration": schema.StringAttribute{Computed: true, Description: "Current resource configuration."}, "desired_value": schema.StringAttribute{Computed: true, Description: "Expected resource configuration."}, "operator": schema.StringAttribute{Computed: true, Description: "Comparison operator."}, "property": schema.StringAttribute{Computed: true, Description: "JSON path of the configuration property."}}}},
		"config_rule_invoked_time":  schema.StringAttribute{Computed: true, Description: "Last rule evaluation time."},
		"config_rule_id":            schema.StringAttribute{Computed: true, Description: "Rule ID."},
		"identifier_type":           schema.StringAttribute{Computed: true, Description: "CUSTOMIZE: custom rule, SYSTEM: managed rule."},
		"compliance_pack_id":        schema.StringAttribute{Computed: true, Description: "Compliance pack ID."},
		"trigger_type":              schema.ListNestedAttribute{Computed: true, Description: "Trigger types.", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{"message_type": schema.StringAttribute{Computed: true, Description: "Trigger type."}, "maximum_execution_frequency": schema.StringAttribute{Computed: true, Description: "Trigger time period."}}}},
		"manage_input_parameter":    schema.ListNestedAttribute{Computed: true, Description: "Managed parameter detail.", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{"value_type": schema.StringAttribute{Computed: true, Description: "Value type."}, "parameter_key": schema.StringAttribute{Computed: true, Description: "Parameter key."}, "type": schema.StringAttribute{Computed: true, Description: "Parameter type."}, "default_value": schema.StringAttribute{Computed: true, Description: "Default value."}, "description": schema.StringAttribute{Computed: true, Description: "Description."}}}},
		"compliance_pack_name":      schema.StringAttribute{Computed: true, Description: "Compliance pack name."},
		"regions_scope":             schema.SetAttribute{Computed: true, ElementType: types.StringType, Description: "Region scope."},
		"tags_scope":                schema.ListNestedAttribute{Computed: true, Description: "Tag scope.", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{"tag_key": schema.StringAttribute{Computed: true, Description: "Tag key."}, "tag_value": schema.StringAttribute{Computed: true, Description: "Tag value."}}}},
		"exclude_resource_ids_scope": schema.SetAttribute{Computed: true, ElementType: types.StringType, Description: "Excluded resource IDs."},
		"account_group_id":          schema.StringAttribute{Computed: true, Description: "Account group ID."},
		"account_group_name":        schema.StringAttribute{Computed: true, Description: "Account group name."},
		"rule_owner_id":             schema.Int64Attribute{Computed: true, Description: "Rule owner user ID."},
		"manage_trigger_type":       schema.SetAttribute{Computed: true, ElementType: types.StringType, Description: "Trigger types supported by preset rules."},
	}
}

// sourceConditionSchemaAttrs returns the schema attributes for a source_condition element.
func sourceConditionSchemaAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"empty_as":      schema.StringAttribute{Computed: true, Description: "Condition when empty."},
		"select_path":   schema.StringAttribute{Computed: true, Description: "Configuration path."},
		"operator":      schema.StringAttribute{Computed: true, Description: "Operator."},
		"required":      schema.BoolAttribute{Computed: true, Description: "Whether required."},
		"desired_value": schema.StringAttribute{Computed: true, Description: "Desired value."},
	}
}

// Create handles `terraform apply` for new resources.
func (r *ConfigAggregateConfigRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := configv20220802.NewAddAggregateConfigRuleRequest()
	request.Identifier = plan.Identifier.ValueStringPointer()
	request.IdentifierType = plan.IdentifierType.ValueStringPointer()
	request.RuleName = plan.RuleName.ValueStringPointer()
	request.ResourceType = stringSetValueToSlice(plan.ResourceType)
	request.TriggerType = triggerTypeModelsToSDK(plan.TriggerType)
	request.RiskLevel = helper.Int64Uint64(plan.RiskLevel.ValueInt64())
	request.AccountGroupId = plan.AccountGroupId.ValueStringPointer()

	if !plan.InputParameter.IsNull() && !plan.InputParameter.IsUnknown() {
		request.InputParameter = inputParameterModelsToSDK(plan.InputParameter)
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}
	if !plan.RegionScope.IsNull() && !plan.RegionScope.IsUnknown() {
		request.RegionScope = stringSetValueToSlice(plan.RegionScope)
	}
	if !plan.TagsScope.IsNull() && !plan.TagsScope.IsUnknown() {
		request.TagsScope = tagsScopeModelsToSDK(plan.TagsScope)
	}
	if !plan.ExcludeResourceIdsScope.IsNull() && !plan.ExcludeResourceIdsScope.IsUnknown() {
		request.ExcludeResourceIdsScope = stringSetValueToSlice(plan.ExcludeResourceIdsScope)
	}

	service := NewConfigService(r.Client())
	_, err := service.AddAggregateConfigRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Config aggregate config rule",
			fmt.Sprintf("AddAggregateConfigRule failed: %v", err),
		)
		return
	}

	// AddAggregateConfigRule only returns a RequestId (no resource-level RuleId).
	// Discover the ConfigRuleId of the newly created rule by listing rules under
	// the account group filtered by the unique rule name.
	items, err := service.ListAggregateConfigRules(ctx, plan.AccountGroupId.ValueString(), plan.RuleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error discovering Config aggregate config rule after create",
			err.Error(),
		)
		return
	}
	var ruleId string
	for _, item := range items {
		if item != nil && item.ConfigRuleId != nil && *item.ConfigRuleId != "" {
			ruleId = *item.ConfigRuleId
			break
		}
	}
	if ruleId == "" {
		resp.Diagnostics.AddError(
			"Error creating Config aggregate config rule",
			"AddAggregateConfigRule succeeded but the new rule could not be found by name; the rule ID is unavailable.",
		)
		return
	}

	plan.RuleId = types.StringValue(ruleId)
	plan.Id = types.StringValue(aggregateConfigRuleId(plan.AccountGroupId.ValueString(), ruleId))

	// Persist the identity BEFORE any read-back: a partial failure here keeps
	// the id in state so `terraform destroy` can clean it up.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := service.DescribeAggregateConfigRule(ctx, ruleId, plan.AccountGroupId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Config aggregate config rule after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"Config aggregate config rule not readable after create",
			fmt.Sprintf("rule %s was created but cannot be read back", ruleId),
		)
		return
	}
	r.refreshModel(ctx, &plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *ConfigAggregateConfigRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, ruleId, ok := splitAggregateConfigRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Config aggregate config rule id is broken, id is %q (expected `<account_group_id>#<rule_id>`)", state.Id.ValueString()),
		)
		return
	}

	service := NewConfigService(r.Client())
	item, err := service.DescribeAggregateConfigRule(ctx, ruleId, accountGroupId)
	if err != nil {
		// Confirmed NotFound => drop from state so Terraform plans a recreation.
		if isAggregateConfigRuleNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Config aggregate config rule", err.Error())
		return
	}
	if item == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.AccountGroupId = types.StringValue(accountGroupId)
	state.RuleId = types.StringValue(ruleId)
	r.refreshModel(ctx, &state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *ConfigAggregateConfigRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, ruleId, ok := splitAggregateConfigRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Config aggregate config rule id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	request := configv20220802.NewUpdateAggregateConfigRuleRequest()
	request.RuleId = helper.String(ruleId)
	request.AccountGroupId = helper.String(accountGroupId)

	// TriggerType, RiskLevel and RuleName are sent from the plan (TriggerType
	// and RiskLevel are required by the Update API). Optional fields are only
	// sent when present in the plan; an omitted (null) field is left unset so
	// the cloud keeps the previous value. The clear-semantics of these fields
	// (absent = keep vs absent = clear) is not documented by the cloud API;
	// "keep" is the safe default and matches the sibling resource behavior.
	request.RuleName = plan.RuleName.ValueStringPointer()
	request.RiskLevel = helper.Int64Uint64(plan.RiskLevel.ValueInt64())

	if !plan.TriggerType.IsNull() && !plan.TriggerType.IsUnknown() {
		request.TriggerType = triggerTypeModelsToSDK(plan.TriggerType)
	}

	if !plan.InputParameter.IsNull() && !plan.InputParameter.IsUnknown() {
		request.InputParameter = inputParameterModelsToSDK(plan.InputParameter)
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}
	if !plan.RegionScope.IsNull() && !plan.RegionScope.IsUnknown() {
		request.RegionScope = stringSetValueToSlice(plan.RegionScope)
	}
	if !plan.TagsScope.IsNull() && !plan.TagsScope.IsUnknown() {
		request.TagsScope = tagsScopeModelsToSDK(plan.TagsScope)
	}
	if !plan.ExcludeResourceIdsScope.IsNull() && !plan.ExcludeResourceIdsScope.IsUnknown() {
		request.ExcludeResourceIdsScope = stringSetValueToSlice(plan.ExcludeResourceIdsScope)
	}

	service := NewConfigService(r.Client())
	if err := service.UpdateAggregateConfigRule(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating Config aggregate config rule",
			fmt.Sprintf("UpdateAggregateConfigRule failed: %v", err),
		)
		return
	}

	item, err := service.DescribeAggregateConfigRule(ctx, ruleId, accountGroupId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Config aggregate config rule after update", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"Config aggregate config rule not readable after update",
			fmt.Sprintf("rule %s was updated but cannot be read back", ruleId),
		)
		return
	}
	plan.Id = state.Id
	plan.RuleId = state.RuleId
	// Immutable fields are not part of the Update request; preserve them from state.
	plan.Identifier = state.Identifier
	plan.IdentifierType = state.IdentifierType
	plan.ResourceType = state.ResourceType
	r.refreshModel(ctx, &plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *ConfigAggregateConfigRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, ruleId, ok := splitAggregateConfigRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Config aggregate config rule id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	service := NewConfigService(r.Client())
	err := service.DeleteAggregateConfigRule(ctx, ruleId, accountGroupId)
	if err != nil {
		// Only a confirmed "already gone" error may be swallowed.
		if !isAggregateConfigRuleNotFound(err) {
			resp.Diagnostics.AddError(
				"Error deleting Config aggregate config rule",
				fmt.Sprintf("DeleteAggregateConfigRule failed: %v", err),
			)
			return
		}
	}
}

// ImportState supports `terraform import` using the `<account_group_id>#<rule_id>` form.
func (r *ConfigAggregateConfigRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := splitAggregateConfigRuleId(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<account_group_id>%s<rule_id>`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// refreshModel copies cloud-owned fields into the model from a ConfigRule.
func (r *ConfigAggregateConfigRuleResource) refreshModel(ctx context.Context, m *ConfigAggregateConfigRuleResourceModel, item *configv20220802.ConfigRule) {
	if item == nil {
		return
	}

	m.Identifier = types.StringPointerValue(item.Identifier)
	m.IdentifierType = types.StringPointerValue(item.IdentifierType)
	m.RuleName = types.StringPointerValue(item.RuleName)
	m.Description = types.StringPointerValue(item.Description)
	m.RuleId = types.StringPointerValue(item.ConfigRuleId)

	if item.RiskLevel != nil {
		m.RiskLevel = types.Int64Value(*item.RiskLevel)
	} else {
		m.RiskLevel = types.Int64Null()
	}

	m.ResourceType = stringSetFromSDK(item.ResourceType)
	m.TriggerType = triggerTypeListFromSDK(ctx, item.TriggerType)
	m.InputParameter = inputParameterListFromSDK(ctx, item.InputParameter)
	m.RegionScope = stringSetFromSDK(item.RegionsScope)
	m.TagsScope = tagsScopeListFromSDK(ctx, item.TagsScope)
	m.ExcludeResourceIdsScope = stringSetFromSDK(item.ExcludeResourceIdsScope)

	m.ConfigRule = buildConfigRuleObject(ctx, item)
}

// stringSetValueToSlice extracts []*string from a types.Set of StringType.
func stringSetValueToSlice(s types.Set) []*string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	var elems []types.String
	_ = s.ElementsAs(context.Background(), &elems, false)
	out := make([]*string, 0, len(elems))
	for _, e := range elems {
		v := e.ValueString()
		out = append(out, &v)
	}
	return out
}

// triggerTypeModelsToSDK converts the plan's trigger_type list to SDK structs.
func triggerTypeModelsToSDK(l types.List) []*configv20220802.TriggerType {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	var elems []TriggerTypeModel
	_ = l.ElementsAs(context.Background(), &elems, false)
	out := make([]*configv20220802.TriggerType, 0, len(elems))
	for _, e := range elems {
		tt := &configv20220802.TriggerType{}
		if !e.MessageType.IsNull() && !e.MessageType.IsUnknown() {
			tt.MessageType = e.MessageType.ValueStringPointer()
		}
		if !e.MaximumExecutionFrequency.IsNull() && !e.MaximumExecutionFrequency.IsUnknown() && e.MaximumExecutionFrequency.ValueString() != "" {
			tt.MaximumExecutionFrequency = e.MaximumExecutionFrequency.ValueStringPointer()
		}
		out = append(out, tt)
	}
	return out
}

// inputParameterModelsToSDK converts the plan's input_parameter list to SDK structs.
func inputParameterModelsToSDK(l types.List) []*configv20220802.InputParameter {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	var elems []InputParameterModel
	_ = l.ElementsAs(context.Background(), &elems, false)
	out := make([]*configv20220802.InputParameter, 0, len(elems))
	for _, e := range elems {
		p := &configv20220802.InputParameter{}
		if !e.ParameterKey.IsNull() && !e.ParameterKey.IsUnknown() {
			p.ParameterKey = e.ParameterKey.ValueStringPointer()
		}
		if !e.Type.IsNull() && !e.Type.IsUnknown() && e.Type.ValueString() != "" {
			p.Type = e.Type.ValueStringPointer()
		}
		if !e.Value.IsNull() && !e.Value.IsUnknown() && e.Value.ValueString() != "" {
			p.Value = e.Value.ValueStringPointer()
		}
		out = append(out, p)
	}
	return out
}

// tagsScopeModelsToSDK converts the plan's tags_scope list to SDK Tag structs.
func tagsScopeModelsToSDK(l types.List) []*configv20220802.Tag {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	var elems []TagsScopeModel
	_ = l.ElementsAs(context.Background(), &elems, false)
	out := make([]*configv20220802.Tag, 0, len(elems))
	for _, e := range elems {
		t := &configv20220802.Tag{}
		if !e.TagKey.IsNull() && !e.TagKey.IsUnknown() && e.TagKey.ValueString() != "" {
			t.TagKey = e.TagKey.ValueStringPointer()
		}
		if !e.TagValue.IsNull() && !e.TagValue.IsUnknown() && e.TagValue.ValueString() != "" {
			t.TagValue = e.TagValue.ValueStringPointer()
		}
		out = append(out, t)
	}
	return out
}

// triggerTypeListFromSDK builds a types.List of trigger_type objects from SDK structs.
func triggerTypeListFromSDK(ctx context.Context, in []*configv20220802.TriggerType) types.List {
	if in == nil {
		return types.ListNull(types.ObjectType{AttrTypes: triggerTypeAttrTypes()})
	}
	elems := make([]TriggerTypeModel, 0, len(in))
	for _, tt := range in {
		model := TriggerTypeModel{MessageType: types.StringPointerValue(tt.MessageType)}
		if tt.MaximumExecutionFrequency != nil {
			model.MaximumExecutionFrequency = types.StringPointerValue(tt.MaximumExecutionFrequency)
		}
		elems = append(elems, model)
	}
	listVal, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: triggerTypeAttrTypes()}, elems)
	return listVal
}

// inputParameterListFromSDK builds a types.List of input_parameter objects from SDK structs.
func inputParameterListFromSDK(ctx context.Context, in []*configv20220802.InputParameter) types.List {
	if in == nil {
		return types.ListNull(types.ObjectType{AttrTypes: inputParameterAttrTypes()})
	}
	elems := make([]InputParameterModel, 0, len(in))
	for _, p := range in {
		elems = append(elems, InputParameterModel{
			ParameterKey: types.StringPointerValue(p.ParameterKey),
			Type:         types.StringPointerValue(p.Type),
			Value:        types.StringPointerValue(p.Value),
		})
	}
	listVal, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: inputParameterAttrTypes()}, elems)
	return listVal
}

// tagsScopeListFromSDK builds a types.List of tags_scope objects from SDK Tag structs.
func tagsScopeListFromSDK(ctx context.Context, in []*configv20220802.Tag) types.List {
	if in == nil {
		return types.ListNull(types.ObjectType{AttrTypes: tagsScopeAttrTypes()})
	}
	elems := make([]TagsScopeModel, 0, len(in))
	for _, t := range in {
		elems = append(elems, TagsScopeModel{
			TagKey:   types.StringPointerValue(t.TagKey),
			TagValue: types.StringPointerValue(t.TagValue),
		})
	}
	listVal, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: tagsScopeAttrTypes()}, elems)
	return listVal
}

// buildConfigRuleObject builds the computed config_rule types.Object from a SDK ConfigRule.
func buildConfigRuleObject(ctx context.Context, item *configv20220802.ConfigRule) types.Object {
	if item == nil {
		return types.ObjectNull(configRuleAttrTypes())
	}

	m := ConfigRuleModel{
		Identifier:            types.StringPointerValue(item.Identifier),
		RuleName:              types.StringPointerValue(item.RuleName),
		ServiceFunction:       types.StringPointerValue(item.ServiceFunction),
		CreateTime:            types.StringPointerValue(item.CreateTime),
		Description:           types.StringPointerValue(item.Description),
		Status:                types.StringPointerValue(item.Status),
		ComplianceResult:      types.StringPointerValue(item.ComplianceResult),
		ConfigRuleInvokedTime: types.StringPointerValue(item.ConfigRuleInvokedTime),
		ConfigRuleId:          types.StringPointerValue(item.ConfigRuleId),
		IdentifierType:        types.StringPointerValue(item.IdentifierType),
		CompliancePackId:      types.StringPointerValue(item.CompliancePackId),
		CompliancePackName:    types.StringPointerValue(item.CompliancePackName),
		AccountGroupId:        types.StringPointerValue(item.AccountGroupId),
		AccountGroupName:      types.StringPointerValue(item.AccountGroupName),
	}

	if item.RiskLevel != nil {
		m.RiskLevel = types.Int64Value(*item.RiskLevel)
	} else {
		m.RiskLevel = types.Int64Null()
	}

	if item.RuleOwnerId != nil {
		m.RuleOwnerId = types.Int64Value(int64(*item.RuleOwnerId))
	} else {
		m.RuleOwnerId = types.Int64Null()
	}

	// input_parameter
	m.InputParameter = inputParameterListFromSDK(ctx, item.InputParameter)

	// source_condition
	if item.SourceCondition == nil {
		m.SourceCondition = types.ListNull(types.ObjectType{AttrTypes: sourceConditionAttrTypes()})
	} else {
		scModels := make([]SourceConditionModel, 0, len(item.SourceCondition))
		for _, sc := range item.SourceCondition {
			scModels = append(scModels, SourceConditionModel{
				EmptyAs:      types.StringPointerValue(sc.EmptyAs),
				SelectPath:   types.StringPointerValue(sc.SelectPath),
				Operator:     types.StringPointerValue(sc.Operator),
				Required:     types.BoolPointerValue(sc.Required),
				DesiredValue: types.StringPointerValue(sc.DesiredValue),
			})
		}
		scList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: sourceConditionAttrTypes()}, scModels)
		m.SourceCondition = scList
	}

	// resource_type
	m.ResourceType = stringSetFromSDK(item.ResourceType)
	// labels
	m.Labels = stringSetFromSDK(item.Labels)
	// trigger_type
	m.TriggerType = triggerTypeListFromSDK(ctx, item.TriggerType)

	// manage_input_parameter
	if item.ManageInputParameter == nil {
		m.ManageInputParameter = types.ListNull(types.ObjectType{AttrTypes: manageInputParameterAttrTypes()})
	} else {
		mipModels := make([]ManageInputParameterModel, 0, len(item.ManageInputParameter))
		for _, mip := range item.ManageInputParameter {
			mipModels = append(mipModels, ManageInputParameterModel{
				ValueType:     types.StringPointerValue(mip.ValueType),
				ParameterKey:  types.StringPointerValue(mip.ParameterKey),
				Type:          types.StringPointerValue(mip.Type),
				DefaultValue:  types.StringPointerValue(mip.DefaultValue),
				Description:   types.StringPointerValue(mip.Description),
			})
		}
		mipList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: manageInputParameterAttrTypes()}, mipModels)
		m.ManageInputParameter = mipList
	}

	// annotation
	if item.Annotation == nil {
		m.Annotation = types.ListNull(types.ObjectType{AttrTypes: annotationAttrTypes()})
	} else {
		annModels := []AnnotationModel{{
			Configuration: types.StringPointerValue(item.Annotation.Configuration),
			DesiredValue:  types.StringPointerValue(item.Annotation.DesiredValue),
			Operator:      types.StringPointerValue(item.Annotation.Operator),
			Property:      types.StringPointerValue(item.Annotation.Property),
		}}
		annList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: annotationAttrTypes()}, annModels)
		m.Annotation = annList
	}

	// regions_scope
	m.RegionsScope = stringSetFromSDK(item.RegionsScope)
	// tags_scope
	m.TagsScope = tagsScopeListFromSDK(ctx, item.TagsScope)
	// exclude_resource_ids_scope
	m.ExcludeResourceIdsScope = stringSetFromSDK(item.ExcludeResourceIdsScope)
	// manage_trigger_type
	m.ManageTriggerType = stringSetFromSDK(item.ManageTriggerType)

	objVal, _ := types.ObjectValueFrom(ctx, configRuleAttrTypes(), &m)
	return objVal
}

// stringSetFromSDK builds a types.Set of StringType from []*string.
func stringSetFromSDK(in []*string) types.Set {
	if in == nil {
		return types.SetNull(types.StringType)
	}
	elems := make([]attr.Value, 0, len(in))
	for _, v := range in {
		elems = append(elems, types.StringPointerValue(v))
	}
	setVal, _ := types.SetValue(types.StringType, elems)
	return setVal
}

// aggregateConfigRuleId joins accountGroupId and ruleId with the FILED_SP separator.
func aggregateConfigRuleId(accountGroupId, ruleId string) string {
	return strings.Join([]string{accountGroupId, ruleId}, tccommon.FILED_SP)
}

// splitAggregateConfigRuleId parses `<account_group_id>#<rule_id>`.
func splitAggregateConfigRuleId(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// isAggregateConfigRuleNotFound reports whether the error means the rule is gone.
func isAggregateConfigRuleNotFound(err error) bool {
	if err == nil {
		return false
	}
	sdkErr, ok := err.(*sdkerr.TencentCloudSDKError)
	if !ok {
		return false
	}
	switch sdkErr.Code {
	case "ResourceNotFound":
		return true
	default:
		return false
	}
}

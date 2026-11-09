package config

import (
	"context"
	"errors"
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
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions (repo convention: only assert the
// extension interfaces; the base interface is checked implicitly by the
// factory return value).
var (
	_ resource.ResourceWithConfigure   = &ConfigAggregateCompliancePackResource{}
	_ resource.ResourceWithImportState = &ConfigAggregateCompliancePackResource{}

	// errAggregateCompliancePackNotFound marks an empty Describe response. The
	// Config API answers with an empty body once the pack is gone, so the case
	// must be told apart from a transport or API failure.
	errAggregateCompliancePackNotFound = errors.New("aggregate compliance pack not found")

	// errAggregateCompliancePackStatusPending marks a status transition that has
	// been requested but is not visible in Describe yet, so the wait keeps
	// polling instead of failing.
	errAggregateCompliancePackStatusPending = errors.New("compliance pack status not converged")
)

// Status values accepted by UpdateAggregateCompliancePackStatus. A pack is
// created in ACTIVE state; UN_ACTIVE disables it.
const (
	aggregateCompliancePackStatusActive   = "ACTIVE"
	aggregateCompliancePackStatusUnActive = "UN_ACTIVE"
)

// NewConfigAggregateCompliancePackResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewConfigAggregateCompliancePackResource() resource.Resource {
	return &ConfigAggregateCompliancePackResource{}
}

// ConfigAggregateCompliancePackResource implements resource.Resource for
// tencentcloud_config_aggregate_compliance_pack.
//
// The cloud side addresses a pack with an aggregator account group id plus a
// compliance pack id, so the Terraform id is both parts joined with `#`
// (tccommon.FILED_SP) and split again for every CRUD call.
type ConfigAggregateCompliancePackResource struct {
	fw.ResourceWithConfigure
}

// ConfigAggregateCompliancePackResourceModel maps the top-level schema
// attributes. Every schema attribute appears here exactly once.
type ConfigAggregateCompliancePackResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	CompliancePackName types.String `tfsdk:"compliance_pack_name"`
	AccountGroupId     types.String `tfsdk:"account_group_id"`
	RiskLevel          types.Int64  `tfsdk:"risk_level"`
	Description        types.String `tfsdk:"description"`
	ConfigRules        types.List   `tfsdk:"config_rules"`
	Status             types.String `tfsdk:"status"`
	CompliancePackId   types.String `tfsdk:"compliance_pack_id"`
	CreateTime         types.String `tfsdk:"create_time"`
}

// aggregateCompliancePackRuleModel maps one element of config_rules.
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

// aggregateCompliancePackInputParameterModel maps one element of
// config_rules[*].input_parameter.
type aggregateCompliancePackInputParameterModel struct {
	ParameterKey types.String `tfsdk:"parameter_key"`
	Type         types.String `tfsdk:"type"`
	Value        types.String `tfsdk:"value"`
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
				Description: "Account group id. Immutable: changing it forces resource replacement.",
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
				Computed:    true,
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
										Computed:    true,
										Description: "Parameter type. Required: `Require`, optional: `Optional`.",
									},
									"value": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
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
							Computed:    true,
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
				Optional: true,
				Computed: true,
				Description: "Compliance pack status. Valid values: `ACTIVE` (enabled) and `UN_ACTIVE` (disabled). " +
					"A pack is created as `ACTIVE`; set `UN_ACTIVE` to disable it.",
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

// Create handles `terraform apply` for new resources.
func (r *ConfigAggregateCompliancePackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConfigAggregateCompliancePackResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	accountGroupId := plan.AccountGroupId.ValueString()

	request := configv20220802.NewAddAggregateCompliancePackRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.CompliancePackName = fwhelper.StringPointerFromValue(plan.CompliancePackName)
	request.RiskLevel = fwhelper.Uint64PointerFromValue(plan.RiskLevel)
	request.Description = fwhelper.StringPointerFromValue(plan.Description)

	rules, ruleDiags := expandAggregateCompliancePackRules(ctx, plan.ConfigRules)
	resp.Diagnostics.Append(ruleDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ConfigRules = rules

	compliancePackId, err := service.AddAggregateCompliancePack(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Config aggregate compliance pack",
			fmt.Sprintf("AddAggregateCompliancePack failed (%s): %s", aggregateCompliancePackScope(accountGroupId, ""), err),
		)
		return
	}

	plan.CompliancePackId = types.StringValue(compliancePackId)
	plan.Id = types.StringValue(aggregateCompliancePackId(accountGroupId, compliancePackId))

	// AddAggregateCompliancePack only guarantees an id, so wait until the pack
	// is actually readable before the status API is allowed to touch it.
	if err := waitAggregateCompliancePackCreated(ctx, service, accountGroupId, compliancePackId); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Config aggregate compliance pack creation",
			fmt.Sprintf("The pack was created but never became readable (%s): %s", aggregateCompliancePackScope(accountGroupId, compliancePackId), err),
		)
		return
	}

	// A new pack is always created in ACTIVE state, so only an explicit
	// UN_ACTIVE request needs the follow-up status call.
	if plan.Status.ValueString() == aggregateCompliancePackStatusUnActive {
		if err := updateAggregateCompliancePackStatus(ctx, service, accountGroupId, compliancePackId, aggregateCompliancePackStatusUnActive); err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Config aggregate compliance pack after create",
				fmt.Sprintf("UpdateAggregateCompliancePackStatus failed after the pack was created (%s): %s", aggregateCompliancePackScope(accountGroupId, compliancePackId), err),
			)
			return
		}
	}

	// Read the pack back so the very first state already holds the cloud-owned
	// attributes (status, create_time, normalized rules).
	item, err := describeAggregateCompliancePack(ctx, service, accountGroupId, compliancePackId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Config aggregate compliance pack after create",
			fmt.Sprintf(
				"AddAggregateCompliancePack succeeded but the pack cannot be read back (%s). "+
					"Import it with `terraform import tencentcloud_config_aggregate_compliance_pack.<name> %s` to keep it managed: %s",
				aggregateCompliancePackScope(accountGroupId, compliancePackId), plan.Id.ValueString(), err,
			),
		)
		return
	}

	resp.Diagnostics.Append(refreshAggregateCompliancePackModel(ctx, &plan, accountGroupId, compliancePackId, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *ConfigAggregateCompliancePackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConfigAggregateCompliancePackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, compliancePackId, ok := parseAggregateCompliancePackId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	item, err := describeAggregateCompliancePack(ctx, service, accountGroupId, compliancePackId)
	if err != nil {
		// A confirmed "not found" means the pack was removed outside Terraform,
		// so it is dropped from state instead of failing the refresh.
		if errors.Is(err, errAggregateCompliancePackNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("Config aggregate compliance pack %s no longer exists, removing it from state.", state.Id.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading Config aggregate compliance pack",
			fmt.Sprintf("DescribeAggregateCompliancePack failed (%s): %s", aggregateCompliancePackScope(accountGroupId, compliancePackId), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAggregateCompliancePackModel(ctx, &state, accountGroupId, compliancePackId, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	accountGroupId, compliancePackId, ok := parseAggregateCompliancePackId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	request := configv20220802.NewUpdateAggregateCompliancePackRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.CompliancePackId = helper.String(compliancePackId)
	request.CompliancePackName = fwhelper.StringPointerFromValue(plan.CompliancePackName)
	request.RiskLevel = fwhelper.Uint64PointerFromValue(plan.RiskLevel)
	request.Description = fwhelper.StringPointerFromValue(plan.Description)

	rules, ruleDiags := expandAggregateCompliancePackRules(ctx, plan.ConfigRules)
	resp.Diagnostics.Append(ruleDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ConfigRules = rules

	if err := service.UpdateAggregateCompliancePack(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating Config aggregate compliance pack",
			fmt.Sprintf("UpdateAggregateCompliancePack failed (%s): %s", aggregateCompliancePackScope(accountGroupId, compliancePackId), err),
		)
		return
	}

	// status is driven by a dedicated API, so it is only called when the
	// practitioner actually changed it.
	if aggregateCompliancePackStatusChanged(plan.Status, state.Status) {
		status := plan.Status.ValueString()
		if err := updateAggregateCompliancePackStatus(ctx, service, accountGroupId, compliancePackId, status); err != nil {
			resp.Diagnostics.AddError(
				"Error updating Config aggregate compliance pack status",
				fmt.Sprintf("UpdateAggregateCompliancePackStatus to %q failed (%s): %s", status, aggregateCompliancePackScope(accountGroupId, compliancePackId), err),
			)
			return
		}
	}

	item, err := describeAggregateCompliancePack(ctx, service, accountGroupId, compliancePackId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Config aggregate compliance pack after update",
			fmt.Sprintf("UpdateAggregateCompliancePack succeeded but the pack cannot be read back (%s): %s", aggregateCompliancePackScope(accountGroupId, compliancePackId), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAggregateCompliancePackModel(ctx, &plan, accountGroupId, compliancePackId, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *ConfigAggregateCompliancePackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConfigAggregateCompliancePackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, compliancePackId, ok := parseAggregateCompliancePackId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	request := configv20220802.NewDeleteAggregateCompliancePackRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.CompliancePackId = helper.String(compliancePackId)

	if err := service.DeleteAggregateCompliancePack(ctx, request); err != nil {
		// An already deleted pack is a successful deletion.
		if tccommon.IsExpectError(err, []string{"ResourceNotFound", "ResourceInexistent"}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting Config aggregate compliance pack",
			fmt.Sprintf("DeleteAggregateCompliancePack failed (%s): %s", aggregateCompliancePackScope(accountGroupId, compliancePackId), err),
		)
	}
}

// ImportState supports `terraform import` using the
// `<account_group_id>#<compliance_pack_id>` form. The remaining attributes are
// filled in by the Read that follows.
func (r *ConfigAggregateCompliancePackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := splitAggregateCompliancePackId(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Import id %q is invalid, expected `<account_group_id>%s<compliance_pack_id>`.", req.ID, tccommon.FILED_SP),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// service returns the Config service bound to the shared client. When the
// provider has not been configured it records a diagnostic and reports false,
// so the CRUD methods never dereference a nil SDK client.
func (r *ConfigAggregateCompliancePackResource) service(diags *diag.Diagnostics) (ConfigService, bool) {
	if r.Client() == nil {
		diags.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return ConfigService{}, false
	}
	return NewConfigService(r.Client()), true
}

// describeAggregateCompliancePack reads the pack from the cloud. An empty
// response means the pack no longer exists and is reported as
// errAggregateCompliancePackNotFound; API and transport failures are returned
// verbatim.
func describeAggregateCompliancePack(ctx context.Context, service ConfigService, accountGroupId, compliancePackId string) (*configv20220802.DescribeAggregateCompliancePackResponseParams, error) {
	item, err := service.DescribeAggregateCompliancePackById(ctx, accountGroupId, compliancePackId)
	if err != nil {
		return nil, err
	}
	if item == nil || item.CompliancePackId == nil {
		return nil, errAggregateCompliancePackNotFound
	}
	return item, nil
}

// aggregateCompliancePackStatusChanged reports whether the planned status is a
// usable value that differs from the one recorded in state.
func aggregateCompliancePackStatusChanged(plan, state types.String) bool {
	if plan.IsNull() || plan.IsUnknown() {
		return false
	}
	return !plan.Equal(state)
}

// waitAggregateCompliancePackCreated polls Describe until the pack created by
// AddAggregateCompliancePack is readable. The Add call only returns an id, so
// the status API must not be invoked before the pack actually exists.
func waitAggregateCompliancePackCreated(ctx context.Context, service ConfigService, accountGroupId, compliancePackId string) error {
	packId := aggregateCompliancePackId(accountGroupId, compliancePackId)

	_, err := fwhelper.RetryFramework(ctx, tccommon.WriteRetryTimeout, func(err error) bool {
		return errors.Is(err, errAggregateCompliancePackNotFound)
	}, func() (string, error) {
		item, err := describeAggregateCompliancePack(ctx, service, accountGroupId, compliancePackId)
		if err != nil {
			return "", err
		}
		return stringValueOf(item.Status), nil
	})
	if err != nil {
		if fwhelper.IsTimeoutError(err) {
			return fmt.Errorf("timed out waiting for compliance pack %s to be created: %w", packId, err)
		}
		return err
	}
	return nil
}

// updateAggregateCompliancePackStatus switches the pack between ACTIVE and
// UN_ACTIVE through the dedicated status API.
//
// The transition is asynchronous, so the pack is polled until the cloud reports
// the requested value: writing back a stale status would make Terraform report
// "provider produced inconsistent result after apply".
func updateAggregateCompliancePackStatus(ctx context.Context, service ConfigService, accountGroupId, compliancePackId, status string) error {
	if status != aggregateCompliancePackStatusActive && status != aggregateCompliancePackStatusUnActive {
		return fmt.Errorf("invalid status %q, expected %q or %q", status, aggregateCompliancePackStatusActive, aggregateCompliancePackStatusUnActive)
	}

	request := configv20220802.NewUpdateAggregateCompliancePackStatusRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.CompliancePackId = helper.String(compliancePackId)
	request.Status = helper.String(status)

	if err := service.UpdateAggregateCompliancePackStatus(ctx, request); err != nil {
		return err
	}

	packId := aggregateCompliancePackId(accountGroupId, compliancePackId)
	_, err := fwhelper.RetryFramework(ctx, tccommon.WriteRetryTimeout, func(err error) bool {
		return errors.Is(err, errAggregateCompliancePackStatusPending)
	}, func() (string, error) {
		item, err := describeAggregateCompliancePack(ctx, service, accountGroupId, compliancePackId)
		if err != nil {
			return "", err
		}

		current := stringValueOf(item.Status)
		if current == status {
			return current, nil
		}
		return current, fmt.Errorf("%w: compliance pack %s is in status %q, waiting for %q", errAggregateCompliancePackStatusPending, packId, current, status)
	})
	if err != nil {
		if fwhelper.IsTimeoutError(err) {
			return fmt.Errorf("timed out waiting for compliance pack %s to reach status %q: %w", packId, status, err)
		}
		return err
	}
	return nil
}

// aggregateCompliancePackId builds the Terraform id from its two cloud parts.
func aggregateCompliancePackId(accountGroupId, compliancePackId string) string {
	return strings.Join([]string{accountGroupId, compliancePackId}, tccommon.FILED_SP)
}

// splitAggregateCompliancePackId parses `<account_group_id>#<compliance_pack_id>`.
func splitAggregateCompliancePackId(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// parseAggregateCompliancePackId resolves the Terraform id into the pair of ids
// the Config API expects, recording a diagnostic when the id is malformed.
func parseAggregateCompliancePackId(diags *diag.Diagnostics, id string) (accountGroupId, compliancePackId string, ok bool) {
	accountGroupId, compliancePackId, ok = splitAggregateCompliancePackId(id)
	if !ok {
		diags.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Resource id %q is invalid, expected `<account_group_id>%s<compliance_pack_id>`.", id, tccommon.FILED_SP),
		)
	}
	return
}

// aggregateCompliancePackScope renders the cloud-side identity of a pack for
// diagnostics, so every failure can be traced back to a concrete account group
// and compliance pack.
func aggregateCompliancePackScope(accountGroupId, compliancePackId string) string {
	if compliancePackId == "" {
		return fmt.Sprintf("account_group_id=%q", accountGroupId)
	}
	return fmt.Sprintf("account_group_id=%q, compliance_pack_id=%q", accountGroupId, compliancePackId)
}

// refreshAggregateCompliancePackModel copies the cloud state into the model.
// Attributes the Describe API may omit keep their plan/state value, so a value
// configured by the practitioner is never nulled out after apply.
func refreshAggregateCompliancePackModel(ctx context.Context, m *ConfigAggregateCompliancePackResourceModel, accountGroupId, compliancePackId string, item *configv20220802.DescribeAggregateCompliancePackResponseParams) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		diags.AddError(
			"Malformed Config aggregate compliance pack response",
			fmt.Sprintf("DescribeAggregateCompliancePack returned an empty response (%s).", aggregateCompliancePackScope(accountGroupId, compliancePackId)),
		)
		return diags
	}

	m.AccountGroupId = types.StringValue(accountGroupId)
	m.CompliancePackId = types.StringValue(compliancePackId)
	m.Id = types.StringValue(aggregateCompliancePackId(accountGroupId, compliancePackId))

	m.CompliancePackName = fwhelper.StringValueOrPrev(item.CompliancePackName, m.CompliancePackName)
	m.Description = fwhelper.StringValueOrPrev(item.Description, m.Description)
	m.RiskLevel = fwhelper.Int64ValueOrPrev(item.RiskLevel, m.RiskLevel)
	m.Status = types.StringPointerValue(item.Status)
	m.CreateTime = types.StringPointerValue(item.CreateTime)

	rules, ruleDiags := flattenAggregateComplianceConfigRules(ctx, item.ConfigRules, m.ConfigRules)
	diags.Append(ruleDiags...)
	m.ConfigRules = rules

	// Required attributes must never be written as null, otherwise Terraform
	// rejects the refreshed object.
	if !diags.HasError() && (m.CompliancePackName.IsNull() || m.RiskLevel.IsNull() || m.ConfigRules.IsNull()) {
		diags.AddError(
			"Malformed Config aggregate compliance pack response",
			fmt.Sprintf("DescribeAggregateCompliancePack returned no compliance_pack_name, risk_level or config_rules (%s).", aggregateCompliancePackScope(accountGroupId, compliancePackId)),
		)
	}
	return diags
}

// expandAggregateCompliancePackRules converts the planned config_rules list
// into the SDK rule slice.
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
		parameters, parameterDiags := expandAggregateInputParameters(ctx, rule.InputParameter)
		diags.Append(parameterDiags...)
		if diags.HasError() {
			return nil, diags
		}

		out = append(out, &configv20220802.CompliancePackRule{
			RuleName:              fwhelper.StringPointerFromValue(rule.RuleName),
			RiskLevel:             fwhelper.Uint64PointerFromValue(rule.RiskLevel),
			InputParameter:        parameters,
			Identifier:            fwhelper.StringPointerFromValue(rule.Identifier),
			Description:           fwhelper.StringPointerFromValue(rule.Description),
			ManagedRuleIdentifier: fwhelper.StringPointerFromValue(rule.ManagedRuleIdentifier),
			ConfigRuleId:          fwhelper.StringPointerFromValue(rule.ConfigRuleId),
			CompliancePackId:      fwhelper.StringPointerFromValue(rule.CompliancePackId),
		})
	}
	return out, diags
}

// expandAggregateInputParameters converts the planned input_parameter list into
// the SDK parameter slice.
func expandAggregateInputParameters(ctx context.Context, list types.List) ([]*configv20220802.InputParameter, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var params []aggregateCompliancePackInputParameterModel
	diags.Append(list.ElementsAs(ctx, &params, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*configv20220802.InputParameter, 0, len(params))
	for _, param := range params {
		out = append(out, &configv20220802.InputParameter{
			ParameterKey: fwhelper.StringPointerFromValue(param.ParameterKey),
			Type:         fwhelper.StringPointerFromValue(param.Type),
			Value:        fwhelper.StringPointerFromValue(param.Value),
		})
	}
	return out, diags
}

// aggregateCompliancePackRuleIndex maps every identifier the cloud may use for
// a rule to its previous state, so fields the API does not echo back can be
// carried over.
type aggregateCompliancePackRuleIndex map[string]aggregateCompliancePackRuleModel

func newAggregateCompliancePackRuleIndex(prev []aggregateCompliancePackRuleModel) aggregateCompliancePackRuleIndex {
	index := make(aggregateCompliancePackRuleIndex, len(prev))
	for _, rule := range prev {
		for _, key := range ruleMatchKeys(rule.Identifier.ValueString(), rule.ConfigRuleId.ValueString(), rule.RuleName.ValueString()) {
			if _, exists := index[key]; !exists {
				index[key] = rule
			}
		}
	}
	return index
}

// lookup returns the previous counterpart of a cloud rule. The zero model is
// returned when nothing matches, which behaves as "every field is null".
func (index aggregateCompliancePackRuleIndex) lookup(keys []string) aggregateCompliancePackRuleModel {
	for _, key := range keys {
		if rule, ok := index[key]; ok {
			return rule
		}
	}
	return aggregateCompliancePackRuleModel{}
}

// ruleMatchKeys lists the non-empty identifiers that correlate a cloud rule
// with its previous counterpart, most specific first.
func ruleMatchKeys(identifier, configRuleId, ruleName string) []string {
	keys := make([]string, 0, 3)
	for _, key := range []string{identifier, configRuleId, ruleName} {
		if key != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

// flattenAggregateComplianceConfigRules builds the config_rules list from the
// Describe response. The previous plan/state list is used to fill in fields the
// cloud does not echo back (managed_rule_identifier, compliance_pack_id) and to
// keep configured values when the API omits them.
func flattenAggregateComplianceConfigRules(ctx context.Context, rules []*configv20220802.ComplianceConfigRule, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := configRulesObjectType()

	prev := make([]aggregateCompliancePackRuleModel, 0)
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prev, false)...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
	}
	index := newAggregateCompliancePackRuleIndex(prev)

	out := make([]aggregateCompliancePackRuleModel, 0, len(rules))
	for i, rule := range rules {
		if rule == nil {
			continue // a malformed response must never panic the provider
		}

		prior := index.lookup(ruleMatchKeys(stringValueOf(rule.Identifier), stringValueOf(rule.ConfigRuleId), stringValueOf(rule.RuleName)))

		parameters, parameterDiags := flattenAggregateInputParameters(ctx, rule.InputParameter, prior.InputParameter)
		diags.Append(parameterDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		ruleModel := aggregateCompliancePackRuleModel{
			RuleName:              fwhelper.StringValueOrPrev(rule.RuleName, prior.RuleName),
			RiskLevel:             fwhelper.Int64ValueOrPrev(rule.RiskLevel, prior.RiskLevel),
			Identifier:            fwhelper.StringValueOrPrev(rule.Identifier, prior.Identifier),
			Description:           fwhelper.StringValueOrPrev(rule.Description, prior.Description),
			ConfigRuleId:          fwhelper.StringValueOrPrev(rule.ConfigRuleId, prior.ConfigRuleId),
			ManagedRuleIdentifier: fwhelper.StringValueOrPrev(nil, prior.ManagedRuleIdentifier),
			CompliancePackId:      fwhelper.StringValueOrPrev(nil, prior.CompliancePackId),
			InputParameter:        parameters,
		}

		// rule_name and risk_level are Required: writing null would make
		// Terraform reject the refreshed object, so fail with a traceable
		// message instead.
		if ruleModel.RuleName.IsNull() || ruleModel.RiskLevel.IsNull() {
			diags.AddError(
				"Malformed Config aggregate compliance pack response",
				fmt.Sprintf("config_rules[%d] carries no rule_name or risk_level in the DescribeAggregateCompliancePack response.", i),
			)
			return types.ListNull(objType), diags
		}

		out = append(out, ruleModel)
	}

	// A Required list must stay consistent with what was planned, so an empty
	// response keeps the previous value instead of clearing configured rules.
	if len(out) == 0 && !prevList.IsNull() && !prevList.IsUnknown() {
		return prevList, diags
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAggregateInputParameters builds the input_parameter list from the
// Describe response. The previous list is used, position by position, to keep
// values the API omits.
func flattenAggregateInputParameters(ctx context.Context, params []*configv20220802.InputParameter, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := inputParameterObjectType()

	prev := make([]aggregateCompliancePackInputParameterModel, 0)
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prev, false)...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
	}

	out := make([]aggregateCompliancePackInputParameterModel, 0, len(params))
	for _, param := range params {
		if param == nil {
			continue // a malformed response must never panic the provider
		}

		// The API keeps parameter order stable, so the entry at the same
		// position holds the previously configured values.
		var prior aggregateCompliancePackInputParameterModel
		if len(out) < len(prev) {
			prior = prev[len(out)]
		}

		parameterModel := aggregateCompliancePackInputParameterModel{
			ParameterKey: fwhelper.StringValueOrPrev(param.ParameterKey, prior.ParameterKey),
			Type:         fwhelper.StringValueOrPrev(param.Type, prior.Type),
			Value:        fwhelper.StringValueOrPrev(param.Value, prior.Value),
		}
		if parameterModel.ParameterKey.IsNull() {
			diags.AddError(
				"Malformed Config aggregate compliance pack response",
				fmt.Sprintf("input_parameter[%d] carries no parameter_key in the DescribeAggregateCompliancePack response.", len(out)),
			)
			return types.ListNull(objType), diags
		}

		out = append(out, parameterModel)
	}

	// Same rule as for config_rules: an empty response must not clear a
	// Required list that was planned with entries.
	if len(out) == 0 && !prevList.IsNull() && !prevList.IsUnknown() {
		return prevList, diags
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// configRulesObjectType is the object type of the config_rules elements.
func configRulesObjectType() attr.Type {
	return types.ObjectType{AttrTypes: configRulesAttrTypes()}
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
		"input_parameter":         types.ListType{ElemType: inputParameterObjectType()},
	}
}

// inputParameterObjectType is the object type of the input_parameter elements.
func inputParameterObjectType() attr.Type {
	return types.ObjectType{AttrTypes: inputParameterAttrTypes()}
}

// inputParameterAttrTypes is the object type of the input_parameter nested
// object, used to build types.List values.
func inputParameterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"parameter_key": types.StringType,
		"type":          types.StringType,
		"value":         types.StringType,
	}
}

// stringValueOf dereferences an optional string, treating nil as empty. It is
// only used to normalise optional SDK fields into request-side match keys.
func stringValueOf(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

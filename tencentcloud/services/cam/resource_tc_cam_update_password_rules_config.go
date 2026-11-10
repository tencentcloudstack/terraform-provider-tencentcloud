package cam

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// Compile-time interface assertions. The base resource.Resource interface is
// checked implicitly by the factory return value.
var (
	_ resource.ResourceWithConfigure   = &CamUpdatePasswordRulesConfigResource{}
	_ resource.ResourceWithImportState = &CamUpdatePasswordRulesConfigResource{}
)

// passwordRulesConfigId is the constant state id of this singleton resource:
// the CAM password rules are a single account-wide configuration, so there is
// no cloud-side identifier to model.
const passwordRulesConfigId = "cam-password-rules"

// errPasswordRulesNotFound marks a GetPasswordRules response that carries no
// rules object. It is reported separately from transport and API failures so the
// caller can describe the situation precisely.
var errPasswordRulesNotFound = errors.New("GetPasswordRules returned no password rules")

// passwordRulesZeroSentinel is the UpdateTime value CAM reports for a password
// rules configuration that has never been modified. Live-account testing
// (2026-11) showed that unmodified accounts do NOT report all-zero rules:
// a fresh account answers GetPasswordRules with e.g. MinimumLength=8,
// MustContain="Aa1!", RetryPasswordLimit=10 while UpdateTime is still the zero
// sentinel and Modifier is empty. Only a real modification (console, API or
// this resource) fills in a timestamp and the modifier UIN.
const passwordRulesZeroSentinel = "0000-00-00 00:00:00"

// NewCamUpdatePasswordRulesConfigResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewCamUpdatePasswordRulesConfigResource() resource.Resource {
	return &CamUpdatePasswordRulesConfigResource{}
}

// CamUpdatePasswordRulesConfigResource implements resource.Resource for
// tencentcloud_cam_update_password_rules.
//
// This is a configuration resource: it never creates or deletes anything on the
// cloud side. Create and Update both call UpdatePasswordRules, Read calls
// GetPasswordRules, and Delete only drops the resource from state so that the
// cloud-side password rules stay untouched.
type CamUpdatePasswordRulesConfigResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// CamUpdatePasswordRulesConfigRuleModel maps the rules nested attribute.
type CamUpdatePasswordRulesConfigRuleModel struct {
	MinimumLength                  types.Int64  `tfsdk:"minimum_length"`
	MustContain                    types.String `tfsdk:"must_contain"`
	ForcePasswordChange            types.Int64  `tfsdk:"force_password_change"`
	ReusePasswordLimit             types.Int64  `tfsdk:"reuse_password_limit"`
	RetryPasswordLimit             types.Int64  `tfsdk:"retry_password_limit"`
	PasswordExpirationInvalidation types.Int64  `tfsdk:"password_expiration_invalidation"`
}

// CamUpdatePasswordRulesConfigResourceModel maps every top-level schema
// attribute; each one appears here exactly once.
type CamUpdatePasswordRulesConfigResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Rules      types.Object `tfsdk:"rules"`
	UpdateTime types.String `tfsdk:"update_time"`
	Modifier   types.String `tfsdk:"modifier"`
}

// passwordRulesAttrTypes returns the attribute types of the rules object. It is
// shared by every conversion so the schema-facing shape never drifts.
func passwordRulesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"minimum_length":                   types.Int64Type,
		"must_contain":                     types.StringType,
		"force_password_change":            types.Int64Type,
		"reuse_password_limit":             types.Int64Type,
		"retry_password_limit":             types.Int64Type,
		"password_expiration_invalidation": types.Int64Type,
	}
}

// Metadata reports the Terraform type name.
func (r *CamUpdatePasswordRulesConfigResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cam_update_password_rules"
}

// Schema declares the password rules configuration.
func (r *CamUpdatePasswordRulesConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage the CAM password rules configuration. " +
			"Create and update both call the UpdatePasswordRules API, read calls GetPasswordRules. " +
			"Destroying this resource only removes it from state; the cloud-side password rules are left unchanged.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the password rules configuration, always `cam-password-rules` because the CAM password rules are a single account-wide configuration.",
			},
			"rules": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Password rules.",
				Attributes: map[string]schema.Attribute{
					"minimum_length": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Minimum password length.",
					},
					"must_contain": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Characters that the password must contain, as a literal character set (e.g. `Aa1!` for upper, lower, digit and special characters). Note that the CAM console re-orders this value (e.g. `Aa1!` is stored and read back as `1!aA`); invalid literals are rejected with `InvalidParameter.PasswordRuleError`.",
					},
					"force_password_change": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Forced password change period.",
					},
					"reuse_password_limit": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Number of previous passwords that cannot be reused.",
					},
					"retry_password_limit": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Number of allowed password retries.",
					},
					"password_expiration_invalidation": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Whether an expired password becomes invalid. Valid values: `1` (the sub-user cannot log in and an administrator must reset the password), `2` (the sub-user can log in and is forced to change the password).",
					},
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last modification time of the password rules.",
			},
			"modifier": schema.StringAttribute{
				Computed:    true,
				Description: "Last modifier of the password rules.",
			},
		},
	}
}

// Create sets the password rules for the first time. It shares the write path
// with Update because the cloud API exposes a single write operation.
func (r *CamUpdatePasswordRulesConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CamUpdatePasswordRulesConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	// Warn instead of error when the cloud-side rules look like they are
	// already managed: the rules are account-wide, so two configurations
	// would silently overwrite each other instead of showing up as a plan
	// diff. The detection relies on the zero-sentinel UpdateTime (see
	// passwordRulesZeroSentinel), which is a heuristic; an unmodified
	// account must never be blocked from adopting the configuration, so the
	// signal is surfaced as a warning and the real values can always be
	// checked with `terraform import` + `terraform plan` first.
	current, err := service.GetPasswordRules(ctx)
	if err != nil && !errors.Is(err, errPasswordRulesNotFound) {
		resp.Diagnostics.AddError("Error reading CAM password rules before create", err.Error())
		return
	}
	if err == nil && passwordRulesTouched(current) {
		resp.Diagnostics.AddWarning(
			"Existing CAM password rules detected",
			fmt.Sprintf("The CAM password rules look like they have already been modified (last modifier %q, last update time %q). "+
				"Applying this configuration will overwrite them. Run `terraform import tencentcloud_cam_update_password_rules.example %s` "+
				"and inspect the plan first if the existing values should be kept.",
				stringOrEmpty(current.Modifier), stringOrEmpty(current.UpdateTime), passwordRulesConfigId),
		)
	}

	r.applyPasswordRules(ctx, service, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(passwordRulesConfigId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the password rules from the cloud.
func (r *CamUpdatePasswordRulesConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CamUpdatePasswordRulesConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	params, err := service.GetPasswordRules(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading CAM password rules", err.Error())
		return
	}

	state.Id = types.StringValue(passwordRulesConfigId)
	state.UpdateTime = stringPointerValue(params.UpdateTime)
	state.Modifier = stringPointerValue(params.Modifier)
	state.Rules = flattenPasswordRules(ctx, params.Rules, state.Rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update applies the planned password rules.
func (r *CamUpdatePasswordRulesConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CamUpdatePasswordRulesConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	r.applyPasswordRules(ctx, service, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(passwordRulesConfigId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is required by the resource.Resource interface. The CAM password rules
// are an account-wide configuration that always exists, so destroying the
// resource only removes it from state and leaves the cloud-side rules
// unchanged.
func (r *CamUpdatePasswordRulesConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	return
}

// service returns the CAM service bound to the shared client. When the provider
// has not been configured it records a diagnostic and reports false, so the
// lifecycle callbacks never dereference a nil SDK client.
func (r *CamUpdatePasswordRulesConfigResource) service(diags *diag.Diagnostics) (CamService, bool) {
	if r.Client() == nil {
		diags.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return CamService{}, false
	}
	return NewCamService(r.Client()), true
}

// applyPasswordRules writes the planned rules through UpdatePasswordRules and
// then reads the configuration back, so every Optional+Computed value is known
// when the state is saved. The state always records the cloud-reported values:
// the CAM API accepts a write and may normalize it (e.g. must_contain), and
// such differences are surfaced by the next plan as a regular diff instead of
// being second-guessed here. Values the cloud rejects outright fail at the
// UpdatePasswordRules call with the API's own InvalidParameter error.
func (r *CamUpdatePasswordRulesConfigResource) applyPasswordRules(ctx context.Context, service CamService, plan *CamUpdatePasswordRulesConfigResourceModel, diags *diag.Diagnostics) {
	rule := expandPasswordRules(ctx, plan.Rules, diags)
	if diags.HasError() {
		return
	}

	if err := service.UpdatePasswordRules(ctx, rule); err != nil {
		diags.AddError("Error updating CAM password rules", err.Error())
		return
	}

	params, err := service.GetPasswordRules(ctx)
	if err != nil {
		diags.AddError("Error reading CAM password rules after update", err.Error())
		return
	}

	observed := flattenPasswordRules(ctx, params.Rules, plan.Rules, diags)
	if diags.HasError() {
		return
	}

	plan.Id = types.StringValue(passwordRulesConfigId)
	plan.Rules = observed
	plan.UpdateTime = stringPointerValue(params.UpdateTime)
	plan.Modifier = stringPointerValue(params.Modifier)
}

// passwordRulesTouched reports whether the cloud-side rules have ever been
// modified. CAM reports UpdateTime "0000-00-00 00:00:00" and an empty Modifier
// for a configuration that was never changed, even though the rule values
// themselves are non-zero defaults (e.g. MinimumLength=8, MustContain="Aa1!").
// A real timestamp is therefore the only reliable signal that somebody already
// manages the rules and a create would silently overwrite their configuration.
func passwordRulesTouched(params *cam.GetPasswordRulesResponseParams) bool {
	if params == nil {
		return false
	}
	if params.UpdateTime == nil || *params.UpdateTime == "" || *params.UpdateTime == passwordRulesZeroSentinel {
		return false
	}
	return true
}

// expandPasswordRules converts the rules object into the SDK request struct.
// Fields left unset in the plan are omitted from the request, so the cloud side
// keeps its current value for them.
func expandPasswordRules(ctx context.Context, value types.Object, diags *diag.Diagnostics) *cam.PassWordRule {
	rule := &cam.PassWordRule{}
	if value.IsNull() || value.IsUnknown() {
		return rule
	}

	var model CamUpdatePasswordRulesConfigRuleModel
	diags.Append(value.As(ctx, &model, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	if !model.MinimumLength.IsNull() && !model.MinimumLength.IsUnknown() {
		rule.MinimumLength = int64Pointer(model.MinimumLength.ValueInt64())
	}
	if !model.MustContain.IsNull() && !model.MustContain.IsUnknown() {
		rule.MustContain = stringPointer(model.MustContain.ValueString())
	}
	if !model.ForcePasswordChange.IsNull() && !model.ForcePasswordChange.IsUnknown() {
		rule.ForcePasswordChange = int64Pointer(model.ForcePasswordChange.ValueInt64())
	}
	if !model.ReusePasswordLimit.IsNull() && !model.ReusePasswordLimit.IsUnknown() {
		rule.ReusePasswordLimit = int64Pointer(model.ReusePasswordLimit.ValueInt64())
	}
	if !model.RetryPasswordLimit.IsNull() && !model.RetryPasswordLimit.IsUnknown() {
		rule.RetryPasswordLimit = int64Pointer(model.RetryPasswordLimit.ValueInt64())
	}
	if !model.PasswordExpirationInvalidation.IsNull() && !model.PasswordExpirationInvalidation.IsUnknown() {
		rule.PasswordExpirationInvalidation = int64Pointer(model.PasswordExpirationInvalidation.ValueInt64())
	}
	return rule
}

// flattenPasswordRules converts the cloud-side rules into the rules object.
// The cloud may omit a field; the previously configured value is then kept so an
// Optional+Computed attribute never flips to null after apply.
func flattenPasswordRules(ctx context.Context, rules *cam.PassWordRule, prev types.Object, diags *diag.Diagnostics) types.Object {
	var previous CamUpdatePasswordRulesConfigRuleModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.As(ctx, &previous, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return types.ObjectNull(passwordRulesAttrTypes())
		}
	} else {
		previous = CamUpdatePasswordRulesConfigRuleModel{
			MinimumLength:                  types.Int64Null(),
			MustContain:                    types.StringNull(),
			ForcePasswordChange:            types.Int64Null(),
			ReusePasswordLimit:             types.Int64Null(),
			RetryPasswordLimit:             types.Int64Null(),
			PasswordExpirationInvalidation: types.Int64Null(),
		}
	}

	current := CamUpdatePasswordRulesConfigRuleModel{
		MinimumLength:                  types.Int64Null(),
		MustContain:                    types.StringNull(),
		ForcePasswordChange:            types.Int64Null(),
		ReusePasswordLimit:             types.Int64Null(),
		RetryPasswordLimit:             types.Int64Null(),
		PasswordExpirationInvalidation: types.Int64Null(),
	}
	if rules != nil {
		current.MinimumLength = int64PointerValue(rules.MinimumLength)
		current.MustContain = stringPointerValue(rules.MustContain)
		current.ForcePasswordChange = int64PointerValue(rules.ForcePasswordChange)
		current.ReusePasswordLimit = int64PointerValue(rules.ReusePasswordLimit)
		current.RetryPasswordLimit = int64PointerValue(rules.RetryPasswordLimit)
		current.PasswordExpirationInvalidation = int64PointerValue(rules.PasswordExpirationInvalidation)
	}

	merged := CamUpdatePasswordRulesConfigRuleModel{
		MinimumLength:                  preferKnownInt64(current.MinimumLength, previous.MinimumLength),
		MustContain:                    preferKnownString(current.MustContain, previous.MustContain),
		ForcePasswordChange:            preferKnownInt64(current.ForcePasswordChange, previous.ForcePasswordChange),
		ReusePasswordLimit:             preferKnownInt64(current.ReusePasswordLimit, previous.ReusePasswordLimit),
		RetryPasswordLimit:             preferKnownInt64(current.RetryPasswordLimit, previous.RetryPasswordLimit),
		PasswordExpirationInvalidation: preferKnownInt64(current.PasswordExpirationInvalidation, previous.PasswordExpirationInvalidation),
	}

	result, objectDiags := types.ObjectValueFrom(ctx, passwordRulesAttrTypes(), merged)
	diags.Append(objectDiags...)
	return result
}

// preferKnownInt64 keeps the cloud value when it is known and falls back to the
// previously configured value otherwise.
func preferKnownInt64(cloud, previous types.Int64) types.Int64 {
	if cloud.IsNull() || cloud.IsUnknown() {
		return previous
	}
	return cloud
}

// preferKnownString keeps the cloud value when it is known and falls back to the
// previously configured value otherwise.
func preferKnownString(cloud, previous types.String) types.String {
	if cloud.IsNull() || cloud.IsUnknown() {
		return previous
	}
	return cloud
}

// int64Pointer returns a pointer to a copy of v.
func int64Pointer(v int64) *int64 {
	return &v
}

// stringPointer returns a pointer to a copy of v.
func stringPointer(v string) *string {
	return &v
}

// int64PointerValue converts a nullable int64 into types.Int64.
func int64PointerValue(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}

// stringPointerValue converts a nullable string into types.String.
func stringPointerValue(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

// stringOrEmpty dereferences an optional string for diagnostics.
func stringOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

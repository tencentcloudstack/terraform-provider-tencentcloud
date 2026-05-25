// Package cvmactions provides framework-side action implementations for the
// CVM product.
//
// Currently includes:
//   - tencentcloud_reboot_instance: a stub action that only validates input
//     parameters and emits log entries; it does **not** call the CVM
//     RebootInstances API.
//
// Wiring: append cvmactions.NewRebootInstanceAction to frameworkActions() in
// tencentcloud/framework/registry.go to expose the action through the
// framework provider.
package cvmactions

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// Compile-time assertions: rebootInstanceAction satisfies both action.Action
// and action.ActionWithConfigure.
var (
	_ action.Action              = &rebootInstanceAction{}
	_ action.ActionWithConfigure = &rebootInstanceAction{}
)

// instanceIDPattern validates that a CVM instance id has the form ins-xxx,
// matching the InstanceId format constraint documented in the TencentCloud
// CVM API reference.
var instanceIDPattern = regexp.MustCompile(`^ins-[a-z0-9]+$`)

// NewRebootInstanceAction is the framework action factory.
func NewRebootInstanceAction() action.Action {
	return &rebootInstanceAction{}
}

// rebootInstanceAction is a stub action that demonstrates the framework
// Action type as a reference implementation in this repository, covering
// the full Metadata / Schema / Configure / Invoke lifecycle.
//
// Behavioural contract:
//   - Only validates the format of the instance_id input (regex
//     `^ins-[a-z0-9]+$`); a validation failure adds an Error diagnostic.
//   - On a successful validation, emits a tflog.Info entry; it does **not**
//     call any CVM API.
//   - Failures may only originate from input validation — never from
//     network/SDK errors.
type rebootInstanceAction struct {
	// client is fetched from *sharedmeta.ProviderMeta during Configure.
	// This stub action does not actually use it; the field is kept solely
	// to demonstrate "how to obtain the shared client".
	client interface{}
}

// rebootInstanceModel is the Go-side counterpart of the action schema.
type rebootInstanceModel struct {
	InstanceID types.String `tfsdk:"instance_id"`
	Force      types.Bool   `tfsdk:"force"`
}

// Metadata sets the action's Terraform type name to
// tencentcloud_reboot_instance.
//
// Note: the framework Action interface lists Schema before Metadata, but
// neither method depends on the other here, so we follow the
// Resource/DataSource convention and place Metadata first for readability.
func (a *rebootInstanceAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reboot_instance"
}

// Schema defines the action's attribute set.
//
// Fields:
//   - instance_id: required string
//   - force:       optional bool
func (a *rebootInstanceAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		Description: "Reference Action that **does NOT call CVM RebootInstances API**. " +
			"Validates instance_id format and emits a tflog.Info entry, then returns success. " +
			"Useful as a template for implementing real CVM action types in this provider.",
		Attributes: map[string]actionschema.Attribute{
			"instance_id": actionschema.StringAttribute{
				Required:    true,
				Description: "CVM instance id to reboot, must match `^ins-[a-z0-9]+$`.",
			},
			"force": actionschema.BoolAttribute{
				Optional:    true,
				Description: "When true, indicates a forced reboot would be requested. Stub implementation only logs this flag.",
			},
		},
	}
}

// Configure retrieves the shared client that the framework provider wrote
// into ActionData during its own Configure phase. The framework requires
// provider.ConfigureResponse to set ActionData for
// action.ConfigureRequest.ProviderData to be non-nil.
func (a *rebootInstanceAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			"Expected *sharedmeta.ProviderMeta, please report this issue to the provider maintainers.",
		)
		return
	}
	a.client = meta.Client
}

// Invoke runs the action logic. Note: the framework v1.19 Action interface
// names this method `Invoke` (not `Run`, despite earlier propose/spec
// drafts).
//
// This stub merely:
//  1. Validates instance_id with a regex;
//  2. Emits a tflog.Info "would-be reboot request" entry;
//  3. Never returns a network or SDK error.
func (a *rebootInstanceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var plan rebootInstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.InstanceID.ValueString()
	if !instanceIDPattern.MatchString(id) {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_id"),
			"Invalid Instance ID Format",
			"instance_id must match the regex `^ins-[a-z0-9]+$`. Got: "+id,
		)
		return
	}

	tflog.Info(ctx, "reference action stub: would reboot instance", map[string]any{
		"instance_id": id,
		"force":       plan.Force.ValueBool(),
	})
}

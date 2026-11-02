package tco

import (
	"context"
	"fmt"
	"time"

	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ action.Action is intentionally omitted: the factory function
// NewTcoProvisionRoleConfiguration returns action.Action, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &TcoProvisionRoleConfiguration{}

// NewTcoProvisionRoleConfiguration is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewTcoProvisionRoleConfiguration() action.Action {
	return &TcoProvisionRoleConfiguration{}
}

// defaultProvisionRoleConfigurationInvokeTimeout is the built-in maximum
// duration to wait for the asynchronous provisioning task to finish.
// Practitioners can override it through the `timeouts` block, e.g.
// `timeouts { invoke = "30m" }`.
const defaultProvisionRoleConfigurationInvokeTimeout = 10 * time.Minute

// TcoProvisionRoleConfiguration implements action.Action for
// tencentcloud_organization_provision_role_configuration.
type TcoProvisionRoleConfiguration struct {
	fw.ActionWithConfigure
}

// TcoProvisionRoleConfigurationModel maps the schema attributes.
type TcoProvisionRoleConfigurationModel struct {
	ZoneId              types.String `tfsdk:"zone_id"`
	RoleConfigurationId types.String `tfsdk:"role_configuration_id"`
	TargetType          types.String `tfsdk:"target_type"`
	TargetUin           types.Int64  `tfsdk:"target_uin"`

	// Timeouts is the value type defined by the official
	// terraform-plugin-framework-timeouts module. Only the `invoke` stage
	// exists because an action has a single lifecycle stage.
	Timeouts actiontimeouts.Value `tfsdk:"timeouts"`
}

func (a *TcoProvisionRoleConfiguration) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_organization_provision_role_configuration"
}

func (a *TcoProvisionRoleConfiguration) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to provision (deploy) a permission configuration to a target account of the " +
			"Tencent Cloud Organization via the ProvisionRoleConfiguration API. The action waits for the asynchronous " +
			"task to finish. This is a one-time operation; no cloud-side state is persisted after the action completes.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Space ID.",
			},
			"role_configuration_id": schema.StringAttribute{
				Required:    true,
				Description: "Permission configuration ID.",
			},
			"target_type": schema.StringAttribute{
				Required: true,
				Description: "Type of the synchronized target account of the Tencent Cloud Organization. " +
					"Valid values: `ManagerUin` (management account), `MemberUin` (member account).",
			},
			"target_uin": schema.Int64Attribute{
				Required:    true,
				Description: "UIN of the target account of the Tencent Cloud Organization.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": tcoProvisionRoleConfigurationTimeoutsBlock(ctx),
		},
	}
}

// tcoProvisionRoleConfigurationTimeoutsBlock builds the `timeouts` block from
// the official terraform-plugin-framework-timeouts module, whose single
// optional `invoke` attribute is validated as a time.Duration at plan time.
// That module leaves the block itself undescribed (only the attribute carries
// a description), so the block description is restored here to keep the schema
// - and the documentation generated from it by gendoc - self-explanatory.
func tcoProvisionRoleConfigurationTimeoutsBlock(ctx context.Context) schema.Block {
	block := actiontimeouts.BlockWithOpts(ctx, actiontimeouts.Opts{
		InvokeDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `10m`. " +
			"It limits how long the action waits for the asynchronous provisioning task to complete. Default is `10m`.",
	})

	if nested, ok := block.(schema.SingleNestedBlock); ok {
		nested.Description = "The timeouts block allows you to specify the timeout for the invoke operation."
		return nested
	}

	return block
}

// Invoke is called to run the logic of the action.
func (a *TcoProvisionRoleConfiguration) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TcoProvisionRoleConfigurationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ZoneId.IsNull() || data.ZoneId.IsUnknown() || data.ZoneId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing zone_id",
			"zone_id is required and cannot be empty.",
		)
		return
	}

	if data.RoleConfigurationId.IsNull() || data.RoleConfigurationId.IsUnknown() || data.RoleConfigurationId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing role_configuration_id",
			"role_configuration_id is required and cannot be empty.",
		)
		return
	}

	if data.TargetType.IsNull() || data.TargetType.IsUnknown() || data.TargetType.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing target_type",
			"target_type is required and cannot be empty.",
		)
		return
	}

	if data.TargetUin.IsNull() || data.TargetUin.IsUnknown() || data.TargetUin.ValueInt64() == 0 {
		resp.Diagnostics.AddError(
			"Missing target_uin",
			"target_uin is required and cannot be empty.",
		)
		return
	}

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	// An unset/null timeouts block falls back to the built-in default, and a
	// malformed duration is reported as an error diagnostic. The schema
	// validator normally rejects malformed values at plan time; this is the
	// runtime safety net.
	timeout, timeoutDiags := data.Timeouts.Invoke(ctx, defaultProvisionRoleConfigurationInvokeTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := data.ZoneId.ValueString()
	roleConfigurationId := data.RoleConfigurationId.ValueString()
	targetType := data.TargetType.ValueString()
	targetUin := data.TargetUin.ValueInt64()

	service := NewOrganizationService(a.Client())
	taskId, err := service.ProvisionRoleConfiguration(ctx, zoneId, roleConfigurationId, targetType, targetUin)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("provisioning role configuration (%s) to target account (%d)", roleConfigurationId, targetUin),
			err.Error(),
		)
		return
	}

	if err := service.GetOrganizationTaskStatusById(ctx, zoneId, taskId, timeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("waiting for the role configuration provisioning task (%s) to complete", taskId),
			err.Error(),
		)
		return
	}
}

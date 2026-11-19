package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ action.Action is intentionally omitted: the factory function
// NewTeoPlanForZoneAction returns action.Action, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &TeoPlanForZoneAction{}

// NewTeoPlanForZoneAction is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewTeoPlanForZoneAction() action.Action {
	return &TeoPlanForZoneAction{}
}

// TeoPlanForZoneAction implements action.Action for
// tencentcloud_teo_plan_for_zone.
type TeoPlanForZoneAction struct {
	fw.ActionWithConfigure
}

// TeoPlanForZoneActionModel maps the action schema attributes.
type TeoPlanForZoneActionModel struct {
	ZoneId   types.String `tfsdk:"zone_id"`
	PlanType types.String `tfsdk:"plan_type"`
}

func (a *TeoPlanForZoneAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_plan_for_zone"
}

func (a *TeoPlanForZoneAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Purchase an EdgeOne (TEO) plan for a zone that has not yet bound a plan via the CreatePlanForZone API. " +
			"This is a one-time operation: the plan purchase is irreversible and cannot be cancelled.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Zone ID.",
			},
			"plan_type": schema.StringAttribute{
				Required:    true,
				Description: "Plan type to purchase. Valid values: `sta`, `sta_with_bot`, `sta_cm`, `sta_cm_with_bot`, `sta_global`, `sta_global_with_bot`, `ent`, `ent_with_bot`, `ent_cm`, `ent_cm_with_bot`, `ent_global`, `ent_global_with_bot`.",
			},
		},
	}
}

// Invoke is called to run the logic of the action.
func (a *TeoPlanForZoneAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TeoPlanForZoneActionModel
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

	if data.PlanType.IsNull() || data.PlanType.IsUnknown() || data.PlanType.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing plan_type",
			"plan_type is required and cannot be empty.",
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

	zoneId := data.ZoneId.ValueString()
	planType := data.PlanType.ValueString()
	service := NewTeoService(a.Client())
	if _, _, err := service.TeoPlanForZone(zoneId, planType); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("purchasing TEO plan for zone (%s)", zoneId),
			err.Error(),
		)
		return
	}
}

package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var _ action.Action = &TeoPlanForZoneAction{}
var _ action.ActionWithConfigure = &TeoPlanForZoneAction{}

// NewTeoPlanForZoneAction is the factory referenced by
// tencentcloud/framework/registry.go to register this framework action.
func NewTeoPlanForZoneAction() action.Action {
	return &TeoPlanForZoneAction{}
}

// TeoPlanForZoneAction implements action.Action for purchasing a TEO plan for
// a zone via the CreatePlanForZone API. It is a one-time operation: it does not
// manage any long-lived resource state and produces side-effects only; the
// response's ResourceNames and DealNames are logged for diagnostics but not
// returned to the caller (framework actions cannot surface output values).
type TeoPlanForZoneAction struct {
	client TeoService
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
	resp.Schema = actionschema.Schema{
		Description: "Purchase an EdgeOne (TEO) plan for a zone that has not yet bound a plan via the CreatePlanForZone API. This is a one-time operation action; the plan purchase is irreversible and cannot be cancelled.",
		Attributes: map[string]actionschema.Attribute{
			"zone_id": actionschema.StringAttribute{
				Required:    true,
				Description: "Zone ID.",
			},
			"plan_type": actionschema.StringAttribute{
				Required:    true,
				Description: "Plan type to purchase. Valid values: `sta`, `sta_with_bot`, `sta_cm`, `sta_cm_with_bot`, `sta_global`, `sta_global_with_bot`, `ent`, `ent_with_bot`, `ent_cm`, `ent_cm_with_bot`, `ent_global`, `ent_global_with_bot`.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure (which itself reuses the SDKv2-built
// *connectivity.TencentCloudClient). We then construct the SDKv2-style
// TeoService directly from the shared client, mirroring how SDKv2 resources
// in this package construct it.
func (a *TeoPlanForZoneAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData type",
			fmt.Sprintf("Expected *sharedmeta.ProviderMeta, got: %T. "+
				"This is a bug in the provider; please report it.", req.ProviderData),
		)
		return
	}

	a.client = NewTeoService(meta.Client)
}

// Invoke runs the CreatePlanForZone API call to purchase a plan for the zone.
// The response's ResourceNames and DealNames are logged for diagnostics.
func (a *TeoPlanForZoneAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TeoPlanForZoneActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := data.ZoneId.ValueString()
	if zoneId == "" {
		resp.Diagnostics.AddError("Missing zone_id", "zone_id is required and cannot be empty.")
		return
	}

	planType := data.PlanType.ValueString()
	if planType == "" {
		resp.Diagnostics.AddError("Missing plan_type", "plan_type is required and cannot be empty.")
		return
	}

	if a.client.client == nil {
		resp.Diagnostics.AddError(
			"TEO client not configured",
			"The TeoPlanForZone action was invoked before the provider client was configured. "+
				"This is a bug in the provider; please report it.",
		)
		return
	}

	resourceNames, dealNames, err := a.client.TeoPlanForZone(zoneId, planType)
	if err != nil {
		resp.Diagnostics.AddError(
			"CreatePlanForZone failed",
			fmt.Sprintf("Failed to purchase plan for zone %q: %v", zoneId, err),
		)
		return
	}

	if resourceNames == nil && dealNames == nil {
		resp.Diagnostics.AddError(
			"CreatePlanForZone returned empty result",
			fmt.Sprintf("CreatePlanForZone returned an empty result for zone %q.", zoneId),
		)
		return
	}

	log.Printf("[DEBUG] teo plan_for_zone invoke success, zone_id=%s, resource_names: %+v, deal_names: %+v",
		zoneId, resourceNames, dealNames)
}

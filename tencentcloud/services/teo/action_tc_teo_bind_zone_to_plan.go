package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var (
	_ action.Action              = &TeoBindZoneToPlanAction{}
	_ action.ActionWithConfigure = &TeoBindZoneToPlanAction{}
)

// NewTeoBindZoneToPlanAction is the factory referenced by
// tencentcloud/framework/registry.go to register this framework action.
func NewTeoBindZoneToPlanAction() action.Action {
	return &TeoBindZoneToPlanAction{}
}

// TeoBindZoneToPlanAction implements action.Action for
// tencentcloud_teo_bind_zone_to_plan. It is a one-shot operation that calls
// the TEO BindZoneToPlan API to bind an unbound zone to an existing plan.
// Actions do not manage state, so there is no Read/Update/Delete lifecycle.
type TeoBindZoneToPlanAction struct {
	Client *connectivity.TencentCloudClient
}

// TeoBindZoneToPlanActionModel maps the action schema attributes.
type TeoBindZoneToPlanActionModel struct {
	ZoneId types.String `tfsdk:"zone_id"`
	PlanId types.String `tfsdk:"plan_id"`
}

// Metadata returns the action type name.
func (a *TeoBindZoneToPlanAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_bind_zone_to_plan"
}

// Schema returns the action schema. Action attributes cannot be Computed or
// Sensitive; both zone_id and plan_id are Required inputs.
func (a *TeoBindZoneToPlanAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		Description: "Provides an action to bind a TEO (EdgeOne) zone to an existing plan. " +
			"After creating an EdgeOne zone, you can invoke this action to bind the unbound zone to an existing plan so that the zone takes effect.",
		Attributes: map[string]actionschema.Attribute{
			"zone_id": actionschema.StringAttribute{
				Required:    true,
				Description: "ID of the unbound zone that needs to be bound to a plan.",
			},
			"plan_id": actionschema.StringAttribute{
				Required:    true,
				Description: "ID of the target plan to be bound.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by the framework
// provider and stores the shared SDK client for use in Invoke.
func (a *TeoBindZoneToPlanAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
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

	a.Client = meta.Client
}

// Invoke runs the action. It builds a BindZoneToPlanRequest from the supplied
// config, calls BindZoneToPlan with retry handling, and reports diagnostics.
func (a *TeoBindZoneToPlanAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config TeoBindZoneToPlanActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := config.ZoneId.ValueString()
	planId := config.PlanId.ValueString()

	if zoneId == "" {
		resp.Diagnostics.AddError("Missing zone_id", "zone_id is required and cannot be empty.")
		return
	}
	if planId == "" {
		resp.Diagnostics.AddError("Missing plan_id", "plan_id is required and cannot be empty.")
		return
	}

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := teov20220901.NewBindZoneToPlanRequest()
	request.ZoneId = helper.String(zoneId)
	request.PlanId = helper.String(planId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := a.Client.UseTeoV20220901Client().BindZoneToPlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s bind zone to plan failed, reason:%+v", logId, err)
		resp.Diagnostics.AddError(
			"Error invoking teo_bind_zone_to_plan action",
			fmt.Sprintf("BindZoneToPlan failed for zone_id %q plan_id %q: %v", zoneId, planId, err),
		)
		return
	}
}

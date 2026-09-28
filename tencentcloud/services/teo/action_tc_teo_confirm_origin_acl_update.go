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
// NewTeoConfirmOriginAclUpdate returns action.Action, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &TeoConfirmOriginAclUpdate{}

// NewTeoConfirmOriginAclUpdate is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewTeoConfirmOriginAclUpdate() action.Action {
	return &TeoConfirmOriginAclUpdate{}
}

// TeoConfirmOriginAclUpdate implements action.Action for
// tencentcloud_teo_confirm_origin_acl_update_action.
type TeoConfirmOriginAclUpdate struct {
	fw.ActionWithConfigure
}

// TeoConfirmOriginAclUpdateModel maps the schema attributes.
type TeoConfirmOriginAclUpdateModel struct {
	ZoneId types.String `tfsdk:"zone_id"`
}

func (a *TeoConfirmOriginAclUpdate) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_confirm_origin_acl_update"
}

func (a *TeoConfirmOriginAclUpdate) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to confirm TEO origin ACL update for a zone. " +
			"When the origin IP ranges of TEO change, you can use this action to confirm that the latest " +
			"origin IP ranges have been updated to the origin firewall, and the change notification will stop being pushed.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Zone ID.",
			},
		},
	}
}

// Invoke is called to run the logic of the action.
func (a *TeoConfirmOriginAclUpdate) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TeoConfirmOriginAclUpdateModel
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

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	zoneId := data.ZoneId.ValueString()
	service := NewTeoService(a.Client())
	if err := service.ConfirmOriginACLUpdate(ctx, zoneId); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("confirming TEO origin ACL update for zone (%s)", zoneId),
			err.Error(),
		)
		return
	}
}

package antiddos

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ action.ActionWithConfigure = &AntiddosUnblockResources{}

// NewAntiddosUnblockResources is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewAntiddosUnblockResources() action.Action {
	return &AntiddosUnblockResources{}
}

// AntiddosUnblockResources implements action.Action for
// tencentcloud_antiddos_unblock_resources.
type AntiddosUnblockResources struct {
	fw.ActionWithConfigure
}

// AntiddosUnblockResourcesModel maps the schema attributes.
type AntiddosUnblockResourcesModel struct {
	Resources types.List `tfsdk:"resources"`
}

func (a *AntiddosUnblockResources) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_antiddos_unblock_resources"
}

func (a *AntiddosUnblockResources) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to apply for unblocking blocked AntiDDoS (DDoS protection) resources " +
			"(public IP list) via the UnblockResources API. This is a one-time operation; no cloud-side state is persisted.",
		Attributes: map[string]schema.Attribute{
			"resources": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of public IPs to unblock. The list length is limited to 10.",
			},
		},
	}
}

// Invoke is called to run the logic of the action.
func (a *AntiddosUnblockResources) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data AntiddosUnblockResourcesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var resources []*string
	if !data.Resources.IsNull() && !data.Resources.IsUnknown() {
		elems := make([]types.String, 0, len(data.Resources.Elements()))
		resp.Diagnostics.Append(data.Resources.ElementsAs(ctx, &elems, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range elems {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			ip := v.ValueString()
			resources = append(resources, &ip)
		}
	}
	if len(resources) == 0 {
		resp.Diagnostics.AddError(
			"Missing resources",
			"resources is required and cannot be empty.",
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

	service := NewAntiddosService(a.Client())
	if err := service.UnblockResources(ctx, resources); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("unblocking AntiDDoS resources"),
			err.Error(),
		)
		return
	}
}

package bdrc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ action.Action is intentionally omitted: the factory function
// NewBdrcCopyPairTasks returns action.Action, which already performs the
// implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &BdrcCopyPairTasks{}

// NewBdrcCopyPairTasks is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewBdrcCopyPairTasks() action.Action {
	return &BdrcCopyPairTasks{}
}

// BdrcCopyPairTasks implements action.Action for
// tencentcloud_bdrc_copy_pair_tasks.
type BdrcCopyPairTasks struct {
	fw.ActionWithConfigure
}

// BdrcCopyPairTasksModel maps the schema attributes.
type BdrcCopyPairTasksModel struct {
	CopyPairIds  types.List   `tfsdk:"copy_pair_ids"`
	CopyPairType types.String `tfsdk:"copy_pair_type"`
}

func (a *BdrcCopyPairTasks) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_bdrc_copy_pair_tasks"
}

func (a *BdrcCopyPairTasks) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to launch a group of BDRC (Business Data Resilient Disaster Recovery) copy pair tasks " +
			"via the RunCopyPairTasks API. This is a one-time operation; no cloud-side state is persisted.",
		Attributes: map[string]schema.Attribute{
			"copy_pair_ids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of copy pair IDs to launch.",
			},
			"copy_pair_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of the copy pairs to launch. Valid values: DISK, INSTANCE, CFS.",
			},
		},
	}
}

// Invoke is called to run the logic of the action.
func (a *BdrcCopyPairTasks) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data BdrcCopyPairTasksModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.CopyPairType.IsNull() || data.CopyPairType.IsUnknown() || data.CopyPairType.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing copy_pair_type",
			"copy_pair_type is required and cannot be empty.",
		)
		return
	}

	copyPairType := data.CopyPairType.ValueString()

	var copyPairIds []*string
	if !data.CopyPairIds.IsNull() && !data.CopyPairIds.IsUnknown() {
		elems := make([]types.String, 0, len(data.CopyPairIds.Elements()))
		resp.Diagnostics.Append(data.CopyPairIds.ElementsAs(ctx, &elems, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range elems {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			id := v.ValueString()
			copyPairIds = append(copyPairIds, &id)
		}
	}
	if len(copyPairIds) == 0 {
		resp.Diagnostics.AddError(
			"Missing copy_pair_ids",
			"copy_pair_ids is required and cannot be empty.",
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

	service := NewBdrcService(a.Client())
	if err := service.RunCopyPairTasks(ctx, copyPairIds, copyPairType); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("launching BDRC copy pair tasks with type (%s)", copyPairType),
			err.Error(),
		)
		return
	}
}

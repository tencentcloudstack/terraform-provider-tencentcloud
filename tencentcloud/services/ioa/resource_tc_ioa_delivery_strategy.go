package ioa

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	ioav20220601 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ioa/v20220601"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions (repo convention: only assert the extension
// interfaces; the base interface is checked implicitly by the factory return value).
var (
	_ resource.ResourceWithConfigure   = &IoaDeliveryStrategyResource{}
	_ resource.ResourceWithImportState = &IoaDeliveryStrategyResource{}
)

// NewIoaDeliveryStrategyResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewIoaDeliveryStrategyResource() resource.Resource {
	return &IoaDeliveryStrategyResource{}
}

// IoaDeliveryStrategyResource implements resource.Resource for
// tencentcloud_ioa_delivery_strategy.
type IoaDeliveryStrategyResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// IoaDeliveryStrategyResourceModel maps the schema attributes.
type IoaDeliveryStrategyResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	LogTypes    types.Set    `tfsdk:"log_types"`
	Type        types.Int64  `tfsdk:"type"`
	ConfigId    types.Int64  `tfsdk:"config_id"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

func (r *IoaDeliveryStrategyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ioa_delivery_strategy"
}

func (r *IoaDeliveryStrategyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create or edit an IOA log delivery strategy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the delivery strategy, from the ID returned by ModifyDeliveryStrategy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Delivery strategy name.",
			},
			"log_types": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "List of forwarded log types. Valid values: identity security `040201000000`, resource access log `040203000000`, client operation log `040202000000`, terminal security log `040205000000`, terminal control log `040206000000`, terminal compliance log `040207000000`, license authorization log `040212000000`.",
			},
			"type": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Delivery method. Valid values: `0` syslog, `2` kafka, `3` http (32-bit only).",
			},
			"config_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Delivery server id (32-bit only).",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Delivery strategy description.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the strategy is enabled. Note: DescribeDeliveryStrategy does not return this field, so drift detection is not available and the state value is kept on refresh.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
//
// ModifyDeliveryStrategy is a create+update combined API: on create the ID
// field is omitted and the API returns the new strategy id in Response.Data.ID.
func (r *IoaDeliveryStrategyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IoaDeliveryStrategyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := ioav20220601.NewModifyDeliveryStrategyRequest()
	request.Name = plan.Name.ValueStringPointer()
	request.LogTypes = ioaDeliveryStrategyLogTypesFromSet(plan.LogTypes)
	request.Type = plan.Type.ValueInt64Pointer()
	request.ConfigId = plan.ConfigId.ValueInt64Pointer()
	request.Description = plan.Description.ValueStringPointer()
	request.Enabled = plan.Enabled.ValueBoolPointer()

	var response *ioav20220601.ModifyDeliveryStrategyResponse
	retryErr := sdkresource.Retry(tccommon.WriteRetryTimeout, func() *sdkresource.RetryError {
		result, e := r.Client().UseIoaV20220601Client().ModifyDeliveryStrategy(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if retryErr != nil {
		resp.Diagnostics.AddError(
			"Error creating ioa delivery strategy",
			fmt.Sprintf("ModifyDeliveryStrategy failed: %v", retryErr),
		)
		return
	}

	// A nil response (or a nil id field) is an API failure, not an empty result.
	// Never write an empty id into state: the resource would become unmanageable.
	if response == nil || response.Response == nil || response.Response.Data == nil || response.Response.Data.ID == nil || *response.Response.Data.ID == 0 {
		resp.Diagnostics.AddError(
			"Error creating ioa delivery strategy",
			"ModifyDeliveryStrategy returned an empty response or an empty strategy id.",
		)
		return
	}

	// Persist the identity BEFORE any read back.
	plan.Id = types.StringValue(strconv.FormatUint(*response.Response.Data.ID, 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ioa delivery strategy after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"ioa delivery strategy not readable after create",
			fmt.Sprintf("strategy %s was created but cannot be read back", plan.Id.ValueString()),
		)
		return
	}
	r.refreshModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *IoaDeliveryStrategyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IoaDeliveryStrategyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ioa delivery strategy", err.Error())
		return
	}

	// Confirmed "not found" on the cloud side => drop it from state so Terraform
	// plans a recreation. A nil response is an API failure and must not land here.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ioa delivery strategy %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	r.refreshModel(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
//
// ModifyDeliveryStrategy is a create+update combined API: on update the ID
// field carries the strategy id from state. Null/unknown attributes are not
// sent, so the cloud keeps its previous value; removing an optional attribute
// from the configuration therefore does not clear it on the cloud side.
func (r *IoaDeliveryStrategyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IoaDeliveryStrategyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IoaDeliveryStrategyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := ioav20220601.NewModifyDeliveryStrategyRequest()
	request.ID = helper.StrToInt64Point(state.Id.ValueString())
	request.Name = plan.Name.ValueStringPointer()
	request.LogTypes = ioaDeliveryStrategyLogTypesFromSet(plan.LogTypes)
	request.Type = plan.Type.ValueInt64Pointer()
	request.ConfigId = plan.ConfigId.ValueInt64Pointer()
	request.Description = plan.Description.ValueStringPointer()
	request.Enabled = plan.Enabled.ValueBoolPointer()

	retryErr := sdkresource.Retry(tccommon.WriteRetryTimeout, func() *sdkresource.RetryError {
		_, e := r.Client().UseIoaV20220601Client().ModifyDeliveryStrategy(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if retryErr != nil {
		resp.Diagnostics.AddError(
			"Error updating ioa delivery strategy",
			fmt.Sprintf("ModifyDeliveryStrategy failed: %v", retryErr),
		)
		return
	}

	item, err := r.read(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ioa delivery strategy after update", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"ioa delivery strategy not readable after update",
			fmt.Sprintf("strategy %s was updated but cannot be read back", state.Id.ValueString()),
		)
		return
	}
	plan.Id = state.Id
	r.refreshModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *IoaDeliveryStrategyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IoaDeliveryStrategyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := ioav20220601.NewDeleteDeliveryStrategyRequest()
	request.ID = helper.StrToUint64Point(state.Id.ValueString())

	retryErr := sdkresource.Retry(tccommon.WriteRetryTimeout, func() *sdkresource.RetryError {
		_, e := r.Client().UseIoaV20220601Client().DeleteDeliveryStrategy(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if retryErr != nil {
		// No error code is documented as "already gone"; everything must fail
		// so the id stays in state for a retry.
		resp.Diagnostics.AddError(
			"Error deleting ioa delivery strategy",
			fmt.Sprintf("DeleteDeliveryStrategy failed: %v", retryErr),
		)
		return
	}
}

// ImportState is embedded via fw.WithImportByID for `terraform import <id>` passthrough.

// read fetches the strategy detail by ID.
//
// Two different "empty" outcomes stay distinguishable:
//   - API failure (nil response / nil Response / nil Data) -> error
//   - queried successfully but nothing matched (Data.Data is null) -> (nil, nil)
func (r *IoaDeliveryStrategyResource) read(ctx context.Context, id string) (*ioav20220601.DeliveryStrategyItem, error) {
	request := ioav20220601.NewDescribeDeliveryStrategyRequest()
	request.ID = helper.StrToUint64Point(id)

	var response *ioav20220601.DescribeDeliveryStrategyResponse
	retryErr := sdkresource.Retry(tccommon.ReadRetryTimeout, func() *sdkresource.RetryError {
		result, e := r.Client().UseIoaV20220601Client().DescribeDeliveryStrategy(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if retryErr != nil {
		return nil, retryErr
	}
	if response == nil || response.Response == nil || response.Response.Data == nil {
		return nil, fmt.Errorf("DescribeDeliveryStrategy returned an empty response for %s", id)
	}
	if response.Response.Data.Data == nil {
		return nil, nil
	}
	return response.Response.Data.Data, nil
}

// refreshModel copies cloud-owned fields into the model.
//
// Enabled has no counterpart in the DescribeDeliveryStrategy response, so its
// state value is kept as-is (drift detection is not possible for this field).
func (r *IoaDeliveryStrategyResource) refreshModel(m *IoaDeliveryStrategyResourceModel, item *ioav20220601.DeliveryStrategyItem) {
	if item == nil {
		return
	}
	if item.Name != nil {
		m.Name = types.StringValue(*item.Name)
	}
	if item.Type != nil {
		m.Type = types.Int64Value(*item.Type)
	}
	m.LogTypes = ioaDeliveryStrategySetFromLogTypes(item.LogTypes)
	if item.Description != nil {
		m.Description = types.StringValue(*item.Description)
	}
	if item.ConfigId != nil {
		m.ConfigId = types.Int64Value(int64(*item.ConfigId))
	}
}

// ioaDeliveryStrategySetFromLogTypes converts the SDK string pointers into a
// Terraform set value; a nil/empty slice becomes a null set.
func ioaDeliveryStrategySetFromLogTypes(values []*string) types.Set {
	if len(values) == 0 {
		return types.SetNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringPointerValue(v))
	}
	set, diags := types.SetValue(types.StringType, elements)
	if diags.HasError() {
		log.Printf("[ERROR] failed to build ioa delivery strategy log_types set: %v", diags.Errors())
		return types.SetNull(types.StringType)
	}
	return set
}

// ioaDeliveryStrategyLogTypesFromSet converts a Terraform set value into the
// SDK []*string; null/unknown sets are sent as nil (field omitted).
func ioaDeliveryStrategyLogTypesFromSet(set types.Set) []*string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	elements := set.Elements()
	values := make([]*string, 0, len(elements))
	for _, element := range elements {
		str, ok := element.(types.String)
		if !ok {
			continue
		}
		values = append(values, str.ValueStringPointer())
	}
	if len(values) == 0 {
		return nil
	}
	return values
}

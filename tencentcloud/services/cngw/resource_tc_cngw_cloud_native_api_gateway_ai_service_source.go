// Resource tencentcloud_cngw_cloud_native_api_gateway_ai_service_source manages
// an AI service source of a cloud native API gateway instance.
package cngw

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	basetypes "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	cngwv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cngw/v20230418"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions (repo convention: only assert the
// extension interfaces; the base interface is checked implicitly by the
// factory return value).
var (
	_ resource.ResourceWithConfigure   = &CngwCloudNativeAPIGatewayAIServiceSourceResource{}
	_ resource.ResourceWithImportState = &CngwCloudNativeAPIGatewayAIServiceSourceResource{}
)

// NewCngwCloudNativeAPIGatewayAIServiceSourceResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewCngwCloudNativeAPIGatewayAIServiceSourceResource() resource.Resource {
	return &CngwCloudNativeAPIGatewayAIServiceSourceResource{}
}

// CngwCloudNativeAPIGatewayAIServiceSourceResource implements resource.Resource for
// tencentcloud_cngw_cloud_native_api_gateway_ai_service_source.
type CngwCloudNativeAPIGatewayAIServiceSourceResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// CngwCloudNativeAPIGatewayAIServiceSourceResourceModel maps the schema attributes;
// each attribute appears exactly once, in the same order as the schema.
type CngwCloudNativeAPIGatewayAIServiceSourceResourceModel struct {
	Id            types.String `tfsdk:"id"`
	GatewayId     types.String `tfsdk:"gateway_id"`
	SourceType    types.String `tfsdk:"source_type"`
	SourceName    types.String `tfsdk:"source_name"`
	SourceProduct types.String `tfsdk:"source_product"`
	SourceInfo    types.Object `tfsdk:"source_info"`
	Description   types.String `tfsdk:"description"`
}

// cngwAIServiceSourceInfoAttributes lists the schema attributes of the
// source_info object attribute.
var cngwAIServiceSourceInfoAttributes = map[string]schema.Attribute{
	"instance_id": schema.StringAttribute{
		Optional:    true,
		Description: "Instance ID of the service source.",
	},
	"auth": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Authorization information.",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Account.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password.",
			},
			"access_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Access token.",
			},
		},
	},
	"addresses": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "List of addresses.",
	},
}

// cngwAIServiceSourceAuthModel maps the nested auth object attribute.
type cngwAIServiceSourceAuthModel struct {
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	AccessToken types.String `tfsdk:"access_token"`
}

// cngwAIServiceSourceInfoSdkModel maps the source_info object attribute for
// the object.As conversion in expand/flatten helpers.
type cngwAIServiceSourceInfoSdkModel struct {
	InstanceId types.String `tfsdk:"instance_id"`
	Auth       types.Object `tfsdk:"auth"`
	Addresses  types.List   `tfsdk:"addresses"`
}

// cngwAIServiceSourceInfoAttrTypes mirrors the attribute types of the
// source_info object attribute; cngwAIServiceSourceAuthAttrTypes mirrors the
// attribute types of the nested auth object attribute.
var (
	cngwAIServiceSourceAuthAttrTypes = map[string]attr.Type{
		"username":     types.StringType,
		"password":     types.StringType,
		"access_token": types.StringType,
	}
	cngwAIServiceSourceInfoAttrTypes = map[string]attr.Type{
		"instance_id": types.StringType,
		"auth":        types.ObjectType{AttrTypes: cngwAIServiceSourceAuthAttrTypes},
		"addresses":   types.ListType{ElemType: types.StringType},
	}
)

func (r *CngwCloudNativeAPIGatewayAIServiceSourceResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cngw_cloud_native_api_gateway_ai_service_source"
}

func (r *CngwCloudNativeAPIGatewayAIServiceSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage an AI service source of a cloud native API gateway.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gateway_id": schema.StringAttribute{
				Required:    true,
				Description: "Instance ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_type": schema.StringAttribute{
				Required:    true,
				Description: "Source type. Valid values: `MCPRegistry`: MCP registry, `Registry`: registry, `DNS`: DNS service. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the service source.",
			},
			"source_product": schema.StringAttribute{
				Optional:    true,
				Description: "Source product. Valid values: `TSFNacos`: TSF Nacos.",
			},
			"source_info": schema.SingleNestedAttribute{
				Optional:    true,
				Attributes:  cngwAIServiceSourceInfoAttributes,
				Description: "Source detail.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *CngwCloudNativeAPIGatewayAIServiceSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CngwCloudNativeAPIGatewayAIServiceSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createAIServiceSourceRequest := cngwv20230418.NewCreateCloudNativeAPIGatewayAIServiceSourceRequest()
	createAIServiceSourceRequest.GatewayId = plan.GatewayId.ValueStringPointer()
	createAIServiceSourceRequest.SourceType = plan.SourceType.ValueStringPointer()
	createAIServiceSourceRequest.SourceName = plan.SourceName.ValueStringPointer()
	if !plan.SourceProduct.IsNull() && !plan.SourceProduct.IsUnknown() {
		createAIServiceSourceRequest.SourceProduct = plan.SourceProduct.ValueStringPointer()
	}
	sourceInfo, sourceInfoDiags := expandCngwAIServiceSourceInfo(ctx, plan.SourceInfo)
	resp.Diagnostics.Append(sourceInfoDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createAIServiceSourceRequest.SourceInfo = sourceInfo
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createAIServiceSourceRequest.Description = plan.Description.ValueStringPointer()
	}

	service := NewCngwService(r.Client())

	sourceId, err := service.CreateCloudNativeAPIGatewayAIServiceSource(ctx, createAIServiceSourceRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cloud native API gateway AI service source",
			fmt.Sprintf("CreateCloudNativeAPIGatewayAIServiceSource failed: %v", err),
		)

		return
	}

	// Never write an empty id into state.
	if sourceId == "" {
		resp.Diagnostics.AddError(
			"Error creating cloud native API gateway AI service source",
			"CreateCloudNativeAPIGatewayAIServiceSource returned an empty source id.",
		)

		return
	}

	plan.Id = types.StringValue(sourceId)

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeAIServiceSourceListRequest := cngwv20230418.NewDescribeCloudNativeAPIGatewayAIServiceSourceListRequest()
	describeAIServiceSourceListRequest.GatewayId = plan.GatewayId.ValueStringPointer()
	describeAIServiceSourceListRequest.Filters = []*cngwv20230418.Filter{
		{
			Name:   fwHelper.String("SourceId"),
			Values: []*string{fwHelper.String(sourceId)},
		},
	}

	sources, err := service.DescribeCloudNativeAPIGatewayAIServiceSourceList(ctx, describeAIServiceSourceListRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cloud native API gateway AI service source after create", err.Error())
		return
	}

	item := findCngwAIServiceSource(sources, sourceId)
	if item == nil {
		resp.Diagnostics.AddError(
			"Cloud native API gateway AI service source not readable after create",
			fmt.Sprintf("AI service source %s was created but cannot be read back", plan.Id.ValueString()),
		)

		return
	}

	resp.Diagnostics.Append(refreshCngwAIServiceSourceModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *CngwCloudNativeAPIGatewayAIServiceSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CngwCloudNativeAPIGatewayAIServiceSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceId := state.Id.ValueString()

	describeAIServiceSourceListRequest := cngwv20230418.NewDescribeCloudNativeAPIGatewayAIServiceSourceListRequest()
	describeAIServiceSourceListRequest.GatewayId = state.GatewayId.ValueStringPointer()
	describeAIServiceSourceListRequest.Filters = []*cngwv20230418.Filter{
		{
			Name:   fwHelper.String("SourceId"),
			Values: []*string{fwHelper.String(sourceId)},
		},
	}

	service := NewCngwService(r.Client())

	sources, err := service.DescribeCloudNativeAPIGatewayAIServiceSourceList(ctx, describeAIServiceSourceListRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cloud native API gateway AI service source", err.Error())
		return
	}

	item := findCngwAIServiceSource(sources, sourceId)

	// Confirmed "not found" on the cloud side => drop it from state so
	// Terraform plans a recreation. The not-found judgement stays INLINE in
	// Read (RemoveResource followed by an immediate return).
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("Cloud native API gateway AI service source %s not found, removing from state.", state.Id.ValueString()),
		)

		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshCngwAIServiceSourceModel(ctx, &state, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
func (r *CngwCloudNativeAPIGatewayAIServiceSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CngwCloudNativeAPIGatewayAIServiceSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CngwCloudNativeAPIGatewayAIServiceSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceId := state.Id.ValueString()

	service := NewCngwService(r.Client())

	// Call a write API ONLY when an attribute it covers changed. A plan that
	// differs from the state in the `timeouts` block alone also reaches
	// Update; sending the write API then is a pointless no-change request
	// that many cloud APIs reject. The read-back below runs either way.
	if !plan.SourceName.Equal(state.SourceName) || !plan.SourceProduct.Equal(state.SourceProduct) || !plan.SourceInfo.Equal(state.SourceInfo) || !plan.Description.Equal(state.Description) {
		modifyAIServiceSourceRequest := cngwv20230418.NewModifyCloudNativeAPIGatewayAIServiceSourceRequest()
		modifyAIServiceSourceRequest.GatewayId = state.GatewayId.ValueStringPointer()
		modifyAIServiceSourceRequest.SourceId = fwHelper.String(sourceId)
		if !plan.SourceName.IsNull() && !plan.SourceName.IsUnknown() {
			modifyAIServiceSourceRequest.SourceName = plan.SourceName.ValueStringPointer()
		}
		if !plan.SourceType.IsNull() && !plan.SourceType.IsUnknown() {
			modifyAIServiceSourceRequest.SourceType = plan.SourceType.ValueStringPointer()
		}
		sourceInfo, sourceInfoDiags := expandCngwAIServiceSourceInfo(ctx, plan.SourceInfo)
		resp.Diagnostics.Append(sourceInfoDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		modifyAIServiceSourceRequest.SourceInfo = sourceInfo
		if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
			modifyAIServiceSourceRequest.Description = plan.Description.ValueStringPointer()
		}

		if err := service.ModifyCloudNativeAPIGatewayAIServiceSource(ctx, modifyAIServiceSourceRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error updating cloud native API gateway AI service source",
				fmt.Sprintf("ModifyCloudNativeAPIGatewayAIServiceSource failed: %v", err),
			)

			return
		}
	}

	describeAIServiceSourceListRequest := cngwv20230418.NewDescribeCloudNativeAPIGatewayAIServiceSourceListRequest()
	describeAIServiceSourceListRequest.GatewayId = state.GatewayId.ValueStringPointer()
	describeAIServiceSourceListRequest.Filters = []*cngwv20230418.Filter{
		{
			Name:   fwHelper.String("SourceId"),
			Values: []*string{fwHelper.String(sourceId)},
		},
	}

	sources, err := service.DescribeCloudNativeAPIGatewayAIServiceSourceList(ctx, describeAIServiceSourceListRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cloud native API gateway AI service source after update", err.Error())
		return
	}

	item := findCngwAIServiceSource(sources, sourceId)
	if item == nil {
		resp.Diagnostics.AddError(
			"Cloud native API gateway AI service source not readable after update",
			fmt.Sprintf("AI service source %s was updated but cannot be read back", state.Id.ValueString()),
		)

		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshCngwAIServiceSourceModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *CngwCloudNativeAPIGatewayAIServiceSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CngwCloudNativeAPIGatewayAIServiceSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceId := state.Id.ValueString()

	deleteAIServiceSourceRequest := cngwv20230418.NewDeleteCloudNativeAPIGatewayAIServiceSourceRequest()
	deleteAIServiceSourceRequest.GatewayId = state.GatewayId.ValueStringPointer()
	deleteAIServiceSourceRequest.SourceId = fwHelper.String(sourceId)
	if !state.SourceType.IsNull() && !state.SourceType.IsUnknown() {
		deleteAIServiceSourceRequest.SourceType = state.SourceType.ValueStringPointer()
	}

	service := NewCngwService(r.Client())

	if err := service.DeleteCloudNativeAPIGatewayAIServiceSource(ctx, deleteAIServiceSourceRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting cloud native API gateway AI service source",
			fmt.Sprintf("DeleteCloudNativeAPIGatewayAIServiceSource failed: %v", err),
		)

		return
	}
}

// refreshCngwAIServiceSourceModel copies cloud-owned fields into the model.
//
// User-configurable attributes keep the previous value where a cloud-omitted
// field would otherwise fail the post-apply consistency check with
// "was <value>, but now null"; the source_info object is refreshed through
// its own flatten helper, which keeps the previous object when the cloud
// omits it entirely.
func refreshCngwAIServiceSourceModel(ctx context.Context, m *CngwCloudNativeAPIGatewayAIServiceSourceResourceModel, item *cngwv20230418.CNAPIGwAIServiceSource) diag.Diagnostics {
	if item == nil {
		return nil
	}

	var diags diag.Diagnostics

	m.SourceName = fwHelper.StringValueOrPrev(item.SourceName, m.SourceName)
	m.SourceProduct = fwHelper.StringValueOrPrev(item.SourceProduct, m.SourceProduct)

	sourceInfo, sourceInfoDiags := flattenCngwAIServiceSourceInfo(ctx, item.SourceInfo, m.SourceInfo)
	diags.Append(sourceInfoDiags...)
	if diags.HasError() {
		return diags
	}
	m.SourceInfo = sourceInfo

	m.Description = fwHelper.StringValueOrPrev(item.Description, m.Description)

	return diags
}

// findCngwAIServiceSource picks this resource's own item from the describe
// list (list-query shape).
func findCngwAIServiceSource(sources []*cngwv20230418.CNAPIGwAIServiceSource, sourceId string) *cngwv20230418.CNAPIGwAIServiceSource {
	for _, s := range sources {
		if s.SourceId != nil && *s.SourceId == sourceId {
			return s
		}
	}

	return nil
}

// expandCngwAIServiceSourceInfo converts the planned source_info object into
// the SDK request struct; null/unknown yields nil ("do not send this field").
func expandCngwAIServiceSourceInfo(ctx context.Context, object types.Object) (*cngwv20230418.CNAPIGwAIServiceSourceInfo, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	// TODO(human): the SDK field types (CNAPIGwAIServiceSourceInfo,
	// CNAPIGwAIServiceSourceAuth) were verified from the vendored SDK; the
	// source_info/auth/addresses shape was NOT verified against a live cloud
	// and the read API is assumed to return the same structure it accepts.
	var sourceInfo cngwAIServiceSourceInfoSdkModel
	diags := object.As(ctx, &sourceInfo, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	info := &cngwv20230418.CNAPIGwAIServiceSourceInfo{}

	if !sourceInfo.InstanceId.IsNull() && !sourceInfo.InstanceId.IsUnknown() {
		info.InstanceId = sourceInfo.InstanceId.ValueStringPointer()
	}

	info.Addresses = fwHelper.PointersFromList[string](ctx, sourceInfo.Addresses, &diags)
	if diags.HasError() {
		return nil, diags
	}

	if !sourceInfo.Auth.IsNull() && !sourceInfo.Auth.IsUnknown() {
		var authModel cngwAIServiceSourceAuthModel
		authDiags := sourceInfo.Auth.As(ctx, &authModel, basetypes.ObjectAsOptions{})
		diags.Append(authDiags...)
		if diags.HasError() {
			return nil, diags
		}

		auth := &cngwv20230418.CNAPIGwAIServiceSourceAuth{}
		auth.Username = fwHelper.StringPointerFromValue(authModel.Username)
		auth.Password = fwHelper.StringPointerFromValue(authModel.Password)
		auth.AccessToken = fwHelper.StringPointerFromValue(authModel.AccessToken)
		info.Auth = auth
	}

	return info, nil
}

// flattenCngwAIServiceSourceInfo converts the cloud-side source info into the
// source_info object attribute. User-configurable attributes keep the previous
// value when the cloud omits them (plain helpers with *OrPrev semantics);
// the whole object stays the previous one when the cloud does not report it.
func flattenCngwAIServiceSourceInfo(ctx context.Context, item *cngwv20230418.CNAPIGwAIServiceSourceInfo, prev types.Object) (types.Object, diag.Diagnostics) {
	if item == nil {
		return prev, nil
	}

	var diags diag.Diagnostics

	// TODO(human): sensitive auth fields are read back verbatim from the
	// cloud answer; confirm with the service owner that the describe API is
	// expected to return them (or mask them) before shipping.
	var model cngwAIServiceSourceInfoSdkModel

	if item.InstanceId == nil {
		model.InstanceId = types.StringNull()
	} else {
		model.InstanceId = types.StringValue(*item.InstanceId)
	}

	model.Addresses = fwHelper.ListFromPointersOrPrev(ctx, types.StringType, item.Addresses, prevAddressList(ctx, prev, &diags), &diags)
	if diags.HasError() {
		return prev, diags
	}

	authObject := types.ObjectNull(cngwAIServiceSourceAuthAttrTypes)
	if item.Auth != nil {
		authModel := cngwAIServiceSourceAuthModel{
			Username:    fwHelper.StringPointerValueOrNull(item.Auth.Username),
			Password:    fwHelper.StringPointerValueOrNull(item.Auth.Password),
			AccessToken: fwHelper.StringPointerValueOrNull(item.Auth.AccessToken),
		}
		authObj, authDiags := types.ObjectValueFrom(ctx, cngwAIServiceSourceAuthAttrTypes, authModel)
		if authDiags.HasError() {
			diags.Append(authDiags...)

			return prev, diags
		}
		authObject = authObj
	}
	model.Auth = authObject

	sourceInfo, sourceInfoDiags := types.ObjectValueFrom(ctx, cngwAIServiceSourceInfoAttrTypes, model)
	diags.Append(sourceInfoDiags...)
	if diags.HasError() {
		return prev, diags
	}

	return sourceInfo, nil
}

// prevAddressList extracts the previous addresses list from the previous
// source_info object (null when absent) for the *OrPrev collection helper.
func prevAddressList(ctx context.Context, prev types.Object, diags *diag.Diagnostics) types.List {
	if prev.IsNull() || prev.IsUnknown() {
		return types.ListNull(types.StringType)
	}

	var model cngwAIServiceSourceInfoSdkModel
	d := prev.As(ctx, &model, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	return model.Addresses
}

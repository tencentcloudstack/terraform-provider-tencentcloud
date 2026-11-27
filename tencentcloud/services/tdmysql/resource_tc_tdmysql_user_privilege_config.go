package tdmysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions. The base resource.Resource interface is
// checked implicitly by the factory return value.
var (
	_ resource.ResourceWithConfigure   = &TdmysqlUserPrivilegeConfigResource{}
	_ resource.ResourceWithImportState = &TdmysqlUserPrivilegeConfigResource{}
)

// NewTdmysqlUserPrivilegeConfigResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTdmysqlUserPrivilegeConfigResource() resource.Resource {
	return &TdmysqlUserPrivilegeConfigResource{}
}

// TdmysqlUserPrivilegeConfigResource implements resource.Resource for
// tencentcloud_tdmysql_user_privilege_config.
type TdmysqlUserPrivilegeConfigResource struct {
	fw.ResourceWithConfigure
}

// TdmysqlUserPrivilegeConfigResourceModel maps every top-level schema attribute; each
// one appears here exactly once.
type TdmysqlUserPrivilegeConfigResourceModel struct {
	Id         types.String `tfsdk:"id"`
	InstanceId types.String `tfsdk:"instance_id"`
	Host       types.String `tfsdk:"host"`
	UserName   types.String `tfsdk:"user_name"`
	DbName     types.String `tfsdk:"db_name"`
	Object     types.String `tfsdk:"object"`
	ObjectType types.String `tfsdk:"object_type"`
	ColName    types.String `tfsdk:"col_name"`
	Privileges types.List   `tfsdk:"privileges"`
}

// Metadata reports the Terraform type name.
func (r *TdmysqlUserPrivilegeConfigResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_tdmysql_user_privilege_config"
}

// Schema declares the user privileges of a TDSQL Boundless instance.
func (r *TdmysqlUserPrivilegeConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage user privileges of a TDSQL Boundless instance. " +
			"Create and update both call the ModifyUserPrivileges API, read calls DescribeUserPrivileges. " +
			"Destroying this resource only removes it from state; the cloud-side privileges live and die with the parent instance and user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the user privilege, in the format `instance_id#host#user_name#db_name`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "Instance ID, such as tdsql3-5baee8df.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host": schema.StringAttribute{
				Required:    true,
				Description: "The access host allowed for the user. The username and host uniquely identify an account.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_name": schema.StringAttribute{
				Required:    true,
				Description: "Login username.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"db_name": schema.StringAttribute{
				Required:    true,
				Description: "Database name. If it is `*`, it represents global privileges (that is, `*.*`), in which case the object_type and object parameters are ignored.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"object": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the specific type, for example the specific table name when object_type is table. When db_name and object_type are both concrete names, object represents the specific object name, which cannot be `*` or empty.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"object_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Type, which can be table, view, proc, func and `*`. When db_name is a concrete database name and object_type is `*`, it represents the privileges of that database (that is, db.*), in which case the object parameter is ignored.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"col_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "When object_type is table, col_name as `*` means querying the privileges of the table, and a specific column name means querying the privileges of the corresponding column.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"privileges": schema.ListAttribute{
				Required:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Privilege list.",
			},
		},
	}
}

// Create writes the privileges for the first time. It shares the write
// path with Update because the cloud API exposes a single write operation.
func (r *TdmysqlUserPrivilegeConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TdmysqlUserPrivilegeConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The state id is the composite locator input — never a fixed constant.
	id := fwHelper.IdFormat(plan.InstanceId.ValueString(), plan.Host.ValueString(), plan.UserName.ValueString(), plan.DbName.ValueString())
	plan.Id = types.StringValue(id)

	service := NewTdmysqlService(r.Client())

	r.modifyTdmysqlUserPrivileges(ctx, service, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	r.refreshTdmysqlUserPrivilege(ctx, service, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the privileges from the cloud.
func (r *TdmysqlUserPrivilegeConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TdmysqlUserPrivilegeConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Split the composite locator back out of the state id.
	parts, err := fwHelper.IdParseErr(state.Id.ValueString(), 4)
	if err != nil {
		resp.Diagnostics.AddError("Invalid TDSQL Boundless user privilege id", fmt.Sprintf("expected instance_id#host#user_name#db_name: %v", err))

		return
	}

	service := NewTdmysqlService(r.Client())

	describeUserPrivilegesRequest := tdmysqlv20211122.NewDescribeUserPrivilegesRequest()
	describeUserPrivilegesRequest.InstanceId = fwHelper.String(parts[0])
	describeUserPrivilegesRequest.Host = fwHelper.String(parts[1])
	describeUserPrivilegesRequest.UserName = fwHelper.String(parts[2])
	describeUserPrivilegesRequest.DbName = fwHelper.String(parts[3])
	describeUserPrivilegesRequest.Object = fwHelper.StringPointerFromValue(state.Object)
	describeUserPrivilegesRequest.ObjectType = fwHelper.StringPointerFromValue(state.ObjectType)
	describeUserPrivilegesRequest.ColName = fwHelper.StringPointerFromValue(state.ColName)

	params, err := service.DescribeTdmysqlUserPrivileges(ctx, describeUserPrivilegesRequest)
	if err != nil {
		// "Privileges gone" => drop from state so Terraform plans a
		// recreation; transport errors still surface as errors.
		if errors.Is(err, fwHelper.ErrNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("TDSQL Boundless user privileges %s not found, removing from state.", state.Id.ValueString()),
			)

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError("Error reading TDSQL Boundless user privileges", err.Error())

		return
	}

	state.InstanceId = types.StringValue(parts[0])
	state.Host = types.StringValue(parts[1])
	state.UserName = types.StringValue(parts[2])
	state.DbName = types.StringValue(parts[3])
	state.Privileges = fwHelper.ListFromPointers(ctx, types.StringType, params.Privileges, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update applies the planned privileges. Shares the write path with Create.
func (r *TdmysqlUserPrivilegeConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TdmysqlUserPrivilegeConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TdmysqlUserPrivilegeConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The locator attributes are RequiresReplace, so the id can never change here.
	plan.Id = state.Id

	service := NewTdmysqlService(r.Client())

	// The write API only covers the privilege list; the read-back always runs.
	if !plan.Privileges.Equal(state.Privileges) {
		r.modifyTdmysqlUserPrivileges(ctx, service, &plan, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	r.refreshTdmysqlUserPrivilege(ctx, service, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is required by the resource.Resource interface. The privileges live
// on the parent instance and user, so destroying this resource only removes
// it from state — the cloud-side privileges die with their parent.
func (r *TdmysqlUserPrivilegeConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// modifyTdmysqlUserPrivileges writes the planned privileges through the write
// API. The privilege group follows the scope carried by the locator: a
// wildcard db_name targets global privileges (*.*), a wildcard object_type on
// a concrete database targets database-level privileges (db.*), and a
// concrete table object targets table-level privileges.
func (r *TdmysqlUserPrivilegeConfigResource) modifyTdmysqlUserPrivileges(ctx context.Context, service TdmysqlService, plan *TdmysqlUserPrivilegeConfigResourceModel, diags *diag.Diagnostics) {
	modifyUserPrivilegesRequest := tdmysqlv20211122.NewModifyUserPrivilegesRequest()
	modifyUserPrivilegesRequest.InstanceId = fwHelper.String(plan.InstanceId.ValueString())
	modifyUserPrivilegesRequest.Users = []*tdmysqlv20211122.User{
		{
			UserName: fwHelper.String(plan.UserName.ValueString()),
			Host:     fwHelper.String(plan.Host.ValueString()),
		},
	}

	privileges := fwHelper.PointersFromList[string](ctx, plan.Privileges, diags)
	if diags.HasError() {
		return
	}

	// sent tracks whether the plan's scope matched a dedicated write group:
	// only a matched group is sent to the cloud.
	sent := false

	switch {
	case plan.DbName.ValueString() == "*":
		// Global privileges (*.*).
		modifyUserPrivilegesRequest.GlobalPrivileges = privileges
		sent = true
	case plan.ObjectType.ValueString() == "*":
		// Database level privileges (db.*).
		modifyUserPrivilegesRequest.DatabasePrivileges = []*tdmysqlv20211122.DatabasePrivileges{
			{
				Database:   fwHelper.String(plan.DbName.ValueString()),
				Privileges: privileges,
			},
		}
		sent = true
	case plan.ObjectType.ValueString() == "table" && plan.Object.ValueString() != "":
		// Table level privileges.
		modifyUserPrivilegesRequest.TablePrivileges = []*tdmysqlv20211122.TablePrivileges{
			{
				Database:   fwHelper.String(plan.DbName.ValueString()),
				Table:      fwHelper.String(plan.Object.ValueString()),
				Privileges: privileges,
			},
		}
		sent = true
	default:
		// TODO(human): scopes without a dedicated write group in
		// ModifyUserPrivileges (object_type view/proc/func, or a table scope
		// without an object name) are skipped conservatively — only the
		// read-back runs. Decide the intended semantics for these scopes.
	}

	if !sent {
		return
	}

	if err := service.ModifyTdmysqlUserPrivileges(ctx, modifyUserPrivilegesRequest); err != nil {
		diags.AddError("Error updating TDSQL Boundless user privileges", err.Error())
	}
}

// refreshTdmysqlUserPrivilege reads the privileges back from the cloud, so
// the Required+Computed privilege list always records the cloud-reported
// values: when the cloud normalizes a write, the difference surfaces in the
// next plan instead of being second-guessed here.
func (r *TdmysqlUserPrivilegeConfigResource) refreshTdmysqlUserPrivilege(ctx context.Context, service TdmysqlService, plan *TdmysqlUserPrivilegeConfigResourceModel, diags *diag.Diagnostics) {
	describeUserPrivilegesRequest := tdmysqlv20211122.NewDescribeUserPrivilegesRequest()
	describeUserPrivilegesRequest.InstanceId = fwHelper.String(plan.InstanceId.ValueString())
	describeUserPrivilegesRequest.Host = fwHelper.String(plan.Host.ValueString())
	describeUserPrivilegesRequest.UserName = fwHelper.String(plan.UserName.ValueString())
	describeUserPrivilegesRequest.DbName = fwHelper.String(plan.DbName.ValueString())
	describeUserPrivilegesRequest.Object = fwHelper.StringPointerFromValue(plan.Object)
	describeUserPrivilegesRequest.ObjectType = fwHelper.StringPointerFromValue(plan.ObjectType)
	describeUserPrivilegesRequest.ColName = fwHelper.StringPointerFromValue(plan.ColName)

	params, err := service.DescribeTdmysqlUserPrivileges(ctx, describeUserPrivilegesRequest)
	if err != nil {
		diags.AddError("Error reading TDSQL Boundless user privileges after update", err.Error())

		return
	}

	plan.Privileges = fwHelper.ListFromPointers(ctx, types.StringType, params.Privileges, diags)
}

// ImportState splits the composite id into the four locator attributes. The
// id segments are parent-first: instance_id#host#user_name#db_name.
func (r *TdmysqlUserPrivilegeConfigResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	parts, err := fwHelper.IdParseErr(request.ID, 4)
	if err != nil {
		response.Diagnostics.AddError("Invalid TDSQL Boundless user privilege id", fmt.Sprintf("expected instance_id#host#user_name#db_name: %v", err))

		return
	}

	response.State.SetAttribute(ctx, path.Root("id"), request.ID)
	response.State.SetAttribute(ctx, path.Root("instance_id"), parts[0])
	response.State.SetAttribute(ctx, path.Root("host"), parts[1])
	response.State.SetAttribute(ctx, path.Root("user_name"), parts[2])
	response.State.SetAttribute(ctx, path.Root("db_name"), parts[3])
}

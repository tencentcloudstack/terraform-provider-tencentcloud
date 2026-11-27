package tdmysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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

// TdmysqlUserPrivilegeConfigDatabasePrivilegesModel maps one entry of the
// database_privileges list attribute.
type TdmysqlUserPrivilegeConfigDatabasePrivilegesModel struct {
	Database   types.String `tfsdk:"database"`
	Privileges types.Set    `tfsdk:"privileges"`
}

// TdmysqlUserPrivilegeConfigTablePrivilegesModel maps one entry of the
// table_privileges list attribute.
type TdmysqlUserPrivilegeConfigTablePrivilegesModel struct {
	Database   types.String `tfsdk:"database"`
	Table      types.String `tfsdk:"table"`
	Privileges types.Set    `tfsdk:"privileges"`
}

// TdmysqlUserPrivilegeConfigResourceModel maps every top-level schema
// attribute; each one appears here exactly once. The schema mirrors the
// ModifyUserPrivileges WRITE API input parameters (InstanceId, Users,
// GlobalPrivileges, DatabasePrivileges, TablePrivileges) — the read API's
// addressing parameters (db_name/object/object_type/col_name) are NOT schema
// attributes; the read-back derives them from the managed groups.
type TdmysqlUserPrivilegeConfigResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	InstanceId         types.String `tfsdk:"instance_id"`
	UserName           types.String `tfsdk:"user_name"`
	Host               types.String `tfsdk:"host"`
	GlobalPrivileges   types.Set    `tfsdk:"global_privileges"`
	DatabasePrivileges types.Set    `tfsdk:"database_privileges"`
	TablePrivileges    types.Set    `tfsdk:"table_privileges"`
}

// TdmysqlUserPrivilegeConfigDatabasePrivilegesAttrTypes returns the attribute
// types of one database_privileges entry. It is shared by every conversion so
// the schema-facing shape never drifts, and it is exported so the external
// test package can build the same object value.
func TdmysqlUserPrivilegeConfigDatabasePrivilegesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"database":   types.StringType,
		"privileges": types.SetType{ElemType: types.StringType},
	}
}

// TdmysqlUserPrivilegeConfigTablePrivilegesAttrTypes is
// TdmysqlUserPrivilegeConfigDatabasePrivilegesAttrTypes for table_privileges.
func TdmysqlUserPrivilegeConfigTablePrivilegesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"database":   types.StringType,
		"table":      types.StringType,
		"privileges": types.SetType{ElemType: types.StringType},
	}
}

// Metadata reports the Terraform type name.
func (r *TdmysqlUserPrivilegeConfigResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_tdmysql_user_privilege_config"
}

// Schema declares the user privileges of a TDSQL Boundless instance. Every
// attribute name is a mechanical snake_case conversion of a
// ModifyUserPrivileges request field key.
func (r *TdmysqlUserPrivilegeConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage user privileges of a TDSQL Boundless instance. " +
			"Create and update both call the ModifyUserPrivileges API, read calls DescribeUserPrivileges. " +
			"Destroying this resource only removes it from state; the cloud-side privileges live and die with the parent instance and user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the user privilege config, in the format `instance_id#user_name#host`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// Locator attributes — Required + RequiresReplace (InstanceId and
			// the Users entry of the write API request).
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "Instance ID, such as tdsql3-5baee8df.",
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
			"host": schema.StringAttribute{
				Required:    true,
				Description: "The access host allowed for the user. The username and host uniquely identify an account.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// Configuration body — the write API's three privilege groups.
			// Every group is Optional+Computed: a group left unset in the
			// plan is not sent to the cloud, and the read-back records the
			// cloud-reported values for the managed groups. UseStateForUnknown
			// keeps a user-omitted group diff-free on update; on create the
			// apply path normalizes the core-planned unknown back to null
			// before persisting (see applyTdmysqlUserPrivilegeConfig).
			"global_privileges": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Global privileges.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"database_privileges": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Database level privileges.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"database": schema.StringAttribute{
							Required:    true,
							Description: "Database name.",
						},
						"privileges": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "Privilege list.",
						},
					},
				},
			},
			"table_privileges": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Table level privileges.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"database": schema.StringAttribute{
							Required:    true,
							Description: "Database name.",
						},
						"table": schema.StringAttribute{
							Required:    true,
							Description: "Table name.",
						},
						"privileges": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "Privilege list.",
						},
					},
				},
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
	// The segments are parent-first: instance_id#user_name#host.
	id := fwHelper.IdFormat(plan.InstanceId.ValueString(), plan.UserName.ValueString(), plan.Host.ValueString())
	plan.Id = types.StringValue(id)

	service := NewTdmysqlService(r.Client())

	r.applyTdmysqlUserPrivilegeConfig(ctx, service, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the managed privilege groups from the cloud.
func (r *TdmysqlUserPrivilegeConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TdmysqlUserPrivilegeConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Split the composite locator back out of the state id.
	parts, err := fwHelper.IdParseErr(state.Id.ValueString(), 3)
	if err != nil {
		resp.Diagnostics.AddError("Invalid TDSQL Boundless user privilege config id", fmt.Sprintf("expected instance_id#user_name#host: %v", err))

		return
	}

	state.InstanceId = types.StringValue(parts[0])
	state.UserName = types.StringValue(parts[1])
	state.Host = types.StringValue(parts[2])

	service := NewTdmysqlService(r.Client())

	queried, notFound := r.refreshTdmysqlUserPrivilegeConfig(ctx, service, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Every managed scope confirmed absent on the cloud side => the parent
	// user (or the whole privilege set) is gone => drop the resource from
	// state so Terraform plans a recreation. Transport errors already
	// surfaced as errors above.
	if queried > 0 && queried == notFound {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("TDSQL Boundless user privileges %s not found, removing from state.", state.Id.ValueString()),
		)

		resp.State.RemoveResource(ctx)

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

	r.applyTdmysqlUserPrivilegeConfig(ctx, service, &plan, &resp.Diagnostics)
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

// ImportState splits the composite id into the three locator attributes. The
// id segments are parent-first: instance_id#user_name#host.
func (r *TdmysqlUserPrivilegeConfigResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	parts, err := fwHelper.IdParseErr(request.ID, 3)
	if err != nil {
		response.Diagnostics.AddError("Invalid TDSQL Boundless user privilege config id", fmt.Sprintf("expected instance_id#user_name#host: %v", err))

		return
	}

	response.State.SetAttribute(ctx, path.Root("id"), request.ID)
	response.State.SetAttribute(ctx, path.Root("instance_id"), parts[0])
	response.State.SetAttribute(ctx, path.Root("user_name"), parts[1])
	response.State.SetAttribute(ctx, path.Root("host"), parts[2])
}

// applyTdmysqlUserPrivilegeConfig writes the planned privilege groups through
// the ModifyUserPrivileges WRITE API and reads every managed group back. The
// request fields mirror the write API input one to one: a group left null or
// unknown in the plan is not sent, so the cloud keeps its current value for
// it.
func (r *TdmysqlUserPrivilegeConfigResource) applyTdmysqlUserPrivilegeConfig(ctx context.Context, service TdmysqlService, plan *TdmysqlUserPrivilegeConfigResourceModel, diags *diag.Diagnostics) {
	// Terraform core plans a user-omitted Optional+Computed group as
	// UNKNOWN. On create there is no prior state, so UseStateForUnknown has
	// nothing to substitute and the plan value stays unknown; normalize it
	// to the "user omitted" null form here, otherwise the unknown would
	// leak into state ("Provider returned invalid result object after
	// apply") and the managed-group guards below would never see it as
	// omitted. On update the modifier already substituted the prior state
	// value, so an unknown here cannot occur.
	if plan.GlobalPrivileges.IsUnknown() {
		plan.GlobalPrivileges = types.SetNull(types.StringType)
	}

	if plan.DatabasePrivileges.IsUnknown() {
		plan.DatabasePrivileges = types.SetNull(types.ObjectType{AttrTypes: TdmysqlUserPrivilegeConfigDatabasePrivilegesAttrTypes()})
	}

	if plan.TablePrivileges.IsUnknown() {
		plan.TablePrivileges = types.SetNull(types.ObjectType{AttrTypes: TdmysqlUserPrivilegeConfigTablePrivilegesAttrTypes()})
	}

	modifyUserPrivilegesRequest := tdmysqlv20211122.NewModifyUserPrivilegesRequest()
	modifyUserPrivilegesRequest.InstanceId = fwHelper.String(plan.InstanceId.ValueString())
	modifyUserPrivilegesRequest.Users = []*tdmysqlv20211122.User{
		{
			UserName: fwHelper.String(plan.UserName.ValueString()),
			Host:     fwHelper.String(plan.Host.ValueString()),
		},
	}

	if !plan.GlobalPrivileges.IsNull() && !plan.GlobalPrivileges.IsUnknown() {
		modifyUserPrivilegesRequest.GlobalPrivileges = fwHelper.PointersFromSet[string](ctx, plan.GlobalPrivileges, diags)
		if diags.HasError() {
			return
		}
	}

	if !plan.DatabasePrivileges.IsNull() && !plan.DatabasePrivileges.IsUnknown() {
		modifyUserPrivilegesRequest.DatabasePrivileges = expandTdmysqlUserPrivilegeConfigDatabasePrivileges(ctx, plan.DatabasePrivileges, diags)
		if diags.HasError() {
			return
		}
	}

	if !plan.TablePrivileges.IsNull() && !plan.TablePrivileges.IsUnknown() {
		modifyUserPrivilegesRequest.TablePrivileges = expandTdmysqlUserPrivilegeConfigTablePrivileges(ctx, plan.TablePrivileges, diags)
		if diags.HasError() {
			return
		}
	}

	// A plan without any managed group carries no write intent; only the
	// read-back runs.
	if modifyUserPrivilegesRequest.GlobalPrivileges != nil || modifyUserPrivilegesRequest.DatabasePrivileges != nil || modifyUserPrivilegesRequest.TablePrivileges != nil {
		if err := service.ModifyTdmysqlUserPrivileges(ctx, modifyUserPrivilegesRequest); err != nil {
			diags.AddError("Error updating TDSQL Boundless user privileges", err.Error())

			return
		}
	}

	r.refreshTdmysqlUserPrivilegeConfig(ctx, service, plan, diags)
}

// refreshTdmysqlUserPrivilegeConfig reads back every privilege group carried
// by the model (the plan during apply, the state during Read), so the
// Required-for-the-group + Computed privilege values always record the
// cloud-reported values. The read API addresses ONE scope per call; the
// addressing is derived from the managed groups themselves:
//   - global_privileges      -> db_name `*` (global, i.e. *.*)
//   - database_privileges[i] -> db_name <database>, object_type `*` (db.*)
//   - table_privileges[i]    -> db_name <database>, object_type `table`,
//     object <table>, col_name `*` (the table's own privileges)
//
// A group confirmed absent on the cloud side keeps the previously configured
// value (the OrPrev anti-diff contract for Optional+Computed attributes).
// The write direction is ModifyUserPrivileges; no async flow id exists to
// expose.
func (r *TdmysqlUserPrivilegeConfigResource) refreshTdmysqlUserPrivilegeConfig(ctx context.Context, service TdmysqlService, model *TdmysqlUserPrivilegeConfigResourceModel, diags *diag.Diagnostics) (queried int, notFound int) {
	if !model.GlobalPrivileges.IsNull() && !model.GlobalPrivileges.IsUnknown() {
		privileges, empty := r.describeTdmysqlUserPrivilegeConfigScope(ctx, service, model, "*", "", "", "", diags)
		queried++
		if empty {
			notFound++
		}

		if diags.HasError() {
			return
		}

		model.GlobalPrivileges = fwHelper.SetFromPointersOrPrev(ctx, types.StringType, privileges, model.GlobalPrivileges, diags)
		if diags.HasError() {
			return
		}
	}

	if !model.DatabasePrivileges.IsNull() && !model.DatabasePrivileges.IsUnknown() {
		var prevEntries []TdmysqlUserPrivilegeConfigDatabasePrivilegesModel
		diags.Append(model.DatabasePrivileges.ElementsAs(ctx, &prevEntries, false)...)
		if diags.HasError() {
			return
		}

		refreshed := make([]TdmysqlUserPrivilegeConfigDatabasePrivilegesModel, 0, len(prevEntries))
		for _, prev := range prevEntries {
			privileges, empty := r.describeTdmysqlUserPrivilegeConfigScope(ctx, service, model, prev.Database.ValueString(), "*", "", "", diags)
			queried++
			if empty {
				notFound++
			}

			if diags.HasError() {
				return
			}

			refreshed = append(refreshed, TdmysqlUserPrivilegeConfigDatabasePrivilegesModel{
				Database:   prev.Database,
				Privileges: fwHelper.SetFromPointersOrPrev(ctx, types.StringType, privileges, prev.Privileges, diags),
			})
			if diags.HasError() {
				return
			}
		}

		list, listDiags := types.SetValueFrom(ctx, types.SetType{ElemType: types.ObjectType{AttrTypes: TdmysqlUserPrivilegeConfigDatabasePrivilegesAttrTypes()}}, refreshed)
		diags.Append(listDiags...)
		if diags.HasError() {
			return
		}

		model.DatabasePrivileges = list
	}

	if !model.TablePrivileges.IsNull() && !model.TablePrivileges.IsUnknown() {
		var prevEntries []TdmysqlUserPrivilegeConfigTablePrivilegesModel
		diags.Append(model.TablePrivileges.ElementsAs(ctx, &prevEntries, false)...)
		if diags.HasError() {
			return
		}

		refreshed := make([]TdmysqlUserPrivilegeConfigTablePrivilegesModel, 0, len(prevEntries))
		for _, prev := range prevEntries {
			privileges, empty := r.describeTdmysqlUserPrivilegeConfigScope(ctx, service, model, prev.Database.ValueString(), "table", prev.Table.ValueString(), "*", diags)
			queried++
			if empty {
				notFound++
			}

			if diags.HasError() {
				return
			}

			refreshed = append(refreshed, TdmysqlUserPrivilegeConfigTablePrivilegesModel{
				Database:   prev.Database,
				Table:      prev.Table,
				Privileges: fwHelper.SetFromPointersOrPrev(ctx, types.StringType, privileges, prev.Privileges, diags),
			})
			if diags.HasError() {
				return
			}
		}

		list, listDiags := types.SetValueFrom(ctx, types.SetType{ElemType: types.ObjectType{AttrTypes: TdmysqlUserPrivilegeConfigTablePrivilegesAttrTypes()}}, refreshed)
		diags.Append(listDiags...)
		if diags.HasError() {
			return
		}

		model.TablePrivileges = list
	}

	return
}

// describeTdmysqlUserPrivilegeConfigScope queries ONE privilege scope through
// the read API. The notFound result reports the shared ErrNotFound sentinel
// ("confirmed absent" — the describe answer carries no privilege item);
// transport errors surface as diagnostics and a nil list.
func (r *TdmysqlUserPrivilegeConfigResource) describeTdmysqlUserPrivilegeConfigScope(ctx context.Context, service TdmysqlService, model *TdmysqlUserPrivilegeConfigResourceModel, dbName, objectType, object, colName string, diags *diag.Diagnostics) (privileges []*string, notFound bool) {
	describeUserPrivilegesRequest := tdmysqlv20211122.NewDescribeUserPrivilegesRequest()
	describeUserPrivilegesRequest.InstanceId = fwHelper.String(model.InstanceId.ValueString())
	describeUserPrivilegesRequest.Host = fwHelper.String(model.Host.ValueString())
	describeUserPrivilegesRequest.UserName = fwHelper.String(model.UserName.ValueString())
	describeUserPrivilegesRequest.DbName = fwHelper.String(dbName)
	if objectType != "" {
		describeUserPrivilegesRequest.ObjectType = fwHelper.String(objectType)
	}

	if object != "" {
		describeUserPrivilegesRequest.Object = fwHelper.String(object)
	}

	if colName != "" {
		describeUserPrivilegesRequest.ColName = fwHelper.String(colName)
	}

	params, err := service.DescribeTdmysqlUserPrivileges(ctx, describeUserPrivilegesRequest)
	if err != nil {
		if errors.Is(err, fwHelper.ErrNotFound) {
			return nil, true
		}

		diags.AddError("Error reading TDSQL Boundless user privileges", err.Error())

		return nil, false
	}

	return params.Privileges, false
}

// expandTdmysqlUserPrivilegeConfigDatabasePrivileges converts the planned
// database_privileges list into the SDK request entries.
func expandTdmysqlUserPrivilegeConfigDatabasePrivileges(ctx context.Context, value types.Set, diags *diag.Diagnostics) []*tdmysqlv20211122.DatabasePrivileges {
	var entries []TdmysqlUserPrivilegeConfigDatabasePrivilegesModel
	diags.Append(value.ElementsAs(ctx, &entries, false)...)
	if diags.HasError() {
		return nil
	}

	result := make([]*tdmysqlv20211122.DatabasePrivileges, 0, len(entries))
	for _, entry := range entries {
		result = append(result, &tdmysqlv20211122.DatabasePrivileges{
			Database:   fwHelper.StringPointerFromValue(entry.Database),
			Privileges: fwHelper.PointersFromSet[string](ctx, entry.Privileges, diags),
		})
		if diags.HasError() {
			return nil
		}
	}

	return result
}

// expandTdmysqlUserPrivilegeConfigTablePrivileges converts the planned
// table_privileges list into the SDK request entries.
func expandTdmysqlUserPrivilegeConfigTablePrivileges(ctx context.Context, value types.Set, diags *diag.Diagnostics) []*tdmysqlv20211122.TablePrivileges {
	var entries []TdmysqlUserPrivilegeConfigTablePrivilegesModel
	diags.Append(value.ElementsAs(ctx, &entries, false)...)
	if diags.HasError() {
		return nil
	}

	result := make([]*tdmysqlv20211122.TablePrivileges, 0, len(entries))
	for _, entry := range entries {
		result = append(result, &tdmysqlv20211122.TablePrivileges{
			Database:   fwHelper.StringPointerFromValue(entry.Database),
			Table:      fwHelper.StringPointerFromValue(entry.Table),
			Privileges: fwHelper.PointersFromSet[string](ctx, entry.Privileges, diags),
		})
		if diags.HasError() {
			return nil
		}
	}

	return result
}

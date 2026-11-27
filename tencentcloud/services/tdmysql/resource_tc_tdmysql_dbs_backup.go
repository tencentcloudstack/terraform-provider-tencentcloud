package tdmysql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions. Only assert the extension interfaces;
// the base resource.Resource interface is checked implicitly by the factory
// return value.
var (
	_ resource.ResourceWithConfigure   = &TdmysqlDbsBackupResource{}
	_ resource.ResourceWithImportState = &TdmysqlDbsBackupResource{}
)

// NewTdmysqlDbsBackupResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTdmysqlDbsBackupResource() resource.Resource {
	return &TdmysqlDbsBackupResource{}
}

// TdmysqlDbsBackupResource implements resource.Resource for
// tencentcloud_tdmysql_dbs_backup.
type TdmysqlDbsBackupResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// TdmysqlDbsBackupResourceModel maps the schema attributes; each attribute
// appears exactly once. Keep the struct fields in the same order as the schema.
type TdmysqlDbsBackupResourceModel struct {
	Id           types.String `tfsdk:"id"`
	InstanceId   types.String `tfsdk:"instance_id"`   // Required, RequiresReplace
	BackupMethod types.String `tfsdk:"backup_method"` // Required, RequiresReplace
	BackupType   types.String `tfsdk:"backup_type"`   // Required, RequiresReplace
	// BackupName is Optional+Computed: when the user omits it the cloud
	// generates and returns a comment automatically.
	BackupName types.String `tfsdk:"backup_name"` // Optional+Computed, covered by ModifyDBSBackupSetComment
}

func (r *TdmysqlDbsBackupResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_tdmysql_dbs_backup"
}

func (r *TdmysqlDbsBackupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a manual backup set of a TDSQL Boundless (tdmysql) instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, in the form of the backup set ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "Instance ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backup_method": schema.StringAttribute{
				Required:    true,
				Description: "Backup method. Valid values: `physical` (physical backup), `snapshot` (snapshot backup). This value should stay consistent with the `backupMethod` returned by the DescribeDBSBackupPolicy API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backup_type": schema.StringAttribute{
				Required:    true,
				Description: "Backup type. Only `full` is supported for now.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// Optional+Computed: when omitted, the cloud generates and
			// returns a comment automatically. UseStateForUnknown keeps a
			// user-omitted value diff-free — on update core plans the
			// attribute as unknown and the modifier substitutes the prior
			// state value; on create the plan stays unknown, which accepts
			// the cloud-generated value the read-back records.
			"backup_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Backup comment. Modifiable through the ModifyDBSBackupSetComment API. If omitted, the cloud side generates one automatically.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TdmysqlDbsBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TdmysqlDbsBackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createDBSBackupRequest := tdmysqlv20211122.NewCreateDBSBackupRequest()
	createDBSBackupRequest.InstanceId = plan.InstanceId.ValueStringPointer()
	createDBSBackupRequest.BackupMethod = plan.BackupMethod.ValueStringPointer()
	createDBSBackupRequest.BackupType = plan.BackupType.ValueStringPointer()
	if !plan.BackupName.IsNull() && !plan.BackupName.IsUnknown() {
		createDBSBackupRequest.BackupName = plan.BackupName.ValueStringPointer()
	}

	service := NewTdmysqlService(r.Client())

	backupSetId, err := service.CreateDBSBackup(ctx, createDBSBackupRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tdmysql dbs backup",
			fmt.Sprintf("CreateDBSBackup failed: %v", err),
		)

		return
	}

	if backupSetId == 0 {
		resp.Diagnostics.AddError(
			"Error creating tdmysql dbs backup",
			"CreateDBSBackup returned an empty backup set id.",
		)

		return
	}

	plan.Id = types.StringValue(strconv.FormatInt(backupSetId, 10))

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeDBSBackupSetsRequest := tdmysqlv20211122.NewDescribeDBSBackupSetsRequest()
	describeDBSBackupSetsRequest.InstanceId = plan.InstanceId.ValueStringPointer()
	describeDBSBackupSetsRequest.BackupSetId = fwHelper.Int64(backupSetId)
	items, err := service.DescribeDBSBackupSets(ctx, describeDBSBackupSetsRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading tdmysql dbs backup after create", err.Error())
		return
	}

	item := findTdmysqlDbsBackup(items, plan.InstanceId.ValueString(), backupSetId)
	if item == nil {
		resp.Diagnostics.AddError(
			"tdmysql dbs backup not readable after create",
			fmt.Sprintf("backup set %s was created but cannot be read back", plan.Id.ValueString()),
		)

		return
	}

	resp.Diagnostics.Append(refreshTdmysqlDbsBackupModel(&plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TdmysqlDbsBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TdmysqlDbsBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backupSetId, err := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tdmysql dbs backup",
			fmt.Sprintf("invalid backup set id %s: %v", state.Id.ValueString(), err),
		)

		return
	}

	describeDBSBackupSetsRequest := tdmysqlv20211122.NewDescribeDBSBackupSetsRequest()
	describeDBSBackupSetsRequest.InstanceId = state.InstanceId.ValueStringPointer()
	describeDBSBackupSetsRequest.BackupSetId = fwHelper.Int64(backupSetId)

	service := NewTdmysqlService(r.Client())

	items, err := service.DescribeDBSBackupSets(ctx, describeDBSBackupSetsRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading tdmysql dbs backup", err.Error())
		return
	}

	item := findTdmysqlDbsBackup(items, state.InstanceId.ValueString(), backupSetId)

	// Confirmed "not found" on the cloud side => drop it from state so
	// Terraform plans a recreation. The judgement stays INLINE in Read:
	// RemoveResource is followed by an immediate return.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("tdmysql dbs backup %s not found, removing from state.", state.Id.ValueString()),
		)

		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshTdmysqlDbsBackupModel(&state, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. Only `backup_name` is covered by
// the ModifyDBSBackupSetComment API; every other attribute is RequiresReplace.
func (r *TdmysqlDbsBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TdmysqlDbsBackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TdmysqlDbsBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backupSetId, err := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tdmysql dbs backup",
			fmt.Sprintf("invalid backup set id %s: %v", state.Id.ValueString(), err),
		)

		return
	}

	service := NewTdmysqlService(r.Client())

	// Call the write API ONLY when the attribute it covers changed. A plan
	// differing from the state in no covered attribute reaches Update too;
	// sending the write API then would be a pointless no-change request.
	// The read-back below runs either way. With UseStateForUnknown on
	// backup_name, a user-omitted value plans as the prior state value and
	// is equal to the state here — removing the attribute from the config
	// therefore never calls the write API (the cloud-side comment is kept
	// as-is; the API has no clearing semantics).
	if !plan.BackupName.Equal(state.BackupName) {
		modifyDBSBackupSetCommentRequest := tdmysqlv20211122.NewModifyDBSBackupSetCommentRequest()
		modifyDBSBackupSetCommentRequest.InstanceId = state.InstanceId.ValueStringPointer()
		modifyDBSBackupSetCommentRequest.BackupSetId = fwHelper.Int64(backupSetId)
		// A change reaching this point always carries a known, non-null
		// plan value (an unknown/null plan equals the state through
		// UseStateForUnknown and never enters this branch); the guard
		// below stays as defense in depth.
		if !plan.BackupName.IsNull() && !plan.BackupName.IsUnknown() {
			modifyDBSBackupSetCommentRequest.NewBackupName = plan.BackupName.ValueStringPointer()
		}

		if err := service.ModifyDBSBackupSetComment(ctx, modifyDBSBackupSetCommentRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error updating tdmysql dbs backup",
				fmt.Sprintf("ModifyDBSBackupSetComment failed: %v", err),
			)

			return
		}
	}

	describeDBSBackupSetsRequest := tdmysqlv20211122.NewDescribeDBSBackupSetsRequest()
	describeDBSBackupSetsRequest.InstanceId = state.InstanceId.ValueStringPointer()
	describeDBSBackupSetsRequest.BackupSetId = fwHelper.Int64(backupSetId)
	items, err := service.DescribeDBSBackupSets(ctx, describeDBSBackupSetsRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading tdmysql dbs backup after update", err.Error())
		return
	}

	item := findTdmysqlDbsBackup(items, state.InstanceId.ValueString(), backupSetId)
	if item == nil {
		resp.Diagnostics.AddError(
			"tdmysql dbs backup not readable after update",
			fmt.Sprintf("backup set %s was updated but cannot be read back", state.Id.ValueString()),
		)

		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshTdmysqlDbsBackupModel(&plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *TdmysqlDbsBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TdmysqlDbsBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backupSetId, err := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting tdmysql dbs backup",
			fmt.Sprintf("invalid backup set id %s: %v", state.Id.ValueString(), err),
		)

		return
	}

	deleteDBSBackupSetsRequest := tdmysqlv20211122.NewDeleteDBSBackupSetsRequest()
	deleteDBSBackupSetsRequest.InstanceId = state.InstanceId.ValueStringPointer()
	deleteDBSBackupSetsRequest.BackupSetIdList = []*int64{fwHelper.Int64(backupSetId)}

	service := NewTdmysqlService(r.Client())

	if err := service.DeleteDBSBackupSets(ctx, deleteDBSBackupSetsRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting tdmysql dbs backup",
			fmt.Sprintf("DeleteDBSBackupSets failed: %v", err),
		)

		return
	}
}

// refreshTdmysqlDbsBackupModel copies cloud-owned fields into the model.
// All attributes are user-configurable, so the *OrPrev family keeps the
// previous value when the cloud omits a field.
func refreshTdmysqlDbsBackupModel(m *TdmysqlDbsBackupResourceModel, item *tdmysqlv20211122.BackupSetModel) diag.Diagnostics {
	if item == nil {
		return nil
	}

	var diags diag.Diagnostics

	m.InstanceId = fwHelper.StringValueOrPrev(item.InstanceId, m.InstanceId)
	m.BackupMethod = fwHelper.StringValueOrPrev(item.BackupMethod, m.BackupMethod)
	m.BackupType = fwHelper.StringValueOrPrev(item.BackupType, m.BackupType)
	m.BackupName = fwHelper.StringValueOrPrev(item.BackupName, m.BackupName)

	return diags
}

// findTdmysqlDbsBackup picks this resource's own item from the describe list
// by instance id and backup set id — "no match in the list" is not-found.
func findTdmysqlDbsBackup(items []*tdmysqlv20211122.BackupSetModel, instanceId string, backupSetId int64) *tdmysqlv20211122.BackupSetModel {
	for _, item := range items {
		if item.BackupSetId != nil && *item.BackupSetId == backupSetId && item.InstanceId != nil && *item.InstanceId == instanceId {
			return item
		}
	}

	return nil
}

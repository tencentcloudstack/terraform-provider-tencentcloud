package tdmysql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions. The base resource.Resource interface is
// checked implicitly by the factory return value.
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
// tencentcloud_tdmysql_dbs_backup. It manages a manual backup set of a
// TDSQL-C (tdmysql) instance.
type TdmysqlDbsBackupResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// TdmysqlDbsBackupResourceModel maps the schema attributes.
type TdmysqlDbsBackupResourceModel struct {
	Id           types.String `tfsdk:"id"`
	InstanceId   types.String `tfsdk:"instance_id"`
	BackupMethod types.String `tfsdk:"backup_method"`
	BackupType   types.String `tfsdk:"backup_type"`
	BackupName   types.String `tfsdk:"backup_name"`
	BackupSetId  types.Int64  `tfsdk:"backup_set_id"`
	Items        types.List   `tfsdk:"items"`
}

// TdmysqlDbsBackupItemModel maps one element of the `items` list attribute,
// which carries the cloud-side backup set detail returned by
// DescribeDBSBackupSets.
type TdmysqlDbsBackupItemModel struct {
	BackupDuration types.Int64   `tfsdk:"backup_duration"`
	BackupMethod   types.String  `tfsdk:"backup_method"`
	BackupName     types.String  `tfsdk:"backup_name"`
	BackupProgress types.Float64 `tfsdk:"backup_progress"`
	BackupSetId    types.Int64   `tfsdk:"backup_set_id"`
	BackupStatus   types.String  `tfsdk:"backup_status"`
	BackupType     types.String  `tfsdk:"backup_type"`
	DBVersion      types.String  `tfsdk:"db_version"`
	EndTime        types.String  `tfsdk:"end_time"`
	EndTrxTime     types.String  `tfsdk:"end_trx_time"`
	ErrorMessage   types.String  `tfsdk:"error_message"`
	ExpiredTime    types.String  `tfsdk:"expired_time"`
	FileName       types.String  `tfsdk:"file_name"`
	FileSize       types.Int64   `tfsdk:"file_size"`
	InstanceId     types.String  `tfsdk:"instance_id"`
	ManualBackup   types.Int64   `tfsdk:"manual_backup"`
	StartTime      types.String  `tfsdk:"start_time"`
}

func (r *TdmysqlDbsBackupResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_tdmysql_dbs_backup"
}

func (r *TdmysqlDbsBackupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Provides a resource to create a manual backup set for a TDSQL-C (tdmysql) instance.",
		MarkdownDescription: "Provides a resource to create a manual backup set for a TDSQL-C (tdmysql) instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the backup set in the format `instance_id#backup_set_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the TDSQL-C instance to back up. Changing this value forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backup_method": schema.StringAttribute{
				Required:    true,
				Description: "Backup method. It must be consistent with the `backupMethod` returned by DescribeDBSBackupPolicy. Valid values: `physical`, `snapshot`. Changing this value forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backup_type": schema.StringAttribute{
				Required:    true,
				Description: "Backup type. Currently only `full` is supported. Changing this value forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backup_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Backup comment name. It can be updated in place through ModifyDBSBackupSetComment.",
			},
			"backup_set_id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the backup set returned by CreateDBSBackup.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Backup set list returned by DescribeDBSBackupSets.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"backup_duration": schema.Int64Attribute{
							Computed:    true,
							Description: "Backup set duration in seconds.",
						},
						"backup_method": schema.StringAttribute{
							Computed:    true,
							Description: "Backup method.",
						},
						"backup_name": schema.StringAttribute{
							Computed:    true,
							Description: "Backup comment name.",
						},
						"backup_progress": schema.Float64Attribute{
							Computed:    true,
							Description: "Backup progress percentage.",
						},
						"backup_set_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Backup set ID.",
						},
						"backup_status": schema.StringAttribute{
							Computed:    true,
							Description: "Backup set status. Valid values: `pending`, `running`, `success`, `failed`.",
						},
						"backup_type": schema.StringAttribute{
							Computed:    true,
							Description: "Backup set type.",
						},
						"db_version": schema.StringAttribute{
							Computed:    true,
							Description: "Instance version.",
						},
						"end_time": schema.StringAttribute{
							Computed:    true,
							Description: "Backup end time.",
						},
						"end_trx_time": schema.StringAttribute{
							Computed:    true,
							Description: "Transaction commit end time.",
						},
						"error_message": schema.StringAttribute{
							Computed:    true,
							Description: "Error message.",
						},
						"expired_time": schema.StringAttribute{
							Computed:    true,
							Description: "Backup expiration time.",
						},
						"file_name": schema.StringAttribute{
							Computed:    true,
							Description: "Backup file name.",
						},
						"file_size": schema.Int64Attribute{
							Computed:    true,
							Description: "Backup file size in bytes.",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "Instance ID.",
						},
						"manual_backup": schema.Int64Attribute{
							Computed:    true,
							Description: "Backup trigger. Valid values: `0` (auto backup), `1` (manual backup).",
						},
						"start_time": schema.StringAttribute{
							Computed:    true,
							Description: "Backup start time.",
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for new backup sets.
func (r *TdmysqlDbsBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TdmysqlDbsBackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewTdmysqlService(r.Client())

	request := tdmysqlv20211122.NewCreateDBSBackupRequest()
	request.InstanceId = plan.InstanceId.ValueStringPointer()
	request.BackupMethod = plan.BackupMethod.ValueStringPointer()
	request.BackupType = plan.BackupType.ValueStringPointer()
	if !plan.BackupName.IsNull() && !plan.BackupName.IsUnknown() {
		request.BackupName = plan.BackupName.ValueStringPointer()
	}

	response, err := service.CreateTdmysqlDBSBackup(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating tdmysql DBS backup", fmt.Sprintf("CreateDBSBackup failed: %v", err))
		return
	}

	// A nil response or a nil backup set id is an API failure, not an empty
	// result. Never write an empty id into state: the resource would become
	// unmanageable.
	if response == nil || response.BackupSetId == nil {
		resp.Diagnostics.AddError(
			"Error creating tdmysql DBS backup",
			"CreateDBSBackup returned an empty response or an empty backup set id.",
		)
		return
	}

	plan.BackupSetId = types.Int64Value(*response.BackupSetId)
	plan.Id = types.StringValue(fmt.Sprintf("%s%s%d",
		plan.InstanceId.ValueString(), tccommon.FILED_SP, *response.BackupSetId))

	// Persist the identity BEFORE any follow-up read: if the read fails, the
	// resource is still in state and `terraform destroy` can clean it up.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	items, err := service.DescribeTdmysqlDBSBackupSets(ctx, newDescribeDBSBackupSetsRequest(
		plan.InstanceId.ValueString(), *response.BackupSetId))
	if err != nil {
		resp.Diagnostics.AddError("Error reading tdmysql DBS backup after create", err.Error())
		return
	}

	item := findBackupSetItem(items, *response.BackupSetId)
	if item == nil {
		resp.Diagnostics.AddError(
			"tdmysql DBS backup not readable after create",
			fmt.Sprintf("backup set %d was created but cannot be read back", *response.BackupSetId),
		)
		return
	}
	refreshBackupModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TdmysqlDbsBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TdmysqlDbsBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewTdmysqlService(r.Client())

	items, err := service.DescribeTdmysqlDBSBackupSets(ctx, newDescribeDBSBackupSetsRequest(
		state.InstanceId.ValueString(), state.BackupSetId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading tdmysql DBS backup", err.Error())
		return
	}

	// Queried successfully but nothing matched => the backup set no longer
	// exists on the cloud side; drop it from state so Terraform plans a
	// recreation.
	item := findBackupSetItem(items, state.BackupSetId.ValueInt64())
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("tdmysql DBS backup set %d of instance %s not found, removing from state.",
				state.BackupSetId.ValueInt64(), state.InstanceId.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	refreshBackupModel(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates. Only the backup comment name is mutable;
// it is updated through ModifyDBSBackupSetComment.
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

	service := NewTdmysqlService(r.Client())

	if !plan.BackupName.Equal(state.BackupName) {
		request := tdmysqlv20211122.NewModifyDBSBackupSetCommentRequest()
		request.InstanceId = state.InstanceId.ValueStringPointer()
		request.BackupSetId = state.BackupSetId.ValueInt64Pointer()
		request.NewBackupName = plan.BackupName.ValueStringPointer()

		_, err := service.ModifyTdmysqlDBSBackupSetComment(ctx, request)
		if err != nil {
			resp.Diagnostics.AddError("Error updating tdmysql DBS backup", fmt.Sprintf("ModifyDBSBackupSetComment failed: %v", err))
			return
		}
	}

	items, err := service.DescribeTdmysqlDBSBackupSets(ctx, newDescribeDBSBackupSetsRequest(
		state.InstanceId.ValueString(), state.BackupSetId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading tdmysql DBS backup after update", err.Error())
		return
	}

	item := findBackupSetItem(items, state.BackupSetId.ValueInt64())
	if item == nil {
		resp.Diagnostics.AddError(
			"tdmysql DBS backup not readable after update",
			fmt.Sprintf("backup set %d was updated but cannot be read back", state.BackupSetId.ValueInt64()),
		)
		return
	}
	plan.Id = state.Id
	plan.BackupSetId = state.BackupSetId
	refreshBackupModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *TdmysqlDbsBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TdmysqlDbsBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewTdmysqlService(r.Client())

	request := tdmysqlv20211122.NewDeleteDBSBackupSetsRequest()
	request.InstanceId = state.InstanceId.ValueStringPointer()
	request.BackupSetIdList = []*int64{state.BackupSetId.ValueInt64Pointer()}

	_, err := service.DeleteTdmysqlDBSBackupSets(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting tdmysql DBS backup", fmt.Sprintf("DeleteDBSBackupSets failed: %v", err))
		return
	}
}

// ImportState parses the composite id `instance_id#backup_set_id`.
func (r *TdmysqlDbsBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := helper.IdParseErr(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("tdmysql DBS backup id must be in the format `instance_id#backup_set_id`, got %q", req.ID),
		)
		return
	}

	backupSetId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("tdmysql DBS backup id segment `backup_set_id` must be an integer, got %q", parts[1]),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("backup_set_id"), backupSetId)...)
}

// newDescribeDBSBackupSetsRequest builds the describe request used by every
// read-back path.
func newDescribeDBSBackupSetsRequest(instanceId string, backupSetId int64) *tdmysqlv20211122.DescribeDBSBackupSetsRequest {
	request := tdmysqlv20211122.NewDescribeDBSBackupSetsRequest()
	request.InstanceId = tccommon.String(instanceId)
	request.BackupSetId = helper.Int64(backupSetId)
	return request
}

// findBackupSetItem picks the item whose BackupSetId matches. The describe
// request is already filtered by instance and backup set id, so the returned
// list is expected to hold at most one item.
func findBackupSetItem(items []*tdmysqlv20211122.BackupSetModel, backupSetId int64) *tdmysqlv20211122.BackupSetModel {
	for _, item := range items {
		if item != nil && item.BackupSetId != nil && *item.BackupSetId == backupSetId {
			return item
		}
	}
	return nil
}

// refreshBackupModel copies the cloud-side backup set detail into the model.
// The items element type is taken from the existing list value so that both
// populated and null lists can be rebuilt consistently.
func refreshBackupModel(m *TdmysqlDbsBackupResourceModel, item *tdmysqlv20211122.BackupSetModel) {
	if item == nil {
		return
	}

	itemModel := TdmysqlDbsBackupItemModel{
		BackupDuration: types.Int64PointerValue(item.BackupDuration),
		BackupMethod:   types.StringPointerValue(item.BackupMethod),
		BackupName:     types.StringPointerValue(item.BackupName),
		BackupProgress: types.Float64PointerValue(item.BackupProgress),
		BackupSetId:    types.Int64PointerValue(item.BackupSetId),
		BackupStatus:   types.StringPointerValue(item.BackupStatus),
		BackupType:     types.StringPointerValue(item.BackupType),
		DBVersion:      types.StringPointerValue(item.DBVersion),
		EndTime:        types.StringPointerValue(item.EndTime),
		EndTrxTime:     types.StringPointerValue(item.EndTrxTime),
		ErrorMessage:   types.StringPointerValue(item.ErrorMessage),
		ExpiredTime:    types.StringPointerValue(item.ExpiredTime),
		FileName:       types.StringPointerValue(item.FileName),
		FileSize:       types.Int64PointerValue(item.FileSize),
		InstanceId:     types.StringPointerValue(item.InstanceId),
		ManualBackup:   types.Int64PointerValue(item.ManualBackup),
		StartTime:      types.StringPointerValue(item.StartTime),
	}
	elemType := m.Items.ElementType(context.Background())
	list, diags := types.ListValueFrom(context.Background(), elemType, []TdmysqlDbsBackupItemModel{itemModel})
	if diags.HasError() {
		return
	}
	m.Items = list

	if item.BackupName != nil && m.BackupName.IsNull() {
		m.BackupName = types.StringValue(*item.BackupName)
	}
}

package ags

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// agsSandboxInstancePollInterval is the interval used while waiting for a
// sandbox instance to reach a stable status.
const agsSandboxInstancePollInterval = 3 * time.Second

var (
	_ resource.ResourceWithConfigure   = &AgsSandboxInstanceResource{}
	_ resource.ResourceWithImportState = &AgsSandboxInstanceResource{}
)

// NewAgsSandboxInstanceResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAgsSandboxInstanceResource() resource.Resource {
	return &AgsSandboxInstanceResource{}
}

// AgsSandboxInstanceResource implements resource.Resource for
// tencentcloud_ags_sandbox_instance.
type AgsSandboxInstanceResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AgsSandboxInstanceResourceModel maps the schema attributes.
//
// Timeout, SourceInstanceId, SnapshotId, ClientToken, StorageMounts and
// StartMode are write-only inputs: DescribeSandboxInstances does not report
// them (the describe API only reports TimeoutSeconds as an integer, which
// cannot be converted back into the original "5m"/"300s"/"1h" form without
// guessing the unit), so Read keeps whatever value is already stored in state.
type AgsSandboxInstanceResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	InstanceId          types.String `tfsdk:"instance_id"`
	ToolId              types.String `tfsdk:"tool_id"`
	ToolName            types.String `tfsdk:"tool_name"`
	SourceInstanceId    types.String `tfsdk:"source_instance_id"`
	SnapshotId          types.String `tfsdk:"snapshot_id"`
	Timeout             types.String `tfsdk:"timeout"`
	ClientToken         types.String `tfsdk:"client_token"`
	MountOptions        types.List   `tfsdk:"mount_options"`
	CustomConfiguration types.Object `tfsdk:"custom_configuration"`
	AuthMode            types.String `tfsdk:"auth_mode"`
	Metadata            types.List   `tfsdk:"metadata"`
	StorageMounts       types.List   `tfsdk:"storage_mounts"`
	StartMode           types.String `tfsdk:"start_mode"`
	Status              types.String `tfsdk:"status"`
	Persistent          types.Bool   `tfsdk:"persistent"`
	TimeoutSeconds      types.Int64  `tfsdk:"timeout_seconds"`
	ExpiresAt           types.String `tfsdk:"expires_at"`
	StopReason          types.String `tfsdk:"stop_reason"`
	CreateTime          types.String `tfsdk:"create_time"`
	UpdateTime          types.String `tfsdk:"update_time"`
	NetworkMode         types.String `tfsdk:"network_mode"`
}

type agsMountOptionModel struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SubPath   types.String `tfsdk:"sub_path"`
	ReadOnly  types.Bool   `tfsdk:"read_only"`
}

type agsMetadataModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

var agsMountOptionAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Name of the storage configuration in the sandbox tool.",
	},
	"mount_path": schema.StringAttribute{
		Optional:    true,
		Description: "Local mount path inside the sandbox instance. It inherits the storage configuration of the tool when it is not set.",
	},
	"sub_path": schema.StringAttribute{
		Optional:    true,
		Description: "Sub path of the storage mount inside the sandbox instance.",
	},
	"read_only": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Whether the mount is read only. It inherits the storage configuration of the tool when it is not set.",
	},
}

var agsMetadataAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Metadata name.",
	},
	"value": schema.StringAttribute{
		Optional:    true,
		Description: "Metadata value.",
	},
}

var (
	agsMountOptionAttrTypes = attributeTypesOf(agsMountOptionAttributes)
	agsMetadataAttrTypes    = attributeTypesOf(agsMetadataAttributes)

	agsMountOptionObjectType = types.ObjectType{AttrTypes: agsMountOptionAttrTypes}
	agsMetadataObjectType    = types.ObjectType{AttrTypes: agsMetadataAttrTypes}
)

func (r *AgsSandboxInstanceResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_instance"
}

func (r *AgsSandboxInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to start an AGS sandbox instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the resource, which is the same as `instance_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the sandbox instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Id of the sandbox tool. Either `tool_id` or `tool_name` must be set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the sandbox tool. Either `tool_id` or `tool_name` must be set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "Id of the source sandbox instance to clone from.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_id": schema.StringAttribute{
				Optional:    true,
				Description: "Id of the sandbox snapshot to start from.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeout": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Timeout of the sandbox instance. The instance is recycled automatically after it expires. " +
					"Supported formats are `5m`, `300s`, `1h`, with a minimum of 30s and a maximum of 24h. " +
					"It is not reported by the describe API as a duration string, so the configured value is kept in state " +
					"and drift on it cannot be detected.",
			},
			"client_token": schema.StringAttribute{
				Optional: true,
				Description: "Idempotency token of the start request, up to 64 characters. When it is not set, a " +
					"random token is generated for each apply so that provider level retries never create a duplicate instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mount_options": schema.ListNestedAttribute{
				Optional:     true,
				Computed:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsMountOptionAttributes},
				Description:  "Storage mount options of the sandbox instance.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Attributes:  agsCustomConfigurationAttributes,
				Description: "Custom configuration of the sandbox instance.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"auth_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Access authentication mode of the sandbox. Valid values: `DEFAULT`, `TOKEN`, `NONE`, `PUBLIC`. " +
					"Default value: `DEFAULT`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.ListNestedAttribute{
				Optional:     true,
				Computed:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsMetadataAttributes},
				Description:  "Metadata of the sandbox instance.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"storage_mounts": schema.ListNestedAttribute{
				Optional:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsStorageMountAttributes},
				Description: "Additional storage mount rules applied when the sandbox instance starts. Only `volume`, `cfs`, " +
					"`cos` and `image` sources are supported. It is not reported by the describe API, so the configured " +
					"value is kept in state and drift on it cannot be detected.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"start_mode": schema.StringAttribute{
				Optional: true,
				Description: "Start mode of the sandbox instance. Valid values: `DEFAULT`, `WAIT_UNTIL_HEALTHY`. " +
					"`WAIT_UNTIL_HEALTHY` uses the persistent background creation flow and the instance may be reported as " +
					"`STARTING` right after the request is accepted; it is currently only supported by WAA sandboxes and " +
					"GPU sandboxes. Default value: `DEFAULT`. It is not reported by the describe API, so the configured " +
					"value is kept in state and drift on it cannot be detected.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Status of the sandbox instance. Valid values: `STARTING`, `RUNNING`, `STOPPING`, `STOPPED`, " +
					"`STOP_FAILED`, `FAILED`, `RESIZING`.",
			},
			"persistent": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the sandbox instance is persistent.",
			},
			"timeout_seconds": schema.Int64Attribute{
				Computed:    true,
				Description: "Timeout of the sandbox instance in seconds. Null means no timeout is set.",
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "Expiry time of the sandbox instance in ISO8601 format. Null means no expiry time.",
			},
			"stop_reason": schema.StringAttribute{
				Computed: true,
				Description: "Reason why the sandbox instance was stopped. Valid values: `manual`, `timeout`, `error`, `system`. " +
					"It is only reported when the status is `STOPPED`, `STOP_FAILED` or `FAILED`.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the sandbox instance in ISO8601 format.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the sandbox instance in ISO8601 format.",
			},
			"network_mode": schema.StringAttribute{
				Computed: true,
				Description: "Network mode of the sandbox instance. Valid values: `PUBLIC`, `SANDBOX`, `INTERNAL_SERVICE`. " +
					"It overrides the network configuration of the tool.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *AgsSandboxInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, diags := r.buildStartRequest(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	instance, err := service.StartSandboxInstance(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error starting AGS sandbox instance",
			fmt.Sprintf("StartSandboxInstance failed: %v", err),
		)
		return
	}

	instanceId := *instance.InstanceId
	plan.Id = types.StringValue(instanceId)
	plan.InstanceId = types.StringValue(instanceId)
	normalizeSandboxInstanceUnknownAttributes(&plan)

	// Persist the identity before waiting: if the wait below fails, the resource
	// is still in state and `terraform destroy` can stop it.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.waitSandboxInstanceRunning(ctx, service, instanceId, tccommon.WriteRetryTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for AGS sandbox instance to become running",
			fmt.Sprintf("sandbox instance %s was started but is not running yet: %v", instanceId, err),
		)
		return
	}

	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and the post-apply state refresh.
func (r *AgsSandboxInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, diags := sandboxInstanceIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	item, err := service.DescribeSandboxInstanceById(ctx, instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox instance",
			fmt.Sprintf("DescribeSandboxInstances failed for %s: %v", instanceId, err),
		)
		return
	}

	// A nil result means the query succeeded and the instance no longer exists.
	// A nil API response is handled as an error inside the service and never
	// reaches this branch, so state is not dropped on a failed request.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"AGS sandbox instance not found",
			fmt.Sprintf("sandbox instance %s does not exist any more, removing it from state.", instanceId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(r.refreshModel(ctx, &state, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
//
// UpdateSandboxInstance only accepts Timeout and Metadata; every other
// attribute is create-only and marked with RequiresReplace in the schema.
func (r *AgsSandboxInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, diags := sandboxInstanceIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := agsv20250920.NewUpdateSandboxInstanceRequest()
	request.InstanceId = helper.String(instanceId)

	// The API documentation states that an omitted Timeout keeps the current
	// setting, so it is only transmitted when the configured value changed.
	// Re-sending an identical value would restart the countdown, which is not
	// what an unchanged attribute should do.
	if !plan.Timeout.IsNull() && !plan.Timeout.IsUnknown() && !plan.Timeout.Equal(state.Timeout) {
		request.Timeout = plan.Timeout.ValueStringPointer()
	}

	// Metadata is an updatable field; it is only transmitted when the configured
	// value changed so that a no-change apply never triggers cloud side work.
	if !plan.Metadata.IsNull() && !plan.Metadata.IsUnknown() && !plan.Metadata.Equal(state.Metadata) {
		metadata, metadataDiags := expandMetadataVars(ctx, plan.Metadata)
		resp.Diagnostics.Append(metadataDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.Metadata = metadata
	}

	service := NewAgsService(r.Client())
	if err := service.UpdateSandboxInstance(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating AGS sandbox instance",
			fmt.Sprintf("UpdateSandboxInstance failed for %s: %v", instanceId, err),
		)
		return
	}

	item, err := service.DescribeSandboxInstanceById(ctx, instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox instance after update",
			fmt.Sprintf("DescribeSandboxInstances failed for %s: %v", instanceId, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"AGS sandbox instance not readable after update",
			fmt.Sprintf("sandbox instance %s was updated but cannot be read back", instanceId),
		)
		return
	}

	plan.Id = state.Id
	plan.InstanceId = state.InstanceId

	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`. It only stops the sandbox instance
// itself; it does not touch the source instance, the sandbox tool or any
// storage that the instance references.
func (r *AgsSandboxInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, diags := sandboxInstanceIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	if err := service.StopSandboxInstance(ctx, instanceId); err != nil {
		// The API documentation does not define an "already stopped" error code,
		// so existence is checked again instead of guessing codes: an instance
		// that is really gone makes the delete idempotent, everything else is
		// an error and keeps the resource in state for a retry.
		item, describeErr := service.DescribeSandboxInstanceById(ctx, instanceId)
		if describeErr != nil || item != nil {
			resp.Diagnostics.AddError(
				"Error stopping AGS sandbox instance",
				fmt.Sprintf("StopSandboxInstance failed for %s: %v", instanceId, err),
			)
		}
		return
	}

	// The stop request is asynchronous: the instance first goes through the
	// STOPPING status and may be recycled afterwards. Wait until it is either
	// STOPPED or no longer visible so that a later apply behaves predictably.
	// A terminal failure of the stop request is reported by the cloud side
	// through the STOP_FAILED status, which is an error rather than a warning.
	if err := r.waitSandboxInstanceStopped(ctx, service, instanceId, tccommon.WriteRetryTimeout); err != nil {
		var terminalStopError *agsStopTerminalError
		if errors.As(err, &terminalStopError) {
			resp.Diagnostics.AddError(
				"Error stopping AGS sandbox instance",
				fmt.Sprintf("sandbox instance %s failed to stop: %v", instanceId, err),
			)
			return
		}
		resp.Diagnostics.AddWarning(
			"AGS sandbox instance is still being stopped",
			fmt.Sprintf("sandbox instance %s was stopped but is not in a terminal status yet: %v", instanceId, err),
		)
	}
}

// ImportState supports `terraform import` with the sandbox instance id.
func (r *AgsSandboxInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Invalid import id",
			"The import id must be the id of an existing sandbox instance.",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), req.ID)...)
}

// buildStartRequest maps the plan onto a StartSandboxInstance request.
func (r *AgsSandboxInstanceResource) buildStartRequest(ctx context.Context, plan *AgsSandboxInstanceResourceModel) (*agsv20250920.StartSandboxInstanceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	request := agsv20250920.NewStartSandboxInstanceRequest()
	if !plan.ToolId.IsNull() && !plan.ToolId.IsUnknown() {
		request.ToolId = plan.ToolId.ValueStringPointer()
	}
	if !plan.ToolName.IsNull() && !plan.ToolName.IsUnknown() {
		request.ToolName = plan.ToolName.ValueStringPointer()
	}
	if !plan.SourceInstanceId.IsNull() && !plan.SourceInstanceId.IsUnknown() {
		request.SourceInstanceId = plan.SourceInstanceId.ValueStringPointer()
	}
	if !plan.SnapshotId.IsNull() && !plan.SnapshotId.IsUnknown() {
		request.SnapshotId = plan.SnapshotId.ValueStringPointer()
	}
	if !plan.Timeout.IsNull() && !plan.Timeout.IsUnknown() {
		request.Timeout = plan.Timeout.ValueStringPointer()
	}
	if !plan.AuthMode.IsNull() && !plan.AuthMode.IsUnknown() {
		request.AuthMode = plan.AuthMode.ValueStringPointer()
	}
	if !plan.StartMode.IsNull() && !plan.StartMode.IsUnknown() {
		request.StartMode = plan.StartMode.ValueStringPointer()
	}

	// The idempotency token protects the provider level retry inside the
	// service layer: a token fixed per apply makes a retried StartSandboxInstance
	// call return the same instance instead of creating a duplicate one. The
	// user supplied token, when present, is forwarded as it is.
	if !plan.ClientToken.IsNull() && !plan.ClientToken.IsUnknown() {
		request.ClientToken = plan.ClientToken.ValueStringPointer()
	} else {
		request.ClientToken = helper.String(helper.BuildUUID())
	}

	mountOptions, mountDiags := expandMountOptions(ctx, plan.MountOptions)
	diags.Append(mountDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.MountOptions = mountOptions

	metadata, metadataDiags := expandMetadataVars(ctx, plan.Metadata)
	diags.Append(metadataDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.Metadata = metadata

	storageMounts, storageDiags := expandStorageMounts(ctx, plan.StorageMounts)
	diags.Append(storageDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.StorageMounts = storageMounts

	customConfiguration, customDiags := expandCustomConfiguration(ctx, plan.CustomConfiguration)
	diags.Append(customDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.CustomConfiguration = customConfiguration

	return request, diags
}

// refreshModel copies the cloud state into the model. Only fields reported by
// the describe API are overwritten, so create-only attributes that cannot be
// read back (source_instance_id, snapshot_id, client_token, storage_mounts,
// start_mode and the timeout duration string) keep the value already stored in
// the model.
func (r *AgsSandboxInstanceResource) refreshModel(ctx context.Context, model *AgsSandboxInstanceResourceModel, item *agsv20250920.SandboxInstance) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		return diags
	}

	if item.InstanceId != nil && *item.InstanceId != "" {
		model.InstanceId = types.StringValue(*item.InstanceId)
		model.Id = types.StringValue(*item.InstanceId)
	}

	// ToolId and ToolName are only overwritten when the describe API reports
	// them, so an instance cloned by source_instance_id keeps the values that
	// were used to start it instead of being reset to null.
	if item.ToolId != nil {
		model.ToolId = types.StringValue(*item.ToolId)
	}
	if item.ToolName != nil {
		model.ToolName = types.StringValue(*item.ToolName)
	}

	model.Status = stringPointerValue(item.Status)
	model.Persistent = boolPointerValue(item.Persistent)
	model.ExpiresAt = stringPointerValue(item.ExpiresAt)
	model.StopReason = stringPointerValue(item.StopReason)
	model.CreateTime = stringPointerValue(item.CreateTime)
	model.UpdateTime = stringPointerValue(item.UpdateTime)
	model.NetworkMode = stringPointerValue(item.NetworkMode)
	model.AuthMode = stringPointerValue(item.AuthMode)
	model.TimeoutSeconds = uint64PointerValue(item.TimeoutSeconds)

	mountOptions, mountDiags := flattenMountOptions(ctx, item.MountOptions)
	diags.Append(mountDiags...)
	if diags.HasError() {
		return diags
	}
	model.MountOptions = mountOptions

	metadata, metadataDiags := flattenMetadataVars(ctx, item.Metadata)
	diags.Append(metadataDiags...)
	if diags.HasError() {
		return diags
	}
	model.Metadata = metadata

	customConfiguration, customDiags := flattenCustomConfiguration(ctx, item.CustomConfiguration)
	diags.Append(customDiags...)
	if diags.HasError() {
		return diags
	}
	model.CustomConfiguration = customConfiguration

	return diags
}

// waitSandboxInstanceRunning polls the instance status until it is RUNNING.
// The start request is asynchronous: with StartMode WAIT_UNTIL_HEALTHY the
// cloud side may report STARTING right after accepting the request.
func (r *AgsSandboxInstanceResource) waitSandboxInstanceRunning(ctx context.Context, service AgsService, instanceId string, timeout time.Duration) (*agsv20250920.SandboxInstance, error) {
	deadline := time.Now().Add(timeout)
	for {
		item, err := service.DescribeSandboxInstanceById(ctx, instanceId)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, fmt.Errorf("sandbox instance %s disappeared while waiting for it to become running", instanceId)
		}

		status := ""
		if item.Status != nil {
			status = *item.Status
		}
		switch status {
		case AGS_SANDBOX_INSTANCE_STATUS_RUNNING:
			return item, nil
		case AGS_SANDBOX_INSTANCE_STATUS_FAILED,
			AGS_SANDBOX_INSTANCE_STATUS_STOP_FAILED,
			AGS_SANDBOX_INSTANCE_STATUS_STOPPED:
			return nil, fmt.Errorf("sandbox instance %s is in %s status", instanceId, status)
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for sandbox instance %s to become RUNNING, current status is %q", instanceId, status)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(agsSandboxInstancePollInterval):
		}
	}
}

// agsStopTerminalError marks a stop that reached a failure status instead of
// the expected STOPPED status, so the caller can report it as an error instead
// of a warning.
type agsStopTerminalError struct {
	status string
}

func (e *agsStopTerminalError) Error() string {
	return fmt.Sprintf("instance stopped with an unexpected status %q", e.status)
}

// waitSandboxInstanceStopped polls the instance status until the stop request
// completes. A stopped instance may either stay in the STOPPED status or be
// recycled by the cloud side, so both outcomes are treated as done.
func (r *AgsSandboxInstanceResource) waitSandboxInstanceStopped(ctx context.Context, service AgsService, instanceId string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		item, err := service.DescribeSandboxInstanceById(ctx, instanceId)
		if err != nil {
			return err
		}
		if item == nil {
			return nil
		}

		status := ""
		if item.Status != nil {
			status = *item.Status
		}
		switch status {
		case AGS_SANDBOX_INSTANCE_STATUS_STOPPED:
			return nil
		case AGS_SANDBOX_INSTANCE_STATUS_FAILED,
			AGS_SANDBOX_INSTANCE_STATUS_STOP_FAILED:
			return &agsStopTerminalError{status: status}
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for sandbox instance %s to stop, current status is %q", instanceId, status)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(agsSandboxInstancePollInterval):
		}
	}
}

// sandboxInstanceIdOf returns the identity stored in the model.
func sandboxInstanceIdOf(model *AgsSandboxInstanceResourceModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !model.InstanceId.IsNull() && !model.InstanceId.IsUnknown() && model.InstanceId.ValueString() != "" {
		return model.InstanceId.ValueString(), diags
	}
	if !model.Id.IsNull() && !model.Id.IsUnknown() && model.Id.ValueString() != "" {
		return model.Id.ValueString(), diags
	}

	diags.AddError(
		"Invalid AGS sandbox instance state",
		"The resource state does not contain a sandbox instance id.",
	)
	return "", diags
}

// normalizeSandboxInstanceUnknownAttributes replaces the optional+computed
// attributes that are still unknown at create time with null, so the identity
// can be stored in state before the cloud has reported those values.
func normalizeSandboxInstanceUnknownAttributes(model *AgsSandboxInstanceResourceModel) {
	if model.ToolId.IsUnknown() {
		model.ToolId = types.StringNull()
	}
	if model.ToolName.IsUnknown() {
		model.ToolName = types.StringNull()
	}
	if model.Timeout.IsUnknown() {
		model.Timeout = types.StringNull()
	}
	if model.MountOptions.IsUnknown() {
		model.MountOptions = types.ListNull(agsMountOptionObjectType)
	}
	if model.CustomConfiguration.IsUnknown() {
		model.CustomConfiguration = types.ObjectNull(agsCustomConfigAttrTypes)
	}
	if model.AuthMode.IsUnknown() {
		model.AuthMode = types.StringNull()
	}
	if model.Metadata.IsUnknown() {
		model.Metadata = types.ListNull(agsMetadataObjectType)
	}
}

// expandMountOptions converts the mount options attribute into the SDK form.
func expandMountOptions(ctx context.Context, list types.List) ([]*agsv20250920.MountOption, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsMountOptionModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	mountOptions := make([]*agsv20250920.MountOption, 0, len(models))
	for _, model := range models {
		mountOption := &agsv20250920.MountOption{}
		if !model.Name.IsNull() && !model.Name.IsUnknown() {
			mountOption.Name = model.Name.ValueStringPointer()
		}
		if !model.MountPath.IsNull() && !model.MountPath.IsUnknown() {
			mountOption.MountPath = model.MountPath.ValueStringPointer()
		}
		if !model.SubPath.IsNull() && !model.SubPath.IsUnknown() {
			mountOption.SubPath = model.SubPath.ValueStringPointer()
		}
		if !model.ReadOnly.IsNull() && !model.ReadOnly.IsUnknown() {
			mountOption.ReadOnly = model.ReadOnly.ValueBoolPointer()
		}
		mountOptions = append(mountOptions, mountOption)
	}

	return mountOptions, diags
}

// flattenMountOptions converts the SDK mount options into the attribute form.
func flattenMountOptions(ctx context.Context, items []*agsv20250920.MountOption) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsMountOptionObjectType), nil
	}

	models := make([]agsMountOptionModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsMountOptionModel{
			Name:      stringPointerValue(item.Name),
			MountPath: stringPointerValue(item.MountPath),
			SubPath:   stringPointerValue(item.SubPath),
			ReadOnly:  boolPointerValue(item.ReadOnly),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsMountOptionObjectType), nil
	}

	return types.ListValueFrom(ctx, agsMountOptionObjectType, models)
}

// expandMetadataVars converts the metadata attribute into the SDK form.
func expandMetadataVars(ctx context.Context, list types.List) ([]*agsv20250920.MetadataVar, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsMetadataModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	metadata := make([]*agsv20250920.MetadataVar, 0, len(models))
	for _, model := range models {
		metadataVar := &agsv20250920.MetadataVar{}
		if !model.Name.IsNull() && !model.Name.IsUnknown() {
			metadataVar.Name = model.Name.ValueStringPointer()
		}
		if !model.Value.IsNull() && !model.Value.IsUnknown() {
			metadataVar.Value = model.Value.ValueStringPointer()
		}
		metadata = append(metadata, metadataVar)
	}

	return metadata, diags
}

// flattenMetadataVars converts the SDK metadata into the attribute form.
func flattenMetadataVars(ctx context.Context, items []*agsv20250920.MetadataVar) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsMetadataObjectType), nil
	}

	models := make([]agsMetadataModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsMetadataModel{
			Name:  stringPointerValue(item.Name),
			Value: stringPointerValue(item.Value),
		})
	}

	return types.ListValueFrom(ctx, agsMetadataObjectType, models)
}

// uint64PointerValue converts an SDK uint64 pointer into an int64 attribute.
func uint64PointerValue(value *uint64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*value))
}

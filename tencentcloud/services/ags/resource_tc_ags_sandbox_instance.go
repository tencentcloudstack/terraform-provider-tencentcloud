package ags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

const (
	// AGS_SANDBOX_INSTANCE_STATUS_STOPPED and AGS_SANDBOX_INSTANCE_STATUS_STOP_FAILED
	// are the terminal states reported by DescribeSandboxInstances: a sandbox
	// instance in one of them is no longer running, so the resource is dropped
	// from state instead of being kept forever.
	AGS_SANDBOX_INSTANCE_STATUS_STOPPED     = "STOPPED"
	AGS_SANDBOX_INSTANCE_STATUS_STOP_FAILED = "STOP_FAILED"
)

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

// AgsSandboxInstanceResourceModel maps the schema attributes; each attribute
// appears exactly once.
//
// Timeout, SourceInstanceId, SnapshotId, MountOptions and CustomConfiguration
// are create-only inputs whose request form is not reported back by
// DescribeSandboxInstances (the timeout is reported in seconds through
// TimeoutSeconds and the nested blocks are answered with the cloud-side
// defaults and the storage values inherited from the tool), so Read keeps the
// configured values instead of rewriting them.
type AgsSandboxInstanceResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	InstanceId            types.String `tfsdk:"instance_id"`
	ToolId                types.String `tfsdk:"tool_id"`
	ToolName              types.String `tfsdk:"tool_name"`
	SourceInstanceId      types.String `tfsdk:"source_instance_id"`
	SnapshotId            types.String `tfsdk:"snapshot_id"`
	Timeout               types.String `tfsdk:"timeout"`
	MountOptions          types.List   `tfsdk:"mount_options"`
	CustomConfiguration   types.Object `tfsdk:"custom_configuration"`
	AuthMode              types.String `tfsdk:"auth_mode"`
	Metadata              types.List   `tfsdk:"metadata"`
	NetworkMode           types.String `tfsdk:"network_mode"`
	Persistent            types.Bool   `tfsdk:"persistent"`
	Status                types.String `tfsdk:"status"`
	TimeoutSeconds        types.Int64  `tfsdk:"timeout_seconds"`
	ExpiresAt             types.String `tfsdk:"expires_at"`
	StopReason            types.String `tfsdk:"stop_reason"`
	CreateTime            types.String `tfsdk:"create_time"`
	UpdateTime            types.String `tfsdk:"update_time"`
	ComputerConfiguration types.Object `tfsdk:"computer_configuration"`
}

type agsSandboxInstanceMountOptionModel struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SubPath   types.String `tfsdk:"sub_path"`
	ReadOnly  types.Bool   `tfsdk:"read_only"`
}

type agsSandboxInstanceMetadataModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

var agsSandboxInstanceMountOptionAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Name of the storage configuration in the sandbox tool.",
	},
	"mount_path": schema.StringAttribute{
		Optional:    true,
		Description: "Local mount path inside the sandbox instance. By default it inherits the mount path of the storage configuration in the tool.",
	},
	"sub_path": schema.StringAttribute{
		Optional:    true,
		Description: "Sub path of the storage mounted into the sandbox instance.",
	},
	"read_only": schema.BoolAttribute{
		Optional:    true,
		Description: "Read and write permission of the mount. By default it inherits the storage configuration in the tool.",
	},
}

var agsSandboxInstanceMetadataAttributes = map[string]schema.Attribute{
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
	agsSandboxInstanceMetadataAttrTypes  = attributeTypesOf(agsSandboxInstanceMetadataAttributes)
	agsSandboxInstanceMetadataObjectType = types.ObjectType{AttrTypes: agsSandboxInstanceMetadataAttrTypes}
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
				Description: "Unique identifier of the sandbox instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Id of the sandbox tool. At least one of `tool_id` and `tool_name` must be specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tool_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the sandbox tool. At least one of `tool_id` and `tool_name` must be specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "Id of the source sandbox instance to clone.",
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
				Description: "How long the sandbox instance keeps running before it is reclaimed automatically, for " +
					"example `5m`, `300s` or `1h`. Default value: `5m`. Minimum `30s`, maximum `24h`.",
			},
			"mount_options": schema.ListNestedAttribute{
				Optional:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsSandboxInstanceMountOptionAttributes},
				Description:  "Storage mount configuration of the sandbox instance.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"custom_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Attributes:  agsCustomConfigurationAttributes,
				Description: "Custom configuration of the sandbox instance.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"auth_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Authentication mode of the sandbox access. Valid values: `DEFAULT` (token " +
					"authentication, the default value), `TOKEN` (token authentication, every port access requires a " +
					"token), `NONE` (no authentication, no port access requires a token), `PUBLIC` (public mode, the " +
					"ENVD management port 49983 requires a token while other ports do not).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata": schema.ListNestedAttribute{
				Optional:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsSandboxInstanceMetadataAttributes},
				Description:  "Metadata of the sandbox instance.",
			},
			"network_mode": schema.StringAttribute{
				Computed: true,
				Description: "Network mode. Valid values: `PUBLIC` (public network access), `SANDBOX` (no network), " +
					"`INTERNAL_SERVICE` (Tencent Cloud internal public service). It can override the network " +
					"configuration of the tool level, but when a tool itself does not support VPC networks, the VPC " +
					"mode chosen for the instance is invalid.",
			},
			"persistent": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the sandbox instance is a persistent instance.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Status of the instance. Valid values: `STARTING` (starting), `RUNNING` (running), " +
					"`STOPPING` (stopping), `STOPPED` (stopped), `STOP_FAILED` (stop failed), `FAILED` (failed), " +
					"`RESIZING` (the specification is being updated).",
			},
			"timeout_seconds": schema.Int64Attribute{
				Computed:    true,
				Description: "Timeout in seconds. `null` means no timeout is configured.",
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "Expiration time in ISO 8601 format. `null` means the instance never expires.",
			},
			"stop_reason": schema.StringAttribute{
				Computed: true,
				Description: "Stop reason: `manual` (manual), `timeout` (timeout), `error` (error) or `system` " +
					"(system). It only has a value when the status is `STOPPED`, `STOP_FAILED` or `FAILED`. When the " +
					"provider fails to stop the instance, the status is `STOP_FAILED` and the reason is `error`.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time in ISO 8601 format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time in ISO 8601 format.",
			},
			"computer_configuration": schema.SingleNestedAttribute{
				Computed:    true,
				Attributes:  agsComputerConfigurationAttributes,
				Description: "Desktop environment configuration of the sandbox instance.",
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

	startSandboxInstanceRequest, diags := buildStartSandboxInstanceRequest(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	instanceId, err := service.StartSandboxInstance(ctx, startSandboxInstanceRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error starting AGS sandbox instance",
			fmt.Sprintf("StartSandboxInstance failed: %v", err),
		)
		return
	}

	// Never write an empty id into state.
	if instanceId == "" {
		resp.Diagnostics.AddError(
			"Error starting AGS sandbox instance",
			"StartSandboxInstance returned an empty instance id.",
		)
		return
	}

	plan.Id = types.StringValue(instanceId)
	plan.InstanceId = types.StringValue(instanceId)

	// Computed attributes are unknown at this point and a state cannot hold
	// unknown values, so they are normalized to null before the defensive
	// persist below.
	normalizeAgsSandboxInstanceComputedUnknowns(&plan)

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.describeAgsSandboxInstance(ctx, service, instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox instance after create",
			fmt.Sprintf("DescribeSandboxInstances failed for %s: %v", instanceId, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"AGS sandbox instance not readable after create",
			fmt.Sprintf("sandbox instance %s was started but cannot be read back", instanceId),
		)
		return
	}

	resp.Diagnostics.Append(refreshAgsSandboxInstanceModel(ctx, &plan, item)...)
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

	instanceId := agsSandboxInstanceIdOf(&state)
	if instanceId == "" {
		resp.Diagnostics.AddError(
			"Invalid AGS sandbox instance state",
			"The resource state does not contain a sandbox instance id.",
		)
		return
	}

	service := NewAgsService(r.Client())
	item, err := r.describeAgsSandboxInstance(ctx, service, instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox instance",
			fmt.Sprintf("DescribeSandboxInstances failed for %s: %v", instanceId, err),
		)
		return
	}

	// Confirmed "not found" on the cloud side (the instance was reclaimed after
	// its timeout or stopped) => drop it from state so Terraform plans a
	// recreation. The judgement stays INLINE in Read: a later unconditional
	// State.Set would overwrite the removal.
	if item == nil || agsSandboxInstanceStopped(item.Status) {
		resp.Diagnostics.AddWarning(
			"AGS sandbox instance not found",
			fmt.Sprintf("sandbox instance %s does not exist any more, removing it from state.", instanceId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshAgsSandboxInstanceModel(ctx, &state, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. UpdateSandboxInstance only accepts
// Timeout and Metadata; every other attribute is create-only and carries
// RequiresReplace in the schema.
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

	instanceId := agsSandboxInstanceIdOf(&state)
	if instanceId == "" {
		resp.Diagnostics.AddError(
			"Invalid AGS sandbox instance state",
			"The resource state does not contain a sandbox instance id.",
		)
		return
	}

	service := NewAgsService(r.Client())

	// Call the write API only when an attribute it covers changed.
	if !plan.Timeout.Equal(state.Timeout) || !plan.Metadata.Equal(state.Metadata) {
		updateSandboxInstanceRequest, diags := buildUpdateSandboxInstanceRequest(ctx, &plan, instanceId)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if err := service.UpdateSandboxInstance(ctx, updateSandboxInstanceRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error updating AGS sandbox instance",
				fmt.Sprintf("UpdateSandboxInstance failed for %s: %v", instanceId, err),
			)
			return
		}
	}

	item, err := r.describeAgsSandboxInstance(ctx, service, instanceId)
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
	resp.Diagnostics.Append(refreshAgsSandboxInstanceModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy` by stopping the sandbox instance.
func (r *AgsSandboxInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := agsSandboxInstanceIdOf(&state)
	if instanceId == "" {
		resp.Diagnostics.AddError(
			"Invalid AGS sandbox instance state",
			"The resource state does not contain a sandbox instance id.",
		)
		return
	}

	stopSandboxInstanceRequest := agsv20250920.NewStopSandboxInstanceRequest()
	stopSandboxInstanceRequest.InstanceId = fwHelper.String(instanceId)

	service := NewAgsService(r.Client())
	if err := service.StopSandboxInstance(ctx, stopSandboxInstanceRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error stopping AGS sandbox instance",
			fmt.Sprintf("StopSandboxInstance failed for %s: %v", instanceId, err),
		)
		return
	}
}

// buildStartSandboxInstanceRequest maps the plan onto a StartSandboxInstance
// request.
func buildStartSandboxInstanceRequest(ctx context.Context, plan *AgsSandboxInstanceResourceModel) (*agsv20250920.StartSandboxInstanceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	startSandboxInstanceRequest := agsv20250920.NewStartSandboxInstanceRequest()
	if !plan.ToolId.IsNull() && !plan.ToolId.IsUnknown() {
		startSandboxInstanceRequest.ToolId = plan.ToolId.ValueStringPointer()
	}
	if !plan.ToolName.IsNull() && !plan.ToolName.IsUnknown() {
		startSandboxInstanceRequest.ToolName = plan.ToolName.ValueStringPointer()
	}
	if !plan.SourceInstanceId.IsNull() && !plan.SourceInstanceId.IsUnknown() {
		startSandboxInstanceRequest.SourceInstanceId = plan.SourceInstanceId.ValueStringPointer()
	}
	if !plan.SnapshotId.IsNull() && !plan.SnapshotId.IsUnknown() {
		startSandboxInstanceRequest.SnapshotId = plan.SnapshotId.ValueStringPointer()
	}
	if !plan.Timeout.IsNull() && !plan.Timeout.IsUnknown() {
		startSandboxInstanceRequest.Timeout = plan.Timeout.ValueStringPointer()
	}
	if !plan.AuthMode.IsNull() && !plan.AuthMode.IsUnknown() {
		startSandboxInstanceRequest.AuthMode = plan.AuthMode.ValueStringPointer()
	}

	// client_token is generated internally as an idempotency key for create: the
	// value is fixed once per apply, so a retried StartSandboxInstance call (see
	// the retry loop in the service layer) cannot start the same sandbox
	// instance twice.
	startSandboxInstanceRequest.ClientToken = fwHelper.String(fwHelper.BuildUUID())

	mountOptions, mountDiags := expandAgsSandboxInstanceMountOptions(ctx, plan.MountOptions)
	diags.Append(mountDiags...)
	if diags.HasError() {
		return nil, diags
	}
	startSandboxInstanceRequest.MountOptions = mountOptions

	customConfiguration, customDiags := expandCustomConfiguration(ctx, plan.CustomConfiguration)
	diags.Append(customDiags...)
	if diags.HasError() {
		return nil, diags
	}
	startSandboxInstanceRequest.CustomConfiguration = customConfiguration

	metadata, metadataDiags := expandAgsSandboxInstanceMetadata(ctx, plan.Metadata)
	diags.Append(metadataDiags...)
	if diags.HasError() {
		return nil, diags
	}
	startSandboxInstanceRequest.Metadata = metadata

	return startSandboxInstanceRequest, diags
}

// buildUpdateSandboxInstanceRequest maps the plan onto an UpdateSandboxInstance
// request. The request struct omits empty values on the wire, so an attribute
// absent from the configuration is transmitted as "no change" and the cloud
// keeps the current value.
func buildUpdateSandboxInstanceRequest(ctx context.Context, plan *AgsSandboxInstanceResourceModel, instanceId string) (*agsv20250920.UpdateSandboxInstanceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	updateSandboxInstanceRequest := agsv20250920.NewUpdateSandboxInstanceRequest()
	updateSandboxInstanceRequest.InstanceId = fwHelper.String(instanceId)

	if !plan.Timeout.IsNull() && !plan.Timeout.IsUnknown() {
		updateSandboxInstanceRequest.Timeout = plan.Timeout.ValueStringPointer()
	}

	metadata, metadataDiags := expandAgsSandboxInstanceMetadata(ctx, plan.Metadata)
	diags.Append(metadataDiags...)
	if diags.HasError() {
		return nil, diags
	}
	updateSandboxInstanceRequest.Metadata = metadata

	return updateSandboxInstanceRequest, diags
}

// describeAgsSandboxInstance reads the sandbox instance with the given id and
// returns nil when the describe answer does not contain it.
func (r *AgsSandboxInstanceResource) describeAgsSandboxInstance(ctx context.Context, service AgsService, instanceId string) (*agsv20250920.SandboxInstance, error) {
	describeSandboxInstancesRequest := agsv20250920.NewDescribeSandboxInstancesRequest()
	describeSandboxInstancesRequest.InstanceIds = []*string{fwHelper.String(instanceId)}

	instances, err := service.DescribeSandboxInstances(ctx, describeSandboxInstancesRequest)
	if err != nil {
		return nil, err
	}

	return findAgsSandboxInstance(instances, instanceId), nil
}

// refreshAgsSandboxInstanceModel copies the cloud state into the model.
func refreshAgsSandboxInstanceModel(ctx context.Context, model *AgsSandboxInstanceResourceModel, item *agsv20250920.SandboxInstance) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		return diags
	}

	if item.InstanceId != nil && *item.InstanceId != "" {
		model.Id = types.StringValue(*item.InstanceId)
		model.InstanceId = types.StringValue(*item.InstanceId)
	}

	// tool_id, tool_name and auth_mode are user-configurable attributes reported
	// by the describe API, so a cloud-omitted value keeps the previous one.
	model.ToolId = fwHelper.StringValueOrPrev(item.ToolId, model.ToolId)
	model.ToolName = fwHelper.StringValueOrPrev(item.ToolName, model.ToolName)
	model.AuthMode = fwHelper.StringValueOrPrev(item.AuthMode, model.AuthMode)

	// The remaining attributes are read-only outputs, so the truthful refreshed
	// value (including null) is recorded as is.
	model.NetworkMode = fwHelper.StringPointerValueOrNull(item.NetworkMode)
	model.Persistent = fwHelper.BoolValueOrNull(item.Persistent)
	model.Status = fwHelper.StringPointerValueOrNull(item.Status)
	model.TimeoutSeconds = fwHelper.Int64ValueFromUint(item.TimeoutSeconds)
	model.ExpiresAt = fwHelper.StringPointerValueOrNull(item.ExpiresAt)
	model.StopReason = fwHelper.StringPointerValueOrNull(item.StopReason)
	model.CreateTime = fwHelper.StringPointerValueOrNull(item.CreateTime)
	model.UpdateTime = fwHelper.StringPointerValueOrNull(item.UpdateTime)

	computerConfiguration, computerDiags := flattenComputerConfiguration(ctx, item.ComputerConfiguration)
	diags.Append(computerDiags...)
	if diags.HasError() {
		return diags
	}
	model.ComputerConfiguration = computerConfiguration

	// metadata is the only configurable attribute the describe API reports back
	// in its request form; a cloud-omitted answer keeps the configured value.
	if item.Metadata != nil {
		metadata, metadataDiags := flattenAgsSandboxInstanceMetadata(ctx, item.Metadata)
		diags.Append(metadataDiags...)
		if diags.HasError() {
			return diags
		}
		model.Metadata = metadata
	}

	return diags
}

// normalizeAgsSandboxInstanceComputedUnknowns rewrites unknown attributes to
// null so the model can be persisted into state before the read-back fills them
// with the cloud-side values.
func normalizeAgsSandboxInstanceComputedUnknowns(model *AgsSandboxInstanceResourceModel) {
	if model.ToolId.IsUnknown() {
		model.ToolId = types.StringNull()
	}
	if model.ToolName.IsUnknown() {
		model.ToolName = types.StringNull()
	}
	if model.AuthMode.IsUnknown() {
		model.AuthMode = types.StringNull()
	}
	if model.NetworkMode.IsUnknown() {
		model.NetworkMode = types.StringNull()
	}
	if model.Persistent.IsUnknown() {
		model.Persistent = types.BoolNull()
	}
	if model.Status.IsUnknown() {
		model.Status = types.StringNull()
	}
	if model.TimeoutSeconds.IsUnknown() {
		model.TimeoutSeconds = types.Int64Null()
	}
	if model.ExpiresAt.IsUnknown() {
		model.ExpiresAt = types.StringNull()
	}
	if model.StopReason.IsUnknown() {
		model.StopReason = types.StringNull()
	}
	if model.CreateTime.IsUnknown() {
		model.CreateTime = types.StringNull()
	}
	if model.UpdateTime.IsUnknown() {
		model.UpdateTime = types.StringNull()
	}
	if model.ComputerConfiguration.IsUnknown() {
		model.ComputerConfiguration = types.ObjectNull(agsComputerConfigAttrTypes)
	}
}

// agsSandboxInstanceIdOf returns the identity stored in the model.
func agsSandboxInstanceIdOf(model *AgsSandboxInstanceResourceModel) string {
	if !model.InstanceId.IsNull() && !model.InstanceId.IsUnknown() && model.InstanceId.ValueString() != "" {
		return model.InstanceId.ValueString()
	}
	if !model.Id.IsNull() && !model.Id.IsUnknown() {
		return model.Id.ValueString()
	}

	return ""
}

// agsSandboxInstanceStopped reports whether the instance reached a terminal
// stopped state, which means it is no longer a running sandbox instance.
func agsSandboxInstanceStopped(status *string) bool {
	if status == nil {
		return false
	}

	return *status == AGS_SANDBOX_INSTANCE_STATUS_STOPPED || *status == AGS_SANDBOX_INSTANCE_STATUS_STOP_FAILED
}

// findAgsSandboxInstance picks this resource's own item from a describe list.
func findAgsSandboxInstance(instances []*agsv20250920.SandboxInstance, instanceId string) *agsv20250920.SandboxInstance {
	for _, instance := range instances {
		if instance != nil && instance.InstanceId != nil && *instance.InstanceId == instanceId {
			return instance
		}
	}

	return nil
}

// expandAgsSandboxInstanceMountOptions converts the mount_options attribute
// into the SDK pointer slice form. A null or unknown list is transmitted as
// "absent".
func expandAgsSandboxInstanceMountOptions(ctx context.Context, list types.List) ([]*agsv20250920.MountOption, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsSandboxInstanceMountOptionModel
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

// expandAgsSandboxInstanceMetadata converts the metadata attribute into the SDK
// pointer slice form. A null or unknown list is transmitted as "absent".
func expandAgsSandboxInstanceMetadata(ctx context.Context, list types.List) ([]*agsv20250920.MetadataVar, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsSandboxInstanceMetadataModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	metadatas := make([]*agsv20250920.MetadataVar, 0, len(models))
	for _, model := range models {
		metadata := &agsv20250920.MetadataVar{}
		if !model.Name.IsNull() && !model.Name.IsUnknown() {
			metadata.Name = model.Name.ValueStringPointer()
		}
		if !model.Value.IsNull() && !model.Value.IsUnknown() {
			metadata.Value = model.Value.ValueStringPointer()
		}
		metadatas = append(metadatas, metadata)
	}

	return metadatas, diags
}

// flattenAgsSandboxInstanceMetadata converts the SDK metadata list into the
// attribute form.
func flattenAgsSandboxInstanceMetadata(ctx context.Context, items []*agsv20250920.MetadataVar) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsSandboxInstanceMetadataObjectType), nil
	}

	models := make([]agsSandboxInstanceMetadataModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsSandboxInstanceMetadataModel{
			Name:  stringPointerValue(item.Name),
			Value: stringPointerValue(item.Value),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsSandboxInstanceMetadataObjectType), nil
	}

	return types.ListValueFrom(ctx, agsSandboxInstanceMetadataObjectType, models)
}

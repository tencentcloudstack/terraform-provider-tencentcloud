package ags

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

const (
	// agsSandboxInstanceStatusStopped and agsSandboxInstanceStatusStopFailed are
	// the terminal states reported by DescribeSandboxInstances: a sandbox
	// instance in one of them is no longer running, so the resource is dropped
	// from state instead of being kept forever.
	agsSandboxInstanceStatusStopped    = "STOPPED"
	agsSandboxInstanceStatusStopFailed = "STOP_FAILED"

	// agsSandboxInstanceReservedMetadataPrefix marks the metadata names the
	// system owns. The cloud side adds such entries on its own (for example
	// x-ags-precreate-match=true when an instance is started), so they are never
	// part of the configuration: they are filtered out of the state and sent
	// back only when an update replaces the whole metadata list.
	agsSandboxInstanceReservedMetadataPrefix = "x-ags-"
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
// are the create-only inputs. The describe API reports the timeout in seconds
// only and answers the nested blocks in normalized form (the cloud-side
// defaults, the quantities rewritten as 1000m/1024Mi and the storage values
// inherited from the tool), so create, update and read all keep the spelling of
// the configuration and the semantic RequiresReplaceIf modifiers stop that
// normalized answer from forcing a replacement.
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
	ExpiresAt             types.String `tfsdk:"expires_at"`
	StopReason            types.String `tfsdk:"stop_reason"`
	CreateTime            types.String `tfsdk:"create_time"`
	UpdateTime            types.String `tfsdk:"update_time"`
	ComputerConfiguration types.Object `tfsdk:"computer_configuration"`
}

// agsSandboxInstanceMountOptionModel maps one mount_options element.
type agsSandboxInstanceMountOptionModel struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SubPath   types.String `tfsdk:"sub_path"`
	ReadOnly  types.Bool   `tfsdk:"read_only"`
}

// agsSandboxInstanceMetadataModel maps one metadata element.
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
		Optional: true,
		Description: "Local mount path inside the sandbox instance. By default it inherits the mount " +
			"path of the storage configuration in the tool.",
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
		Optional: true,
		Description: "Metadata name. The `x-ags-` prefix is reserved by the system and must not be used " +
			"here, otherwise a permanent diff is reported.",
	},
	"value": schema.StringAttribute{
		Optional:    true,
		Description: "Metadata value.",
	},
}

var (
	agsSandboxInstanceMountOptionAttrTypes  = agsSandboxInstanceAttributeTypes(agsSandboxInstanceMountOptionAttributes)
	agsSandboxInstanceMountOptionObjectType = types.ObjectType{AttrTypes: agsSandboxInstanceMountOptionAttrTypes}
	agsSandboxInstanceMetadataAttrTypes     = agsSandboxInstanceAttributeTypes(agsSandboxInstanceMetadataAttributes)
	agsSandboxInstanceMetadataObjectType    = types.ObjectType{AttrTypes: agsSandboxInstanceMetadataAttrTypes}
)

// The custom configuration and the computer configuration of a sandbox instance
// are also used by the sandbox tool resource. Their definitions are duplicated
// here on purpose: this resource must not depend on
// resource_tc_ags_sandbox_tool.go.

type agsSandboxInstanceCustomConfigurationModel struct {
	Image             types.String `tfsdk:"image"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	Command           types.Set    `tfsdk:"command"`
	Args              types.List   `tfsdk:"args"`
	Env               types.List   `tfsdk:"env"`
	Ports             types.List   `tfsdk:"ports"`
	Resources         types.Object `tfsdk:"resources"`
	Probe             types.Object `tfsdk:"probe"`
	DnsConfig         types.Object `tfsdk:"dns_config"`
	SecurityContext   types.Object `tfsdk:"security_context"`
}

type agsSandboxInstanceEnvVarModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type agsSandboxInstancePortModel struct {
	Name     types.String `tfsdk:"name"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

type agsSandboxInstanceResourceModel struct {
	Cpu     types.String `tfsdk:"cpu"`
	Memory  types.String `tfsdk:"memory"`
	Storage types.String `tfsdk:"storage"`
}

type agsSandboxInstanceProbeModel struct {
	HttpGet          types.Object `tfsdk:"http_get"`
	ReadyTimeoutMs   types.Int64  `tfsdk:"ready_timeout_ms"`
	ProbeTimeoutMs   types.Int64  `tfsdk:"probe_timeout_ms"`
	ProbePeriodMs    types.Int64  `tfsdk:"probe_period_ms"`
	SuccessThreshold types.Int64  `tfsdk:"success_threshold"`
	FailureThreshold types.Int64  `tfsdk:"failure_threshold"`
}

type agsSandboxInstanceHttpGetModel struct {
	Path   types.String `tfsdk:"path"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`
}

type agsSandboxInstanceDnsConfigModel struct {
	Servers  types.Set `tfsdk:"servers"`
	Searches types.Set `tfsdk:"searches"`
	Options  types.Set `tfsdk:"options"`
}

type agsSandboxInstanceSecurityContextModel struct {
	Privileged   types.Bool   `tfsdk:"privileged"`
	Capabilities types.Object `tfsdk:"capabilities"`
}

type agsSandboxInstanceCapabilitiesModel struct {
	Add  types.Set `tfsdk:"add"`
	Drop types.Set `tfsdk:"drop"`
}

type agsSandboxInstanceComputerConfigurationModel struct {
	WaaConfiguration     types.Object `tfsdk:"waa_configuration"`
	OsWorldConfiguration types.Object `tfsdk:"os_world_configuration"`
}

type agsSandboxInstanceWaaConfigurationModel struct {
	ImageId types.String `tfsdk:"image_id"`
}

type agsSandboxInstanceOsWorldConfigurationModel struct {
	Version types.String `tfsdk:"version"`
}

var agsSandboxInstanceEnvVarAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Environment variable name.",
	},
	"value": schema.StringAttribute{
		Optional:    true,
		Description: "Environment variable value.",
	},
}

var agsSandboxInstancePortAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Port name.",
	},
	"port": schema.Int64Attribute{
		Optional:    true,
		Description: "Port number.",
	},
	"protocol": schema.StringAttribute{
		Optional:    true,
		Description: "Protocol.",
	},
}

var agsSandboxInstanceResourceAttributes = map[string]schema.Attribute{
	"cpu": schema.StringAttribute{
		Optional:    true,
		Description: "Amount of cpu resources.",
	},
	"memory": schema.StringAttribute{
		Optional:    true,
		Description: "Amount of memory resources.",
	},
	"storage": schema.StringAttribute{
		Optional:    true,
		Description: "Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`.",
	},
}

var agsSandboxInstanceHttpGetAttributes = map[string]schema.Attribute{
	"path": schema.StringAttribute{
		Optional:    true,
		Description: "Request path.",
	},
	"port": schema.Int64Attribute{
		Optional:    true,
		Description: "Request port.",
	},
	"scheme": schema.StringAttribute{
		Optional:    true,
		Description: "Request scheme.",
	},
}

var agsSandboxInstanceProbeAttributes = map[string]schema.Attribute{
	"http_get": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsSandboxInstanceHttpGetAttributes,
		Description: "HTTP GET probe configuration.",
	},
	"ready_timeout_ms": schema.Int64Attribute{
		Optional:    true,
		Description: "Timeout in milliseconds for the probe to become ready.",
	},
	"probe_timeout_ms": schema.Int64Attribute{
		Optional:    true,
		Description: "Timeout in milliseconds of a single probe.",
	},
	"probe_period_ms": schema.Int64Attribute{
		Optional:    true,
		Description: "Interval in milliseconds between two probes.",
	},
	"success_threshold": schema.Int64Attribute{
		Optional:    true,
		Description: "Number of consecutive successes required to mark the probe as healthy.",
	},
	"failure_threshold": schema.Int64Attribute{
		Optional:    true,
		Description: "Number of consecutive failures required to mark the probe as unhealthy.",
	},
}

var agsSandboxInstanceDnsConfigAttributes = map[string]schema.Attribute{
	"servers": schema.SetAttribute{
		Required:    true,
		ElementType: types.StringType,
		Description: "DNS server addresses. Valid IP addresses are required. Default value: `10.0.0.1`.",
	},
	"searches": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Search domains, corresponding to the `search` directive of resolv.conf.",
	},
	"options": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Options, corresponding to the options of resolv.conf.",
	},
}

var agsSandboxInstanceCapabilitiesAttributes = map[string]schema.Attribute{
	"add": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Capabilities to add. Reference value: `[\"SYS_ADMIN\"]`.",
	},
	"drop": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Capabilities to drop. Reference value: `[\"ALL\"]`.",
	},
}

var agsSandboxInstanceSecurityContextAttributes = map[string]schema.Attribute{
	"privileged": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Whether the sandbox runs in privileged mode. Default value: `true`.",
	},
	"capabilities": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsSandboxInstanceCapabilitiesAttributes,
		Description: "POSIX capabilities that can be added to or removed from the running container.",
	},
}

var agsSandboxInstanceCustomConfigurationAttributes = map[string]schema.Attribute{
	"image": schema.StringAttribute{
		Optional:    true,
		Description: "Image address.",
	},
	"image_registry_type": schema.StringAttribute{
		Optional:    true,
		Description: "Image registry type. Valid values: `enterprise`, `personal`, `custom`.",
	},
	"command": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Startup command.",
	},
	"args": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Startup arguments. The order is preserved.",
	},
	"env": schema.ListNestedAttribute{
		Optional:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: agsSandboxInstanceEnvVarAttributes},
		Description:  "Environment variables.",
	},
	"ports": schema.ListNestedAttribute{
		Optional:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: agsSandboxInstancePortAttributes},
		Description:  "Port configuration.",
	},
	"resources": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsSandboxInstanceResourceAttributes,
		Description: "Resource configuration.",
	},
	"probe": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsSandboxInstanceProbeAttributes,
		Description: "Probe configuration.",
	},
	"dns_config": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsSandboxInstanceDnsConfigAttributes,
		Description: "Sandbox DNS configuration.",
	},
	"security_context": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsSandboxInstanceSecurityContextAttributes,
		Description: "Sandbox security context.",
	},
}

var agsSandboxInstanceWaaConfigurationAttributes = map[string]schema.Attribute{
	"image_id": schema.StringAttribute{
		Optional:    true,
		Description: "Custom WAA image id.",
	},
}

var agsSandboxInstanceOsWorldConfigurationAttributes = map[string]schema.Attribute{
	"version": schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Built-in OSWorld version. Valid values: `osworld1`, `osworld2`. Default value: `osworld1`.",
	},
}

var agsSandboxInstanceComputerConfigurationAttributes = map[string]schema.Attribute{
	"waa_configuration": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsSandboxInstanceWaaConfigurationAttributes,
		Description: "WAA sandbox configuration.",
	},
	"os_world_configuration": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsSandboxInstanceOsWorldConfigurationAttributes,
		Description: "Built-in OSWorld configuration.",
	},
}

var (
	agsSandboxInstanceEnvVarAttrTypes               = agsSandboxInstanceAttributeTypes(agsSandboxInstanceEnvVarAttributes)
	agsSandboxInstancePortAttrTypes                 = agsSandboxInstanceAttributeTypes(agsSandboxInstancePortAttributes)
	agsSandboxInstanceResourceAttrTypes             = agsSandboxInstanceAttributeTypes(agsSandboxInstanceResourceAttributes)
	agsSandboxInstanceHttpGetAttrTypes              = agsSandboxInstanceAttributeTypes(agsSandboxInstanceHttpGetAttributes)
	agsSandboxInstanceProbeAttrTypes                = agsSandboxInstanceAttributeTypes(agsSandboxInstanceProbeAttributes)
	agsSandboxInstanceDnsConfigAttrTypes            = agsSandboxInstanceAttributeTypes(agsSandboxInstanceDnsConfigAttributes)
	agsSandboxInstanceCapabilitiesAttrTypes         = agsSandboxInstanceAttributeTypes(agsSandboxInstanceCapabilitiesAttributes)
	agsSandboxInstanceSecurityContextAttrTypes      = agsSandboxInstanceAttributeTypes(agsSandboxInstanceSecurityContextAttributes)
	agsSandboxInstanceCustomConfigAttrTypes         = agsSandboxInstanceAttributeTypes(agsSandboxInstanceCustomConfigurationAttributes)
	agsSandboxInstanceWaaConfigurationAttrTypes     = agsSandboxInstanceAttributeTypes(agsSandboxInstanceWaaConfigurationAttributes)
	agsSandboxInstanceOsWorldConfigurationAttrTypes = agsSandboxInstanceAttributeTypes(agsSandboxInstanceOsWorldConfigurationAttributes)
	agsSandboxInstanceComputerConfigAttrTypes       = agsSandboxInstanceAttributeTypes(agsSandboxInstanceComputerConfigurationAttributes)

	agsSandboxInstanceEnvVarObjectType = types.ObjectType{AttrTypes: agsSandboxInstanceEnvVarAttrTypes}
	agsSandboxInstancePortObjectType   = types.ObjectType{AttrTypes: agsSandboxInstancePortAttrTypes}
)

// agsSandboxInstanceAttributeTypes mirrors the value types of a nested attribute
// map so the conversion helpers stay in sync with the schema definition.
func agsSandboxInstanceAttributeTypes(attributes map[string]schema.Attribute) map[string]attr.Type {
	result := make(map[string]attr.Type, len(attributes))
	for name, attribute := range attributes {
		result[name] = attribute.GetType()
	}
	return result
}
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
			"tool_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Id of the sandbox tool. At least one of `tool_id` and `tool_name` must be specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"tool_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the sandbox tool. At least one of `tool_id` and `tool_name` must be specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
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
				Computed: true,
				Description: "How long the sandbox instance keeps running before it is reclaimed automatically, for " +
					"example `5m`, `300s` or `1h`. Default value: `5m`. Minimum `30s`, maximum `24h`. It is " +
					"refreshed from the cloud side timeout in seconds on read.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mount_options": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Storage mount configuration of the sandbox instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: agsSandboxInstanceMountOptionAttributes,
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIf(
						agsSandboxInstanceMountOptionsChanged,
						"If the storage mount configuration changes, Terraform will destroy and recreate the "+
							"resource. The normalized cloud answer does not count as a change.",
						"If the storage mount configuration changes, Terraform will destroy and recreate the resource.",
					),
				},
			},
			"custom_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Custom configuration of the sandbox instance.",
				Attributes:  agsSandboxInstanceCustomConfigurationAttributes,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplaceIf(
						agsSandboxInstanceCustomConfigurationChanged,
						"If the custom configuration changes, Terraform will destroy and recreate the resource. "+
							"The normalized cloud answer does not count as a change.",
						"If the custom configuration changes, Terraform will destroy and recreate the resource.",
					),
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
				Optional: true,
				Description: "Metadata of the sandbox instance. The `x-ags-` name prefix is reserved by the " +
					"system: such entries are added by the cloud side, are not tracked in the state and must not " +
					"be configured, otherwise a permanent diff is reported.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: agsSandboxInstanceMetadataAttributes,
				},
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the sandbox instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "Expiration time in ISO 8601 format. `null` means the instance never expires.",
			},
			"stop_reason": schema.StringAttribute{
				Computed: true,
				Description: "Stop reason: `manual` (manual), `timeout` (timeout), `error` (error) or `system` " +
					"(system). It only has a value when the status is `STOPPED`, `STOP_FAILED` or `FAILED`. When the " +
					"provider fails to stop the instance, the status is `STOP_FAILED` and the reason is `error`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Attributes:  agsSandboxInstanceComputerConfigurationAttributes,
				Description: "Desktop environment configuration of the sandbox instance.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
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

	// The create-only inputs are adopted only when the state carries no value
	// yet (an import) or the cloud answer is not equivalent to it (a real
	// drift), so create, update and read all keep the configuration spelling.
	// The cloud answers them in a normalized form (cpu 1 -> 1000m, memory
	// 1Gi -> 1024Mi, a default probe is added, the mount values are inherited
	// from the tool); a state rewritten with that form is different from the
	// configuration and makes every plan an in-place update.
	customConfiguration, customDiags := flattenAgsSandboxInstanceCustomConfiguration(ctx, item.CustomConfiguration)
	resp.Diagnostics.Append(customDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.CustomConfiguration = agsAdoptEquivalent(state.CustomConfiguration, customConfiguration)

	mountOptions, mountOptionDiags := flattenAgsSandboxInstanceMountOptions(ctx, item.MountOptions)
	resp.Diagnostics.Append(mountOptionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.MountOptions = agsAdoptEquivalent(state.MountOptions, mountOptions)

	// timeout is answered in seconds only and is managed through the
	// configuration, so it is rendered back only when the state carries no value
	// yet (an import, or a state written before the attribute was refreshed).
	if state.Timeout.IsNull() || state.Timeout.IsUnknown() {
		state.Timeout = agsSandboxInstanceTimeoutValue(item.TimeoutSeconds)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. UpdateSandboxInstance only accepts
// Timeout and Metadata; every other attribute is create-only and carries a
// RequiresReplaceIf modifier in the schema.
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
		// UpdateSandboxInstance replaces the whole metadata list, while the
		// entries owned by the system (x-ags- prefix) are filtered out of the
		// state. They are read back from the cloud and sent again, so an update
		// never drops a system marker the configuration cannot express.
		var reservedMetadata []*agsv20250920.MetadataVar
		if !plan.Metadata.IsNull() && !plan.Metadata.IsUnknown() {
			current, err := r.describeAgsSandboxInstance(ctx, service, instanceId)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error reading AGS sandbox instance before update",
					fmt.Sprintf("DescribeSandboxInstances failed for %s: %v", instanceId, err),
				)
				return
			}
			if current != nil {
				reservedMetadata = reservedAgsSandboxInstanceMetadata(current.Metadata)
			}
		}

		updateSandboxInstanceRequest, diags := buildUpdateSandboxInstanceRequest(ctx, &plan, instanceId, reservedMetadata)
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

	customConfiguration, customDiags := expandAgsSandboxInstanceCustomConfiguration(ctx, plan.CustomConfiguration)
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
//
// reservedMetadata carries the system owned entries (x-ags- prefix) read back
// from the cloud: the request replaces the whole list, so they are appended to
// the configured entries to keep the cloud side bookkeeping intact.
func buildUpdateSandboxInstanceRequest(ctx context.Context, plan *AgsSandboxInstanceResourceModel, instanceId string, reservedMetadata []*agsv20250920.MetadataVar) (*agsv20250920.UpdateSandboxInstanceRequest, diag.Diagnostics) {
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
	updateSandboxInstanceRequest.Metadata = mergeAgsSandboxInstanceReservedMetadata(metadata, reservedMetadata)

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
	model.ExpiresAt = fwHelper.StringPointerValueOrNull(item.ExpiresAt)
	model.StopReason = fwHelper.StringPointerValueOrNull(item.StopReason)
	model.CreateTime = fwHelper.StringPointerValueOrNull(item.CreateTime)
	model.UpdateTime = fwHelper.StringPointerValueOrNull(item.UpdateTime)

	computerConfiguration, computerDiags := flattenAgsSandboxInstanceComputerConfiguration(ctx, item.ComputerConfiguration)
	diags.Append(computerDiags...)
	if diags.HasError() {
		return diags
	}
	model.ComputerConfiguration = computerConfiguration

	// metadata is the only configurable attribute the describe API reports back
	// in its request form; a cloud-omitted answer keeps the configured value. The
	// system owned entries (x-ags- prefix) are filtered out by the flatten helper,
	// so the state keeps being written from the configuration only.
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
	if model.Timeout.IsUnknown() {
		model.Timeout = types.StringNull()
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
		model.ComputerConfiguration = types.ObjectNull(agsSandboxInstanceComputerConfigAttrTypes)
	}
	// mount_options and custom_configuration are optional and computed: when the
	// configuration leaves them unset the plan value is unknown, which a state
	// cannot hold.
	if model.MountOptions.IsUnknown() {
		model.MountOptions = types.ListNull(agsSandboxInstanceMountOptionObjectType)
	}
	if model.CustomConfiguration.IsUnknown() {
		model.CustomConfiguration = types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes)
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

	if *status == agsSandboxInstanceStatusStopped {
		return true
	}

	return *status == agsSandboxInstanceStatusStopFailed
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
// attribute form. The entries owned by the system (x-ags- prefix) are dropped:
// metadata is an Optional-only attribute, so the state has to stay equal to the
// plan and the cloud side entries would otherwise make Terraform fail the apply
// with "new element N has appeared".
func flattenAgsSandboxInstanceMetadata(ctx context.Context, items []*agsv20250920.MetadataVar) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsSandboxInstanceMetadataObjectType), nil
	}

	models := make([]agsSandboxInstanceMetadataModel, 0, len(items))
	for _, item := range items {
		if item == nil || isAgsSandboxInstanceReservedMetadataName(item.Name) {
			continue
		}
		models = append(models, agsSandboxInstanceMetadataModel{
			Name:  fwHelper.StringPointerValueOrNull(item.Name),
			Value: fwHelper.StringPointerValueOrNull(item.Value),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsSandboxInstanceMetadataObjectType), nil
	}

	return types.ListValueFrom(ctx, agsSandboxInstanceMetadataObjectType, models)
}

// isAgsSandboxInstanceReservedMetadataName reports whether the metadata name is
// owned by the system. Such entries are created by the cloud side and can never
// be expressed through the configuration, so they are excluded from the state
// and only passed through on update.
func isAgsSandboxInstanceReservedMetadataName(name *string) bool {
	if name == nil {
		return false
	}

	return strings.HasPrefix(strings.ToLower(*name), agsSandboxInstanceReservedMetadataPrefix)
}

// reservedAgsSandboxInstanceMetadata keeps the system owned entries of the cloud
// answer, so an update can send them back.
func reservedAgsSandboxInstanceMetadata(items []*agsv20250920.MetadataVar) []*agsv20250920.MetadataVar {
	var reserved []*agsv20250920.MetadataVar
	for _, item := range items {
		if item == nil || !isAgsSandboxInstanceReservedMetadataName(item.Name) {
			continue
		}
		reserved = append(reserved, item)
	}

	return reserved
}

// mergeAgsSandboxInstanceReservedMetadata appends the system owned entries to
// the configured ones, skipping the names the configuration already carries. A
// nil configured list stays nil: metadata absent from the plan is transmitted as
// "no change" and must not be turned into an update.
func mergeAgsSandboxInstanceReservedMetadata(configured, reserved []*agsv20250920.MetadataVar) []*agsv20250920.MetadataVar {
	if len(configured) == 0 || len(reserved) == 0 {
		return configured
	}

	names := make(map[string]struct{}, len(configured))
	for _, item := range configured {
		if item != nil && item.Name != nil {
			names[*item.Name] = struct{}{}
		}
	}

	merged := make([]*agsv20250920.MetadataVar, len(configured), len(configured)+len(reserved))
	copy(merged, configured)
	for _, item := range reserved {
		if item == nil || item.Name == nil {
			continue
		}
		if _, ok := names[*item.Name]; ok {
			continue
		}
		names[*item.Name] = struct{}{}
		merged = append(merged, item)
	}

	return merged
}

// flattenAgsSandboxInstanceMountOptions converts the SDK mount option list into
// the attribute form. It is used by Read only.
func flattenAgsSandboxInstanceMountOptions(ctx context.Context, items []*agsv20250920.MountOption) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsSandboxInstanceMountOptionObjectType), nil
	}

	models := make([]agsSandboxInstanceMountOptionModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsSandboxInstanceMountOptionModel{
			Name:      fwHelper.StringPointerValueOrNull(item.Name),
			MountPath: fwHelper.StringPointerValueOrNull(item.MountPath),
			SubPath:   fwHelper.StringPointerValueOrNull(item.SubPath),
			ReadOnly:  fwHelper.BoolValueOrNull(item.ReadOnly),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsSandboxInstanceMountOptionObjectType), nil
	}

	return types.ListValueFrom(ctx, agsSandboxInstanceMountOptionObjectType, models)
}

// agsSandboxInstanceMountOptionsChanged reports whether the storage mount
// configuration really changed. The cloud answers the mount values in a
// normalized form (the defaults and the storage inherited from the tool), so
// only a difference the configuration itself spells out requires a replacement.
func agsSandboxInstanceMountOptionsChanged(_ context.Context, req planmodifier.ListRequest, resp *listplanmodifier.RequiresReplaceIfFuncResponse) {
	resp.RequiresReplace = agsSandboxInstanceCreateOnlyInputChanged(req.ConfigValue, req.StateValue)
}

// agsSandboxInstanceCustomConfigurationChanged reports whether the custom
// configuration really changed. The cloud answers it in a normalized form
// (cpu 1 -> 1000m, memory 1Gi -> 1024Mi, a default probe and the storage
// inherited from the tool), so only a difference the configuration itself
// spells out requires a replacement.
func agsSandboxInstanceCustomConfigurationChanged(_ context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
	resp.RequiresReplace = agsSandboxInstanceCreateOnlyInputChanged(req.ConfigValue, req.StateValue)
}

// agsSandboxInstanceCreateOnlyInputChanged reports whether the configuration of
// a create-only input differs from the state in a way the cloud side cannot
// apply in place.
func agsSandboxInstanceCreateOnlyInputChanged(config, state attr.Value) bool {
	if config.IsNull() || config.IsUnknown() {
		return false
	}
	if state.IsNull() || state.IsUnknown() {
		return true
	}
	return !agsNestedCompatible(config, state)
}

// agsAdoptEquivalent returns the cloud answer only when the state
// carries no value yet (an import) or is not equivalent to it (a real drift),
// and keeps the state value otherwise. Create, update and read all keep the
// configuration spelling of the create-only inputs: a state rewritten with the
// normalized cloud form is different from the configuration, and Terraform
// compares the two before the plan modifiers run, so it turns every plan into an
// in-place update where only the computed attributes are reported as changing.
func agsAdoptEquivalent[T attr.Value](current, cloud T) T {
	if current.IsNull() || current.IsUnknown() {
		return cloud
	}
	if cloud.IsNull() || cloud.IsUnknown() {
		return cloud
	}
	if agsNestedCompatible(current, cloud) {
		return current
	}
	return cloud
}

// agsNestedCompatible reports whether every value the
// configuration spells out already matches the state.
func agsNestedCompatible(config, state attr.Value) bool {
	return agsValueCompatible("", config, state)
}

// agsValueCompatible reports whether every value the
// configuration spells out already matches the state. A null configuration
// attribute means "not configured" and always matches, which is what makes the
// check usable against a normalized cloud answer. The name is used to compare
// quantity attributes (cpu, memory, storage) by their value instead of their
// textual form.
func agsValueCompatible(name string, config, state attr.Value) bool {
	if config.IsNull() || config.IsUnknown() {
		return true
	}
	if state.IsUnknown() {
		return true
	}
	if state.IsNull() {
		return false
	}

	switch configValue := config.(type) {
	case types.String:
		stateValue, ok := state.(types.String)
		if !ok {
			return false
		}
		if name == "cpu" || name == "memory" || name == "storage" {
			return agsQuantityEqual(name, configValue.ValueString(), stateValue.ValueString())
		}
		return configValue.ValueString() == stateValue.ValueString()
	case types.Bool:
		stateValue, ok := state.(types.Bool)
		return ok && configValue.ValueBool() == stateValue.ValueBool()
	case types.Int64:
		stateValue, ok := state.(types.Int64)
		return ok && configValue.ValueInt64() == stateValue.ValueInt64()
	case types.Object:
		stateValue, ok := state.(types.Object)
		if !ok {
			return false
		}
		stateAttributes := stateValue.Attributes()
		for attributeName, configAttribute := range configValue.Attributes() {
			stateAttribute, exists := stateAttributes[attributeName]
			if !exists {
				return false
			}
			if !agsValueCompatible(attributeName, configAttribute, stateAttribute) {
				return false
			}
		}
		return true
	case types.List:
		stateValue, ok := state.(types.List)
		if !ok {
			return false
		}
		configElements := configValue.Elements()
		stateElements := stateValue.Elements()
		if len(configElements) != len(stateElements) {
			return false
		}
		for index, configElement := range configElements {
			if !agsValueCompatible("", configElement, stateElements[index]) {
				return false
			}
		}
		return true
	case types.Set:
		stateValue, ok := state.(types.Set)
		if !ok {
			return false
		}
		return configValue.Equal(stateValue)
	}

	return config.Equal(state)
}

// agsQuantityEqual compares two Kubernetes style quantities. The
// AGS API normalizes what it is given, so a textual comparison would report a
// change for an unchanged value.
func agsQuantityEqual(name, config, state string) bool {
	if config == state {
		return true
	}

	if name == "cpu" {
		configMilliCores, configOk := agsParseCPU(config)
		stateMilliCores, stateOk := agsParseCPU(state)

		return configOk && stateOk && configMilliCores == stateMilliCores
	}

	configBytes, configOk := agsParseBytes(config)
	stateBytes, stateOk := agsParseBytes(state)

	return configOk && stateOk && configBytes == stateBytes
}

// agsParseCPU converts a cpu quantity into millicores.
func agsParseCPU(value string) (int64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}

	if strings.HasSuffix(value, "m") {
		milliCores, err := strconv.ParseInt(strings.TrimSuffix(value, "m"), 10, 64)
		if err != nil {
			return 0, false
		}

		return milliCores, true
	}

	cores, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}

	return int64(math.Round(cores * 1000)), true
}

// agsParseBytes converts a memory or storage quantity into bytes.
func agsParseBytes(value string) (int64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}

	suffixes := []struct {
		suffix string
		factor float64
	}{
		{"Ki", 1 << 10}, {"Mi", 1 << 20}, {"Gi", 1 << 30}, {"Ti", 1 << 40}, {"Pi", 1 << 50},
		{"K", 1e3}, {"M", 1e6}, {"G", 1e9}, {"T", 1e12}, {"P", 1e15},
	}
	for _, item := range suffixes {
		if !strings.HasSuffix(value, item.suffix) {
			continue
		}
		number, err := strconv.ParseFloat(strings.TrimSuffix(value, item.suffix), 64)
		if err != nil {
			return 0, false
		}

		return int64(math.Round(number * item.factor)), true
	}

	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}

	return int64(math.Round(number)), true
}

// agsSandboxInstanceTimeoutValue renders the cloud side timeout in seconds in the
// duration form accepted by the configuration. A missing value means the
// instance never expires.
func agsSandboxInstanceTimeoutValue(seconds *uint64) types.String {
	if seconds == nil {
		return types.StringNull()
	}

	switch total := *seconds; {
	case total == 0:
		return types.StringValue("0s")
	case total%3600 == 0:
		return types.StringValue(fmt.Sprintf("%dh", total/3600))
	case total%60 == 0:
		return types.StringValue(fmt.Sprintf("%dm", total/60))
	default:
		return types.StringValue(fmt.Sprintf("%ds", total))
	}
}
func expandAgsSandboxInstanceCustomConfiguration(ctx context.Context, object types.Object) (*agsv20250920.CustomConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsSandboxInstanceCustomConfigurationModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	customConfiguration := &agsv20250920.CustomConfiguration{}
	if !model.Image.IsNull() && !model.Image.IsUnknown() {
		customConfiguration.Image = model.Image.ValueStringPointer()
	}
	if !model.ImageRegistryType.IsNull() && !model.ImageRegistryType.IsUnknown() {
		customConfiguration.ImageRegistryType = model.ImageRegistryType.ValueStringPointer()
	}

	command, commandDiags := agsSandboxInstanceExpandStringSet(ctx, model.Command)
	diags.Append(commandDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Command = command

	args, argsDiags := agsSandboxInstanceExpandStringList(ctx, model.Args)
	diags.Append(argsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Args = args

	env, envDiags := expandAgsSandboxInstanceEnvVars(ctx, model.Env)
	diags.Append(envDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Env = env

	ports, portDiags := expandAgsSandboxInstancePorts(ctx, model.Ports)
	diags.Append(portDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Ports = ports

	resources, resourceDiags := expandAgsSandboxInstanceResources(ctx, model.Resources)
	diags.Append(resourceDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Resources = resources

	probe, probeDiags := expandAgsSandboxInstanceProbe(ctx, model.Probe)
	diags.Append(probeDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Probe = probe

	dnsConfig, dnsDiags := expandAgsSandboxInstanceDnsConfig(ctx, model.DnsConfig)
	diags.Append(dnsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.DNSConfig = dnsConfig

	securityContext, securityDiags := expandAgsSandboxInstanceSecurityContext(ctx, model.SecurityContext)
	diags.Append(securityDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.SecurityContext = securityContext

	return customConfiguration, diags
}

func flattenAgsSandboxInstanceCustomConfiguration(ctx context.Context, item *agsv20250920.CustomConfigurationDetail) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes), nil
	}

	var diags diag.Diagnostics

	command, commandDiags := agsSandboxInstanceFlattenStringSet(ctx, item.Command)
	diags.Append(commandDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes), diags
	}

	args, argsDiags := agsSandboxInstanceFlattenStringList(ctx, item.Args)
	diags.Append(argsDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes), diags
	}

	env, envDiags := flattenAgsSandboxInstanceEnvVars(ctx, item.Env)
	diags.Append(envDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes), diags
	}

	ports, portDiags := flattenAgsSandboxInstancePorts(ctx, item.Ports)
	diags.Append(portDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes), diags
	}

	resources, resourceDiags := flattenAgsSandboxInstanceResources(ctx, item.Resources)
	diags.Append(resourceDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes), diags
	}

	probe, probeDiags := flattenAgsSandboxInstanceProbe(ctx, item.Probe)
	diags.Append(probeDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes), diags
	}

	dnsConfig, dnsDiags := flattenAgsSandboxInstanceDnsConfig(ctx, item.DNSConfig)
	diags.Append(dnsDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes), diags
	}

	securityContext, securityDiags := flattenAgsSandboxInstanceSecurityContext(ctx, item.SecurityContext)
	diags.Append(securityDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCustomConfigAttrTypes), diags
	}

	// ImageDigest is a read-only output of CustomConfigurationDetail and is not
	// part of the configurable schema, so it is not written back.
	model := agsSandboxInstanceCustomConfigurationModel{
		Image:             fwHelper.StringPointerValueOrNull(item.Image),
		ImageRegistryType: fwHelper.StringPointerValueOrNull(item.ImageRegistryType),
		Command:           command,
		Args:              args,
		Env:               env,
		Ports:             ports,
		Resources:         resources,
		Probe:             probe,
		DnsConfig:         dnsConfig,
		SecurityContext:   securityContext,
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceCustomConfigAttrTypes, model)
}

func expandAgsSandboxInstanceEnvVars(ctx context.Context, list types.List) ([]*agsv20250920.EnvVar, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsSandboxInstanceEnvVarModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	envVars := make([]*agsv20250920.EnvVar, 0, len(models))
	for _, model := range models {
		envVar := &agsv20250920.EnvVar{}
		if !model.Name.IsNull() && !model.Name.IsUnknown() {
			envVar.Name = model.Name.ValueStringPointer()
		}
		if !model.Value.IsNull() && !model.Value.IsUnknown() {
			envVar.Value = model.Value.ValueStringPointer()
		}
		envVars = append(envVars, envVar)
	}

	return envVars, diags
}

func flattenAgsSandboxInstanceEnvVars(ctx context.Context, items []*agsv20250920.EnvVar) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsSandboxInstanceEnvVarObjectType), nil
	}

	models := make([]agsSandboxInstanceEnvVarModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsSandboxInstanceEnvVarModel{
			Name:  fwHelper.StringPointerValueOrNull(item.Name),
			Value: fwHelper.StringPointerValueOrNull(item.Value),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsSandboxInstanceEnvVarObjectType), nil
	}

	return types.ListValueFrom(ctx, agsSandboxInstanceEnvVarObjectType, models)
}

func expandAgsSandboxInstancePorts(ctx context.Context, list types.List) ([]*agsv20250920.PortConfiguration, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsSandboxInstancePortModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	ports := make([]*agsv20250920.PortConfiguration, 0, len(models))
	for _, model := range models {
		port := &agsv20250920.PortConfiguration{}
		if !model.Name.IsNull() && !model.Name.IsUnknown() {
			port.Name = model.Name.ValueStringPointer()
		}
		if !model.Port.IsNull() && !model.Port.IsUnknown() {
			port.Port = model.Port.ValueInt64Pointer()
		}
		if !model.Protocol.IsNull() && !model.Protocol.IsUnknown() {
			port.Protocol = model.Protocol.ValueStringPointer()
		}
		ports = append(ports, port)
	}

	return ports, diags
}

func flattenAgsSandboxInstancePorts(ctx context.Context, items []*agsv20250920.PortConfiguration) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsSandboxInstancePortObjectType), nil
	}

	models := make([]agsSandboxInstancePortModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsSandboxInstancePortModel{
			Name:     fwHelper.StringPointerValueOrNull(item.Name),
			Port:     fwHelper.Int64ValueOrNull(item.Port),
			Protocol: fwHelper.StringPointerValueOrNull(item.Protocol),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsSandboxInstancePortObjectType), nil
	}

	return types.ListValueFrom(ctx, agsSandboxInstancePortObjectType, models)
}

func expandAgsSandboxInstanceResources(ctx context.Context, object types.Object) (*agsv20250920.ResourceConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsSandboxInstanceResourceModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	resources := &agsv20250920.ResourceConfiguration{}
	if !model.Cpu.IsNull() && !model.Cpu.IsUnknown() {
		resources.CPU = model.Cpu.ValueStringPointer()
	}
	if !model.Memory.IsNull() && !model.Memory.IsUnknown() {
		resources.Memory = model.Memory.ValueStringPointer()
	}
	if !model.Storage.IsNull() && !model.Storage.IsUnknown() {
		resources.Storage = model.Storage.ValueStringPointer()
	}

	return resources, diags
}

func flattenAgsSandboxInstanceResources(ctx context.Context, item *agsv20250920.ResourceConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceResourceAttrTypes), nil
	}

	model := agsSandboxInstanceResourceModel{
		Cpu:     fwHelper.StringPointerValueOrNull(item.CPU),
		Memory:  fwHelper.StringPointerValueOrNull(item.Memory),
		Storage: fwHelper.StringPointerValueOrNull(item.Storage),
	}
	// The cloud side reports an unset storage as an empty string instead of
	// omitting the field. A null configuration then ends up as "" in state,
	// which fails Terraform's post-apply consistency check. Normalize an
	// empty string back to null so "unset" stays null end to end.
	if !model.Storage.IsNull() && model.Storage.ValueString() == "" {
		model.Storage = types.StringNull()
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceResourceAttrTypes, model)
}

func expandAgsSandboxInstanceProbe(ctx context.Context, object types.Object) (*agsv20250920.ProbeConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsSandboxInstanceProbeModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	probe := &agsv20250920.ProbeConfiguration{}
	if !model.ReadyTimeoutMs.IsNull() && !model.ReadyTimeoutMs.IsUnknown() {
		probe.ReadyTimeoutMs = model.ReadyTimeoutMs.ValueInt64Pointer()
	}
	if !model.ProbeTimeoutMs.IsNull() && !model.ProbeTimeoutMs.IsUnknown() {
		probe.ProbeTimeoutMs = model.ProbeTimeoutMs.ValueInt64Pointer()
	}
	if !model.ProbePeriodMs.IsNull() && !model.ProbePeriodMs.IsUnknown() {
		probe.ProbePeriodMs = model.ProbePeriodMs.ValueInt64Pointer()
	}
	if !model.SuccessThreshold.IsNull() && !model.SuccessThreshold.IsUnknown() {
		probe.SuccessThreshold = model.SuccessThreshold.ValueInt64Pointer()
	}
	if !model.FailureThreshold.IsNull() && !model.FailureThreshold.IsUnknown() {
		probe.FailureThreshold = model.FailureThreshold.ValueInt64Pointer()
	}

	httpGet, httpGetDiags := expandAgsSandboxInstanceHttpGet(ctx, model.HttpGet)
	diags.Append(httpGetDiags...)
	if diags.HasError() {
		return nil, diags
	}
	probe.HttpGet = httpGet

	return probe, diags
}

func flattenAgsSandboxInstanceProbe(ctx context.Context, item *agsv20250920.ProbeConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceProbeAttrTypes), nil
	}

	httpGet, httpGetDiags := flattenAgsSandboxInstanceHttpGet(ctx, item.HttpGet)
	if httpGetDiags.HasError() {
		return types.ObjectNull(agsSandboxInstanceProbeAttrTypes), httpGetDiags
	}

	model := agsSandboxInstanceProbeModel{
		HttpGet:          httpGet,
		ReadyTimeoutMs:   fwHelper.Int64ValueOrNull(item.ReadyTimeoutMs),
		ProbeTimeoutMs:   fwHelper.Int64ValueOrNull(item.ProbeTimeoutMs),
		ProbePeriodMs:    fwHelper.Int64ValueOrNull(item.ProbePeriodMs),
		SuccessThreshold: fwHelper.Int64ValueOrNull(item.SuccessThreshold),
		FailureThreshold: fwHelper.Int64ValueOrNull(item.FailureThreshold),
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceProbeAttrTypes, model)
}

func expandAgsSandboxInstanceHttpGet(ctx context.Context, object types.Object) (*agsv20250920.HttpGetAction, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsSandboxInstanceHttpGetModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	httpGet := &agsv20250920.HttpGetAction{}
	if !model.Path.IsNull() && !model.Path.IsUnknown() {
		httpGet.Path = model.Path.ValueStringPointer()
	}
	if !model.Port.IsNull() && !model.Port.IsUnknown() {
		httpGet.Port = model.Port.ValueInt64Pointer()
	}
	if !model.Scheme.IsNull() && !model.Scheme.IsUnknown() {
		httpGet.Scheme = model.Scheme.ValueStringPointer()
	}

	return httpGet, diags
}

func flattenAgsSandboxInstanceHttpGet(ctx context.Context, item *agsv20250920.HttpGetAction) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceHttpGetAttrTypes), nil
	}

	model := agsSandboxInstanceHttpGetModel{
		Path:   fwHelper.StringPointerValueOrNull(item.Path),
		Port:   fwHelper.Int64ValueOrNull(item.Port),
		Scheme: fwHelper.StringPointerValueOrNull(item.Scheme),
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceHttpGetAttrTypes, model)
}

func expandAgsSandboxInstanceDnsConfig(ctx context.Context, object types.Object) (*agsv20250920.DNSConfig, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsSandboxInstanceDnsConfigModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	dnsConfig := &agsv20250920.DNSConfig{}

	servers, serverDiags := agsSandboxInstanceExpandStringSet(ctx, model.Servers)
	diags.Append(serverDiags...)
	if diags.HasError() {
		return nil, diags
	}
	dnsConfig.Servers = servers

	searches, searchDiags := agsSandboxInstanceExpandStringSet(ctx, model.Searches)
	diags.Append(searchDiags...)
	if diags.HasError() {
		return nil, diags
	}
	dnsConfig.Searches = searches

	options, optionDiags := agsSandboxInstanceExpandStringSet(ctx, model.Options)
	diags.Append(optionDiags...)
	if diags.HasError() {
		return nil, diags
	}
	dnsConfig.Options = options

	return dnsConfig, diags
}

func flattenAgsSandboxInstanceDnsConfig(ctx context.Context, item *agsv20250920.DNSConfig) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceDnsConfigAttrTypes), nil
	}

	var diags diag.Diagnostics

	servers, serverDiags := agsSandboxInstanceFlattenStringSet(ctx, item.Servers)
	diags.Append(serverDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceDnsConfigAttrTypes), diags
	}

	searches, searchDiags := agsSandboxInstanceFlattenStringSet(ctx, item.Searches)
	diags.Append(searchDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceDnsConfigAttrTypes), diags
	}

	options, optionDiags := agsSandboxInstanceFlattenStringSet(ctx, item.Options)
	diags.Append(optionDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceDnsConfigAttrTypes), diags
	}

	model := agsSandboxInstanceDnsConfigModel{
		Servers:  servers,
		Searches: searches,
		Options:  options,
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceDnsConfigAttrTypes, model)
}

func expandAgsSandboxInstanceSecurityContext(ctx context.Context, object types.Object) (*agsv20250920.SecurityContext, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsSandboxInstanceSecurityContextModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	securityContext := &agsv20250920.SecurityContext{}
	if !model.Privileged.IsNull() && !model.Privileged.IsUnknown() {
		securityContext.Privileged = model.Privileged.ValueBoolPointer()
	}

	capabilities, capabilityDiags := expandAgsSandboxInstanceCapabilities(ctx, model.Capabilities)
	diags.Append(capabilityDiags...)
	if diags.HasError() {
		return nil, diags
	}
	securityContext.Capabilities = capabilities

	return securityContext, diags
}

func flattenAgsSandboxInstanceSecurityContext(ctx context.Context, item *agsv20250920.SecurityContext) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceSecurityContextAttrTypes), nil
	}

	capabilities, capabilityDiags := flattenAgsSandboxInstanceCapabilities(ctx, item.Capabilities)
	if capabilityDiags.HasError() {
		return types.ObjectNull(agsSandboxInstanceSecurityContextAttrTypes), capabilityDiags
	}

	model := agsSandboxInstanceSecurityContextModel{
		Privileged:   fwHelper.BoolValueOrNull(item.Privileged),
		Capabilities: capabilities,
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceSecurityContextAttrTypes, model)
}

func expandAgsSandboxInstanceCapabilities(ctx context.Context, object types.Object) (*agsv20250920.Capabilities, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsSandboxInstanceCapabilitiesModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	capabilities := &agsv20250920.Capabilities{}

	add, addDiags := agsSandboxInstanceExpandStringSet(ctx, model.Add)
	diags.Append(addDiags...)
	if diags.HasError() {
		return nil, diags
	}
	capabilities.Add = add

	drop, dropDiags := agsSandboxInstanceExpandStringSet(ctx, model.Drop)
	diags.Append(dropDiags...)
	if diags.HasError() {
		return nil, diags
	}
	capabilities.Drop = drop

	return capabilities, diags
}

func flattenAgsSandboxInstanceCapabilities(ctx context.Context, item *agsv20250920.Capabilities) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceCapabilitiesAttrTypes), nil
	}

	var diags diag.Diagnostics

	add, addDiags := agsSandboxInstanceFlattenStringSet(ctx, item.Add)
	diags.Append(addDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCapabilitiesAttrTypes), diags
	}

	drop, dropDiags := agsSandboxInstanceFlattenStringSet(ctx, item.Drop)
	diags.Append(dropDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstanceCapabilitiesAttrTypes), diags
	}

	model := agsSandboxInstanceCapabilitiesModel{
		Add:  add,
		Drop: drop,
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceCapabilitiesAttrTypes, model)
}

func flattenAgsSandboxInstanceComputerConfiguration(ctx context.Context, item *agsv20250920.ComputerConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceComputerConfigAttrTypes), nil
	}

	waaConfiguration, waaDiags := flattenAgsSandboxInstanceWaaConfiguration(ctx, item.WAAConfiguration)
	if waaDiags.HasError() {
		return types.ObjectNull(agsSandboxInstanceComputerConfigAttrTypes), waaDiags
	}

	osWorldConfiguration, osWorldDiags := flattenAgsSandboxInstanceOsWorldConfiguration(ctx, item.OSWorldConfiguration)
	if osWorldDiags.HasError() {
		return types.ObjectNull(agsSandboxInstanceComputerConfigAttrTypes), osWorldDiags
	}

	model := agsSandboxInstanceComputerConfigurationModel{
		WaaConfiguration:     waaConfiguration,
		OsWorldConfiguration: osWorldConfiguration,
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceComputerConfigAttrTypes, model)
}

func flattenAgsSandboxInstanceWaaConfiguration(ctx context.Context, item *agsv20250920.WAAConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceWaaConfigurationAttrTypes), nil
	}

	model := agsSandboxInstanceWaaConfigurationModel{
		ImageId: fwHelper.StringPointerValueOrNull(item.ImageId),
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceWaaConfigurationAttrTypes, model)
}

func flattenAgsSandboxInstanceOsWorldConfiguration(ctx context.Context, item *agsv20250920.OSWorldConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstanceOsWorldConfigurationAttrTypes), nil
	}

	model := agsSandboxInstanceOsWorldConfigurationModel{
		Version: fwHelper.StringPointerValueOrNull(item.Version),
	}
	return types.ObjectValueFrom(ctx, agsSandboxInstanceOsWorldConfigurationAttrTypes, model)
}

func agsSandboxInstanceExpandStringSet(ctx context.Context, value types.Set) ([]*string, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var elements []string
	diags := value.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(elements) == 0 {
		return nil, diags
	}

	return agsSandboxInstanceHelperStrings(elements), diags
}

func agsSandboxInstanceFlattenStringSet(ctx context.Context, values []*string) (types.Set, diag.Diagnostics) {
	if len(values) == 0 {
		return types.SetNull(types.StringType), nil
	}
	return types.SetValueFrom(ctx, types.StringType, agsSandboxInstanceHelperPStrings(values))
}

// agsSandboxInstanceExpandStringList converts an ordered string list into the SDK pointer slice
// form. A null or unknown list is transmitted as "absent".
func agsSandboxInstanceExpandStringList(ctx context.Context, value types.List) ([]*string, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var elements []string
	diags := value.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(elements) == 0 {
		return nil, diags
	}

	return agsSandboxInstanceHelperStrings(elements), diags
}

// agsSandboxInstanceFlattenStringList converts the SDK pointer slice into an ordered string list,
// so an argument sequence keeps the order the API reports.
func agsSandboxInstanceFlattenStringList(ctx context.Context, values []*string) (types.List, diag.Diagnostics) {
	if len(values) == 0 {
		return types.ListNull(types.StringType), nil
	}
	return types.ListValueFrom(ctx, types.StringType, agsSandboxInstanceHelperPStrings(values))
}

// agsSandboxInstanceHelperStrings converts a string slice into the SDK pointer slice form.
func agsSandboxInstanceHelperStrings(values []string) []*string {
	result := make([]*string, 0, len(values))
	for i := range values {
		value := values[i]
		result = append(result, &value)
	}
	return result
}

// agsSandboxInstanceHelperPStrings converts an SDK pointer slice into a string slice, skipping
// nil entries so a malformed response cannot panic the provider.
func agsSandboxInstanceHelperPStrings(values []*string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		result = append(result, *value)
	}
	return result
}

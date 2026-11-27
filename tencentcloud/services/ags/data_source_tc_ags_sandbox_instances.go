package ags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &AgsSandboxInstancesDataSource{}

// NewAgsSandboxInstancesDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewAgsSandboxInstancesDataSource() datasource.DataSource {
	return &AgsSandboxInstancesDataSource{}
}

// AgsSandboxInstancesDataSource implements datasource.DataSource for
// tencentcloud_ags_sandbox_instances.
type AgsSandboxInstancesDataSource struct {
	fw.DataSourceWithConfigure
}

// AgsSandboxInstancesDataSourceModel maps the schema attributes.
type AgsSandboxInstancesDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	InstanceIds types.List   `tfsdk:"instance_ids"`
	InstanceSet types.List   `tfsdk:"instance_set"`
}

// AgsSandboxInstancesItemModel is the element type of instance_set.
type AgsSandboxInstancesItemModel struct {
	InstanceId            types.String `tfsdk:"instance_id"`
	ToolId                types.String `tfsdk:"tool_id"`
	ToolName              types.String `tfsdk:"tool_name"`
	Status                types.String `tfsdk:"status"`
	Persistent            types.Bool   `tfsdk:"persistent"`
	TimeoutSeconds        types.Int64  `tfsdk:"timeout_seconds"`
	ExpiresAt             types.String `tfsdk:"expires_at"`
	StopReason            types.String `tfsdk:"stop_reason"`
	CreateTime            types.String `tfsdk:"create_time"`
	UpdateTime            types.String `tfsdk:"update_time"`
	MountOptions          types.List   `tfsdk:"mount_options"`
	CustomConfiguration   types.Object `tfsdk:"custom_configuration"`
	ComputerConfiguration types.Object `tfsdk:"computer_configuration"`
	NetworkMode           types.String `tfsdk:"network_mode"`
	Metadata              types.List   `tfsdk:"metadata"`
	AuthMode              types.String `tfsdk:"auth_mode"`
}

type AgsSandboxInstancesMountOptionModel struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SubPath   types.String `tfsdk:"sub_path"`
	ReadOnly  types.Bool   `tfsdk:"read_only"`
}

type AgsSandboxInstancesMetadataVarModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type AgsSandboxInstancesEnvVarModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type AgsSandboxInstancesPortModel struct {
	Name     types.String `tfsdk:"name"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

type AgsSandboxInstancesResourceModel struct {
	Cpu     types.String `tfsdk:"cpu"`
	Memory  types.String `tfsdk:"memory"`
	Storage types.String `tfsdk:"storage"`
}

type AgsSandboxInstancesHttpGetModel struct {
	Path   types.String `tfsdk:"path"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`
}

type AgsSandboxInstancesProbeModel struct {
	HttpGet          types.Object `tfsdk:"http_get"`
	ReadyTimeoutMs   types.Int64  `tfsdk:"ready_timeout_ms"`
	ProbeTimeoutMs   types.Int64  `tfsdk:"probe_timeout_ms"`
	ProbePeriodMs    types.Int64  `tfsdk:"probe_period_ms"`
	SuccessThreshold types.Int64  `tfsdk:"success_threshold"`
	FailureThreshold types.Int64  `tfsdk:"failure_threshold"`
}

type AgsSandboxInstancesDnsConfigModel struct {
	Servers  types.List `tfsdk:"servers"`
	Searches types.List `tfsdk:"searches"`
	Options  types.List `tfsdk:"options"`
}

type AgsSandboxInstancesCapabilitiesModel struct {
	Add  types.List `tfsdk:"add"`
	Drop types.List `tfsdk:"drop"`
}

type AgsSandboxInstancesSecurityContextModel struct {
	Privileged   types.Bool   `tfsdk:"privileged"`
	Capabilities types.Object `tfsdk:"capabilities"`
}

type AgsSandboxInstancesCustomConfigurationModel struct {
	Image             types.String `tfsdk:"image"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	ImageDigest       types.String `tfsdk:"image_digest"`
	Command           types.List   `tfsdk:"command"`
	Args              types.List   `tfsdk:"args"`
	Env               types.List   `tfsdk:"env"`
	Ports             types.List   `tfsdk:"ports"`
	Resources         types.Object `tfsdk:"resources"`
	Probe             types.Object `tfsdk:"probe"`
	DnsConfig         types.Object `tfsdk:"dns_config"`
	SecurityContext   types.Object `tfsdk:"security_context"`
}

type AgsSandboxInstancesComputerConfigurationModel struct {
	WaaConfiguration     types.Object `tfsdk:"waa_configuration"`
	OsWorldConfiguration types.Object `tfsdk:"os_world_configuration"`
}

type AgsSandboxInstancesWaaConfigurationModel struct {
	ImageId types.String `tfsdk:"image_id"`
}

type AgsSandboxInstancesOsWorldConfigurationModel struct {
	Version types.String `tfsdk:"version"`
}

var agsSandboxInstancesMountOptionAttributes = map[string]datasource_schema.Attribute{
	"name": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Name of the storage configuration in the sandbox tool.",
	},
	"mount_path": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Local mount path of the sandbox instance (optional). By default it inherits the storage configuration of the tool.",
	},
	"sub_path": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Sub path of the sandbox instance storage mount (optional).",
	},
	"read_only": datasource_schema.BoolAttribute{
		Computed:    true,
		Description: "Read and write permission of the sandbox instance storage mount (optional). By default it inherits the storage configuration of the tool.",
	},
}

var agsSandboxInstancesMetadataVarAttributes = map[string]datasource_schema.Attribute{
	"name": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Metadata name.",
	},
	"value": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Metadata value.",
	},
}

var agsSandboxInstancesEnvVarAttributes = map[string]datasource_schema.Attribute{
	"name": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Environment variable name.",
	},
	"value": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Environment variable value.",
	},
}

var agsSandboxInstancesPortAttributes = map[string]datasource_schema.Attribute{
	"name": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Port name.",
	},
	"port": datasource_schema.Int64Attribute{
		Computed:    true,
		Description: "Port number.",
	},
	"protocol": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Protocol.",
	},
}

var agsSandboxInstancesResourceAttributes = map[string]datasource_schema.Attribute{
	"cpu": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Amount of cpu resources.",
	},
	"memory": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Amount of memory resources.",
	},
	"storage": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`. Note: this field may return null, which indicates no valid value.",
	},
}

var agsSandboxInstancesHttpGetAttributes = map[string]datasource_schema.Attribute{
	"path": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Request path.",
	},
	"port": datasource_schema.Int64Attribute{
		Computed:    true,
		Description: "Request port.",
	},
	"scheme": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Request scheme.",
	},
}

var agsSandboxInstancesProbeAttributes = map[string]datasource_schema.Attribute{
	"http_get": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesHttpGetAttributes,
		Description: "HTTP GET probe configuration.",
	},
	"ready_timeout_ms": datasource_schema.Int64Attribute{
		Computed:    true,
		Description: "Timeout in milliseconds for the health check to become ready.",
	},
	"probe_timeout_ms": datasource_schema.Int64Attribute{
		Computed:    true,
		Description: "Timeout in milliseconds of a single health check probe.",
	},
	"probe_period_ms": datasource_schema.Int64Attribute{
		Computed:    true,
		Description: "Interval in milliseconds between two health check probes.",
	},
	"success_threshold": datasource_schema.Int64Attribute{
		Computed:    true,
		Description: "Success threshold of the health check.",
	},
	"failure_threshold": datasource_schema.Int64Attribute{
		Computed:    true,
		Description: "Failure threshold of the health check.",
	},
}

var agsSandboxInstancesDnsConfigAttributes = map[string]datasource_schema.Attribute{
	"servers": datasource_schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Description: "DNS server addresses. Valid IP addresses are required. Default value: `10.0.0.1`.",
	},
	"searches": datasource_schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Description: "Search domains, corresponding to the `search` directive of resolv.conf.",
	},
	"options": datasource_schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Description: "Configuration options, corresponding to the options of resolv.conf.",
	},
}

var agsSandboxInstancesCapabilitiesAttributes = map[string]datasource_schema.Attribute{
	"add": datasource_schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Description: "Capabilities to add. Reference value: `[\"SYS_ADMIN\"]`.",
	},
	"drop": datasource_schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Description: "Capabilities to drop. Reference value: `[\"ALL\"]`.",
	},
}

var agsSandboxInstancesSecurityContextAttributes = map[string]datasource_schema.Attribute{
	"privileged": datasource_schema.BoolAttribute{
		Computed:    true,
		Description: "Whether privileged mode is enabled. Default value: `true`.",
	},
	"capabilities": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesCapabilitiesAttributes,
		Description: "POSIX capabilities that can be added to or removed from the running container.",
	},
}

var agsSandboxInstancesCustomConfigurationAttributes = map[string]datasource_schema.Attribute{
	"image": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Image address.",
	},
	"image_registry_type": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Image registry type. Valid values: `enterprise`, `personal`, `custom`. `enterprise`: TCR enterprise container image service; `personal`: CCR personal container image service.",
	},
	"image_digest": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Image digest.",
	},
	"command": datasource_schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Description: "Startup command.",
	},
	"args": datasource_schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Description: "Startup arguments.",
	},
	"env": datasource_schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: datasource_schema.NestedAttributeObject{Attributes: agsSandboxInstancesEnvVarAttributes},
		Description:  "Environment variables.",
	},
	"ports": datasource_schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: datasource_schema.NestedAttributeObject{Attributes: agsSandboxInstancesPortAttributes},
		Description:  "Port configuration.",
	},
	"resources": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesResourceAttributes,
		Description: "Resource configuration.",
	},
	"probe": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesProbeAttributes,
		Description: "Probe configuration.",
	},
	"dns_config": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesDnsConfigAttributes,
		Description: "Sandbox DNS configuration.",
	},
	"security_context": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesSecurityContextAttributes,
		Description: "Sandbox security context.",
	},
}

var agsSandboxInstancesWaaConfigurationAttributes = map[string]datasource_schema.Attribute{
	"image_id": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Id of the custom WAA image.",
	},
}

var agsSandboxInstancesOsWorldConfigurationAttributes = map[string]datasource_schema.Attribute{
	"version": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Built-in OSWorld version. Valid values: `osworld1` (osworld v1), `osworld2` (osworld v2). Default value: `osworld1`.",
	},
}

var agsSandboxInstancesComputerConfigurationAttributes = map[string]datasource_schema.Attribute{
	"waa_configuration": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesWaaConfigurationAttributes,
		Description: "WAA sandbox tool configuration.",
	},
	"os_world_configuration": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesOsWorldConfigurationAttributes,
		Description: "Built-in OSWorld configuration.",
	},
}

var agsSandboxInstancesItemAttributes = map[string]datasource_schema.Attribute{
	"instance_id": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Unique identifier of the sandbox instance.",
	},
	"tool_id": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Id of the sandbox tool the instance belongs to.",
	},
	"tool_name": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Name of the sandbox tool the instance belongs to.",
	},
	"status": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Instance status. Valid values: `STARTING` (starting), `RUNNING` (running), `STOPPING` (stopping), `STOPPED` (stopped), `STOP_FAILED` (stop failed), `FAILED` (failed), `RESIZING` (updating instance specifications).",
	},
	"persistent": datasource_schema.BoolAttribute{
		Computed:    true,
		Description: "Whether the instance is persistent.",
	},
	"timeout_seconds": datasource_schema.Int64Attribute{
		Computed:    true,
		Description: "Timeout in seconds. Null means no timeout is set.",
	},
	"expires_at": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Expiration time in ISO 8601 format. Null means no expiration time.",
	},
	"stop_reason": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Stop reason: `manual` (manual), `timeout` (timeout), `error` (error), `system` (system). It has a value only when the status is `STOPPED`, `STOP_FAILED` or `FAILED`. When the provider fails to stop the instance, the status is `STOP_FAILED` and the reason is `error`.",
	},
	"create_time": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Creation time in ISO 8601 format.",
	},
	"update_time": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Last updated time in ISO 8601 format.",
	},
	"mount_options": datasource_schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: datasource_schema.NestedAttributeObject{Attributes: agsSandboxInstancesMountOptionAttributes},
		Description:  "Storage mount options.",
	},
	"custom_configuration": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesCustomConfigurationAttributes,
		Description: "Custom configuration of the sandbox instance.",
	},
	"computer_configuration": datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  agsSandboxInstancesComputerConfigurationAttributes,
		Description: "Sandbox configuration of desktop computer environments.",
	},
	"network_mode": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Network mode. Valid values: `PUBLIC` (public network access), `SANDBOX` (no network), `INTERNAL_SERVICE` (Tencent Cloud internal public services). It can override the network configuration at the tool level, but it is ineffective if the tool itself does not support VPC networks.",
	},
	"metadata": datasource_schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: datasource_schema.NestedAttributeObject{Attributes: agsSandboxInstancesMetadataVarAttributes},
		Description:  "Sandbox instance metadata.",
	},
	"auth_mode": datasource_schema.StringAttribute{
		Computed:    true,
		Description: "Sandbox access authentication mode. Valid values: `DEFAULT` (default, i.e. TOKEN authentication), `TOKEN` (Token authentication, all port access requires a TOKEN), `NONE` (no authentication, all port access does not require a TOKEN), `PUBLIC` (public mode, the ENVD management port (49983) access requires a TOKEN while other ports do not). Default value: `DEFAULT`.",
	},
}

var (
	agsSandboxInstancesMountOptionAttrTypes      = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesMountOptionAttributes)
	agsSandboxInstancesMetadataVarAttrTypes      = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesMetadataVarAttributes)
	agsSandboxInstancesEnvVarAttrTypes           = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesEnvVarAttributes)
	agsSandboxInstancesPortAttrTypes             = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesPortAttributes)
	agsSandboxInstancesResourceAttrTypes         = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesResourceAttributes)
	agsSandboxInstancesHttpGetAttrTypes          = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesHttpGetAttributes)
	agsSandboxInstancesProbeAttrTypes            = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesProbeAttributes)
	agsSandboxInstancesDnsConfigAttrTypes        = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesDnsConfigAttributes)
	agsSandboxInstancesCapabilitiesAttrTypes     = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesCapabilitiesAttributes)
	agsSandboxInstancesSecurityContextAttrTypes  = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesSecurityContextAttributes)
	agsSandboxInstancesCustomConfigAttrTypes     = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesCustomConfigurationAttributes)
	agsSandboxInstancesWaaConfigurationAttrTypes = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesWaaConfigurationAttributes)
	agsSandboxInstancesOsWorldConfigAttrTypes    = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesOsWorldConfigurationAttributes)
	agsSandboxInstancesComputerConfigAttrTypes   = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesComputerConfigurationAttributes)
	agsSandboxInstancesItemAttrTypes             = agsSandboxInstancesAttributeTypesOf(agsSandboxInstancesItemAttributes)
)

// agsSandboxInstancesAttributeTypesOf mirrors the value types of a nested
// data source attribute map so the conversion helpers stay in sync with the
// schema definition.
func agsSandboxInstancesAttributeTypesOf(attributes map[string]datasource_schema.Attribute) map[string]attr.Type {
	result := make(map[string]attr.Type, len(attributes))
	for name, attribute := range attributes {
		result[name] = attribute.GetType()
	}
	return result
}

func (d *AgsSandboxInstancesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_instances"
}

func (d *AgsSandboxInstancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_schema.Schema{
		Description: "Use this data source to query the list of sandbox instances.",
		Attributes: map[string]datasource_schema.Attribute{
			// The id carries no cloud meaning — it is a random token
			// generated at read time (see the end of Read).
			"id": datasource_schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			"instance_ids": datasource_schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of sandbox instance ids to query. If empty, all sandbox instances are queried. At most 100 ids are supported.",
			},
			"instance_set": datasource_schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: datasource_schema.NestedAttributeObject{Attributes: agsSandboxInstancesItemAttributes},
				Description:  "Sandbox instance list.",
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *AgsSandboxInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AgsSandboxInstancesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeSandboxInstancesRequest := agsv20250920.NewDescribeSandboxInstancesRequest()
	if !config.InstanceIds.IsNull() && !config.InstanceIds.IsUnknown() {
		describeSandboxInstancesRequest.InstanceIds = fwHelper.PointersFromList[string](ctx, config.InstanceIds, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	service := NewAgsService(d.Client())

	sandboxInstances, err := service.DescribeSandboxInstances(ctx, describeSandboxInstancesRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying AGS sandbox instances",
			fmt.Sprintf("DescribeSandboxInstances failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	items := make([]AgsSandboxInstancesItemModel, 0, len(sandboxInstances))
	var diags diag.Diagnostics
	for _, item := range sandboxInstances {
		if item == nil {
			continue
		}

		itemModel, itemDiags := flattenSandboxInstancesItem(ctx, item)
		diags.Append(itemDiags...)
		if diags.HasError() {
			break
		}
		items = append(items, itemModel)
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: agsSandboxInstancesItemAttrTypes}, items)
	resp.Diagnostics.Append(listDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.InstanceSet = list

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func flattenSandboxInstancesItem(ctx context.Context, item *agsv20250920.SandboxInstance) (AgsSandboxInstancesItemModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := AgsSandboxInstancesItemModel{
		InstanceId:     types.StringPointerValue(item.InstanceId),
		ToolId:         types.StringPointerValue(item.ToolId),
		ToolName:       types.StringPointerValue(item.ToolName),
		Status:         types.StringPointerValue(item.Status),
		Persistent:     types.BoolPointerValue(item.Persistent),
		TimeoutSeconds: fwHelper.Int64ValueFromUint(item.TimeoutSeconds),
		ExpiresAt:      types.StringPointerValue(item.ExpiresAt),
		StopReason:     types.StringPointerValue(item.StopReason),
		CreateTime:     types.StringPointerValue(item.CreateTime),
		UpdateTime:     types.StringPointerValue(item.UpdateTime),
		NetworkMode:    types.StringPointerValue(item.NetworkMode),
		AuthMode:       types.StringPointerValue(item.AuthMode),
	}

	mountOptions, mountDiags := flattenSandboxInstancesMountOptions(ctx, item.MountOptions)
	diags.Append(mountDiags...)
	model.MountOptions = mountOptions

	customConfiguration, customDiags := flattenSandboxInstancesCustomConfiguration(ctx, item.CustomConfiguration)
	diags.Append(customDiags...)
	model.CustomConfiguration = customConfiguration

	computerConfiguration, computerDiags := flattenSandboxInstancesComputerConfiguration(ctx, item.ComputerConfiguration)
	diags.Append(computerDiags...)
	model.ComputerConfiguration = computerConfiguration

	metadata, metadataDiags := flattenSandboxInstancesMetadataVars(ctx, item.Metadata)
	diags.Append(metadataDiags...)
	model.Metadata = metadata

	return model, diags
}

func flattenSandboxInstancesMountOptions(ctx context.Context, items []*agsv20250920.MountOption) (types.List, diag.Diagnostics) {
	models := make([]AgsSandboxInstancesMountOptionModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, AgsSandboxInstancesMountOptionModel{
			Name:      types.StringPointerValue(item.Name),
			MountPath: types.StringPointerValue(item.MountPath),
			SubPath:   types.StringPointerValue(item.SubPath),
			ReadOnly:  types.BoolPointerValue(item.ReadOnly),
		})
	}

	var diags diag.Diagnostics
	list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: agsSandboxInstancesMountOptionAttrTypes}, models)
	diags.Append(listDiags...)
	return list, diags
}

func flattenSandboxInstancesMetadataVars(ctx context.Context, items []*agsv20250920.MetadataVar) (types.List, diag.Diagnostics) {
	models := make([]AgsSandboxInstancesMetadataVarModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, AgsSandboxInstancesMetadataVarModel{
			Name:  types.StringPointerValue(item.Name),
			Value: types.StringPointerValue(item.Value),
		})
	}

	var diags diag.Diagnostics
	list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: agsSandboxInstancesMetadataVarAttrTypes}, models)
	diags.Append(listDiags...)
	return list, diags
}

func flattenSandboxInstancesEnvVars(ctx context.Context, items []*agsv20250920.EnvVar) (types.List, diag.Diagnostics) {
	models := make([]AgsSandboxInstancesEnvVarModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, AgsSandboxInstancesEnvVarModel{
			Name:  types.StringPointerValue(item.Name),
			Value: types.StringPointerValue(item.Value),
		})
	}

	var diags diag.Diagnostics
	list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: agsSandboxInstancesEnvVarAttrTypes}, models)
	diags.Append(listDiags...)
	return list, diags
}

func flattenSandboxInstancesPorts(ctx context.Context, items []*agsv20250920.PortConfiguration) (types.List, diag.Diagnostics) {
	models := make([]AgsSandboxInstancesPortModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, AgsSandboxInstancesPortModel{
			Name:     types.StringPointerValue(item.Name),
			Port:     types.Int64PointerValue(item.Port),
			Protocol: types.StringPointerValue(item.Protocol),
		})
	}

	var diags diag.Diagnostics
	list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: agsSandboxInstancesPortAttrTypes}, models)
	diags.Append(listDiags...)
	return list, diags
}

func flattenSandboxInstancesResources(ctx context.Context, item *agsv20250920.ResourceConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesResourceAttrTypes), nil
	}

	model := AgsSandboxInstancesResourceModel{
		Cpu:     types.StringPointerValue(item.CPU),
		Memory:  types.StringPointerValue(item.Memory),
		Storage: types.StringPointerValue(item.Storage),
	}

	var diags diag.Diagnostics
	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesResourceAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

func flattenSandboxInstancesHttpGet(ctx context.Context, item *agsv20250920.HttpGetAction) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesHttpGetAttrTypes), nil
	}

	model := AgsSandboxInstancesHttpGetModel{
		Path:   types.StringPointerValue(item.Path),
		Port:   types.Int64PointerValue(item.Port),
		Scheme: types.StringPointerValue(item.Scheme),
	}

	var diags diag.Diagnostics
	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesHttpGetAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

func flattenSandboxInstancesProbe(ctx context.Context, item *agsv20250920.ProbeConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesProbeAttrTypes), nil
	}

	var diags diag.Diagnostics

	httpGet, httpGetDiags := flattenSandboxInstancesHttpGet(ctx, item.HttpGet)
	diags.Append(httpGetDiags...)

	model := AgsSandboxInstancesProbeModel{
		HttpGet:          httpGet,
		ReadyTimeoutMs:   types.Int64PointerValue(item.ReadyTimeoutMs),
		ProbeTimeoutMs:   types.Int64PointerValue(item.ProbeTimeoutMs),
		ProbePeriodMs:    types.Int64PointerValue(item.ProbePeriodMs),
		SuccessThreshold: types.Int64PointerValue(item.SuccessThreshold),
		FailureThreshold: types.Int64PointerValue(item.FailureThreshold),
	}

	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesProbeAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

func flattenSandboxInstancesDnsConfig(ctx context.Context, item *agsv20250920.DNSConfig) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesDnsConfigAttrTypes), nil
	}

	var diags diag.Diagnostics

	model := AgsSandboxInstancesDnsConfigModel{
		Servers:  fwHelper.ListFromPointers[string](ctx, types.StringType, item.Servers, &diags),
		Searches: fwHelper.ListFromPointers[string](ctx, types.StringType, item.Searches, &diags),
		Options:  fwHelper.ListFromPointers[string](ctx, types.StringType, item.Options, &diags),
	}

	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesDnsConfigAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

func flattenSandboxInstancesCapabilities(ctx context.Context, item *agsv20250920.Capabilities) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesCapabilitiesAttrTypes), nil
	}

	var diags diag.Diagnostics

	model := AgsSandboxInstancesCapabilitiesModel{
		Add:  fwHelper.ListFromPointers[string](ctx, types.StringType, item.Add, &diags),
		Drop: fwHelper.ListFromPointers[string](ctx, types.StringType, item.Drop, &diags),
	}

	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesCapabilitiesAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

func flattenSandboxInstancesSecurityContext(ctx context.Context, item *agsv20250920.SecurityContext) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesSecurityContextAttrTypes), nil
	}

	var diags diag.Diagnostics

	capabilities, capabilitiesDiags := flattenSandboxInstancesCapabilities(ctx, item.Capabilities)
	diags.Append(capabilitiesDiags...)

	model := AgsSandboxInstancesSecurityContextModel{
		Privileged:   types.BoolPointerValue(item.Privileged),
		Capabilities: capabilities,
	}

	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesSecurityContextAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

func flattenSandboxInstancesCustomConfiguration(ctx context.Context, item *agsv20250920.CustomConfigurationDetail) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesCustomConfigAttrTypes), nil
	}

	var diags diag.Diagnostics

	env, envDiags := flattenSandboxInstancesEnvVars(ctx, item.Env)
	diags.Append(envDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstancesCustomConfigAttrTypes), diags
	}

	ports, portsDiags := flattenSandboxInstancesPorts(ctx, item.Ports)
	diags.Append(portsDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstancesCustomConfigAttrTypes), diags
	}

	resources, resourcesDiags := flattenSandboxInstancesResources(ctx, item.Resources)
	diags.Append(resourcesDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstancesCustomConfigAttrTypes), diags
	}

	probe, probeDiags := flattenSandboxInstancesProbe(ctx, item.Probe)
	diags.Append(probeDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstancesCustomConfigAttrTypes), diags
	}

	dnsConfig, dnsConfigDiags := flattenSandboxInstancesDnsConfig(ctx, item.DNSConfig)
	diags.Append(dnsConfigDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstancesCustomConfigAttrTypes), diags
	}

	securityContext, securityContextDiags := flattenSandboxInstancesSecurityContext(ctx, item.SecurityContext)
	diags.Append(securityContextDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstancesCustomConfigAttrTypes), diags
	}

	model := AgsSandboxInstancesCustomConfigurationModel{
		Image:             types.StringPointerValue(item.Image),
		ImageRegistryType: types.StringPointerValue(item.ImageRegistryType),
		ImageDigest:       types.StringPointerValue(item.ImageDigest),
		Command:           fwHelper.ListFromPointers[string](ctx, types.StringType, item.Command, &diags),
		Args:              fwHelper.ListFromPointers[string](ctx, types.StringType, item.Args, &diags),
		Env:               env,
		Ports:             ports,
		Resources:         resources,
		Probe:             probe,
		DnsConfig:         dnsConfig,
		SecurityContext:   securityContext,
	}

	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesCustomConfigAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

func flattenSandboxInstancesWaaConfiguration(ctx context.Context, item *agsv20250920.WAAConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesWaaConfigurationAttrTypes), nil
	}

	model := AgsSandboxInstancesWaaConfigurationModel{
		ImageId: types.StringPointerValue(item.ImageId),
	}

	var diags diag.Diagnostics
	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesWaaConfigurationAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

func flattenSandboxInstancesOsWorldConfiguration(ctx context.Context, item *agsv20250920.OSWorldConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesOsWorldConfigAttrTypes), nil
	}

	model := AgsSandboxInstancesOsWorldConfigurationModel{
		Version: types.StringPointerValue(item.Version),
	}

	var diags diag.Diagnostics
	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesOsWorldConfigAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

func flattenSandboxInstancesComputerConfiguration(ctx context.Context, item *agsv20250920.ComputerConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSandboxInstancesComputerConfigAttrTypes), nil
	}

	var diags diag.Diagnostics

	waaConfiguration, waaDiags := flattenSandboxInstancesWaaConfiguration(ctx, item.WAAConfiguration)
	diags.Append(waaDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstancesComputerConfigAttrTypes), diags
	}

	osWorldConfiguration, osWorldDiags := flattenSandboxInstancesOsWorldConfiguration(ctx, item.OSWorldConfiguration)
	diags.Append(osWorldDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstancesComputerConfigAttrTypes), diags
	}

	model := AgsSandboxInstancesComputerConfigurationModel{
		WaaConfiguration:     waaConfiguration,
		OsWorldConfiguration: osWorldConfiguration,
	}

	object, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesComputerConfigAttrTypes, model)
	diags.Append(objectDiags...)
	return object, diags
}

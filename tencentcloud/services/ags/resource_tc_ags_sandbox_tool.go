package ags

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// agsSandboxToolPollInterval is the interval used while waiting for a sandbox
// tool to reach a stable state.
const agsSandboxToolPollInterval = 3 * time.Second

var (
	_ resource.ResourceWithConfigure   = &AgsSandboxToolResource{}
	_ resource.ResourceWithImportState = &AgsSandboxToolResource{}
)

// NewAgsSandboxToolResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAgsSandboxToolResource() resource.Resource {
	return &AgsSandboxToolResource{}
}

// AgsSandboxToolResource implements resource.Resource for
// tencentcloud_ags_sandbox_tool.
type AgsSandboxToolResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AgsSandboxToolResourceModel maps the schema attributes.
//
// DefaultTimeout and ClientToken are create-only inputs: DescribeSandboxToolList
// does not return them (it only returns DefaultTimeoutSeconds, which cannot be
// converted back into the "5m"/"300s" form without guessing the original unit),
// so Read keeps whatever value is already stored in state.
type AgsSandboxToolResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	ToolId                types.String `tfsdk:"tool_id"`
	ToolName              types.String `tfsdk:"tool_name"`
	ToolType              types.String `tfsdk:"tool_type"`
	NetworkConfiguration  types.Object `tfsdk:"network_configuration"`
	Description           types.String `tfsdk:"description"`
	DefaultTimeout        types.String `tfsdk:"default_timeout"`
	Tags                  types.List   `tfsdk:"tags"`
	ClientToken           types.String `tfsdk:"client_token"`
	RoleArn               types.String `tfsdk:"role_arn"`
	StorageMounts         types.List   `tfsdk:"storage_mounts"`
	CustomConfiguration   types.Object `tfsdk:"custom_configuration"`
	ComputerConfiguration types.Object `tfsdk:"computer_configuration"`
	LogConfiguration      types.Object `tfsdk:"log_configuration"`
	Persistent            types.Bool   `tfsdk:"persistent"`
	Status                types.String `tfsdk:"status"`
	StatusReason          types.String `tfsdk:"status_reason"`
	CreateTime            types.String `tfsdk:"create_time"`
	UpdateTime            types.String `tfsdk:"update_time"`
}

type agsNetworkConfigurationModel struct {
	NetworkMode types.String `tfsdk:"network_mode"`
	VpcConfig   types.Object `tfsdk:"vpc_config"`
}

type agsVpcConfigModel struct {
	SubnetIds        types.Set `tfsdk:"subnet_ids"`
	SecurityGroupIds types.Set `tfsdk:"security_group_ids"`
}

type agsTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type agsStorageMountModel struct {
	Name           types.String `tfsdk:"name"`
	StorageSource  types.Object `tfsdk:"storage_source"`
	MountPath      types.String `tfsdk:"mount_path"`
	ReadOnly       types.Bool   `tfsdk:"read_only"`
	Volume         types.Object `tfsdk:"volume"`
	VolumeTemplate types.Object `tfsdk:"volume_template"`
}

type agsStorageSourceModel struct {
	Cos         types.Object `tfsdk:"cos"`
	Image       types.Object `tfsdk:"image"`
	Cfs         types.Object `tfsdk:"cfs"`
	AgentBucket types.Object `tfsdk:"agent_bucket"`
}

type agsCosStorageSourceModel struct {
	Endpoint   types.String `tfsdk:"endpoint"`
	BucketName types.String `tfsdk:"bucket_name"`
	BucketPath types.String `tfsdk:"bucket_path"`
}

type agsImageStorageSourceModel struct {
	Reference         types.String `tfsdk:"reference"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	SubPath           types.String `tfsdk:"sub_path"`
	Digest            types.String `tfsdk:"digest"`
}

type agsCfsStorageSourceModel struct {
	FileSystemId types.String `tfsdk:"file_system_id"`
	Path         types.String `tfsdk:"path"`
}

type agsAgentBucketStorageSourceModel struct {
	LibraryId    types.String `tfsdk:"library_id"`
	SpaceId      types.String `tfsdk:"space_id"`
	AccessDomain types.String `tfsdk:"access_domain"`
}

type agsVolumeModel struct {
	VolumeName types.String `tfsdk:"volume_name"`
	VolumeId   types.String `tfsdk:"volume_id"`
}

type agsVolumeTemplateModel struct {
	VolumeTemplateName types.String `tfsdk:"volume_template_name"`
	ReuseKey           types.String `tfsdk:"reuse_key"`
}

type agsCustomConfigurationModel struct {
	Image             types.String `tfsdk:"image"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	Command           types.Set    `tfsdk:"command"`
	Args              types.Set    `tfsdk:"args"`
	Env               types.List   `tfsdk:"env"`
	Ports             types.List   `tfsdk:"ports"`
	Resources         types.Object `tfsdk:"resources"`
	Probe             types.Object `tfsdk:"probe"`
	DnsConfig         types.Object `tfsdk:"dns_config"`
	SecurityContext   types.Object `tfsdk:"security_context"`
}

type agsEnvVarModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type agsPortModel struct {
	Name     types.String `tfsdk:"name"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

type agsResourceModel struct {
	Cpu     types.String `tfsdk:"cpu"`
	Memory  types.String `tfsdk:"memory"`
	Storage types.String `tfsdk:"storage"`
}

type agsProbeModel struct {
	HttpGet          types.Object `tfsdk:"http_get"`
	ReadyTimeoutMs   types.Int64  `tfsdk:"ready_timeout_ms"`
	ProbeTimeoutMs   types.Int64  `tfsdk:"probe_timeout_ms"`
	ProbePeriodMs    types.Int64  `tfsdk:"probe_period_ms"`
	SuccessThreshold types.Int64  `tfsdk:"success_threshold"`
	FailureThreshold types.Int64  `tfsdk:"failure_threshold"`
}

type agsHttpGetModel struct {
	Path   types.String `tfsdk:"path"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`
}

type agsDnsConfigModel struct {
	Servers  types.Set `tfsdk:"servers"`
	Searches types.Set `tfsdk:"searches"`
	Options  types.Set `tfsdk:"options"`
}

type agsSecurityContextModel struct {
	Privileged   types.Bool   `tfsdk:"privileged"`
	Capabilities types.Object `tfsdk:"capabilities"`
}

type agsCapabilitiesModel struct {
	Add  types.Set `tfsdk:"add"`
	Drop types.Set `tfsdk:"drop"`
}

type agsComputerConfigurationModel struct {
	WaaConfiguration     types.Object `tfsdk:"waa_configuration"`
	OsWorldConfiguration types.Object `tfsdk:"os_world_configuration"`
}

type agsWaaConfigurationModel struct {
	ImageId types.String `tfsdk:"image_id"`
}

type agsOsWorldConfigurationModel struct {
	Version types.String `tfsdk:"version"`
}

type agsLogConfigurationModel struct {
	ClsConfig  types.Object `tfsdk:"cls_config"`
	LogSources types.Object `tfsdk:"log_sources"`
}

type agsClsConfigModel struct {
	TopicId types.String `tfsdk:"topic_id"`
}

type agsLogSourcesModel struct {
	Files types.Set `tfsdk:"files"`
}

var agsVpcConfigAttributes = map[string]schema.Attribute{
	"subnet_ids": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "List of VPC subnet ids.",
	},
	"security_group_ids": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "List of security group ids.",
	},
}

var agsNetworkConfigurationAttributes = map[string]schema.Attribute{
	"network_mode": schema.StringAttribute{
		Required:    true,
		Description: "Network mode. Valid values: `PUBLIC`, `VPC`, `SANDBOX`.",
	},
	"vpc_config": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsVpcConfigAttributes,
		Description: "VPC related configuration. It is used when `network_mode` is `VPC`.",
	},
}

var agsTagAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Optional:    true,
		Description: "Tag key.",
	},
	"value": schema.StringAttribute{
		Optional:    true,
		Description: "Tag value.",
	},
}

var agsCosStorageSourceAttributes = map[string]schema.Attribute{
	"endpoint": schema.StringAttribute{
		Optional:    true,
		Description: "Object storage access endpoint.",
	},
	"bucket_name": schema.StringAttribute{
		Optional:    true,
		Description: "Object storage bucket name.",
	},
	"bucket_path": schema.StringAttribute{
		Optional:    true,
		Description: "Object storage bucket path. It must be an absolute path starting with `/`.",
	},
}

var agsImageStorageSourceAttributes = map[string]schema.Attribute{
	"reference": schema.StringAttribute{
		Optional:    true,
		Description: "Image address.",
	},
	"image_registry_type": schema.StringAttribute{
		Optional:    true,
		Description: "Image registry type. Valid values: `enterprise`, `personal`.",
	},
	"sub_path": schema.StringAttribute{
		Optional:    true,
		Description: "Path inside the image.",
	},
	"digest": schema.StringAttribute{
		Computed:    true,
		Description: "Image digest resolved by the cloud side. It is not required in requests and is only reported by the cloud.",
	},
}

var agsCfsStorageSourceAttributes = map[string]schema.Attribute{
	"file_system_id": schema.StringAttribute{
		Optional:    true,
		Description: "Id of the CFS file system.",
	},
	"path": schema.StringAttribute{
		Optional:    true,
		Description: "CFS mount path.",
	},
}

var agsAgentBucketStorageSourceAttributes = map[string]schema.Attribute{
	"library_id": schema.StringAttribute{
		Optional:    true,
		Description: "Library id of the AgentBucket.",
	},
	"space_id": schema.StringAttribute{
		Optional:    true,
		Description: "Space id of the AgentBucket.",
	},
	"access_domain": schema.StringAttribute{
		Optional:    true,
		Description: "Access domain of the AgentBucket.",
	},
}

var agsStorageSourceAttributes = map[string]schema.Attribute{
	"cos": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsCosStorageSourceAttributes,
		Description: "Object storage bucket configuration.",
	},
	"image": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsImageStorageSourceAttributes,
		Description: "Image volume configuration.",
	},
	"cfs": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsCfsStorageSourceAttributes,
		Description: "File storage configuration.",
	},
	"agent_bucket": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsAgentBucketStorageSourceAttributes,
		Description: "AgentBucket storage configuration.",
	},
}

var agsVolumeAttributes = map[string]schema.Attribute{
	"volume_name": schema.StringAttribute{
		Optional:    true,
		Description: "Volume name. The cloud side resolves it into a stable volume id under the current AppId.",
	},
	"volume_id": schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Stable volume id in the form `vol-xxxxxxxx`, resolved by the cloud side.",
	},
}

var agsVolumeTemplateAttributes = map[string]schema.Attribute{
	"volume_template_name": schema.StringAttribute{
		Required:    true,
		Description: "Name of the volume template.",
	},
	"reuse_key": schema.StringAttribute{
		Optional:    true,
		Description: "Reuse key. It can be a literal value or a `${metadata.xxx}` expression.",
	},
}

var agsStorageMountAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Name of the storage mount configuration.",
	},
	"storage_source": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsStorageSourceAttributes,
		Description: "Storage configuration.",
	},
	"mount_path": schema.StringAttribute{
		Optional:    true,
		Description: "Local mount path inside the sandbox instance.",
	},
	"read_only": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Whether the mount is read only. Default value: `false`.",
	},
	"volume": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsVolumeAttributes,
		Description: "Volume reference.",
	},
	"volume_template": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsVolumeTemplateAttributes,
		Description: "Volume template reference and derivation rule.",
	},
}

var agsEnvVarAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Environment variable name.",
	},
	"value": schema.StringAttribute{
		Optional:    true,
		Description: "Environment variable value.",
	},
}

var agsPortAttributes = map[string]schema.Attribute{
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

var agsResourceAttributes = map[string]schema.Attribute{
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

var agsHttpGetAttributes = map[string]schema.Attribute{
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

var agsProbeAttributes = map[string]schema.Attribute{
	"http_get": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsHttpGetAttributes,
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

var agsDnsConfigAttributes = map[string]schema.Attribute{
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

var agsCapabilitiesAttributes = map[string]schema.Attribute{
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

var agsSecurityContextAttributes = map[string]schema.Attribute{
	"privileged": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Whether the sandbox runs in privileged mode. Default value: `true`.",
	},
	"capabilities": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsCapabilitiesAttributes,
		Description: "POSIX capabilities that can be added to or removed from the running container.",
	},
}

var agsCustomConfigurationAttributes = map[string]schema.Attribute{
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
	"args": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Startup arguments.",
	},
	"env": schema.ListNestedAttribute{
		Optional:    true,
		NestedObject: schema.NestedAttributeObject{Attributes: agsEnvVarAttributes},
		Description: "Environment variables.",
	},
	"ports": schema.ListNestedAttribute{
		Optional:    true,
		NestedObject: schema.NestedAttributeObject{Attributes: agsPortAttributes},
		Description: "Port configuration.",
	},
	"resources": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsResourceAttributes,
		Description: "Resource configuration.",
	},
	"probe": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsProbeAttributes,
		Description: "Probe configuration.",
	},
	"dns_config": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsDnsConfigAttributes,
		Description: "Sandbox DNS configuration.",
	},
	"security_context": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsSecurityContextAttributes,
		Description: "Sandbox security context.",
	},
}

var agsWaaConfigurationAttributes = map[string]schema.Attribute{
	"image_id": schema.StringAttribute{
		Optional:    true,
		Description: "Custom WAA image id.",
	},
}

var agsOsWorldConfigurationAttributes = map[string]schema.Attribute{
	"version": schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Built-in OSWorld version. Valid values: `osworld1`, `osworld2`. Default value: `osworld1`.",
	},
}

var agsComputerConfigurationAttributes = map[string]schema.Attribute{
	"waa_configuration": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsWaaConfigurationAttributes,
		Description: "WAA sandbox configuration.",
	},
	"os_world_configuration": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsOsWorldConfigurationAttributes,
		Description: "Built-in OSWorld configuration.",
	},
}

var agsClsConfigAttributes = map[string]schema.Attribute{
	"topic_id": schema.StringAttribute{
		Optional:    true,
		Description: "Id of the CLS log topic used to push sandbox tool logs.",
	},
}

var agsLogSourcesAttributes = map[string]schema.Attribute{
	"files": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Log files to collect. They must live under the `/logs/` directory, sub directories are not " +
			"supported and at most 10 files are allowed.",
	},
}

var agsLogConfigurationAttributes = map[string]schema.Attribute{
	"cls_config": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsClsConfigAttributes,
		Description: "Log push configuration of CLS.",
	},
	"log_sources": schema.SingleNestedAttribute{
		Optional:    true,
		Attributes:  agsLogSourcesAttributes,
		Description: "Log source configuration.",
	},
}

var (
	agsVpcConfigAttrTypes            = attributeTypesOf(agsVpcConfigAttributes)
	agsNetworkConfigurationAttrTypes = attributeTypesOf(agsNetworkConfigurationAttributes)
	agsTagAttrTypes                  = attributeTypesOf(agsTagAttributes)
	agsCosStorageSourceAttrTypes     = attributeTypesOf(agsCosStorageSourceAttributes)
	agsImageStorageSourceAttrTypes   = attributeTypesOf(agsImageStorageSourceAttributes)
	agsCfsStorageSourceAttrTypes     = attributeTypesOf(agsCfsStorageSourceAttributes)
	agsAgentBucketSourceAttrTypes    = attributeTypesOf(agsAgentBucketStorageSourceAttributes)
	agsStorageSourceAttrTypes        = attributeTypesOf(agsStorageSourceAttributes)
	agsVolumeAttrTypes               = attributeTypesOf(agsVolumeAttributes)
	agsVolumeTemplateAttrTypes       = attributeTypesOf(agsVolumeTemplateAttributes)
	agsStorageMountAttrTypes         = attributeTypesOf(agsStorageMountAttributes)
	agsEnvVarAttrTypes               = attributeTypesOf(agsEnvVarAttributes)
	agsPortAttrTypes                 = attributeTypesOf(agsPortAttributes)
	agsResourceAttrTypes             = attributeTypesOf(agsResourceAttributes)
	agsHttpGetAttrTypes              = attributeTypesOf(agsHttpGetAttributes)
	agsProbeAttrTypes                = attributeTypesOf(agsProbeAttributes)
	agsDnsConfigAttrTypes            = attributeTypesOf(agsDnsConfigAttributes)
	agsCapabilitiesAttrTypes         = attributeTypesOf(agsCapabilitiesAttributes)
	agsSecurityContextAttrTypes      = attributeTypesOf(agsSecurityContextAttributes)
	agsCustomConfigAttrTypes         = attributeTypesOf(agsCustomConfigurationAttributes)
	agsWaaConfigurationAttrTypes     = attributeTypesOf(agsWaaConfigurationAttributes)
	agsOsWorldConfigurationAttrTypes = attributeTypesOf(agsOsWorldConfigurationAttributes)
	agsComputerConfigAttrTypes       = attributeTypesOf(agsComputerConfigurationAttributes)
	agsClsConfigAttrTypes            = attributeTypesOf(agsClsConfigAttributes)
	agsLogSourcesAttrTypes           = attributeTypesOf(agsLogSourcesAttributes)
	agsLogConfigurationAttrTypes     = attributeTypesOf(agsLogConfigurationAttributes)

	agsTagObjectType          = types.ObjectType{AttrTypes: agsTagAttrTypes}
	agsStorageMountObjectType = types.ObjectType{AttrTypes: agsStorageMountAttrTypes}
	agsEnvVarObjectType       = types.ObjectType{AttrTypes: agsEnvVarAttrTypes}
	agsPortObjectType         = types.ObjectType{AttrTypes: agsPortAttrTypes}
)

// attributeTypesOf mirrors the value types of a nested attribute map so the
// conversion helpers stay in sync with the schema definition.
func attributeTypesOf(attributes map[string]schema.Attribute) map[string]attr.Type {
	result := make(map[string]attr.Type, len(attributes))
	for name, attribute := range attributes {
		result[name] = attribute.GetType()
	}
	return result
}

func (r *AgsSandboxToolResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_tool"
}

func (r *AgsSandboxToolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create an AGS sandbox tool.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the resource, which is the same as `tool_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the sandbox tool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_name": schema.StringAttribute{
				Required: true,
				Description: "Name of the sandbox tool. Length 1-50 characters and it must be unique under the same " +
					"AppId. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tool_type": schema.StringAttribute{
				Required: true,
				Description: "Type of the sandbox tool. Valid values: `browser`, `code-interpreter`, `mobile`, " +
					"`osworld`, `custom`, `swebench`, `aio`, `android-world`, `waa`. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_configuration": schema.SingleNestedAttribute{
				Required:    true,
				Attributes:  agsNetworkConfigurationAttributes,
				Description: "Network configuration of the sandbox tool.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the sandbox tool. Maximum length 200 characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_timeout": schema.StringAttribute{
				Optional: true,
				Description: "Default timeout of sandbox instances created from this tool, for example `5m`, `300s` or " +
					"`1h`. Maximum 24 hours. When it is not specified, the system default of 5 minutes is used. " +
					"It cannot be changed after creation and it is not reported by the describe API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.ListNestedAttribute{
				Optional:     true,
				Computed:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsTagAttributes},
				Description:  "Tags of the sandbox tool.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"client_token": schema.StringAttribute{
				Optional: true,
				Description: "Idempotency token used by the create API. Maximum length 64 characters. " +
					"It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_arn": schema.StringAttribute{
				Optional:    true,
				Description: "Arn of the role bound to the sandbox tool. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_mounts": schema.ListNestedAttribute{
				Optional:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsStorageMountAttributes},
				Description:  "Storage mount configuration of sandbox instances. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"custom_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Attributes:  agsCustomConfigurationAttributes,
				Description: "Custom configuration of the sandbox tool.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"computer_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Attributes:  agsComputerConfigurationAttributes,
				Description: "Configuration of desktop-like sandbox tools.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"log_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Attributes:  agsLogConfigurationAttributes,
				Description: "Log push configuration of the sandbox tool. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"persistent": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Description: "Whether the sandbox tool is persistent. The cloud side applies a default value when " +
					"it is not specified. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Status of the sandbox tool. Valid values: `CREATING`, `ACTIVE`, `DELETING`, " +
					"`FAILED`.",
			},
			"status_reason": schema.StringAttribute{
				Computed:    true,
				Description: "Reason why the sandbox tool is in its current status.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the sandbox tool in ISO8601 format.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the sandbox tool in ISO8601 format.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *AgsSandboxToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsSandboxToolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, diags := r.buildCreateRequest(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	toolId, err := service.CreateSandboxTool(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AGS sandbox tool",
			fmt.Sprintf("CreateSandboxTool failed: %v", err),
		)
		return
	}

	plan.Id = types.StringValue(toolId)
	plan.ToolId = types.StringValue(toolId)
	plan.Status = types.StringNull()
	plan.StatusReason = types.StringNull()
	plan.CreateTime = types.StringNull()
	plan.UpdateTime = types.StringNull()
	normalizeUnknownAttributes(&plan)

	// Persist the identity before waiting: if the wait below fails, the resource
	// is still in state and `terraform destroy` can delete it.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.waitSandboxToolActive(ctx, service, toolId, tccommon.WriteRetryTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for AGS sandbox tool to become active",
			fmt.Sprintf("sandbox tool %s was created but is not active yet: %v", toolId, err),
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
func (r *AgsSandboxToolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsSandboxToolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	toolId, diags := sandboxToolIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	item, err := service.DescribeSandboxToolById(ctx, toolId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox tool",
			fmt.Sprintf("DescribeSandboxToolList failed for %s: %v", toolId, err),
		)
		return
	}

	// A nil result means the query succeeded and the tool no longer exists. A
	// nil API response is handled as an error inside the service and never
	// reaches this branch, so state is not dropped on a failed request.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"AGS sandbox tool not found",
			fmt.Sprintf("sandbox tool %s does not exist any more, removing it from state.", toolId),
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
func (r *AgsSandboxToolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgsSandboxToolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AgsSandboxToolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	toolId, diags := sandboxToolIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, diags := r.buildUpdateRequest(ctx, &plan, toolId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	if err := service.UpdateSandboxTool(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating AGS sandbox tool",
			fmt.Sprintf("UpdateSandboxTool failed for %s: %v", toolId, err),
		)
		return
	}

	item, err := service.DescribeSandboxToolById(ctx, toolId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox tool after update",
			fmt.Sprintf("DescribeSandboxToolList failed for %s: %v", toolId, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"AGS sandbox tool not readable after update",
			fmt.Sprintf("sandbox tool %s was updated but cannot be read back", toolId),
		)
		return
	}

	plan.Id = state.Id
	plan.ToolId = state.ToolId
	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`. It only deletes the sandbox tool itself;
// it does not touch sandbox instances or any storage that the tool references.
func (r *AgsSandboxToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsSandboxToolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	toolId, diags := sandboxToolIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	if err := service.DeleteSandboxTool(ctx, toolId); err != nil {
		// The API documentation does not define an "already deleted" error code,
		// so existence is checked again instead of guessing codes: a tool that is
		// really gone makes the delete idempotent, everything else is an error.
		item, describeErr := service.DescribeSandboxToolById(ctx, toolId)
		if describeErr != nil || item != nil {
			resp.Diagnostics.AddError(
				"Error deleting AGS sandbox tool",
				fmt.Sprintf("DeleteSandboxTool failed for %s: %v", toolId, err),
			)
		}
		return
	}

	// The tool may stay visible in the deleting state for a while. Wait until it
	// is gone so that a later apply can reuse the same tool name.
	if err := r.waitSandboxToolGone(ctx, service, toolId, tccommon.WriteRetryTimeout); err != nil {
		resp.Diagnostics.AddWarning(
			"AGS sandbox tool is still being deleted",
			fmt.Sprintf("sandbox tool %s was deleted but is still visible: %v", toolId, err),
		)
	}
}

// ImportState supports `terraform import` with the sandbox tool id.
func (r *AgsSandboxToolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Invalid import id",
			"The import id must be the id of an existing sandbox tool.",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tool_id"), req.ID)...)
}

// buildCreateRequest maps the plan onto a CreateSandboxTool request.
func (r *AgsSandboxToolResource) buildCreateRequest(ctx context.Context, plan *AgsSandboxToolResourceModel) (*agsv20250920.CreateSandboxToolRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	request := agsv20250920.NewCreateSandboxToolRequest()
	if !plan.ToolName.IsNull() && !plan.ToolName.IsUnknown() {
		request.ToolName = plan.ToolName.ValueStringPointer()
	}
	if !plan.ToolType.IsNull() && !plan.ToolType.IsUnknown() {
		request.ToolType = plan.ToolType.ValueStringPointer()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}
	if !plan.DefaultTimeout.IsNull() && !plan.DefaultTimeout.IsUnknown() {
		request.DefaultTimeout = plan.DefaultTimeout.ValueStringPointer()
	}
	if !plan.ClientToken.IsNull() && !plan.ClientToken.IsUnknown() {
		request.ClientToken = plan.ClientToken.ValueStringPointer()
	}
	if !plan.RoleArn.IsNull() && !plan.RoleArn.IsUnknown() {
		request.RoleArn = plan.RoleArn.ValueStringPointer()
	}
	if !plan.Persistent.IsNull() && !plan.Persistent.IsUnknown() {
		request.Persistent = plan.Persistent.ValueBoolPointer()
	}

	networkConfiguration, networkDiags := expandNetworkConfiguration(ctx, plan.NetworkConfiguration)
	diags.Append(networkDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.NetworkConfiguration = networkConfiguration

	tags, tagDiags := expandTags(ctx, plan.Tags)
	diags.Append(tagDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.Tags = tags

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

	computerConfiguration, computerDiags := expandComputerConfiguration(ctx, plan.ComputerConfiguration)
	diags.Append(computerDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.ComputerConfiguration = computerConfiguration

	logConfiguration, logDiags := expandLogConfiguration(ctx, plan.LogConfiguration)
	diags.Append(logDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.LogConfiguration = logConfiguration

	return request, diags
}

// buildUpdateRequest maps the plan onto an UpdateSandboxTool request.
//
// UpdateSandboxTool only accepts Description, NetworkConfiguration, Tags,
// CustomConfiguration and ComputerConfiguration; every other attribute is
// create-only and marked with RequiresReplace in the schema.
//
// The request structs omit empty values on the wire, so "attribute absent" is
// transmitted as "no change" and the cloud keeps the current value. The
// optional blocks below are therefore Optional+Computed with
// UseStateForUnknown: removing them from the configuration keeps the value
// reported by the cloud instead of producing a permanent diff.
func (r *AgsSandboxToolResource) buildUpdateRequest(ctx context.Context, plan *AgsSandboxToolResourceModel, toolId string) (*agsv20250920.UpdateSandboxToolRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	request := agsv20250920.NewUpdateSandboxToolRequest()
	request.ToolId = &toolId

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}

	networkConfiguration, networkDiags := expandNetworkConfiguration(ctx, plan.NetworkConfiguration)
	diags.Append(networkDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.NetworkConfiguration = networkConfiguration

	tags, tagDiags := expandTags(ctx, plan.Tags)
	diags.Append(tagDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.Tags = tags

	customConfiguration, customDiags := expandCustomConfiguration(ctx, plan.CustomConfiguration)
	diags.Append(customDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.CustomConfiguration = customConfiguration

	computerConfiguration, computerDiags := expandComputerConfiguration(ctx, plan.ComputerConfiguration)
	diags.Append(computerDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.ComputerConfiguration = computerConfiguration

	return request, diags
}

// refreshModel copies the cloud state into the model. Only fields reported by
// the cloud are overwritten, so create-only attributes that cannot be read back
// (default_timeout, client_token) keep the value already stored in the model.
func (r *AgsSandboxToolResource) refreshModel(ctx context.Context, model *AgsSandboxToolResourceModel, item *agsv20250920.SandboxTool) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		return diags
	}

	if item.ToolId != nil && *item.ToolId != "" {
		model.ToolId = types.StringValue(*item.ToolId)
		model.Id = types.StringValue(*item.ToolId)
	}
	if item.ToolName != nil {
		model.ToolName = types.StringValue(*item.ToolName)
	}
	if item.ToolType != nil {
		model.ToolType = types.StringValue(*item.ToolType)
	}

	model.Description = stringPointerValue(item.Description)
	model.RoleArn = stringPointerValue(item.RoleArn)
	model.Status = stringPointerValue(item.Status)
	model.StatusReason = stringPointerValue(item.StatusReason)
	model.CreateTime = stringPointerValue(item.CreateTime)
	model.UpdateTime = stringPointerValue(item.UpdateTime)
	if item.Persistent != nil {
		model.Persistent = types.BoolValue(*item.Persistent)
	}

	if item.NetworkConfiguration != nil {
		networkConfiguration, networkDiags := flattenNetworkConfiguration(ctx, item.NetworkConfiguration)
		diags.Append(networkDiags...)
		if diags.HasError() {
			return diags
		}
		model.NetworkConfiguration = networkConfiguration
	}

	tags, tagDiags := flattenTags(ctx, item.Tags)
	diags.Append(tagDiags...)
	if diags.HasError() {
		return diags
	}
	model.Tags = tags

	storageMounts, storageDiags := flattenStorageMounts(ctx, item.StorageMounts)
	diags.Append(storageDiags...)
	if diags.HasError() {
		return diags
	}
	model.StorageMounts = storageMounts

	customConfiguration, customDiags := flattenCustomConfiguration(ctx, item.CustomConfiguration)
	diags.Append(customDiags...)
	if diags.HasError() {
		return diags
	}
	model.CustomConfiguration = customConfiguration

	computerConfiguration, computerDiags := flattenComputerConfiguration(ctx, item.ComputerConfiguration)
	diags.Append(computerDiags...)
	if diags.HasError() {
		return diags
	}
	model.ComputerConfiguration = computerConfiguration

	logConfiguration, logDiags := flattenLogConfiguration(ctx, item.LogConfiguration)
	diags.Append(logDiags...)
	if diags.HasError() {
		return diags
	}
	model.LogConfiguration = logConfiguration

	return diags
}

// waitSandboxToolActive waits until the sandbox tool leaves the creating state.
// DescribeSandboxToolList documents CREATING, ACTIVE, DELETING and FAILED as the
// possible status values.
func (r *AgsSandboxToolResource) waitSandboxToolActive(ctx context.Context, service AgsService, toolId string, timeout time.Duration) (*agsv20250920.SandboxTool, error) {
	deadline := time.Now().Add(timeout)
	for {
		item, err := service.DescribeSandboxToolById(ctx, toolId)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, fmt.Errorf("sandbox tool %s disappeared while waiting for it to become active", toolId)
		}

		status := ""
		if item.Status != nil {
			status = *item.Status
		}
		switch status {
		case AGS_SANDBOX_TOOL_STATUS_ACTIVE:
			return item, nil
		case AGS_SANDBOX_TOOL_STATUS_FAILED:
			reason := ""
			if item.StatusReason != nil {
				reason = *item.StatusReason
			}
			return nil, fmt.Errorf("sandbox tool %s is in FAILED status: %s", toolId, reason)
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for sandbox tool %s to become ACTIVE, current status is %q", toolId, status)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(agsSandboxToolPollInterval):
		}
	}
}

// waitSandboxToolGone waits until the deleted sandbox tool is no longer visible.
func (r *AgsSandboxToolResource) waitSandboxToolGone(ctx context.Context, service AgsService, toolId string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		item, err := service.DescribeSandboxToolById(ctx, toolId)
		if err != nil {
			return err
		}
		if item == nil {
			return nil
		}

		if time.Now().After(deadline) {
			status := ""
			if item.Status != nil {
				status = *item.Status
			}
			return fmt.Errorf("sandbox tool %s is still present with status %q after %s", toolId, status, timeout)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(agsSandboxToolPollInterval):
		}
	}
}

// sandboxToolIdOf returns the identity stored in the model.
func sandboxToolIdOf(model *AgsSandboxToolResourceModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !model.ToolId.IsNull() && !model.ToolId.IsUnknown() && model.ToolId.ValueString() != "" {
		return model.ToolId.ValueString(), diags
	}
	if !model.Id.IsNull() && !model.Id.IsUnknown() && model.Id.ValueString() != "" {
		return model.Id.ValueString(), diags
	}

	diags.AddError(
		"Invalid AGS sandbox tool state",
		"The resource state does not contain a sandbox tool id.",
	)
	return "", diags
}

// normalizeUnknownAttributes replaces the optional+computed attributes that are
// still unknown at create time with null, so the identity can be stored in state
// before the cloud has reported those values.
func normalizeUnknownAttributes(model *AgsSandboxToolResourceModel) {
	if model.Description.IsUnknown() {
		model.Description = types.StringNull()
	}
	if model.Tags.IsUnknown() {
		model.Tags = types.ListNull(agsTagObjectType)
	}
	if model.CustomConfiguration.IsUnknown() {
		model.CustomConfiguration = types.ObjectNull(agsCustomConfigAttrTypes)
	}
	if model.ComputerConfiguration.IsUnknown() {
		model.ComputerConfiguration = types.ObjectNull(agsComputerConfigAttrTypes)
	}
	if model.Persistent.IsUnknown() {
		model.Persistent = types.BoolNull()
	}
}

func stringPointerValue(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(*value)
}

// expandStringSet converts a set attribute into the SDK pointer slice form. A
// null or unknown set is transmitted as "absent".
func expandStringSet(ctx context.Context, value types.Set) ([]*string, diag.Diagnostics) {
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

	return helperStrings(elements), diags
}

func flattenStringSet(ctx context.Context, values []*string) (types.Set, diag.Diagnostics) {
	if len(values) == 0 {
		return types.SetNull(types.StringType), nil
	}
	return types.SetValueFrom(ctx, types.StringType, helperPStrings(values))
}

// helperStrings converts a string slice into the SDK pointer slice form.
func helperStrings(values []string) []*string {
	result := make([]*string, 0, len(values))
	for i := range values {
		value := values[i]
		result = append(result, &value)
	}
	return result
}

// helperPStrings converts an SDK pointer slice into a string slice, skipping
// nil entries so a malformed response cannot panic the provider.
func helperPStrings(values []*string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		result = append(result, *value)
	}
	return result
}

func expandNetworkConfiguration(ctx context.Context, object types.Object) (*agsv20250920.NetworkConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsNetworkConfigurationModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	networkConfiguration := &agsv20250920.NetworkConfiguration{}
	if !model.NetworkMode.IsNull() && !model.NetworkMode.IsUnknown() {
		networkConfiguration.NetworkMode = model.NetworkMode.ValueStringPointer()
	}

	vpcConfig, vpcDiags := expandVpcConfig(ctx, model.VpcConfig)
	diags.Append(vpcDiags...)
	if diags.HasError() {
		return nil, diags
	}
	networkConfiguration.VpcConfig = vpcConfig

	return networkConfiguration, diags
}

func flattenNetworkConfiguration(ctx context.Context, item *agsv20250920.NetworkConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsNetworkConfigurationAttrTypes), nil
	}

	model := agsNetworkConfigurationModel{
		NetworkMode: types.StringNull(),
		VpcConfig:   types.ObjectNull(agsVpcConfigAttrTypes),
	}
	if item.NetworkMode != nil {
		model.NetworkMode = types.StringValue(*item.NetworkMode)
	}

	vpcConfig, diags := flattenVpcConfig(ctx, item.VpcConfig)
	if diags.HasError() {
		return types.ObjectNull(agsNetworkConfigurationAttrTypes), diags
	}
	model.VpcConfig = vpcConfig

	return types.ObjectValueFrom(ctx, agsNetworkConfigurationAttrTypes, model)
}

func expandVpcConfig(ctx context.Context, object types.Object) (*agsv20250920.VPCConfig, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsVpcConfigModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	vpcConfig := &agsv20250920.VPCConfig{}

	subnetIds, subnetDiags := expandStringSet(ctx, model.SubnetIds)
	diags.Append(subnetDiags...)
	if diags.HasError() {
		return nil, diags
	}
	vpcConfig.SubnetIds = subnetIds

	securityGroupIds, securityDiags := expandStringSet(ctx, model.SecurityGroupIds)
	diags.Append(securityDiags...)
	if diags.HasError() {
		return nil, diags
	}
	vpcConfig.SecurityGroupIds = securityGroupIds

	return vpcConfig, diags
}

func flattenVpcConfig(ctx context.Context, item *agsv20250920.VPCConfig) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsVpcConfigAttrTypes), nil
	}

	subnetIds, subnetDiags := flattenStringSet(ctx, item.SubnetIds)
	if subnetDiags.HasError() {
		return types.ObjectNull(agsVpcConfigAttrTypes), subnetDiags
	}
	securityGroupIds, securityDiags := flattenStringSet(ctx, item.SecurityGroupIds)
	if securityDiags.HasError() {
		return types.ObjectNull(agsVpcConfigAttrTypes), securityDiags
	}

	model := agsVpcConfigModel{
		SubnetIds:        subnetIds,
		SecurityGroupIds: securityGroupIds,
	}
	return types.ObjectValueFrom(ctx, agsVpcConfigAttrTypes, model)
}

func expandTags(ctx context.Context, list types.List) ([]*agsv20250920.Tag, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsTagModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	tags := make([]*agsv20250920.Tag, 0, len(models))
	for _, model := range models {
		tag := &agsv20250920.Tag{}
		if !model.Key.IsNull() && !model.Key.IsUnknown() {
			tag.Key = model.Key.ValueStringPointer()
		}
		if !model.Value.IsNull() && !model.Value.IsUnknown() {
			tag.Value = model.Value.ValueStringPointer()
		}
		tags = append(tags, tag)
	}

	return tags, diags
}

func flattenTags(ctx context.Context, items []*agsv20250920.Tag) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsTagObjectType), nil
	}

	models := make([]agsTagModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsTagModel{
			Key:   stringPointerValue(item.Key),
			Value: stringPointerValue(item.Value),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsTagObjectType), nil
	}

	return types.ListValueFrom(ctx, agsTagObjectType, models)
}

func expandStorageMounts(ctx context.Context, list types.List) ([]*agsv20250920.StorageMount, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsStorageMountModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	storageMounts := make([]*agsv20250920.StorageMount, 0, len(models))
	for _, model := range models {
		storageMount := &agsv20250920.StorageMount{}

		if !model.Name.IsNull() && !model.Name.IsUnknown() {
			storageMount.Name = model.Name.ValueStringPointer()
		}
		if !model.MountPath.IsNull() && !model.MountPath.IsUnknown() {
			storageMount.MountPath = model.MountPath.ValueStringPointer()
		}
		if !model.ReadOnly.IsNull() && !model.ReadOnly.IsUnknown() {
			storageMount.ReadOnly = model.ReadOnly.ValueBoolPointer()
		}

		storageSource, sourceDiags := expandStorageSource(ctx, model.StorageSource)
		diags.Append(sourceDiags...)
		if diags.HasError() {
			return nil, diags
		}
		storageMount.StorageSource = storageSource

		volume, volumeDiags := expandVolume(ctx, model.Volume)
		diags.Append(volumeDiags...)
		if diags.HasError() {
			return nil, diags
		}
		storageMount.Volume = volume

		volumeTemplate, templateDiags := expandVolumeTemplate(ctx, model.VolumeTemplate)
		diags.Append(templateDiags...)
		if diags.HasError() {
			return nil, diags
		}
		storageMount.VolumeTemplate = volumeTemplate

		storageMounts = append(storageMounts, storageMount)
	}

	return storageMounts, diags
}

func flattenStorageMounts(ctx context.Context, items []*agsv20250920.StorageMount) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsStorageMountObjectType), nil
	}

	var diags diag.Diagnostics
	models := make([]agsStorageMountModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		storageSource, sourceDiags := flattenStorageSource(ctx, item.StorageSource)
		diags.Append(sourceDiags...)
		if diags.HasError() {
			return types.ListNull(agsStorageMountObjectType), diags
		}

		volume, volumeDiags := flattenVolume(ctx, item.Volume)
		diags.Append(volumeDiags...)
		if diags.HasError() {
			return types.ListNull(agsStorageMountObjectType), diags
		}

		volumeTemplate, templateDiags := flattenVolumeTemplate(ctx, item.VolumeTemplate)
		diags.Append(templateDiags...)
		if diags.HasError() {
			return types.ListNull(agsStorageMountObjectType), diags
		}

		models = append(models, agsStorageMountModel{
			Name:           stringPointerValue(item.Name),
			StorageSource:  storageSource,
			MountPath:      stringPointerValue(item.MountPath),
			ReadOnly:       boolPointerValue(item.ReadOnly),
			Volume:         volume,
			VolumeTemplate: volumeTemplate,
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsStorageMountObjectType), diags
	}

	list, listDiags := types.ListValueFrom(ctx, agsStorageMountObjectType, models)
	diags.Append(listDiags...)
	return list, diags
}

func expandStorageSource(ctx context.Context, object types.Object) (*agsv20250920.StorageSource, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsStorageSourceModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	storageSource := &agsv20250920.StorageSource{}

	cos, cosDiags := expandCosStorageSource(ctx, model.Cos)
	diags.Append(cosDiags...)
	if diags.HasError() {
		return nil, diags
	}
	storageSource.Cos = cos

	image, imageDiags := expandImageStorageSource(ctx, model.Image)
	diags.Append(imageDiags...)
	if diags.HasError() {
		return nil, diags
	}
	storageSource.Image = image

	cfs, cfsDiags := expandCfsStorageSource(ctx, model.Cfs)
	diags.Append(cfsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	storageSource.Cfs = cfs

	agentBucket, agentBucketDiags := expandAgentBucketStorageSource(ctx, model.AgentBucket)
	diags.Append(agentBucketDiags...)
	if diags.HasError() {
		return nil, diags
	}
	storageSource.AgentBucket = agentBucket

	return storageSource, diags
}

func flattenStorageSource(ctx context.Context, item *agsv20250920.StorageSource) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsStorageSourceAttrTypes), nil
	}

	var diags diag.Diagnostics

	cos, cosDiags := flattenCosStorageSource(ctx, item.Cos)
	diags.Append(cosDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsStorageSourceAttrTypes), diags
	}

	image, imageDiags := flattenImageStorageSource(ctx, item.Image)
	diags.Append(imageDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsStorageSourceAttrTypes), diags
	}

	cfs, cfsDiags := flattenCfsStorageSource(ctx, item.Cfs)
	diags.Append(cfsDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsStorageSourceAttrTypes), diags
	}

	agentBucket, agentBucketDiags := flattenAgentBucketStorageSource(ctx, item.AgentBucket)
	diags.Append(agentBucketDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsStorageSourceAttrTypes), diags
	}

	model := agsStorageSourceModel{
		Cos:         cos,
		Image:       image,
		Cfs:         cfs,
		AgentBucket: agentBucket,
	}
	return types.ObjectValueFrom(ctx, agsStorageSourceAttrTypes, model)
}

func expandCosStorageSource(ctx context.Context, object types.Object) (*agsv20250920.CosStorageSource, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsCosStorageSourceModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	cos := &agsv20250920.CosStorageSource{}
	if !model.Endpoint.IsNull() && !model.Endpoint.IsUnknown() {
		cos.Endpoint = model.Endpoint.ValueStringPointer()
	}
	if !model.BucketName.IsNull() && !model.BucketName.IsUnknown() {
		cos.BucketName = model.BucketName.ValueStringPointer()
	}
	if !model.BucketPath.IsNull() && !model.BucketPath.IsUnknown() {
		cos.BucketPath = model.BucketPath.ValueStringPointer()
	}

	return cos, diags
}

func flattenCosStorageSource(ctx context.Context, item *agsv20250920.CosStorageSource) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsCosStorageSourceAttrTypes), nil
	}

	model := agsCosStorageSourceModel{
		Endpoint:   stringPointerValue(item.Endpoint),
		BucketName: stringPointerValue(item.BucketName),
		BucketPath: stringPointerValue(item.BucketPath),
	}
	return types.ObjectValueFrom(ctx, agsCosStorageSourceAttrTypes, model)
}

func expandImageStorageSource(ctx context.Context, object types.Object) (*agsv20250920.ImageStorageSource, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsImageStorageSourceModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	image := &agsv20250920.ImageStorageSource{}
	if !model.Reference.IsNull() && !model.Reference.IsUnknown() {
		image.Reference = model.Reference.ValueStringPointer()
	}
	if !model.ImageRegistryType.IsNull() && !model.ImageRegistryType.IsUnknown() {
		image.ImageRegistryType = model.ImageRegistryType.ValueStringPointer()
	}
	if !model.SubPath.IsNull() && !model.SubPath.IsUnknown() {
		image.SubPath = model.SubPath.ValueStringPointer()
	}
	// Digest is resolved by the cloud side and is intentionally not sent.

	return image, diags
}

func flattenImageStorageSource(ctx context.Context, item *agsv20250920.ImageStorageSource) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsImageStorageSourceAttrTypes), nil
	}

	model := agsImageStorageSourceModel{
		Reference:         stringPointerValue(item.Reference),
		ImageRegistryType: stringPointerValue(item.ImageRegistryType),
		SubPath:           stringPointerValue(item.SubPath),
		Digest:            stringPointerValue(item.Digest),
	}
	return types.ObjectValueFrom(ctx, agsImageStorageSourceAttrTypes, model)
}

func expandCfsStorageSource(ctx context.Context, object types.Object) (*agsv20250920.CfsStorageSource, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsCfsStorageSourceModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	cfs := &agsv20250920.CfsStorageSource{}
	if !model.FileSystemId.IsNull() && !model.FileSystemId.IsUnknown() {
		cfs.FileSystemId = model.FileSystemId.ValueStringPointer()
	}
	if !model.Path.IsNull() && !model.Path.IsUnknown() {
		cfs.Path = model.Path.ValueStringPointer()
	}

	return cfs, diags
}

func flattenCfsStorageSource(ctx context.Context, item *agsv20250920.CfsStorageSource) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsCfsStorageSourceAttrTypes), nil
	}

	model := agsCfsStorageSourceModel{
		FileSystemId: stringPointerValue(item.FileSystemId),
		Path:         stringPointerValue(item.Path),
	}
	return types.ObjectValueFrom(ctx, agsCfsStorageSourceAttrTypes, model)
}

func expandAgentBucketStorageSource(ctx context.Context, object types.Object) (*agsv20250920.AgentBucketStorageSource, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsAgentBucketStorageSourceModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	agentBucket := &agsv20250920.AgentBucketStorageSource{}
	if !model.LibraryId.IsNull() && !model.LibraryId.IsUnknown() {
		agentBucket.LibraryId = model.LibraryId.ValueStringPointer()
	}
	if !model.SpaceId.IsNull() && !model.SpaceId.IsUnknown() {
		agentBucket.SpaceId = model.SpaceId.ValueStringPointer()
	}
	if !model.AccessDomain.IsNull() && !model.AccessDomain.IsUnknown() {
		agentBucket.AccessDomain = model.AccessDomain.ValueStringPointer()
	}

	return agentBucket, diags
}

func flattenAgentBucketStorageSource(ctx context.Context, item *agsv20250920.AgentBucketStorageSource) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsAgentBucketSourceAttrTypes), nil
	}

	model := agsAgentBucketStorageSourceModel{
		LibraryId:    stringPointerValue(item.LibraryId),
		SpaceId:      stringPointerValue(item.SpaceId),
		AccessDomain: stringPointerValue(item.AccessDomain),
	}
	return types.ObjectValueFrom(ctx, agsAgentBucketSourceAttrTypes, model)
}

func expandVolume(ctx context.Context, object types.Object) (*agsv20250920.VolumeRef, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsVolumeModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	volume := &agsv20250920.VolumeRef{}
	if !model.VolumeName.IsNull() && !model.VolumeName.IsUnknown() {
		volume.VolumeName = model.VolumeName.ValueStringPointer()
	}
	if !model.VolumeId.IsNull() && !model.VolumeId.IsUnknown() {
		volume.VolumeId = model.VolumeId.ValueStringPointer()
	}

	return volume, diags
}

func flattenVolume(ctx context.Context, item *agsv20250920.VolumeRef) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsVolumeAttrTypes), nil
	}

	model := agsVolumeModel{
		VolumeName: stringPointerValue(item.VolumeName),
		VolumeId:   stringPointerValue(item.VolumeId),
	}
	return types.ObjectValueFrom(ctx, agsVolumeAttrTypes, model)
}

func expandVolumeTemplate(ctx context.Context, object types.Object) (*agsv20250920.VolumeTemplateRef, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsVolumeTemplateModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	volumeTemplate := &agsv20250920.VolumeTemplateRef{}
	if !model.VolumeTemplateName.IsNull() && !model.VolumeTemplateName.IsUnknown() {
		volumeTemplate.VolumeTemplateName = model.VolumeTemplateName.ValueStringPointer()
	}
	if !model.ReuseKey.IsNull() && !model.ReuseKey.IsUnknown() {
		volumeTemplate.ReuseKey = model.ReuseKey.ValueStringPointer()
	}

	return volumeTemplate, diags
}

func flattenVolumeTemplate(ctx context.Context, item *agsv20250920.VolumeTemplateRef) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsVolumeTemplateAttrTypes), nil
	}

	model := agsVolumeTemplateModel{
		VolumeTemplateName: stringPointerValue(item.VolumeTemplateName),
		ReuseKey:           stringPointerValue(item.ReuseKey),
	}
	return types.ObjectValueFrom(ctx, agsVolumeTemplateAttrTypes, model)
}

func expandCustomConfiguration(ctx context.Context, object types.Object) (*agsv20250920.CustomConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsCustomConfigurationModel
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

	command, commandDiags := expandStringSet(ctx, model.Command)
	diags.Append(commandDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Command = command

	args, argsDiags := expandStringSet(ctx, model.Args)
	diags.Append(argsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Args = args

	env, envDiags := expandEnvVars(ctx, model.Env)
	diags.Append(envDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Env = env

	ports, portDiags := expandPorts(ctx, model.Ports)
	diags.Append(portDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Ports = ports

	resources, resourceDiags := expandResources(ctx, model.Resources)
	diags.Append(resourceDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Resources = resources

	probe, probeDiags := expandProbe(ctx, model.Probe)
	diags.Append(probeDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.Probe = probe

	dnsConfig, dnsDiags := expandDnsConfig(ctx, model.DnsConfig)
	diags.Append(dnsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.DNSConfig = dnsConfig

	securityContext, securityDiags := expandSecurityContext(ctx, model.SecurityContext)
	diags.Append(securityDiags...)
	if diags.HasError() {
		return nil, diags
	}
	customConfiguration.SecurityContext = securityContext

	return customConfiguration, diags
}

func flattenCustomConfiguration(ctx context.Context, item *agsv20250920.CustomConfigurationDetail) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsCustomConfigAttrTypes), nil
	}

	var diags diag.Diagnostics

	command, commandDiags := flattenStringSet(ctx, item.Command)
	diags.Append(commandDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCustomConfigAttrTypes), diags
	}

	args, argsDiags := flattenStringSet(ctx, item.Args)
	diags.Append(argsDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCustomConfigAttrTypes), diags
	}

	env, envDiags := flattenEnvVars(ctx, item.Env)
	diags.Append(envDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCustomConfigAttrTypes), diags
	}

	ports, portDiags := flattenPorts(ctx, item.Ports)
	diags.Append(portDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCustomConfigAttrTypes), diags
	}

	resources, resourceDiags := flattenResources(ctx, item.Resources)
	diags.Append(resourceDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCustomConfigAttrTypes), diags
	}

	probe, probeDiags := flattenProbe(ctx, item.Probe)
	diags.Append(probeDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCustomConfigAttrTypes), diags
	}

	dnsConfig, dnsDiags := flattenDnsConfig(ctx, item.DNSConfig)
	diags.Append(dnsDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCustomConfigAttrTypes), diags
	}

	securityContext, securityDiags := flattenSecurityContext(ctx, item.SecurityContext)
	diags.Append(securityDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCustomConfigAttrTypes), diags
	}

	// ImageDigest is a read-only output of CustomConfigurationDetail and is not
	// part of the configurable schema, so it is not written back.
	model := agsCustomConfigurationModel{
		Image:             stringPointerValue(item.Image),
		ImageRegistryType: stringPointerValue(item.ImageRegistryType),
		Command:           command,
		Args:              args,
		Env:               env,
		Ports:             ports,
		Resources:         resources,
		Probe:             probe,
		DnsConfig:         dnsConfig,
		SecurityContext:   securityContext,
	}
	return types.ObjectValueFrom(ctx, agsCustomConfigAttrTypes, model)
}

func expandEnvVars(ctx context.Context, list types.List) ([]*agsv20250920.EnvVar, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsEnvVarModel
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

func flattenEnvVars(ctx context.Context, items []*agsv20250920.EnvVar) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsEnvVarObjectType), nil
	}

	models := make([]agsEnvVarModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsEnvVarModel{
			Name:  stringPointerValue(item.Name),
			Value: stringPointerValue(item.Value),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsEnvVarObjectType), nil
	}

	return types.ListValueFrom(ctx, agsEnvVarObjectType, models)
}

func expandPorts(ctx context.Context, list types.List) ([]*agsv20250920.PortConfiguration, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsPortModel
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

func flattenPorts(ctx context.Context, items []*agsv20250920.PortConfiguration) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsPortObjectType), nil
	}

	models := make([]agsPortModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsPortModel{
			Name:     stringPointerValue(item.Name),
			Port:     int64PointerValue(item.Port),
			Protocol: stringPointerValue(item.Protocol),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsPortObjectType), nil
	}

	return types.ListValueFrom(ctx, agsPortObjectType, models)
}

func expandResources(ctx context.Context, object types.Object) (*agsv20250920.ResourceConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsResourceModel
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

func flattenResources(ctx context.Context, item *agsv20250920.ResourceConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsResourceAttrTypes), nil
	}

	model := agsResourceModel{
		Cpu:     stringPointerValue(item.CPU),
		Memory:  stringPointerValue(item.Memory),
		Storage: stringPointerValue(item.Storage),
	}
	return types.ObjectValueFrom(ctx, agsResourceAttrTypes, model)
}

func expandProbe(ctx context.Context, object types.Object) (*agsv20250920.ProbeConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsProbeModel
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

	httpGet, httpGetDiags := expandHttpGet(ctx, model.HttpGet)
	diags.Append(httpGetDiags...)
	if diags.HasError() {
		return nil, diags
	}
	probe.HttpGet = httpGet

	return probe, diags
}

func flattenProbe(ctx context.Context, item *agsv20250920.ProbeConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsProbeAttrTypes), nil
	}

	httpGet, httpGetDiags := flattenHttpGet(ctx, item.HttpGet)
	if httpGetDiags.HasError() {
		return types.ObjectNull(agsProbeAttrTypes), httpGetDiags
	}

	model := agsProbeModel{
		HttpGet:          httpGet,
		ReadyTimeoutMs:   int64PointerValue(item.ReadyTimeoutMs),
		ProbeTimeoutMs:   int64PointerValue(item.ProbeTimeoutMs),
		ProbePeriodMs:    int64PointerValue(item.ProbePeriodMs),
		SuccessThreshold: int64PointerValue(item.SuccessThreshold),
		FailureThreshold: int64PointerValue(item.FailureThreshold),
	}
	return types.ObjectValueFrom(ctx, agsProbeAttrTypes, model)
}

func expandHttpGet(ctx context.Context, object types.Object) (*agsv20250920.HttpGetAction, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsHttpGetModel
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

func flattenHttpGet(ctx context.Context, item *agsv20250920.HttpGetAction) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsHttpGetAttrTypes), nil
	}

	model := agsHttpGetModel{
		Path:   stringPointerValue(item.Path),
		Port:   int64PointerValue(item.Port),
		Scheme: stringPointerValue(item.Scheme),
	}
	return types.ObjectValueFrom(ctx, agsHttpGetAttrTypes, model)
}

func expandDnsConfig(ctx context.Context, object types.Object) (*agsv20250920.DNSConfig, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsDnsConfigModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	dnsConfig := &agsv20250920.DNSConfig{}

	servers, serverDiags := expandStringSet(ctx, model.Servers)
	diags.Append(serverDiags...)
	if diags.HasError() {
		return nil, diags
	}
	dnsConfig.Servers = servers

	searches, searchDiags := expandStringSet(ctx, model.Searches)
	diags.Append(searchDiags...)
	if diags.HasError() {
		return nil, diags
	}
	dnsConfig.Searches = searches

	options, optionDiags := expandStringSet(ctx, model.Options)
	diags.Append(optionDiags...)
	if diags.HasError() {
		return nil, diags
	}
	dnsConfig.Options = options

	return dnsConfig, diags
}

func flattenDnsConfig(ctx context.Context, item *agsv20250920.DNSConfig) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsDnsConfigAttrTypes), nil
	}

	var diags diag.Diagnostics

	servers, serverDiags := flattenStringSet(ctx, item.Servers)
	diags.Append(serverDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsDnsConfigAttrTypes), diags
	}

	searches, searchDiags := flattenStringSet(ctx, item.Searches)
	diags.Append(searchDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsDnsConfigAttrTypes), diags
	}

	options, optionDiags := flattenStringSet(ctx, item.Options)
	diags.Append(optionDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsDnsConfigAttrTypes), diags
	}

	model := agsDnsConfigModel{
		Servers:  servers,
		Searches: searches,
		Options:  options,
	}
	return types.ObjectValueFrom(ctx, agsDnsConfigAttrTypes, model)
}

func expandSecurityContext(ctx context.Context, object types.Object) (*agsv20250920.SecurityContext, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsSecurityContextModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	securityContext := &agsv20250920.SecurityContext{}
	if !model.Privileged.IsNull() && !model.Privileged.IsUnknown() {
		securityContext.Privileged = model.Privileged.ValueBoolPointer()
	}

	capabilities, capabilityDiags := expandCapabilities(ctx, model.Capabilities)
	diags.Append(capabilityDiags...)
	if diags.HasError() {
		return nil, diags
	}
	securityContext.Capabilities = capabilities

	return securityContext, diags
}

func flattenSecurityContext(ctx context.Context, item *agsv20250920.SecurityContext) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsSecurityContextAttrTypes), nil
	}

	capabilities, capabilityDiags := flattenCapabilities(ctx, item.Capabilities)
	if capabilityDiags.HasError() {
		return types.ObjectNull(agsSecurityContextAttrTypes), capabilityDiags
	}

	model := agsSecurityContextModel{
		Privileged:   boolPointerValue(item.Privileged),
		Capabilities: capabilities,
	}
	return types.ObjectValueFrom(ctx, agsSecurityContextAttrTypes, model)
}

func expandCapabilities(ctx context.Context, object types.Object) (*agsv20250920.Capabilities, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsCapabilitiesModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	capabilities := &agsv20250920.Capabilities{}

	add, addDiags := expandStringSet(ctx, model.Add)
	diags.Append(addDiags...)
	if diags.HasError() {
		return nil, diags
	}
	capabilities.Add = add

	drop, dropDiags := expandStringSet(ctx, model.Drop)
	diags.Append(dropDiags...)
	if diags.HasError() {
		return nil, diags
	}
	capabilities.Drop = drop

	return capabilities, diags
}

func flattenCapabilities(ctx context.Context, item *agsv20250920.Capabilities) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsCapabilitiesAttrTypes), nil
	}

	var diags diag.Diagnostics

	add, addDiags := flattenStringSet(ctx, item.Add)
	diags.Append(addDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCapabilitiesAttrTypes), diags
	}

	drop, dropDiags := flattenStringSet(ctx, item.Drop)
	diags.Append(dropDiags...)
	if diags.HasError() {
		return types.ObjectNull(agsCapabilitiesAttrTypes), diags
	}

	model := agsCapabilitiesModel{
		Add:  add,
		Drop: drop,
	}
	return types.ObjectValueFrom(ctx, agsCapabilitiesAttrTypes, model)
}

func expandComputerConfiguration(ctx context.Context, object types.Object) (*agsv20250920.ComputerConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsComputerConfigurationModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	computerConfiguration := &agsv20250920.ComputerConfiguration{}

	waaConfiguration, waaDiags := expandWaaConfiguration(ctx, model.WaaConfiguration)
	diags.Append(waaDiags...)
	if diags.HasError() {
		return nil, diags
	}
	computerConfiguration.WAAConfiguration = waaConfiguration

	osWorldConfiguration, osWorldDiags := expandOsWorldConfiguration(ctx, model.OsWorldConfiguration)
	diags.Append(osWorldDiags...)
	if diags.HasError() {
		return nil, diags
	}
	computerConfiguration.OSWorldConfiguration = osWorldConfiguration

	return computerConfiguration, diags
}

func flattenComputerConfiguration(ctx context.Context, item *agsv20250920.ComputerConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsComputerConfigAttrTypes), nil
	}

	waaConfiguration, waaDiags := flattenWaaConfiguration(ctx, item.WAAConfiguration)
	if waaDiags.HasError() {
		return types.ObjectNull(agsComputerConfigAttrTypes), waaDiags
	}

	osWorldConfiguration, osWorldDiags := flattenOsWorldConfiguration(ctx, item.OSWorldConfiguration)
	if osWorldDiags.HasError() {
		return types.ObjectNull(agsComputerConfigAttrTypes), osWorldDiags
	}

	model := agsComputerConfigurationModel{
		WaaConfiguration:     waaConfiguration,
		OsWorldConfiguration: osWorldConfiguration,
	}
	return types.ObjectValueFrom(ctx, agsComputerConfigAttrTypes, model)
}

func expandWaaConfiguration(ctx context.Context, object types.Object) (*agsv20250920.WAAConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsWaaConfigurationModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	waaConfiguration := &agsv20250920.WAAConfiguration{}
	if !model.ImageId.IsNull() && !model.ImageId.IsUnknown() {
		waaConfiguration.ImageId = model.ImageId.ValueStringPointer()
	}

	return waaConfiguration, diags
}

func flattenWaaConfiguration(ctx context.Context, item *agsv20250920.WAAConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsWaaConfigurationAttrTypes), nil
	}

	model := agsWaaConfigurationModel{
		ImageId: stringPointerValue(item.ImageId),
	}
	return types.ObjectValueFrom(ctx, agsWaaConfigurationAttrTypes, model)
}

func expandOsWorldConfiguration(ctx context.Context, object types.Object) (*agsv20250920.OSWorldConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsOsWorldConfigurationModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	osWorldConfiguration := &agsv20250920.OSWorldConfiguration{}
	if !model.Version.IsNull() && !model.Version.IsUnknown() {
		osWorldConfiguration.Version = model.Version.ValueStringPointer()
	}

	return osWorldConfiguration, diags
}

func flattenOsWorldConfiguration(ctx context.Context, item *agsv20250920.OSWorldConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsOsWorldConfigurationAttrTypes), nil
	}

	model := agsOsWorldConfigurationModel{
		Version: stringPointerValue(item.Version),
	}
	return types.ObjectValueFrom(ctx, agsOsWorldConfigurationAttrTypes, model)
}

func expandLogConfiguration(ctx context.Context, object types.Object) (*agsv20250920.LogConfiguration, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsLogConfigurationModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	logConfiguration := &agsv20250920.LogConfiguration{}

	clsConfig, clsDiags := expandClsConfig(ctx, model.ClsConfig)
	diags.Append(clsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	logConfiguration.CLSConfig = clsConfig

	logSources, sourceDiags := expandLogSources(ctx, model.LogSources)
	diags.Append(sourceDiags...)
	if diags.HasError() {
		return nil, diags
	}
	logConfiguration.LogSources = logSources

	return logConfiguration, diags
}

func flattenLogConfiguration(ctx context.Context, item *agsv20250920.LogConfiguration) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsLogConfigurationAttrTypes), nil
	}

	clsConfig, clsDiags := flattenClsConfig(ctx, item.CLSConfig)
	if clsDiags.HasError() {
		return types.ObjectNull(agsLogConfigurationAttrTypes), clsDiags
	}

	logSources, sourceDiags := flattenLogSources(ctx, item.LogSources)
	if sourceDiags.HasError() {
		return types.ObjectNull(agsLogConfigurationAttrTypes), sourceDiags
	}

	model := agsLogConfigurationModel{
		ClsConfig:  clsConfig,
		LogSources: logSources,
	}
	return types.ObjectValueFrom(ctx, agsLogConfigurationAttrTypes, model)
}

func expandClsConfig(ctx context.Context, object types.Object) (*agsv20250920.CLSConfig, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsClsConfigModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	clsConfig := &agsv20250920.CLSConfig{}
	if !model.TopicId.IsNull() && !model.TopicId.IsUnknown() {
		clsConfig.TopicId = model.TopicId.ValueStringPointer()
	}

	return clsConfig, diags
}

func flattenClsConfig(ctx context.Context, item *agsv20250920.CLSConfig) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsClsConfigAttrTypes), nil
	}

	model := agsClsConfigModel{
		TopicId: stringPointerValue(item.TopicId),
	}
	return types.ObjectValueFrom(ctx, agsClsConfigAttrTypes, model)
}

func expandLogSources(ctx context.Context, object types.Object) (*agsv20250920.LogSources, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model agsLogSourcesModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	logSources := &agsv20250920.LogSources{}
	files, fileDiags := expandStringSet(ctx, model.Files)
	diags.Append(fileDiags...)
	if diags.HasError() {
		return nil, diags
	}
	logSources.Files = files

	return logSources, diags
}

func flattenLogSources(ctx context.Context, item *agsv20250920.LogSources) (types.Object, diag.Diagnostics) {
	if item == nil {
		return types.ObjectNull(agsLogSourcesAttrTypes), nil
	}

	files, fileDiags := flattenStringSet(ctx, item.Files)
	if fileDiags.HasError() {
		return types.ObjectNull(agsLogSourcesAttrTypes), fileDiags
	}

	model := agsLogSourcesModel{
		Files: files,
	}
	return types.ObjectValueFrom(ctx, agsLogSourcesAttrTypes, model)
}

func boolPointerValue(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

func int64PointerValue(value *int64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*value)
}

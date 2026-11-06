package ags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	sdkerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var (
	_ resource.ResourceWithConfigure   = &AgsSandboxToolResource{}
	_ resource.ResourceWithImportState = &AgsSandboxToolResource{}
)

func NewAgsSandboxToolResource() resource.Resource {
	return &AgsSandboxToolResource{}
}

type AgsSandboxToolResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// ---------- leaf models ----------

type agsVpcConfigModel struct {
	SubnetIds        types.Set `tfsdk:"subnet_ids"`
	SecurityGroupIds types.Set `tfsdk:"security_group_ids"`
}

type agsNetworkConfigModel struct {
	NetworkMode types.String       `tfsdk:"network_mode"`
	VpcConfig   *agsVpcConfigModel `tfsdk:"vpc_config"`
}

type agsTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type agsCosStorageModel struct {
	Endpoint   types.String `tfsdk:"endpoint"`
	BucketName types.String `tfsdk:"bucket_name"`
	BucketPath types.String `tfsdk:"bucket_path"`
}

type agsImageStorageModel struct {
	Reference         types.String `tfsdk:"reference"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	SubPath           types.String `tfsdk:"sub_path"`
	Digest            types.String `tfsdk:"digest"`
}

type agsCfsStorageModel struct {
	FileSystemId types.String `tfsdk:"file_system_id"`
	Path         types.String `tfsdk:"path"`
}

type agsAgentBucketStorageModel struct {
	LibraryId    types.String `tfsdk:"library_id"`
	SpaceId      types.String `tfsdk:"space_id"`
	AccessDomain types.String `tfsdk:"access_domain"`
}

type agsStorageSourceModel struct {
	Cos         *agsCosStorageModel         `tfsdk:"cos"`
	Image       *agsImageStorageModel       `tfsdk:"image"`
	Cfs         *agsCfsStorageModel         `tfsdk:"cfs"`
	AgentBucket *agsAgentBucketStorageModel `tfsdk:"agent_bucket"`
}

type agsVolumeModel struct {
	VolumeName types.String `tfsdk:"volume_name"`
	VolumeId   types.String `tfsdk:"volume_id"`
}

type agsVolumeTemplateModel struct {
	VolumeTemplateName types.String `tfsdk:"volume_template_name"`
	ReuseKey           types.String `tfsdk:"reuse_key"`
}

type agsStorageMountModel struct {
	Name           types.String            `tfsdk:"name"`
	StorageSource  *agsStorageSourceModel  `tfsdk:"storage_source"`
	MountPath      types.String            `tfsdk:"mount_path"`
	ReadOnly       types.Bool              `tfsdk:"read_only"`
	Volume         *agsVolumeModel         `tfsdk:"volume"`
	VolumeTemplate *agsVolumeTemplateModel `tfsdk:"volume_template"`
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

type agsResourceConfigModel struct {
	CPU     types.String `tfsdk:"cpu"`
	Memory  types.String `tfsdk:"memory"`
	Storage types.String `tfsdk:"storage"`
}

type agsHttpGetModel struct {
	Path   types.String `tfsdk:"path"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`
}

type agsProbeModel struct {
	HttpGet          *agsHttpGetModel `tfsdk:"http_get"`
	ReadyTimeoutMs   types.Int64      `tfsdk:"ready_timeout_ms"`
	ProbeTimeoutMs   types.Int64      `tfsdk:"probe_timeout_ms"`
	ProbePeriodMs    types.Int64      `tfsdk:"probe_period_ms"`
	SuccessThreshold types.Int64      `tfsdk:"success_threshold"`
	FailureThreshold types.Int64      `tfsdk:"failure_threshold"`
}

type agsDNSConfigModel struct {
	Servers  types.Set `tfsdk:"servers"`
	Searches types.Set `tfsdk:"searches"`
	Options  types.Set `tfsdk:"options"`
}

type agsCapabilitiesModel struct {
	Add  types.Set `tfsdk:"add"`
	Drop types.Set `tfsdk:"drop"`
}

type agsSecurityContextModel struct {
	Privileged   types.Bool            `tfsdk:"privileged"`
	Capabilities *agsCapabilitiesModel `tfsdk:"capabilities"`
}

type agsCustomConfigModel struct {
	Image             types.String             `tfsdk:"image"`
	ImageRegistryType types.String             `tfsdk:"image_registry_type"`
	Command           types.Set                `tfsdk:"command"`
	Args              types.Set                `tfsdk:"args"`
	Env               types.List               `tfsdk:"env"`
	Ports             types.List               `tfsdk:"ports"`
	Resources         *agsResourceConfigModel  `tfsdk:"resources"`
	Probe             *agsProbeModel           `tfsdk:"probe"`
	DNSConfig         *agsDNSConfigModel       `tfsdk:"dns_config"`
	SecurityContext   *agsSecurityContextModel `tfsdk:"security_context"`
}

type agsWAAConfigModel struct {
	ImageId types.String `tfsdk:"image_id"`
}

type agsOSWorldConfigModel struct {
	Version types.String `tfsdk:"version"`
}

type agsComputerConfigModel struct {
	WAAConfiguration     *agsWAAConfigModel     `tfsdk:"waa_configuration"`
	OSWorldConfiguration *agsOSWorldConfigModel `tfsdk:"os_world_configuration"`
}

type agsCLSConfigModel struct {
	TopicId types.String `tfsdk:"topic_id"`
}

type agsLogSourcesModel struct {
	Files types.Set `tfsdk:"files"`
}

type agsLogConfigModel struct {
	CLSConfig  *agsCLSConfigModel  `tfsdk:"cls_config"`
	LogSources *agsLogSourcesModel `tfsdk:"log_sources"`
}

// ---------- AgsSandboxToolResourceModel ----------

type AgsSandboxToolResourceModel struct {
	Id                    types.String            `tfsdk:"id"`
	ToolName              types.String            `tfsdk:"tool_name"`
	ToolType              types.String            `tfsdk:"tool_type"`
	NetworkConfiguration  *agsNetworkConfigModel  `tfsdk:"network_configuration"`
	Description           types.String            `tfsdk:"description"`
	DefaultTimeout        types.String            `tfsdk:"default_timeout"`
	Tags                  types.List              `tfsdk:"tags"`
	ClientToken           types.String            `tfsdk:"client_token"`
	RoleArn               types.String            `tfsdk:"role_arn"`
	StorageMounts         types.List              `tfsdk:"storage_mounts"`
	CustomConfiguration   *agsCustomConfigModel   `tfsdk:"custom_configuration"`
	ComputerConfiguration *agsComputerConfigModel `tfsdk:"computer_configuration"`
	LogConfiguration      *agsLogConfigModel      `tfsdk:"log_configuration"`
	Persistent            types.Bool              `tfsdk:"persistent"`
	// computed outputs
	ToolId     types.String `tfsdk:"tool_id"`
	Status     types.String `tfsdk:"status"`
	CreateTime types.String `tfsdk:"create_time"`
	UpdateTime types.String `tfsdk:"update_time"`
}

// ---------- schema ----------

func (r *AgsSandboxToolResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_tool"
}

func (r *AgsSandboxToolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create and manage an AGS sandbox tool.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier (same as tool_id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_name": schema.StringAttribute{
				Required: true,
				Description: "Sandbox tool name, 1-50 characters, supports English letters, digits, " +
					"underscores and hyphens. Must be unique under the same AppId.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tool_type": schema.StringAttribute{
				Required: true,
				Description: "Sandbox tool type. Supported values: browser, code-interpreter, mobile, " +
					"osworld, custom, swebench, aio, android-world, waa.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_configuration": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"network_mode": schema.StringAttribute{
						Required:    true,
						Description: "Network mode. Supported values: PUBLIC, VPC, SANDBOX.",
					},
					"vpc_config": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "VPC network configuration.",
						Attributes: map[string]schema.Attribute{
							"subnet_ids": schema.SetAttribute{
								Optional:    true,
								ElementType: types.StringType,
								Description: "List of VPC subnet IDs.",
							},
							"security_group_ids": schema.SetAttribute{
								Optional:    true,
								ElementType: types.StringType,
								Description: "List of security group IDs.",
							},
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Sandbox tool description, maximum 200 characters.",
			},
			"default_timeout": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Default timeout, supports formats like 5m, 300s, 1h. Maximum 24 hours. " +
					"If not specified, the system default (5 minutes) is used.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Tags for the sandbox tool.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Optional:    true,
							Description: "Tag key.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"client_token": schema.StringAttribute{
				Optional:    true,
				Description: "Idempotency token, maximum 64 characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_arn": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Role ARN.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_mounts": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Storage mount configuration for sandbox tool instances.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "Storage mount configuration name.",
						},
						"storage_source": schema.SingleNestedAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Storage source configuration.",
							Attributes: map[string]schema.Attribute{
								"cos": schema.SingleNestedAttribute{
									Optional:    true,
									Computed:    true,
									Description: "COS storage configuration.",
									Attributes: map[string]schema.Attribute{
										"endpoint": schema.StringAttribute{
											Optional:    true,
											Description: "COS access domain.",
										},
										"bucket_name": schema.StringAttribute{
											Optional:    true,
											Description: "COS bucket name.",
										},
										"bucket_path": schema.StringAttribute{
											Optional:    true,
											Description: "COS bucket path, must be an absolute path starting with /.",
										},
									},
								},
								"image": schema.SingleNestedAttribute{
									Optional:    true,
									Computed:    true,
									Description: "Image volume configuration.",
									Attributes: map[string]schema.Attribute{
										"reference": schema.StringAttribute{
											Optional:    true,
											Description: "Image address.",
										},
										"image_registry_type": schema.StringAttribute{
											Optional:    true,
											Description: "Image registry type: enterprise, personal.",
										},
										"sub_path": schema.StringAttribute{
											Optional:    true,
											Description: "Path inside the image.",
										},
										"digest": schema.StringAttribute{
											Optional:    true,
											Computed:    true,
											Description: "Image digest, not required in requests.",
										},
									},
								},
								"cfs": schema.SingleNestedAttribute{
									Optional:    true,
									Computed:    true,
									Description: "CFS storage configuration.",
									Attributes: map[string]schema.Attribute{
										"file_system_id": schema.StringAttribute{
											Optional:    true,
											Description: "CFS resource ID.",
										},
										"path": schema.StringAttribute{
											Optional:    true,
											Description: "CFS mount path.",
										},
									},
								},
								"agent_bucket": schema.SingleNestedAttribute{
									Optional:    true,
									Computed:    true,
									Description: "AgentBucket storage configuration.",
									Attributes: map[string]schema.Attribute{
										"library_id": schema.StringAttribute{
											Optional:    true,
											Description: "AgentBucket LibraryID.",
										},
										"space_id": schema.StringAttribute{
											Optional:    true,
											Description: "AgentBucket spaceId.",
										},
										"access_domain": schema.StringAttribute{
											Optional:    true,
											Description: "AgentBucket AccessDomain.",
										},
									},
								},
							},
						},
						"mount_path": schema.StringAttribute{
							Optional:    true,
							Description: "Local mount path in the sandbox instance.",
						},
						"read_only": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Read-only permission configuration, defaults to false.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"volume": schema.SingleNestedAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Volume reference.",
							Attributes: map[string]schema.Attribute{
								"volume_name": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "Volume name; resolved to a stable VolumeId under the current AppId.",
								},
								"volume_id": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "Stable volume ID, format vol-xxxxxxxx.",
								},
							},
						},
						"volume_template": schema.SingleNestedAttribute{
							Optional:    true,
							Computed:    true,
							Description: "VolumeTemplate reference and derivation rules.",
							Attributes: map[string]schema.Attribute{
								"volume_template_name": schema.StringAttribute{
									Required:    true,
									Description: "VolumeTemplate name.",
								},
								"reuse_key": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "Reuse key, can be a literal or ${metadata.xxx}.",
								},
							},
						},
					},
				},
			},
			"custom_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Custom configuration for sandbox tool.",
				Attributes: map[string]schema.Attribute{
					"image": schema.StringAttribute{
						Optional:    true,
						Description: "Image address.",
					},
					"image_registry_type": schema.StringAttribute{
						Optional:    true,
						Description: "Image registry type: enterprise, personal, custom.",
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
						Description: "Environment variables.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Optional:    true,
									Description: "Environment variable name.",
								},
								"value": schema.StringAttribute{
									Optional:    true,
									Description: "Environment variable value.",
								},
							},
						},
					},
					"ports": schema.ListNestedAttribute{
						Optional:    true,
						Description: "Port configuration.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
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
							},
						},
					},
					"resources": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Resource configuration.",
						Attributes: map[string]schema.Attribute{
							"cpu": schema.StringAttribute{
								Optional:    true,
								Description: "CPU resource amount.",
							},
							"memory": schema.StringAttribute{
								Optional:    true,
								Description: "Memory resource amount.",
							},
							"storage": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Custom disk size. Supported values: 1Gi, 5Gi, 10Gi, 20Gi.",
							},
						},
					},
					"probe": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Probe configuration.",
						Attributes: map[string]schema.Attribute{
							"http_get": schema.SingleNestedAttribute{
								Optional:    true,
								Computed:    true,
								Description: "HTTP GET probe configuration.",
								Attributes: map[string]schema.Attribute{
									"path": schema.StringAttribute{
										Optional:    true,
										Description: "Path.",
									},
									"port": schema.Int64Attribute{
										Optional:    true,
										Description: "Port.",
									},
									"scheme": schema.StringAttribute{
										Optional:    true,
										Description: "Scheme.",
									},
								},
							},
							"ready_timeout_ms": schema.Int64Attribute{
								Optional:    true,
								Description: "Health check readiness timeout.",
							},
							"probe_timeout_ms": schema.Int64Attribute{
								Optional:    true,
								Description: "Single probe timeout.",
							},
							"probe_period_ms": schema.Int64Attribute{
								Optional:    true,
								Description: "Probe interval.",
							},
							"success_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Health check success threshold.",
							},
							"failure_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Health check failure threshold.",
							},
						},
					},
					"dns_config": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Sandbox DNS configuration.",
						Attributes: map[string]schema.Attribute{
							"servers": schema.SetAttribute{
								Required:    true,
								ElementType: types.StringType,
								Description: "DNS server addresses (valid IP addresses). Default: 10.0.0.1.",
							},
							"searches": schema.SetAttribute{
								Optional:    true,
								ElementType: types.StringType,
								Description: "Search domains (corresponds to resolv.conf search directive).",
							},
							"options": schema.SetAttribute{
								Optional:    true,
								ElementType: types.StringType,
								Description: "Configuration items (corresponds to resolv.conf options).",
							},
						},
					},
					"security_context": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Sandbox security context.",
						Attributes: map[string]schema.Attribute{
							"privileged": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Whether to enable privileged mode. Default: true.",
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"capabilities": schema.SingleNestedAttribute{
								Optional:    true,
								Computed:    true,
								Description: "POSIX capabilities to add or remove from the container.",
								Attributes: map[string]schema.Attribute{
									"add": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Capabilities to add. Example: [\"SYS_ADMIN\"].",
									},
									"drop": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Capabilities to drop. Example: [\"ALL\"].",
									},
								},
							},
						},
					},
				},
			},
			"computer_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Desktop/computer environment sandbox configuration.",
				Attributes: map[string]schema.Attribute{
					"waa_configuration": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "WAA sandbox tool configuration.",
						Attributes: map[string]schema.Attribute{
							"image_id": schema.StringAttribute{
								Optional:    true,
								Description: "Custom WAA image ID.",
							},
						},
					},
					"os_world_configuration": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Built-in OSWorld configuration.",
						Attributes: map[string]schema.Attribute{
							"version": schema.StringAttribute{
								Optional:    true,
								Description: "Built-in OSWorld version. Supported: osworld1, osworld2. Default: osworld1.",
							},
						},
					},
				},
			},
			"log_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Log push configuration for sandbox tool.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"cls_config": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "CLS log push configuration.",
						Attributes: map[string]schema.Attribute{
							"topic_id": schema.StringAttribute{
								Optional:    true,
								Description: "CLS log topic ID for sandbox tool log push.",
							},
						},
					},
					"log_sources": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Log source configuration.",
						Attributes: map[string]schema.Attribute{
							"files": schema.SetAttribute{
								Optional:    true,
								ElementType: types.StringType,
								Description: "Log file paths to collect, must be under /logs/ directory, no subdirectories, max 10 files.",
							},
						},
					},
				},
			},
			"persistent": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Persistent sandbox flag.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			// computed outputs
			"tool_id": schema.StringAttribute{
				Computed:    true,
				Description: "Sandbox tool ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Sandbox tool status: CREATING, ACTIVE, DELETING, FAILED.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time in ISO8601 format.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Update time in ISO8601 format.",
			},
		},
	}
}

// ---------- CRUD ----------

func (r *AgsSandboxToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsSandboxToolResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := ags.NewCreateSandboxToolRequest()
	request.ToolName = plan.ToolName.ValueStringPointer()
	request.ToolType = plan.ToolType.ValueStringPointer()

	request.NetworkConfiguration = buildNetworkConfigSDK(plan.NetworkConfiguration)

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}
	if !plan.DefaultTimeout.IsNull() && !plan.DefaultTimeout.IsUnknown() {
		request.DefaultTimeout = plan.DefaultTimeout.ValueStringPointer()
	}
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		request.Tags = buildTagsSDK(ctx, plan.Tags, &resp.Diagnostics)
	}
	if !plan.ClientToken.IsNull() && !plan.ClientToken.IsUnknown() {
		request.ClientToken = plan.ClientToken.ValueStringPointer()
	}
	if !plan.RoleArn.IsNull() && !plan.RoleArn.IsUnknown() {
		request.RoleArn = plan.RoleArn.ValueStringPointer()
	}
	if !plan.StorageMounts.IsNull() && !plan.StorageMounts.IsUnknown() {
		request.StorageMounts = buildStorageMountsSDK(ctx, plan.StorageMounts, &resp.Diagnostics)
	}
	if plan.CustomConfiguration != nil {
		request.CustomConfiguration = buildCustomConfigSDK(ctx, plan.CustomConfiguration, &resp.Diagnostics)
	}
	if plan.ComputerConfiguration != nil {
		request.ComputerConfiguration = buildComputerConfigSDK(plan.ComputerConfiguration)
	}
	if plan.LogConfiguration != nil {
		request.LogConfiguration = buildLogConfigSDK(plan.LogConfiguration)
	}
	if !plan.Persistent.IsNull() && !plan.Persistent.IsUnknown() {
		request.Persistent = plan.Persistent.ValueBoolPointer()
	}
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.Client().UseAgsClient().CreateSandboxToolWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AGS sandbox tool",
			fmt.Sprintf("CreateSandboxTool failed: %v", err),
		)
		return
	}

	if response == nil || response.Response == nil || response.Response.ToolId == nil || *response.Response.ToolId == "" {
		resp.Diagnostics.AddError(
			"Error creating AGS sandbox tool",
			"CreateSandboxTool returned an empty response or an empty tool id.",
		)
		return
	}

	toolId := *response.Response.ToolId
	plan.Id = types.StringValue(toolId)
	plan.ToolId = types.StringValue(toolId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tool, err := r.read(ctx, toolId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading AGS sandbox tool after create", err.Error())
		return
	}
	if tool == nil {
		resp.Diagnostics.AddError(
			"AGS sandbox tool not readable after create",
			fmt.Sprintf("tool %s was created but cannot be read back", toolId),
		)
		return
	}
	refreshModel(ctx, &plan, tool, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AgsSandboxToolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsSandboxToolResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tool, err := r.read(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading AGS sandbox tool", err.Error())
		return
	}

	if tool == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("AGS sandbox tool %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	refreshModel(ctx, &state, tool, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

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

	request := ags.NewUpdateSandboxToolRequest()
	request.ToolId = state.Id.ValueStringPointer()

	if !plan.Description.IsUnknown() {
		if plan.Description.IsNull() {
			request.Description = helper.String("")
		} else {
			request.Description = plan.Description.ValueStringPointer()
		}
	}

	if plan.NetworkConfiguration != nil {
		request.NetworkConfiguration = buildNetworkConfigSDK(plan.NetworkConfiguration)
	}

	if !plan.Tags.IsUnknown() {
		if plan.Tags.IsNull() {
			request.Tags = []*ags.Tag{}
		} else {
			request.Tags = buildTagsSDK(ctx, plan.Tags, &resp.Diagnostics)
		}
	}

	if plan.CustomConfiguration != nil {
		request.CustomConfiguration = buildCustomConfigSDK(ctx, plan.CustomConfiguration, &resp.Diagnostics)
	}

	if plan.ComputerConfiguration != nil {
		request.ComputerConfiguration = buildComputerConfigSDK(plan.ComputerConfiguration)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Client().UseAgsClient().UpdateSandboxToolWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating AGS sandbox tool",
			fmt.Sprintf("UpdateSandboxTool failed: %v", err),
		)
		return
	}

	tool, err := r.read(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading AGS sandbox tool after update", err.Error())
		return
	}
	if tool == nil {
		resp.Diagnostics.AddError(
			"AGS sandbox tool not readable after update",
			fmt.Sprintf("tool %s was updated but cannot be read back", state.Id.ValueString()),
		)
		return
	}
	plan.Id = state.Id
	refreshModel(ctx, &plan, tool, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AgsSandboxToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsSandboxToolResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := ags.NewDeleteSandboxToolRequest()
	request.ToolId = state.Id.ValueStringPointer()

	_, err := r.Client().UseAgsClient().DeleteSandboxToolWithContext(ctx, request)
	if err != nil {
		if !isAgsResourceNotFound(err) {
			resp.Diagnostics.AddError(
				"Error deleting AGS sandbox tool",
				fmt.Sprintf("DeleteSandboxTool failed: %v", err),
			)
			return
		}
	}
}

// ---------- helper: read ----------

func (r *AgsSandboxToolResource) read(ctx context.Context, toolId string) (*ags.SandboxTool, error) {
	request := ags.NewDescribeSandboxToolListRequest()
	request.ToolIds = []*string{&toolId}
	request.Limit = helper.Int64(1)

	response, err := r.Client().UseAgsClient().DescribeSandboxToolListWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("DescribeSandboxToolList returned an empty response for %s", toolId)
	}
	if len(response.Response.SandboxToolSet) == 0 {
		return nil, nil
	}
	return response.Response.SandboxToolSet[0], nil
}

// ---------- helper: refresh model from SDK response ----------

func refreshModel(ctx context.Context, m *AgsSandboxToolResourceModel, tool *ags.SandboxTool, diags *diag.Diagnostics) {
	if tool == nil {
		return
	}
	m.ToolId = types.StringPointerValue(tool.ToolId)
	m.ToolName = types.StringPointerValue(tool.ToolName)
	m.ToolType = types.StringPointerValue(tool.ToolType)
	m.Status = types.StringPointerValue(tool.Status)
	m.Description = types.StringPointerValue(tool.Description)
	m.Persistent = types.BoolPointerValue(tool.Persistent)
	m.RoleArn = types.StringPointerValue(tool.RoleArn)
	m.CreateTime = types.StringPointerValue(tool.CreateTime)
	m.UpdateTime = types.StringPointerValue(tool.UpdateTime)

	if tool.DefaultTimeoutSeconds != nil && m.DefaultTimeout.IsNull() {
		m.DefaultTimeout = types.StringValue(fmt.Sprintf("%ds", *tool.DefaultTimeoutSeconds))
	}

	if tool.NetworkConfiguration != nil {
		m.NetworkConfiguration = buildNetworkConfigState(tool.NetworkConfiguration)
	}
	if tool.Tags != nil {
		tagList := make([]agsTagModel, 0, len(tool.Tags))
		for _, t := range tool.Tags {
			tagList = append(tagList, agsTagModel{
				Key:   types.StringPointerValue(t.Key),
				Value: types.StringPointerValue(t.Value),
			})
		}
		list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"key": types.StringType, "value": types.StringType,
		}}, tagList)
		diags.Append(listDiags...)
		if !listDiags.HasError() {
			m.Tags = list
		}
	}
	if tool.StorageMounts != nil {
		m.StorageMounts = buildStorageMountsState(ctx, tool.StorageMounts, diags)
	}
	if tool.CustomConfiguration != nil {
		m.CustomConfiguration = buildCustomConfigState(tool.CustomConfiguration)
	}
	if tool.LogConfiguration != nil {
		m.LogConfiguration = buildLogConfigState(tool.LogConfiguration)
	}
	if tool.ComputerConfiguration != nil {
		m.ComputerConfiguration = buildComputerConfigObjState(tool.ComputerConfiguration)
	}
}

// ---------- build SDK request from model ----------

func buildNetworkConfigSDK(nc *agsNetworkConfigModel) *ags.NetworkConfiguration {
	if nc == nil {
		return nil
	}
	out := &ags.NetworkConfiguration{
		NetworkMode: nc.NetworkMode.ValueStringPointer(),
	}
	if nc.VpcConfig != nil {
		vpc := nc.VpcConfig
		vpcOut := &ags.VPCConfig{}
		if !vpc.SubnetIds.IsNull() && !vpc.SubnetIds.IsUnknown() {
			elems := make([]types.String, 0)
			_ = vpc.SubnetIds.ElementsAs(context.TODO(), &elems, false)
			for _, s := range elems {
				vpcOut.SubnetIds = append(vpcOut.SubnetIds, s.ValueStringPointer())
			}
		}
		if !vpc.SecurityGroupIds.IsNull() && !vpc.SecurityGroupIds.IsUnknown() {
			elems := make([]types.String, 0)
			_ = vpc.SecurityGroupIds.ElementsAs(context.TODO(), &elems, false)
			for _, s := range elems {
				vpcOut.SecurityGroupIds = append(vpcOut.SecurityGroupIds, s.ValueStringPointer())
			}
		}
		out.VpcConfig = vpcOut
	}
	return out
}

func buildTagsSDK(ctx context.Context, list types.List, diags *diag.Diagnostics) []*ags.Tag {
	var tags []agsTagModel
	diags.Append(list.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil
	}
	out := make([]*ags.Tag, 0, len(tags))
	for _, t := range tags {
		out = append(out, &ags.Tag{
			Key: t.Key.ValueStringPointer(), Value: t.Value.ValueStringPointer(),
		})
	}
	return out
}

func buildStorageSourceSDK(src *agsStorageSourceModel) *ags.StorageSource {
	if src == nil {
		return nil
	}
	out := &ags.StorageSource{}
	if src.Cos != nil {
		out.Cos = &ags.CosStorageSource{
			Endpoint: src.Cos.Endpoint.ValueStringPointer(), BucketName: src.Cos.BucketName.ValueStringPointer(),
			BucketPath: src.Cos.BucketPath.ValueStringPointer(),
		}
	}
	if src.Image != nil {
		out.Image = &ags.ImageStorageSource{
			Reference:         src.Image.Reference.ValueStringPointer(),
			ImageRegistryType: src.Image.ImageRegistryType.ValueStringPointer(),
			SubPath:           src.Image.SubPath.ValueStringPointer(),
		}
	}
	if src.Cfs != nil {
		out.Cfs = &ags.CfsStorageSource{
			FileSystemId: src.Cfs.FileSystemId.ValueStringPointer(), Path: src.Cfs.Path.ValueStringPointer(),
		}
	}
	if src.AgentBucket != nil {
		out.AgentBucket = &ags.AgentBucketStorageSource{
			LibraryId:    src.AgentBucket.LibraryId.ValueStringPointer(),
			SpaceId:      src.AgentBucket.SpaceId.ValueStringPointer(),
			AccessDomain: src.AgentBucket.AccessDomain.ValueStringPointer(),
		}
	}
	return out
}

func buildStorageMountsSDK(ctx context.Context, list types.List, diags *diag.Diagnostics) []*ags.StorageMount {
	var mounts []agsStorageMountModel
	diags.Append(list.ElementsAs(ctx, &mounts, false)...)
	if diags.HasError() {
		return nil
	}
	out := make([]*ags.StorageMount, 0, len(mounts))
	for _, m := range mounts {
		sm := &ags.StorageMount{
			Name: m.Name.ValueStringPointer(), MountPath: m.MountPath.ValueStringPointer(),
		}
		if !m.ReadOnly.IsNull() && !m.ReadOnly.IsUnknown() {
			sm.ReadOnly = m.ReadOnly.ValueBoolPointer()
		}
		if m.StorageSource != nil {
			sm.StorageSource = buildStorageSourceSDK(m.StorageSource)
		}
		if m.Volume != nil {
			sm.Volume = &ags.VolumeRef{
				VolumeName: m.Volume.VolumeName.ValueStringPointer(),
				VolumeId:   m.Volume.VolumeId.ValueStringPointer(),
			}
		}
		if m.VolumeTemplate != nil {
			sm.VolumeTemplate = &ags.VolumeTemplateRef{
				VolumeTemplateName: m.VolumeTemplate.VolumeTemplateName.ValueStringPointer(),
				ReuseKey:           m.VolumeTemplate.ReuseKey.ValueStringPointer(),
			}
		}
		out = append(out, sm)
	}
	return out
}

func buildCustomConfigSDK(ctx context.Context, cfg *agsCustomConfigModel, diags *diag.Diagnostics) *ags.CustomConfiguration {
	if cfg == nil {
		return nil
	}
	out := &ags.CustomConfiguration{
		Image: cfg.Image.ValueStringPointer(), ImageRegistryType: cfg.ImageRegistryType.ValueStringPointer(),
	}
	if !cfg.Command.IsNull() && !cfg.Command.IsUnknown() {
		elems := make([]types.String, 0)
		_ = cfg.Command.ElementsAs(context.TODO(), &elems, false)
		for _, c := range elems {
			out.Command = append(out.Command, c.ValueStringPointer())
		}
	}
	if !cfg.Args.IsNull() && !cfg.Args.IsUnknown() {
		elems := make([]types.String, 0)
		_ = cfg.Args.ElementsAs(context.TODO(), &elems, false)
		for _, a := range elems {
			out.Args = append(out.Args, a.ValueStringPointer())
		}
	}
	if !cfg.Env.IsNull() && !cfg.Env.IsUnknown() {
		var envs []agsEnvVarModel
		diags.Append(cfg.Env.ElementsAs(ctx, &envs, false)...)
		for _, e := range envs {
			out.Env = append(out.Env, &ags.EnvVar{
				Name: e.Name.ValueStringPointer(), Value: e.Value.ValueStringPointer(),
			})
		}
	}
	if !cfg.Ports.IsNull() && !cfg.Ports.IsUnknown() {
		var ports []agsPortModel
		diags.Append(cfg.Ports.ElementsAs(ctx, &ports, false)...)
		for _, p := range ports {
			pc := &ags.PortConfiguration{
				Name: p.Name.ValueStringPointer(), Protocol: p.Protocol.ValueStringPointer(),
			}
			if !p.Port.IsNull() && !p.Port.IsUnknown() {
				pc.Port = p.Port.ValueInt64Pointer()
			}
			out.Ports = append(out.Ports, pc)
		}
	}
	if cfg.Resources != nil {
		out.Resources = &ags.ResourceConfiguration{
			CPU: cfg.Resources.CPU.ValueStringPointer(), Memory: cfg.Resources.Memory.ValueStringPointer(),
			Storage: cfg.Resources.Storage.ValueStringPointer(),
		}
	}
	if cfg.Probe != nil {
		p := cfg.Probe
		out.Probe = &ags.ProbeConfiguration{
			ReadyTimeoutMs: p.ReadyTimeoutMs.ValueInt64Pointer(), ProbeTimeoutMs: p.ProbeTimeoutMs.ValueInt64Pointer(),
			ProbePeriodMs: p.ProbePeriodMs.ValueInt64Pointer(), SuccessThreshold: p.SuccessThreshold.ValueInt64Pointer(),
			FailureThreshold: p.FailureThreshold.ValueInt64Pointer(),
		}
		if p.HttpGet != nil {
			out.Probe.HttpGet = &ags.HttpGetAction{
				Path: p.HttpGet.Path.ValueStringPointer(), Port: p.HttpGet.Port.ValueInt64Pointer(),
				Scheme: p.HttpGet.Scheme.ValueStringPointer(),
			}
		}
	}
	if cfg.DNSConfig != nil {
		dns := cfg.DNSConfig
		out.DNSConfig = &ags.DNSConfig{}
		if !dns.Servers.IsNull() && !dns.Servers.IsUnknown() {
			elems := make([]types.String, 0)
			_ = dns.Servers.ElementsAs(context.TODO(), &elems, false)
			for _, s := range elems {
				out.DNSConfig.Servers = append(out.DNSConfig.Servers, s.ValueStringPointer())
			}
		}
		if !dns.Searches.IsNull() && !dns.Searches.IsUnknown() {
			elems := make([]types.String, 0)
			_ = dns.Searches.ElementsAs(context.TODO(), &elems, false)
			for _, s := range elems {
				out.DNSConfig.Searches = append(out.DNSConfig.Searches, s.ValueStringPointer())
			}
		}
		if !dns.Options.IsNull() && !dns.Options.IsUnknown() {
			elems := make([]types.String, 0)
			_ = dns.Options.ElementsAs(context.TODO(), &elems, false)
			for _, o := range elems {
				out.DNSConfig.Options = append(out.DNSConfig.Options, o.ValueStringPointer())
			}
		}
	}
	if cfg.SecurityContext != nil {
		sc := cfg.SecurityContext
		out.SecurityContext = &ags.SecurityContext{Privileged: sc.Privileged.ValueBoolPointer()}
		if sc.Capabilities != nil {
			caps := sc.Capabilities
			out.SecurityContext.Capabilities = &ags.Capabilities{}
			if !caps.Add.IsNull() && !caps.Add.IsUnknown() {
				elems := make([]types.String, 0)
				_ = caps.Add.ElementsAs(context.TODO(), &elems, false)
				for _, a := range elems {
					out.SecurityContext.Capabilities.Add = append(out.SecurityContext.Capabilities.Add, a.ValueStringPointer())
				}
			}
			if !caps.Drop.IsNull() && !caps.Drop.IsUnknown() {
				elems := make([]types.String, 0)
				_ = caps.Drop.ElementsAs(context.TODO(), &elems, false)
				for _, d := range elems {
					out.SecurityContext.Capabilities.Drop = append(out.SecurityContext.Capabilities.Drop, d.ValueStringPointer())
				}
			}
		}
	}
	return out
}

func buildComputerConfigSDK(cfg *agsComputerConfigModel) *ags.ComputerConfiguration {
	if cfg == nil {
		return nil
	}
	out := &ags.ComputerConfiguration{}
	if cfg.WAAConfiguration != nil {
		out.WAAConfiguration = &ags.WAAConfiguration{ImageId: cfg.WAAConfiguration.ImageId.ValueStringPointer()}
	}
	if cfg.OSWorldConfiguration != nil {
		out.OSWorldConfiguration = &ags.OSWorldConfiguration{Version: cfg.OSWorldConfiguration.Version.ValueStringPointer()}
	}
	return out
}

func buildLogConfigSDK(cfg *agsLogConfigModel) *ags.LogConfiguration {
	if cfg == nil {
		return nil
	}
	out := &ags.LogConfiguration{}
	if cfg.CLSConfig != nil {
		out.CLSConfig = &ags.CLSConfig{TopicId: cfg.CLSConfig.TopicId.ValueStringPointer()}
	}
	if cfg.LogSources != nil {
		out.LogSources = &ags.LogSources{}
		if !cfg.LogSources.Files.IsNull() && !cfg.LogSources.Files.IsUnknown() {
			elems := make([]types.String, 0)
			_ = cfg.LogSources.Files.ElementsAs(context.TODO(), &elems, false)
			for _, f := range elems {
				out.LogSources.Files = append(out.LogSources.Files, f.ValueStringPointer())
			}
		}
	}
	return out
}

// ---------- build state model from SDK response ----------

func buildNetworkConfigState(nc *ags.NetworkConfiguration) *agsNetworkConfigModel {
	if nc == nil {
		return nil
	}
	out := &agsNetworkConfigModel{NetworkMode: types.StringPointerValue(nc.NetworkMode)}
	if nc.VpcConfig != nil {
		subnetIds, _ := types.SetValueFrom(context.TODO(), types.StringType, nc.VpcConfig.SubnetIds)
		sgIds, _ := types.SetValueFrom(context.TODO(), types.StringType, nc.VpcConfig.SecurityGroupIds)
		out.VpcConfig = &agsVpcConfigModel{SubnetIds: subnetIds, SecurityGroupIds: sgIds}
	}
	return out
}

func buildStorageMountsState(ctx context.Context, mounts []*ags.StorageMount, diags *diag.Diagnostics) types.List {
	items := make([]agsStorageMountModel, 0, len(mounts))
	for _, m := range mounts {
		item := agsStorageMountModel{
			Name: types.StringPointerValue(m.Name), MountPath: types.StringPointerValue(m.MountPath),
			ReadOnly: types.BoolPointerValue(m.ReadOnly),
		}
		if m.StorageSource != nil {
			item.StorageSource = buildStorageSourceState(m.StorageSource)
		}
		if m.Volume != nil {
			item.Volume = &agsVolumeModel{
				VolumeName: types.StringPointerValue(m.Volume.VolumeName),
				VolumeId:   types.StringPointerValue(m.Volume.VolumeId),
			}
		}
		if m.VolumeTemplate != nil {
			item.VolumeTemplate = &agsVolumeTemplateModel{
				VolumeTemplateName: types.StringPointerValue(m.VolumeTemplate.VolumeTemplateName),
				ReuseKey:           types.StringPointerValue(m.VolumeTemplate.ReuseKey),
			}
		}
		items = append(items, item)
	}
	smType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name": types.StringType, "storage_source": types.ObjectType{AttrTypes: storageSourceAttrTypesMap()},
		"mount_path": types.StringType, "read_only": types.BoolType,
		"volume":          types.ObjectType{AttrTypes: volumeAttrTypesMap()},
		"volume_template": types.ObjectType{AttrTypes: volumeTemplateAttrTypesMap()},
	}}
	list, listDiags := types.ListValueFrom(ctx, smType, items)
	diags.Append(listDiags...)
	return list
}

func storageSourceAttrTypesMap() map[string]attr.Type {
	return map[string]attr.Type{
		"cos":          types.ObjectType{AttrTypes: cosStorageAttrTypesMap()},
		"image":        types.ObjectType{AttrTypes: imageStorageAttrTypesMap()},
		"cfs":          types.ObjectType{AttrTypes: cfsStorageAttrTypesMap()},
		"agent_bucket": types.ObjectType{AttrTypes: agentBucketStorageAttrTypesMap()},
	}
}

func cosStorageAttrTypesMap() map[string]attr.Type {
	return map[string]attr.Type{"endpoint": types.StringType, "bucket_name": types.StringType, "bucket_path": types.StringType}
}

func imageStorageAttrTypesMap() map[string]attr.Type {
	return map[string]attr.Type{"reference": types.StringType, "image_registry_type": types.StringType, "sub_path": types.StringType, "digest": types.StringType}
}

func cfsStorageAttrTypesMap() map[string]attr.Type {
	return map[string]attr.Type{"file_system_id": types.StringType, "path": types.StringType}
}

func agentBucketStorageAttrTypesMap() map[string]attr.Type {
	return map[string]attr.Type{"library_id": types.StringType, "space_id": types.StringType, "access_domain": types.StringType}
}

func volumeAttrTypesMap() map[string]attr.Type {
	return map[string]attr.Type{"volume_name": types.StringType, "volume_id": types.StringType}
}

func volumeTemplateAttrTypesMap() map[string]attr.Type {
	return map[string]attr.Type{"volume_template_name": types.StringType, "reuse_key": types.StringType}
}

func buildStorageSourceState(src *ags.StorageSource) *agsStorageSourceModel {
	if src == nil {
		return nil
	}
	out := &agsStorageSourceModel{}
	if src.Cos != nil {
		out.Cos = &agsCosStorageModel{
			Endpoint: types.StringPointerValue(src.Cos.Endpoint), BucketName: types.StringPointerValue(src.Cos.BucketName),
			BucketPath: types.StringPointerValue(src.Cos.BucketPath),
		}
	}
	if src.Image != nil {
		out.Image = &agsImageStorageModel{
			Reference:         types.StringPointerValue(src.Image.Reference),
			ImageRegistryType: types.StringPointerValue(src.Image.ImageRegistryType),
			SubPath:           types.StringPointerValue(src.Image.SubPath), Digest: types.StringPointerValue(src.Image.Digest),
		}
	}
	if src.Cfs != nil {
		out.Cfs = &agsCfsStorageModel{
			FileSystemId: types.StringPointerValue(src.Cfs.FileSystemId), Path: types.StringPointerValue(src.Cfs.Path),
		}
	}
	if src.AgentBucket != nil {
		out.AgentBucket = &agsAgentBucketStorageModel{
			LibraryId:    types.StringPointerValue(src.AgentBucket.LibraryId),
			SpaceId:      types.StringPointerValue(src.AgentBucket.SpaceId),
			AccessDomain: types.StringPointerValue(src.AgentBucket.AccessDomain),
		}
	}
	return out
}

func buildCustomConfigState(cfg *ags.CustomConfigurationDetail) *agsCustomConfigModel {
	if cfg == nil {
		return nil
	}
	out := &agsCustomConfigModel{
		Image: types.StringPointerValue(cfg.Image), ImageRegistryType: types.StringPointerValue(cfg.ImageRegistryType),
	}
	cmdSet, _ := types.SetValueFrom(context.TODO(), types.StringType, cfg.Command)
	out.Command = cmdSet
	argsSet, _ := types.SetValueFrom(context.TODO(), types.StringType, cfg.Args)
	out.Args = argsSet
	if cfg.Resources != nil {
		out.Resources = &agsResourceConfigModel{
			CPU: types.StringPointerValue(cfg.Resources.CPU), Memory: types.StringPointerValue(cfg.Resources.Memory),
			Storage: types.StringPointerValue(cfg.Resources.Storage),
		}
	}
	if cfg.Probe != nil {
		probe := &agsProbeModel{
			ReadyTimeoutMs:   types.Int64PointerValue(cfg.Probe.ReadyTimeoutMs),
			ProbeTimeoutMs:   types.Int64PointerValue(cfg.Probe.ProbeTimeoutMs),
			ProbePeriodMs:    types.Int64PointerValue(cfg.Probe.ProbePeriodMs),
			SuccessThreshold: types.Int64PointerValue(cfg.Probe.SuccessThreshold),
			FailureThreshold: types.Int64PointerValue(cfg.Probe.FailureThreshold),
		}
		if cfg.Probe.HttpGet != nil {
			probe.HttpGet = &agsHttpGetModel{
				Path:   types.StringPointerValue(cfg.Probe.HttpGet.Path),
				Port:   types.Int64PointerValue(cfg.Probe.HttpGet.Port),
				Scheme: types.StringPointerValue(cfg.Probe.HttpGet.Scheme),
			}
		}
		out.Probe = probe
	}
	if cfg.DNSConfig != nil {
		serversSet, _ := types.SetValueFrom(context.TODO(), types.StringType, cfg.DNSConfig.Servers)
		searchesSet, _ := types.SetValueFrom(context.TODO(), types.StringType, cfg.DNSConfig.Searches)
		optionsSet, _ := types.SetValueFrom(context.TODO(), types.StringType, cfg.DNSConfig.Options)
		out.DNSConfig = &agsDNSConfigModel{Servers: serversSet, Searches: searchesSet, Options: optionsSet}
	}
	if cfg.SecurityContext != nil {
		sc := &agsSecurityContextModel{Privileged: types.BoolPointerValue(cfg.SecurityContext.Privileged)}
		if cfg.SecurityContext.Capabilities != nil {
			addSet, _ := types.SetValueFrom(context.TODO(), types.StringType, cfg.SecurityContext.Capabilities.Add)
			dropSet, _ := types.SetValueFrom(context.TODO(), types.StringType, cfg.SecurityContext.Capabilities.Drop)
			sc.Capabilities = &agsCapabilitiesModel{Add: addSet, Drop: dropSet}
		}
		out.SecurityContext = sc
	}
	return out
}

func buildComputerConfigObjState(cfg *ags.ComputerConfiguration) *agsComputerConfigModel {
	if cfg == nil {
		return nil
	}
	out := &agsComputerConfigModel{}
	if cfg.WAAConfiguration != nil {
		out.WAAConfiguration = &agsWAAConfigModel{ImageId: types.StringPointerValue(cfg.WAAConfiguration.ImageId)}
	}
	if cfg.OSWorldConfiguration != nil {
		out.OSWorldConfiguration = &agsOSWorldConfigModel{Version: types.StringPointerValue(cfg.OSWorldConfiguration.Version)}
	}
	return out
}

func buildLogConfigState(cfg *ags.LogConfiguration) *agsLogConfigModel {
	if cfg == nil {
		return nil
	}
	out := &agsLogConfigModel{}
	if cfg.CLSConfig != nil {
		out.CLSConfig = &agsCLSConfigModel{TopicId: types.StringPointerValue(cfg.CLSConfig.TopicId)}
	}
	if cfg.LogSources != nil {
		filesSet, _ := types.SetValueFrom(context.TODO(), types.StringType, cfg.LogSources.Files)
		out.LogSources = &agsLogSourcesModel{Files: filesSet}
	}
	return out
}

// ---------- error classification ----------

func isAgsResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	sdkErr, ok := err.(*sdkerr.TencentCloudSDKError)
	if !ok {
		return false
	}
	return sdkErr.Code == "ResourceNotFound"
}

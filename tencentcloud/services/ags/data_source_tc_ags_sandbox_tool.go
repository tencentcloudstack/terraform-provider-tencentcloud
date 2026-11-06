package ags

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

type SandboxToolDataSource struct {
	fw.DataSourceWithConfigure
	testService *Service
}

var _ datasource.DataSourceWithValidateConfig = &SandboxToolDataSource{}

func NewAgsSandboxToolDataSource() datasource.DataSource { return &SandboxToolDataSource{} }
func (d *SandboxToolDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_tool"
}
func (d *SandboxToolDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Reads exactly one visible AGS Tool by ID or exact name. Does not manage, mutate, import or start resources. Zero visible matches or ambiguous results are errors.", Attributes: map[string]schema.Attribute{
		"tool_id":       schema.StringAttribute{Optional: true, Description: "Exact Tool ID selector. Specify exactly one of tool_id or tool_name; the unused selector is not populated from the result."},
		"tool_name":     schema.StringAttribute{Optional: true, Description: "Exact Tool name selector, never a prefix or fuzzy match."},
		"id":            schema.StringAttribute{Computed: true, Description: "Actual Tool ID suitable for the runtime SDK."},
		"status":        schema.StringAttribute{Computed: true, Description: "Actual Tool status, including creating or failed Tools."},
		"status_reason": schema.StringAttribute{Computed: true, Description: "Service-reported status reason."},
		"image_digest":  schema.StringAttribute{Computed: true, Description: "Resolved custom image digest when available."},
		"configuration": schema.ObjectAttribute{Computed: true, AttributeTypes: configurationTypes(), Description: "Supported read-only actual configuration, including remote defaults and system tags. Environment values are stored in Terraform state; use non-sensitive variables only."},
	}}
}
func (d *SandboxToolDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var config types.Object
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() || config.IsNull() || config.IsUnknown() {
		return
	}
	f := config.Attributes()
	id, name := f["tool_id"], f["tool_name"]
	if id.IsUnknown() || name.IsUnknown() {
		return
	}
	if id.IsNull() == name.IsNull() {
		resp.Diagnostics.AddError("Invalid Tool selector", "Specify exactly one of tool_id or tool_name.")
	}
	if id := stringInput(f, "tool_id"); id != nil && *id == "" {
		resp.Diagnostics.AddError("Invalid Tool selector", "tool_id cannot be empty.")
	}
	if name := stringInput(f, "tool_name"); name != nil && *name == "" {
		resp.Diagnostics.AddError("Invalid Tool selector", "tool_name cannot be empty.")
	}
}
func (d *SandboxToolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config types.Object
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	fields := config.Attributes()
	id, name := stringInput(fields, "tool_id"), stringInput(fields, "tool_name")
	if (id == nil) == (name == nil) || id != nil && *id == "" || name != nil && *name == "" {
		resp.Diagnostics.AddError("Invalid Tool selector", "Exactly one known, nonempty tool_id or tool_name is required.")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	svc := d.testService
	if svc == nil {
		client, err := d.Client().UseAgsClient()
		if err != nil {
			resp.Diagnostics.AddError("AGS client unavailable", err.Error())
			return
		}
		svc = NewService(client)
	}
	var toolID string
	if id != nil {
		toolID = *id
	} else {
		tool, err := svc.FindByName(ctx, *name)
		if err != nil {
			resp.Diagnostics.AddError("AGS Tool query failed", err.Error())
			return
		}
		toolID = value(tool.ToolId)
	}
	tool, err := svc.Lookup(ctx, toolID)
	if err != nil {
		resp.Diagnostics.AddError("AGS Tool query failed", err.Error())
		return
	}
	projection, err := effectiveProjection(tool)
	if err != nil {
		resp.Diagnostics.AddError("Incomplete AGS Tool configuration", err.Error())
		return
	}
	response := map[string]any{"tool_id": optionalScalar(id), "tool_name": optionalScalar(name), "id": toolID, "status": value(tool.Status), "status_reason": value(tool.StatusReason), "image_digest": nil, "configuration": projection}
	if tool.CustomConfiguration != nil {
		response["image_digest"] = optionalScalar(tool.CustomConfiguration.ImageDigest)
	}
	resp.State.Schema = req.Config.Schema
	state, err := reflectedValue(ctx, req.Config.Schema.Type(), response)
	if err != nil {
		resp.Diagnostics.AddError("Cannot store AGS query result", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

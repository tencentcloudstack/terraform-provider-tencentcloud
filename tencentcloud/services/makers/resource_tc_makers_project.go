package makers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions (repo convention: only assert the
// extension interfaces; the base interface is checked implicitly by the
// factory return value).
var (
	_ resource.ResourceWithConfigure   = &MakersProjectResource{}
	_ resource.ResourceWithImportState = &MakersProjectResource{}
)

// NewMakersProjectResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewMakersProjectResource() resource.Resource {
	return &MakersProjectResource{}
}

// MakersProjectResource implements resource.Resource for
// tencentcloud_makers_project.
type MakersProjectResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// MakersProjectResourceModel maps the schema attributes.
type MakersProjectResourceModel struct {
	Id              types.String `tfsdk:"id"`
	ZoneId          types.String `tfsdk:"zone_id"`
	Name            types.String `tfsdk:"name"`
	Channel         types.String `tfsdk:"channel"`
	Provider        types.String `tfsdk:"code_provider"`
	Source          types.String `tfsdk:"source"`
	Area            types.String `tfsdk:"area"`
	NodeJsVersion   types.String `tfsdk:"node_js_version"`
	RootDir         types.String `tfsdk:"root_dir"`
	OutputDir       types.String `tfsdk:"output_dir"`
	BuildCmd        types.String `tfsdk:"build_cmd"`
	InstallCmd      types.String `tfsdk:"install_cmd"`
	Framework       types.String `tfsdk:"framework"`
	TemplateId      types.String `tfsdk:"template_id"`
	RepoUrl         types.String `tfsdk:"repo_url"`
	RepoOwner       types.String `tfsdk:"repo_owner"`
	RepoName        types.String `tfsdk:"repo_name"`
	RepoBranch      types.String `tfsdk:"repo_branch"`
	Namespace       types.String `tfsdk:"namespace"`
	CodingProjectId types.Int64  `tfsdk:"coding_project_id"`
	Status          types.String `tfsdk:"status"`
	PresetDomain    types.String `tfsdk:"preset_domain"`
	CreatedOn       types.String `tfsdk:"created_on"`
	ModifiedOn      types.String `tfsdk:"modified_on"`
}

func (r *MakersProjectResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_makers_project"
}

func (r *MakersProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a Makers project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, the Makers project ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Makers zone ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Project name.",
			},
			"channel": schema.StringAttribute{
				Required:    true,
				Description: "Creation channel. Valid values: `Template`, `Custom`, `CLI`, `Guest`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// The attribute is named code_provider instead of provider because
			// "provider" is a reserved root attribute name in Terraform.
			"code_provider": schema.StringAttribute{
				Optional:    true,
				Description: "Code source. Valid values: `Upload`, `Github`, `Gitee`, `Gitlab`, `Coding`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source": schema.StringAttribute{
				Optional:    true,
				Description: "Integrator identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"area": schema.StringAttribute{
				Optional:    true,
				Description: "Acceleration area. Valid values: `mainland`, `overseas`, `global`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_js_version": schema.StringAttribute{
				Optional: true,
				Description: "Node.js version, for example `20.18.0`. The value is echoed back by the cloud API " +
					"as `nodejs_version`.",
			},
			"root_dir": schema.StringAttribute{
				Optional:    true,
				Description: "Project root directory.",
			},
			"output_dir": schema.StringAttribute{
				Optional:    true,
				Description: "Build output directory.",
			},
			"build_cmd": schema.StringAttribute{
				Optional:    true,
				Description: "Build command.",
			},
			"install_cmd": schema.StringAttribute{
				Optional:    true,
				Description: "Install command.",
			},
			"framework": schema.StringAttribute{
				Optional:    true,
				Description: "Framework.",
			},
			"template_id": schema.StringAttribute{
				Optional:    true,
				Description: "Template ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repo_url": schema.StringAttribute{
				Optional:    true,
				Description: "Repository URL.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repo_owner": schema.StringAttribute{
				Optional:    true,
				Description: "Repository owner.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repo_name": schema.StringAttribute{
				Optional:    true,
				Description: "Repository name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repo_branch": schema.StringAttribute{
				Optional:    true,
				Description: "Production branch.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				Optional:    true,
				Description: "Coding namespace. Used when `provider` is `Coding`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"coding_project_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Coding project ID. Used when `provider` is `Coding`.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Project status. Valid values: `Normal`, `Ban`, `Unlink`, `Process`, `Pending`, `Offline`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"preset_domain": schema.StringAttribute{
				Computed:    true,
				Description: "Preset domain of the project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_on": schema.StringAttribute{
				Computed:    true,
				Description: "Last modification time of the project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *MakersProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MakersProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdk.NewCreateProjectRequest()
	request.ZoneId = plan.ZoneId.ValueStringPointer()
	request.Name = plan.Name.ValueStringPointer()
	request.Channel = plan.Channel.ValueStringPointer()
	if !plan.Provider.IsNull() && !plan.Provider.IsUnknown() {
		request.Provider = plan.Provider.ValueStringPointer()
	}
	if !plan.Source.IsNull() && !plan.Source.IsUnknown() {
		request.Source = plan.Source.ValueStringPointer()
	}
	if !plan.Area.IsNull() && !plan.Area.IsUnknown() {
		request.Area = plan.Area.ValueStringPointer()
	}
	if !plan.NodeJsVersion.IsNull() && !plan.NodeJsVersion.IsUnknown() {
		request.NodeJsVersion = plan.NodeJsVersion.ValueStringPointer()
	}
	if !plan.RootDir.IsNull() && !plan.RootDir.IsUnknown() {
		request.RootDir = plan.RootDir.ValueStringPointer()
	}
	if !plan.OutputDir.IsNull() && !plan.OutputDir.IsUnknown() {
		request.OutputDir = plan.OutputDir.ValueStringPointer()
	}
	if !plan.BuildCmd.IsNull() && !plan.BuildCmd.IsUnknown() {
		request.BuildCmd = plan.BuildCmd.ValueStringPointer()
	}
	if !plan.InstallCmd.IsNull() && !plan.InstallCmd.IsUnknown() {
		request.InstallCmd = plan.InstallCmd.ValueStringPointer()
	}
	if !plan.Framework.IsNull() && !plan.Framework.IsUnknown() {
		request.Framework = plan.Framework.ValueStringPointer()
	}
	if !plan.TemplateId.IsNull() && !plan.TemplateId.IsUnknown() {
		request.TemplateId = plan.TemplateId.ValueStringPointer()
	}
	if !plan.RepoUrl.IsNull() && !plan.RepoUrl.IsUnknown() {
		request.RepoUrl = plan.RepoUrl.ValueStringPointer()
	}
	if !plan.RepoOwner.IsNull() && !plan.RepoOwner.IsUnknown() {
		request.RepoOwner = plan.RepoOwner.ValueStringPointer()
	}
	if !plan.RepoName.IsNull() && !plan.RepoName.IsUnknown() {
		request.RepoName = plan.RepoName.ValueStringPointer()
	}
	if !plan.RepoBranch.IsNull() && !plan.RepoBranch.IsUnknown() {
		request.RepoBranch = plan.RepoBranch.ValueStringPointer()
	}
	if !plan.Namespace.IsNull() && !plan.Namespace.IsUnknown() {
		request.Namespace = plan.Namespace.ValueStringPointer()
	}
	if !plan.CodingProjectId.IsNull() && !plan.CodingProjectId.IsUnknown() {
		request.CodingProjectId = plan.CodingProjectId.ValueInt64Pointer()
	}

	service := NewMakersService(r.Client())
	projectId, err := service.CreateProject(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Makers project",
			fmt.Sprintf("CreateProject failed: %v", err),
		)
		return
	}

	// A nil id is an API failure, not an empty result. Never write an empty
	// id into state.
	if projectId == nil || *projectId == "" {
		resp.Diagnostics.AddError(
			"Error creating Makers project",
			"CreateProject returned an empty project id.",
		)
		return
	}

	plan.Id = types.StringValue(*projectId)

	// Computed attributes are unknown at this point; a state cannot hold
	// unknown values, so they are normalized to null before the defensive
	// persist below. The read-back afterwards fills them with cloud values.
	normalizeComputedUnknowns(&plan)

	// Persist the identity before reading back so that a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, plan.ZoneId.ValueString(), plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Makers project after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"Makers project not readable after create",
			fmt.Sprintf("project %s was created but cannot be read back", plan.Id.ValueString()),
		)
		return
	}
	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *MakersProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MakersProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, state.ZoneId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Makers project", err.Error())
		return
	}

	// Confirmed "not found" on the cloud side => drop it from state so
	// Terraform plans a recreation.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("Makers project %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(r.refreshModel(ctx, &state, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
func (r *MakersProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MakersProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state MakersProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdk.NewModifyProjectRequest()
	request.ZoneId = plan.ZoneId.ValueStringPointer()
	request.ProjectId = state.Id.ValueStringPointer()
	// ModifyProject only accepts a subset of the create-time fields; fields
	// it does not support are marked RequiresReplace in the schema.
	//
	// The cloud API does not document whether an omitted optional field means
	// "keep the previous value" or "clear it", so removing one of these
	// attributes from the configuration is rejected loudly instead of
	// silently keeping a stale value (permanent diff) or clearing user data.
	updatable := []struct {
		name     string
		planVal  types.String
		stateVal types.String
	}{
		{"name", plan.Name, state.Name},
		{"root_dir", plan.RootDir, state.RootDir},
		{"output_dir", plan.OutputDir, state.OutputDir},
		{"build_cmd", plan.BuildCmd, state.BuildCmd},
		{"install_cmd", plan.InstallCmd, state.InstallCmd},
		{"framework", plan.Framework, state.Framework},
		{"node_js_version", plan.NodeJsVersion, state.NodeJsVersion},
	}
	for _, field := range updatable {
		if field.planVal.IsNull() && !field.stateVal.IsNull() {
			resp.Diagnostics.AddError(
				"Clear semantics of "+field.name+" is not confirmed",
				fmt.Sprintf(
					"Removing `%s` from the configuration requires knowing whether ModifyProject treats an "+
						"omitted value as \"keep\" or \"clear\". This behavior is not documented by the cloud API. "+
						"Either keep the attribute in the configuration or confirm the clear semantics first.",
					field.name,
				),
			)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		request.Name = plan.Name.ValueStringPointer()
	}
	if !plan.RootDir.IsNull() && !plan.RootDir.IsUnknown() {
		request.RootDir = plan.RootDir.ValueStringPointer()
	}
	if !plan.OutputDir.IsNull() && !plan.OutputDir.IsUnknown() {
		request.OutputDir = plan.OutputDir.ValueStringPointer()
	}
	if !plan.BuildCmd.IsNull() && !plan.BuildCmd.IsUnknown() {
		request.BuildCmd = plan.BuildCmd.ValueStringPointer()
	}
	if !plan.InstallCmd.IsNull() && !plan.InstallCmd.IsUnknown() {
		request.InstallCmd = plan.InstallCmd.ValueStringPointer()
	}
	if !plan.Framework.IsNull() && !plan.Framework.IsUnknown() {
		request.Framework = plan.Framework.ValueStringPointer()
	}
	if !plan.NodeJsVersion.IsNull() && !plan.NodeJsVersion.IsUnknown() {
		request.NodeJsVersion = plan.NodeJsVersion.ValueStringPointer()
	}

	service := NewMakersService(r.Client())
	err := service.ModifyProject(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Makers project",
			fmt.Sprintf("ModifyProject failed: %v", err),
		)
		return
	}

	item, err := r.read(ctx, plan.ZoneId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Makers project after update", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"Makers project not readable after update",
			fmt.Sprintf("project %s was updated but cannot be read back", state.Id.ValueString()),
		)
		return
	}
	plan.Id = state.Id
	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *MakersProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MakersProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdk.NewDeleteProjectRequest()
	request.ZoneId = state.ZoneId.ValueStringPointer()
	request.ProjectId = state.Id.ValueStringPointer()

	service := NewMakersService(r.Client())
	err := service.DeleteProject(ctx, request)
	if err != nil {
		// No documented "already gone" error code exists for DeleteProject, so
		// every failure is surfaced to keep the id in state for a retry.
		resp.Diagnostics.AddError(
			"Error deleting Makers project",
			fmt.Sprintf("DeleteProject failed: %v", err),
		)
		return
	}
}

// ImportState is embedded via fw.WithImportByID for `terraform import <id>`
// passthrough.

// read fetches the project by zone and project id.
//
// Two different "empty" outcomes must stay distinguishable:
//   - API failure (nil response / nil Response) -> error
//   - queried successfully but nothing matched  -> (nil, nil) => "does not exist"
func (r *MakersProjectResource) read(ctx context.Context, zoneId, projectId string) (*sdk.Project, error) {
	request := sdk.NewDescribeProjectsRequest()
	request.ZoneId = helper.String(zoneId)
	request.ProjectIds = []*string{helper.String(projectId)}

	service := NewMakersService(r.Client())
	return service.DescribeProject(ctx, request)
}

// refreshModel copies cloud-owned fields into the model while keeping the
// user-supplied values for fields the cloud API does not echo back.
func (r *MakersProjectResource) refreshModel(ctx context.Context, m *MakersProjectResourceModel, item *sdk.Project) diag.Diagnostics {
	if item == nil {
		return nil
	}
	// Fields the cloud API returns: refresh them from the response.
	m.Name = types.StringPointerValue(item.Name)
	m.Area = types.StringPointerValue(item.Area)
	m.Provider = types.StringPointerValue(item.Provider)
	m.RootDir = types.StringPointerValue(item.RootDir)
	m.OutputDir = types.StringPointerValue(item.OutputDir)
	m.BuildCmd = types.StringPointerValue(item.BuildCmd)
	m.InstallCmd = types.StringPointerValue(item.InstallCmd)
	m.Framework = types.StringPointerValue(item.Framework)
	m.RepoUrl = types.StringPointerValue(item.RepoUrl)
	m.RepoOwner = types.StringPointerValue(item.RepoOwner)
	m.RepoName = types.StringPointerValue(item.RepoName)
	m.RepoBranch = types.StringPointerValue(item.RepoBranch)
	m.Status = types.StringPointerValue(item.Status)
	m.PresetDomain = types.StringPointerValue(item.PresetDomain)
	m.CreatedOn = types.StringPointerValue(item.CreatedOn)
	m.ModifiedOn = types.StringPointerValue(item.ModifiedOn)

	// node_js_version: the create request field is NodeJsVersion while the
	// response field is NodejsVersion (same value, different casing). Only
	// overwrite when the cloud actually returned a value, otherwise keep the
	// user-supplied value to avoid a spurious null.
	if item.NodejsVersion != nil {
		m.NodeJsVersion = types.StringPointerValue(item.NodejsVersion)
	}
	return nil
}

// normalizeComputedUnknowns rewrites unknown computed attributes to null so
// the model can be persisted into state before the read-back fills them with
// the cloud-side values.
func normalizeComputedUnknowns(m *MakersProjectResourceModel) {
	if m.Status.IsUnknown() {
		m.Status = types.StringNull()
	}
	if m.PresetDomain.IsUnknown() {
		m.PresetDomain = types.StringNull()
	}
	if m.CreatedOn.IsUnknown() {
		m.CreatedOn = types.StringNull()
	}
	if m.ModifiedOn.IsUnknown() {
		m.ModifiedOn = types.StringNull()
	}
}

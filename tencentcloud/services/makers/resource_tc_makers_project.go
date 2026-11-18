package makers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
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

// MakersProjectEnvVarModel maps the nested env_vars attributes.
type MakersProjectEnvVarModel struct {
	Key     types.String `tfsdk:"key"`
	Value   types.String `tfsdk:"value"`
	Id      types.Int64  `tfsdk:"id"`
	Comment types.String `tfsdk:"comment"`
	Env     types.Set    `tfsdk:"env"`
}

// MakersProjectResourceModel maps the schema attributes.
type MakersProjectResourceModel struct {
	Id              types.String `tfsdk:"id"`
	ZoneId          types.String `tfsdk:"zone_id"`
	Name            types.String `tfsdk:"name"`
	Channel         types.String `tfsdk:"channel"`
	Provider        types.String `tfsdk:"provider"`
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
	EnvVars         types.List   `tfsdk:"env_vars"`
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
			"provider": schema.StringAttribute{
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
			"env_vars": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Initial environment variables.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:    true,
							Description: "Variable name, up to 255 characters.",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "Variable value, up to 1000 characters.",
						},
						"id": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Variable ID. It is not sent on create.",
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Remark, up to 255 characters.",
						},
						"env": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Environment slugs the variable takes effect on, for example `Production`, " +
								"`Preview`. Multiple rows of the same key cannot overlap environments.",
						},
					},
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
	if !plan.EnvVars.IsNull() && !plan.EnvVars.IsUnknown() {
		envVars, diags := envVarModelsToSdk(ctx, plan.EnvVars)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.EnvVars = envVars
	}
	if !plan.Namespace.IsNull() && !plan.Namespace.IsUnknown() {
		request.Namespace = plan.Namespace.ValueStringPointer()
	}
	if !plan.CodingProjectId.IsNull() && !plan.CodingProjectId.IsUnknown() {
		request.CodingProjectId = plan.CodingProjectId.ValueInt64Pointer()
	}

	response, err := r.Client().UseMakersClient().CreateProjectWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Makers project",
			fmt.Sprintf("CreateProject failed: %v", err),
		)
		return
	}

	// A nil response (or a nil id field) is an API failure, not an empty
	// result. Never write an empty id into state.
	if response == nil || response.Response == nil || response.Response.ProjectId == nil || *response.Response.ProjectId == "" {
		resp.Diagnostics.AddError(
			"Error creating Makers project",
			"CreateProject returned an empty response or an empty project id.",
		)
		return
	}

	plan.Id = types.StringValue(*response.Response.ProjectId)

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

	_, err := r.Client().UseMakersClient().ModifyProjectWithContext(ctx, request)
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

	_, err := r.Client().UseMakersClient().DeleteProjectWithContext(ctx, request)
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

	response, err := r.Client().UseMakersClient().DescribeProjectsWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("DescribeProjects returned an empty response for project %s", projectId)
	}
	if len(response.Response.Projects) == 0 {
		return nil, nil
	}
	for _, project := range response.Response.Projects {
		if project.ProjectId != nil && *project.ProjectId == projectId {
			return project, nil
		}
	}
	return nil, nil
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

	// env_vars: the response only echoes Key/Value; Id, Comment and Env are
	// preserved from the model so a partial response does not silently clear
	// the user's configuration.
	if item.EnvVars != nil {
		envVars, diags := projectEnvVarsFromSdk(ctx, item.EnvVars, m.EnvVars)
		if diags.HasError() {
			return diags
		}
		m.EnvVars = envVars
	} else if !m.EnvVars.IsNull() {
		// The cloud returned no env vars: keep the configured list (the API
		// response is known to omit the field) but normalize unknown leaf
		// values so they can be persisted.
		envVars, diags := normalizeEnvVarUnknowns(ctx, m.EnvVars)
		if diags.HasError() {
			return diags
		}
		m.EnvVars = envVars
	}
	return nil
}

// normalizeEnvVarUnknowns rewrites unknown id/comment/env values inside the
// env_vars list to null so the list can be persisted into state.
func normalizeEnvVarUnknowns(ctx context.Context, list types.List) (types.List, diag.Diagnostics) {
	var models []MakersProjectEnvVarModel
	if diags := list.ElementsAs(ctx, &models, false); diags.HasError() {
		return types.ListNull(EnvVarObjectType()), diags
	}
	for i := range models {
		if models[i].Id.IsUnknown() {
			models[i].Id = types.Int64Null()
		}
		if models[i].Comment.IsUnknown() {
			models[i].Comment = types.StringNull()
		}
		if models[i].Env.IsUnknown() {
			models[i].Env = types.SetNull(types.StringType)
		}
	}
	return types.ListValueFrom(ctx, EnvVarObjectType(), models)
}

// envVarModelsToSdk converts the terraform list into the SDK EnvVar slice.
func envVarModelsToSdk(ctx context.Context, list types.List) ([]*sdk.EnvVar, diag.Diagnostics) {
	var models []MakersProjectEnvVarModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	envVars := make([]*sdk.EnvVar, 0, len(models))
	for _, model := range models {
		envVar := &sdk.EnvVar{
			Key:   model.Key.ValueStringPointer(),
			Value: model.Value.ValueStringPointer(),
		}
		if !model.Comment.IsNull() && !model.Comment.IsUnknown() {
			envVar.Comment = model.Comment.ValueStringPointer()
		}
		if !model.Env.IsNull() && !model.Env.IsUnknown() {
			envs := make([]*string, 0, len(model.Env.Elements()))
			for _, e := range model.Env.Elements() {
				envs = append(envs, helper.String(e.(types.String).ValueString()))
			}
			envVar.Env = envs
		}
		// Id is the cloud-assigned variable id: never sent on create.
		envVars = append(envVars, envVar)
	}
	return envVars, diags
}

// projectEnvVarsFromSdk converts the SDK ProjectEnvVar slice into the
// terraform list, preserving Id/Comment/Env from the previous model values
// because DescribeProjects only returns Key and Value.
func projectEnvVarsFromSdk(ctx context.Context, items []*sdk.ProjectEnvVar, previous types.List) (types.List, diag.Diagnostics) {
	var prevModels []MakersProjectEnvVarModel
	if !previous.IsNull() && !previous.IsUnknown() {
		if diags := previous.ElementsAs(ctx, &prevModels, false); diags.HasError() {
			return types.ListNull(EnvVarObjectType()), diags
		}
	}
	prevByKey := make(map[string]MakersProjectEnvVarModel, len(prevModels))
	for _, prev := range prevModels {
		prevByKey[prev.Key.ValueString()] = prev
	}

	out := make([]MakersProjectEnvVarModel, 0, len(items))
	for _, item := range items {
		key := types.StringPointerValue(item.Key)
		value := types.StringPointerValue(item.Value)
		prev, ok := prevByKey[key.ValueString()]
		next := MakersProjectEnvVarModel{
			Key:   key,
			Value: value,
		}
		if ok {
			next.Id = prev.Id
			next.Comment = prev.Comment
			next.Env = prev.Env
		} else {
			next.Id = types.Int64Null()
			next.Comment = types.StringNull()
			next.Env = types.SetNull(types.StringType)
		}
		// The id is Optional+Computed: on create it is unknown in the plan and
		// DescribeProjects does not echo it back, so it cannot be persisted as
		// unknown — normalize to null.
		if next.Id.IsUnknown() {
			next.Id = types.Int64Null()
		}
		if next.Comment.IsUnknown() {
			next.Comment = types.StringNull()
		}
		if next.Env.IsUnknown() {
			next.Env = types.SetNull(types.StringType)
		}
		out = append(out, next)
	}
	return types.ListValueFrom(ctx, EnvVarObjectType(), out)
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

// EnvVarObjectType returns the object type of one env_vars element.
func EnvVarObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":     types.StringType,
			"value":   types.StringType,
			"id":      types.Int64Type,
			"comment": types.StringType,
			"env":     types.SetType{ElemType: types.StringType},
		},
	}
}

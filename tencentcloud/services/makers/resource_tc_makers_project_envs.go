package makers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	makersv20261113 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// Compile-time interface assertions.
var (
	_ resource.ResourceWithConfigure   = &MakersProjectEnvsResource{}
	_ resource.ResourceWithImportState = &MakersProjectEnvsResource{}
)

// NewMakersProjectEnvsResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewMakersProjectEnvsResource() resource.Resource {
	return &MakersProjectEnvsResource{}
}

// MakersProjectEnvsResource implements resource.Resource for
// tencentcloud_makers_project_envs.
//
// The cloud API reuses ModifyProjectEnvs for both create and update: the
// request carries the full desired set of environment variables scoped by
// ZoneId + ProjectId, and there is no server-assigned single object id.
// The Terraform resource id is therefore composite: `<zone_id>#<project_id>`.
type MakersProjectEnvsResource struct {
	fw.ResourceWithConfigure
}

// envVarAttrTypes is the attribute type set of one `env_vars` element,
// used to build object values when refreshing state from the cloud response.
var envVarAttrTypes = map[string]attr.Type{
	"key":     types.StringType,
	"value":   types.StringType,
	"id":      types.Int64Type,
	"comment": types.StringType,
	"env":     types.SetType{ElemType: types.StringType},
}

// envVarObjectType is the object type of one `env_vars` element.
var envVarObjectType = types.ObjectType{AttrTypes: envVarAttrTypes}

// MakersProjectEnvVarModel maps one element of the `env_vars` attribute.
type MakersProjectEnvVarModel struct {
	Key     types.String `tfsdk:"key"`
	Value   types.String `tfsdk:"value"`
	Id      types.Int64  `tfsdk:"id"`
	Comment types.String `tfsdk:"comment"`
	Env     types.Set    `tfsdk:"env"`
}

// MakersProjectEnvsResourceModel maps the schema attributes.
type MakersProjectEnvsResourceModel struct {
	Id        types.String               `tfsdk:"id"`
	ZoneId    types.String               `tfsdk:"zone_id"`
	ProjectId types.String               `tfsdk:"project_id"`
	EnvVars   []MakersProjectEnvVarModel `tfsdk:"env_vars"`
}

func (r *MakersProjectEnvsResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_makers_project_envs"
}

func (r *MakersProjectEnvsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage environment variables of a Makers project.\n\n" +
			"Example Usage\n\n" +
			"```hcl\n" +
			"resource \"tencentcloud_makers_project_envs\" \"example\" {\n" +
			"  zone_id    = \"zone-xxxx\"\n" +
			"  project_id = \"proj-yyyy\"\n" +
			"  env_vars {\n" +
			"    key     = \"API_KEY\"\n" +
			"    value   = \"secret-value\"\n" +
			"    comment = \"production key\"\n" +
			"    env     = [\"Production\"]\n" +
			"  }\n" +
			"}\n" +
			"```\n\n" +
			"Import\n\n" +
			"tencentcloud_makers_project_envs can be imported using the `<zone_id>#<project_id>`, e.g.\n\n" +
			"```\n" +
			"$ terraform import tencentcloud_makers_project_envs.example zone-xxxx#proj-yyyy\n" +
			"```",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource in the form `<zone_id>#<project_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// Both locating fields are immutable: the whole variable set is
			// scoped by ZoneId + ProjectId, and a different scope is a
			// different managed object.
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Makers site ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "Project ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"env_vars": schema.ListNestedAttribute{
				Required:    true,
				Description: "The environment variables to write.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:    true,
							Description: "Variable name, at most 255 characters.",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "Variable value, at most 1000 characters.",
						},
						"id": schema.Int64Attribute{
							Optional:    true,
							Description: "Variable ID. Not required when creating.",
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Comment, at most 255 characters.",
						},
						"env": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Slugs of the environments the variable applies to, e.g. `Production`, `Preview`. Env sets of variables sharing the same key must not overlap.",
						},
					},
				},
			},
		},
	}
}

// compositeID builds the resource id from its locating fields.
func compositeID(zoneId, projectId string) string {
	return fmt.Sprintf("%s#%s", zoneId, projectId)
}

// splitCompositeID parses a `<zone_id>#<project_id>` resource id.
func splitCompositeID(id string) (string, string, error) {
	var zoneId, projectId string
	if n, err := fmt.Sscanf(id, "%[^#]#%s", &zoneId, &projectId); err != nil || n != 2 {
		return "", "", fmt.Errorf("invalid id %q, expected format <zone_id>#<project_id>", id)
	}
	return zoneId, projectId, nil
}

// requestEnvVars converts the plan/state list into SDK EnvVar items.
func requestEnvVars(models []MakersProjectEnvVarModel) []*makersv20261113.EnvVar {
	items := make([]*makersv20261113.EnvVar, 0, len(models))
	for _, m := range models {
		item := &makersv20261113.EnvVar{
			Key:   m.Key.ValueStringPointer(),
			Value: m.Value.ValueStringPointer(),
		}
		if !m.Id.IsNull() && !m.Id.IsUnknown() {
			item.Id = m.Id.ValueInt64Pointer()
		}
		if !m.Comment.IsNull() && !m.Comment.IsUnknown() {
			item.Comment = m.Comment.ValueStringPointer()
		}
		if !m.Env.IsNull() && !m.Env.IsUnknown() {
			var envs []string
			_ = m.Env.ElementsAs(context.Background(), &envs, false)
			for _, e := range envs {
				env := e
				item.Env = append(item.Env, &env)
			}
		}
		items = append(items, item)
	}
	return items
}

// refreshEnvVars rebuilds the model list from the Describe response.
// The response does not echo ZoneId/ProjectId back, so those keep the
// state values.
func (r *MakersProjectEnvsResource) refreshEnvVars(m *MakersProjectEnvsResourceModel, items []*makersv20261113.EnvVar) {
	models := make([]MakersProjectEnvVarModel, 0, len(items))
	for _, item := range items {
		model := MakersProjectEnvVarModel{
			Key:     types.StringPointerValue(item.Key),
			Value:   types.StringPointerValue(item.Value),
			Id:      types.Int64PointerValue(item.Id),
			Comment: types.StringPointerValue(item.Comment),
		}
		if len(item.Env) == 0 {
			model.Env = types.SetNull(types.StringType)
		} else {
			elems := make([]attr.Value, 0, len(item.Env))
			for _, e := range item.Env {
				elems = append(elems, types.StringPointerValue(e))
			}
			set, _ := types.SetValue(types.StringType, elems)
			model.Env = set
		}
		models = append(models, model)
	}
	m.EnvVars = models
}

// read fetches the environment variables of the project.
//
//   - API failure (nil response / nil Response) -> error
//   - queried successfully but no variables exist -> nil, meaning "does not exist"
func (r *MakersProjectEnvsResource) read(ctx context.Context, zoneId, projectId string) ([]*makersv20261113.EnvVar, error) {
	request := makersv20261113.NewDescribeProjectEnvsRequest()
	request.ZoneId = &zoneId
	request.ProjectId = &projectId

	response, err := r.Client().UseMakersClient().DescribeProjectEnvsWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("DescribeProjectEnvs returned an empty response for %s", compositeID(zoneId, projectId))
	}
	if len(response.Response.EnvVars) == 0 {
		return nil, nil
	}
	return response.Response.EnvVars, nil
}

// Create handles `terraform apply` for new resources.
func (r *MakersProjectEnvsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MakersProjectEnvsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := makersv20261113.NewModifyProjectEnvsRequest()
	request.ZoneId = plan.ZoneId.ValueStringPointer()
	request.ProjectId = plan.ProjectId.ValueStringPointer()
	request.EnvVars = requestEnvVars(plan.EnvVars)

	_, err := r.Client().UseMakersClient().ModifyProjectEnvsWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating makers project envs",
			fmt.Sprintf("ModifyProjectEnvs failed: %v", err),
		)
		return
	}

	// The API returns no resource-level id; the composite id comes from the
	// locating fields. Persist the identity BEFORE the read-back so the
	// resource stays destroyable even if the read fails.
	plan.Id = types.StringValue(compositeID(plan.ZoneId.ValueString(), plan.ProjectId.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	items, err := r.read(ctx, plan.ZoneId.ValueString(), plan.ProjectId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading makers project envs after create", err.Error())
		return
	}
	if items == nil {
		resp.Diagnostics.AddError(
			"makers project envs not readable after create",
			fmt.Sprintf("resource %s was created but cannot be read back", plan.Id.ValueString()),
		)
		return
	}
	r.refreshEnvVars(&plan, items)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *MakersProjectEnvsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MakersProjectEnvsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	items, err := r.read(ctx, state.ZoneId.ValueString(), state.ProjectId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading makers project envs", err.Error())
		return
	}

	// Confirmed "no variables exist" => drop it from state so Terraform
	// plans a recreation. A nil response is an API failure and must not
	// land here.
	if items == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("makers project envs %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	r.refreshEnvVars(&state, items)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *MakersProjectEnvsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MakersProjectEnvsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state MakersProjectEnvsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ModifyProjectEnvs writes the full desired set of variables, so the
	// plan list is sent as-is; variables removed from the configuration are
	// not part of the write and should be removed explicitly beforehand.
	// Variables present in state but dropped from the plan are deleted
	// through DeleteProjectEnvs, using their cloud-assigned ids when known.
	for _, old := range state.EnvVars {
		found := false
		for _, newVar := range plan.EnvVars {
			if old.Key.Equal(newVar.Key) {
				found = true
				break
			}
		}
		if !found {
			delReq := makersv20261113.NewDeleteProjectEnvsRequest()
			delReq.ZoneId = plan.ZoneId.ValueStringPointer()
			delReq.ProjectId = plan.ProjectId.ValueStringPointer()
			delReq.EnvVars = requestEnvVars([]MakersProjectEnvVarModel{old})
			if _, err := r.Client().UseMakersClient().DeleteProjectEnvsWithContext(ctx, delReq); err != nil {
				resp.Diagnostics.AddError(
					"Error updating makers project envs",
					fmt.Sprintf("DeleteProjectEnvs failed while removing variable %s: %v", old.Key.ValueString(), err),
				)
				return
			}
		}
	}

	request := makersv20261113.NewModifyProjectEnvsRequest()
	request.ZoneId = plan.ZoneId.ValueStringPointer()
	request.ProjectId = plan.ProjectId.ValueStringPointer()
	request.EnvVars = requestEnvVars(plan.EnvVars)

	if _, err := r.Client().UseMakersClient().ModifyProjectEnvsWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating makers project envs",
			fmt.Sprintf("ModifyProjectEnvs failed: %v", err),
		)
		return
	}

	items, err := r.read(ctx, plan.ZoneId.ValueString(), plan.ProjectId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading makers project envs after update", err.Error())
		return
	}
	if items == nil {
		resp.Diagnostics.AddError(
			"makers project envs not readable after update",
			fmt.Sprintf("resource %s was updated but cannot be read back", state.Id.ValueString()),
		)
		return
	}
	plan.Id = state.Id
	r.refreshEnvVars(&plan, items)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *MakersProjectEnvsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MakersProjectEnvsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := makersv20261113.NewDeleteProjectEnvsRequest()
	request.ZoneId = state.ZoneId.ValueStringPointer()
	request.ProjectId = state.ProjectId.ValueStringPointer()
	request.EnvVars = requestEnvVars(state.EnvVars)

	if _, err := r.Client().UseMakersClient().DeleteProjectEnvsWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting makers project envs",
			fmt.Sprintf("DeleteProjectEnvs failed: %v", err),
		)
		return
	}
}

// ImportState parses the composite `<zone_id>#<project_id>` id.
func (r *MakersProjectEnvsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	zoneId, projectId, err := splitCompositeID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format <zone_id>#<project_id>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), zoneId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectId)...)
}

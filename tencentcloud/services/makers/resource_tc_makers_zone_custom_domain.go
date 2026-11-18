package makers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	makerssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var (
	_ resource.ResourceWithConfigure   = &MakersZoneCustomDomainResource{}
	_ resource.ResourceWithImportState = &MakersZoneCustomDomainResource{}
)

// NewMakersZoneCustomDomainResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewMakersZoneCustomDomainResource() resource.Resource {
	return &MakersZoneCustomDomainResource{}
}

// MakersZoneCustomDomainResource implements resource.Resource for
// tencentcloud_makers_zone_custom_domain.
//
// The custom domain is located by the composite key (zone id, project id,
// domain name): the create API returns a cloud-side id but the describe,
// update and delete APIs all address the domain by this composite key, so
// the Terraform id is built from the three segments.
type MakersZoneCustomDomainResource struct {
	fw.ResourceWithConfigure
}

// MakersZoneCustomDomainResourceModel maps the schema attributes.
type MakersZoneCustomDomainResourceModel struct {
	Id        types.String           `tfsdk:"id"`
	ZoneId    types.String           `tfsdk:"zone_id"`
	ProjectId types.String           `tfsdk:"project_id"`
	Domain    types.String           `tfsdk:"domain"`
	Env       types.String           `tfsdk:"env"`
	Timeouts  resourcetimeouts.Value `tfsdk:"timeouts"`
}

func (r *MakersZoneCustomDomainResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_makers_zone_custom_domain"
}

func (r *MakersZoneCustomDomainResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to bind a custom domain to a Makers zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource in the format `zone_id#project_id#domain`.",
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
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "Project ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "Custom domain name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"env": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Environment slug the domain is bound to. If omitted on create, the backend decides the environment.",
			},
		},
		// The timeouts definition uses a block (rather than an object
		// attribute) so the familiar `timeouts { ... }` HCL block syntax
		// works, matching the official Terraform examples and the SDKv2
		// providers.
		Blocks: map[string]schema.Block{
			"timeouts": resourcetimeouts.Block(ctx, resourcetimeouts.Opts{
				Create:            true,
				Update:            true,
				Delete:            true,
				CreateDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `30m`. It limits how long to wait for the async custom domain binding to converge. Default is `60m`.",
				UpdateDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `30m`. It limits how long to wait for the async environment switch to converge. Default is `60m`.",
				DeleteDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `30m`. It limits how long to wait for the async custom domain unbinding to converge. Default is `60m`.",
			}),
		},
	}
}

// Create binds the custom domain.
func (r *MakersZoneCustomDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MakersZoneCustomDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := makerssdk.NewCreateZoneCustomDomainRequest()
	request.ZoneId = plan.ZoneId.ValueStringPointer()
	request.ProjectId = plan.ProjectId.ValueStringPointer()
	request.Domain = plan.Domain.ValueStringPointer()
	if !plan.Env.IsNull() && !plan.Env.IsUnknown() {
		request.Env = plan.Env.ValueStringPointer()
	}

	service := NewMakersService(r.Client())
	domain, err := service.CreateZoneCustomDomain(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Makers zone custom domain",
			fmt.Sprintf("CreateZoneCustomDomain failed: %v", err),
		)
		return
	}
	if domain == "" {
		resp.Diagnostics.AddError(
			"Error creating Makers zone custom domain",
			"CreateZoneCustomDomain returned an empty response.",
		)
		return
	}

	// Persist the identity BEFORE any asynchronous wait: if the wait fails,
	// the resource is still in state and `terraform destroy` can clean it up.
	plan.Id = types.StringValue(buildZoneCustomDomainId(
		plan.ZoneId.ValueString(), plan.ProjectId.ValueString(), plan.Domain.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The binding is async on the cloud side: poll until the domain becomes
	// readable before the final read back.
	createTimeout, timeoutDiags := plan.Timeouts.Create(ctx, tccommon.WriteRetryTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := service.WaitForZoneCustomDomain(ctx,
		plan.ZoneId.ValueString(), plan.ProjectId.ValueString(), plan.Domain.ValueString(),
		createTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for Makers zone custom domain", err.Error())
		return
	}

	item, err := service.DescribeZoneCustomDomain(ctx, plan.ZoneId.ValueString(), plan.ProjectId.ValueString(), plan.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Makers zone custom domain after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"Makers zone custom domain not readable after create",
			fmt.Sprintf("domain %s was created but cannot be read back", plan.Domain.ValueString()),
		)
		return
	}
	refreshZoneCustomDomainModel(&plan, item, plan.Env)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *MakersZoneCustomDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MakersZoneCustomDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewMakersService(r.Client())
	item, err := service.DescribeZoneCustomDomain(ctx, state.ZoneId.ValueString(), state.ProjectId.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Makers zone custom domain", err.Error())
		return
	}

	// Confirmed "not found" on the cloud side => drop it from state so Terraform
	// plans a recreation. A nil response is an API failure and must not land here.
	if item == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	refreshZoneCustomDomainModel(&state, item, state.Env)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates. Only the environment can be changed
// (ModifyZoneCustomDomainEnv); any other change is a replacement.
func (r *MakersZoneCustomDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MakersZoneCustomDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state MakersZoneCustomDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewMakersService(r.Client())

	if !plan.Env.IsNull() && !plan.Env.IsUnknown() && plan.Env != state.Env {
		request := makerssdk.NewModifyZoneCustomDomainEnvRequest()
		request.ZoneId = plan.ZoneId.ValueStringPointer()
		request.ProjectId = plan.ProjectId.ValueStringPointer()
		request.Domain = plan.Domain.ValueStringPointer()
		request.Env = plan.Env.ValueStringPointer()

		if err := service.ModifyZoneCustomDomainEnv(ctx, request); err != nil {
			resp.Diagnostics.AddError(
				"Error updating Makers zone custom domain",
				fmt.Sprintf("ModifyZoneCustomDomainEnv failed: %v", err),
			)
			return
		}

		// The environment switch is async on the cloud side: poll until the
		// domain reports the expected environment slug.
		updateTimeout, timeoutDiags := plan.Timeouts.Update(ctx, tccommon.WriteRetryTimeout)
		resp.Diagnostics.Append(timeoutDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := service.WaitForZoneCustomDomainEnv(ctx,
			plan.ZoneId.ValueString(), plan.ProjectId.ValueString(), plan.Domain.ValueString(),
			plan.Env.ValueString(), updateTimeout); err != nil {
			resp.Diagnostics.AddError("Error waiting for Makers zone custom domain env", err.Error())
			return
		}
	}

	item, err := service.DescribeZoneCustomDomain(ctx, plan.ZoneId.ValueString(), plan.ProjectId.ValueString(), plan.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Makers zone custom domain after update", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"Makers zone custom domain not readable after update",
			fmt.Sprintf("domain %s was updated but cannot be read back", plan.Domain.ValueString()),
		)
		return
	}
	plan.Id = state.Id
	refreshZoneCustomDomainModel(&plan, item, plan.Env)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete unbinds the custom domain via the batch delete API, restricted to
// the single managed domain.
func (r *MakersZoneCustomDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MakersZoneCustomDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewMakersService(r.Client())
	if err := service.DeleteProjectCustomDomain(ctx, state.ZoneId.ValueString(), state.ProjectId.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Makers zone custom domain",
			fmt.Sprintf("DeleteProjectCustomDomains failed: %v", err),
		)
		return
	}

	// The unbinding is async on the cloud side: poll until the domain is no
	// longer visible so destroy does not report success too early.
	deleteTimeout, timeoutDiags := state.Timeouts.Delete(ctx, tccommon.WriteRetryTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := service.WaitForZoneCustomDomainDeleted(ctx,
		state.ZoneId.ValueString(), state.ProjectId.ValueString(), state.Domain.ValueString(),
		deleteTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for Makers zone custom domain deletion", err.Error())
		return
	}
}

// ImportState parses the composite id `zone_id#project_id#domain`.
func (r *MakersZoneCustomDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, tccommon.FILED_SP)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Makers zone custom domain id must be in the format `zone_id#project_id#domain`, got %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[2])...)
}

// refreshZoneCustomDomainModel copies the cloud-side domain into the model.
// The env attribute is kept from the plan/state when the cloud side does not
// report one, to avoid a permanent diff on the optional-computed attribute.
func refreshZoneCustomDomainModel(m *MakersZoneCustomDomainResourceModel, item *makerssdk.ZoneDomain, fallbackEnv types.String) {
	if item == nil {
		return
	}
	if item.Env != nil {
		m.Env = types.StringPointerValue(item.Env)
	} else {
		m.Env = fallbackEnv
	}
}

func buildZoneCustomDomainId(zoneId, projectId, domain string) string {
	return zoneId + tccommon.FILED_SP + projectId + tccommon.FILED_SP + domain
}

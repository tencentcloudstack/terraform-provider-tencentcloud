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

	makersv20261113 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

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
	Id           types.String       `tfsdk:"id"`
	ZoneId       types.String       `tfsdk:"zone_id"`
	ProjectId    types.String       `tfsdk:"project_id"`
	Domain       types.String       `tfsdk:"domain"`
	Env          types.String       `tfsdk:"env"`
	PagesDomains []*ZoneDomainModel `tfsdk:"pages_domains"`
}

// ZoneDomainModel maps the computed pages_domains list element.
type ZoneDomainModel struct {
	Type                  types.String                `tfsdk:"type"`
	Domain                types.String                `tfsdk:"domain"`
	Status                types.String                `tfsdk:"status"`
	Env                   types.String                `tfsdk:"env"`
	ZoneId                types.String                `tfsdk:"zone_id"`
	Area                  types.String                `tfsdk:"area"`
	Cname                 types.String                `tfsdk:"cname"`
	MainDomain            types.String                `tfsdk:"main_domain"`
	Ipv6                  *SwitchConfigModel          `tfsdk:"ipv6"`
	Http2                 *SwitchConfigModel          `tfsdk:"http2"`
	OcspStapling          *SwitchConfigModel          `tfsdk:"ocsp_stapling"`
	ForceRedirect         *ForceRedirectConfigModel   `tfsdk:"force_redirect"`
	Hsts                  *HstsConfigModel            `tfsdk:"hsts"`
	OwnershipVerification *OwnershipVerificationModel `tfsdk:"ownership_verification"`
}

// SwitchConfigModel maps the switch-only nested configs.
type SwitchConfigModel struct {
	Switch types.String `tfsdk:"switch"`
}

// ForceRedirectConfigModel maps the force redirect nested config.
type ForceRedirectConfigModel struct {
	Switch             types.String `tfsdk:"switch"`
	RedirectStatusCode types.Int64  `tfsdk:"redirect_status_code"`
}

// HstsConfigModel maps the HSTS nested config.
type HstsConfigModel struct {
	Switch            types.String `tfsdk:"switch"`
	MaxAge            types.Int64  `tfsdk:"max_age"`
	IncludeSubDomains types.String `tfsdk:"include_sub_domains"`
	Preload           types.String `tfsdk:"preload"`
}

// OwnershipVerificationModel maps the ownership verification nested config.
type OwnershipVerificationModel struct {
	DnsVerification  *DnsVerificationModel  `tfsdk:"dns_verification"`
	FileVerification *FileVerificationModel `tfsdk:"file_verification"`
	NsVerification   *NsVerificationModel   `tfsdk:"ns_verification"`
}

// DnsVerificationModel maps the DNS ownership verification.
type DnsVerificationModel struct {
	RecordType  types.String `tfsdk:"record_type"`
	Subdomain   types.String `tfsdk:"subdomain"`
	RecordValue types.String `tfsdk:"record_value"`
}

// FileVerificationModel maps the file ownership verification.
type FileVerificationModel struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
}

// NsVerificationModel maps the NS ownership verification.
type NsVerificationModel struct {
	NameServers types.Set `tfsdk:"name_servers"`
}

func (r *MakersZoneCustomDomainResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_makers_zone_custom_domain"
}

func (r *MakersZoneCustomDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	switchConfigAttributes := func() map[string]schema.Attribute {
		return map[string]schema.Attribute{
			"switch": schema.StringAttribute{
				Computed:    true,
				Description: "Switch state. Valid values: `on`, `off`.",
			},
		}
	}

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
			"pages_domains": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The current domain configuration returned by the cloud side. Only the managed custom domain is reported here.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Domain type. Valid values: `Custom`, `Preset`.",
						},
						"domain": schema.StringAttribute{
							Computed:    true,
							Description: "Domain name.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Domain status.",
						},
						"env": schema.StringAttribute{
							Computed:    true,
							Description: "Bound environment slug, defaults to Production.",
						},
						"zone_id": schema.StringAttribute{
							Computed:    true,
							Description: "Zone ID.",
						},
						"area": schema.StringAttribute{
							Computed:    true,
							Description: "Acceleration area.",
						},
						"cname": schema.StringAttribute{
							Computed:    true,
							Description: "CNAME.",
						},
						"main_domain": schema.StringAttribute{
							Computed:    true,
							Description: "Main domain.",
						},
						"ipv6": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "IPv6 configuration.",
							Attributes:  switchConfigAttributes(),
						},
						"http2": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "HTTP/2 configuration.",
							Attributes:  switchConfigAttributes(),
						},
						"ocsp_stapling": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "OCSP stapling configuration.",
							Attributes:  switchConfigAttributes(),
						},
						"force_redirect": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Force HTTPS redirect configuration.",
							Attributes: map[string]schema.Attribute{
								"switch": schema.StringAttribute{
									Computed:    true,
									Description: "Switch state. Valid values: `on`, `off`.",
								},
								"redirect_status_code": schema.Int64Attribute{
									Computed:    true,
									Description: "Redirect status code. Valid values: 301, 302, 307, 308.",
								},
							},
						},
						"hsts": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "HSTS configuration.",
							Attributes: map[string]schema.Attribute{
								"switch": schema.StringAttribute{
									Computed:    true,
									Description: "Switch state. Valid values: `on`, `off`.",
								},
								"max_age": schema.Int64Attribute{
									Computed:    true,
									Description: "Validity period in seconds.",
								},
								"include_sub_domains": schema.StringAttribute{
									Computed:    true,
									Description: "Whether to include sub domains. Valid values: `on`, `off`.",
								},
								"preload": schema.StringAttribute{
									Computed:    true,
									Description: "Whether to preload. Valid values: `on`, `off`.",
								},
							},
						},
						"ownership_verification": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Ownership verification information.",
							Attributes: map[string]schema.Attribute{
								"dns_verification": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "DNS verification.",
									Attributes: map[string]schema.Attribute{
										"record_type": schema.StringAttribute{
											Computed:    true,
											Description: "Record type, e.g. TXT.",
										},
										"subdomain": schema.StringAttribute{
											Computed:    true,
											Description: "Host record.",
										},
										"record_value": schema.StringAttribute{
											Computed:    true,
											Description: "Record value.",
										},
									},
								},
								"file_verification": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "File verification.",
									Attributes: map[string]schema.Attribute{
										"path": schema.StringAttribute{
											Computed:    true,
											Description: "Verification file path.",
										},
										"content": schema.StringAttribute{
											Computed:    true,
											Description: "Verification file content.",
										},
									},
								},
								"ns_verification": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "NS verification.",
									Attributes: map[string]schema.Attribute{
										"name_servers": schema.SetAttribute{
											Computed:    true,
											ElementType: types.StringType,
											Description: "Name servers.",
										},
									},
								},
							},
						},
					},
				},
			},
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

	request := makersv20261113.NewCreateZoneCustomDomainRequest()
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

	// The resource is addressed by the composite key everywhere, so build the
	// Terraform id from the plan values and persist it before the read back.
	plan.Id = types.StringValue(buildZoneCustomDomainId(
		plan.ZoneId.ValueString(), plan.ProjectId.ValueString(), plan.Domain.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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
	refreshZoneCustomDomainModel(ctx, &plan, item, plan.Env)
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
	refreshZoneCustomDomainModel(ctx, &state, item, state.Env)
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
		request := makersv20261113.NewModifyZoneCustomDomainEnvRequest()
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
	refreshZoneCustomDomainModel(ctx, &plan, item, plan.Env)
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

	request := makersv20261113.NewDeleteProjectCustomDomainsRequest()
	request.ZoneId = state.ZoneId.ValueStringPointer()
	request.ProjectId = state.ProjectId.ValueStringPointer()
	service := NewMakersService(r.Client())
	if err := service.DeleteProjectCustomDomain(ctx, state.ZoneId.ValueString(), state.ProjectId.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Makers zone custom domain",
			fmt.Sprintf("DeleteProjectCustomDomains failed: %v", err),
		)
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
	resp.State.SetAttribute(ctx, path.Root("id"), req.ID)
	resp.State.SetAttribute(ctx, path.Root("zone_id"), parts[0])
	resp.State.SetAttribute(ctx, path.Root("project_id"), parts[1])
	resp.State.SetAttribute(ctx, path.Root("domain"), parts[2])
}

// refreshZoneCustomDomainModel copies the cloud-side domain into the model.
// The env attribute is kept from the plan/state when the cloud side does not
// report one, to avoid a permanent diff on the optional-computed attribute.
func refreshZoneCustomDomainModel(ctx context.Context, m *MakersZoneCustomDomainResourceModel, item *makersv20261113.ZoneDomain, fallbackEnv types.String) {
	if item == nil {
		return
	}
	if item.Env != nil {
		m.Env = types.StringPointerValue(item.Env)
	} else {
		m.Env = fallbackEnv
	}
	m.PagesDomains = []*ZoneDomainModel{zoneDomainModelFromSdk(ctx, item)}
}

func zoneDomainModelFromSdk(ctx context.Context, item *makersv20261113.ZoneDomain) *ZoneDomainModel {
	model := &ZoneDomainModel{
		Type:         types.StringPointerValue(item.Type),
		Domain:       types.StringPointerValue(item.Domain),
		Status:       types.StringPointerValue(item.Status),
		Env:          types.StringPointerValue(item.Env),
		ZoneId:       types.StringPointerValue(item.ZoneId),
		Area:         types.StringPointerValue(item.Area),
		Cname:        types.StringPointerValue(item.Cname),
		MainDomain:   types.StringPointerValue(item.MainDomain),
		Ipv6:         switchConfigModelFromSdk(item.Ipv6),
		Http2:        switchConfigModelFromSdk(item.Http2),
		OcspStapling: switchConfigModelFromSdk(item.OcspStapling),
	}
	if item.ForceRedirect != nil {
		model.ForceRedirect = &ForceRedirectConfigModel{
			Switch:             types.StringPointerValue(item.ForceRedirect.Switch),
			RedirectStatusCode: types.Int64PointerValue(item.ForceRedirect.RedirectStatusCode),
		}
	}
	if item.Hsts != nil {
		model.Hsts = &HstsConfigModel{
			Switch:            types.StringPointerValue(item.Hsts.Switch),
			MaxAge:            types.Int64PointerValue(item.Hsts.MaxAge),
			IncludeSubDomains: types.StringPointerValue(item.Hsts.IncludeSubDomains),
			Preload:           types.StringPointerValue(item.Hsts.Preload),
		}
	}
	if item.OwnershipVerification != nil {
		ownership := &OwnershipVerificationModel{}
		if item.OwnershipVerification.DnsVerification != nil {
			ownership.DnsVerification = &DnsVerificationModel{
				RecordType:  types.StringPointerValue(item.OwnershipVerification.DnsVerification.RecordType),
				Subdomain:   types.StringPointerValue(item.OwnershipVerification.DnsVerification.Subdomain),
				RecordValue: types.StringPointerValue(item.OwnershipVerification.DnsVerification.RecordValue),
			}
		}
		if item.OwnershipVerification.FileVerification != nil {
			ownership.FileVerification = &FileVerificationModel{
				Path:    types.StringPointerValue(item.OwnershipVerification.FileVerification.Path),
				Content: types.StringPointerValue(item.OwnershipVerification.FileVerification.Content),
			}
		}
		if item.OwnershipVerification.NsVerification != nil {
			nameServers := types.SetNull(types.StringType)
			if len(item.OwnershipVerification.NsVerification.NameServers) > 0 {
				set, diags := types.SetValueFrom(ctx, types.StringType, item.OwnershipVerification.NsVerification.NameServers)
				if !diags.HasError() {
					nameServers = set
				}
			}
			ownership.NsVerification = &NsVerificationModel{
				NameServers: nameServers,
			}
		}
		model.OwnershipVerification = ownership
	}
	return model
}

func switchConfigModelFromSdk(item *makersv20261113.SwitchConfig) *SwitchConfigModel {
	if item == nil {
		return nil
	}
	return &SwitchConfigModel{
		Switch: types.StringPointerValue(item.Switch),
	}
}

func buildZoneCustomDomainId(zoneId, projectId, domain string) string {
	return zoneId + tccommon.FILED_SP + projectId + tccommon.FILED_SP + domain
}

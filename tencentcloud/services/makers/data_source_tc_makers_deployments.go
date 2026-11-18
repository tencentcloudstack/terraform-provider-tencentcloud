package makers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	makersv20261113 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &MakersDeploymentsDataSource{}

// NewMakersDeploymentsDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewMakersDeploymentsDataSource() datasource.DataSource {
	return &MakersDeploymentsDataSource{}
}

// MakersDeploymentsDataSource implements datasource.DataSource for
// the tencentcloud_makers_deployments data source.
type MakersDeploymentsDataSource struct {
	fw.DataSourceWithConfigure
}

// MakersDeploymentsDataSourceModel maps the schema attributes.
type MakersDeploymentsDataSourceModel struct {
	ZoneId        types.String `tfsdk:"zone_id"`
	ProjectId     types.String `tfsdk:"project_id"`
	DeploymentIds types.List   `tfsdk:"deployment_ids"`
	Filters       types.List   `tfsdk:"filters"`
	OrderBy       types.String `tfsdk:"order_by"`
	OrderType     types.String `tfsdk:"order_type"`
	StartTime     types.String `tfsdk:"start_time"`
	EndTime       types.String `tfsdk:"end_time"`
	Limit         types.Int64  `tfsdk:"limit"`
	Offset        types.Int64  `tfsdk:"offset"`
	TotalCount    types.Int64  `tfsdk:"total_count"`
	Deployments   types.List   `tfsdk:"deployments"`
}

// MakersDeploymentsFilterModel is the element type of filters.
type MakersDeploymentsFilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

// MakersDeploymentItemModel is the element type of deployments.
type MakersDeploymentItemModel struct {
	DeploymentId   types.String `tfsdk:"deployment_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	Status         types.String `tfsdk:"status"`
	OriginalStatus types.String `tfsdk:"original_status"`
	Code           types.Int64  `tfsdk:"code"`
	BuildCost      types.String `tfsdk:"build_cost"`
	ViaMeta        types.String `tfsdk:"via_meta"`
	Env            types.String `tfsdk:"env"`
	RepoBranch     types.String `tfsdk:"repo_branch"`
	RepoCommitHash types.String `tfsdk:"repo_commit_hash"`
	RepoCommitMsg  types.String `tfsdk:"repo_commit_msg"`
	PreviewUrl     types.String `tfsdk:"preview_url"`
	ProjectUrl     types.String `tfsdk:"project_url"`
	UsedInProd     types.Bool   `tfsdk:"used_in_prod"`
	CoverUrl       types.String `tfsdk:"cover_url"`
	MetaData       types.String `tfsdk:"meta_data"`
	DplRepoUrl     types.String `tfsdk:"dpl_repo_url"`
	CreatedOn      types.String `tfsdk:"created_on"`
	ModifiedOn     types.String `tfsdk:"modified_on"`
}

// MakersDeploymentItemAttrTypes is the object type of the deployments elements.
var MakersDeploymentItemAttrTypes = map[string]attr.Type{
	"deployment_id":    types.StringType,
	"project_id":       types.StringType,
	"status":           types.StringType,
	"original_status":  types.StringType,
	"code":             types.Int64Type,
	"build_cost":       types.StringType,
	"via_meta":         types.StringType,
	"env":              types.StringType,
	"repo_branch":      types.StringType,
	"repo_commit_hash": types.StringType,
	"repo_commit_msg":  types.StringType,
	"preview_url":      types.StringType,
	"project_url":      types.StringType,
	"used_in_prod":     types.BoolType,
	"cover_url":        types.StringType,
	"meta_data":        types.StringType,
	"dpl_repo_url":     types.StringType,
	"created_on":       types.StringType,
	"modified_on":      types.StringType,
}

func (d *MakersDeploymentsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_makers_deployments"
}

func (d *MakersDeploymentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of Makers (TEO Makers site) deployment records.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs ----
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Makers site ID.",
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "Project ID.",
			},
			"deployment_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of deployment IDs to query.",
			},
			"filters": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Filter conditions. Supported filter names: ProjectId, DeploymentId, Name, Status, StartTime, EndTime, RepoBranch.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Filter field. Valid values: ProjectId, DeploymentId, Name, Status, StartTime, EndTime, RepoBranch.",
						},
						"values": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "Filter values.",
						},
					},
				},
			},
			"order_by": schema.StringAttribute{
				Optional:    true,
				Description: "Sort field. Valid values: CreatedOn, ModifiedOn.",
			},
			"order_type": schema.StringAttribute{
				Optional:    true,
				Description: "Sort type. Valid values: Asc, Desc.",
			},
			"start_time": schema.StringAttribute{
				Optional:    true,
				Description: "Start time.",
			},
			"end_time": schema.StringAttribute{
				Optional:    true,
				Description: "End time.",
			},
			"limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of items per page, maximum 100.",
			},
			"offset": schema.Int64Attribute{
				Optional:    true,
				Description: "Pagination offset.",
			},
			// ---- Query outputs: Computed ----
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of deployment records.",
			},
			"deployments": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Deployment list.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"deployment_id": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment ID.",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "Project ID.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment status. Valid values: Success, Failed, Process, Invalid, Pending, Timeout, Cancelled.",
						},
						"original_status": schema.StringAttribute{
							Computed:    true,
							Description: "Original status of the deployment before it was removed.",
						},
						"code": schema.Int64Attribute{
							Computed:    true,
							Description: "Build exit code, 0 means success.",
						},
						"build_cost": schema.StringAttribute{
							Computed:    true,
							Description: "Build cost time, in milliseconds.",
						},
						"via_meta": schema.StringAttribute{
							Computed:    true,
							Description: "Trigger source. Valid values: Github, ReDeploy, Api, Coding, Gitee, Upload, Gitlab.",
						},
						"env": schema.StringAttribute{
							Computed:    true,
							Description: "Environment slug.",
						},
						"repo_branch": schema.StringAttribute{
							Computed:    true,
							Description: "Repository branch.",
						},
						"repo_commit_hash": schema.StringAttribute{
							Computed:    true,
							Description: "Commit hash.",
						},
						"repo_commit_msg": schema.StringAttribute{
							Computed:    true,
							Description: "Commit message.",
						},
						"preview_url": schema.StringAttribute{
							Computed:    true,
							Description: "Access URL of this deployment.",
						},
						"project_url": schema.StringAttribute{
							Computed:    true,
							Description: "Default domain URL of the project.",
						},
						"used_in_prod": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this is the current production version.",
						},
						"cover_url": schema.StringAttribute{
							Computed:    true,
							Description: "Cover image URL.",
						},
						"meta_data": schema.StringAttribute{
							Computed:    true,
							Description: "Route configuration, JSON string.",
						},
						"dpl_repo_url": schema.StringAttribute{
							Computed:    true,
							Description: "Repository URL used by this deployment.",
						},
						"created_on": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time.",
						},
						"modified_on": schema.StringAttribute{
							Computed:    true,
							Description: "Modification time.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *MakersDeploymentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MakersDeploymentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := makersv20261113.NewDescribeDeploymentsRequest()
	request.ZoneId = config.ZoneId.ValueStringPointer()
	request.ProjectId = config.ProjectId.ValueStringPointer()

	if !config.DeploymentIds.IsNull() && !config.DeploymentIds.IsUnknown() {
		var deploymentIds []types.String
		resp.Diagnostics.Append(config.DeploymentIds.ElementsAs(ctx, &deploymentIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ids := make([]*string, 0, len(deploymentIds))
		for _, v := range deploymentIds {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			val := v.ValueString()
			ids = append(ids, &val)
		}
		request.DeploymentIds = ids
	}

	if !config.Filters.IsNull() && !config.Filters.IsUnknown() {
		var filterModels []MakersDeploymentsFilterModel
		resp.Diagnostics.Append(config.Filters.ElementsAs(ctx, &filterModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		filters := make([]*makersv20261113.DeploymentFilter, 0, len(filterModels))
		for _, fm := range filterModels {
			filter := &makersv20261113.DeploymentFilter{
				Name: fm.Name.ValueStringPointer(),
			}
			var values []types.String
			resp.Diagnostics.Append(fm.Values.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			filterValues := make([]*string, 0, len(values))
			for _, v := range values {
				if v.IsNull() || v.IsUnknown() {
					continue
				}
				val := v.ValueString()
				filterValues = append(filterValues, &val)
			}
			filter.Values = filterValues
			filters = append(filters, filter)
		}
		request.Filters = filters
	}

	if !config.OrderBy.IsNull() && !config.OrderBy.IsUnknown() {
		request.OrderBy = config.OrderBy.ValueStringPointer()
	}
	if !config.OrderType.IsNull() && !config.OrderType.IsUnknown() {
		request.OrderType = config.OrderType.ValueStringPointer()
	}
	if !config.StartTime.IsNull() && !config.StartTime.IsUnknown() {
		request.StartTime = config.StartTime.ValueStringPointer()
	}
	if !config.EndTime.IsNull() && !config.EndTime.IsUnknown() {
		request.EndTime = config.EndTime.ValueStringPointer()
	}
	if !config.Limit.IsNull() && !config.Limit.IsUnknown() {
		request.Limit = config.Limit.ValueInt64Pointer()
	}
	if !config.Offset.IsNull() && !config.Offset.IsUnknown() {
		request.Offset = config.Offset.ValueInt64Pointer()
	}

	response, err := d.Client().UseMakersClient().DescribeDeploymentsWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying Makers deployments",
			fmt.Sprintf("DescribeDeployments failed: %v", err),
		)
		return
	}

	// Two different "empty" outcomes must stay distinguishable:
	//   - nil response -> API failure -> error
	//   - len == 0     -> genuinely nothing matched -> empty list, normal return
	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error querying Makers deployments",
			"DescribeDeployments returned an empty response.",
		)
		return
	}

	config.TotalCount = types.Int64PointerValue(response.Response.TotalCount)

	deployments := response.Response.Deployments
	items := make([]MakersDeploymentItemModel, 0, len(deployments))
	for _, item := range deployments {
		items = append(items, MakersDeploymentItemModel{
			DeploymentId:   types.StringPointerValue(item.DeploymentId),
			ProjectId:      types.StringPointerValue(item.ProjectId),
			Status:         types.StringPointerValue(item.Status),
			OriginalStatus: types.StringPointerValue(item.OriginalStatus),
			Code:           types.Int64PointerValue(item.Code),
			BuildCost:      types.StringPointerValue(item.BuildCost),
			ViaMeta:        types.StringPointerValue(item.ViaMeta),
			Env:            types.StringPointerValue(item.Env),
			RepoBranch:     types.StringPointerValue(item.RepoBranch),
			RepoCommitHash: types.StringPointerValue(item.RepoCommitHash),
			RepoCommitMsg:  types.StringPointerValue(item.RepoCommitMsg),
			PreviewUrl:     types.StringPointerValue(item.PreviewUrl),
			ProjectUrl:     types.StringPointerValue(item.ProjectUrl),
			UsedInProd:     types.BoolPointerValue(item.UsedInProd),
			CoverUrl:       types.StringPointerValue(item.CoverUrl),
			MetaData:       types.StringPointerValue(item.MetaData),
			DplRepoUrl:     types.StringPointerValue(item.DplRepoUrl),
			CreatedOn:      types.StringPointerValue(item.CreatedOn),
			ModifiedOn:     types.StringPointerValue(item.ModifiedOn),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: MakersDeploymentItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Deployments = list

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

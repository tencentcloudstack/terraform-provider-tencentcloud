package makers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	makers "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &MakersDeploymentLogDataSource{}

// NewMakersDeploymentLogDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewMakersDeploymentLogDataSource() datasource.DataSource {
	return &MakersDeploymentLogDataSource{}
}

// MakersDeploymentLogDataSource implements datasource.DataSource for
// the tencentcloud_makers_deployment_log data source.
type MakersDeploymentLogDataSource struct {
	fw.DataSourceWithConfigure
}

// MakersDeploymentLogDataSourceModel maps the schema attributes.
type MakersDeploymentLogDataSourceModel struct {
	ZoneId       types.String `tfsdk:"zone_id"`
	ProjectId    types.String `tfsdk:"project_id"`
	DeploymentId types.String `tfsdk:"deployment_id"`
	LogUrl       types.String `tfsdk:"log_url"`
}

func (d *MakersDeploymentLogDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_makers_deployment_log"
}

func (d *MakersDeploymentLogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the build log of a Makers deployment.",
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
			"deployment_id": schema.StringAttribute{
				Required:    true,
				Description: "Deployment ID.",
			},
			// ---- Query outputs: Computed ----
			"log_url": schema.StringAttribute{
				Computed:    true,
				Description: "Build log URL. This field may return null, which means no valid value is available.",
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *MakersDeploymentLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MakersDeploymentLogDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := makers.NewDescribeDeploymentLogRequest()
	request.ZoneId = config.ZoneId.ValueStringPointer()
	request.ProjectId = config.ProjectId.ValueStringPointer()
	request.DeploymentId = config.DeploymentId.ValueStringPointer()

	response, err := d.Client().UseMakersClient().DescribeDeploymentLogWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying Makers deployment log",
			fmt.Sprintf("DescribeDeploymentLog failed: %v", err),
		)
		return
	}
	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error querying Makers deployment log",
			"DescribeDeploymentLog returned an empty response.",
		)
		return
	}

	config.LogUrl = types.StringPointerValue(response.Response.LogUrl)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

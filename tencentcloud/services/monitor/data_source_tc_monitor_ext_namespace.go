package monitor

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &MonitorExtNamespaceDataSource{}

// NewMonitorExtNamespaceDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewMonitorExtNamespaceDataSource() datasource.DataSource {
	return &MonitorExtNamespaceDataSource{}
}

// MonitorExtNamespaceDataSource implements datasource.DataSource for
// the tencentcloud_monitor_ext_namespace data source.
type MonitorExtNamespaceDataSource struct {
	fw.DataSourceWithConfigure
}

// MonitorExtNamespaceDataSourceModel maps the schema attributes.
type MonitorExtNamespaceDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	ExtNamespaceList types.Set    `tfsdk:"ext_namespace_list"`
}

func (d *MonitorExtNamespaceDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_monitor_ext_namespace"
}

func (d *MonitorExtNamespaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query Monitor ext namespaces.",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random
			// token generated at read time (see the end of Read).
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query outputs: Computed (DescribeExtNamespace has no request parameters) ----
			"ext_namespace_list": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of ext namespaces.",
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *MonitorExtNamespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MonitorExtNamespaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	service := NewMonitorService(d.Client())
	names, err := service.DescribeExtNamespace(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying Monitor ext namespaces",
			fmt.Sprintf("DescribeExtNamespace failed: %v", err),
		)
		return
	}

	// An empty result is a genuinely empty list, not a failure.
	extNamespaceItems := make([]attr.Value, 0, len(names))
	for _, name := range names {
		extNamespaceItems = append(extNamespaceItems, types.StringPointerValue(name))
	}

	extNamespaceList, diags := types.SetValue(types.StringType, extNamespaceItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.ExtNamespaceList = extNamespaceList

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

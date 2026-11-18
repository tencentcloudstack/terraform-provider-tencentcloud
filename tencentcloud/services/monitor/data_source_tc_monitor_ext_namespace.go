package monitor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
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
	ExtNamespaceList types.Set `tfsdk:"ext_namespace_list"`
}

func (d *MonitorExtNamespaceDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_monitor_ext_namespace"
}

func (d *MonitorExtNamespaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query Monitor ext namespaces.",
		Attributes: map[string]schema.Attribute{
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

	svc := NewMonitorService(d.Client())
	names, err := svc.DescribeExtNamespace(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying Monitor ext namespaces",
			"DescribeExtNamespace failed: "+err.Error(),
		)
		return
	}
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

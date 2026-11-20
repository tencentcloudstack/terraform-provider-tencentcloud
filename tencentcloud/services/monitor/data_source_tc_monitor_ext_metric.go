package monitor

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	monitorv20230616 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &MonitorExtMetricDataSource{}

// NewMonitorExtMetricDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewMonitorExtMetricDataSource() datasource.DataSource {
	return &MonitorExtMetricDataSource{}
}

// MonitorExtMetricDataSource implements datasource.DataSource for
// the tencentcloud_monitor_ext_metric data source.
type MonitorExtMetricDataSource struct {
	fw.DataSourceWithConfigure
}

// MonitorExtMetricDataSourceModel maps the schema attributes.
type MonitorExtMetricDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	ExtNamespace  types.String `tfsdk:"ext_namespace"`
	ExtMetricList types.List   `tfsdk:"ext_metric_list"`
}

// MonitorExtMetricItemModel is the element type of ext_metric_list.
type MonitorExtMetricItemModel struct {
	MetricName    types.String `tfsdk:"metric_name"`
	MetricCName   types.String `tfsdk:"metric_c_name"`
	CNMeaning     types.String `tfsdk:"cn_meaning"`
	EnMeaning     types.String `tfsdk:"en_meaning"`
	Unit          types.String `tfsdk:"unit"`
	DimensionFlag types.Bool   `tfsdk:"dimension_flag"`
}

// MonitorExtMetricItemAttrTypes is the object type of the ext_metric_list elements.
var MonitorExtMetricItemAttrTypes = map[string]attr.Type{
	"metric_name":    types.StringType,
	"metric_c_name":  types.StringType,
	"cn_meaning":     types.StringType,
	"en_meaning":     types.StringType,
	"unit":           types.StringType,
	"dimension_flag": types.BoolType,
}

func (d *MonitorExtMetricDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_monitor_ext_metric"
}

func (d *MonitorExtMetricDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query Cloud Monitor external metrics of an external namespace.",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random
			// token generated at read time (see the end of Read).
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query inputs ----
			"ext_namespace": schema.StringAttribute{
				Required:    true,
				Description: "External namespace to query.",
			},
			// ---- Query outputs: Computed ----
			"ext_metric_list": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of external metrics under the namespace.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"metric_name": schema.StringAttribute{
							Computed:    true,
							Description: "Metric name. This field may return null, which means no valid value is available.",
						},
						"metric_c_name": schema.StringAttribute{
							Computed:    true,
							Description: "Chinese metric name. This field may return null, which means no valid value is available.",
						},
						"cn_meaning": schema.StringAttribute{
							Computed:    true,
							Description: "Chinese meaning of the metric. This field may return null, which means no valid value is available.",
						},
						"en_meaning": schema.StringAttribute{
							Computed:    true,
							Description: "English meaning of the metric. This field may return null, which means no valid value is available.",
						},
						"unit": schema.StringAttribute{
							Computed:    true,
							Description: "Metric unit. This field may return null, which means no valid value is available.",
						},
						"dimension_flag": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether external dimensions are configured.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *MonitorExtMetricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MonitorExtMetricDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := monitorv20230616.NewDescribeExtMetricRequest()
	request.ExtNamespace = config.ExtNamespace.ValueStringPointer()

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	service := NewMonitorService(d.Client())
	response, err := service.DescribeExtMetric(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying monitor external metrics",
			fmt.Sprintf("DescribeExtMetric failed: %v", err),
		)
		return
	}

	// Two different "empty" outcomes must stay distinguishable:
	//   - nil response  -> API failure -> error
	//   - len == 0      -> genuinely nothing matched -> empty list, normal return
	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error querying monitor external metrics",
			"DescribeExtMetric returned an empty response.",
		)
		return
	}

	items := make([]MonitorExtMetricItemModel, 0, len(response.Response.ExtMetricList))
	for _, item := range response.Response.ExtMetricList {
		if item == nil {
			continue
		}
		items = append(items, MonitorExtMetricItemModel{
			MetricName:    types.StringPointerValue(item.MetricName),
			MetricCName:   types.StringPointerValue(item.MetricCName),
			CNMeaning:     types.StringPointerValue(item.CNMeaning),
			EnMeaning:     types.StringPointerValue(item.EnMeaning),
			Unit:          types.StringPointerValue(item.Unit),
			DimensionFlag: types.BoolPointerValue(item.DimensionFlag),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: MonitorExtMetricItemAttrTypes,
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.ExtMetricList = list

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

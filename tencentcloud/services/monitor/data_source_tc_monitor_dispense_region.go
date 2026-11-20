package monitor

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &MonitorDispenseRegionDataSource{}

// NewMonitorDispenseRegionDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewMonitorDispenseRegionDataSource() datasource.DataSource {
	return &MonitorDispenseRegionDataSource{}
}

// MonitorDispenseRegionDataSource implements datasource.DataSource for
// the tencentcloud_monitor_dispense_region data source.
type MonitorDispenseRegionDataSource struct {
	fw.DataSourceWithConfigure
}

// MonitorDispenseRegionDataSourceModel maps the schema attributes.
type MonitorDispenseRegionDataSourceModel struct {
	RegionList types.List `tfsdk:"region_list"`
}

// MonitorDispenseRegionItemModel is the element type of region_list.
type MonitorDispenseRegionItemModel struct {
	Region       types.String `tfsdk:"region"`
	RegionCnName types.String `tfsdk:"region_cn_name"`
	RegionEnName types.String `tfsdk:"region_en_name"`
	RuleNumber   types.Int64  `tfsdk:"rule_number"`
}

// MonitorDispenseRegionItemAttrTypes is the object type of the
// region_list elements.
var MonitorDispenseRegionItemAttrTypes = map[string]attr.Type{
	"region":         types.StringType,
	"region_cn_name": types.StringType,
	"region_en_name": types.StringType,
	"rule_number":    types.Int64Type,
}

func (d *MonitorDispenseRegionDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_monitor_dispense_region"
}

func (d *MonitorDispenseRegionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the dispense (forwarding) deployment regions of the Monitor product.",
		Attributes: map[string]schema.Attribute{
			"region_list": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of dispense deployment regions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "Region abbreviation, e.g. `ap-guangzhou`.",
						},
						"region_cn_name": schema.StringAttribute{
							Computed:    true,
							Description: "Chinese name of the region.",
						},
						"region_en_name": schema.StringAttribute{
							Computed:    true,
							Description: "English name of the region.",
						},
						"rule_number": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of rules deployed in this region.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *MonitorDispenseRegionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MonitorDispenseRegionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewMonitorService(d.Client())
	regions, err := service.DescribeDispenseRegions(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying Monitor dispense regions",
			fmt.Sprintf("DescribeDispenseRegions failed: %v", err),
		)
		return
	}

	// A nil response is reported as an error by the service layer above;
	// an empty region list is a genuinely empty result, not a failure.
	items := make([]MonitorDispenseRegionItemModel, 0, len(regions))
	for _, item := range regions {
		items = append(items, MonitorDispenseRegionItemModel{
			Region:       types.StringPointerValue(item.Region),
			RegionCnName: types.StringPointerValue(item.RegionCnName),
			RegionEnName: types.StringPointerValue(item.RegionEnName),
			RuleNumber:   types.Int64PointerValue(item.RuleNumber),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: MonitorDispenseRegionItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.RegionList = list

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

package tdmysql

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &TdmysqlSaleInfoDataSource{}

// NewTdmysqlSaleInfoDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTdmysqlSaleInfoDataSource() datasource.DataSource {
	return &TdmysqlSaleInfoDataSource{}
}

// TdmysqlSaleInfoDataSource implements datasource.DataSource for
// tencentcloud_tdmysql_sale_info.
type TdmysqlSaleInfoDataSource struct {
	fw.DataSourceWithConfigure
}

// TdmysqlSaleInfoDataSourceModel maps the schema attributes.
type TdmysqlSaleInfoDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	SrcRegion    types.String `tfsdk:"src_region"`
	InstanceId   types.String `tfsdk:"instance_id"`
	InstanceType types.String `tfsdk:"instance_type"`
	RegionList   types.List   `tfsdk:"region_list"`
}

// TdmysqlSaleRegionItemModel is the element type of region_list.
type TdmysqlSaleRegionItemModel struct {
	Region                types.String `tfsdk:"region"`
	ZoneList              types.List   `tfsdk:"zone_list"`
	RegionName            types.String `tfsdk:"region_name"`
	Enable                types.Int64  `tfsdk:"enable"`
	AvailableZoneNum      types.Int64  `tfsdk:"available_zone_num"`
	IsSupportedLogReplica types.Bool   `tfsdk:"is_supported_log_replica"`
	ZoneGroup             types.List   `tfsdk:"zone_group"`
	IsSupportServerless   types.Bool   `tfsdk:"is_support_serverless"`
}

// TdmysqlSaleRegionItemAttrTypes is the object type of the region_list elements.
var TdmysqlSaleRegionItemAttrTypes = map[string]attr.Type{
	"region":                   types.StringType,
	"zone_list":                types.ListType{ElemType: types.ObjectType{AttrTypes: TdmysqlSaleZoneItemAttrTypes}},
	"region_name":              types.StringType,
	"enable":                   types.Int64Type,
	"available_zone_num":       types.Int64Type,
	"is_supported_log_replica": types.BoolType,
	"zone_group":               types.ListType{ElemType: types.ObjectType{AttrTypes: TdmysqlSaleZoneGroupItemAttrTypes}},
	"is_support_serverless":    types.BoolType,
}

// TdmysqlSaleZoneItemModel is the element type of zone_list.
type TdmysqlSaleZoneItemModel struct {
	Zone                types.String `tfsdk:"zone"`
	ZoneName            types.String `tfsdk:"zone_name"`
	Enable              types.Int64  `tfsdk:"enable"`
	IsDefaultMaster     types.Int64  `tfsdk:"is_default_master"`
	AvailableDiskTypes  types.List   `tfsdk:"available_disk_types"`
	SupportTypes        types.List   `tfsdk:"support_types"`
	IsSupportServerless types.Bool   `tfsdk:"is_support_serverless"`
}

// TdmysqlSaleZoneItemAttrTypes is the object type of the zone_list elements.
var TdmysqlSaleZoneItemAttrTypes = map[string]attr.Type{
	"zone":                  types.StringType,
	"zone_name":             types.StringType,
	"enable":                types.Int64Type,
	"is_default_master":     types.Int64Type,
	"available_disk_types":  types.ListType{ElemType: types.StringType},
	"support_types":         types.ListType{ElemType: types.StringType},
	"is_support_serverless": types.BoolType,
}

// TdmysqlSaleZoneGroupItemModel is the element type of zone_group.
type TdmysqlSaleZoneGroupItemModel struct {
	ZoneNum             types.Int64 `tfsdk:"zone_num"`
	Zones               types.List  `tfsdk:"zones"`
	SupportTypes        types.List  `tfsdk:"support_types"`
	AvailableDiskTypes  types.List  `tfsdk:"available_disk_types"`
	IsSupportServerless types.Bool  `tfsdk:"is_support_serverless"`
}

// TdmysqlSaleZoneGroupItemAttrTypes is the object type of the zone_group elements.
var TdmysqlSaleZoneGroupItemAttrTypes = map[string]attr.Type{
	"zone_num":              types.Int64Type,
	"zones":                 types.ListType{ElemType: types.StringType},
	"support_types":         types.ListType{ElemType: types.StringType},
	"available_disk_types":  types.ListType{ElemType: types.StringType},
	"is_support_serverless": types.BoolType,
}

func (d *TdmysqlSaleInfoDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_tdmysql_sale_info"
}

func (d *TdmysqlSaleInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the TDSQL-C for MySQL sale information of regions and zones.",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random token
			// generated at read time (see the end of Read).
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query inputs ----
			"src_region": schema.StringAttribute{
				Optional:    true,
				Description: "Disaster recovery primary instance region.",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "Instance ID. Passing this parameter returns the sale information of the region and zones where this instance is located.",
			},
			"instance_type": schema.StringAttribute{
				Optional:    true,
				Description: "Specify the sale information of regions that support a certain instance type. Enumerated values: serverless: return all regions that support serverless instances. Default value: none. Currently only serverless is supported; passing other values or leaving it empty returns the sale information of all regions.",
			},
			// ---- Query outputs: Computed ----
			"region_list": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Sale region information.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "Region string in English.",
						},
						"zone_list": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Sale zone list.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"zone": schema.StringAttribute{
										Computed:    true,
										Description: "Zone string in English.",
									},
									"zone_name": schema.StringAttribute{
										Computed:    true,
										Description: "Zone string in Chinese.",
									},
									"enable": schema.Int64Attribute{
										Computed:    true,
										Description: "Whether for sale. 1: for sale; 0: not for sale.",
									},
									"is_default_master": schema.Int64Attribute{
										Computed:    true,
										Description: "Whether it is the default primary zone. 0: no; 1: yes. Note: this field may return null, indicating that no valid value is obtained.",
									},
									"available_disk_types": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "Optional disk types of the zone. Note: this field may return null, indicating that no valid value is obtained.",
									},
									"support_types": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "Whether it is an exclusive zone.",
									},
									"is_support_serverless": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether serverless is supported.",
									},
								},
							},
						},
						"region_name": schema.StringAttribute{
							Computed:    true,
							Description: "Region string in Chinese.",
						},
						"enable": schema.Int64Attribute{
							Computed:    true,
							Description: "Whether for sale. 1: for sale; 0: not for sale.",
						},
						"available_zone_num": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of optional multi-availability zones. Note: this field may return null, indicating that no valid value is obtained.",
						},
						"is_supported_log_replica": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether log replica is allowed. Note: this field may return null, indicating that no valid value is obtained.",
						},
						"zone_group": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Zone group combinations.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"zone_num": schema.Int64Attribute{
										Computed:    true,
										Description: "Number of zones.",
									},
									"zones": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "Zone combination.",
									},
									"support_types": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "Supported types.",
									},
									"available_disk_types": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "Supported disk types. Enumerated values: CLOUD_TCS: local disk; CLOUD_HSSD: enhanced cloud disk. Note: this field may return null, indicating that no valid value is obtained.",
									},
									"is_support_serverless": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether serverless is supported.",
									},
								},
							},
						},
						"is_support_serverless": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether serverless is supported.",
						},
					},
				},
			},
		},
	}
}

// flattenTdmysqlSaleZonesInfo maps one DescribeSaleZonesInfo record to its
// state model.
func flattenTdmysqlSaleZonesInfo(ctx context.Context, zone *tdmysqlv20211122.DescribeSaleZonesInfo, diags *diag.Diagnostics) TdmysqlSaleZoneItemModel {
	item := TdmysqlSaleZoneItemModel{
		Zone:                types.StringPointerValue(zone.Zone),
		ZoneName:            types.StringPointerValue(zone.ZoneName),
		Enable:              fwHelper.Int64ValueOrNull(zone.Enable),
		IsDefaultMaster:     fwHelper.Int64ValueOrNull(zone.IsDefaultMaster),
		IsSupportServerless: fwHelper.BoolValueOrNull(zone.IsSupportServerless),
	}

	item.AvailableDiskTypes = fwHelper.ListFromPointers[string](ctx, types.StringType, zone.AvailableDiskTypes, diags)
	item.SupportTypes = fwHelper.ListFromPointers[string](ctx, types.StringType, zone.SupportTypes, diags)

	return item
}

// flattenTdmysqlSaleZonesGroup maps one DescribeSaleZonesGroup record to its
// state model.
func flattenTdmysqlSaleZonesGroup(ctx context.Context, group *tdmysqlv20211122.DescribeSaleZonesGroup, diags *diag.Diagnostics) TdmysqlSaleZoneGroupItemModel {
	item := TdmysqlSaleZoneGroupItemModel{
		ZoneNum:             fwHelper.Int64ValueOrNull(group.ZoneNum),
		IsSupportServerless: fwHelper.BoolValueOrNull(group.IsSupportServerless),
	}

	item.Zones = fwHelper.ListFromPointers[string](ctx, types.StringType, group.Zones, diags)
	item.SupportTypes = fwHelper.ListFromPointers[string](ctx, types.StringType, group.SupportTypes, diags)
	item.AvailableDiskTypes = fwHelper.ListFromPointers[string](ctx, types.StringType, group.AvailableDiskTypes, diags)

	return item
}

// flattenTdmysqlSaleRegionInfo maps one DescribeSaleRegionInfo record to its
// state model.
func flattenTdmysqlSaleRegionInfo(ctx context.Context, region *tdmysqlv20211122.DescribeSaleRegionInfo, diags *diag.Diagnostics) TdmysqlSaleRegionItemModel {
	item := TdmysqlSaleRegionItemModel{
		Region:                types.StringPointerValue(region.Region),
		RegionName:            types.StringPointerValue(region.RegionName),
		Enable:                fwHelper.Int64ValueOrNull(region.Enable),
		AvailableZoneNum:      fwHelper.Int64ValueOrNull(region.AvailableZoneNum),
		IsSupportedLogReplica: fwHelper.BoolValueOrNull(region.IsSupportedLogReplica),
		IsSupportServerless:   fwHelper.BoolValueOrNull(region.IsSupportServerless),
	}

	zoneItems := make([]TdmysqlSaleZoneItemModel, 0, len(region.ZoneList))
	for _, zone := range region.ZoneList {
		if zone == nil {
			continue
		}

		zoneItems = append(zoneItems, flattenTdmysqlSaleZonesInfo(ctx, zone, diags))
	}

	item.ZoneList, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlSaleZoneItemAttrTypes}, zoneItems)

	groupItems := make([]TdmysqlSaleZoneGroupItemModel, 0, len(region.ZoneGroup))
	for _, group := range region.ZoneGroup {
		if group == nil {
			continue
		}

		groupItems = append(groupItems, flattenTdmysqlSaleZonesGroup(ctx, group, diags))
	}

	item.ZoneGroup, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlSaleZoneGroupItemAttrTypes}, groupItems)

	return item
}

// Read queries the cloud API and writes the result into resp.State.
func (d *TdmysqlSaleInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TdmysqlSaleInfoDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeSaleInfoRequest := tdmysqlv20211122.NewDescribeSaleInfoRequest()
	if !config.SrcRegion.IsNull() && !config.SrcRegion.IsUnknown() {
		describeSaleInfoRequest.SrcRegion = config.SrcRegion.ValueStringPointer()
	}
	if !config.InstanceId.IsNull() && !config.InstanceId.IsUnknown() {
		describeSaleInfoRequest.InstanceId = config.InstanceId.ValueStringPointer()
	}
	if !config.InstanceType.IsNull() && !config.InstanceType.IsUnknown() {
		describeSaleInfoRequest.InstanceType = config.InstanceType.ValueStringPointer()
	}

	service := NewTdmysqlService(d.Client())

	regionList, err := service.DescribeSaleInfo(ctx, describeSaleInfoRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TDSQL-C for MySQL sale info",
			fmt.Sprintf("DescribeSaleInfo failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	regions := make([]TdmysqlSaleRegionItemModel, 0, len(regionList))
	for _, region := range regionList {
		if region == nil {
			continue
		}

		regions = append(regions, flattenTdmysqlSaleRegionInfo(ctx, region, &resp.Diagnostics))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	regionListValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TdmysqlSaleRegionItemAttrTypes}, regions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.RegionList = regionListValue

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

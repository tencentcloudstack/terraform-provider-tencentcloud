package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teosdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &TeoCustomGroupsDataSource{}

// NewTeoCustomGroupsDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTeoCustomGroupsDataSource() datasource.DataSource {
	return &TeoCustomGroupsDataSource{}
}

// TeoCustomGroupsDataSource implements datasource.DataSource for
// the tencentcloud_teo_custom_groups data source.
type TeoCustomGroupsDataSource struct {
	fw.DataSourceWithConfigure
}

// TeoCustomGroupsDataSourceModel maps the schema attributes.
type TeoCustomGroupsDataSourceModel struct {
	ZoneId         types.String `tfsdk:"zone_id"`
	GroupTypes     types.List   `tfsdk:"group_types"`
	GroupTotalCount types.Int64 `tfsdk:"group_total_count"`
	CustomGroups   types.List   `tfsdk:"custom_groups"`
}

// TeoCustomGroupItemModel is the element type of custom_groups.
type TeoCustomGroupItemModel struct {
	GroupId          types.String `tfsdk:"group_id"`
	GroupName        types.String `tfsdk:"group_name"`
	GroupType        types.String `tfsdk:"group_type"`
	GroupItems       types.List   `tfsdk:"group_items"`
	GroupItemTotalCount types.Int64 `tfsdk:"group_item_total_count"`
	RefCount        types.Int64  `tfsdk:"ref_count"`
}

// TeoCustomGroupEntryModel is the element type of group_items.
type TeoCustomGroupEntryModel struct {
	Content  types.String `tfsdk:"content"`
	ItemType types.String `tfsdk:"item_type"`
}

// TeoCustomGroupEntryAttrTypes is the object type of the group_items elements.
var TeoCustomGroupEntryAttrTypes = map[string]attr.Type{
	"content":   types.StringType,
	"item_type": types.StringType,
}

// TeoCustomGroupItemAttrTypes is the object type of the custom_groups elements.
var TeoCustomGroupItemAttrTypes = map[string]attr.Type{
	"group_id":             types.StringType,
	"group_name":           types.StringType,
	"group_type":           types.StringType,
	"group_items":          types.List{ElementType: types.ObjectType{AttrTypes: TeoCustomGroupEntryAttrTypes}},
	"group_item_total_count": types.Int64Type,
	"ref_count":            types.Int64Type,
}

func (d *TeoCustomGroupsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_custom_groups"
}

func (d *TeoCustomGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of EdgeOne (TEO) custom groups.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs ----
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Site ID.",
			},
			"group_types": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of custom group types. If not specified, all types of custom groups are queried. Valid values: `TLSFingerprint`, `UA`, `Referer`.",
			},
			// ---- Query outputs: Computed ----
			"group_total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of custom groups that match the group type conditions.",
			},
			"custom_groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of custom groups that match the group type conditions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique ID of the custom group.",
						},
						"group_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the custom group.",
						},
						"group_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the custom group. Valid values: `TLSFingerprint`, `UA`, `Referer`.",
						},
						"group_items": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of custom group entries.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"content": schema.StringAttribute{
										Computed:    true,
										Description: "Entry content.",
									},
									"item_type": schema.StringAttribute{
										Computed:    true,
										Description: "Entry type. Valid values: `JA3` (only available when the group type is TLSFingerprint), `JA4` (only available when the group type is TLSFingerprint), `UA` (only available when the group type is UA), `Referer` (only available when the group type is Referer).",
									},
								},
							},
						},
						"group_item_total_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Total number of entries in the custom group.",
						},
						"ref_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of times the custom group is referenced.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *TeoCustomGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoCustomGroupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := teosdk.NewDescribeCustomGroupsRequest()
	request.ZoneId = config.ZoneId.ValueStringPointer()

	if !config.GroupTypes.IsNull() && !config.GroupTypes.IsUnknown() {
		var groupTypes []types.String
		resp.Diagnostics.Append(config.GroupTypes.ElementsAs(ctx, &groupTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		groupTypePtrs := make([]*string, 0, len(groupTypes))
		for _, v := range groupTypes {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			val := v.ValueString()
			groupTypePtrs = append(groupTypePtrs, &val)
		}
		request.GroupTypes = groupTypePtrs
	}

	service := NewTeoService(d.Client())
	customGroups, groupTotalCount, err := service.DescribeTeoCustomGroups(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TEO custom groups",
			fmt.Sprintf("DescribeCustomGroups failed: %v", err),
		)
		return
	}

	// A nil response is reported by the service layer above; an empty
	// custom groups list is a genuinely empty result, not a failure.
	config.GroupTotalCount = types.Int64PointerValue(groupTotalCount)

	items := make([]TeoCustomGroupItemModel, 0, len(customGroups))
	for _, group := range customGroups {
		entryModels := make([]TeoCustomGroupEntryModel, 0, len(group.GroupItems))
		for _, entry := range group.GroupItems {
			entryModels = append(entryModels, TeoCustomGroupEntryModel{
				Content:  types.StringPointerValue(entry.Content),
				ItemType: types.StringPointerValue(entry.ItemType),
			})
		}
		entryList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TeoCustomGroupEntryAttrTypes}, entryModels)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		items = append(items, TeoCustomGroupItemModel{
			GroupId:             types.StringPointerValue(group.GroupId),
			GroupName:           types.StringPointerValue(group.GroupName),
			GroupType:           types.StringPointerValue(group.GroupType),
			GroupItems:          entryList,
			GroupItemTotalCount: types.Int64PointerValue(group.GroupItemTotalCount),
			RefCount:            types.Int64PointerValue(group.RefCount),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TeoCustomGroupItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.CustomGroups = list

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

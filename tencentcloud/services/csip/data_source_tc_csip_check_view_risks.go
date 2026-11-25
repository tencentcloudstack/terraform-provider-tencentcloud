package csip

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	csipv20221121 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/csip/v20221121"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &CsipCheckViewRisksDataSource{}

// NewCsipCheckViewRisksDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewCsipCheckViewRisksDataSource() datasource.DataSource {
	return &CsipCheckViewRisksDataSource{}
}

// CsipCheckViewRisksDataSource implements datasource.DataSource for
// tencentcloud_csip_check_view_risks.
type CsipCheckViewRisksDataSource struct {
	fw.DataSourceWithConfigure
}

// CsipCheckViewRisksDataSourceModel maps the schema attributes.
type CsipCheckViewRisksDataSourceModel struct {
	Id                types.String `tfsdk:"id"`
	MemberId          types.List   `tfsdk:"member_id"`
	Filters           types.List   `tfsdk:"filters"`
	Order             types.String `tfsdk:"order"`
	By                types.String `tfsdk:"by"`
	CheckViewRiskList types.List   `tfsdk:"check_view_risk_list"`
	StandardNameList  types.List   `tfsdk:"standard_name_list"`
	AssetTypeList     types.List   `tfsdk:"asset_type_list"`
	ProviderList      types.List   `tfsdk:"provider_list"`
}

// CsipCheckViewRisksFilterModel is the element type of filters.
type CsipCheckViewRisksFilterModel struct {
	Name       types.String `tfsdk:"name"`
	Values     types.List   `tfsdk:"values"`
	ExactMatch types.String `tfsdk:"exact_match"`
}

// CsipCheckViewRiskItemModel is the element type of check_view_risk_list.
type CsipCheckViewRiskItemModel struct {
	RiskRuleId       types.String `tfsdk:"risk_rule_id"`
	RiskTitle        types.String `tfsdk:"risk_title"`
	CheckType        types.String `tfsdk:"check_type"`
	Severity         types.String `tfsdk:"severity"`
	RiskDesc         types.String `tfsdk:"risk_desc"`
	CreateTime       types.String `tfsdk:"create_time"`
	UpdateTime       types.String `tfsdk:"update_time"`
	Provider         types.String `tfsdk:"provider"`
	RiskStatus       types.Int64  `tfsdk:"risk_status"`
	AssetCount       types.Int64  `tfsdk:"asset_count"`
	RiskCount        types.Int64  `tfsdk:"risk_count"`
	AssetType        types.String `tfsdk:"asset_type"`
	EventType        types.String `tfsdk:"event_type"`
	Classify         types.String `tfsdk:"classify"`
	StandardTerms    types.List   `tfsdk:"standard_terms"`
	AssetTypeIconUrl types.String `tfsdk:"asset_type_icon_url"`
}

// CsipCheckViewRiskItemAttrTypes is the object type of the
// check_view_risk_list elements.
var CsipCheckViewRiskItemAttrTypes = map[string]attr.Type{
	"risk_rule_id":        types.StringType,
	"risk_title":          types.StringType,
	"check_type":          types.StringType,
	"severity":            types.StringType,
	"risk_desc":           types.StringType,
	"create_time":         types.StringType,
	"update_time":         types.StringType,
	"provider":            types.StringType,
	"risk_status":         types.Int64Type,
	"asset_count":         types.Int64Type,
	"risk_count":          types.Int64Type,
	"asset_type":          types.StringType,
	"event_type":          types.StringType,
	"classify":            types.StringType,
	"standard_terms":      types.ListType{ElemType: types.ObjectType{AttrTypes: CsipCheckViewRiskStandardTermAttrTypes}},
	"asset_type_icon_url": types.StringType,
}

// CsipCheckViewRiskStandardTermModel is the element type of standard_terms.
type CsipCheckViewRiskStandardTermModel struct {
	Tag   types.String `tfsdk:"tag"`
	Terms types.List   `tfsdk:"terms"`
}

// CsipCheckViewRiskStandardTermAttrTypes is the object type of the
// standard_terms elements.
var CsipCheckViewRiskStandardTermAttrTypes = map[string]attr.Type{
	"tag":   types.StringType,
	"terms": types.ListType{ElemType: types.StringType},
}

// CsipCheckViewStandardNameItemModel is the element type of
// standard_name_list.
type CsipCheckViewStandardNameItemModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// CsipCheckViewStandardNameItemAttrTypes is the object type of the
// standard_name_list elements.
var CsipCheckViewStandardNameItemAttrTypes = map[string]attr.Type{
	"id":   types.Int64Type,
	"name": types.StringType,
}

// CsipCheckViewAttributeOptionItemModel is the element type of
// asset_type_list and provider_list (the shared AttributeOptionSet shape).
type CsipCheckViewAttributeOptionItemModel struct {
	Text  types.String `tfsdk:"text"`
	Value types.String `tfsdk:"value"`
}

// CsipCheckViewAttributeOptionItemAttrTypes is the object type of the
// asset_type_list and provider_list elements.
var CsipCheckViewAttributeOptionItemAttrTypes = map[string]attr.Type{
	"text":  types.StringType,
	"value": types.StringType,
}

func (d *CsipCheckViewRisksDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_csip_check_view_risks"
}

func (d *CsipCheckViewRisksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the cloud asset configuration risks under the check item view.",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random token
			// generated at read time (see the end of Read).
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query inputs ----
			"member_id": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Member IDs of the organization account.",
			},
			"filters": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Filter conditions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Filter name. Valid value: Status (execution result, Values: SUCCESS/FAILED/USER_CANCELED/CHECKING).",
						},
						"values": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "List of filter values.",
						},
						"exact_match": schema.StringAttribute{
							Optional:    true,
							Description: "Whether to match exactly: 1 means exact match; fuzzy match by default.",
						},
					},
				},
			},
			"order": schema.StringAttribute{
				Optional:    true,
				Description: "Sort type.",
			},
			"by": schema.StringAttribute{
				Optional:    true,
				Description: "Sort field.",
			},
			// ---- Query outputs: Computed ----
			"check_view_risk_list": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of risks under the check item view.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"risk_rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "Check item rule ID.",
						},
						"risk_title": schema.StringAttribute{
							Computed:    true,
							Description: "Risk name.",
						},
						"check_type": schema.StringAttribute{
							Computed:    true,
							Description: "Check type.",
						},
						"severity": schema.StringAttribute{
							Computed:    true,
							Description: "Risk level.",
						},
						"risk_desc": schema.StringAttribute{
							Computed:    true,
							Description: "1 risk item exists.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "First discovery time.",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "Risk update time.",
						},
						"provider": schema.StringAttribute{
							Computed:    true,
							Description: "Cloud provider.",
						},
						"risk_status": schema.Int64Attribute{
							Computed:    true,
							Description: "Risk status.",
						},
						"asset_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of affected assets.",
						},
						"risk_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of risks.",
						},
						"asset_type": schema.StringAttribute{
							Computed:    true,
							Description: "Asset type.",
						},
						"event_type": schema.StringAttribute{
							Computed:    true,
							Description: "Event type.",
						},
						"classify": schema.StringAttribute{
							Computed:    true,
							Description: "Disposition classification.",
						},
						"standard_terms": schema.ListNestedAttribute{
							Computed:    true,
							Description: "CSPM standard terms.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"tag": schema.StringAttribute{
										Computed:    true,
										Description: "Tag.",
									},
									"terms": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "Terms.",
									},
								},
							},
						},
						"asset_type_icon_url": schema.StringAttribute{
							Computed:    true,
							Description: "Asset type icon.",
						},
					},
				},
			},
			"standard_name_list": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of CSPM standard tags under the check item view.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "Standard ID.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Standard name.",
						},
					},
				},
			},
			"asset_type_list": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Collection of asset types.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"text": schema.StringAttribute{
							Computed:    true,
							Description: "CVM instance type.",
						},
						"value": schema.StringAttribute{
							Computed:    true,
							Description: "CVM instance name.",
						},
					},
				},
			},
			"provider_list": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Collection of cloud provider types.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"text": schema.StringAttribute{
							Computed:    true,
							Description: "CVM instance type.",
						},
						"value": schema.StringAttribute{
							Computed:    true,
							Description: "CVM instance name.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *CsipCheckViewRisksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CsipCheckViewRisksDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeCheckViewRisksRequest := csipv20221121.NewDescribeCheckViewRisksRequest()

	if !config.MemberId.IsNull() && !config.MemberId.IsUnknown() {
		describeCheckViewRisksRequest.MemberId = fwHelper.PointersFromList[string](ctx, config.MemberId, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !config.Filters.IsNull() && !config.Filters.IsUnknown() {
		var filterModels []CsipCheckViewRisksFilterModel
		resp.Diagnostics.Append(config.Filters.ElementsAs(ctx, &filterModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		filters := make([]*csipv20221121.Filters, 0, len(filterModels))
		for _, fm := range filterModels {
			filter := &csipv20221121.Filters{
				Name: fm.Name.ValueStringPointer(),
			}
			if !fm.Values.IsNull() && !fm.Values.IsUnknown() {
				filter.Values = fwHelper.PointersFromList[string](ctx, fm.Values, &resp.Diagnostics)
				if resp.Diagnostics.HasError() {
					return
				}
			}
			if !fm.ExactMatch.IsNull() && !fm.ExactMatch.IsUnknown() {
				filter.ExactMatch = fm.ExactMatch.ValueStringPointer()
			}
			filters = append(filters, filter)
		}
		describeCheckViewRisksRequest.Filters = filters
	}

	if !config.Order.IsNull() && !config.Order.IsUnknown() {
		describeCheckViewRisksRequest.Order = config.Order.ValueStringPointer()
	}
	if !config.By.IsNull() && !config.By.IsUnknown() {
		describeCheckViewRisksRequest.By = config.By.ValueStringPointer()
	}

	service := NewCsipService(d.Client())

	// The service method pages internally (Offset/Limit); this file never
	// loops by itself and never surfaces the total count.
	checkViewRisks, standardNameList, assetTypeList, providerList, err := service.DescribeCheckViewRisks(ctx, describeCheckViewRisksRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying CSIP check view risks",
			fmt.Sprintf("DescribeCheckViewRisks failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	riskItems := make([]CsipCheckViewRiskItemModel, 0, len(checkViewRisks))
	for _, item := range checkViewRisks {
		if item == nil {
			continue
		}

		standardTerms := types.ListNull(types.ObjectType{AttrTypes: CsipCheckViewRiskStandardTermAttrTypes})
		if item.StandardTerms != nil {
			terms := make([]CsipCheckViewRiskStandardTermModel, 0, len(item.StandardTerms))
			for _, term := range item.StandardTerms {
				if term == nil {
					continue
				}

				termValues := types.ListNull(types.StringType)
				if term.Terms != nil {
					values := make([]attr.Value, 0, len(term.Terms))
					for _, v := range term.Terms {
						values = append(values, types.StringPointerValue(v))
					}
					list, diags := types.ListValue(types.StringType, values)
					resp.Diagnostics.Append(diags...)
					if resp.Diagnostics.HasError() {
						return
					}
					termValues = list
				}

				terms = append(terms, CsipCheckViewRiskStandardTermModel{
					Tag:   types.StringPointerValue(term.Tag),
					Terms: termValues,
				})
			}
			list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CsipCheckViewRiskStandardTermAttrTypes}, terms)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			standardTerms = list
		}

		riskItems = append(riskItems, CsipCheckViewRiskItemModel{
			RiskRuleId:       types.StringPointerValue(item.RiskRuleId),
			RiskTitle:        types.StringPointerValue(item.RiskTitle),
			CheckType:        types.StringPointerValue(item.CheckType),
			Severity:         types.StringPointerValue(item.Severity),
			RiskDesc:         types.StringPointerValue(item.RiskDesc),
			CreateTime:       types.StringPointerValue(item.CreateTime),
			UpdateTime:       types.StringPointerValue(item.UpdateTime),
			Provider:         types.StringPointerValue(item.Provider),
			RiskStatus:       fwHelper.Int64ValueFromUint(item.RiskStatus),
			AssetCount:       fwHelper.Int64ValueFromUint(item.AssetCount),
			RiskCount:        fwHelper.Int64ValueFromUint(item.RiskCount),
			AssetType:        types.StringPointerValue(item.AssetType),
			EventType:        types.StringPointerValue(item.EventType),
			Classify:         types.StringPointerValue(item.Classify),
			StandardTerms:    standardTerms,
			AssetTypeIconUrl: types.StringPointerValue(item.AssetTypeIconURL),
		})
	}

	riskList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CsipCheckViewRiskItemAttrTypes}, riskItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.CheckViewRiskList = riskList

	standardNameItems := make([]CsipCheckViewStandardNameItemModel, 0, len(standardNameList))
	for _, item := range standardNameList {
		if item == nil {
			continue
		}

		standardNameItems = append(standardNameItems, CsipCheckViewStandardNameItemModel{
			Id:   fwHelper.Int64ValueFromUint(item.ID),
			Name: types.StringPointerValue(item.Name),
		})
	}

	standardNameListValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CsipCheckViewStandardNameItemAttrTypes}, standardNameItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.StandardNameList = standardNameListValue

	assetTypeItems := make([]CsipCheckViewAttributeOptionItemModel, 0, len(assetTypeList))
	for _, item := range assetTypeList {
		if item == nil {
			continue
		}

		assetTypeItems = append(assetTypeItems, CsipCheckViewAttributeOptionItemModel{
			Text:  types.StringPointerValue(item.Text),
			Value: types.StringPointerValue(item.Value),
		})
	}

	assetTypeListValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CsipCheckViewAttributeOptionItemAttrTypes}, assetTypeItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.AssetTypeList = assetTypeListValue

	providerItems := make([]CsipCheckViewAttributeOptionItemModel, 0, len(providerList))
	for _, item := range providerList {
		if item == nil {
			continue
		}

		providerItems = append(providerItems, CsipCheckViewAttributeOptionItemModel{
			Text:  types.StringPointerValue(item.Text),
			Value: types.StringPointerValue(item.Value),
		})
	}

	providerListValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CsipCheckViewAttributeOptionItemAttrTypes}, providerItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.ProviderList = providerListValue

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

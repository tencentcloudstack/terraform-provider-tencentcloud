package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &ConfigListAggregateCompliancePacksDataSource{}

// NewConfigListAggregateCompliancePacksDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewConfigListAggregateCompliancePacksDataSource() datasource.DataSource {
	return &ConfigListAggregateCompliancePacksDataSource{}
}

// ConfigListAggregateCompliancePacksDataSource implements datasource.DataSource for
// the tencentcloud_config_list_aggregate_compliance_packs data source.
type ConfigListAggregateCompliancePacksDataSource struct {
	fw.DataSourceWithConfigure
}

// ConfigListAggregateCompliancePacksDataSourceModel maps the schema attributes.
type ConfigListAggregateCompliancePacksDataSourceModel struct {
	AccountGroupId                  types.String `tfsdk:"account_group_id"`
	CompliancePackName              types.String `tfsdk:"compliance_pack_name"`
	RiskLevel                       types.List   `tfsdk:"risk_level"`
	Status                          types.String `tfsdk:"status"`
	ComplianceResult                types.List   `tfsdk:"compliance_result"`
	OrderType                       types.String `tfsdk:"order_type"`
	IncludeCompliancePackRuleResult types.String `tfsdk:"include_compliance_pack_rule_result"`

	Total types.Int64 `tfsdk:"total"`
	Items types.List  `tfsdk:"items"`
}

// ConfigAggregateCompliancePackItemModel is the element type of items.
type ConfigAggregateCompliancePackItemModel struct {
	Status             types.String `tfsdk:"status"`
	CompliancePackId   types.String `tfsdk:"compliance_pack_id"`
	Description        types.String `tfsdk:"description"`
	CompliancePackName types.String `tfsdk:"compliance_pack_name"`
	RiskLevel          types.Int64  `tfsdk:"risk_level"`
	ComplianceResult   types.String `tfsdk:"compliance_result"`
	CreateTime         types.String `tfsdk:"create_time"`
	NoCompliantNames   types.List   `tfsdk:"no_compliant_names"`
	RuleCount          types.Int64  `tfsdk:"rule_count"`
}

// ConfigAggregateCompliancePackItemAttrTypes is the object type of the items elements.
var ConfigAggregateCompliancePackItemAttrTypes = map[string]attr.Type{
	"status":               types.StringType,
	"compliance_pack_id":   types.StringType,
	"description":          types.StringType,
	"compliance_pack_name": types.StringType,
	"risk_level":           types.Int64Type,
	"compliance_result":    types.StringType,
	"create_time":          types.StringType,
	"no_compliant_names":   types.ListType{ElemType: types.StringType},
	"rule_count":           types.Int64Type,
}

func (d *ConfigListAggregateCompliancePacksDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_config_list_aggregate_compliance_packs"
}

func (d *ConfigListAggregateCompliancePacksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of Config aggregate compliance packs by account group.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs ----
			"account_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Account group ID.",
			},
			"compliance_pack_name": schema.StringAttribute{
				Optional:    true,
				Description: "Compliance pack name.",
			},
			"risk_level": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Compliance pack status. Valid values: ACTIVE, NO_ACTIVE.",
			},
			"compliance_result": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Compliance result. Valid values: COMPLIANT, NON_COMPLIANT.",
			},
			"order_type": schema.StringAttribute{
				Optional:    true,
				Description: "Sort type. Valid values: desc (descending), asc (ascending).",
			},
			"include_compliance_pack_rule_result": schema.StringAttribute{
				Optional:    true,
				Description: "Whether to include compliance pack rule result. Valid values: NO (not included). If not set or set to YES, the compliance result data is returned by default.",
			},
			// ---- Query outputs: Computed ----
			"total": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of compliance packs.",
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Compliance pack list.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance pack status.",
						},
						"compliance_pack_id": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance pack ID.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance pack description.",
						},
						"compliance_pack_name": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance pack name.",
						},
						"risk_level": schema.Int64Attribute{
							Computed:    true,
							Description: "Risk level.",
						},
						"compliance_result": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance result.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time.",
						},
						"no_compliant_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "List of non-compliant rule names.",
						},
						"rule_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of rules in the compliance pack.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *ConfigListAggregateCompliancePacksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ConfigListAggregateCompliancePacksDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramMap := make(map[string]interface{})
	if !config.AccountGroupId.IsNull() && !config.AccountGroupId.IsUnknown() {
		paramMap["AccountGroupId"] = config.AccountGroupId.ValueString()
	}
	if !config.CompliancePackName.IsNull() && !config.CompliancePackName.IsUnknown() {
		paramMap["CompliancePackName"] = config.CompliancePackName.ValueString()
	}
	if !config.RiskLevel.IsNull() && !config.RiskLevel.IsUnknown() {
		var riskLevels []types.Int64
		resp.Diagnostics.Append(config.RiskLevel.ElementsAs(ctx, &riskLevels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		levels := make([]*uint64, 0, len(riskLevels))
		for _, v := range riskLevels {
			if u := fwhelper.Uint64PointerFromValue(v); u != nil {
				levels = append(levels, u)
			}
		}
		paramMap["RiskLevel"] = levels
	}
	if !config.Status.IsNull() && !config.Status.IsUnknown() {
		paramMap["Status"] = config.Status.ValueString()
	}
	if !config.ComplianceResult.IsNull() && !config.ComplianceResult.IsUnknown() {
		var results []types.String
		resp.Diagnostics.Append(config.ComplianceResult.ElementsAs(ctx, &results, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		rr := make([]*string, 0, len(results))
		for _, v := range results {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			val := v.ValueString()
			rr = append(rr, &val)
		}
		paramMap["ComplianceResult"] = rr
	}
	if !config.OrderType.IsNull() && !config.OrderType.IsUnknown() {
		paramMap["OrderType"] = config.OrderType.ValueString()
	}
	if !config.IncludeCompliancePackRuleResult.IsNull() && !config.IncludeCompliancePackRuleResult.IsUnknown() {
		paramMap["IncludeCompliancePackRuleResult"] = config.IncludeCompliancePackRuleResult.ValueString()
	}

	service := NewConfigService(d.Client())
	respItems, total, err := service.DescribeConfigListAggregateCompliancePacksByFilter(ctx, paramMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying Config aggregate compliance packs",
			fmt.Sprintf("ListAggregateCompliancePacks failed: %v", err),
		)
		return
	}

	// Two different "empty" outcomes must stay distinguishable:
	//   - nil response error is reported by the service layer above.
	//   - len == 0 is a genuinely empty result -> empty list, normal return.
	items := make([]ConfigAggregateCompliancePackItemModel, 0, len(respItems))
	for _, pack := range respItems {
		noCompliantNames := types.ListNull(types.StringType)
		if pack.NoCompliantNames != nil {
			names := make([]attr.Value, 0, len(pack.NoCompliantNames))
			for _, n := range pack.NoCompliantNames {
				names = append(names, types.StringPointerValue(n))
			}
			list, diags := types.ListValue(types.StringType, names)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			noCompliantNames = list
		}
		items = append(items, ConfigAggregateCompliancePackItemModel{
			Status:             types.StringPointerValue(pack.Status),
			CompliancePackId:   types.StringPointerValue(pack.CompliancePackId),
			Description:        types.StringPointerValue(pack.Description),
			CompliancePackName: types.StringPointerValue(pack.CompliancePackName),
			RiskLevel:          fwhelper.Int64ValueFromUint(pack.RiskLevel),
			ComplianceResult:   types.StringPointerValue(pack.ComplianceResult),
			CreateTime:         types.StringPointerValue(pack.CreateTime),
			NoCompliantNames:   noCompliantNames,
			RuleCount:          fwhelper.Int64ValueFromUint(pack.RuleCount),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: ConfigAggregateCompliancePackItemAttrTypes,
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Items = list
	config.Total = types.Int64Value(int64(total))

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

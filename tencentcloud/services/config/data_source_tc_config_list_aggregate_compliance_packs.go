// Copyright (c) 2018-2025 Tencent. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &ConfigListAggregateCompliancePacksDataSource{}

// NewConfigListAggregateCompliancePacksDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewConfigListAggregateCompliancePacksDataSource() datasource.DataSource {
	return &ConfigListAggregateCompliancePacksDataSource{}
}

// ConfigListAggregateCompliancePacksDataSource implements datasource.DataSource
// for tencentcloud_config_list_aggregate_compliance_packs.
type ConfigListAggregateCompliancePacksDataSource struct {
	fw.DataSourceWithConfigure
}

// ConfigListAggregateCompliancePacksDataSourceModel maps the schema attributes.
// Every schema attribute appears here exactly once.
type ConfigListAggregateCompliancePacksDataSourceModel struct {
	AccountGroupId                   types.String `tfsdk:"account_group_id"`
	CompliancePackName               types.String `tfsdk:"compliance_pack_name"`
	RiskLevel                        types.Set    `tfsdk:"risk_level"`
	Status                           types.String `tfsdk:"status"`
	ComplianceResult                 types.Set    `tfsdk:"compliance_result"`
	OrderType                        types.String `tfsdk:"order_type"`
	IncludeCompliancePackRuleResult  types.String `tfsdk:"include_compliance_pack_rule_result"`
	Total                            types.Int64  `tfsdk:"total"`
	Items                            types.List   `tfsdk:"items"`
}

// configListAggregateCompliancePacksItemModel is the element type of items.
type configListAggregateCompliancePacksItemModel struct {
	Status             types.String `tfsdk:"status"`
	CompliancePackId   types.String `tfsdk:"compliance_pack_id"`
	Description        types.String `tfsdk:"description"`
	CompliancePackName types.String `tfsdk:"compliance_pack_name"`
	RiskLevel          types.Int64  `tfsdk:"risk_level"`
	ComplianceResult   types.String `tfsdk:"compliance_result"`
	CreateTime         types.String `tfsdk:"create_time"`
	NoCompliantNames   types.Set    `tfsdk:"no_compliant_names"`
	RuleCount          types.Int64  `tfsdk:"rule_count"`
}

func (d *ConfigListAggregateCompliancePacksDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_config_list_aggregate_compliance_packs"
}

func (d *ConfigListAggregateCompliancePacksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query detailed information of config aggregate compliance packs by account group.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs: Required / Optional only ----
			"account_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Account group ID.",
			},
			"compliance_pack_name": schema.StringAttribute{
				Optional:    true,
				Description: "Compliance pack name used for filtering.",
			},
			"risk_level": schema.SetAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Risk levels used for filtering. Valid values: 1 (high), 2 (medium), 3 (low).",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Compliance pack status used for filtering. Valid values: ACTIVE, NO_ACTIVE.",
			},
			"compliance_result": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Evaluation results used for filtering. Valid values: COMPLIANT, NON_COMPLIANT.",
			},
			"order_type": schema.StringAttribute{
				Optional:    true,
				Description: "Sort type. Valid values: desc, asc.",
			},
			"include_compliance_pack_rule_result": schema.StringAttribute{
				Optional:    true,
				Description: "Whether to include compliance pack rule result. Valid values: NO. Default: empty. Note that empty or YES both return the compliance result data.",
			},
			// ---- Query outputs: Computed (datasource schema has no plan modifiers) ----
			"total": schema.Int64Attribute{
				Computed:    true,
				Description: "Total count of compliance packs matched by the filter.",
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of compliance packs.",
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
							Description: "Risk level. Valid values: 1 (high), 2 (medium), 3 (low).",
						},
						"compliance_result": schema.StringAttribute{
							Computed:    true,
							Description: "Evaluation result of the compliance pack.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time of the compliance pack.",
						},
						"no_compliant_names": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Names of non-compliant rules.",
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

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Error querying config aggregate compliance packs",
			"Provider client is not configured, please report this bug to the provider developer.",
		)
		return
	}

	paramMap := map[string]interface{}{
		"AccountGroupId": config.AccountGroupId.ValueString(),
	}

	if !config.CompliancePackName.IsNull() && !config.CompliancePackName.IsUnknown() {
		paramMap["CompliancePackName"] = config.CompliancePackName.ValueString()
	}

	if !config.RiskLevel.IsNull() && !config.RiskLevel.IsUnknown() {
		var levels []int64
		resp.Diagnostics.Append(config.RiskLevel.ElementsAs(ctx, &levels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		levelsPtr := make([]*uint64, 0, len(levels))
		for _, v := range levels {
			u := uint64(v)
			levelsPtr = append(levelsPtr, &u)
		}
		paramMap["RiskLevel"] = levelsPtr
	}

	if !config.Status.IsNull() && !config.Status.IsUnknown() {
		paramMap["Status"] = config.Status.ValueString()
	}

	if !config.ComplianceResult.IsNull() && !config.ComplianceResult.IsUnknown() {
		var results []string
		resp.Diagnostics.Append(config.ComplianceResult.ElementsAs(ctx, &results, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resultsPtr := make([]*string, 0, len(results))
		for _, v := range results {
			s := v
			resultsPtr = append(resultsPtr, &s)
		}
		paramMap["ComplianceResult"] = resultsPtr
	}

	if !config.OrderType.IsNull() && !config.OrderType.IsUnknown() {
		paramMap["OrderType"] = config.OrderType.ValueString()
	}

	if !config.IncludeCompliancePackRuleResult.IsNull() && !config.IncludeCompliancePackRuleResult.IsUnknown() {
		paramMap["IncludeCompliancePackRuleResult"] = config.IncludeCompliancePackRuleResult.ValueString()
	}

	service := NewConfigService(d.Client())
	packs, total, err := service.DescribeConfigListAggregateCompliancePacksByFilter(ctx, paramMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying config aggregate compliance packs",
			fmt.Sprintf("ListAggregateCompliancePacks failed: %v", err),
		)
		return
	}

	items := make([]configListAggregateCompliancePacksItemModel, 0, len(packs))
	for _, pack := range packs {
		item := configListAggregateCompliancePacksItemModel{
			Status:             types.StringPointerValue(pack.Status),
			CompliancePackId:   types.StringPointerValue(pack.CompliancePackId),
			Description:        types.StringPointerValue(pack.Description),
			CompliancePackName: types.StringPointerValue(pack.CompliancePackName),
			ComplianceResult:   types.StringPointerValue(pack.ComplianceResult),
			CreateTime:         types.StringPointerValue(pack.CreateTime),
			NoCompliantNames:   types.SetNull(types.StringType),
			RuleCount:          types.Int64Null(),
		}
		if pack.RiskLevel != nil {
			item.RiskLevel = types.Int64Value(int64(*pack.RiskLevel))
		} else {
			item.RiskLevel = types.Int64Null()
		}
		if pack.RuleCount != nil {
			item.RuleCount = types.Int64Value(int64(*pack.RuleCount))
		}
		if pack.NoCompliantNames != nil {
			elements := make([]attr.Value, 0, len(pack.NoCompliantNames))
			for _, name := range pack.NoCompliantNames {
				elements = append(elements, types.StringPointerValue(name))
			}
			namesSet, diags := types.SetValue(types.StringType, elements)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			item.NoCompliantNames = namesSet
		}
		items = append(items, item)
	}

	itemsList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"status":              types.StringType,
			"compliance_pack_id":  types.StringType,
			"description":         types.StringType,
			"compliance_pack_name": types.StringType,
			"risk_level":          types.Int64Type,
			"compliance_result":   types.StringType,
			"create_time":         types.StringType,
			"no_compliant_names":  types.SetType{ElemType: types.StringType},
			"rule_count":          types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Total = types.Int64Value(int64(total))
	config.Items = itemsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

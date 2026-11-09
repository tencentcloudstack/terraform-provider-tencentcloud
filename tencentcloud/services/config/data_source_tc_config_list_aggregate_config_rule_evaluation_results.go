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

var _ datasource.DataSourceWithConfigure = &ConfigListAggregateConfigRuleEvaluationResultsDataSource{}

// NewConfigListAggregateConfigRuleEvaluationResultsDataSource is the factory
// referenced by tencentcloud/framework/registry.go to register this data source.
func NewConfigListAggregateConfigRuleEvaluationResultsDataSource() datasource.DataSource {
	return &ConfigListAggregateConfigRuleEvaluationResultsDataSource{}
}

// ConfigListAggregateConfigRuleEvaluationResultsDataSource implements
// datasource.DataSource for the
// tencentcloud_config_list_aggregate_config_rule_evaluation_results data source.
type ConfigListAggregateConfigRuleEvaluationResultsDataSource struct {
	fw.DataSourceWithConfigure
}

// ConfigListAggregateConfigRuleEvaluationResultsDataSourceModel maps the
// schema attributes.
type ConfigListAggregateConfigRuleEvaluationResultsDataSourceModel struct {
	ConfigRuleId     types.String `tfsdk:"config_rule_id"`
	AccountGroupId     types.String `tfsdk:"account_group_id"`
	ResourceType      types.List    `tfsdk:"resource_type"`
	ComplianceType   types.List    `tfsdk:"compliance_type"`
	ResourceOwnerId   types.Int64 `tfsdk:"resource_owner_id"`
	ResourceOwnerIds types.List  `tfsdk:"resource_owner_ids"`
	Total             types.Int64 `tfsdk:"total"`
	Items             types.List    `tfsdk:"items"`
}

// AggregateEvaluationResultItemModel is the element type of items.
type AggregateEvaluationResultItemModel struct {
	ResourceId                types.String `tfsdk:"resource_id"`
	ResourceType             types.String `tfsdk:"resource_type"`
	ResourceRegion           types.String `tfsdk:"resource_region"`
	ResourceName            types.String `tfsdk:"resource_name"`
	ConfigRuleId            types.String `tfsdk:"config_rule_id"`
	ConfigRuleName          types.String `tfsdk:"config_rule_name"`
	CompliancePackId         types.String `tfsdk:"compliance_pack_id"`
	RiskLevel               types.Int64 `tfsdk:"risk_level"`
	Annotation              types.List `tfsdk:"annotation"`
	ComplianceType          types.String `tfsdk:"compliance_type"`
	InvokingEventMessageType types.String `tfsdk:"invoking_event_message_type"`
	ConfigRuleInvokedTime        types.String `tfsdk:"config_rule_invoked_time"`
	ResultRecordedTime     types.String `tfsdk:"result_recorded_time"`
	ResourceOwnerId         types.Int64 `tfsdk:"resource_owner_id"`
	ResourceOwnerName       types.String `tfsdk:"resource_owner_name"`
}

// AggregateAnnotationItemModel is the element type of annotation.
type AggregateAnnotationItemModel struct {
	Configuration types.String `tfsdk:"configuration"`
	DesiredValue  types.String `tfsdk:"desired_value"`
	Operator      types.String `tfsdk:"operator"`
	Property      types.String `tfsdk:"property"`
}

func (d *ConfigListAggregateConfigRuleEvaluationResultsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_config_list_aggregate_config_rule_evaluation_results"
}

func (d *ConfigListAggregateConfigRuleEvaluationResultsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query Config aggregate config rule evaluation results " +
			"(rule dimension) for an account group.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs ----
			"config_rule_id": schema.StringAttribute{
				Required:    true,
				Description: "Config rule ID.",
			},
			"account_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Account group ID (aggregator ID).",
			},
			"resource_type": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Resource type list for filtering (e.g. QCS::CVM::Instance).",
			},
			"compliance_type": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Compliance type list for filtering. Valid values: COMPLIANT, NON_COMPLIANT.",
			},
			"resource_owner_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Resource owner UIN.",
			},
			"resource_owner_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Resource owner UIN list for filtering.",
			},
			// ---- Query outputs ----
			"total": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of matched evaluation results.",
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Evaluation result list.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_id": schema.StringAttribute{
							Computed:    true,
							Description: "Resource ID.",
						},
						"resource_type": schema.StringAttribute{
							Computed:    true,
							Description: "Resource type.",
						},
						"resource_region": schema.StringAttribute{
							Computed:    true,
							Description: "Resource region.",
						},
						"resource_name": schema.StringAttribute{
							Computed:    true,
							Description: "Resource name.",
						},
						"config_rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "Config rule ID.",
						},
						"config_rule_name": schema.StringAttribute{
							Computed:    true,
							Description: "Config rule name.",
						},
						"compliance_pack_id": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance pack ID. Note: this field may return null, indicating that no valid value was found.",
						},
						"risk_level": schema.Int64Attribute{
							Computed:    true,
							Description: "Risk level.",
						},
						"annotation": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Evaluation annotation detail. Note: this field may return null, indicating that no valid value was found.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"configuration": schema.StringAttribute{
										Computed:    true,
										Description: "Actual resource configuration (non-compliant configuration). Length: 0-256 characters.",
									},
									"desired_value": schema.StringAttribute{
										Computed:    true,
										Description: "Expected resource configuration (compliant configuration). Length: 0-256 characters.",
									},
									"operator": schema.StringAttribute{
										Computed:    true,
										Description: "Comparison operator between actual and expected configuration. Length: 0-16 characters. May be empty for custom rule evaluation results.",
									},
									"property": schema.StringAttribute{
										Computed:    true,
										Description: "JSON path of the current configuration in the resource attribute structure. Length: 0-256 characters. May be empty for custom rule evaluation results.",
									},
								},
							},
						},
						"compliance_type": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance type. Note: this field may return null, indicating that no valid value was found.",
						},
						"invoking_event_message_type": schema.StringAttribute{
							Computed:    true,
							Description: "Rule invocation type.",
						},
						"config_rule_invoked_time": schema.StringAttribute{
							Computed:    true,
							Description: "Evaluation invocation time.",
						},
						"result_recorded_time": schema.StringAttribute{
							Computed:    true,
							Description: "Evaluation result recorded time.",
						},
						"resource_owner_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Resource owner UIN. Note: this field may return null, indicating that no valid value was found.",
						},
						"resource_owner_name": schema.StringAttribute{
							Computed:    true,
							Description: "Resource owner name. Note: this field may return null, indicating that no valid value was found.",
						},
					},
				},
			},
		},
	}
}

// resultItemObjectType returns the attr.Type describing an items element.
func resultItemObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"resource_id":                 types.StringType,
			"resource_type":               types.StringType,
			"resource_region":             types.StringType,
			"resource_name":              types.StringType,
			"config_rule_id":            types.StringType,
			"config_rule_name":         types.StringType,
			"compliance_pack_id":      types.StringType,
			"risk_level":              types.Int64Type,
			"annotation":               annotationListType(),
			"compliance_type":         types.StringType,
			"invoking_event_message_type": types.StringType,
			"config_rule_invoked_time":    types.StringType,
			"result_recorded_time":       types.StringType,
			"resource_owner_id":          types.Int64Type,
			"resource_owner_name":       types.StringType,
		},
	}
}

// annotationListType returns the attr.Type describing the annotation list.
func annotationListType() attr.Type {
	return types.ListType{ElemType: annotationObjectType()}
}

// annotationObjectType returns the attr.Type describing an annotation element.
func annotationObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"configuration": types.StringType,
			"desired_value": types.StringType,
			"operator":      types.StringType,
			"property":      types.StringType,
		},
	}
}

// Read queries the aggregate config rule evaluation results and writes the
// result into resp.State.
func (d *ConfigListAggregateConfigRuleEvaluationResultsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ConfigListAggregateConfigRuleEvaluationResultsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramMap := make(map[string]interface{})
	paramMap["ConfigRuleId"] = config.ConfigRuleId.ValueString()
	paramMap["AccountGroupId"] = config.AccountGroupId.ValueString()

	if !config.ResourceType.IsNull() && !config.ResourceType.IsUnknown() {
		var resourceTypes []types.String
		resp.Diagnostics.Append(config.ResourceType.ElementsAs(ctx, &resourceTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ptrs := make([]*string, 0, len(resourceTypes))
		for _, rt := range resourceTypes {
			v := rt.ValueString()
			ptrs = append(ptrs, &v)
		}
		paramMap["ResourceType"] = ptrs
	}

	if !config.ComplianceType.IsNull() && !config.ComplianceType.IsUnknown() {
		var complianceTypes []types.String
		resp.Diagnostics.Append(config.ComplianceType.ElementsAs(ctx, &complianceTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ptrs := make([]*string, 0, len(complianceTypes))
		for _, ct := range complianceTypes {
			v := ct.ValueString()
			ptrs = append(ptrs, &v)
		}
		paramMap["ComplianceType"] = ptrs
	}

	if !config.ResourceOwnerId.IsNull() && !config.ResourceOwnerId.IsUnknown() {
		paramMap["ResourceOwnerId"] = uint64(config.ResourceOwnerId.ValueInt64())
	}

	if !config.ResourceOwnerIds.IsNull() && !config.ResourceOwnerIds.IsUnknown() {
		var ownerIDs []types.Int64
		resp.Diagnostics.Append(config.ResourceOwnerIds.ElementsAs(ctx, &ownerIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ptrs := make([]*uint64, 0, len(ownerIDs))
		for _, oid := range ownerIDs {
			v := uint64(oid.ValueInt64())
			ptrs = append(ptrs, &v)
		}
		paramMap["ResourceOwnerIds"] = ptrs
	}

	service := NewConfigService(d.Client())
	respItems, total, err := service.DescribeAggregateConfigRuleEvaluationResultsByFilter(ctx, paramMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying Config aggregate config rule evaluation results",
			fmt.Sprintf("ListAggregateConfigRuleEvaluationResults failed: %v", err),
		)
		return
	}

	items := make([]AggregateEvaluationResultItemModel, 0, len(respItems))
	for _, result := range respItems {
		item := AggregateEvaluationResultItemModel{
			ResourceId:                types.StringPointerValue(result.ResourceId),
			ResourceType:              types.StringPointerValue(result.ResourceType),
			ResourceRegion:            types.StringPointerValue(result.ResourceRegion),
			ResourceName:             types.StringPointerValue(result.ResourceName),
			ConfigRuleId:            types.StringPointerValue(result.ConfigRuleId),
			ConfigRuleName:         types.StringPointerValue(result.ConfigRuleName),
			CompliancePackId:      types.StringPointerValue(result.CompliancePackId),
			RiskLevel:              types.Int64PointerValue(result.RiskLevel),
			ComplianceType:         types.StringPointerValue(result.ComplianceType),
			InvokingEventMessageType: types.StringPointerValue(result.InvokingEventMessageType),
			ConfigRuleInvokedTime:    types.StringPointerValue(result.ConfigRuleInvokedTime),
			ResultRecordedTime:       types.StringPointerValue(result.ResultRecordedTime),
			ResourceOwnerId:          types.Int64PointerValue(result.ResourceOwnerId),
			ResourceOwnerName:       types.StringPointerValue(result.ResourceOwnerName),
		}

		if result.Annotation != nil {
			ann := AggregateAnnotationItemModel{
				Configuration: types.StringPointerValue(result.Annotation.Configuration),
				DesiredValue:  types.StringPointerValue(result.Annotation.DesiredValue),
				Operator:      types.StringPointerValue(result.Annotation.Operator),
				Property:      types.StringPointerValue(result.Annotation.Property),
			}
			annList, diags := types.ListValueFrom(ctx, annotationObjectType(), []AggregateAnnotationItemModel{ann})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			item.Annotation = annList
		} else {
			item.Annotation = types.ListNull(annotationObjectType())
		}

		items = append(items, item)
	}

	list, diags := types.ListValueFrom(ctx, resultItemObjectType(), items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Items = list
	config.Total = types.Int64Value(int64(total))

	// Clear input-only fields so they do not leak into state.
	config.ResourceType = types.ListNull(types.StringType)
	config.ComplianceType = types.ListNull(types.StringType)
	config.ResourceOwnerId = types.Int64Null()
	config.ResourceOwnerIds = types.ListNull(types.Int64Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
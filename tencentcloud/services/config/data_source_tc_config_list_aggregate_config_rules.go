package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var _ datasource.DataSourceWithConfigure = &ConfigListAggregateConfigRulesDataSource{}

// NewConfigListAggregateConfigRulesDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewConfigListAggregateConfigRulesDataSource() datasource.DataSource {
	return &ConfigListAggregateConfigRulesDataSource{}
}

// ConfigListAggregateConfigRulesDataSource implements datasource.DataSource for
// the tencentcloud_config_list_aggregate_config_rules data source.
type ConfigListAggregateConfigRulesDataSource struct {
	fw.DataSourceWithConfigure
}

// ConfigListAggregateConfigRulesDataSourceModel maps the schema attributes.
type ConfigListAggregateConfigRulesDataSourceModel struct {
	AccountGroupId   types.String `tfsdk:"account_group_id"`
	OrderType        types.String `tfsdk:"order_type"`
	RiskLevel        types.List   `tfsdk:"risk_level"`
	State            types.String `tfsdk:"state"`
	ComplianceResult types.List   `tfsdk:"compliance_result"`
	RuleName         types.String `tfsdk:"rule_name"`
	RuleOwnerId      types.Int64  `tfsdk:"rule_owner_id"`
	Total            types.Int64  `tfsdk:"total"`
	Items            types.List   `tfsdk:"items"`
}

// aggregateConfigRuleItemModel is the element type of items.
type aggregateConfigRuleItemModel struct {
	Identifier              types.String `tfsdk:"identifier"`
	RuleName                types.String `tfsdk:"rule_name"`
	InputParameter          types.List   `tfsdk:"input_parameter"`
	SourceCondition         types.List   `tfsdk:"source_condition"`
	ResourceType            types.List   `tfsdk:"resource_type"`
	Labels                  types.List   `tfsdk:"labels"`
	RiskLevel               types.Int64  `tfsdk:"risk_level"`
	ServiceFunction         types.String `tfsdk:"service_function"`
	CreateTime              types.String `tfsdk:"create_time"`
	Description             types.String `tfsdk:"description"`
	Status                  types.String `tfsdk:"status"`
	ComplianceResult        types.String `tfsdk:"compliance_result"`
	Annotation              types.List   `tfsdk:"annotation"`
	ConfigRuleInvokedTime   types.String `tfsdk:"config_rule_invoked_time"`
	ConfigRuleId            types.String `tfsdk:"config_rule_id"`
	IdentifierType          types.String `tfsdk:"identifier_type"`
	CompliancePackId        types.String `tfsdk:"compliance_pack_id"`
	TriggerType             types.List   `tfsdk:"trigger_type"`
	ManageInputParameter    types.List   `tfsdk:"manage_input_parameter"`
	CompliancePackName      types.String `tfsdk:"compliance_pack_name"`
	RegionsScope            types.List   `tfsdk:"regions_scope"`
	TagsScope               types.List   `tfsdk:"tags_scope"`
	ExcludeResourceIdsScope types.List   `tfsdk:"exclude_resource_ids_scope"`
	AccountGroupId          types.String `tfsdk:"account_group_id"`
	AccountGroupName        types.String `tfsdk:"account_group_name"`
	RuleOwnerId             types.Int64  `tfsdk:"rule_owner_id"`
	ManageTriggerType       types.List   `tfsdk:"manage_trigger_type"`
}

// aggregateInputParameterModel is the element type of input_parameter.
type aggregateInputParameterModel struct {
	ParameterKey types.String `tfsdk:"parameter_key"`
	Type         types.String `tfsdk:"type"`
	Value        types.String `tfsdk:"value"`
}

// aggregateSourceConditionModel is the element type of source_condition.
type aggregateSourceConditionModel struct {
	EmptyAs      types.String `tfsdk:"empty_as"`
	SelectPath   types.String `tfsdk:"select_path"`
	Operator     types.String `tfsdk:"operator"`
	Required     types.Bool   `tfsdk:"required"`
	DesiredValue types.String `tfsdk:"desired_value"`
}

// aggregateAnnotationModel is the element type of annotation.
type aggregateAnnotationModel struct {
	Configuration types.String `tfsdk:"configuration"`
	DesiredValue  types.String `tfsdk:"desired_value"`
	Operator      types.String `tfsdk:"operator"`
	Property      types.String `tfsdk:"property"`
}

// aggregateTriggerTypeModel is the element type of trigger_type.
type aggregateTriggerTypeModel struct {
	MessageType               types.String `tfsdk:"message_type"`
	MaximumExecutionFrequency types.String `tfsdk:"maximum_execution_frequency"`
}

// aggregateManageInputParameterModel is the element type of manage_input_parameter.
type aggregateManageInputParameterModel struct {
	ValueType    types.String `tfsdk:"value_type"`
	ParameterKey types.String `tfsdk:"parameter_key"`
	Type         types.String `tfsdk:"type"`
	DefaultValue types.String `tfsdk:"default_value"`
	Description  types.String `tfsdk:"description"`
}

// aggregateTagsScopeModel is the element type of tags_scope.
type aggregateTagsScopeModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

func (d *ConfigListAggregateConfigRulesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_config_list_aggregate_config_rules"
}

func (d *ConfigListAggregateConfigRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of Config aggregate config rules under an account group.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs ----
			"account_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Account group ID.",
			},
			"order_type": schema.StringAttribute{
				Optional:    true,
				Description: "Sort order of the rule list. Valid values: `desc` (descending), `asc` (ascending).",
			},
			"risk_level": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Risk level list for filtering. Valid values: `1` (high risk), `2` (medium risk), `3` (low risk).",
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Description: "Rule state for filtering.",
			},
			"compliance_result": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Compliance result list for filtering.",
			},
			"rule_name": schema.StringAttribute{
				Optional:    true,
				Description: "Rule name for filtering.",
			},
			"rule_owner_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Rule owner account ID.",
			},
			// ---- Query outputs ----
			"total": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of matched rules.",
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of matched aggregate config rules.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							Computed:    true,
							Description: "Rule identifier.",
						},
						"rule_name": schema.StringAttribute{
							Computed:    true,
							Description: "Rule name.",
						},
						"input_parameter": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Rule parameters.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"parameter_key": schema.StringAttribute{
										Computed:    true,
										Description: "Parameter name.",
									},
									"type": schema.StringAttribute{
										Computed:    true,
										Description: "Parameter type. `Require` for required, `Optional` for optional.",
									},
									"value": schema.StringAttribute{
										Computed:    true,
										Description: "Parameter value.",
									},
								},
							},
						},
						"source_condition": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Rule trigger conditions.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"empty_as": schema.StringAttribute{
										Computed:    true,
										Description: "Result when the condition is empty. `COMPLIANT` for compliant, `NON_COMPLIANT` for non-compliant, `NOT_APPLICABLE` for not applicable.",
									},
									"select_path": schema.StringAttribute{
										Computed:    true,
										Description: "Configuration path.",
									},
									"operator": schema.StringAttribute{
										Computed:    true,
										Description: "Comparison operator.",
									},
									"required": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether the condition is required.",
									},
									"desired_value": schema.StringAttribute{
										Computed:    true,
										Description: "Expected value.",
									},
								},
							},
						},
						"resource_type": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Supported resource types. The rule only applies to the specified resource types.",
						},
						"labels": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Rule labels.",
						},
						"risk_level": schema.Int64Attribute{
							Computed:    true,
							Description: "Rule risk level. `1` for low risk, `2` for medium risk, `3` for high risk.",
						},
						"service_function": schema.StringAttribute{
							Computed:    true,
							Description: "Corresponding service function.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time in the format YYYY-MM-DD h:i:s.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Rule description.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Rule status. `ACTIVE` for enabled, `NO_ACTIVE` for stopped.",
						},
						"compliance_result": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance result. `COMPLIANT` for compliant, `NON_COMPLIANT` for non-compliant, `NOT_APPLICABLE` for not applicable.",
						},
						"annotation": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Annotation of the non-compliant resource.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"configuration": schema.StringAttribute{
										Computed:    true,
										Description: "Actual configuration of the resource, i.e. the non-compliant part. Up to 256 characters.",
									},
									"desired_value": schema.StringAttribute{
										Computed:    true,
										Description: "Expected configuration of the resource. Up to 256 characters.",
									},
									"operator": schema.StringAttribute{
										Computed:    true,
										Description: "Comparison operator between the actual and expected configuration. Up to 16 characters. May be empty for custom rules.",
									},
									"property": schema.StringAttribute{
										Computed:    true,
										Description: "JSON path of the current configuration in the resource attribute structure. Up to 256 characters. May be empty for custom rules.",
									},
								},
							},
						},
						"config_rule_invoked_time": schema.StringAttribute{
							Computed:    true,
							Description: "Rule evaluation time in the format YYYY-MM-DD h:i:s.",
						},
						"config_rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "Rule ID.",
						},
						"identifier_type": schema.StringAttribute{
							Computed:    true,
							Description: "Rule type. `CUSTOMIZE` for custom rule, `SYSTEM` for managed rule.",
						},
						"compliance_pack_id": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance pack ID.",
						},
						"trigger_type": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Trigger types. `ScheduledNotification` for periodic trigger, `ConfigurationItemChangeNotification` for change trigger.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"message_type": schema.StringAttribute{
										Computed:    true,
										Description: "Trigger type.",
									},
									"maximum_execution_frequency": schema.StringAttribute{
										Computed:    true,
										Description: "Trigger period.",
									},
								},
							},
						},
						"manage_input_parameter": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Parameter details.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"value_type": schema.StringAttribute{
										Computed:    true,
										Description: "Value type. `Integer` for number, `String` for string.",
									},
									"parameter_key": schema.StringAttribute{
										Computed:    true,
										Description: "Parameter key.",
									},
									"type": schema.StringAttribute{
										Computed:    true,
										Description: "Parameter type. `Require` for required, `Optional` for optional.",
									},
									"default_value": schema.StringAttribute{
										Computed:    true,
										Description: "Default value.",
									},
									"description": schema.StringAttribute{
										Computed:    true,
										Description: "Description.",
									},
								},
							},
						},
						"compliance_pack_name": schema.StringAttribute{
							Computed:    true,
							Description: "Compliance pack name.",
						},
						"regions_scope": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Associated regions.",
						},
						"tags_scope": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Associated tags.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"tag_key": schema.StringAttribute{
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": schema.StringAttribute{
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"exclude_resource_ids_scope": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Resource IDs excluded from evaluation.",
						},
						"account_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "Account group ID.",
						},
						"account_group_name": schema.StringAttribute{
							Computed:    true,
							Description: "Account group name.",
						},
						"rule_owner_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Rule owner user ID.",
						},
						"manage_trigger_type": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Trigger types supported by preset rules. `ScheduledNotification` for periodic, `ConfigurationItemChangeNotification` for change.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *ConfigListAggregateConfigRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ConfigListAggregateConfigRulesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramMap := make(map[string]interface{})
	paramMap["AccountGroupId"] = config.AccountGroupId.ValueString()

	if !config.OrderType.IsNull() && !config.OrderType.IsUnknown() {
		paramMap["OrderType"] = config.OrderType.ValueString()
	}

	if !config.RiskLevel.IsNull() && !config.RiskLevel.IsUnknown() {
		riskLevels := make([]types.Int64, 0)
		resp.Diagnostics.Append(config.RiskLevel.ElementsAs(ctx, &riskLevels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ptrs := make([]*uint64, 0, len(riskLevels))
		for _, v := range riskLevels {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			ptrs = append(ptrs, helper.Uint64(uint64(v.ValueInt64())))
		}
		paramMap["RiskLevel"] = ptrs
	}

	if !config.State.IsNull() && !config.State.IsUnknown() {
		paramMap["State"] = config.State.ValueString()
	}

	if !config.ComplianceResult.IsNull() && !config.ComplianceResult.IsUnknown() {
		results := make([]types.String, 0)
		resp.Diagnostics.Append(config.ComplianceResult.ElementsAs(ctx, &results, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ptrs := make([]*string, 0, len(results))
		for _, v := range results {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			ptrs = append(ptrs, helper.String(v.ValueString()))
		}
		paramMap["ComplianceResult"] = ptrs
	}

	if !config.RuleName.IsNull() && !config.RuleName.IsUnknown() {
		paramMap["RuleName"] = config.RuleName.ValueString()
	}

	if !config.RuleOwnerId.IsNull() && !config.RuleOwnerId.IsUnknown() {
		paramMap["RuleOwnerId"] = uint64(config.RuleOwnerId.ValueInt64())
	}

	service := NewConfigService(d.Client())
	items, total, errRet := service.DescribeConfigAggregateConfigRulesByFilter(ctx, paramMap)
	if errRet != nil {
		resp.Diagnostics.AddError(
			"Error querying Config aggregate config rules",
			fmt.Sprintf("ListAggregateConfigRules failed: %v", errRet),
		)
		return
	}

	config.Total = types.Int64Value(int64(total))
	ruleItems := make([]aggregateConfigRuleItemModel, 0, len(items))
	for _, item := range items {
		ruleItems = append(ruleItems, flattenAggregateConfigRule(ctx, resp, item))
		if resp.Diagnostics.HasError() {
			return
		}
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: aggregateConfigRuleItemAttrTypes()}, ruleItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Items = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

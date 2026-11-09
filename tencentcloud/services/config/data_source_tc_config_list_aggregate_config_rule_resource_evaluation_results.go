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

var _ datasource.DataSourceWithConfigure = &ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource{}

// NewConfigListAggregateConfigRuleResourceEvaluationResultsDataSource is the
// factory referenced by tencentcloud/framework/registry.go to register this
// data source.
func NewConfigListAggregateConfigRuleResourceEvaluationResultsDataSource() datasource.DataSource {
	return &ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource{}
}

// ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource implements
// datasource.DataSource for the
// tencentcloud_config_list_aggregate_config_rule_resource_evaluation_results
// data source.
type ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource struct {
	fw.DataSourceWithConfigure
}

// ConfigListAggregateConfigRuleResourceEvaluationResultsModel maps the schema
// attributes of the data source.
type ConfigListAggregateConfigRuleResourceEvaluationResultsModel struct {
	AccountGroupId types.String `tfsdk:"account_group_id"`
	Items          types.List   `tfsdk:"items"`
}

// AggregateConfigRuleResourceEvaluationResultItemModel is the element type of items.
type AggregateConfigRuleResourceEvaluationResultItemModel struct {
	RuleId          types.String `tfsdk:"rule_id"`
	RuleName        types.String `tfsdk:"rule_name"`
	RuleIdentifier  types.String `tfsdk:"rule_identifier"`
	RuleDescription types.String `tfsdk:"rule_description"`
	RuleOwnerId     types.Int64  `tfsdk:"rule_owner_id"`
	RuleRiskLevel   types.Int64  `tfsdk:"rule_risk_level"`
	ResourceId      types.String `tfsdk:"resource_id"`
	ResourceType    types.String `tfsdk:"resource_type"`
	ResourceName    types.String `tfsdk:"resource_name"`
	ResourceRegion  types.String `tfsdk:"resource_region"`
	ResourceTags    types.List   `tfsdk:"resource_tags"`
	Annotation      types.List   `tfsdk:"annotation"`
}

// ResourceTagItemModel is the element type of items.resource_tags.
type ResourceTagItemModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// AnnotationItemModel is the element type of items.annotation.
type AnnotationItemModel struct {
	Configuration types.String `tfsdk:"configuration"`
	DesiredValue  types.String `tfsdk:"desired_value"`
	Operator      types.String `tfsdk:"operator"`
	Property      types.String `tfsdk:"property"`
}

func (d *ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_config_list_aggregate_config_rule_resource_evaluation_results"
}

func (d *ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query Config aggregate config rule resource evaluation results.",
		Attributes: map[string]schema.Attribute{
			"account_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Account group ID. When omitted, the aggregate full results are queried.",
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of config rule resource evaluation results.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "Rule ID.",
						},
						"rule_name": schema.StringAttribute{
							Computed:    true,
							Description: "Rule name.",
						},
						"rule_identifier": schema.StringAttribute{
							Computed:    true,
							Description: "Rule identifier.",
						},
						"rule_description": schema.StringAttribute{
							Computed:    true,
							Description: "Rule description.",
						},
						"rule_owner_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Owner account ID of the rule.",
						},
						"rule_risk_level": schema.Int64Attribute{
							Computed:    true,
							Description: "Risk level of the rule.",
						},
						"resource_id": schema.StringAttribute{
							Computed:    true,
							Description: "Resource ID.",
						},
						"resource_type": schema.StringAttribute{
							Computed:    true,
							Description: "Resource type.",
						},
						"resource_name": schema.StringAttribute{
							Computed:    true,
							Description: "Resource name.",
						},
						"resource_region": schema.StringAttribute{
							Computed:    true,
							Description: "Resource region.",
						},
						"resource_tags": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Resource tags.",
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
						"annotation": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Compliance detail.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"configuration": schema.StringAttribute{
										Computed:    true,
										Description: "Actual resource configuration, i.e. the non-compliant configuration. Up to 256 characters.",
									},
									"desired_value": schema.StringAttribute{
										Computed:    true,
										Description: "Expected resource configuration, i.e. the compliant configuration. Up to 256 characters.",
									},
									"operator": schema.StringAttribute{
										Computed:    true,
										Description: "Comparison operator between the actual and expected configuration. Up to 16 characters. This field may be empty for custom rules.",
									},
									"property": schema.StringAttribute{
										Computed:    true,
										Description: "JSON path of the current configuration in the resource attribute structure. Up to 256 characters. This field may be empty for custom rules.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ConfigListAggregateConfigRuleResourceEvaluationResultsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramMap := make(map[string]interface{})
	if !config.AccountGroupId.IsNull() && !config.AccountGroupId.IsUnknown() {
		paramMap["AccountGroupId"] = config.AccountGroupId.ValueString()
	}

	service := NewConfigService(d.Client())
	results, err := service.DescribeConfigListAggregateConfigRuleResourceEvaluationResultsByFilter(ctx, paramMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying Config aggregate config rule resource evaluation results",
			fmt.Sprintf("ListAggregateConfigRuleResourceEvaluationResults failed: %v", err),
		)
		return
	}

	items := make([]AggregateConfigRuleResourceEvaluationResultItemModel, 0, len(results))
	for _, result := range results {
		item := AggregateConfigRuleResourceEvaluationResultItemModel{
			RuleId:          types.StringPointerValue(result.RuleId),
			RuleName:        types.StringPointerValue(result.RuleName),
			RuleIdentifier:  types.StringPointerValue(result.RuleIdentifier),
			RuleDescription: types.StringPointerValue(result.RuleDescription),
			RuleOwnerId:     types.Int64PointerValue(result.RuleOwnerId),
			RuleRiskLevel:   types.Int64PointerValue(result.RuleRiskLevel),
			ResourceId:      types.StringPointerValue(result.ResourceId),
			ResourceType:    types.StringPointerValue(result.ResourceType),
			ResourceName:    types.StringPointerValue(result.ResourceName),
			ResourceRegion:  types.StringPointerValue(result.ResourceRegion),
		}

		// resource_tags: optional list in the response, may be nil.
		tagElements := make([]ResourceTagItemModel, 0, len(result.ResourceTags))
		for _, tag := range result.ResourceTags {
			tagElements = append(tagElements, ResourceTagItemModel{
				TagKey:   types.StringPointerValue(tag.TagKey),
				TagValue: types.StringPointerValue(tag.TagValue),
			})
		}
		tagList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: resourceTagAttrTypes()}, tagElements)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		item.ResourceTags = tagList

		// annotation: a single object in the SDK; rendered as a list of at most
		// one element to match the schema. nil annotation yields an empty list.
		annotationElements := make([]AnnotationItemModel, 0)
		if result.Annotation != nil {
			annotationElements = append(annotationElements, AnnotationItemModel{
				Configuration: types.StringPointerValue(result.Annotation.Configuration),
				DesiredValue:  types.StringPointerValue(result.Annotation.DesiredValue),
				Operator:      types.StringPointerValue(result.Annotation.Operator),
				Property:      types.StringPointerValue(result.Annotation.Property),
			})
		}
		annotationList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: annotationAttrTypes()}, annotationElements)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		item.Annotation = annotationList

		items = append(items, item)
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: evaluationResultItemAttrTypes()}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Items = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// evaluationResultItemAttrTypes returns the attr.Type mapping for a single
// items element, used when building types.ListValueFrom.
func evaluationResultItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"rule_id":          types.StringType,
		"rule_name":        types.StringType,
		"rule_identifier":  types.StringType,
		"rule_description": types.StringType,
		"rule_owner_id":    types.Int64Type,
		"rule_risk_level":  types.Int64Type,
		"resource_id":      types.StringType,
		"resource_type":    types.StringType,
		"resource_name":    types.StringType,
		"resource_region":  types.StringType,
		"resource_tags":    types.ListType{ElemType: types.ObjectType{AttrTypes: resourceTagAttrTypes()}},
		"annotation":       types.ListType{ElemType: types.ObjectType{AttrTypes: annotationAttrTypes()}},
	}
}

// resourceTagAttrTypes returns the attr.Type mapping for a single
// resource_tags element.
func resourceTagAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}
}

// annotationAttrTypes returns the attr.Type mapping for a single annotation
// element.
func annotationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"configuration": types.StringType,
		"desired_value": types.StringType,
		"operator":      types.StringType,
		"property":      types.StringType,
	}
}

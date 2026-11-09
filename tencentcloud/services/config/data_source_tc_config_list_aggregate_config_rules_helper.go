package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"
)

// aggregateConfigRuleItemAttrTypes returns the attribute types of the items element object.
func aggregateConfigRuleItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"identifier":               types.StringType,
		"rule_name":               types.StringType,
		"input_parameter":         types.ListType{ElemType: types.ObjectType{AttrTypes: aggregateInputParameterAttrTypes()}},
		"source_condition":       types.ListType{ElemType: types.ObjectType{AttrTypes: aggregateSourceConditionAttrTypes()}},
		"resource_type":           types.ListType{ElemType: types.StringType},
		"labels":                   types.ListType{ElemType: types.StringType},
		"risk_level":               types.Int64Type,
		"service_function":        types.StringType,
		"create_time":              types.StringType,
		"description":              types.StringType,
		"status":                   types.StringType,
		"compliance_result":      types.StringType,
		"annotation":             types.ListType{ElemType: types.ObjectType{AttrTypes: aggregateAnnotationAttrTypes()}},
		"config_rule_invoked_time": types.StringType,
		"config_rule_id":           types.StringType,
		"identifier_type":          types.StringType,
		"compliance_pack_id":       types.StringType,
		"trigger_type":              types.ListType{ElemType: types.ObjectType{AttrTypes: aggregateTriggerTypeAttrTypes()}},
		"manage_input_parameter":    types.ListType{ElemType: types.ObjectType{AttrTypes: aggregateManageInputParameterAttrTypes()}},
		"compliance_pack_name":      types.StringType,
		"regions_scope":             types.ListType{ElemType: types.StringType},
		"tags_scope":                types.ListType{ElemType: types.ObjectType{AttrTypes: aggregateTagsScopeAttrTypes()}},
		"exclude_resource_ids_scope": types.ListType{ElemType: types.StringType},
		"account_group_id":        types.StringType,
		"account_group_name":     types.StringType,
		"rule_owner_id":           types.Int64Type,
		"manage_trigger_type":    types.ListType{ElemType: types.StringType},
	}
}

func aggregateInputParameterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"parameter_key": types.StringType,
		"type":          types.StringType,
		"value":         types.StringType,
	}
}

func aggregateSourceConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"empty_as":      types.StringType,
		"select_path":   types.StringType,
		"operator":     types.StringType,
		"required":      types.BoolType,
		"desired_value": types.StringType,
	}
}

func aggregateAnnotationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"configuration": types.StringType,
		"desired_value": types.StringType,
		"operator":      types.StringType,
		"property":      types.StringType,
	}
}

func aggregateTriggerTypeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"message_type":              types.StringType,
		"maximum_execution_frequency": types.StringType,
	}
}

func aggregateManageInputParameterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value_type":    types.StringType,
		"parameter_key": types.StringType,
		"type":          types.StringType,
		"default_value": types.StringType,
		"description":   types.StringType,
	}
}

func aggregateTagsScopeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}
}

// flattenAggregateConfigRule converts a ConfigRule into the data source model.
func flattenAggregateConfigRule(ctx context.Context, resp *datasource.ReadResponse, item *configv20220802.ConfigRule) aggregateConfigRuleItemModel {
	if item == nil {
		return aggregateConfigRuleItemModel{}
	}

	out := aggregateConfigRuleItemModel{
		Identifier:           types.StringPointerValue(item.Identifier),
		RuleName:             types.StringPointerValue(item.RuleName),
		RiskLevel:            types.Int64PointerValue(item.RiskLevel),
		ServiceFunction:     types.StringPointerValue(item.ServiceFunction),
		CreateTime:          types.StringPointerValue(item.CreateTime),
		Description:         types.StringPointerValue(item.Description),
		Status:               types.StringPointerValue(item.Status),
		ComplianceResult:     types.StringPointerValue(item.ComplianceResult),
		ConfigRuleInvokedTime: types.StringPointerValue(item.ConfigRuleInvokedTime),
		ConfigRuleId:         types.StringPointerValue(item.ConfigRuleId),
		IdentifierType:       types.StringPointerValue(item.IdentifierType),
		CompliancePackId:     types.StringPointerValue(item.CompliancePackId),
		CompliancePackName:   types.StringPointerValue(item.CompliancePackName),
		AccountGroupId:      types.StringPointerValue(item.AccountGroupId),
		AccountGroupName:    types.StringPointerValue(item.AccountGroupName),
		RuleOwnerId:          types.Int64PointerValue(int64PtrFromUint64(item.RuleOwnerId)),
	}

	out.InputParameter = flattenAggregateInputParameter(ctx, resp, item.InputParameter)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.SourceCondition = flattenAggregateSourceCondition(ctx, resp, item.SourceCondition)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.ResourceType = flattenStringList(ctx, resp, item.ResourceType)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.Labels = flattenStringList(ctx, resp, item.Labels)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.Annotation = flattenAggregateAnnotation(ctx, resp, item.Annotation)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.TriggerType = flattenAggregateTriggerType(ctx, resp, item.TriggerType)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.ManageInputParameter = flattenAggregateManageInputParameter(ctx, resp, item.ManageInputParameter)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.RegionsScope = flattenStringList(ctx, resp, item.RegionsScope)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.TagsScope = flattenAggregateTagsScope(ctx, resp, item.TagsScope)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.ExcludeResourceIdsScope = flattenStringList(ctx, resp, item.ExcludeResourceIdsScope)
	if resp.Diagnostics.HasError() {
		return out
	}
	out.ManageTriggerType = flattenStringList(ctx, resp, item.ManageTriggerType)

	return out
}

func flattenAggregateInputParameter(ctx context.Context, resp *datasource.ReadResponse, items []*configv20220802.InputParameter) types.List {
	out := make([]aggregateInputParameterModel, 0, len(items))
	for _, it := range items {
		if it == nil {
			continue
		}
		out = append(out, aggregateInputParameterModel{
			ParameterKey: types.StringPointerValue(it.ParameterKey),
			Type:       types.StringPointerValue(it.Type),
			Value:        types.StringPointerValue(it.Value),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: aggregateInputParameterAttrTypes()}, out)
	resp.Diagnostics.Append(diags...)
	return list
}

func flattenAggregateSourceCondition(ctx context.Context, resp *datasource.ReadResponse, items []*configv20220802.SourceConditionForManage) types.List {
	out := make([]aggregateSourceConditionModel, 0, len(items))
	for _, it := range items {
		if it == nil {
			continue
		}
		out = append(out, aggregateSourceConditionModel{
			EmptyAs:     types.StringPointerValue(it.EmptyAs),
			SelectPath:  types.StringPointerValue(it.SelectPath),
			Operator:    types.StringPointerValue(it.Operator),
			Required:   types.BoolPointerValue(it.Required),
			DesiredValue: types.StringPointerValue(it.DesiredValue),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: aggregateSourceConditionAttrTypes()}, out)
	resp.Diagnostics.Append(diags...)
	return list
}

func flattenAggregateAnnotation(ctx context.Context, resp *datasource.ReadResponse, annotation *configv20220802.Annotation) types.List {
	if annotation == nil {
		list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: aggregateAnnotationAttrTypes()}, []aggregateAnnotationModel{})
		resp.Diagnostics.Append(diags...)
		return list
	}
	out := []aggregateAnnotationModel{{
		Configuration: types.StringPointerValue(annotation.Configuration),
		DesiredValue:  types.StringPointerValue(annotation.DesiredValue),
		Operator:       types.StringPointerValue(annotation.Operator),
		Property:       types.StringPointerValue(annotation.Property),
	}}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: aggregateAnnotationAttrTypes()}, out)
	resp.Diagnostics.Append(diags...)
	return list
}

func flattenAggregateTriggerType(ctx context.Context, resp *datasource.ReadResponse, items []*configv20220802.TriggerType) types.List {
	out := make([]aggregateTriggerTypeModel, 0, len(items))
	for _, it := range items {
		if it == nil {
			continue
		}
		out = append(out, aggregateTriggerTypeModel{
			MessageType:              types.StringPointerValue(it.MessageType),
			MaximumExecutionFrequency: types.StringPointerValue(it.MaximumExecutionFrequency),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: aggregateTriggerTypeAttrTypes()}, out)
	resp.Diagnostics.Append(diags...)
	return list
}

func flattenAggregateManageInputParameter(ctx context.Context, resp *datasource.ReadResponse, items []*configv20220802.InputParameterForManage) types.List {
	out := make([]aggregateManageInputParameterModel, 0, len(items))
	for _, it := range items {
		if it == nil {
			continue
		}
		out = append(out, aggregateManageInputParameterModel{
			ValueType:    types.StringPointerValue(it.ValueType),
			ParameterKey: types.StringPointerValue(it.ParameterKey),
			Type:        types.StringPointerValue(it.Type),
			DefaultValue: types.StringPointerValue(it.DefaultValue),
			Description:   types.StringPointerValue(it.Description),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: aggregateManageInputParameterAttrTypes()}, out)
	resp.Diagnostics.Append(diags...)
	return list
}

func flattenAggregateTagsScope(ctx context.Context, resp *datasource.ReadResponse, items []*configv20220802.Tag) types.List {
	out := make([]aggregateTagsScopeModel, 0, len(items))
	for _, it := range items {
		if it == nil {
			continue
		}
		out = append(out, aggregateTagsScopeModel{
			TagKey:   types.StringPointerValue(it.TagKey),
			TagValue: types.StringPointerValue(it.TagValue),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: aggregateTagsScopeAttrTypes()}, out)
	resp.Diagnostics.Append(diags...)
	return list
}

// flattenStringList converts a []*string into a types.List of strings.
func flattenStringList(ctx context.Context, resp *datasource.ReadResponse, items []*string) types.List {
	out := make([]types.String, 0, len(items))
	for _, it := range items {
		if it == nil {
			continue
		}
		out = append(out, types.StringPointerValue(it))
	}
	list, diags := types.ListValueFrom(ctx, types.StringType, out)
	resp.Diagnostics.Append(diags...)
	return list
}

// int64PtrFromUint64 converts a *uint64 to *int64, returning nil for a nil input.
func int64PtrFromUint64(v *uint64) *int64 {
	if v == nil {
		return nil
	}
	val := int64(*v)
	return &val
}
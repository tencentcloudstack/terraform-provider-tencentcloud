package monitor

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	monitorv20230616 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var (
	_ resource.ResourceWithConfigure   = &MonitorDispenseExternalRuleResource{}
	_ resource.ResourceWithImportState = &MonitorDispenseExternalRuleResource{}
)

// NewMonitorDispenseExternalRuleResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewMonitorDispenseExternalRuleResource() resource.Resource {
	return &MonitorDispenseExternalRuleResource{}
}

// MonitorDispenseExternalRuleResource implements resource.Resource for
// tencentcloud_monitor_dispense_external_rule.
type MonitorDispenseExternalRuleResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// MonitorDispenseExternalRuleResourceModel maps the schema attributes.
//
// The writable attributes (name, ext_namespace, producer, dispense_regions,
// ext_metrics, period, dispense_conditions) are echoed back by
// DescribeDispenseExternalRule through the computed `rule` block only. Read
// therefore refreshes `rule` and keeps every writable attribute as stored in
// state: it never overwrites them from the response.
type MonitorDispenseExternalRuleResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	RuleId             types.Int64  `tfsdk:"rule_id"`
	Name               types.String `tfsdk:"name"`
	ExtNamespace       types.String `tfsdk:"ext_namespace"`
	Producer           types.Object `tfsdk:"producer"`
	DispenseRegions    types.Set    `tfsdk:"dispense_regions"`
	ExtMetrics         types.Set    `tfsdk:"ext_metrics"`
	Period             types.Set    `tfsdk:"period"`
	DispenseConditions types.Object `tfsdk:"dispense_conditions"`
	Rule               types.Object `tfsdk:"rule"`
}

type monitorDispenseProducerModel struct {
	ProtocolType types.Int64  `tfsdk:"protocol_type"`
	Type         types.String `tfsdk:"type"`
	Brokers      types.String `tfsdk:"brokers"`
	Topic        types.String `tfsdk:"topic"`
	Merge        types.Int64  `tfsdk:"merge"`
	GlobalTags   types.List   `tfsdk:"global_tags"`
	DefaultTags  types.Set    `tfsdk:"default_tags"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
}

type monitorDispenseGlobalTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type monitorDispenseConditionModel struct {
	ExtMetric       types.String `tfsdk:"ext_metric"`
	DispenseFilters types.List   `tfsdk:"dispense_filters"`
	ConditionId     types.Int64  `tfsdk:"condition_id"`
}

type monitorDispenseFilterModel struct {
	Key        types.String `tfsdk:"key"`
	Values     types.Set    `tfsdk:"values"`
	Expression types.String `tfsdk:"expression"`
}

type monitorDispenseExtMetricModel struct {
	MetricName    types.String `tfsdk:"metric_name"`
	MetricCName   types.String `tfsdk:"metric_c_name"`
	CNMeaning     types.String `tfsdk:"cn_meaning"`
	EnMeaning     types.String `tfsdk:"en_meaning"`
	Unit          types.String `tfsdk:"unit"`
	DimensionFlag types.Bool   `tfsdk:"dimension_flag"`
}

type monitorDispenseRuleModel struct {
	RuleId             types.Int64  `tfsdk:"rule_id"`
	Name               types.String `tfsdk:"name"`
	ExtNamespace       types.String `tfsdk:"ext_namespace"`
	ExtMetric          types.List   `tfsdk:"ext_metric"`
	Producer           types.Object `tfsdk:"producer"`
	UpdateTime         types.Int64  `tfsdk:"update_time"`
	Status             types.Int64  `tfsdk:"status"`
	Period             types.Set    `tfsdk:"period"`
	DispenseConditions types.Object `tfsdk:"dispense_conditions"`
	DispenseRegions    types.Set    `tfsdk:"dispense_regions"`
}

var monitorDispenseGlobalTagAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Optional:    true,
		Description: "Dimension key.",
	},
	"value": schema.StringAttribute{
		Optional:    true,
		Description: "Dimension value.",
	},
}

var monitorDispenseFilterAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Optional:    true,
		Description: "Dimension name.",
	},
	"values": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Dimension value list.",
	},
	"expression": schema.StringAttribute{
		Optional:    true,
		Description: "Expression.",
	},
}

var monitorDispenseConditionAttributes = map[string]schema.Attribute{
	"ext_metric": schema.StringAttribute{
		Optional:    true,
		Description: "External metric name.",
	},
	"dispense_filters": schema.ListNestedAttribute{
		Optional:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: monitorDispenseFilterAttributes},
		Description:  "Filter condition list.",
	},
	"condition_id": schema.Int64Attribute{
		Optional:    true,
		Description: "Filter condition id.",
	},
}

var monitorDispenseProducerAttributes = map[string]schema.Attribute{
	"protocol_type": schema.Int64Attribute{
		Required:    true,
		Description: "Dispense protocol type. Valid values: `0` stormRetPb, `1` tcbDispensePb, `2` stormRetJson, `3` ADPPb (deprecated), `4` zhongtai pb.",
	},
	"type": schema.StringAttribute{
		Required:    true,
		Description: "Target type.",
	},
	"brokers": schema.StringAttribute{
		Required:    true,
		Description: "Kafka broker address to dispense to.",
	},
	"topic": schema.StringAttribute{
		Required:    true,
		Description: "Kafka topic to dispense to.",
	},
	"merge": schema.Int64Attribute{
		Optional:    true,
		Description: "Whether to merge metrics. Default value: `1` (merged).",
	},
	"global_tags": schema.ListNestedAttribute{
		Optional:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: monitorDispenseGlobalTagAttributes},
		Description:  "Global dimension groups.",
	},
	"default_tags": schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Default dimension groups. Only dimensions are provided.",
	},
	"username": schema.StringAttribute{
		Optional:    true,
		Description: "Kafka user name.",
	},
	"password": schema.StringAttribute{
		Optional:  true,
		Sensitive: true,
		// Password is echoed back by Describe only inside the computed `rule`
		// block (possibly masked). The writable attribute is never refreshed
		// from the response, so a masked echo does not create a diff.
		Description: "Kafka password.",
	},
}

var monitorDispenseExtMetricAttributes = map[string]schema.Attribute{
	"metric_name": schema.StringAttribute{
		Computed:    true,
		Description: "Metric name.",
	},
	"metric_c_name": schema.StringAttribute{
		Computed:    true,
		Description: "Chinese metric name.",
	},
	"cn_meaning": schema.StringAttribute{
		Computed:    true,
		Description: "Chinese meaning of the metric.",
	},
	"en_meaning": schema.StringAttribute{
		Computed:    true,
		Description: "English meaning of the metric.",
	},
	"unit": schema.StringAttribute{
		Computed:    true,
		Description: "Metric unit.",
	},
	"dimension_flag": schema.BoolAttribute{
		Computed:    true,
		Description: "Whether external dimensions are configured for the metric.",
	},
}

var monitorDispenseRuleAttributes = map[string]schema.Attribute{
	"rule_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Rule id.",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Rule name.",
	},
	"ext_namespace": schema.StringAttribute{
		Computed:    true,
		Description: "External namespace.",
	},
	"ext_metric": schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: monitorDispenseExtMetricAttributes},
		Description:  "External metric list.",
	},
	"producer": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  monitorDispenseProducerAttributes,
		Description: "Producer information.",
	},
	"update_time": schema.Int64Attribute{
		Computed:    true,
		Description: "Update time.",
	},
	"status": schema.Int64Attribute{
		Computed:    true,
		Description: "Rule trigger status.",
	},
	"period": schema.SetAttribute{
		Computed:    true,
		ElementType: types.Int64Type,
		Description: "Metric statistical period list.",
	},
	"dispense_conditions": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  monitorDispenseConditionAttributes,
		Description: "Dispense filter conditions.",
	},
	"dispense_regions": schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Description: "Dispense region list.",
	},
}

var (
	monitorDispenseGlobalTagAttrTypes   = monitorAttrTypesOf(monitorDispenseGlobalTagAttributes)
	monitorDispenseFilterAttrTypes      = monitorAttrTypesOf(monitorDispenseFilterAttributes)
	monitorDispenseConditionAttrTypes   = monitorAttrTypesOf(monitorDispenseConditionAttributes)
	monitorDispenseProducerAttrTypes    = monitorAttrTypesOf(monitorDispenseProducerAttributes)
	monitorDispenseExtMetricAttrTypes   = monitorAttrTypesOf(monitorDispenseExtMetricAttributes)
	monitorDispenseRuleAttrTypes        = monitorAttrTypesOf(monitorDispenseRuleAttributes)
	monitorDispenseGlobalTagObjectType  = types.ObjectType{AttrTypes: monitorDispenseGlobalTagAttrTypes}
	monitorDispenseFilterObjectType     = types.ObjectType{AttrTypes: monitorDispenseFilterAttrTypes}
	monitorDispenseExtMetricObjectType  = types.ObjectType{AttrTypes: monitorDispenseExtMetricAttrTypes}
	monitorDispenseConditionObjectType  = types.ObjectType{AttrTypes: monitorDispenseConditionAttrTypes}
	monitorDispenseProducerObjectType   = types.ObjectType{AttrTypes: monitorDispenseProducerAttrTypes}
	monitorDispenseConditionListNull    = types.ListNull(monitorDispenseConditionObjectType)
	monitorDispenseGlobalTagListNull    = types.ListNull(monitorDispenseGlobalTagObjectType)
	monitorDispenseExtMetricListNull    = types.ListNull(monitorDispenseExtMetricObjectType)
	monitorDispenseFilterListNull       = types.ListNull(monitorDispenseFilterObjectType)
	monitorDispenseStringSetNull        = types.SetNull(types.StringType)
	monitorDispenseInt64SetNull         = types.SetNull(types.Int64Type)
	monitorDispenseNullProducerObject   = types.ObjectNull(monitorDispenseProducerAttrTypes)
	monitorDispenseNullConditionObject  = types.ObjectNull(monitorDispenseConditionAttrTypes)
)

// monitorAttrTypesOf mirrors the value types of a nested attribute map so the
// conversion helpers stay in sync with the schema definition.
func monitorAttrTypesOf(attributes map[string]schema.Attribute) map[string]attr.Type {
	result := make(map[string]attr.Type, len(attributes))
	for name, attribute := range attributes {
		result[name] = attribute.GetType()
	}
	return result
}

func (r *MonitorDispenseExternalRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_monitor_dispense_external_rule"
}

func (r *MonitorDispenseExternalRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a monitor dispense external rule (dispense rule 3.0).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the resource, which is the string form of `rule_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rule_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Dispense rule id assigned by the cloud side on creation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Rule name.",
			},
			"ext_namespace": schema.StringAttribute{
				Required:    true,
				Description: "Cloud monitor external namespace.",
			},
			"producer": schema.SingleNestedAttribute{
				Required:    true,
				Attributes:  monitorDispenseProducerAttributes,
				Description: "Dispense target producer information.",
			},
			"dispense_regions": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Dispense deployment region list.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"ext_metrics": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Cloud monitor external metric name list. The describe API reports metric details " +
					"through `rule.ext_metric` only; the name-to-detail mapping is not confirmed, so this " +
					"attribute is not refreshed from the response and drift on it is not detected.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"period": schema.SetAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Metric statistical period list.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"dispense_conditions": schema.SingleNestedAttribute{
				Optional:    true,
				Attributes:  monitorDispenseConditionAttributes,
				Description: "Dispense filter condition information.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"rule": schema.SingleNestedAttribute{
				Computed:    true,
				Attributes:  monitorDispenseRuleAttributes,
				Description: "Rule detail reported by the describe API.",
			},
		},
	}
}

func (r *MonitorDispenseExternalRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MonitorDispenseExternalRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, diags := r.buildRuleRequest(ctx, plan.Name, plan.ExtNamespace, plan.Producer,
		plan.DispenseRegions, plan.ExtMetrics, plan.Period, plan.DispenseConditions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewMonitorService(r.Client())
	ruleId, err := service.CreateDispenseExternalRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating monitor dispense external rule",
			fmt.Sprintf("CreateDispenseExternalRule failed: %v", err),
		)
		return
	}

	plan.Id = types.StringValue(strconv.FormatInt(ruleId, 10))
	plan.RuleId = types.Int64Value(ruleId)
	plan.Rule = types.ObjectNull(monitorDispenseRuleAttrTypes)

	// Persist the identity before reading back: a later failure still leaves the
	// resource in state so `terraform destroy` can delete it.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := service.DescribeDispenseExternalRuleById(ctx, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading monitor dispense external rule after create",
			fmt.Sprintf("DescribeDispenseExternalRule failed for rule %d: %v", ruleId, err),
		)
		return
	}
	if rule == nil {
		resp.Diagnostics.AddError(
			"Monitor dispense external rule not readable after create",
			fmt.Sprintf("rule %d was created but cannot be read back", ruleId),
		)
		return
	}

	resp.Diagnostics.Append(refreshRuleBlock(ctx, &plan, rule)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and the post-apply state refresh. Only the
// computed `rule` block is refreshed; writable attributes keep their state
// values so masked or transformed echoes never create diffs.
func (r *MonitorDispenseExternalRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MonitorDispenseExternalRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleId, diags := monitorDispenseRuleIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewMonitorService(r.Client())
	rule, err := service.DescribeDispenseExternalRuleById(ctx, ruleId)
	if err != nil {
		// The API documentation does not define a "not found" error code, so a
		// failed describe is treated as an error: state is not dropped on a
		// failed request.
		resp.Diagnostics.AddError(
			"Error reading monitor dispense external rule",
			fmt.Sprintf("DescribeDispenseExternalRule failed for rule %d: %v", ruleId, err),
		)
		return
	}
	if rule == nil {
		resp.Diagnostics.AddWarning(
			"Monitor dispense external rule not found",
			fmt.Sprintf("rule %d does not exist any more, removing it from state.", ruleId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshRuleBlock(ctx, &state, rule)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates through ModifyDispenseExternalRule.
func (r *MonitorDispenseExternalRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MonitorDispenseExternalRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state MonitorDispenseExternalRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Removing an optional attribute from the configuration makes it Null in the
	// plan. Whether the API treats an omitted list as "keep" or "clear" is not
	// documented, so instead of guessing (and either silently clearing cloud
	// data or producing a permanent diff) the update is rejected.
	resp.Diagnostics.Append(monitorDispenseCheckClearSemantics(&plan, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleId, diags := monitorDispenseRuleIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, diags := r.buildRuleRequest(ctx, plan.Name, plan.ExtNamespace, plan.Producer,
		plan.DispenseRegions, plan.ExtMetrics, plan.Period, plan.DispenseConditions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.RuleId = helper.Int64(ruleId)

	service := NewMonitorService(r.Client())
	if err := service.ModifyDispenseExternalRule(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating monitor dispense external rule",
			fmt.Sprintf("ModifyDispenseExternalRule failed for rule %d: %v", ruleId, err),
		)
		return
	}

	rule, err := service.DescribeDispenseExternalRuleById(ctx, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading monitor dispense external rule after update",
			fmt.Sprintf("DescribeDispenseExternalRule failed for rule %d: %v", ruleId, err),
		)
		return
	}
	if rule == nil {
		resp.Diagnostics.AddError(
			"Monitor dispense external rule not readable after update",
			fmt.Sprintf("rule %d was updated but cannot be read back", ruleId),
		)
		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshRuleBlock(ctx, &plan, rule)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`. It deletes only the dispense rule
// identified by `rule_id`; the effect on the Kafka side is not documented.
func (r *MonitorDispenseExternalRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MonitorDispenseExternalRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleId, diags := monitorDispenseRuleIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewMonitorService(r.Client())
	if err := service.DeleteDispenseExternalRule(ctx, ruleId); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting monitor dispense external rule",
			fmt.Sprintf("DeleteDispenseExternalRule failed for rule %d: %v", ruleId, err),
		)
		return
	}
}

// buildRuleRequest maps the writable attributes onto a Create or Modify
// request. A null or unknown value is transmitted as "absent".
func (r *MonitorDispenseExternalRuleResource) buildRuleRequest(ctx context.Context,
	name, extNamespace types.String, producer types.Object,
	dispenseRegions, extMetrics, period types.Set, dispenseConditions types.Object,
) (*monitorv20230616.CreateDispenseExternalRuleRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	request := monitorv20230616.NewCreateDispenseExternalRuleRequest()
	request.Name = name.ValueStringPointer()
	request.ExtNamespace = extNamespace.ValueStringPointer()

	producerRequest, producerDiags := monitorDispenseExpandProducer(ctx, producer)
	diags.Append(producerDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.Producer = producerRequest

	regions, regionDiags := monitorDispenseExpandStringSet(ctx, dispenseRegions)
	diags.Append(regionDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.DispenseRegions = regions

	metrics, metricDiags := monitorDispenseExpandStringSet(ctx, extMetrics)
	diags.Append(metricDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.ExtMetrics = metrics

	periods, periodDiags := monitorDispenseExpandInt64Set(ctx, period)
	diags.Append(periodDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.Period = periods

	conditions, conditionDiags := monitorDispenseExpandConditions(ctx, dispenseConditions)
	diags.Append(conditionDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.DispenseConditions = conditions

	return request, diags
}

// monitorDispenseCheckClearSemantics rejects updates that removed optional
// attributes from the configuration, because the keep-vs-clear contract of the
// modify API is undocumented.
func monitorDispenseCheckClearSemantics(plan, state *MonitorDispenseExternalRuleResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	reject := func(field string) {
		diags.AddError(
			fmt.Sprintf("Clear semantics of %s is not confirmed", field),
			fmt.Sprintf("Removing %s from the configuration requires knowing whether the "+
				"ModifyDispenseExternalRule API treats an omitted value as \"keep\" or \"clear\". "+
				"Keep the attribute in the configuration or recreate the resource instead.", field),
		)
	}

	// Top-level optional sets.
	if monitorDispenseValueRemoved(plan.DispenseRegions, state.DispenseRegions) {
		reject("dispense_regions")
	}
	if monitorDispenseValueRemoved(plan.ExtMetrics, state.ExtMetrics) {
		reject("ext_metrics")
	}
	if monitorDispenseValueRemoved(plan.Period, state.Period) {
		reject("period")
	}
	if plan.DispenseConditions.IsNull() && !state.DispenseConditions.IsNull() {
		reject("dispense_conditions")
	}

	// Optional producer sub fields.
	if !plan.Producer.IsNull() && !state.Producer.IsNull() {
		planProducer, planDiags := monitorDispenseExpandProducer(context.Background(), plan.Producer)
		if planDiags.HasError() {
			diags.Append(planDiags...)
			return diags
		}
		stateProducer, stateDiags := monitorDispenseExpandProducer(context.Background(), state.Producer)
		if stateDiags.HasError() {
			diags.Append(stateDiags...)
			return diags
		}
		if planProducer != nil && stateProducer != nil {
			if planProducer.Merge == nil && stateProducer.Merge != nil {
				reject("producer.merge")
			}
			if planProducer.Username == nil && stateProducer.Username != nil {
				reject("producer.username")
			}
			if planProducer.Password == nil && stateProducer.Password != nil {
				reject("producer.password")
			}
			if len(planProducer.GlobalTags) == 0 && len(stateProducer.GlobalTags) > 0 {
				reject("producer.global_tags")
			}
			if len(planProducer.DefaultTags) == 0 && len(stateProducer.DefaultTags) > 0 {
				reject("producer.default_tags")
			}
		}
	}

	return diags
}

// monitorDispenseValueRemoved reports whether an optional set was dropped from
// the configuration (Null or empty in the plan but present in state).
func monitorDispenseValueRemoved(plan, state types.Set) bool {
	if !plan.IsNull() && !plan.IsUnknown() {
		return false
	}
	if state.IsNull() || state.IsUnknown() {
		return false
	}
	return true
}

// monitorDispenseRuleIdOf returns the numeric rule id stored in the model.
func monitorDispenseRuleIdOf(model *MonitorDispenseExternalRuleResourceModel) (int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !model.RuleId.IsNull() && !model.RuleId.IsUnknown() {
		return model.RuleId.ValueInt64(), diags
	}
	if !model.Id.IsNull() && !model.Id.IsUnknown() {
		ruleId, err := strconv.ParseInt(model.Id.ValueString(), 10, 64)
		if err != nil {
			diags.AddError(
				"Invalid monitor dispense external rule state",
				fmt.Sprintf("the resource id %q is not a valid rule id: %v", model.Id.ValueString(), err),
			)
			return 0, diags
		}
		return ruleId, diags
	}

	diags.AddError(
		"Invalid monitor dispense external rule state",
		"The resource state does not contain a rule id.",
	)
	return 0, diags
}

// refreshRuleBlock copies the cloud-reported rule into the computed `rule`
// block of the model. Writable attributes are left untouched.
func refreshRuleBlock(ctx context.Context, model *MonitorDispenseExternalRuleResourceModel, rule *monitorv20230616.Rule) diag.Diagnostics {
	var diags diag.Diagnostics

	if rule == nil {
		model.Rule = types.ObjectNull(monitorDispenseRuleAttrTypes)
		return diags
	}

	ruleModel := monitorDispenseRuleModel{
		RuleId:             monitorInt64Pointer(rule.RuleId),
		Name:               monitorStringPointer(rule.Name),
		ExtNamespace:       monitorStringPointer(rule.ExtNamespace),
		ExtMetric:          monitorDispenseExtMetricListNull,
		Producer:           monitorDispenseNullProducerObject,
		UpdateTime:         monitorInt64Pointer(rule.UpdateTime),
		Status:             monitorInt64Pointer(rule.Status),
		Period:             monitorDispenseInt64SetNull,
		DispenseConditions: monitorDispenseNullConditionObject,
		DispenseRegions:    monitorDispenseStringSetNull,
	}

	extMetric, extMetricDiags := monitorDispenseFlattenExtMetrics(ctx, rule.ExtMetric)
	diags.Append(extMetricDiags...)
	if diags.HasError() {
		return diags
	}
	ruleModel.ExtMetric = extMetric

	if rule.Producer != nil {
		producer, producerDiags := monitorDispenseFlattenProducer(ctx, rule.Producer)
		diags.Append(producerDiags...)
		if diags.HasError() {
			return diags
		}
		ruleModel.Producer = producer
	}

	period, periodDiags := monitorDispenseFlattenInt64Set(ctx, rule.Period)
	diags.Append(periodDiags...)
	if diags.HasError() {
		return diags
	}
	ruleModel.Period = period

	conditions, conditionDiags := monitorDispenseFlattenConditions(ctx, rule.DispenseConditions)
	diags.Append(conditionDiags...)
	if diags.HasError() {
		return diags
	}
	ruleModel.DispenseConditions = conditions

	regions, regionDiags := monitorDispenseFlattenStringSet(ctx, rule.DispenseRegions)
	diags.Append(regionDiags...)
	if diags.HasError() {
		return diags
	}
	ruleModel.DispenseRegions = regions

	value, valueDiags := types.ObjectValueFrom(ctx, monitorDispenseRuleAttrTypes, ruleModel)
	diags.Append(valueDiags...)
	if diags.HasError() {
		return diags
	}
	model.Rule = value

	return diags
}

func monitorStringPointer(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(*value)
}

func monitorInt64Pointer(value *int64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*value)
}

// monitorDispenseExpandStringSet converts a string set into the SDK pointer
// slice form. A null or unknown set is transmitted as "absent".
func monitorDispenseExpandStringSet(ctx context.Context, value types.Set) ([]*string, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var elements []string
	diags := value.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(elements) == 0 {
		return nil, diags
	}

	return helper.Strings(elements), diags
}

func monitorDispenseFlattenStringSet(ctx context.Context, values []*string) (types.Set, diag.Diagnostics) {
	if len(values) == 0 {
		return monitorDispenseStringSetNull, nil
	}
	return types.SetValueFrom(ctx, types.StringType, helper.PStrings(values))
}

func monitorDispenseExpandInt64Set(ctx context.Context, value types.Set) ([]*int64, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var elements []int64
	diags := value.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(elements) == 0 {
		return nil, diags
	}

	result := make([]*int64, 0, len(elements))
	for _, element := range elements {
		result = append(result, helper.Int64(element))
	}
	return result, diags
}

func monitorDispenseFlattenInt64Set(ctx context.Context, values []*int64) (types.Set, diag.Diagnostics) {
	if len(values) == 0 {
		return monitorDispenseInt64SetNull, nil
	}

	elements := make([]int64, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		elements = append(elements, *value)
	}
	if len(elements) == 0 {
		return monitorDispenseInt64SetNull, nil
	}
	return types.SetValueFrom(ctx, types.Int64Type, elements)
}

func monitorDispenseExpandProducer(ctx context.Context, object types.Object) (*monitorv20230616.Producer, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model monitorDispenseProducerModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	producer := &monitorv20230616.Producer{}
	if !model.ProtocolType.IsNull() && !model.ProtocolType.IsUnknown() {
		producer.ProtocolType = helper.Int64(model.ProtocolType.ValueInt64())
	}
	if !model.Type.IsNull() && !model.Type.IsUnknown() {
		producer.Type = model.Type.ValueStringPointer()
	}
	if !model.Brokers.IsNull() && !model.Brokers.IsUnknown() {
		producer.Brokers = model.Brokers.ValueStringPointer()
	}
	if !model.Topic.IsNull() && !model.Topic.IsUnknown() {
		producer.Topic = model.Topic.ValueStringPointer()
	}
	if !model.Merge.IsNull() && !model.Merge.IsUnknown() {
		producer.Merge = helper.Int64(model.Merge.ValueInt64())
	}
	if !model.Username.IsNull() && !model.Username.IsUnknown() {
		producer.Username = model.Username.ValueStringPointer()
	}
	if !model.Password.IsNull() && !model.Password.IsUnknown() {
		producer.Password = model.Password.ValueStringPointer()
	}

	globalTags, globalTagsDiags := monitorDispenseExpandGlobalTags(ctx, model.GlobalTags)
	diags.Append(globalTagsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	producer.GlobalTags = globalTags

	defaultTags, defaultTagsDiags := monitorDispenseExpandStringSet(ctx, model.DefaultTags)
	diags.Append(defaultTagsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	producer.DefaultTags = defaultTags

	return producer, diags
}

func monitorDispenseFlattenProducer(ctx context.Context, item *monitorv20230616.Producer) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := monitorDispenseProducerModel{
		ProtocolType: monitorInt64Pointer(item.ProtocolType),
		Type:         monitorStringPointer(item.Type),
		Brokers:      monitorStringPointer(item.Brokers),
		Topic:        monitorStringPointer(item.Topic),
		Merge:        monitorInt64Pointer(item.Merge),
		GlobalTags:   monitorDispenseGlobalTagListNull,
		DefaultTags:  monitorDispenseStringSetNull,
		Username:     monitorStringPointer(item.Username),
		Password:     monitorStringPointer(item.Password),
	}

	globalTags, globalTagsDiags := monitorDispenseFlattenGlobalTags(ctx, item.GlobalTags)
	diags.Append(globalTagsDiags...)
	if diags.HasError() {
		return monitorDispenseNullProducerObject, diags
	}
	model.GlobalTags = globalTags

	defaultTags, defaultTagsDiags := monitorDispenseFlattenStringSet(ctx, item.DefaultTags)
	diags.Append(defaultTagsDiags...)
	if diags.HasError() {
		return monitorDispenseNullProducerObject, diags
	}
	model.DefaultTags = defaultTags

	value, valueDiags := types.ObjectValueFrom(ctx, monitorDispenseProducerAttrTypes, model)
	diags.Append(valueDiags...)
	if diags.HasError() {
		return monitorDispenseNullProducerObject, diags
	}
	return value, diags
}

func monitorDispenseExpandGlobalTags(ctx context.Context, list types.List) ([]*monitorv20230616.DispenseGlobalTag, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []monitorDispenseGlobalTagModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	tags := make([]*monitorv20230616.DispenseGlobalTag, 0, len(models))
	for _, model := range models {
		tag := &monitorv20230616.DispenseGlobalTag{}
		if !model.Key.IsNull() && !model.Key.IsUnknown() {
			tag.Key = model.Key.ValueStringPointer()
		}
		if !model.Value.IsNull() && !model.Value.IsUnknown() {
			tag.Value = model.Value.ValueStringPointer()
		}
		tags = append(tags, tag)
	}

	return tags, diags
}

func monitorDispenseFlattenGlobalTags(ctx context.Context, items []*monitorv20230616.DispenseGlobalTag) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return monitorDispenseGlobalTagListNull, nil
	}

	models := make([]monitorDispenseGlobalTagModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, monitorDispenseGlobalTagModel{
			Key:   monitorStringPointer(item.Key),
			Value: monitorStringPointer(item.Value),
		})
	}
	if len(models) == 0 {
		return monitorDispenseGlobalTagListNull, nil
	}

	return types.ListValueFrom(ctx, monitorDispenseGlobalTagObjectType, models)
}

func monitorDispenseExpandConditions(ctx context.Context, object types.Object) ([]*monitorv20230616.DispenseCondition, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model monitorDispenseConditionModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	condition := &monitorv20230616.DispenseCondition{}
	if !model.ExtMetric.IsNull() && !model.ExtMetric.IsUnknown() {
		condition.ExtMetric = model.ExtMetric.ValueStringPointer()
	}
	if !model.ConditionId.IsNull() && !model.ConditionId.IsUnknown() {
		condition.ConditionId = helper.Int64(model.ConditionId.ValueInt64())
	}

	filters, filterDiags := monitorDispenseExpandFilters(ctx, model.DispenseFilters)
	diags.Append(filterDiags...)
	if diags.HasError() {
		return nil, diags
	}
	condition.DispenseFilters = filters

	return []*monitorv20230616.DispenseCondition{condition}, diags
}

func monitorDispenseFlattenConditions(ctx context.Context, items []*monitorv20230616.DispenseCondition) (types.Object, diag.Diagnostics) {
	if len(items) == 0 {
		return monitorDispenseNullConditionObject, nil
	}

	item := items[0]
	if item == nil {
		return monitorDispenseNullConditionObject, nil
	}

	model := monitorDispenseConditionModel{
		ExtMetric:       monitorStringPointer(item.ExtMetric),
		DispenseFilters: monitorDispenseFilterListNull,
		ConditionId:     monitorInt64Pointer(item.ConditionId),
	}

	filters, filterDiags := monitorDispenseFlattenFilters(ctx, item.DispenseFilters)
	if filterDiags.HasError() {
		return monitorDispenseNullConditionObject, filterDiags
	}
	model.DispenseFilters = filters

	value, valueDiags := types.ObjectValueFrom(ctx, monitorDispenseConditionAttrTypes, model)
	if valueDiags.HasError() {
		return monitorDispenseNullConditionObject, valueDiags
	}
	return value, nil
}

func monitorDispenseExpandFilters(ctx context.Context, list types.List) ([]*monitorv20230616.DispenseFilter, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []monitorDispenseFilterModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	filters := make([]*monitorv20230616.DispenseFilter, 0, len(models))
	for _, model := range models {
		filter := &monitorv20230616.DispenseFilter{}
		if !model.Key.IsNull() && !model.Key.IsUnknown() {
			filter.Key = model.Key.ValueStringPointer()
		}
		if !model.Expression.IsNull() && !model.Expression.IsUnknown() {
			filter.Expression = model.Expression.ValueStringPointer()
		}
		values, valuesDiags := monitorDispenseExpandStringSet(ctx, model.Values)
		diags.Append(valuesDiags...)
		if diags.HasError() {
			return nil, diags
		}
		filter.Values = values
		filters = append(filters, filter)
	}

	return filters, diags
}

func monitorDispenseFlattenFilters(ctx context.Context, items []*monitorv20230616.DispenseFilter) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return monitorDispenseFilterListNull, nil
	}

	var diags diag.Diagnostics
	models := make([]monitorDispenseFilterModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		values, valuesDiags := monitorDispenseFlattenStringSet(ctx, item.Values)
		diags.Append(valuesDiags...)
		if diags.HasError() {
			return monitorDispenseFilterListNull, diags
		}
		models = append(models, monitorDispenseFilterModel{
			Key:        monitorStringPointer(item.Key),
			Values:     values,
			Expression: monitorStringPointer(item.Expression),
		})
	}
	if len(models) == 0 {
		return monitorDispenseFilterListNull, diags
	}

	list, listDiags := types.ListValueFrom(ctx, monitorDispenseFilterObjectType, models)
	diags.Append(listDiags...)
	if diags.HasError() {
		return monitorDispenseFilterListNull, diags
	}
	return list, diags
}

func monitorDispenseFlattenExtMetrics(ctx context.Context, items []*monitorv20230616.ExtMetric) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return monitorDispenseExtMetricListNull, nil
	}

	models := make([]monitorDispenseExtMetricModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, monitorDispenseExtMetricModel{
			MetricName:    monitorStringPointer(item.MetricName),
			MetricCName:   monitorStringPointer(item.MetricCName),
			CNMeaning:     monitorStringPointer(item.CNMeaning),
			EnMeaning:     monitorStringPointer(item.EnMeaning),
			Unit:          monitorStringPointer(item.Unit),
			DimensionFlag: monitorBoolPointer(item.DimensionFlag),
		})
	}
	if len(models) == 0 {
		return monitorDispenseExtMetricListNull, nil
	}

	return types.ListValueFrom(ctx, monitorDispenseExtMetricObjectType, models)
}

func monitorBoolPointer(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

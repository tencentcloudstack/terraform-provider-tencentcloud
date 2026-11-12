package monitor_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	monitorv20230616 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
)

const (
	testDispenseRuleId   = int64(12345)
	testDispenseRuleName = "tf-rule"
)

func testInt64Pointer(value int64) *int64 { return &value }

func testStringPointer(value string) *string { return &value }

func testBoolPointer(value bool) *bool { return &value }

// monitorDispenseTestContext wires a resource instance to a patched SDK client.
type monitorDispenseTestContext struct {
	resource resource.Resource
	client   *monitorv20230616.Client
	schema   schema.Schema
}

func newMonitorDispenseTestContext(t *testing.T, patches *gomonkey.Patches) *monitorDispenseTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	monitorClient := &monitorv20230616.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseMonitorV20230616Client", monitorClient)

	res := svcmonitor.NewMonitorDispenseExternalRuleResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}

	var configureResp resource.ConfigureResponse
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}},
		&configureResp,
	)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("failed to configure resource: %v", configureResp.Diagnostics.Errors())
	}

	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	return &monitorDispenseTestContext{resource: res, client: monitorClient, schema: schemaResp.Schema}
}

func (c *monitorDispenseTestContext) plan(t *testing.T, model svcmonitor.MonitorDispenseExternalRuleResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *monitorDispenseTestContext) state(t *testing.T, model svcmonitor.MonitorDispenseExternalRuleResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func (c *monitorDispenseTestContext) attrTypes(t *testing.T, name string) map[string]attr.Type {
	t.Helper()

	nested, ok := c.schema.Attributes[name].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a single nested attribute", name)
	}
	objectType, ok := nested.GetType().(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object type", name)
	}
	return objectType.AttrTypes
}

func (c *monitorDispenseTestContext) listElementAttrTypes(t *testing.T, parent map[string]attr.Type, name string) map[string]attr.Type {
	t.Helper()

	listType, ok := parent[name].(types.ListType)
	if !ok {
		t.Fatalf("attribute %q does not have a list type", name)
	}
	objectType, ok := listType.ElemType.(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object element type", name)
	}
	return objectType.AttrTypes
}

func mustObjectValue(t *testing.T, attributeTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	t.Helper()

	object, diags := types.ObjectValue(attributeTypes, values)
	if diags.HasError() {
		t.Fatalf("failed to build object value: %v", diags.Errors())
	}
	return object
}

func mustListValue(t *testing.T, elementTypes map[string]attr.Type, elements ...map[string]attr.Value) types.List {
	t.Helper()

	objectType := types.ObjectType{AttrTypes: elementTypes}
	values := make([]attr.Value, 0, len(elements))
	for _, element := range elements {
		object, diags := types.ObjectValue(elementTypes, element)
		if diags.HasError() {
			t.Fatalf("failed to build list element: %v", diags.Errors())
		}
		values = append(values, object)
	}

	list, diags := types.ListValue(objectType, values)
	if diags.HasError() {
		t.Fatalf("failed to build list value: %v", diags.Errors())
	}
	return list
}

func mustStringSet(t *testing.T, values ...string) types.Set {
	t.Helper()

	set, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		t.Fatalf("failed to build string set: %v", diags.Errors())
	}
	return set
}

func mustInt64Set(t *testing.T, values ...int64) types.Set {
	t.Helper()

	set, diags := types.SetValueFrom(context.Background(), types.Int64Type, values)
	if diags.HasError() {
		t.Fatalf("failed to build int64 set: %v", diags.Errors())
	}
	return set
}

// sampleModel returns a model with every configurable leaf populated.
func (c *monitorDispenseTestContext) sampleModel(t *testing.T) svcmonitor.MonitorDispenseExternalRuleResourceModel {
	t.Helper()

	producerTypes := c.attrTypes(t, "producer")
	conditionTypes := c.attrTypes(t, "dispense_conditions")
	tagTypes := c.listElementAttrTypes(t, producerTypes, "global_tags")
	filterTypes := c.listElementAttrTypes(t, conditionTypes, "dispense_filters")

	return svcmonitor.MonitorDispenseExternalRuleResourceModel{
		Id:           types.StringNull(),
		RuleId:       types.Int64Null(),
		Name:         types.StringValue(testDispenseRuleName),
		ExtNamespace: types.StringValue("ns-tf"),
		Producer: mustObjectValue(t, producerTypes, map[string]attr.Value{
			"protocol_type": types.Int64Value(0),
			"type":          types.StringValue("kafka"),
			"brokers":       types.StringValue("broker-1:9092"),
			"topic":         types.StringValue("topic-1"),
			"merge":         types.Int64Value(1),
			"global_tags": mustListValue(t, tagTypes, map[string]attr.Value{
				"key":   types.StringValue("env"),
				"value": types.StringValue("prod"),
			}),
			"default_tags": mustStringSet(t, "dim-1"),
			"username":     types.StringValue("user"),
			"password":     types.StringValue("pass"),
		}),
		DispenseRegions: mustStringSet(t, "ap-guangzhou"),
		ExtMetrics:      mustStringSet(t, "metric-1"),
		Period:          mustInt64Set(t, 60),
		DispenseConditions: mustObjectValue(t, conditionTypes, map[string]attr.Value{
			"ext_metric": types.StringValue("metric-1"),
			"dispense_filters": mustListValue(t, filterTypes, map[string]attr.Value{
				"key":        types.StringValue("dim"),
				"values":     mustStringSet(t, "v1", "v2"),
				"expression": types.StringValue("eq"),
			}),
			"condition_id": types.Int64Null(),
		}),
		Rule: types.ObjectNull(c.attrTypes(t, "rule")),
	}
}

// testProducer struct mirrors the producer object for state assertions.
type testProducer struct {
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

// testRule struct mirrors the computed rule object for state assertions.
type testRule struct {
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

// patchCreate patches the SDK create method and captures the request.
func patchCreate(patches *gomonkey.Patches, captured **monitorv20230616.CreateDispenseExternalRuleRequest, err error) *gomonkey.Patches {
	response := &monitorv20230616.CreateDispenseExternalRuleResponse{
		Response: &monitorv20230616.CreateDispenseExternalRuleResponseParams{
			RuleId: testInt64Pointer(testDispenseRuleId),
		},
	}
	return patches.ApplyMethodReplace(monitorClientReceiver(), "CreateDispenseExternalRule",
		func(_ *monitorv20230616.Client, request *monitorv20230616.CreateDispenseExternalRuleRequest) (*monitorv20230616.CreateDispenseExternalRuleResponse, error) {
			if captured != nil {
				*captured = request
			}
			return response, err
		})
}

// monitorClientReceiver returns a throwaway receiver for gomonkey method
// patching; only the method symbol matters, not the instance.
func monitorClientReceiver() *monitorv20230616.Client {
	return &monitorv20230616.Client{}
}

// patchDescribe patches the SDK describe method with the given rule (nil means gone).
func patchDescribe(patches *gomonkey.Patches, rule *monitorv20230616.Rule, err error) *gomonkey.Patches {
	response := &monitorv20230616.DescribeDispenseExternalRuleResponse{
		Response: &monitorv20230616.DescribeDispenseExternalRuleResponseParams{Rule: rule},
	}
	return patches.ApplyMethodReplace(monitorClientReceiver(), "DescribeDispenseExternalRule",
		func(_ *monitorv20230616.Client, _ *monitorv20230616.DescribeDispenseExternalRuleRequest) (*monitorv20230616.DescribeDispenseExternalRuleResponse, error) {
			return response, err
		})
}


// TestMonitorDispenseExternalRuleImportState verifies the id only import path.
func TestMonitorDispenseExternalRuleImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMonitorDispenseTestContext(t, patches)

	importer, ok := c.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: c.schema}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: "12345"},
		importResp,
	)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("import failed: %v", importResp.Diagnostics.Errors())
	}

	var state svcmonitor.MonitorDispenseExternalRuleResourceModel
	if diags := importResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read imported state: %v", diags.Errors())
	}
	if state.Id.ValueString() != "12345" {
		t.Fatalf("imported id is %q, want %q", state.Id.ValueString(), "12345")
	}
}


func sampleCloudRule() *monitorv20230616.Rule {
	return &monitorv20230616.Rule{
		RuleId:       testInt64Pointer(testDispenseRuleId),
		Name:         testStringPointer(testDispenseRuleName),
		ExtNamespace: testStringPointer("ns-tf"),
		ExtMetric: []*monitorv20230616.ExtMetric{
			{
				MetricName:    testStringPointer("metric-1"),
				MetricCName:   testStringPointer("metric-1-cn"),
				CNMeaning:     testStringPointer("cn"),
				EnMeaning:     testStringPointer("en"),
				Unit:          testStringPointer("ms"),
				DimensionFlag: testBoolPointer(true),
			},
		},
		Producer: &monitorv20230616.Producer{
			ProtocolType: testInt64Pointer(0),
			Type:         testStringPointer("kafka"),
			Brokers:      testStringPointer("broker-1:9092"),
			Topic:        testStringPointer("topic-1"),
			Merge:        testInt64Pointer(1),
			GlobalTags: []*monitorv20230616.DispenseGlobalTag{
				{Key: testStringPointer("env"), Value: testStringPointer("prod")},
			},
			DefaultTags: []*string{testStringPointer("dim-1")},
			Username:    testStringPointer("user"),
			Password:    testStringPointer("pass"),
		},
		UpdateTime: testInt64Pointer(1700000000),
		Status:     testInt64Pointer(1),
		Period:     []*int64{testInt64Pointer(60)},
		DispenseConditions: []*monitorv20230616.DispenseCondition{
			{
				ExtMetric: testStringPointer("metric-1"),
				DispenseFilters: []*monitorv20230616.DispenseFilter{
					{
						Key:        testStringPointer("dim"),
						Values:     []*string{testStringPointer("v1"), testStringPointer("v2")},
						Expression: testStringPointer("eq"),
					},
				},
				ConditionId: testInt64Pointer(7),
			},
		},
		DispenseRegions: []*string{testStringPointer("ap-guangzhou")},
	}
}

func TestMonitorDispenseExternalRuleCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *monitorv20230616.CreateDispenseExternalRuleRequest
	patchCreate(patches, &captured, nil)
	patchDescribe(patches, sampleCloudRule(), nil)

	c := newMonitorDispenseTestContext(t, patches)
	model := c.sampleModel(t)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Create(context.Background(), resource.CreateRequest{Plan: c.plan(t, model)}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateDispenseExternalRule was not called")
	}
	if captured.Name == nil || *captured.Name != testDispenseRuleName {
		t.Errorf("request.Name = %v, want %q", captured.Name, testDispenseRuleName)
	}
	if captured.ExtNamespace == nil || *captured.ExtNamespace != "ns-tf" {
		t.Errorf("request.ExtNamespace = %v, want %q", captured.ExtNamespace, "ns-tf")
	}
	if captured.Producer == nil {
		t.Fatal("request.Producer is nil")
	}
	if captured.Producer.ProtocolType == nil || *captured.Producer.ProtocolType != 0 {
		t.Errorf("request.Producer.ProtocolType = %v, want 0", captured.Producer.ProtocolType)
	}
	if captured.Producer.Brokers == nil || *captured.Producer.Brokers != "broker-1:9092" {
		t.Errorf("request.Producer.Brokers = %v, want %q", captured.Producer.Brokers, "broker-1:9092")
	}
	if len(captured.DispenseRegions) != 1 || *captured.DispenseRegions[0] != "ap-guangzhou" {
		t.Errorf("request.DispenseRegions = %v, want [ap-guangzhou]", captured.DispenseRegions)
	}
	if len(captured.ExtMetrics) != 1 || *captured.ExtMetrics[0] != "metric-1" {
		t.Errorf("request.ExtMetrics = %v, want [metric-1]", captured.ExtMetrics)
	}
	if len(captured.Period) != 1 || *captured.Period[0] != 60 {
		t.Errorf("request.Period = %v, want [60]", captured.Period)
	}
	if len(captured.DispenseConditions) != 1 || captured.DispenseConditions[0].ExtMetric == nil || *captured.DispenseConditions[0].ExtMetric != "metric-1" {
		t.Errorf("request.DispenseConditions = %v, want one condition with metric-1", captured.DispenseConditions)
	}

	var result svcmonitor.MonitorDispenseExternalRuleResourceModel
	if diags := createResp.State.Get(context.Background(), &result); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if result.Id.ValueString() != "12345" {
		t.Errorf("state id = %q, want %q", result.Id.ValueString(), "12345")
	}
	if result.RuleId.ValueInt64() != testDispenseRuleId {
		t.Errorf("state rule_id = %d, want %d", result.RuleId.ValueInt64(), testDispenseRuleId)
	}

	var rule testRule
	if diags := result.Rule.As(context.Background(), &rule, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("failed to read rule block: %v", diags.Errors())
	}
	if rule.RuleId.ValueInt64() != testDispenseRuleId {
		t.Errorf("rule.rule_id = %d, want %d", rule.RuleId.ValueInt64(), testDispenseRuleId)
	}
	if rule.Name.ValueString() != testDispenseRuleName {
		t.Errorf("rule.name = %q, want %q", rule.Name.ValueString(), testDispenseRuleName)
	}
	if rule.ExtNamespace.ValueString() != "ns-tf" {
		t.Errorf("rule.ext_namespace = %q, want %q", rule.ExtNamespace.ValueString(), "ns-tf")
	}
	if rule.UpdateTime.ValueInt64() != 1700000000 {
		t.Errorf("rule.update_time = %d, want 1700000000", rule.UpdateTime.ValueInt64())
	}
	if rule.Status.ValueInt64() != 1 {
		t.Errorf("rule.status = %d, want 1", rule.Status.ValueInt64())
	}
	if len(rule.Period.Elements()) != 1 {
		t.Errorf("rule.period length = %d, want 1", len(rule.Period.Elements()))
	}
	if len(rule.DispenseRegions.Elements()) != 1 {
		t.Errorf("rule.dispense_regions length = %d, want 1", len(rule.DispenseRegions.Elements()))
	}

	var producer testProducer
	if diags := rule.Producer.As(context.Background(), &producer, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("failed to read rule.producer block: %v", diags.Errors())
	}
	if producer.ProtocolType.ValueInt64() != 0 {
		t.Errorf("rule.producer.protocol_type = %d, want 0", producer.ProtocolType.ValueInt64())
	}
	if producer.Topic.ValueString() != "topic-1" {
		t.Errorf("rule.producer.topic = %q, want %q", producer.Topic.ValueString(), "topic-1")
	}
	if producer.Merge.ValueInt64() != 1 {
		t.Errorf("rule.producer.merge = %d, want 1", producer.Merge.ValueInt64())
	}
	if producer.Username.ValueString() != "user" {
		t.Errorf("rule.producer.username = %q, want %q", producer.Username.ValueString(), "user")
	}
	if producer.Password.ValueString() != "pass" {
		t.Errorf("rule.producer.password = %q, want %q", producer.Password.ValueString(), "pass")
	}
	if len(producer.GlobalTags.Elements()) != 1 {
		t.Errorf("rule.producer.global_tags length = %d, want 1", len(producer.GlobalTags.Elements()))
	}
	if len(producer.DefaultTags.Elements()) != 1 {
		t.Errorf("rule.producer.default_tags length = %d, want 1", len(producer.DefaultTags.Elements()))
	}
	if len(rule.ExtMetric.Elements()) != 1 {
		t.Errorf("rule.ext_metric length = %d, want 1", len(rule.ExtMetric.Elements()))
	}
	if len(rule.DispenseConditions.Elements()) != 1 {
		t.Errorf("rule.dispense_conditions length = %d, want 1", len(rule.DispenseConditions.Elements()))
	}
}

func TestMonitorDispenseExternalRuleCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchCreate(patches, nil, fmt.Errorf("api error"))

	c := newMonitorDispenseTestContext(t, patches)
	model := c.sampleModel(t)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Create(context.Background(), resource.CreateRequest{Plan: c.plan(t, model)}, createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("create should report an error when the API call fails")
	}
}

func TestMonitorDispenseExternalRuleReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchDescribe(patches, sampleCloudRule(), nil)

	c := newMonitorDispenseTestContext(t, patches)
	model := c.sampleModel(t)
	model.Id = types.StringValue("12345")
	model.RuleId = types.Int64Value(testDispenseRuleId)
	model.Rule = types.ObjectNull(c.attrTypes(t, "rule"))

	readResp := &resource.ReadResponse{State: c.state(t, model)}
	c.resource.Read(context.Background(), resource.ReadRequest{State: c.state(t, model)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("read failed: %v", readResp.Diagnostics.Errors())
	}

	var result svcmonitor.MonitorDispenseExternalRuleResourceModel
	if diags := readResp.State.Get(context.Background(), &result); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	// Writable attributes keep their state values.
	if result.Name.ValueString() != testDispenseRuleName {
		t.Errorf("state name = %q, want %q", result.Name.ValueString(), testDispenseRuleName)
	}
	if result.Producer.IsNull() {
		t.Fatal("state producer should not be null")
	}

	var rule testRule
	if diags := result.Rule.As(context.Background(), &rule, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("failed to read rule block: %v", diags.Errors())
	}
	if rule.RuleId.ValueInt64() != testDispenseRuleId {
		t.Errorf("rule.rule_id = %d, want %d", rule.RuleId.ValueInt64(), testDispenseRuleId)
	}
	if rule.Status.ValueInt64() != 1 {
		t.Errorf("rule.status = %d, want 1", rule.Status.ValueInt64())
	}
}

func TestMonitorDispenseExternalRuleReadGone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchDescribe(patches, nil, nil)

	c := newMonitorDispenseTestContext(t, patches)
	model := c.sampleModel(t)
	model.Id = types.StringValue("12345")
	model.RuleId = types.Int64Value(testDispenseRuleId)

	readResp := &resource.ReadResponse{State: c.state(t, model)}
	c.resource.Read(context.Background(), resource.ReadRequest{State: c.state(t, model)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("read should not fail when the rule is gone: %v", readResp.Diagnostics.Errors())
	}
}

func TestMonitorDispenseExternalRuleReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchDescribe(patches, nil, fmt.Errorf("api error"))

	c := newMonitorDispenseTestContext(t, patches)
	model := c.sampleModel(t)
	model.Id = types.StringValue("12345")
	model.RuleId = types.Int64Value(testDispenseRuleId)

	readResp := &resource.ReadResponse{State: c.state(t, model)}
	c.resource.Read(context.Background(), resource.ReadRequest{State: c.state(t, model)}, readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("read should report an error when the API call fails")
	}
}

func TestMonitorDispenseExternalRuleUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *monitorv20230616.ModifyDispenseExternalRuleRequest
	response := &monitorv20230616.ModifyDispenseExternalRuleResponse{Response: &monitorv20230616.ModifyDispenseExternalRuleResponseParams{}}
	patches.ApplyMethodReplace(monitorClientReceiver(), "ModifyDispenseExternalRule",
		func(_ *monitorv20230616.Client, request *monitorv20230616.ModifyDispenseExternalRuleRequest) (*monitorv20230616.ModifyDispenseExternalRuleResponse, error) {
			captured = request
			return response, nil
		})
	patchDescribe(patches, sampleCloudRule(), nil)

	c := newMonitorDispenseTestContext(t, patches)
	planModel := c.sampleModel(t)
	planModel.Id = types.StringValue("12345")
	planModel.RuleId = types.Int64Value(testDispenseRuleId)
	stateModel := planModel

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Update(context.Background(), resource.UpdateRequest{Plan: c.plan(t, planModel), State: c.state(t, stateModel)}, updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyDispenseExternalRule was not called")
	}
	if captured.RuleId == nil || *captured.RuleId != testDispenseRuleId {
		t.Errorf("request.RuleId = %v, want %d", captured.RuleId, testDispenseRuleId)
	}
	if captured.Name == nil || *captured.Name != testDispenseRuleName {
		t.Errorf("request.Name = %v, want %q", captured.Name, testDispenseRuleName)
	}

	var result svcmonitor.MonitorDispenseExternalRuleResourceModel
	if diags := updateResp.State.Get(context.Background(), &result); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if result.Id.ValueString() != "12345" {
		t.Errorf("state id = %q, want %q", result.Id.ValueString(), "12345")
	}
	var rule testRule
	if diags := result.Rule.As(context.Background(), &rule, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("failed to read rule block: %v", diags.Errors())
	}
	if rule.RuleId.ValueInt64() != testDispenseRuleId {
		t.Errorf("rule.rule_id = %d, want %d", rule.RuleId.ValueInt64(), testDispenseRuleId)
	}
}

func TestMonitorDispenseExternalRuleUpdateRejectsClearingDispenseRegions(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchDescribe(patches, sampleCloudRule(), nil)

	c := newMonitorDispenseTestContext(t, patches)
	planModel := c.sampleModel(t)
	planModel.Id = types.StringValue("12345")
	planModel.RuleId = types.Int64Value(testDispenseRuleId)
	planModel.DispenseRegions = types.SetNull(types.StringType)
	stateModel := c.sampleModel(t)
	stateModel.Id = types.StringValue("12345")
	stateModel.RuleId = types.Int64Value(testDispenseRuleId)

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Update(context.Background(), resource.UpdateRequest{Plan: c.plan(t, planModel), State: c.state(t, stateModel)}, updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("update should be rejected when an optional attribute is removed from the configuration")
	}
}

func TestMonitorDispenseExternalRuleDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodReplace(monitorClientReceiver(), "DeleteDispenseExternalRule",
		func(_ *monitorv20230616.Client, request *monitorv20230616.DeleteDispenseExternalRuleRequest) (*monitorv20230616.DeleteDispenseExternalRuleResponse, error) {
			if len(request.RuleIdList) != 1 || *request.RuleIdList[0] != testDispenseRuleId {
				t.Errorf("request.RuleIdList = %v, want [%d]", request.RuleIdList, testDispenseRuleId)
			}
			return &monitorv20230616.DeleteDispenseExternalRuleResponse{Response: &monitorv20230616.DeleteDispenseExternalRuleResponseParams{}}, nil
		})

	c := newMonitorDispenseTestContext(t, patches)
	model := c.sampleModel(t)
	model.Id = types.StringValue("12345")
	model.RuleId = types.Int64Value(testDispenseRuleId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Delete(context.Background(), resource.DeleteRequest{State: c.state(t, model)}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("delete failed: %v", deleteResp.Diagnostics.Errors())
	}
}

func TestMonitorDispenseExternalRuleDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodReturn(monitorClientReceiver(), "DeleteDispenseExternalRule",
		(*monitorv20230616.DeleteDispenseExternalRuleResponse)(nil), fmt.Errorf("api error"))

	c := newMonitorDispenseTestContext(t, patches)
	model := c.sampleModel(t)
	model.Id = types.StringValue("12345")
	model.RuleId = types.Int64Value(testDispenseRuleId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Delete(context.Background(), resource.DeleteRequest{State: c.state(t, model)}, deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("delete should report an error when the API call fails")
	}
}

package ioa_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	ioav20220601 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ioa/v20220601"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcioa "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ioa"
)

const (
	ioaTestStrategyId   = "42"
	ioaTestStrategyName = "tf-example"
)

func ioaTestString(value string) *string { return &value }

func ioaTestInt64(value int64) *int64 { return &value }

func ioaTestUint64(value uint64) *uint64 { return &value }

// ioaDeliveryStrategyTestContext wires a resource instance to a patched SDK client.
type ioaDeliveryStrategyTestContext struct {
	resource resource.Resource
	client   *ioav20220601.Client
	schema   schema.Schema
}

func newIoaDeliveryStrategyTestContext(t *testing.T, patches *gomonkey.Patches) *ioaDeliveryStrategyTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	ioaClient := &ioav20220601.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseIoaV20220601Client", ioaClient)

	res := svcioa.NewIoaDeliveryStrategyResource()
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

	return &ioaDeliveryStrategyTestContext{resource: res, client: ioaClient, schema: schemaResp.Schema}
}

func (c *ioaDeliveryStrategyTestContext) plan(t *testing.T, model svcioa.IoaDeliveryStrategyResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *ioaDeliveryStrategyTestContext) state(t *testing.T, model svcioa.IoaDeliveryStrategyResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func ioaTestStringSet(t *testing.T, values ...string) types.Set {
	t.Helper()

	set, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		t.Fatalf("failed to build string set: %v", diags.Errors())
	}
	return set
}

func ioaTestDataList(t *testing.T, item ioav20220601.DeliveryStrategyItem) types.List {
	t.Helper()

	res := svcioa.NewIoaDeliveryStrategyResource()
	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	dataAttr := schemaResp.Schema.Attributes["data"].(schema.ListNestedAttribute)
	dataObjectType := dataAttr.GetType().(types.ListType).ElemType.(types.ObjectType)
	innerAttr := dataObjectType.AttrTypes["data"].(types.ListType)
	detailObjectType := innerAttr.ElemType.(types.ObjectType)

	detail := map[string]attr.Value{
		"name":        types.StringPointerValue(item.Name),
		"type":        types.Int64PointerValue(item.Type),
		"log_types":   ioaTestStringSetFromPointers(t, item.LogTypes),
		"config_name": types.StringPointerValue(item.ConfigName),
		"description": types.StringPointerValue(item.Description),
	}
	if item.ID != nil {
		detail["id"] = types.Int64Value(int64(*item.ID))
	} else {
		detail["id"] = types.Int64Null()
	}
	if item.ConfigId != nil {
		detail["config_id"] = types.Int64Value(int64(*item.ConfigId))
	} else {
		detail["config_id"] = types.Int64Null()
	}

	inner, diags := types.ListValue(detailObjectType, []attr.Value{mustIoaObject(t, detailObjectType.AttrTypes, detail)})
	if diags.HasError() {
		t.Fatalf("failed to build inner data list: %v", diags.Errors())
	}

	outer, diags := types.ListValue(dataObjectType, []attr.Value{
		mustIoaObject(t, dataObjectType.AttrTypes, map[string]attr.Value{"data": inner}),
	})
	if diags.HasError() {
		t.Fatalf("failed to build outer data list: %v", diags.Errors())
	}
	return outer
}

func ioaTestDataObjectType(t *testing.T) types.ObjectType {
	t.Helper()

	res := svcioa.NewIoaDeliveryStrategyResource()
	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}
	dataAttr := schemaResp.Schema.Attributes["data"].(schema.ListNestedAttribute)
	return dataAttr.GetType().(types.ListType).ElemType.(types.ObjectType)
}

func ioaTestNullDataList(t *testing.T) types.List {
	t.Helper()

	return types.ListNull(ioaTestDataObjectType(t))
}

func ioaTestStringSetFromPointers(t *testing.T, values []*string) types.Set {
	t.Helper()

	if len(values) == 0 {
		return types.SetNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringPointerValue(v))
	}
	set, diags := types.SetValue(types.StringType, elements)
	if diags.HasError() {
		t.Fatalf("failed to build string set: %v", diags.Errors())
	}
	return set
}

func mustIoaObject(t *testing.T, attributeTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	t.Helper()

	object, diags := types.ObjectValue(attributeTypes, values)
	if diags.HasError() {
		t.Fatalf("failed to build object value: %v", diags.Errors())
	}
	return object
}

// samplePlanModel returns a plan model with every configurable leaf populated
// and the computed data list left null.
func ioaTestSamplePlanModel(t *testing.T) svcioa.IoaDeliveryStrategyResourceModel {
	t.Helper()

	return svcioa.IoaDeliveryStrategyResourceModel{
		Id:          types.StringNull(),
		Name:        types.StringValue(ioaTestStrategyName),
		LogTypes:    ioaTestStringSet(t, "040201000000", "040203000000"),
		Type:        types.Int64Value(0),
		ConfigId:    types.Int64Value(100),
		Description: types.StringValue("managed by terraform"),
		Enabled:     types.BoolValue(true),
		Data:        ioaTestNullDataList(t),
	}
}

func ioaTestDescribeResponseItem() *ioav20220601.DeliveryStrategyItem {
	return &ioav20220601.DeliveryStrategyItem{
		Name:        ioaTestString(ioaTestStrategyName),
		Type:        ioaTestInt64(0),
		LogTypes:    []*string{ioaTestString("040201000000"), ioaTestString("040203000000")},
		ConfigName:  ioaTestString("test-server"),
		Description: ioaTestString("managed by terraform"),
		ID:          ioaTestUint64(42),
		ConfigId:    ioaTestUint64(100),
	}
}

func ioaTestDescribeResponse(item *ioav20220601.DeliveryStrategyItem) *ioav20220601.DescribeDeliveryStrategyResponse {
	return &ioav20220601.DescribeDeliveryStrategyResponse{
		Response: &ioav20220601.DescribeDeliveryStrategyResponseParams{
			Data: &ioav20220601.DeliveryStrategyData{Data: item},
		},
	}
}

func TestIoaDeliveryStrategyResourceCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)

	modifyResp := &ioav20220601.ModifyDeliveryStrategyResponse{
		Response: &ioav20220601.ModifyDeliveryStrategyResponseParams{
			Data: &ioav20220601.ModifyDeliveryStrategyData{ID: ioaTestUint64(42)},
		},
	}
	var capturedRequest *ioav20220601.ModifyDeliveryStrategyRequest
	patches.ApplyMethod(c.client, "ModifyDeliveryStrategy", func(request *ioav20220601.ModifyDeliveryStrategyRequest) (*ioav20220601.ModifyDeliveryStrategyResponse, error) {
		capturedRequest = request
		return modifyResp, nil
	})
	patches.ApplyMethodReturn(c.client, "DescribeDeliveryStrategy", ioaTestDescribeResponse(ioaTestDescribeResponseItem()), nil)

	plan := c.plan(t, ioaTestSamplePlanModel(t))

	var createResp resource.CreateResponse
	createResp.State = tfsdk.State{Schema: c.schema}
	c.resource.Create(context.Background(), resource.CreateRequest{Plan: plan}, &createResp)

	if createResp.Diagnostics.HasError() {
		t.Fatalf("unexpected create diagnostics: %v", createResp.Diagnostics.Errors())
	}

	// The create request must not carry an ID (create+update combined API).
	assert.Nil(t, capturedRequest.ID)
	assert.Equal(t, ioaTestStringPointer(ioaTestStrategyName), capturedRequest.Name)
	assert.NotNil(t, capturedRequest.Enabled)
	assert.True(t, *capturedRequest.Enabled)

	var state svcioa.IoaDeliveryStrategyResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assert.Equal(t, ioaTestStrategyId, state.Id.ValueString())
	assert.Equal(t, ioaTestStrategyName, state.Name.ValueString())
	assert.Equal(t, int64(0), state.Type.ValueInt64())
	assert.Equal(t, int64(100), state.ConfigId.ValueInt64())
	assert.Equal(t, "managed by terraform", state.Description.ValueString())
	assert.True(t, state.Enabled.ValueBool())
	assert.False(t, state.Data.IsNull())

	elements := state.Data.Elements()
	assert.Len(t, elements, 1)
	inner := elements[0].(types.Object)
	innerAttrs := inner.Attributes()
	innerList := innerAttrs["data"].(types.List)
	assert.Len(t, innerList.Elements(), 1)
	detail := innerList.Elements()[0].(types.Object).Attributes()
	assert.Equal(t, ioaTestStrategyName, detail["name"].(types.String).ValueString())
	assert.Equal(t, int64(42), detail["id"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(100), detail["config_id"].(types.Int64).ValueInt64())
	assert.Equal(t, "test-server", detail["config_name"].(types.String).ValueString())
	assert.Equal(t, int64(0), detail["type"].(types.Int64).ValueInt64())
	assert.Equal(t, "managed by terraform", detail["description"].(types.String).ValueString())
	assert.False(t, detail["log_types"].(types.Set).IsNull())
	assert.Len(t, detail["log_types"].(types.Set).Elements(), 2)
}

func ioaTestStringPointer(value string) *string { return &value }

func TestIoaDeliveryStrategyResourceCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)
	patches.ApplyMethodReturn(c.client, "ModifyDeliveryStrategy", (*ioav20220601.ModifyDeliveryStrategyResponse)(nil), context.DeadlineExceeded)

	plan := c.plan(t, ioaTestSamplePlanModel(t))

	var createResp resource.CreateResponse
	createResp.State = tfsdk.State{Schema: c.schema}
	c.resource.Create(context.Background(), resource.CreateRequest{Plan: plan}, &createResp)

	assert.True(t, createResp.Diagnostics.HasError())
}

func TestIoaDeliveryStrategyResourceCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)
	patches.ApplyMethodReturn(c.client, "ModifyDeliveryStrategy", &ioav20220601.ModifyDeliveryStrategyResponse{
		Response: &ioav20220601.ModifyDeliveryStrategyResponseParams{
			Data: &ioav20220601.ModifyDeliveryStrategyData{ID: ioaTestUint64(0)},
		},
	}, nil)

	plan := c.plan(t, ioaTestSamplePlanModel(t))

	var createResp resource.CreateResponse
	createResp.State = tfsdk.State{Schema: c.schema}
	c.resource.Create(context.Background(), resource.CreateRequest{Plan: plan}, &createResp)

	assert.True(t, createResp.Diagnostics.HasError())
}

func TestIoaDeliveryStrategyResourceRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)
	patches.ApplyMethodReturn(c.client, "DescribeDeliveryStrategy", ioaTestDescribeResponse(ioaTestDescribeResponseItem()), nil)

	model := ioaTestSamplePlanModel(t)
	model.Id = types.StringValue(ioaTestStrategyId)
	model.Enabled = types.BoolValue(false)

	var readResp resource.ReadResponse
	readResp.State = c.state(t, model)
	c.resource.Read(context.Background(), resource.ReadRequest{State: readResp.State}, &readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected read diagnostics: %v", readResp.Diagnostics.Errors())
	}

	var state svcioa.IoaDeliveryStrategyResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assert.Equal(t, ioaTestStrategyId, state.Id.ValueString())
	assert.Equal(t, ioaTestStrategyName, state.Name.ValueString())
	assert.Equal(t, int64(0), state.Type.ValueInt64())
	assert.Equal(t, int64(100), state.ConfigId.ValueInt64())
	assert.Equal(t, "managed by terraform", state.Description.ValueString())
	assert.False(t, state.Data.IsNull())
	// Enabled has no counterpart in the describe response, the state value is kept.
	assert.False(t, state.Enabled.ValueBool())
}

func TestIoaDeliveryStrategyResourceReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)
	patches.ApplyMethodReturn(c.client, "DescribeDeliveryStrategy", ioaTestDescribeResponse(nil), nil)

	model := ioaTestSamplePlanModel(t)
	model.Id = types.StringValue(ioaTestStrategyId)

	var readResp resource.ReadResponse
	readResp.State = c.state(t, model)
	c.resource.Read(context.Background(), resource.ReadRequest{State: readResp.State}, &readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected read diagnostics: %v", readResp.Diagnostics.Errors())
	}
	assert.Nil(t, readResp.State.Raw)
}

func TestIoaDeliveryStrategyResourceReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)
	patches.ApplyMethodReturn(c.client, "DescribeDeliveryStrategy", (*ioav20220601.DescribeDeliveryStrategyResponse)(nil), context.DeadlineExceeded)

	model := ioaTestSamplePlanModel(t)
	model.Id = types.StringValue(ioaTestStrategyId)

	var readResp resource.ReadResponse
	readResp.State = c.state(t, model)
	c.resource.Read(context.Background(), resource.ReadRequest{State: readResp.State}, &readResp)

	assert.True(t, readResp.Diagnostics.HasError())
}

func TestIoaDeliveryStrategyResourceUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)

	var capturedRequest *ioav20220601.ModifyDeliveryStrategyRequest
	patches.ApplyMethod(c.client, "ModifyDeliveryStrategy", func(request *ioav20220601.ModifyDeliveryStrategyRequest) (*ioav20220601.ModifyDeliveryStrategyResponse, error) {
		capturedRequest = request
		return &ioav20220601.ModifyDeliveryStrategyResponse{
			Response: &ioav20220601.ModifyDeliveryStrategyResponseParams{
				Data: &ioav20220601.ModifyDeliveryStrategyData{ID: ioaTestUint64(42)},
			},
		}, nil
	})
	updated := ioaTestDescribeResponseItem()
	updated.Name = ioaTestString("tf-example-updated")
	patches.ApplyMethodReturn(c.client, "DescribeDeliveryStrategy", ioaTestDescribeResponse(updated), nil)

	plan := ioaTestSamplePlanModel(t)
	plan.Id = types.StringValue(ioaTestStrategyId)
	plan.Name = types.StringValue("tf-example-updated")
	plan.Data = ioaTestNullDataList(t)
	state := c.state(t, plan)

	var updateResp resource.UpdateResponse
	updateResp.State = tfsdk.State{Schema: c.schema}
	c.resource.Update(context.Background(), resource.UpdateRequest{Plan: c.plan(t, plan), State: state}, &updateResp)

	if updateResp.Diagnostics.HasError() {
		t.Fatalf("unexpected update diagnostics: %v", updateResp.Diagnostics.Errors())
	}

	// The update request must carry the strategy id from state.
	assert.NotNil(t, capturedRequest.ID)
	assert.Equal(t, int64(42), *capturedRequest.ID)
	assert.Equal(t, ioaTestStringPointer("tf-example-updated"), capturedRequest.Name)

	var newState svcioa.IoaDeliveryStrategyResourceModel
	if diags := updateResp.State.Get(context.Background(), &newState); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assert.Equal(t, ioaTestStrategyId, newState.Id.ValueString())
	assert.Equal(t, "tf-example-updated", newState.Name.ValueString())
	assert.False(t, newState.Data.IsNull())
}

func TestIoaDeliveryStrategyResourceUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)
	patches.ApplyMethodReturn(c.client, "ModifyDeliveryStrategy", (*ioav20220601.ModifyDeliveryStrategyResponse)(nil), context.DeadlineExceeded)

	plan := ioaTestSamplePlanModel(t)
	plan.Id = types.StringValue(ioaTestStrategyId)
	plan.Data = ioaTestNullDataList(t)

	var updateResp resource.UpdateResponse
	updateResp.State = tfsdk.State{Schema: c.schema}
	c.resource.Update(context.Background(), resource.UpdateRequest{Plan: c.plan(t, plan), State: c.state(t, plan)}, &updateResp)

	assert.True(t, updateResp.Diagnostics.HasError())
}

func TestIoaDeliveryStrategyResourceDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)

	var capturedRequest *ioav20220601.DeleteDeliveryStrategyRequest
	patches.ApplyMethod(c.client, "DeleteDeliveryStrategy", func(request *ioav20220601.DeleteDeliveryStrategyRequest) (*ioav20220601.DeleteDeliveryStrategyResponse, error) {
		capturedRequest = request
		return &ioav20220601.DeleteDeliveryStrategyResponse{
			Response: &ioav20220601.DeleteDeliveryStrategyResponseParams{},
		}, nil
	})

	model := ioaTestSamplePlanModel(t)
	model.Id = types.StringValue(ioaTestStrategyId)

	var deleteResp resource.DeleteResponse
	deleteResp.State = c.state(t, model)
	c.resource.Delete(context.Background(), resource.DeleteRequest{State: deleteResp.State}, &deleteResp)

	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("unexpected delete diagnostics: %v", deleteResp.Diagnostics.Errors())
	}
	assert.NotNil(t, capturedRequest.ID)
	assert.Equal(t, uint64(42), *capturedRequest.ID)
}

func TestIoaDeliveryStrategyResourceDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newIoaDeliveryStrategyTestContext(t, patches)
	patches.ApplyMethodReturn(c.client, "DeleteDeliveryStrategy", (*ioav20220601.DeleteDeliveryStrategyResponse)(nil), context.DeadlineExceeded)

	model := ioaTestSamplePlanModel(t)
	model.Id = types.StringValue(ioaTestStrategyId)

	var deleteResp resource.DeleteResponse
	deleteResp.State = c.state(t, model)
	c.resource.Delete(context.Background(), resource.DeleteRequest{State: deleteResp.State}, &deleteResp)

	assert.True(t, deleteResp.Diagnostics.HasError())
}

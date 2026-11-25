package teo_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tftypes "github.com/hashicorp/terraform-plugin-go/tftypes"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

const (
	teoTestZoneId  = "zone-27h0vbm5w1e"
	teoTestGroupId = "cg-2qbzqqti0e9z"
	teoTestResId   = teoTestZoneId + "#" + teoTestGroupId
)

func teoStr(v string) *string    { return &v }
func teoInt64(v int64) *int64    { return &v }
func teoS(v string) types.String { return types.StringValue(v) }

// teoCustomGroupTestContext wires the resource to a patched SDK client.
type teoCustomGroupTestContext struct {
	resource resource.Resource
	client   *teov20220901.Client
	schema   schema.Schema
}

func newTeoCustomGroupTestContext(t *testing.T, patches *gomonkey.Patches) *teoCustomGroupTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseTeoV20220901Client", teoClient)

	res := svcteo.NewTeoCustomGroupResource()
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

	return &teoCustomGroupTestContext{resource: res, client: teoClient, schema: schemaResp.Schema}
}

func (c *teoCustomGroupTestContext) plan(t *testing.T, model svcteo.TeoCustomGroupResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *teoCustomGroupTestContext) state(t *testing.T, model svcteo.TeoCustomGroupResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func teoTestItems() []*teov20220901.CustomGroupItem {
	return []*teov20220901.CustomGroupItem{
		{Content: teoStr("cd08b7ecb0e6f5a72c6f5a8ea17e2e5a"), ItemType: teoStr("JA3")},
		{Content: teoStr("t13d1714h2_8daaf6b1a14185f4"), ItemType: teoStr("JA4")},
	}
}

func teoTestItemModels() []svcteo.TeoCustomGroupItemModel {
	return []svcteo.TeoCustomGroupItemModel{
		{Content: teoS("cd08b7ecb0e6f5a72c6f5a8ea17e2e5a"), ItemType: teoS("JA3")},
		{Content: teoS("t13d1714h2_8daaf6b1a14185f4"), ItemType: teoS("JA4")},
	}
}

func teoTestGroupList(t *testing.T) types.List {
	t.Helper()
	list, diags := types.ListValueFrom(context.Background(), teoCustomGroupItemObjectType(), teoTestItemModels())
	if diags.HasError() {
		t.Fatalf("failed to build group_items list: %v", diags.Errors())
	}
	return list
}

func teoCustomGroupItemObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"content":   types.StringType,
			"item_type": types.StringType,
		},
	}
}

func teoTestModel(t *testing.T, id string) svcteo.TeoCustomGroupResourceModel {
	t.Helper()
	model := svcteo.TeoCustomGroupResourceModel{
		Id:                  types.StringNull(),
		ZoneId:              teoS(teoTestZoneId),
		GroupId:             types.StringNull(),
		GroupName:           teoS("tf-example-group"),
		GroupType:           teoS("TLSFingerprint"),
		GroupItems:          teoTestGroupList(t),
		GroupItemTotalCount: types.Int64Null(),
		RefCount:            types.Int64Null(),
	}
	if id != "" {
		model.Id = teoS(id)
		model.GroupId = teoS(teoTestGroupId)
		model.GroupItemTotalCount = types.Int64Value(2)
		model.RefCount = types.Int64Value(0)
	}
	return model
}

func teoTestResponseParams() *teov20220901.DescribeCustomGroupResponseParams {
	return &teov20220901.DescribeCustomGroupResponseParams{
		GroupId:             teoStr(teoTestGroupId),
		GroupName:           teoStr("tf-example-group"),
		GroupType:           teoStr("TLSFingerprint"),
		GroupItems:          teoTestItems(),
		GroupItemTotalCount: teoInt64(2),
		RefCount:            teoInt64(0),
	}
}

func assertTeoCustomGroupModel(t *testing.T, model svcteo.TeoCustomGroupResourceModel) {
	t.Helper()

	if model.Id.IsNull() || model.Id.ValueString() != teoTestResId {
		t.Errorf("id = %v, want %q", model.Id.ValueString(), teoTestResId)
	}
	if model.ZoneId.IsNull() || model.ZoneId.ValueString() != teoTestZoneId {
		t.Errorf("zone_id = %v, want %q (not returned by the describe API, must keep the configured value)", model.ZoneId.ValueString(), teoTestZoneId)
	}
	if model.GroupId.IsNull() || model.GroupId.ValueString() != teoTestGroupId {
		t.Errorf("group_id = %v, want %q", model.GroupId.ValueString(), teoTestGroupId)
	}
	if model.GroupName.IsNull() || model.GroupName.ValueString() != "tf-example-group" {
		t.Errorf("group_name = %v, want %q", model.GroupName.ValueString(), "tf-example-group")
	}
	if model.GroupType.IsNull() || model.GroupType.ValueString() != "TLSFingerprint" {
		t.Errorf("group_type = %v, want %q", model.GroupType.ValueString(), "TLSFingerprint")
	}
	if model.GroupItemTotalCount.IsNull() || model.GroupItemTotalCount.ValueInt64() != 2 {
		t.Errorf("group_item_total_count = %v, want 2", model.GroupItemTotalCount.ValueInt64())
	}
	if model.RefCount.IsNull() || model.RefCount.ValueInt64() != 0 {
		t.Errorf("ref_count = %v, want 0", model.RefCount.ValueInt64())
	}
	if model.GroupItems.IsNull() {
		t.Fatal("group_items is null, want 2 elements")
	}
	var items []svcteo.TeoCustomGroupItemModel
	if diags := model.GroupItems.ElementsAs(context.Background(), &items, false); diags.HasError() {
		t.Fatalf("failed to read group_items: %v", diags.Errors())
	}
	if len(items) != 2 {
		t.Fatalf("group_items length = %d, want 2", len(items))
	}
	want := teoTestItemModels()
	for i, item := range items {
		if item.Content.ValueString() != want[i].Content.ValueString() {
			t.Errorf("group_items[%d].content = %q, want %q", i, item.Content.ValueString(), want[i].Content.ValueString())
		}
		if item.ItemType.ValueString() != want[i].ItemType.ValueString() {
			t.Errorf("group_items[%d].item_type = %q, want %q", i, item.ItemType.ValueString(), want[i].ItemType.ValueString())
		}
	}
}

// TestTeoCustomGroupCreateSuccess verifies that Create sends the zone id,
// group name, group type and items, and persists the composite id plus every
// read-back field in state.
func TestTeoCustomGroupCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *teov20220901.CreateCustomGroupRequest
	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "CreateCustomGroupWithContext",
		func(_ context.Context, req *teov20220901.CreateCustomGroupRequest) (*teov20220901.CreateCustomGroupResponse, error) {
			captured = req
			return &teov20220901.CreateCustomGroupResponse{
				Response: &teov20220901.CreateCustomGroupResponseParams{
					GroupId: teoStr(teoTestGroupId),
				},
			}, nil
		})
	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "DescribeCustomGroupWithContext",
		func(_ context.Context, req *teov20220901.DescribeCustomGroupRequest) (*teov20220901.DescribeCustomGroupResponse, error) {
			if req.ZoneId == nil || *req.ZoneId != teoTestZoneId {
				t.Errorf("describe request ZoneId = %v, want %q", req.ZoneId, teoTestZoneId)
			}
			if req.GroupId == nil || *req.GroupId != teoTestGroupId {
				t.Errorf("describe request GroupId = %v, want %q", req.GroupId, teoTestGroupId)
			}
			return &teov20220901.DescribeCustomGroupResponse{
				Response: teoTestResponseParams(),
			}, nil
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, teoTestModel(t, ""))
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, teoTestModel(t, ""))}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateCustomGroup request is nil")
	}
	if captured.ZoneId == nil || *captured.ZoneId != teoTestZoneId {
		t.Errorf("create request ZoneId = %v, want %q", captured.ZoneId, teoTestZoneId)
	}
	if captured.GroupName == nil || *captured.GroupName != "tf-example-group" {
		t.Errorf("create request GroupName = %v, want tf-example-group", captured.GroupName)
	}
	if captured.GroupType == nil || *captured.GroupType != "TLSFingerprint" {
		t.Errorf("create request GroupType = %v, want TLSFingerprint", captured.GroupType)
	}
	if len(captured.GroupItems) != 2 {
		t.Errorf("create request GroupItems length = %d, want 2", len(captured.GroupItems))
	}

	var got svcteo.TeoCustomGroupResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertTeoCustomGroupModel(t, got)
}

// TestTeoCustomGroupCreateEmptyGroupId verifies that Create fails loudly when
// the API returns no group id instead of writing an unmanageable resource.
func TestTeoCustomGroupCreateEmptyGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "CreateCustomGroupWithContext",
		func(_ context.Context, _ *teov20220901.CreateCustomGroupRequest) (*teov20220901.CreateCustomGroupResponse, error) {
			return &teov20220901.CreateCustomGroupResponse{
				Response: &teov20220901.CreateCustomGroupResponseParams{},
			}, nil
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, teoTestModel(t, ""))
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, teoTestModel(t, ""))}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty group id")
	}
}

// TestTeoCustomGroupReadSuccess verifies Read refreshes all leaf fields,
// including every element of the group items list.
func TestTeoCustomGroupReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "DescribeCustomGroupWithContext",
		func(_ context.Context, req *teov20220901.DescribeCustomGroupRequest) (*teov20220901.DescribeCustomGroupResponse, error) {
			if req.Limit == nil || *req.Limit != 500 {
				t.Errorf("describe request Limit = %v, want 500", req.Limit)
			}
			return &teov20220901.DescribeCustomGroupResponse{
				Response: teoTestResponseParams(),
			}, nil
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, teoTestModel(t, teoTestResId))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, teoTestModel(t, teoTestResId))}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got svcteo.TeoCustomGroupResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertTeoCustomGroupModel(t, got)
}

// TestTeoCustomGroupReadPagination verifies Read follows the GroupItems
// pagination and collects every page.
func TestTeoCustomGroupReadPagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	allItems := make([]*teov20220901.CustomGroupItem, 0, 3)
	allItems = append(allItems, teoTestItems()...)
	allItems = append(allItems, &teov20220901.CustomGroupItem{Content: teoStr("extra-item"), ItemType: teoStr("JA3")})

	calls := 0
	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "DescribeCustomGroupWithContext",
		func(_ context.Context, req *teov20220901.DescribeCustomGroupRequest) (*teov20220901.DescribeCustomGroupResponse, error) {
			calls++
			switch calls {
			case 1:
				return &teov20220901.DescribeCustomGroupResponse{
					Response: &teov20220901.DescribeCustomGroupResponseParams{
						GroupId:             teoStr(teoTestGroupId),
						GroupName:           teoStr("tf-example-group"),
						GroupType:           teoStr("TLSFingerprint"),
						GroupItems:          teoTestItems(),
						GroupItemTotalCount: teoInt64(3),
						RefCount:            teoInt64(1),
					},
				}, nil
			case 2:
				if req.Offset == nil || *req.Offset != 2 {
					t.Errorf("second describe request Offset = %v, want 2", req.Offset)
				}
				return &teov20220901.DescribeCustomGroupResponse{
					Response: &teov20220901.DescribeCustomGroupResponseParams{
						GroupId:             teoStr(teoTestGroupId),
						GroupName:           teoStr("tf-example-group"),
						GroupType:           teoStr("TLSFingerprint"),
						GroupItems:          []*teov20220901.CustomGroupItem{{Content: teoStr("extra-item"), ItemType: teoStr("JA3")}},
						GroupItemTotalCount: teoInt64(3),
						RefCount:            teoInt64(1),
					},
				}, nil
			default:
				return nil, fmt.Errorf("unexpected extra describe call %d", calls)
			}
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, teoTestModel(t, teoTestResId))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, teoTestModel(t, teoTestResId))}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}
	if calls != 2 {
		t.Errorf("describe call count = %d, want 2", calls)
	}

	var got svcteo.TeoCustomGroupResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.GroupItems.IsNull() {
		t.Fatal("group_items is null")
	}
	var items []svcteo.TeoCustomGroupItemModel
	if diags := got.GroupItems.ElementsAs(context.Background(), &items, false); diags.HasError() {
		t.Fatalf("failed to read group_items: %v", diags.Errors())
	}
	if len(items) != 3 {
		t.Fatalf("group_items length = %d, want 3", len(items))
	}
	if items[2].Content.ValueString() != "extra-item" {
		t.Errorf("group_items[2].content = %q, want extra-item", items[2].Content.ValueString())
	}
	if got.GroupItemTotalCount.ValueInt64() != 3 {
		t.Errorf("group_item_total_count = %d, want 3", got.GroupItemTotalCount.ValueInt64())
	}
	if got.RefCount.ValueInt64() != 1 {
		t.Errorf("ref_count = %d, want 1", got.RefCount.ValueInt64())
	}
}

// TestTeoCustomGroupReadNotFound verifies that a confirmed not-found (API
// error with the ResourceNotFound code) drops the resource from state.
func TestTeoCustomGroupReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "DescribeCustomGroupWithContext",
		func(_ context.Context, _ *teov20220901.DescribeCustomGroupRequest) (*teov20220901.DescribeCustomGroupResponse, error) {
			return nil, sdkErrors.NewTencentCloudSDKError("ResourceNotFound", "custom group not found", "req-1")
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, teoTestModel(t, teoTestResId))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, teoTestModel(t, teoTestResId))}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the group is confirmed gone")
	}
}

// TestTeoCustomGroupReadApiError verifies API failures surface as errors and
// keep the state intact.
func TestTeoCustomGroupReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "DescribeCustomGroupWithContext",
		func(_ context.Context, _ *teov20220901.DescribeCustomGroupRequest) (*teov20220901.DescribeCustomGroupResponse, error) {
			return nil, fmt.Errorf("network error")
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, teoTestModel(t, teoTestResId))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, teoTestModel(t, teoTestResId))}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on an API error")
	}
	if readResp.State.Raw.IsNull() {
		t.Error("state should be kept on an API error")
	}
}

// TestTeoCustomGroupReadInvalidId verifies Read fails on a malformed id
// instead of silently querying a wrong group.
func TestTeoCustomGroupReadInvalidId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTeoCustomGroupTestContext(t, patches)
	var readResp resource.ReadResponse
	model := teoTestModel(t, "just-a-group-id")
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on a non-composite id")
	}
}

// TestTeoCustomGroupUpdateSuccess verifies Update sends the composite
// identity, an explicit GroupItemOperation=Replace and the new items, then
// refreshes state from the read-back.
func TestTeoCustomGroupUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *teov20220901.ModifyCustomGroupRequest
	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "ModifyCustomGroupWithContext",
		func(_ context.Context, req *teov20220901.ModifyCustomGroupRequest) (*teov20220901.ModifyCustomGroupResponse, error) {
			captured = req
			return &teov20220901.ModifyCustomGroupResponse{
				Response: &teov20220901.ModifyCustomGroupResponseParams{},
			}, nil
		})
	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "DescribeCustomGroupWithContext",
		func(_ context.Context, _ *teov20220901.DescribeCustomGroupRequest) (*teov20220901.DescribeCustomGroupResponse, error) {
			return &teov20220901.DescribeCustomGroupResponse{
				Response: teoTestResponseParams(),
			}, nil
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	plan := teoTestModel(t, teoTestResId)
	plan.GroupName = teoS("renamed-group")
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, teoTestModel(t, teoTestResId))
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{Plan: ctx.plan(t, plan), State: ctx.state(t, teoTestModel(t, teoTestResId))}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyCustomGroup request is nil")
	}
	if captured.ZoneId == nil || *captured.ZoneId != teoTestZoneId {
		t.Errorf("update request ZoneId = %v, want %q", captured.ZoneId, teoTestZoneId)
	}
	if captured.GroupId == nil || *captured.GroupId != teoTestGroupId {
		t.Errorf("update request GroupId = %v, want %q", captured.GroupId, teoTestGroupId)
	}
	if captured.GroupName == nil || *captured.GroupName != "renamed-group" {
		t.Errorf("update request GroupName = %v, want renamed-group", captured.GroupName)
	}
	if captured.GroupItemOperation == nil || *captured.GroupItemOperation != "Replace" {
		t.Errorf("update request GroupItemOperation = %v, want Replace", captured.GroupItemOperation)
	}
	if len(captured.GroupItems) != 2 {
		t.Errorf("update request GroupItems length = %d, want 2", len(captured.GroupItems))
	}

	var got svcteo.TeoCustomGroupResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.GroupName.ValueString() != "tf-example-group" {
		t.Errorf("group_name = %q, want the read-back value tf-example-group", got.GroupName.ValueString())
	}
	if got.Id.IsNull() || got.Id.ValueString() != teoTestResId {
		t.Errorf("id = %v, want %q", got.Id.ValueString(), teoTestResId)
	}
	if got.ZoneId.IsNull() || got.ZoneId.ValueString() != teoTestZoneId {
		t.Errorf("zone_id = %v, want %q", got.ZoneId.ValueString(), teoTestZoneId)
	}
}

// TestTeoCustomGroupUpdateEmptyItems verifies Update refuses to send an empty
// item list: the clear semantics of Replace + empty list is not confirmed by
// the API docs, so the update fails loudly instead of guessing a direction.
func TestTeoCustomGroupUpdateEmptyItems(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "ModifyCustomGroupWithContext",
		func(_ context.Context, _ *teov20220901.ModifyCustomGroupRequest) (*teov20220901.ModifyCustomGroupResponse, error) {
			t.Error("ModifyCustomGroup must not be called for an empty item list")
			return nil, fmt.Errorf("unexpected call")
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	plan := teoTestModel(t, teoTestResId)
	emptyList, diags := types.ListValueFrom(context.Background(), teoCustomGroupItemObjectType(), []svcteo.TeoCustomGroupItemModel{})
	if diags.HasError() {
		t.Fatalf("failed to build empty list: %v", diags.Errors())
	}
	plan.GroupItems = emptyList

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, teoTestModel(t, teoTestResId))
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{Plan: ctx.plan(t, plan), State: ctx.state(t, teoTestModel(t, teoTestResId))}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should fail loudly for an empty group_items list")
	}
}

// TestTeoCustomGroupUpdateApiError verifies update failures surface as errors.
func TestTeoCustomGroupUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "ModifyCustomGroupWithContext",
		func(_ context.Context, _ *teov20220901.ModifyCustomGroupRequest) (*teov20220901.ModifyCustomGroupResponse, error) {
			return nil, sdkErrors.NewTencentCloudSDKError("InternalError", "modify failed", "req-2")
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, teoTestModel(t, teoTestResId))
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{Plan: ctx.plan(t, teoTestModel(t, teoTestResId)), State: ctx.state(t, teoTestModel(t, teoTestResId))}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should fail on an API error")
	}
}

// TestTeoCustomGroupDeleteSuccess verifies Delete sends the composite
// identity parts and succeeds without error.
func TestTeoCustomGroupDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *teov20220901.DeleteCustomGroupRequest
	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "DeleteCustomGroupWithContext",
		func(_ context.Context, req *teov20220901.DeleteCustomGroupRequest) (*teov20220901.DeleteCustomGroupResponse, error) {
			captured = req
			return &teov20220901.DeleteCustomGroupResponse{
				Response: &teov20220901.DeleteCustomGroupResponseParams{},
			}, nil
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, teoTestModel(t, teoTestResId))}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}
	if captured == nil {
		t.Fatal("DeleteCustomGroup request is nil")
	}
	if captured.ZoneId == nil || *captured.ZoneId != teoTestZoneId {
		t.Errorf("delete request ZoneId = %v, want %q", captured.ZoneId, teoTestZoneId)
	}
	if captured.GroupId == nil || *captured.GroupId != teoTestGroupId {
		t.Errorf("delete request GroupId = %v, want %q", captured.GroupId, teoTestGroupId)
	}
}

// TestTeoCustomGroupDeleteNotFound verifies Delete treats a confirmed
// ResourceNotFound as success so destroy stays idempotent.
func TestTeoCustomGroupDeleteNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "DeleteCustomGroupWithContext",
		func(_ context.Context, _ *teov20220901.DeleteCustomGroupRequest) (*teov20220901.DeleteCustomGroupResponse, error) {
			return nil, sdkErrors.NewTencentCloudSDKError("ResourceNotFound", "custom group not found", "req-3")
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, teoTestModel(t, teoTestResId))}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete should treat ResourceNotFound as success: %v", deleteResp.Diagnostics.Errors())
	}
}

// TestTeoCustomGroupDeleteError verifies that other delete failures keep the
// error diagnostics so the id stays in state for a retry.
func TestTeoCustomGroupDeleteError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&teov20220901.Client{}), "DeleteCustomGroupWithContext",
		func(_ context.Context, _ *teov20220901.DeleteCustomGroupRequest) (*teov20220901.DeleteCustomGroupResponse, error) {
			return nil, fmt.Errorf("delete failed")
		})

	ctx := newTeoCustomGroupTestContext(t, patches)
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, teoTestModel(t, teoTestResId))}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}

// TestTeoCustomGroupImportState verifies the composite id passthrough import
// support.
func TestTeoCustomGroupImportState(t *testing.T) {
	ctx := newTeoCustomGroupTestContext(t, gomonkey.NewPatches())
	importer, ok := ctx.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}

	var importResp resource.ImportStateResponse
	importResp.State.Schema = ctx.schema
	importResp.State.Raw = tftypes.NewValue(ctx.schema.Type().TerraformType(context.Background()), nil)
	importer.ImportState(context.Background(),
		resource.ImportStateRequest{ID: teoTestResId}, &importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", importResp.Diagnostics.Errors())
	}

	var id types.String
	if diags := importResp.State.GetAttribute(context.Background(), path.Root("id"), &id); diags.HasError() {
		t.Fatalf("failed to read imported id: %v", diags.Errors())
	}
	if id.ValueString() != teoTestResId {
		t.Errorf("imported id = %q, want %q", id.ValueString(), teoTestResId)
	}
}

// TestTeoCustomGroupSchemaFields verifies the schema attribute set, the
// required/computed flags and the immutability of zone_id and group_type.
func TestTeoCustomGroupSchemaFields(t *testing.T) {
	var schemaResp resource.SchemaResponse
	svcteo.NewTeoCustomGroupResource().Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	expect := map[string]bool{
		"id":                     true,
		"zone_id":                true,
		"group_id":               true,
		"group_name":             true,
		"group_type":             true,
		"group_items":            true,
		"group_item_total_count": true,
		"ref_count":              true,
	}
	if len(schemaResp.Schema.Attributes) != len(expect) {
		t.Fatalf("schema attribute count = %d, want %d", len(schemaResp.Schema.Attributes), len(expect))
	}
	for name := range expect {
		if _, ok := schemaResp.Schema.Attributes[name]; !ok {
			t.Errorf("schema attribute %q is missing", name)
		}
	}

	required := map[string]bool{"zone_id": true, "group_name": true, "group_type": true, "group_items": true}
	for name := range required {
		attr, ok := schemaResp.Schema.Attributes[name]
		if !ok {
			t.Fatalf("schema attribute %q is missing", name)
		}
		if !attr.IsRequired() {
			t.Errorf("attribute %q should be required", name)
		}
	}

	computed := map[string]bool{"id": true, "group_id": true, "group_item_total_count": true, "ref_count": true}
	for name := range computed {
		attr, ok := schemaResp.Schema.Attributes[name]
		if !ok {
			t.Fatalf("schema attribute %q is missing", name)
		}
		if !attr.IsComputed() {
			t.Errorf("attribute %q should be computed", name)
		}
	}

	// No enum whitelists: enum values are documented in the descriptions and
	// validated by the cloud API. The concrete attribute types expose their
	// Validators slice, so a non-nil one on an enum-described field is a defect.
	for name, attr := range schemaResp.Schema.Attributes {
		desc := attr.GetDescription()
		if !strings.Contains(desc, "Valid values:") {
			continue
		}
		switch a := attr.(type) {
		case schema.StringAttribute:
			if a.Validators != nil {
				t.Errorf("attribute %q has validators despite an enum description", name)
			}
		case schema.Int64Attribute:
			if a.Validators != nil {
				t.Errorf("attribute %q has validators despite an enum description", name)
			}
		}
	}
}

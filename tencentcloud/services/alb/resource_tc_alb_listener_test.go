package alb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svAlb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/alb"
)

const (
	albListenerTestId          = "lst-12345678"
	albListenerTestLbId        = "alb-12345678"
	albListenerTestTargetGroup = "lbtg-12345678"
	albListenerTestPort        = int64(443)
	albListenerTestProtocol    = "HTTPS"
	albListenerTestActionType  = "TargetGroup"
	albListenerTestName        = "tf-example-listener"
	albListenerTestRenamedName = "tf-example-listener-renamed"
	albListenerTestCloudName   = "tf-example-listener-normalized"
	albListenerTestTagKey      = "env"
	albListenerTestTagValue    = "tf-example"
	albListenerTestWeight      = int64(10)
)

func albListenerStr(v string) *string    { return &v }
func albListenerBool(v bool) *bool       { return &v }
func albListenerUint64(v uint64) *uint64 { return &v }

// albListenerTestItem builds the DescribeListeners response item.
func albListenerTestItem() *albv20251030.ListenerOutput {
	return &albv20251030.ListenerOutput{
		ListenerId:       albListenerStr(albListenerTestId),
		ListenerName:     albListenerStr(albListenerTestName),
		ListenerPort:     albListenerUint64(uint64(albListenerTestPort)),
		ListenerProtocol: albListenerStr(albListenerTestProtocol),
		GzipEnabled:      albListenerBool(false),
		Http2Enable:      albListenerBool(false),
		IdleTimeout:      albListenerUint64(15),
		RequestTimeout:   albListenerUint64(60),
		CaEnable:         albListenerBool(false),
		Tags: []*albv20251030.TagInfo{
			{TagKey: albListenerStr(albListenerTestTagKey), TagValue: albListenerStr(albListenerTestTagValue)},
		},
	}
}

// albListenerTestModel builds the terraform model used for plan/state.
// Every unset attribute carries an EXPLICIT TYPED null value: the raw zero
// value (types.List{} etc.) breaks tfsdk Plan/State reflection.
func albListenerTestModel() svAlb.AlbListenerResourceModel {
	return svAlb.AlbListenerResourceModel{
		Id:               types.StringValue(albListenerTestId),
		LoadBalancerId:   types.StringValue(albListenerTestLbId),
		ListenerPort:     types.Int64Value(albListenerTestPort),
		ListenerProtocol: types.StringValue(albListenerTestProtocol),
		DefaultActions: []svAlb.AlbListenerDefaultActionModel{
			{
				Type: types.StringValue(albListenerTestActionType),
				TargetGroupConfig: &svAlb.AlbListenerTargetGroupConfigModel{
					TargetGroups: []svAlb.AlbListenerTargetGroupModel{
						{
							TargetGroupId: types.StringValue(albListenerTestTargetGroup),
							Weight:        types.Int64Value(albListenerTestWeight),
						},
					},
				},
			},
		},
		CaCertificateIds: types.ListNull(types.StringType),
		CaEnabled:        types.BoolNull(),
		CertificateIds:   types.ListNull(types.StringType),
		ClientToken:      types.StringNull(),
		GzipEnabled:      types.BoolNull(),
		Http2Enabled:     types.BoolNull(),
		IdleTimeout:      types.Int64Null(),
		ListenerName:     types.StringValue(albListenerTestName),
		RequestTimeout:   types.Int64Null(),
		SecurityPolicyId: types.StringNull(),
		Tags: []svAlb.AlbListenerTagModel{
			{
				TagKey:   types.StringValue(albListenerTestTagKey),
				TagValue: types.StringValue(albListenerTestTagValue),
			},
		},
		XForwardedForConfig: nil,
	}
}

// assertAlbListenerModel asserts the user-facing model fields refreshed from
// the describe response.
func assertAlbListenerModel(t *testing.T, m svAlb.AlbListenerResourceModel) {
	t.Helper()

	if m.Id.ValueString() != albListenerTestId {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), albListenerTestId)
	}
	if m.LoadBalancerId.ValueString() != albListenerTestLbId {
		t.Errorf("load_balancer_id = %q, want %q", m.LoadBalancerId.ValueString(), albListenerTestLbId)
	}
	if m.ListenerPort.ValueInt64() != albListenerTestPort {
		t.Errorf("listener_port = %d, want %d", m.ListenerPort.ValueInt64(), albListenerTestPort)
	}
	if m.ListenerProtocol.ValueString() != albListenerTestProtocol {
		t.Errorf("listener_protocol = %q, want %q", m.ListenerProtocol.ValueString(), albListenerTestProtocol)
	}
	if m.ListenerName.ValueString() != albListenerTestName {
		t.Errorf("listener_name = %q, want %q", m.ListenerName.ValueString(), albListenerTestName)
	}
	if m.IdleTimeout.ValueInt64() != 15 {
		t.Errorf("idle_timeout = %d, want 15", m.IdleTimeout.ValueInt64())
	}
	if m.RequestTimeout.ValueInt64() != 60 {
		t.Errorf("request_timeout = %d, want 60", m.RequestTimeout.ValueInt64())
	}
	if len(m.DefaultActions) != 1 {
		t.Fatalf("default_actions has %d item(s), want 1", len(m.DefaultActions))
	}
	if m.DefaultActions[0].Type.ValueString() != albListenerTestActionType {
		t.Errorf("default_actions[0].type = %q, want %q", m.DefaultActions[0].Type.ValueString(), albListenerTestActionType)
	}
	if len(m.Tags) != 1 || m.Tags[0].TagKey.ValueString() != albListenerTestTagKey || m.Tags[0].TagValue.ValueString() != albListenerTestTagValue {
		t.Errorf("tags = %v, want key %q value %q", m.Tags, albListenerTestTagKey, albListenerTestTagValue)
	}
}

// albListenerTestContext wires the resource to a configured client; the SDK
// client methods are stubbed per test case.
type albListenerTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newAlbListenerTestContext(t *testing.T) *albListenerTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svAlb.NewAlbListenerResource()
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

	return &albListenerTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *albListenerTestContext) plan(t *testing.T, m svAlb.AlbListenerResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *albListenerTestContext) state(t *testing.T, m svAlb.AlbListenerResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// patchAlbListenerSdkClient stubs the versioned client accessor and returns
// the mock SDK client whose methods the test patches next.
func patchAlbListenerSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *albv20251030.Client {
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)

	return albClient
}

// stubAlbListenerDescribe stubs DescribeListeners with the given listeners.
func stubAlbListenerDescribe(patches *gomonkey.Patches, albClient *albv20251030.Client, listeners []*albv20251030.ListenerOutput) {
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(req *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
			Listeners: listeners,
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})
}

// TestAlbListenerCreateSuccess verifies that Create sends all request fields,
// tracks the returned listener id and refreshes the state.
func TestAlbListenerCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.CreateListenerRequest
	patches.ApplyMethodFunc(albClient, "CreateListener", func(req *albv20251030.CreateListenerRequest) (*albv20251030.CreateListenerResponse, error) {
		captured = req

		return &albv20251030.CreateListenerResponse{Response: &albv20251030.CreateListenerResponseParams{
			ListenerId: albListenerStr(albListenerTestId),
			RequestId:  albListenerStr("fake-request-id"),
		}}, nil
	})
	stubAlbListenerDescribe(patches, albClient, []*albv20251030.ListenerOutput{albListenerTestItem()})

	model := albListenerTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateListener request is nil")
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albListenerTestLbId {
		t.Errorf("create request load_balancer_id = %v, want %q", captured.LoadBalancerId, albListenerTestLbId)
	}
	if captured.ListenerPort == nil || *captured.ListenerPort != uint64(albListenerTestPort) {
		t.Errorf("create request listener_port = %v, want %d", captured.ListenerPort, albListenerTestPort)
	}
	if captured.ListenerProtocol == nil || *captured.ListenerProtocol != albListenerTestProtocol {
		t.Errorf("create request listener_protocol = %v, want %q", captured.ListenerProtocol, albListenerTestProtocol)
	}
	if captured.ListenerName == nil || *captured.ListenerName != albListenerTestName {
		t.Errorf("create request listener_name = %v, want %q", captured.ListenerName, albListenerTestName)
	}
	if len(captured.DefaultActions) != 1 {
		t.Fatalf("create request default_actions has %d item(s), want 1", len(captured.DefaultActions))
	}
	action := captured.DefaultActions[0]
	if action.Type == nil || *action.Type != albListenerTestActionType {
		t.Errorf("create request default_actions[0].type = %v, want %q", action.Type, albListenerTestActionType)
	}
	if action.TargetGroupConfig == nil || len(action.TargetGroupConfig.TargetGroups) != 1 {
		t.Fatalf("create request default_actions[0].target_group_config.target_groups is missing")
	}
	tuple := action.TargetGroupConfig.TargetGroups[0]
	if tuple.TargetGroupId == nil || *tuple.TargetGroupId != albListenerTestTargetGroup {
		t.Errorf("create request target_group_id = %v, want %q", tuple.TargetGroupId, albListenerTestTargetGroup)
	}
	if tuple.Weight == nil || *tuple.Weight != albListenerTestWeight {
		t.Errorf("create request weight = %v, want %d", tuple.Weight, albListenerTestWeight)
	}
	if len(captured.Tags) != 1 || captured.Tags[0].TagKey == nil || *captured.Tags[0].TagKey != albListenerTestTagKey {
		t.Errorf("create request tags = %v, want key %q", captured.Tags, albListenerTestTagKey)
	}

	var got svAlb.AlbListenerResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertAlbListenerModel(t, got)
}

// TestAlbListenerCreateEmptyId verifies that Create fails loudly when the API
// returns no listener id.
func TestAlbListenerCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "CreateListener", func(req *albv20251030.CreateListenerRequest) (*albv20251030.CreateListenerResponse, error) {
		return &albv20251030.CreateListenerResponse{Response: &albv20251030.CreateListenerResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, albListenerTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty listener id")
	}
}

// TestAlbListenerCreateApiError verifies API failures surface as errors.
func TestAlbListenerCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "CreateListener", func(req *albv20251030.CreateListenerRequest) (*albv20251030.CreateListenerResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, albListenerTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestAlbListenerReadSuccess verifies Read refreshes all leaf fields and
// addresses the listener by its id.
func TestAlbListenerReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var gotListenerId, gotLbId string
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(req *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		if len(req.ListenerIds) != 1 || req.ListenerIds[0] == nil {
			t.Errorf("expected the listener id filter to be passed")
		} else {
			gotListenerId = *req.ListenerIds[0]
		}
		if req.LoadBalancerId != nil {
			gotLbId = *req.LoadBalancerId
		}

		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
			Listeners: []*albv20251030.ListenerOutput{albListenerTestItem()},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, albListenerTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotListenerId != albListenerTestId {
		t.Errorf("DescribeListeners called with listener id %q, want %q", gotListenerId, albListenerTestId)
	}
	if gotLbId != albListenerTestLbId {
		t.Errorf("DescribeListeners called with load balancer id %q, want %q", gotLbId, albListenerTestLbId)
	}

	var got svAlb.AlbListenerResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertAlbListenerModel(t, got)
}

// TestAlbListenerReadApiError verifies describe failures surface as errors.
func TestAlbListenerReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(req *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		return nil, fmt.Errorf("describe failed")
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, albListenerTestModel())}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on an API error")
	}
}

// TestAlbListenerReadNotFound verifies that a describe list without the
// tracked listener drops the resource from state (confirmed not-found).
func TestAlbListenerReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	stubAlbListenerDescribe(patches, albClient, nil)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, albListenerTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the listener is confirmed gone")
	}
}

// TestAlbListenerUpdateSuccess verifies Update sends the PLANNED value to the
// modify request and records the POST-update cloud state from the read-back —
// not the plan.
func TestAlbListenerUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.ModifyListenerAttributesRequest
	patches.ApplyMethodFunc(albClient, "ModifyListenerAttributes", func(req *albv20251030.ModifyListenerAttributesRequest) (*albv20251030.ModifyListenerAttributesResponse, error) {
		captured = req

		return &albv20251030.ModifyListenerAttributesResponse{Response: &albv20251030.ModifyListenerAttributesResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(req *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		normalized := albListenerTestItem()
		normalized.ListenerName = albListenerStr(albListenerTestCloudName)

		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
			Listeners: []*albv20251030.ListenerOutput{normalized},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	plannedModel := albListenerTestModel()
	plannedModel.ListenerName = types.StringValue(albListenerTestRenamedName)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, albListenerTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyListenerAttributes request is nil")
	}
	if captured.ListenerId == nil || *captured.ListenerId != albListenerTestId {
		t.Errorf("update request listener_id = %v, want %q", captured.ListenerId, albListenerTestId)
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albListenerTestLbId {
		t.Errorf("update request load_balancer_id = %v, want %q", captured.LoadBalancerId, albListenerTestLbId)
	}
	if captured.ListenerName == nil || *captured.ListenerName != albListenerTestRenamedName {
		t.Errorf("update request listener_name = %v, want the planned %q", captured.ListenerName, albListenerTestRenamedName)
	}

	var got svAlb.AlbListenerResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != albListenerTestId {
		t.Errorf("id = %q, want the unchanged %q", got.Id.ValueString(), albListenerTestId)
	}
	if got.ListenerName.ValueString() != albListenerTestCloudName {
		t.Errorf("listener_name = %q, want the read-back cloud value %q", got.ListenerName.ValueString(), albListenerTestCloudName)
	}
}

// TestAlbListenerUpdateNoChange verifies that a plan identical to the state
// calls NO write API: no attribute ModifyListenerAttributes covers changed,
// so the call would be a pointless no-change request. The read-back still
// runs and the state is refreshed.
func TestAlbListenerUpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	modifyCalls := 0
	patches.ApplyMethodFunc(albClient, "ModifyListenerAttributes", func(req *albv20251030.ModifyListenerAttributesRequest) (*albv20251030.ModifyListenerAttributesResponse, error) {
		modifyCalls++

		return nil, fmt.Errorf("the modify API must not be called")
	})
	stubAlbListenerDescribe(patches, albClient, []*albv20251030.ListenerOutput{albListenerTestItem()})

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, albListenerTestModel()),
			State: ctx.state(t, albListenerTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if modifyCalls != 0 {
		t.Errorf("ModifyListenerAttributes was called %d time(s) for a no-change plan", modifyCalls)
	}

	var got svAlb.AlbListenerResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertAlbListenerModel(t, got)
}

// TestAlbListenerUpdateApiError verifies modify failures surface as errors.
func TestAlbListenerUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "ModifyListenerAttributes", func(req *albv20251030.ModifyListenerAttributesRequest) (*albv20251030.ModifyListenerAttributesResponse, error) {
		return nil, fmt.Errorf("modify failed")
	})

	plannedModel := albListenerTestModel()
	plannedModel.ListenerName = types.StringValue(albListenerTestRenamedName)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, albListenerTestModel()),
		}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should fail on an API error")
	}
}

// TestAlbListenerDeleteSuccess verifies Delete sends the delete request with
// the tracked id.
func TestAlbListenerDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.DeleteListenerRequest
	patches.ApplyMethodFunc(albClient, "DeleteListener", func(req *albv20251030.DeleteListenerRequest) (*albv20251030.DeleteListenerResponse, error) {
		captured = req

		return &albv20251030.DeleteListenerResponse{Response: &albv20251030.DeleteListenerResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, albListenerTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteListener request is nil")
	}
	if len(captured.ListenerIds) != 1 || captured.ListenerIds[0] == nil || *captured.ListenerIds[0] != albListenerTestId {
		t.Errorf("delete request listener_ids = %v, want [%q]", captured.ListenerIds, albListenerTestId)
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albListenerTestLbId {
		t.Errorf("delete request load_balancer_id = %v, want %q", captured.LoadBalancerId, albListenerTestLbId)
	}
}

// TestAlbListenerDescribeListenersPagination verifies the service method's
// built-in NextToken paging: two pages are fetched and combined.
func TestAlbListenerDescribeListenersPagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	first := albListenerTestItem()
	second := albListenerTestItem()
	second.ListenerId = albListenerStr("lst-87654321")

	callCount := 0
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(req *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		callCount++
		if callCount == 1 {
			return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
				Listeners: []*albv20251030.ListenerOutput{first},
				NextToken: albListenerStr("next-page-token"),
				RequestId: albListenerStr("fake-request-id"),
			}}, nil
		}

		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
			Listeners: []*albv20251030.ListenerOutput{second},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = albListenerStr(albListenerTestLbId)

	service := svAlb.NewAlbService(ctx.client)
	listeners, err := service.DescribeListeners(context.Background(), describeListenersRequest)
	if err != nil {
		t.Fatalf("DescribeListeners failed: %v", err)
	}

	if callCount != 2 {
		t.Errorf("DescribeListeners was called %d time(s), want 2 pages", callCount)
	}
	if len(listeners) != 2 {
		t.Fatalf("listeners has %d item(s), want 2", len(listeners))
	}
	if listeners[0].ListenerId == nil || *listeners[0].ListenerId != albListenerTestId {
		t.Errorf("listeners[0].ListenerId = %v, want %q", listeners[0].ListenerId, albListenerTestId)
	}
	if listeners[1].ListenerId == nil || *listeners[1].ListenerId != "lst-87654321" {
		t.Errorf("listeners[1].ListenerId = %v, want lst-87654321", listeners[1].ListenerId)
	}
}

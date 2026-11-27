package tdmysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svtdmysql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql"
)

const (
	tdmysqlUserTestId          = "tf-test-instance#tf-example-user#10.0.0.1"
	tdmysqlUserTestInstanceId  = "tf-test-instance"
	tdmysqlUserTestUserName    = "tf-example-user"
	tdmysqlUserTestHost        = "10.0.0.1"
	tdmysqlUserTestPassword    = "TF@example-12345"
	tdmysqlUserTestDescription = "created by terraform"
	tdmysqlUserTestFlowId      = int64(12345)
)

// tdmysqlUserStr is the string pointer helper for this resource's tests.
func tdmysqlUserStr(v string) *string { return &v }

func tdmysqlUserFlowId() *int64 { v := tdmysqlUserTestFlowId; return &v }

// tdmysqlUserTestItem builds the DescribeUsers response item.
func tdmysqlUserTestItem() *tdmysqlv20211122.UserInfo {
	return &tdmysqlv20211122.UserInfo{
		UserName:    tdmysqlUserStr(tdmysqlUserTestUserName),
		Host:        tdmysqlUserStr(tdmysqlUserTestHost),
		Description: tdmysqlUserStr(tdmysqlUserTestDescription),
	}
}

// tdmysqlUserTestModel builds the terraform model used for plan/state.
// Every unset attribute carries an EXPLICIT TYPED null value — the raw zero
// value (types.List{} etc.) breaks tfsdk Plan/State reflection. Cases
// exercising a wait path pass a SECOND-SCALE short user value (5s): the test
// verifies the code path is REACHED, never hard-waits a minute-scale default.
func tdmysqlUserTestModel() svtdmysql.TdmysqlUserResourceModel {
	return svtdmysql.TdmysqlUserResourceModel{
		Id:          types.StringValue(tdmysqlUserTestId),
		InstanceId:  types.StringValue(tdmysqlUserTestInstanceId),
		UserName:    types.StringValue(tdmysqlUserTestUserName),
		Host:        types.StringValue(tdmysqlUserTestHost),
		Password:    types.StringValue(tdmysqlUserTestPassword),
		Description: types.StringNull(),
		Timeouts: resourcetimeouts.Value{Object: types.ObjectValueMust(
			tdmysqlUserTimeoutsAttrTypes(),
			map[string]attr.Value{
				"create": types.StringValue("5s"),
				"update": types.StringValue("5s"),
				"delete": types.StringValue("5s"),
			},
		)},
	}
}

// tdmysqlUserTimeoutsAttrTypes mirrors the attribute types of the timeouts
// block declared by the resource schema.
func tdmysqlUserTimeoutsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"create": types.StringType,
		"update": types.StringType,
		"delete": types.StringType,
	}
}

// assertTdmysqlUserModel asserts the model fields that survive a refresh.
func assertTdmysqlUserModel(t *testing.T, m svtdmysql.TdmysqlUserResourceModel) {
	t.Helper()

	if m.Id.ValueString() != tdmysqlUserTestId {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), tdmysqlUserTestId)
	}
	if m.InstanceId.ValueString() != tdmysqlUserTestInstanceId {
		t.Errorf("instance_id = %q, want %q", m.InstanceId.ValueString(), tdmysqlUserTestInstanceId)
	}
	if m.UserName.ValueString() != tdmysqlUserTestUserName {
		t.Errorf("user_name = %q, want %q", m.UserName.ValueString(), tdmysqlUserTestUserName)
	}
	if m.Host.ValueString() != tdmysqlUserTestHost {
		t.Errorf("host = %q, want %q", m.Host.ValueString(), tdmysqlUserTestHost)
	}
	if m.Description.ValueString() != tdmysqlUserTestDescription {
		t.Errorf("description = %q, want %q", m.Description.ValueString(), tdmysqlUserTestDescription)
	}
}

// patchTdmysqlUserSdkClient stubs the versioned client accessor and returns
// the mock SDK client whose methods the test patches next (two-level mode).
func patchTdmysqlUserSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *tdmysqlv20211122.Client {
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)

	return tdmysqlClient
}

// stubTdmysqlUserFlowStatus stubs the async flow-query API with a fixed
// status. A wait case whose stub reports the terminal state on the FIRST
// call finishes immediately; a stub that never reports it ends at the short
// user timeout instead of the minute-scale default. Stub errors are plain
// fmt.Errorf values: tccommon.RetryError classifies them as non-retryable,
// so error cases finish immediately too.
func stubTdmysqlUserFlowStatus(patches *gomonkey.Patches, tdmysqlClient *tdmysqlv20211122.Client, status string) {
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeFlow", func(req *tdmysqlv20211122.DescribeFlowRequest) (*tdmysqlv20211122.DescribeFlowResponse, error) {
		return &tdmysqlv20211122.DescribeFlowResponse{Response: &tdmysqlv20211122.DescribeFlowResponseParams{
			Status:    tdmysqlUserStr(status),
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})
}

// tdmysqlUserTestContext wires the resource to a configured client; the SDK
// client methods are stubbed per test case.
type tdmysqlUserTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newTdmysqlUserTestContext(t *testing.T) *tdmysqlUserTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svtdmysql.NewTdmysqlUserResource()
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

	return &tdmysqlUserTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *tdmysqlUserTestContext) plan(t *testing.T, m svtdmysql.TdmysqlUserResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *tdmysqlUserTestContext) state(t *testing.T, m svtdmysql.TdmysqlUserResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// TestTdmysqlUserCreateSuccess verifies that Create sends all request fields,
// tracks the returned composite id and refreshes the state.
func TestTdmysqlUserCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.CreateUsersRequest
	patches.ApplyMethodFunc(tdmysqlClient, "CreateUsers", func(req *tdmysqlv20211122.CreateUsersRequest) (*tdmysqlv20211122.CreateUsersResponse, error) {
		captured = req

		return &tdmysqlv20211122.CreateUsersResponse{Response: &tdmysqlv20211122.CreateUsersResponseParams{
			FlowId:    tdmysqlUserFlowId(),
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserFlowStatus(patches, tdmysqlClient, svtdmysql.FLOW_STATUS_SUCCESS)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUsers", func(req *tdmysqlv20211122.DescribeUsersRequest) (*tdmysqlv20211122.DescribeUsersResponse, error) {
		return &tdmysqlv20211122.DescribeUsersResponse{Response: &tdmysqlv20211122.DescribeUsersResponseParams{
			Users:     []*tdmysqlv20211122.UserInfo{tdmysqlUserTestItem()},
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})

	model := tdmysqlUserTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateUsers request is nil")
	}
	if captured.InstanceId == nil || *captured.InstanceId != tdmysqlUserTestInstanceId {
		t.Errorf("create request instance_id = %v, want %q", captured.InstanceId, tdmysqlUserTestInstanceId)
	}
	if len(captured.Users) != 1 || captured.Users[0].UserName == nil || *captured.Users[0].UserName != tdmysqlUserTestUserName {
		t.Errorf("create request user_name = %v, want %q", captured.Users, tdmysqlUserTestUserName)
	}
	if len(captured.Users) != 1 || captured.Users[0].Host == nil || *captured.Users[0].Host != tdmysqlUserTestHost {
		t.Errorf("create request host = %v, want %q", captured.Users, tdmysqlUserTestHost)
	}
	if captured.Password == nil || *captured.Password != tdmysqlUserTestPassword {
		t.Errorf("create request password = %v, want the planned value", captured.Password)
	}

	var got svtdmysql.TdmysqlUserResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertTdmysqlUserModel(t, got)
}

// TestTdmysqlUserCreateEmptyId verifies that Create fails loudly when the
// request carries no addressing fields.
func TestTdmysqlUserCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "CreateUsers", func(req *tdmysqlv20211122.CreateUsersRequest) (*tdmysqlv20211122.CreateUsersResponse, error) {
		return &tdmysqlv20211122.CreateUsersResponse{Response: &tdmysqlv20211122.CreateUsersResponseParams{
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})

	model := tdmysqlUserTestModel()
	model.UserName = types.StringNull()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the request carries no addressing fields")
	}
}

// TestTdmysqlUserCreateWaitTimeout verifies that a flow which never reaches
// its terminal state ends the apply with an error once the (short) user
// timeout elapses — this exercises the `timeouts` block end to end.
func TestTdmysqlUserCreateWaitTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "CreateUsers", func(req *tdmysqlv20211122.CreateUsersRequest) (*tdmysqlv20211122.CreateUsersResponse, error) {
		return &tdmysqlv20211122.CreateUsersResponse{Response: &tdmysqlv20211122.CreateUsersResponseParams{
			FlowId:    tdmysqlUserFlowId(),
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserFlowStatus(patches, tdmysqlClient, "running")

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, tdmysqlUserTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, tdmysqlUserTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the flow never reaches its terminal state")
	}
}

// TestTdmysqlUserCreateApiError verifies API failures surface as errors.
func TestTdmysqlUserCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "CreateUsers", func(req *tdmysqlv20211122.CreateUsersRequest) (*tdmysqlv20211122.CreateUsersResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, tdmysqlUserTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, tdmysqlUserTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestTdmysqlUserReadSuccess verifies Read refreshes all leaf fields.
func TestTdmysqlUserReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)

	var gotInstanceId string
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUsers", func(req *tdmysqlv20211122.DescribeUsersRequest) (*tdmysqlv20211122.DescribeUsersResponse, error) {
		if req.InstanceId == nil {
			t.Errorf("expected the instance id to be passed")
		} else {
			gotInstanceId = *req.InstanceId
		}

		return &tdmysqlv20211122.DescribeUsersResponse{Response: &tdmysqlv20211122.DescribeUsersResponseParams{
			Users:     []*tdmysqlv20211122.UserInfo{tdmysqlUserTestItem()},
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlUserTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlUserTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotInstanceId != tdmysqlUserTestInstanceId {
		t.Errorf("DescribeUsers called with %q, want %q", gotInstanceId, tdmysqlUserTestInstanceId)
	}

	var got svtdmysql.TdmysqlUserResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertTdmysqlUserModel(t, got)
}

// TestTdmysqlUserReadDetectsDescriptionDrift verifies a user-configurable
// attribute is refreshed from the cloud: the state holds no description, the
// cloud now answers a value, and Read must record the cloud value (not keep
// the previous one), so a change made outside Terraform shows up in the next
// plan.
func TestTdmysqlUserReadDetectsDescriptionDrift(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUsers", func(req *tdmysqlv20211122.DescribeUsersRequest) (*tdmysqlv20211122.DescribeUsersResponse, error) {
		drifted := tdmysqlUserTestItem()
		drifted.Description = tdmysqlUserStr("drifted-description")

		return &tdmysqlv20211122.DescribeUsersResponse{Response: &tdmysqlv20211122.DescribeUsersResponseParams{
			Users:     []*tdmysqlv20211122.UserInfo{drifted},
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlUserTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlUserTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got svtdmysql.TdmysqlUserResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Description.ValueString() != "drifted-description" {
		t.Errorf("description = %q, want the cloud value %q", got.Description.ValueString(), "drifted-description")
	}
}

// TestTdmysqlUserReadNotFound verifies that a describe list without the
// tracked user drops the resource from state (confirmed not-found).
func TestTdmysqlUserReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUsers", func(req *tdmysqlv20211122.DescribeUsersRequest) (*tdmysqlv20211122.DescribeUsersResponse, error) {
		return &tdmysqlv20211122.DescribeUsersResponse{Response: &tdmysqlv20211122.DescribeUsersResponseParams{
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlUserTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlUserTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the user is confirmed gone")
	}
}

// TestTdmysqlUserUpdateTimeoutsOnly verifies that a plan differing from the
// state in the `timeouts` block alone calls NO write API: the cloud has no
// user-update API, so the read-back runs and the new timeouts are persisted.
func TestTdmysqlUserUpdateTimeoutsOnly(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)

	writeCalls := 0
	patches.ApplyMethodFunc(tdmysqlClient, "CreateUsers", func(req *tdmysqlv20211122.CreateUsersRequest) (*tdmysqlv20211122.CreateUsersResponse, error) {
		writeCalls++

		return nil, fmt.Errorf("no write API must be called on update")
	})
	patches.ApplyMethodFunc(tdmysqlClient, "DeleteUsers", func(req *tdmysqlv20211122.DeleteUsersRequest) (*tdmysqlv20211122.DeleteUsersResponse, error) {
		writeCalls++

		return nil, fmt.Errorf("no write API must be called on update")
	})
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUsers", func(req *tdmysqlv20211122.DescribeUsersRequest) (*tdmysqlv20211122.DescribeUsersResponse, error) {
		return &tdmysqlv20211122.DescribeUsersResponse{Response: &tdmysqlv20211122.DescribeUsersResponseParams{
			Users:     []*tdmysqlv20211122.UserInfo{tdmysqlUserTestItem()},
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})

	plannedModel := tdmysqlUserTestModel()
	plannedModel.Timeouts = resourcetimeouts.Value{Object: types.ObjectValueMust(
		tdmysqlUserTimeoutsAttrTypes(),
		map[string]attr.Value{
			"create": types.StringValue("5s"),
			"update": types.StringValue("7s"),
			"delete": types.StringValue("5s"),
		},
	)}

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, tdmysqlUserTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, tdmysqlUserTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if writeCalls != 0 {
		t.Errorf("a write API was called %d time(s) for a timeouts-only change", writeCalls)
	}

	var got svtdmysql.TdmysqlUserResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertTdmysqlUserModel(t, got)
}

// TestTdmysqlUserDeleteSuccess verifies Delete sends the delete request with
// the tracked user and then waits until the flow succeeds.
func TestTdmysqlUserDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.DeleteUsersRequest
	patches.ApplyMethodFunc(tdmysqlClient, "DeleteUsers", func(req *tdmysqlv20211122.DeleteUsersRequest) (*tdmysqlv20211122.DeleteUsersResponse, error) {
		captured = req

		return &tdmysqlv20211122.DeleteUsersResponse{Response: &tdmysqlv20211122.DeleteUsersResponseParams{
			FlowId:    tdmysqlUserFlowId(),
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})
	// The flow wait polls DescribeFlow; the terminal success status is
	// reported on the first call.
	stubTdmysqlUserFlowStatus(patches, tdmysqlClient, svtdmysql.FLOW_STATUS_SUCCESS)

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, tdmysqlUserTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteUsers request is nil")
	}
	if captured.InstanceId == nil || *captured.InstanceId != tdmysqlUserTestInstanceId {
		t.Errorf("delete request instance_id = %v, want %q", captured.InstanceId, tdmysqlUserTestInstanceId)
	}
	if len(captured.Users) != 1 || captured.Users[0].UserName == nil || *captured.Users[0].UserName != tdmysqlUserTestUserName {
		t.Errorf("delete request user_name = %v, want %q", captured.Users, tdmysqlUserTestUserName)
	}
	if len(captured.Users) != 1 || captured.Users[0].Host == nil || *captured.Users[0].Host != tdmysqlUserTestHost {
		t.Errorf("delete request host = %v, want %q", captured.Users, tdmysqlUserTestHost)
	}
}

// TestTdmysqlUserUpdatePasswordReset verifies that an in-place password
// change calls ResetUsersPassword with the tracked user and the planned
// plain-text password (EncryptedPassword stays unset), then waits until the
// flow succeeds.
func TestTdmysqlUserUpdatePasswordReset(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserTestContext(t)
	tdmysqlClient := patchTdmysqlUserSdkClient(patches, ctx.client)

	newPassword := "TF@example-67890"

	var captured *tdmysqlv20211122.ResetUsersPasswordRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ResetUsersPassword", func(req *tdmysqlv20211122.ResetUsersPasswordRequest) (*tdmysqlv20211122.ResetUsersPasswordResponse, error) {
		captured = req

		return &tdmysqlv20211122.ResetUsersPasswordResponse{Response: &tdmysqlv20211122.ResetUsersPasswordResponseParams{
			FlowId:    tdmysqlUserFlowId(),
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})
	// The flow wait polls DescribeFlow; the terminal success status is
	// reported on the first call.
	stubTdmysqlUserFlowStatus(patches, tdmysqlClient, svtdmysql.FLOW_STATUS_SUCCESS)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUsers", func(req *tdmysqlv20211122.DescribeUsersRequest) (*tdmysqlv20211122.DescribeUsersResponse, error) {
		return &tdmysqlv20211122.DescribeUsersResponse{Response: &tdmysqlv20211122.DescribeUsersResponseParams{
			Users:     []*tdmysqlv20211122.UserInfo{tdmysqlUserTestItem()},
			RequestId: tdmysqlUserStr("fake-request-id"),
		}}, nil
	})

	plannedModel := tdmysqlUserTestModel()
	plannedModel.Password = types.StringValue(newPassword)

	var updateResp resource.UpdateResponse
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, tdmysqlUserTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ResetUsersPassword request is nil")
	}
	if captured.InstanceId == nil || *captured.InstanceId != tdmysqlUserTestInstanceId {
		t.Errorf("reset request instance_id = %v, want %q", captured.InstanceId, tdmysqlUserTestInstanceId)
	}
	if len(captured.Users) != 1 || captured.Users[0].UserName == nil || *captured.Users[0].UserName != tdmysqlUserTestUserName {
		t.Errorf("reset request user_name = %v, want %q", captured.Users, tdmysqlUserTestUserName)
	}
	if len(captured.Users) != 1 || captured.Users[0].Host == nil || *captured.Users[0].Host != tdmysqlUserTestHost {
		t.Errorf("reset request host = %v, want %q", captured.Users, tdmysqlUserTestHost)
	}
	if len(captured.Users) != 1 || captured.Users[0].Password == nil || *captured.Users[0].Password != newPassword {
		t.Errorf("reset request password = %v, want the planned value", captured.Users)
	}
	if len(captured.Users) == 1 && captured.Users[0].EncryptedPassword != nil {
		t.Errorf("reset request encrypted_password = %v, want unset", *captured.Users[0].EncryptedPassword)
	}

	var got svtdmysql.TdmysqlUserResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Password.ValueString() != newPassword {
		t.Errorf("state password = %q, want %q", got.Password.ValueString(), newPassword)
	}
}

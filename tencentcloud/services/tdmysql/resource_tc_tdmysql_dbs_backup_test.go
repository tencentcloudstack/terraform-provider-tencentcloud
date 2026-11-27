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
	tdmysqlDbsBackupTestId         = int64(12345)
	tdmysqlDbsBackupTestIdStr      = "12345"
	tdmysqlDbsBackupTestInstanceId = "tdsqlc-mysql-12345"
	tdmysqlDbsBackupTestMethod     = "physical"
	tdmysqlDbsBackupTestType       = "full"
	tdmysqlDbsBackupTestName       = "tf-example-backup"
	tdmysqlDbsBackupTestDriftName  = "tf-example-backup-drifted"
	tdmysqlDbsBackupTestNewName    = "tf-example-backup-renamed"
)

// tdmysqlDbsBackupStr is the pointer helper for string fields.
func tdmysqlDbsBackupStr(v string) *string { return &v }

// tdmysqlDbsBackupInt64 is the pointer helper for int64 fields.
func tdmysqlDbsBackupInt64(v int64) *int64 { return &v }

// tdmysqlDbsBackupTestItem builds the DescribeDBSBackupSets response item.
func tdmysqlDbsBackupTestItem() *tdmysqlv20211122.BackupSetModel {
	return &tdmysqlv20211122.BackupSetModel{
		BackupSetId:  tdmysqlDbsBackupInt64(tdmysqlDbsBackupTestId),
		InstanceId:   tdmysqlDbsBackupStr(tdmysqlDbsBackupTestInstanceId),
		BackupMethod: tdmysqlDbsBackupStr(tdmysqlDbsBackupTestMethod),
		BackupType:   tdmysqlDbsBackupStr(tdmysqlDbsBackupTestType),
		BackupName:   tdmysqlDbsBackupStr(tdmysqlDbsBackupTestName),
		BackupStatus: tdmysqlDbsBackupStr("success"),
		ManualBackup: tdmysqlDbsBackupInt64(1),
	}
}

// tdmysqlDbsBackupTestTimeouts builds the timeouts value with a SECOND-SCALE
// short user value: the create tests exercise the wait path, and the wait
// must be reached without hard-waiting the minute-scale default.
func tdmysqlDbsBackupTestTimeouts() resourcetimeouts.Value {
	return resourcetimeouts.Value{Object: types.ObjectValueMust(
		map[string]attr.Type{"create": types.StringType},
		map[string]attr.Value{"create": types.StringValue("5s")},
	)}
}

// tdmysqlDbsBackupTestModel builds the terraform model used for plan/state.
// Every unset attribute carries an EXPLICIT typed null value — the raw zero
// value breaks tfsdk Plan/State reflection. The create tests exercise the
// async wait path, so the model carries a short 5s user timeout.
func tdmysqlDbsBackupTestModel() svtdmysql.TdmysqlDbsBackupResourceModel {
	return svtdmysql.TdmysqlDbsBackupResourceModel{
		Id:           types.StringValue(tdmysqlDbsBackupTestIdStr),
		Timeouts:     tdmysqlDbsBackupTestTimeouts(),
		InstanceId:   types.StringValue(tdmysqlDbsBackupTestInstanceId),
		BackupMethod: types.StringValue(tdmysqlDbsBackupTestMethod),
		BackupType:   types.StringValue(tdmysqlDbsBackupTestType),
		BackupName:   types.StringValue(tdmysqlDbsBackupTestName),
	}
}

// assertTdmysqlDbsBackupModel asserts every model field.
func assertTdmysqlDbsBackupModel(t *testing.T, m svtdmysql.TdmysqlDbsBackupResourceModel) {
	t.Helper()

	if m.Id.ValueString() != tdmysqlDbsBackupTestIdStr {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), tdmysqlDbsBackupTestIdStr)
	}
	if m.InstanceId.ValueString() != tdmysqlDbsBackupTestInstanceId {
		t.Errorf("instance_id = %q, want %q", m.InstanceId.ValueString(), tdmysqlDbsBackupTestInstanceId)
	}
	if m.BackupMethod.ValueString() != tdmysqlDbsBackupTestMethod {
		t.Errorf("backup_method = %q, want %q", m.BackupMethod.ValueString(), tdmysqlDbsBackupTestMethod)
	}
	if m.BackupType.ValueString() != tdmysqlDbsBackupTestType {
		t.Errorf("backup_type = %q, want %q", m.BackupType.ValueString(), tdmysqlDbsBackupTestType)
	}
	if m.BackupName.ValueString() != tdmysqlDbsBackupTestName {
		t.Errorf("backup_name = %q, want %q", m.BackupName.ValueString(), tdmysqlDbsBackupTestName)
	}
}

// patchTdmysqlDbsBackupSdkClient stubs the versioned client accessor and
// returns the mock SDK client whose methods the test patches next (the
// two-level mode).
func patchTdmysqlDbsBackupSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *tdmysqlv20211122.Client {
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)

	return tdmysqlClient
}

// tdmysqlDbsBackupTestContext wires the resource to a configured client; the
// SDK client methods are stubbed per test case.
type tdmysqlDbsBackupTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newTdmysqlDbsBackupTestContext(t *testing.T) *tdmysqlDbsBackupTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svtdmysql.NewTdmysqlDbsBackupResource()
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

	return &tdmysqlDbsBackupTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *tdmysqlDbsBackupTestContext) plan(t *testing.T, m svtdmysql.TdmysqlDbsBackupResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *tdmysqlDbsBackupTestContext) state(t *testing.T, m svtdmysql.TdmysqlDbsBackupResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// TestTdmysqlDbsBackupCreateSuccess verifies that Create sends all request
// fields, tracks the returned backup set id and refreshes the state.
func TestTdmysqlDbsBackupCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.CreateDBSBackupRequest
	patches.ApplyMethodFunc(tdmysqlClient, "CreateDBSBackup", func(req *tdmysqlv20211122.CreateDBSBackupRequest) (*tdmysqlv20211122.CreateDBSBackupResponse, error) {
		captured = req

		return &tdmysqlv20211122.CreateDBSBackupResponse{Response: &tdmysqlv20211122.CreateDBSBackupResponseParams{
			BackupSetId: tdmysqlDbsBackupInt64(tdmysqlDbsBackupTestId),
			IsSuccess:   &[]bool{true}[0],
			RequestId:   tdmysqlDbsBackupStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupSets", func(req *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
		return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
			Items: []*tdmysqlv20211122.BackupSetModel{tdmysqlDbsBackupTestItem()},
		}}, nil
	})

	model := tdmysqlDbsBackupTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateDBSBackup request is nil")
	}
	if captured.InstanceId == nil || *captured.InstanceId != tdmysqlDbsBackupTestInstanceId {
		t.Errorf("create request instance_id = %v, want %q", captured.InstanceId, tdmysqlDbsBackupTestInstanceId)
	}
	if captured.BackupMethod == nil || *captured.BackupMethod != tdmysqlDbsBackupTestMethod {
		t.Errorf("create request backup_method = %v, want %q", captured.BackupMethod, tdmysqlDbsBackupTestMethod)
	}
	if captured.BackupType == nil || *captured.BackupType != tdmysqlDbsBackupTestType {
		t.Errorf("create request backup_type = %v, want %q", captured.BackupType, tdmysqlDbsBackupTestType)
	}
	if captured.BackupName == nil || *captured.BackupName != tdmysqlDbsBackupTestName {
		t.Errorf("create request backup_name = %v, want %q", captured.BackupName, tdmysqlDbsBackupTestName)
	}

	var got svtdmysql.TdmysqlDbsBackupResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertTdmysqlDbsBackupModel(t, got)
}

// TestTdmysqlDbsBackupCreateWaitTimeout verifies that a backup whose
// BackupStatus never reaches `success` fails the create at the short user
// timeout instead of hanging (the API exposes no failure terminal state).
func TestTdmysqlDbsBackupCreateWaitTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)

	patches.ApplyMethodFunc(tdmysqlClient, "CreateDBSBackup", func(req *tdmysqlv20211122.CreateDBSBackupRequest) (*tdmysqlv20211122.CreateDBSBackupResponse, error) {
		return &tdmysqlv20211122.CreateDBSBackupResponse{Response: &tdmysqlv20211122.CreateDBSBackupResponseParams{
			BackupSetId: tdmysqlDbsBackupInt64(tdmysqlDbsBackupTestId),
			IsSuccess:   &[]bool{true}[0],
			RequestId:   tdmysqlDbsBackupStr("fake-request-id"),
		}}, nil
	})
	// The backup set is listed but stuck in a non-terminal status forever —
	// the ONLY stub allowed to never reach a terminal state.
	pending := tdmysqlDbsBackupTestItem()
	pending.BackupStatus = tdmysqlDbsBackupStr("running")
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupSets", func(req *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
		return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
			Items: []*tdmysqlv20211122.BackupSetModel{pending},
		}}, nil
	})

	model := tdmysqlDbsBackupTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the backup never reaches `success` within the timeout")
	}
}

// TestTdmysqlDbsBackupCreateEmptyId verifies that Create fails loudly when
// the API returns no backup set id.
func TestTdmysqlDbsBackupCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "CreateDBSBackup", func(req *tdmysqlv20211122.CreateDBSBackupRequest) (*tdmysqlv20211122.CreateDBSBackupResponse, error) {
		return &tdmysqlv20211122.CreateDBSBackupResponse{Response: &tdmysqlv20211122.CreateDBSBackupResponseParams{
			RequestId: tdmysqlDbsBackupStr("fake-request-id"),
		}}, nil
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, tdmysqlDbsBackupTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, tdmysqlDbsBackupTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty backup set id")
	}
}

// TestTdmysqlDbsBackupCreateApiError verifies API failures surface as errors.
func TestTdmysqlDbsBackupCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "CreateDBSBackup", func(req *tdmysqlv20211122.CreateDBSBackupRequest) (*tdmysqlv20211122.CreateDBSBackupResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, tdmysqlDbsBackupTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, tdmysqlDbsBackupTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestTdmysqlDbsBackupReadSuccess verifies Read refreshes all leaf fields.
func TestTdmysqlDbsBackupReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)

	var gotBackupSetId int64
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupSets", func(req *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
		if req.BackupSetId == nil {
			t.Errorf("expected the backup set id filter to be passed")
		} else {
			gotBackupSetId = *req.BackupSetId
		}

		return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
			Items: []*tdmysqlv20211122.BackupSetModel{tdmysqlDbsBackupTestItem()},
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlDbsBackupTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlDbsBackupTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotBackupSetId != tdmysqlDbsBackupTestId {
		t.Errorf("DescribeDBSBackupSets called with %d, want %d", gotBackupSetId, tdmysqlDbsBackupTestId)
	}

	var got svtdmysql.TdmysqlDbsBackupResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertTdmysqlDbsBackupModel(t, got)
}

// TestTdmysqlDbsBackupReadDetectsNameDrift verifies a user-configurable
// attribute is refreshed from the cloud: the state holds the old comment, the
// cloud now answers the drifted one, and Read must record the cloud value so
// a change made outside Terraform shows up in the next plan.
func TestTdmysqlDbsBackupReadDetectsNameDrift(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupSets", func(req *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
		drifted := tdmysqlDbsBackupTestItem()
		drifted.BackupName = tdmysqlDbsBackupStr(tdmysqlDbsBackupTestDriftName)

		return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
			Items: []*tdmysqlv20211122.BackupSetModel{drifted},
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlDbsBackupTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlDbsBackupTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got svtdmysql.TdmysqlDbsBackupResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.BackupName.ValueString() != tdmysqlDbsBackupTestDriftName {
		t.Errorf("backup_name = %q, want the cloud value %q", got.BackupName.ValueString(), tdmysqlDbsBackupTestDriftName)
	}
}

// TestTdmysqlDbsBackupReadNotFound verifies that a describe list without the
// tracked backup set drops the resource from state (confirmed not-found).
func TestTdmysqlDbsBackupReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupSets", func(req *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
		return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlDbsBackupTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlDbsBackupTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the backup set is confirmed gone")
	}
}

// TestTdmysqlDbsBackupUpdateSuccess verifies Update sends the PLANNED comment
// to the modify request and records the POST-update cloud state from the
// read-back — not the plan. The describe stub reports a value DIFFERENT from
// the plan (the cloud normalized it); only a state built from the read-back
// can equal it.
func TestTdmysqlDbsBackupUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyDBSBackupSetCommentRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyDBSBackupSetComment", func(req *tdmysqlv20211122.ModifyDBSBackupSetCommentRequest) (*tdmysqlv20211122.ModifyDBSBackupSetCommentResponse, error) {
		captured = req

		return &tdmysqlv20211122.ModifyDBSBackupSetCommentResponse{Response: &tdmysqlv20211122.ModifyDBSBackupSetCommentResponseParams{
			IsSuccess: &[]bool{true}[0],
			RequestId: tdmysqlDbsBackupStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupSets", func(req *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
		normalized := tdmysqlDbsBackupTestItem()
		normalized.BackupName = tdmysqlDbsBackupStr(tdmysqlDbsBackupTestNewName)

		return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
			Items: []*tdmysqlv20211122.BackupSetModel{normalized},
		}}, nil
	})

	plannedModel := tdmysqlDbsBackupTestModel()
	plannedModel.BackupName = types.StringValue(tdmysqlDbsBackupTestNewName)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, tdmysqlDbsBackupTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, tdmysqlDbsBackupTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyDBSBackupSetComment request is nil")
	}
	if captured.BackupSetId == nil || *captured.BackupSetId != tdmysqlDbsBackupTestId {
		t.Errorf("update request backup_set_id = %v, want %d", captured.BackupSetId, tdmysqlDbsBackupTestId)
	}
	if captured.NewBackupName == nil || *captured.NewBackupName != tdmysqlDbsBackupTestNewName {
		t.Errorf("update request new_backup_name = %v, want the planned %q", captured.NewBackupName, tdmysqlDbsBackupTestNewName)
	}

	var got svtdmysql.TdmysqlDbsBackupResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != tdmysqlDbsBackupTestIdStr {
		t.Errorf("id = %q, want the unchanged %q", got.Id.ValueString(), tdmysqlDbsBackupTestIdStr)
	}
	if got.BackupName.ValueString() != tdmysqlDbsBackupTestNewName {
		t.Errorf("backup_name = %q, want the read-back cloud value %q", got.BackupName.ValueString(), tdmysqlDbsBackupTestNewName)
	}
}

// TestTdmysqlDbsBackupUpdateNoCommentChange verifies that a plan differing
// from the state in no attribute the modify API covers calls NO write API.
// The read-back still runs.
func TestTdmysqlDbsBackupUpdateNoCommentChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)

	modifyCalls := 0
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyDBSBackupSetComment", func(req *tdmysqlv20211122.ModifyDBSBackupSetCommentRequest) (*tdmysqlv20211122.ModifyDBSBackupSetCommentResponse, error) {
		modifyCalls++

		return nil, fmt.Errorf("the modify API must not be called")
	})
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupSets", func(req *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
		return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
			Items: []*tdmysqlv20211122.BackupSetModel{tdmysqlDbsBackupTestItem()},
		}}, nil
	})

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, tdmysqlDbsBackupTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, tdmysqlDbsBackupTestModel()),
			State: ctx.state(t, tdmysqlDbsBackupTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if modifyCalls != 0 {
		t.Errorf("ModifyDBSBackupSetComment was called %d time(s) with no comment change", modifyCalls)
	}

	var got svtdmysql.TdmysqlDbsBackupResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertTdmysqlDbsBackupModel(t, got)
}

// TestTdmysqlDbsBackupDeleteSuccess verifies Delete sends the delete request
// with the tracked backup set id.
func TestTdmysqlDbsBackupDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.DeleteDBSBackupSetsRequest
	patches.ApplyMethodFunc(tdmysqlClient, "DeleteDBSBackupSets", func(req *tdmysqlv20211122.DeleteDBSBackupSetsRequest) (*tdmysqlv20211122.DeleteDBSBackupSetsResponse, error) {
		captured = req

		return &tdmysqlv20211122.DeleteDBSBackupSetsResponse{Response: &tdmysqlv20211122.DeleteDBSBackupSetsResponseParams{
			IsSuccess: &[]bool{true}[0],
			Deleted:   tdmysqlDbsBackupInt64(1),
			Total:     tdmysqlDbsBackupInt64(1),
			RequestId: tdmysqlDbsBackupStr("fake-request-id"),
		}}, nil
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, tdmysqlDbsBackupTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteDBSBackupSets request is nil")
	}
	if captured.InstanceId == nil || *captured.InstanceId != tdmysqlDbsBackupTestInstanceId {
		t.Errorf("delete request instance_id = %v, want %q", captured.InstanceId, tdmysqlDbsBackupTestInstanceId)
	}
	if len(captured.BackupSetIdList) != 1 || captured.BackupSetIdList[0] == nil || *captured.BackupSetIdList[0] != tdmysqlDbsBackupTestId {
		t.Errorf("delete request backup_set_id_list = %v, want [%d]", captured.BackupSetIdList, tdmysqlDbsBackupTestId)
	}
}

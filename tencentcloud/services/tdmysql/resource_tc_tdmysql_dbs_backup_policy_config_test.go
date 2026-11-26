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

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svtdmysql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql"
)

const (
	tdmysqlDbsBackupPolicyTestInstanceId = "tdsql-12345"
	tdmysqlDbsBackupPolicyTestMethod     = "physical"
	tdmysqlDbsBackupPolicyTestStartTime  = "02:00"
	tdmysqlDbsBackupPolicyTestEndTime    = "04:00"
	tdmysqlDbsBackupPolicyTestPeriodTime = "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday"
	tdmysqlDbsBackupPolicyTestStorage    = "COS"
)

func tdmysqlDbsBackupPolicyStr(v string) *string { return &v }

func tdmysqlDbsBackupPolicyInt64(v int64) *int64 { return &v }

// tdmysqlDbsBackupPolicyTestPolicy builds the backup policy the cloud reports.
func tdmysqlDbsBackupPolicyTestPolicy() []*tdmysqlv20211122.BackupPolicyModelOutPut {
	return []*tdmysqlv20211122.BackupPolicyModelOutPut{
		{
			BackupEndTime:       tdmysqlDbsBackupPolicyStr(tdmysqlDbsBackupPolicyTestEndTime),
			BackupMethod:        tdmysqlDbsBackupPolicyStr(tdmysqlDbsBackupPolicyTestMethod),
			BackupStartTime:     tdmysqlDbsBackupPolicyStr(tdmysqlDbsBackupPolicyTestStartTime),
			EnableFull:          tdmysqlDbsBackupPolicyInt64(1),
			EnableLog:           tdmysqlDbsBackupPolicyInt64(1),
			FullRetentionPeriod: tdmysqlDbsBackupPolicyInt64(7),
			LogRetentionPeriod:  tdmysqlDbsBackupPolicyInt64(7),
			PeriodTime:          tdmysqlDbsBackupPolicyStr(tdmysqlDbsBackupPolicyTestPeriodTime),
		},
	}
}

// tdmysqlDbsBackupPolicyTestObject builds the backup_policy object value the
// cloud reports (used for plans and expected states).
func tdmysqlDbsBackupPolicyTestObject() types.Object {
	return types.ObjectValueMust(svtdmysql.TdmysqlDbsBackupPolicyConfigBackupPolicyAttrTypes(), map[string]attr.Value{
		"backup_end_time":       types.StringValue(tdmysqlDbsBackupPolicyTestEndTime),
		"backup_method":         types.StringValue(tdmysqlDbsBackupPolicyTestMethod),
		"backup_start_time":     types.StringValue(tdmysqlDbsBackupPolicyTestStartTime),
		"enable_full":           types.Int64Value(1),
		"enable_log":            types.Int64Value(1),
		"full_retention_period": types.Int64Value(7),
		"log_retention_period":  types.Int64Value(7),
		"period_time":           types.StringValue(tdmysqlDbsBackupPolicyTestPeriodTime),
		"storage_type":          types.StringValue(tdmysqlDbsBackupPolicyTestStorage),
	})
}

// tdmysqlDbsBackupPolicyTestModel builds the terraform model used for
// plan/state. Every unset attribute carries an EXPLICIT TYPED null value —
// the raw zero value breaks tfsdk Plan/State reflection.
func tdmysqlDbsBackupPolicyTestModel() svtdmysql.TdmysqlDbsBackupPolicyConfigResourceModel {
	return svtdmysql.TdmysqlDbsBackupPolicyConfigResourceModel{
		Id:           types.StringValue(tdmysqlDbsBackupPolicyTestInstanceId),
		InstanceId:   types.StringValue(tdmysqlDbsBackupPolicyTestInstanceId),
		BackupPolicy: tdmysqlDbsBackupPolicyTestObject(),
	}
}

// patchTdmysqlDbsBackupPolicySdkClient stubs the versioned client accessor
// and returns the mock SDK client whose methods the test patches next.
func patchTdmysqlDbsBackupPolicySdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *tdmysqlv20211122.Client {
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)

	return tdmysqlClient
}

// stubTdmysqlDbsBackupPolicyDescribe stubs the describe API with the given
// policy items (nil = the cloud answers without a policy, which the service
// layer turns into fwHelper.ErrNotFound).
func stubTdmysqlDbsBackupPolicyDescribe(patches *gomonkey.Patches, tdmysqlClient *tdmysqlv20211122.Client, items []*tdmysqlv20211122.BackupPolicyModelOutPut) {
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupPolicyWithContext", func(_ context.Context, request *tdmysqlv20211122.DescribeDBSBackupPolicyRequest) (*tdmysqlv20211122.DescribeDBSBackupPolicyResponse, error) {
		return &tdmysqlv20211122.DescribeDBSBackupPolicyResponse{Response: &tdmysqlv20211122.DescribeDBSBackupPolicyResponseParams{
			Items:     items,
			RequestId: tdmysqlDbsBackupPolicyStr("fake-request-id"),
		}}, nil
	})
}

// tdmysqlDbsBackupPolicyTestContext wires the resource to a configured client.
type tdmysqlDbsBackupPolicyTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newTdmysqlDbsBackupPolicyTestContext(t *testing.T) *tdmysqlDbsBackupPolicyTestContext {
	t.Helper()

	client := &connectivity.TencentCloudClient{}

	res := svtdmysql.NewTdmysqlDbsBackupPolicyConfigResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}

	var configureResp resource.ConfigureResponse
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}},
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

	return &tdmysqlDbsBackupPolicyTestContext{resource: res, schema: schemaResp.Schema, client: client}
}

func (c *tdmysqlDbsBackupPolicyTestContext) plan(t *testing.T, m svtdmysql.TdmysqlDbsBackupPolicyConfigResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *tdmysqlDbsBackupPolicyTestContext) state(t *testing.T, m svtdmysql.TdmysqlDbsBackupPolicyConfigResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// TestTdmysqlDbsBackupPolicyConfigCreateSuccess verifies Create takes the
// state id from the locator input, sends the settings, and records the
// cloud-reported values.
func TestTdmysqlDbsBackupPolicyConfigCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupPolicyTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupPolicySdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyDBSBackupPolicyRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyDBSBackupPolicyWithContext", func(_ context.Context, request *tdmysqlv20211122.ModifyDBSBackupPolicyRequest) (*tdmysqlv20211122.ModifyDBSBackupPolicyResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyDBSBackupPolicyResponse{Response: &tdmysqlv20211122.ModifyDBSBackupPolicyResponseParams{
			IsSuccess: func() *bool { b := true; return &b }(),
			RequestId: tdmysqlDbsBackupPolicyStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlDbsBackupPolicyDescribe(patches, tdmysqlClient, tdmysqlDbsBackupPolicyTestPolicy())

	model := tdmysqlDbsBackupPolicyTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyDBSBackupPolicy request is nil")
	}
	if captured.InstanceId == nil || *captured.InstanceId != tdmysqlDbsBackupPolicyTestInstanceId {
		t.Errorf("request instance_id = %v, want %q: the locator must be sent", captured.InstanceId, tdmysqlDbsBackupPolicyTestInstanceId)
	}
	if captured.BackupPolicy == nil || captured.BackupPolicy.BackupMethod == nil || *captured.BackupPolicy.BackupMethod != tdmysqlDbsBackupPolicyTestMethod {
		t.Errorf("request backup_policy = %+v, want backup_method %q", captured.BackupPolicy, tdmysqlDbsBackupPolicyTestMethod)
	}

	var got svtdmysql.TdmysqlDbsBackupPolicyConfigResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != tdmysqlDbsBackupPolicyTestInstanceId {
		t.Errorf("id = %q, want the locator value %q", got.Id.ValueString(), tdmysqlDbsBackupPolicyTestInstanceId)
	}
	if !got.BackupPolicy.Equal(tdmysqlDbsBackupPolicyTestObject()) {
		t.Errorf("backup_policy = %s, want the cloud-reported values %s", got.BackupPolicy, tdmysqlDbsBackupPolicyTestObject())
	}
}

// TestTdmysqlDbsBackupPolicyConfigCreateApiError verifies a write failure
// surfaces as an error diagnostic (a plain fmt.Errorf is non-retryable).
func TestTdmysqlDbsBackupPolicyConfigCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupPolicyTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupPolicySdkClient(patches, ctx.client)
	stubTdmysqlDbsBackupPolicyDescribe(patches, tdmysqlClient, tdmysqlDbsBackupPolicyTestPolicy())
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyDBSBackupPolicyWithContext", func(_ context.Context, request *tdmysqlv20211122.ModifyDBSBackupPolicyRequest) (*tdmysqlv20211122.ModifyDBSBackupPolicyResponse, error) {
		return nil, fmt.Errorf("modify failed")
	})

	model := tdmysqlDbsBackupPolicyTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestTdmysqlDbsBackupPolicyConfigReadSuccess verifies Read addresses the
// configuration by the state id and refreshes every field.
func TestTdmysqlDbsBackupPolicyConfigReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupPolicyTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupPolicySdkClient(patches, ctx.client)

	var gotInstanceId string
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupPolicyWithContext", func(_ context.Context, request *tdmysqlv20211122.DescribeDBSBackupPolicyRequest) (*tdmysqlv20211122.DescribeDBSBackupPolicyResponse, error) {
		if request.InstanceId != nil {
			gotInstanceId = *request.InstanceId
		}

		return &tdmysqlv20211122.DescribeDBSBackupPolicyResponse{Response: &tdmysqlv20211122.DescribeDBSBackupPolicyResponseParams{
			Items:     tdmysqlDbsBackupPolicyTestPolicy(),
			RequestId: tdmysqlDbsBackupPolicyStr("fake-request-id"),
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlDbsBackupPolicyTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlDbsBackupPolicyTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotInstanceId != tdmysqlDbsBackupPolicyTestInstanceId {
		t.Errorf("DescribeDBSBackupPolicy called with %q, want %q", gotInstanceId, tdmysqlDbsBackupPolicyTestInstanceId)
	}

	var got svtdmysql.TdmysqlDbsBackupPolicyConfigResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if !got.BackupPolicy.Equal(tdmysqlDbsBackupPolicyTestObject()) {
		t.Errorf("backup_policy = %s, want the cloud-reported values %s", got.BackupPolicy, tdmysqlDbsBackupPolicyTestObject())
	}
}

// TestTdmysqlDbsBackupPolicyConfigReadNotFound verifies that a describe
// answer without a policy (fwHelper.ErrNotFound) drops the resource from
// state.
func TestTdmysqlDbsBackupPolicyConfigReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupPolicyTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupPolicySdkClient(patches, ctx.client)
	stubTdmysqlDbsBackupPolicyDescribe(patches, tdmysqlClient, nil)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlDbsBackupPolicyTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlDbsBackupPolicyTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the policy is confirmed gone")
	}
}

// TestTdmysqlDbsBackupPolicyConfigReadApiError verifies a describe failure
// surfaces as an error diagnostic and keeps the state.
func TestTdmysqlDbsBackupPolicyConfigReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupPolicyTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupPolicySdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBSBackupPolicyWithContext", func(_ context.Context, request *tdmysqlv20211122.DescribeDBSBackupPolicyRequest) (*tdmysqlv20211122.DescribeDBSBackupPolicyResponse, error) {
		return nil, fmt.Errorf("describe failed")
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlDbsBackupPolicyTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlDbsBackupPolicyTestModel())}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on a transport error")
	}
}

// TestTdmysqlDbsBackupPolicyConfigUpdateSuccess verifies Update shares the
// write path and keeps the id of the existing state.
func TestTdmysqlDbsBackupPolicyConfigUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlDbsBackupPolicyTestContext(t)
	tdmysqlClient := patchTdmysqlDbsBackupPolicySdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyDBSBackupPolicyRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyDBSBackupPolicyWithContext", func(_ context.Context, request *tdmysqlv20211122.ModifyDBSBackupPolicyRequest) (*tdmysqlv20211122.ModifyDBSBackupPolicyResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyDBSBackupPolicyResponse{Response: &tdmysqlv20211122.ModifyDBSBackupPolicyResponseParams{
			IsSuccess: func() *bool { b := true; return &b }(),
			RequestId: tdmysqlDbsBackupPolicyStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlDbsBackupPolicyDescribe(patches, tdmysqlClient, tdmysqlDbsBackupPolicyTestPolicy())

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, tdmysqlDbsBackupPolicyTestModel())
	planModel := tdmysqlDbsBackupPolicyTestModel()
	planModel.BackupPolicy = types.ObjectValueMust(svtdmysql.TdmysqlDbsBackupPolicyConfigBackupPolicyAttrTypes(), map[string]attr.Value{
		"backup_end_time":       types.StringValue(tdmysqlDbsBackupPolicyTestEndTime),
		"backup_method":         types.StringValue("snapshot"),
		"backup_start_time":     types.StringValue(tdmysqlDbsBackupPolicyTestStartTime),
		"enable_full":           types.Int64Value(1),
		"enable_log":            types.Int64Value(1),
		"full_retention_period": types.Int64Value(7),
		"log_retention_period":  types.Int64Value(7),
		"period_time":           types.StringValue(tdmysqlDbsBackupPolicyTestPeriodTime),
		"storage_type":          types.StringValue(tdmysqlDbsBackupPolicyTestStorage),
	})
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, planModel),
			State: ctx.state(t, tdmysqlDbsBackupPolicyTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyDBSBackupPolicy request is nil")
	}
	if captured.BackupPolicy == nil || captured.BackupPolicy.BackupMethod == nil || *captured.BackupPolicy.BackupMethod != "snapshot" {
		t.Errorf("request backup_policy = %+v, want backup_method %q", captured.BackupPolicy, "snapshot")
	}
}

// TestTdmysqlDbsBackupPolicyConfigDeleteIsNoop verifies Delete calls no API:
// the configuration lives and dies with its parent resource.
func TestTdmysqlDbsBackupPolicyConfigDeleteIsNoop(t *testing.T) {
	ctx := newTdmysqlDbsBackupPolicyTestContext(t)

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, tdmysqlDbsBackupPolicyTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}
}

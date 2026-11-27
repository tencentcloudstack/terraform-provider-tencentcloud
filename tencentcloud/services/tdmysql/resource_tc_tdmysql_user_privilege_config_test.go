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
	svctdmysql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql"
)

const (
	tdmysqlUserPrivilegeConfigTestInstanceId = "tdsql3-5baee8df"
	tdmysqlUserPrivilegeConfigTestUserName   = "tf-example-user"
	tdmysqlUserPrivilegeConfigTestHost       = "localhost"
	tdmysqlUserPrivilegeConfigTestComposite  = tdmysqlUserPrivilegeConfigTestInstanceId + "#" + tdmysqlUserPrivilegeConfigTestUserName + "#" + tdmysqlUserPrivilegeConfigTestHost
	tdmysqlUserPrivilegeConfigTestDatabase   = "tf-example-db"
	tdmysqlUserPrivilegeConfigTestTable      = "tf-example-table"
)

func tdmysqlUserPrivilegeConfigStr(v string) *string { return &v }

// tdmysqlUserPrivilegeConfigTestPrivileges builds the DescribeUserPrivileges response payload.
func tdmysqlUserPrivilegeConfigTestPrivileges() []*string {
	return []*string{
		tdmysqlUserPrivilegeConfigStr("SELECT"),
		tdmysqlUserPrivilegeConfigStr("INSERT"),
	}
}

// tdmysqlUserPrivilegeConfigTestGlobalPrivileges builds the global_privileges
// list value used by the default model.
func tdmysqlUserPrivilegeConfigTestGlobalPrivileges() types.List {
	return types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("SELECT"),
		types.StringValue("INSERT"),
	})
}

// tdmysqlUserPrivilegeConfigTestDatabasePrivileges builds a one-entry
// database_privileges list value.
func tdmysqlUserPrivilegeConfigTestDatabasePrivileges() types.List {
	entry, diags := types.ObjectValue(svctdmysql.TdmysqlUserPrivilegeConfigDatabasePrivilegesAttrTypes(), map[string]attr.Value{
		"database":   types.StringValue(tdmysqlUserPrivilegeConfigTestDatabase),
		"privileges": tdmysqlUserPrivilegeConfigTestGlobalPrivileges(),
	})
	if diags.HasError() {
		panic(fmt.Sprintf("failed to build the database_privileges entry: %v", diags.Errors()))
	}

	return types.ListValueMust(
		types.ListType{ElemType: types.ObjectType{AttrTypes: svctdmysql.TdmysqlUserPrivilegeConfigDatabasePrivilegesAttrTypes()}},
		[]attr.Value{entry},
	)
}

// tdmysqlUserPrivilegeConfigTestTablePrivileges builds a one-entry
// table_privileges list value.
func tdmysqlUserPrivilegeConfigTestTablePrivileges() types.List {
	entry, diags := types.ObjectValue(svctdmysql.TdmysqlUserPrivilegeConfigTablePrivilegesAttrTypes(), map[string]attr.Value{
		"database":   types.StringValue(tdmysqlUserPrivilegeConfigTestDatabase),
		"table":      types.StringValue(tdmysqlUserPrivilegeConfigTestTable),
		"privileges": tdmysqlUserPrivilegeConfigTestGlobalPrivileges(),
	})
	if diags.HasError() {
		panic(fmt.Sprintf("failed to build the table_privileges entry: %v", diags.Errors()))
	}

	return types.ListValueMust(
		types.ListType{ElemType: types.ObjectType{AttrTypes: svctdmysql.TdmysqlUserPrivilegeConfigTablePrivilegesAttrTypes()}},
		[]attr.Value{entry},
	)
}

func tdmysqlUserPrivilegeConfigTestDatabasePrivilegesNull() types.List {
	return types.ListNull(types.ObjectType{AttrTypes: svctdmysql.TdmysqlUserPrivilegeConfigDatabasePrivilegesAttrTypes()})
}

func tdmysqlUserPrivilegeConfigTestTablePrivilegesNull() types.List {
	return types.ListNull(types.ObjectType{AttrTypes: svctdmysql.TdmysqlUserPrivilegeConfigTablePrivilegesAttrTypes()})
}

// tdmysqlUserPrivilegeConfigTestModel builds the terraform model used for plan/state.
// Every attribute carries an explicit typed value: the raw zero value breaks
// tfsdk Plan/State reflection with "Value Conversion Error". The default model
// manages the GLOBAL privilege group only.
func tdmysqlUserPrivilegeConfigTestModel() svctdmysql.TdmysqlUserPrivilegeConfigResourceModel {
	return svctdmysql.TdmysqlUserPrivilegeConfigResourceModel{
		Id:                 types.StringValue(tdmysqlUserPrivilegeConfigTestComposite),
		InstanceId:         types.StringValue(tdmysqlUserPrivilegeConfigTestInstanceId),
		UserName:           types.StringValue(tdmysqlUserPrivilegeConfigTestUserName),
		Host:               types.StringValue(tdmysqlUserPrivilegeConfigTestHost),
		GlobalPrivileges:   tdmysqlUserPrivilegeConfigTestGlobalPrivileges(),
		DatabasePrivileges: tdmysqlUserPrivilegeConfigTestDatabasePrivilegesNull(),
		TablePrivileges:    tdmysqlUserPrivilegeConfigTestTablePrivilegesNull(),
	}
}

// patchTdmysqlUserPrivilegeConfigSdkClient stubs the versioned client accessor and
// returns the mock SDK client whose methods the test patches next.
func patchTdmysqlUserPrivilegeConfigSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *tdmysqlv20211122.Client {
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)

	return tdmysqlClient
}

// stubTdmysqlUserPrivilegeConfigDescribe stubs the describe API with the given
// privileges (nil = the cloud answers without any privilege item, which the
// service layer turns into fwHelper.ErrNotFound). The addressed request is
// captured for scope assertions.
func stubTdmysqlUserPrivilegeConfigDescribe(patches *gomonkey.Patches, tdmysqlClient *tdmysqlv20211122.Client, privileges []*string) **tdmysqlv20211122.DescribeUserPrivilegesRequest {
	captured := new(*tdmysqlv20211122.DescribeUserPrivilegesRequest)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUserPrivileges", func(request *tdmysqlv20211122.DescribeUserPrivilegesRequest) (*tdmysqlv20211122.DescribeUserPrivilegesResponse, error) {
		*captured = request

		return &tdmysqlv20211122.DescribeUserPrivilegesResponse{Response: &tdmysqlv20211122.DescribeUserPrivilegesResponseParams{
			Privileges: privileges,
			RequestId:  tdmysqlUserPrivilegeConfigStr("fake-request-id"),
		}}, nil
	})

	return captured
}

// tdmysqlUserPrivilegeConfigTestContext wires the resource to a configured client.
type tdmysqlUserPrivilegeConfigTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newTdmysqlUserPrivilegeConfigTestContext(t *testing.T) *tdmysqlUserPrivilegeConfigTestContext {
	t.Helper()

	client := &connectivity.TencentCloudClient{}

	res := svctdmysql.NewTdmysqlUserPrivilegeConfigResource()
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

	return &tdmysqlUserPrivilegeConfigTestContext{resource: res, schema: schemaResp.Schema, client: client}
}

func (c *tdmysqlUserPrivilegeConfigTestContext) plan(t *testing.T, m svctdmysql.TdmysqlUserPrivilegeConfigResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *tdmysqlUserPrivilegeConfigTestContext) state(t *testing.T, m svctdmysql.TdmysqlUserPrivilegeConfigResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// TestTdmysqlUserPrivilegeConfigCreateSuccess verifies Create sends the locator
// and the planned global privilege group to the write API, and records the
// cloud-reported values in the state.
func TestTdmysqlUserPrivilegeConfigCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeConfigTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeConfigSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyUserPrivilegesRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyUserPrivilegesResponse{Response: &tdmysqlv20211122.ModifyUserPrivilegesResponseParams{
			RequestId: tdmysqlUserPrivilegeConfigStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserPrivilegeConfigDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeConfigTestPrivileges())

	model := tdmysqlUserPrivilegeConfigTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyUserPrivileges request is nil")
	}
	if captured.InstanceId == nil || *captured.InstanceId != tdmysqlUserPrivilegeConfigTestInstanceId {
		t.Errorf("request instance_id = %v, want %q: the locator must be sent", captured.InstanceId, tdmysqlUserPrivilegeConfigTestInstanceId)
	}
	if len(captured.Users) != 1 || captured.Users[0].UserName == nil || *captured.Users[0].UserName != tdmysqlUserPrivilegeConfigTestUserName ||
		captured.Users[0].Host == nil || *captured.Users[0].Host != tdmysqlUserPrivilegeConfigTestHost {
		t.Errorf("request users = %+v, want user_name %q with host %q", captured.Users, tdmysqlUserPrivilegeConfigTestUserName, tdmysqlUserPrivilegeConfigTestHost)
	}
	if len(captured.GlobalPrivileges) != 2 {
		t.Errorf("request global_privileges = %+v, want the 2 planned items", captured.GlobalPrivileges)
	}
	if captured.DatabasePrivileges != nil || captured.TablePrivileges != nil {
		t.Errorf("request groups = %+v, want the unmanaged groups omitted", captured)
	}

	var got svctdmysql.TdmysqlUserPrivilegeConfigResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != tdmysqlUserPrivilegeConfigTestComposite {
		t.Errorf("id = %q, want the composite locator %q", got.Id.ValueString(), tdmysqlUserPrivilegeConfigTestComposite)
	}
	if len(got.GlobalPrivileges.Elements()) != 2 {
		t.Errorf("global_privileges = %v, want 2 cloud-reported items", got.GlobalPrivileges)
	}
}

// TestTdmysqlUserPrivilegeConfigCreateApiError verifies a write failure surfaces as
// an error diagnostic (a plain fmt.Errorf is non-retryable: no waiting).
func TestTdmysqlUserPrivilegeConfigCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeConfigTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeConfigSdkClient(patches, ctx.client)
	stubTdmysqlUserPrivilegeConfigDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeConfigTestPrivileges())
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		return nil, fmt.Errorf("modify failed")
	})

	model := tdmysqlUserPrivilegeConfigTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestTdmysqlUserPrivilegeConfigCreateDatabaseScope verifies the database-level
// write group carries the planned entries and the unmanaged global group is
// omitted from the request.
func TestTdmysqlUserPrivilegeConfigCreateDatabaseScope(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeConfigTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeConfigSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyUserPrivilegesRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyUserPrivilegesResponse{Response: &tdmysqlv20211122.ModifyUserPrivilegesResponseParams{
			RequestId: tdmysqlUserPrivilegeConfigStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserPrivilegeConfigDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeConfigTestPrivileges())

	model := tdmysqlUserPrivilegeConfigTestModel()
	model.GlobalPrivileges = types.ListNull(types.StringType)
	model.DatabasePrivileges = tdmysqlUserPrivilegeConfigTestDatabasePrivileges()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyUserPrivileges request is nil")
	}
	if captured.GlobalPrivileges != nil {
		t.Errorf("request global_privileges = %+v, want nil for the unmanaged group", captured.GlobalPrivileges)
	}
	if len(captured.DatabasePrivileges) != 1 ||
		captured.DatabasePrivileges[0].Database == nil || *captured.DatabasePrivileges[0].Database != tdmysqlUserPrivilegeConfigTestDatabase ||
		len(captured.DatabasePrivileges[0].Privileges) != 2 {
		t.Errorf("request database_privileges = %+v, want one planned entry for db %q", captured.DatabasePrivileges, tdmysqlUserPrivilegeConfigTestDatabase)
	}

	var got svctdmysql.TdmysqlUserPrivilegeConfigResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.DatabasePrivileges.IsNull() || len(got.DatabasePrivileges.Elements()) != 1 {
		t.Errorf("database_privileges = %v, want the refreshed one-entry list", got.DatabasePrivileges)
	}
}

// TestTdmysqlUserPrivilegeConfigCreateTableScope verifies the table-level write
// group carries the planned entries.
func TestTdmysqlUserPrivilegeConfigCreateTableScope(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeConfigTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeConfigSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyUserPrivilegesRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyUserPrivilegesResponse{Response: &tdmysqlv20211122.ModifyUserPrivilegesResponseParams{
			RequestId: tdmysqlUserPrivilegeConfigStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserPrivilegeConfigDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeConfigTestPrivileges())

	model := tdmysqlUserPrivilegeConfigTestModel()
	model.GlobalPrivileges = types.ListNull(types.StringType)
	model.TablePrivileges = tdmysqlUserPrivilegeConfigTestTablePrivileges()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyUserPrivileges request is nil")
	}
	if len(captured.TablePrivileges) != 1 ||
		captured.TablePrivileges[0].Database == nil || *captured.TablePrivileges[0].Database != tdmysqlUserPrivilegeConfigTestDatabase ||
		captured.TablePrivileges[0].Table == nil || *captured.TablePrivileges[0].Table != tdmysqlUserPrivilegeConfigTestTable ||
		len(captured.TablePrivileges[0].Privileges) != 2 {
		t.Errorf("request table_privileges = %+v, want one planned entry for %q.%q", captured.TablePrivileges, tdmysqlUserPrivilegeConfigTestDatabase, tdmysqlUserPrivilegeConfigTestTable)
	}
}

// TestTdmysqlUserPrivilegeConfigReadSuccess verifies Read refreshes the managed
// global group through the read API's global addressing (db_name `*`).
func TestTdmysqlUserPrivilegeConfigReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeConfigTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeConfigSdkClient(patches, ctx.client)

	captured := stubTdmysqlUserPrivilegeConfigDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeConfigTestPrivileges())

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlUserPrivilegeConfigTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlUserPrivilegeConfigTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if *captured == nil {
		t.Fatal("DescribeUserPrivileges was not called")
	}
	if (*captured).DbName == nil || *(*captured).DbName != "*" {
		t.Errorf("read addressing db_name = %v, want %q for the global group", (*captured).DbName, "*")
	}

	var got svctdmysql.TdmysqlUserPrivilegeConfigResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.InstanceId.ValueString() != tdmysqlUserPrivilegeConfigTestInstanceId {
		t.Errorf("instance_id = %q, want %q", got.InstanceId.ValueString(), tdmysqlUserPrivilegeConfigTestInstanceId)
	}
	if got.UserName.ValueString() != tdmysqlUserPrivilegeConfigTestUserName {
		t.Errorf("user_name = %q, want %q", got.UserName.ValueString(), tdmysqlUserPrivilegeConfigTestUserName)
	}
	if got.Host.ValueString() != tdmysqlUserPrivilegeConfigTestHost {
		t.Errorf("host = %q, want %q", got.Host.ValueString(), tdmysqlUserPrivilegeConfigTestHost)
	}
	if len(got.GlobalPrivileges.Elements()) != 2 {
		t.Errorf("global_privileges = %v, want 2 cloud-reported items", got.GlobalPrivileges)
	}
}

// TestTdmysqlUserPrivilegeConfigReadNotFound verifies that a describe answer
// without any privilege item (fwHelper.ErrNotFound) on EVERY managed scope
// drops the resource from state.
func TestTdmysqlUserPrivilegeConfigReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeConfigTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeConfigSdkClient(patches, ctx.client)
	stubTdmysqlUserPrivilegeConfigDescribe(patches, tdmysqlClient, nil)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlUserPrivilegeConfigTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlUserPrivilegeConfigTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the privileges are confirmed gone")
	}
}

// TestTdmysqlUserPrivilegeConfigReadApiError verifies a describe failure surfaces as
// an error diagnostic and keeps the state.
func TestTdmysqlUserPrivilegeConfigReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeConfigTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeConfigSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUserPrivileges", func(request *tdmysqlv20211122.DescribeUserPrivilegesRequest) (*tdmysqlv20211122.DescribeUserPrivilegesResponse, error) {
		return nil, fmt.Errorf("describe failed")
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlUserPrivilegeConfigTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlUserPrivilegeConfigTestModel())}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on a transport error")
	}
}

// TestTdmysqlUserPrivilegeConfigUpdateSuccess verifies Update shares the write path
// and keeps the id of the existing state.
func TestTdmysqlUserPrivilegeConfigUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeConfigTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeConfigSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyUserPrivilegesRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyUserPrivilegesResponse{Response: &tdmysqlv20211122.ModifyUserPrivilegesResponseParams{
			RequestId: tdmysqlUserPrivilegeConfigStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserPrivilegeConfigDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeConfigTestPrivileges())

	model := tdmysqlUserPrivilegeConfigTestModel()
	model.GlobalPrivileges = types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("SELECT"),
		types.StringValue("UPDATE"),
		types.StringValue("DELETE"),
	})

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, tdmysqlUserPrivilegeConfigTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, model),
			State: ctx.state(t, tdmysqlUserPrivilegeConfigTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyUserPrivileges request is nil")
	}
	if len(captured.GlobalPrivileges) != 3 {
		t.Errorf("request global_privileges = %+v, want the 3 planned items", captured.GlobalPrivileges)
	}

	var got svctdmysql.TdmysqlUserPrivilegeConfigResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != tdmysqlUserPrivilegeConfigTestComposite {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), tdmysqlUserPrivilegeConfigTestComposite)
	}
	if len(got.GlobalPrivileges.Elements()) != 2 {
		t.Errorf("global_privileges = %v, want the 2 cloud-reported items", got.GlobalPrivileges)
	}
}

// TestTdmysqlUserPrivilegeConfigDeleteIsNoop verifies Delete calls no API: the
// privileges live and die with the parent instance and user.
func TestTdmysqlUserPrivilegeConfigDeleteIsNoop(t *testing.T) {
	ctx := newTdmysqlUserPrivilegeConfigTestContext(t)

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, tdmysqlUserPrivilegeConfigTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}
}

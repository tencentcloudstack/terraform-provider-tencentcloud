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
	tdmysqlUserPrivilegeTestInstanceId = "tdsql3-5baee8df"
	tdmysqlUserPrivilegeTestHost       = "localhost"
	tdmysqlUserPrivilegeTestUserName   = "tf-example-user"
	tdmysqlUserPrivilegeTestDbName     = "*"
	tdmysqlUserPrivilegeTestComposite  = tdmysqlUserPrivilegeTestInstanceId + "#" + tdmysqlUserPrivilegeTestHost + "#" + tdmysqlUserPrivilegeTestUserName + "#" + tdmysqlUserPrivilegeTestDbName
)

func tdmysqlUserPrivilegeStr(v string) *string { return &v }

// tdmysqlUserPrivilegeTestPrivileges builds the DescribeUserPrivileges response payload.
func tdmysqlUserPrivilegeTestPrivileges() []*string {
	return []*string{
		tdmysqlUserPrivilegeStr("SELECT"),
		tdmysqlUserPrivilegeStr("INSERT"),
	}
}

// tdmysqlUserPrivilegeTestModel builds the terraform model used for plan/state.
// Every attribute carries an explicit typed value: the raw zero value breaks
// tfsdk Plan/State reflection with "Value Conversion Error".
func tdmysqlUserPrivilegeTestModel() svctdmysql.TdmysqlUserPrivilegeResourceModel {
	return svctdmysql.TdmysqlUserPrivilegeResourceModel{
		Id:         types.StringValue(tdmysqlUserPrivilegeTestComposite),
		InstanceId: types.StringValue(tdmysqlUserPrivilegeTestInstanceId),
		Host:       types.StringValue(tdmysqlUserPrivilegeTestHost),
		UserName:   types.StringValue(tdmysqlUserPrivilegeTestUserName),
		DbName:     types.StringValue(tdmysqlUserPrivilegeTestDbName),
		Object:     types.StringNull(),
		ObjectType: types.StringNull(),
		ColName:    types.StringNull(),
		Privileges: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("SELECT"),
			types.StringValue("INSERT"),
		}),
	}
}

// patchTdmysqlUserPrivilegeSdkClient stubs the versioned client accessor and
// returns the mock SDK client whose methods the test patches next.
func patchTdmysqlUserPrivilegeSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *tdmysqlv20211122.Client {
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)

	return tdmysqlClient
}

// stubTdmysqlUserPrivilegeDescribe stubs the describe API with the given
// privileges (nil = the cloud answers without any privilege item, which the
// service layer turns into fwHelper.ErrNotFound).
func stubTdmysqlUserPrivilegeDescribe(patches *gomonkey.Patches, tdmysqlClient *tdmysqlv20211122.Client, privileges []*string) {
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUserPrivileges", func(request *tdmysqlv20211122.DescribeUserPrivilegesRequest) (*tdmysqlv20211122.DescribeUserPrivilegesResponse, error) {
		return &tdmysqlv20211122.DescribeUserPrivilegesResponse{Response: &tdmysqlv20211122.DescribeUserPrivilegesResponseParams{
			Privileges: privileges,
			RequestId:  tdmysqlUserPrivilegeStr("fake-request-id"),
		}}, nil
	})
}

// tdmysqlUserPrivilegeTestContext wires the resource to a configured client.
type tdmysqlUserPrivilegeTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newTdmysqlUserPrivilegeTestContext(t *testing.T) *tdmysqlUserPrivilegeTestContext {
	t.Helper()

	client := &connectivity.TencentCloudClient{}

	res := svctdmysql.NewTdmysqlUserPrivilegeResource()
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

	return &tdmysqlUserPrivilegeTestContext{resource: res, schema: schemaResp.Schema, client: client}
}

func (c *tdmysqlUserPrivilegeTestContext) plan(t *testing.T, m svctdmysql.TdmysqlUserPrivilegeResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *tdmysqlUserPrivilegeTestContext) state(t *testing.T, m svctdmysql.TdmysqlUserPrivilegeResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// TestTdmysqlUserPrivilegeCreateSuccess verifies Create sends the locator and
// the privilege list to the write API (global scope, db_name = "*"), and
// records the cloud-reported values.
func TestTdmysqlUserPrivilegeCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyUserPrivilegesRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyUserPrivilegesResponse{Response: &tdmysqlv20211122.ModifyUserPrivilegesResponseParams{
			RequestId: tdmysqlUserPrivilegeStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserPrivilegeDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeTestPrivileges())

	model := tdmysqlUserPrivilegeTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyUserPrivileges request is nil")
	}
	if captured.InstanceId == nil || *captured.InstanceId != tdmysqlUserPrivilegeTestInstanceId {
		t.Errorf("request instance_id = %v, want %q: the locator must be sent", captured.InstanceId, tdmysqlUserPrivilegeTestInstanceId)
	}
	if len(captured.Users) != 1 || captured.Users[0].UserName == nil || *captured.Users[0].UserName != tdmysqlUserPrivilegeTestUserName ||
		captured.Users[0].Host == nil || *captured.Users[0].Host != tdmysqlUserPrivilegeTestHost {
		t.Errorf("request users = %+v, want user_name %q with host %q", captured.Users, tdmysqlUserPrivilegeTestUserName, tdmysqlUserPrivilegeTestHost)
	}
	if len(captured.GlobalPrivileges) != 2 {
		t.Errorf("request global_privileges = %+v, want 2 items (db_name is %q)", captured.GlobalPrivileges, tdmysqlUserPrivilegeTestDbName)
	}

	var got svctdmysql.TdmysqlUserPrivilegeResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != tdmysqlUserPrivilegeTestComposite {
		t.Errorf("id = %q, want the composite locator %q", got.Id.ValueString(), tdmysqlUserPrivilegeTestComposite)
	}
	if len(got.Privileges.Elements()) != 2 {
		t.Errorf("privileges = %v, want 2 cloud-reported items", got.Privileges)
	}
}

// TestTdmysqlUserPrivilegeCreateApiError verifies a write failure surfaces as
// an error diagnostic (a plain fmt.Errorf is non-retryable: no waiting).
func TestTdmysqlUserPrivilegeCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeSdkClient(patches, ctx.client)
	stubTdmysqlUserPrivilegeDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeTestPrivileges())
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		return nil, fmt.Errorf("modify failed")
	})

	model := tdmysqlUserPrivilegeTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestTdmysqlUserPrivilegeCreateDatabaseScope verifies the database-level
// write group is used when object_type is "*" on a concrete database.
func TestTdmysqlUserPrivilegeCreateDatabaseScope(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyUserPrivilegesRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyUserPrivilegesResponse{Response: &tdmysqlv20211122.ModifyUserPrivilegesResponseParams{
			RequestId: tdmysqlUserPrivilegeStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserPrivilegeDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeTestPrivileges())

	model := tdmysqlUserPrivilegeTestModel()
	model.DbName = types.StringValue("tf-example-db")
	model.ObjectType = types.StringValue("*")
	model.Id = types.StringValue(tdmysqlUserPrivilegeTestInstanceId + "#" + tdmysqlUserPrivilegeTestHost + "#" + tdmysqlUserPrivilegeTestUserName + "#tf-example-db")

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyUserPrivileges request is nil")
	}
	if len(captured.DatabasePrivileges) != 1 || captured.DatabasePrivileges[0].Database == nil || *captured.DatabasePrivileges[0].Database != "tf-example-db" {
		t.Errorf("request database_privileges = %+v, want one entry for db %q", captured.DatabasePrivileges, "tf-example-db")
	}
}

// TestTdmysqlUserPrivilegeCreateTableScope verifies the table-level write
// group is used for a concrete table object.
func TestTdmysqlUserPrivilegeCreateTableScope(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyUserPrivilegesRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyUserPrivilegesResponse{Response: &tdmysqlv20211122.ModifyUserPrivilegesResponseParams{
			RequestId: tdmysqlUserPrivilegeStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserPrivilegeDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeTestPrivileges())

	model := tdmysqlUserPrivilegeTestModel()
	model.DbName = types.StringValue("tf-example-db")
	model.Object = types.StringValue("tf-example-table")
	model.ObjectType = types.StringValue("table")
	model.Id = types.StringValue(tdmysqlUserPrivilegeTestInstanceId + "#" + tdmysqlUserPrivilegeTestHost + "#" + tdmysqlUserPrivilegeTestUserName + "#tf-example-db")

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyUserPrivileges request is nil")
	}
	if len(captured.TablePrivileges) != 1 || captured.TablePrivileges[0].Table == nil || *captured.TablePrivileges[0].Table != "tf-example-table" {
		t.Errorf("request table_privileges = %+v, want one entry for table %q", captured.TablePrivileges, "tf-example-table")
	}
}

// TestTdmysqlUserPrivilegeReadSuccess verifies Read addresses the privileges
// by the composite id and refreshes the privilege list.
func TestTdmysqlUserPrivilegeReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeSdkClient(patches, ctx.client)

	var gotInstanceId string
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUserPrivileges", func(request *tdmysqlv20211122.DescribeUserPrivilegesRequest) (*tdmysqlv20211122.DescribeUserPrivilegesResponse, error) {
		if request.InstanceId != nil {
			gotInstanceId = *request.InstanceId
		}

		return &tdmysqlv20211122.DescribeUserPrivilegesResponse{Response: &tdmysqlv20211122.DescribeUserPrivilegesResponseParams{
			Privileges: tdmysqlUserPrivilegeTestPrivileges(),
			RequestId:  tdmysqlUserPrivilegeStr("fake-request-id"),
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlUserPrivilegeTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlUserPrivilegeTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotInstanceId != tdmysqlUserPrivilegeTestInstanceId {
		t.Errorf("DescribeUserPrivileges called with %q, want %q", gotInstanceId, tdmysqlUserPrivilegeTestInstanceId)
	}

	var got svctdmysql.TdmysqlUserPrivilegeResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.InstanceId.ValueString() != tdmysqlUserPrivilegeTestInstanceId {
		t.Errorf("instance_id = %q, want %q", got.InstanceId.ValueString(), tdmysqlUserPrivilegeTestInstanceId)
	}
	if got.Host.ValueString() != tdmysqlUserPrivilegeTestHost {
		t.Errorf("host = %q, want %q", got.Host.ValueString(), tdmysqlUserPrivilegeTestHost)
	}
	if got.UserName.ValueString() != tdmysqlUserPrivilegeTestUserName {
		t.Errorf("user_name = %q, want %q", got.UserName.ValueString(), tdmysqlUserPrivilegeTestUserName)
	}
	if got.DbName.ValueString() != tdmysqlUserPrivilegeTestDbName {
		t.Errorf("db_name = %q, want %q", got.DbName.ValueString(), tdmysqlUserPrivilegeTestDbName)
	}
	if len(got.Privileges.Elements()) != 2 {
		t.Errorf("privileges = %v, want 2 items", got.Privileges)
	}
}

// TestTdmysqlUserPrivilegeReadNotFound verifies that a describe answer without
// any privilege item (fwHelper.ErrNotFound) drops the resource from state.
func TestTdmysqlUserPrivilegeReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeSdkClient(patches, ctx.client)
	stubTdmysqlUserPrivilegeDescribe(patches, tdmysqlClient, nil)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlUserPrivilegeTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlUserPrivilegeTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the privileges are confirmed gone")
	}
}

// TestTdmysqlUserPrivilegeReadApiError verifies a describe failure surfaces as
// an error diagnostic and keeps the state.
func TestTdmysqlUserPrivilegeReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeUserPrivileges", func(request *tdmysqlv20211122.DescribeUserPrivilegesRequest) (*tdmysqlv20211122.DescribeUserPrivilegesResponse, error) {
		return nil, fmt.Errorf("describe failed")
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, tdmysqlUserPrivilegeTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, tdmysqlUserPrivilegeTestModel())}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on a transport error")
	}
}

// TestTdmysqlUserPrivilegeUpdateSuccess verifies Update shares the write path
// and keeps the id of the existing state.
func TestTdmysqlUserPrivilegeUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newTdmysqlUserPrivilegeTestContext(t)
	tdmysqlClient := patchTdmysqlUserPrivilegeSdkClient(patches, ctx.client)

	var captured *tdmysqlv20211122.ModifyUserPrivilegesRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyUserPrivileges", func(request *tdmysqlv20211122.ModifyUserPrivilegesRequest) (*tdmysqlv20211122.ModifyUserPrivilegesResponse, error) {
		captured = request

		return &tdmysqlv20211122.ModifyUserPrivilegesResponse{Response: &tdmysqlv20211122.ModifyUserPrivilegesResponseParams{
			RequestId: tdmysqlUserPrivilegeStr("fake-request-id"),
		}}, nil
	})
	stubTdmysqlUserPrivilegeDescribe(patches, tdmysqlClient, tdmysqlUserPrivilegeTestPrivileges())

	model := tdmysqlUserPrivilegeTestModel()
	model.Privileges = types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("SELECT"),
		types.StringValue("UPDATE"),
		types.StringValue("DELETE"),
	})

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, tdmysqlUserPrivilegeTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, model),
			State: ctx.state(t, tdmysqlUserPrivilegeTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyUserPrivileges request is nil")
	}

	var got svctdmysql.TdmysqlUserPrivilegeResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != tdmysqlUserPrivilegeTestComposite {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), tdmysqlUserPrivilegeTestComposite)
	}
	if len(got.Privileges.Elements()) != 2 {
		t.Errorf("privileges = %v, want the 2 cloud-reported items", got.Privileges)
	}
}

// TestTdmysqlUserPrivilegeDeleteIsNoop verifies Delete calls no API: the
// privileges live and die with the parent instance and user.
func TestTdmysqlUserPrivilegeDeleteIsNoop(t *testing.T) {
	ctx := newTdmysqlUserPrivilegeTestContext(t)

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, tdmysqlUserPrivilegeTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}
}

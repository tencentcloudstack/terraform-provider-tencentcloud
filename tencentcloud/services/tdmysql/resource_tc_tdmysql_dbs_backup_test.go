package tdmysql_test

import (
	"context"
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
	tdmysqlTestInstanceId   = "tdsqlc-abc123"
	tdmysqlTestBackupSetId  = int64(10086)
	tdmysqlTestBackupId     = "tdsqlc-abc123#10086"
	tdmysqlTestBackupMethod = "physical"
	tdmysqlTestBackupType   = "full"
)

func tdmysqlString(value string) *string { return &value }

func tdmysqlInt64(value int64) *int64 { return &value }

func tdmysqlFloat64(value float64) *float64 { return &value }

// tdmysqlDbsBackupTestContext wires a resource instance to a patched SDK client.
type tdmysqlDbsBackupTestContext struct {
	resource      resource.Resource
	client        *tdmysqlv20211122.Client
	schema        schema.Schema
	itemAttrTypes map[string]attr.Type
}

func newTdmysqlDbsBackupTestContext(t *testing.T, patches *gomonkey.Patches) *tdmysqlDbsBackupTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	sdkClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseTdmysqlV20211122Client", sdkClient)

	res := svctdmysql.NewTdmysqlDbsBackupResource()
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

	itemsAttr, ok := schemaResp.Schema.Attributes["items"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a list nested attribute", "items")
	}
	listType, ok := itemsAttr.GetType().(types.ListType)
	if !ok {
		t.Fatalf("attribute items does not have a list type")
	}
	objectType, ok := listType.ElemType.(types.ObjectType)
	if !ok {
		t.Fatalf("attribute items does not have an object element type")
	}

	return &tdmysqlDbsBackupTestContext{
		resource:      res,
		client:        sdkClient,
		schema:        schemaResp.Schema,
		itemAttrTypes: objectType.AttrTypes,
	}
}

func (c *tdmysqlDbsBackupTestContext) plan(t *testing.T, model svctdmysql.TdmysqlDbsBackupResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *tdmysqlDbsBackupTestContext) state(t *testing.T, model svctdmysql.TdmysqlDbsBackupResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func (c *tdmysqlDbsBackupTestContext) emptyItems(t *testing.T) types.List {
	t.Helper()

	return types.ListNull(types.ObjectType{AttrTypes: c.itemAttrTypes})
}

func tdmysqlSampleBackupSetItem() *tdmysqlv20211122.BackupSetModel {
	return &tdmysqlv20211122.BackupSetModel{
		BackupDuration: tdmysqlInt64(120),
		BackupMethod:   tdmysqlString(tdmysqlTestBackupMethod),
		BackupName:     tdmysqlString("tf-backup"),
		BackupProgress: tdmysqlFloat64(100.0),
		BackupSetId:    tdmysqlInt64(tdmysqlTestBackupSetId),
		BackupStatus:   tdmysqlString("success"),
		BackupType:     tdmysqlString(tdmysqlTestBackupType),
		DBVersion:      tdmysqlString("8.0"),
		EndTime:        tdmysqlString("2026-01-01 00:02:00"),
		EndTrxTime:     tdmysqlString("2026-01-01 00:01:59"),
		ErrorMessage:   tdmysqlString(""),
		ExpiredTime:    tdmysqlString("2026-02-01 00:00:00"),
		FileName:       tdmysqlString("backup-10086.file"),
		FileSize:       tdmysqlInt64(1024),
		InstanceId:     tdmysqlString(tdmysqlTestInstanceId),
		ManualBackup:   tdmysqlInt64(1),
		StartTime:      tdmysqlString("2026-01-01 00:00:00"),
	}
}

func tdmysqlSampleModel(t *testing.T, ctx *tdmysqlDbsBackupTestContext) svctdmysql.TdmysqlDbsBackupResourceModel {
	t.Helper()

	return svctdmysql.TdmysqlDbsBackupResourceModel{
		Id:           types.StringValue(tdmysqlTestBackupId),
		InstanceId:   types.StringValue(tdmysqlTestInstanceId),
		BackupMethod: types.StringValue(tdmysqlTestBackupMethod),
		BackupType:   types.StringValue(tdmysqlTestBackupType),
		BackupName:   types.StringValue("tf-backup"),
		BackupSetId:  types.Int64Value(tdmysqlTestBackupSetId),
		Items:        ctx.emptyItems(t),
	}
}

func assertNoErrorDiagnostics(t *testing.T, diags interface{ HasError() bool }) {
	t.Helper()

	if diags.HasError() {
		t.Fatalf("unexpected error diagnostics")
	}
}

// TestTdmysqlDbsBackupSchema verifies the attribute set and required flags.
func TestTdmysqlDbsBackupSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newTdmysqlDbsBackupTestContext(t, patches)
	top := testContext.schema.Attributes

	if len(top) != 7 {
		t.Fatalf("schema has %d attributes, want 7", len(top))
	}

	required := []string{"instance_id", "backup_method", "backup_type"}
	for _, name := range required {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsRequired() {
			t.Fatalf("attribute %q must be required", name)
		}
	}

	optional := []string{"backup_name"}
	for _, name := range optional {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsOptional() || attribute.IsRequired() {
			t.Fatalf("attribute %q must be optional", name)
		}
	}

	computed := []string{"id", "backup_set_id", "items"}
	for _, name := range computed {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsComputed() {
			t.Fatalf("attribute %q must be computed", name)
		}
	}
}

// TestTdmysqlDbsBackupCreate verifies the create request and the read back state.
func TestTdmysqlDbsBackupCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newTdmysqlDbsBackupTestContext(t, patches)

	var capturedCreate *tdmysqlv20211122.CreateDBSBackupRequest
	var capturedDescribe *tdmysqlv20211122.DescribeDBSBackupSetsRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"CreateDBSBackupWithContext",
		func(_ context.Context, request *tdmysqlv20211122.CreateDBSBackupRequest) (*tdmysqlv20211122.CreateDBSBackupResponse, error) {
			capturedCreate = request
			return &tdmysqlv20211122.CreateDBSBackupResponse{
				Response: &tdmysqlv20211122.CreateDBSBackupResponseParams{
					BackupSetId: tdmysqlInt64(tdmysqlTestBackupSetId),
					IsSuccess:   tdmysqlBool(true),
				},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeDBSBackupSetsWithContext",
		func(_ context.Context, request *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
			capturedDescribe = request
			return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{
				Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
					Items:      []*tdmysqlv20211122.BackupSetModel{tdmysqlSampleBackupSetItem()},
					TotalCount: tdmysqlInt64(1),
				},
			}, nil
		},
	)

	model := tdmysqlSampleModel(t, testContext)
	model.Id = types.StringNull()
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	assertNoErrorDiagnostics(t, createResp.Diagnostics)

	if capturedCreate == nil {
		t.Fatalf("CreateDBSBackup was not called")
	}
	if capturedCreate.InstanceId == nil || *capturedCreate.InstanceId != tdmysqlTestInstanceId {
		t.Fatalf("CreateDBSBackup InstanceId is not %q", tdmysqlTestInstanceId)
	}
	if capturedCreate.BackupMethod == nil || *capturedCreate.BackupMethod != tdmysqlTestBackupMethod {
		t.Fatalf("CreateDBSBackup BackupMethod is not %q", tdmysqlTestBackupMethod)
	}
	if capturedCreate.BackupType == nil || *capturedCreate.BackupType != tdmysqlTestBackupType {
		t.Fatalf("CreateDBSBackup BackupType is not %q", tdmysqlTestBackupType)
	}
	if capturedCreate.BackupName == nil || *capturedCreate.BackupName != "tf-backup" {
		t.Fatalf("CreateDBSBackup BackupName is not tf-backup")
	}

	if capturedDescribe == nil {
		t.Fatalf("DescribeDBSBackupSets was not called")
	}
	if capturedDescribe.InstanceId == nil || *capturedDescribe.InstanceId != tdmysqlTestInstanceId {
		t.Fatalf("DescribeDBSBackupSets InstanceId is not %q", tdmysqlTestInstanceId)
	}
	if capturedDescribe.BackupSetId == nil || *capturedDescribe.BackupSetId != tdmysqlTestBackupSetId {
		t.Fatalf("DescribeDBSBackupSets BackupSetId is not %d", tdmysqlTestBackupSetId)
	}

	var state svctdmysql.TdmysqlDbsBackupResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read created state: %v", diags.Errors())
	}
	if state.Id.ValueString() != tdmysqlTestBackupId {
		t.Fatalf("state id is %q, want %q", state.Id.ValueString(), tdmysqlTestBackupId)
	}
	if state.BackupSetId.ValueInt64() != tdmysqlTestBackupSetId {
		t.Fatalf("state backup_set_id is %d, want %d", state.BackupSetId.ValueInt64(), tdmysqlTestBackupSetId)
	}

	var items []svctdmysql.TdmysqlDbsBackupItemModel
	diags := state.Items.ElementsAs(context.Background(), &items, false)
	if diags.HasError() || len(items) != 1 {
		t.Fatalf("state items must have 1 element, diags %v", diags.Errors())
	}
	item := items[0]
	if item.BackupDuration.ValueInt64() != 120 {
		t.Fatalf("items.backup_duration is %d, want 120", item.BackupDuration.ValueInt64())
	}
	if item.BackupMethod.ValueString() != tdmysqlTestBackupMethod {
		t.Fatalf("items.backup_method is %q", item.BackupMethod.ValueString())
	}
	if item.BackupName.ValueString() != "tf-backup" {
		t.Fatalf("items.backup_name is %q", item.BackupName.ValueString())
	}
	if item.BackupProgress.ValueFloat64() != 100.0 {
		t.Fatalf("items.backup_progress is %v", item.BackupProgress.ValueFloat64())
	}
	if item.BackupSetId.ValueInt64() != tdmysqlTestBackupSetId {
		t.Fatalf("items.backup_set_id is %d", item.BackupSetId.ValueInt64())
	}
	if item.BackupStatus.ValueString() != "success" {
		t.Fatalf("items.backup_status is %q", item.BackupStatus.ValueString())
	}
	if item.BackupType.ValueString() != tdmysqlTestBackupType {
		t.Fatalf("items.backup_type is %q", item.BackupType.ValueString())
	}
	if item.DBVersion.ValueString() != "8.0" {
		t.Fatalf("items.db_version is %q", item.DBVersion.ValueString())
	}
	if item.EndTime.ValueString() != "2026-01-01 00:02:00" {
		t.Fatalf("items.end_time is %q", item.EndTime.ValueString())
	}
	if item.EndTrxTime.ValueString() != "2026-01-01 00:01:59" {
		t.Fatalf("items.end_trx_time is %q", item.EndTrxTime.ValueString())
	}
	if item.ErrorMessage.ValueString() != "" {
		t.Fatalf("items.error_message is %q", item.ErrorMessage.ValueString())
	}
	if item.ExpiredTime.ValueString() != "2026-02-01 00:00:00" {
		t.Fatalf("items.expired_time is %q", item.ExpiredTime.ValueString())
	}
	if item.FileName.ValueString() != "backup-10086.file" {
		t.Fatalf("items.file_name is %q", item.FileName.ValueString())
	}
	if item.FileSize.ValueInt64() != 1024 {
		t.Fatalf("items.file_size is %d", item.FileSize.ValueInt64())
	}
	if item.InstanceId.ValueString() != tdmysqlTestInstanceId {
		t.Fatalf("items.instance_id is %q", item.InstanceId.ValueString())
	}
	if item.ManualBackup.ValueInt64() != 1 {
		t.Fatalf("items.manual_backup is %d", item.ManualBackup.ValueInt64())
	}
	if item.StartTime.ValueString() != "2026-01-01 00:00:00" {
		t.Fatalf("items.start_time is %q", item.StartTime.ValueString())
	}
}

func tdmysqlBool(value bool) *bool { return &value }

// TestTdmysqlDbsBackupCreateEmptyBackupSetId verifies that a nil BackupSetId in
// the create response is treated as an API failure.
func TestTdmysqlDbsBackupCreateEmptyBackupSetId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newTdmysqlDbsBackupTestContext(t, patches)
	patches.ApplyMethodFunc(
		testContext.client,
		"CreateDBSBackupWithContext",
		func(_ context.Context, _ *tdmysqlv20211122.CreateDBSBackupRequest) (*tdmysqlv20211122.CreateDBSBackupResponse, error) {
			return &tdmysqlv20211122.CreateDBSBackupResponse{
				Response: &tdmysqlv20211122.CreateDBSBackupResponseParams{},
			}, nil
		},
	)

	model := tdmysqlSampleModel(t, testContext)
	model.Id = types.StringNull()
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)

	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected error diagnostics for empty backup set id")
	}
}

// TestTdmysqlDbsBackupRead verifies the read refresh and request parameters.
func TestTdmysqlDbsBackupRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newTdmysqlDbsBackupTestContext(t, patches)

	var capturedDescribe *tdmysqlv20211122.DescribeDBSBackupSetsRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeDBSBackupSetsWithContext",
		func(_ context.Context, request *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
			capturedDescribe = request
			return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{
				Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
					Items:      []*tdmysqlv20211122.BackupSetModel{tdmysqlSampleBackupSetItem()},
					TotalCount: tdmysqlInt64(1),
				},
			}, nil
		},
	)

	model := tdmysqlSampleModel(t, testContext)
	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: readResp.State}, readResp)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	if capturedDescribe == nil {
		t.Fatalf("DescribeDBSBackupSets was not called")
	}
	if capturedDescribe.BackupSetId == nil || *capturedDescribe.BackupSetId != tdmysqlTestBackupSetId {
		t.Fatalf("DescribeDBSBackupSets BackupSetId is not %d", tdmysqlTestBackupSetId)
	}

	var state svctdmysql.TdmysqlDbsBackupResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags.Errors())
	}
	if len(state.Items.Elements()) != 1 {
		t.Fatalf("state items must have 1 element, got %d", len(state.Items.Elements()))
	}
	if state.BackupName.ValueString() != "tf-backup" {
		t.Fatalf("state backup_name is %q", state.BackupName.ValueString())
	}
}

// TestTdmysqlDbsBackupReadNotFound verifies that an empty describe result
// removes the resource from state.
func TestTdmysqlDbsBackupReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newTdmysqlDbsBackupTestContext(t, patches)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeDBSBackupSetsWithContext",
		func(_ context.Context, _ *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
			return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{
				Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
					Items:      []*tdmysqlv20211122.BackupSetModel{},
					TotalCount: tdmysqlInt64(0),
				},
			}, nil
		},
	)

	model := tdmysqlSampleModel(t, testContext)
	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: readResp.State}, readResp)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	if readResp.State.Raw.IsNull() {
		return
	}
	t.Fatalf("state must be removed when the backup set no longer exists")
}

// TestTdmysqlDbsBackupUpdate verifies the modify request when backup_name changes.
func TestTdmysqlDbsBackupUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newTdmysqlDbsBackupTestContext(t, patches)

	var capturedModify *tdmysqlv20211122.ModifyDBSBackupSetCommentRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"ModifyDBSBackupSetCommentWithContext",
		func(_ context.Context, request *tdmysqlv20211122.ModifyDBSBackupSetCommentRequest) (*tdmysqlv20211122.ModifyDBSBackupSetCommentResponse, error) {
			capturedModify = request
			return &tdmysqlv20211122.ModifyDBSBackupSetCommentResponse{
				Response: &tdmysqlv20211122.ModifyDBSBackupSetCommentResponseParams{
					IsSuccess: tdmysqlBool(true),
					Msg:       tdmysqlString("ok"),
				},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeDBSBackupSetsWithContext",
		func(_ context.Context, _ *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (*tdmysqlv20211122.DescribeDBSBackupSetsResponse, error) {
			item := tdmysqlSampleBackupSetItem()
			item.BackupName = tdmysqlString("renamed-backup")
			return &tdmysqlv20211122.DescribeDBSBackupSetsResponse{
				Response: &tdmysqlv20211122.DescribeDBSBackupSetsResponseParams{
					Items:      []*tdmysqlv20211122.BackupSetModel{item},
					TotalCount: tdmysqlInt64(1),
				},
			}, nil
		},
	)

	stateModel := tdmysqlSampleModel(t, testContext)
	planModel := tdmysqlSampleModel(t, testContext)
	planModel.BackupName = types.StringValue("renamed-backup")

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{Plan: testContext.plan(t, planModel), State: testContext.state(t, stateModel)},
		updateResp,
	)
	assertNoErrorDiagnostics(t, updateResp.Diagnostics)

	if capturedModify == nil {
		t.Fatalf("ModifyDBSBackupSetComment was not called")
	}
	if capturedModify.InstanceId == nil || *capturedModify.InstanceId != tdmysqlTestInstanceId {
		t.Fatalf("ModifyDBSBackupSetComment InstanceId is not %q", tdmysqlTestInstanceId)
	}
	if capturedModify.BackupSetId == nil || *capturedModify.BackupSetId != tdmysqlTestBackupSetId {
		t.Fatalf("ModifyDBSBackupSetComment BackupSetId is not %d", tdmysqlTestBackupSetId)
	}
	if capturedModify.NewBackupName == nil || *capturedModify.NewBackupName != "renamed-backup" {
		t.Fatalf("ModifyDBSBackupSetComment NewBackupName is not renamed-backup")
	}

	var state svctdmysql.TdmysqlDbsBackupResourceModel
	if diags := updateResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags.Errors())
	}
	if state.BackupName.ValueString() != "renamed-backup" {
		t.Fatalf("state backup_name is %q, want renamed-backup", state.BackupName.ValueString())
	}
	if state.Id.ValueString() != tdmysqlTestBackupId {
		t.Fatalf("state id is %q, want %q", state.Id.ValueString(), tdmysqlTestBackupId)
	}
}

// TestTdmysqlDbsBackupDelete verifies the delete request.
func TestTdmysqlDbsBackupDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newTdmysqlDbsBackupTestContext(t, patches)

	var capturedDelete *tdmysqlv20211122.DeleteDBSBackupSetsRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteDBSBackupSetsWithContext",
		func(_ context.Context, request *tdmysqlv20211122.DeleteDBSBackupSetsRequest) (*tdmysqlv20211122.DeleteDBSBackupSetsResponse, error) {
			capturedDelete = request
			return &tdmysqlv20211122.DeleteDBSBackupSetsResponse{
				Response: &tdmysqlv20211122.DeleteDBSBackupSetsResponseParams{
					Deleted:   tdmysqlInt64(1),
					IsSuccess: tdmysqlBool(true),
					Total:     tdmysqlInt64(1),
				},
			}, nil
		},
	)

	model := tdmysqlSampleModel(t, testContext)
	deleteResp := &resource.DeleteResponse{}
	testContext.resource.Delete(context.Background(), resource.DeleteRequest{State: testContext.state(t, model)}, deleteResp)
	assertNoErrorDiagnostics(t, deleteResp.Diagnostics)

	if capturedDelete == nil {
		t.Fatalf("DeleteDBSBackupSets was not called")
	}
	if capturedDelete.InstanceId == nil || *capturedDelete.InstanceId != tdmysqlTestInstanceId {
		t.Fatalf("DeleteDBSBackupSets InstanceId is not %q", tdmysqlTestInstanceId)
	}
	if len(capturedDelete.BackupSetIdList) != 1 || *capturedDelete.BackupSetIdList[0] != tdmysqlTestBackupSetId {
		t.Fatalf("DeleteDBSBackupSets BackupSetIdList is not [%d]", tdmysqlTestBackupSetId)
	}
}

// TestTdmysqlDbsBackupImportState verifies the composite id parsing.
func TestTdmysqlDbsBackupImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newTdmysqlDbsBackupTestContext(t, patches)

	importer, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: tdmysqlTestBackupId},
		importResp,
	)
	assertNoErrorDiagnostics(t, importResp.Diagnostics)

	var state svctdmysql.TdmysqlDbsBackupResourceModel
	if diags := importResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read imported state: %v", diags.Errors())
	}
	if state.Id.ValueString() != tdmysqlTestBackupId {
		t.Fatalf("imported id is %q", state.Id.ValueString())
	}
	if state.InstanceId.ValueString() != tdmysqlTestInstanceId {
		t.Fatalf("imported instance_id is %q", state.InstanceId.ValueString())
	}
	if state.BackupSetId.ValueInt64() != tdmysqlTestBackupSetId {
		t.Fatalf("imported backup_set_id is %d", state.BackupSetId.ValueInt64())
	}
}

// TestTdmysqlDbsBackupImportStateInvalid verifies that a malformed import id
// produces an error diagnostic.
func TestTdmysqlDbsBackupImportStateInvalid(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newTdmysqlDbsBackupTestContext(t, patches)

	importer, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: "not-a-composite-id"},
		importResp,
	)

	if !importResp.Diagnostics.HasError() {
		t.Fatalf("expected error diagnostics for malformed import id")
	}
}

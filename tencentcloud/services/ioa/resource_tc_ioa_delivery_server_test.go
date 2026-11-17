package ioa_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tftypes "github.com/hashicorp/terraform-plugin-go/tftypes"
	ioav20220601 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ioa/v20220601"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcioa "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ioa"
)

const (
	ioaTestServerId   = uint64(12345)
	ioaTestServerName = "tf-example-server"
)

func ioaStr(v string) *string { return &v }

func ioaInt64(v int64) *int64 { return &v }

func ioaUint64(v uint64) *uint64 { return &v }

func ioaBool(v bool) *bool { return &v }

func stringPtrValue(v string) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}

func int64PtrValue(v int64) types.Int64 { return types.Int64Value(v) }

func boolPtrValue(v bool) types.Bool { return types.BoolValue(v) }

// ioaDeliveryServerTestContext wires the resource to a patched SDK client.
type ioaDeliveryServerTestContext struct {
	resource resource.Resource
	client   *ioav20220601.Client
	schema   schema.Schema
}

func newIoaDeliveryServerTestContext(t *testing.T, patches *gomonkey.Patches) *ioaDeliveryServerTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	ioaClient := &ioav20220601.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseIoaV20220601Client", ioaClient)

	res := svcioa.NewIoaDeliveryServerResource()
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

	return &ioaDeliveryServerTestContext{resource: res, client: ioaClient, schema: schemaResp.Schema}
}

func (c *ioaDeliveryServerTestContext) plan(t *testing.T, model svcioa.IoaDeliveryServerResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *ioaDeliveryServerTestContext) state(t *testing.T, model svcioa.IoaDeliveryServerResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func ioaTestItem() *ioav20220601.DeliverServerItem {
	return &ioav20220601.DeliverServerItem{
		ID:                   ioaUint64(ioaTestServerId),
		Name:                 ioaStr(ioaTestServerName),
		Type:                 ioaInt64(4),
		Description:          ioaStr("managed by terraform"),
		TargetSyslogAddress:  ioaStr("10.0.0.1"),
		TargetSyslogProtocol: ioaStr("udp"),
		TargetSyslogPort:     ioaUint64(514),
		SyslogFacility:       ioaStr("local0"),
		TargetKafkaAddress:   ioaStr("10.0.0.2:9092"),
		TargetKafkaTopic:     ioaStr("ioa-log"),
		KafkaVersion:         ioaStr("2.8.0"),
		KafkaSaslUserName:    ioaStr("sasl-user"),
		KafkaSaslPassword:    ioaStr("sasl-pass"),
		TargetHttpUrl:        ioaStr("http://10.0.0.3/collect"),
		TargetHttpTimeout:    ioaUint64(3000),
		ConnectorGroupId:     ioaStr("cg-1"),
		ConnectorGroupName:   ioaStr("group-one"),
		KafkaSSLEnable:       ioaBool(true),
		SSLSkipVerify:        ioaBool(false),
	}
}

func ioaTestModel(id string) svcioa.IoaDeliveryServerResourceModel {
	return svcioa.IoaDeliveryServerResourceModel{
		Id: stringPtrValue(id),
		Data: &svcioa.IoaDeliveryServerDataModel{
			Name:                 stringPtrValue(ioaTestServerName),
			Type:                 int64PtrValue(4),
			Description:          stringPtrValue("managed by terraform"),
			TargetSyslogAddress:  stringPtrValue("10.0.0.1"),
			TargetSyslogProtocol: stringPtrValue("udp"),
			TargetSyslogPort:     int64PtrValue(514),
			SyslogFacility:       stringPtrValue("local0"),
			TargetKafkaAddress:   stringPtrValue("10.0.0.2:9092"),
			TargetKafkaTopic:     stringPtrValue("ioa-log"),
			KafkaVersion:         stringPtrValue("2.8.0"),
			KafkaSaslUserName:    stringPtrValue("sasl-user"),
			KafkaSaslPassword:    stringPtrValue("sasl-pass"),
			TargetHttpUrl:        stringPtrValue("http://10.0.0.3/collect"),
			TargetHttpTimeout:    int64PtrValue(3000),
			ConnectorGroupId:     stringPtrValue("cg-1"),
			ConnectorGroupName:   stringPtrValue("group-one"),
			KafkaSslEnable:       boolPtrValue(true),
			SslSkipVerify:        boolPtrValue(false),
		},
	}
}

func assertDeliveryServerModel(t *testing.T, model svcioa.IoaDeliveryServerResourceModel) {
	t.Helper()

	if model.Id.IsNull() || model.Id.ValueString() != "12345" {
		t.Errorf("id = %v, want 12345", model.Id.ValueString())
	}
	if model.Data == nil {
		t.Fatal("data is nil")
	}
	expectStr := map[string]func() string{
		"name":                   model.Data.Name.ValueString,
		"description":            model.Data.Description.ValueString,
		"target_syslog_address":  model.Data.TargetSyslogAddress.ValueString,
		"target_syslog_protocol": model.Data.TargetSyslogProtocol.ValueString,
		"syslog_facility":        model.Data.SyslogFacility.ValueString,
		"target_kafka_address":   model.Data.TargetKafkaAddress.ValueString,
		"target_kafka_topic":     model.Data.TargetKafkaTopic.ValueString,
		"kafka_version":          model.Data.KafkaVersion.ValueString,
		"kafka_sasl_user_name":   model.Data.KafkaSaslUserName.ValueString,
		"target_http_url":        model.Data.TargetHttpUrl.ValueString,
		"connector_group_id":     model.Data.ConnectorGroupId.ValueString,
		"connector_group_name":   model.Data.ConnectorGroupName.ValueString,
	}
	for name, getter := range expectStr {
		if getter() == "" {
			t.Errorf("data.%s is empty, want a non-empty value", name)
		}
	}
	// kafka_sasl_password is write-only and must never be persisted in state.
	if !model.Data.KafkaSaslPassword.IsNull() {
		t.Errorf("data.kafka_sasl_password = %v, want null (write-only)", model.Data.KafkaSaslPassword)
	}
	if model.Data.Name.ValueString() != ioaTestServerName {
		t.Errorf("data.name = %q, want %q", model.Data.Name.ValueString(), ioaTestServerName)
	}
	if model.Data.Type.ValueInt64() != 4 {
		t.Errorf("data.type = %d, want 4", model.Data.Type.ValueInt64())
	}
	if model.Data.TargetSyslogPort.ValueInt64() != 514 {
		t.Errorf("data.target_syslog_port = %d, want 514", model.Data.TargetSyslogPort.ValueInt64())
	}
	if model.Data.TargetHttpTimeout.ValueInt64() != 3000 {
		t.Errorf("data.target_http_timeout = %d, want 3000", model.Data.TargetHttpTimeout.ValueInt64())
	}
	if model.Data.Id.ValueInt64() != int64(ioaTestServerId) {
		t.Errorf("data.id = %d, want %d", model.Data.Id.ValueInt64(), ioaTestServerId)
	}
	if !model.Data.KafkaSslEnable.ValueBool() {
		t.Errorf("data.kafka_ssl_enable = false, want true")
	}
	if model.Data.SslSkipVerify.ValueBool() {
		t.Errorf("data.ssl_skip_verify = true, want false")
	}
}

// TestIoaDeliveryServerCreateSuccess verifies that Create sends the data item
// without a server id, reads back the generated id and persists the full state.
func TestIoaDeliveryServerCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *ioav20220601.ModifyDeliveryServerRequest
	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "ModifyDeliveryServerWithContext",
		func(_ context.Context, req *ioav20220601.ModifyDeliveryServerRequest) (*ioav20220601.ModifyDeliveryServerResponse, error) {
			captured = req
			return &ioav20220601.ModifyDeliveryServerResponse{
				Response: &ioav20220601.ModifyDeliveryServerResponseParams{
					Data: &ioav20220601.ModifyDeliveryServerData{ID: ioaUint64(ioaTestServerId)},
				},
			}, nil
		})
	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "DescribeDeliveryServerWithContext",
		func(_ context.Context, _ *ioav20220601.DescribeDeliveryServerRequest) (*ioav20220601.DescribeDeliveryServerResponse, error) {
			return &ioav20220601.DescribeDeliveryServerResponse{
				Response: &ioav20220601.DescribeDeliveryServerResponseParams{Data: ioaTestItem()},
			}, nil
		})

	ctx := newIoaDeliveryServerTestContext(t, patches)
	model := ioaTestModel("")
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil || captured.Data == nil {
		t.Fatal("ModifyDeliveryServer request or Data is nil")
	}
	if captured.Data.ID != nil {
		t.Errorf("create request Data.Id = %v, want nil (ignored on create)", *captured.Data.ID)
	}
	if captured.Data.Name == nil || *captured.Data.Name != ioaTestServerName {
		t.Errorf("create request Data.Name = %v, want %q", captured.Data.Name, ioaTestServerName)
	}

	var got svcioa.IoaDeliveryServerResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertDeliveryServerModel(t, got)
}

// TestIoaDeliveryServerCreateEmptyId verifies that Create fails loudly when
// the API returns no server id instead of writing an unmanageable resource.
func TestIoaDeliveryServerCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "ModifyDeliveryServerWithContext",
		func(_ context.Context, _ *ioav20220601.ModifyDeliveryServerRequest) (*ioav20220601.ModifyDeliveryServerResponse, error) {
			return &ioav20220601.ModifyDeliveryServerResponse{
				Response: &ioav20220601.ModifyDeliveryServerResponseParams{},
			}, nil
		})

	ctx := newIoaDeliveryServerTestContext(t, patches)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, ioaTestModel(""))
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, ioaTestModel(""))}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty server id")
	}
}

// TestIoaDeliveryServerReadSuccess verifies Read refreshes all leaf fields.
func TestIoaDeliveryServerReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "DescribeDeliveryServerWithContext",
		func(_ context.Context, req *ioav20220601.DescribeDeliveryServerRequest) (*ioav20220601.DescribeDeliveryServerResponse, error) {
			if req.ID == nil || *req.ID != ioaTestServerId {
				t.Errorf("DescribeDeliveryServer request id = %v, want %d", req.ID, ioaTestServerId)
			}
			return &ioav20220601.DescribeDeliveryServerResponse{
				Response: &ioav20220601.DescribeDeliveryServerResponseParams{Data: ioaTestItem()},
			}, nil
		})

	ctx := newIoaDeliveryServerTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, ioaTestModel("12345"))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, ioaTestModel("12345"))}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got svcioa.IoaDeliveryServerResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertDeliveryServerModel(t, got)
}

// TestIoaDeliveryServerReadNotFound verifies that a successful Describe
// response with a null Data drops the resource from state.
func TestIoaDeliveryServerReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "DescribeDeliveryServerWithContext",
		func(_ context.Context, _ *ioav20220601.DescribeDeliveryServerRequest) (*ioav20220601.DescribeDeliveryServerResponse, error) {
			return &ioav20220601.DescribeDeliveryServerResponse{
				Response: &ioav20220601.DescribeDeliveryServerResponseParams{},
			}, nil
		})

	ctx := newIoaDeliveryServerTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, ioaTestModel("12345"))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, ioaTestModel("12345"))}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the server is confirmed gone")
	}
}

// TestIoaDeliveryServerReadApiError verifies API failures surface as errors
// and keep the state intact.
func TestIoaDeliveryServerReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "DescribeDeliveryServerWithContext",
		func(_ context.Context, _ *ioav20220601.DescribeDeliveryServerRequest) (*ioav20220601.DescribeDeliveryServerResponse, error) {
			return nil, fmt.Errorf("network error")
		})

	ctx := newIoaDeliveryServerTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, ioaTestModel("12345"))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, ioaTestModel("12345"))}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on an API error")
	}
	if readResp.State.Raw.IsNull() {
		t.Error("state should be kept on an API error")
	}
}

// TestIoaDeliveryServerUpdateSuccess verifies Update sends the create-time id
// as Data.ID and refreshes state from the read-back.
func TestIoaDeliveryServerUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *ioav20220601.ModifyDeliveryServerRequest
	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "ModifyDeliveryServerWithContext",
		func(_ context.Context, req *ioav20220601.ModifyDeliveryServerRequest) (*ioav20220601.ModifyDeliveryServerResponse, error) {
			captured = req
			return &ioav20220601.ModifyDeliveryServerResponse{
				Response: &ioav20220601.ModifyDeliveryServerResponseParams{
					Data: &ioav20220601.ModifyDeliveryServerData{ID: ioaUint64(ioaTestServerId)},
				},
			}, nil
		})
	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "DescribeDeliveryServerWithContext",
		func(_ context.Context, _ *ioav20220601.DescribeDeliveryServerRequest) (*ioav20220601.DescribeDeliveryServerResponse, error) {
			return &ioav20220601.DescribeDeliveryServerResponse{
				Response: &ioav20220601.DescribeDeliveryServerResponseParams{Data: ioaTestItem()},
			}, nil
		})

	ctx := newIoaDeliveryServerTestContext(t, patches)
	plan := ioaTestModel("")
	plan.Data.Name = stringPtrValue("renamed-server")
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, ioaTestModel("12345"))
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{Plan: ctx.plan(t, plan), State: ctx.state(t, ioaTestModel("12345"))}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil || captured.Data == nil {
		t.Fatal("ModifyDeliveryServer request or Data is nil")
	}
	if captured.Data.ID == nil || *captured.Data.ID != ioaTestServerId {
		t.Errorf("update request Data.Id = %v, want %d", captured.Data.ID, ioaTestServerId)
	}
	if captured.Data.Name == nil || *captured.Data.Name != "renamed-server" {
		t.Errorf("update request Data.Name = %v, want %q", captured.Data.Name, "renamed-server")
	}

	var got svcioa.IoaDeliveryServerResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	// Non-computed values from the plan must be preserved in state after update.
	if got.Data == nil || got.Data.Name.ValueString() != "renamed-server" {
		t.Errorf("data.name = %v, want %q", got.Data.Name, "renamed-server")
	}
	if got.Data.Id.ValueInt64() != int64(ioaTestServerId) {
		t.Errorf("data.id = %d, want %d", got.Data.Id.ValueInt64(), ioaTestServerId)
	}
	if !got.Data.KafkaSaslPassword.IsNull() {
		t.Errorf("data.kafka_sasl_password = %v, want null (write-only)", got.Data.KafkaSaslPassword)
	}
	if got.Id.IsNull() || got.Id.ValueString() != "12345" {
		t.Errorf("id = %v, want 12345", got.Id.ValueString())
	}
}

// TestIoaDeliveryServerDeleteSuccess verifies Delete sends the server id and
// succeeds without error.
func TestIoaDeliveryServerDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *ioav20220601.DeleteDeliveryServerRequest
	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "DeleteDeliveryServerWithContext",
		func(_ context.Context, req *ioav20220601.DeleteDeliveryServerRequest) (*ioav20220601.DeleteDeliveryServerResponse, error) {
			captured = req
			return &ioav20220601.DeleteDeliveryServerResponse{
				Response: &ioav20220601.DeleteDeliveryServerResponseParams{},
			}, nil
		})

	ctx := newIoaDeliveryServerTestContext(t, patches)
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, ioaTestModel("12345"))}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}
	if captured == nil || captured.ID == nil || *captured.ID != ioaTestServerId {
		t.Errorf("delete request id = %v, want %d", captured, ioaTestServerId)
	}
}

// TestIoaDeliveryServerDeleteError verifies that delete failures keep the
// error diagnostics so the id stays in state for a retry.
func TestIoaDeliveryServerDeleteError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(reflect.TypeOf(&ioav20220601.Client{}), "DeleteDeliveryServerWithContext",
		func(_ context.Context, _ *ioav20220601.DeleteDeliveryServerRequest) (*ioav20220601.DeleteDeliveryServerResponse, error) {
			return nil, fmt.Errorf("delete failed")
		})

	ctx := newIoaDeliveryServerTestContext(t, patches)
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, ioaTestModel("12345"))}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}

// TestIoaDeliveryServerImportState verifies the id-passthrough import support.
func TestIoaDeliveryServerImportState(t *testing.T) {
	ctx := newIoaDeliveryServerTestContext(t, gomonkey.NewPatches())
	importer, ok := ctx.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}

	var importResp resource.ImportStateResponse
	importResp.State.Schema = ctx.schema
	importResp.State.Raw = tftypes.NewValue(ctx.schema.Type().TerraformType(context.Background()), nil)
	importer.ImportState(context.Background(),
		resource.ImportStateRequest{ID: "12345"}, &importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", importResp.Diagnostics.Errors())
	}

	var id types.String
	if diags := importResp.State.GetAttribute(context.Background(), path.Root("id"), &id); diags.HasError() {
		t.Fatalf("failed to read imported id: %v", diags.Errors())
	}
	if id.ValueString() != "12345" {
		t.Errorf("imported id = %q, want 12345", id.ValueString())
	}
}

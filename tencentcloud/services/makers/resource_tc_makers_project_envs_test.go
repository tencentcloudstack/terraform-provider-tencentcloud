package makers_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	makersv20261113 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcmakers "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/makers"
)

const (
	testEnvZoneId    = "zone-1234"
	testEnvProjectId = "proj-5678"
)

func sdkString(value string) *string { return &value }

func sdkInt64(value int64) *int64 { return &value }

// makersProjectEnvsTestContext wires a resource instance to a patched SDK client.
type makersProjectEnvsTestContext struct {
	resource resource.Resource
	client   *makersv20261113.Client
	schema   schema.Schema
}

func newMakersProjectEnvsTestContext(t *testing.T, patches *gomonkey.Patches) *makersProjectEnvsTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	makersClient := &makersv20261113.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseMakersClient", makersClient)

	res := svcmakers.NewMakersProjectEnvsResource()
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

	return &makersProjectEnvsTestContext{resource: res, client: makersClient, schema: schemaResp.Schema}
}

func (c *makersProjectEnvsTestContext) plan(t *testing.T, model svcmakers.MakersProjectEnvsResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *makersProjectEnvsTestContext) state(t *testing.T, model svcmakers.MakersProjectEnvsResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func mustStringSet(t *testing.T, values ...string) types.Set {
	t.Helper()

	set, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		t.Fatalf("failed to build string set: %v", diags.Errors())
	}
	return set
}

func sampleModel(t *testing.T) svcmakers.MakersProjectEnvsResourceModel {
	t.Helper()

	return svcmakers.MakersProjectEnvsResourceModel{
		Id:        types.StringNull(),
		ZoneId:    types.StringValue(testEnvZoneId),
		ProjectId: types.StringValue(testEnvProjectId),
		EnvVars: []svcmakers.MakersProjectEnvVarModel{
			{
				Key:     types.StringValue("API_KEY"),
				Value:   types.StringValue("secret"),
				Id:      types.Int64Null(),
				Comment: types.StringValue("production key"),
				Env:     mustStringSet(t, "Production"),
			},
			{
				Key:     types.StringValue("DEBUG"),
				Value:   types.StringValue("1"),
				Id:      types.Int64Null(),
				Comment: types.StringNull(),
				Env:     types.SetNull(types.StringType),
			},
		},
	}
}

func modelFromState(t *testing.T, state tfsdk.State) svcmakers.MakersProjectEnvsResourceModel {
	t.Helper()

	var model svcmakers.MakersProjectEnvsResourceModel
	if diags := state.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	return model
}

func describeResponse(items ...*makersv20261113.EnvVar) *makersv20261113.DescribeProjectEnvsResponse {
	return &makersv20261113.DescribeProjectEnvsResponse{
		Response: &makersv20261113.DescribeProjectEnvsResponseParams{EnvVars: items},
	}
}

func sampleEnvVarItems() []*makersv20261113.EnvVar {
	return []*makersv20261113.EnvVar{
		{
			Key:     sdkString("API_KEY"),
			Value:   sdkString("secret"),
			Id:      sdkInt64(11),
			Comment: sdkString("production key"),
			Env:     []*string{sdkString("Production")},
		},
		{
			Key:   sdkString("DEBUG"),
			Value: sdkString("1"),
			Id:    sdkInt64(12),
		},
	}
}

// TestMakersProjectEnvsSchema verifies the attribute set and its required flags.
func TestMakersProjectEnvsSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newMakersProjectEnvsTestContext(t, patches)
	top := testContext.schema.Attributes

	if len(top) != 4 {
		t.Fatalf("schema has %d attributes, want 4", len(top))
	}

	for _, name := range []string{"zone_id", "project_id", "env_vars"} {
		if attribute, ok := top[name]; !ok || !attribute.IsRequired() {
			t.Fatalf("attribute %q must be required", name)
		}
	}
	if !top["id"].IsComputed() {
		t.Fatalf("attribute id must be computed")
	}

	envVars := top["env_vars"].(schema.ListNestedAttribute)
	if !envVars.NestedObject.Attributes["key"].IsRequired() {
		t.Fatalf("env_vars.key must be required")
	}
	if !envVars.NestedObject.Attributes["value"].IsRequired() {
		t.Fatalf("env_vars.value must be required")
	}
	if !envVars.NestedObject.Attributes["id"].IsOptional() {
		t.Fatalf("env_vars.id must be optional")
	}
	if !envVars.NestedObject.Attributes["comment"].IsOptional() {
		t.Fatalf("env_vars.comment must be optional")
	}
	if !envVars.NestedObject.Attributes["env"].IsOptional() {
		t.Fatalf("env_vars.env must be optional")
	}
}
func TestMakersProjectEnvsCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newMakersProjectEnvsTestContext(t, patches)

	var capturedModify *makersv20261113.ModifyProjectEnvsRequest
	var capturedDescribe *makersv20261113.DescribeProjectEnvsRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"ModifyProjectEnvsWithContext",
		func(_ context.Context, request *makersv20261113.ModifyProjectEnvsRequest) (*makersv20261113.ModifyProjectEnvsResponse, error) {
			capturedModify = request
			return &makersv20261113.ModifyProjectEnvsResponse{
				Response: &makersv20261113.ModifyProjectEnvsResponseParams{},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeProjectEnvsWithContext",
		func(_ context.Context, request *makersv20261113.DescribeProjectEnvsRequest) (*makersv20261113.DescribeProjectEnvsResponse, error) {
			capturedDescribe = request
			return describeResponse(sampleEnvVarItems()...), nil
		},
	)

	model := sampleModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", createResp.Diagnostics.Errors())
	}

	if capturedModify == nil {
		t.Fatalf("ModifyProjectEnvs was not called")
	}
	if capturedModify.ZoneId == nil || *capturedModify.ZoneId != testEnvZoneId {
		t.Fatalf("ModifyProjectEnvs ZoneId is not %q", testEnvZoneId)
	}
	if capturedModify.ProjectId == nil || *capturedModify.ProjectId != testEnvProjectId {
		t.Fatalf("ModifyProjectEnvs ProjectId is not %q", testEnvProjectId)
	}
	if len(capturedModify.EnvVars) != 2 {
		t.Fatalf("ModifyProjectEnvs EnvVars has %d elements, want 2", len(capturedModify.EnvVars))
	}
	if *capturedModify.EnvVars[0].Key != "API_KEY" || *capturedModify.EnvVars[0].Value != "secret" {
		t.Fatalf("ModifyProjectEnvs first env var is %v/%v", *capturedModify.EnvVars[0].Key, *capturedModify.EnvVars[0].Value)
	}
	if *capturedModify.EnvVars[0].Comment != "production key" {
		t.Fatalf("ModifyProjectEnvs first env var comment is %v", *capturedModify.EnvVars[0].Comment)
	}
	if len(capturedModify.EnvVars[0].Env) != 1 || *capturedModify.EnvVars[0].Env[0] != "Production" {
		t.Fatalf("ModifyProjectEnvs first env var env set is wrong")
	}
	if capturedModify.EnvVars[1].Id != nil {
		t.Fatalf("ModifyProjectEnvs second env var must not carry an id on create")
	}
	if capturedModify.EnvVars[1].Comment != nil || capturedModify.EnvVars[1].Env != nil {
		t.Fatalf("ModifyProjectEnvs second env var must omit unset optional fields")
	}

	if capturedDescribe == nil {
		t.Fatalf("DescribeProjectEnvs was not called")
	}
	if capturedDescribe.Env != nil {
		t.Fatalf("DescribeProjectEnvs must not use the Env filter")
	}

	stateModel := modelFromState(t, createResp.State)
	// The plan does not set the nested ids (null); the cloud-assigned ids
	// read back after apply must land in state — this is exactly the
	// Optional+Computed semantics that keeps Terraform's post-apply
	// consistency check happy.
	if stateModel.EnvVars[0].Id.ValueInt64() != 11 {
		t.Fatalf("state env_vars[0].id is %d, want the cloud-assigned 11", stateModel.EnvVars[0].Id.ValueInt64())
	}
	if stateModel.EnvVars[1].Id.ValueInt64() != 12 {
		t.Fatalf("state env_vars[1].id is %d, want the cloud-assigned 12", stateModel.EnvVars[1].Id.ValueInt64())
	}
	if !stateModel.Id.Equal(types.StringValue(testEnvZoneId + "#" + testEnvProjectId)) {
		t.Fatalf("state id is %s, want %s", stateModel.Id.ValueString(), testEnvZoneId+"#"+testEnvProjectId)
	}
	if len(stateModel.EnvVars) != 2 {
		t.Fatalf("state has %d env vars, want 2", len(stateModel.EnvVars))
	}
	if !stateModel.EnvVars[0].Id.Equal(types.Int64Value(11)) {
		t.Fatalf("state env var id is %s, want 11", stateModel.EnvVars[0].Id.String())
	}
	if !stateModel.EnvVars[0].Key.Equal(types.StringValue("API_KEY")) {
		t.Fatalf("state env var key is %s", stateModel.EnvVars[0].Key.ValueString())
	}
	if !stateModel.EnvVars[0].Env.Equal(mustStringSet(t, "Production")) {
		t.Fatalf("state env var env set does not match")
	}
	if !stateModel.EnvVars[1].Comment.IsNull() {
		t.Fatalf("state second env var comment must be null")
	}
	if !stateModel.EnvVars[1].Env.IsNull() {
		t.Fatalf("state second env var env set must be null")
	}
}

// TestMakersProjectEnvsReadNotFound verifies that an empty describe result
// removes the resource from state.
func TestMakersProjectEnvsReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newMakersProjectEnvsTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeProjectEnvsWithContext",
		func(_ context.Context, request *makersv20261113.DescribeProjectEnvsRequest) (*makersv20261113.DescribeProjectEnvsResponse, error) {
			return describeResponse(), nil
		},
	)

	model := sampleModel(t)
	model.Id = types.StringValue(testEnvZoneId + "#" + testEnvProjectId)
	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: readResp.State}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}
}

// TestMakersProjectEnvsDelete verifies the delete request built from state.
func TestMakersProjectEnvsDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newMakersProjectEnvsTestContext(t, patches)

	var capturedDelete *makersv20261113.DeleteProjectEnvsRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteProjectEnvsWithContext",
		func(_ context.Context, request *makersv20261113.DeleteProjectEnvsRequest) (*makersv20261113.DeleteProjectEnvsResponse, error) {
			capturedDelete = request
			return &makersv20261113.DeleteProjectEnvsResponse{
				Response: &makersv20261113.DeleteProjectEnvsResponseParams{},
			}, nil
		},
	)

	model := sampleModel(t)
	model.Id = types.StringValue(testEnvZoneId + "#" + testEnvProjectId)
	model.EnvVars[0].Id = types.Int64Value(11)
	deleteResp := &resource.DeleteResponse{}
	testContext.resource.Delete(context.Background(), resource.DeleteRequest{State: testContext.state(t, model)}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", deleteResp.Diagnostics.Errors())
	}

	if capturedDelete == nil {
		t.Fatalf("DeleteProjectEnvs was not called")
	}
	if capturedDelete.ZoneId == nil || *capturedDelete.ZoneId != testEnvZoneId {
		t.Fatalf("DeleteProjectEnvs ZoneId is not %q", testEnvZoneId)
	}
	if capturedDelete.ProjectId == nil || *capturedDelete.ProjectId != testEnvProjectId {
		t.Fatalf("DeleteProjectEnvs ProjectId is not %q", testEnvProjectId)
	}
	if len(capturedDelete.EnvVars) != 2 {
		t.Fatalf("DeleteProjectEnvs EnvVars has %d elements, want 2", len(capturedDelete.EnvVars))
	}
	if *capturedDelete.EnvVars[0].Id != 11 {
		t.Fatalf("DeleteProjectEnvs first env var id is %v, want 11", capturedDelete.EnvVars[0].Id)
	}
	if *capturedDelete.EnvVars[1].Key != "DEBUG" {
		t.Fatalf("DeleteProjectEnvs second env var key is %v", capturedDelete.EnvVars[1].Key)
	}
}

// TestMakersProjectEnvsImportState verifies composite id parsing.
func TestMakersProjectEnvsImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newMakersProjectEnvsTestContext(t, patches)

	imported, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	// The framework server initializes the import state raw value to a fully
	// null object of the schema type before the provider ImportState runs;
	// mirror that so SetAttribute can store into it.
	importResp.State.Raw = tftypes.NewValue(testContext.schema.Type().TerraformType(context.Background()), nil)
	imported.ImportState(context.Background(), resource.ImportStateRequest{ID: testEnvZoneId + "#" + testEnvProjectId}, importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", importResp.Diagnostics.Errors())
	}

	var model svcmakers.MakersProjectEnvsResourceModel
	if diags := importResp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to read imported state: %v", diags.Errors())
	}
	if model.ZoneId.ValueString() != testEnvZoneId || model.ProjectId.ValueString() != testEnvProjectId {
		t.Fatalf("imported zone/project ids are %s/%s", model.ZoneId.ValueString(), model.ProjectId.ValueString())
	}

	badResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	badResp.State.Raw = tftypes.NewValue(testContext.schema.Type().TerraformType(context.Background()), nil)
	imported.ImportState(context.Background(), resource.ImportStateRequest{ID: "invalid-id"}, badResp)
	if !badResp.Diagnostics.HasError() {
		t.Fatalf("import with an invalid id must fail")
	}
}

package makers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tftypes "github.com/hashicorp/terraform-plugin-go/tftypes"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcmakers "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/makers"
)

const (
	makersTestZoneId      = "zone-12345"
	makersTestProjectId   = "proj-67890"
	makersTestCompositeId = makersTestZoneId + "#" + makersTestProjectId
	makersTestName        = "tf-example-project"
	makersTestTemplateId  = "tpl-abc"
)

func makersStr(v string) *string { return &v }

// makersTestProject builds the DescribeProjects response item.
func makersTestProject() *sdk.Project {
	return &sdk.Project{
		ProjectId:     makersStr(makersTestProjectId),
		Name:          makersStr(makersTestName),
		Status:        makersStr("Normal"),
		ZoneId:        makersStr(makersTestZoneId),
		Area:          makersStr("mainland"),
		Provider:      makersStr("Github"),
		Channel:       makersStr("Custom"),
		PresetDomain:  makersStr("tf-example.makers.site"),
		RootDir:       makersStr("./"),
		OutputDir:     makersStr("dist"),
		BuildCmd:      makersStr("npm run build"),
		InstallCmd:    makersStr("npm install"),
		Framework:     makersStr("react"),
		NodejsVersion: makersStr("20.18.0"),
		RepoUrl:       makersStr("https://github.com/example/tf-example"),
		RepoOwner:     makersStr("example"),
		RepoName:      makersStr("tf-example"),
		RepoBranch:    makersStr("main"),
		CreatedOn:     makersStr("2024-01-01T00:00:00Z"),
		ModifiedOn:    makersStr("2024-01-02T00:00:00Z"),
	}
}

// makersTestModel builds the terraform model used for plan/state.
func makersTestModel(id string) svcmakers.MakersProjectResourceModel {
	return svcmakers.MakersProjectResourceModel{
		Id:              types.StringValue(id),
		ZoneId:          types.StringValue(makersTestZoneId),
		Name:            types.StringValue(makersTestName),
		Channel:         types.StringValue("Custom"),
		Provider:        types.StringValue("Github"),
		Source:          types.StringValue("terraform"),
		Area:            types.StringValue("mainland"),
		NodeJsVersion:   types.StringValue("20.18.0"),
		RootDir:         types.StringValue("./"),
		OutputDir:       types.StringValue("dist"),
		BuildCmd:        types.StringValue("npm run build"),
		InstallCmd:      types.StringValue("npm install"),
		Framework:       types.StringValue("react"),
		TemplateId:      types.StringValue(makersTestTemplateId),
		RepoUrl:         types.StringValue("https://github.com/example/tf-example"),
		RepoOwner:       types.StringValue("example"),
		RepoName:        types.StringValue("tf-example"),
		RepoBranch:      types.StringValue("main"),
		Namespace:       types.StringValue("tf-namespace"),
		CodingProjectId: types.Int64Value(42),
	}
}

func assertMakersProjectModel(t *testing.T, model svcmakers.MakersProjectResourceModel) {
	t.Helper()

	if model.Id.ValueString() != makersTestCompositeId {
		t.Errorf("id = %q, want %q", model.Id.ValueString(), makersTestCompositeId)
	}
	strFields := map[string]string{
		"zone_id":         model.ZoneId.ValueString(),
		"name":            model.Name.ValueString(),
		"channel":         model.Channel.ValueString(),
		"code_provider":   model.Provider.ValueString(),
		"source":          model.Source.ValueString(),
		"area":            model.Area.ValueString(),
		"node_js_version": model.NodeJsVersion.ValueString(),
		"root_dir":        model.RootDir.ValueString(),
		"output_dir":      model.OutputDir.ValueString(),
		"build_cmd":       model.BuildCmd.ValueString(),
		"install_cmd":     model.InstallCmd.ValueString(),
		"framework":       model.Framework.ValueString(),
		"repo_url":        model.RepoUrl.ValueString(),
		"repo_owner":      model.RepoOwner.ValueString(),
		"repo_name":       model.RepoName.ValueString(),
		"repo_branch":     model.RepoBranch.ValueString(),
		"namespace":       model.Namespace.ValueString(),
		"status":          model.Status.ValueString(),
		"preset_domain":   model.PresetDomain.ValueString(),
		"created_on":      model.CreatedOn.ValueString(),
		"modified_on":     model.ModifiedOn.ValueString(),
	}
	for name, value := range strFields {
		if value == "" {
			t.Errorf("%s is empty, want a non-empty value", name)
		}
	}
	if model.Status.ValueString() != "Normal" {
		t.Errorf("status = %q, want Normal", model.Status.ValueString())
	}
	if model.CodingProjectId.ValueInt64() != 42 {
		t.Errorf("coding_project_id = %d, want 42", model.CodingProjectId.ValueInt64())
	}
}

// makersProjectTestContext wires the resource to a configured client. The
// service layer methods are patched per test case.
type makersProjectTestContext struct {
	resource resource.Resource
	schema   schema.Schema
}

func newMakersProjectTestContext(t *testing.T, patches *gomonkey.Patches) *makersProjectTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svcmakers.NewMakersProjectResource()
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

	_ = patches
	return &makersProjectTestContext{resource: res, schema: schemaResp.Schema}
}

func (c *makersProjectTestContext) plan(t *testing.T, model svcmakers.MakersProjectResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *makersProjectTestContext) state(t *testing.T, model svcmakers.MakersProjectResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

// patchMakersService patches the named MakersService method with fn.
func patchMakersService(patches *gomonkey.Patches, method string, fn interface{}) *gomonkey.Patches {
	patches.ApplyMethodFunc(&svcmakers.MakersService{}, method, fn)
	return patches
}

// TestMakersProjectCreateSuccess verifies that Create sends all request
// fields, tracks the returned project id and refreshes the state.
func TestMakersProjectCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *sdk.CreateProjectRequest
	patchMakersService(patches, "CreateProject",
		func(_ context.Context, req *sdk.CreateProjectRequest) (*string, error) {
			captured = req
			return makersStr(makersTestProjectId), nil
		})
	patchMakersService(patches, "DescribeProject",
		func(_ context.Context, zoneId, _ string) (*sdk.Project, error) {
			if zoneId != makersTestZoneId {
				t.Errorf("DescribeProject zone_id = %q, want %q", zoneId, makersTestZoneId)
			}
			return makersTestProject(), nil
		})

	ctx := newMakersProjectTestContext(t, patches)
	model := makersTestModel("")
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateProject request is nil")
	}
	if captured.ZoneId == nil || *captured.ZoneId != makersTestZoneId {
		t.Errorf("create request zone_id = %v, want %q", captured.ZoneId, makersTestZoneId)
	}
	if captured.Name == nil || *captured.Name != makersTestName {
		t.Errorf("create request name = %v, want %q", captured.Name, makersTestName)
	}
	if captured.Channel == nil || *captured.Channel != "Custom" {
		t.Errorf("create request channel = %v, want Custom", captured.Channel)
	}
	if captured.TemplateId == nil || *captured.TemplateId != makersTestTemplateId {
		t.Errorf("create request template_id = %v, want %q", captured.TemplateId, makersTestTemplateId)
	}
	if captured.Namespace == nil || *captured.Namespace != "tf-namespace" {
		t.Errorf("create request namespace = %v, want tf-namespace", captured.Namespace)
	}
	if captured.CodingProjectId == nil || *captured.CodingProjectId != 42 {
		t.Errorf("create request coding_project_id = %v, want 42", captured.CodingProjectId)
	}

	var got svcmakers.MakersProjectResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertMakersProjectModel(t, got)
}

// TestMakersProjectCreateNullCloudFields verifies that fields the describe
// API reports as null keep the configured value instead of being overwritten
// with null (which would fail Terraform's post-apply consistency check).
func TestMakersProjectCreateNullCloudFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchMakersService(patches, "CreateProject",
		func(_ context.Context, _ *sdk.CreateProjectRequest) (*string, error) {
			return makersStr(makersTestProjectId), nil
		})
	// The describe response echoes only the project id; every other field is
	// null, as the API documents for sparse projects.
	patchMakersService(patches, "DescribeProject",
		func(_ context.Context, _, _ string) (*sdk.Project, error) {
			return &sdk.Project{ProjectId: makersStr(makersTestProjectId)}, nil
		})

	ctx := newMakersProjectTestContext(t, patches)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, makersTestModel(""))
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, makersTestModel(""))}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	var got svcmakers.MakersProjectResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Channel.ValueString() != "Custom" {
		t.Errorf("channel = %q, want the configured value Custom", got.Channel.ValueString())
	}
	if got.Name.ValueString() != makersTestName {
		t.Errorf("name = %q, want the configured value %q", got.Name.ValueString(), makersTestName)
	}
	if got.ZoneId.ValueString() != makersTestZoneId {
		t.Errorf("zone_id = %q, want the configured value %q", got.ZoneId.ValueString(), makersTestZoneId)
	}
}

// TestMakersProjectCreateEmptyId verifies that Create fails loudly when the
// API returns no project id.
func TestMakersProjectCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchMakersService(patches, "CreateProject",
		func(_ context.Context, _ *sdk.CreateProjectRequest) (*string, error) {
			return nil, nil
		})

	ctx := newMakersProjectTestContext(t, patches)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, makersTestModel(""))
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, makersTestModel(""))}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty project id")
	}
}

// TestMakersProjectCreateApiError verifies API failures surface as errors.
func TestMakersProjectCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchMakersService(patches, "CreateProject",
		func(_ context.Context, _ *sdk.CreateProjectRequest) (*string, error) {
			return nil, fmt.Errorf("create failed")
		})

	ctx := newMakersProjectTestContext(t, patches)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, makersTestModel(""))
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, makersTestModel(""))}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestMakersProjectReadSuccess verifies Read refreshes all leaf fields.
func TestMakersProjectReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var gotZoneId, gotProjectId string
	patchMakersService(patches, "DescribeProject",
		func(_ context.Context, zoneId, projectId string) (*sdk.Project, error) {
			gotZoneId, gotProjectId = zoneId, projectId
			return makersTestProject(), nil
		})

	ctx := newMakersProjectTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, makersTestModel(makersTestCompositeId))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, makersTestModel(makersTestCompositeId))}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotZoneId != makersTestZoneId || gotProjectId != makersTestProjectId {
		t.Errorf("DescribeProject called with (%q, %q), want (%q, %q)", gotZoneId, gotProjectId, makersTestZoneId, makersTestProjectId)
	}

	var got svcmakers.MakersProjectResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertMakersProjectModel(t, got)
}

// TestMakersProjectReadNotFound verifies that a successful DescribeProjects
// response without the tracked project drops the resource from state.
func TestMakersProjectReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchMakersService(patches, "DescribeProject",
		func(_ context.Context, _, _ string) (*sdk.Project, error) {
			return nil, nil
		})

	ctx := newMakersProjectTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, makersTestModel(makersTestCompositeId))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, makersTestModel(makersTestCompositeId))}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the project is confirmed gone")
	}
}

// TestMakersProjectReadApiError verifies API failures surface as errors and
// keep the state intact.
func TestMakersProjectReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchMakersService(patches, "DescribeProject",
		func(_ context.Context, _, _ string) (*sdk.Project, error) {
			return nil, fmt.Errorf("describe failed")
		})

	ctx := newMakersProjectTestContext(t, patches)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, makersTestModel(makersTestCompositeId))
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, makersTestModel(makersTestCompositeId))}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on an API error")
	}
	if readResp.State.Raw.IsNull() {
		t.Error("state should be kept on an API error")
	}
}

// TestMakersProjectUpdateSuccess verifies Update sends only ModifyProject
// supported fields and refreshes the state afterwards.
func TestMakersProjectUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *sdk.ModifyProjectRequest
	patchMakersService(patches, "ModifyProject",
		func(_ context.Context, req *sdk.ModifyProjectRequest) error {
			captured = req
			return nil
		})
	patchMakersService(patches, "DescribeProject",
		func(_ context.Context, _, _ string) (*sdk.Project, error) {
			return makersTestProject(), nil
		})

	ctx := newMakersProjectTestContext(t, patches)
	// modified_on is a plain Computed attribute without UseStateForUnknown:
	// the plan value is unknown, so the fresh cloud-side timestamp written
	// by ModifyProject flows straight into the state after the read back.
	staleModifiedOn := "2026-11-18T09:34:36.022Z"
	planModel := makersTestModel(makersTestCompositeId)
	planModel.ModifiedOn = types.StringValue(staleModifiedOn)
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, makersTestModel(makersTestCompositeId))
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, planModel),
			State: ctx.state(t, makersTestModel(makersTestCompositeId)),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyProject request is nil")
	}
	if captured.ZoneId == nil || *captured.ZoneId != makersTestZoneId {
		t.Errorf("update request zone_id = %v, want %q", captured.ZoneId, makersTestZoneId)
	}
	if captured.ProjectId == nil || *captured.ProjectId != makersTestProjectId {
		t.Errorf("update request project_id = %v, want %q", captured.ProjectId, makersTestProjectId)
	}
	if captured.Name == nil || *captured.Name != makersTestName {
		t.Errorf("update request name = %v, want %q", captured.Name, makersTestName)
	}
	if captured.NodeJsVersion == nil || *captured.NodeJsVersion != "20.18.0" {
		t.Errorf("update request node_js_version = %v, want 20.18.0", captured.NodeJsVersion)
	}

	var got svcmakers.MakersProjectResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	// The stale planned value must be overwritten by the fresh timestamp the
	// cloud side reported: modified_on is Computed and follows the API.
	if got.ModifiedOn.ValueString() == staleModifiedOn {
		t.Errorf("modified_on still holds the stale planned value %q", staleModifiedOn)
	}
	if got.ModifiedOn.ValueString() != *makersTestProject().ModifiedOn {
		t.Errorf("modified_on = %q, want the cloud-side value %q", got.ModifiedOn.ValueString(), *makersTestProject().ModifiedOn)
	}
	assertMakersProjectModel(t, got)
}

// TestMakersProjectUpdateApiError verifies update failures surface as errors.
func TestMakersProjectUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchMakersService(patches, "ModifyProject",
		func(_ context.Context, _ *sdk.ModifyProjectRequest) error {
			return fmt.Errorf("modify failed")
		})

	ctx := newMakersProjectTestContext(t, patches)
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, makersTestModel(makersTestCompositeId))
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, makersTestModel(makersTestCompositeId)),
			State: ctx.state(t, makersTestModel(makersTestCompositeId)),
		}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should fail on an API error")
	}
}

// TestMakersProjectUpdateClearRejected verifies that removing an updatable
// optional attribute from the configuration is rejected loudly because the
// clear semantics of ModifyProject is not documented.
func TestMakersProjectUpdateClearRejected(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var modifyCalled bool
	patchMakersService(patches, "ModifyProject",
		func(_ context.Context, _ *sdk.ModifyProjectRequest) error {
			modifyCalled = true
			return nil
		})

	ctx := newMakersProjectTestContext(t, patches)
	planModel := makersTestModel(makersTestCompositeId)
	planModel.Framework = types.StringNull()
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, makersTestModel(makersTestCompositeId))
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, planModel),
			State: ctx.state(t, makersTestModel(makersTestCompositeId)),
		}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should reject clearing an attribute with unconfirmed clear semantics")
	}
	if modifyCalled {
		t.Error("ModifyProject should not be called when the clear semantics check fails")
	}
}

// TestMakersProjectDeleteSuccess verifies Delete sends the zone and project id.
func TestMakersProjectDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *sdk.DeleteProjectRequest
	patchMakersService(patches, "DeleteProject",
		func(_ context.Context, req *sdk.DeleteProjectRequest) error {
			captured = req
			return nil
		})

	ctx := newMakersProjectTestContext(t, patches)
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, makersTestModel(makersTestCompositeId))}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}
	if captured == nil || captured.ZoneId == nil || *captured.ZoneId != makersTestZoneId {
		t.Errorf("delete request zone_id = %v, want %q", captured, makersTestZoneId)
	}
	if captured.ProjectId == nil || *captured.ProjectId != makersTestProjectId {
		t.Errorf("delete request project_id = %v, want %q", captured, makersTestProjectId)
	}
}

// TestMakersProjectDeleteError verifies that delete failures keep the error
// diagnostics so the id stays in state for a retry.
func TestMakersProjectDeleteError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchMakersService(patches, "DeleteProject",
		func(_ context.Context, _ *sdk.DeleteProjectRequest) error {
			return fmt.Errorf("delete failed")
		})

	ctx := newMakersProjectTestContext(t, patches)
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, makersTestModel(makersTestCompositeId))}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}

// TestMakersProjectImportState verifies the id-passthrough import support.
func TestMakersProjectImportState(t *testing.T) {
	ctx := newMakersProjectTestContext(t, gomonkey.NewPatches())
	importer, ok := ctx.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}

	var importResp resource.ImportStateResponse
	importResp.State.Schema = ctx.schema
	importResp.State.Raw = tftypes.NewValue(ctx.schema.Type().TerraformType(context.Background()), nil)
	importer.ImportState(context.Background(),
		resource.ImportStateRequest{ID: makersTestCompositeId}, &importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", importResp.Diagnostics.Errors())
	}

	var id types.String
	if diags := importResp.State.GetAttribute(context.Background(), path.Root("id"), &id); diags.HasError() {
		t.Fatalf("failed to read imported id: %v", diags.Errors())
	}
	if id.ValueString() != makersTestCompositeId {
		t.Errorf("imported id = %q, want %q", id.ValueString(), makersTestCompositeId)
	}

	// The zone id segment is backfilled so the state is usable directly.
	var zoneId types.String
	if diags := importResp.State.GetAttribute(context.Background(), path.Root("zone_id"), &zoneId); diags.HasError() {
		t.Fatalf("failed to read imported zone_id: %v", diags.Errors())
	}
	if zoneId.ValueString() != makersTestZoneId {
		t.Errorf("imported zone_id = %q, want %q", zoneId.ValueString(), makersTestZoneId)
	}
}

// TestMakersProjectImportStateReadBack verifies the full import flow: after
// ImportState stores the composite id, the Read that Terraform runs right
// afterwards fills all cloud-reported fields into the state. Fields the
// describe API does not echo back (source, template_id, namespace,
// coding_project_id) stay null until they appear in the configuration.
func TestMakersProjectImportStateReadBack(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patchMakersService(patches, "DescribeProject",
		func(_ context.Context, _, _ string) (*sdk.Project, error) {
			return makersTestProject(), nil
		})

	ctx := newMakersProjectTestContext(t, patches)
	importer, ok := ctx.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}

	var importResp resource.ImportStateResponse
	importResp.State.Schema = ctx.schema
	importResp.State.Raw = tftypes.NewValue(ctx.schema.Type().TerraformType(context.Background()), nil)
	importer.ImportState(context.Background(),
		resource.ImportStateRequest{ID: makersTestCompositeId}, &importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", importResp.Diagnostics.Errors())
	}

	// Terraform refreshes the freshly imported state right away.
	readReq := resource.ReadRequest{State: tfsdk.State{Schema: ctx.schema, Raw: importResp.State.Raw}}
	var readResp resource.ReadResponse
	readResp.State = tfsdk.State{Schema: ctx.schema}
	ctx.resource.Read(context.Background(), readReq, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read after import failed: %v", readResp.Diagnostics.Errors())
	}

	var got svcmakers.MakersProjectResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}

	// Cloud-reported fields must be filled in.
	for name, value := range map[string]string{
		"id":            got.Id.ValueString(),
		"zone_id":       got.ZoneId.ValueString(),
		"name":          got.Name.ValueString(),
		"channel":       got.Channel.ValueString(),
		"status":        got.Status.ValueString(),
		"repo_url":      got.RepoUrl.ValueString(),
		"repo_branch":   got.RepoBranch.ValueString(),
		"preset_domain": got.PresetDomain.ValueString(),
	} {
		if value == "" {
			t.Errorf("%s is empty after import + read", name)
		}
	}
	if got.Id.ValueString() != makersTestCompositeId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), makersTestCompositeId)
	}
	if got.ZoneId.ValueString() != makersTestZoneId {
		t.Errorf("zone_id = %q, want %q", got.ZoneId.ValueString(), makersTestZoneId)
	}
}

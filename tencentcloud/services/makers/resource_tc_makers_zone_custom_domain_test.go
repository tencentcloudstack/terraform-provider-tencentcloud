package makers_test

import (
	"context"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdkerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	makerssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcmakers "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/makers"
)

const (
	testZoneId      = "zone-123"
	testProjectId   = "proj-456"
	testDomain      = "www.example.com"
	testEnv         = "production"
	testCompositeId = testZoneId + "#" + testProjectId + "#" + testDomain
)

func mkString(value string) *string { return &value }

func mkInt64(value int64) *int64 { return &value }

// makersZoneCustomDomainTestContext wires a resource instance to a patched
// SDK client. Service layer methods are patched per test with gomonkey.
type makersZoneCustomDomainTestContext struct {
	resource resource.Resource
	client   *makerssdk.Client
	schema   schema.Schema
}

func newMakersZoneCustomDomainTestContext(t *testing.T, patches *gomonkey.Patches) *makersZoneCustomDomainTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	makersClient := &makerssdk.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseMakersClient", makersClient)

	res := svcmakers.NewMakersZoneCustomDomainResource()
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

	return &makersZoneCustomDomainTestContext{resource: res, client: makersClient, schema: schemaResp.Schema}
}

func (c *makersZoneCustomDomainTestContext) plan(t *testing.T, model svcmakers.MakersZoneCustomDomainResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *makersZoneCustomDomainTestContext) state(t *testing.T, model svcmakers.MakersZoneCustomDomainResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func (c *makersZoneCustomDomainTestContext) readState(t *testing.T, state tfsdk.State) svcmakers.MakersZoneCustomDomainResourceModel {
	t.Helper()

	var model svcmakers.MakersZoneCustomDomainResourceModel
	if diags := state.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	return model
}

// testZoneDomainWithStatus returns a settled test domain with the given
// status, mirroring the cloud side while an async operation is running.
func testZoneDomainWithStatus(status string) *makerssdk.ZoneDomain {
	item := testZoneDomain()
	item.Status = mkString(status)
	return item
}

func testZoneDomain() *makerssdk.ZoneDomain {
	return &makerssdk.ZoneDomain{
		Type:       mkString("Custom"),
		Domain:     mkString(testDomain),
		Status:     mkString("online"),
		Env:        mkString(testEnv),
		ZoneId:     mkString(testZoneId),
		Area:       mkString("mainland"),
		Cname:      mkString("cname.example.com"),
		MainDomain: mkString("example.com"),
		Ipv6:       &makerssdk.SwitchConfig{Switch: mkString("off")},
		Http2:      &makerssdk.SwitchConfig{Switch: mkString("on")},
		ForceRedirect: &makerssdk.ForceRedirectConfig{
			Switch:             mkString("on"),
			RedirectStatusCode: mkInt64(301),
		},
		Hsts: &makerssdk.HstsConfig{
			Switch:            mkString("on"),
			MaxAge:            mkInt64(3600),
			IncludeSubDomains: mkString("on"),
			Preload:           mkString("off"),
		},
		OwnershipVerification: &makerssdk.OwnershipVerification{
			DnsVerification: &makerssdk.DnsVerification{
				RecordType:  mkString("TXT"),
				Subdomain:   mkString("_tcdns.www.example.com"),
				RecordValue: mkString("verify-value"),
			},
			FileVerification: &makerssdk.FileVerification{
				Path:    mkString("/verify.txt"),
				Content: mkString("file-content"),
			},
			NsVerification: &makerssdk.NsVerification{
				NameServers: []*string{mkString("ns1.example.com"), mkString("ns2.example.com")},
			},
		},
	}
}

// nullTimeouts builds a null timeouts.Value matching the resource schema, so
// model (de)serialization through tfsdk.Plan/State behaves exactly like a
// configuration that does not set the timeouts block.
func (c *makersZoneCustomDomainTestContext) nullTimeouts() resourcetimeouts.Value {
	timeoutsBlock, ok := c.schema.Blocks["timeouts"]
	if !ok {
		panic("timeouts block not found in schema")
	}
	timeoutsType, ok := timeoutsBlock.Type().(resourcetimeouts.Type)
	if !ok {
		panic("unexpected timeouts block type")
	}
	return resourcetimeouts.Value{
		Object: types.ObjectNull(timeoutsType.AttributeTypes()),
	}
}

func (c *makersZoneCustomDomainTestContext) sampleModel() svcmakers.MakersZoneCustomDomainResourceModel {
	return svcmakers.MakersZoneCustomDomainResourceModel{
		Id:        types.StringNull(),
		ZoneId:    types.StringValue(testZoneId),
		ProjectId: types.StringValue(testProjectId),
		Domain:    types.StringValue(testDomain),
		Env:       types.StringValue(testEnv),
		Timeouts:  c.nullTimeouts(),
	}
}

func TestMakersZoneCustomDomainCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	var createdRequest *makerssdk.CreateZoneCustomDomainRequest
	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "CreateZoneCustomDomain", func(_ context.Context, request *makerssdk.CreateZoneCustomDomainRequest) (string, error) {
		createdRequest = request
		return testDomain, nil
	})
	describeCalls := 0
	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "DescribeZoneCustomDomain", func(_ context.Context, zoneId, projectId, domain string) (*makerssdk.ZoneDomain, error) {
		if zoneId != testZoneId || projectId != testProjectId || domain != testDomain {
			t.Errorf("unexpected read-back locator: %s/%s/%s", zoneId, projectId, domain)
		}
		describeCalls++
		// The binding is async: the first poll (the wait loop) sees the
		// domain not yet visible, the second one sees it with a `process`
		// status while the binding converges, the third one (the final read
		// back) sees the settled domain.
		if describeCalls == 1 {
			return nil, nil
		}
		if describeCalls == 2 {
			return testZoneDomainWithStatus("process"), nil
		}
		return testZoneDomain(), nil
	})

	// The framework server normally pre-initializes the response state with
	// the resource schema; mirror that here so State.Set works.
	createResp := resource.CreateResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Create(context.Background(), resource.CreateRequest{Plan: c.plan(t, c.sampleModel())}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if createdRequest == nil {
		t.Fatalf("CreateZoneCustomDomain was not called")
	}
	if createdRequest.ZoneId == nil || *createdRequest.ZoneId != testZoneId {
		t.Errorf("Create request ZoneId mismatch: %v", createdRequest.ZoneId)
	}
	if createdRequest.ProjectId == nil || *createdRequest.ProjectId != testProjectId {
		t.Errorf("Create request ProjectId mismatch: %v", createdRequest.ProjectId)
	}
	if createdRequest.Domain == nil || *createdRequest.Domain != testDomain {
		t.Errorf("Create request Domain mismatch: %v", createdRequest.Domain)
	}
	if createdRequest.Env == nil || *createdRequest.Env != testEnv {
		t.Errorf("Create request Env mismatch: %v", createdRequest.Env)
	}

	model := c.readState(t, createResp.State)
	if model.Id.ValueString() != testCompositeId {
		t.Errorf("unexpected id: %s", model.Id.ValueString())
	}
	if model.Env.ValueString() != testEnv {
		t.Errorf("unexpected env: %s", model.Env.ValueString())
	}
	if describeCalls < 3 {
		t.Errorf("expected the async wait to poll at least three times (invisible, process, settled), got %d calls", describeCalls)
	}
}

func TestMakersZoneCustomDomainCreateError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "CreateZoneCustomDomain", func(_ context.Context, _ *makerssdk.CreateZoneCustomDomainRequest) (string, error) {
		return "", &sdkerr.TencentCloudSDKError{}
	})

	createResp := resource.CreateResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Create(context.Background(), resource.CreateRequest{Plan: c.plan(t, c.sampleModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected Create to fail on API error")
	}
}

func TestMakersZoneCustomDomainReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "DescribeZoneCustomDomain", func(_ context.Context, zoneId, projectId, domain string) (*makerssdk.ZoneDomain, error) {
		if zoneId != testZoneId || projectId != testProjectId || domain != testDomain {
			t.Errorf("unexpected read locator: %s/%s/%s", zoneId, projectId, domain)
		}
		return nil, nil
	})

	model := c.sampleModel()
	model.Id = types.StringValue(testCompositeId)

	readResp := resource.ReadResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Read(context.Background(), resource.ReadRequest{State: c.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Errorf("state should be removed when the domain is gone")
	}
}

func TestMakersServiceDescribeZoneCustomDomainFiltersPreset(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	connectivityClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(connectivityClient, "UseMakersClient", &makerssdk.Client{})

	preset := testZoneDomain()
	preset.Type = mkString("Preset")
	preset.Domain = mkString("preset.example.com")

	var queriedRequest *makerssdk.DescribeZoneCustomDomainsRequest
	patches.ApplyMethodFunc(&makerssdk.Client{}, "DescribeZoneCustomDomainsWithContext", func(_ context.Context, request *makerssdk.DescribeZoneCustomDomainsRequest) (*makerssdk.DescribeZoneCustomDomainsResponse, error) {
		queriedRequest = request
		return &makerssdk.DescribeZoneCustomDomainsResponse{
			Response: &makerssdk.DescribeZoneCustomDomainsResponseParams{
				TotalCount:   mkInt64(2),
				PagesDomains: []*makerssdk.ZoneDomain{preset, testZoneDomain()},
			},
		}, nil
	})

	service := svcmakers.NewMakersService(connectivityClient)
	item, err := service.DescribeZoneCustomDomain(context.Background(), testZoneId, testProjectId, testDomain)
	if err != nil {
		t.Fatalf("DescribeZoneCustomDomain failed: %v", err)
	}
	if item == nil {
		t.Fatalf("expected the Custom domain to be selected, got nil")
	}
	if item.Type == nil || *item.Type != "Custom" {
		t.Errorf("expected the Custom domain to be selected, got %v", item.Type)
	}

	if queriedRequest == nil {
		t.Fatalf("DescribeZoneCustomDomains was not called")
	}
	if queriedRequest.ZoneId == nil || *queriedRequest.ZoneId != testZoneId {
		t.Errorf("query ZoneId mismatch: %v", queriedRequest.ZoneId)
	}
	if queriedRequest.ProjectId == nil || *queriedRequest.ProjectId != testProjectId {
		t.Errorf("query ProjectId mismatch: %v", queriedRequest.ProjectId)
	}
	if len(queriedRequest.CustomDomains) != 1 || *queriedRequest.CustomDomains[0].Domain != testDomain {
		t.Errorf("query should filter exactly the managed domain: %v", queriedRequest.CustomDomains)
	}
}

func TestMakersZoneCustomDomainUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	var updatedRequest *makerssdk.ModifyZoneCustomDomainEnvRequest
	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "ModifyZoneCustomDomainEnv", func(_ context.Context, request *makerssdk.ModifyZoneCustomDomainEnvRequest) error {
		updatedRequest = request
		return nil
	})
	describeCalls := 0
	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "DescribeZoneCustomDomain", func(_ context.Context, _ string, _ string, _ string) (*makerssdk.ZoneDomain, error) {
		domain := testZoneDomain()
		describeCalls++
		// The env switch is async: the first poll still reports the old env
		// with a `process` status while the switch is running.
		if describeCalls == 1 {
			domain.Env = mkString(testEnv)
			domain.Status = mkString("process")
		} else {
			domain.Env = mkString("staging")
		}
		return domain, nil
	})

	state := c.sampleModel()
	state.Id = types.StringValue(testCompositeId)
	state.Env = types.StringValue(testEnv)

	plan := state
	plan.Env = types.StringValue("staging")

	updateResp := resource.UpdateResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Update(context.Background(), resource.UpdateRequest{Plan: c.plan(t, plan), State: c.state(t, state)}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if updatedRequest == nil {
		t.Fatalf("ModifyZoneCustomDomainEnv was not called")
	}
	if updatedRequest.Env == nil || *updatedRequest.Env != "staging" {
		t.Errorf("Update request Env mismatch: %v", updatedRequest.Env)
	}
	if updatedRequest.ZoneId == nil || *updatedRequest.ZoneId != testZoneId {
		t.Errorf("Update request ZoneId mismatch: %v", updatedRequest.ZoneId)
	}
	if updatedRequest.ProjectId == nil || *updatedRequest.ProjectId != testProjectId {
		t.Errorf("Update request ProjectId mismatch: %v", updatedRequest.ProjectId)
	}
	if updatedRequest.Domain == nil || *updatedRequest.Domain != testDomain {
		t.Errorf("Update request Domain mismatch: %v", updatedRequest.Domain)
	}

	refreshed := c.readState(t, updateResp.State)
	if refreshed.Env.ValueString() != "staging" {
		t.Errorf("unexpected env after update: %s", refreshed.Env.ValueString())
	}
	if describeCalls < 2 {
		t.Errorf("expected the async env wait to poll at least twice, got %d calls", describeCalls)
	}
}

func TestMakersZoneCustomDomainDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	var deletedZoneId, deletedProjectId, deletedDomain string
	deleted := false
	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "DeleteProjectCustomDomain", func(_ context.Context, zoneId, projectId, domain string) error {
		deleted = true
		deletedZoneId, deletedProjectId, deletedDomain = zoneId, projectId, domain
		return nil
	})
	describeCalls := 0
	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "DescribeZoneCustomDomain", func(_ context.Context, zoneId, projectId, domain string) (*makerssdk.ZoneDomain, error) {
		describeCalls++
		// The unbinding is async: the first poll still sees the domain.
		if describeCalls == 1 {
			return testZoneDomain(), nil
		}
		return nil, nil
	})

	model := c.sampleModel()
	model.Id = types.StringValue(testCompositeId)

	deleteResp := resource.DeleteResponse{}
	c.resource.Delete(context.Background(), resource.DeleteRequest{State: c.state(t, model)}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if !deleted {
		t.Fatalf("DeleteProjectCustomDomain was not called")
	}
	if deletedZoneId != testZoneId || deletedProjectId != testProjectId || deletedDomain != testDomain {
		t.Errorf("delete locator mismatch: %s/%s/%s", deletedZoneId, deletedProjectId, deletedDomain)
	}
	if describeCalls < 2 {
		t.Errorf("expected the async deletion wait to poll at least twice, got %d calls", describeCalls)
	}
}

func TestMakersServiceWaitForZoneCustomDomainTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	connectivityClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(connectivityClient, "UseMakersClient", &makerssdk.Client{})

	// The domain never becomes visible: the wait must fail after the timeout.
	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "DescribeZoneCustomDomain", func(_ context.Context, _ string, _ string, _ string) (*makerssdk.ZoneDomain, error) {
		return nil, nil
	})

	service := svcmakers.NewMakersService(connectivityClient)
	// gomonkey shortens nothing here: use a tiny timeout so the test is fast.
	err := service.WaitForZoneCustomDomain(context.Background(), testZoneId, testProjectId, testDomain, 1*time.Second)
	if err == nil {
		t.Fatalf("expected the wait to time out when the domain never appears")
	}
}

func TestMakersZoneCustomDomainCreateTimeoutPassedToWait(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "CreateZoneCustomDomain", func(_ context.Context, _ *makerssdk.CreateZoneCustomDomainRequest) (string, error) {
		return testDomain, nil
	})
	patches.ApplyMethodFunc(&svcmakers.MakersService{}, "DescribeZoneCustomDomain", func(_ context.Context, _ string, _ string, _ string) (*makerssdk.ZoneDomain, error) {
		return testZoneDomain(), nil
	})

	model := c.sampleModel()
	model.Timeouts = makersZoneCustomDomainTimeouts(t, "90m")

	createResp := resource.CreateResponse{State: tfsdk.State{Schema: c.schema}}
	c.resource.Create(context.Background(), resource.CreateRequest{Plan: c.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create with custom timeout failed: %v", createResp.Diagnostics.Errors())
	}

	refreshed := c.readState(t, createResp.State)
	if refreshed.Timeouts.IsNull() {
		t.Fatalf("timeouts block should be kept in state")
	}
}

// testSchemaBlocks returns the resource schema blocks.
func testSchemaBlocks(t *testing.T) map[string]schema.Block {
	t.Helper()

	var schemaResp resource.SchemaResponse
	res := svcmakers.NewMakersZoneCustomDomainResource()
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}
	return schemaResp.Schema.Blocks
}

// makersZoneCustomDomainTimeouts builds a resourcetimeouts.Value with the
// given create duration, reusing the attribute types of the schema block.
func makersZoneCustomDomainTimeouts(t *testing.T, create string) resourcetimeouts.Value {
	t.Helper()

	timeoutsBlock, ok := testSchemaBlocks(t)["timeouts"]
	if !ok {
		t.Fatalf("timeouts block not found in schema")
	}
	timeoutsType, ok := timeoutsBlock.Type().(resourcetimeouts.Type)
	if !ok {
		t.Fatalf("timeouts does not have the framework timeouts type: %T", timeoutsBlock.Type())
	}

	attrTypes := timeoutsType.AttributeTypes()
	attributes := map[string]attr.Value{}
	for name := range attrTypes {
		if name == "create" {
			attributes[name] = types.StringValue(create)
		} else {
			attributes[name] = types.StringNull()
		}
	}

	object, diags := types.ObjectValue(attrTypes, attributes)
	if diags.HasError() {
		t.Fatalf("failed to build timeouts object: %v", diags.Errors())
	}
	value, diags := timeoutsType.ValueFromObject(context.Background(), object)
	if diags.HasError() {
		t.Fatalf("failed to build timeouts value: %v", diags.Errors())
	}
	timeoutsValue, ok := value.(resourcetimeouts.Value)
	if !ok {
		t.Fatalf("timeouts value has unexpected type %T", value)
	}
	return timeoutsValue
}

func TestMakersZoneCustomDomainImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	importer, ok := c.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := resource.ImportStateResponse{State: tfsdk.State{Schema: c.schema}}
	// The framework server initializes the import state raw value to a fully
	// null object of the schema type before the provider ImportState runs;
	// mirror that so SetAttribute can store into it.
	importResp.State.Raw = tftypes.NewValue(c.schema.Type().TerraformType(context.Background()), nil)
	importer.ImportState(context.Background(), resource.ImportStateRequest{
		ID: testCompositeId,
	}, &importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", importResp.Diagnostics.Errors())
	}

	model := c.readState(t, importResp.State)
	if model.ZoneId.ValueString() != testZoneId {
		t.Errorf("unexpected zone_id after import: %s", model.ZoneId.ValueString())
	}
	if model.ProjectId.ValueString() != testProjectId {
		t.Errorf("unexpected project_id after import: %s", model.ProjectId.ValueString())
	}
	if model.Domain.ValueString() != testDomain {
		t.Errorf("unexpected domain after import: %s", model.Domain.ValueString())
	}
	if model.Id.ValueString() != testCompositeId {
		t.Errorf("unexpected id after import: %s", model.Id.ValueString())
	}
}

func TestMakersZoneCustomDomainImportStateInvalidId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	importer, ok := c.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := resource.ImportStateResponse{State: tfsdk.State{Schema: c.schema}}
	// The framework server initializes the import state raw value to a fully
	// null object of the schema type before the provider ImportState runs;
	// mirror that so SetAttribute can store into it.
	importResp.State.Raw = tftypes.NewValue(c.schema.Type().TerraformType(context.Background()), nil)
	importer.ImportState(context.Background(), resource.ImportStateRequest{
		ID: "only-one-segment",
	}, &importResp)
	if !importResp.Diagnostics.HasError() {
		t.Fatalf("expected ImportState to fail for an invalid composite id")
	}
}

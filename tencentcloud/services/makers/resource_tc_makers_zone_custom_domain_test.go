package makers_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	makersv20261113 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcmakers "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/makers"
)

const (
	testZoneId    = "zone-123"
	testProjectId = "proj-456"
	testDomain    = "www.example.com"
	testEnv       = "production"
)

func mkString(value string) *string { return &value }

func mkInt64(value int64) *int64 { return &value }

// makersZoneCustomDomainTestContext wires a resource instance to a patched SDK client.
type makersZoneCustomDomainTestContext struct {
	resource resource.Resource
	client   *makersv20261113.Client
	schema   schema.Schema
}

func newMakersZoneCustomDomainTestContext(t *testing.T, patches *gomonkey.Patches) *makersZoneCustomDomainTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	makersClient := &makersv20261113.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseMakersV20261113Client", makersClient)

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

func testZoneDomain() *makersv20261113.ZoneDomain {
	return &makersv20261113.ZoneDomain{
		Type:       mkString("Custom"),
		Domain:     mkString(testDomain),
		Status:     mkString("Pass"),
		Env:        mkString(testEnv),
		ZoneId:     mkString(testZoneId),
		Area:       mkString("mainland"),
		Cname:      mkString("cname.example.com"),
		MainDomain: mkString("example.com"),
		Ipv6:       &makersv20261113.SwitchConfig{Switch: mkString("off")},
		Http2:      &makersv20261113.SwitchConfig{Switch: mkString("on")},
		ForceRedirect: &makersv20261113.ForceRedirectConfig{
			Switch:             mkString("on"),
			RedirectStatusCode: mkInt64(301),
		},
		Hsts: &makersv20261113.HstsConfig{
			Switch:            mkString("on"),
			MaxAge:            mkInt64(3600),
			IncludeSubDomains: mkString("on"),
			Preload:           mkString("off"),
		},
		OwnershipVerification: &makersv20261113.OwnershipVerification{
			DnsVerification: &makersv20261113.DnsVerification{
				RecordType:  mkString("TXT"),
				Subdomain:   mkString("_tcdns.www.example.com"),
				RecordValue: mkString("verify-value"),
			},
			FileVerification: &makersv20261113.FileVerification{
				Path:    mkString("/verify.txt"),
				Content: mkString("file-content"),
			},
			NsVerification: &makersv20261113.NsVerification{
				NameServers: []*string{mkString("ns1.example.com"), mkString("ns2.example.com")},
			},
		},
	}
}

func testSchemaAttributes(t *testing.T) map[string]schema.Attribute {
	t.Helper()

	var schemaResp resource.SchemaResponse
	res := svcmakers.NewMakersZoneCustomDomainResource()
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}
	return schemaResp.Schema.Attributes
}

func mustObject(t *testing.T, attributeTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	t.Helper()

	object, diags := types.ObjectValue(attributeTypes, values)
	if diags.HasError() {
		t.Fatalf("failed to build object value: %v", diags.Errors())
	}
	return object
}

func (c *makersZoneCustomDomainTestContext) sampleModel() svcmakers.MakersZoneCustomDomainResourceModel {
	return svcmakers.MakersZoneCustomDomainResourceModel{
		Id:           types.StringNull(),
		ZoneId:       types.StringValue(testZoneId),
		ProjectId:    types.StringValue(testProjectId),
		Domain:       types.StringValue(testDomain),
		Env:          types.StringValue(testEnv),
		PagesDomains: nil,
	}
}

func TestMakersZoneCustomDomainCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	var createdRequest *makersv20261113.CreateZoneCustomDomainRequest
	patches.ApplyMethodFunc(c.client, "CreateZoneCustomDomainWithContext", func(_ context.Context, request *makersv20261113.CreateZoneCustomDomainRequest) (*makersv20261113.CreateZoneCustomDomainResponse, error) {
		createdRequest = request
		return &makersv20261113.CreateZoneCustomDomainResponse{
			Response: &makersv20261113.CreateZoneCustomDomainResponseParams{
				Domain: mkString(testDomain),
				Id:     mkString("cd-001"),
			},
		}, nil
	})
	patches.ApplyMethodFunc(c.client, "DescribeZoneCustomDomainsWithContext", func(_ context.Context, request *makersv20261113.DescribeZoneCustomDomainsRequest) (*makersv20261113.DescribeZoneCustomDomainsResponse, error) {
		return &makersv20261113.DescribeZoneCustomDomainsResponse{
			Response: &makersv20261113.DescribeZoneCustomDomainsResponseParams{
				TotalCount:   mkInt64(1),
				PagesDomains: []*makersv20261113.ZoneDomain{testZoneDomain()},
			},
		}, nil
	})

	var createResp resource.CreateResponse
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
	if model.Id.ValueString() != testZoneId+"#"+testProjectId+"#"+testDomain {
		t.Errorf("unexpected id: %s", model.Id.ValueString())
	}
	if model.Env.ValueString() != testEnv {
		t.Errorf("unexpected env: %s", model.Env.ValueString())
	}
	if model.PagesDomains == nil || len(model.PagesDomains) != 1 {
		t.Fatalf("expected one pages domain, got %v", model.PagesDomains)
	}
	pd := model.PagesDomains[0]
	if pd.Type.ValueString() != "Custom" {
		t.Errorf("unexpected pages domain type: %s", pd.Type.ValueString())
	}
	if pd.Status.ValueString() != "Pass" {
		t.Errorf("unexpected pages domain status: %s", pd.Status.ValueString())
	}
	if pd.Cname.ValueString() != "cname.example.com" {
		t.Errorf("unexpected cname: %s", pd.Cname.ValueString())
	}
	if pd.Ipv6 == nil || pd.Ipv6.Switch.ValueString() != "off" {
		t.Errorf("unexpected ipv6 switch: %v", pd.Ipv6)
	}
	if pd.ForceRedirect == nil || pd.ForceRedirect.RedirectStatusCode.ValueInt64() != 301 {
		t.Errorf("unexpected force redirect: %v", pd.ForceRedirect)
	}
	if pd.Hsts == nil || pd.Hsts.MaxAge.ValueInt64() != 3600 {
		t.Errorf("unexpected hsts: %v", pd.Hsts)
	}
	if pd.OwnershipVerification == nil || pd.OwnershipVerification.NsVerification == nil {
		t.Fatalf("missing ownership verification: %v", pd.OwnershipVerification)
	}
	if pd.OwnershipVerification.NsVerification.NameServers.IsNull() {
		t.Errorf("name servers should not be null")
	}
}

func TestMakersZoneCustomDomainCreateError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	patches.ApplyMethodFunc(c.client, "CreateZoneCustomDomainWithContext", func(_ context.Context, _ *makersv20261113.CreateZoneCustomDomainRequest) (*makersv20261113.CreateZoneCustomDomainResponse, error) {
		return nil, &sdkerr.TencentCloudSDKError{}
	})

	var createResp resource.CreateResponse
	c.resource.Create(context.Background(), resource.CreateRequest{Plan: c.plan(t, c.sampleModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected Create to fail on API error")
	}
}

func TestMakersZoneCustomDomainReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	patches.ApplyMethodFunc(c.client, "DescribeZoneCustomDomainsWithContext", func(_ context.Context, request *makersv20261113.DescribeZoneCustomDomainsRequest) (*makersv20261113.DescribeZoneCustomDomainsResponse, error) {
		if request.CustomDomains == nil || len(request.CustomDomains) != 1 || *request.CustomDomains[0].Domain != testDomain {
			t.Errorf("Read should query exactly the managed domain")
		}
		return &makersv20261113.DescribeZoneCustomDomainsResponse{
			Response: &makersv20261113.DescribeZoneCustomDomainsResponseParams{
				TotalCount:   mkInt64(0),
				PagesDomains: nil,
			},
		}, nil
	})

	model := c.sampleModel()
	model.Id = types.StringValue(testZoneId + "#" + testProjectId + "#" + testDomain)

	var readResp resource.ReadResponse
	c.resource.Read(context.Background(), resource.ReadRequest{State: c.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Errorf("state should be removed when the domain is gone")
	}
}

func TestMakersZoneCustomDomainReadFiltersPreset(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	preset := testZoneDomain()
	preset.Type = mkString("Preset")
	preset.Domain = mkString("preset.example.com")

	patches.ApplyMethodFunc(c.client, "DescribeZoneCustomDomainsWithContext", func(_ context.Context, _ *makersv20261113.DescribeZoneCustomDomainsRequest) (*makersv20261113.DescribeZoneCustomDomainsResponse, error) {
		return &makersv20261113.DescribeZoneCustomDomainsResponse{
			Response: &makersv20261113.DescribeZoneCustomDomainsResponseParams{
				TotalCount:   mkInt64(2),
				PagesDomains: []*makersv20261113.ZoneDomain{preset, testZoneDomain()},
			},
		}, nil
	})

	model := c.sampleModel()
	model.Id = types.StringValue(testZoneId + "#" + testProjectId + "#" + testDomain)

	var readResp resource.ReadResponse
	c.resource.Read(context.Background(), resource.ReadRequest{State: c.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	refreshed := c.readState(t, readResp.State)
	if len(refreshed.PagesDomains) != 1 {
		t.Fatalf("expected one pages domain after filtering preset domains, got %d", len(refreshed.PagesDomains))
	}
	if refreshed.PagesDomains[0].Type.ValueString() != "Custom" {
		t.Errorf("expected the Custom domain to be selected, got %s", refreshed.PagesDomains[0].Type.ValueString())
	}
}

func TestMakersZoneCustomDomainUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	var updatedRequest *makersv20261113.ModifyZoneCustomDomainEnvRequest
	patches.ApplyMethodFunc(c.client, "ModifyZoneCustomDomainEnvWithContext", func(_ context.Context, request *makersv20261113.ModifyZoneCustomDomainEnvRequest) (*makersv20261113.ModifyZoneCustomDomainEnvResponse, error) {
		updatedRequest = request
		return &makersv20261113.ModifyZoneCustomDomainEnvResponse{
			Response: &makersv20261113.ModifyZoneCustomDomainEnvResponseParams{
				Domain: mkString(testDomain),
				Env:    mkString("staging"),
			},
		}, nil
	})
	patches.ApplyMethodFunc(c.client, "DescribeZoneCustomDomainsWithContext", func(_ context.Context, _ *makersv20261113.DescribeZoneCustomDomainsRequest) (*makersv20261113.DescribeZoneCustomDomainsResponse, error) {
		domain := testZoneDomain()
		domain.Env = mkString("staging")
		return &makersv20261113.DescribeZoneCustomDomainsResponse{
			Response: &makersv20261113.DescribeZoneCustomDomainsResponseParams{
				TotalCount:   mkInt64(1),
				PagesDomains: []*makersv20261113.ZoneDomain{domain},
			},
		}, nil
	})

	state := c.sampleModel()
	state.Id = types.StringValue(testZoneId + "#" + testProjectId + "#" + testDomain)
	state.Env = types.StringValue(testEnv)
	state.PagesDomains = nil

	plan := state
	plan.Env = types.StringValue("staging")

	var updateResp resource.UpdateResponse
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

	refreshed := c.readState(t, updateResp.State)
	if refreshed.Env.ValueString() != "staging" {
		t.Errorf("unexpected env after update: %s", refreshed.Env.ValueString())
	}
}

func TestMakersZoneCustomDomainDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	var deletedRequest *makersv20261113.DeleteProjectCustomDomainsRequest
	patches.ApplyMethodFunc(c.client, "DeleteProjectCustomDomainsWithContext", func(_ context.Context, request *makersv20261113.DeleteProjectCustomDomainsRequest) (*makersv20261113.DeleteProjectCustomDomainsResponse, error) {
		deletedRequest = request
		return &makersv20261113.DeleteProjectCustomDomainsResponse{
			Response: &makersv20261113.DeleteProjectCustomDomainsResponseParams{},
		}, nil
	})

	model := c.sampleModel()
	model.Id = types.StringValue(testZoneId + "#" + testProjectId + "#" + testDomain)

	var deleteResp resource.DeleteResponse
	c.resource.Delete(context.Background(), resource.DeleteRequest{State: c.state(t, model)}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if deletedRequest == nil {
		t.Fatalf("DeleteProjectCustomDomains was not called")
	}
	if deletedRequest.CustomDomains == nil || len(deletedRequest.CustomDomains) != 1 || *deletedRequest.CustomDomains[0] != testDomain {
		t.Errorf("delete should send exactly the managed domain, got %v", deletedRequest.CustomDomains)
	}
	if deletedRequest.ZoneId == nil || *deletedRequest.ZoneId != testZoneId {
		t.Errorf("delete request ZoneId mismatch: %v", deletedRequest.ZoneId)
	}
}

func TestMakersZoneCustomDomainImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	c := newMakersZoneCustomDomainTestContext(t, patches)

	importer, ok := c.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	var importResp resource.ImportStateResponse
	emptyModel := svcmakers.MakersZoneCustomDomainResourceModel{}
	importer.ImportState(context.Background(), resource.ImportStateRequest{
		ID:    testZoneId + "#" + testProjectId + "#" + testDomain,
		State: c.state(t, emptyModel),
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
	if model.Id.ValueString() != testZoneId+"#"+testProjectId+"#"+testDomain {
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

	var importResp resource.ImportStateResponse
	importer.ImportState(context.Background(), resource.ImportStateRequest{
		ID:    "only-one-segment",
		State: tfsdk.State{Schema: c.schema},
	}, &importResp)
	if !importResp.Diagnostics.HasError() {
		t.Fatalf("expected ImportState to fail for an invalid composite id")
	}
}

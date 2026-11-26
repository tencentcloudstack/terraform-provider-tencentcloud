package ags_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

const (
	agsSandboxInstanceTestId     = "instance-001"
	agsSandboxInstanceTestToolId = "tool-001"
	agsSandboxInstanceTestTool   = "tf-example"
)

// agsSandboxInstanceTestContext wires a resource instance to a patched SDK
// client.
type agsSandboxInstanceTestContext struct {
	resource resource.Resource
	client   *agsv20250920.Client
	schema   schema.Schema
}

func newAgsSandboxInstanceTestContext(t *testing.T, patches *gomonkey.Patches) *agsSandboxInstanceTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseAgsClient", agsClient)

	res := svcags.NewAgsSandboxInstanceResource()
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

	return &agsSandboxInstanceTestContext{resource: res, client: agsClient, schema: schemaResp.Schema}
}

func (c *agsSandboxInstanceTestContext) plan(t *testing.T, model svcags.AgsSandboxInstanceResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *agsSandboxInstanceTestContext) state(t *testing.T, model svcags.AgsSandboxInstanceResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

// agsSandboxInstanceAttributeTypes returns the nested attribute types of the
// given top level attributes, so the test models stay in sync with the schema.
func (c *agsSandboxInstanceTestContext) agsSandboxInstanceAttributeTypes(t *testing.T) (custom map[string]attr.Type, computer map[string]attr.Type, mountOption map[string]attr.Type, metadata map[string]attr.Type) {
	t.Helper()

	top := c.schema.Attributes
	custom = objectAttributeTypes(t, top, "custom_configuration")
	computer = objectAttributeTypes(t, top, "computer_configuration")
	mountOption = listElementAttributeTypes(t, top, "mount_options")
	metadata = listElementAttributeTypes(t, top, "metadata")
	return
}

// sampleModel returns a plan model with every configurable leaf populated.
func (c *agsSandboxInstanceTestContext) sampleModel(t *testing.T) svcags.AgsSandboxInstanceResourceModel {
	t.Helper()

	customTypes, computerTypes, mountOptionTypes, metadataTypes := c.agsSandboxInstanceAttributeTypes(t)

	customAttributes := c.schema.Attributes["custom_configuration"].(schema.SingleNestedAttribute).Attributes
	envTypes := listElementAttributeTypes(t, customAttributes, "env")
	portTypes := listElementAttributeTypes(t, customAttributes, "ports")
	resourcesTypes := objectAttributeTypes(t, customAttributes, "resources")
	probeTypes := objectAttributeTypes(t, customAttributes, "probe")
	dnsConfigTypes := objectAttributeTypes(t, customAttributes, "dns_config")
	securityContextTypes := objectAttributeTypes(t, customAttributes, "security_context")

	return svcags.AgsSandboxInstanceResourceModel{
		Id:               types.StringNull(),
		InstanceId:       types.StringNull(),
		ToolId:           types.StringNull(),
		ToolName:         types.StringValue(agsSandboxInstanceTestTool),
		SourceInstanceId: types.StringNull(),
		SnapshotId:       types.StringNull(),
		Timeout:          types.StringValue("30m"),
		MountOptions: mustList(t, mountOptionTypes, map[string]attr.Value{
			"name":       types.StringValue("workspace"),
			"mount_path": types.StringValue("/workspace"),
			"sub_path":   types.StringValue("/sub"),
			"read_only":  types.BoolValue(true),
		}),
		CustomConfiguration: mustObject(t, customTypes, map[string]attr.Value{
			"image":               types.StringValue("img:1"),
			"image_registry_type": types.StringValue("personal"),
			"command":             mustStringSet(t, "/bin/sh"),
			"args":                mustStringList(t, "-c"),
			"env": mustList(t, envTypes, map[string]attr.Value{
				"name":  types.StringValue("K"),
				"value": types.StringValue("V"),
			}),
			"ports": mustList(t, portTypes, map[string]attr.Value{
				"name":     types.StringValue("http"),
				"port":     types.Int64Value(8080),
				"protocol": types.StringValue("TCP"),
			}),
			"resources":        types.ObjectNull(resourcesTypes),
			"probe":            types.ObjectNull(probeTypes),
			"dns_config":       types.ObjectNull(dnsConfigTypes),
			"security_context": types.ObjectNull(securityContextTypes),
		}),
		AuthMode: types.StringValue("DEFAULT"),
		Metadata: mustList(t, metadataTypes, map[string]attr.Value{
			"name":  types.StringValue("run_id"),
			"value": types.StringValue("tf-example"),
		}),
		NetworkMode:           types.StringNull(),
		Persistent:            types.BoolNull(),
		Status:                types.StringNull(),
		ExpiresAt:             types.StringNull(),
		StopReason:            types.StringNull(),
		CreateTime:            types.StringNull(),
		UpdateTime:            types.StringNull(),
		ComputerConfiguration: types.ObjectNull(computerTypes),
	}
}

// sampleInstance returns the cloud side representation used by the describe
// mock. Every value differs from the plan where the cloud may normalize it.
func sampleInstance() *agsv20250920.SandboxInstance {
	return &agsv20250920.SandboxInstance{
		InstanceId:     agsString(agsSandboxInstanceTestId),
		ToolId:         agsString(agsSandboxInstanceTestToolId),
		ToolName:       agsString(agsSandboxInstanceTestTool),
		Status:         agsString("RUNNING"),
		Persistent:     agsBool(true),
		TimeoutSeconds: helperUint64(3600),
		ExpiresAt:      agsString("2026-01-01T01:00:00Z"),
		StopReason:     agsString("manual"),
		CreateTime:     agsString("2026-01-01T00:00:00Z"),
		UpdateTime:     agsString("2026-01-01T00:30:00Z"),
		NetworkMode:    agsString("PUBLIC"),
		AuthMode:       agsString("TOKEN"),
		Metadata: []*agsv20250920.MetadataVar{
			{Name: agsString("run_id"), Value: agsString("cloud-value")},
		},
		ComputerConfiguration: &agsv20250920.ComputerConfiguration{
			OSWorldConfiguration: &agsv20250920.OSWorldConfiguration{Version: agsString("osworld2")},
		},
	}
}

func describeSandboxInstancesResponse(instances ...*agsv20250920.SandboxInstance) *agsv20250920.DescribeSandboxInstancesResponse {
	return &agsv20250920.DescribeSandboxInstancesResponse{
		Response: &agsv20250920.DescribeSandboxInstancesResponseParams{InstanceSet: instances},
	}
}

// nullModel returns a model whose configurable attributes are all unset, so a
// read only fills what the cloud reports.
func (c *agsSandboxInstanceTestContext) nullModel(t *testing.T) svcags.AgsSandboxInstanceResourceModel {
	t.Helper()

	customTypes, computerTypes, mountOptionTypes, metadataTypes := c.agsSandboxInstanceAttributeTypes(t)

	return svcags.AgsSandboxInstanceResourceModel{
		Id:                    types.StringValue(agsSandboxInstanceTestId),
		InstanceId:            types.StringValue(agsSandboxInstanceTestId),
		ToolId:                types.StringNull(),
		ToolName:              types.StringValue(agsSandboxInstanceTestTool),
		SourceInstanceId:      types.StringNull(),
		SnapshotId:            types.StringNull(),
		Timeout:               types.StringNull(),
		MountOptions:          types.ListNull(types.ObjectType{AttrTypes: mountOptionTypes}),
		CustomConfiguration:   types.ObjectNull(customTypes),
		AuthMode:              types.StringNull(),
		Metadata:              types.ListNull(types.ObjectType{AttrTypes: metadataTypes}),
		NetworkMode:           types.StringNull(),
		Persistent:            types.BoolNull(),
		Status:                types.StringNull(),
		ExpiresAt:             types.StringNull(),
		StopReason:            types.StringNull(),
		CreateTime:            types.StringNull(),
		UpdateTime:            types.StringNull(),
		ComputerConfiguration: types.ObjectNull(computerTypes),
	}
}

// TestAgsSandboxInstanceSchema verifies the attribute set and its flags.
func TestAgsSandboxInstanceSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	top := testContext.schema.Attributes

	if len(top) != 19 {
		t.Fatalf("schema has %d attributes, want 19", len(top))
	}

	computed := []string{
		"id", "instance_id", "network_mode", "persistent", "status",
		"expires_at", "stop_reason", "create_time", "update_time", "computer_configuration",
	}
	for _, name := range computed {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsComputed() {
			t.Fatalf("attribute %q must be computed", name)
		}
		if attribute.IsRequired() || attribute.IsOptional() {
			t.Fatalf("attribute %q must not be configurable", name)
		}
	}

	// Only the attributes covered by UpdateSandboxInstance may stay in place.
	for _, name := range []string{"timeout", "metadata"} {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsOptional() {
			t.Fatalf("attribute %q must be optional", name)
		}
	}
	// timeout is refreshed from the cloud side timeout in seconds: it is computed
	// and carries the compatibility modifier, which never replaces the resource.
	timeout := top["timeout"].(schema.StringAttribute)
	if !timeout.Computed {
		t.Fatalf("attribute timeout must be computed")
	}
	if len(timeout.PlanModifiers) == 0 {
		t.Fatalf("attribute timeout must carry a plan modifier")
	}
	if modifiers := top["metadata"].(schema.ListNestedAttribute).PlanModifiers; len(modifiers) != 0 {
		t.Fatalf("attribute metadata must not require replacement")
	}

	mountOption := top["mount_options"].(schema.ListNestedAttribute)
	if modifiers := mountOption.PlanModifiers; len(modifiers) == 0 {
		t.Fatalf("attribute mount_options must require replacement")
	}
}

// TestAgsSandboxInstanceCreate verifies the start request and the read back
// state.
func TestAgsSandboxInstanceCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	var capturedRequest *agsv20250920.StartSandboxInstanceRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"StartSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StartSandboxInstanceRequest) (*agsv20250920.StartSandboxInstanceResponse, error) {
			capturedRequest = request
			return &agsv20250920.StartSandboxInstanceResponse{
				Response: &agsv20250920.StartSandboxInstanceResponseParams{
					Instance: &agsv20250920.SandboxInstance{InstanceId: agsString(agsSandboxInstanceTestId)},
				},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			assertStringPointerList(t, "describe request InstanceIds", request.InstanceIds, agsSandboxInstanceTestId)
			return describeSandboxInstancesResponse(sampleInstance()), nil
		},
	)

	model := testContext.sampleModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", createResp.Diagnostics.Errors())
	}

	if capturedRequest == nil {
		t.Fatalf("StartSandboxInstance was not called")
	}
	if capturedRequest.ToolName == nil || *capturedRequest.ToolName != agsSandboxInstanceTestTool {
		t.Fatalf("StartSandboxInstance ToolName is not %q", agsSandboxInstanceTestTool)
	}
	if capturedRequest.ToolId != nil {
		t.Fatalf("StartSandboxInstance ToolId must not be sent")
	}
	if capturedRequest.Timeout == nil || *capturedRequest.Timeout != "30m" {
		t.Fatalf("StartSandboxInstance Timeout is not 30m")
	}
	if capturedRequest.AuthMode == nil || *capturedRequest.AuthMode != "DEFAULT" {
		t.Fatalf("StartSandboxInstance AuthMode is not DEFAULT")
	}
	// client_token is generated internally as an idempotency key.
	if capturedRequest.ClientToken == nil || *capturedRequest.ClientToken == "" {
		t.Fatalf("StartSandboxInstance ClientToken is empty")
	}
	if len(capturedRequest.MountOptions) != 1 {
		t.Fatalf("StartSandboxInstance MountOptions has %d elements, want 1", len(capturedRequest.MountOptions))
	}
	mountOption := capturedRequest.MountOptions[0]
	if mountOption.Name == nil || *mountOption.Name != "workspace" {
		t.Fatalf("MountOption.Name is not workspace")
	}
	if mountOption.MountPath == nil || *mountOption.MountPath != "/workspace" {
		t.Fatalf("MountOption.MountPath is not /workspace")
	}
	if mountOption.SubPath == nil || *mountOption.SubPath != "/sub" {
		t.Fatalf("MountOption.SubPath is not /sub")
	}
	if mountOption.ReadOnly == nil || !*mountOption.ReadOnly {
		t.Fatalf("MountOption.ReadOnly is not true")
	}
	if capturedRequest.CustomConfiguration == nil || capturedRequest.CustomConfiguration.Image == nil ||
		*capturedRequest.CustomConfiguration.Image != "img:1" {
		t.Fatalf("CustomConfiguration.Image is not img:1")
	}
	assertStringPointerList(t, "CustomConfiguration.Command", capturedRequest.CustomConfiguration.Command, "/bin/sh")
	if len(capturedRequest.Metadata) != 1 {
		t.Fatalf("StartSandboxInstance Metadata has %d elements, want 1", len(capturedRequest.Metadata))
	}
	if capturedRequest.Metadata[0].Name == nil || *capturedRequest.Metadata[0].Name != "run_id" {
		t.Fatalf("Metadata.Name is not run_id")
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if state.Id.ValueString() != agsSandboxInstanceTestId {
		t.Fatalf("state id is %q, want %q", state.Id.ValueString(), agsSandboxInstanceTestId)
	}
	if state.InstanceId.ValueString() != agsSandboxInstanceTestId {
		t.Fatalf("state instance_id is %q, want %q", state.InstanceId.ValueString(), agsSandboxInstanceTestId)
	}
	// status and auth_mode come from the describe answer, not from the plan.
	if state.Status.ValueString() != "RUNNING" {
		t.Fatalf("state status is %q, want RUNNING", state.Status.ValueString())
	}
	if state.AuthMode.ValueString() != "TOKEN" {
		t.Fatalf("state auth_mode is %q, want TOKEN", state.AuthMode.ValueString())
	}
	if state.NetworkMode.ValueString() != "PUBLIC" {
		t.Fatalf("state network_mode is %q, want PUBLIC", state.NetworkMode.ValueString())
	}
	if state.Timeout.ValueString() != "30m" {
		t.Fatalf("state timeout is %q, want the configured 30m", state.Timeout.ValueString())
	}
	if !state.Persistent.ValueBool() {
		t.Fatalf("state persistent is not true")
	}
	// metadata is reported back by the cloud, so the cloud value wins for the
	// configured entries (the system owned ones are filtered out).
	elements := listElements(t, map[string]attr.Value{"metadata": state.Metadata}, "metadata")
	if len(elements) != 1 {
		t.Fatalf("state metadata has %d elements, want 1", len(elements))
	}
	if value := elements[0].(types.Object).Attributes()["value"].(types.String); value.ValueString() != "cloud-value" {
		t.Fatalf("state metadata value is %q, want cloud-value", value.ValueString())
	}
}

// TestAgsSandboxInstanceCreateIgnoresReservedMetadata verifies that the entries
// the cloud adds on its own (x-ags- prefix) are filtered out of the state. They
// are not part of the configuration, so keeping them makes Terraform fail the
// apply with "new element N has appeared".
func TestAgsSandboxInstanceCreateIgnoresReservedMetadata(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	patches.ApplyMethodReturn(
		testContext.client,
		"StartSandboxInstanceWithContext",
		&agsv20250920.StartSandboxInstanceResponse{
			Response: &agsv20250920.StartSandboxInstanceResponseParams{
				Instance: &agsv20250920.SandboxInstance{InstanceId: agsString(agsSandboxInstanceTestId)},
			},
		},
		nil,
	)

	instance := sampleInstance()
	instance.Metadata = []*agsv20250920.MetadataVar{
		{Name: agsString("run_id"), Value: agsString("tf-example")},
		{Name: agsString("x-ags-precreate-match"), Value: agsString("skipped")},
	}
	patches.ApplyMethodReturn(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		describeSandboxInstancesResponse(instance),
		nil,
	)

	model := testContext.sampleModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", createResp.Diagnostics.Errors())
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	elements := listElements(t, map[string]attr.Value{"metadata": state.Metadata}, "metadata")
	if len(elements) != 1 {
		t.Fatalf("state metadata has %d elements, want the single planned one", len(elements))
	}
	if name := elements[0].(types.Object).Attributes()["name"].(types.String); name.ValueString() != "run_id" {
		t.Fatalf("state metadata name is %q, want run_id", name.ValueString())
	}
}

// TestAgsSandboxInstanceCreateEmptyInstanceId verifies that an empty id from
// the API is reported instead of being written into state.
func TestAgsSandboxInstanceCreateEmptyInstanceId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	patches.ApplyMethodReturn(
		testContext.client,
		"StartSandboxInstanceWithContext",
		&agsv20250920.StartSandboxInstanceResponse{
			Response: &agsv20250920.StartSandboxInstanceResponseParams{
				Instance: &agsv20250920.SandboxInstance{},
			},
		},
		nil,
	)

	model := testContext.sampleModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)

	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for an empty instance id")
	}
	if !createResp.State.Raw.IsNull() {
		t.Fatalf("state must stay empty when the instance id is missing")
	}
}

// TestAgsSandboxInstanceCreateApiError verifies the API error path.
func TestAgsSandboxInstanceCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	patches.ApplyMethodFunc(
		testContext.client,
		"StartSandboxInstanceWithContext",
		func(_ context.Context, _ *agsv20250920.StartSandboxInstanceRequest) (*agsv20250920.StartSandboxInstanceResponse, error) {
			return nil, fmt.Errorf("api error")
		},
	)

	model := testContext.sampleModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)

	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic on API failure")
	}
}

// TestAgsSandboxInstanceRead verifies the read back of every reported field.
func TestAgsSandboxInstanceRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			assertStringPointerList(t, "describe request InstanceIds", request.InstanceIds, agsSandboxInstanceTestId)
			return describeSandboxInstancesResponse(sampleInstance()), nil
		},
	)

	model := testContext.nullModel(t)

	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if state.Status.ValueString() != "RUNNING" {
		t.Fatalf("state status is %q, want RUNNING", state.Status.ValueString())
	}
	// timeout is answered in seconds only: the read renders it back.
	if state.Timeout.ValueString() != "1h" {
		t.Fatalf("state timeout is %q, want the rendered 1h", state.Timeout.ValueString())
	}
	if state.StopReason.ValueString() != "manual" {
		t.Fatalf("state stop_reason is %q, want manual", state.StopReason.ValueString())
	}
	if state.ExpiresAt.ValueString() != "2026-01-01T01:00:00Z" {
		t.Fatalf("state expires_at is %q", state.ExpiresAt.ValueString())
	}
	if state.ToolId.ValueString() != agsSandboxInstanceTestToolId {
		t.Fatalf("state tool_id is %q, want %q", state.ToolId.ValueString(), agsSandboxInstanceTestToolId)
	}
	if state.ComputerConfiguration.IsNull() {
		t.Fatalf("state computer_configuration must be filled from the cloud answer")
	}
	// tool_name is configurable: the cloud answer wins here, the previous value
	// is only kept when the API omits the field.
	if state.ToolName.ValueString() != agsSandboxInstanceTestTool {
		t.Fatalf("state tool_name is %q, want %q", state.ToolName.ValueString(), agsSandboxInstanceTestTool)
	}
}

// TestAgsSandboxInstanceReadNotFound verifies that a missing instance is
// removed from state.
func TestAgsSandboxInstanceReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	patches.ApplyMethodReturn(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		describeSandboxInstancesResponse(),
		nil,
	)

	model := testContext.nullModel(t)

	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the instance does not exist any more")
	}
}

// TestAgsSandboxInstanceReadStopped verifies that a stopped instance is removed
// from state so Terraform plans a recreation.
func TestAgsSandboxInstanceReadStopped(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	stopped := sampleInstance()
	stopped.Status = agsString("STOPPED")
	patches.ApplyMethodReturn(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		describeSandboxInstancesResponse(stopped),
		nil,
	)

	model := testContext.nullModel(t)

	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the instance is stopped")
	}
}

// TestAgsSandboxInstanceReadApiError verifies the API error path of Read.
func TestAgsSandboxInstanceReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, _ *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return nil, fmt.Errorf("api error")
		},
	)

	model := testContext.nullModel(t)

	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic on API failure")
	}
}

// TestAgsSandboxInstanceUpdate verifies that only the update covered attributes
// are transmitted and that the read back wins over the plan.
func TestAgsSandboxInstanceUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	var capturedRequest *agsv20250920.UpdateSandboxInstanceRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"UpdateSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.UpdateSandboxInstanceRequest) (*agsv20250920.UpdateSandboxInstanceResponse, error) {
			capturedRequest = request
			return &agsv20250920.UpdateSandboxInstanceResponse{
				Response: &agsv20250920.UpdateSandboxInstanceResponseParams{},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, _ *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(sampleInstance()), nil
		},
	)

	stateModel := testContext.sampleModel(t)
	stateModel.Id = types.StringValue(agsSandboxInstanceTestId)
	stateModel.InstanceId = types.StringValue(agsSandboxInstanceTestId)

	planModel := testContext.sampleModel(t)
	planModel.Id = types.StringValue(agsSandboxInstanceTestId)
	planModel.InstanceId = types.StringValue(agsSandboxInstanceTestId)
	planModel.Timeout = types.StringValue("1h")

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{Plan: testContext.plan(t, planModel), State: testContext.state(t, stateModel)},
		updateResp,
	)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", updateResp.Diagnostics.Errors())
	}

	if capturedRequest == nil {
		t.Fatalf("UpdateSandboxInstance was not called")
	}
	if capturedRequest.InstanceId == nil || *capturedRequest.InstanceId != agsSandboxInstanceTestId {
		t.Fatalf("UpdateSandboxInstance InstanceId is not %q", agsSandboxInstanceTestId)
	}
	if capturedRequest.Timeout == nil || *capturedRequest.Timeout != "1h" {
		t.Fatalf("UpdateSandboxInstance Timeout is not 1h")
	}
	if len(capturedRequest.Metadata) != 1 {
		t.Fatalf("UpdateSandboxInstance Metadata has %d elements, want 1", len(capturedRequest.Metadata))
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := updateResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if state.Status.ValueString() != "RUNNING" {
		t.Fatalf("state status is %q, want the read back value RUNNING", state.Status.ValueString())
	}
	if state.NetworkMode.ValueString() != "PUBLIC" {
		t.Fatalf("state network_mode is %q, want the read back value PUBLIC", state.NetworkMode.ValueString())
	}
}

// TestAgsSandboxInstanceUpdateResendsReservedMetadata verifies that the system
// owned metadata entries (x-ags- prefix) are sent back with an update. The
// request replaces the whole list, so dropping them would break the bookkeeping
// the cloud side keeps in the metadata.
func TestAgsSandboxInstanceUpdateResendsReservedMetadata(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	var capturedRequest *agsv20250920.UpdateSandboxInstanceRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"UpdateSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.UpdateSandboxInstanceRequest) (*agsv20250920.UpdateSandboxInstanceResponse, error) {
			capturedRequest = request
			return &agsv20250920.UpdateSandboxInstanceResponse{
				Response: &agsv20250920.UpdateSandboxInstanceResponseParams{},
			}, nil
		},
	)

	instance := sampleInstance()
	instance.Metadata = []*agsv20250920.MetadataVar{
		{Name: agsString("run_id"), Value: agsString("tf-example")},
		{Name: agsString("x-ags-precreate-match"), Value: agsString("skipped")},
	}
	patches.ApplyMethodReturn(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		describeSandboxInstancesResponse(instance),
		nil,
	)

	stateModel := testContext.sampleModel(t)
	stateModel.Id = types.StringValue(agsSandboxInstanceTestId)
	stateModel.InstanceId = types.StringValue(agsSandboxInstanceTestId)

	planModel := testContext.sampleModel(t)
	planModel.Id = types.StringValue(agsSandboxInstanceTestId)
	planModel.InstanceId = types.StringValue(agsSandboxInstanceTestId)
	planModel.Timeout = types.StringValue("1h")

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{Plan: testContext.plan(t, planModel), State: testContext.state(t, stateModel)},
		updateResp,
	)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", updateResp.Diagnostics.Errors())
	}

	if capturedRequest == nil {
		t.Fatalf("UpdateSandboxInstance was not called")
	}
	if len(capturedRequest.Metadata) != 2 {
		t.Fatalf("UpdateSandboxInstance Metadata has %d elements, want the configured one plus the system one", len(capturedRequest.Metadata))
	}
	if name := capturedRequest.Metadata[0].Name; name == nil || *name != "run_id" {
		t.Fatalf("UpdateSandboxInstance Metadata[0].Name is not run_id")
	}
	if name := capturedRequest.Metadata[1].Name; name == nil || *name != "x-ags-precreate-match" {
		t.Fatalf("UpdateSandboxInstance Metadata[1].Name is not x-ags-precreate-match")
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := updateResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	elements := listElements(t, map[string]attr.Value{"metadata": state.Metadata}, "metadata")
	if len(elements) != 1 {
		t.Fatalf("state metadata has %d elements, want the single configured one", len(elements))
	}
}

// TestAgsSandboxInstanceUpdateWithoutChange verifies that the write API is not
// called when no update covered attribute changed.
func TestAgsSandboxInstanceUpdateWithoutChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	called := false
	patches.ApplyMethodFunc(
		testContext.client,
		"UpdateSandboxInstanceWithContext",
		func(_ context.Context, _ *agsv20250920.UpdateSandboxInstanceRequest) (*agsv20250920.UpdateSandboxInstanceResponse, error) {
			called = true
			return &agsv20250920.UpdateSandboxInstanceResponse{
				Response: &agsv20250920.UpdateSandboxInstanceResponseParams{},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, _ *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(sampleInstance()), nil
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsSandboxInstanceTestId)
	model.InstanceId = types.StringValue(agsSandboxInstanceTestId)

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{Plan: testContext.plan(t, model), State: testContext.state(t, model)},
		updateResp,
	)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", updateResp.Diagnostics.Errors())
	}
	if called {
		t.Fatalf("UpdateSandboxInstance must not be called when nothing changed")
	}
}

// TestAgsSandboxInstanceUpdateApiError verifies the API error path of Update.
func TestAgsSandboxInstanceUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	patches.ApplyMethodFunc(
		testContext.client,
		"UpdateSandboxInstanceWithContext",
		func(_ context.Context, _ *agsv20250920.UpdateSandboxInstanceRequest) (*agsv20250920.UpdateSandboxInstanceResponse, error) {
			return nil, fmt.Errorf("api error")
		},
	)

	stateModel := testContext.sampleModel(t)
	stateModel.Id = types.StringValue(agsSandboxInstanceTestId)
	stateModel.InstanceId = types.StringValue(agsSandboxInstanceTestId)

	planModel := testContext.sampleModel(t)
	planModel.Id = types.StringValue(agsSandboxInstanceTestId)
	planModel.InstanceId = types.StringValue(agsSandboxInstanceTestId)
	planModel.Timeout = types.StringValue("1h")

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{Plan: testContext.plan(t, planModel), State: testContext.state(t, stateModel)},
		updateResp,
	)
	if !updateResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic on API failure")
	}
}

// TestAgsSandboxInstanceDelete verifies the stop request.
func TestAgsSandboxInstanceDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	var capturedRequest *agsv20250920.StopSandboxInstanceRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"StopSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StopSandboxInstanceRequest) (*agsv20250920.StopSandboxInstanceResponse, error) {
			capturedRequest = request
			return &agsv20250920.StopSandboxInstanceResponse{
				Response: &agsv20250920.StopSandboxInstanceResponseParams{},
			}, nil
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsSandboxInstanceTestId)
	model.InstanceId = types.StringValue(agsSandboxInstanceTestId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(context.Background(), resource.DeleteRequest{State: testContext.state(t, model)}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", deleteResp.Diagnostics.Errors())
	}

	if capturedRequest == nil {
		t.Fatalf("StopSandboxInstance was not called")
	}
	if capturedRequest.InstanceId == nil || *capturedRequest.InstanceId != agsSandboxInstanceTestId {
		t.Fatalf("StopSandboxInstance InstanceId is not %q", agsSandboxInstanceTestId)
	}
}

// TestAgsSandboxInstanceDeleteApiError verifies the API error path of Delete.
func TestAgsSandboxInstanceDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	patches.ApplyMethodFunc(
		testContext.client,
		"StopSandboxInstanceWithContext",
		func(_ context.Context, _ *agsv20250920.StopSandboxInstanceRequest) (*agsv20250920.StopSandboxInstanceResponse, error) {
			return nil, fmt.Errorf("api error")
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsSandboxInstanceTestId)
	model.InstanceId = types.StringValue(agsSandboxInstanceTestId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(context.Background(), resource.DeleteRequest{State: testContext.state(t, model)}, deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic on API failure")
	}
}

// customConfigurationValue builds a custom_configuration object, so the plan
// modifier tests can spell out only the attributes that matter for them.
func (c *agsSandboxInstanceTestContext) customConfigurationValue(t *testing.T, overrides map[string]attr.Value) types.Object {
	t.Helper()

	customTypes, _, _, _ := c.agsSandboxInstanceAttributeTypes(t)
	customAttributes := c.schema.Attributes["custom_configuration"].(schema.SingleNestedAttribute).Attributes
	envTypes := listElementAttributeTypes(t, customAttributes, "env")
	portTypes := listElementAttributeTypes(t, customAttributes, "ports")

	values := map[string]attr.Value{
		"image":               types.StringValue("img:1"),
		"image_registry_type": types.StringValue("personal"),
		"command":             mustStringSet(t, "/bin/sh"),
		"args":                mustStringList(t, "-c"),
		"env": mustList(t, envTypes, map[string]attr.Value{
			"name":  types.StringValue("K"),
			"value": types.StringValue("V"),
		}),
		"ports": mustList(t, portTypes, map[string]attr.Value{
			"name":     types.StringValue("http"),
			"port":     types.Int64Value(8080),
			"protocol": types.StringValue("TCP"),
		}),
		"resources":        types.ObjectNull(objectAttributeTypes(t, customAttributes, "resources")),
		"probe":            types.ObjectNull(objectAttributeTypes(t, customAttributes, "probe")),
		"dns_config":       types.ObjectNull(objectAttributeTypes(t, customAttributes, "dns_config")),
		"security_context": types.ObjectNull(objectAttributeTypes(t, customAttributes, "security_context")),
	}
	for name, value := range overrides {
		values[name] = value
	}

	return mustObject(t, customTypes, values)
}

// TestAgsSandboxInstanceCustomConfigurationPlanModifier verifies that the
// normalized cloud answer (1000m instead of 1, 1024Mi instead of 1Gi, the
// storage inherited from the tool) does not require a replacement, while a real
// change still does. The planned value keeps the configuration spelling, so the
// state converges back to it instead of drifting to the normalized cloud form on
// every plan.
func TestAgsSandboxInstanceCustomConfigurationPlanModifier(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	customAttributes := testContext.schema.Attributes["custom_configuration"].(schema.SingleNestedAttribute).Attributes
	modifiers := testContext.schema.Attributes["custom_configuration"].(schema.SingleNestedAttribute).PlanModifiers
	if len(modifiers) == 0 {
		t.Fatalf("custom_configuration has no plan modifiers")
	}
	resourcesTypes := objectAttributeTypes(t, customAttributes, "resources")
	customTypes, _, _, _ := testContext.agsSandboxInstanceAttributeTypes(t)

	configValue := testContext.customConfigurationValue(t, map[string]attr.Value{
		"resources": mustObject(t, resourcesTypes, map[string]attr.Value{
			"cpu":     types.StringValue("1"),
			"memory":  types.StringValue("1Gi"),
			"storage": types.StringNull(),
		}),
	})
	stateValue := testContext.customConfigurationValue(t, map[string]attr.Value{
		"resources": mustObject(t, resourcesTypes, map[string]attr.Value{
			"cpu":     types.StringValue("1000m"),
			"memory":  types.StringValue("1024Mi"),
			"storage": types.StringValue("1024Mi"),
		}),
	})
	changedValue := testContext.customConfigurationValue(t, map[string]attr.Value{
		"image": types.StringValue("img:2"),
	})

	cases := []struct {
		name    string
		config  types.Object
		state   types.Object
		replace bool
	}{
		{name: "normalized cloud answer", config: configValue, state: stateValue, replace: false},
		{name: "changed image", config: changedValue, state: stateValue, replace: true},
		{name: "unset configuration", config: types.ObjectNull(customTypes), state: stateValue, replace: false},
		{name: "no prior state", config: configValue, state: types.ObjectNull(customTypes), replace: true},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			resp := &planmodifier.ObjectResponse{PlanValue: testCase.config}
			modifiers[0].PlanModifyObject(
				context.Background(),
				planmodifier.ObjectRequest{
					Config:      tfsdk.Config{Schema: testContext.schema},
					Plan:        testContext.plan(t, testContext.sampleModel(t)),
					State:       testContext.state(t, testContext.sampleModel(t)),
					ConfigValue: testCase.config,
					PlanValue:   testCase.config,
					StateValue:  testCase.state,
				},
				resp,
			)

			if resp.RequiresReplace != testCase.replace {
				t.Fatalf("requires replacement is %t, want %t", resp.RequiresReplace, testCase.replace)
			}
			if !resp.PlanValue.Equal(testCase.config) {
				t.Fatalf("the plan modifier must keep the configured value")
			}
		})
	}
}

// TestAgsSandboxInstanceMountOptionsPlanModifier verifies that the mount values
// inherited from the tool are ignored when the configuration only spells out the
// storage name.
func TestAgsSandboxInstanceMountOptionsPlanModifier(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	_, _, mountOptionTypes, _ := testContext.agsSandboxInstanceAttributeTypes(t)
	modifiers := testContext.schema.Attributes["mount_options"].(schema.ListNestedAttribute).PlanModifiers
	if len(modifiers) == 0 {
		t.Fatalf("mount_options has no plan modifiers")
	}

	configValue := mustList(t, mountOptionTypes, map[string]attr.Value{
		"name":       types.StringValue("workspace"),
		"mount_path": types.StringNull(),
		"sub_path":   types.StringNull(),
		"read_only":  types.BoolNull(),
	})
	stateValue := mustList(t, mountOptionTypes, map[string]attr.Value{
		"name":       types.StringValue("workspace"),
		"mount_path": types.StringValue("/workspace"),
		"sub_path":   types.StringValue("/sub"),
		"read_only":  types.BoolValue(true),
	})

	resp := &planmodifier.ListResponse{PlanValue: configValue}
	modifiers[0].PlanModifyList(
		context.Background(),
		planmodifier.ListRequest{
			Config:      tfsdk.Config{Schema: testContext.schema},
			Plan:        testContext.plan(t, testContext.sampleModel(t)),
			State:       testContext.state(t, testContext.sampleModel(t)),
			ConfigValue: configValue,
			PlanValue:   configValue,
			StateValue:  stateValue,
		},
		resp,
	)
	if resp.RequiresReplace {
		t.Fatalf("the inherited mount values must not require a replacement")
	}
	if !resp.PlanValue.Equal(configValue) {
		t.Fatalf("the plan modifier must keep the configured value")
	}
}

// TestAgsSandboxInstanceReadFillsCreateOnlyInputs verifies that a read on a
// state built from the id only (the import path) adopts the create-only inputs,
// so the next plan does not propose a replacement.
func TestAgsSandboxInstanceReadFillsCreateOnlyInputs(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	instance := sampleInstance()
	instance.CustomConfiguration = &agsv20250920.CustomConfigurationDetail{
		Image: agsString("img:1"),
		Resources: &agsv20250920.ResourceConfiguration{
			CPU:     agsString("1000m"),
			Memory:  agsString("1024Mi"),
			Storage: agsString("1024Mi"),
		},
	}
	instance.MountOptions = []*agsv20250920.MountOption{
		{Name: agsString("workspace"), MountPath: agsString("/workspace"), ReadOnly: agsBool(false)},
	}
	patches.ApplyMethodReturn(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		describeSandboxInstancesResponse(instance),
		nil,
	)

	model := testContext.nullModel(t)
	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	// timeout is rendered from the cloud side seconds when the state has none.
	if state.Timeout.ValueString() != "1h" {
		t.Fatalf("state timeout is %q, want the rendered 1h", state.Timeout.ValueString())
	}
	if state.CustomConfiguration.IsNull() {
		t.Fatalf("custom_configuration must be filled from the cloud answer")
	}
	if image := state.CustomConfiguration.Attributes()["image"].(types.String); image.ValueString() != "img:1" {
		t.Fatalf("state image is %q, want img:1", image.ValueString())
	}

	elements := listElements(t, map[string]attr.Value{"mount_options": state.MountOptions}, "mount_options")
	if len(elements) != 1 {
		t.Fatalf("state mount_options has %d elements, want 1", len(elements))
	}
	if name := elements[0].(types.Object).Attributes()["name"].(types.String); name.ValueString() != "workspace" {
		t.Fatalf("state mount option name is %q, want workspace", name.ValueString())
	}
}

// TestAgsSandboxInstanceReadKeepsConfiguredTimeout verifies that a read does not
// overwrite the configured timeout: the cloud answer is adopted only when the
// state carries no value yet.
func TestAgsSandboxInstanceReadKeepsConfiguredTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	patches.ApplyMethodReturn(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		describeSandboxInstancesResponse(sampleInstance()),
		nil,
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsSandboxInstanceTestId)
	model.InstanceId = types.StringValue(agsSandboxInstanceTestId)
	model.Timeout = types.StringValue("30m")

	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if state.Timeout.ValueString() != "30m" {
		t.Fatalf("state timeout is %q, want the configured 30m", state.Timeout.ValueString())
	}
}

// TestAgsSandboxInstanceReadKeepsConfiguredCustomConfiguration verifies that a
// read keeps the configuration spelling of custom_configuration while the cloud
// answer is equivalent to it. Writing the normalized answer back would leave the
// state different from the configuration and make every plan an in-place update.
func TestAgsSandboxInstanceReadKeepsConfiguredCustomConfiguration(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsSandboxInstanceTestId)
	model.InstanceId = types.StringValue(agsSandboxInstanceTestId)
	customAttributes := testContext.schema.Attributes["custom_configuration"].(schema.SingleNestedAttribute).Attributes
	resourcesTypes := objectAttributeTypes(t, customAttributes, "resources")
	model.CustomConfiguration = testContext.customConfigurationValue(t, map[string]attr.Value{
		"resources": mustObject(t, resourcesTypes, map[string]attr.Value{
			"cpu":     types.StringValue("1"),
			"memory":  types.StringValue("1Gi"),
			"storage": types.StringNull(),
		}),
	})

	// The cloud answers the same values in their normalized form plus the
	// storage inherited from the tool.
	instance := sampleInstance()
	instance.CustomConfiguration = &agsv20250920.CustomConfigurationDetail{
		Image:             agsString("img:1"),
		ImageRegistryType: agsString("personal"),
		Command:           []*string{agsString("/bin/sh")},
		Args:              []*string{agsString("-c")},
		Env: []*agsv20250920.EnvVar{
			{Name: agsString("K"), Value: agsString("V")},
		},
		Ports: []*agsv20250920.PortConfiguration{
			{Name: agsString("http"), Port: agsInt64(8080), Protocol: agsString("TCP")},
		},
		Resources: &agsv20250920.ResourceConfiguration{
			CPU:     agsString("1000m"),
			Memory:  agsString("1024Mi"),
			Storage: agsString("1024Mi"),
		},
	}
	patches.ApplyMethodReturn(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		describeSandboxInstancesResponse(instance),
		nil,
	)

	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if !state.CustomConfiguration.Equal(model.CustomConfiguration) {
		t.Fatalf("state custom_configuration was rewritten from the cloud answer")
	}
}

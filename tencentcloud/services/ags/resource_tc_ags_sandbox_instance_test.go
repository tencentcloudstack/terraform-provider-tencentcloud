package ags_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

const (
	agsTestInstanceId     = "instance-001"
	agsTestSourceInstance = "instance-000"
	agsTestSnapshotId     = "snap-001"
)

// agsSandboxInstanceTestContext wires a resource instance to a patched SDK client.
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

// sampleInstanceModel returns a plan model with every configurable leaf populated.
func (c *agsSandboxInstanceTestContext) sampleInstanceModel(t *testing.T) svcags.AgsSandboxInstanceResourceModel {
	t.Helper()

	top := c.schema.Attributes

	mountOptionTypes := listElementAttributeTypes(t, top, "mount_options")
	metadataTypes := listElementAttributeTypes(t, top, "metadata")

	customAttributes := top["custom_configuration"].(schema.SingleNestedAttribute).Attributes
	probeAttributes := customAttributes["probe"].(schema.SingleNestedAttribute).Attributes
	securityContextAttributes := customAttributes["security_context"].(schema.SingleNestedAttribute).Attributes
	resourcesTypes := objectAttributeTypes(t, customAttributes, "resources")
	httpGetTypes := objectAttributeTypes(t, probeAttributes, "http_get")
	dnsConfigTypes := objectAttributeTypes(t, customAttributes, "dns_config")
	capabilitiesTypes := objectAttributeTypes(t, securityContextAttributes, "capabilities")
	securityContextTypes := objectAttributeTypes(t, customAttributes, "security_context")
	probeTypes := objectAttributeTypes(t, customAttributes, "probe")
	customConfigurationTypes := objectAttributeTypes(t, top, "custom_configuration")

	storageMountAttributes := top["storage_mounts"].(schema.ListNestedAttribute).NestedObject.Attributes
	storageSourceAttributes := storageMountAttributes["storage_source"].(schema.SingleNestedAttribute).Attributes
	cosTypes := objectAttributeTypes(t, storageSourceAttributes, "cos")
	imageTypes := objectAttributeTypes(t, storageSourceAttributes, "image")
	cfsTypes := objectAttributeTypes(t, storageSourceAttributes, "cfs")
	agentBucketTypes := objectAttributeTypes(t, storageSourceAttributes, "agent_bucket")
	storageSourceTypes := objectAttributeTypes(t, storageMountAttributes, "storage_source")
	volumeTypes := objectAttributeTypes(t, storageMountAttributes, "volume")
	volumeTemplateTypes := objectAttributeTypes(t, storageMountAttributes, "volume_template")
	storageMountTypes := listElementAttributeTypes(t, top, "storage_mounts")

	envTypes := listElementAttributeTypes(t, customAttributes, "env")
	portTypes := listElementAttributeTypes(t, customAttributes, "ports")

	return svcags.AgsSandboxInstanceResourceModel{
		Id:               types.StringNull(),
		InstanceId:       types.StringNull(),
		ToolId:           types.StringValue(agsTestToolId),
		ToolName:         types.StringValue(agsTestToolName),
		SourceInstanceId: types.StringValue(agsTestSourceInstance),
		SnapshotId:       types.StringValue(agsTestSnapshotId),
		Timeout:          types.StringValue("10m"),
		ClientToken:      types.StringValue("token-001"),
		MountOptions: mustList(t, mountOptionTypes, map[string]attr.Value{
			"name":       types.StringValue("workspace"),
			"mount_path": types.StringValue("/workspace"),
			"sub_path":   types.StringValue("/sub"),
			"read_only":  types.BoolValue(true),
		}),
		CustomConfiguration: mustObject(t, customConfigurationTypes, map[string]attr.Value{
			"image":               types.StringValue("img:1"),
			"image_registry_type": types.StringValue("personal"),
			"command":             mustStringSet(t, "/bin/sh"),
			"args":                mustStringSet(t, "-c"),
			"env": mustList(t, envTypes, map[string]attr.Value{
				"name":  types.StringValue("K"),
				"value": types.StringValue("V"),
			}),
			"ports": mustList(t, portTypes, map[string]attr.Value{
				"name":     types.StringValue("http"),
				"port":     types.Int64Value(8080),
				"protocol": types.StringValue("TCP"),
			}),
			"resources": mustObject(t, resourcesTypes, map[string]attr.Value{
				"cpu":                   types.StringValue("1"),
				"memory":                types.StringValue("2Gi"),
				"storage":               types.StringValue("5Gi"),
				"gpu_resource_group_id": types.StringValue("gpusg-1"),
			}),
			"probe": mustObject(t, probeTypes, map[string]attr.Value{
				"http_get": mustObject(t, httpGetTypes, map[string]attr.Value{
					"path":   types.StringValue("/healthz"),
					"port":   types.Int64Value(8080),
					"scheme": types.StringValue("HTTP"),
				}),
				"ready_timeout_ms":  types.Int64Value(1000),
				"probe_timeout_ms":  types.Int64Value(2000),
				"probe_period_ms":   types.Int64Value(3000),
				"success_threshold": types.Int64Value(1),
				"failure_threshold": types.Int64Value(3),
			}),
			"dns_config": mustObject(t, dnsConfigTypes, map[string]attr.Value{
				"servers":  mustStringSet(t, "10.0.0.1"),
				"searches": mustStringSet(t, "example.com"),
				"options":  mustStringSet(t, "ndots:2"),
			}),
			"security_context": mustObject(t, securityContextTypes, map[string]attr.Value{
				"privileged": types.BoolValue(true),
				"capabilities": mustObject(t, capabilitiesTypes, map[string]attr.Value{
					"add":  mustStringSet(t, "SYS_ADMIN"),
					"drop": mustStringSet(t, "ALL"),
				}),
			}),
		}),
		AuthMode: types.StringValue("TOKEN"),
		Metadata: mustList(t, metadataTypes, map[string]attr.Value{
			"name":  types.StringValue("run_id"),
			"value": types.StringValue("42"),
		}),
		StorageMounts: mustList(t, storageMountTypes, map[string]attr.Value{
			"name": types.StringValue("extra"),
			"storage_source": mustObject(t, storageSourceTypes, map[string]attr.Value{
				"cos": mustObject(t, cosTypes, map[string]attr.Value{
					"endpoint":    types.StringValue("cos.example.com"),
					"bucket_name": types.StringValue("bucket"),
					"bucket_path": types.StringValue("/path"),
				}),
				"image": mustObject(t, imageTypes, map[string]attr.Value{
					"reference":           types.StringValue("img:1"),
					"image_registry_type": types.StringValue("personal"),
					"sub_path":            types.StringValue("/sub"),
					"digest":              types.StringNull(),
				}),
				"cfs": mustObject(t, cfsTypes, map[string]attr.Value{
					"file_system_id": types.StringValue("cfs-1"),
					"path":           types.StringValue("/cfs"),
				}),
				"agent_bucket": mustObject(t, agentBucketTypes, map[string]attr.Value{
					"library_id":    types.StringValue("lib-1"),
					"space_id":      types.StringValue("space-1"),
					"access_domain": types.StringValue("domain.example.com"),
				}),
			}),
			"mount_path": types.StringValue("/extra"),
			"read_only":  types.BoolValue(true),
			"volume": mustObject(t, volumeTypes, map[string]attr.Value{
				"volume_name": types.StringValue("vol-name"),
				"volume_id":   types.StringNull(),
			}),
			"volume_template": mustObject(t, volumeTemplateTypes, map[string]attr.Value{
				"volume_template_name": types.StringValue("tpl"),
				"reuse_key":            types.StringValue("${metadata.run_id}"),
			}),
		}),
		StartMode:      types.StringValue("DEFAULT"),
		Status:         types.StringNull(),
		Persistent:     types.BoolNull(),
		TimeoutSeconds: types.Int64Null(),
		ExpiresAt:      types.StringNull(),
		StopReason:     types.StringNull(),
		CreateTime:     types.StringNull(),
		UpdateTime:     types.StringNull(),
		NetworkMode:    types.StringNull(),
	}
}

// sampleInstanceItem returns the cloud side representation used by the mocks.
func sampleInstanceItem() *agsv20250920.SandboxInstance {
	return &agsv20250920.SandboxInstance{
		InstanceId:     agsString(agsTestInstanceId),
		ToolId:         agsString(agsTestToolId),
		ToolName:       agsString(agsTestToolName),
		Status:         agsString("RUNNING"),
		Persistent:     agsBool(false),
		TimeoutSeconds: helperUint64(600),
		ExpiresAt:      agsString("2026-01-01T00:10:00Z"),
		CreateTime:     agsString(agsTestCreateTime),
		UpdateTime:     agsString(agsTestUpdateTime),
		NetworkMode:    agsString("PUBLIC"),
		AuthMode:       agsString("TOKEN"),
		MountOptions: []*agsv20250920.MountOption{
			{
				Name:      agsString("workspace"),
				MountPath: agsString("/workspace"),
				SubPath:   agsString("/sub"),
				ReadOnly:  agsBool(true),
			},
		},
		CustomConfiguration: &agsv20250920.CustomConfigurationDetail{
			Image:             agsString("img:1"),
			ImageRegistryType: agsString("personal"),
			ImageDigest:       agsString("sha256:abc"),
			Command:           []*string{agsString("/bin/sh")},
			Args:              []*string{agsString("-c")},
			Env: []*agsv20250920.EnvVar{
				{Name: agsString("K"), Value: agsString("V")},
			},
			Ports: []*agsv20250920.PortConfiguration{
				{Name: agsString("http"), Port: agsInt64(8080), Protocol: agsString("TCP")},
			},
			Resources: &agsv20250920.ResourceConfiguration{
				CPU:                agsString("1"),
				Memory:             agsString("2Gi"),
				Storage:            agsString("5Gi"),
				GPUResourceGroupId: agsString("gpusg-1"),
			},
			Probe: &agsv20250920.ProbeConfiguration{
				HttpGet: &agsv20250920.HttpGetAction{
					Path:   agsString("/healthz"),
					Port:   agsInt64(8080),
					Scheme: agsString("HTTP"),
				},
				ReadyTimeoutMs:   agsInt64(1000),
				ProbeTimeoutMs:   agsInt64(2000),
				ProbePeriodMs:    agsInt64(3000),
				SuccessThreshold: agsInt64(1),
				FailureThreshold: agsInt64(3),
			},
			DNSConfig: &agsv20250920.DNSConfig{
				Servers:  []*string{agsString("10.0.0.1")},
				Searches: []*string{agsString("example.com")},
				Options:  []*string{agsString("ndots:2")},
			},
			SecurityContext: &agsv20250920.SecurityContext{
				Privileged: agsBool(true),
				Capabilities: &agsv20250920.Capabilities{
					Add:  []*string{agsString("SYS_ADMIN")},
					Drop: []*string{agsString("ALL")},
				},
			},
		},
		Metadata: []*agsv20250920.MetadataVar{
			{Name: agsString("run_id"), Value: agsString("42")},
		},
	}
}

func describeSandboxInstancesResponse(instances ...*agsv20250920.SandboxInstance) *agsv20250920.DescribeSandboxInstancesResponse {
	return &agsv20250920.DescribeSandboxInstancesResponse{
		Response: &agsv20250920.DescribeSandboxInstancesResponseParams{
			InstanceSet: instances,
		},
	}
}

// TestAgsSandboxInstanceSchema verifies the attribute set and its flags.
func TestAgsSandboxInstanceSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	top := testContext.schema.Attributes

	if len(top) != 22 {
		t.Fatalf("schema has %d attributes, want 22", len(top))
	}

	optional := []string{
		"tool_id", "tool_name", "source_instance_id", "snapshot_id", "timeout", "client_token",
		"mount_options", "custom_configuration", "auth_mode", "metadata", "storage_mounts", "start_mode",
	}
	for _, name := range optional {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsOptional() {
			t.Fatalf("attribute %q must be optional", name)
		}
	}

	computed := []string{
		"id", "instance_id", "status", "persistent", "timeout_seconds", "expires_at", "stop_reason",
		"create_time", "update_time", "network_mode",
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

	// Updatable fields are optional+computed and do not force a replacement.
	for _, name := range []string{"timeout", "metadata"} {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsOptional() || !attribute.IsComputed() {
			t.Fatalf("attribute %q must be optional and computed", name)
		}
	}

	// Create-only fields are optional and force a replacement on change.
	for _, name := range []string{"tool_id", "tool_name", "source_instance_id", "snapshot_id", "client_token", "mount_options", "custom_configuration", "auth_mode", "storage_mounts", "start_mode"} {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsOptional() {
			t.Fatalf("attribute %q must be optional", name)
		}
	}
}

// TestAgsSandboxInstanceCreate verifies the start request and the read back state.
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
				Response: &agsv20250920.StartSandboxInstanceResponseParams{Instance: sampleInstanceItem()},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			assertStringPointerList(t, "describe request InstanceIds", request.InstanceIds, agsTestInstanceId)
			return describeSandboxInstancesResponse(sampleInstanceItem()), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	assertNoErrorDiagnostics(t, createResp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("StartSandboxInstance was not called")
	}
	if capturedRequest.ToolId == nil || *capturedRequest.ToolId != agsTestToolId {
		t.Fatalf("StartSandboxInstance ToolId is not %q", agsTestToolId)
	}
	if capturedRequest.ToolName == nil || *capturedRequest.ToolName != agsTestToolName {
		t.Fatalf("StartSandboxInstance ToolName is not %q", agsTestToolName)
	}
	if capturedRequest.SourceInstanceId == nil || *capturedRequest.SourceInstanceId != agsTestSourceInstance {
		t.Fatalf("StartSandboxInstance SourceInstanceId is not %q", agsTestSourceInstance)
	}
	if capturedRequest.SnapshotId == nil || *capturedRequest.SnapshotId != agsTestSnapshotId {
		t.Fatalf("StartSandboxInstance SnapshotId is not %q", agsTestSnapshotId)
	}
	if capturedRequest.Timeout == nil || *capturedRequest.Timeout != "10m" {
		t.Fatalf("StartSandboxInstance Timeout is not 10m")
	}
	// A configured client token is forwarded as it is.
	if capturedRequest.ClientToken == nil || *capturedRequest.ClientToken != "token-001" {
		t.Fatalf("StartSandboxInstance ClientToken is not token-001")
	}
	if capturedRequest.AuthMode == nil || *capturedRequest.AuthMode != "TOKEN" {
		t.Fatalf("StartSandboxInstance AuthMode is not TOKEN")
	}
	if capturedRequest.StartMode == nil || *capturedRequest.StartMode != "DEFAULT" {
		t.Fatalf("StartSandboxInstance StartMode is not DEFAULT")
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

	if len(capturedRequest.Metadata) != 1 || capturedRequest.Metadata[0].Name == nil ||
		*capturedRequest.Metadata[0].Name != "run_id" {
		t.Fatalf("StartSandboxInstance Metadata is not [run_id]")
	}

	customConfiguration := capturedRequest.CustomConfiguration
	if customConfiguration == nil {
		t.Fatalf("StartSandboxInstance CustomConfiguration is nil")
	}
	if customConfiguration.Image == nil || *customConfiguration.Image != "img:1" {
		t.Fatalf("CustomConfiguration.Image is not img:1")
	}
	assertStringPointerList(t, "CustomConfiguration.Command", customConfiguration.Command, "/bin/sh")
	assertStringPointerList(t, "CustomConfiguration.Args", customConfiguration.Args, "-c")
	if len(customConfiguration.Env) != 1 || customConfiguration.Env[0].Name == nil || *customConfiguration.Env[0].Name != "K" {
		t.Fatalf("CustomConfiguration.Env is not [K]")
	}
	if len(customConfiguration.Ports) != 1 || customConfiguration.Ports[0].Port == nil || *customConfiguration.Ports[0].Port != 8080 {
		t.Fatalf("CustomConfiguration.Ports is not [8080]")
	}
	if customConfiguration.Resources == nil || customConfiguration.Resources.GPUResourceGroupId == nil ||
		*customConfiguration.Resources.GPUResourceGroupId != "gpusg-1" {
		t.Fatalf("CustomConfiguration.Resources.GPUResourceGroupId is not gpusg-1")
	}
	if customConfiguration.Probe == nil || customConfiguration.Probe.HttpGet == nil ||
		customConfiguration.Probe.HttpGet.Path == nil || *customConfiguration.Probe.HttpGet.Path != "/healthz" {
		t.Fatalf("CustomConfiguration.Probe.HttpGet.Path is not /healthz")
	}
	if customConfiguration.DNSConfig == nil {
		t.Fatalf("CustomConfiguration.DNSConfig is nil")
	}
	assertStringPointerList(t, "DNSConfig.Servers", customConfiguration.DNSConfig.Servers, "10.0.0.1")
	if customConfiguration.SecurityContext == nil || customConfiguration.SecurityContext.Privileged == nil ||
		!*customConfiguration.SecurityContext.Privileged {
		t.Fatalf("CustomConfiguration.SecurityContext.Privileged is not true")
	}

	if len(capturedRequest.StorageMounts) != 1 {
		t.Fatalf("StartSandboxInstance StorageMounts has %d elements, want 1", len(capturedRequest.StorageMounts))
	}
	storageMount := capturedRequest.StorageMounts[0]
	if storageMount.Name == nil || *storageMount.Name != "extra" {
		t.Fatalf("StorageMount.Name is not extra")
	}
	if storageMount.StorageSource == nil || storageMount.StorageSource.Cos == nil ||
		storageMount.StorageSource.Cos.BucketName == nil || *storageMount.StorageSource.Cos.BucketName != "bucket" {
		t.Fatalf("StorageMount.StorageSource.Cos.BucketName is not bucket")
	}
	if storageMount.StorageSource.Image == nil || storageMount.StorageSource.Image.Digest != nil {
		t.Fatalf("StorageMount.StorageSource.Image.Digest must not be sent")
	}
	if storageMount.VolumeTemplate == nil || storageMount.VolumeTemplate.VolumeTemplateName == nil ||
		*storageMount.VolumeTemplate.VolumeTemplateName != "tpl" {
		t.Fatalf("StorageMount.VolumeTemplate.VolumeTemplateName is not tpl")
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read created state: %v", diags.Errors())
	}

	if state.Id.ValueString() != agsTestInstanceId {
		t.Fatalf("state id is %q, want %q", state.Id.ValueString(), agsTestInstanceId)
	}
	if state.InstanceId.ValueString() != agsTestInstanceId {
		t.Fatalf("state instance_id is %q, want %q", state.InstanceId.ValueString(), agsTestInstanceId)
	}
	if state.ToolId.ValueString() != agsTestToolId {
		t.Fatalf("state tool_id is %q, want %q", state.ToolId.ValueString(), agsTestToolId)
	}
	if state.Status.ValueString() != "RUNNING" {
		t.Fatalf("state status is %q, want RUNNING", state.Status.ValueString())
	}
	if state.Persistent.ValueBool() {
		t.Fatalf("state persistent is true, want false")
	}
	if state.TimeoutSeconds.ValueInt64() != 600 {
		t.Fatalf("state timeout_seconds is %d, want 600", state.TimeoutSeconds.ValueInt64())
	}
	if state.NetworkMode.ValueString() != "PUBLIC" {
		t.Fatalf("state network_mode is %q, want PUBLIC", state.NetworkMode.ValueString())
	}
	if state.AuthMode.ValueString() != "TOKEN" {
		t.Fatalf("state auth_mode is %q, want TOKEN", state.AuthMode.ValueString())
	}
	// Write-only inputs are not reported by the describe API and must be kept.
	if state.Timeout.ValueString() != "10m" {
		t.Fatalf("state timeout is %q, want 10m", state.Timeout.ValueString())
	}
	if state.SourceInstanceId.ValueString() != agsTestSourceInstance {
		t.Fatalf("state source_instance_id is %q, want %q", state.SourceInstanceId.ValueString(), agsTestSourceInstance)
	}
	if state.SnapshotId.ValueString() != agsTestSnapshotId {
		t.Fatalf("state snapshot_id is %q, want %q", state.SnapshotId.ValueString(), agsTestSnapshotId)
	}
	if state.ClientToken.ValueString() != "token-001" {
		t.Fatalf("state client_token is %q, want token-001", state.ClientToken.ValueString())
	}
	if state.StartMode.ValueString() != "DEFAULT" {
		t.Fatalf("state start_mode is %q, want DEFAULT", state.StartMode.ValueString())
	}

	mountOptions := listElements(t, map[string]attr.Value{"mount_options": state.MountOptions}, "mount_options")
	if len(mountOptions) != 1 {
		t.Fatalf("state mount_options has %d elements, want 1", len(mountOptions))
	}
	mountOptionAttributes := attributesOf(t, mountOptions[0])
	assertString(t, mountOptionAttributes, "name", "workspace")
	assertString(t, mountOptionAttributes, "mount_path", "/workspace")
	assertString(t, mountOptionAttributes, "sub_path", "/sub")
	assertBool(t, mountOptionAttributes, "read_only", true)

	metadata := listElements(t, map[string]attr.Value{"metadata": state.Metadata}, "metadata")
	if len(metadata) != 1 {
		t.Fatalf("state metadata has %d elements, want 1", len(metadata))
	}
	metadataAttributes := attributesOf(t, metadata[0])
	assertString(t, metadataAttributes, "name", "run_id")
	assertString(t, metadataAttributes, "value", "42")

	customConfigurationAttributes := attributesOf(t, state.CustomConfiguration)
	assertString(t, customConfigurationAttributes, "image", "img:1")
	assertStringSet(t, customConfigurationAttributes, "command", "/bin/sh")
	assertStringSet(t, customConfigurationAttributes, "args", "-c")

	resources := nestedAttributes(t, customConfigurationAttributes, "resources")
	assertString(t, resources, "cpu", "1")
	assertString(t, resources, "memory", "2Gi")
	assertString(t, resources, "storage", "5Gi")
	assertString(t, resources, "gpu_resource_group_id", "gpusg-1")

	storageMounts := listElements(t, map[string]attr.Value{"storage_mounts": state.StorageMounts}, "storage_mounts")
	if len(storageMounts) != 1 {
		t.Fatalf("state storage_mounts has %d elements, want 1", len(storageMounts))
	}
	storageMountState := attributesOf(t, storageMounts[0])
	assertString(t, storageMountState, "name", "extra")
}

// TestAgsSandboxInstanceCreateGeneratedToken verifies the idempotency token fallback.
func TestAgsSandboxInstanceCreateGeneratedToken(t *testing.T) {
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
				Response: &agsv20250920.StartSandboxInstanceResponseParams{Instance: sampleInstanceItem()},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(sampleInstanceItem()), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.ClientToken = types.StringNull()
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	assertNoErrorDiagnostics(t, createResp.Diagnostics)

	if capturedRequest == nil || capturedRequest.ClientToken == nil || *capturedRequest.ClientToken == "" {
		t.Fatalf("StartSandboxInstance must generate a client token when none is configured")
	}
}

// TestAgsSandboxInstanceCreateEmptyInstanceId verifies an empty create response is an error.
func TestAgsSandboxInstanceCreateEmptyInstanceId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"StartSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StartSandboxInstanceRequest) (*agsv20250920.StartSandboxInstanceResponse, error) {
			return &agsv20250920.StartSandboxInstanceResponse{
				Response: &agsv20250920.StartSandboxInstanceResponseParams{},
			}, nil
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleInstanceModel(t))},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the created instance id is empty")
	}
}

// TestAgsSandboxInstanceCreateFailedStatus verifies a failed start is reported.
func TestAgsSandboxInstanceCreateFailedStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"StartSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StartSandboxInstanceRequest) (*agsv20250920.StartSandboxInstanceResponse, error) {
			return &agsv20250920.StartSandboxInstanceResponse{
				Response: &agsv20250920.StartSandboxInstanceResponseParams{Instance: sampleInstanceItem()},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			item := sampleInstanceItem()
			item.Status = agsString("FAILED")
			return describeSandboxInstancesResponse(item), nil
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleInstanceModel(t))},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the instance is in FAILED status")
	}
}

// TestAgsSandboxInstanceCreateApiError verifies a start API failure is reported.
func TestAgsSandboxInstanceCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"StartSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StartSandboxInstanceRequest) (*agsv20250920.StartSandboxInstanceResponse, error) {
			return nil, &mockSdkError{message: "start failed"}
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleInstanceModel(t))},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the start API fails")
	}
}

// TestAgsSandboxInstanceRead verifies state is refreshed from the describe response.
func TestAgsSandboxInstanceRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(sampleInstanceItem()), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, model)},
		readResp,
	)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags.Errors())
	}

	if state.Id.ValueString() != agsTestInstanceId || state.InstanceId.ValueString() != agsTestInstanceId {
		t.Fatalf("identity was not refreshed")
	}
	if state.Status.ValueString() != "RUNNING" {
		t.Fatalf("status is %q, want RUNNING", state.Status.ValueString())
	}
	if state.ToolId.ValueString() != agsTestToolId {
		t.Fatalf("tool_id is %q, want %q", state.ToolId.ValueString(), agsTestToolId)
	}
	// Write-only inputs keep the value stored in state.
	if state.Timeout.ValueString() != "10m" {
		t.Fatalf("timeout is %q, want 10m", state.Timeout.ValueString())
	}
	if state.StorageMounts.IsNull() {
		t.Fatalf("storage_mounts must keep the configured value")
	}
	if state.StartMode.ValueString() != "DEFAULT" {
		t.Fatalf("start_mode is %q, want DEFAULT", state.StartMode.ValueString())
	}
}

// TestAgsSandboxInstanceReadNotFound verifies a recycled instance is removed from state.
func TestAgsSandboxInstanceReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, model)},
		readResp,
	)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the sandbox instance does not exist")
	}
}

// TestAgsSandboxInstanceReadApiError verifies an API failure never drops state.
func TestAgsSandboxInstanceReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return nil, &mockSdkError{message: "describe failed"}
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, model)},
		readResp,
	)
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the describe API fails")
	}
	if readResp.State.Raw.IsNull() {
		t.Fatalf("state must not be removed when the describe API fails")
	}
}

// TestAgsSandboxInstanceReadEmptyResponse verifies a nil API response is an error.
func TestAgsSandboxInstanceReadEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return &agsv20250920.DescribeSandboxInstancesResponse{}, nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, model)},
		readResp,
	)
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the describe response is empty")
	}
	if readResp.State.Raw.IsNull() {
		t.Fatalf("state must not be removed when the describe response is empty")
	}
}

// TestAgsSandboxInstanceUpdate verifies the update request and the refreshed state.
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
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			item := sampleInstanceItem()
			item.Metadata = []*agsv20250920.MetadataVar{
				{Name: agsString("run_id"), Value: agsString("43")},
			}
			return describeSandboxInstancesResponse(item), nil
		},
	)

	stateModel := testContext.sampleInstanceModel(t)
	stateModel.Id = types.StringValue(agsTestInstanceId)
	stateModel.InstanceId = types.StringValue(agsTestInstanceId)

	planModel := stateModel
	planModel.Timeout = types.StringValue("1h")
	planModel.Metadata = mustList(t, listElementAttributeTypes(t, testContext.schema.Attributes, "metadata"), map[string]attr.Value{
		"name":  types.StringValue("run_id"),
		"value": types.StringValue("43"),
	})

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{
			Plan:  testContext.plan(t, planModel),
			State: testContext.state(t, stateModel),
		},
		updateResp,
	)
	assertNoErrorDiagnostics(t, updateResp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("UpdateSandboxInstance was not called")
	}
	if capturedRequest.InstanceId == nil || *capturedRequest.InstanceId != agsTestInstanceId {
		t.Fatalf("UpdateSandboxInstance InstanceId is not %q", agsTestInstanceId)
	}
	if capturedRequest.Timeout == nil || *capturedRequest.Timeout != "1h" {
		t.Fatalf("UpdateSandboxInstance Timeout is not 1h")
	}
	if len(capturedRequest.Metadata) != 1 || capturedRequest.Metadata[0].Value == nil ||
		*capturedRequest.Metadata[0].Value != "43" {
		t.Fatalf("UpdateSandboxInstance Metadata is not [run_id=43]")
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := updateResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags.Errors())
	}
	if state.Id.ValueString() != agsTestInstanceId || state.InstanceId.ValueString() != agsTestInstanceId {
		t.Fatalf("identity changed during update")
	}
	if state.Timeout.ValueString() != "1h" {
		t.Fatalf("timeout is %q, want 1h", state.Timeout.ValueString())
	}
	// Write-only inputs keep the value stored in state.
	if state.SourceInstanceId.ValueString() != agsTestSourceInstance {
		t.Fatalf("source_instance_id is %q, want %q", state.SourceInstanceId.ValueString(), agsTestSourceInstance)
	}
	if state.StorageMounts.IsNull() {
		t.Fatalf("storage_mounts must keep the configured value")
	}
}

// TestAgsSandboxInstanceUpdateUnchanged verifies an unchanged apply sends nothing.
func TestAgsSandboxInstanceUpdateUnchanged(t *testing.T) {
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
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(sampleInstanceItem()), nil
		},
	)

	stateModel := testContext.sampleInstanceModel(t)
	stateModel.Id = types.StringValue(agsTestInstanceId)
	stateModel.InstanceId = types.StringValue(agsTestInstanceId)

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{
			Plan:  testContext.plan(t, stateModel),
			State: testContext.state(t, stateModel),
		},
		updateResp,
	)
	assertNoErrorDiagnostics(t, updateResp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("UpdateSandboxInstance was not called")
	}
	if capturedRequest.Timeout != nil {
		t.Fatalf("UpdateSandboxInstance Timeout must not be sent when unchanged")
	}
	if capturedRequest.Metadata != nil {
		t.Fatalf("UpdateSandboxInstance Metadata must not be sent when unchanged")
	}
}

// TestAgsSandboxInstanceUpdateApiError verifies an update API failure is reported.
func TestAgsSandboxInstanceUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"UpdateSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.UpdateSandboxInstanceRequest) (*agsv20250920.UpdateSandboxInstanceResponse, error) {
			return nil, &mockSdkError{message: "update failed"}
		},
	)

	stateModel := testContext.sampleInstanceModel(t)
	stateModel.Id = types.StringValue(agsTestInstanceId)
	stateModel.InstanceId = types.StringValue(agsTestInstanceId)

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{
			Plan:  testContext.plan(t, stateModel),
			State: testContext.state(t, stateModel),
		},
		updateResp,
	)
	if !updateResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the update API fails")
	}
}

// TestAgsSandboxInstanceDelete verifies the stop API and the follow up wait.
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
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			item := sampleInstanceItem()
			item.Status = agsString("STOPPED")
			return describeSandboxInstancesResponse(item), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, model)},
		deleteResp,
	)
	assertNoErrorDiagnostics(t, deleteResp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("StopSandboxInstance was not called")
	}
	if capturedRequest.InstanceId == nil || *capturedRequest.InstanceId != agsTestInstanceId {
		t.Fatalf("StopSandboxInstance InstanceId is not %q", agsTestInstanceId)
	}
}

// TestAgsSandboxInstanceDeleteRecycled verifies a recycled instance stops cleanly.
func TestAgsSandboxInstanceDeleteRecycled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"StopSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StopSandboxInstanceRequest) (*agsv20250920.StopSandboxInstanceResponse, error) {
			return &agsv20250920.StopSandboxInstanceResponse{
				Response: &agsv20250920.StopSandboxInstanceResponseParams{},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, model)},
		deleteResp,
	)
	assertNoErrorDiagnostics(t, deleteResp.Diagnostics)
}

// TestAgsSandboxInstanceDeleteAlreadyGone verifies a repeated delete is idempotent.
func TestAgsSandboxInstanceDeleteAlreadyGone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"StopSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StopSandboxInstanceRequest) (*agsv20250920.StopSandboxInstanceResponse, error) {
			return nil, &mockSdkError{message: "instance does not exist"}
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, model)},
		deleteResp,
	)
	assertNoErrorDiagnostics(t, deleteResp.Diagnostics)
}

// TestAgsSandboxInstanceDeleteApiError verifies a failed delete keeps the resource.
func TestAgsSandboxInstanceDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"StopSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StopSandboxInstanceRequest) (*agsv20250920.StopSandboxInstanceResponse, error) {
			return nil, &mockSdkError{message: "stop failed"}
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(sampleInstanceItem()), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, model)},
		deleteResp,
	)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the stop API fails")
	}
}

// TestAgsSandboxInstanceDeleteStopFailed verifies a stop failure surfaces a warning.
func TestAgsSandboxInstanceDeleteStopFailed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"StopSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StopSandboxInstanceRequest) (*agsv20250920.StopSandboxInstanceResponse, error) {
			return &agsv20250920.StopSandboxInstanceResponse{
				Response: &agsv20250920.StopSandboxInstanceResponseParams{},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			item := sampleInstanceItem()
			item.Status = agsString("STOP_FAILED")
			return describeSandboxInstancesResponse(item), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, model)},
		deleteResp,
	)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the instance ends in STOP_FAILED status")
	}
}

// TestAgsSandboxInstanceImportState verifies the id only import path.
func TestAgsSandboxInstanceImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	importer, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: agsTestInstanceId},
		importResp,
	)
	assertNoErrorDiagnostics(t, importResp.Diagnostics)

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := importResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read imported state: %v", diags.Errors())
	}
	if state.Id.ValueString() != agsTestInstanceId {
		t.Fatalf("imported id is %q, want %q", state.Id.ValueString(), agsTestInstanceId)
	}
	if state.InstanceId.ValueString() != agsTestInstanceId {
		t.Fatalf("imported instance_id is %q, want %q", state.InstanceId.ValueString(), agsTestInstanceId)
	}
}

// TestAgsSandboxInstanceImportStateEmptyId verifies an empty import id is rejected.
func TestAgsSandboxInstanceImportStateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	importer := testContext.resource.(resource.ResourceWithImportState)
	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: ""},
		importResp,
	)
	if !importResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for an empty import id")
	}
}

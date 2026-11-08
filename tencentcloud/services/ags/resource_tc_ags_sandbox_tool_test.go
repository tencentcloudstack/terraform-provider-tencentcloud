package ags_test

import (
	"context"
	"sort"
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
	agsTestToolId     = "tool-001"
	agsTestToolName   = "tf-example"
	agsTestRoleArn    = "qcs::cam::uin/123:role/test"
	agsTestCreateTime = "2026-01-01T00:00:00Z"
	agsTestUpdateTime = "2026-01-02T00:00:00Z"
)

func agsString(value string) *string { return &value }

func agsBool(value bool) *bool { return &value }

func agsInt64(value int64) *int64 { return &value }

// agsSandboxToolTestContext wires a resource instance to a patched SDK client.
type agsSandboxToolTestContext struct {
	resource resource.Resource
	client   *agsv20250920.Client
	schema   schema.Schema
}

func newAgsSandboxToolTestContext(t *testing.T, patches *gomonkey.Patches) *agsSandboxToolTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseAgsClient", agsClient)

	res := svcags.NewAgsSandboxToolResource()
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

	return &agsSandboxToolTestContext{resource: res, client: agsClient, schema: schemaResp.Schema}
}

func (c *agsSandboxToolTestContext) plan(t *testing.T, model svcags.AgsSandboxToolResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *agsSandboxToolTestContext) state(t *testing.T, model svcags.AgsSandboxToolResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func objectAttributeTypes(t *testing.T, attributes map[string]schema.Attribute, name string) map[string]attr.Type {
	t.Helper()

	nested, ok := attributes[name].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a single nested attribute", name)
	}
	objectType, ok := nested.GetType().(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object type", name)
	}
	return objectType.AttrTypes
}

func listElementAttributeTypes(t *testing.T, attributes map[string]schema.Attribute, name string) map[string]attr.Type {
	t.Helper()

	nested, ok := attributes[name].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a list nested attribute", name)
	}
	listType, ok := nested.GetType().(types.ListType)
	if !ok {
		t.Fatalf("attribute %q does not have a list type", name)
	}
	objectType, ok := listType.ElemType.(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object element type", name)
	}
	return objectType.AttrTypes
}

func mustObject(t *testing.T, attributeTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	t.Helper()

	object, diags := types.ObjectValue(attributeTypes, values)
	if diags.HasError() {
		t.Fatalf("failed to build object value: %v", diags.Errors())
	}
	return object
}

func mustList(t *testing.T, elementTypes map[string]attr.Type, elements ...map[string]attr.Value) types.List {
	t.Helper()

	objectType := types.ObjectType{AttrTypes: elementTypes}
	values := make([]attr.Value, 0, len(elements))
	for _, element := range elements {
		object, diags := types.ObjectValue(elementTypes, element)
		if diags.HasError() {
			t.Fatalf("failed to build list element: %v", diags.Errors())
		}
		values = append(values, object)
	}

	list, diags := types.ListValue(objectType, values)
	if diags.HasError() {
		t.Fatalf("failed to build list value: %v", diags.Errors())
	}
	return list
}

func mustStringSet(t *testing.T, values ...string) types.Set {
	t.Helper()

	set, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		t.Fatalf("failed to build string set: %v", diags.Errors())
	}
	return set
}

// sampleModel returns a plan model with every configurable leaf populated.
func (c *agsSandboxToolTestContext) sampleModel(t *testing.T) svcags.AgsSandboxToolResourceModel {
	t.Helper()

	top := c.schema.Attributes

	networkAttributes := top["network_configuration"].(schema.SingleNestedAttribute).Attributes
	vpcConfigTypes := objectAttributeTypes(t, networkAttributes, "vpc_config")
	networkConfigurationTypes := objectAttributeTypes(t, top, "network_configuration")

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
	tagTypes := listElementAttributeTypes(t, top, "tags")
	envTypes := listElementAttributeTypes(t, customAttributes, "env")
	portTypes := listElementAttributeTypes(t, customAttributes, "ports")

	computerAttributes := top["computer_configuration"].(schema.SingleNestedAttribute).Attributes
	waaTypes := objectAttributeTypes(t, computerAttributes, "waa_configuration")
	osWorldTypes := objectAttributeTypes(t, computerAttributes, "os_world_configuration")
	computerConfigurationTypes := objectAttributeTypes(t, top, "computer_configuration")

	logAttributes := top["log_configuration"].(schema.SingleNestedAttribute).Attributes
	clsConfigTypes := objectAttributeTypes(t, logAttributes, "cls_config")
	logSourcesTypes := objectAttributeTypes(t, logAttributes, "log_sources")
	logConfigurationTypes := objectAttributeTypes(t, top, "log_configuration")

	return svcags.AgsSandboxToolResourceModel{
		Id:       types.StringNull(),
		ToolId:   types.StringNull(),
		ToolName: types.StringValue(agsTestToolName),
		ToolType: types.StringValue("custom"),
		NetworkConfiguration: mustObject(t, networkConfigurationTypes, map[string]attr.Value{
			"network_mode": types.StringValue("VPC"),
			"vpc_config": mustObject(t, vpcConfigTypes, map[string]attr.Value{
				"subnet_ids":         mustStringSet(t, "subnet-1"),
				"security_group_ids": mustStringSet(t, "sg-1"),
			}),
		}),
		Description:    types.StringValue("desc"),
		DefaultTimeout: types.StringValue("10m"),
		Tags: mustList(t, tagTypes, map[string]attr.Value{
			"key":   types.StringValue("k1"),
			"value": types.StringValue("v1"),
		}),
		ClientToken: types.StringValue("token-1"),
		RoleArn:     types.StringValue(agsTestRoleArn),
		StorageMounts: mustList(t, storageMountTypes, map[string]attr.Value{
			"name": types.StringValue("workspace"),
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
			"mount_path": types.StringValue("/workspace"),
			"read_only":  types.BoolValue(true),
			"volume": mustObject(t, volumeTypes, map[string]attr.Value{
				"volume_name": types.StringValue("vol-name"),
				"volume_id":   types.StringValue("vol-12345678"),
			}),
			"volume_template": mustObject(t, volumeTemplateTypes, map[string]attr.Value{
				"volume_template_name": types.StringValue("tpl"),
				"reuse_key":            types.StringValue("${metadata.run_id}"),
			}),
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
				"cpu":     types.StringValue("1"),
				"memory":  types.StringValue("2Gi"),
				"storage": types.StringValue("5Gi"),
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
		ComputerConfiguration: mustObject(t, computerConfigurationTypes, map[string]attr.Value{
			"waa_configuration": mustObject(t, waaTypes, map[string]attr.Value{
				"image_id": types.StringValue("img-waa"),
			}),
			"os_world_configuration": mustObject(t, osWorldTypes, map[string]attr.Value{
				"version": types.StringValue("osworld1"),
			}),
		}),
		LogConfiguration: mustObject(t, logConfigurationTypes, map[string]attr.Value{
			"cls_config": mustObject(t, clsConfigTypes, map[string]attr.Value{
				"topic_id": types.StringValue("topic-1"),
			}),
			"log_sources": mustObject(t, logSourcesTypes, map[string]attr.Value{
				"files": mustStringSet(t, "/logs/app.log"),
			}),
		}),
		Persistent:   types.BoolValue(true),
		Status:       types.StringNull(),
		StatusReason: types.StringNull(),
		CreateTime:   types.StringNull(),
		UpdateTime:   types.StringNull(),
	}
}

// sampleItem returns the cloud side representation used by the describe mock.
func sampleItem() *agsv20250920.SandboxTool {
	return &agsv20250920.SandboxTool{
		ToolId:                agsString(agsTestToolId),
		ToolName:              agsString(agsTestToolName),
		ToolType:              agsString("custom"),
		Status:                agsString("ACTIVE"),
		Description:           agsString("desc"),
		Persistent:            agsBool(true),
		DefaultTimeoutSeconds: helperUint64(600),
		NetworkConfiguration: &agsv20250920.NetworkConfiguration{
			NetworkMode: agsString("VPC"),
			VpcConfig: &agsv20250920.VPCConfig{
				SubnetIds:        []*string{agsString("subnet-1")},
				SecurityGroupIds: []*string{agsString("sg-1")},
			},
		},
		Tags: []*agsv20250920.Tag{
			{Key: agsString("k1"), Value: agsString("v1")},
		},
		CreateTime: agsString(agsTestCreateTime),
		UpdateTime: agsString(agsTestUpdateTime),
		RoleArn:    agsString(agsTestRoleArn),
		StorageMounts: []*agsv20250920.StorageMount{
			{
				Name: agsString("workspace"),
				StorageSource: &agsv20250920.StorageSource{
					Cos: &agsv20250920.CosStorageSource{
						Endpoint:   agsString("cos.example.com"),
						BucketName: agsString("bucket"),
						BucketPath: agsString("/path"),
					},
					Image: &agsv20250920.ImageStorageSource{
						Reference:         agsString("img:1"),
						ImageRegistryType: agsString("personal"),
						SubPath:           agsString("/sub"),
						Digest:            agsString("sha256:abc"),
					},
					Cfs: &agsv20250920.CfsStorageSource{
						FileSystemId: agsString("cfs-1"),
						Path:         agsString("/cfs"),
					},
					AgentBucket: &agsv20250920.AgentBucketStorageSource{
						LibraryId:    agsString("lib-1"),
						SpaceId:      agsString("space-1"),
						AccessDomain: agsString("domain.example.com"),
					},
				},
				MountPath: agsString("/workspace"),
				ReadOnly:  agsBool(true),
				Volume: &agsv20250920.VolumeRef{
					VolumeName: agsString("vol-name"),
					VolumeId:   agsString("vol-12345678"),
				},
				VolumeTemplate: &agsv20250920.VolumeTemplateRef{
					VolumeTemplateName: agsString("tpl"),
					ReuseKey:           agsString("${metadata.run_id}"),
				},
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
				CPU:     agsString("1"),
				Memory:  agsString("2Gi"),
				Storage: agsString("5Gi"),
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
		LogConfiguration: &agsv20250920.LogConfiguration{
			CLSConfig: &agsv20250920.CLSConfig{TopicId: agsString("topic-1")},
			LogSources: &agsv20250920.LogSources{
				Files: []*string{agsString("/logs/app.log")},
			},
		},
		ComputerConfiguration: &agsv20250920.ComputerConfiguration{
			WAAConfiguration:     &agsv20250920.WAAConfiguration{ImageId: agsString("img-waa")},
			OSWorldConfiguration: &agsv20250920.OSWorldConfiguration{Version: agsString("osworld1")},
		},
		StatusReason: agsString("running"),
	}
}

func helperUint64(value uint64) *uint64 { return &value }

func describeSandboxToolResponse(tools ...*agsv20250920.SandboxTool) *agsv20250920.DescribeSandboxToolListResponse {
	return &agsv20250920.DescribeSandboxToolListResponse{
		Response: &agsv20250920.DescribeSandboxToolListResponseParams{
			SandboxToolSet: tools,
			TotalCount:     agsInt64(int64(len(tools))),
		},
	}
}

func attributesOf(t *testing.T, value attr.Value) map[string]attr.Value {
	t.Helper()

	object, ok := value.(types.Object)
	if !ok {
		t.Fatalf("value is not an object: %T", value)
	}
	return object.Attributes()
}

func nestedAttributes(t *testing.T, attributes map[string]attr.Value, name string) map[string]attr.Value {
	t.Helper()

	value, ok := attributes[name]
	if !ok {
		t.Fatalf("attribute %q is missing", name)
	}
	return attributesOf(t, value)
}

func listElements(t *testing.T, attributes map[string]attr.Value, name string) []attr.Value {
	t.Helper()

	value, ok := attributes[name]
	if !ok {
		t.Fatalf("attribute %q is missing", name)
	}
	list, ok := value.(types.List)
	if !ok {
		t.Fatalf("attribute %q is not a list: %T", name, value)
	}
	if list.IsNull() || list.IsUnknown() {
		t.Fatalf("attribute %q has no value", name)
	}
	return list.Elements()
}

func setStrings(t *testing.T, value attr.Value) []string {
	t.Helper()

	set, ok := value.(types.Set)
	if !ok {
		t.Fatalf("value is not a set: %T", value)
	}
	if set.IsNull() || set.IsUnknown() {
		t.Fatalf("set has no value")
	}

	var result []string
	if diags := set.ElementsAs(context.Background(), &result, false); diags.HasError() {
		t.Fatalf("failed to read set elements: %v", diags.Errors())
	}
	sort.Strings(result)
	return result
}

func assertString(t *testing.T, attributes map[string]attr.Value, name, expected string) {
	t.Helper()

	value, ok := attributes[name].(types.String)
	if !ok {
		t.Fatalf("attribute %q is not a string", name)
	}
	if value.IsNull() || value.IsUnknown() {
		t.Fatalf("attribute %q has no value", name)
	}
	if value.ValueString() != expected {
		t.Fatalf("attribute %q is %q, want %q", name, value.ValueString(), expected)
	}
}

func assertBool(t *testing.T, attributes map[string]attr.Value, name string, expected bool) {
	t.Helper()

	value, ok := attributes[name].(types.Bool)
	if !ok {
		t.Fatalf("attribute %q is not a bool", name)
	}
	if value.IsNull() || value.IsUnknown() {
		t.Fatalf("attribute %q has no value", name)
	}
	if value.ValueBool() != expected {
		t.Fatalf("attribute %q is %v, want %v", name, value.ValueBool(), expected)
	}
}

func assertInt64(t *testing.T, attributes map[string]attr.Value, name string, expected int64) {
	t.Helper()

	value, ok := attributes[name].(types.Int64)
	if !ok {
		t.Fatalf("attribute %q is not an int64", name)
	}
	if value.IsNull() || value.IsUnknown() {
		t.Fatalf("attribute %q has no value", name)
	}
	if value.ValueInt64() != expected {
		t.Fatalf("attribute %q is %d, want %d", name, value.ValueInt64(), expected)
	}
}

func assertStringSet(t *testing.T, attributes map[string]attr.Value, name string, expected ...string) {
	t.Helper()

	value, ok := attributes[name]
	if !ok {
		t.Fatalf("attribute %q is missing", name)
	}
	actual := setStrings(t, value)

	sort.Strings(expected)
	if len(actual) != len(expected) {
		t.Fatalf("attribute %q has %d elements, want %d", name, len(actual), len(expected))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("attribute %q element %d is %q, want %q", name, i, actual[i], expected[i])
		}
	}
}

func assertStringPointerList(t *testing.T, name string, actual []*string, expected ...string) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("%s has %d elements, want %d", name, len(actual), len(expected))
	}
	for i := range expected {
		if actual[i] == nil {
			t.Fatalf("%s element %d is nil", name, i)
		}
		if *actual[i] != expected[i] {
			t.Fatalf("%s element %d is %q, want %q", name, i, *actual[i], expected[i])
		}
	}
}

func assertNoErrorDiagnostics(t *testing.T, diags interface{ HasError() bool }) {
	t.Helper()

	if diags.HasError() {
		t.Fatalf("unexpected error diagnostics")
	}
}

// TestAgsSandboxToolSchema verifies the attribute set and its required flags.
func TestAgsSandboxToolSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)
	top := testContext.schema.Attributes

	if len(top) != 19 {
		t.Fatalf("schema has %d attributes, want 19", len(top))
	}

	required := []string{"tool_name", "tool_type", "network_configuration"}
	for _, name := range required {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsRequired() {
			t.Fatalf("attribute %q must be required", name)
		}
	}

	computed := []string{"id", "tool_id", "status", "status_reason", "create_time", "update_time"}
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

	network := top["network_configuration"].(schema.SingleNestedAttribute)
	if !network.Attributes["network_mode"].IsRequired() {
		t.Fatalf("network_configuration.network_mode must be required")
	}

	custom := top["custom_configuration"].(schema.SingleNestedAttribute)
	dnsConfig := custom.Attributes["dns_config"].(schema.SingleNestedAttribute)
	if !dnsConfig.Attributes["servers"].IsRequired() {
		t.Fatalf("custom_configuration.dns_config.servers must be required")
	}

	storageMount := top["storage_mounts"].(schema.ListNestedAttribute)
	volumeTemplate := storageMount.NestedObject.Attributes["volume_template"].(schema.SingleNestedAttribute)
	if !volumeTemplate.Attributes["volume_template_name"].IsRequired() {
		t.Fatalf("storage_mounts.volume_template.volume_template_name must be required")
	}
}

// TestAgsSandboxToolCreate verifies the create request and the read back state.
func TestAgsSandboxToolCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	var capturedRequest *agsv20250920.CreateSandboxToolRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"CreateSandboxToolWithContext",
		func(_ context.Context, request *agsv20250920.CreateSandboxToolRequest) (*agsv20250920.CreateSandboxToolResponse, error) {
			capturedRequest = request
			return &agsv20250920.CreateSandboxToolResponse{
				Response: &agsv20250920.CreateSandboxToolResponseParams{ToolId: agsString(agsTestToolId)},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			assertStringPointerList(t, "describe request ToolIds", request.ToolIds, agsTestToolId)
			return describeSandboxToolResponse(sampleItem()), nil
		},
	)

	model := testContext.sampleModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	assertNoErrorDiagnostics(t, createResp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("CreateSandboxTool was not called")
	}
	if capturedRequest.ToolName == nil || *capturedRequest.ToolName != agsTestToolName {
		t.Fatalf("CreateSandboxTool ToolName is not %q", agsTestToolName)
	}
	if capturedRequest.ToolType == nil || *capturedRequest.ToolType != "custom" {
		t.Fatalf("CreateSandboxTool ToolType is not custom")
	}
	if capturedRequest.Description == nil || *capturedRequest.Description != "desc" {
		t.Fatalf("CreateSandboxTool Description is not desc")
	}
	if capturedRequest.DefaultTimeout == nil || *capturedRequest.DefaultTimeout != "10m" {
		t.Fatalf("CreateSandboxTool DefaultTimeout is not 10m")
	}
	if capturedRequest.ClientToken == nil || *capturedRequest.ClientToken != "token-1" {
		t.Fatalf("CreateSandboxTool ClientToken is not token-1")
	}
	if capturedRequest.RoleArn == nil || *capturedRequest.RoleArn != agsTestRoleArn {
		t.Fatalf("CreateSandboxTool RoleArn is not %q", agsTestRoleArn)
	}
	if capturedRequest.Persistent == nil || !*capturedRequest.Persistent {
		t.Fatalf("CreateSandboxTool Persistent is not true")
	}
	if capturedRequest.NetworkConfiguration == nil || capturedRequest.NetworkConfiguration.NetworkMode == nil ||
		*capturedRequest.NetworkConfiguration.NetworkMode != "VPC" {
		t.Fatalf("CreateSandboxTool network mode is not VPC")
	}
	if capturedRequest.NetworkConfiguration.VpcConfig == nil {
		t.Fatalf("CreateSandboxTool VpcConfig is nil")
	}
	assertStringPointerList(t, "VpcConfig.SubnetIds", capturedRequest.NetworkConfiguration.VpcConfig.SubnetIds, "subnet-1")
	assertStringPointerList(t, "VpcConfig.SecurityGroupIds", capturedRequest.NetworkConfiguration.VpcConfig.SecurityGroupIds, "sg-1")

	if len(capturedRequest.Tags) != 1 || capturedRequest.Tags[0].Key == nil || *capturedRequest.Tags[0].Key != "k1" ||
		capturedRequest.Tags[0].Value == nil || *capturedRequest.Tags[0].Value != "v1" {
		t.Fatalf("CreateSandboxTool Tags are not [k1=v1]")
	}

	if len(capturedRequest.StorageMounts) != 1 {
		t.Fatalf("CreateSandboxTool StorageMounts has %d elements, want 1", len(capturedRequest.StorageMounts))
	}
	storageMount := capturedRequest.StorageMounts[0]
	if storageMount.Name == nil || *storageMount.Name != "workspace" {
		t.Fatalf("StorageMount.Name is not workspace")
	}
	if storageMount.MountPath == nil || *storageMount.MountPath != "/workspace" {
		t.Fatalf("StorageMount.MountPath is not /workspace")
	}
	if storageMount.ReadOnly == nil || !*storageMount.ReadOnly {
		t.Fatalf("StorageMount.ReadOnly is not true")
	}
	if storageMount.StorageSource == nil || storageMount.StorageSource.Cos == nil ||
		storageMount.StorageSource.Cos.BucketName == nil || *storageMount.StorageSource.Cos.BucketName != "bucket" {
		t.Fatalf("StorageMount.StorageSource.Cos.BucketName is not bucket")
	}
	if storageMount.StorageSource.Image == nil || storageMount.StorageSource.Image.Digest != nil {
		t.Fatalf("StorageMount.StorageSource.Image.Digest must not be sent")
	}
	if storageMount.StorageSource.Cfs == nil || storageMount.StorageSource.Cfs.FileSystemId == nil ||
		*storageMount.StorageSource.Cfs.FileSystemId != "cfs-1" {
		t.Fatalf("StorageMount.StorageSource.Cfs.FileSystemId is not cfs-1")
	}
	if storageMount.StorageSource.AgentBucket == nil || storageMount.StorageSource.AgentBucket.LibraryId == nil ||
		*storageMount.StorageSource.AgentBucket.LibraryId != "lib-1" {
		t.Fatalf("StorageMount.StorageSource.AgentBucket.LibraryId is not lib-1")
	}
	if storageMount.Volume == nil || storageMount.Volume.VolumeName == nil || *storageMount.Volume.VolumeName != "vol-name" {
		t.Fatalf("StorageMount.Volume.VolumeName is not vol-name")
	}
	if storageMount.VolumeTemplate == nil || storageMount.VolumeTemplate.VolumeTemplateName == nil ||
		*storageMount.VolumeTemplate.VolumeTemplateName != "tpl" {
		t.Fatalf("StorageMount.VolumeTemplate.VolumeTemplateName is not tpl")
	}

	customConfiguration := capturedRequest.CustomConfiguration
	if customConfiguration == nil {
		t.Fatalf("CreateSandboxTool CustomConfiguration is nil")
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
	if customConfiguration.Resources == nil || customConfiguration.Resources.Storage == nil ||
		*customConfiguration.Resources.Storage != "5Gi" {
		t.Fatalf("CustomConfiguration.Resources.Storage is not 5Gi")
	}
	if customConfiguration.Probe == nil || customConfiguration.Probe.HttpGet == nil ||
		customConfiguration.Probe.HttpGet.Path == nil || *customConfiguration.Probe.HttpGet.Path != "/healthz" {
		t.Fatalf("CustomConfiguration.Probe.HttpGet.Path is not /healthz")
	}
	if customConfiguration.Probe.FailureThreshold == nil || *customConfiguration.Probe.FailureThreshold != 3 {
		t.Fatalf("CustomConfiguration.Probe.FailureThreshold is not 3")
	}
	if customConfiguration.DNSConfig == nil {
		t.Fatalf("CustomConfiguration.DNSConfig is nil")
	}
	assertStringPointerList(t, "DNSConfig.Servers", customConfiguration.DNSConfig.Servers, "10.0.0.1")
	if customConfiguration.SecurityContext == nil || customConfiguration.SecurityContext.Privileged == nil ||
		!*customConfiguration.SecurityContext.Privileged {
		t.Fatalf("CustomConfiguration.SecurityContext.Privileged is not true")
	}
	if customConfiguration.SecurityContext.Capabilities == nil {
		t.Fatalf("CustomConfiguration.SecurityContext.Capabilities is nil")
	}
	assertStringPointerList(t, "Capabilities.Add", customConfiguration.SecurityContext.Capabilities.Add, "SYS_ADMIN")

	if capturedRequest.ComputerConfiguration == nil || capturedRequest.ComputerConfiguration.WAAConfiguration == nil ||
		capturedRequest.ComputerConfiguration.WAAConfiguration.ImageId == nil ||
		*capturedRequest.ComputerConfiguration.WAAConfiguration.ImageId != "img-waa" {
		t.Fatalf("ComputerConfiguration.WAAConfiguration.ImageId is not img-waa")
	}
	if capturedRequest.ComputerConfiguration.OSWorldConfiguration == nil ||
		capturedRequest.ComputerConfiguration.OSWorldConfiguration.Version == nil ||
		*capturedRequest.ComputerConfiguration.OSWorldConfiguration.Version != "osworld1" {
		t.Fatalf("ComputerConfiguration.OSWorldConfiguration.Version is not osworld1")
	}

	if capturedRequest.LogConfiguration == nil || capturedRequest.LogConfiguration.CLSConfig == nil ||
		capturedRequest.LogConfiguration.CLSConfig.TopicId == nil ||
		*capturedRequest.LogConfiguration.CLSConfig.TopicId != "topic-1" {
		t.Fatalf("LogConfiguration.CLSConfig.TopicId is not topic-1")
	}
	if capturedRequest.LogConfiguration.LogSources == nil {
		t.Fatalf("LogConfiguration.LogSources is nil")
	}
	assertStringPointerList(t, "LogSources.Files", capturedRequest.LogConfiguration.LogSources.Files, "/logs/app.log")

	var state svcags.AgsSandboxToolResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read created state: %v", diags.Errors())
	}

	if state.Id.ValueString() != agsTestToolId {
		t.Fatalf("state id is %q, want %q", state.Id.ValueString(), agsTestToolId)
	}
	if state.ToolId.ValueString() != agsTestToolId {
		t.Fatalf("state tool_id is %q, want %q", state.ToolId.ValueString(), agsTestToolId)
	}
	if state.ToolName.ValueString() != agsTestToolName {
		t.Fatalf("state tool_name is %q, want %q", state.ToolName.ValueString(), agsTestToolName)
	}
	if state.ToolType.ValueString() != "custom" {
		t.Fatalf("state tool_type is %q, want custom", state.ToolType.ValueString())
	}
	if state.Status.ValueString() != "ACTIVE" {
		t.Fatalf("state status is %q, want ACTIVE", state.Status.ValueString())
	}
	if state.StatusReason.ValueString() != "running" {
		t.Fatalf("state status_reason is %q, want running", state.StatusReason.ValueString())
	}
	if state.CreateTime.ValueString() != agsTestCreateTime {
		t.Fatalf("state create_time is %q", state.CreateTime.ValueString())
	}
	if state.UpdateTime.ValueString() != agsTestUpdateTime {
		t.Fatalf("state update_time is %q", state.UpdateTime.ValueString())
	}
	// create only inputs are not reported by the describe API and must be kept.
	if state.DefaultTimeout.ValueString() != "10m" {
		t.Fatalf("state default_timeout is %q, want 10m", state.DefaultTimeout.ValueString())
	}
	if state.ClientToken.ValueString() != "token-1" {
		t.Fatalf("state client_token is %q, want token-1", state.ClientToken.ValueString())
	}

	top := map[string]attr.Value{
		"id":                      state.Id,
		"tool_id":                 state.ToolId,
		"tool_name":               state.ToolName,
		"tool_type":               state.ToolType,
		"network_configuration":   state.NetworkConfiguration,
		"description":             state.Description,
		"default_timeout":         state.DefaultTimeout,
		"tags":                    state.Tags,
		"client_token":            state.ClientToken,
		"role_arn":                state.RoleArn,
		"storage_mounts":          state.StorageMounts,
		"custom_configuration":    state.CustomConfiguration,
		"computer_configuration":  state.ComputerConfiguration,
		"log_configuration":       state.LogConfiguration,
		"persistent":              state.Persistent,
		"status":                  state.Status,
		"status_reason":           state.StatusReason,
		"create_time":             state.CreateTime,
		"update_time":             state.UpdateTime,
	}

	assertString(t, top, "description", "desc")
	assertString(t, top, "role_arn", agsTestRoleArn)
	assertBool(t, top, "persistent", true)

	networkConfiguration := attributesOf(t, top["network_configuration"])
	assertString(t, networkConfiguration, "network_mode", "VPC")
	vpcConfig := nestedAttributes(t, networkConfiguration, "vpc_config")
	assertStringSet(t, vpcConfig, "subnet_ids", "subnet-1")
	assertStringSet(t, vpcConfig, "security_group_ids", "sg-1")

	tags := listElements(t, top, "tags")
	if len(tags) != 1 {
		t.Fatalf("state tags has %d elements, want 1", len(tags))
	}
	tag := attributesOf(t, tags[0])
	assertString(t, tag, "key", "k1")
	assertString(t, tag, "value", "v1")

	storageMounts := listElements(t, top, "storage_mounts")
	if len(storageMounts) != 1 {
		t.Fatalf("state storage_mounts has %d elements, want 1", len(storageMounts))
	}
	storageMountAttributes := attributesOf(t, storageMounts[0])
	assertString(t, storageMountAttributes, "name", "workspace")
	assertString(t, storageMountAttributes, "mount_path", "/workspace")
	assertBool(t, storageMountAttributes, "read_only", true)

	storageSource := nestedAttributes(t, storageMountAttributes, "storage_source")
	cos := nestedAttributes(t, storageSource, "cos")
	assertString(t, cos, "endpoint", "cos.example.com")
	assertString(t, cos, "bucket_name", "bucket")
	assertString(t, cos, "bucket_path", "/path")
	image := nestedAttributes(t, storageSource, "image")
	assertString(t, image, "reference", "img:1")
	assertString(t, image, "sub_path", "/sub")
	assertString(t, image, "digest", "sha256:abc")
	cfs := nestedAttributes(t, storageSource, "cfs")
	assertString(t, cfs, "file_system_id", "cfs-1")
	agentBucket := nestedAttributes(t, storageSource, "agent_bucket")
	assertString(t, agentBucket, "library_id", "lib-1")
	assertString(t, agentBucket, "space_id", "space-1")
	assertString(t, agentBucket, "access_domain", "domain.example.com")

	volume := nestedAttributes(t, storageMountAttributes, "volume")
	assertString(t, volume, "volume_name", "vol-name")
	assertString(t, volume, "volume_id", "vol-12345678")
	volumeTemplate := nestedAttributes(t, storageMountAttributes, "volume_template")
	assertString(t, volumeTemplate, "volume_template_name", "tpl")
	assertString(t, volumeTemplate, "reuse_key", "${metadata.run_id}")

	customConfigurationAttributes := attributesOf(t, top["custom_configuration"])
	assertString(t, customConfigurationAttributes, "image", "img:1")
	assertStringSet(t, customConfigurationAttributes, "command", "/bin/sh")
	assertStringSet(t, customConfigurationAttributes, "args", "-c")

	env := listElements(t, customConfigurationAttributes, "env")
	if len(env) != 1 {
		t.Fatalf("state custom_configuration.env has %d elements, want 1", len(env))
	}
	assertString(t, attributesOf(t, env[0]), "name", "K")

	ports := listElements(t, customConfigurationAttributes, "ports")
	if len(ports) != 1 {
		t.Fatalf("state custom_configuration.ports has %d elements, want 1", len(ports))
	}
	port := attributesOf(t, ports[0])
	assertString(t, port, "name", "http")
	assertInt64(t, port, "port", 8080)
	assertString(t, port, "protocol", "TCP")

	resources := nestedAttributes(t, customConfigurationAttributes, "resources")
	assertString(t, resources, "cpu", "1")
	assertString(t, resources, "memory", "2Gi")
	assertString(t, resources, "storage", "5Gi")

	probe := nestedAttributes(t, customConfigurationAttributes, "probe")
	httpGet := nestedAttributes(t, probe, "http_get")
	assertString(t, httpGet, "path", "/healthz")
	assertInt64(t, httpGet, "port", 8080)
	assertString(t, httpGet, "scheme", "HTTP")
	assertInt64(t, probe, "ready_timeout_ms", 1000)
	assertInt64(t, probe, "probe_timeout_ms", 2000)
	assertInt64(t, probe, "probe_period_ms", 3000)
	assertInt64(t, probe, "success_threshold", 1)
	assertInt64(t, probe, "failure_threshold", 3)

	dnsConfig := nestedAttributes(t, customConfigurationAttributes, "dns_config")
	assertStringSet(t, dnsConfig, "servers", "10.0.0.1")
	assertStringSet(t, dnsConfig, "searches", "example.com")
	assertStringSet(t, dnsConfig, "options", "ndots:2")

	securityContext := nestedAttributes(t, customConfigurationAttributes, "security_context")
	assertBool(t, securityContext, "privileged", true)
	capabilities := nestedAttributes(t, securityContext, "capabilities")
	assertStringSet(t, capabilities, "add", "SYS_ADMIN")
	assertStringSet(t, capabilities, "drop", "ALL")

	computerConfiguration := attributesOf(t, top["computer_configuration"])
	waaConfiguration := nestedAttributes(t, computerConfiguration, "waa_configuration")
	assertString(t, waaConfiguration, "image_id", "img-waa")
	osWorldConfiguration := nestedAttributes(t, computerConfiguration, "os_world_configuration")
	assertString(t, osWorldConfiguration, "version", "osworld1")

	logConfiguration := attributesOf(t, top["log_configuration"])
	clsConfig := nestedAttributes(t, logConfiguration, "cls_config")
	assertString(t, clsConfig, "topic_id", "topic-1")
	logSources := nestedAttributes(t, logConfiguration, "log_sources")
	assertStringSet(t, logSources, "files", "/logs/app.log")
}

// TestAgsSandboxToolCreateEmptyToolId verifies an empty create response is an error.
func TestAgsSandboxToolCreateEmptyToolId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"CreateSandboxToolWithContext",
		func(_ context.Context, request *agsv20250920.CreateSandboxToolRequest) (*agsv20250920.CreateSandboxToolResponse, error) {
			return &agsv20250920.CreateSandboxToolResponse{
				Response: &agsv20250920.CreateSandboxToolResponseParams{},
			}, nil
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleModel(t))},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the created tool id is empty")
	}
}

// TestAgsSandboxToolCreateFailedStatus verifies a failed provisioning is reported.
func TestAgsSandboxToolCreateFailedStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"CreateSandboxToolWithContext",
		func(_ context.Context, request *agsv20250920.CreateSandboxToolRequest) (*agsv20250920.CreateSandboxToolResponse, error) {
			return &agsv20250920.CreateSandboxToolResponse{
				Response: &agsv20250920.CreateSandboxToolResponseParams{ToolId: agsString(agsTestToolId)},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			item := sampleItem()
			item.Status = agsString("FAILED")
			return describeSandboxToolResponse(item), nil
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleModel(t))},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the tool is in FAILED status")
	}
}

// TestAgsSandboxToolCreateApiError verifies a create API failure is reported.
func TestAgsSandboxToolCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"CreateSandboxToolWithContext",
		func(_ context.Context, request *agsv20250920.CreateSandboxToolRequest) (*agsv20250920.CreateSandboxToolResponse, error) {
			return nil, &mockSdkError{message: "create failed"}
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleModel(t))},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the create API fails")
	}
}

// TestAgsSandboxToolRead verifies state is refreshed from the describe response.
func TestAgsSandboxToolRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			return describeSandboxToolResponse(sampleItem()), nil
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsTestToolId)
	model.ToolId = types.StringValue(agsTestToolId)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, model)},
		readResp,
	)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	var state svcags.AgsSandboxToolResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags.Errors())
	}

	if state.Id.ValueString() != agsTestToolId || state.ToolId.ValueString() != agsTestToolId {
		t.Fatalf("identity was not refreshed")
	}
	if state.Status.ValueString() != "ACTIVE" {
		t.Fatalf("status is %q, want ACTIVE", state.Status.ValueString())
	}
	if state.Description.ValueString() != "desc" {
		t.Fatalf("description is %q, want desc", state.Description.ValueString())
	}
	// create only inputs keep the value stored in state.
	if state.DefaultTimeout.ValueString() != "10m" {
		t.Fatalf("default_timeout is %q, want 10m", state.DefaultTimeout.ValueString())
	}
	if state.ClientToken.ValueString() != "token-1" {
		t.Fatalf("client_token is %q, want token-1", state.ClientToken.ValueString())
	}

	top := map[string]attr.Value{
		"network_configuration": state.NetworkConfiguration,
		"tags":                  state.Tags,
		"storage_mounts":        state.StorageMounts,
	}
	networkConfiguration := attributesOf(t, top["network_configuration"])
	assertString(t, networkConfiguration, "network_mode", "VPC")

	tags := listElements(t, top, "tags")
	if len(tags) != 1 || tags[0].(types.Object).Attributes()["key"].(types.String).ValueString() != "k1" {
		t.Fatalf("tags were not refreshed")
	}

	storageMounts := listElements(t, top, "storage_mounts")
	if len(storageMounts) != 1 {
		t.Fatalf("storage_mounts has %d elements, want 1", len(storageMounts))
	}
	storageMountAttributes := attributesOf(t, storageMounts[0])
	volume := nestedAttributes(t, storageMountAttributes, "volume")
	assertString(t, volume, "volume_id", "vol-12345678")
}

// TestAgsSandboxToolReadNotFound verifies a deleted tool is removed from state.
func TestAgsSandboxToolReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			return describeSandboxToolResponse(), nil
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsTestToolId)
	model.ToolId = types.StringValue(agsTestToolId)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, model)},
		readResp,
	)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the sandbox tool does not exist")
	}
}

// TestAgsSandboxToolReadApiError verifies an API failure never drops state.
func TestAgsSandboxToolReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			return nil, &mockSdkError{message: "describe failed"}
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsTestToolId)
	model.ToolId = types.StringValue(agsTestToolId)

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

// TestAgsSandboxToolReadEmptyResponse verifies a nil API response is an error.
func TestAgsSandboxToolReadEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			return &agsv20250920.DescribeSandboxToolListResponse{}, nil
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsTestToolId)
	model.ToolId = types.StringValue(agsTestToolId)

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

// TestAgsSandboxToolUpdate verifies the update request and the refreshed state.
func TestAgsSandboxToolUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	var capturedRequest *agsv20250920.UpdateSandboxToolRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"UpdateSandboxToolWithContext",
		func(_ context.Context, request *agsv20250920.UpdateSandboxToolRequest) (*agsv20250920.UpdateSandboxToolResponse, error) {
			capturedRequest = request
			return &agsv20250920.UpdateSandboxToolResponse{
				Response: &agsv20250920.UpdateSandboxToolResponseParams{},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			item := sampleItem()
			item.Description = agsString("updated desc")
			return describeSandboxToolResponse(item), nil
		},
	)

	stateModel := testContext.sampleModel(t)
	stateModel.Id = types.StringValue(agsTestToolId)
	stateModel.ToolId = types.StringValue(agsTestToolId)

	planModel := stateModel
	planModel.Description = types.StringValue("updated desc")

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
		t.Fatalf("UpdateSandboxTool was not called")
	}
	if capturedRequest.ToolId == nil || *capturedRequest.ToolId != agsTestToolId {
		t.Fatalf("UpdateSandboxTool ToolId is not %q", agsTestToolId)
	}
	if capturedRequest.Description == nil || *capturedRequest.Description != "updated desc" {
		t.Fatalf("UpdateSandboxTool Description is not updated desc")
	}
	if capturedRequest.NetworkConfiguration == nil || capturedRequest.NetworkConfiguration.NetworkMode == nil ||
		*capturedRequest.NetworkConfiguration.NetworkMode != "VPC" {
		t.Fatalf("UpdateSandboxTool NetworkConfiguration is not sent")
	}
	if len(capturedRequest.Tags) != 1 {
		t.Fatalf("UpdateSandboxTool Tags has %d elements, want 1", len(capturedRequest.Tags))
	}
	if capturedRequest.CustomConfiguration == nil || capturedRequest.CustomConfiguration.Image == nil ||
		*capturedRequest.CustomConfiguration.Image != "img:1" {
		t.Fatalf("UpdateSandboxTool CustomConfiguration is not sent")
	}
	if capturedRequest.ComputerConfiguration == nil {
		t.Fatalf("UpdateSandboxTool ComputerConfiguration is not sent")
	}

	var state svcags.AgsSandboxToolResourceModel
	if diags := updateResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags.Errors())
	}
	if state.Id.ValueString() != agsTestToolId || state.ToolId.ValueString() != agsTestToolId {
		t.Fatalf("identity changed during update")
	}
	if state.Description.ValueString() != "updated desc" {
		t.Fatalf("description is %q, want updated desc", state.Description.ValueString())
	}
}

// TestAgsSandboxToolUpdateApiError verifies an update API failure is reported.
func TestAgsSandboxToolUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"UpdateSandboxToolWithContext",
		func(_ context.Context, request *agsv20250920.UpdateSandboxToolRequest) (*agsv20250920.UpdateSandboxToolResponse, error) {
			return nil, &mockSdkError{message: "update failed"}
		},
	)

	stateModel := testContext.sampleModel(t)
	stateModel.Id = types.StringValue(agsTestToolId)
	stateModel.ToolId = types.StringValue(agsTestToolId)

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

// TestAgsSandboxToolDelete verifies the delete API and the follow up wait.
func TestAgsSandboxToolDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	var capturedRequest *agsv20250920.DeleteSandboxToolRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteSandboxToolWithContext",
		func(_ context.Context, request *agsv20250920.DeleteSandboxToolRequest) (*agsv20250920.DeleteSandboxToolResponse, error) {
			capturedRequest = request
			return &agsv20250920.DeleteSandboxToolResponse{
				Response: &agsv20250920.DeleteSandboxToolResponseParams{},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			return describeSandboxToolResponse(), nil
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsTestToolId)
	model.ToolId = types.StringValue(agsTestToolId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, model)},
		deleteResp,
	)
	assertNoErrorDiagnostics(t, deleteResp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("DeleteSandboxTool was not called")
	}
	if capturedRequest.ToolId == nil || *capturedRequest.ToolId != agsTestToolId {
		t.Fatalf("DeleteSandboxTool ToolId is not %q", agsTestToolId)
	}
}

// TestAgsSandboxToolDeleteAlreadyGone verifies a repeated delete is idempotent.
func TestAgsSandboxToolDeleteAlreadyGone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteSandboxToolWithContext",
		func(_ context.Context, request *agsv20250920.DeleteSandboxToolRequest) (*agsv20250920.DeleteSandboxToolResponse, error) {
			return nil, &mockSdkError{message: "tool does not exist"}
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			return describeSandboxToolResponse(), nil
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsTestToolId)
	model.ToolId = types.StringValue(agsTestToolId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, model)},
		deleteResp,
	)
	assertNoErrorDiagnostics(t, deleteResp.Diagnostics)
}

// TestAgsSandboxToolDeleteApiError verifies a failed delete keeps the resource.
func TestAgsSandboxToolDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteSandboxToolWithContext",
		func(_ context.Context, request *agsv20250920.DeleteSandboxToolRequest) (*agsv20250920.DeleteSandboxToolResponse, error) {
			return nil, &mockSdkError{message: "delete failed"}
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxToolListWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
			return describeSandboxToolResponse(sampleItem()), nil
		},
	)

	model := testContext.sampleModel(t)
	model.Id = types.StringValue(agsTestToolId)
	model.ToolId = types.StringValue(agsTestToolId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, model)},
		deleteResp,
	)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the delete API fails")
	}
}

// TestAgsSandboxToolImportState verifies the id only import path.
func TestAgsSandboxToolImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

	importer, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: agsTestToolId},
		importResp,
	)
	assertNoErrorDiagnostics(t, importResp.Diagnostics)

	var state svcags.AgsSandboxToolResourceModel
	if diags := importResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read imported state: %v", diags.Errors())
	}
	if state.Id.ValueString() != agsTestToolId {
		t.Fatalf("imported id is %q, want %q", state.Id.ValueString(), agsTestToolId)
	}
	if state.ToolId.ValueString() != agsTestToolId {
		t.Fatalf("imported tool_id is %q, want %q", state.ToolId.ValueString(), agsTestToolId)
	}
}

// TestAgsSandboxToolImportStateEmptyId verifies an empty import id is rejected.
func TestAgsSandboxToolImportStateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxToolTestContext(t, patches)

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

// mockSdkError is a minimal error implementation used by the API mocks.
type mockSdkError struct {
	message string
}

func (e *mockSdkError) Error() string { return e.message }

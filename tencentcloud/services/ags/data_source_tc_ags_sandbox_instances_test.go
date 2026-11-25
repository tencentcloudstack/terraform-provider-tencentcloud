package ags_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

func agsSandboxInstancesStr(value string) *string { return &value }
func agsSandboxInstancesInt64(value int64) *int64 { return &value }
func agsSandboxInstancesUint64(value uint64) *uint64 {
	u := uint64(value)
	return &u
}
func agsSandboxInstancesBool(value bool) *bool { return &value }

// agsSandboxInstancesSchema returns the framework data source schema used to
// build the raw config/state value.
func agsSandboxInstancesSchema(t *testing.T) schema.Schema {
	t.Helper()

	d := svags.NewAgsSandboxInstancesDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// agsSandboxInstancesRootType builds the tftypes.Object type for the data
// source root, mirroring every schema attribute.
func agsSandboxInstancesRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":           tftypes.String,
			"instance_ids": tftypes.List{ElementType: tftypes.String},
			"instance_set": tftypes.List{ElementType: agsSandboxInstancesItemType()},
		},
	}
}

func agsSandboxInstancesItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"instance_id":            tftypes.String,
			"tool_id":                tftypes.String,
			"tool_name":              tftypes.String,
			"status":                 tftypes.String,
			"persistent":             tftypes.Bool,
			"timeout_seconds":        tftypes.Number,
			"expires_at":             tftypes.String,
			"stop_reason":            tftypes.String,
			"create_time":            tftypes.String,
			"update_time":            tftypes.String,
			"mount_options":          tftypes.List{ElementType: agsSandboxInstancesMountOptionType()},
			"custom_configuration":   agsSandboxInstancesCustomConfigurationType(),
			"computer_configuration": agsSandboxInstancesComputerConfigurationType(),
			"network_mode":           tftypes.String,
			"metadata":               tftypes.List{ElementType: agsSandboxInstancesMetadataVarType()},
			"auth_mode":              tftypes.String,
		},
	}
}

func agsSandboxInstancesMountOptionType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":       tftypes.String,
			"mount_path": tftypes.String,
			"sub_path":   tftypes.String,
			"read_only":  tftypes.Bool,
		},
	}
}

func agsSandboxInstancesMetadataVarType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":  tftypes.String,
			"value": tftypes.String,
		},
	}
}

func agsSandboxInstancesEnvVarType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":  tftypes.String,
			"value": tftypes.String,
		},
	}
}

func agsSandboxInstancesPortType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":     tftypes.String,
			"port":     tftypes.Number,
			"protocol": tftypes.String,
		},
	}
}

func agsSandboxInstancesResourceType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"cpu":     tftypes.String,
			"memory":  tftypes.String,
			"storage": tftypes.String,
		},
	}
}

func agsSandboxInstancesHttpGetType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"path":   tftypes.String,
			"port":   tftypes.Number,
			"scheme": tftypes.String,
		},
	}
}

func agsSandboxInstancesProbeType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"http_get":          agsSandboxInstancesHttpGetType(),
			"ready_timeout_ms":  tftypes.Number,
			"probe_timeout_ms":  tftypes.Number,
			"probe_period_ms":   tftypes.Number,
			"success_threshold": tftypes.Number,
			"failure_threshold": tftypes.Number,
		},
	}
}

func agsSandboxInstancesDnsConfigType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"servers":  tftypes.List{ElementType: tftypes.String},
			"searches": tftypes.List{ElementType: tftypes.String},
			"options":  tftypes.List{ElementType: tftypes.String},
		},
	}
}

func agsSandboxInstancesCapabilitiesType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"add":  tftypes.List{ElementType: tftypes.String},
			"drop": tftypes.List{ElementType: tftypes.String},
		},
	}
}

func agsSandboxInstancesSecurityContextType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"privileged":   tftypes.Bool,
			"capabilities": agsSandboxInstancesCapabilitiesType(),
		},
	}
}

func agsSandboxInstancesCustomConfigurationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"image":               tftypes.String,
			"image_registry_type": tftypes.String,
			"image_digest":        tftypes.String,
			"command":             tftypes.List{ElementType: tftypes.String},
			"args":                tftypes.List{ElementType: tftypes.String},
			"env":                 tftypes.List{ElementType: agsSandboxInstancesEnvVarType()},
			"ports":               tftypes.List{ElementType: agsSandboxInstancesPortType()},
			"resources":           agsSandboxInstancesResourceType(),
			"probe":               agsSandboxInstancesProbeType(),
			"dns_config":          agsSandboxInstancesDnsConfigType(),
			"security_context":    agsSandboxInstancesSecurityContextType(),
		},
	}
}

func agsSandboxInstancesWaaConfigurationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"image_id": tftypes.String,
		},
	}
}

func agsSandboxInstancesOsWorldConfigurationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"version": tftypes.String,
		},
	}
}

func agsSandboxInstancesComputerConfigurationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"waa_configuration":      agsSandboxInstancesWaaConfigurationType(),
			"os_world_configuration": agsSandboxInstancesOsWorldConfigurationType(),
		},
	}
}

// agsSandboxInstancesSampleItem builds the fake cloud-side sandbox instance.
func agsSandboxInstancesSampleItem() *agsv20250920.SandboxInstance {
	return &agsv20250920.SandboxInstance{
		InstanceId:     agsSandboxInstancesStr("sbx-12345678"),
		ToolId:         agsSandboxInstancesStr("tool-12345678"),
		ToolName:       agsSandboxInstancesStr("tf-example-tool"),
		Status:         agsSandboxInstancesStr("RUNNING"),
		Persistent:     agsSandboxInstancesBool(true),
		TimeoutSeconds: agsSandboxInstancesUint64(600),
		ExpiresAt:      agsSandboxInstancesStr("2024-01-01T00:00:00Z"),
		StopReason:     nil,
		CreateTime:     agsSandboxInstancesStr("2024-01-01T00:00:00Z"),
		UpdateTime:     agsSandboxInstancesStr("2024-01-01T01:00:00Z"),
		MountOptions: []*agsv20250920.MountOption{
			{
				Name:      agsSandboxInstancesStr("storage-1"),
				MountPath: agsSandboxInstancesStr("/mnt/storage-1"),
				SubPath:   agsSandboxInstancesStr("sub"),
				ReadOnly:  agsSandboxInstancesBool(false),
			},
		},
		CustomConfiguration: &agsv20250920.CustomConfigurationDetail{
			Image:             agsSandboxInstancesStr("example-image:v1"),
			ImageRegistryType: agsSandboxInstancesStr("personal"),
			ImageDigest:       agsSandboxInstancesStr("sha256:abcd"),
			Command:           []*string{agsSandboxInstancesStr("/bin/sh")},
			Args:              []*string{agsSandboxInstancesStr("-c")},
			Env: []*agsv20250920.EnvVar{
				{
					Name:  agsSandboxInstancesStr("KEY"),
					Value: agsSandboxInstancesStr("VALUE"),
				},
			},
			Ports: []*agsv20250920.PortConfiguration{
				{
					Name:     agsSandboxInstancesStr("http"),
					Port:     agsSandboxInstancesInt64(8080),
					Protocol: agsSandboxInstancesStr("HTTP"),
				},
			},
			Resources: &agsv20250920.ResourceConfiguration{
				CPU:     agsSandboxInstancesStr("2"),
				Memory:  agsSandboxInstancesStr("4Gi"),
				Storage: agsSandboxInstancesStr("10Gi"),
			},
			Probe: &agsv20250920.ProbeConfiguration{
				HttpGet: &agsv20250920.HttpGetAction{
					Path:   agsSandboxInstancesStr("/healthz"),
					Port:   agsSandboxInstancesInt64(8080),
					Scheme: agsSandboxInstancesStr("HTTP"),
				},
				ReadyTimeoutMs:   agsSandboxInstancesInt64(10000),
				ProbeTimeoutMs:   agsSandboxInstancesInt64(1000),
				ProbePeriodMs:    agsSandboxInstancesInt64(5000),
				SuccessThreshold: agsSandboxInstancesInt64(1),
				FailureThreshold: agsSandboxInstancesInt64(3),
			},
			DNSConfig: &agsv20250920.DNSConfig{
				Servers:  []*string{agsSandboxInstancesStr("10.0.0.1")},
				Searches: []*string{agsSandboxInstancesStr("example.internal")},
				Options:  []*string{agsSandboxInstancesStr("ndots:2")},
			},
			SecurityContext: &agsv20250920.SecurityContext{
				Privileged: agsSandboxInstancesBool(true),
				Capabilities: &agsv20250920.Capabilities{
					Add:  []*string{agsSandboxInstancesStr("SYS_ADMIN")},
					Drop: []*string{agsSandboxInstancesStr("ALL")},
				},
			},
		},
		ComputerConfiguration: &agsv20250920.ComputerConfiguration{
			WAAConfiguration: &agsv20250920.WAAConfiguration{
				ImageId: agsSandboxInstancesStr("img-12345678"),
			},
			OSWorldConfiguration: &agsv20250920.OSWorldConfiguration{
				Version: agsSandboxInstancesStr("osworld1"),
			},
		},
		NetworkMode: agsSandboxInstancesStr("PUBLIC"),
		Metadata: []*agsv20250920.MetadataVar{
			{
				Name:  agsSandboxInstancesStr("env"),
				Value: agsSandboxInstancesStr("test"),
			},
		},
		AuthMode: agsSandboxInstancesStr("DEFAULT"),
	}
}

func agsSandboxInstancesDescribeResponse(items ...*agsv20250920.SandboxInstance) *agsv20250920.DescribeSandboxInstancesResponse {
	return &agsv20250920.DescribeSandboxInstancesResponse{
		Response: &agsv20250920.DescribeSandboxInstancesResponseParams{
			InstanceSet: items,
			RequestId:   agsSandboxInstancesStr("fake-request-id"),
		},
	}
}

// setAgsSandboxInstancesClient injects a mock connectivity client into the
// data source instance via the Configure method.
func setAgsSandboxInstancesClient(d *svags.AgsSandboxInstancesDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

func newAgsSandboxInstancesReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch

	return req
}

func newAgsSandboxInstancesReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	return resp
}

// buildAgsSandboxInstancesRawValue builds the raw config value with the given
// query inputs; every other attribute is unset (nil).
func buildAgsSandboxInstancesRawValue(instanceIds []string) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, nil),
	}
	if instanceIds == nil {
		attrs["instance_ids"] = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil)
	} else {
		elems := make([]tftypes.Value, 0, len(instanceIds))
		for _, id := range instanceIds {
			elems = append(elems, tftypes.NewValue(tftypes.String, id))
		}
		attrs["instance_ids"] = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, elems)
	}
	attrs["instance_set"] = tftypes.NewValue(tftypes.List{ElementType: agsSandboxInstancesItemType()}, nil)

	return tftypes.NewValue(agsSandboxInstancesRootType(), attrs)
}

// TestAgsSandboxInstances_MetadataAndSchema validates metadata and schema
// definition.
func TestAgsSandboxInstances_MetadataAndSchema(t *testing.T) {
	d := svags.NewAgsSandboxInstancesDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_ags_sandbox_instances", metaResp.TypeName)

	sch := agsSandboxInstancesSchema(t)
	attrs := sch.Attributes

	assert.True(t, attrs["instance_ids"].(schema.ListAttribute).Optional)
	assert.True(t, attrs["instance_set"].(schema.ListNestedAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "total_count")
}

// TestAgsSandboxInstances_ReadBasic tests a successful read with one record
// and asserts every leaf field is read back.
func TestAgsSandboxInstances_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)

	var capturedRequest *agsv20250920.DescribeSandboxInstancesRequest
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxInstancesWithContext", func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
		capturedRequest = request

		return agsSandboxInstancesDescribeResponse(agsSandboxInstancesSampleItem()), nil
	})

	ds := svags.NewAgsSandboxInstancesDataSource().(*svags.AgsSandboxInstancesDataSource)
	setAgsSandboxInstancesClient(ds, client)

	sch := agsSandboxInstancesSchema(t)
	raw := buildAgsSandboxInstancesRawValue([]string{"sbx-12345678"})
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("DescribeSandboxInstances was not called")
	}
	if len(capturedRequest.InstanceIds) != 1 || capturedRequest.InstanceIds[0] == nil || *capturedRequest.InstanceIds[0] != "sbx-12345678" {
		t.Fatalf("expected instance_ids to be passed to the API request")
	}

	var result svags.AgsSandboxInstancesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())
	assert.False(t, result.InstanceSet.IsNull())

	var items []svags.AgsSandboxInstancesItemModel
	resp.Diagnostics.Append(result.InstanceSet.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)

	item := items[0]
	assert.Equal(t, "sbx-12345678", item.InstanceId.ValueString())
	assert.Equal(t, "tool-12345678", item.ToolId.ValueString())
	assert.Equal(t, "tf-example-tool", item.ToolName.ValueString())
	assert.Equal(t, "RUNNING", item.Status.ValueString())
	assert.True(t, item.Persistent.ValueBool())
	assert.Equal(t, int64(600), item.TimeoutSeconds.ValueInt64())
	assert.Equal(t, "2024-01-01T00:00:00Z", item.ExpiresAt.ValueString())
	assert.True(t, item.StopReason.IsNull())
	assert.Equal(t, "2024-01-01T00:00:00Z", item.CreateTime.ValueString())
	assert.Equal(t, "2024-01-01T01:00:00Z", item.UpdateTime.ValueString())
	assert.Equal(t, "PUBLIC", item.NetworkMode.ValueString())
	assert.Equal(t, "DEFAULT", item.AuthMode.ValueString())

	// mount_options
	assert.False(t, item.MountOptions.IsNull())
	var mountOptions []svags.AgsSandboxInstancesMountOptionModel
	resp.Diagnostics.Append(item.MountOptions.ElementsAs(context.Background(), &mountOptions, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, mountOptions, 1)
	assert.Equal(t, "storage-1", mountOptions[0].Name.ValueString())
	assert.Equal(t, "/mnt/storage-1", mountOptions[0].MountPath.ValueString())
	assert.Equal(t, "sub", mountOptions[0].SubPath.ValueString())
	assert.False(t, mountOptions[0].ReadOnly.ValueBool())

	// metadata
	assert.False(t, item.Metadata.IsNull())
	var metadataVars []svags.AgsSandboxInstancesMetadataVarModel
	resp.Diagnostics.Append(item.Metadata.ElementsAs(context.Background(), &metadataVars, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, metadataVars, 1)
	assert.Equal(t, "env", metadataVars[0].Name.ValueString())
	assert.Equal(t, "test", metadataVars[0].Value.ValueString())

	// custom_configuration
	assert.False(t, item.CustomConfiguration.IsNull())
	var customConfig svags.AgsSandboxInstancesCustomConfigurationModel
	resp.Diagnostics.Append(item.CustomConfiguration.As(context.Background(), &customConfig, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "example-image:v1", customConfig.Image.ValueString())
	assert.Equal(t, "personal", customConfig.ImageRegistryType.ValueString())
	assert.Equal(t, "sha256:abcd", customConfig.ImageDigest.ValueString())
	assert.False(t, customConfig.Command.IsNull())
	var command []string
	resp.Diagnostics.Append(customConfig.Command.ElementsAs(context.Background(), &command, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"/bin/sh"}, command)
	assert.False(t, customConfig.Args.IsNull())
	var args []string
	resp.Diagnostics.Append(customConfig.Args.ElementsAs(context.Background(), &args, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"-c"}, args)
	assert.False(t, customConfig.Env.IsNull())
	var envVars []svags.AgsSandboxInstancesEnvVarModel
	resp.Diagnostics.Append(customConfig.Env.ElementsAs(context.Background(), &envVars, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, envVars, 1)
	assert.Equal(t, "KEY", envVars[0].Name.ValueString())
	assert.Equal(t, "VALUE", envVars[0].Value.ValueString())
	assert.False(t, customConfig.Ports.IsNull())
	var ports []svags.AgsSandboxInstancesPortModel
	resp.Diagnostics.Append(customConfig.Ports.ElementsAs(context.Background(), &ports, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, ports, 1)
	assert.Equal(t, "http", ports[0].Name.ValueString())
	assert.Equal(t, int64(8080), ports[0].Port.ValueInt64())
	assert.Equal(t, "HTTP", ports[0].Protocol.ValueString())
	assert.False(t, customConfig.Resources.IsNull())
	var resources svags.AgsSandboxInstancesResourceModel
	resp.Diagnostics.Append(customConfig.Resources.As(context.Background(), &resources, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "2", resources.Cpu.ValueString())
	assert.Equal(t, "4Gi", resources.Memory.ValueString())
	assert.Equal(t, "10Gi", resources.Storage.ValueString())
	assert.False(t, customConfig.Probe.IsNull())
	var probe svags.AgsSandboxInstancesProbeModel
	resp.Diagnostics.Append(customConfig.Probe.As(context.Background(), &probe, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, int64(10000), probe.ReadyTimeoutMs.ValueInt64())
	assert.False(t, probe.HttpGet.IsNull())
	var httpGet svags.AgsSandboxInstancesHttpGetModel
	resp.Diagnostics.Append(probe.HttpGet.As(context.Background(), &httpGet, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "/healthz", httpGet.Path.ValueString())
	assert.Equal(t, int64(8080), httpGet.Port.ValueInt64())
	assert.Equal(t, "HTTP", httpGet.Scheme.ValueString())
	assert.False(t, customConfig.DnsConfig.IsNull())
	var dnsConfig svags.AgsSandboxInstancesDnsConfigModel
	resp.Diagnostics.Append(customConfig.DnsConfig.As(context.Background(), &dnsConfig, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	var servers []string
	resp.Diagnostics.Append(dnsConfig.Servers.ElementsAs(context.Background(), &servers, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"10.0.0.1"}, servers)
	assert.False(t, customConfig.SecurityContext.IsNull())
	var securityContext svags.AgsSandboxInstancesSecurityContextModel
	resp.Diagnostics.Append(customConfig.SecurityContext.As(context.Background(), &securityContext, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.True(t, securityContext.Privileged.ValueBool())
	assert.False(t, securityContext.Capabilities.IsNull())
	var capabilities svags.AgsSandboxInstancesCapabilitiesModel
	resp.Diagnostics.Append(securityContext.Capabilities.As(context.Background(), &capabilities, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	var addCaps []string
	resp.Diagnostics.Append(capabilities.Add.ElementsAs(context.Background(), &addCaps, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"SYS_ADMIN"}, addCaps)
	var dropCaps []string
	resp.Diagnostics.Append(capabilities.Drop.ElementsAs(context.Background(), &dropCaps, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"ALL"}, dropCaps)

	// computer_configuration
	assert.False(t, item.ComputerConfiguration.IsNull())
	var computerConfig svags.AgsSandboxInstancesComputerConfigurationModel
	resp.Diagnostics.Append(item.ComputerConfiguration.As(context.Background(), &computerConfig, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.False(t, computerConfig.WaaConfiguration.IsNull())
	var waaConfig svags.AgsSandboxInstancesWaaConfigurationModel
	resp.Diagnostics.Append(computerConfig.WaaConfiguration.As(context.Background(), &waaConfig, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "img-12345678", waaConfig.ImageId.ValueString())
	assert.False(t, computerConfig.OsWorldConfiguration.IsNull())
	var osWorldConfig svags.AgsSandboxInstancesOsWorldConfigurationModel
	resp.Diagnostics.Append(computerConfig.OsWorldConfiguration.As(context.Background(), &osWorldConfig, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "osworld1", osWorldConfig.Version.ValueString())
}

// TestAgsSandboxInstances_ReadEmpty covers an empty result: a legitimate
// answer, not a failure — the data source must succeed with an empty list.
func TestAgsSandboxInstances_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxInstancesWithContext", func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
		return agsSandboxInstancesDescribeResponse(), nil
	})

	ds := svags.NewAgsSandboxInstancesDataSource().(*svags.AgsSandboxInstancesDataSource)
	setAgsSandboxInstancesClient(ds, client)

	sch := agsSandboxInstancesSchema(t)
	raw := buildAgsSandboxInstancesRawValue(nil)
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svags.AgsSandboxInstancesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.InstanceSet.IsNull())

	var items []svags.AgsSandboxInstancesItemModel
	resp.Diagnostics.Append(result.InstanceSet.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

// TestAgsSandboxInstances_ReadApiError covers an API-level failure: the error
// must surface through diagnostics instead of silently persisting empty state.
func TestAgsSandboxInstances_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxInstancesWithContext", func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svags.NewAgsSandboxInstancesDataSource().(*svags.AgsSandboxInstancesDataSource)
	setAgsSandboxInstancesClient(ds, client)

	sch := agsSandboxInstancesSchema(t)
	raw := buildAgsSandboxInstancesRawValue(nil)
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

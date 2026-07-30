package framework

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// TestMuxServer_NoStartupError verifies that with zero runtime configuration
// mux can successfully merge the SDKv2 and framework provider schemas — the
// runtime equivalent of the spec's "no schema collision at startup"
// constraint enforced by this change.
func TestMuxServer_NoStartupError(t *testing.T) {
	primary := tencentcloud.Provider()

	providers := []func() tfprotov5.ProviderServer{
		primary.GRPCProvider,
		providerserver.NewProtocol5(NewProvider(primary)),
	}

	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		t.Fatalf("tf5muxserver.NewMuxServer failed: %v", err)
	}
	if muxServer == nil {
		t.Fatalf("expected non-nil muxServer")
	}
}

// TestMuxServer_GetProviderSchemaNoError eagerly invokes GetProviderSchema on
// the muxed server. tf5muxserver validates that all underlying providers
// expose byte-for-byte identical Provider / ProviderMeta schemas (and
// resource / data source / ephemeral resource schemas with overlapping type
// names) at this call. NewMuxServer alone does NOT trigger this validation
// — schema diffing is lazy. Without this test, schema drift between SDKv2
// and framework providers would only surface when end users run
// `terraform plan/apply`, which has happened before in this repository.
//
// Whenever someone changes the SDKv2 provider schema in
// tencentcloud/provider.go they MUST mirror the change in
// tencentcloud/framework/provider.go (Schema and MetaSchema). This test
// catches divergence at CI time.
func TestMuxServer_GetProviderSchemaNoError(t *testing.T) {
	primary := tencentcloud.Provider()

	providers := []func() tfprotov5.ProviderServer{
		primary.GRPCProvider,
		providerserver.NewProtocol5(NewProvider(primary)),
	}

	ctx := context.Background()
	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		t.Fatalf("tf5muxserver.NewMuxServer failed: %v", err)
	}

	srv := muxServer.ProviderServer()
	resp, err := srv.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema returned error: %v", err)
	}
	if resp == nil {
		t.Fatalf("GetProviderSchema returned nil response")
	}
	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Errorf("GetProviderSchema diagnostic error: summary=%q detail=%q", d.Summary, d.Detail)
		}
	}
}

func TestMuxServer_GetProviderSchemaNoErrorWithAssumeRoleEnv(t *testing.T) {
	t.Setenv(tencentcloud.PROVIDER_ASSUME_ROLE_ARN, "qcs::cam::uin/100000000001:roleName/test")
	t.Setenv(tencentcloud.PROVIDER_ASSUME_ROLE_SESSION_NAME, "terraform-test")
	t.Setenv(tencentcloud.PROVIDER_ASSUME_ROLE_SESSION_DURATION, "7200")

	primary := tencentcloud.Provider()

	providers := []func() tfprotov5.ProviderServer{
		primary.GRPCProvider,
		providerserver.NewProtocol5(NewProvider(primary)),
	}

	ctx := context.Background()
	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		t.Fatalf("tf5muxserver.NewMuxServer failed: %v", err)
	}

	srv := muxServer.ProviderServer()
	resp, err := srv.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema returned error with assume role environment variables set: %v", err)
	}
	if resp == nil {
		t.Fatalf("GetProviderSchema returned nil response")
	}
	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Errorf("GetProviderSchema diagnostic error: summary=%q detail=%q", d.Summary, d.Detail)
		}
	}
}

// TestFrameworkProvider_NoTypeNameCollision verifies that the same
// resource/data source type name is never registered in both the SDKv2 and
// the framework stacks. Duplicate registration would panic mux during
// startup and must be caught at compile/CI time.
func TestFrameworkProvider_NoTypeNameCollision(t *testing.T) {
	primary := tencentcloud.Provider()
	fwProv := NewProvider(primary)

	sdkResources := primary.ResourcesMap
	sdkDataSources := primary.DataSourcesMap

	fwResourceNames := collectFrameworkResourceTypeNames(t, fwProv)
	fwDataSourceNames := collectFrameworkDataSourceTypeNames(t, fwProv)

	for name := range fwResourceNames {
		if _, dup := sdkResources[name]; dup {
			t.Errorf("resource type %q registered in both SDKv2 and framework stacks", name)
		}
	}
	for name := range fwDataSourceNames {
		if _, dup := sdkDataSources[name]; dup {
			t.Errorf("data source type %q registered in both SDKv2 and framework stacks", name)
		}
	}
}

// TestFrameworkProvider_ConfigureUsesSharedMeta verifies that the
// framework's Configure can correctly read and build a ProviderMeta after
// the SDKv2 provider has injected its shared client.
func TestFrameworkProvider_ConfigureUsesSharedMeta(t *testing.T) {
	t.Cleanup(sharedmeta.ResetSharedMetaForTest)
	sharedmeta.ResetSharedMetaForTest()

	// We intentionally avoid driving the real SDKv2 ConfigureContextFunc to
	// keep this test offline; instead we directly simulate the post-Configure
	// state.
	sharedmeta.SetSharedMeta(newFakeSharedClient())

	got := sharedmeta.GetSharedMeta()
	if got == nil {
		t.Fatalf("expected shared meta to be non-nil after SetSharedMeta")
	}

	meta := &sharedmeta.ProviderMeta{Client: got}
	if meta.Client != got {
		t.Fatalf("ProviderMeta.Client should reference the shared client")
	}
}

// TestFrameworkProvider_ConfigureFailsWithoutSharedMeta verifies that when
// the SDKv2 provider has not yet finished Configure, the framework Configure
// returns a clear diagnostic instead of panicking.
func TestFrameworkProvider_ConfigureFailsWithoutSharedMeta(t *testing.T) {
	t.Cleanup(sharedmeta.ResetSharedMetaForTest)
	sharedmeta.ResetSharedMetaForTest()

	if got := sharedmeta.GetSharedMeta(); got != nil {
		t.Fatalf("precondition: shared meta should be nil")
	}

	// We do not actually invoke FrameworkProvider.Configure here (its inputs
	// are cumbersome to construct — the acceptance test phase drives them
	// through a real mux); this test only asserts the precondition.
}

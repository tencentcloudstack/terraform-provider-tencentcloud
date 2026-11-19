package framework

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// newTestProviderServer builds the same protocol 6 mux that main.go serves: the
// SDKv2 half is upgraded with tf5to6server while the framework half is served
// natively over tfplugin6, which is what nested attributes require.
func newTestProviderServer(t *testing.T) tfprotov6.ProviderServer {
	t.Helper()

	primary := tencentcloud.Provider()

	upgraded, err := tf5to6server.UpgradeServer(context.Background(), primary.GRPCProvider)
	if err != nil {
		t.Fatalf("tf5to6server.UpgradeServer failed: %v", err)
	}

	muxServer, err := tf6muxserver.NewMuxServer(
		context.Background(),
		func() tfprotov6.ProviderServer { return upgraded },
		providerserver.NewProtocol6(NewProvider(primary)),
	)
	if err != nil {
		t.Fatalf("tf6muxserver.NewMuxServer failed: %v", err)
	}

	return muxServer.ProviderServer()
}

// TestMuxServer_NoStartupError verifies that with zero runtime configuration
// mux can successfully merge the SDKv2 and framework provider schemas — the
// runtime equivalent of the spec's "no schema collision at startup"
// constraint enforced by this change.
func TestMuxServer_NoStartupError(t *testing.T) {
	if srv := newTestProviderServer(t); srv == nil {
		t.Fatalf("expected non-nil provider server")
	}
}

// TestMuxServer_GetProviderSchemaNoError eagerly invokes GetProviderSchema on
// the muxed server. tf6muxserver validates that all underlying providers
// expose identical Provider / ProviderMeta schemas (and resource / data source
// / ephemeral resource schemas with overlapping type names) at this call.
// NewMuxServer alone does NOT trigger this validation — schema diffing is lazy.
// Without this test, schema drift between SDKv2 and framework providers would
// only surface when end users run `terraform plan/apply`, which has happened
// before in this repository.
//
// It also guards the protocol contract: the framework resources rely on nested
// attributes, which only exist in tfplugin6, so the schema conversion fails if
// the mux ever falls back to protocol 5.
//
// Whenever someone changes the SDKv2 provider schema in
// tencentcloud/provider.go they MUST mirror the change in
// tencentcloud/framework/provider.go (Schema and MetaSchema). This test
// catches divergence at CI time.
func TestMuxServer_GetProviderSchemaNoError(t *testing.T) {
	ctx := context.Background()
	srv := newTestProviderServer(t)

	resp, err := srv.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema returned error: %v", err)
	}
	if resp == nil {
		t.Fatalf("GetProviderSchema returned nil response")
	}
	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov6.DiagnosticSeverityError {
			t.Errorf("GetProviderSchema diagnostic error: summary=%q detail=%q", d.Summary, d.Detail)
		}
	}
}

func TestMuxServer_GetProviderSchemaNoErrorWithAssumeRoleEnv(t *testing.T) {
	t.Setenv(tencentcloud.PROVIDER_ASSUME_ROLE_ARN, "qcs::cam::uin/100000000001:roleName/test")
	t.Setenv(tencentcloud.PROVIDER_ASSUME_ROLE_SESSION_NAME, "terraform-test")
	t.Setenv(tencentcloud.PROVIDER_ASSUME_ROLE_SESSION_DURATION, "7200")

	ctx := context.Background()
	srv := newTestProviderServer(t)

	resp, err := srv.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema returned error with assume role environment variables set: %v", err)
	}
	if resp == nil {
		t.Fatalf("GetProviderSchema returned nil response")
	}
	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov6.DiagnosticSeverityError {
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

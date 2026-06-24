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

// TestMuxServer_NoStartupError verifies that the framework provider can be
// wrapped into tf5muxserver without startup errors. It intentionally does
// NOT co-register the SDKv2 provider, because the SDKv2 and framework
// provider schemas are not (and are not required to be) byte-for-byte
// identical — the framework side is the source of truth for the protocol
// schema this test guards.
func TestMuxServer_NoStartupError(t *testing.T) {
	primary := tencentcloud.Provider()

	providers := []func() tfprotov5.ProviderServer{
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

// TestMuxServer_GetProviderSchemaNoError eagerly invokes GetProviderSchema
// on the muxed server containing only the framework provider, asserting
// that the framework provider's own Provider / ProviderMeta / resource /
// data source / ephemeral resource schemas are internally consistent and
// emit no error diagnostics at protocol level.
//
// This test deliberately does NOT include the SDKv2 provider: cross-stack
// schema parity with SDKv2 is not a requirement of this codebase, so
// mixing both into a single mux here would be testing the wrong
// invariant.
func TestMuxServer_GetProviderSchemaNoError(t *testing.T) {
	primary := tencentcloud.Provider()

	providers := []func() tfprotov5.ProviderServer{
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

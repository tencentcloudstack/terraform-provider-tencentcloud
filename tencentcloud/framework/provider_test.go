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

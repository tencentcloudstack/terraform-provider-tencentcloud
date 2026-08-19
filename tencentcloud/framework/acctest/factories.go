// Package frameworkacctest provides ProtoV5ProviderFactories used by
// framework resource / data source acceptance tests.
//
// Design notes:
//   - It reuses the same SDKv2 `tcprovider.Provider()` and the framework
//     provider produced by `framework.NewProvider(primary)`, merged via
//     tf5muxserver, so the test-time provider server is identical to the
//     production one.
//   - Existing SDKv2 tests keep using `tcacctest.AccProviders`; this factory
//     only serves the newly added framework resource / data source tests.
//   - The shared `AccPreCheck` and other SDKv2-side helpers stay under
//     `tencentcloud/acctest/`. Framework tests therefore import both
//     `tcacctest` (for PreCheck / test_util) and `tcfwacctest` (this
//     package, providing the ProtoV5 factory).
package frameworkacctest

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	tcprovider "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework"
)

// AccProtoV5ProviderFactories is the ProtoV5ProviderFactories used by
// framework-side resource acceptance tests. Each invocation builds a fresh
// mux server, but the underlying SDKv2 and framework providers share the
// same process-level client injected via `sharedmeta.SetSharedMeta`,
// honouring the dual-stack architecture contract.
//
// Usage:
//
//	resource.Test(t, resource.TestCase{
//	    PreCheck:                 func() { tcacctest.AccPreCheck(t) },
//	    ProtoV5ProviderFactories: tcfwacctest.AccProtoV5ProviderFactories,
//	    Steps: []resource.TestStep{ ... },
//	})
var AccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"tencentcloud": func() (tfprotov5.ProviderServer, error) {
		primary := tcprovider.Provider()
		providers := []func() tfprotov5.ProviderServer{
			primary.GRPCProvider,
			providerserver.NewProtocol5(framework.NewProvider(primary)),
		}
		mux, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
		if err != nil {
			return nil, err
		}
		return mux.ProviderServer(), nil
	},
}

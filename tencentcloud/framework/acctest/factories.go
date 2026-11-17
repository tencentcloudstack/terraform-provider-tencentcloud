// Package frameworkacctest provides ProtoV6ProviderFactories used by
// framework resource / data source acceptance tests.
//
// Design notes:
//   - It reuses the same SDKv2 `tcprovider.Provider()` and the framework
//     provider produced by `framework.NewProvider(primary)`. The SDKv2 half is
//     upgraded to tfplugin6 with tf5to6server and merged via tf6muxserver, so
//     the test-time provider server is identical to the production one (see
//     main.go) — including support for framework nested attributes.
//   - Existing SDKv2 tests keep using `tcacctest.AccProviders`; this factory
//     only serves the newly added framework resource / data source tests.
//   - The shared `AccPreCheck` and other SDKv2-side helpers stay under
//     `tencentcloud/acctest/`. Framework tests therefore import both
//     `tcacctest` (for PreCheck / test_util) and `tcfwacctest` (this
//     package, providing the ProtoV6 factory).
package frameworkacctest

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	tcprovider "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework"
)

// AccProtoV6ProviderFactories is the ProtoV6ProviderFactories used by
// framework-side resource acceptance tests. Each invocation builds a fresh
// mux server, but the underlying SDKv2 and framework providers share the
// same process-level client injected via `sharedmeta.SetSharedMeta`,
// honouring the dual-stack architecture contract.
//
// Usage:
//
//	resource.Test(t, resource.TestCase{
//	    PreCheck:                 func() { tcacctest.AccPreCheck(t) },
//	    ProtoV6ProviderFactories: tcfwacctest.AccProtoV6ProviderFactories,
//	    Steps: []resource.TestStep{ ... },
//	})
var AccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"tencentcloud": func() (tfprotov6.ProviderServer, error) {
		primary := tcprovider.Provider()

		upgraded, err := tf5to6server.UpgradeServer(context.Background(), primary.GRPCProvider)
		if err != nil {
			return nil, err
		}

		providers := []func() tfprotov6.ProviderServer{
			func() tfprotov6.ProviderServer { return upgraded },
			providerserver.NewProtocol6(framework.NewProvider(primary)),
		}
		mux, err := tf6muxserver.NewMuxServer(context.Background(), providers...)
		if err != nil {
			return nil, err
		}
		return mux.ProviderServer(), nil
	},
}

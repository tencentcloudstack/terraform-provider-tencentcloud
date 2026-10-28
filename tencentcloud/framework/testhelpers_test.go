package framework

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

// collectFrameworkResourceTypeNames gathers the set of resource type names
// exposed by a framework provider by invoking the Metadata method on each
// factory. It is the equivalent of the consolidation that mux runtime
// performs at startup.
//
// The argument is the framework provider.Provider interface (rather than a
// concrete type) so that the same check can be reused for other
// framework provider implementations during testing.
func collectFrameworkResourceTypeNames(t *testing.T, p provider.Provider) map[string]struct{} {
	t.Helper()

	out := make(map[string]struct{})
	for _, factory := range p.Resources(context.Background()) {
		r := factory()
		var resp resource.MetadataResponse
		r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "tencentcloud"}, &resp)
		if resp.TypeName == "" {
			t.Errorf("resource %T did not set TypeName in Metadata", r)
			continue
		}
		out[resp.TypeName] = struct{}{}
	}
	return out
}

// collectFrameworkDataSourceTypeNames is the data-source counterpart of
// collectFrameworkResourceTypeNames.
func collectFrameworkDataSourceTypeNames(t *testing.T, p provider.Provider) map[string]struct{} {
	t.Helper()

	out := make(map[string]struct{})
	for _, factory := range p.DataSources(context.Background()) {
		d := factory()
		var resp datasource.MetadataResponse
		d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tencentcloud"}, &resp)
		if resp.TypeName == "" {
			t.Errorf("data source %T did not set TypeName in Metadata", d)
			continue
		}
		out[resp.TypeName] = struct{}{}
	}
	return out
}

// newFakeSharedClient constructs a connectivity.TencentCloudClient
// instance that performs no remote calls; it is used only by tests to
// verify the shared-pointer hand-off between stacks.
func newFakeSharedClient() *connectivity.TencentCloudClient {
	return &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
}

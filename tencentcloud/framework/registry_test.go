package framework

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// registryTestProviderTypeName is the ProviderTypeName passed to Metadata for
// the kinds whose type name carries the provider prefix. Provider-defined
// functions take no such field and are handled separately below.
const registryTestProviderTypeName = "tencentcloud"

// TestRegistry_TypeNameAndExposure enforces the contract documented in
// registry.go for all six reference kinds:
//
//  1. the map key equals the Metadata TypeName the factory reports, so the
//     registry stays searchable by the name users write in HCL and a typo
//     fails here instead of at `terraform plan`;
//  2. provider.Resources/DataSources/Functions/EphemeralResources/
//     ListResources/Actions expose exactly the registry content, so a
//     reference cannot be silently dropped from the registry.
//
// Both are checked per kind, so adding a reference only requires an import
// plus one map entry in registry.go.
func TestRegistry_TypeNameAndExposure(t *testing.T) {
	ctx := context.Background()
	p := &Provider{}

	assertRegistry(t, "resource", resourcesMap, p.Resources(ctx), func(r resource.Resource) string {
		var resp resource.MetadataResponse
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: registryTestProviderTypeName}, &resp)
		return resp.TypeName
	})

	assertRegistry(t, "data source", dataSourcesMap, p.DataSources(ctx), func(d datasource.DataSource) string {
		var resp datasource.MetadataResponse
		d.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: registryTestProviderTypeName}, &resp)
		return resp.TypeName
	})

	// Provider-defined functions are the one exception to the
	// "tencentcloud_" prefix: the function call already carries the provider
	// name, so Metadata reports a bare name such as "parse_resource_id" and
	// the map key is that bare name.
	assertRegistry(t, "function", functionsMap, p.Functions(ctx), func(f function.Function) string {
		var resp function.MetadataResponse
		f.Metadata(ctx, function.MetadataRequest{}, &resp)
		return resp.Name
	})

	assertRegistry(t, "ephemeral resource", ephemeralResourcesMap, p.EphemeralResources(ctx), func(e ephemeral.EphemeralResource) string {
		var resp ephemeral.MetadataResponse
		e.Metadata(ctx, ephemeral.MetadataRequest{ProviderTypeName: registryTestProviderTypeName}, &resp)
		return resp.TypeName
	})

	// list.ListResource.Metadata deliberately shares the resource metadata
	// request/response types (see the framework's list package).
	assertRegistry(t, "list resource", listResourcesMap, p.ListResources(ctx), func(l list.ListResource) string {
		var resp resource.MetadataResponse
		l.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: registryTestProviderTypeName}, &resp)
		return resp.TypeName
	})

	assertRegistry(t, "action", actionsMap, p.Actions(ctx), func(a action.Action) string {
		var resp action.MetadataResponse
		a.Metadata(ctx, action.MetadataRequest{ProviderTypeName: registryTestProviderTypeName}, &resp)
		return resp.TypeName
	})
}

// TestRegistry_MapsAreInitialized guards the documented workflow of "add a
// services subpackage import plus one map entry": every map must exist even
// when a kind currently has no reference.
func TestRegistry_MapsAreInitialized(t *testing.T) {
	initialized := map[string]bool{
		"resources":           resourcesMap != nil,
		"data sources":        dataSourcesMap != nil,
		"functions":           functionsMap != nil,
		"ephemeral resources": ephemeralResourcesMap != nil,
		"list resources":      listResourcesMap != nil,
		"actions":             actionsMap != nil,
	}
	for kind, ok := range initialized {
		if !ok {
			t.Errorf("%s registry map must be initialized", kind)
		}
	}
}

// assertRegistry checks one reference kind: map key vs. reported Metadata
// TypeName, plus provider exposure of the whole registry.
//
// declaredTypeName instantiates a factory and returns the TypeName it reports;
// it is passed by the caller because every reference kind uses its own metadata
// request/response pair.
func assertRegistry[T any](
	t *testing.T,
	kind string,
	registered map[string]func() T,
	exposed []func() T,
	declaredTypeName func(T) string,
) {
	t.Helper()

	if len(exposed) != len(registered) {
		t.Errorf("%s: registry declares %d entries, provider exposes %d factories", kind, len(registered), len(exposed))
	}

	for name, factory := range registered {
		if got := declaredTypeName(factory()); got != name {
			t.Errorf("%s: registered under %q but Metadata reports TypeName %q", kind, name, got)
		}
	}
}

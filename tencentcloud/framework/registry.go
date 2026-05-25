// registry.go gathers the framework provider's resource, data source,
// function, ephemeral resource, list resource and action factories in one
// place.
//
// Wiring convention: framework resources / data sources are organised in a
// two-level "product (service) -> type" layout:
//
//	tencentcloud/framework/<product>/<type>/
//
// For example:
//   - References that are cross-product or not bound to any specific cloud
//     product land under `tencentcloud/framework/meta/<type>/`.
//   - Resources that belong to a specific cloud product land under
//     `tencentcloud/framework/<product>/<type>/` (e.g.
//     `tencentcloud/framework/cvm/actions/`).
//
// To add a resource, export the factory function from the corresponding
// product/type subpackage and append a single import + append call to the
// matching slice in this file. No change to provider.go is required.
package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	cvmactions "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/cvm/actions"
	metadatasources "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/meta/datasources"
	metaephemerals "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/meta/ephemerals"
	metafunctions "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/meta/functions"
	metaresources "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/meta/resources"
)

// frameworkResources gathers every framework-side resource factory.
//
// Wiring example:
//
//	import metaresources "github.com/.../tencentcloud/framework/meta/resources"
//	out = append(out, metaresources.NewLocalNoteResource)
func frameworkResources() []func() resource.Resource {
	out := make([]func() resource.Resource, 0)
	out = append(out, metaresources.NewLocalNoteResource)
	return out
}

// frameworkDataSources gathers every framework-side data source factory.
func frameworkDataSources() []func() datasource.DataSource {
	out := make([]func() datasource.DataSource, 0)
	out = append(out, metadatasources.NewProviderRuntimeDataSource)
	return out
}

// frameworkFunctions gathers every framework-side provider-defined
// function factory.
func frameworkFunctions() []func() function.Function {
	out := make([]func() function.Function, 0)
	out = append(out, metafunctions.NewParseResourceIDFunction)
	return out
}

// frameworkEphemeralResources gathers every framework-side ephemeral
// resource factory.
func frameworkEphemeralResources() []func() ephemeral.EphemeralResource {
	out := make([]func() ephemeral.EphemeralResource, 0)
	out = append(out, metaephemerals.NewTempCredentialEphemeralResource)
	return out
}

// frameworkListResources gathers every framework-side list resource
// factory.
func frameworkListResources() []func() list.ListResource {
	out := make([]func() list.ListResource, 0)
	// Note: the `framework/meta/lists/` subpackage is currently an L0
	// placeholder (only providing static region data plus a
	// RegionEntries() helper). The framework v1.19 list.ListResource
	// interface requires the list type name to strictly match an
	// already-registered managed resource and demands ResourceIdentity
	// plus an `iter.Seq[ListResult]` iterator; a full integration is
	// beyond the scope of this change and will be carried out in a
	// separate follow-up change.
	return out
}

// frameworkActions gathers every framework-side action factory.
func frameworkActions() []func() action.Action {
	out := make([]func() action.Action, 0)
	out = append(out, cvmactions.NewRebootInstanceAction)
	return out
}

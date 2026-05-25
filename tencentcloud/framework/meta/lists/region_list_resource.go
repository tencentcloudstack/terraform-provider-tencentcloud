// Package metalists is a placeholder package for the framework provider's
// "meta" product ListResource.
//
// **Important**: this package is currently an L0-tier reference
// placeholder — it only provides a static region list helper
// (regionEntries). It does **not** implement the framework
// list.ListResource interface and is **not** wired into
// frameworkListResources() in tencentcloud/framework/registry.go.
//
// Reasons:
//  1. The framework v1.19 list.ListResource interface requires the list's
//     type name to strictly match an **already-registered managed
//     resource** type name (see the Metadata documentation in
//     vendor/.../list/list_resource.go), i.e. before wiring up a list this
//     repo first needs a tencentcloud_region resource.
//  2. The List interface also demands ListResourceConfigSchema and
//     List(ctx, ListRequest, *ListResultsStream); the latter has to build
//     a *tfsdk.ResourceIdentity, which depends on the resource identity
//     schema introduced in framework v1.19.
//  3. The scope of this change (restructure-framework-types-and-naming) is
//     "naming cleanup + flattened directory layout + product hierarchy +
//     6 type references". Implementing ResourceIdentity together with the
//     companion region resource is beyond that scope, so we ship the L0
//     placeholder first and let the directory shape and file naming come
//     into existence.
//
// When the full list reference is added later, please:
//  1. Implement the framework managed resource tencentcloud_region first
//     (with IdentitySchema).
//  2. Then implement list.ListResource with a type name strictly identical
//     to that of the managed resource.
//  3. Append the factory to frameworkListResources() in
//     tencentcloud/framework/registry.go.
package metalists

// regionEntry is a single placeholder region row.
//
// When the actual list reference is wired up, this type can be reused as
// the input row for the listed managed resource when constructing
// *tfsdk.ResourceIdentity inside its List method.
type regionEntry struct {
	ID          string
	Name        string
	Description string
}

// regionEntries is the static placeholder region list. It covers at least
// five common regions and the data is kept compatible with the
// TencentCloud region constants in the connectivity package.
var regionEntries = []regionEntry{
	{ID: "ap-guangzhou", Name: "Guangzhou", Description: "South China region (Guangzhou)."},
	{ID: "ap-shanghai", Name: "Shanghai", Description: "East China region (Shanghai)."},
	{ID: "ap-beijing", Name: "Beijing", Description: "North China region (Beijing)."},
	{ID: "ap-chengdu", Name: "Chengdu", Description: "Southwest China region (Chengdu)."},
	{ID: "ap-hongkong", Name: "Hong Kong", Description: "Hong Kong region."},
	{ID: "ap-singapore", Name: "Singapore", Description: "Southeast Asia region (Singapore)."},
	{ID: "na-siliconvalley", Name: "Silicon Valley", Description: "North America region (Silicon Valley)."},
}

// RegionEntries returns a copy of regionEntries for use by a future real
// list.ListResource implementation. Returning a copy rather than the
// underlying slice prevents external code from accidentally mutating the
// data.
func RegionEntries() []regionEntry {
	out := make([]regionEntry, len(regionEntries))
	copy(out, regionEntries)
	return out
}

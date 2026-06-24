// registry.go is the central, SDKv2-style manifest of every
// terraform-plugin-framework reference shipped by this provider.
//
// Layout mirrors tencentcloud/provider.go's ResourcesMap / DataSourcesMap:
// each framework reference type has its own append-only block of one
// entry per line. To add a new reference, edit ONLY this file by adding
// the corresponding services subpackage import (if not already present)
// and a single new entry in the matching block. No edit to provider.go
// is required.
package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// resourceFactories lists every framework Resource factory.
var resourceFactories = []func() resource.Resource{}

// dataSourceFactories lists every framework DataSource factory.
var dataSourceFactories = []func() datasource.DataSource{}

// functionFactories lists every framework Function factory.
var functionFactories = []func() function.Function{}

// ephemeralResourceFactories lists every framework EphemeralResource factory.
var ephemeralResourceFactories = []func() ephemeral.EphemeralResource{}

// listResourceFactories lists every framework ListResource factory.
var listResourceFactories = []func() list.ListResource{}

// actionFactories lists every framework Action factory.
var actionFactories = []func() action.Action{}

// frameworkResources returns every framework Resource factory.
func frameworkResources() []func() resource.Resource { return resourceFactories }

// frameworkDataSources returns every framework DataSource factory.
func frameworkDataSources() []func() datasource.DataSource { return dataSourceFactories }

// frameworkFunctions returns every framework Function factory.
func frameworkFunctions() []func() function.Function { return functionFactories }

// frameworkEphemeralResources returns every framework EphemeralResource factory.
func frameworkEphemeralResources() []func() ephemeral.EphemeralResource {
	return ephemeralResourceFactories
}

// frameworkListResources returns every framework ListResource factory.
func frameworkListResources() []func() list.ListResource { return listResourceFactories }

// frameworkActions returns every framework Action factory.
func frameworkActions() []func() action.Action { return actionFactories }

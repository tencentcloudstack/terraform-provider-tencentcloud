// Package sharedmeta provides framework-side resources/data sources with the
// same runtime context as SDKv2 (credentials and a shared SDK client).
//
// This package never parses credentials. Credentials are parsed exclusively
// by the SDKv2 provider's ConfigureContextFunc and then injected via
// SetSharedMeta; this package only carries and holds the value.
package sharedmeta

import (
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

// ProviderMeta is the concrete type of req.ProviderData that framework-side
// resources/data sources receive in their Configure phase. Each framework
// resource should type-assert to *ProviderMeta in Configure to obtain Client.
//
// The field set is deliberately minimal. Future extensions (e.g. region
// override or module name) can be appended here without affecting existing
// framework resources.
type ProviderMeta struct {
	// Client is the same TencentCloud SDK client instance shared with the
	// SDKv2 provider; it owns credentials, SDK client cache, UA, and retry
	// behaviour.
	Client *connectivity.TencentCloudClient
}

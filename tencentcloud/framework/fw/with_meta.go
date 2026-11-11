package fw

import (
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

// withMeta holds the shared TencentCloud SDK client injected by the
// framework provider during Configure. It is embedded by the various
// *WithConfigure helper types so that every framework reference can access
// the same client via the Client() accessor.
type withMeta struct {
	client *connectivity.TencentCloudClient
}

// Client returns the shared TencentCloud SDK client.
func (w *withMeta) Client() *connectivity.TencentCloudClient {
	return w.client
}

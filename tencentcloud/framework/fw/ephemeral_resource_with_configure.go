package fw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// EphemeralResourceWithConfigure is a structure to be embedded within an
// EphemeralResource that implements the ephemeral.EphemeralResourceWithConfigure
// interface. It stores the shared TencentCloud SDK client so the ephemeral
// resource can use it during Open/Close.
type EphemeralResourceWithConfigure struct {
	withMeta
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure and stores the shared client.
func (e *EphemeralResourceWithConfigure) Configure(_ context.Context, request ephemeral.ConfigureRequest, _ *ephemeral.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if v, ok := request.ProviderData.(*sharedmeta.ProviderMeta); ok {
		e.client = v.Client
	}
}

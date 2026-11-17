package fw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// ActionWithConfigure is a structure to be embedded within an Action that
// implements the action.ActionWithConfigure interface. It stores the shared
// TencentCloud SDK client so the action can use it during Invoke.
type ActionWithConfigure struct {
	withMeta
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure and stores the shared client.
func (a *ActionWithConfigure) Configure(_ context.Context, request action.ConfigureRequest, _ *action.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if v, ok := request.ProviderData.(*sharedmeta.ProviderMeta); ok {
		a.client = v.Client
	}
}

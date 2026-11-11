package fw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// ResourceWithConfigure is a structure to be embedded within a Resource that
// implements the resource.ResourceWithConfigure interface. It stores the
// shared TencentCloud SDK client so the resource can use it during CRUD.
type ResourceWithConfigure struct {
	withMeta
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure and stores the shared client.
func (r *ResourceWithConfigure) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if v, ok := request.ProviderData.(*sharedmeta.ProviderMeta); ok {
		r.client = v.Client
	}
}

package fw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// DataSourceWithConfigure is a structure to be embedded within a DataSource
// that implements the datasource.DataSourceWithConfigure interface. It stores
// the shared TencentCloud SDK client so the data source can use it during Read.
type DataSourceWithConfigure struct {
	withMeta
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure and stores the shared client.
func (d *DataSourceWithConfigure) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if v, ok := request.ProviderData.(*sharedmeta.ProviderMeta); ok {
		d.client = v.Client
	}
}

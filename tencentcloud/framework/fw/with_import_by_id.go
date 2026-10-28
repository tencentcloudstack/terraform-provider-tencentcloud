package fw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// WithImportByID is intended to be embedded in resources which import state via the "id" attribute.
// See https://developer.hashicorp.com/terraform/plugin/framework/resources/import.
type WithImportByID struct{}

func (w *WithImportByID) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

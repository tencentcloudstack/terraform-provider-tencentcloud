package fw

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OverwriteString overwrites target only when the cloud side reported a
// value; a nil value keeps the existing (configured or state) value.
//
// Cloud describe APIs commonly mark response fields as "may return null".
// Unconditionally copying such a field into a model overwrites a configured
// value with null and fails Terraform's post-apply consistency check with
// "was <value>, but now null". Framework resources should use this helper in
// their refresh/flatten helpers for every cloud-reported string field.
func OverwriteString(target *types.String, value *string) {
	if value != nil {
		*target = types.StringValue(*value)
	}
}

package helper

import (
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TimeoutsBlock returns the generic timeouts block declaration used by
// framework resource schemas. It mirrors the user-facing semantics of the
// SDKv2 helper/schema Timeouts block:
//
//	resource "tencentcloud_xxx" "demo" {
//	  timeouts {
//	    create = "30m"
//	    read   = "10m"
//	    update = "30m"
//	    delete = "30m"
//	  }
//	}
//
// We deliberately avoid pulling in terraform-plugin-framework-timeouts as an
// external dependency to keep the vendor graph small (the project is
// sensitive to vendoring). All fields are Optional.
//
// stages may be set to any of "create" / "read" / "update" / "delete"
// (lowercase). For instance, a pure data source only needs "read" while a
// resource typically uses create+read+update+delete. Stages that are not
// listed will not appear in the schema.
func TimeoutsBlock(stages ...string) schema.Block {
	if len(stages) == 0 {
		stages = []string{"create", "read", "update", "delete"}
	}

	attrs := make(map[string]schema.Attribute, len(stages))
	for _, s := range stages {
		attrs[s] = schema.StringAttribute{
			Optional: true,
			Description: "A string that can be " +
				"[parsed as a duration](https://pkg.go.dev/time#ParseDuration). " +
				"Default is the resource's built-in timeout.",
		}
	}
	return schema.SingleNestedBlock{
		Attributes:  attrs,
		Description: "The timeouts block allows the user to specify timeouts for create/read/update/delete operations.",
	}
}

// TimeoutsModel is the Go counterpart of TimeoutsBlock. Resource models
// should embed it with the tag `tfsdk:"timeouts"`. Inside CRUD methods, use
// ParseTimeout to recover the time.Duration for a given stage.
type TimeoutsModel struct {
	Create types.String `tfsdk:"create"`
	Read   types.String `tfsdk:"read"`
	Update types.String `tfsdk:"update"`
	Delete types.String `tfsdk:"delete"`
}

// ErrTimeoutNotSet indicates that the user did not explicitly provide a
// timeout for the requested stage; callers should fall back to the
// resource's built-in default. (time.ParseDuration never returns this
// error.)
var ErrTimeoutNotSet = errors.New("timeout not set")

// ParseTimeout parses a duration string from a timeouts block field.
//
// When the input is Null/Unknown or an empty string the function returns
// (0, ErrTimeoutNotSet) and the caller is expected to substitute the
// resource's built-in default. An invalid string surfaces the error from
// ParseDuration directly.
func ParseTimeout(v types.String) (time.Duration, error) {
	if v.IsNull() || v.IsUnknown() {
		return 0, ErrTimeoutNotSet
	}
	s := v.ValueString()
	if s == "" {
		return 0, ErrTimeoutNotSet
	}
	return time.ParseDuration(s)
}

// TimeoutOrDefault falls back to def when ParseTimeout fails or no value is
// configured, suitable for one-line use at the entry point of CRUD methods:
//
//	timeout := helper.TimeoutOrDefault(plan.Timeouts.Create, 30*time.Minute)
func TimeoutOrDefault(v types.String, def time.Duration) time.Duration {
	d, err := ParseTimeout(v)
	if err != nil {
		return def
	}
	return d
}

// Package helper is the shared utility package for resources/data sources
// on the terraform-plugin-framework side. It centralises common type
// conversions, retries, error normalisation and the timeouts block
// declaration so that individual resources do not re-implement them and
// drift apart in behaviour.
//
// This package lives under `tencentcloud/framework/internal/`, where Go's
// `internal/` visibility rule limits imports to the `tencentcloud/framework/...`
// subtree. It is path-isolated from the project's existing SDKv2-facing
// `tencentcloud/internal/helper` (which depends on many SDKv2 types) so that
// neither symbols nor semantics can collide. The package depends only on
// framework standard types and the TencentCloud SDK error types.
package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringValueOrNull safely converts a Go string to the framework's
// types.String:
//   - An empty string is treated as "value not present on the business
//     side" and is returned as types.StringNull().
//   - A non-empty string is wrapped as types.StringValue.
//
// In a resource Read, the common "the API did not return this field"
// scenario should produce types.StringNull() rather than an empty value;
// otherwise the plan comparison generates noisy diffs against
// Optional+Computed fields. If the empty string is itself a legal business
// value, do not use this helper — construct types.StringValue("")
// directly.
func StringValueOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// StringPointerValueOrNull converts a *string to types.String:
//   - A nil pointer returns types.StringNull().
//   - A non-nil pointer is preserved as types.StringValue("") even when the
//     pointee is the empty string, because the very presence of the pointer
//     means the API explicitly returned the field.
func StringPointerValueOrNull(p *string) types.String {
	if p == nil {
		return types.StringNull()
	}
	return types.StringValue(*p)
}

// Int64ValueOrNull converts an *int64 to types.Int64. Unlike the SDKv2
// world, in framework a 0 is often a legal value, so this helper only
// performs nil handling at the pointer level.
func Int64ValueOrNull(p *int64) types.Int64 {
	if p == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*p)
}

// Int64ValueFromUint safely converts the *uint64 fields commonly seen in
// the TencentCloud SDK to types.Int64.
//
// Note: when the source value exceeds math.MaxInt64 the conversion loses
// precision; however, business fields in TencentCloud (page size, quota,
// capacity, etc.) are not expected to surpass that ceiling for the
// foreseeable future, so this helper does not perform an overflow check.
// Callers must judge for themselves.
func Int64ValueFromUint(p *uint64) types.Int64 {
	if p == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*p))
}

// BoolValueOrNull converts a *bool to types.Bool.
func BoolValueOrNull(p *bool) types.Bool {
	if p == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*p)
}

// StringPointerFromValue converts a types.String back to *string:
//   - Null/Unknown returns nil (which usually means "do not send this
//     field" when invoking the API).
//   - Otherwise it returns a pointer to the string value.
//
// SDK request structs use *string heavily to express "optional, omitted".
// This helper keeps Create/Update code paths concise.
func StringPointerFromValue(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

// Int64PointerFromValue converts a types.Int64 back to *int64.
func Int64PointerFromValue(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	n := v.ValueInt64()
	return &n
}

// BoolPointerFromValue converts a types.Bool back to *bool.
func BoolPointerFromValue(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

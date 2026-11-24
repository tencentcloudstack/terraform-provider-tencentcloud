package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// This file converts between the SDK's pointer-element slices ([]*string,
// []*int64, []*uint64, []*float64, []*bool — by far the most common
// collection shape of the TencentCloud SDK) and the framework's types.List /
// types.Set. One generic pair per direction and collection kind covers every
// scalar element type; the element type E is the Go type the SDK uses
// (uint64 for a []*uint64 field) and elemType is the matching framework type
// (types.Int64Type for an Int64 attribute).
//
// Two directions, two different "nothing here" meanings:
//
//   - WRITE (state -> API): null/unknown yields nil, which means "do not send
//     this field". A known empty collection yields an empty non-nil slice.
//   - READ (API -> state): a nil slice means the cloud omitted the field; a
//     non-nil empty slice is a truthful "empty".

// elements is the read surface shared by types.List and types.Set.
type elements interface {
	IsNull() bool
	IsUnknown() bool
	ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics
}

// pointersFromElements expands a List/Set into a pointer slice.
func pointersFromElements[E any](ctx context.Context, v elements, diags *diag.Diagnostics) []*E {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	var raw []*E

	d := v.ElementsAs(ctx, &raw, false)
	diags.Append(d...)

	if d.HasError() {
		return nil
	}

	out := make([]*E, 0, len(raw))

	for _, p := range raw {
		if p != nil {
			out = append(out, p)
		}
	}

	return out
}

// valuesFromPointers dereferences a pointer slice, dropping nil elements.
func valuesFromPointers[E any](cloud []*E) []E {
	out := make([]E, 0, len(cloud))

	for _, p := range cloud {
		if p != nil {
			out = append(out, *p)
		}
	}

	return out
}

// PointersFromList converts a types.List into the []*E an SDK request field
// expects (WRITE direction).
//   - null / unknown list: nil ("do not send this field");
//   - known list: a non-nil slice (empty for an empty list);
//   - null elements are skipped;
//   - an unknown element, or an element that cannot be represented as E (for
//     example a negative number for E = uint64), is reported as an error
//     diagnostic and nil is returned — never silently truncated.
//
// E is not inferable: write PointersFromList[string](ctx, plan.Zones, diags).
func PointersFromList[E any](ctx context.Context, v types.List, diags *diag.Diagnostics) []*E {
	return pointersFromElements[E](ctx, v, diags)
}

// PointersFromSet is PointersFromList for a types.Set.
func PointersFromSet[E any](ctx context.Context, v types.Set, diags *diag.Diagnostics) []*E {
	return pointersFromElements[E](ctx, v, diags)
}

// ListFromPointers converts a []*E read from the API into a types.List whose
// elements have the framework type elemType (READ direction, for attributes
// whose truthful value may be absent: Computed-only attributes).
//   - nil slice: null list (the cloud omitted the field);
//   - non-nil slice: a list holding the non-nil elements (empty stays an
//     EMPTY list, not null).
func ListFromPointers[E any](ctx context.Context, elemType attr.Type, cloud []*E, diags *diag.Diagnostics) types.List {
	if cloud == nil {
		return types.ListNull(elemType)
	}

	list, d := types.ListValueFrom(ctx, elemType, valuesFromPointers(cloud))
	diags.Append(d...)

	return list
}

// SetFromPointers is ListFromPointers for a types.Set.
func SetFromPointers[E any](ctx context.Context, elemType attr.Type, cloud []*E, diags *diag.Diagnostics) types.Set {
	if cloud == nil {
		return types.SetNull(elemType)
	}

	set, d := types.SetValueFrom(ctx, elemType, valuesFromPointers(cloud))
	diags.Append(d...)

	return set
}

// ListFromPointersOrPrev is ListFromPointers for USER-CONFIGURABLE attributes
// (Required, Optional, Optional+Computed): it keeps the previous plan/state
// value where taking the cloud answer verbatim would make Terraform fail the
// post-apply consistency check.
//   - nil slice (cloud omitted the field): prev; an unknown prev becomes null
//     because state must never hold unknown values;
//   - empty answer while prev is null/unknown (the practitioner never set the
//     attribute): null, not an empty list ("was null, but now []");
//   - otherwise the cloud answer, including a truthful empty list when prev
//     held elements (that is real drift and must stay visible).
func ListFromPointersOrPrev[E any](ctx context.Context, elemType attr.Type, cloud []*E, prev types.List, diags *diag.Diagnostics) types.List {
	prevUnset := prev.IsNull() || prev.IsUnknown()

	if cloud == nil {
		if prevUnset {
			return types.ListNull(elemType)
		}

		return prev
	}

	vals := valuesFromPointers(cloud)
	if len(vals) == 0 && prevUnset {
		return types.ListNull(elemType)
	}

	list, d := types.ListValueFrom(ctx, elemType, vals)
	diags.Append(d...)

	return list
}

// SetFromPointersOrPrev is ListFromPointersOrPrev for a types.Set.
func SetFromPointersOrPrev[E any](ctx context.Context, elemType attr.Type, cloud []*E, prev types.Set, diags *diag.Diagnostics) types.Set {
	prevUnset := prev.IsNull() || prev.IsUnknown()

	if cloud == nil {
		if prevUnset {
			return types.SetNull(elemType)
		}

		return prev
	}

	vals := valuesFromPointers(cloud)
	if len(vals) == 0 && prevUnset {
		return types.SetNull(elemType)
	}

	set, d := types.SetValueFrom(ctx, elemType, vals)
	diags.Append(d...)

	return set
}

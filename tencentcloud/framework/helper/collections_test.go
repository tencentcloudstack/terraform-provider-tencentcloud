package helper

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func strList(vals ...string) types.List {
	elems := make([]attr.Value, 0, len(vals))
	for _, v := range vals {
		elems = append(elems, types.StringValue(v))
	}

	return types.ListValueMust(types.StringType, elems)
}

func strSet(vals ...string) types.Set {
	elems := make([]attr.Value, 0, len(vals))
	for _, v := range vals {
		elems = append(elems, types.StringValue(v))
	}

	return types.SetValueMust(types.StringType, elems)
}

func intList(vals ...int64) types.List {
	elems := make([]attr.Value, 0, len(vals))
	for _, v := range vals {
		elems = append(elems, types.Int64Value(v))
	}

	return types.ListValueMust(types.Int64Type, elems)
}

func ptrs[E any](vals ...E) []*E {
	out := make([]*E, 0, len(vals))
	for i := range vals {
		out = append(out, &vals[i])
	}

	return out
}

func derefAll[E comparable](t *testing.T, got []*E, want []E) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("got %d elements, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] == nil || *got[i] != want[i] {
			t.Fatalf("element %d = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestPointersFromList(t *testing.T) {
	ctx := context.Background()

	t.Run("null yields nil (do not send)", func(t *testing.T) {
		var diags diag.Diagnostics
		if got := PointersFromList[string](ctx, types.ListNull(types.StringType), &diags); got != nil || diags.HasError() {
			t.Fatalf("got %v, diags %v", got, diags)
		}
	})

	t.Run("unknown yields nil (do not send)", func(t *testing.T) {
		var diags diag.Diagnostics
		if got := PointersFromList[string](ctx, types.ListUnknown(types.StringType), &diags); got != nil || diags.HasError() {
			t.Fatalf("got %v, diags %v", got, diags)
		}
	})

	t.Run("known empty yields an empty NON-nil slice", func(t *testing.T) {
		var diags diag.Diagnostics

		got := PointersFromList[string](ctx, strList(), &diags)
		if got == nil || len(got) != 0 || diags.HasError() {
			t.Fatalf("got %#v, diags %v", got, diags)
		}
	})

	t.Run("strings keep order", func(t *testing.T) {
		var diags diag.Diagnostics

		derefAll(t, PointersFromList[string](ctx, strList("a", "b", "c"), &diags), []string{"a", "b", "c"})
	})

	t.Run("null elements are skipped", func(t *testing.T) {
		var diags diag.Diagnostics

		in := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringNull(), types.StringValue("b")})
		derefAll(t, PointersFromList[string](ctx, in, &diags), []string{"a", "b"})

		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
	})

	t.Run("unknown element is an error, never truncated", func(t *testing.T) {
		var diags diag.Diagnostics

		in := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringUnknown()})
		if got := PointersFromList[string](ctx, in, &diags); got != nil || !diags.HasError() {
			t.Fatalf("got %v, diags %v: want nil + error", got, diags)
		}
	})

	t.Run("int64 elements", func(t *testing.T) {
		var diags diag.Diagnostics

		derefAll(t, PointersFromList[int64](ctx, intList(1, 2, 3), &diags), []int64{1, 2, 3})
	})

	t.Run("uint64 elements from an Int64 attribute", func(t *testing.T) {
		var diags diag.Diagnostics

		derefAll(t, PointersFromList[uint64](ctx, intList(1, 2), &diags), []uint64{1, 2})

		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
	})

	t.Run("negative number for uint64 is an error, never wrapped", func(t *testing.T) {
		var diags diag.Diagnostics

		if got := PointersFromList[uint64](ctx, intList(1, -1), &diags); got != nil || !diags.HasError() {
			t.Fatalf("got %v, diags %v: want nil + error", got, diags)
		}
	})

	t.Run("float64 and bool elements", func(t *testing.T) {
		var diags diag.Diagnostics

		fl := types.ListValueMust(types.Float64Type, []attr.Value{types.Float64Value(1.5), types.Float64Value(2)})
		derefAll(t, PointersFromList[float64](ctx, fl, &diags), []float64{1.5, 2})

		bl := types.ListValueMust(types.BoolType, []attr.Value{types.BoolValue(true), types.BoolValue(false)})
		derefAll(t, PointersFromList[bool](ctx, bl, &diags), []bool{true, false})
	})

	t.Run("wrong element type is an error diagnostic, not a panic", func(t *testing.T) {
		var diags diag.Diagnostics

		if got := PointersFromList[int64](ctx, strList("a"), &diags); got != nil || !diags.HasError() {
			t.Fatalf("got %v, diags %v: want nil + error", got, diags)
		}
	})
}

func TestPointersFromSet(t *testing.T) {
	ctx := context.Background()

	var diags diag.Diagnostics

	if got := PointersFromSet[string](ctx, types.SetNull(types.StringType), &diags); got != nil {
		t.Fatalf("null set: got %v, want nil", got)
	}

	if got := PointersFromSet[string](ctx, types.SetUnknown(types.StringType), &diags); got != nil {
		t.Fatalf("unknown set: got %v, want nil", got)
	}

	if got := PointersFromSet[string](ctx, strSet(), &diags); got == nil || len(got) != 0 {
		t.Fatalf("empty set: got %#v, want empty non-nil", got)
	}

	got := PointersFromSet[string](ctx, strSet("x"), &diags)
	derefAll(t, got, []string{"x"})

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
}

func TestListFromPointers(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name  string
		cloud []*string
		want  types.List
	}{
		{"nil is null", nil, types.ListNull(types.StringType)},
		{"empty stays an EMPTY list", []*string{}, strList()},
		{"values", ptrs("a", "b"), strList("a", "b")},
		{"nil elements dropped", []*string{nil}, strList()},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := ListFromPointers(ctx, types.StringType, c.cloud, &diags)
			if !got.Equal(c.want) || diags.HasError() {
				t.Fatalf("got %s, want %s, diags %v", got, c.want, diags)
			}
		})
	}

	t.Run("uint64 into an Int64 list", func(t *testing.T) {
		var diags diag.Diagnostics

		got := ListFromPointers(ctx, types.Int64Type, ptrs[uint64](7, 9), &diags)
		if !got.Equal(intList(7, 9)) || diags.HasError() {
			t.Fatalf("got %s, diags %v", got, diags)
		}
	})

	t.Run("float64 and bool", func(t *testing.T) {
		var diags diag.Diagnostics

		got := ListFromPointers(ctx, types.Float64Type, ptrs(1.5), &diags)
		if !got.Equal(types.ListValueMust(types.Float64Type, []attr.Value{types.Float64Value(1.5)})) {
			t.Fatalf("float list: got %s", got)
		}

		gotB := ListFromPointers(ctx, types.BoolType, ptrs(true), &diags)
		if !gotB.Equal(types.ListValueMust(types.BoolType, []attr.Value{types.BoolValue(true)})) {
			t.Fatalf("bool list: got %s", gotB)
		}
	})

	t.Run("element type mismatch is an error diagnostic", func(t *testing.T) {
		var diags diag.Diagnostics

		_ = ListFromPointers(ctx, types.Int64Type, ptrs("a"), &diags)
		if !diags.HasError() {
			t.Fatal("expected an error diagnostic")
		}
	})
}

func TestSetFromPointers(t *testing.T) {
	ctx := context.Background()

	var diags diag.Diagnostics

	if got := SetFromPointers[string](ctx, types.StringType, nil, &diags); !got.IsNull() {
		t.Fatalf("nil: got %s, want null", got)
	}

	if got := SetFromPointers(ctx, types.StringType, []*string{}, &diags); got.IsNull() || len(got.Elements()) != 0 {
		t.Fatalf("empty: got %s, want an empty set", got)
	}

	if got := SetFromPointers(ctx, types.StringType, ptrs("a", "b"), &diags); !got.Equal(strSet("a", "b")) {
		t.Fatalf("values: got %s", got)
	}

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
}

func TestListFromPointersOrPrev(t *testing.T) {
	ctx := context.Background()
	st := types.StringType

	cases := []struct {
		name  string
		cloud []*string
		prev  types.List
		want  types.List
	}{
		{"cloud values win", ptrs("c"), strList("p"), strList("c")},
		{"cloud values win over null prev", ptrs("c"), types.ListNull(st), strList("c")},
		{"cloud omitted keeps prev", nil, strList("p"), strList("p")},
		{"cloud omitted, null prev stays null", nil, types.ListNull(st), types.ListNull(st)},
		{"cloud omitted, unknown prev becomes null", nil, types.ListUnknown(st), types.ListNull(st)},
		{"cloud empty, never configured: null (no was-null-now-[])", []*string{}, types.ListNull(st), types.ListNull(st)},
		{"cloud empty, unknown prev: null", []*string{}, types.ListUnknown(st), types.ListNull(st)},
		{"cloud empty but prev had elements: real drift stays visible", []*string{}, strList("p"), strList()},
		{"cloud empty, prev known empty", []*string{}, strList(), strList()},
		{"only nil elements counts as empty", []*string{nil}, types.ListNull(st), types.ListNull(st)},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := ListFromPointersOrPrev(ctx, st, c.cloud, c.prev, &diags)
			if !got.Equal(c.want) || diags.HasError() {
				t.Fatalf("got %s, want %s, diags %v", got, c.want, diags)
			}
		})
	}

	t.Run("int64 elements", func(t *testing.T) {
		var diags diag.Diagnostics

		got := ListFromPointersOrPrev(ctx, types.Int64Type, ptrs[int64](5), types.ListNull(types.Int64Type), &diags)
		if !got.Equal(intList(5)) {
			t.Fatalf("got %s", got)
		}
	})
}

func TestSetFromPointersOrPrev(t *testing.T) {
	ctx := context.Background()
	st := types.StringType

	cases := []struct {
		name  string
		cloud []*string
		prev  types.Set
		want  types.Set
	}{
		{"cloud values win", ptrs("c"), strSet("p"), strSet("c")},
		{"cloud omitted keeps prev", nil, strSet("p"), strSet("p")},
		{"cloud omitted, unknown prev becomes null", nil, types.SetUnknown(st), types.SetNull(st)},
		{"cloud empty, never configured: null", []*string{}, types.SetNull(st), types.SetNull(st)},
		{"cloud empty but prev had elements: drift visible", []*string{}, strSet("p"), strSet()},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := SetFromPointersOrPrev(ctx, st, c.cloud, c.prev, &diags)
			if !got.Equal(c.want) || diags.HasError() {
				t.Fatalf("got %s, want %s, diags %v", got, c.want, diags)
			}
		})
	}
}

func TestFloat64PointerFromValue(t *testing.T) {
	if got := Float64PointerFromValue(types.Float64Null()); got != nil {
		t.Fatalf("null: got %v, want nil", got)
	}

	if got := Float64PointerFromValue(types.Float64Unknown()); got != nil {
		t.Fatalf("unknown: got %v, want nil", got)
	}

	if got := Float64PointerFromValue(types.Float64Value(0)); got == nil || *got != 0 {
		t.Fatalf("zero is a real value and must be sent: got %v", got)
	}

	if got := Float64PointerFromValue(types.Float64Value(2.5)); got == nil || *got != 2.5 {
		t.Fatalf("got %v, want 2.5", got)
	}
}

func TestFloat64ValueOrPrev(t *testing.T) {
	one, zero := 1.5, 0.0

	cases := []struct {
		name  string
		cloud *float64
		prev  types.Float64
		want  types.Float64
	}{
		{"cloud value wins", &one, types.Float64Value(9), types.Float64Value(1.5)},
		{"cloud zero is a real value", &zero, types.Float64Value(9), types.Float64Value(0)},
		{"cloud omitted keeps prev", nil, types.Float64Value(9), types.Float64Value(9)},
		{"cloud omitted, null prev stays null", nil, types.Float64Null(), types.Float64Null()},
		{"cloud omitted, unknown prev becomes null", nil, types.Float64Unknown(), types.Float64Null()},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			if got := Float64ValueOrPrev(c.cloud, c.prev); !got.Equal(c.want) {
				t.Fatalf("got %s, want %s", got, c.want)
			}
		})
	}
}

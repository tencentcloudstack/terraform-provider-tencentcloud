package helper

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringValueOrNull(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want types.String
	}{
		{"empty becomes null", "", types.StringNull()},
		{"non-empty wraps", "hello", types.StringValue("hello")},
		{"whitespace preserved", "  ", types.StringValue("  ")},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := StringValueOrNull(c.in)
			if !got.Equal(c.want) {
				t.Fatalf("got %s, want %s", got, c.want)
			}
		})
	}
}

func TestStringPointerValueOrNull(t *testing.T) {
	empty := ""
	hello := "hello"

	cases := []struct {
		name string
		in   *string
		want types.String
	}{
		{"nil becomes null", nil, types.StringNull()},
		{"pointer to empty stays empty", &empty, types.StringValue("")},
		{"pointer to value wraps", &hello, types.StringValue("hello")},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := StringPointerValueOrNull(c.in)
			if !got.Equal(c.want) {
				t.Fatalf("got %s, want %s", got, c.want)
			}
		})
	}
}

func TestInt64ValueOrNull(t *testing.T) {
	zero := int64(0)
	pos := int64(42)
	cases := []struct {
		name string
		in   *int64
		want types.Int64
	}{
		{"nil null", nil, types.Int64Null()},
		{"zero kept", &zero, types.Int64Value(0)},
		{"positive kept", &pos, types.Int64Value(42)},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := Int64ValueOrNull(c.in)
			if !got.Equal(c.want) {
				t.Fatalf("got %s, want %s", got, c.want)
			}
		})
	}
}

func TestInt64ValueFromUint(t *testing.T) {
	v := uint64(7)
	got := Int64ValueFromUint(&v)
	if got.ValueInt64() != 7 {
		t.Fatalf("expected 7, got %d", got.ValueInt64())
	}
	if !Int64ValueFromUint(nil).IsNull() {
		t.Fatalf("nil pointer should be null")
	}
}

func TestBoolValueOrNull(t *testing.T) {
	tt := true
	ff := false
	if !BoolValueOrNull(nil).IsNull() {
		t.Fatalf("nil should be null")
	}
	if !BoolValueOrNull(&tt).Equal(types.BoolValue(true)) {
		t.Fatalf("true wrapping failed")
	}
	if !BoolValueOrNull(&ff).Equal(types.BoolValue(false)) {
		t.Fatalf("false wrapping failed")
	}
}

func TestStringPointerFromValue(t *testing.T) {
	if StringPointerFromValue(types.StringNull()) != nil {
		t.Fatalf("null should produce nil pointer")
	}
	if StringPointerFromValue(types.StringUnknown()) != nil {
		t.Fatalf("unknown should produce nil pointer")
	}
	got := StringPointerFromValue(types.StringValue("hello"))
	if got == nil || *got != "hello" {
		t.Fatalf("expected pointer to %q, got %v", "hello", got)
	}
}

func TestInt64PointerFromValue(t *testing.T) {
	if Int64PointerFromValue(types.Int64Null()) != nil {
		t.Fatalf("null should produce nil pointer")
	}
	got := Int64PointerFromValue(types.Int64Value(99))
	if got == nil || *got != 99 {
		t.Fatalf("expected pointer to 99, got %v", got)
	}
}

func TestBoolPointerFromValue(t *testing.T) {
	if BoolPointerFromValue(types.BoolNull()) != nil {
		t.Fatalf("null should produce nil pointer")
	}
	got := BoolPointerFromValue(types.BoolValue(true))
	if got == nil || *got != true {
		t.Fatalf("expected pointer to true, got %v", got)
	}
}

func TestUint64PointerFromValue(t *testing.T) {
	zero := uint64(0)
	three := uint64(3)
	maxInt64 := uint64(1<<63 - 1)

	cases := []struct {
		name string
		in   types.Int64
		want *uint64
	}{
		{"null becomes nil", types.Int64Null(), nil},
		{"unknown becomes nil", types.Int64Unknown(), nil},
		{"negative becomes nil", types.Int64Value(-1), nil},
		{"zero is kept", types.Int64Value(0), &zero},
		{"positive is kept", types.Int64Value(3), &three},
		{"int64 max fits uint64", types.Int64Value(1<<63 - 1), &maxInt64},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := Uint64PointerFromValue(c.in)
			if (got == nil) != (c.want == nil) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
			if got != nil && *got != *c.want {
				t.Fatalf("got %d, want %d", *got, *c.want)
			}
		})
	}
}

func TestStringValueOrPrev(t *testing.T) {
	cloud := "from-cloud"

	cases := []struct {
		name string
		in   *string
		prev types.String
		want types.String
	}{
		{"cloud wins", &cloud, types.StringValue("from-prev"), types.StringValue("from-cloud")},
		{"cloud wins over null prev", &cloud, types.StringNull(), types.StringValue("from-cloud")},
		{"nil keeps prev", nil, types.StringValue("from-prev"), types.StringValue("from-prev")},
		{"nil and null prev becomes null", nil, types.StringNull(), types.StringNull()},
		{"unknown prev becomes null", nil, types.StringUnknown(), types.StringNull()},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := StringValueOrPrev(c.in, c.prev); !got.Equal(c.want) {
				t.Fatalf("got %s, want %s", got, c.want)
			}
		})
	}
}

func TestInt64ValueOrPrev(t *testing.T) {
	cloud := uint64(2)

	cases := []struct {
		name string
		in   *uint64
		prev types.Int64
		want types.Int64
	}{
		{"cloud wins", &cloud, types.Int64Value(1), types.Int64Value(2)},
		{"nil keeps prev", nil, types.Int64Value(1), types.Int64Value(1)},
		{"nil and null prev becomes null", nil, types.Int64Null(), types.Int64Null()},
		{"unknown prev becomes null", nil, types.Int64Unknown(), types.Int64Null()},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := Int64ValueOrPrev(c.in, c.prev); !got.Equal(c.want) {
				t.Fatalf("got %s, want %s", got, c.want)
			}
		})
	}
}

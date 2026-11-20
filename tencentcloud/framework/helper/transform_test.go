package helper

import (
	"testing"
	"time"
)

func TestValueToPointers(t *testing.T) {
	if v := String("x"); *v != "x" {
		t.Fatalf("String: got %q", *v)
	}
	if v := Int(7); *v != 7 {
		t.Fatalf("Int: got %d", *v)
	}
	if v := Int64(7); *v != 7 {
		t.Fatalf("Int64: got %d", *v)
	}
	if v := Uint64(7); *v != 7 {
		t.Fatalf("Uint64: got %d", *v)
	}
	if v := Float64(1.5); *v != 1.5 {
		t.Fatalf("Float64: got %f", *v)
	}
	if v := Bool(true); !*v {
		t.Fatal("Bool: got false")
	}
}

func TestStrings(t *testing.T) {
	if got := Strings(nil); got != nil {
		t.Fatalf("nil input must yield nil, got %v", got)
	}
	if got := Strings([]string{}); got != nil {
		t.Fatalf("empty input must yield nil, got %v", got)
	}

	got := Strings([]string{"a", "b"})
	if len(got) != 2 || *got[0] != "a" || *got[1] != "b" {
		t.Fatalf("got %v, want [a b]", got)
	}
}

func TestFormatUnixTime(t *testing.T) {
	got := FormatUnixTime(0)
	if got != "1970-01-01T00:00:00Z" {
		t.Fatalf("got %q", got)
	}
}

func TestParseTime(t *testing.T) {
	want := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	got, err := ParseTime("2026-01-02T03:04:05Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	if _, err := ParseTime("not-a-time"); err == nil {
		t.Fatal("expected error for malformed input")
	}
}

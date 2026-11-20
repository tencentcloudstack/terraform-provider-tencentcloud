package helper

import (
	"strings"
	"testing"
)

func TestIdFormat(t *testing.T) {
	cases := []struct {
		name     string
		in       []string
		expected string
	}{
		{"single segment", []string{"ins-123"}, "ins-123"},
		{"two segments parent-first", []string{"ins-123", "disk-456"}, "ins-123#disk-456"},
		{"three segments", []string{"a", "b", "c"}, "a#b#c"},
		{"no segments", nil, ""},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := IdFormat(c.in...); got != c.expected {
				t.Fatalf("got %q, want %q", got, c.expected)
			}
		})
	}
}

func TestIdParse(t *testing.T) {
	got := IdParse("a#b#c")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("got %v, want [a b c]", got)
	}
}

func TestIdParseErr(t *testing.T) {
	cases := []struct {
		name      string
		id        string
		segments  int
		wantParts []string
		wantErr   bool
	}{
		{"valid two segments", "ins-1#disk-2", 2, []string{"ins-1", "disk-2"}, false},
		{"valid three segments", "a#b#c", 3, []string{"a", "b", "c"}, false},
		{"wrong count", "a#b", 3, nil, true},
		{"empty segment", "a##c", 3, nil, true},
		{"single empty segment", "", 1, nil, true},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			parts, err := IdParseErr(c.id, c.segments)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected error, got parts %v", parts)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.Join(parts, "|") != strings.Join(c.wantParts, "|") {
				t.Fatalf("got %v, want %v", parts, c.wantParts)
			}
		})
	}
}

func TestBuildToken(t *testing.T) {
	a := BuildToken()
	b := BuildToken()
	if a == "" {
		t.Fatal("token must not be empty")
	}
	if a == b {
		t.Fatal("two tokens must differ")
	}
}

func TestBuildUUID(t *testing.T) {
	a := BuildUUID()
	b := BuildUUID()
	if len(a) != 36 {
		t.Fatalf("uuid has unexpected length %d: %q", len(a), a)
	}
	if a == b {
		t.Fatal("two uuids must differ")
	}
}

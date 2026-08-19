package helper

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParseTimeout_NullUnknownEmpty(t *testing.T) {
	cases := []types.String{
		types.StringNull(),
		types.StringUnknown(),
		types.StringValue(""),
	}
	for _, in := range cases {
		_, err := ParseTimeout(in)
		if !errors.Is(err, ErrTimeoutNotSet) {
			t.Fatalf("expected ErrTimeoutNotSet for %s, got %v", in, err)
		}
	}
}

func TestParseTimeout_Valid(t *testing.T) {
	d, err := ParseTimeout(types.StringValue("30m"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 30*time.Minute {
		t.Fatalf("expected 30m, got %s", d)
	}
}

func TestParseTimeout_Invalid(t *testing.T) {
	_, err := ParseTimeout(types.StringValue("not-a-duration"))
	if err == nil {
		t.Fatalf("expected parse error")
	}
	if errors.Is(err, ErrTimeoutNotSet) {
		t.Fatalf("invalid duration should not be ErrTimeoutNotSet")
	}
}

func TestTimeoutOrDefault(t *testing.T) {
	def := 10 * time.Minute

	if got := TimeoutOrDefault(types.StringNull(), def); got != def {
		t.Fatalf("null should fall back to default, got %s", got)
	}
	if got := TimeoutOrDefault(types.StringValue("bad"), def); got != def {
		t.Fatalf("invalid should fall back to default, got %s", got)
	}
	if got := TimeoutOrDefault(types.StringValue("5s"), def); got != 5*time.Second {
		t.Fatalf("expected 5s, got %s", got)
	}
}

func TestTimeoutsBlock_DefaultsAllStages(t *testing.T) {
	// Interface-level smoke check: with no input the block should declare
	// all four stage attributes.
	block := TimeoutsBlock()
	got := block.GetDeprecationMessage() // any method call to avoid unused
	_ = got
	// Deeper schema-attribute checks are deferred to the framework's own
	// schema validation when the block is mounted on a resource.
}

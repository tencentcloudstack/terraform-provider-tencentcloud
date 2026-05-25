package metafunctions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestParseResourceIDFunction_Metadata verifies the function name is
// parse_resource_id.
func TestParseResourceIDFunction_Metadata(t *testing.T) {
	f := NewParseResourceIDFunction()

	var resp function.MetadataResponse
	f.Metadata(context.Background(), function.MetadataRequest{}, &resp)

	if got, want := resp.Name, "parse_resource_id"; got != want {
		t.Fatalf("Metadata.Name = %q, want %q", got, want)
	}
}

// TestParseResourceIDFunction_Definition verifies the function signature:
// two String parameters and one List[String] return.
func TestParseResourceIDFunction_Definition(t *testing.T) {
	f := NewParseResourceIDFunction()

	var resp function.DefinitionResponse
	f.Definition(context.Background(), function.DefinitionRequest{}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Definition returned errors: %v", resp.Diagnostics)
	}

	if n := len(resp.Definition.Parameters); n != 2 {
		t.Fatalf("len(Parameters) = %d, want 2", n)
	}
	if _, ok := resp.Definition.Parameters[0].(function.StringParameter); !ok {
		t.Errorf("Parameters[0] = %T, want function.StringParameter", resp.Definition.Parameters[0])
	}
	if _, ok := resp.Definition.Parameters[1].(function.StringParameter); !ok {
		t.Errorf("Parameters[1] = %T, want function.StringParameter", resp.Definition.Parameters[1])
	}

	listReturn, ok := resp.Definition.Return.(function.ListReturn)
	if !ok {
		t.Fatalf("Return = %T, want function.ListReturn", resp.Definition.Return)
	}
	if listReturn.ElementType != types.StringType {
		t.Errorf("Return.ElementType = %v, want types.StringType", listReturn.ElementType)
	}
}

// TestParseResourceIDFunction_Run_Cases covers typical scenarios:
//   - Splitting "ins-abc#u-xyz" by "#".
//   - A single segment (no separator match).
//   - An empty id string.
func TestParseResourceIDFunction_Run_Cases(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		id   string
		sep  string
		want []string
	}{
		{"normal_split", "ins-abc#u-xyz", "#", []string{"ins-abc", "u-xyz"}},
		{"three_segments", "a-b-c", "-", []string{"a", "b", "c"}},
		{"no_separator_match", "ins-only", "#", []string{"ins-only"}},
		{"empty_id", "", "#", []string{""}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := runParseResourceID(t, tc.id, tc.sep)
			if !equalStringSlice(got, tc.want) {
				t.Fatalf("Run(%q, %q) = %v, want %v", tc.id, tc.sep, got, tc.want)
			}
		})
	}
}

// runParseResourceID is a test helper that builds the framework's standard
// RunRequest / RunResponse, invokes the function implementation, and
// decodes the resulting ListValue into a []string.
func runParseResourceID(t *testing.T, id, sep string) []string {
	t.Helper()

	f := NewParseResourceIDFunction()

	args := function.NewArgumentsData([]attr.Value{
		types.StringValue(id),
		types.StringValue(sep),
	})

	emptyList, diags := types.ListValue(types.StringType, nil)
	if diags.HasError() {
		t.Fatalf("failed to construct empty list value for ResultData: %v", diags)
	}

	resp := function.RunResponse{
		Result: function.NewResultData(emptyList),
	}
	f.Run(context.Background(), function.RunRequest{Arguments: args}, &resp)
	if resp.Error != nil {
		t.Fatalf("Run returned error: %v", resp.Error)
	}

	listValue, ok := resp.Result.Value().(types.List)
	if !ok {
		t.Fatalf("Result is %T, want types.List", resp.Result.Value())
	}

	var out []string
	for _, elem := range listValue.Elements() {
		s, ok := elem.(types.String)
		if !ok {
			t.Fatalf("list element is %T, want types.String", elem)
		}
		out = append(out, s.ValueString())
	}
	return out
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

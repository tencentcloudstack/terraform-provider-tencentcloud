// Package metafunctions provides framework-side provider-defined function
// implementations for the "meta" product family — i.e. pure functions
// that are cross-product or not bound to any specific cloud product.
//
// Currently includes:
//   - parse_resource_id: splits a composite resource id of the form
//     "instanceId#userId" by a user-specified separator and returns the
//     segments as a list of strings. Pure string processing, no IO.
//
// Wiring: append metafunctions.NewParseResourceIDFunction to
// frameworkFunctions() in tencentcloud/framework/registry.go to expose
// the function through the framework provider.
package metafunctions

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Compile-time assertion: parseResourceIDFunction satisfies
// function.Function.
var _ function.Function = &parseResourceIDFunction{}

// NewParseResourceIDFunction is the framework function factory.
func NewParseResourceIDFunction() function.Function {
	return &parseResourceIDFunction{}
}

// parseResourceIDFunction implements a pure string-processing function:
// parse_resource_id(id, sep) -> list[string].
//
// Purpose: serve as the in-repo reference implementation of the framework
// Function type, covering the full Metadata / Definition / Run lifecycle
// so that future business framework functions have a concrete template
// to follow.
type parseResourceIDFunction struct{}

// Metadata sets the function name to parse_resource_id. Users reference
// this function from HCL as:
//
//	provider::tencentcloud::parse_resource_id("ins-x#u-y", "#")
func (f *parseResourceIDFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "parse_resource_id"
}

// Definition describes the function signature: (id, separator) ->
// list[string].
func (f *parseResourceIDFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Split a composite resource id into segments.",
		Description: "Splits an arbitrary composite resource id by a user-specified separator and returns the segments as a list of strings. Performs pure in-memory string processing; calls no TencentCloud API.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "id",
				Description: "The composite resource id, for example \"ins-abc#u-xyz\".",
			},
			function.StringParameter{
				Name:        "separator",
				Description: "The separator to split id by. Typically a single character, but any non-empty string is accepted.",
			},
		},
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

// Run executes the function logic: strings.Split(id, separator).
//
// Behaviour:
//   - When separator is the empty string the function still returns a
//     list of length len(id)+1 (Go's strings.Split splits per character
//     when sep is empty), matching standard Go behaviour.
//   - When id does not contain separator, the function returns a single
//     element ([id]).
//   - No IO and no panics.
func (f *parseResourceIDFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var id, sep string

	if err := req.Arguments.GetArgument(ctx, 0, &id); err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, err)
	}
	if err := req.Arguments.GetArgument(ctx, 1, &sep); err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, err)
	}
	if resp.Error != nil {
		return
	}

	parts := strings.Split(id, sep)

	values := make([]attr.Value, 0, len(parts))
	for _, p := range parts {
		values = append(values, types.StringValue(p))
	}

	listValue, diags := types.ListValue(types.StringType, values)
	resp.Error = function.ConcatFuncErrors(resp.Error, function.FuncErrorFromDiags(ctx, diags))
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, listValue))
}

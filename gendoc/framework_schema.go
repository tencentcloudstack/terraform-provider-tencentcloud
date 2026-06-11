// Package main — framework_schema.go renders terraform-plugin-framework
// schemas to the same markdown shape used by the SDKv2 doc generator
// (see main.go). It is intentionally kept self-contained so it can be
// unit-tested without touching the SDKv2 path.
package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ephschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// fwAttrSpec is a renderer-friendly normalised view of any framework
// attribute (across resource / datasource / ephemeral / action schemas).
type fwAttrSpec struct {
	Name        string
	TypeLabel   string
	Required    bool
	Optional    bool
	Computed    bool
	Sensitive   bool
	ForceNew    bool // synthesised from RequiresReplace plan modifiers when available
	Deprecated  string
	Description string
	// Nested holds the recursively normalised nested attributes (object /
	// list/set/map of nested objects). Empty for primitive attributes.
	Nested []fwAttrSpec
}

// renderResourceSchema renders a resource.Schema (also reused for
// datasource.Schema, ephemeral.Schema and action.Schema by way of the
// schemaShape interface below) into the (arguments, attributes,
// subStruct) string triple consumed by docTPL.
//
// The output format is byte-for-byte compatible with the SDKv2 renderer:
//   - Required and Optional arguments are bullet lines starting with "* ".
//   - Pure Computed attributes go to the "Attributes Reference" section.
//   - Nested objects emit additional "The `xxx` object supports..." blocks.
func renderFrameworkSchema(name string, attrs map[string]fwAttrSpec) (arguments, attributes string) {
	var (
		required []string
		optional []string
		readOnly []string
		nested   []string
	)

	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		a := attrs[k]
		// guard: every attribute MUST carry a non-empty description, just
		// like the SDKv2 generator enforces.
		if a.Description == "" {
			fail(fmt.Sprintf("description for %q is missing in framework schema %q", k, name))
		}
		checkDescription(k, a.Description)

		opt := classify(a)
		desc := a.Description
		if a.Deprecated != "" {
			opt += ", **Deprecated**"
			desc = fmt.Sprintf("%s %s", a.Deprecated, desc)
		}
		if a.ForceNew {
			opt += ", ForceNew"
		}
		line := fmt.Sprintf("* `%s` - (%s) %s", k, opt, desc)

		switch {
		case a.Required:
			required = append(required, line)
			if len(a.Nested) > 0 {
				nested = append(nested, renderNested(k, a.Nested)...)
			}
		case a.Optional:
			optional = append(optional, line)
			if len(a.Nested) > 0 {
				nested = append(nested, renderNested(k, a.Nested)...)
			}
		case a.Computed:
			readOnly = append(readOnly, fmt.Sprintf("* `%s` - %s", k, desc))
			if len(a.Nested) > 0 {
				readOnly = append(readOnly, renderNestedAttrs(1, a.Nested)...)
			}
		}
	}

	sort.Strings(required)
	sort.Strings(optional)
	sort.Strings(readOnly)

	args := append(required, optional...)
	arguments = strings.Join(args, "\n")
	if len(nested) > 0 {
		arguments += "\n" + strings.Join(nested, "\n")
	}
	attributes = strings.Join(readOnly, "\n")
	return arguments, attributes
}

// renderNested produces the "The `foo` object supports the following:"
// section for a nested object attribute.
func renderNested(name string, nested []fwAttrSpec) []string {
	var (
		req []string
		opt []string
	)
	for _, a := range nested {
		if a.Description == "" {
			continue
		}
		c := classify(a)
		desc := a.Description
		if a.Deprecated != "" {
			c += ", **Deprecated**"
			desc = fmt.Sprintf("%s %s", a.Deprecated, desc)
		}
		line := fmt.Sprintf("* `%s` - (%s) %s", a.Name, c, desc)
		if a.Required {
			req = append(req, line)
		} else if a.Optional {
			opt = append(opt, line)
		}
	}
	sort.Strings(req)
	sort.Strings(opt)

	out := []string{fmt.Sprintf("\nThe `%s` object supports the following:\n", name)}
	out = append(out, req...)
	out = append(out, opt...)
	return out
}

// renderNestedAttrs produces indented bullets under a Computed nested
// object (so it lines up under the "Attributes Reference" section).
func renderNestedAttrs(depth int, nested []fwAttrSpec) []string {
	indent := strings.Repeat("  ", depth)
	var out []string
	for _, a := range nested {
		if a.Description == "" {
			continue
		}
		out = append(out, fmt.Sprintf("%s* `%s` - %s", indent, a.Name, a.Description))
		if len(a.Nested) > 0 {
			out = append(out, renderNestedAttrs(depth+1, a.Nested)...)
		}
	}
	return out
}

// classify returns the leading "(Required, String)" / "(Optional, List)"
// fragment for a bullet line.
func classify(a fwAttrSpec) string {
	var role string
	switch {
	case a.Required:
		role = "Required"
	case a.Optional:
		role = "Optional"
	default:
		role = "Computed"
	}
	if a.TypeLabel != "" {
		return role + ", " + a.TypeLabel
	}
	return role
}

// --- schema flattening helpers -----------------------------------------

// flattenResourceSchema converts a resource.Schema into a flat map of
// fwAttrSpec keyed by attribute or block name.
func flattenResourceSchema(s rsschema.Schema) map[string]fwAttrSpec {
	out := map[string]fwAttrSpec{}
	for name, attr := range s.GetAttributes() {
		out[name] = describeAttribute(name, attr)
	}
	for name, blk := range s.GetBlocks() {
		out[name] = describeBlock(name, blk)
	}
	return out
}

// flattenDataSourceSchema converts a datasource.Schema.
func flattenDataSourceSchema(s dsschema.Schema) map[string]fwAttrSpec {
	out := map[string]fwAttrSpec{}
	for name, attr := range s.GetAttributes() {
		out[name] = describeAttribute(name, attr)
	}
	for name, blk := range s.GetBlocks() {
		out[name] = describeBlock(name, blk)
	}
	return out
}

// flattenEphemeralSchema converts an ephemeral.Schema.
func flattenEphemeralSchema(s ephschema.Schema) map[string]fwAttrSpec {
	out := map[string]fwAttrSpec{}
	for name, attr := range s.GetAttributes() {
		out[name] = describeAttribute(name, attr)
	}
	for name, blk := range s.GetBlocks() {
		out[name] = describeBlock(name, blk)
	}
	return out
}

// flattenActionSchema converts an action.Schema.
func flattenActionSchema(s actionschema.Schema) map[string]fwAttrSpec {
	out := map[string]fwAttrSpec{}
	for name, attr := range s.GetAttributes() {
		out[name] = describeAttribute(name, attr)
	}
	for name, blk := range s.GetBlocks() {
		out[name] = describeBlock(name, blk)
	}
	return out
}

// flattenFunctionDefinition converts a function.Definition into a flat
// argument list (parameters become Required attributes, the variadic
// parameter becomes an Optional attribute, the return type becomes the
// sole Computed attribute called "result").
func flattenFunctionDefinition(d function.Definition) map[string]fwAttrSpec {
	out := map[string]fwAttrSpec{}
	for _, p := range d.Parameters {
		name := p.GetName()
		if name == "" {
			continue
		}
		out[name] = fwAttrSpec{
			Name:        name,
			TypeLabel:   typeLabelFromValueType(p.GetType()),
			Required:    true,
			Description: ensureSentence(firstNonEmpty(p.GetMarkdownDescription(), p.GetDescription(), "Function parameter.")),
		}
	}
	if d.VariadicParameter != nil {
		name := d.VariadicParameter.GetName()
		if name == "" {
			name = "variadic"
		}
		out[name] = fwAttrSpec{
			Name:        name,
			TypeLabel:   typeLabelFromValueType(d.VariadicParameter.GetType()) + ", Variadic",
			Optional:    true,
			Description: ensureSentence(firstNonEmpty(d.VariadicParameter.GetMarkdownDescription(), d.VariadicParameter.GetDescription(), "Variadic function parameter.")),
		}
	}
	if d.Return != nil {
		out["result"] = fwAttrSpec{
			Name:        "result",
			TypeLabel:   typeLabelFromValueType(d.Return.GetType()),
			Computed:    true,
			Description: ensureSentence(firstNonEmpty(d.MarkdownDescription, d.Description, d.Summary, "Function result.")),
		}
	}
	return out
}

// describeAttribute uses reflection to extract the common required /
// optional / computed / sensitive / description fields from any
// schema.Attribute implementation across the framework subpackages.
//
// Reflection is used (rather than per-package switches) because the
// concrete attribute types (resource/schema.StringAttribute,
// datasource/schema.StringAttribute, ephemeral/schema.StringAttribute,
// action/schema.StringAttribute, ...) all expose the same set of public
// boolean fields and string fields, but live in different packages and
// implement non-exported fwschema.Attribute interface methods.
func describeAttribute(name string, attr any) fwAttrSpec {
	v := reflect.Indirect(reflect.ValueOf(attr))
	spec := fwAttrSpec{Name: name}

	if v.Kind() != reflect.Struct {
		spec.Description = "Framework attribute."
		return spec
	}

	getBool := func(field string) bool {
		f := v.FieldByName(field)
		if f.IsValid() && f.Kind() == reflect.Bool {
			return f.Bool()
		}
		return false
	}
	getString := func(field string) string {
		f := v.FieldByName(field)
		if f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
		return ""
	}

	spec.Required = getBool("Required")
	spec.Optional = getBool("Optional")
	spec.Computed = getBool("Computed")
	spec.Sensitive = getBool("Sensitive")
	spec.Deprecated = getString("DeprecationMessage")

	desc := getString("MarkdownDescription")
	if desc == "" {
		desc = getString("Description")
	}
	if desc == "" {
		desc = "Framework attribute."
	}
	spec.Description = ensureSentence(desc)

	spec.TypeLabel = typeLabelFromAttribute(attr)

	// Nested attributes (single object / list of object / etc.). Look for
	// well-known field names.
	if nested := nestedAttributesOf(v); len(nested) > 0 {
		spec.Nested = nested
	}

	return spec
}

// describeBlock extracts the same shape from any schema.Block.
func describeBlock(name string, blk any) fwAttrSpec {
	v := reflect.Indirect(reflect.ValueOf(blk))
	spec := fwAttrSpec{
		Name:      name,
		Optional:  true, // blocks are always optional in the framework
		TypeLabel: blockTypeLabel(blk),
	}

	getString := func(field string) string {
		f := v.FieldByName(field)
		if f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
		return ""
	}

	desc := getString("MarkdownDescription")
	if desc == "" {
		desc = getString("Description")
	}
	if desc == "" {
		desc = fmt.Sprintf("The `%s` block.", name)
	}
	spec.Description = ensureSentence(desc)
	spec.Deprecated = getString("DeprecationMessage")

	if nested := nestedAttributesOfBlock(v); len(nested) > 0 {
		spec.Nested = nested
	}
	return spec
}

// nestedAttributesOf inspects an attribute's reflected value for nested
// attribute carriers like `Attributes map[string]ResourceAttribute` (single
// nested object) or `NestedObject.Attributes map[string]...` (list/set/map
// of nested objects).
func nestedAttributesOf(v reflect.Value) []fwAttrSpec {
	if !v.IsValid() {
		return nil
	}
	if f := v.FieldByName("Attributes"); f.IsValid() && f.Kind() == reflect.Map {
		return mapToSpecs(f)
	}
	if f := v.FieldByName("NestedObject"); f.IsValid() {
		return nestedAttributesOf(f)
	}
	return nil
}

// nestedAttributesOfBlock is the analogue for blocks. It inspects the
// `NestedObject` field for `Attributes` and `Blocks`.
func nestedAttributesOfBlock(v reflect.Value) []fwAttrSpec {
	if !v.IsValid() {
		return nil
	}
	no := v.FieldByName("NestedObject")
	if !no.IsValid() {
		return nil
	}
	out := []fwAttrSpec{}
	if attrs := no.FieldByName("Attributes"); attrs.IsValid() && attrs.Kind() == reflect.Map {
		out = append(out, mapToSpecs(attrs)...)
	}
	if blks := no.FieldByName("Blocks"); blks.IsValid() && blks.Kind() == reflect.Map {
		out = append(out, mapToSpecs(blks)...)
	}
	return out
}

// mapToSpecs converts a reflected map[string]Attribute into []fwAttrSpec.
func mapToSpecs(m reflect.Value) []fwAttrSpec {
	keys := m.MapKeys()
	out := make([]fwAttrSpec, 0, len(keys))
	names := make([]string, 0, len(keys))
	for _, k := range keys {
		names = append(names, k.String())
	}
	sort.Strings(names)
	for _, name := range names {
		v := m.MapIndex(reflect.ValueOf(name)).Interface()
		out = append(out, describeAttribute(name, v))
	}
	return out
}

// typeLabelFromAttribute returns a coarse human-readable label for an
// attribute (String / Number / Bool / List / Set / Map / Object).
func typeLabelFromAttribute(attr any) string {
	t := reflect.TypeOf(attr)
	if t == nil {
		return "Object"
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()
	switch {
	case strings.HasPrefix(name, "String"):
		return "String"
	case strings.HasPrefix(name, "Bool"):
		return "Bool"
	case strings.HasPrefix(name, "Int64"), strings.HasPrefix(name, "Int32"):
		return "Int"
	case strings.HasPrefix(name, "Float64"), strings.HasPrefix(name, "Float32"), strings.HasPrefix(name, "Number"):
		return "Number"
	case strings.HasPrefix(name, "List"):
		return "List"
	case strings.HasPrefix(name, "Set"):
		return "Set"
	case strings.HasPrefix(name, "Map"):
		return "Map"
	case strings.HasPrefix(name, "Object"), strings.HasPrefix(name, "SingleNested"):
		return "Object"
	case strings.HasPrefix(name, "Dynamic"):
		return "Dynamic"
	}
	return "Object"
}

// blockTypeLabel returns a coarse label for a framework Block (List / Set
// / Single).
func blockTypeLabel(blk any) string {
	t := reflect.TypeOf(blk)
	if t == nil {
		return "Block"
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()
	switch {
	case strings.HasPrefix(name, "ListNested"):
		return "List of Object"
	case strings.HasPrefix(name, "SetNested"):
		return "Set of Object"
	case strings.HasPrefix(name, "SingleNested"):
		return "Object"
	}
	return "Block"
}

// typeLabelFromValueType maps a function parameter / return value type
// (any concrete attr.Type implementation) into a coarse label.
func typeLabelFromValueType(t any) string {
	if t == nil {
		return "Object"
	}
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	name := rt.Name()
	switch {
	case strings.Contains(name, "String"):
		return "String"
	case strings.Contains(name, "Bool"):
		return "Bool"
	case strings.Contains(name, "Int"):
		return "Int"
	case strings.Contains(name, "Float"), strings.Contains(name, "Number"):
		return "Number"
	case strings.Contains(name, "List"):
		return "List"
	case strings.Contains(name, "Set"):
		return "Set"
	case strings.Contains(name, "Map"):
		return "Map"
	case strings.Contains(name, "Object"):
		return "Object"
	case strings.Contains(name, "Dynamic"):
		return "Dynamic"
	}
	return "Object"
}

// ensureSentence appends a trailing period when the description does not
// end with `.` or `:` so checkDescription does not reject it.
func ensureSentence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	last := s[len(s)-1]
	if last != '.' && last != ':' && last != '!' && last != '?' {
		s += "."
	}
	return s
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

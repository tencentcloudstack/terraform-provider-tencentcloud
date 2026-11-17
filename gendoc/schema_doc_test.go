package main

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestGetSubStructAttrs verifies that Computed sub-fields of a
// Required/Optional nested object are rendered into an "exports the
// following" block, while Required/Optional sub-fields are excluded.
func TestGetSubStructAttrs(t *testing.T) {
	v := &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Parameter name.",
				},
				"expected_value": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The new value.",
				},
				"default_value": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The default value.",
				},
				"param_description_ch": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Parameter Chinese Description.",
				},
			},
		},
	}

	got := getSubStructAttrs(0, "", "param_list", v)
	if len(got) != 1 {
		t.Fatalf("expected 1 block, got %d: %#v", len(got), got)
	}
	block := got[0]
	if !strings.Contains(block, "The `param_list` object exports the following:") {
		t.Errorf("missing exports title: %s", block)
	}
	if !strings.Contains(block, "* `default_value` - The default value.") {
		t.Errorf("missing computed field default_value: %s", block)
	}
	if !strings.Contains(block, "* `param_description_ch` - Parameter Chinese Description.") {
		t.Errorf("missing computed field param_description_ch: %s", block)
	}
	if strings.Contains(block, "* `name` -") || strings.Contains(block, "* `expected_value` -") {
		t.Errorf("required/optional fields must not appear in exports block: %s", block)
	}
}

// TestGetSubStructAttrsNested verifies deep nesting produces a
// "The `b` object of `a` exports the following:" title.
func TestGetSubStructAttrsNested(t *testing.T) {
	v := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"b": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "b object.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"c": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "c value.",
							},
						},
					},
				},
			},
		},
	}

	got := getSubStructAttrs(0, "", "a", v)
	if len(got) != 1 {
		t.Fatalf("expected 1 block, got %d: %#v", len(got), got)
	}
	if !strings.Contains(got[0], "The `b` object of `a` exports the following:") {
		t.Errorf("missing nested exports title: %s", got[0])
	}
	if !strings.Contains(got[0], "* `c` - c value.") {
		t.Errorf("missing nested computed field c: %s", got[0])
	}
}

// TestRenderNestedExports verifies the framework renderer mirrors the SDKv2
// renderer: Computed sub-fields go to an exports block, others are excluded.
func TestRenderNestedExports(t *testing.T) {
	nested := []fwAttrSpec{
		{Name: "name", Required: true, TypeLabel: "String", Description: "Parameter name."},
		{Name: "expected_value", Required: true, TypeLabel: "String", Description: "The new value."},
		{Name: "default_value", Computed: true, Description: "The default value."},
		{Name: "param_description_ch", Computed: true, Description: "Parameter Chinese Description."},
	}

	got := renderNestedExports("", "param_list", nested)
	if len(got) != 1 {
		t.Fatalf("expected 1 block, got %d: %#v", len(got), got)
	}
	block := got[0]
	if !strings.Contains(block, "The `param_list` object exports the following:") {
		t.Errorf("missing exports title: %s", block)
	}
	if !strings.Contains(block, "* `default_value` - The default value.") {
		t.Errorf("missing computed field default_value: %s", block)
	}
	if strings.Contains(block, "* `name` -") || strings.Contains(block, "* `expected_value` -") {
		t.Errorf("required/optional fields must not appear in exports block: %s", block)
	}
}

// TestRenderNestedRecursion verifies renderNested recurses into deeper
// Required/Optional objects with the "object of" title, matching getSubStruct.
func TestRenderNestedRecursion(t *testing.T) {
	nested := []fwAttrSpec{
		{
			Name:        "b",
			Optional:    true,
			TypeLabel:   "List",
			Description: "b object.",
			Nested: []fwAttrSpec{
				{Name: "d", Required: true, TypeLabel: "String", Description: "d value."},
			},
		},
	}

	got := renderNested("", "a", nested)
	joined := strings.Join(got, "\n")
	if !strings.Contains(joined, "The `a` object supports the following:") {
		t.Errorf("missing top-level supports title: %s", joined)
	}
	if !strings.Contains(joined, "The `b` object of `a` supports the following:") {
		t.Errorf("missing nested supports title: %s", joined)
	}
	if !strings.Contains(joined, "* `d` - (Required, String) d value.") {
		t.Errorf("missing nested required field d: %s", joined)
	}
}

// TestGetSubStructAttrsSkipsOptionalComputed verifies that Optional+Computed
// sub-fields are NOT duplicated into the exports block (they are already
// listed as Optional arguments in the supports block).
func TestGetSubStructAttrsSkipsOptionalComputed(t *testing.T) {
	v := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "The name of parameter.",
				},
				"current_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "The value of parameter.",
				},
				"readonly": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Readonly field.",
				},
			},
		},
	}

	got := getSubStructAttrs(0, "", "param_list", v)
	if len(got) != 1 {
		t.Fatalf("expected 1 block, got %d: %#v", len(got), got)
	}
	block := got[0]
	if !strings.Contains(block, "* `readonly` - Readonly field.") {
		t.Errorf("missing pure computed field readonly: %s", block)
	}
	if strings.Contains(block, "`name`") || strings.Contains(block, "`current_value`") {
		t.Errorf("optional+computed fields must not appear in exports block: %s", block)
	}
}

// TestRenderNestedExportsSkipsOptionalComputed is the framework renderer
// counterpart of the SDKv2 Optional+Computed dedup test.
func TestRenderNestedExportsSkipsOptionalComputed(t *testing.T) {
	nested := []fwAttrSpec{
		{Name: "name", Optional: true, Computed: true, TypeLabel: "String", Description: "The name of parameter."},
		{Name: "current_value", Optional: true, Computed: true, TypeLabel: "String", Description: "The value of parameter."},
		{Name: "readonly", Computed: true, Description: "Readonly field."},
	}

	got := renderNestedExports("", "param_list", nested)
	if len(got) != 1 {
		t.Fatalf("expected 1 block, got %d: %#v", len(got), got)
	}
	block := got[0]
	if !strings.Contains(block, "* `readonly` - Readonly field.") {
		t.Errorf("missing pure computed field readonly: %s", block)
	}
	if strings.Contains(block, "`name`") || strings.Contains(block, "`current_value`") {
		t.Errorf("optional+computed fields must not appear in exports block: %s", block)
	}
}

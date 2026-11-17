package cosvectors

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCosVectorConfigurationObject(t *testing.T) {
	typ := types.ObjectType{AttrTypes: cosVectorEncryptionTypes}
	for _, tc := range []struct {
		name  string
		value types.List
		valid bool
	}{
		{"null", types.ListNull(typ), true},
		{"unknown", types.ListUnknown(typ), true},
		{"configured", cosVectorEncryptionValue("AES256"), true},
		{"empty", types.ListValueMust(typ, nil), false},
		{"two", types.ListValueMust(typ, []attr.Value{cosVectorEncryptionValue("AES256").Elements()[0], cosVectorEncryptionValue("AES256").Elements()[0]}), false},
		{"nullElement", types.ListValueMust(typ, []attr.Value{types.ObjectNull(cosVectorEncryptionTypes)}), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := cosVectorConfigurationObject(tc.value)
			if (err == nil) != tc.valid {
				t.Fatal(err)
			}
		})
	}
}

func TestCosVectorMetadataUnknownKeyValidation(t *testing.T) {
	ctx := context.Background()
	r := &CosVectorIndexResource{}
	var schema resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schema)
	keys := []attr.Value{types.StringUnknown()}
	for _, key := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"} {
		keys = append(keys, types.StringValue(key))
	}
	object := types.ObjectValueMust(cosVectorMetadataTypes, map[string]attr.Value{"non_filterable_metadata_keys": types.SetValueMust(types.StringType, keys)})
	data := CosVectorIndexResourceModel{VectorBucketName: types.StringValue("tf-unknown-1250000000"), IndexName: types.StringValue("test-index"), DataType: types.StringValue("float32"), Dimension: types.Int64Value(3), DistanceMetric: types.StringValue("cosine"), MetadataConfiguration: types.ListValueMust(types.ObjectType{AttrTypes: cosVectorMetadataTypes}, []attr.Value{object})}
	state := tfsdk.State{Schema: schema.Schema}
	if diags := state.Set(ctx, &data); diags.HasError() {
		t.Fatal(diags)
	}
	var resp resource.ValidateConfigResponse
	r.ValidateConfig(ctx, resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: state.Raw}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal("unknown key may resolve to an existing set member", resp.Diagnostics)
	}
}

func TestCosVectorPolicyServerFieldNames(t *testing.T) {
	lower := `{"version":"2.0","statement":[{"effect":"allow","principal":{"qcs":["qcs::cam::uin/100000000001:uin/100000000002"]},"action":["name/cos:QueryVectors"],"resource":["qcs::cosvector:ap-shanghai:uid/1250000000:bucket/tf-canary-1250000000/index/test-index"],"condition":{"StringEquals":{"CustomCaseSensitiveKey":"Value"}}}]}`
	server := strings.NewReplacer(`"version":`, `"Version":`, `"statement":`, `"Statement":`, `"effect":`, `"Effect":`, `"principal":`, `"Principal":`, `"action":`, `"Action":`, `"resource":`, `"Resource":`, `"condition":`, `"Condition":`).Replace(lower)
	if err := validateCosVectorPolicy(server); err != nil {
		t.Fatal(err)
	}
	if !cosVectorPoliciesEqual(lower, server) {
		t.Fatal("server field casing caused policy drift")
	}
	if cosVectorPoliciesEqual(lower, strings.Replace(server, "CustomCaseSensitiveKey", "customcasesensitivekey", 1)) {
		t.Fatal("condition key casing was ignored")
	}
	conflict := strings.Replace(server, `"Version":"2.0"`, `"Version":"2.0","version":"2.0"`, 1)
	if _, err := normalizeCosVectorPolicy(conflict); err == nil {
		t.Fatal("ambiguous field aliases accepted")
	}
}

func TestCosVectorPolicyValidation(t *testing.T) {
	valid := `{"version":"2.0","statement":[{"effect":"allow","principal":{"qcs":["qcs::cam::uin/100000000001:uin/100000000002"]},"action":["name/cos:QueryVectors"],"resource":["qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-test-1250000000/index/test-index"]}]}`
	var policy map[string]any
	_ = json.Unmarshal([]byte(valid), &policy)
	if err := validateCosVectorPolicy(valid); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"effect", "principal", "action", "resource"} {
		t.Run(field, func(t *testing.T) {
			var invalid map[string]any
			_ = json.Unmarshal([]byte(valid), &invalid)
			delete(invalid["statement"].([]any)[0].(map[string]any), field)
			encoded, _ := json.Marshal(invalid)
			if validateCosVectorPolicy(string(encoded)) == nil {
				t.Fatal("missing field accepted")
			}
		})
	}
	statement := policy["statement"].([]any)[0].(map[string]any)
	statement["principal"] = map[string]any{"qcs": []string{""}}
	encoded, _ := json.Marshal(policy)
	if validateCosVectorPolicy(string(encoded)) == nil {
		t.Fatal("empty principal accepted")
	}
	statement["principal"] = map[string]any{"qcs": []string{"qcs::cam::uin/100000000001:uin/100000000002"}}
	for _, effect := range []string{"Allow", "ALLOW", "Deny"} {
		statement["effect"] = effect
		encoded, _ = json.Marshal(policy)
		if err := validateCosVectorPolicy(string(encoded)); err != nil {
			t.Fatal("supported effect rejected", err)
		}
		if !cosVectorPoliciesEqual(string(encoded), strings.Replace(string(encoded), effect, strings.ToLower(effect), 1)) {
			t.Fatal("effect casing caused policy drift")
		}
	}
	statement["effect"] = "deny"
	statement["principal"] = map[string]any{"qcs": []string{"*"}}
	encoded, _ = json.Marshal(policy)
	if err := validateCosVectorPolicy(string(encoded)); err != nil {
		t.Fatal("deny wildcard rejected", err)
	}
}

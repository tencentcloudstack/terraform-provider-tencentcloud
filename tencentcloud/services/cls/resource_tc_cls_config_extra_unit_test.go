package cls

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestClsConfigExtraExtractRuleKeysPreserveOrder(t *testing.T) {
	resource := ResourceTencentCloudClsConfigExtra()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"extract_rule": []interface{}{
			map[string]interface{}{
				"keys": []interface{}{"first", "second"},
			},
		},
	})

	extractRules := data.Get("extract_rule").([]interface{})
	keysValue := extractRules[0].(map[string]interface{})["keys"]
	keys, ok := keysValue.([]interface{})
	if !ok {
		t.Fatalf("extract_rule.keys must preserve order as a list, got %T", keysValue)
	}

	want := []interface{}{"first", "second"}
	if !reflect.DeepEqual(keys, want) {
		t.Fatalf("extract_rule.keys = %#v, want %#v", keys, want)
	}
}

func TestClsConfigExtraStateUpgradeV0(t *testing.T) {
	resourceDef := ResourceTencentCloudClsConfigExtra()
	if resourceDef.SchemaVersion != 1 {
		t.Fatalf("SchemaVersion = %d, want 1", resourceDef.SchemaVersion)
	}
	if len(resourceDef.StateUpgraders) != 1 {
		t.Fatalf("StateUpgraders length = %d, want 1", len(resourceDef.StateUpgraders))
	}

	upgrader := resourceDef.StateUpgraders[0]
	legacyKeysType := upgrader.Type.AttributeType("extract_rule").ElementType().AttributeType("keys")
	if !legacyKeysType.IsSetType() {
		t.Fatalf("legacy extract_rule.keys type = %s, want set", legacyKeysType.FriendlyName())
	}

	currentType := resourceDef.CoreConfigSchema().ImpliedType()
	currentKeysType := currentType.AttributeType("extract_rule").ElementType().AttributeType("keys")
	if !currentKeysType.IsListType() {
		t.Fatalf("current extract_rule.keys type = %s, want list", currentKeysType.FriendlyName())
	}

	server := schema.NewGRPCProviderServer(&schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"tencentcloud_cls_config_extra": resourceDef,
		},
	})
	firstKey := fmt.Sprintf("extract_rule.0.keys.%d", schema.HashString("first"))
	secondKey := fmt.Sprintf("extract_rule.0.keys.%d", schema.HashString("second"))

	tests := []struct {
		name     string
		rawState *tfprotov5.RawState
		ordered  bool
	}{
		{
			name: "JSON state",
			rawState: &tfprotov5.RawState{
				JSON: []byte(`{"id":"config-extra-id","extract_rule":[{"keys":["first","second"]}]}`),
			},
			ordered: true,
		},
		{
			name: "legacy flatmap state",
			rawState: &tfprotov5.RawState{
				Flatmap: map[string]string{
					"id":                    "config-extra-id",
					"extract_rule.#":        "1",
					"extract_rule.0.keys.#": "2",
					firstKey:                "first",
					secondKey:               "second",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := server.UpgradeResourceState(t.Context(), &tfprotov5.UpgradeResourceStateRequest{
				TypeName: "tencentcloud_cls_config_extra",
				Version:  0,
				RawState: test.rawState,
			})
			if err != nil {
				t.Fatalf("UpgradeResourceState returned error: %v", err)
			}
			if len(response.Diagnostics) != 0 {
				t.Fatalf("UpgradeResourceState returned diagnostics: %#v", response.Diagnostics)
			}
			if response.UpgradedState == nil || len(response.UpgradedState.MsgPack) == 0 {
				t.Fatal("UpgradeResourceState returned empty state")
			}

			keys := clsConfigExtraUpgradedKeys(t, server, response.UpgradedState)
			if test.ordered {
				want := []string{"first", "second"}
				if !reflect.DeepEqual(keys, want) {
					t.Fatalf("upgraded keys = %#v, want %#v", keys, want)
				}
				return
			}

			if len(keys) != 2 {
				t.Fatalf("upgraded keys = %#v, want two values", keys)
			}
			got := map[string]bool{keys[0]: true, keys[1]: true}
			if !got["first"] || !got["second"] {
				t.Fatalf("upgraded keys = %#v, want first and second", keys)
			}
		})
	}

	if err := resourceDef.InternalValidate(nil, true); err != nil {
		t.Fatalf("resource validation failed: %v", err)
	}
}

func clsConfigExtraUpgradedKeys(t *testing.T, server *schema.GRPCProviderServer, state *tfprotov5.DynamicValue) []string {
	t.Helper()

	schemaResponse, err := server.GetProviderSchema(context.Background(), &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema returned error: %v", err)
	}
	if len(schemaResponse.Diagnostics) != 0 {
		t.Fatalf("GetProviderSchema returned diagnostics: %#v", schemaResponse.Diagnostics)
	}
	resourceSchema := schemaResponse.ResourceSchemas["tencentcloud_cls_config_extra"]
	if resourceSchema == nil {
		t.Fatal("GetProviderSchema did not return tencentcloud_cls_config_extra")
	}

	value, err := state.Unmarshal(resourceSchema.ValueType())
	if err != nil {
		t.Fatalf("cannot decode upgraded state: %v", err)
	}
	var attributes map[string]tftypes.Value
	if err := value.As(&attributes); err != nil {
		t.Fatalf("cannot decode state attributes: %v", err)
	}
	var extractRules []tftypes.Value
	if err := attributes["extract_rule"].As(&extractRules); err != nil {
		t.Fatalf("cannot decode extract_rule: %v", err)
	}
	if len(extractRules) != 1 {
		t.Fatalf("extract_rule length = %d, want 1", len(extractRules))
	}
	var extractRule map[string]tftypes.Value
	if err := extractRules[0].As(&extractRule); err != nil {
		t.Fatalf("cannot decode extract_rule element: %v", err)
	}
	var keyValues []tftypes.Value
	if err := extractRule["keys"].As(&keyValues); err != nil {
		t.Fatalf("cannot decode keys: %v", err)
	}

	keys := make([]string, len(keyValues))
	for i, keyValue := range keyValues {
		if err := keyValue.As(&keys[i]); err != nil {
			t.Fatalf("cannot decode key %d: %v", i, err)
		}
	}
	return keys
}

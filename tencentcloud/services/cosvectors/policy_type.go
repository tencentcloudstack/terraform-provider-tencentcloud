package cosvectors

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// cosVectorPolicyType 忽略 JSON 格式及 COS 策略集合顺序。
type cosVectorPolicyType struct{ basetypes.StringType }
type cosVectorPolicyValue struct{ basetypes.StringValue }

func (cosVectorPolicyType) Equal(other attr.Type) bool {
	_, ok := other.(cosVectorPolicyType)
	return ok
}
func (cosVectorPolicyType) String() string                       { return "cosvectors.VectorPolicyType" }
func (cosVectorPolicyType) ValueType(context.Context) attr.Value { return cosVectorPolicyValue{} }
func (cosVectorPolicyType) ValueFromString(_ context.Context, value basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return cosVectorPolicyValue{StringValue: value}, nil
}
func (t cosVectorPolicyType) ValueFromTerraform(ctx context.Context, value tftypes.Value) (attr.Value, error) {
	base, err := t.StringType.ValueFromTerraform(ctx, value)
	if err != nil {
		return nil, err
	}
	return cosVectorPolicyValue{StringValue: base.(basetypes.StringValue)}, nil
}
func (cosVectorPolicyValue) Type(context.Context) attr.Type { return cosVectorPolicyType{} }
func (v cosVectorPolicyValue) Equal(other attr.Value) bool {
	value, ok := other.(cosVectorPolicyValue)
	return ok && v.StringValue.Equal(value.StringValue)
}
func (v cosVectorPolicyValue) StringSemanticEquals(_ context.Context, other basetypes.StringValuable) (bool, diag.Diagnostics) {
	value, ok := other.(cosVectorPolicyValue)
	if !ok {
		return false, nil
	}
	return cosVectorPoliciesEqual(v.ValueString(), value.ValueString()), nil
}

type cosVectorPolicyPlanModifier struct{}

func (cosVectorPolicyPlanModifier) Description(context.Context) string {
	return "Preserves state for semantically equivalent COS policy JSON."
}
func (m cosVectorPolicyPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}
func (cosVectorPolicyPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.Plan.Raw.IsNull() || req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if cosVectorPoliciesEqual(req.PlanValue.ValueString(), req.StateValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}

func cosVectorPolicyString(value string) cosVectorPolicyValue {
	return cosVectorPolicyValue{StringValue: basetypes.NewStringValue(value)}
}
func cosVectorPolicyNull() cosVectorPolicyValue {
	return cosVectorPolicyValue{StringValue: basetypes.NewStringNull()}
}

func normalizeCosVectorPolicy(text string) (map[string]any, error) {
	if len(text) > 1<<20 {
		return nil, errors.New("policy exceeds the provider response limit")
	}
	var policy map[string]any
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	if !json.Valid([]byte(text)) || decoder.Decode(&policy) != nil || policy == nil {
		return nil, errors.New("policy must be a JSON object")
	}
	policy = normalizeCosVectorPolicyNumbers(policy).(map[string]any)
	if err := normalizeCosVectorPolicyFields(policy, "version", "statement"); err != nil {
		return nil, err
	}
	if policy["version"] != "2.0" {
		return nil, errors.New("COS policy version must be 2.0")
	}
	statements, ok := policy["statement"].([]any)
	if !ok || len(statements) == 0 {
		return nil, errors.New("COS policy statement must be a non-empty array")
	}
	for _, entry := range statements {
		statement, ok := entry.(map[string]any)
		if !ok {
			return nil, errors.New("COS policy statements must be objects")
		}
		if err := normalizeCosVectorPolicyFields(statement, "effect", "principal", "action", "resource", "condition", "sid"); err != nil {
			return nil, err
		}
		if effect, ok := statement["effect"].(string); ok {
			statement["effect"] = strings.ToLower(effect)
		}
		for _, key := range []string{"action", "resource"} {
			sortCosVectorPolicyArray(statement[key])
		}
		if principal, ok := statement["principal"].(map[string]any); ok {
			sortCosVectorPolicyArray(principal["qcs"])
		}
	}
	sortCosVectorPolicyArray(statements)
	return policy, nil
}

func normalizeCosVectorPolicyFields(object map[string]any, fields ...string) error {
	for _, field := range fields {
		for key, value := range object {
			if key == field || !strings.EqualFold(key, field) {
				continue
			}
			if _, exists := object[field]; exists {
				return errors.New("COS policy contains conflicting field names: " + field)
			}
			object[field] = value
			delete(object, key)
		}
	}
	return nil
}

func validateCosVectorPolicy(text string) error {
	policy, err := normalizeCosVectorPolicy(text)
	if err != nil {
		return err
	}
	for _, entry := range policy["statement"].([]any) {
		statement := entry.(map[string]any)
		effect, ok := statement["effect"].(string)
		if !ok || (effect != "allow" && effect != "deny") {
			return errors.New("COS policy effect must be allow or deny")
		}
		principal, ok := statement["principal"].(map[string]any)
		if !ok || !cosVectorPolicyStrings(principal["qcs"]) {
			return errors.New("COS policy principal.qcs must be a non-empty array of non-empty strings")
		}
		for _, field := range []string{"action", "resource"} {
			if !cosVectorPolicyStrings(statement[field]) {
				return errors.New("COS policy " + field + " must be a non-empty array of non-empty strings")
			}
		}
	}
	return nil
}

func cosVectorPolicyStrings(value any) bool {
	items, ok := value.([]any)
	if !ok || len(items) == 0 {
		return false
	}
	for _, item := range items {
		text, ok := item.(string)
		if !ok || strings.TrimSpace(text) == "" {
			return false
		}
	}
	return true
}

type cosVectorPolicyNumber struct{ Exact string }

// normalizeCosVectorPolicyNumbers 使用十进制系数与指数，不展开巨大指数。
func normalizeCosVectorPolicyNumbers(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			v[key] = normalizeCosVectorPolicyNumbers(child)
		}
	case []any:
		for i, child := range v {
			v[i] = normalizeCosVectorPolicyNumbers(child)
		}
	case json.Number:
		text := strings.ToLower(string(v))
		parts := strings.SplitN(text, "e", 2)
		exponent := new(big.Int)
		if len(parts) == 2 {
			exponent.SetString(parts[1], 10)
		}
		coefficient := parts[0]
		if dot := strings.IndexByte(coefficient, '.'); dot >= 0 {
			exponent.Sub(exponent, big.NewInt(int64(len(coefficient)-dot-1)))
			coefficient = strings.ReplaceAll(coefficient, ".", "")
		}
		sign := ""
		if strings.HasPrefix(coefficient, "-") {
			sign, coefficient = "-", coefficient[1:]
		}
		coefficient = strings.TrimLeft(coefficient, "0")
		if coefficient == "" {
			return cosVectorPolicyNumber{Exact: "0"}
		}
		trimmed := strings.TrimRight(coefficient, "0")
		exponent.Add(exponent, big.NewInt(int64(len(coefficient)-len(trimmed))))
		return cosVectorPolicyNumber{Exact: sign + trimmed + "e" + exponent.String()}
	}
	return value
}

func sortCosVectorPolicyArray(value any) {
	items, ok := value.([]any)
	if !ok {
		return
	}
	sort.SliceStable(items, func(i, j int) bool {
		a, _ := json.Marshal(items[i])
		b, _ := json.Marshal(items[j])
		return string(a) < string(b)
	})
}

func cosVectorPoliciesEqual(left, right string) bool {
	a, err := normalizeCosVectorPolicy(left)
	if err != nil {
		return false
	}
	b, err := normalizeCosVectorPolicy(right)
	return err == nil && reflect.DeepEqual(a, b)
}

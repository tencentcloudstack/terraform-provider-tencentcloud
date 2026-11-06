package ags

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dependencyConfig(t *testing.T, ids []attr.Value, logs []attr.Value) types.Object {
	t.Helper()
	ctx := context.Background()
	base := modelFixture(t, map[string]any{"network_configuration": []map[string]any{{"network_mode": "VPC", "vpc_config": []map[string]any{{"subnet_ids": []string{"subnet-owned"}, "security_group_ids": []string{"sg-owned"}}}}}, "log_configuration": []map[string]any{{"topic_id": "topic-owned", "files": []string{"/logs/nonce.log"}}}})
	f := base.Attributes()
	n := f["network_configuration"].(types.List).Elements()[0].(types.Object)
	nf := n.Attributes()
	v := nf["vpc_config"].(types.List).Elements()[0].(types.Object)
	vf := v.Attributes()
	vf["subnet_ids"] = types.SetValueMust(types.StringType, ids)
	vf["security_group_ids"] = types.SetValueMust(types.StringType, ids)
	v = types.ObjectValueMust(v.AttributeTypes(ctx), vf)
	nf["vpc_config"] = types.ListValueMust(v.Type(ctx), []attr.Value{v})
	n = types.ObjectValueMust(n.AttributeTypes(ctx), nf)
	f["network_configuration"] = types.ListValueMust(n.Type(ctx), []attr.Value{n})
	l := f["log_configuration"].(types.List).Elements()[0].(types.Object)
	lf := l.Attributes()
	lf["files"] = types.SetValueMust(types.StringType, logs)
	l = types.ObjectValueMust(l.AttributeTypes(ctx), lf)
	f["log_configuration"] = types.ListValueMust(l.Type(ctx), []attr.Value{l})
	return types.ObjectValueMust(base.AttributeTypes(ctx), f)
}

func TestDependencyUnknownElementsDeferValidationUntilResolved(t *testing.T) {
	config := dependencyConfig(t, []attr.Value{types.StringUnknown()}, []attr.Value{types.StringUnknown()})
	if ds := validateToolConfig(config); ds.HasError() {
		t.Fatalf("dependency references must be permitted during plan: %v", ds)
	}
	if _, err := createRequest(config); err == nil {
		t.Fatal("unresolved dependency values reached an SDK create request")
	}
	resolved := dependencyConfig(t, []attr.Value{types.StringValue("owned-id")}, []attr.Value{types.StringValue("/logs/nonce.log")})
	if ds := validateToolConfig(resolved); ds.HasError() {
		t.Fatal(ds)
	}
	if _, err := createRequest(resolved); err != nil {
		t.Fatal("fully resolved references must be serializable", err)
	}
}

func TestDependencyKnownInvalidElementsAreStillRejectedBesideUnknowns(t *testing.T) {
	for _, kind := range []string{"null-id", "empty-id", "empty-ids", "null-file", "invalid-file", "empty-files"} {
		t.Run(kind, func(t *testing.T) {
			ids, logs := []attr.Value{types.StringUnknown()}, []attr.Value{types.StringUnknown()}
			switch kind {
			case "null-id":
				ids = append(ids, types.StringNull())
			case "empty-id":
				ids = append(ids, types.StringValue(""))
			case "empty-ids":
				ids = nil
			case "null-file":
				logs = append(logs, types.StringNull())
			case "invalid-file":
				logs = append(logs, types.StringValue("/shared/outside.log"))
			case "empty-files":
				logs = nil
			}
			if ds := validateToolConfig(dependencyConfig(t, ids, logs)); !ds.HasError() {
				t.Fatal("a known invalid value was hidden by another unknown element")
			}
		})
	}
}

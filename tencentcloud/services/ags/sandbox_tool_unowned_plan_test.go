package ags

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

func TestUnownedCustomInputsOnlyReuseRefreshedStateForAnOtherwiseUnchangedPlan(t *testing.T) {
	ctx := context.Background()
	for _, name := range []string{"omitted-env", "omitted-args", "omitted-command", "omitted-env-after-remote-drift", "owned-env-drift", "explicit-env-clear", "image-change-with-env-omitted", "description-change-with-env-omitted", "configured-env-unknown"} {
		t.Run(name, func(t *testing.T) {
			tool := remoteCode()
			tool.ToolType = pointer("custom")
			tool.CustomConfiguration = &sdk.CustomConfigurationDetail{Image: pointer("registry.invalid/app:1"), ImageRegistryType: pointer("enterprise"), ImageDigest: pointer("sha256:1111"), Command: []*string{pointer("/app/start")}, Args: []*string{pointer("--one")}, Env: []*sdk.EnvVar{{Name: pointer("TF_CASE"), Value: pointer("actual")}}, Resources: &sdk.ResourceConfiguration{CPU: pointer("1"), Memory: pointer("1Gi")}}
			if name == "omitted-env-after-remote-drift" {
				tool.CustomConfiguration.Env[0].Value = pointer("fresh-external-value")
			}
			prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), tool, true)
			if err != nil {
				t.Fatal(err)
			}
			configFields, planFields := prior.Attributes(), prior.Attributes()
			custom := prior.Attributes()["custom_configuration"].(types.List).Elements()[0].(types.Object)
			configCustom, planCustom := custom.Attributes(), custom.Attributes()
			configFields["id"] = types.StringNull()
			for _, key := range []string{"status", "status_reason", "image_digest"} {
				configFields[key], planFields[key] = types.StringNull(), types.StringUnknown()
			}
			configFields["effective_configuration"] = types.ObjectNull(configurationTypes())
			planFields["effective_configuration"] = types.ObjectUnknown(configurationTypes())
			key := "env"
			if name == "omitted-args" {
				key = "args"
			} else if name == "omitted-command" {
				key = "command"
			}
			if key == "env" {
				configCustom[key], planCustom[key] = types.MapNull(types.StringType), types.MapUnknown(types.StringType)
			} else {
				configCustom[key], planCustom[key] = types.ListNull(types.StringType), types.ListUnknown(types.StringType)
			}
			switch name {
			case "owned-env-drift":
				declared := types.MapValueMust(types.StringType, map[string]attr.Value{"TF_CASE": types.StringValue("declared")})
				configCustom[key], planCustom[key] = declared, declared
			case "explicit-env-clear":
				empty := types.MapValueMust(types.StringType, map[string]attr.Value{})
				configCustom[key], planCustom[key] = empty, empty
			case "image-change-with-env-omitted":
				configCustom["image"], planCustom["image"] = types.StringValue("registry.invalid/app:2"), types.StringValue("registry.invalid/app:2")
			case "description-change-with-env-omitted":
				configFields["description"], planFields["description"] = types.StringValue("new description"), types.StringValue("new description")
			case "configured-env-unknown":
				configCustom[key] = types.MapUnknown(types.StringType)
			}
			configFields["custom_configuration"] = types.ListValueMust(custom.Type(ctx), []attr.Value{types.ObjectValueMust(custom.AttributeTypes(ctx), configCustom)})
			planFields["custom_configuration"] = types.ListValueMust(custom.Type(ctx), []attr.Value{types.ObjectValueMust(custom.AttributeTypes(ctx), planCustom)})
			config, plan := types.ObjectValueMust(inputAttributeTypes(), configFields), types.ObjectValueMust(inputAttributeTypes(), planFields)
			r := &SandboxToolResource{}
			response := resource.ModifyPlanResponse{Plan: planFor(t, plan)}
			r.ModifyPlan(ctx, resource.ModifyPlanRequest{State: stateFor(t, prior), Plan: planFor(t, plan), Config: configFor(t, config)}, &response)
			if response.Diagnostics.HasError() {
				t.Fatal(response.Diagnostics)
			}
			var actual types.Object
			if ds := response.Plan.Get(ctx, &actual); ds.HasError() {
				t.Fatal(ds)
			}
			unchanged := name == "omitted-env" || name == "omitted-args" || name == "omitted-command" || name == "omitted-env-after-remote-drift"
			if unchanged {
				if !actual.Equal(prior) {
					t.Fatal("omitted input manufactured an update instead of accepting the refreshed remote value")
				}
				_, changed, err := updateRequest(prior, actual)
				if err != nil || changed {
					t.Fatal("unowned input requires an SDK mutation", changed, err)
				}
				return
			}
			if !actual.Attributes()["effective_configuration"].IsUnknown() {
				t.Fatal("a real update or configured unknown had its future readback frozen to old state")
			}
			actualCustom := actual.Attributes()["custom_configuration"].(types.List).Elements()[0].(types.Object).Attributes()
			if !actualCustom[key].Equal(planCustom[key]) {
				t.Fatal("real drift, explicit clear or future unowned value was hidden")
			}
			if name == "explicit-env-clear" {
				request, changed, err := updateRequest(prior, actual)
				if err != nil || !changed || request.CustomConfiguration.Env == nil || len(request.CustomConfiguration.Env) != 0 {
					t.Fatal("explicit empty Env no longer requests a clear", err)
				}
			}
		})
	}
}

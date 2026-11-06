package ags

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

// Public Describe returns Resources.Storage even though V1 intentionally has
// only CPU/memory inputs. Import must not silently discard a disk override.
func TestImportRejectsUnsupportedCustomDiskThroughActualSDK(t *testing.T) {
	for _, test := range []struct {
		name    string
		storage *string
		reject  bool
	}{
		{name: "omitted"},
		{name: "empty", storage: pointer("")},
		{name: "10Gi", storage: pointer("10Gi"), reject: true},
		{name: "20Gi", storage: pointer("20Gi"), reject: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			tool := remoteCode()
			tool.ToolType = pointer("custom")
			tool.CustomConfiguration = &sdk.CustomConfigurationDetail{
				Image: pointer("ccr.ccs.tencentyun.com/offline/ags:fixture"), ImageRegistryType: pointer("personal"),
				Resources: &sdk.ResourceConfiguration{CPU: pointer("1"), Memory: pointer("1Gi"), Storage: test.storage},
			}
			reads := 0
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				reads++
				if action != "DescribeSandboxToolList" || len(input["ToolIds"].([]any)) != 1 || input["ToolIds"].([]any)[0] != "sdt-offline" {
					t.Fatal("import sent an unconstrained read or mutation", action)
				}
				return map[string]any{"SandboxToolSet": []any{tool}, "TotalCount": 1}, nil
			})
			r := &SandboxToolResource{testService: fastService(client)}
			response := resource.ImportStateResponse{State: tfsdk.State{Schema: toolSchema()}}
			r.ImportState(context.Background(), resource.ImportStateRequest{ID: "sdt-offline"}, &response)
			if reads != 1 {
				t.Fatal("import did not perform exactly one point query", reads)
			}
			if test.reject {
				if !response.Diagnostics.HasError() || !response.State.Raw.IsNull() || !strings.Contains(response.Diagnostics.Errors()[0].Detail(), "disk") {
					t.Fatal("custom disk was silently accepted/truncated instead of rejecting adoption", response.Diagnostics)
				}
				return
			}
			if response.Diagnostics.HasError() {
				t.Fatal("an omitted/empty unsupported disk should not reject an otherwise supported Tool", response.Diagnostics)
			}
			var imported types.Object
			if response.State.Get(context.Background(), &imported).HasError() || imported.Attributes()["id"].(types.String).ValueString() != "sdt-offline" {
				t.Fatal("supported import did not preserve its actual identity")
			}
		})
	}
}

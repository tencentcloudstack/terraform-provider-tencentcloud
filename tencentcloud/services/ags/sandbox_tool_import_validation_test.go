package ags

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

func TestImportValidatesRepresentableCustomInputsThroughActualSDK(t *testing.T) {
	for _, port := range []int64{8080, 49999} {
		t.Run(stringPortForImportTest(port), func(t *testing.T) {
			tool := remoteCode()
			tool.ToolType = pointer("custom")
			tool.CustomConfiguration = &sdk.CustomConfigurationDetail{
				Image: pointer("ccr.ccs.tencentyun.com/offline/ags:fixture"), ImageRegistryType: pointer("personal"),
				Resources: &sdk.ResourceConfiguration{CPU: pointer("1"), Memory: pointer("1Gi")},
				Ports:     []*sdk.PortConfiguration{{Name: pointer("business"), Port: pointer(port), Protocol: pointer("TCP")}},
			}
			client := testClient(t, func(action string, _ map[string]any) (map[string]any, error) {
				if action != "DescribeSandboxToolList" {
					t.Fatal("import emitted a write", action)
				}
				return map[string]any{"SandboxToolSet": []any{tool}, "TotalCount": 1}, nil
			})
			r := &SandboxToolResource{testService: fastService(client)}
			response := resource.ImportStateResponse{State: tfsdk.State{Schema: toolSchema()}}
			r.ImportState(context.Background(), resource.ImportStateRequest{ID: "sdt-offline"}, &response)
			if port == 49999 {
				if !response.Diagnostics.HasError() || !response.State.Raw.IsNull() {
					t.Fatal("import accepted a remote port that the V1 configuration cannot express", response.Diagnostics)
				}
			} else if response.Diagnostics.HasError() || response.State.Raw.IsNull() {
				t.Fatal("representable business port was incorrectly rejected", response.Diagnostics)
			}
		})
	}
}

func stringPortForImportTest(port int64) string {
	if port == 49999 {
		return "reserved-platform-port"
	}
	return "valid-business-port"
}

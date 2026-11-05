package ags

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDataSourceSelectorsAndReadRemainReadOnly(t *testing.T) {
	ctx := context.Background()
	for _, selector := range []string{"tool_id", "tool_name"} {
		t.Run(selector, func(t *testing.T) {
			d := &SandboxToolDataSource{testService: fastService(testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				if action != "DescribeSandboxToolList" {
					t.Fatal("Data Source issued a mutation")
				}
				return pointResult(testTool("sdt-offline", "owned-code", "ACTIVE")), nil
			}))}
			schema := datasource.SchemaResponse{}
			d.Schema(ctx, datasource.SchemaRequest{}, &schema)
			raw := map[string]any{"tool_id": nil, "tool_name": nil}
			if selector == "tool_id" {
				raw[selector] = "sdt-offline"
			} else {
				raw[selector] = "owned-code"
			}
			v, err := reflectedValue(ctx, schema.Schema.Type(), raw)
			if err != nil {
				t.Fatal(err)
			}
			wire, err := v.ToTerraformValue(ctx)
			if err != nil {
				t.Fatal(err)
			}
			config := tfsdk.Config{Raw: wire, Schema: schema.Schema}
			validation := datasource.ValidateConfigResponse{}
			d.ValidateConfig(ctx, datasource.ValidateConfigRequest{Config: config}, &validation)
			if validation.Diagnostics.HasError() {
				t.Fatal(validation.Diagnostics)
			}
			resp := datasource.ReadResponse{}
			d.Read(ctx, datasource.ReadRequest{Config: config}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatal(resp.Diagnostics)
			}
			var state types.Object
			if ds := resp.State.Get(ctx, &state); ds.HasError() {
				t.Fatal(ds)
			}
			if value(stringInput(state.Attributes(), "id")) != "sdt-offline" {
				t.Fatal("query returned the wrong identity")
			}
			other := "tool_name"
			if selector == "tool_name" {
				other = "tool_id"
			}
			if !state.Attributes()[other].IsNull() {
				t.Fatal("query populated its unused selector")
			}
		})
	}
}

package ags

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

func stateFor(t *testing.T, v types.Object) tfsdk.State {
	t.Helper()
	raw, err := v.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return tfsdk.State{Raw: raw, Schema: toolSchema()}
}
func planFor(t *testing.T, v types.Object) tfsdk.Plan {
	t.Helper()
	state := stateFor(t, v)
	return tfsdk.Plan{Raw: state.Raw, Schema: state.Schema}
}
func configFor(t *testing.T, v types.Object) tfsdk.Config {
	t.Helper()
	state := stateFor(t, v)
	return tfsdk.Config{Raw: state.Raw, Schema: state.Schema}
}
func computedPlan(t *testing.T, v types.Object) types.Object {
	t.Helper()
	fields := v.Attributes()
	for _, key := range []string{"id", "description", "default_timeout", "persistent", "role_arn", "tags", "status", "status_reason", "image_digest", "effective_configuration"} {
		if fields[key].IsNull() {
			unknown, err := fields[key].Type(context.Background()).ValueFromTerraform(context.Background(), tftypes.NewValue(fields[key].Type(context.Background()).TerraformType(context.Background()), tftypes.UnknownValue))
			if err != nil {
				t.Fatal(err)
			}
			fields[key] = unknown
		}
	}
	out, ds := types.ObjectValue(inputAttributeTypes(), fields)
	if ds.HasError() {
		t.Fatal(ds)
	}
	return out
}

func TestResourceReadRetainsStateForDeniedOrFilteredResults(t *testing.T) {
	ctx := context.Background()
	prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), remoteCode(), true)
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"", "UnauthorizedOperation", "AuthFailure.TokenFailure", "ResourceNotFound.SandboxTool"} {
		t.Run(code, func(t *testing.T) {
			calls := 0
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				calls++
				if action != "DescribeSandboxToolList" {
					t.Fatal("Read caused a mutation")
				}
				if code == "" {
					return map[string]any{"SandboxToolSet": []any{}, "TotalCount": 0}, nil
				}
				return errorResult(code), nil
			})
			r := &SandboxToolResource{testService: fastService(client)}
			req := resource.ReadRequest{State: stateFor(t, prior)}
			resp := resource.ReadResponse{}
			r.Read(ctx, req, &resp)
			if code == "ResourceNotFound.SandboxTool" {
				if resp.Diagnostics.HasError() || !resp.State.Raw.IsNull() {
					t.Fatal("confirmed absence was not removed")
				}
			} else {
				if !resp.Diagnostics.HasError() || !resp.State.Raw.Equal(req.State.Raw) {
					t.Fatalf("state was lost on uncertain read: %v", resp.Diagnostics)
				}
			}
			if calls != 1 {
				t.Fatalf("unexpected calls %d", calls)
			}
		})
	}
}

func TestResourceCreateTimeoutReturnsKnownIdentityWithoutCleanup(t *testing.T) {
	ctx := context.Background()
	calls := map[string]int{}
	client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
		calls[action]++
		switch action {
		case "CreateSandboxTool":
			return map[string]any{"ToolId": "sdt-accepted"}, nil
		case "DescribeSandboxToolList":
			return errorResult("InternalError"), nil
		default:
			t.Fatalf("unexpected mutation %s", action)
			return nil, nil
		}
	})
	config := modelFixture(t, map[string]any{"timeouts": []map[string]any{{"create": "10ms", "update": nil, "delete": nil, "read": nil}}})
	plan := computedPlan(t, config)
	r := &SandboxToolResource{testService: fastService(client)}
	resp := resource.CreateResponse{}
	r.Create(ctx, resource.CreateRequest{Plan: planFor(t, plan), Config: configFor(t, config)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("creation timeout was reported as success")
	}
	var state types.Object
	if ds := resp.State.Get(ctx, &state); ds.HasError() {
		t.Fatal(ds)
	}
	if value(stringInput(state.Attributes(), "id")) != "sdt-accepted" {
		t.Fatal("accepted Tool ID was lost")
	}
	if calls["CreateSandboxTool"] != 1 || calls["DeleteSandboxTool"] != 0 {
		t.Fatalf("unsafe creation recovery: %#v", calls)
	}
}

func TestResourceCreateReturnsCompleteReadableState(t *testing.T) {
	ctx := context.Background()
	tool := testTool("sdt-created", "owned-code", "ACTIVE")
	client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
		switch action {
		case "CreateSandboxTool":
			if _, ok := input["ClientToken"]; !ok {
				t.Fatal("missing token")
			}
			return map[string]any{"ToolId": "sdt-created"}, nil
		case "DescribeSandboxToolList":
			return pointResult(tool), nil
		default:
			t.Fatalf("unexpected action %s", action)
			return nil, nil
		}
	})
	config := modelFixture(t, nil)
	plan := computedPlan(t, config)
	r := &SandboxToolResource{testService: fastService(client)}
	resp := resource.CreateResponse{}
	r.Create(ctx, resource.CreateRequest{Plan: planFor(t, plan), Config: configFor(t, config)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	var state types.Object
	if ds := resp.State.Get(ctx, &state); ds.HasError() {
		t.Fatal(ds)
	}
	if value(stringInput(state.Attributes(), "id")) != "sdt-created" || value(stringInput(state.Attributes(), "status")) != "ACTIVE" || state.Attributes()["effective_configuration"].IsNull() {
		t.Fatal("create state is incomplete")
	}
}

func TestResourceCreateKnownDurationMatchesActualReadback(t *testing.T) {
	ctx := context.Background()
	config := modelFixture(t, map[string]any{"description": "v1", "default_timeout": "5m", "tags": map[string]string{"revision": "1"}, "timeouts": []map[string]any{{"create": "20ms", "update": nil, "delete": nil, "read": nil}}})
	plan := computedPlan(t, config)
	var roundTrip types.Object
	if ds := planFor(t, plan).Get(ctx, &roundTrip); ds.HasError() {
		t.Fatal(ds)
	}
	tool := remoteCode()
	tool.Description = pointer("v1")
	tool.Tags = append(tool.Tags, &sdk.Tag{Key: pointer("revision"), Value: pointer("1")})
	actual, err := flattenTool(ctx, roundTrip, tool, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"tool_name", "tool_type", "description", "default_timeout", "persistent", "role_arn", "tags", "network_configuration", "custom_configuration", "storage_mounts", "log_configuration"} {
		if !plannedEqual(roundTrip.Attributes()[key], actual.Attributes()[key]) {
			t.Errorf("mismatch %s: expected %T %s; actual %T %s", key, roundTrip.Attributes()[key], roundTrip.Attributes()[key], actual.Attributes()[key], actual.Attributes()[key])
		}
	}
}

func TestResourceUpdateDoesNotPretendDescriptionWasCleared(t *testing.T) {
	ctx := context.Background()
	prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), remoteCode(), true)
	if err != nil {
		t.Fatal(err)
	}
	fields := prior.Attributes()
	fields["description"] = types.StringValue("")
	timeouts, err := reflectedValue(ctx, inputAttributeTypes()["timeouts"], []map[string]any{{"create": nil, "update": "12ms", "delete": nil, "read": nil}})
	if err != nil {
		t.Fatal(err)
	}
	fields["timeouts"] = timeouts
	plan, ds := types.ObjectValue(inputAttributeTypes(), fields)
	if ds.HasError() {
		t.Fatal(ds)
	}
	calls := map[string]int{}
	client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
		calls[action]++
		if action == "UpdateSandboxTool" {
			if input["Description"] != "" {
				t.Fatal("clear was not sent")
			}
			if len(input) != 2 {
				t.Fatal("update sent unchanged fields")
			}
			return map[string]any{}, nil
		}
		if action == "DescribeSandboxToolList" {
			return pointResult(testTool("sdt-offline", "owned-code", "ACTIVE")), nil
		}
		t.Fatalf("unexpected action %s", action)
		return nil, nil
	})
	r := &SandboxToolResource{testService: fastService(client)}
	resp := resource.UpdateResponse{}
	r.Update(ctx, resource.UpdateRequest{State: stateFor(t, prior), Plan: planFor(t, plan)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("old ACTIVE description was accepted as a successful clear")
	}
	var actual types.Object
	if ds := resp.State.Get(ctx, &actual); ds.HasError() {
		t.Fatal(ds)
	}
	if value(stringInput(actual.Attributes(), "id")) != "sdt-offline" || value(stringInput(actual.Attributes(), "description")) != "old" {
		t.Fatal("recovery state does not contain actual readback")
	}
	if calls["UpdateSandboxTool"] != 1 || calls["CreateSandboxTool"] != 0 || calls["DeleteSandboxTool"] != 0 {
		t.Fatalf("unexpected mutations: %#v", calls)
	}
}

func TestResourceDeleteDependencyAndAmbiguousAbsenceRetainState(t *testing.T) {
	ctx := context.Background()
	prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), remoteCode(), true)
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"active-instance", "ambiguous-empty", "confirmed"} {
		t.Run(mode, func(t *testing.T) {
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				if action == "DeleteSandboxTool" {
					if mode == "active-instance" {
						return errorResult("UnsupportedOperation.SandboxToolHasActiveInstances"), nil
					}
					return map[string]any{}, nil
				}
				if action == "DescribeSandboxToolList" {
					if mode == "confirmed" {
						return errorResult("ResourceNotFound.SandboxTool"), nil
					}
					return map[string]any{"SandboxToolSet": []any{}, "TotalCount": 0}, nil
				}
				t.Fatalf("implicit instance action %s", action)
				return nil, nil
			})
			r := &SandboxToolResource{testService: fastService(client)}
			req := resource.DeleteRequest{State: stateFor(t, prior)}
			resp := resource.DeleteResponse{}
			r.Delete(ctx, req, &resp)
			if mode == "confirmed" {
				if resp.Diagnostics.HasError() || !resp.State.Raw.IsNull() {
					t.Fatal("confirmed delete did not complete")
				}
			} else if !resp.Diagnostics.HasError() || !resp.State.Raw.Equal(req.State.Raw) {
				t.Fatalf("delete failure lost state: %v", resp.Diagnostics)
			}
		})
	}
}

func TestResourceFailedReadAllowsStateRefreshForDestroy(t *testing.T) {
	ctx := context.Background()
	prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), remoteCode(), true)
	if err != nil {
		t.Fatal(err)
	}
	client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
		if action != "DescribeSandboxToolList" {
			t.Fatal("failed-state Read mutated resources")
		}
		return pointResult(testTool("sdt-offline", "owned-code", "FAILED")), nil
	})
	r := &SandboxToolResource{testService: fastService(client)}
	resp := resource.ReadResponse{}
	r.Read(ctx, resource.ReadRequest{State: stateFor(t, prior)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal("FAILED state was made impossible to destroy", resp.Diagnostics)
	}
	var state types.Object
	if ds := resp.State.Get(ctx, &state); ds.HasError() {
		t.Fatal(ds)
	}
	if value(stringInput(state.Attributes(), "status")) != "FAILED" {
		t.Fatal("failed status hidden")
	}
}

func TestResourceValidationRejectsInvalidConfigurationWithoutClient(t *testing.T) {
	for name, overrides := range map[string]map[string]any{"missing-network": {"network_configuration": []map[string]any{}}, "custom-without-image": {"tool_type": "custom"}, "unsupported-persistent": {"persistent": true}, "bad-timeout": {"default_timeout": "0.5s"}, "reserved-tag": {"tags": map[string]string{creatorSystemTag: "not-editable"}}, "duplicate-singleton": {"network_configuration": []map[string]any{{"network_mode": "PUBLIC", "vpc_config": []map[string]any{}}, {"network_mode": "PUBLIC", "vpc_config": []map[string]any{}}}}} {
		t.Run(name, func(t *testing.T) {
			r := &SandboxToolResource{}
			resp := resource.ValidateConfigResponse{}
			r.ValidateConfig(context.Background(), resource.ValidateConfigRequest{Config: configFor(t, modelFixture(t, overrides))}, &resp)
			if !resp.Diagnostics.HasError() {
				t.Fatal("invalid configuration accepted")
			}
		})
	}
}

func TestProvisionalStateContainsNoUnknownValues(t *testing.T) {
	ctx := context.Background()
	plan := computedPlan(t, modelFixture(t, nil))
	state, err := provisionalState(ctx, plan, "sdt-known")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := state.ToTerraformValue(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !raw.IsFullyKnown() {
		t.Fatal("failed create response still contains unknown values")
	}
	if v := state.Attributes()["id"]; v.(interface{ ValueString() string }).ValueString() != "sdt-known" {
		t.Fatal("provisional identity lost")
	}
}

func TestResourceDeleteAlreadyAbsentUsesOriginalReliableEvidence(t *testing.T) {
	ctx := context.Background()
	prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), remoteCode(), true)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
		calls++
		if action != "DeleteSandboxTool" {
			t.Fatal("confirmed absence was discarded and followed by an ambiguous query")
		}
		return errorResult("ResourceNotFound.SandboxTool"), nil
	})
	r := &SandboxToolResource{testService: fastService(client)}
	resp := resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{State: stateFor(t, prior)}, &resp)
	if resp.Diagnostics.HasError() || !resp.State.Raw.IsNull() || calls != 1 {
		t.Fatalf("already absent delete failed: calls=%d diagnostics=%v", calls, resp.Diagnostics)
	}
}

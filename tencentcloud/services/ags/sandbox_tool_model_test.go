package ags

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

func modelFixture(t *testing.T, overrides map[string]any) types.Object {
	t.Helper()
	raw := map[string]any{}
	for key := range inputAttributeTypes() {
		raw[key] = nil
	}
	raw["tool_name"], raw["tool_type"] = "owned-code", "code-interpreter"
	raw["network_configuration"] = []map[string]any{{"network_mode": "PUBLIC", "vpc_config": []map[string]any{}}}
	for _, key := range []string{"custom_configuration", "storage_mounts", "log_configuration", "timeouts"} {
		raw[key] = []map[string]any{}
	}
	for key, value := range overrides {
		raw[key] = value
	}
	v, err := reflectedValue(context.Background(), toolSchema().Type(), raw)
	if err != nil {
		t.Fatal(err)
	}
	return v.(types.Object)
}

func remoteCode() *sdk.SandboxTool {
	return &sdk.SandboxTool{ToolId: pointer("sdt-offline"), ToolName: pointer("owned-code"), ToolType: pointer("code-interpreter"), Status: pointer("ACTIVE"), Description: pointer("description"), DefaultTimeoutSeconds: pointer(uint64(300)), Persistent: pointer(false), NetworkConfiguration: &sdk.NetworkConfiguration{NetworkMode: pointer("PUBLIC")}, Tags: []*sdk.Tag{{Key: pointer(creatorSystemTag), Value: pointer("system")}}}
}

func TestModelImportRoundTripAndTypedDuration(t *testing.T) {
	ctx := context.Background()
	tool := remoteCode()
	state, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), tool, true)
	if err != nil {
		t.Fatal(err)
	}
	if state.Attributes()["default_timeout"].(semanticStringValue).ValueString() != "300s" {
		t.Fatal("missing typed timeout")
	}
	if got := state.Attributes()["tags"].(types.Map); len(got.Elements()) != 0 {
		t.Fatal("system tag leaked into managed tags")
	}
	effective := state.Attributes()["effective_configuration"].(types.Object)
	if len(effective.Attributes()["tags"].(types.Map).Elements()) != 1 {
		t.Fatal("system tag lost from actual configuration")
	}
	req, err := createRequest(state)
	if err != nil {
		t.Fatal(err)
	}
	if value(req.ToolName) != "owned-code" || value(req.DefaultTimeout) != "300s" || value(req.NetworkConfiguration.NetworkMode) != "PUBLIC" {
		t.Fatal("imported input mapping is incomplete")
	}
}

func TestModelReadPreservesEquivalentDurationAndShowsRealDrift(t *testing.T) {
	ctx := context.Background()
	prior := modelFixture(t, map[string]any{"default_timeout": "5m", "description": "old"})
	state, err := flattenTool(ctx, prior, remoteCode(), false)
	if err != nil {
		t.Fatal(err)
	}
	if got := state.Attributes()["default_timeout"].(semanticStringValue).ValueString(); got != "5m" {
		t.Fatalf("equivalent duration changed to %s", got)
	}
	if value(stringInput(state.Attributes(), "description")) != "description" {
		t.Fatal("external description drift was hidden")
	}
	if len(state.Attributes()["custom_configuration"].(types.List).Elements()) != 0 {
		t.Fatal("invented unconfigured custom block")
	}
}

func TestModelExplicitEmptyDescriptionAndTagsRemainWrites(t *testing.T) {
	ctx := context.Background()
	before := modelFixture(t, map[string]any{"id": "sdt-offline", "description": "old", "tags": map[string]string{"owner": "test"}})
	remote := remoteCode()
	remote.Tags = append(remote.Tags, &sdk.Tag{Key: pointer("owner"), Value: pointer("test")})
	before, err := flattenTool(ctx, before, remote, false)
	if err != nil {
		t.Fatal(err)
	}
	fields := before.Attributes()
	fields["description"] = types.StringValue("")
	fields["tags"] = types.MapValueMust(types.StringType, map[string]attr.Value{})
	after, ds := types.ObjectValue(inputAttributeTypes(), fields)
	if ds.HasError() {
		t.Fatal(ds)
	}
	req, changed, err := updateRequest(before, after)
	if err != nil {
		t.Fatal(err)
	}
	if !changed || req.Description == nil || *req.Description != "" || req.Tags == nil || len(req.Tags) != 0 {
		t.Fatal("explicit clears became omissions")
	}
}

func TestModelIncompleteReadCannotOverwriteState(t *testing.T) {
	tool := remoteCode()
	tool.DefaultTimeoutSeconds = nil
	if _, err := flattenTool(context.Background(), modelFixture(t, nil), tool, false); err == nil {
		t.Fatal("accepted missing nonpersistent timeout")
	}
}

func TestModelProtocolFiveNullBlocksMatchAbsentRemoteBlocks(t *testing.T) {
	ctx := context.Background()
	plan := modelFixture(t, map[string]any{
		"description": "description", "default_timeout": "5m", "persistent": false,
		"tags": map[string]string{}, "custom_configuration": nil,
		"storage_mounts": nil, "log_configuration": nil,
		"network_configuration": []map[string]any{{"network_mode": "PUBLIC", "vpc_config": nil}},
	})
	state, err := flattenTool(ctx, plan, remoteCode(), false)
	if err != nil {
		t.Fatal(err)
	}
	if pending, err := pendingPlanFields(ctx, plan, plan, remoteCode()); err != nil || len(pending) != 0 {
		t.Fatalf("Protocol 5 absent blocks prevent creation: %v %v", pending, err)
	}
	for _, key := range []string{"custom_configuration", "storage_mounts", "log_configuration"} {
		if !state.Attributes()[key].IsNull() {
			t.Fatalf("%s did not retain the absent block representation", key)
		}
	}
	// Scalar reset lists are not blocks: null and [] remain distinct.
	if plannedEqual(types.ListNull(types.StringType), types.ListValueMust(types.StringType, nil)) {
		t.Fatal("scalar list reset was suppressed")
	}
}

func TestModelPortRemovalMustVerifyRemoteClear(t *testing.T) {
	ctx := context.Background()
	custom := map[string]any{
		"image": "offline.invalid/fixture:1", "image_registry_type": "enterprise",
		"command": []string{}, "args": []string{}, "env": map[string]string{},
		"resources": []map[string]any{{"cpu": "1", "memory": "1Gi"}},
		"ports":     []map[string]any{{"name": "http", "port": int64(8080), "protocol": "TCP"}},
	}
	prior := modelFixture(t, map[string]any{"id": "sdt-offline", "tool_type": "custom", "custom_configuration": []map[string]any{custom}})
	removed := cloneMap(custom)
	removed["ports"] = nil
	plan := modelFixture(t, map[string]any{"id": "sdt-offline", "tool_type": "custom", "custom_configuration": []map[string]any{removed}})
	request, changed, err := updateRequest(prior, plan)
	if err != nil || !changed || request.CustomConfiguration == nil || request.CustomConfiguration.Ports == nil || len(request.CustomConfiguration.Ports) != 0 {
		t.Fatalf("removal did not send an explicit empty collection: changed=%v err=%v", changed, err)
	}
	shape, err := ownershipForUpdate(ctx, prior, plan)
	if err != nil {
		t.Fatal(err)
	}
	tool := remoteCode()
	tool.ToolType = pointer("custom")
	tool.CustomConfiguration = &sdk.CustomConfigurationDetail{
		Image: pointer("offline.invalid/fixture:1"), ImageRegistryType: pointer("enterprise"),
		Resources: &sdk.ResourceConfiguration{CPU: pointer("1"), Memory: pointer("1Gi")},
		Ports:     []*sdk.PortConfiguration{{Name: pointer("http"), Port: pointer(int64(8080)), Protocol: pointer("TCP")}},
	}
	actual, err := flattenTool(ctx, shape, tool, false)
	if err != nil {
		t.Fatal(err)
	}
	if plannedEqual(plan.Attributes()["custom_configuration"], actual.Attributes()["custom_configuration"]) {
		t.Fatal("old remote port was hidden by the removed block")
	}
	tool.CustomConfiguration.Ports = nil
	actual, err = flattenTool(ctx, shape, tool, false)
	if err != nil {
		t.Fatal(err)
	}
	if !plannedEqual(plan.Attributes()["custom_configuration"], actual.Attributes()["custom_configuration"]) {
		t.Fatal("confirmed empty ports did not converge")
	}
}

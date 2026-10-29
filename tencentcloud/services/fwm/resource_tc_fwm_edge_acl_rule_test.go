package fwm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

// mockMeta implements tccommon.ProviderMeta
type mockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMeta{}

func newMockMeta() *mockMeta {
	return &mockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string {
	return &s
}

func ptrInt64(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/fwm/ -run "TestFwmEdgeAclRule" -v -count=1

const edgeAclRuleId = "cfw-group-xxxx#cfw-rule-xxxx#1"

func buildEdgeAclRuleResp() *fwmv20250611.EdgeAclRuleResp {
	return &fwmv20250611.EdgeAclRuleResp{
		RuleId:        ptrString("cfw-rule-xxxx"),
		Sequence:      ptrInt64(1),
		Direction:     ptrInt64(1),
		SourceContent: ptrString("1.1.1.1"),
		SourceType:    ptrString("ip"),
		TargetContent: ptrString("2.2.2.2"),
		TargetType:    ptrString("ip"),
		Port:          ptrString("80"),
		Protocol:      ptrString("TCP"),
		RuleAction:    ptrString("accept-pass"),
		Scope:         ptrString("all"),
		BelongMember: &fwmv20250611.MemberInfo{
			MemberId: ptrString("mem-xxxx"),
		},
	}
}

func mockDescribeEdgeAclRules(t *testing.T, patches *gomonkey.Patches, client *fwmv20250611.Client, rules ...*fwmv20250611.EdgeAclRuleResp) {
	patches.ApplyMethodFunc(client, "DescribeEdgeAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		if request != nil && request.GroupId != nil {
			assert.Equal(t, "cfw-group-xxxx", *request.GroupId)
		}
		if request != nil && request.Direction != nil {
			assert.Equal(t, int64(1), *request.Direction)
		}
		if request != nil && request.Limit != nil {
			assert.Equal(t, int64(1000), *request.Limit)
		}

		resp := fwmv20250611.NewDescribeEdgeAclRulesResponse()
		total := int64(len(rules))
		resp.Response = &fwmv20250611.DescribeEdgeAclRulesResponseParams{
			TotalCount: &total,
			Rules:      rules,
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})
}

func TestFwmEdgeAclRule_Schema(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	groupId := res.Schema["group_id"]
	assert.NotNil(t, groupId)
	assert.True(t, groupId.Required)
	assert.True(t, groupId.ForceNew)

	direction := res.Schema["direction"]
	assert.NotNil(t, direction)
	assert.Equal(t, schema.TypeInt, direction.Type)
	assert.True(t, direction.Required)
	assert.True(t, direction.ForceNew)

	sourceType := res.Schema["source_type"]
	assert.NotNil(t, sourceType)
	assert.True(t, sourceType.Optional)
	assert.True(t, sourceType.ForceNew)

	targetType := res.Schema["target_type"]
	assert.NotNil(t, targetType)
	assert.True(t, targetType.Optional)
	assert.True(t, targetType.ForceNew)

	ruleId := res.Schema["rule_id"]
	assert.NotNil(t, ruleId)
	assert.True(t, ruleId.Computed)

	sequence := res.Schema["sequence"]
	assert.NotNil(t, sequence)
	assert.True(t, sequence.Computed)

	// order_index is no longer part of the schema; sequence is read-only.
	_, ok := res.Schema["order_index"]
	assert.False(t, ok)
}

// TestFwmEdgeAclRule_Create verifies the composite id now carries direction.
func TestFwmEdgeAclRule_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "cfw-group-xxxx", *request.GroupId)
		assert.Len(t, request.Rules, 1)

		rule := request.Rules[0]
		assert.Equal(t, int64(1), *rule.OrderIndex)
		assert.Equal(t, int64(1), *rule.Direction)
		assert.Equal(t, "ip", *rule.SourceType)
		assert.Equal(t, "accept-pass", *rule.RuleAction)
		assert.Equal(t, "all", *rule.Scope)

		resp := fwmv20250611.NewCreateEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.CreateEdgeAclRuleResponseParams{
			RuleIds:   []*string{ptrString("cfw-rule-xxxx")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribeEdgeAclRules(t, patches, fwmClient, buildEdgeAclRuleResp())

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":       "cfw-group-xxxx",
		"direction":      1,
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "2.2.2.2",
		"target_type":    "ip",
		"port":           "80",
		"protocol":       "TCP",
		"rule_action":    "accept-pass",
		"scope":          "all",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	// groupId#ruleId#direction
	assert.Equal(t, edgeAclRuleId, d.Id())
	assert.Equal(t, "cfw-rule-xxxx", d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("sequence"))
}

// TestFwmEdgeAclRule_Create_EmptyRuleIds verifies an empty RuleIds fails Create.
func TestFwmEdgeAclRule_Create_EmptyRuleIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, _ *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		resp := fwmv20250611.NewCreateEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.CreateEdgeAclRuleResponseParams{
			RuleIds:   []*string{},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

// TestFwmEdgeAclRule_Create_APIError verifies Create surfaces API errors.
func TestFwmEdgeAclRule_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, _ *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid group_id")
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-invalid",
		"direction": 1,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestFwmEdgeAclRule_Read verifies Read parses the 3-segment id, queries with
// direction and maps the response back into state.
func TestFwmEdgeAclRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeEdgeAclRules(t, patches, fwmClient, buildEdgeAclRuleResp())

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId(edgeAclRuleId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, edgeAclRuleId, d.Id())
	assert.Equal(t, "cfw-rule-xxxx", d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("sequence"))
	assert.Equal(t, 1, d.Get("direction"))
	assert.Equal(t, "mem-xxxx", d.Get("belong_member_id"))
}

// TestFwmEdgeAclRule_Read_NotFound verifies Read clears id when the rule is gone.
func TestFwmEdgeAclRule_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeEdgeAclRules(t, patches, fwmClient)

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId(edgeAclRuleId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmEdgeAclRule_Read_LegacyId verifies the old 2-segment id is rejected.
func TestFwmEdgeAclRule_Read_LegacyId(t *testing.T) {
	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId("cfw-group-xxxx#cfw-rule-xxxx")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestFwmEdgeAclRule_Read_InvalidDirection verifies a non-numeric direction
// segment is reported instead of being silently ignored.
func TestFwmEdgeAclRule_Read_InvalidDirection(t *testing.T) {
	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId("cfw-group-xxxx#cfw-rule-xxxx#ingress")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestFwmEdgeAclRule_Update verifies Update parses the 3-segment id and sends
// the mutable fields.
func TestFwmEdgeAclRule_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	var capturedRequest *fwmv20250611.ModifyEdgeAclRuleRequest
	patches.ApplyMethodFunc(fwmClient, "ModifyEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyEdgeAclRuleRequest) (*fwmv20250611.ModifyEdgeAclRuleResponse, error) {
		capturedRequest = request
		resp := fwmv20250611.NewModifyEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyEdgeAclRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribeEdgeAclRules(t, patches, fwmClient, func() *fwmv20250611.EdgeAclRuleResp {
		rule := buildEdgeAclRuleResp()
		rule.Port = ptrString("8080")
		rule.RuleAction = ptrString("drop-block")
		rule.Description = ptrString("updated")
		return rule
	}())

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()

	oldData := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":       "cfw-group-xxxx",
		"direction":      1,
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "2.2.2.2",
		"target_type":    "ip",
		"port":           "80",
		"protocol":       "TCP",
		"rule_action":    "accept-pass",
		"scope":          "all",
	})
	oldData.SetId(edgeAclRuleId)
	state := oldData.State()

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"group_id":       "cfw-group-xxxx",
		"direction":      1,
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "2.2.2.2",
		"target_type":    "ip",
		"port":           "8080",
		"protocol":       "TCP",
		"rule_action":    "drop-block",
		"description":    "updated",
		"scope":          "all",
	})
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "cfw-group-xxxx", *capturedRequest.GroupId)
	assert.NotNil(t, capturedRequest.Rule)
	assert.Equal(t, "cfw-rule-xxxx", *capturedRequest.Rule.RuleId)
	assert.Equal(t, "8080", *capturedRequest.Rule.Port)
	assert.Equal(t, "drop-block", *capturedRequest.Rule.RuleAction)
	assert.Equal(t, "updated", *capturedRequest.Rule.Description)
}

// TestFwmEdgeAclRule_Delete verifies Delete parses the 3-segment id.
func TestFwmEdgeAclRule_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteEdgeAclRuleRequest) (*fwmv20250611.DeleteEdgeAclRuleResponse, error) {
		assert.Equal(t, "cfw-group-xxxx", *request.GroupId)
		assert.Len(t, request.RuleIds, 1)
		assert.Equal(t, "cfw-rule-xxxx", *request.RuleIds[0])

		resp := fwmv20250611.NewDeleteEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.DeleteEdgeAclRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId(edgeAclRuleId)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFwmEdgeAclRule_Delete_NotFound verifies Delete is idempotent on ResourceNotFound.
func TestFwmEdgeAclRule_Delete_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteEdgeAclRuleWithContext", func(_ context.Context, _ *fwmv20250611.DeleteEdgeAclRuleRequest) (*fwmv20250611.DeleteEdgeAclRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=rule not found")
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId(edgeAclRuleId)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

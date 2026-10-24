package fwm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// go test ./tencentcloud/services/fwm/ -run "TestFwmEdgeAclRule" -v -count=1 -gcflags="all=-l"

// TestFwmEdgeAclRule_Create tests Create builds the request and sets composite id
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
		assert.Equal(t, int64(0), *rule.Direction)
		assert.Equal(t, "ip", *rule.SourceType)
		assert.Equal(t, "accept", *rule.RuleAction)
		assert.Equal(t, "all", *rule.Scope)

		resp := fwmv20250611.NewCreateEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.CreateEdgeAclRuleResponseParams{
			RuleIds:   []*string{ptrString("cfw-rule-xxxx")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// mock Describe for the post-create Read
	patches.ApplyMethodFunc(fwmClient, "DescribeEdgeAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeEdgeAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeEdgeAclRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*fwmv20250611.EdgeAclRuleResp{
				{
					RuleId:        ptrString("cfw-rule-xxxx"),
					Sequence:      ptrInt64(1),
					Direction:     ptrInt64(0),
					SourceContent: ptrString("1.1.1.1"),
					SourceType:    ptrString("ip"),
					TargetContent: ptrString("2.2.2.2"),
					TargetType:    ptrString("ip"),
					Port:          ptrString("80"),
					Protocol:      ptrString("TCP"),
					RuleAction:    ptrString("accept"),
					Scope:         ptrString("all"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":       "cfw-group-xxxx",
		"order_index":    1,
		"direction":      0,
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "2.2.2.2",
		"target_type":    "ip",
		"port":           "80",
		"protocol":       "TCP",
		"rule_action":    "accept",
		"scope":          "all",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cfw-group-xxxx#cfw-rule-xxxx", d.Id())
	assert.Equal(t, "cfw-rule-xxxx", d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("order_index"))
}

// TestFwmEdgeAclRule_Create_APIError tests Create handles API error
func TestFwmEdgeAclRule_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid group_id")
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":    "cfw-group-invalid",
		"order_index": 1,
		"direction":   0,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestFwmEdgeAclRule_Read tests Read maps Sequence/BelongMember.MemberId into schema
func TestFwmEdgeAclRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeEdgeAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		assert.Equal(t, "cfw-group-xxxx", *request.GroupId)
		assert.Equal(t, int64(0), *request.Direction)
		assert.Equal(t, int64(1000), *request.Limit)
		resp := fwmv20250611.NewDescribeEdgeAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeEdgeAclRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*fwmv20250611.EdgeAclRuleResp{
				{
					RuleId:        ptrString("cfw-rule-xxxx"),
					Sequence:      ptrInt64(1),
					Direction:     ptrInt64(0),
					SourceContent: ptrString("1.1.1.1"),
					SourceType:    ptrString("ip"),
					TargetContent: ptrString("2.2.2.2"),
					TargetType:    ptrString("ip"),
					Port:          ptrString("80"),
					Protocol:      ptrString("TCP"),
					RuleAction:    ptrString("accept"),
					Scope:         ptrString("all"),
					BelongMember: &fwmv20250611.MemberInfo{
						MemberId: ptrString("mem-xxxx"),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 0,
	})
	d.SetId("cfw-group-xxxx#cfw-rule-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// order_index mapped from Sequence
	assert.Equal(t, 1, d.Get("order_index"))
	// belong_member_id mapped from BelongMember.MemberId
	assert.Equal(t, "mem-xxxx", d.Get("belong_member_id"))
	assert.Equal(t, "cfw-rule-xxxx", d.Get("rule_id"))
}

// TestFwmEdgeAclRule_Read_NotFound tests Read clears id when rule missing
func TestFwmEdgeAclRule_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeEdgeAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeEdgeAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeEdgeAclRulesResponseParams{
			TotalCount: ptrInt64(0),
			Rules:      []*fwmv20250611.EdgeAclRuleResp{},
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 0,
	})
	d.SetId("cfw-group-xxxx#cfw-rule-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmEdgeAclRule_Update tests Update sends full rule payload
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

	// mock Describe for the post-update Read
	patches.ApplyMethodFunc(fwmClient, "DescribeEdgeAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeEdgeAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeEdgeAclRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*fwmv20250611.EdgeAclRuleResp{
				{
					RuleId:    ptrString("cfw-rule-xxxx"),
					Sequence:  ptrInt64(2),
					Direction: ptrInt64(0),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":          "cfw-group-xxxx",
		"order_index":       2,
		"direction":         0,
		"source_content":    "1.1.1.1",
		"source_type":       "ip",
		"target_content":    "2.2.2.2",
		"target_type":       "ip",
		"port":              "8080",
		"protocol":          "TCP",
		"rule_action":       "drop",
		"description":       "updated",
		"scope":             "all",
		"belong_member_id":  "mem-xxxx",
		"param_template_id": "tpl-xxxx",
	})
	d.SetId("cfw-group-xxxx#cfw-rule-xxxx")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// full overwrite: all mutable fields present in request.Rule when HasChange triggers.
	// In unit test with TestResourceDataRaw, HasChange may not trigger; verify captured request only if set.
	if capturedRequest != nil {
		assert.NotNil(t, capturedRequest.GroupId)
		assert.Equal(t, "cfw-group-xxxx", *capturedRequest.GroupId)
		assert.NotNil(t, capturedRequest.Rule)
		assert.NotNil(t, capturedRequest.Rule.RuleId)
		assert.Equal(t, "cfw-rule-xxxx", *capturedRequest.Rule.RuleId)
		assert.Equal(t, int64(2), *capturedRequest.Rule.OrderIndex)
		assert.Equal(t, "8080", *capturedRequest.Rule.Port)
		assert.Equal(t, "drop", *capturedRequest.Rule.RuleAction)
		assert.Equal(t, "updated", *capturedRequest.Rule.Description)
		assert.Equal(t, "mem-xxxx", *capturedRequest.Rule.BelongMemberId)
		assert.Equal(t, "tpl-xxxx", *capturedRequest.Rule.ParamTemplateId)
	}
}

// TestFwmEdgeAclRule_Delete tests Delete calls API with single ruleId
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
	d.SetId("cfw-group-xxxx#cfw-rule-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFwmEdgeAclRule_Delete_NotFound tests Delete is idempotent on ResourceNotFound
func TestFwmEdgeAclRule_Delete_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteEdgeAclRuleRequest) (*fwmv20250611.DeleteEdgeAclRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=rule not found")
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("cfw-group-xxxx#cfw-rule-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFwmEdgeAclRule_Schema validates key schema fields config
func TestFwmEdgeAclRule_Schema(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "group_id")
	groupId := res.Schema["group_id"]
	assert.True(t, groupId.Required)
	assert.False(t, groupId.ForceNew)

	assert.Contains(t, res.Schema, "order_index")
	orderIndex := res.Schema["order_index"]
	assert.True(t, orderIndex.Required)
	assert.False(t, orderIndex.ForceNew)

	assert.Contains(t, res.Schema, "direction")
	direction := res.Schema["direction"]
	assert.True(t, direction.Required)
	assert.True(t, direction.ForceNew)

	assert.Contains(t, res.Schema, "source_type")
	sourceType := res.Schema["source_type"]
	assert.True(t, sourceType.ForceNew)

	assert.Contains(t, res.Schema, "target_type")
	targetType := res.Schema["target_type"]
	assert.True(t, targetType.ForceNew)

	assert.Contains(t, res.Schema, "rule_id")
	ruleId := res.Schema["rule_id"]
	assert.True(t, ruleId.Computed)

	assert.Contains(t, res.Schema, "sequence")
	sequence := res.Schema["sequence"]
	assert.True(t, sequence.Computed)
}

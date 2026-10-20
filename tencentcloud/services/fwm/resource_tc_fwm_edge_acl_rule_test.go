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

type mockMetaFwmEdgeAclRule struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmEdgeAclRule) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmEdgeAclRule{}

func newMockMetaFwmEdgeAclRule() *mockMetaFwmEdgeAclRule {
	return &mockMetaFwmEdgeAclRule{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFwmEdge(s string) *string { return &s }
func ptrInt64FwmEdge(i int64) *int64    { return &i }

// mockDescribeEdgeAclRules patches DescribeEdgeAclRulesWithContext to return
// rules with the given rule ids.
func mockDescribeEdgeAclRules(patches *gomonkey.Patches, fwmClient *fwmv20250611.Client, ruleIds []string) {
	rules := make([]*fwmv20250611.EdgeAclRuleResp, 0, len(ruleIds))
	for _, id := range ruleIds {
		rule := &fwmv20250611.EdgeAclRuleResp{
			RuleId:        ptrStringFwmEdge(id),
			Sequence:      ptrInt64FwmEdge(int64(len(rules) + 1)),
			Direction:     ptrInt64FwmEdge(int64(1)),
			SourceContent: ptrStringFwmEdge("1.1.1.1/32"),
			SourceType:    ptrStringFwmEdge("ip"),
			TargetContent: ptrStringFwmEdge("0.0.0.0/0"),
			TargetType:    ptrStringFwmEdge("ip"),
			Port:          ptrStringFwmEdge("80"),
			Protocol:      ptrStringFwmEdge("TCP"),
			RuleAction:    ptrStringFwmEdge("accept"),
			Description:   ptrStringFwmEdge("example rule"),
			Scope:         ptrStringFwmEdge("serial"),
		}
		rules = append(rules, rule)
	}
	patches.ApplyMethodFunc(fwmClient, "DescribeEdgeAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeEdgeAclRulesResponse()
		totalCount := int64(len(ruleIds))
		resp.Response = &fwmv20250611.DescribeEdgeAclRulesResponseParams{
			TotalCount: &totalCount,
			Rules:      rules,
			RequestId:  ptrStringFwmEdge("fake-request-id"),
		}
		return resp, nil
	})
}

// TestFwmEdgeAclRule_CreateBackfillRuleIds tests Create returns RuleIds into state
func TestFwmEdgeAclRule_CreateBackfillRuleIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmEdgeAclRule().client, "UseFwmV20250611Client", fwmClient)

	ruleId := "rule-001"
	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "group-xxxx", *request.GroupId)
		resp := fwmv20250611.NewCreateEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.CreateEdgeAclRuleResponseParams{
			RuleIds:   []*string{&ruleId},
			RequestId: ptrStringFwmEdge("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribeEdgeAclRules(patches, fwmClient, []string{"rule-001"})

	meta := newMockMetaFwmEdgeAclRule()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id": "group-xxxx",
		"rules": []interface{}{
			map[string]interface{}{
				"order_index":    1,
				"direction":      1,
				"source_content": "1.1.1.1/32",
				"source_type":    "ip",
				"target_content": "0.0.0.0/0",
				"target_type":    "ip",
				"port":           "80",
				"protocol":       "TCP",
				"rule_action":    "accept",
				"description":    "example rule",
				"scope":          "serial",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "group-xxxx#rule-001", d.Id())

	ruleIds := d.Get("rule_ids").([]interface{})
	assert.Len(t, ruleIds, 1)
	assert.Equal(t, "rule-001", ruleIds[0].(string))
}

// TestFwmEdgeAclRule_CreateEmptyRuleIdsError tests Create fails when RuleIds is empty
func TestFwmEdgeAclRule_CreateEmptyRuleIdsError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmEdgeAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		resp := fwmv20250611.NewCreateEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.CreateEdgeAclRuleResponseParams{
			RuleIds:   []*string{},
			RequestId: ptrStringFwmEdge("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmEdgeAclRule()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id": "group-xxxx",
		"rules": []interface{}{
			map[string]interface{}{
				"order_index": 1,
				"direction":   1,
				"protocol":    "TCP",
				"rule_action": "accept",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

// TestFwmEdgeAclRule_ReadRefreshRuleIds tests Read refreshes RuleIds
func TestFwmEdgeAclRule_ReadRefreshRuleIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmEdgeAclRule().client, "UseFwmV20250611Client", fwmClient)

	mockDescribeEdgeAclRules(patches, fwmClient, []string{"rule-001", "rule-002"})

	meta := newMockMetaFwmEdgeAclRule()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id": "group-xxxx",
	})
	d.SetId("group-xxxx#rule-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	ruleIds := d.Get("rule_ids").([]interface{})
	assert.Len(t, ruleIds, 2)
	assert.Equal(t, "rule-001", ruleIds[0].(string))
	assert.Equal(t, "rule-002", ruleIds[1].(string))
}

// TestFwmEdgeAclRule_ReadEmptySetId tests Read clears id when no rules found
func TestFwmEdgeAclRule_ReadEmptySetId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmEdgeAclRule().client, "UseFwmV20250611Client", fwmClient)

	mockDescribeEdgeAclRules(patches, fwmClient, []string{})

	meta := newMockMetaFwmEdgeAclRule()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id": "group-xxxx",
	})
	d.SetId("group-xxxx#rule-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Id())
}

// TestFwmEdgeAclRule_Delete tests Delete calls API with rule ids
func TestFwmEdgeAclRule_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmEdgeAclRule().client, "UseFwmV20250611Client", fwmClient)

	deleteCalled := false
	patches.ApplyMethodFunc(fwmClient, "DeleteEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteEdgeAclRuleRequest) (*fwmv20250611.DeleteEdgeAclRuleResponse, error) {
		deleteCalled = true
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "group-xxxx", *request.GroupId)
		assert.Len(t, request.RuleIds, 1)
		resp := fwmv20250611.NewDeleteEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.DeleteEdgeAclRuleResponseParams{
			RequestId: ptrStringFwmEdge("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmEdgeAclRule()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id": "group-xxxx",
	})
	d.SetId("group-xxxx#rule-001")
	_ = d.Set("rule_ids", []interface{}{"rule-001"})

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}

// TestFwmEdgeAclRule_CreateAPIError tests Create handles API error
func TestFwmEdgeAclRule_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmEdgeAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation, Message=group not exist")
	})

	meta := newMockMetaFwmEdgeAclRule()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id": "group-xxxx",
		"rules": []interface{}{
			map[string]interface{}{
				"order_index": 1,
				"direction":   1,
				"protocol":    "TCP",
				"rule_action": "accept",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FailedOperation")
}

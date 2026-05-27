package ga2_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
)

type mockMetaGa2ForwardingRule struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaGa2ForwardingRule) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaGa2ForwardingRule{}

func newMockMetaGa2ForwardingRule() *mockMetaGa2ForwardingRule {
	return &mockMetaGa2ForwardingRule{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStringGa2FR(s string) *string {
	return &s
}

func ptrBoolGa2FR(b bool) *bool {
	return &b
}

// go test ./tencentcloud/services/ga2/ -run "TestUnitGa2ForwardingRule" -v -count=1 -gcflags="all=-l"

func TestUnitGa2ForwardingRule_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2ForwardingRule()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch CreateForwardingRuleWithContext
	patches.ApplyMethodFunc(ga2Client, "CreateForwardingRuleWithContext", func(_ context.Context, request *ga2v20250115.CreateForwardingRuleRequest) (*ga2v20250115.CreateForwardingRuleResponse, error) {
		assert.Equal(t, "ga2-test001", *request.GlobalAcceleratorId)
		assert.Equal(t, "lis-test001", *request.ListenerId)
		assert.Equal(t, "fp-test001", *request.ForwardingPolicyId)
		assert.Equal(t, 1, len(request.RuleConditions))
		assert.Equal(t, "HOST", *request.RuleConditions[0].RuleConditionType)
		assert.Equal(t, 1, len(request.RuleActions))
		assert.Equal(t, "ForwardGroup", *request.RuleActions[0].RuleActionType)
		assert.Equal(t, true, *request.EnableOriginSni)
		assert.Equal(t, "example.com", *request.OriginSni)
		assert.Equal(t, "origin.example.com", *request.OriginHost)

		resp := &ga2v20250115.CreateForwardingRuleResponse{}
		resp.Response = &ga2v20250115.CreateForwardingRuleResponseParams{
			TaskId:           ptrStringGa2FR("task-001"),
			ForwardingRuleId: ptrStringGa2FR("fr-test001"),
			RequestId:        ptrStringGa2FR("req-001"),
		}
		return resp, nil
	})

	// Patch DescribeTaskResultWithContext for WaitForGa2TaskFinish
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(_ context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := &ga2v20250115.DescribeTaskResultResponse{}
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2FR("SUCCESS"),
			RequestId: ptrStringGa2FR("req-002"),
		}
		return resp, nil
	})

	// Patch DescribeForwardingRuleWithContext for Read after Create
	patches.ApplyMethodFunc(ga2Client, "DescribeForwardingRuleWithContext", func(_ context.Context, request *ga2v20250115.DescribeForwardingRuleRequest) (*ga2v20250115.DescribeForwardingRuleResponse, error) {
		resp := &ga2v20250115.DescribeForwardingRuleResponse{}
		resp.Response = &ga2v20250115.DescribeForwardingRuleResponseParams{
			ForwardingRuleSet: []*ga2v20250115.ForwardingRuleSet{
				{
					ForwardingRuleId:    ptrStringGa2FR("fr-test001"),
					GlobalAcceleratorId: ptrStringGa2FR("ga2-test001"),
					ListenerId:          ptrStringGa2FR("lis-test001"),
					ForwardingPolicyId:  ptrStringGa2FR("fp-test001"),
					RuleCondition: []*ga2v20250115.RuleCondition{
						{
							RuleConditionType:  ptrStringGa2FR("HOST"),
							RuleConditionValue: []*string{ptrStringGa2FR("example.com")},
						},
					},
					RuleAction: []*ga2v20250115.RuleAction{
						{
							RuleActionType:  ptrStringGa2FR("ForwardGroup"),
							RuleActionValue: ptrStringGa2FR("eg-test001"),
						},
					},
					OriginHeaders: []*ga2v20250115.OriginHeader{
						{
							Key:   ptrStringGa2FR("X-Custom-Header"),
							Value: ptrStringGa2FR("custom-value"),
						},
					},
					EnableOriginSni: ptrBoolGa2FR(true),
					OriginSni:       ptrStringGa2FR("example.com"),
					OriginHost:      ptrStringGa2FR("origin.example.com"),
				},
			},
			TotalCount: func() *uint64 { v := uint64(1); return &v }(),
			RequestId:  ptrStringGa2FR("req-003"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2ForwardingRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test001",
		"listener_id":           "lis-test001",
		"forwarding_policy_id":  "fp-test001",
		"rule_conditions": []interface{}{
			map[string]interface{}{
				"rule_condition_type":  "HOST",
				"rule_condition_value": []interface{}{"example.com"},
			},
		},
		"rule_actions": []interface{}{
			map[string]interface{}{
				"rule_action_type":  "ForwardGroup",
				"rule_action_value": "eg-test001",
			},
		},
		"origin_headers": []interface{}{
			map[string]interface{}{
				"key":   "X-Custom-Header",
				"value": "custom-value",
			},
		},
		"enable_origin_sni": true,
		"origin_sni":        "example.com",
		"origin_host":       "origin.example.com",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga2-test001"+tccommon.FILED_SP+"lis-test001"+tccommon.FILED_SP+"fp-test001"+tccommon.FILED_SP+"fr-test001", d.Id())
}

func TestUnitGa2ForwardingRule_CreateNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2ForwardingRule()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch CreateForwardingRuleWithContext to return nil response
	patches.ApplyMethodFunc(ga2Client, "CreateForwardingRuleWithContext", func(_ context.Context, request *ga2v20250115.CreateForwardingRuleRequest) (*ga2v20250115.CreateForwardingRuleResponse, error) {
		resp := &ga2v20250115.CreateForwardingRuleResponse{}
		resp.Response = nil
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2ForwardingRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test001",
		"listener_id":           "lis-test001",
		"forwarding_policy_id":  "fp-test001",
		"rule_conditions": []interface{}{
			map[string]interface{}{
				"rule_condition_type":  "HOST",
				"rule_condition_value": []interface{}{"example.com"},
			},
		},
		"rule_actions": []interface{}{
			map[string]interface{}{
				"rule_action_type":  "ForwardGroup",
				"rule_action_value": "eg-test001",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

func TestUnitGa2ForwardingRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2ForwardingRule()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch DescribeForwardingRuleWithContext
	patches.ApplyMethodFunc(ga2Client, "DescribeForwardingRuleWithContext", func(_ context.Context, request *ga2v20250115.DescribeForwardingRuleRequest) (*ga2v20250115.DescribeForwardingRuleResponse, error) {
		assert.Equal(t, "ga2-test001", *request.GlobalAcceleratorId)
		assert.Equal(t, "lis-test001", *request.ListenerId)
		assert.Equal(t, "fp-test001", *request.ForwardingPolicyId)

		resp := &ga2v20250115.DescribeForwardingRuleResponse{}
		resp.Response = &ga2v20250115.DescribeForwardingRuleResponseParams{
			ForwardingRuleSet: []*ga2v20250115.ForwardingRuleSet{
				{
					ForwardingRuleId:    ptrStringGa2FR("fr-test001"),
					GlobalAcceleratorId: ptrStringGa2FR("ga2-test001"),
					ListenerId:          ptrStringGa2FR("lis-test001"),
					ForwardingPolicyId:  ptrStringGa2FR("fp-test001"),
					RuleCondition: []*ga2v20250115.RuleCondition{
						{
							RuleConditionType:  ptrStringGa2FR("HOST"),
							RuleConditionValue: []*string{ptrStringGa2FR("example.com")},
						},
					},
					RuleAction: []*ga2v20250115.RuleAction{
						{
							RuleActionType:  ptrStringGa2FR("ForwardGroup"),
							RuleActionValue: ptrStringGa2FR("eg-test001"),
						},
					},
					OriginHeaders: []*ga2v20250115.OriginHeader{
						{
							Key:   ptrStringGa2FR("X-Custom-Header"),
							Value: ptrStringGa2FR("custom-value"),
						},
					},
					EnableOriginSni: ptrBoolGa2FR(true),
					OriginSni:       ptrStringGa2FR("example.com"),
					OriginHost:      ptrStringGa2FR("origin.example.com"),
				},
			},
			TotalCount: func() *uint64 { v := uint64(1); return &v }(),
			RequestId:  ptrStringGa2FR("req-001"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2ForwardingRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test001",
		"listener_id":           "lis-test001",
		"forwarding_policy_id":  "fp-test001",
		"rule_conditions":       []interface{}{},
		"rule_actions":          []interface{}{},
	})
	d.SetId("ga2-test001" + tccommon.FILED_SP + "lis-test001" + tccommon.FILED_SP + "fp-test001" + tccommon.FILED_SP + "fr-test001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fr-test001", d.Get("forwarding_rule_id"))
	assert.Equal(t, true, d.Get("enable_origin_sni"))
	assert.Equal(t, "example.com", d.Get("origin_sni"))
	assert.Equal(t, "origin.example.com", d.Get("origin_host"))
}

func TestUnitGa2ForwardingRule_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2ForwardingRule()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch DescribeForwardingRuleWithContext - return empty set
	patches.ApplyMethodFunc(ga2Client, "DescribeForwardingRuleWithContext", func(_ context.Context, request *ga2v20250115.DescribeForwardingRuleRequest) (*ga2v20250115.DescribeForwardingRuleResponse, error) {
		resp := &ga2v20250115.DescribeForwardingRuleResponse{}
		resp.Response = &ga2v20250115.DescribeForwardingRuleResponseParams{
			ForwardingRuleSet: []*ga2v20250115.ForwardingRuleSet{},
			TotalCount:        func() *uint64 { v := uint64(0); return &v }(),
			RequestId:         ptrStringGa2FR("req-001"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2ForwardingRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test001",
		"listener_id":           "lis-test001",
		"forwarding_policy_id":  "fp-test001",
		"rule_conditions":       []interface{}{},
		"rule_actions":          []interface{}{},
	})
	d.SetId("ga2-test001" + tccommon.FILED_SP + "lis-test001" + tccommon.FILED_SP + "fp-test001" + tccommon.FILED_SP + "fr-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestUnitGa2ForwardingRule_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2ForwardingRule()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch ModifyForwardingRuleWithContext
	patches.ApplyMethodFunc(ga2Client, "ModifyForwardingRuleWithContext", func(_ context.Context, request *ga2v20250115.ModifyForwardingRuleRequest) (*ga2v20250115.ModifyForwardingRuleResponse, error) {
		assert.Equal(t, "ga2-test001", *request.GlobalAcceleratorId)
		assert.Equal(t, "lis-test001", *request.ListenerId)
		assert.Equal(t, "fp-test001", *request.ForwardingPolicyId)
		assert.Equal(t, "fr-test001", *request.ForwardingRuleId)

		resp := &ga2v20250115.ModifyForwardingRuleResponse{}
		resp.Response = &ga2v20250115.ModifyForwardingRuleResponseParams{
			TaskId:    ptrStringGa2FR("task-002"),
			RequestId: ptrStringGa2FR("req-001"),
		}
		return resp, nil
	})

	// Patch DescribeTaskResultWithContext
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(_ context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := &ga2v20250115.DescribeTaskResultResponse{}
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2FR("SUCCESS"),
			RequestId: ptrStringGa2FR("req-002"),
		}
		return resp, nil
	})

	// Patch DescribeForwardingRuleWithContext for Read after Update
	patches.ApplyMethodFunc(ga2Client, "DescribeForwardingRuleWithContext", func(_ context.Context, request *ga2v20250115.DescribeForwardingRuleRequest) (*ga2v20250115.DescribeForwardingRuleResponse, error) {
		resp := &ga2v20250115.DescribeForwardingRuleResponse{}
		resp.Response = &ga2v20250115.DescribeForwardingRuleResponseParams{
			ForwardingRuleSet: []*ga2v20250115.ForwardingRuleSet{
				{
					ForwardingRuleId:    ptrStringGa2FR("fr-test001"),
					GlobalAcceleratorId: ptrStringGa2FR("ga2-test001"),
					ListenerId:          ptrStringGa2FR("lis-test001"),
					ForwardingPolicyId:  ptrStringGa2FR("fp-test001"),
					RuleCondition: []*ga2v20250115.RuleCondition{
						{
							RuleConditionType:  ptrStringGa2FR("PATH"),
							RuleConditionValue: []*string{ptrStringGa2FR("/api/*")},
						},
					},
					RuleAction: []*ga2v20250115.RuleAction{
						{
							RuleActionType:  ptrStringGa2FR("ForwardGroup"),
							RuleActionValue: ptrStringGa2FR("eg-test002"),
						},
					},
					EnableOriginSni: ptrBoolGa2FR(false),
				},
			},
			TotalCount: func() *uint64 { v := uint64(1); return &v }(),
			RequestId:  ptrStringGa2FR("req-003"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2ForwardingRule()
	d := res.TestResourceData()
	d.SetId("ga2-test001" + tccommon.FILED_SP + "lis-test001" + tccommon.FILED_SP + "fp-test001" + tccommon.FILED_SP + "fr-test001")

	_ = d.Set("global_accelerator_id", "ga2-test001")
	_ = d.Set("listener_id", "lis-test001")
	_ = d.Set("forwarding_policy_id", "fp-test001")
	_ = d.Set("rule_conditions", []interface{}{
		map[string]interface{}{
			"rule_condition_type":  "PATH",
			"rule_condition_value": []interface{}{"/api/*"},
		},
	})
	_ = d.Set("rule_actions", []interface{}{
		map[string]interface{}{
			"rule_action_type":  "ForwardGroup",
			"rule_action_value": "eg-test002",
		},
	})
	_ = d.Set("enable_origin_sni", false)

	// Mark as not new to enable HasChange
	d.MarkNewResource()

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestUnitGa2ForwardingRule_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2ForwardingRule()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch DeleteForwardingRuleWithContext
	patches.ApplyMethodFunc(ga2Client, "DeleteForwardingRuleWithContext", func(_ context.Context, request *ga2v20250115.DeleteForwardingRuleRequest) (*ga2v20250115.DeleteForwardingRuleResponse, error) {
		assert.Equal(t, "ga2-test001", *request.GlobalAcceleratorId)
		assert.Equal(t, "lis-test001", *request.ListenerId)
		assert.Equal(t, "fp-test001", *request.ForwardingPolicyId)
		assert.Equal(t, "fr-test001", *request.ForwardingRuleId)

		resp := &ga2v20250115.DeleteForwardingRuleResponse{}
		resp.Response = &ga2v20250115.DeleteForwardingRuleResponseParams{
			TaskId:    ptrStringGa2FR("task-003"),
			RequestId: ptrStringGa2FR("req-001"),
		}
		return resp, nil
	})

	// Patch DescribeTaskResultWithContext
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(_ context.Context, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		assert.Equal(t, "task-003", *request.TaskId)

		resp := &ga2v20250115.DescribeTaskResultResponse{}
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2FR("SUCCESS"),
			RequestId: ptrStringGa2FR("req-002"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2ForwardingRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test001",
		"listener_id":           "lis-test001",
		"forwarding_policy_id":  "fp-test001",
		"rule_conditions":       []interface{}{},
		"rule_actions":          []interface{}{},
	})
	d.SetId("ga2-test001" + tccommon.FILED_SP + "lis-test001" + tccommon.FILED_SP + "fp-test001" + tccommon.FILED_SP + "fr-test001")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

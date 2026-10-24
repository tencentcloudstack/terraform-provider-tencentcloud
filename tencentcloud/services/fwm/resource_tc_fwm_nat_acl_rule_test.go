package fwm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

type mockMetaFwmNatAclRule struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmNatAclRule) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmNatAclRule{}

func newMockMetaFwmNatAclRule() *mockMetaFwmNatAclRule {
	return &mockMetaFwmNatAclRule{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFwm(s string) *string {
	return &s
}

func ptrInt64Fwm(i int64) *int64 {
	return &i
}

func baseFwmNatAclRuleInput() map[string]interface{} {
	return map[string]interface{}{
		"group_id":       "fwmrg-abc",
		"direction":      0,
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "0.0.0.0/0",
		"target_type":    "ip",
		"protocol":       "TCP",
		"rule_action":    "accept",
		"order_index":    1,
		"scope":          "ALL",
		"port":           "80",
		"description":    "tf example nat acl rule",
	}
}

func newFwmNatAclRuleResp() *fwmv20250611.NatAclRuleResp {
	return &fwmv20250611.NatAclRuleResp{
		RuleId:        ptrStringFwm("fwmrule-xyz"),
		Sequence:      ptrInt64Fwm(1),
		Direction:     ptrInt64Fwm(0),
		SourceContent: ptrStringFwm("1.1.1.1"),
		SourceType:    ptrStringFwm("ip"),
		TargetContent: ptrStringFwm("0.0.0.0/0"),
		TargetType:    ptrStringFwm("ip"),
		Port:          ptrStringFwm("80"),
		Protocol:      ptrStringFwm("TCP"),
		RuleAction:    ptrStringFwm("accept"),
		Description:   ptrStringFwm("tf example nat acl rule"),
		Scope:         ptrStringFwm("ALL"),
		ScopeDesc:     ptrStringFwm("all scope"),
	}
}

func TestFwmNatAclRule_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-abc", *request.GroupId)
		assert.Equal(t, 1, len(request.Rules))
		rule := request.Rules[0]
		assert.Equal(t, "1.1.1.1", *rule.SourceContent)
		assert.Equal(t, "ip", *rule.SourceType)
		assert.Equal(t, int64(1), *rule.OrderIndex)
		assert.Equal(t, int64(0), *rule.Direction)

		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RuleIds:   []*string{ptrStringFwm("fwmrule-xyz")},
			RequestId: ptrStringFwm("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-abc", *request.GroupId)
		assert.NotNil(t, request.Filters)
		assert.Equal(t, 1, len(request.Filters))

		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{newFwmNatAclRuleResp()},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-abc#fwmrule-xyz", d.Id())
	assert.Equal(t, "fwmrule-xyz", d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("order_index"))
}

func TestFwmNatAclRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{newFwmNatAclRuleResp()},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId("fwmrg-abc#fwmrule-xyz")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-abc#fwmrule-xyz", d.Id())
	assert.Equal(t, "fwmrule-xyz", d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("order_index"))
	assert.Equal(t, 1, d.Get("sequence"))
}

func TestFwmNatAclRule_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId("fwmrg-abc#fwmrule-xyz")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmNatAclRule_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyNatAclRuleRequest) (*fwmv20250611.ModifyNatAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-abc", *request.GroupId)
		assert.NotNil(t, request.Rule)
		assert.NotNil(t, request.Rule.RuleId)
		assert.Equal(t, "fwmrule-xyz", *request.Rule.RuleId)
		assert.Equal(t, "ap-guangzhou", *request.Rule.Scope)

		resp := fwmv20250611.NewModifyNatAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyNatAclRuleResponseParams{
			RequestId: ptrStringFwm("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		respRule := newFwmNatAclRuleResp()
		respRule.Scope = ptrStringFwm("ap-guangzhou")
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{respRule},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	input := baseFwmNatAclRuleInput()
	input["scope"] = "ap-guangzhou"
	d := schema.TestResourceDataRaw(t, res.Schema, input)
	d.SetId("fwmrg-abc#fwmrule-xyz")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ap-guangzhou", d.Get("scope"))
}

func TestFwmNatAclRule_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteNatAclRuleRequest) (*fwmv20250611.DeleteNatAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-abc", *request.GroupId)
		assert.NotNil(t, request.RuleIds)
		assert.Equal(t, 1, len(request.RuleIds))
		assert.Equal(t, "fwmrule-xyz", *request.RuleIds[0])

		resp := fwmv20250611.NewDeleteNatAclRuleResponse()
		resp.Response = &fwmv20250611.DeleteNatAclRuleResponseParams{
			RequestId: ptrStringFwm("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId("fwmrg-abc#fwmrule-xyz")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestFwmNatAclRule_DeleteResourceNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteNatAclRuleRequest) (*fwmv20250611.DeleteNatAclRuleResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{Code: "ResourceNotFound"}
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId("fwmrg-abc#fwmrule-xyz")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

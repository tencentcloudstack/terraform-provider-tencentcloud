package fwm_test

import (
	"context"
	"strconv"
	"strings"
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

const (
	testFwmNatGroupId   = "fwmrg-abc"
	testFwmNatRuleId    = "fwmrule-xyz"
	testFwmNatDirection = int64(0)
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

func ptrStringFwmNatAcl(s string) *string {
	return &s
}

func ptrInt64FwmNatAcl(i int64) *int64 {
	return &i
}

// fwmNatAclRuleStateId builds the expected resource id: <groupId>#<ruleId>#<direction>.
func fwmNatAclRuleStateId(groupId, ruleId string, direction int64) string {
	return strings.Join([]string{groupId, ruleId, strconv.FormatInt(direction, 10)}, tccommon.FILED_SP)
}

func baseFwmNatAclRuleInput() map[string]interface{} {
	return map[string]interface{}{
		"group_id":       testFwmNatGroupId,
		"direction":      int(testFwmNatDirection),
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "0.0.0.0/0",
		"target_type":    "ip",
		"protocol":       "TCP",
		"rule_action":    "accept",
		"scope":          "ALL",
		"port":           "80",
		"description":    "tf example nat acl rule",
	}
}

func newFwmNatAclRuleResp() *fwmv20250611.NatAclRuleResp {
	return &fwmv20250611.NatAclRuleResp{
		RuleId:        ptrStringFwmNatAcl(testFwmNatRuleId),
		Sequence:      ptrInt64FwmNatAcl(1),
		Direction:     ptrInt64FwmNatAcl(testFwmNatDirection),
		SourceContent: ptrStringFwmNatAcl("1.1.1.1"),
		SourceType:    ptrStringFwmNatAcl("ip"),
		TargetContent: ptrStringFwmNatAcl("0.0.0.0/0"),
		TargetType:    ptrStringFwmNatAcl("ip"),
		Port:          ptrStringFwmNatAcl("80"),
		Protocol:      ptrStringFwmNatAcl("TCP"),
		RuleAction:    ptrStringFwmNatAcl("accept"),
		Description:   ptrStringFwmNatAcl("tf example nat acl rule"),
		Scope:         ptrStringFwmNatAcl("ALL"),
	}
}

func TestFwmNatAclRule_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		assert.Equal(t, 1, len(request.Rules))
		rule := request.Rules[0]
		assert.Equal(t, "1.1.1.1", *rule.SourceContent)
		assert.Equal(t, "ip", *rule.SourceType)
		assert.Equal(t, int64(1), *rule.OrderIndex)
		assert.Equal(t, testFwmNatDirection, *rule.Direction)

		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RuleIds:   []*string{ptrStringFwmNatAcl(testFwmNatRuleId)},
			RequestId: ptrStringFwmNatAcl("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		assert.NotNil(t, request.Direction)
		assert.Equal(t, testFwmNatDirection, *request.Direction)
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
	// resource id is <groupId>#<ruleId>#<direction>
	assert.Equal(t, fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection), d.Id())
	assert.Equal(t, testFwmNatRuleId, d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("sequence"))
}

func TestFwmNatAclRule_Create_InboundDirection(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		assert.Equal(t, int64(1), *request.Rules[0].Direction)

		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RuleIds:   []*string{ptrStringFwmNatAcl(testFwmNatRuleId)},
			RequestId: ptrStringFwmNatAcl("fake-request-id-inbound"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		assert.Equal(t, int64(1), *request.Direction)

		respRule := newFwmNatAclRuleResp()
		respRule.Direction = ptrInt64FwmNatAcl(1)
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{respRule},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	input := baseFwmNatAclRuleInput()
	input["direction"] = 1
	d := schema.TestResourceDataRaw(t, res.Schema, input)

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, testFwmNatGroupId+"#"+testFwmNatRuleId+"#1", d.Id())
}

func TestFwmNatAclRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		// direction comes from the resource id, not from state.
		assert.Equal(t, int64(1), *request.Direction)
		assert.Equal(t, int64(0), *request.Offset)
		assert.Equal(t, int64(100), *request.Limit)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "RuleId", *request.Filters[0].Name)
		assert.Equal(t, testFwmNatRuleId, *request.Filters[0].Values[0])

		respRule := newFwmNatAclRuleResp()
		respRule.Direction = ptrInt64FwmNatAcl(1)
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{respRule},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(testFwmNatGroupId + "#" + testFwmNatRuleId + "#1")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, testFwmNatGroupId+"#"+testFwmNatRuleId+"#1", d.Id())
	assert.Equal(t, testFwmNatRuleId, d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("direction"))
	assert.Equal(t, 1, d.Get("sequence"))
}

func TestFwmNatAclRule_ReadBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	called := false
	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		called = true
		return fwmv20250611.NewDescribeNatAclRulesResponse(), nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()

	// two segments only.
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(testFwmNatGroupId + tccommon.FILED_SP + testFwmNatRuleId)
	assert.Error(t, res.Read(d, meta))

	// direction is not a number.
	d2 := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d2.SetId(testFwmNatGroupId + tccommon.FILED_SP + testFwmNatRuleId + tccommon.FILED_SP + "inbound")
	assert.Error(t, res.Read(d2, meta))

	assert.False(t, called)
}

// TestFwmNatAclRule_ReadRetryableError verifies the describe call is retried on a retryable SDK error.
func TestFwmNatAclRule_ReadRetryableError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	callCount := 0
	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		callCount++
		if callCount == 1 {
			return nil, &sdkErrors.TencentCloudSDKError{Code: "RequestLimitExceeded", Message: "too many requests"}
		}

		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{newFwmNatAclRuleResp()},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)
	assert.Equal(t, testFwmNatRuleId, d.Get("rule_id"))
}

// TestFwmNatAclRule_ReadNilResponse verifies a nil describe response surfaces an error instead of being treated as deleted.
func TestFwmNatAclRule_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		return &fwmv20250611.DescribeNatAclRulesResponse{}, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Equal(t, fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection), d.Id())
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
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

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
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		assert.NotNil(t, request.Rule)
		assert.NotNil(t, request.Rule.RuleId)
		assert.Equal(t, testFwmNatRuleId, *request.Rule.RuleId)
		assert.Equal(t, "ap-guangzhou", *request.Rule.Scope)
		// direction is taken from the resource id.
		assert.Equal(t, int64(1), *request.Rule.Direction)

		resp := fwmv20250611.NewModifyNatAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyNatAclRuleResponseParams{
			RequestId: ptrStringFwmNatAcl("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		respRule := newFwmNatAclRuleResp()
		respRule.Scope = ptrStringFwmNatAcl("ap-guangzhou")
		respRule.Direction = ptrInt64FwmNatAcl(1)
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
	d.SetId(testFwmNatGroupId + "#" + testFwmNatRuleId + "#1")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ap-guangzhou", d.Get("scope"))
}

func TestFwmNatAclRule_UpdateBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(testFwmNatGroupId + tccommon.FILED_SP + testFwmNatRuleId)

	assert.Error(t, res.Update(d, newMockMetaFwmNatAclRule()))
}

func TestFwmNatAclRule_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteNatAclRuleRequest) (*fwmv20250611.DeleteNatAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		assert.NotNil(t, request.RuleIds)
		assert.Equal(t, 1, len(request.RuleIds))
		assert.Equal(t, testFwmNatRuleId, *request.RuleIds[0])

		resp := fwmv20250611.NewDeleteNatAclRuleResponse()
		resp.Response = &fwmv20250611.DeleteNatAclRuleResponseParams{
			RequestId: ptrStringFwmNatAcl("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestFwmNatAclRule_DeleteBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(testFwmNatGroupId + tccommon.FILED_SP + testFwmNatRuleId)

	assert.Error(t, res.Delete(d, newMockMetaFwmNatAclRule()))
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
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

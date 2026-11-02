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

type mockMetaFwmSgr struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmSgr) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmSgr{}

func newMockMetaFwmSgr() *mockMetaFwmSgr {
	return &mockMetaFwmSgr{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFwmSgr(s string) *string { return &s }
func ptrInt64FwmSgr(v int64) *int64    { return &v }

const (
	testFwmGroupId           = "group-id-example"
	testFwmRuleId            = "rule-id-example"
	testFwmIpVersion         = "ipv4"
	testFwmSourceType        = "net"
	testFwmSourceContent     = "1.1.1.1"
	testFwmDestType          = "net"
	testFwmDestContent       = "2.2.2.2"
	testFwmProtocol          = "TCP"
	testFwmPort              = "80"
	testFwmRuleAction        = "accept"
	testFwmScope             = "SG"
	testFwmDescription       = "tf example"
	testFwmServiceTemplateId = "svc-tpl-example"
	testFwmMemberId          = "member-example"

	// testFwmSequence is the order index returned by the API.
	testFwmSequence = int64(1)
	// testFwmDefaultOrderIndex is the order index hard-coded when creating a rule.
	testFwmDefaultOrderIndex = int64(1)

	// source/destination type and strategy int values returned by the API.
	testFwmSourceTypeInt = int64(0)  // net
	testFwmDestTypeInt   = int64(3)  // cvm
	testFwmStrategyInt   = int64(2)  // accept
	testFwmUnknownEnum   = int64(99) // out of the enum range, kept as decimal string
)

func fwmSgrStateId(groupId, ruleId string) string {
	return groupId + tccommon.FILED_SP + ruleId
}

func buildFwmSgrResponse(ruleId string) *fwmv20250611.DescribeSecurityGroupRulesResponse {
	resp := fwmv20250611.NewDescribeSecurityGroupRulesResponse()
	resp.Response = &fwmv20250611.DescribeSecurityGroupRulesResponseParams{
		Rules: []*fwmv20250611.SecGroupRuleResp{
			{
				RuleId:            ptrStringFwmSgr(ruleId),
				IpVersion:         ptrStringFwmSgr(testFwmIpVersion),
				OrderIndex:        ptrInt64FwmSgr(testFwmSequence),
				SourceId:          ptrStringFwmSgr(testFwmSourceContent),
				SourceType:        ptrInt64FwmSgr(testFwmSourceTypeInt),
				TargetId:          ptrStringFwmSgr(testFwmDestContent),
				TargetType:        ptrInt64FwmSgr(testFwmDestTypeInt),
				Protocol:          ptrStringFwmSgr(testFwmProtocol),
				Port:              ptrStringFwmSgr(testFwmPort),
				Strategy:          ptrInt64FwmSgr(testFwmStrategyInt),
				Detail:            ptrStringFwmSgr(testFwmDescription),
				Scope:             ptrStringFwmSgr(testFwmScope),
				ServiceTemplateId: ptrStringFwmSgr(testFwmServiceTemplateId),
				BelongMember: &fwmv20250611.MemberInfo{
					MemberId: ptrStringFwmSgr(testFwmMemberId),
				},
			},
		},
		TotalCount:    ptrInt64FwmSgr(1),
		AllTotalCount: ptrInt64FwmSgr(1),
		RequestId:     ptrStringFwmSgr("fake-request-id"),
	}
	return resp
}

func baseFwmSgrResourceData() map[string]interface{} {
	return map[string]interface{}{
		"group_id":       testFwmGroupId,
		"ip_version":     testFwmIpVersion,
		"source_type":    testFwmSourceType,
		"source_content": testFwmSourceContent,
		"dest_type":      testFwmDestType,
		"dest_content":   testFwmDestContent,
		"protocol":       testFwmProtocol,
		"port":           testFwmPort,
		"rule_action":    testFwmRuleAction,
		"scope":          testFwmScope,
		"description":    testFwmDescription,
	}
}

// TestFwmSecurityGroupRule_Create verifies the create request payload and the composite id.
func TestFwmSecurityGroupRule_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	var capturedRequest *fwmv20250611.CreateSecurityGroupRuleRequest
	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		capturedRequest = request

		assert.Equal(t, testFwmGroupId, *request.GroupId)
		assert.Len(t, request.Rules, 1)

		rule := request.Rules[0]
		assert.Equal(t, testFwmIpVersion, *rule.IpVersion)
		assert.Equal(t, testFwmSourceType, *rule.SourceType)
		assert.Equal(t, testFwmSourceContent, *rule.SourceContent)
		assert.Equal(t, testFwmDestType, *rule.DestType)
		assert.Equal(t, testFwmDestContent, *rule.DestContent)
		assert.Equal(t, testFwmProtocol, *rule.Protocol)
		assert.Equal(t, testFwmPort, *rule.Port)
		assert.Equal(t, testFwmRuleAction, *rule.RuleAction)
		assert.Equal(t, testFwmScope, *rule.Scope)
		assert.Equal(t, testFwmDescription, *rule.Description)
		// create always uses the default order index.
		assert.Equal(t, testFwmDefaultOrderIndex, *rule.OrderIndex)
		// optional arguments are not set in base data.
		assert.Nil(t, rule.BelongMemberId)
		assert.Nil(t, rule.ProtocolPortType)

		resp := fwmv20250611.NewCreateSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.CreateSecurityGroupRuleResponseParams{
			RuleIds:   []*string{ptrStringFwmSgr(testFwmRuleId)},
			RequestId: ptrStringFwmSgr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		return buildFwmSgrResponse(testFwmRuleId), nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, fwmSgrStateId(testFwmGroupId, testFwmRuleId), d.Id())
	assert.NotNil(t, capturedRequest)
}

// TestFwmSecurityGroupRule_Create_EmptyRuleIds verifies create fails when the API returns no rule id.
func TestFwmSecurityGroupRule_Create_EmptyRuleIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		resp := fwmv20250611.NewCreateSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.CreateSecurityGroupRuleResponseParams{
			RuleIds:   []*string{},
			RequestId: ptrStringFwmSgr("fake-request-id-empty"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmSecurityGroupRule_Read verifies Read maps the API response into state.
func TestFwmSecurityGroupRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		assert.Equal(t, testFwmGroupId, *request.GroupId)
		assert.Equal(t, int64(0), *request.Offset)
		assert.Equal(t, int64(100), *request.Limit)
		assert.Len(t, request.Filters, 1)
		assert.Equal(t, "RuleId", *request.Filters[0].Name)
		assert.Equal(t, testFwmRuleId, *request.Filters[0].Values[0])
		assert.Equal(t, int64(1), *request.Filters[0].OperatorType)
		return buildFwmSgrResponse(testFwmRuleId), nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())
	d.SetId(fwmSgrStateId(testFwmGroupId, testFwmRuleId))

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, fwmSgrStateId(testFwmGroupId, testFwmRuleId), d.Id())

	assert.Equal(t, testFwmGroupId, d.Get("group_id"))
	assert.Equal(t, testFwmRuleId, d.Get("rule_id"))
	assert.Equal(t, testFwmIpVersion, d.Get("ip_version"))
	assert.Equal(t, int(testFwmSequence), d.Get("sequence"))
	assert.Equal(t, testFwmSourceType, d.Get("source_type"))
	assert.Equal(t, testFwmSourceContent, d.Get("source_content"))
	assert.Equal(t, testFwmSourceContent, d.Get("source_id"))
	// TargetType 3 is mapped to cvm.
	assert.Equal(t, "cvm", d.Get("dest_type"))
	assert.Equal(t, testFwmDestContent, d.Get("dest_content"))
	assert.Equal(t, testFwmDestContent, d.Get("target_id"))
	assert.Equal(t, testFwmProtocol, d.Get("protocol"))
	assert.Equal(t, testFwmPort, d.Get("port"))
	// Strategy 2 is mapped to accept.
	assert.Equal(t, testFwmRuleAction, d.Get("rule_action"))
	assert.Equal(t, testFwmDescription, d.Get("description"))
	assert.Equal(t, testFwmScope, d.Get("scope"))
	assert.Equal(t, testFwmServiceTemplateId, d.Get("service_template_id"))
	assert.Equal(t, testFwmMemberId, d.Get("belong_member_id"))
}

// TestFwmSecurityGroupRule_Read_UnknownEnumCode verifies unknown enum values are kept as decimal strings.
func TestFwmSecurityGroupRule_Read_UnknownEnumCode(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		resp := fwmv20250611.NewDescribeSecurityGroupRulesResponse()
		resp.Response = &fwmv20250611.DescribeSecurityGroupRulesResponseParams{
			Rules: []*fwmv20250611.SecGroupRuleResp{
				{
					RuleId:     ptrStringFwmSgr(testFwmRuleId),
					SourceType: ptrInt64FwmSgr(testFwmUnknownEnum),
					TargetType: ptrInt64FwmSgr(testFwmUnknownEnum),
					Strategy:   ptrInt64FwmSgr(testFwmUnknownEnum),
				},
			},
			RequestId: ptrStringFwmSgr("fake-request-id-unknown"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())
	d.SetId(fwmSgrStateId(testFwmGroupId, testFwmRuleId))

	err := res.Read(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, "99", d.Get("source_type"))
	assert.Equal(t, "99", d.Get("dest_type"))
	assert.Equal(t, "99", d.Get("rule_action"))
}

// TestFwmSecurityGroupRule_Read_NotFound verifies Read clears state when the describe returns empty.
func TestFwmSecurityGroupRule_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		resp := fwmv20250611.NewDescribeSecurityGroupRulesResponse()
		resp.Response = &fwmv20250611.DescribeSecurityGroupRulesResponseParams{
			Rules:      []*fwmv20250611.SecGroupRuleResp{},
			TotalCount: ptrInt64FwmSgr(0),
			RequestId:  ptrStringFwmSgr("fake-request-id-empty"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())
	d.SetId(fwmSgrStateId("group-missing", "rule-missing"))

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmSecurityGroupRule_Read_BrokenId verifies Read rejects a malformed id before calling the API.
func TestFwmSecurityGroupRule_Read_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	called := false
	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		called = true
		return buildFwmSgrResponse(testFwmRuleId), nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())
	d.SetId("broken-id")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.False(t, called)
}

// TestFwmSecurityGroupRule_Update verifies Update triggers ModifySecurityGroupRule with the new port.
func TestFwmSecurityGroupRule_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	var capturedModifyRequest *fwmv20250611.ModifySecurityGroupRuleRequest
	patches.ApplyMethodFunc(fwmClient, "ModifySecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifySecurityGroupRuleRequest) (*fwmv20250611.ModifySecurityGroupRuleResponse, error) {
		capturedModifyRequest = request

		assert.Equal(t, testFwmGroupId, *request.GroupId)
		assert.NotNil(t, request.Rule)
		assert.Equal(t, testFwmRuleId, *request.Rule.RuleId)
		assert.Equal(t, "443", *request.Rule.Port)
		assert.Equal(t, testFwmProtocol, *request.Rule.Protocol)
		// sequence is absent from state, so the default order index is used.
		assert.Equal(t, testFwmDefaultOrderIndex, *request.Rule.OrderIndex)

		resp := fwmv20250611.NewModifySecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.ModifySecurityGroupRuleResponseParams{
			RequestId: ptrStringFwmSgr("fake-request-id-modify"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		return buildFwmSgrResponse(testFwmRuleId), nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	rd := baseFwmSgrResourceData()
	rd["port"] = "443"
	d := schema.TestResourceDataRaw(t, res.Schema, rd)
	d.SetId(fwmSgrStateId(testFwmGroupId, testFwmRuleId))

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "port"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
}

// TestFwmSecurityGroupRule_Update_WithSequence verifies the sequence in state is used as the order index.
func TestFwmSecurityGroupRule_Update_WithSequence(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	var capturedModifyRequest *fwmv20250611.ModifySecurityGroupRuleRequest
	patches.ApplyMethodFunc(fwmClient, "ModifySecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifySecurityGroupRuleRequest) (*fwmv20250611.ModifySecurityGroupRuleResponse, error) {
		capturedModifyRequest = request

		resp := fwmv20250611.NewModifySecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.ModifySecurityGroupRuleResponseParams{
			RequestId: ptrStringFwmSgr("fake-request-id-modify-sequence"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		return buildFwmSgrResponse(testFwmRuleId), nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	rd := baseFwmSgrResourceData()
	rd["sequence"] = 7
	d := schema.TestResourceDataRaw(t, res.Schema, rd)
	d.SetId(fwmSgrStateId(testFwmGroupId, testFwmRuleId))

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "rule_action"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
	assert.Equal(t, int64(7), *capturedModifyRequest.Rule.OrderIndex)
}

// TestFwmSecurityGroupRule_Update_NoChange verifies no modify call is made when nothing changed.
func TestFwmSecurityGroupRule_Update_NoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	modified := false
	patches.ApplyMethodFunc(fwmClient, "ModifySecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifySecurityGroupRuleRequest) (*fwmv20250611.ModifySecurityGroupRuleResponse, error) {
		modified = true
		return fwmv20250611.NewModifySecurityGroupRuleResponse(), nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		return buildFwmSgrResponse(testFwmRuleId), nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())
	d.SetId(fwmSgrStateId(testFwmGroupId, testFwmRuleId))

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, modified)
}

// TestFwmSecurityGroupRule_Delete verifies Delete issues DeleteSecurityGroupRule with the id.
func TestFwmSecurityGroupRule_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	var capturedRequest *fwmv20250611.DeleteSecurityGroupRuleRequest
	patches.ApplyMethodFunc(fwmClient, "DeleteSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteSecurityGroupRuleRequest) (*fwmv20250611.DeleteSecurityGroupRuleResponse, error) {
		capturedRequest = request
		assert.Equal(t, testFwmGroupId, *request.GroupId)
		assert.Len(t, request.RuleIds, 1)
		assert.Equal(t, testFwmRuleId, *request.RuleIds[0])

		resp := fwmv20250611.NewDeleteSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.DeleteSecurityGroupRuleResponseParams{
			RequestId: ptrStringFwmSgr("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())
	d.SetId(fwmSgrStateId(testFwmGroupId, testFwmRuleId))

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
}

// TestFwmSecurityGroupRule_Delete_ResourceNotFound verifies a missing rule is treated as deleted.
func TestFwmSecurityGroupRule_Delete_ResourceNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteSecurityGroupRuleRequest) (*fwmv20250611.DeleteSecurityGroupRuleResponse, error) {
		return nil, sdkErrors.NewTencentCloudSDKError("ResourceNotFound", "rule not found", "fake-request-id-not-found")
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())
	d.SetId(fwmSgrStateId(testFwmGroupId, testFwmRuleId))

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFwmSecurityGroupRule_Delete_BrokenId verifies Delete rejects a malformed id.
func TestFwmSecurityGroupRule_Delete_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())
	d.SetId("broken-id")

	err := res.Delete(d, newMockMetaFwmSgr())
	assert.Error(t, err)
}

// TestFwmSecurityGroupRule_Schema validates the schema definition.
func TestFwmSecurityGroupRule_Schema(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	requiredArgs := []string{
		"group_id",
		"ip_version",
		"source_type",
		"source_content",
		"dest_type",
		"dest_content",
		"protocol",
		"port",
		"rule_action",
	}
	for _, arg := range requiredArgs {
		assert.Contains(t, res.Schema, arg)
		assert.True(t, res.Schema[arg].Required, arg)
	}

	forceNewArgs := []string{"group_id", "ip_version", "source_type", "dest_type", "scope"}
	for _, arg := range forceNewArgs {
		assert.True(t, res.Schema[arg].ForceNew, arg)
	}

	optionalArgs := []string{"scope", "service_template_id", "description", "protocol_port_type", "belong_member_id"}
	for _, arg := range optionalArgs {
		assert.Contains(t, res.Schema, arg)
		assert.True(t, res.Schema[arg].Optional, arg)
		assert.False(t, res.Schema[arg].Required, arg)
	}

	computedArgs := []string{"rule_id", "source_id", "target_id", "sequence"}
	for _, arg := range computedArgs {
		assert.Contains(t, res.Schema, arg)
		assert.True(t, res.Schema[arg].Computed, arg)
		assert.False(t, res.Schema[arg].Optional, arg)
	}

	assert.Equal(t, schema.TypeString, res.Schema["group_id"].Type)
	assert.Equal(t, schema.TypeString, res.Schema["ip_version"].Type)
	assert.Equal(t, schema.TypeInt, res.Schema["sequence"].Type)
	assert.Equal(t, schema.TypeInt, res.Schema["protocol_port_type"].Type)
	assert.Equal(t, schema.TypeString, res.Schema["rule_id"].Type)

	// service_template_id is optional and computed at the same time.
	assert.True(t, res.Schema["service_template_id"].Optional)
	assert.True(t, res.Schema["service_template_id"].Computed)

	// fields removed from the resource must not come back silently.
	assert.NotContains(t, res.Schema, "order_index")
	assert.NotContains(t, res.Schema, "dns_parse_count")
	assert.NotContains(t, res.Schema, "strategy")
	assert.NotContains(t, res.Schema, "detail")
}

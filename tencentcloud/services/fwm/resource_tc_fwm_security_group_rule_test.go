package fwm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
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
	testFwmGroupId       = "group-id-example"
	testFwmRuleId        = "rule-id-example"
	testFwmIpVersion     = "ipv4"
	testFwmOrderIndex    = int64(1)
	testFwmSourceType    = "net"
	testFwmSourceContent = "1.1.1.1"
	testFwmDestType      = "net"
	testFwmDestContent   = "2.2.2.2"
	testFwmProtocol      = "TCP"
	testFwmPort          = "80"
	testFwmRuleAction    = "accept"
	testFwmScope         = "SG"
	testFwmDescription   = "tf example"
)

func buildFwmSgrResponse(ruleId string) *fwmv20250611.DescribeSecurityGroupRulesResponse {
	resp := fwmv20250611.NewDescribeSecurityGroupRulesResponse()
	resp.Response = &fwmv20250611.DescribeSecurityGroupRulesResponseParams{
		Rules: []*fwmv20250611.SecGroupRuleResp{
			{
				RuleId:     ptrStringFwmSgr(ruleId),
				IpVersion:  ptrStringFwmSgr(testFwmIpVersion),
				OrderIndex: ptrInt64FwmSgr(testFwmOrderIndex),
				SourceId:   ptrStringFwmSgr(testFwmSourceContent),
				SourceType: ptrInt64FwmSgr(0),
				TargetId:   ptrStringFwmSgr(testFwmDestContent),
				TargetType: ptrInt64FwmSgr(0),
				Protocol:   ptrStringFwmSgr(testFwmProtocol),
				Port:       ptrStringFwmSgr(testFwmPort),
				Strategy:   ptrInt64FwmSgr(2),
				Detail:     ptrStringFwmSgr(testFwmDescription),
				Scope:      ptrStringFwmSgr(testFwmScope),
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
		"order_index":    int(testFwmOrderIndex),
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

// TestFwmSecurityGroupRule_Create verifies the Create flow sets the composite id from the mock response.
func TestFwmSecurityGroupRule_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	var capturedRequest *fwmv20250611.CreateSecurityGroupRuleRequest
	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		capturedRequest = request
		assert.Equal(t, testFwmGroupId, *request.GroupId)
		assert.Equal(t, 1, len(request.Rules))
		rule := request.Rules[0]
		assert.Equal(t, testFwmIpVersion, *rule.IpVersion)
		assert.Equal(t, testFwmOrderIndex, *rule.OrderIndex)
		assert.Equal(t, testFwmSourceType, *rule.SourceType)
		assert.Equal(t, testFwmSourceContent, *rule.SourceContent)
		assert.Equal(t, testFwmDestType, *rule.DestType)
		assert.Equal(t, testFwmDestContent, *rule.DestContent)
		assert.Equal(t, testFwmProtocol, *rule.Protocol)
		assert.Equal(t, testFwmPort, *rule.Port)
		assert.Equal(t, testFwmRuleAction, *rule.RuleAction)
		assert.Equal(t, testFwmScope, *rule.Scope)

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
	assert.Equal(t, testFwmGroupId+"#"+testFwmRuleId, d.Id())
	assert.NotNil(t, capturedRequest)
}

// TestFwmSecurityGroupRule_Read verifies Read populates computed fields from the response.
func TestFwmSecurityGroupRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmSgr().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		assert.Equal(t, testFwmGroupId, *request.GroupId)
		assert.NotNil(t, request.Filters)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "RuleId", *request.Filters[0].Name)
		assert.Equal(t, testFwmRuleId, *request.Filters[0].Values[0])
		return buildFwmSgrResponse(testFwmRuleId), nil
	})

	meta := newMockMetaFwmSgr()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmSgrResourceData())
	d.SetId(testFwmGroupId + "#" + testFwmRuleId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, testFwmGroupId+"#"+testFwmRuleId, d.Id())

	assert.Equal(t, testFwmRuleId, d.Get("rule_id"))
	assert.Equal(t, testFwmIpVersion, d.Get("ip_version"))
	assert.Equal(t, int(testFwmOrderIndex), d.Get("order_index"))
	assert.Equal(t, testFwmSourceType, d.Get("source_type"))
	assert.Equal(t, testFwmSourceContent, d.Get("source_content"))
	assert.Equal(t, testFwmDestType, d.Get("dest_type"))
	assert.Equal(t, testFwmDestContent, d.Get("dest_content"))
	assert.Equal(t, testFwmProtocol, d.Get("protocol"))
	assert.Equal(t, testFwmPort, d.Get("port"))
	assert.Equal(t, testFwmRuleAction, d.Get("rule_action"))
	assert.Equal(t, testFwmDescription, d.Get("description"))
	assert.Equal(t, testFwmScope, d.Get("scope"))
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
	d.SetId("group-missing#rule-missing")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
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
	d.SetId(testFwmGroupId + "#" + testFwmRuleId)

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "port"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
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
		assert.NotNil(t, request.RuleIds)
		assert.Equal(t, 1, len(request.RuleIds))
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
	d.SetId(testFwmGroupId + "#" + testFwmRuleId)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
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

	assert.Contains(t, res.Schema, "group_id")
	assert.Contains(t, res.Schema, "ip_version")
	assert.Contains(t, res.Schema, "order_index")
	assert.Contains(t, res.Schema, "source_type")
	assert.Contains(t, res.Schema, "source_content")
	assert.Contains(t, res.Schema, "dest_type")
	assert.Contains(t, res.Schema, "dest_content")
	assert.Contains(t, res.Schema, "protocol")
	assert.Contains(t, res.Schema, "port")
	assert.Contains(t, res.Schema, "rule_action")
	assert.Contains(t, res.Schema, "scope")
	assert.Contains(t, res.Schema, "service_template_id")
	assert.Contains(t, res.Schema, "description")
	assert.Contains(t, res.Schema, "protocol_port_type")
	assert.Contains(t, res.Schema, "belong_member_id")
	assert.Contains(t, res.Schema, "rule_id")

	groupId := res.Schema["group_id"]
	assert.Equal(t, schema.TypeString, groupId.Type)
	assert.True(t, groupId.Required)

	ipVersion := res.Schema["ip_version"]
	assert.Equal(t, schema.TypeString, ipVersion.Type)
	assert.True(t, ipVersion.Required)
	assert.True(t, ipVersion.ForceNew)

	sourceType := res.Schema["source_type"]
	assert.Equal(t, schema.TypeString, sourceType.Type)
	assert.True(t, sourceType.Required)
	assert.True(t, sourceType.ForceNew)

	ruleId := res.Schema["rule_id"]
	assert.Equal(t, schema.TypeString, ruleId.Type)
	assert.True(t, ruleId.Computed)

	dnsParseCount := res.Schema["dns_parse_count"]
	assert.Equal(t, schema.TypeList, dnsParseCount.Type)
	assert.True(t, dnsParseCount.Computed)

	assert.NotContains(t, res.Schema, "rule_set")
	assert.NotContains(t, res.Schema, "rule_list")
}

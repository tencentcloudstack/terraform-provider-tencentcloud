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

type mockMetaFwm struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwm) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwm{}

func newMockMetaFwm() *mockMetaFwm {
	return &mockMetaFwm{client: &connectivity.TencentCloudClient{}}
}

func fwmStringPtr(s string) *string {
	return &s
}

func fwmInt64Ptr(i int64) *int64 {
	return &i
}

func fwmDescribeRuleResp(ruleId string) *fwmv20250611.DescribeSecurityGroupRulesResponse {
	resp := fwmv20250611.NewDescribeSecurityGroupRulesResponse()
	resp.Response = &fwmv20250611.DescribeSecurityGroupRulesResponseParams{
		RequestId:     fwmStringPtr("fake-request-id"),
		TotalCount:    fwmInt64Ptr(1),
		AllTotalCount: fwmInt64Ptr(1),
		Rules: []*fwmv20250611.SecGroupRuleResp{
			{
				OrderIndex:        fwmInt64Ptr(1),
				RuleId:            fwmStringPtr(ruleId),
				IpVersion:         fwmStringPtr("4"),
				SourceId:          fwmStringPtr("1.1.1.1/0"),
				SourceType:        fwmInt64Ptr(0),
				TargetId:          fwmStringPtr("0.0.0.0/0"),
				TargetType:        fwmInt64Ptr(0),
				Protocol:          fwmStringPtr("TCP"),
				Port:              fwmStringPtr("-1/-1"),
				Strategy:          fwmInt64Ptr(1),
				Detail:            fwmStringPtr("rule description."),
				Region:            fwmStringPtr("ap-guangzhou"),
				ServiceTemplateId: fwmStringPtr("ppmg-xxxxxxxx"),
				SouInstanceName:   fwmStringPtr("sou-instance"),
				SouPublicIp:       fwmStringPtr("2.2.2.2"),
				SouPrivateIp:      fwmStringPtr("10.0.0.1"),
				SouCidr:           fwmStringPtr("1.1.1.1/0"),
				SouParameterName:  fwmStringPtr("sou-template"),
				InstanceName:      fwmStringPtr("dest-instance"),
				PublicIp:          fwmStringPtr("3.3.3.3"),
				PrivateIp:         fwmStringPtr("10.0.0.2"),
				Cidr:              fwmStringPtr("0.0.0.0/0"),
				ParameterName:     fwmStringPtr("dest-template"),
				ProtocolPortName:  fwmStringPtr("port-template"),
				Id:                fwmInt64Ptr(1234567890),
				Scope:             fwmStringPtr("all"),
				IsNew:             fwmInt64Ptr(0),
				DnsParseCount: &fwmv20250611.SgDnsParseCount{
					ValidCount:   fwmInt64Ptr(2),
					InvalidCount: fwmInt64Ptr(1),
				},
				BelongMember: &fwmv20250611.MemberInfo{
					AppId:    fwmStringPtr("1300000000"),
					Uin:      fwmStringPtr("100000000001"),
					Nickname: fwmStringPtr("tf-example"),
					MemberId: fwmStringPtr("member-xxxxxxxx"),
					NodeName: fwmStringPtr("department"),
				},
			},
		},
	}
	return resp
}

func TestFwmSecurityGroupRule_Create_ByGroupName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "tf-example-fwm-group", *request.GroupName)
		assert.NotNil(t, request.Product)
		assert.Equal(t, "enterprise_sg", *request.Product)

		resp := fwmv20250611.NewCreateRuleGroupResponse()
		resp.Response = &fwmv20250611.CreateRuleGroupResponseParams{
			GroupId:   fwmStringPtr("fwmrg-xxxxxxxx"),
			RequestId: fwmStringPtr("fake-request-id-group"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-xxxxxxxx", *request.GroupId)
		assert.Len(t, request.Rules, 1)
		assert.NotNil(t, request.Rules[0].Protocol)
		assert.Equal(t, "TCP", *request.Rules[0].Protocol)
		assert.NotNil(t, request.Rules[0].Description)
		assert.Equal(t, "rule description.", *request.Rules[0].Description)

		resp := fwmv20250611.NewCreateSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.CreateSecurityGroupRuleResponseParams{
			RuleIds:   []*string{fwmStringPtr("fwmrule-yyyyyyyy")},
			RequestId: fwmStringPtr("fake-request-id-rule"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		return fwmDescribeRuleResp("fwmrule-yyyyyyyy"), nil
	})

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name":     "tf-example-fwm-group",
		"ip_version":     "4",
		"protocol":       "TCP",
		"port":           "-1/-1",
		"rule_action":    "accept",
		"description":    "rule description.",
		"source_type":    "net",
		"dest_type":      "net",
		"source_content": "1.1.1.1/0",
		"dest_content":   "0.0.0.0/0",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-xxxxxxxx#fwmrule-yyyyyyyy", d.Id())
	assert.Equal(t, "fwmrule-yyyyyyyy", d.Get("rule_id"))
	assert.Equal(t, "fwmrg-xxxxxxxx", d.Get("group_id"))
	assert.Equal(t, "TCP", d.Get("protocol"))
}

func TestFwmSecurityGroupRule_Create_ByGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
		t.Fatalf("CreateRuleGroup should not be called when group_id is provided")
		return nil, nil
	})

	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-existed", *request.GroupId)
		assert.Len(t, request.Rules, 1)

		resp := fwmv20250611.NewCreateSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.CreateSecurityGroupRuleResponseParams{
			RuleIds:   []*string{fwmStringPtr("fwmrule-zzzzzzzz")},
			RequestId: fwmStringPtr("fake-request-id-rule"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		return fwmDescribeRuleResp("fwmrule-zzzzzzzz"), nil
	})

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":   "fwmrg-existed",
		"ip_version": "4",
		"protocol":   "TCP",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-existed#fwmrule-zzzzzzzz", d.Id())
	assert.Equal(t, "fwmrule-zzzzzzzz", d.Get("rule_id"))
}

func TestFwmSecurityGroupRule_Create_NoIdentity(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"ip_version": "4",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "either group_id or group_name must be specified")
	assert.Equal(t, "", d.Id())
}

func TestFwmSecurityGroupRule_Create_BothIdentity(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":   "fwmrg-existed",
		"group_name": "tf-example-fwm-group",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mutually exclusive")
	assert.Equal(t, "", d.Id())
}

func TestFwmSecurityGroupRule_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-xxxxxxxx", *request.GroupId)
		assert.Len(t, request.Filters, 1)
		assert.Equal(t, "RuleId", *request.Filters[0].Name)
		assert.Equal(t, "fwmrule-yyyyyyyy", *request.Filters[0].Values[0])
		return fwmDescribeRuleResp("fwmrule-yyyyyyyy"), nil
	})

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("fwmrg-xxxxxxxx#fwmrule-yyyyyyyy")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-xxxxxxxx#fwmrule-yyyyyyyy", d.Id())
	assert.Equal(t, "fwmrg-xxxxxxxx", d.Get("group_id"))
	assert.Equal(t, "fwmrule-yyyyyyyy", d.Get("rule_id"))
	assert.Equal(t, "4", d.Get("ip_version"))
	assert.Equal(t, "TCP", d.Get("protocol"))
	assert.Equal(t, "-1/-1", d.Get("port"))
	assert.Equal(t, "ppmg-xxxxxxxx", d.Get("service_template_id"))
	assert.Equal(t, "all", d.Get("scope"))
	assert.Equal(t, 1, d.Get("order_index"))
	assert.Equal(t, "1.1.1.1/0", d.Get("source_id"))
	assert.Equal(t, "0.0.0.0/0", d.Get("target_id"))
	assert.Equal(t, 0, d.Get("target_type"))
	assert.Equal(t, 1, d.Get("strategy"))
	assert.Equal(t, "rule description.", d.Get("detail"))
	assert.Equal(t, "ap-guangzhou", d.Get("region"))
	assert.Equal(t, "sou-instance", d.Get("sou_instance_name"))
	assert.Equal(t, "2.2.2.2", d.Get("sou_public_ip"))
	assert.Equal(t, "10.0.0.1", d.Get("sou_private_ip"))
	assert.Equal(t, "1.1.1.1/0", d.Get("sou_cidr"))
	assert.Equal(t, "sou-template", d.Get("sou_parameter_name"))
	assert.Equal(t, "dest-instance", d.Get("instance_name"))
	assert.Equal(t, "3.3.3.3", d.Get("public_ip"))
	assert.Equal(t, "10.0.0.2", d.Get("private_ip"))
	assert.Equal(t, "0.0.0.0/0", d.Get("cidr"))
	assert.Equal(t, "dest-template", d.Get("parameter_name"))
	assert.Equal(t, "port-template", d.Get("protocol_port_name"))
	assert.Equal(t, 0, d.Get("is_new"))
	assert.Equal(t, 1234567890, d.Get("rule_uuid"))

	dnsParseCount := d.Get("dns_parse_count").([]interface{})
	assert.Len(t, dnsParseCount, 1)
	assert.Equal(t, 2, dnsParseCount[0].(map[string]interface{})["valid_count"])
	assert.Equal(t, 1, dnsParseCount[0].(map[string]interface{})["invalid_count"])

	belongMember := d.Get("belong_member").([]interface{})
	assert.Len(t, belongMember, 1)
	assert.Equal(t, "1300000000", belongMember[0].(map[string]interface{})["app_id"])
	assert.Equal(t, "100000000001", belongMember[0].(map[string]interface{})["uin"])
	assert.Equal(t, "tf-example", belongMember[0].(map[string]interface{})["nickname"])
	assert.Equal(t, "member-xxxxxxxx", belongMember[0].(map[string]interface{})["member_id"])
	assert.Equal(t, "department", belongMember[0].(map[string]interface{})["node_name"])
}

func TestFwmSecurityGroupRule_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		resp := fwmv20250611.NewDescribeSecurityGroupRulesResponse()
		resp.Response = &fwmv20250611.DescribeSecurityGroupRulesResponseParams{
			RequestId:  fwmStringPtr("fake-request-id"),
			TotalCount: fwmInt64Ptr(0),
			Rules:      []*fwmv20250611.SecGroupRuleResp{},
		}
		return resp, nil
	})

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("fwmrg-xxxxxxxx#fwmrule-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmSecurityGroupRule_Read_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("fwmrule-yyyyyyyy")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

func TestFwmSecurityGroupRule_Update_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifySecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifySecurityGroupRuleRequest) (*fwmv20250611.ModifySecurityGroupRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-xxxxxxxx", *request.GroupId)
		assert.NotNil(t, request.Rule)
		assert.NotNil(t, request.Rule.RuleId)
		assert.Equal(t, "fwmrule-yyyyyyyy", *request.Rule.RuleId)

		resp := fwmv20250611.NewModifySecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.ModifySecurityGroupRuleResponseParams{
			RequestId: fwmStringPtr("fake-request-id-modify"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRulesRequest) (*fwmv20250611.DescribeSecurityGroupRulesResponse, error) {
		return fwmDescribeRuleResp("fwmrule-yyyyyyyy"), nil
	})

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":    "fwmrg-xxxxxxxx",
		"ip_version":  "4",
		"protocol":    "TCP",
		"description": "rule description changed.",
	})
	d.SetId("fwmrg-xxxxxxxx#fwmrule-yyyyyyyy")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-xxxxxxxx#fwmrule-yyyyyyyy", d.Id())
}

func TestFwmSecurityGroupRule_Delete_WithGroupName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	deleteRuleCalled := false
	patches.ApplyMethodFunc(fwmClient, "DeleteSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteSecurityGroupRuleRequest) (*fwmv20250611.DeleteSecurityGroupRuleResponse, error) {
		deleteRuleCalled = true
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-xxxxxxxx", *request.GroupId)
		assert.Len(t, request.RuleIds, 1)
		assert.Equal(t, "fwmrule-yyyyyyyy", *request.RuleIds[0])

		resp := fwmv20250611.NewDeleteSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.DeleteSecurityGroupRuleResponseParams{
			RequestId: fwmStringPtr("fake-request-id-delete-rule"),
		}
		return resp, nil
	})

	deleteGroupCalled := false
	patches.ApplyMethodFunc(fwmClient, "DeleteRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteRuleGroupRequest) (*fwmv20250611.DeleteRuleGroupResponse, error) {
		deleteGroupCalled = true
		assert.Len(t, request.GroupIds, 1)
		assert.Equal(t, "fwmrg-xxxxxxxx", *request.GroupIds[0])

		resp := fwmv20250611.NewDeleteRuleGroupResponse()
		resp.Response = &fwmv20250611.DeleteRuleGroupResponseParams{
			RequestId: fwmStringPtr("fake-request-id-delete-group"),
		}
		return resp, nil
	})

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example-fwm-group",
	})
	d.SetId("fwmrg-xxxxxxxx#fwmrule-yyyyyyyy")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteRuleCalled)
	assert.True(t, deleteGroupCalled)
}

func TestFwmSecurityGroupRule_Delete_WithGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwm().client, "UseFwmV20250611Client", fwmClient)

	deleteRuleCalled := false
	patches.ApplyMethodFunc(fwmClient, "DeleteSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteSecurityGroupRuleRequest) (*fwmv20250611.DeleteSecurityGroupRuleResponse, error) {
		deleteRuleCalled = true
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-existed", *request.GroupId)
		assert.Len(t, request.RuleIds, 1)
		assert.Equal(t, "fwmrule-yyyyyyyy", *request.RuleIds[0])

		resp := fwmv20250611.NewDeleteSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.DeleteSecurityGroupRuleResponseParams{
			RequestId: fwmStringPtr("fake-request-id-delete-rule"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DeleteRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteRuleGroupRequest) (*fwmv20250611.DeleteRuleGroupResponse, error) {
		t.Fatalf("DeleteRuleGroup should not be called when group_id is provided")
		return nil, nil
	})

	meta := newMockMetaFwm()
	res := fwm.ResourceTencentCloudFwmSecurityGroupRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id": "fwmrg-existed",
	})
	d.SetId("fwmrg-existed#fwmrule-yyyyyyyy")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteRuleCalled)
}

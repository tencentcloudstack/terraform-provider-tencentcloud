package fwm_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcfwm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

// mockMetaFwmVpcAclRule implements tccommon.ProviderMeta
type mockMetaFwmVpcAclRule struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmVpcAclRule) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmVpcAclRule{}

func ptrFwmVpcAclRuleString(s string) *string {
	return &s
}

func ptrFwmVpcAclRuleInt64(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/fwm/ -run "TestFwmVpcAclRule" -v -count=1 -gcflags="all=-l"

// TestFwmVpcAclRule_Create tests the Create function
func TestFwmVpcAclRule_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateVpcAclRuleWithContext", func(_ interface{}, request *fwmv20250611.CreateVpcAclRuleRequest) (*fwmv20250611.CreateVpcAclRuleResponse, error) {
		resp := fwmv20250611.NewCreateVpcAclRuleResponse()
		resp.Response = &fwmv20250611.CreateVpcAclRuleResponseParams{
			RuleIds:   []*string{ptrFwmVpcAclRuleString("rule-test001")},
			RequestId: ptrFwmVpcAclRuleString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeVpcAclRulesWithContext", func(_ interface{}, request *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
		rule := &fwmv20250611.VpcAclRuleResp{
			RuleId:        ptrFwmVpcAclRuleString("rule-test001"),
			Sequence:      ptrFwmVpcAclRuleInt64(1),
			IpVersion:     ptrFwmVpcAclRuleString("ipv4"),
			SourceType:    ptrFwmVpcAclRuleString("ip"),
			SourceContent: ptrFwmVpcAclRuleString("10.0.0.0/24"),
			DestType:      ptrFwmVpcAclRuleString("ip"),
			DestContent:   ptrFwmVpcAclRuleString("192.168.1.0/24"),
			Protocol:      ptrFwmVpcAclRuleString("TCP"),
			Port:          ptrFwmVpcAclRuleString("80"),
			RuleAction:    ptrFwmVpcAclRuleString("accept"),
			Description:   ptrFwmVpcAclRuleString("test rule"),
			EdgeId:        ptrFwmVpcAclRuleString("ALL"),
			FwGroupId:     ptrFwmVpcAclRuleString("ALL"),
		}
		resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{
			Rules:         []*fwmv20250611.VpcAclRuleResp{rule},
			AllTotalCount: ptrFwmVpcAclRuleInt64(1),
			RequestId:     ptrFwmVpcAclRuleString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaFwmVpcAclRule{client: mockClient}
	res := svcfwm.ResourceTencentCloudFwmVpcAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":       "group-test001",
		"edge_id":        "ALL",
		"fw_group_id":    "ALL",
		"ip_version":     "ipv4",
		"source_type":    "ip",
		"source_content": "10.0.0.0/24",
		"dest_type":      "ip",
		"dest_content":   "192.168.1.0/24",
		"protocol":       "TCP",
		"port":           "80",
		"rule_action":    "accept",
		"order_index":    1,
		"description":    "test rule",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "group-test001"+tccommon.FILED_SP+"rule-test001", d.Id())
}

// TestFwmVpcAclRule_Read tests the Read function
func TestFwmVpcAclRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeVpcAclRulesWithContext", func(_ interface{}, request *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
		rule := &fwmv20250611.VpcAclRuleResp{
			RuleId:        ptrFwmVpcAclRuleString("rule-test001"),
			Sequence:      ptrFwmVpcAclRuleInt64(1),
			IpVersion:     ptrFwmVpcAclRuleString("ipv4"),
			SourceType:    ptrFwmVpcAclRuleString("ip"),
			SourceContent: ptrFwmVpcAclRuleString("10.0.0.0/24"),
			SourceName:    ptrFwmVpcAclRuleString("src-name"),
			DestType:      ptrFwmVpcAclRuleString("ip"),
			DestContent:   ptrFwmVpcAclRuleString("192.168.1.0/24"),
			DestName:      ptrFwmVpcAclRuleString("dst-name"),
			Protocol:      ptrFwmVpcAclRuleString("TCP"),
			Port:          ptrFwmVpcAclRuleString("80"),
			RuleAction:    ptrFwmVpcAclRuleString("accept"),
			Description:   ptrFwmVpcAclRuleString("test rule"),
			EdgeId:        ptrFwmVpcAclRuleString("ALL"),
			FwGroupId:     ptrFwmVpcAclRuleString("ALL"),
			FwGroupName:   ptrFwmVpcAclRuleString("fw-group-name"),
			Invalid:       ptrFwmVpcAclRuleInt64(0),
			BelongMember: &fwmv20250611.MemberInfo{
				AppId:    ptrFwmVpcAclRuleString("app-id-test"),
				Uin:      ptrFwmVpcAclRuleString("uin-test"),
				Nickname: ptrFwmVpcAclRuleString("nickname-test"),
				MemberId: ptrFwmVpcAclRuleString("member-id-test"),
				NodeName: ptrFwmVpcAclRuleString("node-name-test"),
			},
			CreateTime:  ptrFwmVpcAclRuleString("2024-01-01 00:00:00"),
			UpdateTime:  ptrFwmVpcAclRuleString("2024-01-02 00:00:00"),
			DnsParseCnt: ptrFwmVpcAclRuleInt64(2),
			HitCount:    ptrFwmVpcAclRuleInt64(3),
			LastHitTime: ptrFwmVpcAclRuleString("-"),
		}
		resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{
			Rules:         []*fwmv20250611.VpcAclRuleResp{rule},
			AllTotalCount: ptrFwmVpcAclRuleInt64(1),
			RequestId:     ptrFwmVpcAclRuleString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaFwmVpcAclRule{client: mockClient}
	res := svcfwm.ResourceTencentCloudFwmVpcAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("group-test001" + tccommon.FILED_SP + "rule-test001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ALL", d.Get("edge_id").(string))
	assert.Equal(t, "ALL", d.Get("fw_group_id").(string))
	assert.Equal(t, "ipv4", d.Get("ip_version").(string))
	assert.Equal(t, "ip", d.Get("source_type").(string))
	assert.Equal(t, "10.0.0.0/24", d.Get("source_content").(string))
	assert.Equal(t, "ip", d.Get("dest_type").(string))
	assert.Equal(t, "192.168.1.0/24", d.Get("dest_content").(string))
	assert.Equal(t, "TCP", d.Get("protocol").(string))
	assert.Equal(t, "80", d.Get("port").(string))
	assert.Equal(t, "accept", d.Get("rule_action").(string))
	assert.Equal(t, int64(1), int64(d.Get("order_index").(int)))
	assert.Equal(t, "test rule", d.Get("description").(string))
	assert.Equal(t, "rule-test001", d.Get("rule_id").(string))
	assert.Equal(t, int64(1), int64(d.Get("sequence").(int)))
	assert.Equal(t, "src-name", d.Get("source_name").(string))
	assert.Equal(t, "dst-name", d.Get("dest_name").(string))
	assert.Equal(t, "fw-group-name", d.Get("fw_group_name").(string))
	assert.Equal(t, int64(0), int64(d.Get("invalid").(int)))
	assert.Equal(t, "app-id-test", d.Get("app_id").(string))
	assert.Equal(t, "uin-test", d.Get("uin").(string))
	assert.Equal(t, "nickname-test", d.Get("nickname").(string))
	assert.Equal(t, "member-id-test", d.Get("member_id").(string))
	assert.Equal(t, "member-id-test", d.Get("belong_member_id").(string))
	assert.Equal(t, "node-name-test", d.Get("node_name").(string))
	assert.Equal(t, int64(1), int64(d.Get("all_total_count").(int)))
	assert.Equal(t, int64(2), int64(d.Get("dns_parse_cnt").(int)))
	assert.Equal(t, int64(3), int64(d.Get("hit_count").(int)))
	assert.Equal(t, "-", d.Get("last_hit_time").(string))
}

// TestFwmVpcAclRule_Update tests the Update function
func TestFwmVpcAclRule_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyVpcAclRuleWithContext", func(_ interface{}, request *fwmv20250611.ModifyVpcAclRuleRequest) (*fwmv20250611.ModifyVpcAclRuleResponse, error) {
		resp := fwmv20250611.NewModifyVpcAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyVpcAclRuleResponseParams{
			RequestId: ptrFwmVpcAclRuleString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeVpcAclRulesWithContext", func(_ interface{}, request *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
		rule := &fwmv20250611.VpcAclRuleResp{
			RuleId:        ptrFwmVpcAclRuleString("rule-test001"),
			Sequence:      ptrFwmVpcAclRuleInt64(2),
			IpVersion:     ptrFwmVpcAclRuleString("ipv4"),
			SourceType:    ptrFwmVpcAclRuleString("ip"),
			SourceContent: ptrFwmVpcAclRuleString("10.0.0.0/24"),
			DestType:      ptrFwmVpcAclRuleString("ip"),
			DestContent:   ptrFwmVpcAclRuleString("192.168.1.0/24"),
			Protocol:      ptrFwmVpcAclRuleString("UDP"),
			Port:          ptrFwmVpcAclRuleString("443"),
			RuleAction:    ptrFwmVpcAclRuleString("drop"),
			Description:   ptrFwmVpcAclRuleString("updated rule"),
			EdgeId:        ptrFwmVpcAclRuleString("ALL"),
			FwGroupId:     ptrFwmVpcAclRuleString("ALL"),
		}
		resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{
			Rules:         []*fwmv20250611.VpcAclRuleResp{rule},
			AllTotalCount: ptrFwmVpcAclRuleInt64(1),
			RequestId:     ptrFwmVpcAclRuleString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaFwmVpcAclRule{client: mockClient}
	res := svcfwm.ResourceTencentCloudFwmVpcAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":       "group-test001",
		"edge_id":        "ALL",
		"fw_group_id":    "ALL",
		"ip_version":     "ipv4",
		"source_type":    "ip",
		"source_content": "10.0.0.0/24",
		"dest_type":      "ip",
		"dest_content":   "192.168.1.0/24",
		"protocol":       "UDP",
		"port":           "443",
		"rule_action":    "drop",
		"order_index":    2,
		"description":    "updated rule",
	})
	d.SetId("group-test001" + tccommon.FILED_SP + "rule-test001")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestFwmVpcAclRule_Delete tests the Delete function
func TestFwmVpcAclRule_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteVpcAclRuleWithContext", func(_ interface{}, request *fwmv20250611.DeleteVpcAclRuleRequest) (*fwmv20250611.DeleteVpcAclRuleResponse, error) {
		resp := fwmv20250611.NewDeleteVpcAclRuleResponse()
		resp.Response = &fwmv20250611.DeleteVpcAclRuleResponseParams{
			RequestId: ptrFwmVpcAclRuleString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaFwmVpcAclRule{client: mockClient}
	res := svcfwm.ResourceTencentCloudFwmVpcAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("group-test001" + tccommon.FILED_SP + "rule-test001")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFwmVpcAclRule_Schema tests the schema definition
func TestFwmVpcAclRule_Schema(t *testing.T) {
	res := svcfwm.ResourceTencentCloudFwmVpcAclRule()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "group_id")
	assert.Contains(t, res.Schema, "edge_id")
	assert.Contains(t, res.Schema, "fw_group_id")
	assert.Contains(t, res.Schema, "ip_version")
	assert.Contains(t, res.Schema, "source_type")
	assert.Contains(t, res.Schema, "source_content")
	assert.Contains(t, res.Schema, "dest_type")
	assert.Contains(t, res.Schema, "dest_content")
	assert.Contains(t, res.Schema, "protocol")
	assert.Contains(t, res.Schema, "rule_action")
	assert.Contains(t, res.Schema, "order_index")
	assert.Contains(t, res.Schema, "description")
	assert.Contains(t, res.Schema, "port")
	assert.Contains(t, res.Schema, "param_template_id")
	assert.Contains(t, res.Schema, "belong_member_id")
	assert.Contains(t, res.Schema, "rule_id")
	assert.Contains(t, res.Schema, "all_total_count")

	forceNewFields := []string{"group_id", "edge_id", "fw_group_id", "ip_version", "source_type", "source_content", "dest_type", "dest_content"}
	for _, fieldName := range forceNewFields {
		field := res.Schema[fieldName]
		assert.True(t, field.ForceNew, "field %s should be ForceNew", fieldName)
		assert.True(t, field.Required, "field %s should be Required", fieldName)
	}

	requiredFields := []string{"protocol", "rule_action", "order_index", "description"}
	for _, fieldName := range requiredFields {
		field := res.Schema[fieldName]
		assert.True(t, field.Required, "field %s should be Required", fieldName)
		assert.False(t, field.ForceNew, "field %s should not be ForceNew", fieldName)
	}

	optionalFields := []string{"port", "param_template_id", "belong_member_id"}
	for _, fieldName := range optionalFields {
		field := res.Schema[fieldName]
		assert.True(t, field.Optional, "field %s should be Optional", fieldName)
	}

	computedFields := []string{"rule_id", "sequence", "source_name", "dest_name", "param_template_name", "fw_group_name", "invalid", "app_id", "uin", "nickname", "member_id", "node_name", "create_time", "update_time", "dns_parse_cnt", "hit_count", "last_hit_time", "all_total_count"}
	for _, fieldName := range computedFields {
		field := res.Schema[fieldName]
		assert.True(t, field.Computed, "field %s should be Computed", fieldName)
	}
}

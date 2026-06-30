package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaForL7AccRuleV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForL7AccRuleV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForL7AccRuleV2{}

func newMockMetaForL7AccRuleV2() *mockMetaForL7AccRuleV2 {
	return &mockMetaForL7AccRuleV2{client: &connectivity.TencentCloudClient{}}
}

func ptrStrL7V2(s string) *string {
	return &s
}

func ptrInt64L7V2(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestL7AccRuleV2" -v -count=1 -gcflags="all=-l"

// TestL7AccRuleV2_Create_Success tests Create calls API and sets composite ID
func TestL7AccRuleV2_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, 1, len(request.Rules))
		assert.Equal(t, "enable", *request.Rules[0].Status)
		assert.Equal(t, "Test Rule", *request.Rules[0].RuleName)
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrStrL7V2("rule-abc123")},
			RequestId: ptrStrL7V2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7V2(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrStrL7V2("rule-abc123"),
					Status:       ptrStrL7V2("enable"),
					RuleName:     ptrStrL7V2("Test Rule"),
					Description:  []*string{ptrStrL7V2("test desc")},
					RulePriority: ptrInt64L7V2(1),
					Branches: []*teov20220901.RuleBranch{
						{
							Condition: ptrStrL7V2("${http.request.host} in ['example.com']"),
							Actions:   []*teov20220901.RuleEngineAction{},
						},
					},
				},
			},
			RequestId: ptrStrL7V2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"status":      "enable",
		"rule_name":   "Test Rule",
		"description": []interface{}{"test desc"},
		"branches":    []interface{}{},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#rule-abc123", d.Id())
}

// TestL7AccRuleV2_Create_EmptyResponse tests Create returns error on empty response
func TestL7AccRuleV2_Create_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{},
			RequestId: ptrStrL7V2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"status":    "enable",
		"rule_name": "Test Rule",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty rule id")
}

// TestL7AccRuleV2_Create_APIError tests Create handles API error
func TestL7AccRuleV2_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-invalid",
		"status":    "enable",
		"rule_name": "Test Rule",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestL7AccRuleV2_Read_Success tests Read populates state from API
func TestL7AccRuleV2_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "rule-id", *request.Filters[0].Name)
		assert.Equal(t, "rule-abc123", *request.Filters[0].Values[0])
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7V2(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrStrL7V2("rule-abc123"),
					Status:       ptrStrL7V2("enable"),
					RuleName:     ptrStrL7V2("Test Rule"),
					Description:  []*string{ptrStrL7V2("test desc")},
					RulePriority: ptrInt64L7V2(5),
					Branches: []*teov20220901.RuleBranch{
						{
							Condition: ptrStrL7V2("${http.request.host} in ['example.com']"),
							Actions:   []*teov20220901.RuleEngineAction{},
						},
					},
				},
			},
			RequestId: ptrStrL7V2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"status":    "enable",
		"rule_name": "Test Rule",
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "rule-abc123", d.Get("rule_id"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "Test Rule", d.Get("rule_name"))
	assert.Equal(t, 5, d.Get("rule_priority"))
}

// TestL7AccRuleV2_Read_NotFound tests Read handles resource not found
func TestL7AccRuleV2_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7V2(0),
			Rules:      []*teov20220901.RuleEngineItem{},
			RequestId:  ptrStrL7V2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"status":    "enable",
		"rule_name": "Test Rule",
	})
	d.SetId("zone-test123#rule-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestL7AccRuleV2_Update_Success tests Update calls ModifyL7AccRule API
func TestL7AccRuleV2_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule", func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "rule-abc123", *request.Rule.RuleId)
		assert.Equal(t, "Updated Rule", *request.Rule.RuleName)
		assert.Equal(t, "disable", *request.Rule.Status)
		resp := teov20220901.NewModifyL7AccRuleResponse()
		resp.Response = &teov20220901.ModifyL7AccRuleResponseParams{
			RequestId: ptrStrL7V2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7V2(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrStrL7V2("rule-abc123"),
					Status:       ptrStrL7V2("disable"),
					RuleName:     ptrStrL7V2("Updated Rule"),
					Description:  []*string{ptrStrL7V2("updated desc")},
					RulePriority: ptrInt64L7V2(1),
					Branches: []*teov20220901.RuleBranch{
						{
							Condition: ptrStrL7V2("${http.request.host} in ['example.com']"),
							Actions:   []*teov20220901.RuleEngineAction{},
						},
					},
				},
			},
			RequestId: ptrStrL7V2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"status":      "disable",
		"rule_name":   "Updated Rule",
		"description": []interface{}{"updated desc"},
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestL7AccRuleV2_Update_APIError tests Update handles API error
func TestL7AccRuleV2_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule", func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Rule not found")
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"status":    "disable",
		"rule_name": "Updated Rule",
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestL7AccRuleV2_Delete_Success tests Delete calls DeleteL7AccRules API
func TestL7AccRuleV2_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteL7AccRules", func(request *teov20220901.DeleteL7AccRulesRequest) (*teov20220901.DeleteL7AccRulesResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, 1, len(request.RuleIds))
		assert.Equal(t, "rule-abc123", *request.RuleIds[0])
		resp := teov20220901.NewDeleteL7AccRulesResponse()
		resp.Response = &teov20220901.DeleteL7AccRulesResponseParams{
			RequestId: ptrStrL7V2("fake-request-id"),
		}
		return resp, nil
	})

	// After delete, Read is called which should return empty
	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7V2(0),
			Rules:      []*teov20220901.RuleEngineItem{},
			RequestId:  ptrStrL7V2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"status":    "enable",
		"rule_name": "Test Rule",
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestL7AccRuleV2_Delete_APIError tests Delete handles API error
func TestL7AccRuleV2_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteL7AccRules", func(request *teov20220901.DeleteL7AccRulesRequest) (*teov20220901.DeleteL7AccRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Rule not found")
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"status":    "enable",
		"rule_name": "Test Rule",
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestL7AccRuleV2_Schema validates schema definition
func TestL7AccRuleV2_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "rule_name")
	assert.Contains(t, res.Schema, "description")
	assert.Contains(t, res.Schema, "branches")
	assert.Contains(t, res.Schema, "rule_id")
	assert.Contains(t, res.Schema, "rule_priority")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Optional)

	ruleName := res.Schema["rule_name"]
	assert.Equal(t, schema.TypeString, ruleName.Type)
	assert.True(t, ruleName.Optional)

	description := res.Schema["description"]
	assert.Equal(t, schema.TypeList, description.Type)
	assert.True(t, description.Optional)

	branches := res.Schema["branches"]
	assert.Equal(t, schema.TypeList, branches.Type)
	assert.True(t, branches.Optional)

	ruleId := res.Schema["rule_id"]
	assert.Equal(t, schema.TypeString, ruleId.Type)
	assert.True(t, ruleId.Computed)

	rulePriority := res.Schema["rule_priority"]
	assert.Equal(t, schema.TypeInt, rulePriority.Type)
	assert.True(t, rulePriority.Computed)
}

// TestL7AccRuleV2_Read_APIError tests Read handles API error
func TestL7AccRuleV2_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=Internal error")
	})

	meta := newMockMetaForL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"status":    "enable",
		"rule_name": "Test Rule",
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}

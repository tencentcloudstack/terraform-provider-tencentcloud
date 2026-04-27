package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestSecurityJsInjectionRule" -v -count=1 -gcflags="all=-l"

// TestSecurityJsInjectionRule_Create_Success tests Create calls API and sets ID
func TestSecurityJsInjectionRule_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.CreateSecurityJSInjectionRuleRequest) (*teov20220901.CreateSecurityJSInjectionRuleResponse, error) {
		resp := teov20220901.NewCreateSecurityJSInjectionRuleResponse()
		resp.Response = &teov20220901.CreateSecurityJSInjectionRuleResponseParams{
			JSInjectionRuleIds: []*string{ptrString("rule-001"), ptrString("rule-002")},
			RequestId:          ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityJSInjectionRuleRequest) (*teov20220901.DescribeSecurityJSInjectionRuleResponse, error) {
		resp := teov20220901.NewDescribeSecurityJSInjectionRuleResponse()
		resp.Response = &teov20220901.DescribeSecurityJSInjectionRuleResponseParams{
			TotalCount: ptrInt64(2),
			JSInjectionRules: []*teov20220901.JSInjectionRule{
				{
					RuleId:    ptrString("rule-001"),
					Name:      ptrString("test-rule-1"),
					Priority:  ptrInt64(0),
					Condition: ptrString("${http.request.host} in ['example.com']"),
					InjectJS:  ptrString("inject-sdk-only"),
				},
				{
					RuleId:    ptrString("rule-002"),
					Name:      ptrString("test-rule-2"),
					Priority:  ptrInt64(10),
					Condition: ptrString("${http.request.uri.path} in ['/api/*']"),
					InjectJS:  ptrString("no-injection"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityJsInjectionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"js_injection_rules": []interface{}{
			map[string]interface{}{
				"name":      "test-rule-1",
				"condition": "${http.request.host} in ['example.com']",
			},
			map[string]interface{}{
				"name":      "test-rule-2",
				"priority":  10,
				"condition": "${http.request.uri.path} in ['/api/*']",
				"inject_js": "no-injection",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())
}

// TestSecurityJsInjectionRule_Create_APIError tests Create handles API error
func TestSecurityJsInjectionRule_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.CreateSecurityJSInjectionRuleRequest) (*teov20220901.CreateSecurityJSInjectionRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityJsInjectionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"js_injection_rules": []interface{}{
			map[string]interface{}{
				"name":      "test-rule",
				"condition": "${http.request.host} in ['example.com']",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestSecurityJsInjectionRule_Read_Success tests Read retrieves rule data
func TestSecurityJsInjectionRule_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityJSInjectionRuleRequest) (*teov20220901.DescribeSecurityJSInjectionRuleResponse, error) {
		resp := teov20220901.NewDescribeSecurityJSInjectionRuleResponse()
		resp.Response = &teov20220901.DescribeSecurityJSInjectionRuleResponseParams{
			TotalCount: ptrInt64(1),
			JSInjectionRules: []*teov20220901.JSInjectionRule{
				{
					RuleId:    ptrString("rule-001"),
					Name:      ptrString("test-rule"),
					Priority:  ptrInt64(5),
					Condition: ptrString("${http.request.host} in ['example.com']"),
					InjectJS:  ptrString("inject-sdk-only"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityJsInjectionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())

	rules := d.Get("js_injection_rules").([]interface{})
	assert.Equal(t, 1, len(rules))
	ruleMap := rules[0].(map[string]interface{})
	assert.Equal(t, "rule-001", ruleMap["rule_id"])
	assert.Equal(t, "test-rule", ruleMap["name"])
	assert.Equal(t, 5, ruleMap["priority"])
	assert.Equal(t, "${http.request.host} in ['example.com']", ruleMap["condition"])
	assert.Equal(t, "inject-sdk-only", ruleMap["inject_js"])

	ruleIds := d.Get("js_injection_rule_ids").([]interface{})
	assert.Equal(t, 1, len(ruleIds))
	assert.Equal(t, "rule-001", ruleIds[0])
}

// TestSecurityJsInjectionRule_Read_NotFound tests Read handles resource not found
func TestSecurityJsInjectionRule_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityJSInjectionRuleRequest) (*teov20220901.DescribeSecurityJSInjectionRuleResponse, error) {
		resp := teov20220901.NewDescribeSecurityJSInjectionRuleResponse()
		resp.Response = &teov20220901.DescribeSecurityJSInjectionRuleResponseParams{
			TotalCount:       ptrInt64(0),
			JSInjectionRules: []*teov20220901.JSInjectionRule{},
			RequestId:        ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityJsInjectionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestSecurityJsInjectionRule_Update_Success tests Update calls Modify API
func TestSecurityJsInjectionRule_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifySecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.ModifySecurityJSInjectionRuleRequest) (*teov20220901.ModifySecurityJSInjectionRuleResponse, error) {
		resp := teov20220901.NewModifySecurityJSInjectionRuleResponse()
		resp.Response = &teov20220901.ModifySecurityJSInjectionRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityJSInjectionRuleRequest) (*teov20220901.DescribeSecurityJSInjectionRuleResponse, error) {
		resp := teov20220901.NewDescribeSecurityJSInjectionRuleResponse()
		resp.Response = &teov20220901.DescribeSecurityJSInjectionRuleResponseParams{
			TotalCount: ptrInt64(1),
			JSInjectionRules: []*teov20220901.JSInjectionRule{
				{
					RuleId:    ptrString("rule-001"),
					Name:      ptrString("updated-rule"),
					Priority:  ptrInt64(20),
					Condition: ptrString("${http.request.host} in ['updated.com']"),
					InjectJS:  ptrString("inject-sdk-only"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityJsInjectionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"js_injection_rules": []interface{}{
			map[string]interface{}{
				"rule_id":   "rule-001",
				"name":      "updated-rule",
				"priority":  20,
				"condition": "${http.request.host} in ['updated.com']",
				"inject_js": "inject-sdk-only",
			},
		},
		"js_injection_rule_ids": []interface{}{"rule-001"},
	})
	d.SetId("zone-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestSecurityJsInjectionRule_Update_APIError tests Update handles API error
func TestSecurityJsInjectionRule_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifySecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.ModifySecurityJSInjectionRuleRequest) (*teov20220901.ModifySecurityJSInjectionRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid rule")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityJsInjectionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"js_injection_rules": []interface{}{
			map[string]interface{}{
				"rule_id":   "rule-001",
				"name":      "test-rule",
				"condition": "${http.request.host} in ['example.com']",
			},
		},
		"js_injection_rule_ids": []interface{}{"rule-001"},
	})
	d.SetId("zone-12345678")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestSecurityJsInjectionRule_Delete_Success tests Delete removes rules
func TestSecurityJsInjectionRule_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.DeleteSecurityJSInjectionRuleRequest) (*teov20220901.DeleteSecurityJSInjectionRuleResponse, error) {
		resp := teov20220901.NewDeleteSecurityJSInjectionRuleResponse()
		resp.Response = &teov20220901.DeleteSecurityJSInjectionRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityJsInjectionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":               "zone-12345678",
		"js_injection_rule_ids": []interface{}{"rule-001", "rule-002"},
	})
	d.SetId("zone-12345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestSecurityJsInjectionRule_Delete_APIError tests Delete handles API error
func TestSecurityJsInjectionRule_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSecurityJSInjectionRuleWithContext", func(ctx context.Context, request *teov20220901.DeleteSecurityJSInjectionRuleRequest) (*teov20220901.DeleteSecurityJSInjectionRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Rule not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityJsInjectionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":               "zone-12345678",
		"js_injection_rule_ids": []interface{}{"rule-001"},
	})
	d.SetId("zone-12345678")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestSecurityJsInjectionRule_Schema validates schema definition
func TestSecurityJsInjectionRule_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSecurityJsInjectionRule()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields with ForceNew
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// Check required fields without ForceNew
	assert.Contains(t, res.Schema, "js_injection_rules")
	rules := res.Schema["js_injection_rules"]
	assert.Equal(t, schema.TypeList, rules.Type)
	assert.True(t, rules.Required)
	assert.False(t, rules.ForceNew)

	// Check computed fields
	assert.Contains(t, res.Schema, "js_injection_rule_ids")
	ruleIds := res.Schema["js_injection_rule_ids"]
	assert.Equal(t, schema.TypeList, ruleIds.Type)
	assert.True(t, ruleIds.Computed)

	// Check nested schema
	elem := rules.Elem.(*schema.Resource)
	assert.Contains(t, elem.Schema, "rule_id")
	assert.True(t, elem.Schema["rule_id"].Computed)

	assert.Contains(t, elem.Schema, "name")
	assert.True(t, elem.Schema["name"].Required)

	assert.Contains(t, elem.Schema, "priority")
	assert.True(t, elem.Schema["priority"].Optional)
	assert.True(t, elem.Schema["priority"].Computed)

	assert.Contains(t, elem.Schema, "condition")
	assert.True(t, elem.Schema["condition"].Required)

	assert.Contains(t, elem.Schema, "inject_js")
	assert.True(t, elem.Schema["inject_js"].Optional)
	assert.True(t, elem.Schema["inject_js"].Computed)
}

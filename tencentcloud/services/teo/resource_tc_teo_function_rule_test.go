package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func TestAccTencentCloudTeoFunctionRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "function_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "remark", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.target", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.values.0", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.target", "extension"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.values.0", ".txt"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.operator", "notequal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.target", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.values.0", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.target", "extension"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.values.0", ".png"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function_rule.teo_function_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoFunctionRuleUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "function_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "remark", "bbb"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.operator", "notequal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.target", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.values.0", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.target", "extension"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.values.0", ".txt"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.operator", "notequal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.target", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.values.0", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.target", "extension"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.values.0", ".png"),
				),
			},
		},
	})
}

const testAccTeoFunctionRule = `

resource "tencentcloud_teo_function_rule" "teo_function_rule" {
    function_id   = "ef-txx7fnua"
    remark        = "aaa"
    zone_id       = "zone-2qtuhspy7cr6"

    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".txt",
            ]
        }
    }
    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "notequal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".png",
            ]
        }
    }
}
`

const testAccTeoFunctionRuleUp = `

resource "tencentcloud_teo_function_rule" "teo_function_rule" {
    function_id   = "ef-txx7fnua"
    remark        = "bbb"
    zone_id       = "zone-2qtuhspy7cr6"

    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "notequal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".txt",
            ]
        }
    }
    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "notequal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".png",
            ]
        }
    }
}
`

// ---- Unit Tests (gomonkey mock) ----

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionRule_" -v -count=1 -gcflags="all=-l"

// mockMetaFuncRule implements tccommon.ProviderMeta
type mockMetaFuncRule struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFuncRule) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFuncRule{}

func newMockMetaFuncRule() *mockMetaFuncRule {
	return &mockMetaFuncRule{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFuncRule(s string) *string {
	return &s
}

func ptrInt64FuncRule(i int64) *int64 {
	return &i
}

// TestTeoFunctionRule_Schema validates trigger_type schema definition
func TestTeoFunctionRule_Schema(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoFunctionRule()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "trigger_type")

	triggerTypeField := res.Schema["trigger_type"]
	assert.Equal(t, schema.TypeString, triggerTypeField.Type)
	assert.True(t, triggerTypeField.Optional)
	assert.False(t, triggerTypeField.Required)
	assert.False(t, triggerTypeField.Computed)
	assert.Nil(t, triggerTypeField.Default)
	assert.NotNil(t, triggerTypeField.ValidateFunc)
}

// TestTeoFunctionRule_CreateWithTriggerType tests creating function rule with trigger_type set
func TestTeoFunctionRule_CreateWithTriggerType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFuncRule().client, "UseTeoV20220901Client", teoClient)

	// Mock CreateFunctionRuleWithContext to verify TriggerType is set in request
	patches.ApplyMethodFunc(teoClient, "CreateFunctionRuleWithContext",
		func(_ context.Context, request *teov20220901.CreateFunctionRuleRequest) (*teov20220901.CreateFunctionRuleResponse, error) {
			// Verify TriggerType is set correctly
			assert.NotNil(t, request.TriggerType)
			assert.Equal(t, "direct", *request.TriggerType)

			resp := teov20220901.NewCreateFunctionRuleResponse()
			resp.Response = &teov20220901.CreateFunctionRuleResponseParams{
				RuleId:    ptrStringFuncRule("rule-12345678"),
				RequestId: ptrStringFuncRule("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Mock DescribeTeoFunctionRuleById for the Read call after Create
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoFunctionRuleById",
		func(_ context.Context, zoneId string, functionId string, ruleId string) (*teov20220901.FunctionRule, error) {
			return &teov20220901.FunctionRule{
				RuleId:       ptrStringFuncRule("rule-12345678"),
				FunctionId:   ptrStringFuncRule("ef-testfunc1"),
				FunctionName: ptrStringFuncRule("test-func"),
				TriggerType:  ptrStringFuncRule("direct"),
				Priority:     ptrInt64FuncRule(1),
				Remark:       ptrStringFuncRule("test remark"),
				FunctionRuleConditions: []*teov20220901.FunctionRuleCondition{
					{
						RuleConditions: []*teov20220901.RuleCondition{
							{
								Operator: ptrStringFuncRule("equal"),
								Target:   ptrStringFuncRule("host"),
								Values:   []*string{ptrStringFuncRule("example.com")},
							},
						},
					},
				},
			}, nil
		},
	)

	meta := newMockMetaFuncRule()
	res := svcteo.ResourceTencentCloudTeoFunctionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"function_id":  "ef-testfunc1",
		"remark":       "test remark",
		"trigger_type": "direct",
		"function_rule_conditions": []interface{}{
			map[string]interface{}{
				"rule_conditions": []interface{}{
					map[string]interface{}{
						"operator": "equal",
						"target":   "host",
						"values":   []interface{}{"example.com"},
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify trigger_type is set in state
	triggerType := d.Get("trigger_type").(string)
	assert.Equal(t, "direct", triggerType)

	// Verify composite ID
	assert.Equal(t, "zone-test1234#ef-testfunc1#rule-12345678", d.Id())
}

// TestTeoFunctionRule_CreateWithoutTriggerType tests creating function rule without trigger_type
func TestTeoFunctionRule_CreateWithoutTriggerType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFuncRule().client, "UseTeoV20220901Client", teoClient)

	// Mock CreateFunctionRuleWithContext
	patches.ApplyMethodFunc(teoClient, "CreateFunctionRuleWithContext",
		func(_ context.Context, request *teov20220901.CreateFunctionRuleRequest) (*teov20220901.CreateFunctionRuleResponse, error) {
			// When trigger_type is not set, TriggerType should be nil in request
			assert.Nil(t, request.TriggerType)

			resp := teov20220901.NewCreateFunctionRuleResponse()
			resp.Response = &teov20220901.CreateFunctionRuleResponseParams{
				RuleId:    ptrStringFuncRule("rule-87654321"),
				RequestId: ptrStringFuncRule("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Mock DescribeTeoFunctionRuleById for the Read call
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoFunctionRuleById",
		func(_ context.Context, zoneId string, functionId string, ruleId string) (*teov20220901.FunctionRule, error) {
			return &teov20220901.FunctionRule{
				RuleId:       ptrStringFuncRule("rule-87654321"),
				FunctionId:   ptrStringFuncRule("ef-testfunc2"),
				FunctionName: ptrStringFuncRule("test-func-2"),
				TriggerType:  ptrStringFuncRule("direct"),
				Priority:     ptrInt64FuncRule(1),
				FunctionRuleConditions: []*teov20220901.FunctionRuleCondition{
					{
						RuleConditions: []*teov20220901.RuleCondition{
							{
								Operator: ptrStringFuncRule("equal"),
								Target:   ptrStringFuncRule("host"),
								Values:   []*string{ptrStringFuncRule("example.com")},
							},
						},
					},
				},
			}, nil
		},
	)

	meta := newMockMetaFuncRule()
	res := svcteo.ResourceTencentCloudTeoFunctionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test1234",
		"function_id": "ef-testfunc2",
		"function_rule_conditions": []interface{}{
			map[string]interface{}{
				"rule_conditions": []interface{}{
					map[string]interface{}{
						"operator": "equal",
						"target":   "host",
						"values":   []interface{}{"example.com"},
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test1234#ef-testfunc2#rule-87654321", d.Id())
}

// TestTeoFunctionRule_ReadWithTriggerType tests reading trigger_type from API response
func TestTeoFunctionRule_ReadWithTriggerType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeTeoFunctionRuleById to return a rule with TriggerType
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoFunctionRuleById",
		func(_ context.Context, zoneId string, functionId string, ruleId string) (*teov20220901.FunctionRule, error) {
			assert.Equal(t, "zone-test1234", zoneId)
			assert.Equal(t, "ef-testfunc1", functionId)
			assert.Equal(t, "rule-12345678", ruleId)
			return &teov20220901.FunctionRule{
				RuleId:       ptrStringFuncRule("rule-12345678"),
				FunctionId:   ptrStringFuncRule("ef-testfunc1"),
				FunctionName: ptrStringFuncRule("test-func"),
				TriggerType:  ptrStringFuncRule("weight"),
				Priority:     ptrInt64FuncRule(5),
				Remark:       ptrStringFuncRule("test remark"),
				FunctionRuleConditions: []*teov20220901.FunctionRuleCondition{
					{
						RuleConditions: []*teov20220901.RuleCondition{
							{
								Operator: ptrStringFuncRule("equal"),
								Target:   ptrStringFuncRule("host"),
								Values:   []*string{ptrStringFuncRule("example.com")},
							},
						},
					},
				},
			}, nil
		},
	)

	meta := newMockMetaFuncRule()
	res := svcteo.ResourceTencentCloudTeoFunctionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test1234",
		"function_id": "ef-testfunc1",
		"function_rule_conditions": []interface{}{
			map[string]interface{}{
				"rule_conditions": []interface{}{
					map[string]interface{}{
						"operator": "equal",
						"target":   "host",
						"values":   []interface{}{"example.com"},
					},
				},
			},
		},
	})
	d.SetId("zone-test1234#ef-testfunc1#rule-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify trigger_type is read correctly
	triggerType := d.Get("trigger_type").(string)
	assert.Equal(t, "weight", triggerType)
}

// TestTeoFunctionRule_ReadTriggerTypeNil tests reading when API returns nil TriggerType
func TestTeoFunctionRule_ReadTriggerTypeNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeTeoFunctionRuleById to return a rule without TriggerType
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoFunctionRuleById",
		func(_ context.Context, zoneId string, functionId string, ruleId string) (*teov20220901.FunctionRule, error) {
			return &teov20220901.FunctionRule{
				RuleId:       ptrStringFuncRule("rule-12345678"),
				FunctionId:   ptrStringFuncRule("ef-testfunc1"),
				FunctionName: ptrStringFuncRule("test-func"),
				TriggerType:  nil,
				Priority:     ptrInt64FuncRule(5),
				FunctionRuleConditions: []*teov20220901.FunctionRuleCondition{
					{
						RuleConditions: []*teov20220901.RuleCondition{
							{
								Operator: ptrStringFuncRule("equal"),
								Target:   ptrStringFuncRule("host"),
								Values:   []*string{ptrStringFuncRule("example.com")},
							},
						},
					},
				},
			}, nil
		},
	)

	meta := newMockMetaFuncRule()
	res := svcteo.ResourceTencentCloudTeoFunctionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test1234",
		"function_id": "ef-testfunc1",
		"function_rule_conditions": []interface{}{
			map[string]interface{}{
				"rule_conditions": []interface{}{
					map[string]interface{}{
						"operator": "equal",
						"target":   "host",
						"values":   []interface{}{"example.com"},
					},
				},
			},
		},
	})
	d.SetId("zone-test1234#ef-testfunc1#rule-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// When TriggerType is nil in response, trigger_type should remain empty
	triggerType := d.Get("trigger_type").(string)
	assert.Equal(t, "", triggerType)
}

// TestTeoFunctionRule_UpdateTriggerType tests updating trigger_type
func TestTeoFunctionRule_UpdateTriggerType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFuncRule().client, "UseTeoV20220901Client", teoClient)

	// Mock ModifyFunctionRuleWithContext to verify TriggerType is set in request
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionRuleWithContext",
		func(_ context.Context, request *teov20220901.ModifyFunctionRuleRequest) (*teov20220901.ModifyFunctionRuleResponse, error) {
			// Verify TriggerType is set correctly
			assert.NotNil(t, request.TriggerType)
			assert.Equal(t, "region", *request.TriggerType)

			resp := teov20220901.NewModifyFunctionRuleResponse()
			resp.Response = &teov20220901.ModifyFunctionRuleResponseParams{
				RequestId: ptrStringFuncRule("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Mock DescribeTeoFunctionRuleById for the Read call after Update
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoFunctionRuleById",
		func(_ context.Context, zoneId string, functionId string, ruleId string) (*teov20220901.FunctionRule, error) {
			return &teov20220901.FunctionRule{
				RuleId:       ptrStringFuncRule("rule-12345678"),
				FunctionId:   ptrStringFuncRule("ef-testfunc1"),
				FunctionName: ptrStringFuncRule("test-func"),
				TriggerType:  ptrStringFuncRule("region"),
				Priority:     ptrInt64FuncRule(5),
				Remark:       ptrStringFuncRule("test remark"),
				FunctionRuleConditions: []*teov20220901.FunctionRuleCondition{
					{
						RuleConditions: []*teov20220901.RuleCondition{
							{
								Operator: ptrStringFuncRule("equal"),
								Target:   ptrStringFuncRule("host"),
								Values:   []*string{ptrStringFuncRule("example.com")},
							},
						},
					},
				},
			}, nil
		},
	)

	meta := newMockMetaFuncRule()
	res := svcteo.ResourceTencentCloudTeoFunctionRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"function_id":  "ef-testfunc1",
		"trigger_type": "region",
		"remark":       "test remark",
		"function_rule_conditions": []interface{}{
			map[string]interface{}{
				"rule_conditions": []interface{}{
					map[string]interface{}{
						"operator": "equal",
						"target":   "host",
						"values":   []interface{}{"example.com"},
					},
				},
			},
		},
	})
	d.SetId("zone-test1234#ef-testfunc1#rule-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify trigger_type is updated correctly
	triggerType := d.Get("trigger_type").(string)
	assert.Equal(t, "region", triggerType)
}

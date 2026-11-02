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

func TestAccTencentCloudTeoL7AccRuleV2Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7V2AccRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "zone_id", "zone-3fkff38fyw8s"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "rule_name", "Web Acceleration 1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "status", "enable"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "rule_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "rule_priority"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.condition", "${http.request.host} in ['aaa.makn.cn']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.name", "Cache"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.0.custom_time.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.0.custom_time.0.ignore_cache_control", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.0.custom_time.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.name", "CacheKey"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.query_string.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.query_string.0.values.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.condition", "lower(${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.0.name", "Cache"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.0.no_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.0.no_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.condition", "${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.0.name", "MaxAge"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.0.cache_time", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.0.follow_origin", "off"),
				),
			},
			{
				Config: testAccTeoL7V2AccRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "description.0", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "rule_name", "Web Acceleration 2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.name", "OriginPullProtocol"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.origin_pull_protocol_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.origin_pull_protocol_parameters.0.protocol", "https"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoL7V2AccRule = `
resource "tencentcloud_teo_l7_acc_rule_v2" "teo_l7_acc_rule_v2" {
  zone_id     = "zone-3fkff38fyw8s"
  description = ["1"]
  rule_name   = "Web Acceleration 1"
  status = "enable"
  branches {
    condition = "$${http.request.host} in ['aaa.makn.cn']"
    actions {
      name = "Cache"
      cache_parameters {
        custom_time {
          cache_time           = 2592000
          ignore_cache_control = "off"
          switch               = "on"
        }
      }
    }

    actions {
      name = "CacheKey"
      cache_key_parameters {
        full_url_cache = "on"
        ignore_case    = "off"
        query_string {
          switch = "off"
          values = []
        }
      }
    }

    sub_rules {
      description = ["1-1"]
      branches {
        condition = "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"
        actions {
          name = "Cache"
          cache_parameters {
            no_cache {
              switch = "on"
            }
          }
        }
      }
    }

    sub_rules {
      description = ["1-2"]
      branches {
        condition = "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"
        actions {
          name = "MaxAge"
          max_age_parameters {
            cache_time    = 3600
            follow_origin = "off"
          }
        }
      }
    }
  }
}
`

const testAccTeoL7V2AccRuleUpdate = `
resource "tencentcloud_teo_l7_acc_rule_v2" "teo_l7_acc_rule_v2" {
  zone_id     = "zone-3fkff38fyw8s"
  description = ["2"]
  rule_name   = "Web Acceleration 2"
  status = "enable"
  branches {
    condition = "$${http.request.host} in ['aaa.makn.cn']"
    actions {
      name = "Cache"
      cache_parameters {
        custom_time {
          cache_time           = 2592000
          ignore_cache_control = "off"
          switch               = "on"
        }
      }
    }
    actions {
      name = "OriginPullProtocol"
      origin_pull_protocol_parameters {
          protocol = "https"
      }
    }

    sub_rules {
      description = ["01-1"]
      branches {
        condition = "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"
        actions {
          name = "Cache"
          cache_parameters {
            no_cache {
              switch = "on"
            }
          }
        }
      }
    }

  }
}
`

// go test ./tencentcloud/services/teo/ -run "TestL7AccRuleV2" -v -count=1 -gcflags="all=-l"

// ---- Unit Tests (gomonkey mock) ----

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_" -v -count=1 -gcflags="all=-l"

type mockMetaL7AccRuleV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaL7AccRuleV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaL7AccRuleV2{}

func newMockMetaL7AccRuleV2() *mockMetaL7AccRuleV2 {
	return &mockMetaL7AccRuleV2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringL7AccRuleV2(s string) *string {
	return &s
}

// TestTeoL7AccRuleV2_CustomVariables tests custom_variables parameter in Create/Read flow
func TestTeoL7AccRuleV2_CustomVariables(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRulesWithContext", func(_ context.Context, _ *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrStringL7AccRuleV2("rule-12345678")},
			RequestId: ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL7AccRuleById", func(_ context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		return &teov20220901.DescribeL7AccRulesResponseParams{
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:   ptrStringL7AccRuleV2("rule-12345678"),
					RuleName: ptrStringL7AccRuleV2("test-rule"),
					Status:   ptrStringL7AccRuleV2("enable"),
					CustomVariables: []*teov20220901.CustomVariable{
						{
							Name:         ptrStringL7AccRuleV2("user.rule.var1"),
							InitialValue: ptrStringL7AccRuleV2("init-val"),
							Description:  ptrStringL7AccRuleV2("test variable"),
						},
					},
				},
			},
		}, nil
	})

	meta := newMockMetaL7AccRuleV2()
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test1234",
		"rule_name": "test-rule",
		"status":    "enable",
		"custom_variables": []interface{}{
			map[string]interface{}{
				"name":          "user.rule.var1",
				"initial_value": "init-val",
				"description":   "test variable",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	customVariables := d.Get("custom_variables").([]interface{})
	assert.Equal(t, 1, len(customVariables))
	cvMap := customVariables[0].(map[string]interface{})
	assert.Equal(t, "user.rule.var1", cvMap["name"])
	assert.Equal(t, "init-val", cvMap["initial_value"])
	assert.Equal(t, "test variable", cvMap["description"])
}

// TestTeoL7AccRuleV2_CustomVariableOperations tests custom_variable_operations parameter in Create/Read flow
func TestTeoL7AccRuleV2_CustomVariableOperations(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRulesWithContext", func(_ context.Context, _ *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrStringL7AccRuleV2("rule-87654321")},
			RequestId: ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL7AccRuleById", func(_ context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		return &teov20220901.DescribeL7AccRulesResponseParams{
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:   ptrStringL7AccRuleV2("rule-87654321"),
					RuleName: ptrStringL7AccRuleV2("test-rule"),
					Status:   ptrStringL7AccRuleV2("enable"),
					CustomVariableOperations: []*teov20220901.CustomVariableOperation{
						{
							Branches: []*teov20220901.CustomVariableOperationRuleBranch{
								{
									Condition: ptrStringL7AccRuleV2("${http.request.host} in ['aaa.com']"),
									Actions: []*teov20220901.CustomVariableOperationRuleAction{
										{
											Name: ptrStringL7AccRuleV2("Set"),
											SetParameters: &teov20220901.SetParameters{
												Name:  ptrStringL7AccRuleV2("user.rule.var1"),
												Value: ptrStringL7AccRuleV2("new-val"),
											},
										},
									},
									SubRules: []*teov20220901.CustomVariableOperationSubRule{
										{
											Description: []*string{ptrStringL7AccRuleV2("sub-desc")},
											Branches: []*teov20220901.CustomVariableOperationRuleBranch{
												{
													Condition: ptrStringL7AccRuleV2("${http.request.uri} in ['/path']"),
													Actions: []*teov20220901.CustomVariableOperationRuleAction{
														{
															Name: ptrStringL7AccRuleV2("Set"),
															SetParameters: &teov20220901.SetParameters{
																Name:  ptrStringL7AccRuleV2("user.rule.var2"),
																Value: ptrStringL7AccRuleV2("sub-val"),
															},
														},
													},
												},
											},
										},
									},
								},
							},
							Description: []*string{ptrStringL7AccRuleV2("op-desc")},
						},
					},
				},
			},
		}, nil
	})

	meta := newMockMetaL7AccRuleV2()
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test1234",
		"rule_name": "test-rule",
		"status":    "enable",
		"custom_variable_operations": []interface{}{
			map[string]interface{}{
				"branches": []interface{}{
					map[string]interface{}{
						"condition": "${http.request.host} in ['aaa.com']",
						"actions": []interface{}{
							map[string]interface{}{
								"name": "Set",
								"set_parameters": []interface{}{
									map[string]interface{}{
										"name":  "user.rule.var1",
										"value": "new-val",
									},
								},
							},
						},
						"sub_rules": []interface{}{
							map[string]interface{}{
								"branches": []interface{}{
									map[string]interface{}{
										"condition": "${http.request.uri} in ['/path']",
										"actions": []interface{}{
											map[string]interface{}{
												"name": "Set",
												"set_parameters": []interface{}{
													map[string]interface{}{
														"name":  "user.rule.var2",
														"value": "sub-val",
													},
												},
											},
										},
									},
								},
								"description": []interface{}{"sub-desc"},
							},
						},
					},
				},
				"description": []interface{}{"op-desc"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	operations := d.Get("custom_variable_operations").([]interface{})
	assert.Equal(t, 1, len(operations))
	opMap := operations[0].(map[string]interface{})
	branches := opMap["branches"].([]interface{})
	assert.Equal(t, 1, len(branches))
	branchMap := branches[0].(map[string]interface{})
	assert.Equal(t, "${http.request.host} in ['aaa.com']", branchMap["condition"])
	actions := branchMap["actions"].([]interface{})
	assert.Equal(t, 1, len(actions))
	actionMap := actions[0].(map[string]interface{})
	assert.Equal(t, "Set", actionMap["name"])
	setParams := actionMap["set_parameters"].([]interface{})
	assert.Equal(t, 1, len(setParams))
	setParamMap := setParams[0].(map[string]interface{})
	assert.Equal(t, "user.rule.var1", setParamMap["name"])
	assert.Equal(t, "new-val", setParamMap["value"])
	subRules := branchMap["sub_rules"].([]interface{})
	assert.Equal(t, 1, len(subRules))
	subRuleMap := subRules[0].(map[string]interface{})
	subBranches := subRuleMap["branches"].([]interface{})
	assert.Equal(t, 1, len(subBranches))
	subBranchMap := subBranches[0].(map[string]interface{})
	assert.Equal(t, "${http.request.uri} in ['/path']", subBranchMap["condition"])
}

// TestTeoL7AccRuleV2_CustomActionParameters tests custom_action_parameters in branches.actions Create/Read flow
func TestTeoL7AccRuleV2_CustomActionParameters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRulesWithContext", func(_ context.Context, _ *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrStringL7AccRuleV2("rule-aabbccdd")},
			RequestId: ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL7AccRuleById", func(_ context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		return &teov20220901.DescribeL7AccRulesResponseParams{
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:   ptrStringL7AccRuleV2("rule-aabbccdd"),
					RuleName: ptrStringL7AccRuleV2("test-rule"),
					Status:   ptrStringL7AccRuleV2("enable"),
					Branches: []*teov20220901.RuleBranch{
						{
							Condition: ptrStringL7AccRuleV2("${http.request.host} in ['bbb.com']"),
							Actions: []*teov20220901.RuleEngineAction{
								{
									Name: ptrStringL7AccRuleV2("CustomAction"),
									CustomActionParameters: &teov20220901.CustomActionParameters{
										CustomActions: []*teov20220901.CustomAction{
											{
												ActionId: ptrStringL7AccRuleV2("action-111"),
												Parameters: []*teov20220901.CustomActionParameter{
													{
														Name:        ptrStringL7AccRuleV2("Seconds"),
														ValueType:   ptrStringL7AccRuleV2("String"),
														StringValue: ptrStringL7AccRuleV2("300"),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, nil
	})

	meta := newMockMetaL7AccRuleV2()
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test1234",
		"rule_name": "test-rule",
		"status":    "enable",
		"branches": []interface{}{
			map[string]interface{}{
				"condition": "${http.request.host} in ['bbb.com']",
				"actions": []interface{}{
					map[string]interface{}{
						"name": "CustomAction",
						"custom_action_parameters": []interface{}{
							map[string]interface{}{
								"custom_actions": []interface{}{
									map[string]interface{}{
										"action_id": "action-111",
										"parameters": []interface{}{
											map[string]interface{}{
												"name":  "Seconds",
												"value": "300",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	branches := d.Get("branches").([]interface{})
	assert.Equal(t, 1, len(branches))
	branchMap := branches[0].(map[string]interface{})
	actions := branchMap["actions"].([]interface{})
	assert.Equal(t, 1, len(actions))
	actionMap := actions[0].(map[string]interface{})
	assert.Equal(t, "CustomAction", actionMap["name"])
	customActionParams := actionMap["custom_action_parameters"].([]interface{})
	assert.Equal(t, 1, len(customActionParams))
	capMap := customActionParams[0].(map[string]interface{})
	customActions := capMap["custom_actions"].([]interface{})
	assert.Equal(t, 1, len(customActions))
	customActionMap := customActions[0].(map[string]interface{})
	assert.Equal(t, "action-111", customActionMap["action_id"])
	params := customActionMap["parameters"].([]interface{})
	assert.Equal(t, 1, len(params))
	paramMap := params[0].(map[string]interface{})
	assert.Equal(t, "Seconds", paramMap["name"])
	assert.Equal(t, "300", paramMap["value"])
}

// TestTeoL7AccRuleV2_Schema tests that new parameters are defined in schema
func TestTeoL7AccRuleV2_Schema(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()

	assert.NotNil(t, res)

	assert.Contains(t, res.Schema, "custom_variables")
	customVariablesSchema := res.Schema["custom_variables"]
	assert.Equal(t, schema.TypeList, customVariablesSchema.Type)
	assert.True(t, customVariablesSchema.Optional)

	assert.Contains(t, res.Schema, "custom_variable_operations")
	cvoSchema := res.Schema["custom_variable_operations"]
	assert.Equal(t, schema.TypeList, cvoSchema.Type)
	assert.True(t, cvoSchema.Optional)
	assert.Equal(t, 1, cvoSchema.MaxItems)
}

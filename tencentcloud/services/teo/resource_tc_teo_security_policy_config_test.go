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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func TestAccTencentCloudTeoSecurityPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoSecurityPolicy,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_teo_security_policy_config_config.example", "id"),
			),
		},
			{
				ResourceName:      "tencentcloud_teo_security_policy_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoSecurityPolicy = `
resource "tencentcloud_teo_security_policy_config" "example" {
  zone_id = "zone-37u62pwxfo8s"
  entity  = "ZoneDefaultPolicy"
  security_policy {
    custom_rules {
      rules {
        name      = "rule1"
        condition = "$${http.request.host} contain ['abc']"
        enabled   = "on"
        rule_type = "PreciseMatchRule"
        priority  = 50
        action {
          name = "BlockIP"
          block_ip_action_parameters {
            duration = "120s"
          }
        }
      }

      rules {
        name      = "rule2"
        condition = "$${http.request.ip} in ['119.28.103.58']"
        enabled   = "off"
        id        = "2182252647"
        rule_type = "BasicAccessRule"
        action {
          name = "Deny"
        }
      }
    }

    managed_rules {
      enabled           = "on"
      detection_only    = "off"
      semantic_analysis = "off"
      auto_update {
        auto_update_to_latest_version = "off"
      }

      managed_rule_groups {
        group_id          = "wafgroup-webshell-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xxe-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-non-compliant-protocol-usages"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-file-upload-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-command-and-code-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ldap-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssrf-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xss-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-vulnerability-scanners"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-cms-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-other-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-sql-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-file-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-oa-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssti-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-shiro-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }
    }
  }
}
`

// mockMetaSecurityPolicy implements tccommon.ProviderMeta
type mockMetaSecurityPolicy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaSecurityPolicy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaSecurityPolicy{}

func newMockMetaSecurityPolicy() *mockMetaSecurityPolicy {
	return &mockMetaSecurityPolicy{client: &connectivity.TencentCloudClient{}}
}

func ptrStringSecurityPolicy(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestBotManagementLite_ReadWithBotManagementLite" -v -count=1 -gcflags="all=-l"
// TestBotManagementLite_ReadWithBotManagementLite tests Read flattens BotManagementLite from API response
func TestBotManagementLite_ReadWithBotManagementLite(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					CAPTCHAPageChallenge: &teov20220901.CAPTCHAPageChallenge{
						Enabled: ptrStringSecurityPolicy("on"),
					},
					AICrawlerDetection: &teov20220901.AICrawlerDetection{
						Enabled: ptrStringSecurityPolicy("on"),
						Action: &teov20220901.SecurityAction{
							Name: ptrStringSecurityPolicy("Deny"),
							DenyActionParameters: &teov20220901.DenyActionParameters{
								BlockIp:         ptrStringSecurityPolicy("on"),
								BlockIpDuration: ptrStringSecurityPolicy("120s"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	botMgmtLite := spMap["bot_management_lite"].([]interface{})
	assert.Len(t, botMgmtLite, 1)
	bmlMap := botMgmtLite[0].(map[string]interface{})

	captchaPageChallenge := bmlMap["captcha_page_challenge"].([]interface{})
	assert.Len(t, captchaPageChallenge, 1)
	cpcMap := captchaPageChallenge[0].(map[string]interface{})
	assert.Equal(t, "on", cpcMap["enabled"])

	aiCrawlerDetection := bmlMap["ai_crawler_detection"].([]interface{})
	assert.Len(t, aiCrawlerDetection, 1)
	acdMap := aiCrawlerDetection[0].(map[string]interface{})
	assert.Equal(t, "on", acdMap["enabled"])

	action := acdMap["action"].([]interface{})
	assert.Len(t, action, 1)
	actionMap := action[0].(map[string]interface{})
	assert.Equal(t, "Deny", actionMap["name"])

	denyParams := actionMap["deny_action_parameters"].([]interface{})
	assert.Len(t, denyParams, 1)
	denyMap := denyParams[0].(map[string]interface{})
	assert.Equal(t, "on", denyMap["block_ip"])
	assert.Equal(t, "120s", denyMap["block_ip_duration"])
}

// TestBotManagementLite_ReadWithNilBotManagementLite tests Read when BotManagementLite is nil
func TestBotManagementLite_ReadWithNilBotManagementLite(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{},
			RequestId:      ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	if len(securityPolicy) > 0 && securityPolicy[0] != nil {
		spMap := securityPolicy[0].(map[string]interface{})
		botMgmtLite := spMap["bot_management_lite"].([]interface{})
		assert.Len(t, botMgmtLite, 0)
	}
}

// TestBotManagementLite_ReadWithPartialBotManagementLite tests Read when only CAPTCHAPageChallenge is set
func TestBotManagementLite_ReadWithPartialBotManagementLite(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					CAPTCHAPageChallenge: &teov20220901.CAPTCHAPageChallenge{
						Enabled: ptrStringSecurityPolicy("on"),
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	botMgmtLite := spMap["bot_management_lite"].([]interface{})
	assert.Len(t, botMgmtLite, 1)
	bmlMap := botMgmtLite[0].(map[string]interface{})

	captchaPageChallenge := bmlMap["captcha_page_challenge"].([]interface{})
	assert.Len(t, captchaPageChallenge, 1)
	cpcMap := captchaPageChallenge[0].(map[string]interface{})
	assert.Equal(t, "on", cpcMap["enabled"])

	aiCrawlerDetection := bmlMap["ai_crawler_detection"].([]interface{})
	assert.Len(t, aiCrawlerDetection, 0)
}

// TestBotManagementLite_ReadWithAllowAction tests Read with Allow action parameters
func TestBotManagementLite_ReadWithAllowAction(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					AICrawlerDetection: &teov20220901.AICrawlerDetection{
						Enabled: ptrStringSecurityPolicy("on"),
						Action: &teov20220901.SecurityAction{
							Name: ptrStringSecurityPolicy("Allow"),
							AllowActionParameters: &teov20220901.AllowActionParameters{
								MinDelayTime: ptrStringSecurityPolicy("0s"),
								MaxDelayTime: ptrStringSecurityPolicy("5s"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	botMgmtLite := spMap["bot_management_lite"].([]interface{})
	assert.Len(t, botMgmtLite, 1)
	bmlMap := botMgmtLite[0].(map[string]interface{})

	aiCrawlerDetection := bmlMap["ai_crawler_detection"].([]interface{})
	assert.Len(t, aiCrawlerDetection, 1)
	acdMap := aiCrawlerDetection[0].(map[string]interface{})
	assert.Equal(t, "on", acdMap["enabled"])

	action := acdMap["action"].([]interface{})
	assert.Len(t, action, 1)
	actionMap := action[0].(map[string]interface{})
	assert.Equal(t, "Allow", actionMap["name"])

	allowParams := actionMap["allow_action_parameters"].([]interface{})
	assert.Len(t, allowParams, 1)
	allowMap := allowParams[0].(map[string]interface{})
	assert.Equal(t, "0s", allowMap["min_delay_time"])
	assert.Equal(t, "5s", allowMap["max_delay_time"])
}

// TestBotManagementLite_ReadWithChallengeAction tests Read with Challenge action parameters
func TestBotManagementLite_ReadWithChallengeAction(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					AICrawlerDetection: &teov20220901.AICrawlerDetection{
						Enabled: ptrStringSecurityPolicy("on"),
						Action: &teov20220901.SecurityAction{
							Name: ptrStringSecurityPolicy("Challenge"),
							ChallengeActionParameters: &teov20220901.ChallengeActionParameters{
								ChallengeOption: ptrStringSecurityPolicy("JSChallenge"),
								Interval:        ptrStringSecurityPolicy("300s"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	botMgmtLite := spMap["bot_management_lite"].([]interface{})
	assert.Len(t, botMgmtLite, 1)
	bmlMap := botMgmtLite[0].(map[string]interface{})

	aiCrawlerDetection := bmlMap["ai_crawler_detection"].([]interface{})
	assert.Len(t, aiCrawlerDetection, 1)
	acdMap := aiCrawlerDetection[0].(map[string]interface{})
	assert.Equal(t, "on", acdMap["enabled"])

	action := acdMap["action"].([]interface{})
	assert.Len(t, action, 1)
	actionMap := action[0].(map[string]interface{})
	assert.Equal(t, "Challenge", actionMap["name"])

	challengeParams := actionMap["challenge_action_parameters"].([]interface{})
	assert.Len(t, challengeParams, 1)
	challengeMap := challengeParams[0].(map[string]interface{})
	assert.Equal(t, "JSChallenge", challengeMap["challenge_option"])
	assert.Equal(t, "300s", challengeMap["interval"])
}

// TestBotManagementLite_UpdateExpand tests Update expands bot_management_lite into ModifySecurityPolicy request
func TestBotManagementLite_UpdateExpand(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.ModifySecurityPolicyRequest
	patches.ApplyMethodFunc(teoClient, "ModifySecurityPolicyWithContext", func(_ context.Context, request *teov20220901.ModifySecurityPolicyRequest) (*teov20220901.ModifySecurityPolicyResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewModifySecurityPolicyResponse()
		resp.Response = &teov20220901.ModifySecurityPolicyResponseParams{
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	// Also mock DescribeSecurityPolicy for the Read call after Update
	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					CAPTCHAPageChallenge: &teov20220901.CAPTCHAPageChallenge{
						Enabled: ptrStringSecurityPolicy("on"),
					},
					AICrawlerDetection: &teov20220901.AICrawlerDetection{
						Enabled: ptrStringSecurityPolicy("on"),
						Action: &teov20220901.SecurityAction{
							Name: ptrStringSecurityPolicy("Deny"),
							DenyActionParameters: &teov20220901.DenyActionParameters{
								BlockIp:         ptrStringSecurityPolicy("on"),
								BlockIpDuration: ptrStringSecurityPolicy("120s"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
		"security_policy": []interface{}{
			map[string]interface{}{
				"bot_management_lite": []interface{}{
					map[string]interface{}{
						"captcha_page_challenge": []interface{}{
							map[string]interface{}{
								"enabled": "on",
							},
						},
						"ai_crawler_detection": []interface{}{
							map[string]interface{}{
								"enabled": "on",
								"action": []interface{}{
									map[string]interface{}{
										"name": "Deny",
										"deny_action_parameters": []interface{}{
											map[string]interface{}{
												"block_ip":          "on",
												"block_ip_duration": "120s",
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
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.SecurityPolicy)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite.CAPTCHAPageChallenge)
	assert.Equal(t, "on", *capturedRequest.SecurityPolicy.BotManagementLite.CAPTCHAPageChallenge.Enabled)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection)
	assert.Equal(t, "on", *capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Enabled)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action)
	assert.Equal(t, "Deny", *capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action.Name)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters)
	assert.Equal(t, "on", *capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.BlockIp)
	assert.Equal(t, "120s", *capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.BlockIpDuration)
}

// go test ./tencentcloud/services/teo/ -run "TestSecurityPolicy_ReadWithCustomRules" -v -count=1 -gcflags="all=-l"
// TestSecurityPolicy_ReadWithCustomRules tests Read flattens CustomRules from API response
func TestSecurityPolicy_ReadWithCustomRules(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				CustomRules: &teov20220901.CustomRules{
					Rules: []*teov20220901.CustomRule{
						{
							Name:      ptrStringSecurityPolicy("precise-rule-1"),
							Condition: ptrStringSecurityPolicy("${http.request.host} contain ['test']"),
							Action: &teov20220901.SecurityAction{
								Name: ptrStringSecurityPolicy("BlockIP"),
								BlockIPActionParameters: &teov20220901.BlockIPActionParameters{
									Duration: ptrStringSecurityPolicy("120s"),
								},
							},
							Enabled:  ptrStringSecurityPolicy("on"),
							Id:       ptrStringSecurityPolicy("12345"),
							RuleType: ptrStringSecurityPolicy("PreciseMatchRule"),
							Priority: helper.Int64(50),
						},
						{
							Name:      ptrStringSecurityPolicy("basic-rule-1"),
							Condition: ptrStringSecurityPolicy("${http.request.ip} in ['1.2.3.4']"),
							Action: &teov20220901.SecurityAction{
								Name: ptrStringSecurityPolicy("Deny"),
							},
							Enabled:  ptrStringSecurityPolicy("off"),
							Id:       ptrStringSecurityPolicy("67890"),
							RuleType: ptrStringSecurityPolicy("BasicAccessRule"),
							Priority: helper.Int64(30),
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	customRules := spMap["custom_rules"].([]interface{})
	assert.Len(t, customRules, 1)
	crMap := customRules[0].(map[string]interface{})

	preciseMatchRules := crMap["precise_match_rules"].([]interface{})
	assert.Len(t, preciseMatchRules, 1)
	pmrMap := preciseMatchRules[0].(map[string]interface{})
	assert.Equal(t, "precise-rule-1", pmrMap["name"])
	assert.Equal(t, "${http.request.host} contain ['test']", pmrMap["condition"])
	assert.Equal(t, "on", pmrMap["enabled"])
	assert.Equal(t, "12345", pmrMap["id"])
	assert.Equal(t, "PreciseMatchRule", pmrMap["rule_type"])

	pmrAction := pmrMap["action"].([]interface{})
	assert.Len(t, pmrAction, 1)
	pmrActionMap := pmrAction[0].(map[string]interface{})
	assert.Equal(t, "BlockIP", pmrActionMap["name"])
	blockIPParams := pmrActionMap["block_ip_action_parameters"].([]interface{})
	assert.Len(t, blockIPParams, 1)
	blockIPMap := blockIPParams[0].(map[string]interface{})
	assert.Equal(t, "120s", blockIPMap["duration"])

	basicAccessRules := crMap["basic_access_rules"].([]interface{})
	assert.Len(t, basicAccessRules, 1)
	barMap := basicAccessRules[0].(map[string]interface{})
	assert.Equal(t, "basic-rule-1", barMap["name"])
	assert.Equal(t, "${http.request.ip} in ['1.2.3.4']", barMap["condition"])
	assert.Equal(t, "off", barMap["enabled"])
	assert.Equal(t, "67890", barMap["id"])
	assert.Equal(t, "BasicAccessRule", barMap["rule_type"])

	barAction := barMap["action"].([]interface{})
	assert.Len(t, barAction, 1)
	barActionMap := barAction[0].(map[string]interface{})
	assert.Equal(t, "Deny", barActionMap["name"])
}

// go test ./tencentcloud/services/teo/ -run "TestSecurityPolicy_ReadWithManagedRules" -v -count=1 -gcflags="all=-l"
// TestSecurityPolicy_ReadWithManagedRules tests Read flattens ManagedRules from API response
func TestSecurityPolicy_ReadWithManagedRules(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				ManagedRules: &teov20220901.ManagedRules{
					Enabled:          ptrStringSecurityPolicy("on"),
					DetectionOnly:    ptrStringSecurityPolicy("off"),
					SemanticAnalysis: ptrStringSecurityPolicy("off"),
					AutoUpdate: &teov20220901.ManagedRuleAutoUpdate{
						AutoUpdateToLatestVersion: ptrStringSecurityPolicy("on"),
						RulesetVersion:            ptrStringSecurityPolicy("v2024.01"),
					},
					ManagedRuleGroups: []*teov20220901.ManagedRuleGroup{
						{
							GroupId:          ptrStringSecurityPolicy("wafgroup-sql-injections"),
							SensitivityLevel: ptrStringSecurityPolicy("strict"),
							Action: &teov20220901.SecurityAction{
								Name: ptrStringSecurityPolicy("Deny"),
							},
						},
						{
							GroupId:          ptrStringSecurityPolicy("wafgroup-xss-attacks"),
							SensitivityLevel: ptrStringSecurityPolicy("normal"),
							Action: &teov20220901.SecurityAction{
								Name: ptrStringSecurityPolicy("Monitor"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	managedRules := spMap["managed_rules"].([]interface{})
	assert.Len(t, managedRules, 1)
	mrMap := managedRules[0].(map[string]interface{})
	assert.Equal(t, "on", mrMap["enabled"])
	assert.Equal(t, "off", mrMap["detection_only"])
	assert.Equal(t, "off", mrMap["semantic_analysis"])

	autoUpdate := mrMap["auto_update"].([]interface{})
	assert.Len(t, autoUpdate, 1)
	auMap := autoUpdate[0].(map[string]interface{})
	assert.Equal(t, "on", auMap["auto_update_to_latest_version"])
	assert.Equal(t, "v2024.01", auMap["ruleset_version"])

	managedRuleGroups := mrMap["managed_rule_groups"].(*schema.Set).List()
	assert.Len(t, managedRuleGroups, 2)
}

// go test ./tencentcloud/services/teo/ -run "TestSecurityPolicy_UpdateWithCustomRules" -v -count=1 -gcflags="all=-l"
// TestSecurityPolicy_UpdateWithCustomRules tests Update expands custom_rules into ModifySecurityPolicy request
func TestSecurityPolicy_UpdateWithCustomRules(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.ModifySecurityPolicyRequest
	patches.ApplyMethodFunc(teoClient, "ModifySecurityPolicyWithContext", func(_ context.Context, request *teov20220901.ModifySecurityPolicyRequest) (*teov20220901.ModifySecurityPolicyResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewModifySecurityPolicyResponse()
		resp.Response = &teov20220901.ModifySecurityPolicyResponseParams{
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				CustomRules: &teov20220901.CustomRules{
					Rules: []*teov20220901.CustomRule{
						{
							Name:      ptrStringSecurityPolicy("precise-rule-1"),
							Condition: ptrStringSecurityPolicy("${http.request.host} contain ['test']"),
							Action: &teov20220901.SecurityAction{
								Name: ptrStringSecurityPolicy("BlockIP"),
								BlockIPActionParameters: &teov20220901.BlockIPActionParameters{
									Duration: ptrStringSecurityPolicy("120s"),
								},
							},
							Enabled:  ptrStringSecurityPolicy("on"),
							RuleType: ptrStringSecurityPolicy("PreciseMatchRule"),
							Priority: helper.Int64(50),
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
		"security_policy": []interface{}{
			map[string]interface{}{
				"custom_rules": []interface{}{
					map[string]interface{}{
						"precise_match_rules": []interface{}{
							map[string]interface{}{
								"name":      "precise-rule-1",
								"condition": "${http.request.host} contain ['test']",
								"enabled":   "on",
								"priority":  50,
								"action": []interface{}{
									map[string]interface{}{
										"name": "BlockIP",
										"block_ip_action_parameters": []interface{}{
											map[string]interface{}{
												"duration": "120s",
											},
										},
									},
								},
							},
						},
						"basic_access_rules": []interface{}{
							map[string]interface{}{
								"name":      "basic-rule-1",
								"condition": "${http.request.ip} in ['1.2.3.4']",
								"enabled":   "off",
								"priority":  30,
								"action": []interface{}{
									map[string]interface{}{
										"name": "Deny",
									},
								},
							},
						},
					},
				},
			},
		},
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.SecurityPolicy)
	assert.NotNil(t, capturedRequest.SecurityPolicy.CustomRules)
	assert.Len(t, capturedRequest.SecurityPolicy.CustomRules.Rules, 2)

	// Verify precise match rule
	preciseRule := capturedRequest.SecurityPolicy.CustomRules.Rules[0]
	assert.Equal(t, "precise-rule-1", *preciseRule.Name)
	assert.Equal(t, "${http.request.host} contain ['test']", *preciseRule.Condition)
	assert.Equal(t, "on", *preciseRule.Enabled)
	assert.Equal(t, "PreciseMatchRule", *preciseRule.RuleType)
	assert.Equal(t, int64(50), *preciseRule.Priority)
	assert.NotNil(t, preciseRule.Action)
	assert.Equal(t, "BlockIP", *preciseRule.Action.Name)
	assert.NotNil(t, preciseRule.Action.BlockIPActionParameters)
	assert.Equal(t, "120s", *preciseRule.Action.BlockIPActionParameters.Duration)

	// Verify basic access rule
	basicRule := capturedRequest.SecurityPolicy.CustomRules.Rules[1]
	assert.Equal(t, "basic-rule-1", *basicRule.Name)
	assert.Equal(t, "${http.request.ip} in ['1.2.3.4']", *basicRule.Condition)
	assert.Equal(t, "off", *basicRule.Enabled)
	assert.Equal(t, "BasicAccessRule", *basicRule.RuleType)
	assert.Equal(t, int64(30), *basicRule.Priority)
	assert.NotNil(t, basicRule.Action)
	assert.Equal(t, "Deny", *basicRule.Action.Name)
}

// go test ./tencentcloud/services/teo/ -run "TestSecurityPolicy_UpdateWithManagedRules" -v -count=1 -gcflags="all=-l"
// TestSecurityPolicy_UpdateWithManagedRules tests Update expands managed_rules into ModifySecurityPolicy request
func TestSecurityPolicy_UpdateWithManagedRules(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.ModifySecurityPolicyRequest
	patches.ApplyMethodFunc(teoClient, "ModifySecurityPolicyWithContext", func(_ context.Context, request *teov20220901.ModifySecurityPolicyRequest) (*teov20220901.ModifySecurityPolicyResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewModifySecurityPolicyResponse()
		resp.Response = &teov20220901.ModifySecurityPolicyResponseParams{
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				ManagedRules: &teov20220901.ManagedRules{
					Enabled:          ptrStringSecurityPolicy("on"),
					DetectionOnly:    ptrStringSecurityPolicy("off"),
					SemanticAnalysis: ptrStringSecurityPolicy("off"),
					AutoUpdate: &teov20220901.ManagedRuleAutoUpdate{
						AutoUpdateToLatestVersion: ptrStringSecurityPolicy("off"),
					},
					ManagedRuleGroups: []*teov20220901.ManagedRuleGroup{
						{
							GroupId:          ptrStringSecurityPolicy("wafgroup-sql-injections"),
							SensitivityLevel: ptrStringSecurityPolicy("strict"),
							Action: &teov20220901.SecurityAction{
								Name: ptrStringSecurityPolicy("Deny"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
		"security_policy": []interface{}{
			map[string]interface{}{
				"managed_rules": []interface{}{
					map[string]interface{}{
						"enabled":           "on",
						"detection_only":    "off",
						"semantic_analysis": "off",
						"auto_update": []interface{}{
							map[string]interface{}{
								"auto_update_to_latest_version": "off",
							},
						},
					},
				},
			},
		},
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.SecurityPolicy)
	assert.NotNil(t, capturedRequest.SecurityPolicy.ManagedRules)
	assert.Equal(t, "on", *capturedRequest.SecurityPolicy.ManagedRules.Enabled)
	assert.Equal(t, "off", *capturedRequest.SecurityPolicy.ManagedRules.DetectionOnly)
	assert.Equal(t, "off", *capturedRequest.SecurityPolicy.ManagedRules.SemanticAnalysis)
	assert.NotNil(t, capturedRequest.SecurityPolicy.ManagedRules.AutoUpdate)
	assert.Equal(t, "off", *capturedRequest.SecurityPolicy.ManagedRules.AutoUpdate.AutoUpdateToLatestVersion)
}

// go test ./tencentcloud/services/teo/ -run "TestSecurityPolicy_ReadWithNilSecurityPolicy" -v -count=1 -gcflags="all=-l"
// TestSecurityPolicy_ReadWithNilSecurityPolicy tests Read when SecurityPolicy is nil (resource deleted)
func TestSecurityPolicy_ReadWithNilSecurityPolicy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: nil,
			RequestId:      ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// go test ./tencentcloud/services/teo/ -run "TestSecurityPolicy_ReadWithHostEntity" -v -count=1 -gcflags="all=-l"
// TestSecurityPolicy_ReadWithHostEntity tests Read correctly parses Host entity ID
func TestSecurityPolicy_ReadWithHostEntity(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.DescribeSecurityPolicyRequest
	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{},
			RequestId:      ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "Host",
		"host":    "www.example.com",
	})
	d.SetId("zone-12345678#Host#www.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "zone-12345678", *capturedRequest.ZoneId)
	assert.Equal(t, "Host", *capturedRequest.Entity)
	assert.Equal(t, "www.example.com", *capturedRequest.Host)
}

// TestBotManagementLite_Schema tests the bot_management_lite schema definition
func TestBotManagementLite_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "security_policy")

	spSchema := res.Schema["security_policy"]
	assert.NotNil(t, spSchema.Elem)
	spRes := spSchema.Elem.(*schema.Resource)
	assert.Contains(t, spRes.Schema, "bot_management_lite")

	bmlSchema := spRes.Schema["bot_management_lite"]
	assert.Equal(t, schema.TypeList, bmlSchema.Type)
	assert.True(t, bmlSchema.Optional)
	assert.True(t, bmlSchema.Computed)
	assert.Equal(t, 1, bmlSchema.MaxItems)

	bmlRes := bmlSchema.Elem.(*schema.Resource)
	assert.Contains(t, bmlRes.Schema, "captcha_page_challenge")
	assert.Contains(t, bmlRes.Schema, "ai_crawler_detection")

	cpcSchema := bmlRes.Schema["captcha_page_challenge"]
	assert.Equal(t, schema.TypeList, cpcSchema.Type)
	assert.True(t, cpcSchema.Optional)
	assert.Equal(t, 1, cpcSchema.MaxItems)

	cpcRes := cpcSchema.Elem.(*schema.Resource)
	assert.Contains(t, cpcRes.Schema, "enabled")
	assert.Equal(t, schema.TypeString, cpcRes.Schema["enabled"].Type)
	assert.True(t, cpcRes.Schema["enabled"].Required)

	acdSchema := bmlRes.Schema["ai_crawler_detection"]
	assert.Equal(t, schema.TypeList, acdSchema.Type)
	assert.True(t, acdSchema.Optional)
	assert.Equal(t, 1, acdSchema.MaxItems)

	acdRes := acdSchema.Elem.(*schema.Resource)
	assert.Contains(t, acdRes.Schema, "enabled")
	assert.Contains(t, acdRes.Schema, "action")
	assert.Equal(t, schema.TypeString, acdRes.Schema["enabled"].Type)
	assert.True(t, acdRes.Schema["enabled"].Required)

	actionSchema := acdRes.Schema["action"]
	assert.Equal(t, schema.TypeList, actionSchema.Type)
	assert.True(t, actionSchema.Optional)
	assert.Equal(t, 1, actionSchema.MaxItems)

	actionRes := actionSchema.Elem.(*schema.Resource)
	assert.Contains(t, actionRes.Schema, "name")
	assert.Contains(t, actionRes.Schema, "deny_action_parameters")
	assert.Contains(t, actionRes.Schema, "allow_action_parameters")
	assert.Contains(t, actionRes.Schema, "challenge_action_parameters")

	denySchema := actionRes.Schema["deny_action_parameters"]
	denyRes := denySchema.Elem.(*schema.Resource)
	assert.Contains(t, denyRes.Schema, "block_ip")
	assert.Contains(t, denyRes.Schema, "block_ip_duration")
	assert.Contains(t, denyRes.Schema, "return_custom_page")
	assert.Contains(t, denyRes.Schema, "response_code")
	assert.Contains(t, denyRes.Schema, "error_page_id")
	assert.Contains(t, denyRes.Schema, "stall")

	allowSchema := actionRes.Schema["allow_action_parameters"]
	allowRes := allowSchema.Elem.(*schema.Resource)
	assert.Contains(t, allowRes.Schema, "min_delay_time")
	assert.Contains(t, allowRes.Schema, "max_delay_time")

	challengeSchema := actionRes.Schema["challenge_action_parameters"]
	challengeRes := challengeSchema.Elem.(*schema.Resource)
	assert.Contains(t, challengeRes.Schema, "challenge_option")
	assert.Contains(t, challengeRes.Schema, "interval")
	assert.Contains(t, challengeRes.Schema, "attester_id")
}

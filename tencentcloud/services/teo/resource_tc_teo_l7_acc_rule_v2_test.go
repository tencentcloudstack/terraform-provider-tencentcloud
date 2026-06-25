package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
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

// mockMetaL7AccRuleV2 implements tccommon.ProviderMeta
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

func ptrInt64L7AccRuleV2(i int64) *int64 {
	return &i
}

// TestL7AccRuleV2_flattenRules tests flattenTeoL7AccRuleV2Rules with full RuleEngineItem data
func TestL7AccRuleV2_flattenRules(t *testing.T) {
	rules := []*teov20220901.RuleEngineItem{
		{
			RuleId:       ptrStringL7AccRuleV2("rule-abc123"),
			Status:       ptrStringL7AccRuleV2("enable"),
			RuleName:     ptrStringL7AccRuleV2("Test Rule"),
			Description:  []*string{ptrStringL7AccRuleV2("desc1"), ptrStringL7AccRuleV2("desc2")},
			RulePriority: ptrInt64L7AccRuleV2(100),
		},
	}

	// We need to call the flatten function from the teo package.
	// Since flattenTeoL7AccRuleV2Rules is in package teo, we use it through teo.
	// But it's unexported, so we can't call it directly from teo_test.
	// We'll instead test via the Read path using gomonkey.
	// For now, let's skip direct flatten testing and go through Read path.

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7AccRuleV2(1),
			Rules:      rules,
			RequestId:  ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123#rule-abc123")

	_ = res.Read(d, meta)

	// Verify rules field is populated
	rulesOut := d.Get("rules").([]interface{})
	assert.Len(t, rulesOut, 1)
	ruleMap := rulesOut[0].(map[string]interface{})
	assert.Equal(t, "rule-abc123", ruleMap["rule_id"].(string))
	assert.Equal(t, "enable", ruleMap["status"].(string))
	assert.Equal(t, "Test Rule", ruleMap["rule_name"].(string))
	assert.Equal(t, 100, ruleMap["rule_priority"].(int))

	descriptions := ruleMap["description"].([]interface{})
	assert.Len(t, descriptions, 2)
}

// TestL7AccRuleV2_flattenRulesEmpty tests flattenTeoL7AccRuleV2Rules with empty rules list
func TestL7AccRuleV2_flattenRulesEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7AccRuleV2(0),
			Rules:      nil,
			RequestId:  ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123#rule-notexist")

	_ = res.Read(d, meta)

	// When rules list is empty, the resource should be removed from state
	assert.Equal(t, "", d.Id())
}

// TestL7AccRuleV2_flattenRulesNilFields tests flatten with nil sub-fields
func TestL7AccRuleV2_flattenRulesNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7AccRuleV2(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrStringL7AccRuleV2("rule-minimal"),
					Status:       nil,
					RuleName:     nil,
					Description:  nil,
					Branches:     nil,
					RulePriority: nil,
				},
			},
			RequestId: ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123#rule-minimal")

	_ = res.Read(d, meta)

	// Verify rules field is populated even with nil sub-fields
	rulesOut := d.Get("rules").([]interface{})
	assert.Len(t, rulesOut, 1)
	ruleMap := rulesOut[0].(map[string]interface{})
	assert.Equal(t, "rule-minimal", ruleMap["rule_id"].(string))
	// status should be empty string since it was nil in the API response
	assert.Equal(t, "", ruleMap["status"].(string))
	assert.Equal(t, "", ruleMap["rule_name"].(string))
}

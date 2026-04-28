package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
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

// TestL7AccRuleV2_RuleIds_Create_Success tests that rule_ids is populated after Create
func TestL7AccRuleV2_RuleIds_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrStringL7AccRuleV2("rule-aaa"), ptrStringL7AccRuleV2("rule-bbb")},
			RequestId: ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7AccRuleV2(2),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:   ptrStringL7AccRuleV2("rule-aaa"),
					RuleName: ptrStringL7AccRuleV2("test-rule-1"),
					Status:   ptrStringL7AccRuleV2("enable"),
				},
				{
					RuleId:   ptrStringL7AccRuleV2("rule-bbb"),
					RuleName: ptrStringL7AccRuleV2("test-rule-2"),
					Status:   ptrStringL7AccRuleV2("enable"),
				},
			},
			RequestId: ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-1234567890",
		"rule_name": "test-rule",
		"status":    "enable",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#rule-aaa", d.Id())

	ruleIds := d.Get("rule_ids").([]interface{})
	assert.Equal(t, 2, len(ruleIds))
	assert.Equal(t, "rule-aaa", ruleIds[0].(string))
	assert.Equal(t, "rule-bbb", ruleIds[1].(string))
}

// TestL7AccRuleV2_RuleIds_Read_Success tests that rule_ids is populated during Read
func TestL7AccRuleV2_RuleIds_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7AccRuleV2(2),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:   ptrStringL7AccRuleV2("rule-aaa"),
					RuleName: ptrStringL7AccRuleV2("test-rule-1"),
					Status:   ptrStringL7AccRuleV2("enable"),
				},
				{
					RuleId:   ptrStringL7AccRuleV2("rule-bbb"),
					RuleName: ptrStringL7AccRuleV2("test-rule-2"),
					Status:   ptrStringL7AccRuleV2("disable"),
				},
			},
			RequestId: ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890#rule-aaa")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	ruleIds := d.Get("rule_ids").([]interface{})
	assert.Equal(t, 2, len(ruleIds))
	assert.Equal(t, "rule-aaa", ruleIds[0].(string))
	assert.Equal(t, "rule-bbb", ruleIds[1].(string))
}

// TestL7AccRuleV2_RuleIds_Schema tests the schema definition of rule_ids
func TestL7AccRuleV2_RuleIds_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()

	assert.Contains(t, res.Schema, "rule_ids")
	ruleIds := res.Schema["rule_ids"]
	assert.Equal(t, schema.TypeList, ruleIds.Type)
	assert.True(t, ruleIds.Computed)
	assert.False(t, ruleIds.Required)
	assert.False(t, ruleIds.Optional)
}

// TestL7AccRuleV2_RuleIds_Create_APIError tests Create handles API error
func TestL7AccRuleV2_RuleIds_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-invalid",
		"rule_name": "test-rule",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestL7AccRuleV2_RuleIds_Read_NotFound tests Read handles rule not found
func TestL7AccRuleV2_RuleIds_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7AccRuleV2(0),
			Rules:      []*teov20220901.RuleEngineItem{},
			RequestId:  ptrStringL7AccRuleV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890#rule-nonexistent")

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

func ptrStringL7AccRuleV2(s string) *string {
	return &s
}

func ptrInt64L7AccRuleV2(i int64) *int64 {
	return &i
}

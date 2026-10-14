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

func ptrInt64L7AccRuleV2(n int64) *int64 {
	return &n
}

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

// ---- Unit Tests (gomonkey mock) ----

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_" -v -count=1 -gcflags="all=-l"

// TestTeoL7AccRuleV2_Create tests that Create maps zone_id to request.ZoneId
// and rule fields to request.Rules correctly
func TestTeoL7AccRuleV2_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	// Mock CreateL7AccRules to return a response with RuleIds
	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules",
		func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
			// Verify zone_id is mapped to request.ZoneId
			assert.NotNil(t, request.ZoneId)
			assert.Equal(t, "zone-test1234", *request.ZoneId)

			// Verify rule fields are mapped to request.Rules
			assert.NotNil(t, request.Rules)
			assert.Len(t, request.Rules, 1)
			if len(request.Rules) > 0 {
				rule := request.Rules[0]
				assert.NotNil(t, rule.Status)
				assert.Equal(t, "enable", *rule.Status)
				assert.NotNil(t, rule.RuleName)
				assert.Equal(t, "test-rule", *rule.RuleName)
			}

			resp := teov20220901.NewCreateL7AccRulesResponse()
			resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
				RuleIds:   []*string{ptrStringL7AccRuleV2("rule-abc123")},
				RequestId: ptrStringL7AccRuleV2("fake-request-id"),
			}
			return resp, nil
		})

	// Mock TeoService.DescribeTeoL7AccRuleById for the Read call after Create
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL7AccRuleById",
		func(_ context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
			assert.Equal(t, "zone-test1234", zoneId)
			assert.Equal(t, "rule-abc123", ruleId)
			return &teov20220901.DescribeL7AccRulesResponseParams{
				TotalCount: ptrInt64L7AccRuleV2(1),
				Rules: []*teov20220901.RuleEngineItem{
					{
						Status:       ptrStringL7AccRuleV2("enable"),
						RuleId:       ptrStringL7AccRuleV2("rule-abc123"),
						RuleName:     ptrStringL7AccRuleV2("test-rule"),
						RulePriority: ptrInt64L7AccRuleV2(1),
					},
				},
			}, nil
		})

	meta := newMockMetaL7AccRuleV2()
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test1234",
		"status":    "enable",
		"rule_name": "test-rule",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify rule_id is set correctly
	ruleId := d.Get("rule_id").(string)
	assert.Equal(t, "rule-abc123", ruleId)

	// Verify composite ID
	assert.Equal(t, "zone-test1234#rule-abc123", d.Id())
}

// TestTeoL7AccRuleV2_Read tests that Read uses correct Filters and maps response fields
func TestTeoL7AccRuleV2_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock TeoService.DescribeTeoL7AccRuleById to verify zoneId and ruleId are passed correctly
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL7AccRuleById",
		func(_ context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
			// Verify that zoneId and ruleId are correctly parsed from composite ID
			assert.Equal(t, "zone-test1234", zoneId)
			assert.Equal(t, "rule-xyz789", ruleId)
			return &teov20220901.DescribeL7AccRulesResponseParams{
				TotalCount: ptrInt64L7AccRuleV2(1),
				Rules: []*teov20220901.RuleEngineItem{
					{
						Status:       ptrStringL7AccRuleV2("enable"),
						RuleId:       ptrStringL7AccRuleV2("rule-xyz789"),
						RuleName:     ptrStringL7AccRuleV2("my-rule"),
						Description:  []*string{ptrStringL7AccRuleV2("desc1")},
						RulePriority: ptrInt64L7AccRuleV2(5),
					},
				},
			}, nil
		})

	meta := newMockMetaL7AccRuleV2()
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test1234",
	})
	d.SetId("zone-test1234#rule-xyz789")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify fields are set correctly from Read response
	assert.Equal(t, "enable", d.Get("status").(string))
	assert.Equal(t, "my-rule", d.Get("rule_name").(string))
	assert.Equal(t, "rule-xyz789", d.Get("rule_id").(string))
}

// TestTeoL7AccRuleV2_Read_NotFound tests read when resource is not found
func TestTeoL7AccRuleV2_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock TeoService.DescribeTeoL7AccRuleById to return empty rules (not found)
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL7AccRuleById",
		func(_ context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
			return &teov20220901.DescribeL7AccRulesResponseParams{
				TotalCount: ptrInt64L7AccRuleV2(0),
				Rules:      []*teov20220901.RuleEngineItem{},
			}, nil
		})

	meta := newMockMetaL7AccRuleV2()
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoL7AccRuleV2_Update tests that Update maps zone_id to request.ZoneId
// and rule_id to request.Rule.RuleId, and other fields to request.Rule
func TestTeoL7AccRuleV2_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	// Mock ModifyL7AccRule to verify parameter mapping
	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule",
		func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
			// Verify zone_id is mapped to request.ZoneId
			assert.NotNil(t, request.ZoneId)
			assert.Equal(t, "zone-test1234", *request.ZoneId)

			// Verify rule fields are mapped to request.Rule
			assert.NotNil(t, request.Rule)
			assert.NotNil(t, request.Rule.RuleId)
			assert.Equal(t, "rule-abc123", *request.Rule.RuleId)
			assert.NotNil(t, request.Rule.Status)
			assert.Equal(t, "disable", *request.Rule.Status)
			assert.NotNil(t, request.Rule.RuleName)
			assert.Equal(t, "updated-rule", *request.Rule.RuleName)

			resp := teov20220901.NewModifyL7AccRuleResponse()
			resp.Response = &teov20220901.ModifyL7AccRuleResponseParams{
				RequestId: ptrStringL7AccRuleV2("fake-request-id"),
			}
			return resp, nil
		})

	// Mock TeoService.DescribeTeoL7AccRuleById for the Read call after Update
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL7AccRuleById",
		func(_ context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
			return &teov20220901.DescribeL7AccRulesResponseParams{
				TotalCount: ptrInt64L7AccRuleV2(1),
				Rules: []*teov20220901.RuleEngineItem{
					{
						Status:       ptrStringL7AccRuleV2("disable"),
						RuleId:       ptrStringL7AccRuleV2("rule-abc123"),
						RuleName:     ptrStringL7AccRuleV2("updated-rule"),
						RulePriority: ptrInt64L7AccRuleV2(1),
					},
				},
			}, nil
		})

	meta := newMockMetaL7AccRuleV2()
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test1234",
		"status":    "disable",
		"rule_name": "updated-rule",
	})
	d.SetId("zone-test1234#rule-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoL7AccRuleV2_Delete tests that Delete maps zone_id to request.ZoneId
// and rule_id to request.RuleIds
func TestTeoL7AccRuleV2_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	// Mock DeleteL7AccRules to verify parameter mapping
	patches.ApplyMethodFunc(teoClient, "DeleteL7AccRules",
		func(request *teov20220901.DeleteL7AccRulesRequest) (*teov20220901.DeleteL7AccRulesResponse, error) {
			// Verify zone_id is mapped to request.ZoneId
			assert.NotNil(t, request.ZoneId)
			assert.Equal(t, "zone-test1234", *request.ZoneId)

			// Verify rule_id is mapped to request.RuleIds
			assert.NotNil(t, request.RuleIds)
			assert.Len(t, request.RuleIds, 1)
			if len(request.RuleIds) > 0 {
				assert.Equal(t, "rule-abc123", *request.RuleIds[0])
			}

			resp := teov20220901.NewDeleteL7AccRulesResponse()
			resp.Response = &teov20220901.DeleteL7AccRulesResponseParams{
				RequestId: ptrStringL7AccRuleV2("fake-request-id"),
			}
			return resp, nil
		})

	// Mock TeoService.DescribeTeoL7AccRuleById for the Read call after Delete
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL7AccRuleById",
		func(_ context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
			return &teov20220901.DescribeL7AccRulesResponseParams{
				TotalCount: ptrInt64L7AccRuleV2(0),
				Rules:      []*teov20220901.RuleEngineItem{},
			}, nil
		})

	meta := newMockMetaL7AccRuleV2()
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test1234",
	})
	d.SetId("zone-test1234#rule-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoL7AccRuleV2_Schema tests that the schema has the correct attributes
func TestTeoL7AccRuleV2_Schema(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoL7AccRuleV2()

	assert.NotNil(t, res)

	// Verify zone_id is Required and ForceNew
	assert.Contains(t, res.Schema, "zone_id")
	zoneIdSchema := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneIdSchema.Type)
	assert.True(t, zoneIdSchema.Required)
	assert.True(t, zoneIdSchema.ForceNew)

	// Verify rule_id is Computed
	assert.Contains(t, res.Schema, "rule_id")
	ruleIdSchema := res.Schema["rule_id"]
	assert.Equal(t, schema.TypeString, ruleIdSchema.Type)
	assert.True(t, ruleIdSchema.Computed)

	// Verify rule_priority is Computed
	assert.Contains(t, res.Schema, "rule_priority")
	rulePrioritySchema := res.Schema["rule_priority"]
	assert.Equal(t, schema.TypeInt, rulePrioritySchema.Type)
	assert.True(t, rulePrioritySchema.Computed)

	// Verify status is Optional
	assert.Contains(t, res.Schema, "status")
	statusSchema := res.Schema["status"]
	assert.Equal(t, schema.TypeString, statusSchema.Type)
	assert.True(t, statusSchema.Optional)

	// Verify rule_name is Optional
	assert.Contains(t, res.Schema, "rule_name")
	ruleNameSchema := res.Schema["rule_name"]
	assert.Equal(t, schema.TypeString, ruleNameSchema.Type)
	assert.True(t, ruleNameSchema.Optional)

	// Verify description is Optional
	assert.Contains(t, res.Schema, "description")
	descSchema := res.Schema["description"]
	assert.Equal(t, schema.TypeList, descSchema.Type)
	assert.True(t, descSchema.Optional)

	// Verify branches is Optional
	assert.Contains(t, res.Schema, "branches")
	branchesSchema := res.Schema["branches"]
	assert.Equal(t, schema.TypeList, branchesSchema.Type)
	assert.True(t, branchesSchema.Optional)

	// Verify no rule_ids attribute exists
	assert.NotContains(t, res.Schema, "rule_ids")
}

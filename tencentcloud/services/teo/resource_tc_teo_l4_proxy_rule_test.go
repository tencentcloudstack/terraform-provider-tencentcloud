package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// mockMetaL4ProxyRule implements tccommon.ProviderMeta
type mockMetaL4ProxyRule struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaL4ProxyRule) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaL4ProxyRule{}

func newMockMetaL4ProxyRule() *mockMetaL4ProxyRule {
	return &mockMetaL4ProxyRule{client: &connectivity.TencentCloudClient{}}
}

func ptrStringL4ProxyRule(s string) *string {
	return &s
}

func ptrUint64L4ProxyRule(v uint64) *uint64 {
	return &v
}

// go test -test.run TestAccTencentCloudTeoL4ProxyRuleResource_basic -v -timeout=0
func TestAccTencentCloudTeoL4ProxyRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL4ProxyRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.client_ip_pass_through_mode", "OFF"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.origin_port_range", "1212"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.origin_type", "IP_DOMAIN"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.origin_value.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.origin_value.0", "www.aaa.com"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.port_range.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.port_range.0", "1212"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.rule_tag", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.session_persist", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.session_persist_time", "3600"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoL4ProxyRuleUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.client_ip_pass_through_mode", "OFF"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.origin_port_range", "1213"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.origin_type", "IP_DOMAIN"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.origin_value.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.origin_value.0", "www.bbb.com"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.port_range.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.port_range.0", "1213"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.rule_tag", "bbb"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.session_persist", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule", "l4_proxy_rules.0.session_persist_time", "3600"),
				),
			},
		},
	})
}

const testAccTeoL4ProxyRule = `

resource "tencentcloud_teo_l4_proxy_rule" "teo_l4_proxy_rule" {
    proxy_id = "sid-38hbn26osico"
    zone_id  = "zone-36bjhygh1bxe"

    l4_proxy_rules {
        client_ip_pass_through_mode = "OFF"
        origin_port_range           = "1212"
        origin_type                 = "IP_DOMAIN"
        origin_value                = [
            "www.aaa.com",
        ]
        port_range                  = [
            "1212",
        ]
        protocol                    = "TCP"
        rule_tag                    = "aaa"
        session_persist             = "off"
        session_persist_time        = 3600
    }
}
`

const testAccTeoL4ProxyRuleUp = `

resource "tencentcloud_teo_l4_proxy_rule" "teo_l4_proxy_rule" {
    proxy_id = "sid-38hbn26osico"
    zone_id  = "zone-36bjhygh1bxe"

    l4_proxy_rules {
        client_ip_pass_through_mode = "OFF"
        origin_port_range           = "1213"
        origin_type                 = "IP_DOMAIN"
        origin_value                = [
            "www.bbb.com",
        ]
        port_range                  = [
            "1213",
        ]
        protocol                    = "TCP"
        rule_tag                    = "bbb"
        session_persist             = "off"
        session_persist_time        = 3600
    }
}
`

// go test ./tencentcloud/services/teo/ -run "TestL4ProxyRuleL4proxyRuleIds" -v -count=1 -gcflags="all=-l"

// TestL4ProxyRuleL4proxyRuleIds_Schema validates l4proxy_rule_ids schema definition
func TestL4ProxyRuleL4proxyRuleIds_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoL4ProxyRule()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "l4proxy_rule_ids")

	l4proxyRuleIdsSchema := res.Schema["l4proxy_rule_ids"]
	assert.Equal(t, schema.TypeList, l4proxyRuleIdsSchema.Type)
	assert.True(t, l4proxyRuleIdsSchema.Computed)
	assert.False(t, l4proxyRuleIdsSchema.Required)
	assert.False(t, l4proxyRuleIdsSchema.Optional)
	assert.NotNil(t, l4proxyRuleIdsSchema.Elem)
}

// TestL4ProxyRuleL4proxyRuleIds_Read tests that l4proxy_rule_ids is populated in read from composite ID
func TestL4ProxyRuleL4proxyRuleIds_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL4ProxyRule().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMetaL4ProxyRule().client, "UseTeoClient", teoClient)

	// Patch DescribeL4ProxyRules for the read function's service call
	patches.ApplyMethodFunc(teoClient, "DescribeL4ProxyRules", func(request *teov20220901.DescribeL4ProxyRulesRequest) (*teov20220901.DescribeL4ProxyRulesResponse, error) {
		resp := teov20220901.NewDescribeL4ProxyRulesResponse()
		resp.Response = &teov20220901.DescribeL4ProxyRulesResponseParams{
			L4ProxyRules: []*teov20220901.L4ProxyRule{
				{
					RuleId:                  ptrStringL4ProxyRule("rule-xyz789"),
					Protocol:                ptrStringL4ProxyRule("TCP"),
					PortRange:               []*string{ptrStringL4ProxyRule("80")},
					OriginType:              ptrStringL4ProxyRule("IP_DOMAIN"),
					OriginValue:             []*string{ptrStringL4ProxyRule("8.8.8.8")},
					OriginPortRange:         ptrStringL4ProxyRule("80"),
					ClientIPPassThroughMode: ptrStringL4ProxyRule("OFF"),
					SessionPersist:          ptrStringL4ProxyRule("off"),
					SessionPersistTime:      ptrUint64L4ProxyRule(3600),
					Status:                  ptrStringL4ProxyRule("online"),
				},
			},
			RequestId: ptrStringL4ProxyRule("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL4ProxyRule()
	res := teo.ResourceTencentCloudTeoL4ProxyRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-12345678",
		"proxy_id": "sid-abcdef",
		"l4_proxy_rules": []interface{}{
			map[string]interface{}{
				"protocol":          "TCP",
				"port_range":        []interface{}{"80"},
				"origin_type":       "IP_DOMAIN",
				"origin_value":      []interface{}{"8.8.8.8"},
				"origin_port_range": "80",
			},
		},
	})
	d.SetId("zone-12345678#sid-abcdef#rule-xyz789")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify l4proxy_rule_ids is set from composite ID
	l4proxyRuleIds := d.Get("l4proxy_rule_ids").([]interface{})
	assert.Equal(t, 1, len(l4proxyRuleIds))
	assert.Equal(t, "rule-xyz789", l4proxyRuleIds[0].(string))
}

// TestL4ProxyRuleL4proxyRuleIds_Create tests that l4proxy_rule_ids is set after create
func TestL4ProxyRuleL4proxyRuleIds_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL4ProxyRule().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMetaL4ProxyRule().client, "UseTeoClient", teoClient)

	// Patch CreateL4ProxyRulesWithContext to return L4ProxyRuleIds
	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyRulesWithContext", func(ctx context.Context, request *teov20220901.CreateL4ProxyRulesRequest) (*teov20220901.CreateL4ProxyRulesResponse, error) {
		resp := teov20220901.NewCreateL4ProxyRulesResponse()
		resp.Response = &teov20220901.CreateL4ProxyRulesResponseParams{
			L4ProxyRuleIds: []*string{
				ptrStringL4ProxyRule("rule-abc123"),
			},
			RequestId: ptrStringL4ProxyRule("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeL4ProxyRules for the state refresh and read function
	patches.ApplyMethodFunc(teoClient, "DescribeL4ProxyRules", func(request *teov20220901.DescribeL4ProxyRulesRequest) (*teov20220901.DescribeL4ProxyRulesResponse, error) {
		resp := teov20220901.NewDescribeL4ProxyRulesResponse()
		resp.Response = &teov20220901.DescribeL4ProxyRulesResponseParams{
			L4ProxyRules: []*teov20220901.L4ProxyRule{
				{
					RuleId:                  ptrStringL4ProxyRule("rule-abc123"),
					Protocol:                ptrStringL4ProxyRule("TCP"),
					PortRange:               []*string{ptrStringL4ProxyRule("1212")},
					OriginType:              ptrStringL4ProxyRule("IP_DOMAIN"),
					OriginValue:             []*string{ptrStringL4ProxyRule("8.8.8.8")},
					OriginPortRange:         ptrStringL4ProxyRule("1212"),
					ClientIPPassThroughMode: ptrStringL4ProxyRule("OFF"),
					SessionPersist:          ptrStringL4ProxyRule("off"),
					SessionPersistTime:      ptrUint64L4ProxyRule(3600),
					Status:                  ptrStringL4ProxyRule("online"),
				},
			},
			RequestId: ptrStringL4ProxyRule("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL4ProxyRule()
	res := teo.ResourceTencentCloudTeoL4ProxyRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-12345678",
		"proxy_id": "sid-abcdef",
		"l4_proxy_rules": []interface{}{
			map[string]interface{}{
				"protocol":          "TCP",
				"port_range":        []interface{}{"1212"},
				"origin_type":       "IP_DOMAIN",
				"origin_value":      []interface{}{"8.8.8.8"},
				"origin_port_range": "1212",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify l4proxy_rule_ids is set after create
	l4proxyRuleIds := d.Get("l4proxy_rule_ids").([]interface{})
	assert.Equal(t, 1, len(l4proxyRuleIds))
	assert.Equal(t, "rule-abc123", l4proxyRuleIds[0].(string))
}

// TestL4ProxyRuleL4proxyRuleIds_APIError tests create handles API error
func TestL4ProxyRuleL4proxyRuleIds_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL4ProxyRule().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyRulesWithContext", func(ctx context.Context, request *teov20220901.CreateL4ProxyRulesRequest) (*teov20220901.CreateL4ProxyRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMetaL4ProxyRule()
	res := teo.ResourceTencentCloudTeoL4ProxyRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-invalid",
		"proxy_id": "sid-invalid",
		"l4_proxy_rules": []interface{}{
			map[string]interface{}{
				"protocol":          "TCP",
				"port_range":        []interface{}{"80"},
				"origin_type":       "IP_DOMAIN",
				"origin_value":      []interface{}{"8.8.8.8"},
				"origin_port_range": "80",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

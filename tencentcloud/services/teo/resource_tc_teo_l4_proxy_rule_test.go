package teo_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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

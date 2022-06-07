package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCdnDomainDataSources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainData,
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domains.cdn_domain", "domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domains.cdn_domain", "service_type"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "area", "mainland"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "full_url_cache", "false"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "origin.0.origin_type", "ip"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "origin.0.origin_list.#", "1"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "origin.0.server_name", "test.zhaoshaona.com"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "origin.0.origin_pull_protocol", "follow"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "https_config.0.https_switch", "on"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "https_config.0.http2_switch", "on"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "https_config.0.ocsp_stapling_switch", "on"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "https_config.0.spdy_switch", "off"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "https_config.0.verify_client", "off"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "https_config.0.server_certificate_config.0.message", "test"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domains.cdn_domain", "https_config.0.server_certificate_config.0.deploy_time"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domains.cdn_domain", "https_config.0.server_certificate_config.0.expire_time"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "tags.test", "world"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "range_origin_switch", "off"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.cache_time", "10000"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.rule_paths.#", "1"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.rule_type", "default"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.switch", "off"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.compare_max_age", "off"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.ignore_cache_control", "off"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.ignore_set_cookie", "off"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.no_cache_switch", "on"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.re_validate", "on"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "rule_cache.0.follow_origin_switch", "off"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "request_header.0.switch", "on"),
					//resource.TestCheckResourceAttr("data.tencentcloud_cdn_domains.cdn_domain", "request_header.0.header_rules.#", "1"),
				),
			},
		},
	})
}

const testAccCdnDomainData = `
data "tencentcloud_cdn_domains" "cdn_domain" {
  service_type = "web"
}
`

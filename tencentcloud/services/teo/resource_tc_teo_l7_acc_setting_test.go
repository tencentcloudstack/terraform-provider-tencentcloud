package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoL7AccSettingResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccSetting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_id", "zone-36bjhygh1bxe"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.accelerate_mainland.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_strategy", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.no_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.action", "includeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.cache_time_percent", "90"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_country.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.algorithms.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.redirect_status_code", "302"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.grpc.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.include_sub_domains", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.preload", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.timeout", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ipv6.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.cache_time", "600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.follow_origin", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ocsp_stapling.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.offline_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.max_size", "838860800"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.quic.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.smart_routing.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.expires", "1969-12-31T16:00:00Z"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.cipher_suite", "loose-v2023"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.version.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.upstream_http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.timeout", "30"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoL7AccSettingUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_id", "zone-36bjhygh1bxe"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.accelerate_mainland.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_strategy", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.no_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.action", "includeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.cache_time_percent", "90"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_country.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.algorithms.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.redirect_status_code", "302"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.grpc.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.include_sub_domains", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.preload", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.timeout", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ipv6.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.cache_time", "600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.follow_origin", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ocsp_stapling.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.offline_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.max_size", "838860800"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.quic.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.smart_routing.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.expires", "1969-12-31T16:00:00Z"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.cipher_suite", "loose-v2023"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.version.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.upstream_http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.timeout", "30"),
				),
			},
			{
				Config: testAccTeoL7AccSetting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_id", "zone-36bjhygh1bxe"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.accelerate_mainland.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_strategy", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.no_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.action", "includeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.cache_time_percent", "90"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_country.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.algorithms.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.redirect_status_code", "302"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.grpc.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.include_sub_domains", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.preload", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.timeout", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ipv6.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.cache_time", "600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.follow_origin", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ocsp_stapling.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.offline_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.max_size", "838860800"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.quic.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.smart_routing.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.expires", "1969-12-31T16:00:00Z"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.cipher_suite", "loose-v2023"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.version.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.upstream_http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.timeout", "30"),
				),
			},
		},
	})
}

const testAccTeoL7AccSetting = `

resource "tencentcloud_teo_l7_acc_setting" "teo_l7_acc_setting" {
  zone_id = "zone-36bjhygh1bxe"
  zone_config {
    accelerate_mainland {
      switch = "off"
    }
    cache {
      custom_time {
        cache_time = 2592000
        switch     = "off"
      }
      follow_origin {
        default_cache          = "off"
        default_cache_strategy = "on"
        default_cache_time     = 0
        switch                 = "on"
      }
      no_cache {
        switch = "off"
      }
    }
    cache_key {
      full_url_cache = "on"
      ignore_case    = "off"
      query_string {
        action = "includeCustom"
        switch = "off"
      }
    }
    cache_prefresh {
      cache_time_percent = 90
      switch             = "off"
    }
    client_ip_country {
      switch      = "off"
    }
    client_ip_header {
      switch      = "off"
    }
    compression {
      algorithms = ["brotli", "gzip"]
      switch     = "on"
    }
    force_redirect_https {
      redirect_status_code = 302
      switch               = "off"
    }
    grpc {
      switch = "off"
    }
    hsts {
      include_sub_domains = "off"
      preload             = "off"
      switch              = "off"
      timeout             = 0
    }
    http2 {
      switch = "off"
    }
    ipv6 {
      switch = "off"
    }
    max_age {
      cache_time    = 600
      follow_origin = "on"
    }
    ocsp_stapling {
      switch = "off"
    }
    offline_cache {
      switch = "on"
    }
    post_max_size {
      max_size = 838860800
      switch   = "on"
    }
    quic {
      switch = "off"
    }
    smart_routing {
      switch = "off"
    }
    standard_debug {
      allow_client_ip_list = []
      expires              = "1969-12-31T16:00:00Z"
      switch               = "off"
    }
    tls_config {
      cipher_suite = "loose-v2023"
      version      = ["TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3"]
    }
    upstream_http2 {
      switch = "off"
    }
    web_socket {
      switch  = "off"
      timeout = 30
    }
  }
}
`

const testAccTeoL7AccSettingUp = `

resource "tencentcloud_teo_l7_acc_setting" "teo_l7_acc_setting" {
  zone_id = "zone-36bjhygh1bxe"
  zone_config {
    accelerate_mainland {
      switch = "on"
    }
    cache {
      custom_time {
        cache_time = 2592000
        switch     = "off"
      }
      follow_origin {
        default_cache          = "off"
        default_cache_strategy = "on"
        default_cache_time     = 0
        switch                 = "on"
      }
      no_cache {
        switch = "off"
      }
    }
    cache_key {
      full_url_cache = "on"
      ignore_case    = "off"
      query_string {
        action = "includeCustom"
        switch = "off"
      }
    }
    cache_prefresh {
      cache_time_percent = 90
      switch             = "off"
    }
    client_ip_country {
      switch      = "off"
    }
    client_ip_header {
      switch      = "off"
    }
    compression {
      algorithms = ["brotli", "gzip"]
      switch     = "on"
    }
    force_redirect_https {
      redirect_status_code = 302
      switch               = "off"
    }
    grpc {
      switch = "off"
    }
    hsts {
      include_sub_domains = "off"
      preload             = "off"
      switch              = "off"
      timeout             = 0
    }
    http2 {
      switch = "off"
    }
    ipv6 {
      switch = "off"
    }
    max_age {
      cache_time    = 600
      follow_origin = "on"
    }
    ocsp_stapling {
      switch = "off"
    }
    offline_cache {
      switch = "on"
    }
    post_max_size {
      max_size = 838860800
      switch   = "on"
    }
    quic {
      switch = "off"
    }
    smart_routing {
      switch = "off"
    }
    standard_debug {
      allow_client_ip_list = []
      expires              = "1969-12-31T16:00:00Z"
      switch               = "off"
    }
    tls_config {
      cipher_suite = "loose-v2023"
      version      = ["TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3"]
    }
    upstream_http2 {
      switch = "off"
    }
    web_socket {
      switch  = "off"
      timeout = 30
    }
  }
}
`

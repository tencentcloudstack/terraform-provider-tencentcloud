package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoZoneSetting_basic -v
func TestAccTencentCloudTeoZoneSetting_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZoneSetting,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneSettingExists("tencentcloud_teo_zone_setting.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "zone_id", defaultZoneId),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", " cache.0.cache.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.follow_origin.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.follow_origin.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.no_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.no_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.full_url_cache", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.ignore_case", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.0.action", "excludeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.0.value.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_prefresh.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_prefresh.0.percent", "90"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_prefresh.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "client_ip_header.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "client_ip_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "compression.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "compression.0.algorithms.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "compression.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "force_redirect.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "force_redirect.0.redirect_status_code", "302"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "force_redirect.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.0.include_sub_domains", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.0.preload", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.http2", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.ocsp_stapling", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.tls_version.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "ipv6.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "ipv6.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "max_age.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "max_age.0.follow_origin", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "offline_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "offline_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "origin.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "origin.0.origin_pull_protocol", "follow"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "post_max_size.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "post_max_size.0.max_size", "524288000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "quic.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "quic.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "smart_routing.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "upstream_http2.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "upstream_http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "web_socket.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "web_socket.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "web_socket.0.timeout", "30"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_zone_setting.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckZoneSettingExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		setting, err := service.DescribeTeoZoneSetting(ctx, rs.Primary.ID)
		if setting == nil {
			return fmt.Errorf("zone setting %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoZoneSettingVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}`

const testAccTeoZoneSetting = testAccTeoZoneSettingVar + `

resource "tencentcloud_teo_zone_setting" "basic" {
  zone_id = var.zone_id

  cache {
    follow_origin {
      switch = "off"
    }

    no_cache {
      switch = "off"
    }
  }

  cache_key {
    full_url_cache = "off"
    ignore_case    = "on"

    query_string {
      action = "excludeCustom"
      switch = "on"
      value  = ["test", "apple"]
    }
  }

  cache_prefresh {
    percent = 90
    switch  = "off"
  }

  client_ip_header {
    switch = "off"
  }

  compression {
    switch = "off"
  }

  force_redirect {
    redirect_status_code = 302
    switch               = "on"
  }

  https {
    http2         = "on"
    ocsp_stapling = "off"
    tls_version   = [
      "TLSv1.2",
      "TLSv1.3",
    ]

    hsts {
      include_sub_domains = "off"
      max_age             = 0
      preload             = "off"
      switch              = "off"
    }
  }

  max_age {
    follow_origin = "off"
    max_age_time  = 0
  }

  offline_cache {
    switch = "off"
  }

  origin {
    origin_pull_protocol = "follow"
  }

  post_max_size {
    max_size = 524288000
    switch   = "on"
  }

  quic {
    switch = "on"
  }

  smart_routing {
    switch = "on"
  }

  upstream_http2 {
    switch = "off"
  }

  web_socket {
    switch  = "off"
    timeout = 30
  }
}

`

package teo_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -test.run TestAccTencentCloudTeoZoneSetting_basic -v
func TestAccTencentCloudTeoZoneSetting_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZoneSetting,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneSettingExists("tencentcloud_teo_zone_setting.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone_setting.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", " cache.0.cache.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.follow_origin.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.follow_origin.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.no_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.no_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.0.action", "includeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_prefresh.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_prefresh.0.percent", "90"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_prefresh.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "client_ip_header.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "client_ip_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "compression.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "compression.0.algorithms.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "compression.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "force_redirect.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "force_redirect.0.redirect_status_code", "302"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "force_redirect.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.0.include_sub_domains", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.0.preload", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.http2", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.ocsp_stapling", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.tls_version.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "ipv6.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "ipv6.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "max_age.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "max_age.0.follow_origin", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "max_age.0.max_age_time", "600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "offline_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "offline_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "origin.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "origin.0.origin_pull_protocol", "follow"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "post_max_size.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "post_max_size.0.max_size", "838860800"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "quic.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "quic.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "smart_routing.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "upstream_http2.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "upstream_http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "web_socket.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "web_socket.0.switch", "off"),
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccTeoZoneSetting = testAccTeoZone + `

resource "tencentcloud_teo_zone_setting" "basic" {
  zone_id = tencentcloud_teo_zone.basic.id

    cache {
        follow_origin {
            switch = "on"
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
            value  = []
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
        algorithms = [
            "brotli",
            "gzip",
        ]
        switch     = "on"
    }

    force_redirect {
        redirect_status_code = 302
        switch               = "off"
    }

    https {
        http2         = "on"
        ocsp_stapling = "off"
        tls_version   = [
            "TLSv1",
            "TLSv1.1",
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

    ipv6 {
        switch = "off"
    }

    max_age {
        follow_origin = "on"
        max_age_time  = 600
    }

    offline_cache {
        switch = "on"
    }

    origin {
        backup_origins       = []
        origin_pull_protocol = "follow"
        origins              = []
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

    upstream_http2 {
        switch = "off"
    }

    web_socket {
        switch  = "off"
        timeout = 30
    }
}
`

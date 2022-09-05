package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTeoZoneSetting_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZoneSetting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone_setting.zone_setting", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_zone_setting.zone_setting",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoZoneSetting = `

resource "tencentcloud_teo_zone_setting" "zone_setting" {
  zone_id = tencentcloud_teo_zone.zone.id

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
    max_age_time  = 600
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

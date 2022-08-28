package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoZoneSetting_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZoneSetting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone_setting.zoneSetting", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_zone_setting.zoneSetting",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoZoneSetting = `

resource "tencentcloud_teo_zone_setting" "zoneSetting" {
  zone_id = ""
  cache {
    cache {
      switch               = ""
      cache_time           = ""
      ignore_cache_control = ""
    }
    no_cache {
      switch = ""
    }
    follow_origin {
      switch = ""
    }

  }
  cache_key {
    full_url_cache = ""
    ignore_case    = ""
    query_string {
      switch = ""
      action = ""
      value  = ""
    }

  }
  max_age {
    max_age_time  = ""
    follow_origin = ""

  }
  offline_cache {
    switch = ""

  }
  quic {
    switch = ""

  }
  post_max_size {
    switch   = ""
    max_size = ""

  }
  compression {
    switch = ""

  }
  upstream_http2 {
    switch = ""

  }
  force_redirect {
    switch               = ""
    redirect_status_code = ""

  }
  https {
    http2         = ""
    ocsp_stapling = ""
    tls_version   = ""
    hsts {
      switch              = ""
      max_age             = ""
      include_sub_domains = ""
      preload             = ""
    }

  }
  origin {
    origin_pull_protocol = ""

  }
  smart_routing {
    switch = ""

  }
  web_socket {
    switch  = ""
    timeout = ""

  }
  client_ip_header {
    switch      = ""
    header_name = ""

  }
  cache_prefresh {
    switch  = ""
    percent = ""

  }
}

`
